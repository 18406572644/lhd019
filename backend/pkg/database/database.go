package database

import (
	"log"
	"time"

	"cocktail-bar-system/internal/config"
	"cocktail-bar-system/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dsn := config.GetDSN()

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = DB.AutoMigrate(
		&models.Spirit{},
		&models.Ingredient{},
		&models.Recipe{},
		&models.RecipeIngredient{},
		&models.Order{},
		&models.OrderItem{},
		&models.WasteRecord{},
		&models.SpecialCreation{},
		&models.Purchase{},
		&models.PurchaseItem{},
		&models.OperatingCost{},
		&models.ReconciliationLog{},
		&models.CustomReportConfig{},
		&models.StockBatch{},
		&models.BatchOutRecord{},
		&models.Stocktake{},
		&models.StocktakeItem{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database migration completed successfully")
	log.Println("Database connection established successfully")
}

func GetDB() *gorm.DB {
	return DB
}
