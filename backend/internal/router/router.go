package router

import (
	"cocktail-bar-system/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		spirits := api.Group("/spirits")
		{
			spirits.GET("", handlers.GetSpirits)
			spirits.GET("/:id", handlers.GetSpirit)
			spirits.POST("", handlers.CreateSpirit)
			spirits.PUT("/:id", handlers.UpdateSpirit)
			spirits.DELETE("/:id", handlers.DeleteSpirit)
			spirits.GET("/low-stock", handlers.GetLowStockSpirits)
		}

		ingredients := api.Group("/ingredients")
		{
			ingredients.GET("", handlers.GetIngredients)
			ingredients.GET("/:id", handlers.GetIngredient)
			ingredients.POST("", handlers.CreateIngredient)
			ingredients.PUT("/:id", handlers.UpdateIngredient)
			ingredients.DELETE("/:id", handlers.DeleteIngredient)
			ingredients.GET("/low-stock", handlers.GetLowStockIngredients)
		}

		recipes := api.Group("/recipes")
		{
			recipes.GET("", handlers.GetRecipes)
			recipes.GET("/:id", handlers.GetRecipe)
		recipes.POST("", handlers.CreateRecipe)
			recipes.PUT("/:id", handlers.UpdateRecipe)
			recipes.DELETE("/:id", handlers.DeleteRecipe)
		}

		orders := api.Group("/orders")
		{
			orders.GET("", handlers.GetOrders)
			orders.GET("/:id", handlers.GetOrder)
			orders.POST("", handlers.CreateOrder)
			orders.PUT("/:id", handlers.UpdateOrder)
			orders.DELETE("/:id", handlers.DeleteOrder)
		}

		waste := api.Group("/waste")
		{
			waste.GET("", handlers.GetWasteRecords)
			waste.GET("/:id", handlers.GetWasteRecord)
			waste.POST("", handlers.CreateWasteRecord)
			waste.DELETE("/:id", handlers.DeleteWasteRecord)
		}

		specials := api.Group("/specials")
		{
			specials.GET("", handlers.GetSpecialCreations)
			specials.GET("/:id", handlers.GetSpecialCreation)
			specials.POST("", handlers.CreateSpecialCreation)
			specials.PUT("/:id", handlers.UpdateSpecialCreation)
			specials.DELETE("/:id", handlers.DeleteSpecialCreation)
		}

		purchases := api.Group("/purchases")
		{
			purchases.GET("", handlers.GetPurchases)
			purchases.GET("/:id", handlers.GetPurchase)
			purchases.POST("", handlers.CreatePurchase)
			purchases.DELETE("/:id", handlers.DeletePurchase)
		}

		summary := api.Group("/summary")
		{
			summary.GET("", handlers.GetBusinessSummary)
		}

		finance := api.Group("/finance")
		{
			finance.GET("/revenue", handlers.GetRevenueReport)
			finance.GET("/cost-analysis", handlers.GetCostAnalysisReport)
			finance.GET("/category-sales", handlers.GetCategorySalesReport)
			finance.GET("/payment-reconciliation", handlers.GetPaymentReconciliation)
			finance.GET("/profit", handlers.GetProfitReport)
			finance.GET("/operating-costs", handlers.GetOperatingCosts)
			finance.POST("/operating-costs", handlers.CreateOperatingCost)
			finance.PUT("/operating-costs/:id", handlers.UpdateOperatingCost)
			finance.DELETE("/operating-costs/:id", handlers.DeleteOperatingCost)
			finance.GET("/reconciliation-logs", handlers.GetReconciliationLogs)
			finance.POST("/reconciliation-logs", handlers.CreateReconciliationLog)
			finance.PUT("/reconciliation-logs/:id", handlers.UpdateReconciliationLog)
		}
	}
}
