<template>
  <div class="orders-container">
    <div class="page-header">
      <h1 class="gold-text">订单管理</h1>
      <p class="subtitle">豪华酒吧订单系统</p>
    </div>

    <div class="summary-cards">
      <el-card class="glass-card summary-card">
        <div class="card-content">
          <div class="card-icon" style="background: linear-gradient(135deg, #d4af37, #c9a227);">
            <el-icon :size="28"><Money /></el-icon>
          </div>
          <div class="card-info">
            <p class="card-label">总收入</p>
            <p class="card-value gold-text">¥{{ formatNumber(totalRevenue) }}</p>
            <p class="card-trend trend-up">
              <el-icon><Top /></el-icon>
              共 {{ orders.length }} 单
            </p>
          </div>
        </div>
      </el-card>

      <el-card class="glass-card summary-card">
        <div class="card-content">
          <div class="card-icon" style="background: linear-gradient(135deg, #27ae60, #219a52);">
            <el-icon :size="28"><CircleCheck /></el-icon>
          </div>
          <div class="card-info">
            <p class="card-label">已完成</p>
            <p class="card-value" style="color: #27ae60;">{{ completedCount }} 单</p>
            <p class="card-trend trend-up">
              <el-icon><Money /></el-icon>
              ¥{{ formatNumber(completedRevenue) }}
            </p>
          </div>
        </div>
      </el-card>

      <el-card class="glass-card summary-card">
        <div class="card-content">
          <div class="card-icon" style="background: linear-gradient(135deg, #f39c12, #e67e22);">
            <el-icon :size="28"><Clock /></el-icon>
          </div>
          <div class="card-info">
            <p class="card-label">待处理</p>
            <p class="card-value" style="color: #f39c12;">{{ pendingCount }} 单</p>
            <p class="card-trend trend-up">
              <el-icon><Money /></el-icon>
              ¥{{ formatNumber(pendingRevenue) }}
            </p>
          </div>
        </div>
      </el-card>

      <el-card class="glass-card summary-card">
        <div class="card-content">
          <div class="card-icon" style="background: linear-gradient(135deg, #e74c3c, #c0392b);">
            <el-icon :size="28"><CircleClose /></el-icon>
          </div>
          <div class="card-info">
            <p class="card-label">已取消</p>
            <p class="card-value" style="color: #e74c3c;">{{ cancelledCount }} 单</p>
            <p class="card-trend trend-down">
              <el-icon><Money /></el-icon>
              ¥{{ formatNumber(cancelledRevenue) }}
            </p>
          </div>
        </div>
      </el-card>
    </div>

    <el-card class="glass-card">
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
          />
          <el-select
            v-model="statusFilter"
            placeholder="订单状态"
            clearable
            class="filter-item"
          >
            <el-option label="已完成" value="completed" />
            <el-option label="待处理" value="pending" />
            <el-option label="已取消" value="cancelled" />
          </el-select>
          <el-button class="glow-button" @click="fetchOrders">
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
            新建订单
          </el-button>
        </div>
      </div>

      <el-table
        :data="paginatedOrders"
        style="width: 100%"
        row-key="id"
        v-loading="loading"
        :expand-row-keys="expandRowKeys"
        @expand-change="handleExpandChange"
      >
        <el-table-column type="expand" width="50">
          <template #default="{ row }">
            <div class="order-items-container">
              <div class="order-items-header">
                <span class="items-title">
                  <el-icon><List /></el-icon>
                  订单明细
                </span>
                <span class="items-count">共 {{ row.order_items?.length || 0 }} 项</span>
              </div>
              <el-table :data="row.order_items || []" size="small" style="width: 100%">
                <el-table-column prop="recipe_name" label="饮品名称" min-width="150">
                  <template #default="{ row }">
                    <span class="recipe-name">{{ row.recipe_name }}</span>
                  </template>
                </el-table-column>
                <el-table-column prop="quantity" label="数量" width="100" align="center" />
                <el-table-column prop="unit_price" label="单价" width="120" align="right">
                  <template #default="{ row }">
                    <span>¥{{ row.unit_price.toFixed(2) }}</span>
                  </template>
                </el-table-column>
                <el-table-column prop="subtotal" label="小计" width="120" align="right">
                  <template #default="{ row }">
                    <span class="subtotal">¥{{ row.subtotal.toFixed(2) }}</span>
                  </template>
                </el-table-column>
                <el-table-column prop="remark" label="备注" min-width="150">
                  <template #default="{ row }">
                    <span v-if="row.remark" class="item-remark">{{ row.remark }}</span>
                    <span v-else class="text-muted">-</span>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="order_no" label="订单号" min-width="160">
          <template #default="{ row }">
            <span class="order-no">{{ row.order_no }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="table_no" label="桌号" width="100" align="center">
          <template #default="{ row }">
            <span v-if="row.table_no" class="table-no">{{ row.table_no }}</span>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>

        <el-table-column prop="customer_count" label="人数" width="80" align="center" />

        <el-table-column prop="total_amount" label="总金额" width="120" align="right">
          <template #default="{ row }">
            <span>¥{{ row.total_amount.toFixed(2) }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="discount" label="折扣" width="100" align="right">
          <template #default="{ row }">
            <span v-if="row.discount > 0" class="discount">-¥{{ row.discount.toFixed(2) }}</span>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>

        <el-table-column prop="actual_amount" label="实收金额" width="120" align="right">
          <template #default="{ row }">
            <span class="actual-amount">¥{{ row.actual_amount.toFixed(2) }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="payment_method" label="支付方式" width="100" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.payment_method" size="small" effect="dark">
              {{ getPaymentMethodLabel(row.payment_method) }}
            </el-tag>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>

        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" effect="dark" size="small">
              {{ getStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="created_at" label="创建时间" width="170">
          <template #default="{ row }">
            <span class="created-at">{{ formatDateTime(row.created_at) }}</span>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="100" align="center" fixed="right">
          <template #default="{ row }">
            <el-button
              type="danger"
              link
              size="small"
              :disabled="row.status === 'cancelled'"
              @click="handleDelete(row)"
            >
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
          :page-sizes="[10, 20, 50, 100]"
          :total="filteredOrders.length"
          layout="total, sizes, prev, pager, next, jumper"
          background
        />
      </div>
    </el-card>

    <el-dialog
      v-model="createDialogVisible"
      title="新建订单"
      width="900px"
      :close-on-click-modal="false"
      @close="resetForm"
    >
      <el-form :model="orderForm" :rules="orderRules" ref="orderFormRef" label-width="100px">
        <div class="form-section">
          <div class="form-row">
            <el-form-item label="桌号" prop="table_no">
              <el-input v-model="orderForm.table_no" placeholder="请输入桌号" style="width: 200px;" />
            </el-form-item>
            <el-form-item label="人数" prop="customer_count">
              <el-input-number v-model="orderForm.customer_count" :min="1" :max="20" />
            </el-form-item>
            <el-form-item label="折扣" prop="discount">
              <el-input-number
                v-model="orderForm.discount"
                :min="0"
                :precision="2"
                :step="10"
                placeholder="0"
              />
              <span class="form-tip">元</span>
            </el-form-item>
          </div>
          <div class="form-row">
            <el-form-item label="支付方式" prop="payment_method">
              <el-select v-model="orderForm.payment_method" placeholder="请选择支付方式" style="width: 200px;">
                <el-option label="现金" value="cash" />
                <el-option label="微信支付" value="wechat" />
                <el-option label="支付宝" value="alipay" />
                <el-option label="刷卡" value="card" />
                <el-option label="挂账" value="credit" />
              </el-select>
            </el-form-item>
            <el-form-item label="备注" prop="remark" style="flex: 1;">
              <el-input v-model="orderForm.remark" placeholder="请输入备注" />
            </el-form-item>
          </div>
        </div>

        <div class="order-items-section">
          <div class="section-header">
            <span class="section-title">
              <el-icon><List /></el-icon>
              订单明细
            </span>
            <el-button type="primary" size="small" @click="addOrderItem">
              <el-icon><Plus /></el-icon>
              添加饮品
            </el-button>
          </div>

          <el-table :data="orderForm.items" style="width: 100%" size="small">
            <el-table-column label="饮品名称" min-width="200">
              <template #default="{ row, $index }">
                <el-select
                  v-model="row.recipe_id"
                  placeholder="请选择饮品"
                  style="width: 100%"
                  @change="handleRecipeChange($index)"
                >
                  <el-option
                    v-for="recipe in recipes"
                    :key="recipe.id"
                    :label="recipe.name"
                    :value="recipe.id"
                  >
                    <span>{{ recipe.name }}</span>
                    <span style="float: right; color: #d4af37;">¥{{ recipe.price.toFixed(2) }}</span>
                  </el-option>
                </el-select>
              </template>
            </el-table-column>
            <el-table-column label="单价" width="120" align="right">
              <template #default="{ row }">
                <span>¥{{ getRecipePrice(row.recipe_id).toFixed(2) }}</span>
              </template>
            </el-table-column>
            <el-table-column label="数量" width="120" align="center">
              <template #default="{ row, $index }">
                <el-input-number
                  v-model="row.quantity"
                  :min="1"
                  :max="99"
                  size="small"
                  @change="calculateAmount"
                />
              </template>
            </el-table-column>
            <el-table-column label="小计" width="120" align="right">
              <template #default="{ row }">
                <span class="subtotal">¥{{ (getRecipePrice(row.recipe_id) * row.quantity).toFixed(2) }}</span>
              </template>
            </el-table-column>
            <el-table-column label="备注" min-width="150">
              <template #default="{ row }">
                <el-input v-model="row.remark" placeholder="特殊要求" size="small" />
              </template>
            </el-table-column>
            <el-table-column label="操作" width="70" align="center">
              <template #default="{ $index }">
                <el-button
                  type="danger"
                  link
                  size="small"
                  :disabled="orderForm.items.length <= 1"
                  @click="removeOrderItem($index)"
                >
                  <el-icon><Delete /></el-icon>
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <div class="amount-summary">
          <div class="amount-row">
            <span class="amount-label">总金额：</span>
            <span class="amount-value">¥{{ totalAmount.toFixed(2) }}</span>
          </div>
          <div class="amount-row discount-row">
            <span class="amount-label">折扣：</span>
            <span class="amount-value">-¥{{ orderForm.discount.toFixed(2) }}</span>
          </div>
          <div class="amount-row actual-row">
            <span class="amount-label">实收金额：</span>
            <span class="amount-value gold-text">¥{{ actualAmount.toFixed(2) }}</span>
          </div>
        </div>
      </el-form>

      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button class="glow-button" @click="submitOrder" :loading="submitting">
          确认创建
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import {
  Money,
  Top,
  Clock,
  CircleCheck,
  CircleClose,
  Search,
  Refresh,
  Plus,
  List,
  Delete
} from '@element-plus/icons-vue'
import { api } from '@/api'

interface OrderItem {
  id: number
  order_id: number
  recipe_id: number
  recipe_name: string
  quantity: number
  unit_price: number
  subtotal: number
  remark: string
  created_at: string
}

interface Order {
  id: number
  order_no: string
  table_no: string
  customer_count: number
  total_amount: number
  discount: number
  actual_amount: number
  payment_method: string
  status: 'completed' | 'pending' | 'cancelled'
  remark: string
  order_items?: OrderItem[]
  created_at: string
  updated_at: string
}

interface Recipe {
  id: number
  name: string
  category: string
  price: number
}

interface OrderItemForm {
  recipe_id: number
  quantity: number
  remark: string
}

interface OrderForm {
  table_no: string
  customer_count: number
  discount: number
  payment_method: string
  remark: string
  items: OrderItemForm[]
}

const loading = ref(false)
const submitting = ref(false)
const orders = ref<Order[]>([])
const recipes = ref<Recipe[]>([])
const dateRange = ref<string[]>([])
const statusFilter = ref<string>('')
const currentPage = ref(1)
const pageSize = ref(10)
const expandRowKeys = ref<number[]>([])
const createDialogVisible = ref(false)
const orderFormRef = ref<FormInstance>()

const orderForm = ref<OrderForm>({
  table_no: '',
  customer_count: 1,
  discount: 0,
  payment_method: '',
  remark: '',
  items: [{ recipe_id: 0, quantity: 1, remark: '' }]
})

const orderRules: FormRules = {
  table_no: [{ required: true, message: '请输入桌号', trigger: 'blur' }],
  customer_count: [{ required: true, message: '请输入人数', trigger: 'blur' }],
  payment_method: [{ required: true, message: '请选择支付方式', trigger: 'change' }]
}

const totalRevenue = computed(() => {
  return orders.value
    .filter(o => o.status !== 'cancelled')
    .reduce((sum, o) => sum + o.actual_amount, 0)
})

const completedCount = computed(() => orders.value.filter(o => o.status === 'completed').length)
const pendingCount = computed(() => orders.value.filter(o => o.status === 'pending').length)
const cancelledCount = computed(() => orders.value.filter(o => o.status === 'cancelled').length)

const completedRevenue = computed(() => {
  return orders.value
    .filter(o => o.status === 'completed')
    .reduce((sum, o) => sum + o.actual_amount, 0)
})

const pendingRevenue = computed(() => {
  return orders.value
    .filter(o => o.status === 'pending')
    .reduce((sum, o) => sum + o.actual_amount, 0)
})

const cancelledRevenue = computed(() => {
  return orders.value
    .filter(o => o.status === 'cancelled')
    .reduce((sum, o) => sum + o.actual_amount, 0)
})

const filteredOrders = computed(() => {
  let result = [...orders.value]

  if (statusFilter.value) {
    result = result.filter(o => o.status === statusFilter.value)
  }

  if (dateRange.value && dateRange.value.length === 2) {
    const startDate = new Date(dateRange.value[0])
    startDate.setHours(0, 0, 0, 0)
    const endDate = new Date(dateRange.value[1])
    endDate.setHours(23, 59, 59, 999)

    result = result.filter(o => {
      const orderDate = new Date(o.created_at)
      return orderDate >= startDate && orderDate <= endDate
    })
  }

  return result.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
})

const paginatedOrders = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredOrders.value.slice(start, end)
})

const totalAmount = computed(() => {
  return orderForm.value.items.reduce((sum, item) => {
    return sum + getRecipePrice(item.recipe_id) * item.quantity
  }, 0)
})

const actualAmount = computed(() => {
  return Math.max(0, totalAmount.value - orderForm.value.discount)
})

const formatNumber = (num: number): string => {
  return num.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

const formatDateTime = (dateStr: string): string => {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getStatusLabel = (status: string): string => {
  const map: Record<string, string> = {
    completed: '已完成',
    pending: '待处理',
    cancelled: '已取消'
  }
  return map[status] || status
}

const getStatusType = (status: string): 'success' | 'warning' | 'danger' | 'info' => {
  const map: Record<string, 'success' | 'warning' | 'danger' | 'info'> = {
    completed: 'success',
    pending: 'warning',
    cancelled: 'danger'
  }
  return map[status] || 'info'
}

const getPaymentMethodLabel = (method: string): string => {
  const map: Record<string, string> = {
    cash: '现金',
    wechat: '微信',
    alipay: '支付宝',
    card: '刷卡',
    credit: '挂账'
  }
  return map[method] || method
}

const getRecipePrice = (recipeId: number): number => {
  const recipe = recipes.value.find(r => r.id === recipeId)
  return recipe?.price || 0
}

const fetchOrders = async () => {
  loading.value = true
  try {
    const params: Record<string, any> = {}
    if (dateRange.value && dateRange.value.length === 2) {
      params.start_date = dateRange.value[0]
      params.end_date = dateRange.value[1]
    }
    if (statusFilter.value) {
      params.status = statusFilter.value
    }

    const res = await api.getOrders(params)
    if (res.data.code === 0) {
      orders.value = res.data.data || []
    } else {
      ElMessage.error(res.data.message || '获取订单失败')
    }
  } catch (error) {
    console.error('获取订单失败:', error)
    ElMessage.error('获取订单失败')
  } finally {
    loading.value = false
  }
}

const fetchRecipes = async () => {
  try {
    const res = await api.getRecipes()
    if (res.data.code === 0) {
      recipes.value = res.data.data || []
    }
  } catch (error) {
    console.error('获取配方失败:', error)
  }
}

const resetFilters = () => {
  dateRange.value = []
  statusFilter.value = ''
  currentPage.value = 1
  fetchOrders()
}

const handleExpandChange = (row: Order, expandedRows: Order[]) => {
  expandRowKeys.value = expandedRows.map(r => r.id)
}

const openCreateDialog = () => {
  createDialogVisible.value = true
  fetchRecipes()
}

const resetForm = () => {
  orderForm.value = {
    table_no: '',
    customer_count: 1,
    discount: 0,
    payment_method: '',
    remark: '',
    items: [{ recipe_id: 0, quantity: 1, remark: '' }]
  }
  orderFormRef.value?.resetFields()
}

const addOrderItem = () => {
  orderForm.value.items.push({ recipe_id: 0, quantity: 1, remark: '' })
}

const removeOrderItem = (index: number) => {
  if (orderForm.value.items.length > 1) {
    orderForm.value.items.splice(index, 1)
    calculateAmount()
  }
}

const handleRecipeChange = (index: number) => {
  calculateAmount()
}

const calculateAmount = () => {
}

const validateItems = (): boolean => {
  for (let i = 0; i < orderForm.value.items.length; i++) {
    const item = orderForm.value.items[i]
    if (!item.recipe_id) {
      ElMessage.error(`请选择第 ${i + 1} 项饮品`)
      return false
    }
    if (item.quantity < 1) {
      ElMessage.error(`第 ${i + 1} 项数量不能小于1`)
      return false
    }
  }
  return true
}

const submitOrder = async () => {
  if (!orderFormRef.value) return

  try {
    await orderFormRef.value.validate()
  } catch {
    return
  }

  if (!validateItems()) {
    return
  }

  if (orderForm.value.discount > totalAmount.value) {
    ElMessage.error('折扣金额不能大于总金额')
    return
  }

  submitting.value = true
  try {
    const orderData = {
      table_no: orderForm.value.table_no,
      customer_count: orderForm.value.customer_count,
      discount: orderForm.value.discount,
      payment_method: orderForm.value.payment_method,
      remark: orderForm.value.remark,
      items: orderForm.value.items.map(item => ({
        recipe_id: item.recipe_id,
        quantity: item.quantity,
        remark: item.remark
      }))
    }

    const res = await api.createOrder(orderData)
    if (res.data.code === 0) {
      ElMessage.success('订单创建成功，库存已自动扣减')
      createDialogVisible.value = false
      resetForm()
      fetchOrders()
    } else {
      ElMessage.error(res.data.message || '创建订单失败')
    }
  } catch (error: any) {
    console.error('创建订单失败:', error)
    ElMessage.error(error.response?.data?.message || '创建订单失败')
  } finally {
    submitting.value = false
  }
}

const handleDelete = async (row: Order) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除订单 ${row.order_no} 吗？删除后库存将自动恢复。`,
      '删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const res = await api.deleteOrder(row.id)
    if (res.data.code === 0) {
      ElMessage.success('订单删除成功，库存已恢复')
      fetchOrders()
    } else {
      ElMessage.error(res.data.message || '删除订单失败')
    }
  } catch {
  }
}

onMounted(() => {
  fetchOrders()
})
</script>

<style lang="scss" scoped>
@import '@/styles/variables.scss';

.orders-container {
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
    width: 240px;
  }
}

.order-no {
  font-family: 'Consolas', monospace;
  font-weight: 600;
  color: $primary-gold;
}

.table-no {
  display: inline-block;
  padding: 4px 12px;
  background: rgba(212, 175, 55, 0.2);
  border-radius: $border-radius-sm;
  color: $primary-gold;
  font-weight: 500;
}

.discount {
  color: $danger;
  font-weight: 500;
}

.actual-amount {
  color: $primary-gold;
  font-weight: 600;
}

.created-at {
  color: $text-secondary;
  font-size: 13px;
}

.text-muted {
  color: $text-muted;
}

.pagination-section {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.order-items-container {
  padding: 16px 40px;
  background: rgba(0, 0, 0, 0.2);
  border-radius: $border-radius-md;
  margin: 8px 0;

  .order-items-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
    padding-bottom: 8px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);

    .items-title {
      display: flex;
      align-items: center;
      gap: 8px;
      color: $primary-gold;
      font-weight: 600;
    }

    .items-count {
      color: $text-secondary;
      font-size: 13px;
    }
  }

  .recipe-name {
    color: $text-primary;
    font-weight: 500;
  }

  .subtotal {
    color: $primary-gold;
    font-weight: 600;
  }

  .item-remark {
    color: $text-secondary;
    font-size: 13px;
  }
}

.form-section {
  margin-bottom: 24px;

  .form-row {
    display: flex;
    gap: 24px;
    align-items: center;
    flex-wrap: wrap;
  }

  .form-tip {
    margin-left: 8px;
    color: $text-secondary;
    font-size: 13px;
  }
}

.order-items-section {
  margin-bottom: 24px;

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    padding-bottom: 8px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);

    .section-title {
      display: flex;
      align-items: center;
      gap: 8px;
      color: $primary-gold;
      font-weight: 600;
      font-size: 16px;
    }
  }

  .subtotal {
    color: $primary-gold;
    font-weight: 600;
  }
}

.amount-summary {
  padding: 20px;
  background: rgba(0, 0, 0, 0.3);
  border-radius: $border-radius-md;
  border: 1px solid rgba(212, 175, 55, 0.3);

  .amount-row {
    display: flex;
    justify-content: flex-end;
    align-items: center;
    gap: 16px;
    padding: 8px 0;

    .amount-label {
      color: $text-secondary;
      font-size: 14px;
    }

    .amount-value {
      font-size: 18px;
      font-weight: 600;
      min-width: 120px;
      text-align: right;
    }

    &.discount-row .amount-value {
      color: $danger;
    }

    &.actual-row {
      margin-top: 8px;
      padding-top: 16px;
      border-top: 1px solid rgba(255, 255, 255, 0.1);

      .amount-label {
        font-size: 16px;
        color: $text-primary;
      }

      .amount-value {
        font-size: 24px;
      }
    }
  }
}
</style>
