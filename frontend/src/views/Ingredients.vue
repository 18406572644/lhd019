<template>
  <div class="ingredients-page">
    <el-card class="search-card">
      <div class="search-form">
        <div class="search-left">
          <el-select
            v-model="searchForm.category"
            placeholder="选择分类"
            clearable
            class="search-select"
            @change="handleSearch"
          >
            <el-option
              v-for="cat in categories"
              :key="cat.value"
              :label="cat.label"
              :value="cat.value"
            />
          </el-select>
          <el-input
            v-model="searchForm.keyword"
            placeholder="搜索食材名称、供应商..."
            clearable
            class="search-input"
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-button class="glow-button" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button class="reset-button" @click="handleReset">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
        </div>
        <div class="search-right">
          <el-button class="glow-button" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新增食材
          </el-button>
        </div>
      </div>
    </el-card>

    <el-card class="table-card">
      <div class="table-header">
        <h3 class="table-title">
          <el-icon :size="20"><Food /></el-icon>
          食材库存列表
        </h3>
        <div class="table-stats">
          <span class="stat-item">
            <span class="stat-label">总食材数</span>
            <span class="stat-value gold-text">{{ total }}</span>
          </span>
          <span class="stat-item low-stock">
            <span class="stat-label">低库存</span>
            <span class="stat-value">{{ lowStockCount }}</span>
          </span>
        </div>
      </div>

      <el-table
        :data="paginatedData"
        v-loading="loading"
        stripe
        style="width: 100%"
        :row-class-name="tableRowClassName"
      >
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="name" label="食材名称" min-width="150">
          <template #default="{ row }">
            <span class="name-text">{{ row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="category" label="分类" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="getCategoryTagType(row.category)" effect="dark">
              {{ getCategoryLabel(row.category) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="stock_quantity" label="库存数量" width="120" align="center">
          <template #default="{ row }">
            <span :class="{ 'low-stock-text': isLowStock(row) }">
              {{ row.stock_quantity }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="unit" label="单位" width="100" align="center" />
        <el-table-column prop="min_stock" label="最低库存" width="100" align="center" />
        <el-table-column prop="cost_price" label="成本价" width="120" align="center">
          <template #default="{ row }">
            <span class="price-text">¥{{ row.cost_price?.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="supplier" label="供应商" min-width="150" />
        <el-table-column label="操作" width="180" align="center" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              link
              size="small"
              @click="handleEdit(row)"
            >
              <el-icon><Edit /></el-icon>
              编辑
            </el-button>
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
      :title="isEdit ? '编辑食材' : '新增食材'"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
        class="ingredient-form"
      >
        <el-form-item label="食材名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入食材名称" />
        </el-form-item>
        <el-form-item label="分类" prop="category">
          <el-select v-model="formData.category" placeholder="请选择分类" style="width: 100%">
            <el-option
              v-for="cat in categories"
              :key="cat.value"
              :label="cat.label"
              :value="cat.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="库存数量" prop="stock_quantity">
          <el-input-number
            v-model="formData.stock_quantity"
            :min="0"
            :precision="2"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="单位" prop="unit">
          <el-select v-model="formData.unit" placeholder="请选择单位" style="width: 100%">
            <el-option label="个" value="个" />
            <el-option label="瓶" value="瓶" />
            <el-option label="ml" value="ml" />
            <el-option label="g" value="g" />
            <el-option label="kg" value="kg" />
            <el-option label="杯" value="杯" />
            <el-option label="盒" value="盒" />
            <el-option label="罐" value="罐" />
          </el-select>
        </el-form-item>
        <el-form-item label="最低库存" prop="min_stock">
          <el-input-number
            v-model="formData.min_stock"
            :min="0"
            :precision="2"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="成本价" prop="cost_price">
          <el-input-number
            v-model="formData.cost_price"
            :min="0"
            :precision="2"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="供应商" prop="supplier">
          <el-input v-model="formData.supplier" placeholder="请输入供应商名称" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input
            v-model="formData.remark"
            type="textarea"
            :rows="3"
            placeholder="请输入备注信息"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button class="glow-button" @click="handleSubmit">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Search, Refresh, Plus, Edit, Delete, Food } from '@element-plus/icons-vue'
import { api } from '@/api'

interface Ingredient {
  id: number
  name: string
  category: string
  stock_quantity: number
  unit: string
  min_stock: number
  cost_price: number
  supplier: string
  remark?: string
  created_at?: string
  updated_at?: string
}

const categories = [
  { label: '水果', value: '水果' },
  { label: '果汁', value: '果汁' },
  { label: '糖浆', value: '糖浆' },
  { label: '辅料', value: '辅料' },
  { label: '汽水', value: '汽水' },
  { label: '香草', value: '香草' },
  { label: '调料', value: '调料' }
]

const loading = ref(false)
const ingredients = ref<Ingredient[]>([])

const searchForm = reactive({
  category: '',
  keyword: ''
})

const pagination = reactive({
  currentPage: 1,
  pageSize: 10
})

const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInstance>()

const defaultFormData: Partial<Ingredient> = {
  name: '',
  category: '',
  stock_quantity: 0,
  unit: '',
  min_stock: 0,
  cost_price: 0,
  supplier: '',
  remark: ''
}

const formData = reactive<Partial<Ingredient>>({ ...defaultFormData })

const formRules: FormRules = {
  name: [{ required: true, message: '请输入食材名称', trigger: 'blur' }],
  category: [{ required: true, message: '请选择分类', trigger: 'change' }],
  stock_quantity: [{ required: true, message: '请输入库存数量', trigger: 'blur' }],
  unit: [{ required: true, message: '请选择单位', trigger: 'change' }],
  min_stock: [{ required: true, message: '请输入最低库存', trigger: 'blur' }],
  cost_price: [{ required: true, message: '请输入成本价', trigger: 'blur' }],
  supplier: [{ required: true, message: '请输入供应商', trigger: 'blur' }]
}

const filteredData = computed(() => {
  let result = [...ingredients.value]
  
  if (searchForm.category) {
    result = result.filter(item => item.category === searchForm.category)
  }
  
  if (searchForm.keyword) {
    const keyword = searchForm.keyword.toLowerCase()
    result = result.filter(item =>
      item.name.toLowerCase().includes(keyword) ||
      item.supplier.toLowerCase().includes(keyword)
    )
  }
  
  return result
})

const total = computed(() => filteredData.value.length)

const lowStockCount = computed(() => {
  return ingredients.value.filter(item => isLowStock(item)).length
})

const paginatedData = computed(() => {
  const start = (pagination.currentPage - 1) * pagination.pageSize
  const end = start + pagination.pageSize
  return filteredData.value.slice(start, end)
})

function isLowStock(item: Ingredient): boolean {
  return item.stock_quantity <= item.min_stock
}

function getCategoryLabel(value: string): string {
  const cat = categories.find(c => c.value === value)
  return cat ? cat.label : value
}

function getCategoryTagType(category: string): string {
  const typeMap: Record<string, string> = {
    '水果': 'success',
    '果汁': 'warning',
    '糖浆': 'danger',
    '辅料': 'info',
    '汽水': 'primary',
    '香草': '',
    '调料': 'warning'
  }
  return typeMap[category] || 'info'
}

function tableRowClassName({ row }: { row: Ingredient }): string {
  if (isLowStock(row)) {
    return 'low-stock-row'
  }
  return ''
}

async function fetchIngredients() {
  loading.value = true
  try {
    const response = await api.getIngredients()
    if (response.data.code === 200) {
      ingredients.value = response.data.data || []
    }
  } catch (error) {
    ElMessage.error('获取食材列表失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  pagination.currentPage = 1
}

function handleReset() {
  searchForm.category = ''
  searchForm.keyword = ''
  pagination.currentPage = 1
}

function handleAdd() {
  isEdit.value = false
  Object.assign(formData, defaultFormData)
  dialogVisible.value = true
}

function handleEdit(row: Ingredient) {
  isEdit.value = true
  Object.assign(formData, row)
  dialogVisible.value = true
}

async function handleDelete(row: Ingredient) {
  try {
    await ElMessageBox.confirm(
      `确定要删除食材「${row.name}」吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    const response = await api.deleteIngredient(row.id)
    if (response.data.code === 200) {
      ElMessage.success('删除成功')
      fetchIngredients()
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
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    try {
      if (isEdit.value && formData.id) {
        const response = await api.updateIngredient(formData.id, formData)
        if (response.data.code === 200) {
          ElMessage.success('更新成功')
          dialogVisible.value = false
          fetchIngredients()
        } else {
          ElMessage.error(response.data.message || '更新失败')
        }
      } else {
        const response = await api.createIngredient(formData)
        if (response.data.code === 200) {
          ElMessage.success('创建成功')
          dialogVisible.value = false
          fetchIngredients()
        } else {
          ElMessage.error(response.data.message || '创建失败')
        }
      }
    } catch (error) {
      ElMessage.error('操作失败')
      console.error(error)
    }
  })
}

onMounted(() => {
  fetchIngredients()
})
</script>

<style lang="scss" scoped>
@import '@/styles/variables.scss';

.ingredients-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.search-card {
  .search-form {
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex-wrap: wrap;
    gap: 16px;
  }

  .search-left {
    display: flex;
    align-items: center;
    gap: 12px;
    flex-wrap: wrap;
  }

  .search-select {
    width: 160px;
  }

  .search-input {
    width: 280px;
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

    .table-stats {
      display: flex;
      gap: 24px;

      .stat-item {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 4px;

        .stat-label {
          font-size: 12px;
          color: $text-secondary;
        }

        .stat-value {
          font-size: 20px;
          font-weight: 700;
        }

        &.low-stock .stat-value {
          color: $danger;
        }
      }
    }
  }

  .name-text {
    font-weight: 500;
  }

  .price-text {
    font-weight: 600;
    color: $success;
  }

  .low-stock-text {
    color: $danger !important;
    font-weight: 600;
  }

  :deep(.low-stock-row) {
    td {
      background-color: rgba(231, 76, 60, 0.1) !important;
    }

    &:hover > td {
      background-color: rgba(231, 76, 60, 0.15) !important;
    }
  }

  .pagination-wrapper {
    display: flex;
    justify-content: flex-end;
    margin-top: 20px;
  }
}

.ingredient-form {
  padding: 10px 0;
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
}
</style>
