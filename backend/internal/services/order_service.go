package services

import (
	"errors"
	"fmt"
	"time"

	"cocktail-bar-system/internal/models"
	"cocktail-bar-system/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func generateOrderNo() string {
	now := time.Now()
	prefix := fmt.Sprintf("ORD%s", now.Format("20060102"))
	random := uuid.New().String()[:8]
	return fmt.Sprintf("%s%s", prefix, random)
}

func GetOrders(startDate, endDate, status string) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64
	query := database.DB.Model(&models.Order{}).Preload("OrderItems")

	if startDate != "" {
		query = query.Where("DATE(created_at) >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("DATE(created_at) <= ?", endDate)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	err := query.Order("created_at DESC").Find(&orders).Error
	return orders, total, err
}

func GetOrder(id uint) (*models.Order, error) {
	var order models.Order
	err := database.DB.Preload("OrderItems").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func CreateOrder(req *models.OrderCreateRequest) (*models.Order, error) {
	tx := database.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var totalAmount float64
	var orderItems []models.OrderItem

	for _, item := range req.Items {
		recipe, err := GetRecipe(item.RecipeID)
		if err != nil {
			tx.Rollback()
			return nil, errors.New("recipe not found: " + err.Error())
		}

		subtotal := float64(item.Quantity) * recipe.Price
		totalAmount += subtotal

		orderItem := models.OrderItem{
			RecipeID:   item.RecipeID,
			RecipeName: recipe.Name,
			Quantity:   item.Quantity,
			UnitPrice:  recipe.Price,
			Subtotal:   subtotal,
			Remark:     item.Remark,
		}
		orderItems = append(orderItems, orderItem)

		for _, ri := range recipe.RecipeIngredients {
			totalDeduct := ri.Amount * float64(item.Quantity)
			if ri.IngredientType == "spirit" {
				spirit, err := GetSpirit(ri.IngredientID)
				if err != nil {
					tx.Rollback()
					return nil, errors.New("spirit not found")
				}
				mlPerBottle := float64(spirit.VolumeMl)
				bottlesToDeduct := totalDeduct / mlPerBottle
				if float64(spirit.StockQuantity) < bottlesToDeduct {
					tx.Rollback()
					return nil, fmt.Errorf("insufficient stock for spirit: %s, need %.2f bottles, have %d", spirit.Name, bottlesToDeduct, spirit.StockQuantity)
				}
				if err := tx.Model(&models.Spirit{}).Where("id = ?", ri.IngredientID).
					Update("stock_quantity", gorm.Expr("stock_quantity - ?", bottlesToDeduct)).Error; err != nil {
					tx.Rollback()
					return nil, err
				}
			} else {
				ingredient, err := GetIngredient(ri.IngredientID)
				if err != nil {
					tx.Rollback()
					return nil, errors.New("ingredient not found")
				}
				if ingredient.StockQuantity < totalDeduct {
					tx.Rollback()
					return nil, fmt.Errorf("insufficient stock for ingredient: %s, need %.2f %s, have %.2f", ingredient.Name, totalDeduct, ingredient.Unit, ingredient.StockQuantity)
				}
				if err := tx.Model(&models.Ingredient{}).Where("id = ?", ri.IngredientID).
					Update("stock_quantity", gorm.Expr("stock_quantity - ?", totalDeduct)).Error; err != nil {
					tx.Rollback()
					return nil, err
				}
			}
		}
	}

	actualAmount := totalAmount - req.Discount
	if actualAmount < 0 {
		actualAmount = 0
	}

	order := &models.Order{
		OrderNo:       generateOrderNo(),
		TableNo:       req.TableNo,
		CustomerCount: req.CustomerCount,
		TotalAmount:   totalAmount,
		Discount:      req.Discount,
		ActualAmount:  actualAmount,
		PaymentMethod: req.PaymentMethod,
		Status:        "completed",
		Remark:        req.Remark,
		OrderItems:    orderItems,
	}

	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return GetOrder(order.ID)
}

func UpdateOrder(id uint, req *models.OrderCreateRequest) (*models.Order, error) {
	return nil, errors.New("update order not implemented")
}

func DeleteOrder(id uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var order models.Order
		if err := tx.Preload("OrderItems").First(&order, id).Error; err != nil {
			return err
		}

		for _, item := range order.OrderItems {
			recipe, err := GetRecipe(item.RecipeID)
			if err != nil {
				continue
			}
			for _, ri := range recipe.RecipeIngredients {
				totalAdd := ri.Amount * float64(item.Quantity)
				if ri.IngredientType == "spirit" {
					spirit, _ := GetSpirit(ri.IngredientID)
					if spirit != nil {
						mlPerBottle := float64(spirit.VolumeMl)
						bottlesToAdd := totalAdd / mlPerBottle
						tx.Model(&models.Spirit{}).Where("id = ?", ri.IngredientID).
							Update("stock_quantity", gorm.Expr("stock_quantity + ?", bottlesToAdd))
					}
				} else {
					tx.Model(&models.Ingredient{}).Where("id = ?", ri.IngredientID).
						Update("stock_quantity", gorm.Expr("stock_quantity + ?", totalAdd))
				}
			}
		}

		if err := tx.Where("order_id = ?", id).Delete(&models.OrderItem{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&models.Order{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}
