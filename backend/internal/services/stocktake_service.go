package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"cocktail-bar-system/internal/models"
	"cocktail-bar-system/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func generateStocktakeNo() string {
	now := time.Now()
	prefix := fmt.Sprintf("ST%s", now.Format("20060102"))
	random := uuid.New().String()[:3]
	return fmt.Sprintf("%s%s", prefix, strings.ToUpper(random))
}

func GetStocktakes(startDate, endDate, status, stocktakeType string) ([]models.Stocktake, int64, error) {
	var stocktakes []models.Stocktake
	var total int64
	query := database.DB.Model(&models.Stocktake{}).Preload("StocktakeItems")

	if startDate != "" {
		query = query.Where("stocktake_date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("stocktake_date <= ?", endDate)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if stocktakeType != "" {
		query = query.Where("stocktake_type = ?", stocktakeType)
	}

	query.Count(&total)
	err := query.Order("stocktake_date DESC, created_at DESC").Find(&stocktakes).Error
	return stocktakes, total, err
}

func GetStocktake(id uint) (*models.Stocktake, error) {
	var stocktake models.Stocktake
	err := database.DB.Preload("StocktakeItems").First(&stocktake, id).Error
	if err != nil {
		return nil, err
	}
	return &stocktake, nil
}

func CreateStocktake(req *models.StocktakeCreateRequest) (*models.Stocktake, error) {
	tx := database.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	stocktake := &models.Stocktake{
		StocktakeNo:   generateStocktakeNo(),
		StocktakeDate: req.StocktakeDate,
		StocktakeType: req.StocktakeType,
		Status:        "draft",
		Operator:      req.Operator,
		Remark:        req.Remark,
	}

	if err := tx.Create(stocktake).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	var totalProfit, totalLoss, totalDiff float64
	var stocktakeItems []models.StocktakeItem

	for _, item := range req.Items {
		var ingredientName string
		var unit string
		var unitPrice float64
		var systemQuantity float64

		if item.IngredientType == "spirit" {
			spirit, err := GetSpirit(item.IngredientID)
			if err != nil {
				tx.Rollback()
				return nil, errors.New("spirit not found")
			}
			ingredientName = spirit.Name
			unit = spirit.Unit
			unitPrice = spirit.CostPrice
			systemQuantity = float64(spirit.StockQuantity)
		} else {
			ingredient, err := GetIngredient(item.IngredientID)
			if err != nil {
				tx.Rollback()
				return nil, errors.New("ingredient not found")
			}
			ingredientName = ingredient.Name
			unit = ingredient.Unit
			unitPrice = ingredient.CostPrice
			systemQuantity = ingredient.StockQuantity
		}

		diffQuantity := item.ActualQuantity - systemQuantity
		diffAmount := diffQuantity * unitPrice

		diffType := "normal"
		if diffQuantity > 0.001 {
			diffType = "profit"
			totalProfit += diffAmount
		} else if diffQuantity < -0.001 {
			diffType = "loss"
			totalLoss += -diffAmount
		}
		totalDiff += diffAmount

		stocktakeItem := models.StocktakeItem{
			StocktakeID:    stocktake.ID,
			IngredientType: item.IngredientType,
			IngredientID:   item.IngredientID,
			IngredientName: ingredientName,
			SystemQuantity: systemQuantity,
			ActualQuantity: item.ActualQuantity,
			DiffQuantity:   diffQuantity,
			Unit:           unit,
			UnitPrice:      unitPrice,
			DiffAmount:     diffAmount,
			DiffType:       diffType,
			Remark:         item.Remark,
		}
		stocktakeItems = append(stocktakeItems, stocktakeItem)
	}

	if len(stocktakeItems) > 0 {
		if err := tx.Create(&stocktakeItems).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Model(stocktake).Updates(map[string]interface{}{
		"total_profit": totalProfit,
		"total_loss":   totalLoss,
		"total_diff":   totalDiff,
	}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return GetStocktake(stocktake.ID)
}

func ConfirmStocktake(id uint, req *models.StocktakeConfirmRequest) (*models.Stocktake, error) {
	tx := database.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	stocktake, err := GetStocktake(id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if stocktake.Status == "confirmed" {
		tx.Rollback()
		return nil, errors.New("stocktake already confirmed")
	}

	now := time.Now()
	for _, item := range stocktake.StocktakeItems {
		if item.DiffQuantity == 0 {
			continue
		}

		if item.IngredientType == "spirit" {
			if err := tx.Model(&models.Spirit{}).Where("id = ?", item.IngredientID).
				Update("stock_quantity", gorm.Expr("stock_quantity + ?", item.DiffQuantity)).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		} else {
			if err := tx.Model(&models.Ingredient{}).Where("id = ?", item.IngredientID).
				Update("stock_quantity", gorm.Expr("stock_quantity + ?", item.DiffQuantity)).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	if err := tx.Model(stocktake).Updates(map[string]interface{}{
		"status":       "confirmed",
		"remark":       req.Remark,
		"confirmed_at": now,
	}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return GetStocktake(id)
}

func DeleteStocktake(id uint) error {
	stocktake, err := GetStocktake(id)
	if err != nil {
		return err
	}
	if stocktake.Status == "confirmed" {
		return errors.New("cannot delete confirmed stocktake")
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("stocktake_id = ?", id).Delete(&models.StocktakeItem{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&models.Stocktake{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}

func GenerateStocktakeItems() ([]map[string]interface{}, error) {
	var result []map[string]interface{}

	var spirits []models.Spirit
	if err := database.DB.Find(&spirits).Error; err != nil {
		return nil, err
	}
	for _, s := range spirits {
		result = append(result, map[string]interface{}{
			"ingredient_type": "spirit",
			"ingredient_id":   s.ID,
			"ingredient_name": s.Name,
			"system_quantity": s.StockQuantity,
			"unit":            s.Unit,
			"unit_price":      s.CostPrice,
			"category":        s.Category,
		})
	}

	var ingredients []models.Ingredient
	if err := database.DB.Find(&ingredients).Error; err != nil {
		return nil, err
	}
	for _, i := range ingredients {
		result = append(result, map[string]interface{}{
			"ingredient_type": "ingredient",
			"ingredient_id":   i.ID,
			"ingredient_name": i.Name,
			"system_quantity": i.StockQuantity,
			"unit":            i.Unit,
			"unit_price":      i.CostPrice,
			"category":        i.Category,
		})
	}

	return result, nil
}
