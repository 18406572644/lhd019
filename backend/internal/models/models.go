package models

import (
	"time"
)

type Spirit struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	Name               string    `gorm:"type:varchar(100);not null" json:"name"`
	Category           string    `gorm:"type:varchar(50);not null" json:"category"`
	Brand              string    `gorm:"type:varchar(100)" json:"brand"`
	Origin             string    `gorm:"type:varchar(100)" json:"origin"`
	AlcoholContent     float64   `gorm:"type:decimal(4,1)" json:"alcohol_content"`
	VolumeMl           int       `gorm:"type:int;not null;default:700" json:"volume_ml"`
	Unit               string    `gorm:"type:varchar(20);not null;default:'瓶'" json:"unit"`
	StockQuantity      int       `gorm:"type:int;not null;default:0" json:"stock_quantity"`
	MinStock           int       `gorm:"type:int;not null;default:5" json:"min_stock"`
	CostPrice          float64   `gorm:"type:decimal(10,2);not null" json:"cost_price"`
	SellingPricePerMl  float64   `gorm:"type:decimal(10,4)" json:"selling_price_per_ml"`
	Description        string    `gorm:"type:text" json:"description"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type Ingredient struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"type:varchar(100);not null" json:"name"`
	Category      string    `gorm:"type:varchar(50);not null" json:"category"`
	Unit          string    `gorm:"type:varchar(20);not null" json:"unit"`
	StockQuantity float64   `gorm:"type:decimal(10,2);not null;default:0" json:"stock_quantity"`
	MinStock      float64   `gorm:"type:decimal(10,2);not null;default:0" json:"min_stock"`
	CostPrice     float64   `gorm:"type:decimal(10,2);not null" json:"cost_price"`
	Supplier      string    `gorm:"type:varchar(100)" json:"supplier"`
	Description   string    `gorm:"type:text" json:"description"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Recipe struct {
	ID                 uint              `gorm:"primaryKey" json:"id"`
	Name               string            `gorm:"type:varchar(100);not null" json:"name"`
	Category           string            `gorm:"type:varchar(50);not null" json:"category"`
	GlassType          string            `gorm:"type:varchar(50)" json:"glass_type"`
	ServingMl          int               `gorm:"type:int" json:"serving_ml"`
	Price              float64           `gorm:"type:decimal(10,2);not null" json:"price"`
	Cost               float64           `gorm:"type:decimal(10,2)" json:"cost"`
	PreparationMethod  string            `gorm:"type:text" json:"preparation_method"`
	Garnish            string            `gorm:"type:varchar(200)" json:"garnish"`
	TasteProfile       string            `gorm:"type:varchar(200)" json:"taste_profile"`
	Difficulty         string            `gorm:"type:varchar(20)" json:"difficulty"`
	IsSignature        bool              `gorm:"type:boolean;default:false" json:"is_signature"`
	ImageUrl           string            `gorm:"type:varchar(500)" json:"image_url"`
	Description        string            `gorm:"type:text" json:"description"`
	RecipeIngredients  []RecipeIngredient `gorm:"foreignKey:RecipeID" json:"recipe_ingredients,omitempty"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
}

type RecipeIngredient struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	RecipeID       uint      `gorm:"type:bigint;not null;index" json:"recipe_id"`
	IngredientType string    `gorm:"type:varchar(20);not null" json:"ingredient_type"`
	IngredientID   uint      `gorm:"type:bigint;not null" json:"ingredient_id"`
	Amount         float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	Unit           string    `gorm:"type:varchar(20);not null" json:"unit"`
	CreatedAt      time.Time `json:"created_at"`
}

