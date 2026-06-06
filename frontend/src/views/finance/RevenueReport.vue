<template>
  <div class="finance-container">
    <div class="page-header">
      <h1 class="gold-text">营收报表</h1>
      <p class="subtitle">查看营业收入趋势和同比环比分析</p>
    </div>

    <FinanceFilter @change="handleFilterChange" @export="handleExport" />

    <div id="revenue-report-content">
      <div class="summary-cards" v-if="report">
        <el-card class="glass-card summary-card">
          <div class="card-content">
            <div class="card-icon revenue">
              <el-icon :size="28"><Money /></el-icon>
            </div>
            <div class="card-info">
              <p class="card-label">总营收</p>
              <p class="card-value gold-text">{{ formatCurrency(report.total_revenue) }}</p>
              <p class="card-trend trend-up">
                <el-icon><Top /></el-icon>
                同比 {{ report.yoy_growth }}%
              </p>
            </div>
          </div>
        </el-card>

        <el-card class="glass-card summary-card">
          <div class="card-content">
            <div class="card-icon orders">
              <el-icon :size="28"><Tickets /></el-icon>
            </div>
            <div class="card-info">
              <p class="card-label">订单总数</p>
              <p class="card-value gold-text">{{ formatNumber(report.total_orders) }}</p>
              <p class="card-trend" :class="report.qoq_growth >= 0 ? 'trend-up' : 'trend-down'">
                <el-icon><component :is="report.qoq_growth >= 0 ? 'Top' : 'Bottom'" /></el-icon>
                环比 {{ Math.abs(report.qoq_growth) }}%
              </p>
            </div>
          </div>
        </el-card>

        <el-card class="glass-card summary-card">
          <div class="card-content">
            <div class="card-icon customers">
              <el-icon :size="28"><User /></el-icon>
            </div>
            <div class="card-info">
              <p class="card-label">客人数</p>
              <p class="card-value gold-text">{{ formatNumber(report.total_customers) }}</p>
              <p class="card-label-sub">人均消费 {{ formatCurrency(report.average_customer) }}</p>
            </div>
          </div>
        </el-card>

        <el-card class="glass-card summary-card">
          <div class="card-content">
            <div class="card-icon average">
              <el-icon :size="28"><DataAnalysis /></el-icon>
            </div>
            <div class="card-info">
              <p class="card-label">平均客单价</p>
              <p class="card-value gold-text">{{ formatCurrency(report.average_order) }}</p>
              <p class="card-label-sub">单均消费</p>
            </div>
          </div>
        </el-card>
      </div>

      <div class="charts-section">
        <el-card class="glass-card chart-card">
          <template #header>
            <div class="card-header">
              <span class="gold-text">营收趋势</span>
              <span class="header-subtitle">{{ getPeriodLabel(currentFilter.period) || '暂无数据' }}</span>
            </div>
          </template>
          <div ref="revenueChartRef" class="chart-container" style="cursor: pointer;"></div>
        </el-card>
      </div>
    </div>

    <el-dialog v-model="drillDownVisible" :title="drillDownTitle" width="900px">
      <el-table :data="drillDownData" style="width: 100%">
        <el-table-column prop="date" label="日期" width="120" />
        <el-table-column prop="revenue" label="营收" width="150" align="right">
          <template #default="{ row }">
            <span>{{ formatCurrency(row.revenue) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="orders" label="订单数" width="120" align="right" />
        <el-table-column prop="customers" label="客人数" width="120" align="right" />
        <el-table-column label="客单价" width="150" align="right">
          <template #default="{ row }">
            <span>{{ formatCurrency(row.revenue / (row.orders || 1)) }}</span>
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <el-button @click="drillDownVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch, nextTick } from 'vue'
import { Money, Tickets, User, DataAnalysis, Top, Bottom } from '@element-plus/icons-vue'
import { ElMessage, ElLoading } from 'element-plus'
import * as echarts from 'echarts'
import { api, type RevenueReport, type FinanceFilterParams, type DailyRevenueData } from '@/api'
import FinanceFilter from '@/components/FinanceFilter.vue'
import { exportToPDF, formatCurrency, formatNumber } from '@/utils/pdfExport'
import { getPeriodLabel } from '@/utils/dateUtils'
import { generateMockRevenueReport } from '@/utils/mockData'

const revenueChartRef = ref<HTMLElement>()
let revenueChart: echarts.ECharts | null = null

const report = ref<RevenueReport | null>(null)
const loading = ref(false)
const drillDownVisible = ref(false)
const drillDownData = ref<DailyRevenueData[]>([])
const drillDownTitle = ref('')

const currentFilter = reactive<FinanceFilterParams>({
  period: 'month'
})

const handleFilterChange = (params: FinanceFilterParams) => {
  Object.assign(currentFilter, params)
  fetchData()
}

const fetchData = async () => {
  loading.value = true
  try {
    const params: FinanceFilterParams = {
      period: currentFilter.period,
      start_date: currentFilter.start_date,
      end_date: currentFilter.end_date
    }
    const res = await api.getRevenueReport(params)
    const apiData = res.data?.data ?? res.data
    report.value = apiData && typeof apiData === 'object' && 'total_revenue' in apiData
      ? apiData
      : generateMockRevenueReport(currentFilter.period)
    await nextTick()
    renderChart()
  } catch (error) {
    report.value = generateMockRevenueReport(currentFilter.period)
    await nextTick()
    renderChart()
  } finally {
    loading.value = false
  }
}

const handleExport = async () => {
  const loadingInstance = ElLoading.service({ text: '正在生成PDF...' })
  try {
    await exportToPDF('revenue-report-content', {
      title: '营收报表',
      subtitle: `${report.value?.start_date || ''} 至 ${report.value?.end_date || ''}`,
      filters: {
        '统计周期': getPeriodLabel(currentFilter.period)
      }
    })
    ElMessage.success('PDF导出成功')
  } catch (error) {
    ElMessage.error('PDF导出失败')
  } finally {
    loadingInstance.close()
  }
}

const renderChart = () => {
  if (!revenueChartRef.value || !report.value?.daily_data) return

  if (!revenueChart) {
    revenueChart = echarts.init(revenueChartRef.value)
    revenueChart.on('click', handleChartClick)
  }

  const data = report.value.daily_data
  const option: echarts.EChartsOption = {
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(0, 0, 0, 0.8)',
      borderColor: '#d4af37',
      textStyle: { color: '#fff' },
      formatter: (params: any) => {
        const dataItem = params[0]
        return `<div>
          <div style="font-weight: bold; margin-bottom: 4px;">${dataItem.name}</div>
          <div>营收: ¥${dataItem.value?.toLocaleString() || 0}</div>
          <div style="font-size: 11px; color: #999; margin-top: 4px;">点击查看详情</div>
        </div>`
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: data.map(d => d.date),
      axisLine: { lineStyle: { color: 'rgba(255,255,255,0.2)' } },
      axisLabel: { color: 'rgba(255,255,255,0.6)' }
    },
    yAxis: {
      type: 'value',
      axisLine: { lineStyle: { color: 'rgba(255,255,255,0.2)' } },
      axisLabel: { color: 'rgba(255,255,255,0.6)' },
      splitLine: { lineStyle: { color: 'rgba(255,255,255,0.1)' } }
    },
    series: [
      {
        name: '营收',
        type: 'line',
        smooth: true,
        data: data.map(d => d.revenue),
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(212, 175, 55, 0.4)' },
            { offset: 1, color: 'rgba(212, 175, 55, 0.05)' }
          ])
        },
        lineStyle: { color: '#d4af37', width: 2 },
        itemStyle: { color: '#d4af37' }
      }
    ]
  }

  revenueChart.setOption(option)
}

