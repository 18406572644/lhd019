<template>
  <div class="waste-page">
    <div class="page-header">
      <div class="header-title">
        <el-icon :size="28" color="#d4af37"><Delete /></el-icon>
        <h2 class="title gold-text">原料损耗管理</h2>
      </div>
      <div class="header-stats">
        <div class="stat-card">
          <el-icon class="stat-icon"><Document /></el-icon>
          <div class="stat-content">
            <span class="stat-label">总损耗记录</span>
            <span class="stat-value">{{ totalCount }}</span>
          </div>
        </div>
        <div class="stat-card danger">
          <el-icon class="stat-icon"><Money /></el-icon>
          <div class="stat-content">
            <span class="stat-label">总损耗成本</span>
            <span class="stat-value cost-value">¥{{ totalCost.toFixed(2) }}</span>
          </div>
        </div>
      </div>
    </div>

    <el-card class="filter-card">
      <div class="filter-bar">
        <div class="filter-left">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
            class="filter-date"
            @change="handleFilterChange"
          />
          <el-select
            v-model="filterType"
            placeholder="选择配料类型"
            clearable
            class="filter-select"
            @change="handleFilterChange"
          >
            <el-option label="基酒" value="spirit" />
            <el-option label="食材" value="ingredient" />
          </el-select>
          <el-button class="glow-button" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button class="reset-button" @click="handleReset">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
        </div>
        <div class="filter-right">
          <el-button class="glow-button" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新增损耗
          </el-button>
        </div>
      </div>
    </el-card>

    <el-card class="table-card">
      <div class="table-header">
        <h3 class="table-title">
          <el-icon :size="20"><List /></el-icon>
          损耗记录列表
        </h3>
      </div>

      <el-table
        :data="paginatedData"
        v-loading="loading"
        stripe
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="ingredient_name" label="原料名称" min-width="150">
          <template #default="{ row }">
            <span class="name-text">{{ row.ingredient_name }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="ingredient_type" label="类型" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.ingredient_type === 'spirit' ? 'warning' : 'success'" effect="dark">
              {{ row.ingredient_type === 'spirit' ? '基酒' : '食材' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="数量" width="100" align="center" />
        <el-table-column prop="unit" label="单位" width="80" align="center" />
        <el-table-column prop="reason" label="损耗原因" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="getReasonTagType(row.reason)" effect="dark">
              {{ row.reason }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="cost" label="成本" width="120" align="right">
          <template #default="{ row }">
            <span class="cost-text">¥{{ row.cost?.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="operator" label="操作人" width="100" align="center" />
        <el-table-column prop="remark" label="备注" min-width="150" show-overflow-tooltip />
        <el-table-column prop="created_at" label="创建时间" width="180" align="center" />
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

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.currentPage"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          background
        />
      </div>
    </el-card>

    <el-dialog
      v-model="dialogVisible"
      title="新增损耗记录"
      width="550px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
        class="waste-form"
      >
        <el-form-item label="配料类型" prop="ingredient_type">
          <el-radio-group v-model="formData.ingredient_type" @change="handleTypeChange">
            <el-radio value="spirit">基酒</el-radio>
            <el-radio value="ingredient">食材</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="选择原料" prop="ingredient_id">
          <el-select
            v-model="formData.ingredient_id"
            placeholder="请选择原料"
            style="width: 100%"
            @change="handleIngredientChange"
          >
            <el-option
              v-for="item in ingredientOptions"
              :key="item.id"
              :label="`${item.name} (库存: ${item.stock_quantity}${item.unit})`"
              :value="item.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="损耗数量" prop="amount">
          <el-input-number
            v-model="formData.amount"
            :min="0.01"
            :precision="2"
            :step="1"
            :max="availableStock"
            style="width: 100%"
            @change="calculateCost"
          />
          <div v-if="selectedIngredient" class="stock-info">
            当前库存: {{ selectedIngredient.stock_quantity }} {{ selectedIngredient.unit }}
          </div>
        </el-form-item>

        <el-form-item label="损耗原因" prop="reason">
          <el-select v-model="formData.reason" placeholder="请选择原因" style="width: 100%">
            <el-option
              v-for="reason in reasonOptions"
              :key="reason.value"
              :label="reason.label"
              :value="reason.value"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="操作人" prop="operator">
          <el-input v-model="formData.operator" placeholder="请输入操作人姓名" />
        </el-form-item>

        <el-form-item label="备注">
          <el-input
            v-model="formData.remark"
            type="textarea"
            :rows="3"
            placeholder="请输入备注信息"
          />
        </el-form-item>

        <el-form-item label="预估成本">
          <div class="cost-display">
            <span class="cost-label">¥{{ calculatedCost.toFixed(2) }}</span>
            <span class="cost-hint">
              (单价: ¥{{ selectedIngredient?.cost_price?.toFixed(2) }} × 数量: {{ formData.amount }})
            </span>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button class="glow-button" @click="handleSubmit" :disabled="!formData.amount || formData.amount <= 0">
          确认添加
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Search, Refresh, Plus, Delete, Document, Money, List } from '@element-plus/icons-vue'
import { api, type WasteRecord, type WasteRecordForm, type Spirit, type Ingredient } from '@/api'

const reasonOptions = [
  { label: '腐烂变质', value: '腐烂变质' },
  { label: '操作失误', value: '操作失误' },
  { label: '过期', value: '过期' },
  { label: '保存不当', value: '保存不当' },
  { label: '其他', value: '其他' }
]

const loading = ref(false)
const wasteRecords = ref<WasteRecord[]>([])
const spirits = ref<Spirit[]>([])
const ingredients = ref<Ingredient[]>([])

const dateRange = ref<[string, string] | null>(null)
const filterType = ref<string>('')

const pagination = reactive({
  currentPage: 1,
  pageSize: 10
})

const dialogVisible = ref(false)
const formRef = ref<FormInstance>()

const defaultFormData: WasteRecordForm = {
  ingredient_type: 'spirit',
  ingredient_id: 0,
  amount: 0,
  reason: '',
  operator: '',
  remark: ''
}

const formData = reactive<WasteRecordForm>({ ...defaultFormData })

const formRules: FormRules = {
  ingredient_type: [{ required: true, message: '请选择配料类型', trigger: 'change' }],
  ingredient_id: [{ required: true, message: '请选择原料', trigger: 'change' }],
  amount: [
    { required: true, message: '请输入损耗数量', trigger: 'blur' },
    { type: 'number', min: 0.01, message: '数量必须大于0', trigger: 'blur' }
  ],
  reason: [{ required: true, message: '请选择损耗原因', trigger: 'change' }],
  operator: [{ required: true, message: '请输入操作人', trigger: 'blur' }]
}

const filteredData = computed(() => {
  let result = [...wasteRecords.value]

  if (filterType.value) {
    result = result.filter(item => item.ingredient_type === filterType.value)
  }

  return result
})

const total = computed(() => filteredData.value.length)

const totalCount = computed(() => wasteRecords.value.length)

const totalCost = computed(() => {
  return wasteRecords.value.reduce((sum, item) => sum + (item.cost || 0), 0)
})

const paginatedData = computed(() => {
  const start = (pagination.currentPage - 1) * pagination.pageSize
  const end = start + pagination.pageSize
  return filteredData.value.slice(start, end)
})

const ingredientOptions = computed(() => {
  if (formData.ingredient_type === 'spirit') {
    return spirits.value.map(s => ({
      id: s.id,
      name: s.name,
      stock_quantity: s.stock_quantity,
      unit: s.unit,
      cost_price: s.cost_price
    }))
  } else {
    return ingredients.value.map(i => ({
      id: i.id,
      name: i.name,
      stock_quantity: i.stock_quantity,
      unit: i.unit,
      cost_price: i.cost_price
    }))
  }
})

const selectedIngredient = computed(() => {
  if (!formData.ingredient_id) return null
  return ingredientOptions.value.find(item => item.id === formData.ingredient_id)
})

const availableStock = computed(() => {
  return selectedIngredient.value?.stock_quantity || 0
})

const calculatedCost = computed(() => {
  if (!selectedIngredient.value || !formData.amount) return 0
  return selectedIngredient.value.cost_price * formData.amount
})

function getReasonTagType(reason: string): string {
  const typeMap: Record<string, string> = {
    '腐烂变质': 'danger',
    '操作失误': 'warning',
    '过期': 'info',
    '保存不当': 'warning',
    '其他': ''
  }
  return typeMap[reason] || 'info'
}

async function fetchWasteRecords() {
  loading.value = true
  try {
    const params: Record<string, any> = {}
    if (dateRange.value && dateRange.value[0]) {
      params.start_date = dateRange.value[0]
      params.end_date = dateRange.value[1]
    }
    if (filterType.value) {
      params.ingredient_type = filterType.value
    }

    const response = await api.getWasteRecords(params)
    if (response.data.code === 200) {
      wasteRecords.value = response.data.data || []
    }
  } catch (error) {
    ElMessage.error('获取损耗记录失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

async function fetchSpirits() {
  try {
    const response = await api.getSpirits()
    if (response.data.code === 200) {
      spirits.value = response.data.data || []
    }
  } catch (error) {
    console.error('获取基酒列表失败', error)
  }
}

async function fetchIngredients() {
  try {
    const response = await api.getIngredients()
    if (response.data.code === 200) {
      ingredients.value = response.data.data || []
    }
  } catch (error) {
    console.error('获取食材列表失败', error)
  }
}

function handleFilterChange() {
  pagination.currentPage = 1
}

function handleSearch() {
  fetchWasteRecords()
}

function handleReset() {
  dateRange.value = null
  filterType.value = ''
  pagination.currentPage = 1
  fetchWasteRecords()
}

function handleAdd() {
  Object.assign(formData, defaultFormData)
  dialogVisible.value = true
}

function handleTypeChange() {
  formData.ingredient_id = 0
  formData.amount = 0
}

function handleIngredientChange() {
  formData.amount = 0
}



async function handleDelete(row: WasteRecord) {
  try {
    await ElMessageBox.confirm(
      `确定要删除「${row.ingredient_name}」的损耗记录吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    const response = await api.deleteWasteRecord(row.id)
    if (response.data.code === 200) {
      ElMessage.success('删除成功')
      fetchWasteRecords()
    } else {
      ElMessage.error(response.data.message || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
    }
  }
}

async function handleSubmit() {
  if (!formRef.value) return
  
  if (formData.amount > availableStock.value) {
    ElMessage.error(`损耗数量不能超过当前库存 (${availableStock.value} ${selectedIngredient.value?.unit})`)
    return
  }
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    try {
      const response = await api.createWasteRecord(formData)
      if (response.data.code === 200) {
        ElMessage.success('损耗记录创建成功，库存已自动扣减')
        dialogVisible.value = false
        fetchWasteRecords()
        fetchSpirits()
        fetchIngredients()
      } else {
        ElMessage.error(response.data.message || '创建失败')
      }
    } catch (error) {
      ElMessage.error('操作失败')
      console.error(error)
    }
  })
}

onMounted(() => {
  fetchWasteRecords()
  fetchSpirits()
  fetchIngredients()
})
</script>

<style lang="scss" scoped>
@import '@/styles/variables.scss';

.waste-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 20px;

  .header-title {
    display: flex;
    align-items: center;
    gap: 12px;

    .title {
      margin: 0;
      font-size: 24px;
      font-weight: 700;
    }
  }

  .header-stats {
    display: flex;
    gap: 20px;

    .stat-card {
      display: flex;
      align-items: center;
      gap: 12px;
      padding: 16px 24px;
      background: rgba(255, 255, 255, 0.05);
      border: 1px solid rgba(212, 175, 55, 0.3);
      border-radius: $border-radius-lg;
      min-width: 180px;

      .stat-icon {
        font-size: 32px;
        color: $primary-gold;
      }

      .stat-content {
        display: flex;
        flex-direction: column;
        gap: 4px;

        .stat-label {
          font-size: 12px;
          color: $text-secondary;
        }

        .stat-value {
          font-size: 24px;
          font-weight: 700;
          color: $primary-gold;
        }
      }

      &.danger {
        border-color: rgba(231, 76, 60, 0.3);

        .stat-icon {
          color: $danger;
        }

        .stat-value {
          color: $danger;
        }

        .cost-value {
          color: $danger !important;
        }
      }
    }
  }
}

.filter-card {
  .filter-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex-wrap: wrap;
    gap: 16px;
  }

  .filter-left {
    display: flex;
    align-items: center;
    gap: 12px;
    flex-wrap: wrap;
  }

  .filter-date {
    width: 280px;
  }

  .filter-select {
    width: 160px;
  }

  .reset-button {
    background: rgba(255, 255, 255, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.2);
    color: $text-secondary;

    &:hover {
      border-color: $primary-gold;
      color: $primary-gold;
    }
  }
}

.table-card {
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    padding-bottom: 16px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);

    .table-title {
      display: flex;
      align-items: center;
      gap: 10px;
      font-size: 18px;
      font-weight: 600;
      color: $primary-gold;
      margin: 0;

      .el-icon {
        color: $primary-gold;
      }
    }
  }

  .name-text {
    font-weight: 500;
  }

  .cost-text {
    font-weight: 600;
    color: $danger;
  }

  .pagination-wrapper {
    display: flex;
    justify-content: flex-end;
    margin-top: 20px;
  }
}

.waste-form {
  padding: 10px 0;

  .stock-info {
    margin-top: 8px;
    font-size: 12px;
    color: $text-secondary;
  }

  .cost-display {
    display: flex;
    flex-direction: column;
    gap: 4px;

    .cost-label {
      font-size: 28px;
      font-weight: 700;
      color: $danger;
    }

    .cost-hint {
      font-size: 12px;
      color: $text-secondary;
    }
  }
}

.gold-text {
  color: $primary-gold;
  text-shadow: 0 0 10px rgba(212, 175, 55, 0.5);
}

.glow-button {
  background: linear-gradient(135deg, $primary-gold, $secondary-gold);
  border: none;
  color: $dark-bg;
  font-weight: 600;
  transition: all $transition-normal;
  box-shadow: $shadow-gold;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 0 30px rgba(212, 175, 55, 0.5);
  }

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    transform: none;
    box-shadow: none;
  }
}

:deep(.el-radio) {
  .el-radio__label {
    color: $text-primary !important;
  }

  .el-radio__input.is-checked .el-radio__inner {
    border-color: $primary-gold;
    background: $primary-gold;
  }
}
</style>
