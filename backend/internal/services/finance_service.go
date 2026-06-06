package services

import (
	"time"

	"cocktail-bar-system/internal/models"
	"cocktail-bar-system/pkg/database"

	"gorm.io/gorm"
)

func applyDateFilter(query *gorm.DB, params models.FinanceFilterParams) *gorm.DB {
	if params.StartDate != "" {
		query = query.Where("DATE(created_at) >= ?", params.StartDate)
	}
	if params.EndDate != "" {
		query = query.Where("DATE(created_at) <= ?", params.EndDate)
	}
	return query
}

func calculateGrowthRate(current, previous float64) float64 {
	if previous == 0 {
		return 0
	}
	return ((current - previous) / previous) * 100
}

func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

func getPreviousYearDates(startDate, endDate string) (string, string) {
	start, _ := parseDate(startDate)
	end, _ := parseDate(endDate)
	return start.AddDate(-1, 0, 0).Format("2006-01-02"),
		end.AddDate(-1, 0, 0).Format("2006-01-02")
}

func getPreviousQuarterDates(startDate, endDate string) (string, string) {
	start, _ := parseDate(startDate)
	end, _ := parseDate(endDate)
	return start.AddDate(0, -3, 0).Format("2006-01-02"),
		end.AddDate(0, -3, 0).Format("2006-01-02")
}

