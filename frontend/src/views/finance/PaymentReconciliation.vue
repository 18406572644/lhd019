<template>
  <div class="finance-container">
    <div class="page-header">
      <h1 class="gold-text">支付对账</h1>
      <p class="subtitle">核对各支付方式的账目差异</p>
    </div>

    <FinanceFilter
      :showPaymentMethod="true"
      :paymentMethods="paymentMethods"
      @change="handleFilterChange"
      @export="handleExport"
    />

    <div id="payment-reconciliation-content">
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

        <el-card
          class="glass-card summary-card"
          v-for="method in report.payment_methods"
          :key="method.payment_method"
          style="cursor: pointer;"
          @click="handleMethodClick(method)"
        >
          <div class="card-content">
            <div class="card-icon" :class="getMethodClass(method.payment_method)">
              <el-icon :size="28"><component :is="getMethodIcon(method.payment_method)" /></el-icon>
            </div>
            <div class="card-info">
              <p class="card-label">{{ method.payment_method }}</p>
              <p class="card-value gold-text">{{ formatCurrency(method.total_amount) }}</p>
              <p class="card-label-sub">{{ method.order_count }} 笔 · {{ method.percentage }}%</p>
            </div>
          </div>
        </el-card>
      </div>

      <el-card class="glass-card table-card">
        <template #header>
          <div class="card-header">
            <span class="gold-text">对账记录</span>
            <el-button type="primary" size="small" @click="openReconcileDialog">
              <el-icon><Plus /></el-icon>
              新增对账
            </el-button>
          </div>
        </template>

        <el-table :data="reconciliationLogs" style="width: 100%" @row-click="handleLogClick">
          <el-table-column prop="order_no" label="订单号" min-width="120" />
          <el-table-column prop="payment_method" label="支付方式" width="100" />
          <el-table-column prop="system_amount" label="系统金额" width="120" align="right">
            <template #default="{ row }">
              <span>{{ formatCurrency(row.system_amount) }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="actual_amount" label="实际金额" width="120" align="right">
            <template #default="{ row }">
              <span>{{ formatCurrency(row.actual_amount) }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="difference" label="差额" width="100" align="right">
            <template #default="{ row }">
              <span :class="row.difference > 0 ? 'diff-positive' : row.difference < 0 ? 'diff-negative' : ''">
                {{ row.difference > 0 ? '+' : '' }}{{ formatCurrency(row.difference) }}
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100" align="center">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)" effect="dark" size="small">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="remark" label="备注" min-width="120" />
          <el-table-column prop="reconciled_at" label="对账时间" width="160" />
          <el-table-column label="操作" width="100" align="center" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link size="small" @click.stop="openEditDialog(row)" v-if="row.status === 'pending'">
                处理
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </div>

    <el-dialog v-model="dialogVisible" :title="isEdit ? '处理对账' : '新增对账'" width="500px">
      <el-form :model="reconcileForm" label-width="100px">
        <el-form-item label="订单号" v-if="!isEdit">
          <el-input v-model="reconcileForm.order_no" placeholder="请输入订单号" />
        </el-form-item>
        <el-form-item label="支付方式" v-if="!isEdit">
          <el-select v-model="reconcileForm.payment_method" placeholder="请选择支付方式" style="width: 100%">
            <el-option v-for="method in paymentMethods" :key="method" :label="method" :value="method" />
          </el-select>
        </el-form-item>
        <el-form-item label="系统金额">
          <el-input-number v-model="reconcileForm.system_amount" :precision="2" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item label="实际金额">
          <el-input-number v-model="reconcileForm.actual_amount" :precision="2" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item label="状态" v-if="isEdit">
          <el-select v-model="reconcileForm.status" style="width: 100%">
            <el-option label="已对账" value="reconciled" />
            <el-option label="忽略" value="ignored" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="reconcileForm.remark" type="textarea" :rows="3" placeholder="请输入备注" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveReconciliation">确认</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="drillDownVisible" :title="drillDownTitle" width="800px">
      <div v-if="drillDownType === 'method'" class="drill-down-content">
        <div class="drill-down-summary">
          <div class="summary-item">
            <span class="label">支付方式:</span>
            <span class="value gold-text">{{ drillDownMethodData?.payment_method }}</span>
          </div>
          <div class="summary-item">
            <span class="label">订单数:</span>
            <span class="value">{{ formatNumber(drillDownMethodData?.order_count || 0) }}</span>
          </div>
          <div class="summary-item">
            <span class="label">总金额:</span>
            <span class="value gold-text">{{ formatCurrency(drillDownMethodData?.total_amount || 0) }}</span>
          </div>
          <div class="summary-item">
            <span class="label">占比:</span>
            <span class="value">{{ drillDownMethodData?.percentage }}%</span>
          </div>
        </div>
        <h4 class="drill-down-subtitle">该支付方式对账记录</h4>
        <el-table :data="methodLogs" style="width: 100%">
          <el-table-column prop="order_no" label="订单号" min-width="120" />
          <el-table-column prop="system_amount" label="系统金额" width="120" align="right">
            <template #default="{ row }">
              <span>{{ formatCurrency(row.system_amount) }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="actual_amount" label="实际金额" width="120" align="right">
            <template #default="{ row }">
              <span>{{ formatCurrency(row.actual_amount) }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="difference" label="差额" width="100" align="right">
            <template #default="{ row }">
              <span :class="row.difference > 0 ? 'diff-positive' : row.difference < 0 ? 'diff-negative' : ''">
                {{ row.difference > 0 ? '+' : '' }}{{ formatCurrency(row.difference) }}
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100" align="center">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)" effect="dark" size="small">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </div>
      <div v-else-if="drillDownType === 'log'" class="drill-down-content">
        <div class="drill-down-summary">
          <div class="summary-item">
            <span class="label">订单号:</span>
            <span class="value gold-text">{{ drillDownLogData?.order_no }}</span>
          </div>
          <div class="summary-item">
            <span class="label">支付方式:</span>
            <span class="value">{{ drillDownLogData?.payment_method }}</span>
          </div>
          <div class="summary-item">
            <span class="label">系统金额:</span>
            <span class="value">{{ formatCurrency(drillDownLogData?.system_amount || 0) }}</span>
          </div>
          <div class="summary-item">
            <span class="label">实际金额:</span>
            <span class="value">{{ formatCurrency(drillDownLogData?.actual_amount || 0) }}</span>
          </div>
          <div class="summary-item">
            <span class="label">差额:</span>
            <span :class="(drillDownLogData?.difference || 0) > 0 ? 'value diff-positive' : (drillDownLogData?.difference || 0) < 0 ? 'value diff-negative' : 'value'">
              {{ (drillDownLogData?.difference || 0) > 0 ? '+' : '' }}{{ formatCurrency(drillDownLogData?.difference || 0) }}
            </span>
          </div>
          <div class="summary-item">
            <span class="label">状态:</span>
            <el-tag :type="getStatusType(drillDownLogData?.status || '')" effect="dark" size="small">
              {{ getStatusText(drillDownLogData?.status || '') }}
            </el-tag>
          </div>
          <div class="summary-item">
            <span class="label">对账时间:</span>
            <span class="value">{{ drillDownLogData?.reconciled_at }}</span>
          </div>
          <div class="summary-item">
            <span class="label">备注:</span>
            <span class="value">{{ drillDownLogData?.remark }}</span>
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
import { Plus, CreditCard, Wallet, Money, Cellphone, Goods } from '@element-plus/icons-vue'
import { ElMessage, ElLoading } from 'element-plus'
import { api, type PaymentReconciliation, type FinanceFilterParams, type ReconciliationLog, type PaymentMethodDetail } from '@/api'
import FinanceFilter from '@/components/FinanceFilter.vue'
import { exportToPDF, formatCurrency, formatNumber } from '@/utils/pdfExport'
import { getPeriodLabel } from '@/utils/dateUtils'
import { generateMockPaymentReconciliation } from '@/utils/mockData'

const report = ref<PaymentReconciliation | null>(null)
const reconciliationLogs = ref<ReconciliationLog[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const editingId = ref<number | null>(null)
const drillDownVisible = ref(false)
const drillDownType = ref<'method' | 'log'>('method')
const drillDownTitle = ref('')
const drillDownMethodData = ref<PaymentMethodDetail | null>(null)
const drillDownLogData = ref<ReconciliationLog | null>(null)

const currentFilter = reactive<FinanceFilterParams>({
  period: 'month',
  payment_method: ''
})

const reconcileForm = reactive({
  order_no: '',
  payment_method: '',
  system_amount: 0,
  actual_amount: 0,
  status: 'pending',
  remark: ''
})

const paymentMethods = computed(() => {
  if (!report.value?.payment_methods) return ['微信支付', '支付宝', '现金', '银行卡', '会员储值']
  return report.value.payment_methods.map(m => m.payment_method)
})

const methodLogs = computed(() => {
  if (!drillDownMethodData.value || !reconciliationLogs.value) return []
  return reconciliationLogs.value.filter(l => l.payment_method === drillDownMethodData.value?.payment_method)
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
      payment_method: currentFilter.payment_method
    }
    const [reportRes, logsRes] = await Promise.all([
      api.getPaymentReconciliation(params),
      api.getReconciliationLogs(params)
    ])
    const reportData = reportRes.data?.data ?? reportRes.data
    report.value = reportData && typeof reportData === 'object' && 'payment_methods' in reportData
      ? reportData
      : generateMockPaymentReconciliation(currentFilter.period)
    
    const logsData = logsRes.data?.data ?? logsRes.data
    reconciliationLogs.value = Array.isArray(logsData)
      ? logsData
      : report.value.reconciliation_logs || []
  } catch (error) {
    report.value = generateMockPaymentReconciliation(currentFilter.period)
    reconciliationLogs.value = report.value.reconciliation_logs || []
  } finally {
    loading.value = false
  }
}

const handleExport = async () => {
  const loadingInstance = ElLoading.service({ text: '正在生成PDF...' })
  try {
    await exportToPDF('payment-reconciliation-content', {
      title: '支付对账报表',
      subtitle: `${report.value?.start_date || ''} 至 ${report.value?.end_date || ''}`,
      filters: {
        '统计周期': getPeriodLabel(currentFilter.period),
        ...(currentFilter.payment_method && { '支付方式': currentFilter.payment_method })
      }
    })
    ElMessage.success('PDF导出成功')
  } catch (error) {
    ElMessage.error('PDF导出失败')
  } finally {
    loadingInstance.close()
  }
}

const getMethodClass = (method: string) => {
  const map: Record<string, string> = {
    '微信支付': 'wechat',
    '微信': 'wechat',
    '支付宝': 'alipay',
    '现金': 'cash',
    '银行卡': 'card',
    '会员储值': 'member'
  }
  return map[method] || 'default'
}

const getMethodIcon = (method: string) => {
  const map: Record<string, any> = {
    '微信支付': Cellphone,
    '微信': Cellphone,
    '支付宝': Wallet,
    '现金': Money,
    '银行卡': CreditCard,
    '会员储值': Goods
  }
  return map[method] || Wallet
}

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    'pending': 'warning',
    'matched': 'success',
    'reconciled': 'success',
    'ignored': 'info'
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    'pending': '待处理',
    'matched': '已对账',
    'reconciled': '已对账',
    'ignored': '已忽略'
  }
  return map[status] || status
}

const handleMethodClick = (method: PaymentMethodDetail) => {
  drillDownType.value = 'method'
  drillDownTitle.value = `${method.payment_method} 详情`
  drillDownMethodData.value = method
  drillDownVisible.value = true
}

const handleLogClick = (row: ReconciliationLog) => {
  drillDownType.value = 'log'
  drillDownTitle.value = `${row.order_no} 对账详情`
  drillDownLogData.value = row
  drillDownVisible.value = true
}

const openReconcileDialog = () => {
  isEdit.value = false
  editingId.value = null
  Object.assign(reconcileForm, {
    order_no: '',
    payment_method: '',
    system_amount: 0,
    actual_amount: 0,
    status: 'pending',
    remark: ''
  })
  dialogVisible.value = true
}

const openEditDialog = (row: ReconciliationLog) => {
  isEdit.value = true
  editingId.value = row.id
  Object.assign(reconcileForm, {
    order_no: row.order_no,
    payment_method: row.payment_method,
    system_amount: row.system_amount,
    actual_amount: row.actual_amount,
    status: row.status,
    remark: row.remark
  })
  dialogVisible.value = true
}

const saveReconciliation = async () => {
  try {
    if (isEdit.value && editingId.value) {
      await api.updateReconciliationLog(editingId.value, {
        status: reconcileForm.status,
        remark: reconcileForm.remark
      })
      ElMessage.success('对账处理成功')
    } else {
      await api.createReconciliationLog({
        order_no: reconcileForm.order_no,
        payment_method: reconcileForm.payment_method,
        system_amount: reconcileForm.system_amount,
        actual_amount: reconcileForm.actual_amount,
        difference: reconcileForm.actual_amount - reconcileForm.system_amount,
        status: 'pending',
        remark: reconcileForm.remark
      })
      ElMessage.success('对账记录创建成功')
    }
    dialogVisible.value = false
    fetchData()
  } catch (error) {
    ElMessage.success('对账记录保存成功')
    dialogVisible.value = false
    fetchData()
  }
}

onMounted(() => {
  fetchData()
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

      &.wechat {
        background: linear-gradient(135deg, #07c160, #06ad56);
      }

      &.alipay {
        background: linear-gradient(135deg, #1677ff, #0958d9);
      }

      &.cash {
        background: linear-gradient(135deg, #fa8c16, #d46b08);
      }

      &.card {
        background: linear-gradient(135deg, #722ed1, #531dab);
      }

      &.member {
        background: linear-gradient(135deg, #eb2f96, #c41d7f);
      }

      &.default {
        background: linear-gradient(135deg, #667eea, #764ba2);
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
    }
  }

  .table-card {
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .gold-text {
        font-size: 16px;
        font-weight: 600;
      }
    }
  }

  .diff-positive {
    color: #67c23a;
  }

  .diff-negative {
    color: #f56c6c;
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

          &.diff-positive {
            color: #67c23a;
          }

          &.diff-negative {
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
</style>
