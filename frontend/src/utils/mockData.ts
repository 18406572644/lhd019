import type {
  RevenueReport,
  CostAnalysisReport,
  CategorySalesReport,
  PaymentReconciliation,
  ProfitReport,
  OperatingCost
} from '@/api'

export const generateMockRevenueReport = (period: string = 'month'): RevenueReport => {
  return {
    period,
    start_date: '2024-01-01',
    end_date: '2024-01-31',
    total_revenue: 286500,
    total_orders: 2865,
    total_customers: 1980,
    average_order: 100,
    average_customer: 144.7,
    yoy_growth: 15.8,
    qoq_growth: 8.2,
    yoy_previous: 247409,
    qoq_previous: 264787,
    daily_data: Array.from({ length: 30 }, (_, i) => ({
      date: `01-${String(i + 1).padStart(2, '0')}`,
      revenue: 8000 + Math.random() * 12000,
      orders: 80 + Math.random() * 60,
      customers: 60 + Math.random() * 40
    }))
  }
}

export const generateMockCostAnalysisReport = (period: string = 'month'): CostAnalysisReport => {
  return {
    period,
    start_date: '2024-01-01',
    end_date: '2024-01-31',
    total_revenue: 286500,
    ingredient_cost: 28650,
    spirit_cost: 57300,
    total_material_cost: 85950,
    waste_cost: 4500,
    purchase_cost: 120000,
    operating_cost: 64500,
    total_cost: 274950,
    gross_profit: 200550,
    gross_margin: 70,
    net_profit: 11550,
    net_margin: 4.03,
    cost_breakdown: [
      { name: '原料成本', value: 28650, ratio: 10.42 },
      { name: '基酒成本', value: 57300, ratio: 20.84 },
      { name: '浪费成本', value: 4500, ratio: 1.64 },
      { name: '采购成本', value: 120000, ratio: 43.64 },
      { name: '运营成本', value: 64500, ratio: 23.46 }
    ]
  }
}

export const generateMockCategorySalesReport = (period: string = 'month'): CategorySalesReport => {
  return {
    period,
    start_date: '2024-01-01',
    end_date: '2024-01-31',
    category_sales: [
      { category: '经典鸡尾酒', quantity: 1256, revenue: 98400, percentage: 34.3 },
      { category: '特调鸡尾酒', quantity: 890, revenue: 115600, percentage: 40.3 },
      { category: '纯饮烈酒', quantity: 420, revenue: 52000, percentage: 18.1 },
      { category: '无酒精饮品', quantity: 299, revenue: 20500, percentage: 7.1 }
    ],
    recipe_sales: [
      { recipe_id: 9, recipe_name: '特调-东方茉莉', category: '特调鸡尾酒', quantity: 186, revenue: 23808, cost: 7068, profit: 16740, profit_margin: 70.3 },
      { recipe_id: 10, recipe_name: '特调-烟熏古典', category: '特调鸡尾酒', quantity: 165, revenue: 26070, cost: 10725, profit: 15345, profit_margin: 58.9 },
      { recipe_id: 6, recipe_name: '马天尼', category: '经典鸡尾酒', quantity: 142, revenue: 15336, cost: 6816, profit: 8520, profit_margin: 55.6 },
      { recipe_id: 3, recipe_name: '长岛冰茶', category: '经典鸡尾酒', quantity: 128, revenue: 12544, cost: 4480, profit: 8064, profit_margin: 64.3 },
      { recipe_id: 5, recipe_name: '曼哈顿', category: '经典鸡尾酒', quantity: 105, revenue: 10290, cost: 4725, profit: 5565, profit_margin: 54.1 }
    ],
    time_slot_sales: [
      { time_slot: '18-20', quantity: 320, revenue: 32000, orders: 85 },
      { time_slot: '20-22', quantity: 890, revenue: 95600, orders: 256 },
      { time_slot: '22-24', quantity: 1050, revenue: 128900, orders: 312 },
      { time_slot: '00-02', quantity: 450, revenue: 30000, orders: 125 }
    ]
  }
}