type Order struct {
	ID            uint         `gorm:"primaryKey" json:"id"`
	OrderNo       string       `gorm:"type:varchar(50);not null;unique" json:"order_no"`
	TableNo       string       `gorm:"type:varchar(20)" json:"table_no"`
	CustomerCount int          `gorm:"type:int;default:1" json:"customer_count"`
	TotalAmount   float64      `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	Discount      float64      `gorm:"type:decimal(10,2);default:0" json:"discount"`
	ActualAmount  float64      `gorm:"type:decimal(10,2);not null" json:"actual_amount"`
	PaymentMethod string       `gorm:"type:varchar(30)" json:"payment_method"`
	Status        string       `gorm:"type:varchar(20);not null;default:'completed'" json:"status"`
	Remark        string       `gorm:"type:text" json:"remark"`
	OrderItems    []OrderItem  `gorm:"foreignKey:OrderID" json:"order_items,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

type OrderItem struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	OrderID    uint      `gorm:"type:bigint;not null;index" json:"order_id"`
	RecipeID   uint      `gorm:"type:bigint;not null;index" json:"recipe_id"`
	RecipeName string    `gorm:"type:varchar(100);not null" json:"recipe_name"`
	Quantity   int       `gorm:"type:int;not null" json:"quantity"`
	UnitPrice  float64   `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	Subtotal   float64   `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	Remark     string    `gorm:"type:text" json:"remark"`
	CreatedAt  time.Time `json:"created_at"`
}

type WasteRecord struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	IngredientType string    `gorm:"type:varchar(20);not null" json:"ingredient_type"`
	IngredientID   uint      `gorm:"type:bigint;not null" json:"ingredient_id"`
	IngredientName string    `gorm:"type:varchar(100);not null" json:"ingredient_name"`
	Amount         float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	Unit           string    `gorm:"type:varchar(20);not null" json:"unit"`
	Reason         string    `gorm:"type:varchar(200);not null" json:"reason"`
	Cost           float64   `gorm:"type:decimal(10,2);not null" json:"cost"`
	Operator       string    `gorm:"type:varchar(50)" json:"operator"`
	Remark         string    `gorm:"type:text" json:"remark"`
	CreatedAt      time.Time `json:"created_at"`
}

