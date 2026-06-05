package services

import (
	"cocktail-bar-system/internal/models"
	"cocktail-bar-system/pkg/database"

	"gorm.io/gorm"
)

func GetWasteRecords(startDate, endDate, ingredientType string) ([]models.WasteRecord, int64, error) {
	var records []models.WasteRecord
	var total int64
	query := database.DB.Model(&models.WasteRecord{})

	if startDate != "" {
		query = query.Where("DATE(created_at) >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("DATE(created_at) <= ?", endDate)
	}
	if ingredientType != "" {
		query = query.Where("ingredient_type = ?", ingredientType)
	}

	query.Count(&total)
	err := query.Order("created_at DESC").Find(&records).Error
	return records, total, err
}

func GetWasteRecord(id uint) (*models.WasteRecord, error) {
	var record models.WasteRecord
	err := database.DB.First(&record, id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func CreateWasteRecord(req *models.WasteCreateRequest) (*models.WasteRecord, error) {
	tx := database.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	var ingredientName string
	var cost float64
	var unit string

	if req.IngredientType == "spirit" {
		spirit, err := GetSpirit(req.IngredientID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		ingredientName = spirit.Name
		unit = spirit.Unit
		mlPerBottle := float64(spirit.VolumeMl)
		costPerMl := spirit.CostPrice / mlPerBottle
		cost = costPerMl * req.Amount

		if err := tx.Model(&models.Spirit{}).Where("id = ?", req.IngredientID).
			Update("stock_quantity", gorm.Expr("stock_quantity - ?", req.Amount/mlPerBottle)).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	} else {
		ingredient, err := GetIngredient(req.IngredientID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		ingredientName = ingredient.Name
		unit = ingredient.Unit
		cost = ingredient.CostPrice * req.Amount

		if err := tx.Model(&models.Ingredient{}).Where("id = ?", req.IngredientID).
			Update("stock_quantity", gorm.Expr("stock_quantity - ?", req.Amount)).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	record := &models.WasteRecord{
		IngredientType: req.IngredientType,
		IngredientID:   req.IngredientID,
		IngredientName: ingredientName,
		Amount:         req.Amount,
		Unit:           unit,
		Reason:         req.Reason,
		Cost:           cost,
		Operator:       req.Operator,
		Remark:         req.Remark,
	}

	if err := tx.Create(record).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return record, nil
}

func DeleteWasteRecord(id uint) error {
	return database.DB.Delete(&models.WasteRecord{}, id).Error
}
