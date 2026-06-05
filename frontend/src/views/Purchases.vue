<template>
  <div class="purchases-container">
    <div class="page-header">
      <h1 class="gold-text">采购台账</h1>
      <p class="subtitle">豪华酒吧采购管理系统</p>
    </div>

    <div class="summary-cards">
      <el-card class="glass-card summary-card">
        <div class="card-content">
          <div class="card-icon" style="background: linear-gradient(135deg, #d4af37, #c9a227);">
            <el-icon :size="28"><ShoppingCart /></el-icon>
          </div>
          <div class="card-info">
            <p class="card-label">采购单数</p>
            <p class="card-value gold-text">{{ totalCount }}</p>
            <p class="card-trend trend-up">
              <el-icon><TrendCharts /></el-icon>
              共 {{ filteredPurchases.length }} 条记录
            </p>
          </div>
        </div>
      </el-card>

      <el-card class="glass-card summary-card">
        <div class="card-content">
          <div class="card-icon" style="background: linear-gradient(135deg, #27ae60, #219a52);">
            <el-icon :size="28"><Money /></el-icon>
          </div>
          <div class="card-info">
            <p class="card-label">采购总金额</p>
            <p class="card-value" style="color: #27ae60;">¥{{ formatNumber(totalAmount) }}</p>
            <p class="card-trend trend-up">
              <el-icon><ShoppingCart /></el-icon>
              平均 ¥{{ formatNumber(averageAmount) }}/单
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
          <el-input
            v-model="supplierKeyword"
            placeholder="供应商关键词"
            clearable
            class="filter-item"
            @keyup.enter="fetchPurchases"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-button class="glow-button" @click="fetchPurchases">
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
            新建采购
          </el-button>
        </div>
      </div>

      <el-table
        :data="paginatedPurchases"
        style="width: 100%"
        row-key="id"
        v-loading="loading"
        :expand-row-keys="expandRowKeys"
        @expand-change="handleExpandChange"
      >
        <el-table-column type="expand" width="50">
          <template #default="{ row }">
            <div class="purchase-items-container">
              <div class="purchase-items-header">
                <span class="items-title">
                  <el-icon><List /></el-icon>
                  采购明细
                </span>
                <span class="items-count">共 {{ row.purchase_items?.length || 0 }} 项</span>
              </div>
              <el-table :data="row.purchase_items || []" size="small" style="width: 100%">
                <el-table-column prop="ingredient_type" label="类型" width="100" align="center">
                  <template #default="{ row }">
                    <el-tag :type="row.ingredient_type === 'spirit' ? 'warning' : 'success'" effect="dark" size="small">
                      {{ row.ingredient_type === 'spirit' ? '烈酒' : '辅料' }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="ingredient_name" label="原料名称" min-width="150">
                  <template #default="{ row }">
                    <span class="ingredient-name">{{ row.ingredient_name }}</span>
                  </template>
                </el-table-column>
                <el-table-column prop="quantity" label="数量" width="100" align="center" />
                <el-table-column prop="unit" label="单位" width="80" align="center" />
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
                <el-table-column prop="batch_no" label="批号" width="120" align="center">
                  <template #default="{ row }">
                    <span v-if="row.batch_no" class="batch-no">{{ row.batch_no }}</span>
                    <span v-else class="text-muted">-</span>
                  </template>
                </el-table-column>
                <el-table-column prop="expiry_date" label="有效期" width="120" align="center">
                  <template #default="{ row }">
                    <span v-if="row.expiry_date" class="expiry-date">{{ row.expiry_date }}</span>
                    <span v-else class="text-muted">-</span>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="purchase_no" label="采购单号" min-width="180">
          <template #default="{ row }">
            <span class="purchase-no">{{ row.purchase_no }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="supplier" label="供应商" min-width="150" />

        <el-table-column prop="purchase_date" label="采购日期" width="120" align="center" />

        <el-table-column prop="total_amount" label="总金额" width="130" align="right">
          <template #default="{ row }">
            <span class="total-amount">¥{{ row.total_amount.toFixed(2) }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="operator" label="操作员" width="100" align="center" />

        <el-table-column prop="remark" label="备注" min-width="150">
          <template #default="{ row }">
            <span v-if="row.remark" class="remark">{{ row.remark }}</span>
            <span v-else class="text-muted">-</span>
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
          :total="filteredPurchases.length"
          layout="total, sizes, prev, pager, next, jumper"
          background
        />
      </div>
    </el-card>

    <el-dialog
      v-model="createDialogVisible"
      title="新建采购单"
      width="1100px"
      :close-on-click-modal="false"
      @close="resetForm"
    >
      <el-form :model="purchaseForm" :rules="purchaseRules" ref="purchaseFormRef" label-width="100px">
        <div class="form-section">
          <div class="form-row">
            <el-form-item label="供应商" prop="supplier">
              <el-input v-model="purchaseForm.supplier" placeholder="请输入供应商名称" style="width: 300px;" />
            </el-form-item>
            <el-form-item label="采购日期" prop="purchase_date">
              <el-date-picker
                v-model="purchaseForm.purchase_date"
                type="date"
                placeholder="选择日期"
                value-format="YYYY-MM-DD"
                style="width: 200px;"
              />
            </el-form-item>
          </div>
          <div class="form-row">
            <el-form-item label="操作员" prop="operator">
              <el-input v-model="purchaseForm.operator" placeholder="请输入操作员" style="width: 200px;" />
            </el-form-item>
            <el-form-item label="备注" prop="remark" style="flex: 1;">
              <el-input v-model="purchaseForm.remark" placeholder="请输入备注" />
            </el-form-item>
          </div>
        </div>

        <div class="purchase-items-section">
          <div class="section-header">
            <span class="section-title">
              <el-icon><List /></el-icon>
              采购明细
            </span>
            <el-button type="primary" size="small" @click="addPurchaseItem">
              <el-icon><Plus /></el-icon>
              添加明细
            </el-button>
          </div>

          <el-table :data="purchaseForm.items" style="width: 100%" size="default">
            <el-table-column label="原料类型" width="140" align="center">
              <template #default="{ row, $index }">
                <el-select
                  v-model="row.ingredient_type"
                  placeholder="选择类型"
                  style="width: 100%"
                  @change="handleTypeChange($index)"
                >
                  <el-option label="烈酒" value="spirit" />
                  <el-option label="辅料" value="ingredient" />
                </el-select>
              </template>
            </el-table-column>

            <el-table-column label="原料名称" min-width="200">
              <template #default="{ row, $index }">
                <el-select
                  v-model="row.ingredient_id"
                  placeholder="请选择原料"
                  style="width: 100%"
                  filterable
                  @change="handleIngredientChange($index)"
                >
                  <el-option
                    v-for="item in getIngredientsByType(row.ingredient_type)"
                    :key="item.id"
                    :label="item.name"
                    :value="item.id"
                  >
                    <span>{{ item.name }}</span>
                    <span style="float: right; color: #d4af37; font-size: 12px;">
                      {{ item.category }}
                    </span>
                  </el-option>
                </el-select>
              </template>
            </el-table-column>

            <el-table-column label="数量" width="120" align="center">
              <template #default="{ row, $index }">
                <el-input-number
                  v-model="row.quantity"
                  :min="0"
                  :precision="2"
                  :step="1"
                  size="default"
                  style="width: 100%"
                  @change="calculateAmount"
                  controls-position="right"
                />
              </template>
            </el-table-column>

            <el-table-column label="单位" width="100" align="center">
              <template #default="{ row }">
                <el-select v-model="row.unit" placeholder="单位" style="width: 100%">
                  <el-option label="瓶" value="瓶" />
                  <el-option label="箱" value="箱" />
                  <el-option label="kg" value="kg" />
                  <el-option label="g" value="g" />
                  <el-option label="L" value="L" />
                  <el-option label="ml" value="ml" />
                  <el-option label="个" value="个" />
                  <el-option label="包" value="包" />
                </el-select>
              </template>
            </el-table-column>

            <el-table-column label="单价" width="140" align="right">
              <template #default="{ row }">
                <el-input-number
                  v-model="row.unit_price"
                  :min="0"
                  :precision="2"
                  :step="10"
                  size="default"
                  style="width: 100%"
                  @change="calculateAmount"
                  controls-position="right"
                />
              </template>
            </el-table-column>

            <el-table-column label="小计" width="140" align="right">
              <template #default="{ row }">
                <span class="subtotal">¥{{ (row.quantity * row.unit_price).toFixed(2) }}</span>
              </template>
            </el-table-column>

            <el-table-column label="批号" width="130" align="center">
              <template #default="{ row }">
                <el-input v-model="row.batch_no" placeholder="可选" size="default" />
              </template>
            </el-table-column>

            <el-table-column label="有效期" width="130" align="center">
              <template #default="{ row }">
                <el-date-picker
                  v-model="row.expiry_date"
                  type="date"
                  placeholder="可选"
                  value-format="YYYY-MM-DD"
                  size="default"
                  style="width: 100%"
                />
              </template>
            </el-table-column>

            <el-table-column label="操作" width="70" align="center">
              <template #default="{ $index }">
                <el-button
                  type="danger"
                  link
                  size="small"
                  :disabled="purchaseForm.items.length <= 1"
                  @click="removePurchaseItem($index)"
                >
                  <el-icon><Delete /></el-icon>
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <div class="amount-summary">
          <div class="amount-row actual-row">
            <span class="amount-label">总金额：</span>
            <span class="amount-value gold-text">¥{{ formTotalAmount.toFixed(2) }}</span>
          </div>
        </div>
      </el-form>

      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button class="glow-button" @click="submitPurchase" :loading="submitting">
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
  ShoppingCart,
  Money,
  TrendCharts,
  Search,
  Refresh,
  Plus,
  List,
  Delete
} from '@element-plus/icons-vue'
import { api, type Purchase, type PurchaseItem, type PurchaseForm as PurchaseFormType, type PurchaseItemForm, type Spirit, type Ingredient } from '@/api'

const loading = ref(false)
const submitting = ref(false)
const purchases = ref<Purchase[]>([])
const spirits = ref<Spirit[]>([])
const ingredients = ref<Ingredient[]>([])
const dateRange = ref<string[]>([])
const supplierKeyword = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const expandRowKeys = ref<number[]>([])
const createDialogVisible = ref(false)
const purchaseFormRef = ref<FormInstance>()

const purchaseForm = ref<PurchaseFormType>({
  supplier: '',
  purchase_date: '',
  operator: '',
  remark: '',
  items: [{
    ingredient_type: 'spirit',
    ingredient_id: 0,
    quantity: 0,
    unit: '瓶',
    unit_price: 0,
    batch_no: '',
    expiry_date: ''
  }]
})

const purchaseRules: FormRules = {
  supplier: [{ required: true, message: '请输入供应商', trigger: 'blur' }],
  purchase_date: [{ required: true, message: '请选择采购日期', trigger: 'change' }],
  operator: [{ required: true, message: '请输入操作员', trigger: 'blur' }]
}

const totalCount = computed(() => filteredPurchases.value.length)

const totalAmount = computed(() => {
  return filteredPurchases.value.reduce((sum, p) => sum + p.total_amount, 0)
})

const averageAmount = computed(() => {
  if (filteredPurchases.value.length === 0) return 0
  return totalAmount.value / filteredPurchases.value.length
})

const filteredPurchases = computed(() => {
  let result = [...purchases.value]

  if (dateRange.value && dateRange.value.length === 2) {
    const startDate = new Date(dateRange.value[0])
    startDate.setHours(0, 0, 0, 0)
    const endDate = new Date(dateRange.value[1])
    endDate.setHours(23, 59, 59, 999)

    result = result.filter(p => {
      const purchaseDate = new Date(p.purchase_date)
      return purchaseDate >= startDate && purchaseDate <= endDate
    })
  }

  if (supplierKeyword.value) {
    const keyword = supplierKeyword.value.toLowerCase()
    result = result.filter(p =>
      p.supplier.toLowerCase().includes(keyword)
    )
  }

  return result.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
})

const paginatedPurchases = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredPurchases.value.slice(start, end)
})