type SpecialCreation struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	Name              string    `gorm:"type:varchar(100);not null" json:"name"`
	Creator           string    `gorm:"type:varchar(50)" json:"creator"`
	Inspiration       string    `gorm:"type:text" json:"inspiration"`
	TasteProfile      string    `gorm:"type:varchar(200)" json:"taste_profile"`
	GlassType         string    `gorm:"type:varchar(50)" json:"glass_type"`
	ServingMl         int       `gorm:"type:int" json:"serving_ml"`
	Price             float64   `gorm:"type:decimal(10,2)" json:"price"`
	PreparationMethod string    `gorm:"type:text" json:"preparation_method"`
	Garnish           string    `gorm:"type:varchar(200)" json:"garnish"`
	IngredientsText   string    `gorm:"type:text" json:"ingredients_text"`
	ImageUrl          string    `gorm:"type:varchar(500)" json:"image_url"`
	Status            string    `gorm:"type:varchar(20);default:'draft'" json:"status"`
	TastingNotes      string    `gorm:"type:text" json:"tasting_notes"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type Purchase struct {
	ID             uint            `gorm:"primaryKey" json:"id"`
	PurchaseNo     string          `gorm:"type:varchar(50);not null;unique" json:"purchase_no"`
	Supplier       string          `gorm:"type:varchar(100)" json:"supplier"`
	TotalAmount    float64         `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	PurchaseDate   string          `gorm:"type:date;not null" json:"purchase_date"`
	Operator       string          `gorm:"type:varchar(50)" json:"operator"`
	Remark         string          `gorm:"type:text" json:"remark"`
	PurchaseItems  []PurchaseItem  `gorm:"foreignKey:PurchaseID" json:"purchase_items,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type PurchaseItem struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	PurchaseID     uint      `gorm:"type:bigint;not null;index" json:"purchase_id"`
	IngredientType string    `gorm:"type:varchar(20);not null" json:"ingredient_type"`
	IngredientID   uint      `gorm:"type:bigint;not null" json:"ingredient_id"`
	IngredientName string    `gorm:"type:varchar(100);not null" json:"ingredient_name"`
	Quantity       float64   `gorm:"type:decimal(10,2);not null" json:"quantity"`
	Unit           string    `gorm:"type:varchar(20);not null" json:"unit"`
	UnitPrice      float64   `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	Subtotal       float64   `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	BatchNo        string    `gorm:"type:varchar(100)" json:"batch_no"`
	ExpiryDate     string    `gorm:"type:date" json:"expiry_date"`
	CreatedAt      time.Time `json:"created_at"`
}

type BusinessSummary struct {
	TotalRevenue    float64 `json:"total_revenue"`
	TotalOrders     int64   `json:"total_orders"`
	TotalCustomers  int64   `json:"total_customers"`
	AverageOrder    float64 `json:"average_order"`
	TopDrinks       []TopDrink `json:"top_drinks"`
	DailyRevenue    []DailyRevenue `json:"daily_revenue"`
	LowStockSpirits []Spirit `json:"low_stock_spirits"`
	LowStockIngredients []Ingredient `json:"low_stock_ingredients"`
	TotalWasteCost  float64 `json:"total_waste_cost"`
}

type TopDrink struct {
	RecipeName string `json:"recipe_name"`
	Quantity   int64  `json:"quantity"`
	Revenue    float64 `json:"revenue"`
}

type DailyRevenue struct {
	Date    string  `json:"date"`
	Revenue float64 `json:"revenue"`
	Orders  int64   `json:"orders"`
}

type OrderCreateRequest struct {
	TableNo       string              `json:"table_no"`
	CustomerCount int                 `json:"customer_count"`
	Discount      float64             `json:"discount"`
	PaymentMethod string              `json:"payment_method"`
	Remark        string              `json:"remark"`
	Items         []OrderItemCreate   `json:"items"`
}

type OrderItemCreate struct {
	RecipeID uint   `json:"recipe_id"`
	Quantity int    `json:"quantity"`
	Remark   string `json:"remark"`
}

type PurchaseCreateRequest struct {
	Supplier     string               `json:"supplier"`
	PurchaseDate string               `json:"purchase_date"`
	Operator     string               `json:"operator"`
	Remark       string               `json:"remark"`
	Items        []PurchaseItemCreate `json:"items"`
}

type PurchaseItemCreate struct {
	IngredientType string  `json:"ingredient_type"`
	IngredientID   uint    `json:"ingredient_id"`
	Quantity       float64 `json:"quantity"`
	Unit           string  `json:"unit"`
	UnitPrice      float64 `json:"unit_price"`
	BatchNo        string  `json:"batch_no"`
	ExpiryDate     string  `json:"expiry_date"`
}

type WasteCreateRequest struct {
	IngredientType string  `json:"ingredient_type"`
	IngredientID   uint    `json:"ingredient_id"`
	Amount         float64 `json:"amount"`
	Reason         string  `json:"reason"`
	Operator       string  `json:"operator"`
	Remark         string  `json:"remark"`
}

type OperatingCost struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CostType    string    `gorm:"type:varchar(50);not null" json:"cost_type"`
	CostName    string    `gorm:"type:varchar(100);not null" json:"cost_name"`
	Amount      float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	Period      string    `gorm:"type:varchar(20);not null;default:'monthly'" json:"period"`
	IsFixed     bool      `gorm:"type:boolean;default:true" json:"is_fixed"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RevenueReport struct {
	Period           string             `json:"period"`
	StartDate        string             `json:"start_date"`
	EndDate          string             `json:"end_date"`
	TotalRevenue     float64            `json:"total_revenue"`
	TotalOrders      int64              `json:"total_orders"`
	TotalCustomers   int64              `json:"total_customers"`
	AverageOrder     float64            `json:"average_order"`
	AverageCustomer  float64            `json:"average_customer"`
	YoYGrowth        float64            `json:"yoy_growth"`
	QoQGrowth        float64            `json:"qoq_growth"`
	YoYPrevious      float64            `json:"yoy_previous"`
	QoQPrevious      float64            `json:"qoq_previous"`
	DailyData        []DailyRevenueData `json:"daily_data"`
}

type DailyRevenueData struct {
	Date      string  `json:"date"`
	Revenue   float64 `json:"revenue"`
	Orders    int64   `json:"orders"`
	Customers int64   `json:"customers"`
}

