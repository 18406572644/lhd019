<template>
  <el-container class="layout-container">
    <el-aside width="240px" class="sidebar">
      <div class="logo">
        <el-icon :size="32" color="#d4af37"><Goblet /></el-icon>
        <span class="logo-text">调酒台账</span>
      </div>
      <el-menu
        :default-active="activeMenu"
        class="menu"
        background-color="transparent"
        text-color="#a0a0a0"
        active-text-color="#d4af37"
        router
      >
        <el-menu-item v-for="item in menuItems" :key="item.path" :index="item.path">
          <el-icon><component :is="item.icon" /></el-icon>
          <span>{{ item.title }}</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="header">
        <div class="header-left">
          <h2 class="page-title">{{ currentPageTitle }}</h2>
        </div>
        <div class="header-right">
          <el-dropdown>
            <span class="user-info">
              <el-avatar :size="36" style="background: linear-gradient(135deg, #d4af37, #c9a227)">
                <el-icon><User /></el-icon>
              </el-avatar>
              <span class="username">调酒师</span>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item>个人设置</el-dropdown-item>
                <el-dropdown-item divided>退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <el-main class="main-content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  DataAnalysis,
  Goblet,
  Food,
  Menu,
  Tickets,
  Delete,
  Star,
  ShoppingCart,
  User
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()

const menuItems = [
  { path: '/dashboard', title: '营业概览', icon: DataAnalysis },
  { path: '/spirits', title: '基酒库存', icon: Goblet },
  { path: '/ingredients', title: '配料管理', icon: Food },
  { path: '/recipes', title: '调酒配方', icon: Menu },
  { path: '/orders', title: '客单记录', icon: Tickets },
  { path: '/waste', title: '原料损耗', icon: Delete },
  { path: '/specials', title: '特调新品', icon: Star },
  { path: '/purchases', title: '采购台账', icon: ShoppingCart }
]

const activeMenu = computed(() => route.path)

const currentPageTitle = computed(() => {
  const item = menuItems.find(m => m.path === route.path)
  return item ? item.title : ''
})
</script>

<style lang="scss" scoped>
.layout-container {
  height: 100vh;
}

.sidebar {
  background: rgba(0, 0, 0, 0.3);
  border-right: 1px solid rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);

  .logo {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 24px 20px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);

    .logo-text {
      font-size: 20px;
      font-weight: 700;
      background: linear-gradient(135deg, #d4af37, #c9a227);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
    }
  }

  .menu {
    border-right: none !important;
    padding: 12px 0;

    .el-menu-item {
      margin: 4px 12px;
      border-radius: 8px;
      transition: all 0.3s ease;

      &:hover {
        background: rgba(212, 175, 55, 0.1) !important;
      }

      &.is-active {
        background: rgba(212, 175, 55, 0.15) !important;
        border-left: 3px solid #d4af37;
      }
    }
  }
}

.header {
  background: rgba(0, 0, 0, 0.2);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 32px;

  .page-title {
    font-size: 20px;
    font-weight: 600;
    color: #d4af37;
    margin: 0;
  }

  .user-info {
    display: flex;
    align-items: center;
    gap: 12px;
    cursor: pointer;

    .username {
      color: #f5f5f5;
      font-weight: 500;
    }
  }
}

.main-content {
  padding: 24px;
  overflow-y: auto;
  background: transparent;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.fade-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