const formTotalAmount = computed(() => {
  return purchaseForm.value.items.reduce((sum, item) => {
    return sum + item.quantity * item.unit_price
  }, 0)
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

const getIngredientsByType = (type: 'spirit' | 'ingredient') => {
  if (type === 'spirit') {
    return spirits.value
  }
  return ingredients.value
}

const fetchPurchases = async () => {
  loading.value = true
  try {
    const params: Record<string, any> = {}
    if (dateRange.value && dateRange.value.length === 2) {
      params.start_date = dateRange.value[0]
      params.end_date = dateRange.value[1]
    }
    if (supplierKeyword.value) {
      params.supplier = supplierKeyword.value
    }

    const res = await api.getPurchases(params)
    if (res.data.code === 0) {
      purchases.value = res.data.data || []
    } else {
      ElMessage.error(res.data.message || '获取采购列表失败')
    }
  } catch (error) {
    console.error('获取采购列表失败:', error)
    ElMessage.error('获取采购列表失败')
  } finally {
    loading.value = false
  }
}

const fetchSpirits = async () => {
  try {
    const res = await api.getSpirits()
    if (res.data.code === 0) {
      spirits.value = res.data.data || []
    }
  } catch (error) {
    console.error('获取烈酒列表失败:', error)
    ElMessage.error('获取烈酒列表失败')
  }
}

const fetchIngredients = async () => {
  try {
    const res = await api.getIngredients()
    if (res.data.code === 0) {
      ingredients.value = res.data.data || []
    }
  } catch (error) {
    console.error('获取辅料列表失败:', error)
    ElMessage.error('获取辅料列表失败')
  }
}

const resetFilters = () => {
  dateRange.value = []
  supplierKeyword.value = ''
  currentPage.value = 1
  fetchPurchases()
}

const handleExpandChange = (row: Purchase, expandedRows: Purchase[]) => {
  expandRowKeys.value = expandedRows.map(r => r.id)
}

const openCreateDialog = () => {
  createDialogVisible.value = true
  fetchSpirits()
  fetchIngredients()
}

const resetForm = () => {
  purchaseForm.value = {
    supplier: '',
    purchase_date: '',
    operator: '',
    remark: '',
    items: [{
      ingredient_type: 'spirit',
      ingredient_id: 0,
      quantity: 0,
      unit: '瓶',
      unit_price: 0,
      batch_no: '',
      expiry_date: ''
    }]
  }
  purchaseFormRef.value?.resetFields()
}

const addPurchaseItem = () => {
  purchaseForm.value.items.push({
    ingredient_type: 'spirit',
    ingredient_id: 0,
    quantity: 0,
    unit: '瓶',
    unit_price: 0,
    batch_no: '',
    expiry_date: ''
  })
}

const removePurchaseItem = (index: number) => {
  if (purchaseForm.value.items.length > 1) {
    purchaseForm.value.items.splice(index, 1)
    calculateAmount()
  }
}

const handleTypeChange = (index: number) => {
  purchaseForm.value.items[index].ingredient_id = 0
  calculateAmount()
}

const handleIngredientChange = (index: number) => {
  const item = purchaseForm.value.items[index]
  const ingredientList = getIngredientsByType(item.ingredient_type)
  const ingredient = ingredientList.find(i => i.id === item.ingredient_id)
  if (ingredient) {
    item.unit = ingredient.unit
    if (item.unit_price === 0) {
      item.unit_price = ingredient.cost_price
    }
  }
  calculateAmount()
}

const calculateAmount = () => {
}

const validateItems = (): boolean => {
  for (let i = 0; i < purchaseForm.value.items.length; i++) {
    const item = purchaseForm.value.items[i]
    if (!item.ingredient_id) {
      ElMessage.error(`请选择第 ${i + 1} 项原料`)
      return false
    }
    if (item.quantity <= 0) {
      ElMessage.error(`第 ${i + 1} 项数量必须大于0`)
      return false
    }
    if (item.unit_price < 0) {
      ElMessage.error(`第 ${i + 1} 项单价不能为负数`)
      return false
    }
    if (!item.unit) {
      ElMessage.error(`请选择第 ${i + 1} 项单位`)
      return false
    }
  }
  return true
}

const submitPurchase = async () => {
  if (!purchaseFormRef.value) return

  try {
    await purchaseFormRef.value.validate()
  } catch {
    return
  }

  if (!validateItems()) {
    return
  }

  if (formTotalAmount.value <= 0) {
    ElMessage.error('采购总金额必须大于0')
    return
  }

  submitting.value = true
  try {
    const purchaseData: PurchaseFormType = {
      supplier: purchaseForm.value.supplier,
      purchase_date: purchaseForm.value.purchase_date,
      operator: purchaseForm.value.operator,
      remark: purchaseForm.value.remark,
      items: purchaseForm.value.items.map(item => ({
        ingredient_type: item.ingredient_type,
        ingredient_id: item.ingredient_id,
        quantity: item.quantity,
        unit: item.unit,
        unit_price: item.unit_price,
        batch_no: item.batch_no,
        expiry_date: item.expiry_date
      }))
    }

    const res = await api.createPurchase(purchaseData)
    if (res.data.code === 0) {
      ElMessage.success('采购单创建成功，库存已自动增加')
      createDialogVisible.value = false
      resetForm()
      fetchPurchases()
    } else {
      ElMessage.error(res.data.message || '创建采购单失败')
    }
  } catch (error: any) {
    console.error('创建采购单失败:', error)
    ElMessage.error(error.response?.data?.message || '创建采购单失败')
  } finally {
    submitting.value = false
  }
}

const handleDelete = async (row: Purchase) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除采购单 ${row.purchase_no} 吗？删除后库存将自动扣减。`,
      '删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const res = await api.deletePurchase(row.id)
    if (res.data.code === 0) {
      ElMessage.success('采购单删除成功，库存已扣减')
      fetchPurchases()
    } else {
      ElMessage.error(res.data.message || '删除采购单失败')
    }
  } catch {
  }
}

onMounted(() => {
  fetchPurchases()
})
</script>

<style lang="scss" scoped>
@import '@/styles/variables.scss';

.purchases-container {
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
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
  margin-bottom: 24px;

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

.purchase-no {
  font-family: 'Consolas', monospace;
  font-weight: 600;
  color: $primary-gold;
}

.total-amount {
  color: $primary-gold;
  font-weight: 600;
}

.ingredient-name {
  color: $text-primary;
  font-weight: 500;
}

.batch-no {
  font-family: 'Consolas', monospace;
  color: $text-secondary;
}

.expiry-date {
  color: $text-secondary;
}

.remark {
  color: $text-secondary;
  font-size: 13px;
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

.purchase-items-container {
  padding: 16px 40px;
  background: rgba(0, 0, 0, 0.2);
  border-radius: $border-radius-md;
  margin: 8px 0;

  .purchase-items-header {
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

  .subtotal {
    color: $primary-gold;
    font-weight: 600;
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
}

.purchase-items-section {
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
      min-width: 150px;
      text-align: right;
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
        font-size: 28px;
      }
    }
  }
}
</style>