type CostAnalysisReport struct {
	Period              string  `json:"period"`
	StartDate           string  `json:"start_date"`
	EndDate             string  `json:"end_date"`
	TotalRevenue        float64 `json:"total_revenue"`
	IngredientCost      float64 `json:"ingredient_cost"`
	SpiritCost          float64 `json:"spirit_cost"`
	TotalMaterialCost   float64 `json:"total_material_cost"`
	WasteCost           float64 `json:"waste_cost"`
	PurchaseCost        float64 `json:"purchase_cost"`
	OperatingCost       float64 `json:"operating_cost"`
	TotalCost           float64 `json:"total_cost"`
	GrossProfit         float64 `json:"gross_profit"`
	GrossMargin         float64 `json:"gross_margin"`
	NetProfit           float64 `json:"net_profit"`
	NetMargin           float64 `json:"net_margin"`
	CostBreakdown       []CostBreakdownItem `json:"cost_breakdown"`
}

type CostBreakdownItem struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Ratio float64 `json:"ratio"`
}

type CategorySalesReport struct {
	Period        string          `json:"period"`
	StartDate     string          `json:"start_date"`
	EndDate       string          `json:"end_date"`
	CategorySales []CategorySales `json:"category_sales"`
	RecipeSales   []RecipeSales   `json:"recipe_sales"`
	TimeSlotSales []TimeSlotSales `json:"time_slot_sales"`
}

type CategorySales struct {
	Category   string  `json:"category"`
	Quantity   int64   `json:"quantity"`
	Revenue    float64 `json:"revenue"`
	Percentage float64 `json:"percentage"`
}

type RecipeSales struct {
	RecipeID     uint    `json:"recipe_id"`
	RecipeName   string  `json:"recipe_name"`
	Category     string  `json:"category"`
	Quantity     int64   `json:"quantity"`
	Revenue      float64 `json:"revenue"`
	Cost         float64 `json:"cost"`
	Profit       float64 `json:"profit"`
	ProfitMargin float64 `json:"profit_margin"`
}

type TimeSlotSales struct {
	TimeSlot string  `json:"time_slot"`
	Quantity int64   `json:"quantity"`
	Revenue  float64 `json:"revenue"`
	Orders   int64   `json:"orders"`
}

type PaymentReconciliation struct {
	Period             string                  `json:"period"`
	StartDate          string                  `json:"start_date"`
	EndDate            string                  `json:"end_date"`
	TotalRevenue       float64                 `json:"total_revenue"`
	PaymentMethods     []PaymentMethodDetail   `json:"payment_methods"`
	ReconciliationLogs []ReconciliationLog     `json:"reconciliation_logs"`
}

type PaymentMethodDetail struct {
	PaymentMethod string  `json:"payment_method"`
	OrderCount    int64   `json:"order_count"`
	TotalAmount   float64 `json:"total_amount"`
	Percentage    float64 `json:"percentage"`
}

type ReconciliationLog struct {
	ID            uint      `json:"id"`
	OrderNo       string    `json:"order_no"`
	PaymentMethod string    `json:"payment_method"`
	SystemAmount  float64   `json:"system_amount"`
	ActualAmount  float64   `json:"actual_amount"`
	Difference    float64   `json:"difference"`
	Status        string    `json:"status"`
	ReconciledAt  string    `json:"reconciled_at"`
	Remark        string    `json:"remark"`
}

type ProfitReport struct {
	Period            string  `json:"period"`
	StartDate         string  `json:"start_date"`
	EndDate           string  `json:"end_date"`
	TotalRevenue      float64 `json:"total_revenue"`
	MaterialCost      float64 `json:"material_cost"`
	WasteCost         float64 `json:"waste_cost"`
	OperatingCost     float64 `json:"operating_cost"`
	TotalExpenses     float64 `json:"total_expenses"`
	GrossProfit       float64 `json:"gross_profit"`
	GrossMargin       float64 `json:"gross_margin"`
	NetProfit         float64 `json:"net_profit"`
	NetMargin         float64 `json:"net_margin"`
	ProfitBreakdown   []ProfitBreakdownItem `json:"profit_breakdown"`
}

type ProfitBreakdownItem struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Type  string  `json:"type"`
}