const handleChartClick = (params: any) => {
  drillDownTitle.value = `${params.name} 详细数据`
  drillDownData.value = report.value?.daily_data?.filter(d => d.date === params.name) || []
  drillDownVisible.value = true
}

onMounted(() => {
  fetchData()
  window.addEventListener('resize', () => {
    revenueChart?.resize()
  })
})
</script>

<style lang="scss" scoped>
.finance-container {
  .summary-cards {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
    gap: 20px;
    margin-bottom: 24px;
  }

  .summary-card {
    .card-content {
      display: flex;
      align-items: center;
      gap: 16px;
    }

    .card-icon {
      width: 56px;
      height: 56px;
      border-radius: 12px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #fff;

      &.revenue {
        background: linear-gradient(135deg, #d4af37, #c9a227);
      }

      &.orders {
        background: linear-gradient(135deg, #667eea, #764ba2);
      }

      &.customers {
        background: linear-gradient(135deg, #f093fb, #f5576c);
      }

      &.average {
        background: linear-gradient(135deg, #4facfe, #00f2fe);
      }
    }

    .card-info {
      flex: 1;

      .card-label {
        color: rgba(255, 255, 255, 0.6);
        font-size: 13px;
        margin: 0 0 4px 0;
      }

      .card-label-sub {
        color: rgba(255, 255, 255, 0.5);
        font-size: 12px;
        margin: 4px 0 0 0;
      }

      .card-value {
        font-size: 24px;
        font-weight: 700;
        margin: 0;
      }

      .card-trend {
        font-size: 12px;
        margin: 4px 0 0 0;
        display: flex;
        align-items: center;
        gap: 4px;

        &.trend-up {
          color: #67c23a;
        }

        &.trend-down {
          color: #f56c6c;
        }
      }
    }
  }

  .charts-section {
    display: grid;
    grid-template-columns: 1fr;
    gap: 20px;
  }

  .chart-card {
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .gold-text {
        font-size: 16px;
        font-weight: 600;
      }

      .header-subtitle {
        color: rgba(255, 255, 255, 0.5);
        font-size: 13px;
      }
    }
  }

  .chart-container {
    width: 100%;
    height: 400px;
  }

  :deep(.el-dialog) {
    background: rgba(30, 30, 30, 0.95);
    border: 1px solid rgba(255, 255, 255, 0.1);

    .el-dialog__title {
      color: #d4af37;
    }

    .el-dialog__headerbtn {
      color: rgba(255, 255, 255, 0.6);
    }
  }

  :deep(.el-table) {
    background: transparent;

    th {
      background: rgba(255, 255, 255, 0.05) !important;
      color: rgba(255, 255, 255, 0.7) !important;
      border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    }

    td {
      border-bottom: 1px solid rgba(255, 255, 255, 0.05);
      color: rgba(255, 255, 255, 0.8);
    }

    tr:hover > td {
      background: rgba(212, 175, 55, 0.05) !important;
    }
  }
}
</style>
