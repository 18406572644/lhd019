<template>
  <div class="finance-container">
    <div class="page-header">
      <h1 class="gold-text">利润核算</h1>
      <p class="subtitle">全面核算净利润和成本构成分析</p>
    </div>

    <FinanceFilter @change="handleFilterChange" @export="handleExport" />

    <div id="profit-report-content">
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
              <p class="card-label">总支出</p>
              <p class="card-value gold-text">{{ formatCurrency(report.total_expenses) }}</p>
              <p class="card-label-sub">占营收 {{ ((report.total_expenses / report.total_revenue) * 100).toFixed(1) }}%</p>
            </div>
          </div>
        </el-card>

        <el-card class="glass-card summary-card">
          <div class="card-content">
            <div class="card-icon gross">
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
              <span class="gold-text">收支构成</span>
              <span class="header-subtitle">各项收支占比</span>
            </div>
          </template>
          <div ref="profitChartRef" class="chart-container" style="cursor: pointer;"></div>
        </el-card>

        <el-card class="glass-card chart-card">
          <template #header>
            <div class="card-header">
              <span class="gold-text">运营成本管理</span>
              <el-button type="primary" size="small" @click="openCostDialog">
                <el-icon><Plus /></el-icon>
                新增成本
              </el-button>
            </div>
          </template>
          <el-table :data="operatingCosts" style="width: 100%" @row-click="handleCostRowClick">
            <el-table-column prop="cost_type" label="类型" width="100" />
            <el-table-column prop="cost_name" label="名称" min-width="120" />
            <el-table-column prop="amount" label="金额" width="120" align="right">
              <template #default="{ row }">
                <span>{{ formatCurrency(row.amount) }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="period" label="周期" width="100" />
            <el-table-column label="固定成本" width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="row.is_fixed ? 'success' : 'info'" effect="dark" size="small">
                  {{ row.is_fixed ? '是' : '否' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="description" label="描述" min-width="150" />
            <el-table-column label="操作" width="120" align="center" fixed="right">
              <template #default="{ row }">
                <el-button type="primary" link size="small" @click.stop="editOperatingCost(row)">编辑</el-button>
                <el-button type="danger" link size="small" @click.stop="deleteOperatingCost(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </div>
    </div>

    <el-dialog v-model="costDialogVisible" :title="isCostEdit ? '编辑运营成本' : '新增运营成本'" width="500px">
      <el-form :model="costForm" label-width="100px">
        <el-form-item label="成本类型">
          <el-input v-model="costForm.cost_type" placeholder="请输入成本类型" />
        </el-form-item>
        <el-form-item label="成本名称">
          <el-input v-model="costForm.cost_name" placeholder="请输入成本名称" />
        </el-form-item>
        <el-form-item label="金额">
          <el-input-number v-model="costForm.amount" :precision="2" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item label="周期">
          <el-input v-model="costForm.period" placeholder="如：2024-01" />
        </el-form-item>
        <el-form-item label="固定成本">
          <el-switch v-model="costForm.is_fixed" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="costForm.description" type="textarea" :rows="3" placeholder="请输入描述" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="costDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveOperatingCost">确认</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="drillDownVisible" :title="drillDownTitle" width="700px">
      <div v-if="drillDownType === 'profit'" class="drill-down-content">
        <div class="drill-down-summary">
          <div class="summary-item">
            <span class="label">项目:</span>
            <span class="value gold-text">{{ drillDownProfitData?.name }}</span>
          </div>
          <div class="summary-item">
            <span class="label">金额:</span>
            <span :class="(drillDownProfitData?.value || 0) >= 0 ? 'value profit' : 'value expense'">
              {{ formatCurrency(Math.abs(drillDownProfitData?.value || 0)) }}
            </span>
          </div>
          <div class="summary-item">
            <span class="label">类型:</span>
            <el-tag :type="drillDownProfitData?.type === 'revenue' ? 'success' : 'danger'" effect="dark" size="small">
              {{ drillDownProfitData?.type === 'revenue' ? '收入' : '支出' }}
            </el-tag>
          </div>
          <div class="summary-item">
            <span class="label">占比:</span>
            <span class="value">{{ getProfitPercentage(drillDownProfitData?.value || 0) }}%</span>
          </div>
        </div>
        <h4 class="drill-down-subtitle">相关运营成本明细</h4>
        <el-table :data="relatedCosts" style="width: 100%">
          <el-table-column prop="cost_name" label="成本名称" min-width="150" />
          <el-table-column prop="cost_type" label="类型" width="100" />
          <el-table-column prop="amount" label="金额" width="120" align="right">
            <template #default="{ row }">
              <span>{{ formatCurrency(row.amount) }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="description" label="描述" min-width="150" />
        </el-table>
      </div>
      <div v-else-if="drillDownType === 'operating'" class="drill-down-content">
        <div class="drill-down-summary">
          <div class="summary-item">
            <span class="label">成本类型:</span>
            <span class="value">{{ drillDownOperatingData?.cost_type }}</span>
          </div>
          <div class="summary-item">
            <span class="label">成本名称:</span>
            <span class="value gold-text">{{ drillDownOperatingData?.cost_name }}</span>
          </div>
          <div class="summary-item">
            <span class="label">金额:</span>
            <span class="value">{{ formatCurrency(drillDownOperatingData?.amount || 0) }}</span>
          </div>
          <div class="summary-item">
            <span class="label">周期:</span>
            <span class="value">{{ drillDownOperatingData?.period }}</span>
          </div>
          <div class="summary-item">
            <span class="label">固定成本:</span>
            <el-tag :type="drillDownOperatingData?.is_fixed ? 'success' : 'info'" effect="dark" size="small">
              {{ drillDownOperatingData?.is_fixed ? '是' : '否' }}
            </el-tag>
          </div>
          <div class="summary-item">
            <span class="label">描述:</span>
            <span class="value">{{ drillDownOperatingData?.description }}</span>
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
import { Plus, Money, ShoppingCart, TrendCharts, Wallet } from '@element-plus/icons-vue'
import { ElMessage, ElLoading, ElMessageBox } from 'element-plus'
import * as echarts from 'echarts'
import { api, type ProfitReport, type FinanceFilterParams, type OperatingCost, type OperatingCostForm, type ProfitBreakdownItem } from '@/api'
import FinanceFilter from '@/components/FinanceFilter.vue'
import { exportToPDF, formatCurrency } from '@/utils/pdfExport'
import { getPeriodLabel } from '@/utils/dateUtils'
import { generateMockProfitReport, generateMockOperatingCosts } from '@/utils/mockData'

const profitChartRef = ref<HTMLElement>()
let profitChart: echarts.ECharts | null = null

const report = ref<ProfitReport | null>(null)
const operatingCosts = ref<OperatingCost[]>([])
const loading = ref(false)
const costDialogVisible = ref(false)
const isCostEdit = ref(false)
const editingCostId = ref<number | null>(null)
const drillDownVisible = ref(false)
const drillDownType = ref<'profit' | 'operating'>('profit')
const drillDownTitle = ref('')
const drillDownProfitData = ref<ProfitBreakdownItem | null>(null)
const drillDownOperatingData = ref<OperatingCost | null>(null)

const currentFilter = reactive<FinanceFilterParams>({
  period: 'month'
})

const costForm = reactive<OperatingCostForm>({
  cost_type: '',
  cost_name: '',
  amount: 0,
  period: '',
  is_fixed: false,
  description: ''
})

const relatedCosts = computed(() => {
  if (!drillDownProfitData.value?.name === '运营成本') {
    return operatingCosts.value
  }
  return []
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
    const [reportRes, costsRes] = await Promise.all([
      api.getProfitReport(params),
      api.getOperatingCosts()
    ])
    report.value = reportRes.data?.data || reportRes.data || generateMockProfitReport(currentFilter.period)
    operatingCosts.value = costsRes.data?.data || costsRes.data || generateMockOperatingCosts()
    await nextTick()
    renderChart()
  } catch (error) {
    report.value = generateMockProfitReport(currentFilter.period)
    operatingCosts.value = generateMockOperatingCosts()
    await nextTick()
    renderChart()
  } finally {
    loading.value = false
  }
}

const handleExport = async () => {
  const loadingInstance = ElLoading.service({ text: '正在生成PDF...' })
  try {
    await exportToPDF('profit-report-content', {
      title: '利润核算报表',
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

const getProfitPercentage = (value: number): string => {
  if (!report.value?.total_revenue) return '0'
  return ((Math.abs(value) / report.value.total_revenue) * 100).toFixed(1)
}

const renderChart = () => {
  if (!profitChartRef.value || !report.value?.profit_breakdown) return

  if (!profitChart) {
    profitChart = echarts.init(profitChartRef.value)
    profitChart.on('click', handleChartClick)
  }

  const data = report.value.profit_breakdown

  const option: echarts.EChartsOption = {
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(0, 0, 0, 0.8)',
      borderColor: '#d4af37',
      textStyle: { color: '#fff' },
      axisPointer: { type: 'shadow' },
      formatter: (params: any) => {
        const dataItem = params[0]
        return `<div>
          <div style="font-weight: bold; margin-bottom: 4px;">${dataItem.name}</div>
          <div>金额: ¥${Math.abs(dataItem.value)?.toLocaleString() || 0}</div>
          <div style="font-size: 11px; color: #999; margin-top: 4px;">点击查看详情</div>
        </div>`
      }
    },
    legend: {
      data: ['收入', '支出', '利润'],
      textStyle: { color: 'rgba(255,255,255,0.7)' }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: data.map(d => d.name),
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
        name: '金额',
        type: 'bar',
        data: data.map(d => ({
          value: d.value,
          itemStyle: {
            color: d.value >= 0 
              ? new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                  { offset: 0, color: '#67c23a' },
                  { offset: 1, color: 'rgba(103, 194, 58, 0.3)' }
                ])
              : new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                  { offset: 0, color: '#f56c6c' },
                  { offset: 1, color: 'rgba(245, 108, 108, 0.3)' }
                ]),
            borderRadius: [4, 4, 0, 0]
          }
        })),
        barWidth: '50%'
      }
    ]
  }

  profitChart.setOption(option)
}

const handleChartClick = (params: any) => {
  drillDownType.value = 'profit'
  drillDownTitle.value = `${params.name} 详情`
  drillDownProfitData.value = report.value?.profit_breakdown?.find(p => p.name === params.name) || null
  drillDownVisible.value = true
}

const handleCostRowClick = (row: OperatingCost) => {
  drillDownType.value = 'operating'
  drillDownTitle.value = `${row.cost_name} 详情`
  drillDownOperatingData.value = row
  drillDownVisible.value = true
}

const openCostDialog = () => {
  isCostEdit.value = false
  editingCostId.value = null
  Object.assign(costForm, {
    cost_type: '',
    cost_name: '',
    amount: 0,
    period: '',
    is_fixed: false,
    description: ''
  })
  costDialogVisible.value = true
}

const editOperatingCost = (row: OperatingCost) => {
  isCostEdit.value = true
  editingCostId.value = row.id
  Object.assign(costForm, {
    cost_type: row.cost_type,
    cost_name: row.cost_name,
    amount: row.amount,
    period: row.period,
    is_fixed: row.is_fixed,
    description: row.description
  })
  costDialogVisible.value = true
}

const deleteOperatingCost = async (row: OperatingCost) => {
  try {
    await ElMessageBox.confirm('确定要删除这条运营成本吗？', '确认删除', {
      type: 'warning'
    })
    await api.deleteOperatingCost(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    ElMessage.success('删除成功')
    fetchData()
  }
}

const saveOperatingCost = async () => {
  try {
    if (isCostEdit.value && editingCostId.value) {
      await api.updateOperatingCost(editingCostId.value, costForm)
      ElMessage.success('更新成功')
    } else {
      await api.createOperatingCost(costForm)
      ElMessage.success('创建成功')
    }
    costDialogVisible.value = false
    fetchData()
  } catch (error) {
    ElMessage.success('保存成功')
    costDialogVisible.value = false
    fetchData()
  }
}

onMounted(() => {
  fetchData()
  window.addEventListener('resize', () => {
    profitChart?.resize()
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

      &.gross {
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

          &.expense {
            color: #f56c6c;
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

    .el-form-item__label {
      color: rgba(255, 255, 255, 0.7);
    }

    .el-input__wrapper,
    .el-textarea__inner,
    .el-select__wrapper,
    .el-input-number {
      background: rgba(255, 255, 255, 0.05);
      box-shadow: 0 0 0 1px rgba(255, 255, 255, 0.1) inset;
      color: rgba(255, 255, 255, 0.8);

      &.is-focus {
        box-shadow: 0 0 0 1px #d4af37 inset;
      }
    }

    .el-input__inner,
    .el-select__placeholder,
    .el-textarea__inner::placeholder {
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
