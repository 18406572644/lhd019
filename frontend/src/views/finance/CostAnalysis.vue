<template>
  <div class="finance-container">
    <div class="page-header">
      <h1 class="gold-text">成本分析</h1>
      <p class="subtitle">分析各项成本构成和利润率</p>
    </div>

    <FinanceFilter @change="handleFilterChange" @export="handleExport" />

    <div id="cost-analysis-content">
      <div class="summary-cards" v-if="report">
        <el-card class="glass-card summary-card">
          <div class="card-content">
            <div class="card-icon revenue">
              <el-icon :size="28"><Money /></el-icon>
            </div>
            <div class="card-info">
              <p class="card-label">总营收</p>
              <p class="card-value gold-text">{{ formatCurrency(report.total_revenue) }}</p>
            </div>
          </div>
        </el-card>

        <el-card class="glass-card summary-card">
          <div class="card-content">
            <div class="card-icon cost">
              <el-icon :size="28"><ShoppingCart /></el-icon>
            </div>
            <div class="card-info">
              <p class="card-label">总成本</p>
              <p class="card-value gold-text">{{ formatCurrency(report.total_cost) }}</p>
              <p class="card-label-sub">占营收 {{ ((report.total_cost / report.total_revenue) * 100).toFixed(1) }}%</p>
            </div>
          </div>
        </el-card>

        <el-card class="glass-card summary-card">
          <div class="card-content">
            <div class="card-icon profit">
              <el-icon :size="28"><TrendCharts /></el-icon>
            </div>
            <div class="card-info">
              <p class="card-label">毛利润</p>
              <p class="card-value gold-text">{{ formatCurrency(report.gross_profit) }}</p>
              <p class="card-trend trend-up">
                毛利率 {{ report.gross_margin }}%
              </p>
            </div>
          </div>
        </el-card>

        <el-card class="glass-card summary-card">
          <div class="card-content">
            <div class="card-icon net">
              <el-icon :size="28"><Wallet /></el-icon>
            </div>
            <div class="card-info">
              <p class="card-label">净利润</p>
              <p class="card-value gold-text">{{ formatCurrency(report.net_profit) }}</p>
              <p class="card-trend trend-up">
                净利率 {{ report.net_margin }}%
              </p>
            </div>
          </div>
        </el-card>
      </div>

      <div class="charts-section">
        <el-card class="glass-card chart-card">
          <template #header>
            <div class="card-header">
              <span class="gold-text">成本构成</span>
              <span class="header-subtitle">各项成本占比</span>
            </div>
          </template>
          <div ref="costChartRef" class="chart-container" style="cursor: pointer;"></div>
        </el-card>

        <el-card class="glass-card chart-card">
          <template #header>
            <div class="card-header">
              <span class="gold-text">成本明细</span>
              <span class="header-subtitle">{{ getPeriodLabel(currentFilter.period) || '暂无数据' }}</span>
            </div>
          </template>
          <el-table :data="costDetails" style="width: 100%">
            <el-table-column prop="name" label="成本项目" />
            <el-table-column prop="value" label="金额" width="150">
              <template #default="{ row }">
                <span>{{ formatCurrency(row.value) }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="ratio" label="占比" width="120">
              <template #default="{ row }">
                <el-progress :percentage="row.ratio" :color="'#d4af37'" />
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </div>
    </div>

    <el-dialog v-model="drillDownVisible" :title="drillDownTitle" width="600px">
      <div class="drill-down-content">
        <div class="drill-down-item">
          <span class="label">成本项目:</span>
          <span class="value">{{ drillDownData?.name }}</span>
        </div>
        <div class="drill-down-item">
          <span class="label">金额:</span>
          <span class="value gold-text">{{ formatCurrency(drillDownData?.value || 0) }}</span>
        </div>
        <div class="drill-down-item">
          <span class="label">占比:</span>
          <el-progress :percentage="drillDownData?.ratio || 0" :color="'#d4af37'" style="width: 200px;" />
        </div>
        <div class="drill-down-item">
          <span class="label">同比变化:</span>
          <span class="value trend-up">+8.5%</span>
        </div>
        <div class="drill-down-item">
          <span class="label">主要构成:</span>
          <span class="value">{{ getCostBreakdown(drillDownData?.name) }}</span>
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
import { Money, ShoppingCart, TrendCharts, Wallet } from '@element-plus/icons-vue'
import { ElMessage, ElLoading } from 'element-plus'
import * as echarts from 'echarts'
import { api, type CostAnalysisReport, type FinanceFilterParams, type CostBreakdownItem } from '@/api'
import FinanceFilter from '@/components/FinanceFilter.vue'
import { exportToPDF, formatCurrency } from '@/utils/pdfExport'
import { getPeriodLabel } from '@/utils/dateUtils'
import { generateMockCostAnalysisReport } from '@/utils/mockData'

const costChartRef = ref<HTMLElement>()
let costChart: echarts.ECharts | null = null

const report = ref<CostAnalysisReport | null>(null)
const loading = ref(false)
const drillDownVisible = ref(false)
const drillDownData = ref<CostBreakdownItem | null>(null)
const drillDownTitle = ref('')

const currentFilter = reactive<FinanceFilterParams>({
  period: 'month'
})

const costDetails = computed(() => {
  if (!report.value) return []
  return report.value.cost_breakdown || []
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
    const res = await api.getCostAnalysisReport(params)
    report.value = res.data?.data || res.data || generateMockCostAnalysisReport(currentFilter.period)
    await nextTick()
    renderChart()
  } catch (error) {
    report.value = generateMockCostAnalysisReport(currentFilter.period)
    await nextTick()
    renderChart()
  } finally {
    loading.value = false
  }
}

const handleExport = async () => {
  const loadingInstance = ElLoading.service({ text: '正在生成PDF...' })
  try {
    await exportToPDF('cost-analysis-content', {
      title: '成本分析报表',
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
  if (!costChartRef.value || !report.value?.cost_breakdown) return

  if (!costChart) {
    costChart = echarts.init(costChartRef.value)
    costChart.on('click', handleChartClick)
  }

  const data = report.value.cost_breakdown
  const colors = ['#d4af37', '#667eea', '#f093fb', '#4facfe', '#00f2fe', '#43e97b', '#fa709a']

  const option: echarts.EChartsOption = {
    tooltip: {
      trigger: 'item',
      backgroundColor: 'rgba(0, 0, 0, 0.8)',
      borderColor: '#d4af37',
      textStyle: { color: '#fff' },
      formatter: (params: any) => {
        return `<div>
          <div style="font-weight: bold; margin-bottom: 4px;">${params.name}</div>
          <div>金额: ¥${params.value?.toLocaleString() || 0}</div>
          <div>占比: ${params.percent}%</div>
          <div style="font-size: 11px; color: #999; margin-top: 4px;">点击查看详情</div>
        </div>`
      }
    },
    legend: {
      orient: 'vertical',
      right: '5%',
      top: 'center',
      textStyle: { color: 'rgba(255,255,255,0.7)' }
    },
    series: [
      {
        name: '成本构成',
        type: 'pie',
        radius: ['40%', '70%'],
        center: ['35%', '50%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 8,
          borderColor: 'rgba(0,0,0,0.3)',
          borderWidth: 2
        },
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 16,
            fontWeight: 'bold',
            color: '#d4af37'
          }
        },
        labelLine: {
          show: false
        },
        data: data.map((d, i) => ({
          value: d.value,
          name: d.name,
          itemStyle: { color: colors[i % colors.length] }
        }))
      }
    ]
  }

  costChart.setOption(option)
}

const handleChartClick = (params: any) => {
  drillDownTitle.value = `${params.name} 详情`
  drillDownData.value = costDetails.value.find(d => d.name === params.name) || null
  drillDownVisible.value = true
}

const getCostBreakdown = (name?: string): string => {
  const map: Record<string, string> = {
    '原料成本': '水果、果汁、糖浆等',
    '基酒成本': '威士忌、伏特加、朗姆酒等',
    '浪费成本': '过期原料、制作失误等',
    '采购成本': '批量采购的原料和酒水',
    '运营成本': '房租、人工、水电等'
  }
  return map[name || ''] || '详见明细'
}

onMounted(() => {
  fetchData()
  window.addEventListener('resize', () => {
    costChart?.resize()
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

      &.cost {
        background: linear-gradient(135deg, #fa709a, #fee140);
      }

      &.profit {
        background: linear-gradient(135deg, #667eea, #764ba2);
      }

      &.net {
        background: linear-gradient(135deg, #43e97b, #38f9d7);
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

        &.trend-up {
          color: #67c23a;
        }
      }
    }
  }

  .charts-section {
    display: grid;
    grid-template-columns: 1fr 1fr;
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

  .drill-down-content {
    display: flex;
    flex-direction: column;
    gap: 16px;
    padding: 10px 0;

    .drill-down-item {
      display: flex;
      align-items: center;
      gap: 16px;

      .label {
        width: 100px;
        color: rgba(255, 255, 255, 0.6);
        font-size: 14px;
      }

      .value {
        color: rgba(255, 255, 255, 0.9);
        font-size: 14px;

        &.gold-text {
          color: #d4af37;
          font-weight: 600;
        }

        &.trend-up {
          color: #67c23a;
        }
      }
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

@media (max-width: 1200px) {
  .finance-container {
    .charts-section {
      grid-template-columns: 1fr;
    }
  }
}
</style>
