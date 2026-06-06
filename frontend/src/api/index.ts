import request from './request'

export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
  total?: number
}

export interface Spirit {
  id: number
  name: string
  category: string
  brand: string
  origin: string
  alcohol_content: number
  volume_ml: number
  unit: string
  stock_quantity: number
  min_stock: number
  cost_price: number
  selling_price_per_ml: number
  description: string
  created_at: string
  updated_at: string
}

export type SpiritForm = Omit<Spirit, 'id' | 'created_at' | 'updated_at'>

export interface RecipeIngredient {
  id?: number
  recipe_id?: number
  ingredient_type: 'spirit' | 'ingredient'
  ingredient_id: number
  ingredient_name?: string
  amount: number
  unit: string
  cost?: number
  created_at?: string
  updated_at?: string
}

export interface Recipe {
  id: number
  name: string
  category: string
  glass_type: string
  serving_ml: number
  price: number
  cost: number
  difficulty: string
  is_signature: boolean
  taste_profile: string
  garnish: string
  description?: string
  preparation_method?: string
  image_url?: string
  recipe_ingredients: RecipeIngredient[]
  created_at: string
  updated_at: string
}

export type RecipeForm = Omit<Recipe, 'id' | 'created_at' | 'updated_at'>

export interface Ingredient {
  id: number
  name: string
  category: string
  stock_quantity: number
  unit: string
  min_stock: number
  cost_price: number
  supplier: string
  remark?: string
  created_at?: string
  updated_at?: string
}

export interface WasteRecord {
  id: number
  ingredient_id: number
  ingredient_name: string
  ingredient_type: 'spirit' | 'ingredient'
  amount: number
  unit: string
  reason: string
  cost: number
  operator: string
  remark?: string
  created_at: string
}

export interface WasteRecordForm {
  ingredient_type: 'spirit' | 'ingredient'
  ingredient_id: number
  amount: number
  reason: string
  operator: string
  remark?: string
}

export interface SpecialCreation {
  id: number
  name: string
  creator: string
  inspiration: string
  taste_profile: string
  glass_type: string
  serving_ml: number
  price: number
  preparation_method: string
  garnish: string
  ingredients_text: string
  image_url: string
  status: 'draft' | 'testing' | 'approved'
  tasting_notes: string
  created_at: string
  updated_at: string
}

export type SpecialCreationForm = Omit<SpecialCreation, 'id' | 'created_at' | 'updated_at'>

export interface PurchaseItem {
  id: number
  purchase_id: number
  ingredient_type: 'spirit' | 'ingredient'
  ingredient_id: number
  ingredient_name: string
  quantity: number
  unit: string
  unit_price: number
  subtotal: number
  batch_no: string
  expiry_date: string
  created_at: string
}

export interface Purchase {
  id: number
  purchase_no: string
  supplier: string
  total_amount: number
  purchase_date: string
  operator: string
  remark: string
  purchase_items?: PurchaseItem[]
  created_at: string
  updated_at: string
}

export interface PurchaseItemForm {
  ingredient_type: 'spirit' | 'ingredient'
  ingredient_id: number
  quantity: number
  unit: string
  unit_price: number
  batch_no: string
  expiry_date: string
}

export interface PurchaseForm {
  supplier: string
  purchase_date: string
  operator: string
  remark: string
  items: PurchaseItemForm[]
}

export interface OrderItem {
  id: number
  order_id: number
  recipe_id: number
  recipe_name: string
  quantity: number
  unit_price: number
  subtotal: number
  remark?: string
  created_at: string
}

export interface Order {
  id: number
  order_no: string
  table_no: string
  customer_count: number
  total_amount: number
  discount: number
  actual_amount: number
  payment_method: string
  status: 'completed' | 'pending' | 'cancelled'
  remark?: string
  order_items: OrderItem[]
  created_at: string
  updated_at: string
}

export interface OrderItemForm {
  recipe_id: number
  quantity: number
  remark?: string
}

export interface OrderForm {
  table_no: string
  customer_count: number
  discount: number
  payment_method: string
  remark?: string
  items: OrderItemForm[]
}

export interface DailyRevenueData {
  date: string
  revenue: number
  orders: number
  customers: number
}

export interface RevenueReport {
  period: string
  start_date: string
  end_date: string
  total_revenue: number
  total_orders: number
  total_customers: number
  average_order: number
  average_customer: number
  yoy_growth: number
  qoq_growth: number
  yoy_previous: number
  qoq_previous: number
  daily_data?: DailyRevenueData[]
}