type FinanceFilterParams struct {
	StartDate    string `form:"start_date"`
	EndDate      string `form:"end_date"`
	Period       string `form:"period"`
	Category     string `form:"category"`
	PaymentMethod string `form:"payment_method"`
}

type CustomReportConfig struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ReportName  string    `gorm:"type:varchar(100);not null" json:"report_name"`
	ReportType  string    `gorm:"type:varchar(50);not null" json:"report_type"`
	Config      string    `gorm:"type:text;not null" json:"config"`
	CreatedBy   string    `gorm:"type:varchar(50)" json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type OperatingCostCreateRequest struct {
	CostType    string  `json:"cost_type"`
	CostName    string  `json:"cost_name"`
	Amount      float64 `json:"amount"`
	Period      string  `json:"period"`
	IsFixed     bool    `json:"is_fixed"`
	Description string  `json:"description"`
}

type StockBatch struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	BatchCode        string    `gorm:"type:varchar(50);not null;unique" json:"batch_code"`
	BatchNo          string    `gorm:"type:varchar(100)" json:"batch_no"`
	IngredientType   string    `gorm:"type:varchar(20);not null" json:"ingredient_type"`
	IngredientID     uint      `gorm:"type:bigint;not null" json:"ingredient_id"`
	IngredientName   string    `gorm:"type:varchar(100);not null" json:"ingredient_name"`
	PurchaseItemID   *uint     `gorm:"type:bigint" json:"purchase_item_id"`
	TotalQuantity    float64   `gorm:"type:decimal(10,2);not null" json:"total_quantity"`
	RemainingQuantity float64  `gorm:"type:decimal(10,2);not null" json:"remaining_quantity"`
	Unit             string    `gorm:"type:varchar(20);not null" json:"unit"`
	UnitPrice        float64   `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	ExpiryDate       string    `gorm:"type:date" json:"expiry_date"`
	IsPromotion      bool      `gorm:"type:boolean;default:false" json:"is_promotion"`
	Status           string    `gorm:"type:varchar(20);not null;default:'normal'" json:"status"`
	WarehousePosition string   `gorm:"type:varchar(50)" json:"warehouse_position"`
	Remark           string    `gorm:"type:text" json:"remark"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type BatchOutRecord struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	BatchID        uint      `gorm:"type:bigint;not null;index" json:"batch_id"`
	BatchCode      string    `gorm:"type:varchar(50);not null" json:"batch_code"`
	IngredientType string    `gorm:"type:varchar(20);not null" json:"ingredient_type"`
	IngredientID   uint      `gorm:"type:bigint;not null" json:"ingredient_id"`
	IngredientName string    `gorm:"type:varchar(100);not null" json:"ingredient_name"`
	OutType        string    `gorm:"type:varchar(20);not null" json:"out_type"`
	OutQuantity    float64   `gorm:"type:decimal(10,2);not null" json:"out_quantity"`
	Unit           string    `gorm:"type:varchar(20);not null" json:"unit"`
	UnitPrice      float64   `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	TotalCost      float64   `gorm:"type:decimal(10,2);not null" json:"total_cost"`
	OrderID        *uint     `gorm:"type:bigint;index" json:"order_id"`
	OrderNo        string    `gorm:"type:varchar(50)" json:"order_no"`
	WasteID        *uint     `gorm:"type:bigint;index" json:"waste_id"`
	Operator       string    `gorm:"type:varchar(50)" json:"operator"`
	Remark         string    `gorm:"type:text" json:"remark"`
	CreatedAt      time.Time `json:"created_at"`
}

