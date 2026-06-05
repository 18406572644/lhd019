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
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
