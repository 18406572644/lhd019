package services

import (
	"cocktail-bar-system/internal/models"
	"cocktail-bar-system/pkg/database"

	"gorm.io/gorm"
)

func GetRecipes(category, keyword string, isSignature *bool) ([]models.Recipe, int64, error) {
	var recipes []models.Recipe
	var total int64
	query := database.DB.Model(&models.Recipe{}).Preload("RecipeIngredients")

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR taste_profile LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if isSignature != nil {
		query = query.Where("is_signature = ?", *isSignature)
	}

	query.Count(&total)
	err := query.Order("created_at DESC").Find(&recipes).Error
	return recipes, total, err
}

func GetRecipe(id uint) (*models.Recipe, error) {
	var recipe models.Recipe
	err := database.DB.Preload("RecipeIngredients").First(&recipe, id).Error
	if err != nil {
		return nil, err
	}
	return &recipe, nil
}

func CreateRecipe(recipe *models.Recipe) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(recipe).Error; err != nil {
			return err
		}
		return nil
	})
}

func UpdateRecipe(id uint, recipe *models.Recipe) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Recipe{}).Where("id = ?", id).Omit("RecipeIngredients").Updates(recipe).Error; err != nil {
			return err
		}
		if err := tx.Where("recipe_id = ?", id).Delete(&models.RecipeIngredient{}).Error; err != nil {
			return err
		}
		if len(recipe.RecipeIngredients) > 0 {
			for i := range recipe.RecipeIngredients {
				recipe.RecipeIngredients[i].RecipeID = id
			}
			if err := tx.Create(&recipe.RecipeIngredients).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func DeleteRecipe(id uint) error {
	return database.DB.Delete(&models.Recipe{}, id).Error
}
