package handlers

import (
	"net/http"
	"strconv"

	"cocktail-bar-system/internal/models"
	"cocktail-bar-system/internal/services"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Total   int64       `json:"total,omitempty"`
}

func successResponse(c *gin.Context, data interface{}, total ...int64) {
	resp := Response{
		Code:    0,
		Message: "success",
		Data:    data,
	}
	if len(total) > 0 {
		resp.Total = total[0]
	}
	c.JSON(http.StatusOK, resp)
}

func errorResponse(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

func GetSpirits(c *gin.Context) {
	category := c.Query("category")
	keyword := c.Query("keyword")
	spirits, total, err := services.GetSpirits(category, keyword)
	if err != nil {
		errorResponse(c, 500, err.Error())
		return
	}
	successResponse(c, spirits, total)
}

func GetSpirit(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	spirit, err := services.GetSpirit(uint(id))
	if err != nil {
		errorResponse(c, 404, "Spirit not found")
		return
	}
	successResponse(c, spirit)
}

func CreateSpirit(c *gin.Context) {
	var spirit models.Spirit
	if err := c.ShouldBindJSON(&spirit); err != nil {
		errorResponse(c, 400, err.Error())
		return
	}
	if err := services.CreateSpirit(&spirit); err != nil {
		errorResponse(c, 500, err.Error())
		return
	}
	successResponse(c, spirit)
}

func UpdateSpirit(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var spirit models.Spirit
	if err := c.ShouldBindJSON(&spirit); err != nil {
		errorResponse(c, 400, err.Error())
		return
	}
	if err := services.UpdateSpirit(uint(id), &spirit); err != nil {
		errorResponse(c, 500, err.Error())
		return
	}
	successResponse(c, spirit)
}

func DeleteSpirit(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := services.DeleteSpirit(uint(id)); err != nil {
		errorResponse(c, 500, err.Error())
		return
	}
	successResponse(c, nil)
}

func GetLowStockSpirits(c *gin.Context) {
	spirits, err := services.GetLowStockSpirits()
	if err != nil {
		errorResponse(c, 500, err.Error())
		return
	}
	successResponse(c, spirits)
}

func GetIngredients(c *gin.Context) {
	category := c.Query("category")
	keyword := c.Query("keyword")
	ingredients, total, err := services.GetIngredients(category, keyword)
	if err != nil {
		errorResponse(c, 500, err.Error())
		return
	}
	successResponse(c, ingredients, total)
}

func GetIngredient(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	ingredient, err := services.GetIngredient(uint(id))
	if err != nil {
		errorResponse(c, 404, "Ingredient not found")
		return
	}
	successResponse(c, ingredient)
}

func CreateIngredient(c *gin.Context) {
	var ingredient models.Ingredient
	if err := c.ShouldBindJSON(&ingredient); err != nil {
		errorResponse(c, 400, err.Error())
		return
	}
	if err := services.CreateIngredient(&ingredient); err != nil {
		errorResponse(c, 500, err.Error())
		return
	}
	successResponse(c, ingredient)
}

func UpdateIngredient(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var ingredient models.Ingredient
	if err := c.ShouldBindJSON(&ingredient); err != nil {
		errorResponse(c, 400, err.Error())
		return
	}
	if err := services.UpdateIngredient(uint(id), &ingredient); err != nil {
		errorResponse(c, 500, err.Error())
		return
	}
	successResponse(c, ingredient)
}

func DeleteIngredient(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := services.DeleteIngredient(uint(id)); err != nil {
		errorResponse(c, 500, err.Error())
		return
	}
	successResponse(c, nil)
}

func GetLowStockIngredients(c *gin.Context) {
	ingredients, err := services.GetLowStockIngredients()
	if err != nil {
		errorResponse(c, 500, err.Error())
		return
	}
	successResponse(c, ingredients)
}

func GetRecipes(c *gin.Context) {
	category := c.Query("category")
	keyword := c.Query("keyword")
	isSignatureStr := c.Query("is_signature")
	var isSignature *bool
	if isSignatureStr != "" {
		b := isSignatureStr == "true"
		isSignature = &b
	}
	recipes, total, err := services.GetRecipes(category, keyword, isSignature)
	if err != nil {
		errorResponse(c, 500, err.Error())
		return
	}
	successResponse(c, recipes, total)
}

func GetRecipe(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	recipe, err := services.GetRecipe(uint(id))
	if err != nil {
		errorResponse(c, 404, "Recipe not found")
		return
	}
	successResponse(c, recipe)
}

func CreateRecipe(c *gin.Context) {
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		errorResponse(c, 400, err.Error())
		return
	}
	if err := services.CreateRecipe(&recipe); err != nil {
		errorResponse(c, 500, err.Error())
		return
	}
	successResponse(c, recipe)
}

func UpdateRecipe(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		errorResponse(c, 400, err.Error())
		return
	}
	if err := services.UpdateRecipe(uint(id), &recipe); err != nil {
		errorResponse(c, 500, err.Error())
		return
	}
	successResponse(c, recipe)
}

