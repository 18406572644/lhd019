<template>
  <div class="stocktake-container">
    <div class="page-header">
      <h1 class="gold-text">库存盘点</h1>
      <p class="subtitle">周期盘点、盘盈盘亏管理</p>
    </div>

    <div class="summary-cards">
      <el-card class="glass-card summary-card">
        <div class="card-content">
          <div class="card-icon" style="background: linear-gradient(135deg, #d4af37, #c9a227);">
            <el-icon :size="28"><Document /></el-icon>
          </div>
          <div class="card-info">
            <p class="card-label">盘点单总数</p>
            <p class="card-value gold-text">{{ stocktakes.length }}</p>
            <p class="card-trend trend-up">
              <el-icon><TrendCharts /></el-icon>
              历史记录
            </p>
          </div>
        </div>
      </el-card>

      <el-card class="glass-card summary-card">
        <div class="card-content">
          <div class="card-icon" style="background: linear-gradient(135deg, #3498db, #2980b9);">
            <el-icon :size="28"><Edit /></el-icon>
          </div>
          <div class="card-info">
            <p class="card-label">待确认</p>
            <p class="card-value" style="color: #3498db;">{{ draftCount }}</p>
            <p class="card-trend trend-up">
              <el-icon><Clock /></el-icon>
              草稿状态
            </p>
          </div>
        </div>
      </el-card>

      <el-card class="glass-card summary-card">
        <div class="card-content">
          <div class="card-icon" style="background: linear-gradient(135deg, #27ae60, #229954);">
            <el-icon :size="28"><CircleCheck /></el-icon>
          </div>
          <div class="card-info">
            <p class="card-label">已完成</p>
            <p class="card-value" style="color: #27ae60;">{{ confirmedCount }}</p>
            <p class="card-trend trend-up">
              <el-icon><Check /></el-icon>
              已确认
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
            <p class="card-label">本月盘亏</p>
            <p class="card-value" style="color: #e74c3c;">¥{{ formatNum(monthlyLoss) }}</p>
            <p class="card-trend trend-down">
              <el-icon><TrendCharts /></el-icon>
              总金额
            </p>
          </div>
        </div>
      </el-card>
    </div>

    <el-card class="glass-card" v-if="!showDetail && !showCreate">
      <div class="filter-section">
        <div class="filter-left">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
            class="filter-item"
            @change="fetchStocktakes"
          />
          <el-select v-model="filterStatus" placeholder="盘点状态" clearable class="filter-item" @change="fetchStocktakes">
            <el-option label="草稿" value="draft" />
            <el-option label="已确认" value="confirmed" />
          </el-select>
          <el-select v-model="filterType" placeholder="盘点类型" clearable class="filter-item" @change="fetchStocktakes">
            <el-option label="周期盘点" value="periodic" />
            <el-option label="月度盘点" value="monthly" />
            <el-option label="年度盘点" value="yearly" />
          </el-select>
          <el-button class="glow-button" @click="fetchStocktakes">
            <el-icon><Search /></el-icon>
            查询
          </el-button>
          <el-button @click="resetFilters">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
        </div>
        <div class="filter-right">
          <el-button class="glow-button" @click="openCreateDialog">
            <el-icon><Plus /></el-icon>
            新建盘点
          </el-button>
        </div>
      </div>

      <el-table :data="paginatedStocktakes" style="width: 100%" v-loading="loading" row-key="id">
        <el-table-column prop="stocktake_no" label="盘点单号" min-width="160">
          <template #default="{ row }">
            <span class="stocktake-no">{{ row.stocktake_no }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="stocktake_date" label="盘点日期" width="120" align="center" />
        <el-table-column prop="stocktake_type" label="盘点类型" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStocktakeTypeTag(row.stocktake_type)" effect="dark" size="small">
              {{ getStocktakeTypeText(row.stocktake_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 'draft' ? 'info' : 'success'" effect="dark" size="small">
              {{ row.status === 'draft' ? '草稿' : '已确认' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="盘盈" width="120" align="right">
          <template #default="{ row }">
            <span v-if="row.total_profit > 0" class="profit">+¥{{ formatNum(row.total_profit) }}</span>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>
        <el-table-column label="盘亏" width="120" align="right">
          <template #default="{ row }">
            <span v-if="row.total_loss > 0" class="loss">-¥{{ formatNum(row.total_loss) }}</span>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>
        <el-table-column label="差异" width="120" align="right">
          <template #default="{ row }">
            <span :class="row.total_diff >= 0 ? 'profit' : 'loss'">
              {{ row.total_diff >= 0 ? '+' : '' }}¥{{ formatNum(row.total_diff) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="operator" label="操作人" width="100" align="center">
          <template #default="{ row }">
            <span>{{ row.operator || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160">
          <template #default="{ row }">
            <span class="created-at">{{ formatDateTime(row.created_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" align="center" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="viewDetail(row)">
              <el-icon><View /></el-icon>
              详情
            </el-button>
            <el-button v-if="row.status === 'draft'" type="danger" link size="small" @click="deleteStocktake(row)">
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-section">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50]"
          :total="filteredStocktakes.length"
          layout="total, sizes, prev, pager, next, jumper"
          background
        />
      </div>
    </el-card>

    <el-card class="glass-card" v-if="showCreate">
      <div class="detail-header">
        <div class="detail-title">
          <el-button @click="showCreate = false" class="back-btn">
            <el-icon><Back /></el-icon>
            返回
          </el-button>
          <h2>新建盘点单</h2>
        </div>
      </div>

      <el-form :model="createForm" label-width="100px" class="create-form">
        <div class="form-row">
          <el-form-item label="盘点日期" required>
            <el-date-picker
              v-model="createForm.stocktake_date"
              type="date"
              placeholder="选择盘点日期"
              value-format="YYYY-MM-DD"
              style="width: 200px"
            />
          </el-form-item>
          <el-form-item label="盘点类型" required>
            <el-select v-model="createForm.stocktake_type" placeholder="请选择盘点类型" style="width: 200px">
              <el-option label="周期盘点" value="periodic" />
              <el-option label="月度盘点" value="monthly" />
              <el-option label="年度盘点" value="yearly" />
            </el-select>
          </el-form-item>
          <el-form-item label="操作人">
            <el-input v-model="createForm.operator" placeholder="请输入操作人" style="width: 200px" />
          </el-form-item>
        </div>
        <el-form-item label="备注">
          <el-input v-model="createForm.remark" type="textarea" :rows="2" placeholder="请输入备注" />
        </el-form-item>
        <el-form-item>
          <el-button class="glow-button" @click="generateItems">
            <el-icon><RefreshRight /></el-icon>
            生成盘点项
          </el-button>
        </el-form-item>
      </el-form>

      <div v-if="generatedItems.length > 0">
        <div class="table-header">
          <div class="table-title">
            <el-icon :size="20" class="info-icon"><Tickets /></el-icon>
            <span>盘点明细 (共 {{ generatedItems.length }} 项)</span>
          </div>
          <div class="table-actions">
            <el-button type="success" @click="submitStocktake">
              <el-icon><Check /></el-icon>
              提交盘点单
            </el-button>
          </div>
        </div>

        <el-table :data="generatedItems" style="width: 100%" border>
          <el-table-column prop="ingredient_type" label="类型" width="80" align="center">
            <template #default="{ row }">
              <el-tag :type="row.ingredient_type === 'spirit' ? 'warning' : 'success'" effect="dark" size="small">
                {{ row.ingredient_type === 'spirit' ? '烈酒' : '辅料' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="ingredient_name" label="原料名称" min-width="150" />
          <el-table-column prop="category" label="分类" width="120" />
          <el-table-column label="系统库存" width="120" align="right">
            <template #default="{ row }">
              <span class="system-qty">{{ formatNum(row.system_quantity) }} {{ row.unit }}</span>
            </template>
          </el-table-column>
          <el-table-column label="实际库存" width="140" align="center">
            <template #default="{ row }">
              <el-input-number
                v-model="row.actual_quantity"
                :min="0"
                :precision="2"
                :step="1"
                size="small"
                style="width: 120px"
                @change="calculateDiff(row)"
              />
            </template>
          </el-table-column>
          <el-table-column label="差异数量" width="120" align="right">
            <template #default="{ row }">
              <span :class="getDiffClass(row)">
                {{ getDiffQty(row) >= 0 ? '+' : '' }}{{ formatNum(getDiffQty(row)) }} {{ row.unit }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="差异金额" width="120" align="right">
            <template #default="{ row }">
              <span :class="getDiffClass(row)">
                {{ getDiffAmount(row) >= 0 ? '+' : '' }}¥{{ formatNum(getDiffAmount(row)) }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="差异类型" width="90" align="center">
            <template #default="{ row }">
              <el-tag v-if="getDiffType(row) === 'profit'" type="success" effect="dark" size="small">盘盈</el-tag>
              <el-tag v-else-if="getDiffType(row) === 'loss'" type="danger" effect="dark" size="small">盘亏</el-tag>
              <el-tag v-else type="info" effect="dark" size="small">正常</el-tag>
            </template>
          </el-table-column>
        </el-table>

        <div class="summary-section">
          <div class="summary-item profit">
            <span>盘盈总额:</span>
            <strong>+¥{{ formatNum(totalProfit) }}</strong>
          </div>
          <div class="summary-item loss">
            <span>盘亏总额:</span>
            <strong>-¥{{ formatNum(totalLoss) }}</strong>
          </div>
          <div class="summary-item">
            <span>净差异:</span>
            <strong :class="totalDiff >= 0 ? 'profit' : 'loss'">
              {{ totalDiff >= 0 ? '+' : '' }}¥{{ formatNum(totalDiff) }}
            </strong>
          </div>
        </div>
      </div>
    </el-card>

    <el-card class="glass-card" v-if="showDetail && currentStocktake">
      <div class="detail-header">
        <div class="detail-title">
          <el-button @click="showDetail = false" class="back-btn">
            <el-icon><Back /></el-icon>
            返回
          </el-button>
          <h2>盘点单详情</h2>
          <el-tag :type="currentStocktake.status === 'draft' ? 'info' : 'success'" effect="dark" size="default">
            {{ currentStocktake.status === 'draft' ? '草稿' : '已确认' }}
          </el-tag>
        </div>
        <div v-if="currentStocktake.status === 'draft'" class="detail-actions">
          <el-button type="success" @click="confirmStocktake">
            <el-icon><Check /></el-icon>
            确认盘点
          </el-button>
          <el-button type="danger" @click="deleteStocktake(currentStocktake)">
            <el-icon><Delete /></el-icon>
            删除
          </el-button>
        </div>
      </div>

      <div class="info-grid">
        <div class="info-item">
          <span class="info-label">盘点单号</span>
          <span class="info-value stocktake-no">{{ currentStocktake.stocktake_no }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">盘点日期</span>
          <span class="info-value">{{ currentStocktake.stocktake_date }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">盘点类型</span>
          <span class="info-value">
            <el-tag :type="getStocktakeTypeTag(currentStocktake.stocktake_type)" effect="dark" size="small">
              {{ getStocktakeTypeText(currentStocktake.stocktake_type) }}
            </el-tag>
          </span>
        </div>
        <div class="info-item">
          <span class="info-label">操作人</span>
          <span class="info-value">{{ currentStocktake.operator || '-' }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">盘盈总额</span>
          <span class="info-value profit">+¥{{ formatNum(currentStocktake.total_profit) }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">盘亏总额</span>
          <span class="info-value loss">-¥{{ formatNum(currentStocktake.total_loss) }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">净差异</span>
          <span class="info-value" :class="currentStocktake.total_diff >= 0 ? 'profit' : 'loss'">
            {{ currentStocktake.total_diff >= 0 ? '+' : '' }}¥{{ formatNum(currentStocktake.total_diff) }}
          </span>
        </div>
        <div class="info-item">
          <span class="info-label">确认时间</span>
          <span class="info-value">{{ currentStocktake.confirmed_at ? formatDateTime(currentStocktake.confirmed_at) : '-' }}</span>
        </div>
      </div>

      <div v-if="currentStocktake.remark" class="remark-section">
        <span class="info-label">备注:</span>
        <span>{{ currentStocktake.remark }}</span>
      </div>

      <div class="table-header" style="margin-top: 24px;">
        <div class="table-title">
          <el-icon :size="20" class="info-icon"><Tickets /></el-icon>
          <span>盘点明细 (共 {{ currentStocktake.stocktake_items?.length || 0 }} 项)</span>
        </div>
      </div>

      <el-table :data="currentStocktake.stocktake_items || []" style="width: 100%" border>
        <el-table-column prop="ingredient_type" label="类型" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.ingredient_type === 'spirit' ? 'warning' : 'success'" effect="dark" size="small">
              {{ row.ingredient_type === 'spirit' ? '烈酒' : '辅料' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="ingredient_name" label="原料名称" min-width="150" />
        <el-table-column label="系统库存" width="120" align="right">
          <template #default="{ row }">
            <span class="system-qty">{{ formatNum(row.system_quantity) }} {{ row.unit }}</span>
          </template>
        </el-table-column>
        <el-table-column label="实际库存" width="120" align="right">
          <template #default="{ row }">
            <span>{{ formatNum(row.actual_quantity) }} {{ row.unit }}</span>
          </template>
        </el-table-column>
        <el-table-column label="差异数量" width="120" align="right">
          <template #default="{ row }">
            <span :class="getDiffType(row.diff_type) === 'profit' ? 'profit' : getDiffType(row.diff_type) === 'loss' ? 'loss' : ''">
              {{ row.diff_quantity >= 0 ? '+' : '' }}{{ formatNum(row.diff_quantity) }} {{ row.unit }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="差异金额" width="120" align="right">
          <template #default="{ row }">
            <span :class="getDiffType(row.diff_type) === 'profit' ? 'profit' : getDiffType(row.diff_type) === 'loss' ? 'loss' : ''">
              {{ row.diff_amount >= 0 ? '+' : '' }}¥{{ formatNum(row.diff_amount) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="差异类型" width="90" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.diff_type === 'profit'" type="success" effect="dark" size="small">盘盈</el-tag>
            <el-tag v-else-if="row.diff_type === 'loss'" type="danger" effect="dark" size="small">盘亏</el-tag>
            <el-tag v-else type="info" effect="dark" size="small">正常</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="120">
          <template #default="{ row }">
            <span>{{ row.remark || '-' }}</span>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Document,
  Edit,
  CircleCheck,
  Warning,
  TrendCharts,
  Clock,
  Check,
  Search,
  Refresh,
  Plus,
  View,
  Delete,
  Back,
  RefreshRight,
  Tickets
} from '@element-plus/icons-vue'
import { api, type Stocktake, type StocktakeItemGenerate, type StocktakeCreateRequest } from '@/api'

const loading = ref(false)
const stocktakes = ref<Stocktake[]>([])
const dateRange = ref<string[]>([])
const filterStatus = ref('')
const filterType = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const showDetail = ref(false)
const showCreate = ref(false)
const currentStocktake = ref<Stocktake | null>(null)

const generatedItems = ref<StocktakeItemGenerate[]>([])
const createForm = reactive({
  stocktake_date: '',
  stocktake_type: '',
  operator: '',
  remark: ''
})

const draftCount = computed(() => stocktakes.value.filter(s => s.status === 'draft').length)
const confirmedCount = computed(() => stocktakes.value.filter(s => s.status === 'confirmed').length)
const monthlyLoss = computed(() => {
  const now = new Date()
  const monthStart = new Date(now.getFullYear(), now.getMonth(), 1)
  return stocktakes.value
    .filter(s => s.status === 'confirmed' && new Date(s.created_at) >= monthStart)
    .reduce((sum, s) => sum + s.total_loss, 0)
})

const filteredStocktakes = computed(() => {
  let result = [...stocktakes.value]
  if (dateRange.value.length === 2) {
    const [start, end] = dateRange.value
    result = result.filter(s => s.stocktake_date >= start && s.stocktake_date <= end)
  }
  if (filterStatus.value) {
    result = result.filter(s => s.status === filterStatus.value)
  }
  if (filterType.value) {
    result = result.filter(s => s.stocktake_type === filterType.value)
  }
  return result.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
})

const paginatedStocktakes = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredStocktakes.value.slice(start, end)
})

const totalProfit = computed(() => {
  return generatedItems.value.reduce((sum, item) => {
    const diff = (item.actual_quantity || 0) - item.system_quantity
    return sum + (diff > 0 ? diff * item.unit_price : 0)
  }, 0)
})

const totalLoss = computed(() => {
  return generatedItems.value.reduce((sum, item) => {
    const diff = (item.actual_quantity || 0) - item.system_quantity
    return sum + (diff < 0 ? Math.abs(diff) * item.unit_price : 0)
  }, 0)
})

const totalDiff = computed(() => totalProfit.value - totalLoss.value)

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

const getStocktakeTypeTag = (type: string): string => {
  switch (type) {
    case 'periodic': return 'warning'
    case 'monthly': return 'primary'
    case 'yearly': return 'success'
    default: return 'info'
  }
}

const getStocktakeTypeText = (type: string): string => {
  switch (type) {
    case 'periodic': return '周期盘点'
    case 'monthly': return '月度盘点'
    case 'yearly': return '年度盘点'
    default: return type
  }
}

const getDiffQty = (row: StocktakeItemGenerate): number => {
  return (row.actual_quantity || 0) - row.system_quantity
}

const getDiffAmount = (row: StocktakeItemGenerate): number => {
  return getDiffQty(row) * row.unit_price
}

const getDiffType = (row: StocktakeItemGenerate | string): string => {
  if (typeof row === 'string') return row
  const diff = getDiffQty(row)
  if (diff > 0.001) return 'profit'
  if (diff < -0.001) return 'loss'
  return 'normal'
}

const getDiffClass = (row: StocktakeItemGenerate): string => {
  const type = getDiffType(row)
  if (type === 'profit') return 'profit'
  if (type === 'loss') return 'loss'
  return ''
}

const calculateDiff = (_row: StocktakeItemGenerate) => {
}

const fetchStocktakes = async () => {
  loading.value = true
  try {
    const params: Record<string, any> = {}
    if (dateRange.value.length === 2) {
      params.start_date = dateRange.value[0]
      params.end_date = dateRange.value[1]
    }
    if (filterStatus.value) params.status = filterStatus.value
    if (filterType.value) params.stocktake_type = filterType.value

    const res = await api.getStocktakes(params)
    if (res.data.code === 0) {
      stocktakes.value = res.data.data || []
    } else {
      ElMessage.error(res.data.message || '获取盘点列表失败')
    }
  } catch (error) {
    console.error('获取盘点列表失败:', error)
    ElMessage.error('获取盘点列表失败')
  } finally {
    loading.value = false
  }
}

const resetFilters = () => {
  dateRange.value = []
  filterStatus.value = ''
  filterType.value = ''
  currentPage.value = 1
  fetchStocktakes()
}

const openCreateDialog = () => {
  showCreate.value = true
  showDetail.value = false
  createForm.stocktake_date = new Date().toISOString().split('T')[0]
  createForm.stocktake_type = 'periodic'
  createForm.operator = ''
  createForm.remark = ''
  generatedItems.value = []
}

const generateItems = async () => {
  if (!createForm.stocktake_date || !createForm.stocktake_type) {
    ElMessage.warning('请先填写盘点日期和盘点类型')
    return
  }
  try {
    const res = await api.generateStocktakeItems()
    if (res.data.code === 0) {
      generatedItems.value = (res.data.data || []).map((item: StocktakeItemGenerate) => ({
        ...item,
        actual_quantity: item.system_quantity
      }))
      ElMessage.success(`已生成 ${generatedItems.value.length} 项盘点明细`)
    } else {
      ElMessage.error(res.data.message || '生成盘点项失败')
    }
  } catch (error) {
    console.error('生成盘点项失败:', error)
    ElMessage.error('生成盘点项失败')
  }
}

const submitStocktake = async () => {
  if (!createForm.stocktake_date || !createForm.stocktake_type) {
    ElMessage.warning('请填写完整的盘点信息')
    return
  }
  if (generatedItems.value.length === 0) {
    ElMessage.warning('请先生成盘点项')
    return
  }

  try {
    await ElMessageBox.confirm(
      '确认提交该盘点单？提交后将创建草稿，可在详情页确认。',
      '确认提交',
      { type: 'info' }
    )

    const items = generatedItems.value.map(item => ({
      ingredient_type: item.ingredient_type,
      ingredient_id: item.ingredient_id,
      actual_quantity: item.actual_quantity || 0
    }))

    const req: StocktakeCreateRequest = {
      stocktake_date: createForm.stocktake_date,
      stocktake_type: createForm.stocktake_type,
      operator: createForm.operator || undefined,
      remark: createForm.remark || undefined,
      items
    }

    const res = await api.createStocktake(req)
    if (res.data.code === 0) {
      ElMessage.success('盘点单创建成功')
      showCreate.value = false
      fetchStocktakes()
    } else {
      ElMessage.error(res.data.message || '创建失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('创建失败:', error)
      ElMessage.error('创建失败')
    }
  }
}

const viewDetail = async (row: Stocktake) => {
  try {
    const res = await api.getStocktake(row.id)
    if (res.data.code === 0) {
      currentStocktake.value = res.data.data
      showDetail.value = true
      showCreate.value = false
    } else {
      ElMessage.error(res.data.message || '获取详情失败')
    }
  } catch (error) {
    console.error('获取详情失败:', error)
    ElMessage.error('获取详情失败')
  }
}

const confirmStocktake = async () => {
  if (!currentStocktake.value) return
  try {
    await ElMessageBox.confirm(
      '确认该盘点单？确认后将更新库存并生成盘盈盘亏记录，此操作不可撤销。',
      '确认盘点',
      { type: 'warning' }
    )

    const res = await api.confirmStocktake(currentStocktake.value.id, {
      status: 'confirmed',
      remark: '盘点确认'
    })
    if (res.data.code === 0) {
      ElMessage.success('盘点确认成功')
      currentStocktake.value = res.data.data
      fetchStocktakes()
    } else {
      ElMessage.error(res.data.message || '确认失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('确认失败:', error)
      ElMessage.error('确认失败')
    }
  }
}

const deleteStocktake = async (row: Stocktake) => {
  try {
    await ElMessageBox.confirm(
      '确认删除该盘点单？此操作不可撤销。',
      '确认删除',
      { type: 'warning' }
    )

    const res = await api.deleteStocktake(row.id)
    if (res.data.code === 0) {
      ElMessage.success('删除成功')
      showDetail.value = false
      currentStocktake.value = null
      fetchStocktakes()
    } else {
      ElMessage.error(res.data.message || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

onMounted(() => {
  fetchStocktakes()
})
</script>

<style lang="scss" scoped>
@import '@/styles/variables.scss';

.stocktake-container {
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

.stocktake-no {
  font-family: 'Consolas', monospace;
  font-weight: 600;
  color: $primary-gold;
  font-size: 13px;
}

.profit {
  color: $success;
  font-weight: 600;
}

.loss {
  color: $danger;
  font-weight: 600;
}

.system-qty {
  color: $text-secondary;
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

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid rgba(212, 175, 55, 0.2);
}

.detail-title {
  display: flex;
  align-items: center;
  gap: 16px;
  h2 {
    font-size: 20px;
    font-weight: 600;
    margin: 0;
  }
}

.back-btn {
  padding: 8px 12px;
}

.detail-actions {
  display: flex;
  gap: 12px;
}

.create-form {
  margin-bottom: 24px;
  .form-row {
    display: flex;
    gap: 24px;
    flex-wrap: wrap;
  }
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

.info-icon {
  color: $primary-gold;
}

.summary-section {
  display: flex;
  justify-content: flex-end;
  gap: 32px;
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid rgba(212, 175, 55, 0.2);
  .summary-item {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 14px;
    strong {
      font-size: 18px;
    }
  }
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 20px;
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

.remark-section {
  padding: 16px;
  background: rgba(212, 175, 55, 0.05);
  border-radius: $border-radius-md;
  border: 1px solid rgba(212, 175, 55, 0.1);
  display: flex;
  gap: 8px;
  .info-label {
    flex-shrink: 0;
  }
}
</style>
