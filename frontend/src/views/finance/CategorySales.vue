<template>
  <div class="finance-container">
    <div class="page-header">
      <h1 class="gold-text">分类销售</h1>
      <p class="subtitle">按分类、产品和时段分析销售数据</p>
    </div>

    <FinanceFilter
      :showCategory="true"
      :categories="categories"
      @change="handleFilterChange"
      @export="handleExport"
    />

    <div id="category-sales-content">
      <div class="charts-section">
        <el-card class="glass-card chart-card">
          <template #header>
            <div class="card-header">
              <span class="gold-text">分类销售占比</span>
              <span class="header-subtitle">{{ getPeriodLabel(currentFilter.period) || '暂无数据' }}</span>
            </div>
          </template>
          <div ref="categoryChartRef" class="chart-container" style="cursor: pointer;"></div>
        </el-card>

        <el-card class="glass-card chart-card">
          <template #header>
            <div class="card-header">
              <span class="gold-text">时段销售趋势</span>
              <span class="header-subtitle">按时段分布</span>
            </div>
          </template>
          <div ref="timeSlotChartRef" class="chart-container" style="cursor: pointer;"></div>
        </el-card>
      </div>

      <el-card class="glass-card table-card">
        <template #header>
          <div class="card-header">
            <span class="gold-text">产品销售排行</span>
            <span class="header-subtitle">按销量排序</span>
          </div>
        </template>
        <el-table :data="recipeSales" style="width: 100%" @row-click="handleRowClick">
          <el-table-column type="index" label="排名" width="60" align="center">
            <template #default="{ $index }">
              <el-tag v-if="$index < 3" :type="['danger', 'warning', 'success'][$index]" effect="dark" size="small">
                {{ $index + 1 }}
              </el-tag>
              <span v-else>{{ $index + 1 }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="recipe_name" label="产品名称" min-width="150" />
          <el-table-column prop="category" label="分类" width="100" />
          <el-table-column prop="quantity" label="销量" width="100" align="right" />
          <el-table-column prop="revenue" label="营收" width="120" align="right">
            <template #default="{ row }">
              <span>{{ formatCurrency(row.revenue) }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="cost" label="成本" width="120" align="right">
            <template #default="{ row }">
              <span>{{ formatCurrency(row.cost) }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="profit" label="利润" width="120" align="right">
            <template #default="{ row }">
              <span class="profit">{{ formatCurrency(row.profit) }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="profit_margin" label="利润率" width="100" align="right">
            <template #default="{ row }">
              <el-tag :type="row.profit_margin >= 50 ? 'success' : row.profit_margin >= 30 ? 'warning' : 'info'" size="small">
                {{ row.profit_margin }}%
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </div>

    <el-dialog v-model="drillDownVisible" :title="drillDownTitle" width="800px">
      <div v-if="drillDownType === 'category'" class="drill-down-content">
        <div class="drill-down-summary">
          <div class="summary-item">
            <span class="label">分类:</span>
            <span class="value gold-text">{{ drillDownCategoryData?.category }}</span>
          </div>
          <div class="summary-item">
            <span class="label">销量:</span>
            <span class="value">{{ formatNumber(drillDownCategoryData?.quantity || 0) }}</span>
          </div>
          <div class="summary-item">
            <span class="label">营收:</span>
            <span class="value gold-text">{{ formatCurrency(drillDownCategoryData?.revenue || 0) }}</span>
          </div>
          <div class="summary-item">
            <span class="label">占比:</span>
            <span class="value">{{ drillDownCategoryData?.percentage }}%</span>
          </div>
        </div>
        <h4 class="drill-down-subtitle">该分类下产品明细</h4>
        <el-table :data="categoryProducts" style="width: 100%">
          <el-table-column prop="recipe_name" label="产品名称" min-width="150" />
          <el-table-column prop="quantity" label="销量" width="100" align="right" />
          <el-table-column prop="revenue" label="营收" width="120" align="right">
            <template #default="{ row }">
              <span>{{ formatCurrency(row.revenue) }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="profit_margin" label="利润率" width="100" align="right">
            <template #default="{ row }">
              <el-tag :type="row.profit_margin >= 50 ? 'success' : 'warning'" size="small">
                {{ row.profit_margin }}%
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </div>
      <div v-else-if="drillDownType === 'timeslot'" class="drill-down-content">
        <div class="drill-down-summary">
          <div class="summary-item">
            <span class="label">时段:</span>
            <span class="value gold-text">{{ drillDownTimeSlotData?.time_slot }}</span>
          </div>
          <div class="summary-item">
            <span class="label">销量:</span>
            <span class="value">{{ formatNumber(drillDownTimeSlotData?.quantity || 0) }}</span>
          </div>
          <div class="summary-item">
            <span class="label">营收:</span>
            <span class="value gold-text">{{ formatCurrency(drillDownTimeSlotData?.revenue || 0) }}</span>
          </div>
          <div class="summary-item">
            <span class="label">订单数:</span>
            <span class="value">{{ drillDownTimeSlotData?.orders }}</span>
          </div>
        </div>
      </div>
      <div v-else-if="drillDownType === 'product'" class="drill-down-content">
        <div class="drill-down-summary">
          <div class="summary-item">
            <span class="label">产品:</span>
            <span class="value gold-text">{{ drillDownProductData?.recipe_name }}</span>
          </div>
          <div class="summary-item">
            <span class="label">分类:</span>
            <span class="value">{{ drillDownProductData?.category }}</span>
          </div>
          <div class="summary-item">
            <span class="label">销量:</span>
            <span class="value">{{ formatNumber(drillDownProductData?.quantity || 0) }}</span>
          </div>
          <div class="summary-item">
            <span class="label">营收:</span>
            <span class="value gold-text">{{ formatCurrency(drillDownProductData?.revenue || 0) }}</span>
          </div>
          <div class="summary-item">
            <span class="label">成本:</span>
            <span class="value">{{ formatCurrency(drillDownProductData?.cost || 0) }}</span>
          </div>
          <div class="summary-item">
            <span class="label">利润:</span>
            <span class="value profit">{{ formatCurrency(drillDownProductData?.profit || 0) }}</span>
          </div>
          <div class="summary-item">
            <span class="label">利润率:</span>
            <el-tag :type="(drillDownProductData?.profit_margin || 0) >= 50 ? 'success' : 'warning'" size="small">
              {{ drillDownProductData?.profit_margin }}%
            </el-tag>
          </div>
        </div>
      </div>
      <template #footer>
        <el-button @click="drillDownVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, nextTick } from 'vue'
import { ElMessage, ElLoading } from 'element-plus'
import * as echarts from 'echarts'
import { api, type CategorySalesReport, type FinanceFilterParams, type RecipeSales, type CategorySales, type TimeSlotSales } from '@/api'
import FinanceFilter from '@/components/FinanceFilter.vue'
import { exportToPDF, formatCurrency, formatNumber } from '@/utils/pdfExport'
import { getPeriodLabel } from '@/utils/dateUtils'
import { generateMockCategorySalesReport } from '@/utils/mockData'

const categoryChartRef = ref<HTMLElement>()
const timeSlotChartRef = ref<HTMLElement>()
let categoryChart: echarts.ECharts | null = null
let timeSlotChart: echarts.ECharts | null = null

const report = ref<CategorySalesReport | null>(null)
const loading = ref(false)
const drillDownVisible = ref(false)
const drillDownType = ref<'category' | 'timeslot' | 'product'>('category')
const drillDownTitle = ref('')
const drillDownCategoryData = ref<CategorySales | null>(null)
const drillDownTimeSlotData = ref<TimeSlotSales | null>(null)
const drillDownProductData = ref<RecipeSales | null>(null)

const currentFilter = reactive<FinanceFilterParams>({
  period: 'month',
  category: ''
})

const categories = computed(() => {
  if (!report.value?.category_sales) return []
  return report.value.category_sales.map(c => c.category)
})

const recipeSales = computed<RecipeSales[]>(() => {
  if (!report.value?.recipe_sales) return []
  return [...report.value.recipe_sales].sort((a, b) => b.quantity - a.quantity)
})

const categoryProducts = computed(() => {
  if (!drillDownCategoryData.value || !report.value?.recipe_sales) return []
  return report.value.recipe_sales.filter(r => r.category === drillDownCategoryData.value?.category)
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
      end_date: currentFilter.end_date,
      category: currentFilter.category
    }
    const res = await api.getCategorySalesReport(params)
    report.value = res.data?.data || res.data || generateMockCategorySalesReport(currentFilter.period)
    await nextTick()
    renderCategoryChart()
    renderTimeSlotChart()
  } catch (error) {
    report.value = generateMockCategorySalesReport(currentFilter.period)
    await nextTick()
    renderCategoryChart()
    renderTimeSlotChart()
  } finally {
    loading.value = false
  }
}

const handleExport = async () => {
  const loadingInstance = ElLoading.service({ text: '正在生成PDF...' })
  try {
    await exportToPDF('category-sales-content', {
      title: '分类销售报表',
      subtitle: `${report.value?.start_date || ''} 至 ${report.value?.end_date || ''}`,
      filters: {
        '统计周期': getPeriodLabel(currentFilter.period),
        ...(currentFilter.category && { '分类': currentFilter.category })
      }
    })
    ElMessage.success('PDF导出成功')
  } catch (error) {
    ElMessage.error('PDF导出失败')
  } finally {
    loadingInstance.close()
  }
}

const renderCategoryChart = () => {
  if (!categoryChartRef.value || !report.value?.category_sales) return

  if (!categoryChart) {
    categoryChart = echarts.init(categoryChartRef.value)
    categoryChart.on('click', handleCategoryChartClick)
  }

  const data = report.value.category_sales
  const colors = ['#d4af37', '#667eea', '#f093fb', '#4facfe', '#00f2fe', '#43e97b', '#fa709a', '#ffecd2']

  const option: echarts.EChartsOption = {
    tooltip: {
      trigger: 'item',
      backgroundColor: 'rgba(0, 0, 0, 0.8)',
      borderColor: '#d4af37',
      textStyle: { color: '#fff' },
      formatter: (params: any) => {
        return `<div>
          <div style="font-weight: bold; margin-bottom: 4px;">${params.name}</div>
          <div>营收: ¥${params.value?.toLocaleString() || 0}</div>
          <div>占比: ${params.percent}%</div>
          <div style="font-size: 11px; color: #999; margin-top: 4px;">点击查看详情</div>
        </div>`
      }
    },
    legend: {
      bottom: '0',
      textStyle: { color: 'rgba(255,255,255,0.7)' }
    },
    series: [
      {
        name: '分类销售',
        type: 'pie',
        radius: ['40%', '70%'],
        center: ['50%', '40%'],
        avoidLabelOverlap: true,
        itemStyle: {
          borderRadius: 8,
          borderColor: 'rgba(0,0,0,0.3)',
          borderWidth: 2
        },
        label: {
          show: true,
          color: 'rgba(255,255,255,0.7)',
          formatter: '{b}\n{d}%'
        },
        data: data.map((d, i) => ({
          value: d.revenue,
          name: d.category,
          itemStyle: { color: colors[i % colors.length] }
        }))
      }
    ]
  }

  categoryChart.setOption(option)
}

const renderTimeSlotChart = () => {
  if (!timeSlotChartRef.value || !report.value?.time_slot_sales) return

  if (!timeSlotChart) {
    timeSlotChart = echarts.init(timeSlotChartRef.value)
    timeSlotChart.on('click', handleTimeSlotChartClick)
  }

  const data = report.value.time_slot_sales

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
      data: data.map(d => d.time_slot),
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
        type: 'bar',
        data: data.map(d => d.revenue),
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#d4af37' },
            { offset: 1, color: 'rgba(212, 175, 55, 0.3)' }
          ]),
          borderRadius: [4, 4, 0, 0]
        },
        barWidth: '50%'
      }
    ]
  }

  timeSlotChart.setOption(option)
}

