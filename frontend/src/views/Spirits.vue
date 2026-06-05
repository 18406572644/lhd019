<template>
  <div class="spirits-page">
    <div class="page-header">
      <div class="header-title">
        <el-icon :size="28" color="#d4af37"><Goblet /></el-icon>
        <h2 class="title">基酒库存管理</h2>
      </div>
      <div class="header-stats">
        <div class="stat-card">
          <span class="stat-label">总品类</span>
          <span class="stat-value">{{ total }}</span>
        </div>
        <div class="stat-card warning">
          <span class="stat-label">低库存</span>
          <span class="stat-value">{{ lowStockCount }}</span>
        </div>
      </div>
    </div>

    <div class="filter-bar">
      <div class="filter-left">
        <el-select
          v-model="filterCategory"
          placeholder="选择分类"
          clearable
          class="filter-select"
          @change="fetchSpirits"
        >
          <el-option
            v-for="cat in categories"
            :key="cat"
            :label="cat"
            :value="cat"
          />
        </el-select>
        <el-input
          v-model="searchKeyword"
          placeholder="搜索名称/品牌..."
          clearable
          class="filter-input"
          :prefix-icon="Search"
          @keyup.enter="fetchSpirits"
        />
        <el-button type="primary" :icon="Search" @click="fetchSpirits">
          搜索
        </el-button>
        <el-button :icon="Refresh" @click="resetFilters">
          重置
        </el-button>
      </div>
      <div class="filter-right">
        <el-button type="primary" :icon="Plus" @click="openAddDialog">
          新增基酒
        </el-button>
      </div>
    </div>

    <div class="table-container">
      <el-table
        v-loading="loading"
        :data="displayedSpirits"
        stripe
        class="spirits-table"
        :row-class-name="tableRowClassName"
      >
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="name" label="名称" min-width="140">
          <template #default="{ row }">
            <div class="cell-name">
              <span class="name-text">{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="category" label="分类" width="100">
          <template #default="{ row }">
            <el-tag :type="getCategoryTagType(row.category)" effect="dark">
              {{ row.category }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="brand" label="品牌" min-width="120" />
        <el-table-column prop="stock_quantity" label="库存数量" width="110" align="center">
          <template #default="{ row }">
            <span :class="{ 'low-stock-text': isLowStock(row) }">
              {{ row.stock_quantity }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="min_stock" label="最低库存" width="100" align="center" />
        <el-table-column prop="cost_price" label="成本价" width="110" align="right">
          <template #default="{ row }">
            ¥{{ row.cost_price.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="unit" label="单位" width="70" align="center" />
        <el-table-column prop="volume_ml" label="容量(ml)" width="100" align="center" />
        <el-table-column label="操作" width="160" align="center" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link :icon="Edit" @click="openEditDialog(row)">
              编辑
            </el-button>
            <el-button type="danger" link :icon="Delete" @click="handleDelete(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          background
        />
      </div>
    </div>

    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑基酒' : '新增基酒'"
      width="600px"
      :close-on-click-modal="false"
      class="spirit-dialog"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
        class="spirit-form"
      >
        <el-form-item label="名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入基酒名称" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="分类" prop="category">
              <el-select v-model="formData.category" placeholder="请选择分类" class="form-input">
                <el-option
                  v-for="cat in categories"
                  :key="cat"
                  :label="cat"
                  :value="cat"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="品牌" prop="brand">
              <el-input v-model="formData.brand" placeholder="请输入品牌" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="产地" prop="origin">
              <el-input v-model="formData.origin" placeholder="请输入产地" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="酒精度(%)" prop="alcohol_content">
              <el-input-number
                v-model="formData.alcohol_content"
                :min="0"
                :max="100"
                :step="0.1"
                :precision="1"
                class="form-input"
                placeholder="酒精度"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="容量(ml)" prop="volume_ml">
              <el-input-number
                v-model="formData.volume_ml"
                :min="0"
                class="form-input"
                placeholder="容量"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="单位" prop="unit">
              <el-select v-model="formData.unit" placeholder="单位" class="form-input">
                <el-option label="瓶" value="瓶" />
                <el-option label="罐" value="罐" />
                <el-option label="箱" value="箱" />
                <el-option label="L" value="L" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="成本价(¥)" prop="cost_price">
              <el-input-number
                v-model="formData.cost_price"
                :min="0"
                :precision="2"
                :step="1"
                class="form-input"
                placeholder="成本价"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="库存数量" prop="stock_quantity">
              <el-input-number
                v-model="formData.stock_quantity"
                :min="0"
                class="form-input"
                placeholder="库存数量"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="最低库存" prop="min_stock">
              <el-input-number
                v-model="formData.min_stock"
                :min="0"
                class="form-input"
                placeholder="最低库存预警值"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="每ml售价" prop="selling_price_per_ml">
          <el-input-number
            v-model="formData.selling_price_per_ml"
            :min="0"
            :precision="4"
            :step="0.001"
            class="form-input"
            placeholder="每毫升售价"
          />
        </el-form-item>
        <el-form-item label="备注" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="3"
            placeholder="请输入备注信息"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          确定
        </el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="deleteDialogVisible"
      title="确认删除"
      width="420px"
      class="delete-dialog"
    >
      <div class="delete-content">
        <el-icon :size="48" color="#e74c3c"><Warning /></el-icon>
        <p class="delete-text">
          确定要删除 <span class="highlight">{{ currentDeleteItem?.name }}</span> 吗？
          <br />此操作不可恢复。
        </p>
      </div>
      <template #footer>
        <el-button @click="deleteDialogVisible = false">取消</el-button>
        <el-button type="danger" :loading="deleting" @click="confirmDelete">
          确认删除
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import {
  Goblet,
  Search,
  Plus,
  Edit,
  Delete,
  Refresh,
  Warning
} from '@element-plus/icons-vue'
import { api, type Spirit, type SpiritForm } from '@/api'

const categories = ['威士忌', '白兰地', '伏特加', '朗姆', '金酒', '龙舌兰', '利口酒', '其他']

const mockSpirits: Spirit[] = [
  {
    id: 1,
    name: '尊尼获加黑牌',
    category: '威士忌',
    brand: 'Johnnie Walker',
    origin: '苏格兰',
    alcohol_content: 40.0,
    volume_ml: 700,
    unit: '瓶',
    stock_quantity: 12,
    min_stock: 5,
    cost_price: 280.00,
    selling_price_per_ml: 0.6,
    description: '苏格兰调和威士忌',
    created_at: '2024-01-15T10:00:00Z',
    updated_at: '2024-01-15T10:00:00Z'
  },
  {
    id: 2,
    name: '轩尼诗VSOP',
    category: '白兰地',
    brand: 'Hennessy',
    origin: '法国',
    alcohol_content: 40.0,
    volume_ml: 700,
    unit: '瓶',
    stock_quantity: 3,
    min_stock: 5,
    cost_price: 480.00,
    selling_price_per_ml: 1.0,
    description: '法国干邑白兰地',
    created_at: '2024-01-16T10:00:00Z',
    updated_at: '2024-01-16T10:00:00Z'
  },
  {
    id: 3,
    name: '绝对伏特加',
    category: '伏特加',
    brand: 'Absolut',
    origin: '瑞典',
    alcohol_content: 40.0,
    volume_ml: 700,
    unit: '瓶',
    stock_quantity: 8,
    min_stock: 5,
    cost_price: 120.00,
    selling_price_per_ml: 0.3,
    description: '瑞典伏特加',
    created_at: '2024-01-17T10:00:00Z',
    updated_at: '2024-01-17T10:00:00Z'
  },
  {
    id: 4,
    name: '百加得白朗姆',
    category: '朗姆',
    brand: 'Bacardi',
    origin: '波多黎各',
    alcohol_content: 37.5,
    volume_ml: 700,
    unit: '瓶',
    stock_quantity: 15,
    min_stock: 5,
    cost_price: 95.00,
    selling_price_per_ml: 0.25,
    description: '白朗姆酒',
    created_at: '2024-01-18T10:00:00Z',
    updated_at: '2024-01-18T10:00:00Z'
  },
  {
    id: 5,
    name: '添加利伦敦干金酒',
    category: '金酒',
    brand: 'Tanqueray',
    origin: '英国',
    alcohol_content: 43.1,
    volume_ml: 700,
    unit: '瓶',
    stock_quantity: 2,
    min_stock: 5,
    cost_price: 180.00,
    selling_price_per_ml: 0.45,
    description: '伦敦干金酒',
    created_at: '2024-01-19T10:00:00Z',
    updated_at: '2024-01-19T10:00:00Z'
  },
  {
    id: 6,
    name: '奥美加银标龙舌兰',
    category: '龙舌兰',
    brand: 'Olmeca',
    origin: '墨西哥',
    alcohol_content: 38.0,
    volume_ml: 700,
    unit: '瓶',
    stock_quantity: 6,
    min_stock: 5,
    cost_price: 150.00,
    selling_price_per_ml: 0.4,
    description: '银标龙舌兰',
    created_at: '2024-01-20T10:00:00Z',
    updated_at: '2024-01-20T10:00:00Z'
  },
  {
    id: 7,
    name: '君度橙味力娇酒',
    category: '利口酒',
    brand: 'Cointreau',
    origin: '法国',
    alcohol_content: 40.0,
    volume_ml: 700,
    unit: '瓶',
    stock_quantity: 4,
    min_stock: 3,
    cost_price: 210.00,
    selling_price_per_ml: 0.5,
    description: '橙味利口酒',
    created_at: '2024-01-21T10:00:00Z',
    updated_at: '2024-01-21T10:00:00Z'
  },
  {
    id: 8,
    name: '麦卡伦12年',
    category: '威士忌',
    brand: 'Macallan',
    origin: '苏格兰',
    alcohol_content: 40.0,
    volume_ml: 700,
    unit: '瓶',
    stock_quantity: 1,
    min_stock: 3,
    cost_price: 880.00,
    selling_price_per_ml: 1.8,
    description: '单一麦芽威士忌',
    created_at: '2024-01-22T10:00:00Z',
    updated_at: '2024-01-22T10:00:00Z'
  },
  {
    id: 9,
    name: '灰雁伏特加',
    category: '伏特加',
    brand: 'Grey Goose',
    origin: '法国',
    alcohol_content: 40.0,
    volume_ml: 750,
    unit: '瓶',
    stock_quantity: 10,
    min_stock: 5,
    cost_price: 320.00,
    selling_price_per_ml: 0.7,
    description: '法国高端伏特加',
    created_at: '2024-01-23T10:00:00Z',
    updated_at: '2024-01-23T10:00:00Z'
  },
  {
    id: 10,
    name: '百利甜酒',
    category: '利口酒',
    brand: 'Baileys',
    origin: '爱尔兰',
    alcohol_content: 17.0,
    volume_ml: 700,
    unit: '瓶',
    stock_quantity: 7,
    min_stock: 5,
    cost_price: 110.00,
    selling_price_per_ml: 0.3,
    description: '奶油利口酒',
    created_at: '2024-01-24T10:00:00Z',
    updated_at: '2024-01-24T10:00:00Z'
  },
  {
    id: 11,
    name: '人头马VSOP',
    category: '白兰地',
    brand: 'Remy Martin',
    origin: '法国',
    alcohol_content: 40.0,
    volume_ml: 700,
    unit: '瓶',
    stock_quantity: 5,
    min_stock: 5,
    cost_price: 420.00,
    selling_price_per_ml: 0.9,
    description: '优质香槟区干邑',
    created_at: '2024-01-25T10:00:00Z',
    updated_at: '2024-01-25T10:00:00Z'
  },
  {
    id: 12,
    name: '哈瓦那俱乐部7年',
    category: '朗姆',
    brand: 'Havana Club',
    origin: '古巴',
    alcohol_content: 40.0,
    volume_ml: 700,
    unit: '瓶',
    stock_quantity: 0,
    min_stock: 3,
    cost_price: 180.00,
    selling_price_per_ml: 0.45,
    description: '古巴陈酿朗姆酒',
    created_at: '2024-01-26T10:00:00Z',
    updated_at: '2024-01-26T10:00:00Z'
  }
]

const loading = ref(false)
const submitting = ref(false)
const deleting = ref(false)
const dialogVisible = ref(false)
const deleteDialogVisible = ref(false)
const isEdit = ref(false)
const currentEditId = ref<number | null>(null)
const currentDeleteItem = ref<Spirit | null>(null)

const filterCategory = ref('')
const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = ref(10)

const allSpirits = ref<Spirit[]>([])
const formRef = ref<FormInstance>()

const defaultFormData: SpiritForm = {
  name: '',
  category: '',
  brand: '',
  origin: '',
  alcohol_content: 40.0,
  volume_ml: 700,
  unit: '瓶',
  stock_quantity: 0,
  min_stock: 5,
  cost_price: 0,
  selling_price_per_ml: 0,
  description: ''
}

const formData = reactive<SpiritForm>({ ...defaultFormData })

const formRules: FormRules = {
  name: [{ required: true, message: '请输入基酒名称', trigger: 'blur' }],
  category: [{ required: true, message: '请选择分类', trigger: 'change' }],
  brand: [{ required: true, message: '请输入品牌', trigger: 'blur' }],
  volume_ml: [{ required: true, message: '请输入容量', trigger: 'blur' }],
  unit: [{ required: true, message: '请选择单位', trigger: 'change' }],
  stock_quantity: [{ required: true, message: '请输入库存数量', trigger: 'blur' }],
  min_stock: [{ required: true, message: '请输入最低库存', trigger: 'blur' }],
  cost_price: [{ required: true, message: '请输入成本价', trigger: 'blur' }]
}

const filteredSpirits = computed(() => {
  let result = [...allSpirits.value]
  
  if (filterCategory.value) {
    result = result.filter(s => s.category === filterCategory.value)
  }
  
  if (searchKeyword.value.trim()) {
    const keyword = searchKeyword.value.toLowerCase().trim()
    result = result.filter(s =>
      s.name.toLowerCase().includes(keyword) ||
      s.brand.toLowerCase().includes(keyword)
    )
  }
  
  return result
})

const displayedSpirits = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredSpirits.value.slice(start, end)
})

const total = computed(() => filteredSpirits.value.length)

const lowStockCount = computed(() =>
  allSpirits.value.filter(s => isLowStock(s)).length
)

const isLowStock = (row: Spirit) => row.stock_quantity <= row.min_stock

const tableRowClassName = ({ row }: { row: Spirit }) => {
  if (isLowStock(row)) {
    return 'low-stock-row'
  }
  return ''
}

const getCategoryTagType = (category: string): 'primary' | 'success' | 'warning' | 'danger' | 'info' => {
  const typeMap: Record<string, 'primary' | 'success' | 'warning' | 'danger' | 'info'> = {
    '威士忌': 'primary',
    '白兰地': 'warning',
    '伏特加': 'info',
    '朗姆': 'success',
    '金酒': 'primary',
    '龙舌兰': 'danger',
    '利口酒': 'success',
    '其他': 'info'
  }
  return typeMap[category] || 'info'
}

const fetchSpirits = async () => {
  loading.value = true
  try {
    const response = await api.getSpirits({
      category: filterCategory.value || undefined,
      keyword: searchKeyword.value || undefined
    })
    if (response.code === 0 && Array.isArray(response.data)) {
      allSpirits.value = response.data
    } else {
      throw new Error('API response invalid')
    }
  } catch (error) {
    console.warn('API请求失败，使用mock数据:', error)
    allSpirits.value = mockSpirits
  } finally {
    loading.value = false
  }
}

const resetFilters = () => {
  filterCategory.value = ''
  searchKeyword.value = ''
  currentPage.value = 1
  fetchSpirits()
}

const resetForm = () => {
  Object.assign(formData, defaultFormData)
  formRef.value?.resetFields()
}

const openAddDialog = () => {
  isEdit.value = false
  currentEditId.value = null
  resetForm()
  dialogVisible.value = true
}

const openEditDialog = (row: Spirit) => {
  isEdit.value = true
  currentEditId.value = row.id
  Object.assign(formData, {
    name: row.name,
    category: row.category,
    brand: row.brand,
    origin: row.origin,
    alcohol_content: row.alcohol_content,
    volume_ml: row.volume_ml,
    unit: row.unit,
    stock_quantity: row.stock_quantity,
    min_stock: row.min_stock,
    cost_price: row.cost_price,
    selling_price_per_ml: row.selling_price_per_ml,
    description: row.description
  })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return
  
  submitting.value = true
  try {
    if (isEdit.value && currentEditId.value) {
      await api.updateSpirit(currentEditId.value, formData)
      ElMessage.success('更新成功')
    } else {
      await api.createSpirit(formData)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchSpirits()
  } catch (error) {
    console.error('提交失败:', error)
    if (isEdit.value && currentEditId.value) {
      const index = mockSpirits.findIndex(s => s.id === currentEditId.value)
      if (index !== -1) {
        mockSpirits[index] = { ...mockSpirits[index], ...formData }
      }
    } else {
      const newId = Math.max(...mockSpirits.map(s => s.id), 0) + 1
      mockSpirits.push({
        ...formData,
        id: newId,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString()
      })
    }
    ElMessage.success('操作成功')
    dialogVisible.value = false
    fetchSpirits()
  } finally {
    submitting.value = false
  }
}

const handleDelete = (row: Spirit) => {
  currentDeleteItem.value = row
  deleteDialogVisible.value = true
}

const confirmDelete = async () => {
  if (!currentDeleteItem.value) return
  
  deleting.value = true
  try {
    await api.deleteSpirit(currentDeleteItem.value.id)
    ElMessage.success('删除成功')
    deleteDialogVisible.value = false
    fetchSpirits()
  } catch (error) {
    console.error('删除失败:', error)
    const index = mockSpirits.findIndex(s => s.id === currentDeleteItem.value!.id)
    if (index !== -1) {
      mockSpirits.splice(index, 1)
    }
    ElMessage.success('删除成功')
    deleteDialogVisible.value = false
    fetchSpirits()
  } finally {
    deleting.value = false
    currentDeleteItem.value = null
  }
}

onMounted(() => {
  fetchSpirits()
})
</script>

<style lang="scss" scoped>
.spirits-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  background: linear-gradient(135deg, rgba(212, 175, 55, 0.1), rgba(212, 175, 55, 0.02));
  border: 1px solid rgba(212, 175, 55, 0.2);
  border-radius: 12px;
  backdrop-filter: blur(10px);

  .header-title {
    display: flex;
    align-items: center;
    gap: 12px;

    .title {
      margin: 0;
      font-size: 24px;
      font-weight: 700;
      background: linear-gradient(135deg, #d4af37, #c9a227);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
    }
  }

  .header-stats {
    display: flex;
    gap: 20px;

    .stat-card {
      display: flex;
      flex-direction: column;
      align-items: center;
      padding: 12px 24px;
      background: rgba(26, 26, 46, 0.8);
      border: 1px solid rgba(212, 175, 55, 0.2);
      border-radius: 8px;
      min-width: 100px;

      .stat-label {
        font-size: 13px;
        color: #a0a0a0;
        margin-bottom: 4px;
      }

      .stat-value {
        font-size: 24px;
        font-weight: 700;
        color: #d4af37;
      }

      &.warning .stat-value {
        color: #e94560;
      }
    }
  }
}

.filter-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: rgba(26, 26, 46, 0.6);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 10px;
  backdrop-filter: blur(10px);

  .filter-left {
    display: flex;
    gap: 12px;
    align-items: center;

    .filter-select {
      width: 160px;
    }

    .filter-input {
      width: 240px;
    }
  }
}

.table-container {
  background: rgba(26, 26, 46, 0.6);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  overflow: hidden;
  backdrop-filter: blur(10px);

  .spirits-table {
    :deep(.el-table) {
      background: transparent;

      th.el-table__cell {
        background: rgba(212, 175, 55, 0.08);
        color: #d4af37;
        font-weight: 600;
        border-bottom: 1px solid rgba(212, 175, 55, 0.2);
      }

      td.el-table__cell {
        border-bottom: 1px solid rgba(255, 255, 255, 0.05);
      }

      .el-table__row:hover > td {
        background: rgba(212, 175, 55, 0.05) !important;
      }

      .el-table__row.low-stock-row > td {
        background: rgba(233, 69, 96, 0.1) !important;

        &:hover {
          background: rgba(233, 69, 96, 0.15) !important;
        }
      }

      .low-stock-text {
        color: #e94560;
        font-weight: 600;
      }

      .cell-name {
        display: flex;
        align-items: center;

        .name-text {
          color: #f5f5f5;
          font-weight: 500;
        }
      }
    }
  }

  .pagination-container {
    display: flex;
    justify-content: flex-end;
    padding: 20px;
    border-top: 1px solid rgba(255, 255, 255, 0.08);

    :deep(.el-pagination) {
      .el-pager li {
        background: rgba(255, 255, 255, 0.05);
        color: #a0a0a0;

        &.is-active {
          background: #d4af37;
          color: #1a1a2e;
        }

        &:hover:not(.is-active) {
          color: #d4af37;
        }
      }

      .btn-prev,
      .btn-next {
        background: rgba(255, 255, 255, 0.05);
        color: #a0a0a0;

        &:hover {
          color: #d4af37;
        }
      }

      .el-pagination__total,
      .el-pagination__jump {
        color: #a0a0a0;
      }
    }
  }
}

:deep(.spirit-dialog) {
  .el-dialog {
    background: rgba(26, 26, 46, 0.95);
    border: 1px solid rgba(212, 175, 55, 0.2);
    border-radius: 12px;
    backdrop-filter: blur(20px);

    .el-dialog__header {
      border-bottom: 1px solid rgba(212, 175, 55, 0.15);
      padding: 20px 24px;

      .el-dialog__title {
        color: #d4af37;
        font-size: 18px;
        font-weight: 600;
      }

      .el-dialog__headerbtn .el-dialog__close {
        color: #a0a0a0;

        &:hover {
          color: #d4af37;
        }
      }
    }

    .el-dialog__body {
      padding: 24px;
    }

    .el-dialog__footer {
      border-top: 1px solid rgba(212, 175, 55, 0.15);
      padding: 16px 24px;
    }
  }
}

.spirit-form {
  :deep(.el-form-item__label) {
    color: #a0a0a0;
  }

  :deep(.el-input__wrapper),
  :deep(.el-textarea__inner),
  :deep(.el-select__wrapper),
  :deep(.el-input-number) {
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.1);
    color: #f5f5f5;
    box-shadow: none;

    &:hover,
    &.is-hover {
      border-color: rgba(212, 175, 55, 0.4);
    }

    &.is-focus,
    &.is-focused {
      border-color: #d4af37;
      box-shadow: 0 0 0 2px rgba(212, 175, 55, 0.15);
    }

    .el-input__inner,
    .el-select__placeholder,
    .el-input-number__decrease,
    .el-input-number__increase {
      color: #f5f5f5;
    }

    .el-select__placeholder {
      color: #6c757d;
    }
  }

  :deep(.el-input-number) {
    width: 100%;

    .el-input-number__decrease,
    .el-input-number__increase {
      background: transparent;
      border-color: rgba(255, 255, 255, 0.1);
      color: #a0a0a0;

      &:hover {
        color: #d4af37;
        background: rgba(212, 175, 55, 0.1);
      }

      &.is-disabled {
        color: #6c757d;
      }
    }
  }

  .form-input {
    width: 100%;
  }
}

:deep(.delete-dialog) {
  .el-dialog {
    background: rgba(26, 26, 46, 0.95);
    border: 1px solid rgba(233, 69, 96, 0.2);
    border-radius: 12px;
    backdrop-filter: blur(20px);

    .el-dialog__header {
      border-bottom: 1px solid rgba(233, 69, 96, 0.15);
      padding: 20px 24px;

      .el-dialog__title {
        color: #e94560;
        font-size: 18px;
        font-weight: 600;
      }
    }

    .el-dialog__body {
      padding: 24px;
    }

    .el-dialog__footer {
      border-top: 1px solid rgba(233, 69, 96, 0.15);
      padding: 16px 24px;
    }
  }
}

.delete-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 16px;

  .delete-text {
    margin: 0;
    color: #a0a0a0;
    line-height: 1.8;
    font-size: 15px;

    .highlight {
      color: #d4af37;
      font-weight: 600;
    }
  }
}

:deep(.el-button) {
  border-radius: 6px;
  font-weight: 500;
  transition: all 0.2s ease;

  &.el-button--primary {
    background: linear-gradient(135deg, #d4af37, #c9a227);
    border: none;
    color: #1a1a2e;

    &:hover {
      background: linear-gradient(135deg, #e0bc44, #d4af37);
      box-shadow: 0 4px 12px rgba(212, 175, 55, 0.3);
    }
  }

  &.el-button--danger {
    background: #e94560;
    border: none;

    &:hover {
      background: #ff5e78;
      box-shadow: 0 4px 12px rgba(233, 69, 96, 0.3);
    }
  }

  &.el-button--default {
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.1);
    color: #a0a0a0;

    &:hover {
      border-color: #d4af37;
      color: #d4af37;
    }
  }
}

:deep(.el-tag) {
  border-radius: 4px;
  font-weight: 500;
}
</style>
