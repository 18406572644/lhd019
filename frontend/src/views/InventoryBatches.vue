<template>
  <div class="inventory-batches-container">
    <div class="page-header">
      <h1 class="gold-text">库存批次管理</h1>
      <p class="subtitle">库存批次与保质期管理系统</p>
    </div>

    <div class="summary-cards">
      <el-card class="glass-card summary-card">
        <div class="card-content">
          <div class="card-icon" style="background: linear-gradient(135deg, #d4af37, #c9a227);">
            <el-icon :size="28"><Box /></el-icon>
          </div>
          <div class="card-info">
            <p class="card-label">总批次数量</p>
            <p class="card-value gold-text">{{ totalBatches }}</p>
            <p class="card-trend trend-up">
              <el-icon><TrendCharts /></el-icon>
              共 {{ filteredBatches.length }} 条记录
            </p>
          </div>
        </div>
      </el-card>

      <el-card class="glass-card summary-card">
        <div class="card-content">
          <div class="card-icon" style="background: linear-gradient(135deg, #e74c3c, #c0392b);">
            <el-icon :size="28"><Warning /></el-icon>
          </div>
          <div class="card-info">
            <p class="card-label">临期预警</p>
            <p class="card-value" style="color: #e74c3c;">{{ expiryWarningCount }}</p>
            <p class="card-trend trend-down">
              <el-icon><Clock /></el-icon>
              30天内到期
            </p>
          </div>
        </div>
      </el-card>

      <el-card class="glass-card summary-card">
        <div class="card-content">
          <div class="card-icon" style="background: linear-gradient(135deg, #f39c12, #e67e22);">
            <el-icon :size="28"><Promotion /></el-icon>
          </div>
          <div class="card-info">
            <p class="card-label">临期促销</p>
            <p class="card-value" style="color: #f39c12;">{{ promotionCount }}</p>
            <p class="card-trend trend-up">
              <el-icon><ShoppingCart /></el-icon>
              已标记促销
            </p>
          </div>
        </div>
      </el-card>

      <el-card class="glass-card summary-card">
        <div class="card-content">
          <div class="card-icon" style="background: linear-gradient(135deg, #95a5a6, #7f8c8d);">
            <el-icon :size="28"><Close /></el-icon>
          </div>
          <div class="card-info">
            <p class="card-label">已过期</p>
            <p class="card-value" style="color: #95a5a6;">{{ expiredCount }}</p>
            <p class="card-trend trend-down">
              <el-icon><Delete /></el-icon>
              不可使用
            </p>
          </div>
        </div>
      </el-card>
    </div>

    <el-card class="glass-card">
      <div class="filter-section">
        <div class="filter-left">
          <el-select v-model="filterType" placeholder="原料类型" clearable class="filter-item" @change="fetchBatches">
            <el-option label="烈酒" value="spirit" />
            <el-option label="辅料" value="ingredient" />
          </el-select>
          <el-select v-model="filterStatus" placeholder="批次状态" clearable class="filter-item" @change="fetchBatches">
            <el-option label="正常" value="normal" />
            <el-option label="临期促销" value="promotion" />
            <el-option label="已耗尽" value="depleted" />
            <el-option label="已过期" value="expired" />
          </el-select>
          <el-input
            v-model="keyword"
            placeholder="搜索批次编码/批号/原料名称"
            clearable
            class="filter-item"
            @keyup.enter="fetchBatches"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-button class="glow-button" @click="fetchBatches">
            <el-icon><Search /></el-icon>
            查询
          </el-button>
          <el-button @click="resetFilters">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
        </div>
        <div class="filter-right">
          <el-button @click="fetchExpiryWarnings">
            <el-icon><Warning /></el-icon>
            刷新预警
          </el-button>
          <el-button class="glow-button" @click="showExpiryDialog = true">
            <el-icon><Bell /></el-icon>
            查看预警
          </el-button>
        </div>
      </div>

      <el-table :data="paginatedBatches" style="width: 100%" v-loading="loading" row-key="id">
        <el-table-column prop="batch_code" label="批次编码" min-width="180">
          <template #default="{ row }">
            <span class="batch-code">{{ row.batch_code }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="batch_no" label="供应商批号" min-width="140">
          <template #default="{ row }">
            <span v-if="row.batch_no" class="batch-no">{{ row.batch_no }}</span>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="ingredient_type" label="类型" width="90" align="center">
          <template #default="{ row }">
            <el-tag :type="row.ingredient_type === 'spirit' ? 'warning' : 'success'" effect="dark" size="small">
              {{ row.ingredient_type === 'spirit' ? '烈酒' : '辅料' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="ingredient_name" label="原料名称" min-width="150" />
        <el-table-column label="库存数量" width="140" align="right">
          <template #default="{ row }">
            <span class="qty-normal">{{ formatNum(row.remaining_quantity) }} {{ row.unit }}</span>
            <div class="qty-total">总: {{ formatNum(row.total_quantity) }} {{ row.unit }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="expiry_date" label="保质期" width="120" align="center">
          <template #default="{ row }">
            <span v-if="row.expiry_date" :class="getExpiryClass(row)">
              {{ row.expiry_date }}
            </span>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>
        <el-table-column label="到期天数" width="100" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.expiry_date" :type="getDaysTagType(row)" size="small">
              {{ getDaysToExpiry(row) }}天
            </el-tag>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="is_promotion" label="临期促销" width="100" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.is_promotion"
              :disabled="row.status === 'depleted' || row.status === 'expired'"
              @change="handlePromotionChange(row)"
              active-text="是"
              inactive-text="否"
            />
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" effect="dark" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="入库时间" width="160">
          <template #default="{ row }">
            <span class="created-at">{{ formatDateTime(row.created_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="140" align="center" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="openTraceDialog(row.batch_code)">
              <el-icon><Search /></el-icon>
              追溯
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-section">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="filteredBatches.length"
          layout="total, sizes, prev, pager, next, jumper"
          background
        />
      </div>
    </el-card>

    <el-dialog v-model="showExpiryDialog" title="保质期预警" width="900px">
      <el-table :data="expiryWarnings" style="width: 100%" size="default">
        <el-table-column prop="warning_level" label="预警级别" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getWarningLevelType(row.warning_level)" effect="dark" size="small">
              {{ getWarningLevelText(row.warning_level) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="days_to_expiry" label="剩余天数" width="100" align="center">
          <template #default="{ row }">
            <span :class="row.days_to_expiry < 0 ? 'text-danger' : 'text-warning'">
              {{ row.days_to_expiry < 0 ? '已过期' : row.days_to_expiry + '天' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="stock_batch.batch_code" label="批次编码" min-width="160" />
        <el-table-column prop="stock_batch.ingredient_name" label="原料名称" min-width="140" />
        <el-table-column label="剩余库存" width="120" align="right">
          <template #default="{ row }">
            {{ formatNum(row.stock_batch.remaining_quantity) }} {{ row.stock_batch.unit }}
          </template>
        </el-table-column>
        <el-table-column prop="stock_batch.expiry_date" label="到期日期" width="120" align="center" />
        <el-table-column label="操作" width="140" align="center">
          <template #default="{ row }">
            <el-button
              v-if="row.stock_batch.status !== 'expired'"
              type="primary"
              link
              size="small"
              @click="markPromotion(row.stock_batch)"
            >
              <el-icon><Promotion /></el-icon>
              标记促销
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Box,
  Warning,
  Promotion,
  Close,
  TrendCharts,
  Clock,
  ShoppingCart,
  Delete,
  Search,
  Refresh,
  Bell
} from '@element-plus/icons-vue'
import { api, type StockBatch, type ExpiryWarningResult } from '@/api'

const loading = ref(false)
const batches = ref<StockBatch[]>([])
const expiryWarnings = ref<ExpiryWarningResult[]>([])
const filterType = ref('')
const filterStatus = ref('')
const keyword = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const showExpiryDialog = ref(false)

const totalBatches = computed(() => batches.value.length)
const expiryWarningCount = computed(() => expiryWarnings.value.filter(w => w.warning_level !== 'expired').length)
const promotionCount = computed(() => batches.value.filter(b => b.is_promotion).length)
const expiredCount = computed(() => batches.value.filter(b => b.status === 'expired').length)

const filteredBatches = computed(() => {
  let result = [...batches.value]
  if (filterType.value) {
    result = result.filter(b => b.ingredient_type === filterType.value)
  }
  if (filterStatus.value) {
    if (filterStatus.value === 'promotion') {
      result = result.filter(b => b.is_promotion)
    } else {
      result = result.filter(b => b.status === filterStatus.value)
    }
  }
  if (keyword.value) {
    const kw = keyword.value.toLowerCase()
    result = result.filter(b =>
      b.batch_code.toLowerCase().includes(kw) ||
      (b.batch_no && b.batch_no.toLowerCase().includes(kw)) ||
      b.ingredient_name.toLowerCase().includes(kw)
    )
  }
  return result.sort((a, b) => {
    if (a.expiry_date && b.expiry_date) {
      return new Date(a.expiry_date).getTime() - new Date(b.expiry_date).getTime()
    }
    return new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
  })
})

const paginatedBatches = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredBatches.value.slice(start, end)
})

const formatNum = (num: number): string => {
  if (num === undefined || num === null || isNaN(num)) return '0.00'
  return Number(num).toFixed(2)
}

const formatDateTime = (dateStr: string): string => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  if (isNaN(date.getTime())) return '-'
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getDaysToExpiry = (row: StockBatch): number => {
  if (!row.expiry_date) return 0
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  const expiry = new Date(row.expiry_date)
  expiry.setHours(0, 0, 0, 0)
  return Math.ceil((expiry.getTime() - today.getTime()) / (1000 * 60 * 60 * 24))
}

const getExpiryClass = (row: StockBatch): string => {
  const days = getDaysToExpiry(row)
  if (days < 0) return 'text-danger'
  if (days <= 7) return 'text-danger'
  if (days <= 15) return 'text-warning'
  if (days <= 30) return 'text-orange'
  return ''
}

const getDaysTagType = (row: StockBatch): string => {
  const days = getDaysToExpiry(row)
  if (days < 0) return 'danger'
  if (days <= 7) return 'danger'
  if (days <= 15) return 'warning'
  if (days <= 30) return 'warning'
  return 'success'
}

const getStatusType = (status: string): string => {
  switch (status) {
    case 'normal': return 'success'
    case 'depleted': return 'info'
    case 'expired': return 'danger'
    default: return 'info'
  }
}

const getStatusText = (status: string): string => {
  switch (status) {
    case 'normal': return '正常'
    case 'depleted': return '已耗尽'
    case 'expired': return '已过期'
    default: return status
  }
}

const getWarningLevelType = (level: string): string => {
  switch (level) {
    case 'expired': return 'danger'
    case 'urgent': return 'danger'
    case 'warning': return 'warning'
    case 'attention': return 'warning'
    default: return 'info'
  }
}

const getWarningLevelText = (level: string): string => {
  switch (level) {
    case 'expired': return '已过期'
    case 'urgent': return '紧急'
    case 'warning': return '警告'
    case 'attention': return '注意'
    default: return '正常'
  }
}

const fetchBatches = async () => {
  loading.value = true
  try {
    const params: Record<string, any> = {}
    if (filterType.value) params.ingredient_type = filterType.value
    if (filterStatus.value && filterStatus.value !== 'promotion') params.status = filterStatus.value
    if (keyword.value) params.keyword = keyword.value

    const res = await api.getStockBatches(params)
    if (res.data.code === 0) {
      batches.value = res.data.data || []
    } else {
      ElMessage.error(res.data.message || '获取批次列表失败')
    }
  } catch (error) {
    console.error('获取批次列表失败:', error)
    ElMessage.error('获取批次列表失败')
  } finally {
    loading.value = false
  }
}

const fetchExpiryWarnings = async () => {
  try {
    const res = await api.getExpiryWarnings({ days: 30 })
    if (res.data.code === 0) {
      expiryWarnings.value = res.data.data || []
    }
  } catch (error) {
    console.error('获取预警信息失败:', error)
  }
}

const resetFilters = () => {
  filterType.value = ''
  filterStatus.value = ''
  keyword.value = ''
  currentPage.value = 1
  fetchBatches()
}

const handlePromotionChange = async (row: StockBatch) => {
  try {
    const res = await api.updateStockBatchPromotion(row.id, {
      is_promotion: row.is_promotion,
      remark: row.is_promotion ? '手动标记临期促销' : '取消临期促销标记'
    })
    if (res.data.code === 0) {
      ElMessage.success(row.is_promotion ? '已标记为临期促销' : '已取消临期促销标记')
    } else {
      row.is_promotion = !row.is_promotion
      ElMessage.error(res.data.message || '操作失败')
    }
  } catch (error) {
    row.is_promotion = !row.is_promotion
    console.error('操作失败:', error)
    ElMessage.error('操作失败')
  }
}

const markPromotion = async (batch: StockBatch) => {
  try {
    const res = await api.updateStockBatchPromotion(batch.id, {
      is_promotion: true,
      remark: '临期预警标记促销'
    })
    if (res.data.code === 0) {
      ElMessage.success('已标记为临期促销')
      fetchBatches()
      fetchExpiryWarnings()
    } else {
      ElMessage.error(res.data.message || '操作失败')
    }
  } catch (error) {
    console.error('操作失败:', error)
    ElMessage.error('操作失败')
  }
}

const openTraceDialog = (batchCode: string) => {
  window.location.href = `/#/batch-trace?batch_code=${batchCode}`
}

onMounted(() => {
  fetchBatches()
  fetchExpiryWarnings()
})
</script>

<style lang="scss" scoped>
@import '@/styles/variables.scss';

.inventory-batches-container {
  padding: 24px;
  min-height: calc(100vh - 60px);
  background: $dark-bg;
}

.page-header {
  margin-bottom: 24px;
  h1 {
    font-size: 28px;
    font-weight: 700;
    margin-bottom: 8px;
  }
  .subtitle {
    color: $text-secondary;
    font-size: 14px;
  }
}

.summary-cards {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-bottom: 24px;
  @media (max-width: 1200px) {
    grid-template-columns: repeat(2, 1fr);
  }
  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }
}

.summary-card {
  transition: all $transition-normal;
  &:hover {
    transform: translateY(-4px);
    box-shadow: $shadow-gold, $shadow-lg;
  }
  .card-content {
    display: flex;
    align-items: center;
    gap: 16px;
  }
  .card-icon {
    width: 56px;
    height: 56px;
    border-radius: $border-radius-md;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #fff;
    flex-shrink: 0;
  }
  .card-info {
    flex: 1;
    min-width: 0;
  }
  .card-label {
    color: $text-secondary;
    font-size: 13px;
    margin-bottom: 4px;
  }
  .card-value {
    font-size: 24px;
    font-weight: 700;
    margin-bottom: 4px;
  }
  .card-trend {
    display: flex;
    align-items: center;
    gap: 4px;
    font-size: 12px;
    &.trend-up { color: $success; }
    &.trend-down { color: $danger; }
  }
}

.filter-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  flex-wrap: wrap;
  gap: 16px;
  .filter-left {
    display: flex;
    align-items: center;
    gap: 12px;
    flex-wrap: wrap;
  }
  .filter-item {
    width: 200px;
  }
}

.batch-code {
  font-family: 'Consolas', monospace;
  font-weight: 600;
  color: $primary-gold;
  font-size: 13px;
}

.batch-no {
  font-family: 'Consolas', monospace;
  color: $text-secondary;
  font-size: 13px;
}

.qty-normal {
  font-weight: 600;
  color: $text-primary;
}

.qty-total {
  font-size: 11px;
  color: $text-muted;
  margin-top: 2px;
}

.text-danger {
  color: #e74c3c;
  font-weight: 500;
}

.text-warning {
  color: #f39c12;
  font-weight: 500;
}

.text-orange {
  color: #e67e22;
  font-weight: 500;
}

.created-at {
  color: $text-secondary;
  font-size: 12px;
}

.text-muted {
  color: $text-muted;
}

.pagination-section {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
