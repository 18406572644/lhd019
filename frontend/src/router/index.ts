import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('@/views/Layout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: '营业概览', icon: 'DataAnalysis' }
      },
      {
        path: 'spirits',
        name: 'Spirits',
        component: () => import('@/views/Spirits.vue'),
        meta: { title: '基酒库存', icon: 'WineGlass' }
      },
      {
        path: 'ingredients',
        name: 'Ingredients',
        component: () => import('@/views/Ingredients.vue'),
        meta: { title: '配料管理', icon: 'Food' }
      },
      {
        path: 'recipes',
        name: 'Recipes',
        component: () => import('@/views/Recipes.vue'),
        meta: { title: '调酒配方', icon: 'Menu' }
      },
      {
        path: 'orders',
        name: 'Orders',
        component: () => import('@/views/Orders.vue'),
        meta: { title: '客单记录', icon: 'Tickets' }
      },
      {
        path: 'waste',
        name: 'Waste',
        component: () => import('@/views/Waste.vue'),
        meta: { title: '原料损耗', icon: 'Delete' }
      },
      {
        path: 'specials',
        name: 'Specials',
        component: () => import('@/views/Specials.vue'),
        meta: { title: '特调新品', icon: 'Star' }
      },
      {
        path: 'purchases',
        name: 'Purchases',
        component: () => import('@/views/Purchases.vue'),
        meta: { title: '采购台账', icon: 'ShoppingCart' }
      },
      {
        path: 'inventory-batches',
        name: 'InventoryBatches',
        component: () => import('@/views/InventoryBatches.vue'),
        meta: { title: '库存批次', icon: 'Box' }
      },
      {
        path: 'batch-trace',
        name: 'BatchTrace',
        component: () => import('@/views/BatchTrace.vue'),
        meta: { title: '批次追溯', icon: 'Search' }
      },
      {
        path: 'stocktake',
        name: 'Stocktake',
        component: () => import('@/views/StockTake.vue'),
        meta: { title: '库存盘点', icon: 'Document' }
      },
      {
        path: 'finance',
        name: 'Finance',
        redirect: '/finance/revenue',
        meta: { title: '财务报表', icon: 'Wallet' },
        children: [
          {
            path: 'revenue',
            name: 'FinanceRevenue',
            component: () => import('@/views/finance/RevenueReport.vue'),
            meta: { title: '营收报表', icon: 'TrendCharts' }
          },
          {
            path: 'cost',
            name: 'FinanceCost',
            component: () => import('@/views/finance/CostAnalysis.vue'),
            meta: { title: '成本分析', icon: 'PieChart' }
          },
          {
            path: 'sales',
            name: 'FinanceSales',
            component: () => import('@/views/finance/CategorySales.vue'),
            meta: { title: '分类销售', icon: 'Histogram' }
          },
          {
            path: 'payment',
            name: 'FinancePayment',
            component: () => import('@/views/finance/PaymentReconciliation.vue'),
            meta: { title: '支付对账', icon: 'CreditCard' }
          },
          {
            path: 'profit',
            name: 'FinanceProfit',
            component: () => import('@/views/finance/ProfitReport.vue'),
            meta: { title: '利润核算', icon: 'Money' }
          }
        ]
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
