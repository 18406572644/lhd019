package services

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"cocktail-bar-system/internal/models"
	"cocktail-bar-system/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func generateBatchCode(ingredientType string, ingredientName string) string {
	now := time.Now()
	prefix := fmt.Sprintf("BATCH-%s-%s",
		strings.ToUpper(ingredientType[:3]),
		now.Format("20060102"))
	random := strings.ToUpper(uuid.New().String()[:4])
	return fmt.Sprintf("%s-%s", prefix, random)
}

func GetStockBatches(ingredientType string, ingredientID uint, status string, keyword string) ([]models.StockBatch, int64, error) {
	var batches []models.StockBatch
	var total int64
	query := database.DB.Model(&models.StockBatch{})

	if ingredientType != "" {
		query = query.Where("ingredient_type = ?", ingredientType)
	}
	if ingredientID > 0 {
		query = query.Where("ingredient_id = ?", ingredientID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		query = query.Where("batch_code LIKE ? OR batch_no LIKE ? OR ingredient_name LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	err := query.Order("expiry_date ASC, created_at ASC").Find(&batches).Error
	return batches, total, err
}

func GetStockBatch(id uint) (*models.StockBatch, error) {
	var batch models.StockBatch
	err := database.DB.First(&batch, id).Error
	if err != nil {
		return nil, err
	}
	return &batch, nil
}

func GetStockBatchByCode(batchCode string) (*models.StockBatch, error) {
	var batch models.StockBatch
	err := database.DB.Where("batch_code = ?", batchCode).First(&batch).Error
	if err != nil {
		return nil, err
	}
	return &batch, nil
}

func CreateStockBatch(req *models.StockBatchCreateRequest, purchaseItemID *uint) (*models.StockBatch, error) {
	batch := &models.StockBatch{
		BatchCode:        generateBatchCode(req.IngredientType, req.IngredientName),
		BatchNo:          req.BatchNo,
		IngredientType:   req.IngredientType,
		IngredientID:     req.IngredientID,
		IngredientName:   req.IngredientName,
		PurchaseItemID:   purchaseItemID,
		TotalQuantity:    req.TotalQuantity,
		RemainingQuantity: req.TotalQuantity,
		Unit:             req.Unit,
		UnitPrice:        req.UnitPrice,
		ExpiryDate:       req.ExpiryDate,
		Status:           "normal",
		WarehousePosition: req.WarehousePosition,
		Remark:           req.Remark,
	}

	if err := database.DB.Create(batch).Error; err != nil {
		return nil, err
	}
	return batch, nil
}

func UpdateStockBatchPromotion(id uint, isPromotion bool, remark string) (*models.StockBatch, error) {
	batch, err := GetStockBatch(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"is_promotion": isPromotion,
	}
	if remark != "" {
		updates["remark"] = remark
	}

	if err := database.DB.Model(batch).Updates(updates).Error; err != nil {
		return nil, err
	}
	return GetStockBatch(id)
}

func DeductStockFIFO(tx *gorm.DB, req *models.BatchDeductRequest) ([]models.BatchDeductResult, error) {
	var results []models.BatchDeductResult
	remainingDeduct := req.Quantity

	if remainingDeduct <= 0 {
		return results, nil
	}

	var batches []models.StockBatch
	err := tx.Where("ingredient_type = ? AND ingredient_id = ? AND status = 'normal' AND remaining_quantity > 0",
		req.IngredientType, req.IngredientID).
		Order("expiry_date IS NULL ASC, expiry_date ASC, created_at ASC").
		Set("gorm:query_option", "FOR UPDATE").
		Find(&batches).Error
	if err != nil {
		return nil, err
	}

	var totalAvailable float64
	for _, b := range batches {
		totalAvailable += b.RemainingQuantity
	}

	if totalAvailable < remainingDeduct {
		return nil, fmt.Errorf("insufficient batch stock: need %.2f, have %.2f", remainingDeduct, totalAvailable)
	}

	for _, batch := range batches {
		if remainingDeduct <= 0 {
			break
		}

		deductQty := math.Min(batch.RemainingQuantity, remainingDeduct)
		totalCost := deductQty * batch.UnitPrice

		newRemaining := batch.RemainingQuantity - deductQty
		status := batch.Status
		if newRemaining <= 0.001 {
			newRemaining = 0
			status = "depleted"
		}

		if err := tx.Model(&models.StockBatch{}).Where("id = ?", batch.ID).
			Updates(map[string]interface{}{
				"remaining_quantity": newRemaining,
				"status":             status,
			}).Error; err != nil {
			return nil, err
		}

		outRecord := &models.BatchOutRecord{
			BatchID:        batch.ID,
			BatchCode:      batch.BatchCode,
			IngredientType: batch.IngredientType,
			IngredientID:   batch.IngredientID,
			IngredientName: batch.IngredientName,
			OutType:        req.OutType,
			OutQuantity:    deductQty,
			Unit:           batch.Unit,
			UnitPrice:      batch.UnitPrice,
			TotalCost:      totalCost,
			Operator:       req.Operator,
			Remark:         req.Remark,
		}
		if req.OrderID > 0 {
			outRecord.OrderID = &req.OrderID
			outRecord.OrderNo = req.OrderNo
		}

		if err := tx.Create(outRecord).Error; err != nil {
			return nil, err
		}

		results = append(results, models.BatchDeductResult{
			BatchID:      batch.ID,
			BatchCode:    batch.BatchCode,
			DeductQty:    deductQty,
			RemainingQty: newRemaining,
			UnitPrice:    batch.UnitPrice,
		})

		remainingDeduct -= deductQty
	}

	return results, nil
}

func GetBatchOutRecords(batchID uint, orderID uint, ingredientType string, ingredientID uint) ([]models.BatchOutRecord, int64, error) {
	var records []models.BatchOutRecord
	var total int64
	query := database.DB.Model(&models.BatchOutRecord{})

	if batchID > 0 {
		query = query.Where("batch_id = ?", batchID)
	}
	if orderID > 0 {
		query = query.Where("order_id = ?", orderID)
	}
	if ingredientType != "" {
		query = query.Where("ingredient_type = ?", ingredientType)
	}
	if ingredientID > 0 {
		query = query.Where("ingredient_id = ?", ingredientID)
	}

	query.Count(&total)
	err := query.Order("created_at DESC").Find(&records).Error
	return records, total, err
}

func GetExpiryWarnings(days int) ([]models.ExpiryWarningResult, error) {
	var results []models.ExpiryWarningResult
	var batches []models.StockBatch

	today := time.Now()
	warningDate := today.AddDate(0, 0, days)
	warningDateStr := warningDate.Format("2006-01-02")
	todayStr := today.Format("2006-01-02")

	err := database.DB.Where("status = 'normal' AND remaining_quantity > 0 AND expiry_date IS NOT NULL AND expiry_date <= ?",
		warningDateStr).
		Order("expiry_date ASC").
		Find(&batches).Error
	if err != nil {
		return nil, err
	}

	for _, batch := range batches {
		expiryDate, err := time.Parse("2006-01-02", batch.ExpiryDate)
		if err != nil {
			continue
		}
		daysToExpiry := int(expiryDate.Sub(today).Hours() / 24)

		warningLevel := "normal"
		if daysToExpiry < 0 {
			warningLevel = "expired"
		} else if daysToExpiry <= 7 {
			warningLevel = "urgent"
		} else if daysToExpiry <= 15 {
			warningLevel = "warning"
		} else if daysToExpiry <= 30 {
			warningLevel = "attention"
		}

		_ = todayStr

		results = append(results, models.ExpiryWarningResult{
			StockBatch:  batch,
			DaysToExpiry: daysToExpiry,
			WarningLevel: warningLevel,
		})
	}

	return results, nil
}

func TraceBatch(batchCode string) (*models.BatchTraceResult, error) {
	batch, err := GetStockBatchByCode(batchCode)
	if err != nil {
		return nil, errors.New("batch not found")
	}

	var outRecords []models.BatchOutRecord
	err = database.DB.Where("batch_id = ?", batch.ID).
		Order("created_at DESC").
		Find(&outRecords).Error
	if err != nil {
		return nil, err
	}

	totalOutQty := 0.0
	for _, r := range outRecords {
		totalOutQty += r.OutQuantity
	}

	return &models.BatchTraceResult{
		StockBatch:  *batch,
		OutRecords:  outRecords,
		TotalOutQty: totalOutQty,
	}, nil
}

func UpdateExpiredBatches() error {
	today := time.Now().Format("2006-01-02")
	return database.DB.Model(&models.StockBatch{}).
		Where("status = 'normal' AND expiry_date IS NOT NULL AND expiry_date < ?", today).
		Update("status", "expired").Error
}

func GetBatchOutRecordsByOrder(orderID uint) ([]models.BatchOutRecord, error) {
	var records []models.BatchOutRecord
	err := database.DB.Where("order_id = ?", orderID).Find(&records).Error
	return records, err
}
