package services

import (
	"cocktail-bar-system/internal/models"
	"cocktail-bar-system/pkg/database"

	"gorm.io/gorm"
)

func GetSpirits(category, keyword string) ([]models.Spirit, int64, error) {
	var spirits []models.Spirit
	var total int64
	query := database.DB.Model(&models.Spirit{})

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR brand LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	err := query.Order("created_at DESC").Find(&spirits).Error
	return spirits, total, err
}

func GetSpirit(id uint) (*models.Spirit, error) {
	var spirit models.Spirit
	err := database.DB.First(&spirit, id).Error
	if err != nil {
		return nil, err
	}
	return &spirit, nil
}

func CreateSpirit(spirit *models.Spirit) error {
	return database.DB.Create(spirit).Error
}

func UpdateSpirit(id uint, spirit *models.Spirit) error {
	return database.DB.Model(&models.Spirit{}).Where("id = ?", id).Updates(spirit).Error
}

func DeleteSpirit(id uint) error {
	return database.DB.Delete(&models.Spirit{}, id).Error
}

func GetLowStockSpirits() ([]models.Spirit, error) {
	var spirits []models.Spirit
	err := database.DB.Where("stock_quantity <= min_stock").Order("stock_quantity ASC").Find(&spirits).Error
	return spirits, err
}

func DeductSpiritStock(id uint, amount float64) error {
	return database.DB.Model(&models.Spirit{}).Where("id = ?", id).
		Update("stock_quantity", gorm.Expr("stock_quantity - ?", amount)).Error
}

func AddSpiritStock(id uint, amount float64) error {
	return database.DB.Model(&models.Spirit{}).Where("id = ?", id).
		Update("stock_quantity", gorm.Expr("stock_quantity + ?", amount)).Error
}