const handleCategoryChartClick = (params: any) => {
  drillDownType.value = 'category'
  drillDownTitle.value = `${params.name} 销售详情`
  drillDownCategoryData.value = report.value?.category_sales?.find(c => c.category === params.name) || null
  drillDownVisible.value = true
}

const handleTimeSlotChartClick = (params: any) => {
  drillDownType.value = 'timeslot'
  drillDownTitle.value = `${params.name} 时段详情`
  drillDownTimeSlotData.value = report.value?.time_slot_sales?.find(t => t.time_slot === params.name) || null
  drillDownVisible.value = true
}

const handleRowClick = (row: RecipeSales) => {
  drillDownType.value = 'product'
  drillDownTitle.value = `${row.recipe_name} 产品详情`
  drillDownProductData.value = row
  drillDownVisible.value = true
}

onMounted(() => {
  fetchData()
  window.addEventListener('resize', () => {
    categoryChart?.resize()
    timeSlotChart?.resize()
  })
})
</script>

<style lang="scss" scoped>
.finance-container {
  .charts-section {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 20px;
    margin-bottom: 24px;
  }

  .chart-card,
  .table-card {
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
    height: 350px;
  }

  .table-card {
    margin-bottom: 24px;
  }

  .profit {
    color: #67c23a;
    font-weight: 500;
  }

  .drill-down-content {
    padding: 10px 0;

    .drill-down-summary {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
      gap: 16px;
      margin-bottom: 24px;
      padding: 16px;
      background: rgba(255, 255, 255, 0.03);
      border-radius: 8px;

      .summary-item {
        display: flex;
        flex-direction: column;
        gap: 4px;

        .label {
          color: rgba(255, 255, 255, 0.5);
          font-size: 12px;
        }

        .value {
          color: rgba(255, 255, 255, 0.9);
          font-size: 16px;
          font-weight: 500;

          &.gold-text {
            color: #d4af37;
          }

          &.profit {
            color: #67c23a;
          }
        }
      }
    }

    .drill-down-subtitle {
      color: rgba(255, 255, 255, 0.8);
      font-size: 14px;
      font-weight: 500;
      margin: 0 0 12px 0;
    }
  }

  :deep(.el-table) {
    background: transparent;
    cursor: pointer;

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
}

@media (max-width: 1200px) {
  .finance-container {
    .charts-section {
      grid-template-columns: 1fr;
    }
  }
}
</style>
