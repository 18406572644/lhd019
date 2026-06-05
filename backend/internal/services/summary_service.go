package services

import (
	"cocktail-bar-system/internal/models"
	"cocktail-bar-system/pkg/database"
)

func GetBusinessSummary(startDate, endDate string) (*models.BusinessSummary, error) {
	summary := &models.BusinessSummary{}

	var totalRevenue float64
	var totalOrders int64
	var totalCustomers int64

	orderQuery := database.DB.Model(&models.Order{})
	if startDate != "" {
		orderQuery = orderQuery.Where("DATE(created_at) >= ?", startDate)
	}
	if endDate != "" {
		orderQuery = orderQuery.Where("DATE(created_at) <= ?", endDate)
	}

	orderQuery.Select("COALESCE(SUM(actual_amount), 0)").Scan(&totalRevenue)
	orderQuery.Count(&totalOrders)
	orderQuery.Select("COALESCE(SUM(customer_count), 0)").Scan(&totalCustomers)

	summary.TotalRevenue = totalRevenue
	summary.TotalOrders = totalOrders
	summary.TotalCustomers = totalCustomers
	if totalOrders > 0 {
		summary.AverageOrder = totalRevenue / float64(totalOrders)
	}

	var topDrinks []struct {
		RecipeName string
		Quantity   int64
		Revenue    float64
	}
	drinkQuery := database.DB.Table("order_items").
		Select("recipe_name, SUM(quantity) as quantity, SUM(subtotal) as revenue")

	if startDate != "" {
		drinkQuery = drinkQuery.Joins("JOIN orders ON orders.id = order_items.order_id").
			Where("DATE(orders.created_at) >= ?", startDate)
	}
	if endDate != "" {
		drinkQuery = drinkQuery.Joins("JOIN orders ON orders.id = order_items.order_id").
			Where("DATE(orders.created_at) <= ?", endDate)
	}

	drinkQuery.Group("recipe_id, recipe_name").
		Order("quantity DESC").
		Limit(5).
		Scan(&topDrinks)

	for _, d := range topDrinks {
		summary.TopDrinks = append(summary.TopDrinks, models.TopDrink{
			RecipeName: d.RecipeName,
			Quantity:   d.Quantity,
			Revenue:    d.Revenue,
		})
	}

	var dailyRevenue []struct {
		Date    string
		Revenue float64
		Orders  int64
	}
	dailyQuery := database.DB.Table("orders").
		Select("DATE(created_at) as date, COALESCE(SUM(actual_amount), 0) as revenue, COUNT(*) as orders")

	if startDate != "" {
		dailyQuery = dailyQuery.Where("DATE(created_at) >= ?", startDate)
	}
	if endDate != "" {
		dailyQuery = dailyQuery.Where("DATE(created_at) <= ?", endDate)
	}

	dailyQuery.Group("DATE(created_at)").
		Order("date DESC").
		Limit(7).
		Scan(&dailyRevenue)

	for _, d := range dailyRevenue {
		summary.DailyRevenue = append(summary.DailyRevenue, models.DailyRevenue{
			Date:    d.Date,
			Revenue: d.Revenue,
			Orders:  d.Orders,
		})
	}

	lowStockSpirits, _ := GetLowStockSpirits()
	summary.LowStockSpirits = lowStockSpirits

	lowStockIngredients, _ := GetLowStockIngredients()
	summary.LowStockIngredients = lowStockIngredients

	var totalWasteCost float64
	wasteQuery := database.DB.Model(&models.WasteRecord{})
	if startDate != "" {
		wasteQuery = wasteQuery.Where("DATE(created_at) >= ?", startDate)
	}
	if endDate != "" {
		wasteQuery = wasteQuery.Where("DATE(created_at) <= ?", endDate)
	}
	wasteQuery.Select("COALESCE(SUM(cost), 0)").Scan(&totalWasteCost)
	summary.TotalWasteCost = totalWasteCost

	return summary, nil
}