export interface CostBreakdownItem {
  name: string
  value: number
  ratio: number
}

export interface CostAnalysisReport {
  period: string
  start_date: string
  end_date: string
  total_revenue: number
  ingredient_cost: number
  spirit_cost: number
  total_material_cost: number
  waste_cost: number
  purchase_cost: number
  operating_cost: number
  total_cost: number
  gross_profit: number
  gross_margin: number
  net_profit: number
  net_margin: number
  cost_breakdown: CostBreakdownItem[]
}

export interface CategorySales {
  category: string
  quantity: number
  revenue: number
  percentage: number
}

export interface RecipeSales {
  recipe_id: number
  recipe_name: string
  category: string
  quantity: number
  revenue: number
  cost: number
  profit: number
  profit_margin: number
}

export interface TimeSlotSales {
  time_slot: string
  quantity: number
  revenue: number
  orders: number
}

export interface CategorySalesReport {
  period: string
  start_date: string
  end_date: string
  category_sales: CategorySales[]
  recipe_sales: RecipeSales[]
  time_slot_sales: TimeSlotSales[]
}

export interface PaymentMethodDetail {
  payment_method: string
  order_count: number
  total_amount: number
  percentage: number
}

export interface ReconciliationLog {
  id: number
  order_no: string
  payment_method: string
  system_amount: number
  actual_amount: number
  difference: number
  status: string
  reconciled_at: string
  remark: string
  created_at: string
}

export interface PaymentReconciliation {
  period: string
  start_date: string
  end_date: string
  total_revenue: number
  payment_methods: PaymentMethodDetail[]
  reconciliation_logs: ReconciliationLog[]
}

export interface ProfitBreakdownItem {
  name: string
  value: number
  type: string
}

export interface ProfitReport {
  period: string
  start_date: string
  end_date: string
  total_revenue: number
  material_cost: number
  waste_cost: number
  operating_cost: number
  total_expenses: number
  gross_profit: number
  gross_margin: number
  net_profit: number
  net_margin: number
  profit_breakdown: ProfitBreakdownItem[]
}

export interface OperatingCost {
  id: number
  cost_type: string
  cost_name: string
  amount: number
  period: string
  is_fixed: boolean
  description: string
  created_at: string
  updated_at: string
}

export interface OperatingCostForm {
  cost_type: string
  cost_name: string
  amount: number
  period: string
  is_fixed: boolean
  description: string
}

export interface FinanceFilterParams {
  start_date?: string
  end_date?: string
  period?: string
  category?: string
  payment_method?: string
}

