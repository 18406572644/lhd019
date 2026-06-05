package services

import (
	"cocktail-bar-system/internal/models"
	"cocktail-bar-system/pkg/database"
)

func GetSpecialCreations(status, keyword string) ([]models.SpecialCreation, int64, error) {
	var specials []models.SpecialCreation
	var total int64
	query := database.DB.Model(&models.SpecialCreation{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR creator LIKE ? OR taste_profile LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	err := query.Order("created_at DESC").Find(&specials).Error
	return specials, total, err
}

func GetSpecialCreation(id uint) (*models.SpecialCreation, error) {
	var special models.SpecialCreation
	err := database.DB.First(&special, id).Error
	if err != nil {
		return nil, err
	}
	return &special, nil
}

func CreateSpecialCreation(special *models.SpecialCreation) error {
	return database.DB.Create(special).Error
}

func UpdateSpecialCreation(id uint, special *models.SpecialCreation) error {
	return database.DB.Model(&models.SpecialCreation{}).Where("id = ?", id).Updates(special).Error
}

func DeleteSpecialCreation(id uint) error {
	return database.DB.Delete(&models.SpecialCreation{}, id).Error
}
