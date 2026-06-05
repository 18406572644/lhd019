<template>
  <div class="dashboard-container">
    <div class="page-header">
      <h1 class="gold-text">数据概览</h1>
      <p class="subtitle">实时掌握酒吧运营状况</p>
    </div>

    <div class="summary-cards">
      <el-card class="glass-card summary-card" v-for="card in summaryCards" :key="card.label">
        <div class="card-content">
          <div class="card-icon" :style="{ background: card.bgColor }">
            <el-icon :size="28">
              <component :is="card.icon" />
            </el-icon>
          </div>
          <div class="card-info">
            <p class="card-label">{{ card.label }}</p>
            <p class="card-value gold-text">{{ card.prefix }}{{ formatNumber(card.value) }}{{ card.suffix }}</p>
            <p class="card-trend" :class="card.trend > 0 ? 'trend-up' : 'trend-down'">
              <el-icon>
                <component :is="card.trend > 0 ? 'Top' : 'Bottom'" />
              </el-icon>
              {{ Math.abs(card.trend) }}% 较上周
            </p>
          </div>
        </div>
      </el-card>
    </div>

    <div class="charts-section">
      <el-card class="glass-card chart-card">
        <template #header>
          <div class="card-header">
            <span class="gold-text">每日收入趋势</span>
            <span class="header-subtitle">最近 7 天</span>
          </div>
        </template>
        <div ref="revenueChartRef" class="chart-container"></div>
      </el-card>

      <el-card class="glass-card chart-card">
        <template #header>
          <div class="card-header">
            <span class="gold-text">畅销饮品 TOP 5</span>
            <span class="header-subtitle">本周销量</span>
          </div>
        </template>
        <div ref="drinksChartRef" class="chart-container"></div>
      </el-card>
    </div>

    <div class="alerts-section">
      <el-card class="glass-card alert-card">
        <template #header>
          <div class="card-header">
            <span class="gold-text">
              <el-icon><WarningFilled /></el-icon>
              烈酒库存预警
            </span>
            <el-tag type="danger" effect="dark">{{ lowStockSpirits.length }} 件商品</el-tag>
          </div>
        </template>
        <el-table :data="lowStockSpirits" style="width: 100%" :empty-text="loading ? '加载中...' : '暂无数据'">
          <el-table-column prop="name" label="商品名称">
            <template #default="{ row }">
              <span class="spirit-name">{{ row.name }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="category" label="分类" width="120" />
          <el-table-column prop="stock" label="当前库存" width="120">
            <template #default="{ row }">
              <span class="stock-low">{{ row.stock }} {{ row.unit }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="min_stock" label="最低库存" width="120" />
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.stock <= row.min_stock * 0.5 ? 'danger' : 'warning'" effect="dark" size="small">
                {{ row.stock <= row.min_stock * 0.5 ? '紧急' : '预警' }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </el-card>

      <el-card class="glass-card alert-card">
        <template #header>
          <div class="card-header">
            <span class="gold-text">
              <el-icon><WarningFilled /></el-icon>
              配料库存预警
            </span>
            <el-tag type="danger" effect="dark">{{ lowStockIngredients.length }} 件商品</el-tag>
          </div>
        </template>
        <el-table :data="lowStockIngredients" style="width: 100%" :empty-text="loading ? '加载中...' : '暂无数据'">
          <el-table-column prop="name" label="配料名称">
            <template #default="{ row }">
              <span class="ingredient-name">{{ row.name }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="category" label="分类" width="120" />
          <el-table-column prop="stock" label="当前库存" width="120">
            <template #default="{ row }">
              <span class="stock-low">{{ row.stock }} {{ row.unit }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="min_stock" label="最低库存" width="120" />
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.stock <= row.min_stock * 0.5 ? 'danger' : 'warning'" effect="dark" size="small">
                {{ row.stock <= row.min_stock * 0.5 ? '紧急' : '预警' }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import * as echarts from 'echarts'
import {
  Money,
  ShoppingCart,
  User,
  PriceTag,
  Top,
  Bottom,
  WarningFilled
} from '@element-plus/icons-vue'
import { api, ApiResponse } from '@/api'

interface SummaryData {
  totalRevenue: number
  totalOrders: number
  totalCustomers: number
  avgOrderValue: number
  revenueTrend: number
  ordersTrend: number
  customersTrend: number
  avgValueTrend: number
}

interface DailyRevenueData {
  date: string
  revenue: number
}

interface TopDrinkData {
  name: string
  sales: number
  revenue: number
}

interface LowStockItem {
  id: number
  name: string
  category: string
  stock: number
  unit: string
  min_stock: number
}

interface SummaryCard {
  label: string
  value: number
  prefix: string
  suffix: string
  trend: number
  icon: string
  bgColor: string
}

const revenueChartRef = ref<HTMLElement>()
const drinksChartRef = ref<HTMLElement>()
let revenueChart: echarts.ECharts | null = null
let drinksChart: echarts.ECharts | null = null

const loading = ref(false)
const summaryData = ref<SummaryData>({
  totalRevenue: 0,
  totalOrders: 0,
  totalCustomers: 0,
  avgOrderValue: 0,
  revenueTrend: 0,
  ordersTrend: 0,
  customersTrend: 0,
  avgValueTrend: 0
})

const dailyRevenueData = ref<DailyRevenueData[]>([])
const topDrinksData = ref<TopDrinkData[]>([])
const lowStockSpirits = ref<LowStockItem[]>([])
const lowStockIngredients = ref<LowStockItem[]>([])

const summaryCards = ref<SummaryCard[]>([
  {
    label: '总收入',
    value: 0,
    prefix: '¥',
    suffix: '',
    trend: 0,
    icon: 'Money',
    bgColor: 'linear-gradient(135deg, #d4af37, #c9a227)'
  },
  {
    label: '总订单',
    value: 0,
    prefix: '',
    suffix: ' 单',
    trend: 0,
    icon: 'ShoppingCart',
    bgColor: 'linear-gradient(135deg, #e94560, #c73e54)'
  },
  {
    label: '总客户',
    value: 0,
    prefix: '',
    suffix: ' 人',
    trend: 0,
    icon: 'User',
    bgColor: 'linear-gradient(135deg, #533483, #432a6b)'
  },
  {
    label: '平均客单价',
    value: 0,
    prefix: '¥',
    suffix: '',
    trend: 0,
    icon: 'PriceTag',
    bgColor: 'linear-gradient(135deg, #0f3460, #0c2a4d)'
  }
])

const mockSummaryData: SummaryData = {
  totalRevenue: 128560,
  totalOrders: 1256,
  totalCustomers: 892,
  avgOrderValue: 102.4,
  revenueTrend: 12.5,
  ordersTrend: 8.3,
  customersTrend: 15.2,
  avgValueTrend: 3.8
}

const mockDailyRevenueData: DailyRevenueData[] = [
  { date: '05-30', revenue: 15200 },
  { date: '05-31', revenue: 18500 },
  { date: '06-01', revenue: 22300 },
  { date: '06-02', revenue: 19800 },
  { date: '06-03', revenue: 25600 },
  { date: '06-04', revenue: 31200 },
  { date: '06-05', revenue: 28900 }
]

const mockTopDrinksData: TopDrinkData[] = [
  { name: '经典马天尼', sales: 186, revenue: 27900 },
  { name: '莫吉托', sales: 165, revenue: 19800 },
  { name: '威士忌酸', sales: 142, revenue: 21300 },
  { name: '长岛冰茶', sales: 128, revenue: 16640 },
  { name: '曼哈顿', sales: 105, revenue: 18900 }
]

const mockLowStockSpirits: LowStockItem[] = [
  { id: 1, name: '麦卡伦12年单一麦芽', category: '威士忌', stock: 3, unit: '瓶', min_stock: 10 },
  { id: 2, name: '灰雁伏特加', category: '伏特加', stock: 2, unit: '瓶', min_stock: 8 },
  { id: 3, name: '百加得白朗姆', category: '朗姆酒', stock: 5, unit: '瓶', min_stock: 10 },
  { id: 4, name: '培恩龙舌兰', category: '龙舌兰', stock: 4, unit: '瓶', min_stock: 6 }
]

const mockLowStockIngredients: LowStockItem[] = [
  { id: 1, name: '新鲜薄荷叶', category: '香草', stock: 200, unit: 'g', min_stock: 500 },
  { id: 2, name: '青柠', category: '水果', stock: 15, unit: '个', min_stock: 50 },
  { id: 3, name: '安格斯特拉苦精', category: '调酒配料', stock: 2, unit: '瓶', min_stock: 5 },
  { id: 4, name: '红石榴糖浆', category: '糖浆', stock: 1, unit: '瓶', min_stock: 3 }
]

const formatNumber = (num: number | undefined | null): string => {
  if (num === undefined || num === null || isNaN(num)) {
    return '0'
  }
  const n = Number(num)
  if (n >= 10000) {
    return (n / 10000).toFixed(1) + '万'
  }
  return n.toLocaleString('zh-CN', { minimumFractionDigits: 0, maximumFractionDigits: 2 })
}

const mapSummaryData = (raw: any): SummaryData => ({
  totalRevenue: raw?.total_revenue ?? 0,
  totalOrders: raw?.total_orders ?? 0,
  totalCustomers: raw?.total_customers ?? 0,
  avgOrderValue: raw?.average_order ?? 0,
  revenueTrend: raw?.revenue_trend ?? 0,
  ordersTrend: raw?.orders_trend ?? 0,
  customersTrend: raw?.customers_trend ?? 0,
  avgValueTrend: raw?.avg_value_trend ?? 0
})

const safeName = (name: string | undefined | null, defaultVal: string = '未知饮品'): string => {
  if (!name || name === '' || name === 'null' || name === 'undefined') {
    return defaultVal
  }
  try {
    return String(name).trim()
  } catch {
    return defaultVal
  }
}

const mapTopDrinksData = (raw: any[]): TopDrinkData[] => {
  if (!Array.isArray(raw) || raw.length === 0) return []
  return raw.map((item: any) => ({
    name: safeName(item?.recipe_name ?? item?.name),
    sales: Number(item?.quantity ?? item?.sales ?? 0),
    revenue: Number(item?.revenue ?? 0)
  }))
}

const mapDailyRevenueData = (raw: any[]): DailyRevenueData[] => {
  if (!Array.isArray(raw)) return []
  return raw.map((item: any) => ({
    date: item?.date ? item.date.substring(5) : '未知',
    revenue: Number(item?.revenue ?? 0)
  }))
}

const mapLowStockData = (raw: any[]): LowStockItem[] => {
  if (!Array.isArray(raw)) return []
  return raw.map((item: any) => ({
    id: item?.id ?? 0,
    name: item?.name ?? '未知商品',
    category: item?.category ?? '其他',
    stock: Number(item?.stock_quantity ?? item?.stock ?? 0),
    unit: item?.unit ?? '',
    min_stock: Number(item?.min_stock ?? 0)
  }))
}

const fetchDashboardData = async () => {
  loading.value = true
  try {
    const [summaryRes, spiritsRes, ingredientsRes] = await Promise.all([
      api.getSummary(),
      api.getLowStockSpirits(),
      api.getLowStockIngredients()
    ])

    if (summaryRes.data.code === 0 && summaryRes.data.data) {
      const rawData = summaryRes.data.data
      summaryData.value = mapSummaryData(rawData)
      dailyRevenueData.value = rawData.daily_revenue?.length
        ? mapDailyRevenueData(rawData.daily_revenue)
        : mockDailyRevenueData
      topDrinksData.value = rawData.top_drinks?.length
        ? mapTopDrinksData(rawData.top_drinks)
        : mockTopDrinksData
    } else {
      summaryData.value = mockSummaryData
      dailyRevenueData.value = mockDailyRevenueData
      topDrinksData.value = mockTopDrinksData
    }

    if (spiritsRes.data.code === 0 && spiritsRes.data.data) {
      lowStockSpirits.value = mapLowStockData(spiritsRes.data.data)
    } else {
      lowStockSpirits.value = mockLowStockSpirits
    }

    if (ingredientsRes.data.code === 0 && ingredientsRes.data.data) {
      lowStockIngredients.value = mapLowStockData(ingredientsRes.data.data)
    } else {
      lowStockIngredients.value = mockLowStockIngredients
    }
  } catch (error) {
    console.error('获取仪表板数据失败，使用模拟数据:', error)
    summaryData.value = mockSummaryData
    dailyRevenueData.value = mockDailyRevenueData
    topDrinksData.value = mockTopDrinksData
    lowStockSpirits.value = mockLowStockSpirits
    lowStockIngredients.value = mockLowStockIngredients
  } finally {
    loading.value = false
    updateSummaryCards()
    await nextTick()
    initRevenueChart()
    initDrinksChart()
  }
}

const updateSummaryCards = () => {
  summaryCards.value[0].value = summaryData.value.totalRevenue
  summaryCards.value[0].trend = summaryData.value.revenueTrend
  summaryCards.value[1].value = summaryData.value.totalOrders
  summaryCards.value[1].trend = summaryData.value.ordersTrend
  summaryCards.value[2].value = summaryData.value.totalCustomers
  summaryCards.value[2].trend = summaryData.value.customersTrend
  summaryCards.value[3].value = summaryData.value.avgOrderValue
  summaryCards.value[3].trend = summaryData.value.avgValueTrend
}

const initRevenueChart = () => {
  if (!revenueChartRef.value) return

  if (revenueChart) {
    revenueChart.dispose()
  }

  revenueChart = echarts.init(revenueChartRef.value)

  const option: echarts.EChartsOption = {
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(22, 33, 62, 0.95)',
      borderColor: '#d4af37',
      borderWidth: 1,
      textStyle: {
        color: '#f5f5f5'
      },
      formatter: (params: any) => {
        const data = params[0]
        const value = data?.value ?? 0
        return `${data?.name ?? '未知'}<br/>收入: <span style="color: #d4af37; font-weight: 600;">¥${formatNumber(value)}</span>`
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      top: '10%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: dailyRevenueData.value.map(item => item.date),
      axisLine: {
        lineStyle: {
          color: 'rgba(255, 255, 255, 0.1)'
        }
      },
      axisLabel: {
        color: '#a0a0a0'
      }
    },
    yAxis: {
      type: 'value',
      axisLine: {
        show: false
      },
      axisTick: {
        show: false
      },
      splitLine: {
        lineStyle: {
          color: 'rgba(255, 255, 255, 0.05)'
        }
      },
      axisLabel: {
        color: '#a0a0a0',
        formatter: (value: number) => `¥${(value / 1000).toFixed(0)}k`
      }
    },
    series: [
      {
        name: '收入',
        type: 'line',
        smooth: true,
        symbol: 'circle',
        symbolSize: 8,
        lineStyle: {
          width: 3,
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 1,
            y2: 0,
            colorStops: [
              { offset: 0, color: '#d4af37' },
              { offset: 1, color: '#c9a227' }
            ]
          }
        },
        itemStyle: {
          color: '#d4af37',
          borderColor: '#1a1a2e',
          borderWidth: 2
        },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(212, 175, 55, 0.3)' },
              { offset: 1, color: 'rgba(212, 175, 55, 0)' }
            ]
          }
        },
        data: dailyRevenueData.value.map(item => item.revenue)
      }
    ]
  }

  revenueChart.setOption(option)
}

const initDrinksChart = () => {
  if (!drinksChartRef.value) return

  if (drinksChart) {
    drinksChart.dispose()
  }

  drinksChart = echarts.init(drinksChartRef.value)

  const chartData = topDrinksData.value.length > 0 ? topDrinksData.value : mockTopDrinksData
  const yAxisData = chartData.map(item => safeName(item.name)).reverse()
  const seriesData = chartData.map(item => Number(item.sales) || 0).reverse()

  const option: echarts.EChartsOption = {
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(22, 33, 62, 0.95)',
      borderColor: '#d4af37',
      borderWidth: 1,
      textStyle: {
        color: '#f5f5f5'
      },
      axisPointer: {
        type: 'shadow'
      },
      formatter: (params: any) => {
        const data = params[0]
        const name = safeName(data?.name)
        const value = data?.value ?? 0
        const drink = chartData.find(d => safeName(d.name) === name)
        const revenue = drink?.revenue ?? 0
        return `${name}<br/>销量: <span style="color: #d4af37; font-weight: 600;">${value} 杯</span><br/>收入: <span style="color: #e94560;">¥${formatNumber(revenue)}</span>`
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      top: '10%',
      containLabel: true
    },
    xAxis: {
      type: 'value',
      axisLine: {
        show: false
      },
      axisTick: {
        show: false
      },
      splitLine: {
        lineStyle: {
          color: 'rgba(255, 255, 255, 0.05)'
        }
      },
      axisLabel: {
        color: '#a0a0a0'
      }
    },
    yAxis: {
      type: 'category',
      data: yAxisData,
      axisLine: {
        lineStyle: {
          color: 'rgba(255, 255, 255, 0.1)'
        }
      },
      axisLabel: {
        color: '#f5f5f5',
        fontWeight: 500,
        formatter: (value: string) => safeName(value)
      }
    },
    series: [
      {
        name: '销量',
        type: 'bar',
        barWidth: '50%',
        itemStyle: {
          borderRadius: [0, 6, 6, 0],
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 1,
            y2: 0,
            colorStops: [
              { offset: 0, color: '#533483' },
              { offset: 0.5, color: '#d4af37' },
              { offset: 1, color: '#c9a227' }
            ]
          }
        },
        emphasis: {
          itemStyle: {
            shadowBlur: 20,
            shadowColor: 'rgba(212, 175, 55, 0.5)'
          }
        },
        label: {
          show: true,
          position: 'right',
          color: '#d4af37',
          fontWeight: 600,
          formatter: '{c} 杯'
        },
        data: seriesData
      }
    ]
  }

  drinksChart.setOption(option, true)
}

const handleResize = () => {
  revenueChart?.resize()
  drinksChart?.resize()
}

onMounted(() => {
  fetchDashboardData()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  revenueChart?.dispose()
  drinksChart?.dispose()
})
</script>

<style lang="scss" scoped>
@import '@/styles/variables.scss';

.dashboard-container {
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

    &.trend-up {
      color: $success;
    }

    &.trend-down {
      color: $danger;
    }
  }
}

.charts-section {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  margin-bottom: 24px;

  @media (max-width: 1200px) {
    grid-template-columns: 1fr;
  }
}

.chart-card {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .header-subtitle {
      font-size: 12px;
      color: $text-muted;
      font-weight: normal;
    }
  }
}

.chart-container {
  width: 100%;
  height: 320px;
}

.alerts-section {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;

  @media (max-width: 1200px) {
    grid-template-columns: 1fr;
  }
}

.alert-card {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 8px;

    .el-icon {
      margin-right: 6px;
      color: $warning;
    }
  }

  .spirit-name,
  .ingredient-name {
    font-weight: 500;
    color: $text-primary;
  }

  .stock-low {
    color: $danger;
    font-weight: 600;
  }
}
</style>
