package services

import (
	"fmt"
	"time"

	"cocktail-bar-system/internal/models"
	"cocktail-bar-system/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func generatePurchaseNo() string {
	now := time.Now()
	prefix := fmt.Sprintf("PUR%s", now.Format("20060102"))
	random := uuid.New().String()[:8]
	return fmt.Sprintf("%s%s", prefix, random)
}

func GetPurchases(startDate, endDate, supplier string) ([]models.Purchase, int64, error) {
	var purchases []models.Purchase
	var total int64
	query := database.DB.Model(&models.Purchase{}).Preload("PurchaseItems")

	if startDate != "" {
		query = query.Where("purchase_date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("purchase_date <= ?", endDate)
	}
	if supplier != "" {
		query = query.Where("supplier LIKE ?", "%"+supplier+"%")
	}

	query.Count(&total)
	err := query.Order("purchase_date DESC").Find(&purchases).Error
	return purchases, total, err
}

func GetPurchase(id uint) (*models.Purchase, error) {
	var purchase models.Purchase
	err := database.DB.Preload("PurchaseItems").First(&purchase, id).Error
	if err != nil {
		return nil, err
	}
	return &purchase, nil
}

func CreatePurchase(req *models.PurchaseCreateRequest) (*models.Purchase, error) {
	tx := database.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	var totalAmount float64
	var purchaseItems []models.PurchaseItem

	for idx, item := range req.Items {
		subtotal := item.Quantity * item.UnitPrice
		totalAmount += subtotal

		var ingredientName string
		if item.IngredientType == "spirit" {
			spirit, err := GetSpirit(item.IngredientID)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			ingredientName = spirit.Name
		} else {
			ingredient, err := GetIngredient(item.IngredientID)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			ingredientName = ingredient.Name
		}

		if item.BatchNo == "" || item.ExpiryDate == "" {
			tx.Rollback()
			return nil, fmt.Errorf("采购明细第 %d 项：批次号和保质期为必填项", idx+1)
		}

		purchaseItem := models.PurchaseItem{
			IngredientType: item.IngredientType,
			IngredientID:   item.IngredientID,
			IngredientName: ingredientName,
			Quantity:       item.Quantity,
			Unit:           item.Unit,
			UnitPrice:      item.UnitPrice,
			Subtotal:       subtotal,
			BatchNo:        item.BatchNo,
			ExpiryDate:     item.ExpiryDate,
		}
		purchaseItems = append(purchaseItems, purchaseItem)

		if item.IngredientType == "spirit" {
			if err := tx.Model(&models.Spirit{}).Where("id = ?", item.IngredientID).
				Update("stock_quantity", gorm.Expr("stock_quantity + ?", item.Quantity)).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		} else {
			if err := tx.Model(&models.Ingredient{}).Where("id = ?", item.IngredientID).
				Update("stock_quantity", gorm.Expr("stock_quantity + ?", item.Quantity)).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	purchase := &models.Purchase{
		PurchaseNo:    generatePurchaseNo(),
		Supplier:      req.Supplier,
		TotalAmount:   totalAmount,
		PurchaseDate:  req.PurchaseDate,
		Operator:      req.Operator,
		Remark:        req.Remark,
		PurchaseItems: purchaseItems,
	}

	if err := tx.Create(purchase).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	var savedItems []models.PurchaseItem
	if err := tx.Where("purchase_id = ?", purchase.ID).Find(&savedItems).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, item := range savedItems {
		batchReq := &models.StockBatchCreateRequest{
			IngredientType: item.IngredientType,
			IngredientID:   item.IngredientID,
			IngredientName: item.IngredientName,
			BatchNo:        item.BatchNo,
			TotalQuantity:  item.Quantity,
			Unit:           item.Unit,
			UnitPrice:      item.UnitPrice,
			ExpiryDate:     item.ExpiryDate,
			Remark:         fmt.Sprintf("采购单: %s", purchase.PurchaseNo),
		}
		purchaseItemID := item.ID
		_, err := CreateStockBatchWithTx(tx, batchReq, &purchaseItemID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()
	return GetPurchase(purchase.ID)
}

func CreateStockBatchWithTx(tx *gorm.DB, req *models.StockBatchCreateRequest, purchaseItemID *uint) (*models.StockBatch, error) {
	batch := &models.StockBatch{
		BatchCode:         generateBatchCode(req.IngredientType, req.IngredientName),
		BatchNo:           req.BatchNo,
		IngredientType:    req.IngredientType,
		IngredientID:      req.IngredientID,
		IngredientName:    req.IngredientName,
		PurchaseItemID:    purchaseItemID,
		TotalQuantity:     req.TotalQuantity,
		RemainingQuantity: req.TotalQuantity,
		Unit:              req.Unit,
		UnitPrice:         req.UnitPrice,
		ExpiryDate:        req.ExpiryDate,
		Status:            "normal",
		WarehousePosition: req.WarehousePosition,
		Remark:            req.Remark,
	}

	if err := tx.Create(batch).Error; err != nil {
		return nil, err
	}
	return batch, nil
}

func DeletePurchase(id uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var purchase models.Purchase
		if err := tx.Preload("PurchaseItems").First(&purchase, id).Error; err != nil {
			return err
		}

		for _, item := range purchase.PurchaseItems {
			if item.IngredientType == "spirit" {
				tx.Model(&models.Spirit{}).Where("id = ?", item.IngredientID).
					Update("stock_quantity", gorm.Expr("stock_quantity - ?", item.Quantity))
			} else {
				tx.Model(&models.Ingredient{}).Where("id = ?", item.IngredientID).
					Update("stock_quantity", gorm.Expr("stock_quantity - ?", item.Quantity))
			}

			var batches []models.StockBatch
			if err := tx.Where("purchase_item_id = ?", item.ID).Find(&batches).Error; err != nil {
				return err
			}

			for _, batch := range batches {
				if err := tx.Where("batch_id = ?", batch.ID).Delete(&models.BatchOutRecord{}).Error; err != nil {
					return err
				}
			}

			if err := tx.Where("purchase_item_id = ?", item.ID).Delete(&models.StockBatch{}).Error; err != nil {
				return err
			}
		}

		if err := tx.Where("purchase_id = ?", id).Delete(&models.PurchaseItem{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&models.Purchase{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}