func getTotalRevenueForPeriod(startDate, endDate string) float64 {
	var revenue float64
	query := database.DB.Model(&models.Order{})
	if startDate != "" {
		query = query.Where("DATE(created_at) >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("DATE(created_at) <= ?", endDate)
	}
	query.Select("COALESCE(SUM(actual_amount), 0)").Scan(&revenue)
	return revenue
}

func GetRevenueReport(params models.FinanceFilterParams) (*models.RevenueReport, error) {
	report := &models.RevenueReport{
		Period:    params.Period,
		StartDate: params.StartDate,
		EndDate:   params.EndDate,
	}

	var totalRevenue float64
	var totalOrders int64
	var totalCustomers int64

	orderQuery := database.DB.Model(&models.Order{})
	orderQuery = applyDateFilter(orderQuery, params)

	orderQuery.Select("COALESCE(SUM(actual_amount), 0)").Scan(&totalRevenue)
	orderQuery.Count(&totalOrders)
	orderQuery.Select("COALESCE(SUM(customer_count), 0)").Scan(&totalCustomers)

	report.TotalRevenue = totalRevenue
	report.TotalOrders = totalOrders
	report.TotalCustomers = totalCustomers

	if totalOrders > 0 {
		report.AverageOrder = totalRevenue / float64(totalOrders)
	}
	if totalCustomers > 0 {
		report.AverageCustomer = totalRevenue / float64(totalCustomers)
	}

	if params.StartDate != "" && params.EndDate != "" {
		yoyStart, yoyEnd := getPreviousYearDates(params.StartDate, params.EndDate)
		yoyRevenue := getTotalRevenueForPeriod(yoyStart, yoyEnd)
		report.YoYPrevious = yoyRevenue
		report.YoYGrowth = calculateGrowthRate(totalRevenue, yoyRevenue)

		qoqStart, qoqEnd := getPreviousQuarterDates(params.StartDate, params.EndDate)
		qoqRevenue := getTotalRevenueForPeriod(qoqStart, qoqEnd)
		report.QoQPrevious = qoqRevenue
		report.QoQGrowth = calculateGrowthRate(totalRevenue, qoqRevenue)

		type dailyResult struct {
			Date      string
			Revenue   float64
			Orders    int64
			Customers int64
		}
		var dailyResults []dailyResult

		dailyQuery := database.DB.Model(&models.Order{}).
			Select("DATE(created_at) as date, COALESCE(SUM(actual_amount), 0) as revenue, COUNT(*) as orders, COALESCE(SUM(customer_count), 0) as customers")
		dailyQuery = applyDateFilter(dailyQuery, params)
		dailyQuery.Group("DATE(created_at)").Order("date ASC").Scan(&dailyResults)

		for _, dr := range dailyResults {
			report.DailyData = append(report.DailyData, models.DailyRevenueData{
				Date:      dr.Date,
				Revenue:   dr.Revenue,
				Orders:    dr.Orders,
				Customers: dr.Customers,
			})
		}
	}

	return report, nil
}

func GetCostAnalysisReport(params models.FinanceFilterParams) (*models.CostAnalysisReport, error) {
	report := &models.CostAnalysisReport{
		Period:    params.Period,
		StartDate: params.StartDate,
		EndDate:   params.EndDate,
	}

	orderQuery := database.DB.Model(&models.Order{})
	orderQuery = applyDateFilter(orderQuery, params)
	orderQuery.Select("COALESCE(SUM(actual_amount), 0)").Scan(&report.TotalRevenue)

	var ingredientCost, spiritCost float64
	orderItemQuery := database.DB.Table("order_items oi").
		Joins("JOIN orders o ON o.id = oi.order_id").
		Joins("JOIN recipes r ON r.id = oi.recipe_id")
	orderItemQuery = applyDateFilter(orderItemQuery, params)

	orderItemQuery.Select(`
		COALESCE(SUM(CASE WHEN ri.ingredient_type = 'ingredient' THEN ri.amount * i.cost_price * oi.quantity ELSE 0 END), 0)
	`).
		Joins("JOIN recipe_ingredients ri ON ri.recipe_id = r.id").
		Joins("LEFT JOIN ingredients i ON i.id = ri.ingredient_id AND ri.ingredient_type = 'ingredient'").
		Scan(&ingredientCost)

	orderItemQuery.Select(`
		COALESCE(SUM(CASE WHEN ri.ingredient_type = 'spirit' THEN ri.amount * (s.cost_price / s.volume_ml) * oi.quantity ELSE 0 END), 0)
	`).
		Joins("JOIN recipe_ingredients ri ON ri.recipe_id = r.id").
		Joins("LEFT JOIN spirits s ON s.id = ri.ingredient_id AND ri.ingredient_type = 'spirit'").
		Scan(&spiritCost)

	report.IngredientCost = ingredientCost
	report.SpiritCost = spiritCost
	report.TotalMaterialCost = ingredientCost + spiritCost

	wasteQuery := database.DB.Model(&models.WasteRecord{})
	wasteQuery = applyDateFilter(wasteQuery, params)
	wasteQuery.Select("COALESCE(SUM(cost), 0)").Scan(&report.WasteCost)

	purchaseQuery := database.DB.Model(&models.Purchase{})
	purchaseQuery = applyDateFilter(purchaseQuery, params)
	purchaseQuery.Select("COALESCE(SUM(total_amount), 0)").Scan(&report.PurchaseCost)

	opCostQuery := database.DB.Model(&models.OperatingCost{})
	opCostQuery = applyDateFilter(opCostQuery, params)
	opCostQuery.Select("COALESCE(SUM(amount), 0)").Scan(&report.OperatingCost)

	report.TotalCost = report.TotalMaterialCost + report.WasteCost + report.PurchaseCost + report.OperatingCost
	report.GrossProfit = report.TotalRevenue - report.TotalMaterialCost
	if report.TotalRevenue > 0 {
		report.GrossMargin = (report.GrossProfit / report.TotalRevenue) * 100
	}
	report.NetProfit = report.TotalRevenue - report.TotalCost
	if report.TotalRevenue > 0 {
		report.NetMargin = (report.NetProfit / report.TotalRevenue) * 100
	}

	if report.TotalCost > 0 {
		report.CostBreakdown = []models.CostBreakdownItem{
			{Name: "原料成本", Value: report.IngredientCost, Ratio: (report.IngredientCost / report.TotalCost) * 100},
			{Name: "基酒成本", Value: report.SpiritCost, Ratio: (report.SpiritCost / report.TotalCost) * 100},
			{Name: "浪费成本", Value: report.WasteCost, Ratio: (report.WasteCost / report.TotalCost) * 100},
			{Name: "采购成本", Value: report.PurchaseCost, Ratio: (report.PurchaseCost / report.TotalCost) * 100},
			{Name: "运营成本", Value: report.OperatingCost, Ratio: (report.OperatingCost / report.TotalCost) * 100},
		}
	}

	return report, nil
}

func GetCategorySalesReport(params models.FinanceFilterParams) (*models.CategorySalesReport, error) {
	report := &models.CategorySalesReport{
		Period:    params.Period,
		StartDate: params.StartDate,
		EndDate:   params.EndDate,
	}

	type categoryResult struct {
		Category string
		Quantity int64
		Revenue  float64
	}
	var categoryResults []categoryResult

	categoryQuery := database.DB.Table("order_items oi").
		Joins("JOIN orders o ON o.id = oi.order_id").
		Joins("JOIN recipes r ON r.id = oi.recipe_id").
		Select("r.category, SUM(oi.quantity) as quantity, SUM(oi.subtotal) as revenue")
	categoryQuery = applyDateFilter(categoryQuery, params)
	if params.Category != "" {
		categoryQuery = categoryQuery.Where("r.category = ?", params.Category)
	}
	categoryQuery.Group("r.category").Order("revenue DESC").Scan(&categoryResults)

	var totalRevenue float64
	for _, cr := range categoryResults {
		totalRevenue += cr.Revenue
	}

	for _, cr := range categoryResults {
		percentage := 0.0
		if totalRevenue > 0 {
			percentage = (cr.Revenue / totalRevenue) * 100
		}
		report.CategorySales = append(report.CategorySales, models.CategorySales{
			Category:   cr.Category,
			Quantity:   cr.Quantity,
			Revenue:    cr.Revenue,
			Percentage: percentage,
		})
	}

	type recipeResult struct {
		RecipeID   uint
		RecipeName string
		Category   string
		Quantity   int64
		Revenue    float64
		Cost       float64
	}
	var recipeResults []recipeResult

	recipeQuery := database.DB.Table("order_items oi").
		Joins("JOIN orders o ON o.id = oi.order_id").
		Joins("JOIN recipes r ON r.id = oi.recipe_id").
		Select("r.id as recipe_id, r.name as recipe_name, r.category, SUM(oi.quantity) as quantity, SUM(oi.subtotal) as revenue, SUM(r.cost * oi.quantity) as cost")
	recipeQuery = applyDateFilter(recipeQuery, params)
	if params.Category != "" {
		recipeQuery = recipeQuery.Where("r.category = ?", params.Category)
	}
	recipeQuery.Group("r.id, r.name, r.category").Order("revenue DESC").Scan(&recipeResults)

	for _, rr := range recipeResults {
		profit := rr.Revenue - rr.Cost
		profitMargin := 0.0
		if rr.Revenue > 0 {
			profitMargin = (profit / rr.Revenue) * 100
		}
		report.RecipeSales = append(report.RecipeSales, models.RecipeSales{
			RecipeID:     rr.RecipeID,
			RecipeName:   rr.RecipeName,
			Category:     rr.Category,
			Quantity:     rr.Quantity,
			Revenue:      rr.Revenue,
			Cost:         rr.Cost,
			Profit:       profit,
			ProfitMargin: profitMargin,
		})
	}

	timeSlots := []string{"00-02", "02-04", "04-06", "06-08", "08-10", "10-12", "12-14", "14-16", "16-18", "18-20", "20-22", "22-24"}
	for _, slot := range timeSlots {
		startHour, _ := time.Parse("15", slot[:2])
		endHour, _ := time.Parse("15", slot[3:])

		var tsResult struct {
			Quantity int64
			Revenue  float64
			Orders   int64
		}

		tsQuery := database.DB.Table("order_items oi").
			Joins("JOIN orders o ON o.id = oi.order_id").
			Select("COALESCE(SUM(oi.quantity), 0) as quantity, COALESCE(SUM(oi.subtotal), 0) as revenue, COUNT(DISTINCT o.id) as orders").
			Where("HOUR(o.created_at) >= ? AND HOUR(o.created_at) < ?", startHour.Hour(), endHour.Hour())
		tsQuery = applyDateFilter(tsQuery, params)
		tsQuery.Scan(&tsResult)

		report.TimeSlotSales = append(report.TimeSlotSales, models.TimeSlotSales{
			TimeSlot: slot,
			Quantity: tsResult.Quantity,
			Revenue:  tsResult.Revenue,
			Orders:   tsResult.Orders,
		})
	}

	return report, nil
}

func GetPaymentReconciliation(params models.FinanceFilterParams) (*models.PaymentReconciliation, error) {
	report := &models.PaymentReconciliation{
		Period:    params.Period,
		StartDate: params.StartDate,
		EndDate:   params.EndDate,
	}

	orderQuery := database.DB.Model(&models.Order{})
	orderQuery = applyDateFilter(orderQuery, params)
	orderQuery.Select("COALESCE(SUM(actual_amount), 0)").Scan(&report.TotalRevenue)

	type pmResult struct {
		PaymentMethod string
		OrderCount    int64
		TotalAmount   float64
	}
	var pmResults []pmResult

	pmQuery := database.DB.Model(&models.Order{}).
		Select("payment_method, COUNT(*) as order_count, COALESCE(SUM(actual_amount), 0) as total_amount")
	pmQuery = applyDateFilter(pmQuery, params)
	if params.PaymentMethod != "" {
		pmQuery = pmQuery.Where("payment_method = ?", params.PaymentMethod)
	}
	pmQuery.Group("payment_method").Order("total_amount DESC").Scan(&pmResults)

	for _, pm := range pmResults {
		percentage := 0.0
		if report.TotalRevenue > 0 {
			percentage = (pm.TotalAmount / report.TotalRevenue) * 100
		}
		report.PaymentMethods = append(report.PaymentMethods, models.PaymentMethodDetail{
			PaymentMethod: pm.PaymentMethod,
			OrderCount:    pm.OrderCount,
			TotalAmount:   pm.TotalAmount,
			Percentage:    percentage,
		})
	}

	logs, _, _ := GetReconciliationLogs(params)
	report.ReconciliationLogs = logs

	return report, nil
}

func GetProfitReport(params models.FinanceFilterParams) (*models.ProfitReport, error) {
	report := &models.ProfitReport{
		Period:    params.Period,
		StartDate: params.StartDate,
		EndDate:   params.EndDate,
	}

	orderQuery := database.DB.Model(&models.Order{})
	orderQuery = applyDateFilter(orderQuery, params)
	orderQuery.Select("COALESCE(SUM(actual_amount), 0)").Scan(&report.TotalRevenue)

	var ingredientCost, spiritCost float64
	materialQuery := database.DB.Table("order_items oi").
		Joins("JOIN orders o ON o.id = oi.order_id").
		Joins("JOIN recipes r ON r.id = oi.recipe_id").
		Joins("JOIN recipe_ingredients ri ON ri.recipe_id = r.id")
	materialQuery = applyDateFilter(materialQuery, params)

	materialQuery.Select(`
		COALESCE(SUM(CASE WHEN ri.ingredient_type = 'ingredient' THEN ri.amount * i.cost_price * oi.quantity ELSE 0 END), 0)
	`).
		Joins("LEFT JOIN ingredients i ON i.id = ri.ingredient_id AND ri.ingredient_type = 'ingredient'").
		Scan(&ingredientCost)

	materialQuery.Select(`
		COALESCE(SUM(CASE WHEN ri.ingredient_type = 'spirit' THEN ri.amount * (s.cost_price / s.volume_ml) * oi.quantity ELSE 0 END), 0)
	`).
		Joins("LEFT JOIN spirits s ON s.id = ri.ingredient_id AND ri.ingredient_type = 'spirit'").
		Scan(&spiritCost)

	report.MaterialCost = ingredientCost + spiritCost

	wasteQuery := database.DB.Model(&models.WasteRecord{})
	wasteQuery = applyDateFilter(wasteQuery, params)
	wasteQuery.Select("COALESCE(SUM(cost), 0)").Scan(&report.WasteCost)

	opCostQuery := database.DB.Model(&models.OperatingCost{})
	opCostQuery = applyDateFilter(opCostQuery, params)
	opCostQuery.Select("COALESCE(SUM(amount), 0)").Scan(&report.OperatingCost)

	report.TotalExpenses = report.MaterialCost + report.WasteCost + report.OperatingCost
	report.GrossProfit = report.TotalRevenue - report.MaterialCost
	if report.TotalRevenue > 0 {
		report.GrossMargin = (report.GrossProfit / report.TotalRevenue) * 100
	}
	report.NetProfit = report.TotalRevenue - report.TotalExpenses
	if report.TotalRevenue > 0 {
		report.NetMargin = (report.NetProfit / report.TotalRevenue) * 100
	}

	report.ProfitBreakdown = []models.ProfitBreakdownItem{
		{Name: "营业收入", Value: report.TotalRevenue, Type: "revenue"},
		{Name: "物料成本", Value: report.MaterialCost, Type: "expense"},
		{Name: "浪费成本", Value: report.WasteCost, Type: "expense"},
		{Name: "运营成本", Value: report.OperatingCost, Type: "expense"},
	}

	return report, nil
}

func GetOperatingCosts() ([]models.OperatingCost, error) {
	var costs []models.OperatingCost
	err := database.DB.Order("created_at DESC").Find(&costs).Error
	return costs, err
}

func CreateOperatingCost(req *models.OperatingCostCreateRequest) (*models.OperatingCost, error) {
	cost := &models.OperatingCost{
		CostType:    req.CostType,
		CostName:    req.CostName,
		Amount:      req.Amount,
		Period:      req.Period,
		IsFixed:     req.IsFixed,
		Description: req.Description,
	}

	err := database.DB.Create(cost).Error
	if err != nil {
		return nil, err
	}
	return cost, nil
}

func UpdateOperatingCost(id uint, req *models.OperatingCostCreateRequest) (*models.OperatingCost, error) {
	var cost models.OperatingCost
	if err := database.DB.First(&cost, id).Error; err != nil {
		return nil, err
	}

	cost.CostType = req.CostType
	cost.CostName = req.CostName
	cost.Amount = req.Amount
	cost.Period = req.Period
	cost.IsFixed = req.IsFixed
	cost.Description = req.Description

	err := database.DB.Save(&cost).Error
	if err != nil {
		return nil, err
	}
	return &cost, nil
}

func DeleteOperatingCost(id uint) error {
	return database.DB.Delete(&models.OperatingCost{}, id).Error
}

func GetReconciliationLogs(params models.FinanceFilterParams) ([]models.ReconciliationLog, int64, error) {
	var logs []models.ReconciliationLog
	var total int64

	query := database.DB.Table("reconciliation_logs")
	query = applyDateFilter(query, params)
	if params.PaymentMethod != "" {
		query = query.Where("payment_method = ?", params.PaymentMethod)
	}

	query.Count(&total)
	err := query.Order("created_at DESC").Find(&logs).Error
	return logs, total, err
}

func CreateReconciliationLog(log *models.ReconciliationLog) (*models.ReconciliationLog, error) {
	err := database.DB.Create(log).Error
	if err != nil {
		return nil, err
	}
	return log, nil
}

func UpdateReconciliationLog(id uint, status string, remark string) (*models.ReconciliationLog, error) {
	var log models.ReconciliationLog
	if err := database.DB.First(&log, id).Error; err != nil {
		return nil, err
	}

	log.Status = status
	log.Remark = remark

	err := database.DB.Save(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}