type Stocktake struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	StocktakeNo    string         `gorm:"type:varchar(50);not null;unique" json:"stocktake_no"`
	StocktakeDate  string         `gorm:"type:date;not null" json:"stocktake_date"`
	StocktakeType  string         `gorm:"type:varchar(20);not null;default:'periodic'" json:"stocktake_type"`
	Status         string         `gorm:"type:varchar(20);not null;default:'draft'" json:"status"`
	Operator       string         `gorm:"type:varchar(50)" json:"operator"`
	TotalProfit    float64        `gorm:"type:decimal(10,2);default:0" json:"total_profit"`
	TotalLoss      float64        `gorm:"type:decimal(10,2);default:0" json:"total_loss"`
	TotalDiff      float64        `gorm:"type:decimal(10,2);default:0" json:"total_diff"`
	Remark         string         `gorm:"type:text" json:"remark"`
	ConfirmedAt    *time.Time     `json:"confirmed_at"`
	StocktakeItems []StocktakeItem `gorm:"foreignKey:StocktakeID" json:"stocktake_items,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type StocktakeItem struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	StocktakeID    uint      `gorm:"type:bigint;not null;index" json:"stocktake_id"`
	IngredientType string    `gorm:"type:varchar(20);not null" json:"ingredient_type"`
	IngredientID   uint      `gorm:"type:bigint;not null" json:"ingredient_id"`
	IngredientName string    `gorm:"type:varchar(100);not null" json:"ingredient_name"`
	SystemQuantity float64   `gorm:"type:decimal(10,2);not null" json:"system_quantity"`
	ActualQuantity float64   `gorm:"type:decimal(10,2);not null" json:"actual_quantity"`
	DiffQuantity   float64   `gorm:"type:decimal(10,2);not null" json:"diff_quantity"`
	Unit           string    `gorm:"type:varchar(20);not null" json:"unit"`
	UnitPrice      float64   `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	DiffAmount     float64   `gorm:"type:decimal(10,2);not null" json:"diff_amount"`
	DiffType       string    `gorm:"type:varchar(20)" json:"diff_type"`
	Remark         string    `gorm:"type:text" json:"remark"`
	CreatedAt      time.Time `json:"created_at"`
}

type BatchDeductResult struct {
	BatchID     uint    `json:"batch_id"`
	BatchCode   string  `json:"batch_code"`
	DeductQty   float64 `json:"deduct_qty"`
	RemainingQty float64 `json:"remaining_qty"`
	UnitPrice   float64 `json:"unit_price"`
}

type StockBatchCreateRequest struct {
	IngredientType   string  `json:"ingredient_type"`
	IngredientID     uint    `json:"ingredient_id"`
	IngredientName   string  `json:"ingredient_name"`
	BatchNo          string  `json:"batch_no"`
	TotalQuantity    float64 `json:"total_quantity"`
	Unit             string  `json:"unit"`
	UnitPrice        float64 `json:"unit_price"`
	ExpiryDate       string  `json:"expiry_date"`
	WarehousePosition string `json:"warehouse_position"`
	Remark           string  `json:"remark"`
}

type BatchDeductRequest struct {
	IngredientType string  `json:"ingredient_type"`
	IngredientID   uint    `json:"ingredient_id"`
	Quantity       float64 `json:"quantity"`
	OutType        string  `json:"out_type"`
	OrderID        uint    `json:"order_id"`
	OrderNo        string  `json:"order_no"`
	Operator       string  `json:"operator"`
	Remark         string  `json:"remark"`
}

type BatchTraceResult struct {
	StockBatch  StockBatch      `json:"stock_batch"`
	OutRecords  []BatchOutRecord `json:"out_records"`
	TotalOutQty float64         `json:"total_out_qty"`
}

type StocktakeCreateRequest struct {
	StocktakeDate string                  `json:"stocktake_date"`
	StocktakeType string                  `json:"stocktake_type"`
	Operator      string                  `json:"operator"`
	Remark        string                  `json:"remark"`
	Items         []StocktakeItemCreate   `json:"items"`
}

type StocktakeItemCreate struct {
	IngredientType string  `json:"ingredient_type"`
	IngredientID   uint    `json:"ingredient_id"`
	ActualQuantity float64 `json:"actual_quantity"`
	Remark         string  `json:"remark"`
}

type StocktakeConfirmRequest struct {
	Status string `json:"status"`
	Remark string `json:"remark"`
}

type ExpiryWarningResult struct {
	StockBatch  StockBatch `json:"stock_batch"`
	DaysToExpiry int       `json:"days_to_expiry"`
	WarningLevel string    `json:"warning_level"`
}