func DeleteRecipe(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := services.DeleteRecipe(uint(id)); err != nil {
		errorResponse(c, 500, err.Error())
		return
	}
	successResponse(c, nil)
}

func GetOrders(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	status := c.Query("status")
	orders, total, err := services.GetOrders(startDate, endDate, status)
	if err != nil {
		errorResponse(c, 500, err.Error())
		return
	}
	successResponse(c, orders, total)
}

func GetOrder(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	order, err := services.GetOrder(uint(id))
	if err != nil {
		errorResponse(c, 404, "Order not found")
		return
	}
	successResponse(c, order)
}

func CreateOrder(c *gin.Context) {
	var req models.OrderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, 400, err.Error())
		return
	}
	order, err := services.CreateOrder(&req)
	if err != nil {
		errorResponse(c, 500, err.Error())
		return
	}
	successResponse(c, order)
}

func UpdateOrder(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req models.OrderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, 400, err.Error())
		return
	}
	order, err := services.UpdateOrder(uint(id), &req)
	if err != nil {
		errorResponse(c, 500, err.Error())
		return
	}
	successResponse(c, order)
}

func DeleteOrder(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := services.DeleteOrder(uint(id)); err != nil {
		errorResponse(c, 500, err.Error())
		return
	}
	successResponse(c, nil)
}

func GetWasteRecords(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	ingredientType := c.Query("ingredient_type")
	records, total, err := services.GetWasteRecords(startDate, endDate, ingredientType)
	if err != nil {
		errorResponse(c, 500, err.Error())
		return
	}
	successResponse(c, records, total)
}

func GetWasteRecord(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	record, err := services.GetWasteRecord(uint(id))
	if err != nil {
		errorResponse(c, 404, "Waste record not found")
		return
	}
	successResponse(c, record)
}

func CreateWasteRecord(c *gin.Context) {
	var req models.WasteCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, 400, err.Error())
		return
	}
	record, err := services.CreateWasteRecord(&req)
	if err != nil {
		errorResponse(c, 500, err.Error())
		return
	}
	successResponse(c, record)
}

func DeleteWasteRecord(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := services.DeleteWasteRecord(uint(id)); err != nil {
		errorResponse(c, 500, err.Error())
		return
	}
	successResponse(c, nil)
}

func GetSpecialCreations(c *gin.Context) {
	status := c.Query("status")
	keyword := c.Query("keyword")
	specials, total, err := services.GetSpecialCreations(status, keyword)
	if err != nil {
		errorResponse(c, 500, err.Error())
		return
	}
	successResponse(c, specials, total)
}

func GetSpecialCreation(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	special, err := services.GetSpecialCreation(uint(id))
	if err != nil {
		errorResponse(c, 404, "Special creation not found")
		return
	}
	successResponse(c, special)
}

func CreateSpecialCreation(c *gin.Context) {
	var special models.SpecialCreation
	if err := c.ShouldBindJSON(&special); err != nil {
		errorResponse(c, 400, err.Error())
		return
	}
	if err := services.CreateSpecialCreation(&special); err != nil {
		errorResponse(c, 500, err.Error())
		return
	}
	successResponse(c, special)
}

func UpdateSpecialCreation(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var special models.SpecialCreation
	if err := c.ShouldBindJSON(&special); err != nil {
		errorResponse(c, 400, err.Error())
		return
	}
	if err := services.UpdateSpecialCreation(uint(id), &special); err != nil {
		errorResponse(c, 500, err.Error())
		return
	}
	successResponse(c, special)
}

func DeleteSpecialCreation(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := services.DeleteSpecialCreation(uint(id)); err != nil {
		errorResponse(c, 500, err.Error())
		return
	}
	successResponse(c, nil)
}

func GetPurchases(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	supplier := c.Query("supplier")
	purchases, total, err := services.GetPurchases(startDate, endDate, supplier)
	if err != nil {
		errorResponse(c, 500, err.Error())
		return
	}
	successResponse(c, purchases, total)
}

func GetPurchase(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	purchase, err := services.GetPurchase(uint(id))
	if err != nil {
		errorResponse(c, 404, "Purchase not found")
		return
	}
	successResponse(c, purchase)
}

func CreatePurchase(c *gin.Context) {
	var req models.PurchaseCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, 400, err.Error())
		return
	}
	purchase, err := services.CreatePurchase(&req)
	if err != nil {
		errorResponse(c, 500, err.Error())
		return
	}
	successResponse(c, purchase)
}

func DeletePurchase(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := services.DeletePurchase(uint(id)); err != nil {
		errorResponse(c, 500, err.Error())
		return
	}
	successResponse(c, nil)
}

func GetBusinessSummary(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	summary, err := services.GetBusinessSummary(startDate, endDate)
	if err != nil {
		errorResponse(c, 500, err.Error())
		return
	}
	successResponse(c, summary)
}