export const api = {
  getSpirits: (params?: { category?: string; keyword?: string }) =>
    request.get<ApiResponse<Spirit[]>>('/spirits', { params }),
  getSpirit: (id: number) => request.get<ApiResponse<Spirit>>(`/spirits/${id}`),
  createSpirit: (data: SpiritForm) => request.post<ApiResponse<Spirit>>('/spirits', data),
  updateSpirit: (id: number, data: Partial<SpiritForm>) => request.put<ApiResponse<Spirit>>(`/spirits/${id}`, data),
  deleteSpirit: (id: number) => request.delete<ApiResponse<void>>(`/spirits/${id}`),
  getLowStockSpirits: () => request.get<ApiResponse<Spirit[]>>('/spirits/low-stock'),

  getIngredients: (params?: { category?: string; keyword?: string }) =>
    request.get<ApiResponse<Ingredient[]>>('/ingredients', { params }),
  getIngredient: (id: number) => request.get<ApiResponse<Ingredient>>(`/ingredients/${id}`),
  createIngredient: (data: Ingredient) => request.post<ApiResponse<Ingredient>>('/ingredients', data),
  updateIngredient: (id: number, data: Partial<Ingredient>) => request.put<ApiResponse<Ingredient>>(`/ingredients/${id}`, data),
  deleteIngredient: (id: number) => request.delete<ApiResponse<void>>(`/ingredients/${id}`),
  getLowStockIngredients: () => request.get<ApiResponse<Ingredient[]>>('/ingredients/low-stock'),

  getRecipes: (params?: { category?: string; keyword?: string; is_signature?: boolean }) =>
    request.get<ApiResponse<Recipe[]>>('/recipes', { params }),
  getRecipe: (id: number) => request.get<ApiResponse<Recipe>>(`/recipes/${id}`),
  createRecipe: (data: RecipeForm) => request.post<ApiResponse<Recipe>>('/recipes', data),
  updateRecipe: (id: number, data: Partial<RecipeForm>) => request.put<ApiResponse<Recipe>>(`/recipes/${id}`, data),
  deleteRecipe: (id: number) => request.delete<ApiResponse<void>>(`/recipes/${id}`),

  getOrders: (params?: { start_date?: string; end_date?: string; status?: string }) =>
    request.get<ApiResponse<Order[]>>('/orders', { params }),
  getOrder: (id: number) => request.get<ApiResponse<Order>>(`/orders/${id}`),
  createOrder: (data: OrderForm) => request.post<ApiResponse<Order>>('/orders', data),
  deleteOrder: (id: number) => request.delete<ApiResponse<void>>(`/orders/${id}`),

  getWasteRecords: (params?: { start_date?: string; end_date?: string; ingredient_type?: string }) =>
    request.get<ApiResponse<WasteRecord[]>>('/waste', { params }),
  createWasteRecord: (data: WasteRecordForm) => request.post<ApiResponse<WasteRecord>>('/waste', data),
  deleteWasteRecord: (id: number) => request.delete<ApiResponse<void>>(`/waste/${id}`),

  getSpecials: (params?: { status?: string; keyword?: string }) =>
    request.get<ApiResponse<SpecialCreation[]>>('/specials', { params }),
  getSpecial: (id: number) => request.get<ApiResponse<SpecialCreation>>(`/specials/${id}`),
  createSpecial: (data: SpecialCreationForm) => request.post<ApiResponse<SpecialCreation>>('/specials', data),
  updateSpecial: (id: number, data: Partial<SpecialCreationForm>) => request.put<ApiResponse<SpecialCreation>>(`/specials/${id}`, data),
  deleteSpecial: (id: number) => request.delete<ApiResponse<void>>(`/specials/${id}`),

  getPurchases: (params?: { start_date?: string; end_date?: string; supplier?: string }) =>
    request.get<ApiResponse<Purchase[]>>('/purchases', { params }),
  getPurchase: (id: number) => request.get<ApiResponse<Purchase>>(`/purchases/${id}`),
  createPurchase: (data: PurchaseForm) => request.post<ApiResponse<Purchase>>('/purchases', data),
  deletePurchase: (id: number) => request.delete<ApiResponse<void>>(`/purchases/${id}`),

  getSummary: (params?: { start_date?: string; end_date?: string }) =>
    request.get<ApiResponse<any>>('/summary', { params }),

  getRevenueReport: (params?: FinanceFilterParams) =>
    request.get<ApiResponse<RevenueReport>>('/finance/revenue', { params }),
  getCostAnalysisReport: (params?: FinanceFilterParams) =>
    request.get<ApiResponse<CostAnalysisReport>>('/finance/cost-analysis', { params }),
  getCategorySalesReport: (params?: FinanceFilterParams) =>
    request.get<ApiResponse<CategorySalesReport>>('/finance/category-sales', { params }),
  getPaymentReconciliation: (params?: FinanceFilterParams) =>
    request.get<ApiResponse<PaymentReconciliation>>('/finance/payment-reconciliation', { params }),
  getProfitReport: (params?: FinanceFilterParams) =>
    request.get<ApiResponse<ProfitReport>>('/finance/profit', { params }),
  getOperatingCosts: () =>
    request.get<ApiResponse<OperatingCost[]>>('/finance/operating-costs'),
  createOperatingCost: (data: OperatingCostForm) =>
    request.post<ApiResponse<OperatingCost>>('/finance/operating-costs', data),
  updateOperatingCost: (id: number, data: OperatingCostForm) =>
    request.put<ApiResponse<OperatingCost>>(`/finance/operating-costs/${id}`, data),
  deleteOperatingCost: (id: number) =>
    request.delete<ApiResponse<void>>(`/finance/operating-costs/${id}`),
  getReconciliationLogs: (params?: FinanceFilterParams & { page?: number; page_size?: number }) =>
    request.get<ApiResponse<ReconciliationLog[]>>('/finance/reconciliation-logs', { params }),
  createReconciliationLog: (data: Partial<ReconciliationLog>) =>
    request.post<ApiResponse<ReconciliationLog>>('/finance/reconciliation-logs', data),
  updateReconciliationLog: (id: number, data: { status: string; remark: string }) =>
    request.put<ApiResponse<ReconciliationLog>>(`/finance/reconciliation-logs/${id}`, data)
}