export const generateMockPaymentReconciliation = (period: string = 'month'): PaymentReconciliation => {
  return {
    period,
    start_date: '2024-01-01',
    end_date: '2024-01-31',
    total_revenue: 286500,
    payment_methods: [
      { payment_method: '微信支付', order_count: 1560, total_amount: 158900, percentage: 55.5 },
      { payment_method: '支付宝', order_count: 890, total_amount: 88500, percentage: 30.9 },
      { payment_method: '现金', order_count: 220, total_amount: 22800, percentage: 8.0 },
      { payment_method: '银行卡', order_count: 180, total_amount: 15300, percentage: 5.3 },
      { payment_method: '会员储值', order_count: 15, total_amount: 1000, percentage: 0.3 }
    ],
    reconciliation_logs: [
      { id: 1, order_no: 'ORD202401010001', payment_method: '微信支付', system_amount: 236.00, actual_amount: 236.00, difference: 0.00, status: 'matched', reconciled_at: '2024-01-01 23:59:00', remark: '对账一致', created_at: '2024-01-01' },
      { id: 2, order_no: 'ORD202401010002', payment_method: '支付宝', system_amount: 438.00, actual_amount: 438.00, difference: 0.00, status: 'matched', reconciled_at: '2024-01-01 23:59:00', remark: '对账一致', created_at: '2024-01-01' }
    ]
  }
}

export const generateMockProfitReport = (period: string = 'month'): ProfitReport => {
  return {
    period,
    start_date: '2024-01-01',
    end_date: '2024-01-31',
    total_revenue: 286500,
    material_cost: 85950,
    waste_cost: 4500,
    operating_cost: 64500,
    total_expenses: 154950,
    gross_profit: 200550,
    gross_margin: 70,
    net_profit: 131550,
    net_margin: 45.9,
    profit_breakdown: [
      { name: '营业收入', value: 286500, type: 'revenue' },
      { name: '物料成本', value: -85950, type: 'expense' },
      { name: '浪费成本', value: -4500, type: 'expense' },
      { name: '运营成本', value: -64500, type: 'expense' }
    ]
  }
}

export const generateMockOperatingCosts = (): OperatingCost[] => {
  return [
    { id: 1, cost_type: '房租', cost_name: '店面租金', amount: 15000, period: 'monthly', is_fixed: true, description: '每月固定房租支出', created_at: '2024-01-01', updated_at: '2024-01-01' },
    { id: 2, cost_type: '人工', cost_name: '调酒师工资', amount: 25000, period: 'monthly', is_fixed: true, description: '3名调酒师基本工资', created_at: '2024-01-01', updated_at: '2024-01-01' },
    { id: 3, cost_type: '人工', cost_name: '服务员基本工资', amount: 12000, period: 'monthly', is_fixed: true, description: '2名服务员基本工资', created_at: '2024-01-01', updated_at: '2024-01-01' },
    { id: 4, cost_type: '水电', cost_name: '水电费', amount: 3000, period: 'monthly', is_fixed: false, description: '每月水电费用预估', created_at: '2024-01-01', updated_at: '2024-01-01' },
    { id: 5, cost_type: '设备折旧', cost_name: '设备折旧摊销', amount: 2000, period: 'monthly', is_fixed: true, description: '制冷设备、音响等折旧', created_at: '2024-01-01', updated_at: '2024-01-01' },
    { id: 6, cost_type: '营销', cost_name: '营销推广费用', amount: 5000, period: 'monthly', is_fixed: false, description: '社交媒体、活动推广', created_at: '2024-01-01', updated_at: '2024-01-01' },
    { id: 7, cost_type: '耗材', cost_name: '一次性用品', amount: 1500, period: 'monthly', is_fixed: false, description: '纸巾、吸管、杯垫等', created_at: '2024-01-01', updated_at: '2024-01-01' },
    { id: 8, cost_type: '其他', cost_name: '其他杂费', amount: 1000, period: 'monthly', is_fixed: false, description: '清洁费、维修费等', created_at: '2024-01-01', updated_at: '2024-01-01' }
  ]
}
