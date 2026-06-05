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
    request.get<ApiResponse<any>>('/summary', { params })
}
