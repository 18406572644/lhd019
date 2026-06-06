<template>
  <div class="batch-trace-container">
    <div class="page-header">
      <h1 class="gold-text">批次追溯查询</h1>
      <p class="subtitle">查询批次原料的出库使用记录和关联订单</p>
    </div>

    <el-card class="glass-card search-card">
      <div class="search-section">
        <el-input
          v-model="searchBatchCode"
          placeholder="请输入批次编码进行追溯查询"
          clearable
          class="search-input"
          @keyup.enter="doSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-button class="glow-button" @click="doSearch">
          <el-icon><Search /></el-icon>
          查询
        </el-button>
        <el-button @click="goBack">
          <el-icon><Back /></el-icon>
          返回
        </el-button>
      </div>
    </el-card>

    <div v-if="traceResult" v-loading="loading">
      <el-card class="glass-card batch-info-card">
        <div class="info-header">
          <div class="info-title">
            <el-icon :size="24" class="info-icon"><Box /></el-icon>
            <span>批次基本信息</span>
          </div>
          <el-tag :type="getStatusType(traceResult.stock_batch.status)" effect="dark" size="default">
            {{ getStatusText(traceResult.stock_batch.status) }}
          </el-tag>
        </div>
        <div class="info-grid">
          <div class="info-item">
            <span class="info-label">批次编码</span>
            <span class="info-value batch-code">{{ traceResult.stock_batch.batch_code }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">供应商批号</span>
            <span class="info-value">{{ traceResult.stock_batch.batch_no || '-' }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">原料类型</span>
            <span class="info-value">
              <el-tag :type="traceResult.stock_batch.ingredient_type === 'spirit' ? 'warning' : 'success'" effect="dark" size="small">
                {{ traceResult.stock_batch.ingredient_type === 'spirit' ? '烈酒' : '辅料' }}
              </el-tag>
            </span>
          </div>
          <div class="info-item">
            <span class="info-label">原料名称</span>
            <span class="info-value">{{ traceResult.stock_batch.ingredient_name }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">总入库数量</span>
            <span class="info-value">{{ formatNum(traceResult.stock_batch.total_quantity) }} {{ traceResult.stock_batch.unit }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">剩余库存</span>
            <span class="info-value qty-remaining">{{ formatNum(traceResult.stock_batch.remaining_quantity) }} {{ traceResult.stock_batch.unit }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">已出库数量</span>
            <span class="info-value qty-out">{{ formatNum(traceResult.total_out_qty) }} {{ traceResult.stock_batch.unit }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">库存单价</span>
            <span class="info-value">¥{{ formatNum(traceResult.stock_batch.unit_price) }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">保质期</span>
            <span class="info-value" :class="getExpiryClass(traceResult.stock_batch)">
              {{ traceResult.stock_batch.expiry_date || '-' }}
            </span>
          </div>
          <div class="info-item">
            <span class="info-label">到期天数</span>
            <span class="info-value">
              <el-tag v-if="traceResult.stock_batch.expiry_date" :type="getDaysTagType(traceResult.stock_batch)" size="small">
                {{ getDaysToExpiry(traceResult.stock_batch) }}天
              </el-tag>
              <span v-else>-</span>
            </span>
          </div>
          <div class="info-item">
            <span class="info-label">临期促销</span>
            <span class="info-value">
              <el-tag :type="traceResult.stock_batch.is_promotion ? 'warning' : 'info'" effect="dark" size="small">
                {{ traceResult.stock_batch.is_promotion ? '是' : '否' }}
              </el-tag>
            </span>
          </div>
          <div class="info-item">
            <span class="info-label">入库时间</span>
            <span class="info-value">{{ formatDateTime(traceResult.stock_batch.created_at) }}</span>
          </div>
        </div>
      </el-card>

      <el-card class="glass-card">
        <div class="table-header">
          <div class="table-title">
            <el-icon :size="20" class="info-icon"><Tickets /></el-icon>
            <span>出库记录 (共 {{ traceResult.out_records.length }} 条)</span>
          </div>
          <div class="table-stats">
            <span class="stat-item">
              <el-icon><TrendCharts /></el-icon>
              累计出库: <strong>{{ formatNum(traceResult.total_out_qty) }} {{ traceResult.stock_batch.unit }}</strong>
            </span>
          </div>
        </div>

        <el-table :data="paginatedRecords" style="width: 100%" v-loading="loading" row-key="id">
          <el-table-column prop="id" label="记录ID" width="90" align="center" />
          <el-table-column prop="out_type" label="出库类型" width="100" align="center">
            <template #default="{ row }">
              <el-tag :type="getOutTypeTag(row.out_type)" effect="dark" size="small">
                {{ getOutTypeText(row.out_type) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="关联订单" min-width="180">
            <template #default="{ row }">
              <div v-if="row.order_id">
                <span class="order-no" @click="goToOrder(row.order_id)">
                  <el-icon><Link /></el-icon>
                  {{ row.order_no || ('订单#' + row.order_id) }}
                </span>
              </div>
              <span v-else class="text-muted">-</span>
            </template>
          </el-table-column>
          <el-table-column label="出库数量" width="130" align="right">
            <template #default="{ row }">
              <span class="qty-out">{{ formatNum(row.out_quantity) }} {{ row.unit }}</span>
            </template>
          </el-table-column>
          <el-table-column label="出库成本" width="130" align="right">
            <template #default="{ row }">
              <span>¥{{ formatNum(row.total_cost) }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="operator" label="操作人" width="100" align="center">
            <template #default="{ row }">
              <span>{{ row.operator || '系统' }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="出库时间" width="160">
            <template #default="{ row }">
              <span class="created-at">{{ formatDateTime(row.created_at) }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="remark" label="备注" min-width="150">
            <template #default="{ row }">
              <span>{{ row.remark || '-' }}</span>
            </template>
          </el-table-column>
        </el-table>

        <div class="pagination-section" v-if="traceResult.out_records.length > 0">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[10, 20, 50]"
            :total="traceResult.out_records.length"
            layout="total, sizes, prev, pager, next, jumper"
            background
          />
        </div>

        <el-empty v-if="traceResult.out_records.length === 0" description="该批次暂无出库记录" />
      </el-card>
    </div>

    <el-card v-else-if="!loading" class="glass-card empty-card">
      <el-empty description="请输入批次编码进行追溯查询">
        <template #image>
          <el-icon :size="80" style="color: #d4af37;"><Search /></el-icon>
        </template>
      </el-empty>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Search,
  Back,
  Box,
  TrendCharts,
  Tickets,
  Link
} from '@element-plus/icons-vue'
import { api, type BatchTraceResult, type StockBatch, type BatchOutRecord } from '@/api'

const route = useRoute()
const loading = ref(false)
const searchBatchCode = ref('')
const traceResult = ref<BatchTraceResult | null>(null)
const currentPage = ref(1)
const pageSize = ref(10)

const paginatedRecords = computed(() => {
  if (!traceResult.value) return []
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return traceResult.value.out_records.slice(start, end)
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

const getOutTypeTag = (type: string): string => {
  switch (type) {
    case 'order': return 'primary'
    case 'waste': return 'danger'
    case 'manual': return 'warning'
    default: return 'info'
  }
}

const getOutTypeText = (type: string): string => {
  switch (type) {
    case 'order': return '订单出库'
    case 'waste': return '损耗出库'
    case 'manual': return '手动出库'
    default: return type
  }
}

const doSearch = async () => {
  if (!searchBatchCode.value.trim()) {
    ElMessage.warning('请输入批次编码')
    return
  }
  loading.value = true
  try {
    const res = await api.traceBatch(searchBatchCode.value.trim())
    if (res.data.code === 0) {
      traceResult.value = res.data.data
      currentPage.value = 1
    } else {
      ElMessage.error(res.data.message || '查询失败')
      traceResult.value = null
    }
  } catch (error) {
    console.error('查询失败:', error)
    ElMessage.error('查询失败，请检查批次编码是否正确')
    traceResult.value = null
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  window.location.href = '/#/inventory-batches'
}

const goToOrder = (orderId: number) => {
  window.location.href = `/#/orders?id=${orderId}`
}

onMounted(() => {
  const batchCode = route.query.batch_code as string
  if (batchCode) {
    searchBatchCode.value = batchCode
    doSearch()
  }
})
</script>

<style lang="scss" scoped>
@import '@/styles/variables.scss';

.batch-trace-container {
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

.search-card {
  margin-bottom: 24px;
}

.search-section {
  display: flex;
  align-items: center;
  gap: 12px;
  .search-input {
    flex: 1;
    max-width: 500px;
  }
}

.batch-info-card {
  margin-bottom: 24px;
}

.info-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid rgba(212, 175, 55, 0.2);
}

.info-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 18px;
  font-weight: 600;
  color: $text-primary;
}

.info-icon {
  color: $primary-gold;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  @media (max-width: 1200px) {
    grid-template-columns: repeat(2, 1fr);
  }
  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 12px 16px;
  background: rgba(212, 175, 55, 0.05);
  border-radius: $border-radius-md;
  border: 1px solid rgba(212, 175, 55, 0.1);
}

.info-label {
  font-size: 12px;
  color: $text-secondary;
}

.info-value {
  font-size: 14px;
  color: $text-primary;
  font-weight: 500;
}

.batch-code {
  font-family: 'Consolas', monospace;
  color: $primary-gold;
  font-weight: 600;
}

.qty-remaining {
  color: $success;
  font-weight: 600;
}

.qty-out {
  color: $danger;
  font-weight: 600;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid rgba(212, 175, 55, 0.2);
}

.table-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  color: $text-primary;
}

.table-stats {
  .stat-item {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 13px;
    color: $text-secondary;
    strong {
      color: $primary-gold;
      font-size: 15px;
    }
  }
}

.order-no {
  color: $primary-gold;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 4px;
  font-weight: 500;
  &:hover {
    text-decoration: underline;
  }
}

.created-at {
  color: $text-secondary;
  font-size: 12px;
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

.text-muted {
  color: $text-muted;
}

.pagination-section {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.empty-card {
  min-height: 400px;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
