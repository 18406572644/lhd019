package services

import (
	"cocktail-bar-system/internal/models"
	"cocktail-bar-system/pkg/database"

	"gorm.io/gorm"
)

func GetIngredients(category, keyword string) ([]models.Ingredient, int64, error) {
	var ingredients []models.Ingredient
	var total int64
	query := database.DB.Model(&models.Ingredient{})

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	query.Count(&total)
	err := query.Order("created_at DESC").Find(&ingredients).Error
	return ingredients, total, err
}

func GetIngredient(id uint) (*models.Ingredient, error) {
	var ingredient models.Ingredient
	err := database.DB.First(&ingredient, id).Error
	if err != nil {
		return nil, err
	}
	return &ingredient, nil
}

func CreateIngredient(ingredient *models.Ingredient) error {
	return database.DB.Create(ingredient).Error
}

func UpdateIngredient(id uint, ingredient *models.Ingredient) error {
	return database.DB.Model(&models.Ingredient{}).Where("id = ?", id).Updates(ingredient).Error
}

func DeleteIngredient(id uint) error {
	return database.DB.Delete(&models.Ingredient{}, id).Error
}

func GetLowStockIngredients() ([]models.Ingredient, error) {
	var ingredients []models.Ingredient
	err := database.DB.Where("stock_quantity <= min_stock").Order("stock_quantity ASC").Find(&ingredients).Error
	return ingredients, err
}

func DeductIngredientStock(id uint, amount float64) error {
	return database.DB.Model(&models.Ingredient{}).Where("id = ?", id).
		Update("stock_quantity", gorm.Expr("stock_quantity - ?", amount)).Error
}

func AddIngredientStock(id uint, amount float64) error {
	return database.DB.Model(&models.Ingredient{}).Where("id = ?", id).
		Update("stock_quantity", gorm.Expr("stock_quantity + ?", amount)).Error
}
