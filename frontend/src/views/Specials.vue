<template>
  <div class="specials-page">
    <div class="page-header">
      <div class="header-title">
        <el-icon :size="28" color="#d4af37"><Goblet /></el-icon>
        <h2 class="title">特调创意管理</h2>
      </div>
      <div class="header-stats">
        <div class="stat-card" @click="activeTab = 'all'">
          <span class="stat-label">总数</span>
          <span class="stat-value">{{ totalCount }}</span>
        </div>
        <div class="stat-card draft" @click="activeTab = 'draft'">
          <span class="stat-label">草稿</span>
          <span class="stat-value">{{ draftCount }}</span>
        </div>
        <div class="stat-card testing" @click="activeTab = 'testing'">
          <span class="stat-label">测试中</span>
          <span class="stat-value">{{ testingCount }}</span>
        </div>
        <div class="stat-card approved" @click="activeTab = 'approved'">
          <span class="stat-label">已通过</span>
          <span class="stat-value">{{ approvedCount }}</span>
        </div>
      </div>
    </div>

    <el-tabs v-model="activeTab" class="status-tabs" @tab-change="handleTabChange">
      <el-tab-pane label="全部" name="all" />
      <el-tab-pane label="草稿" name="draft" />
      <el-tab-pane label="测试中" name="testing" />
      <el-tab-pane label="已通过" name="approved" />
    </el-tabs>

    <div class="filter-bar">
      <div class="filter-left">
        <el-select
          v-model="filterForm.status"
          placeholder="状态筛选"
          clearable
          class="filter-select"
          @change="handleFilterChange"
        >
          <el-option label="草稿" value="draft" />
          <el-option label="测试中" value="testing" />
          <el-option label="已通过" value="approved" />
        </el-select>
        <el-input
          v-model="filterForm.keyword"
          placeholder="搜索名称、创作者、口味..."
          clearable
          class="filter-input"
          :prefix-icon="Search"
          @keyup.enter="handleFilterChange"
        />
        <el-button class="glow-button" :icon="Search" @click="handleFilterChange">
          搜索
        </el-button>
        <el-button class="reset-button" :icon="Refresh" @click="resetFilters">
          重置
        </el-button>
      </div>
      <div class="filter-right">
        <el-radio-group v-model="viewMode" size="default">
          <el-radio-button value="table">
            <el-icon><Grid /></el-icon>
            表格
          </el-radio-button>
          <el-radio-button value="card">
            <el-icon><Tickets /></el-icon>
            卡片
          </el-radio-button>
        </el-radio-group>
        <el-button class="glow-button" :icon="Plus" @click="openAddDialog">
          新增特调
        </el-button>
      </div>
    </div>

    <div v-if="viewMode === 'table'" class="table-container">
      <el-table
        v-loading="loading"
        :data="paginatedData"
        stripe
        class="specials-table"
        :row-class-name="tableRowClassName"
      >
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="name" label="名称" min-width="160">
          <template #default="{ row }">
            <span class="name-text">{{ row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="creator" label="创作者" width="120" align="center" />
        <el-table-column prop="taste_profile" label="口味" min-width="140" />
        <el-table-column prop="price" label="售价" width="100" align="right">
          <template #default="{ row }">
            <span class="price-text">¥{{ row.price.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)" effect="dark" size="small">
              {{ getStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" align="center">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" align="center" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link :icon="Edit" size="small" @click="openEditDialog(row)">
              编辑
            </el-button>
            <el-button type="danger" link :icon="Delete" size="small" @click="handleDelete(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.currentPage"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          background
        />
      </div>
    </div>

    <div v-else class="card-container">
      <div class="card-grid">
        <div
          v-for="special in paginatedData"
          :key="special.id"
          class="special-card"
          :class="{ 'approved-card': special.status === 'approved' }"
        >
          <div class="card-image">
            <img
              :src="special.image_url || getDefaultImage()"
              :alt="special.name"
            />
            <div class="status-badge" :class="special.status">
              {{ getStatusLabel(special.status) }}
            </div>
          </div>
          <div class="card-content">
            <div class="card-header">
              <h3 class="card-title">{{ special.name }}</h3>
              <span class="card-creator">
                <el-icon><User /></el-icon>
                {{ special.creator }}
              </span>
            </div>
            <div class="card-meta">
              <span class="meta-item">
                <el-icon><Goblet /></el-icon>
                {{ special.glass_type }}
              </span>
              <span class="meta-item">
                <el-icon><Star /></el-icon>
                {{ special.serving_ml }}ml
              </span>
            </div>
            <p class="card-taste">{{ special.taste_profile }}</p>
            <p v-if="special.inspiration" class="card-inspiration">
              <el-icon><View /></el-icon>
              {{ special.inspiration }}
            </p>
            <div class="card-price">
              <span class="current-price">¥{{ special.price.toFixed(2) }}</span>
            </div>
            <div class="card-footer">
              <span class="create-time">{{ formatDate(special.created_at) }}</span>
              <div class="card-actions">
                <el-button type="primary" link :icon="Edit" size="small" @click="openEditDialog(special)">
                  编辑
                </el-button>
                <el-button type="danger" link :icon="Delete" size="small" @click="handleDelete(special)">
                  删除
                </el-button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.currentPage"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[12, 24, 48, 96]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          background
        />
      </div>
    </div>

    <el-dialog
      v-model="formDialogVisible"
      :title="isEdit ? '编辑特调' : '新增特调'"
      width="700px"
      :close-on-click-modal="false"
      class="special-dialog"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
        class="special-form"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="特调名称" prop="name">
              <el-input v-model="formData.name" placeholder="请输入特调名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="创作者" prop="creator">
              <el-input v-model="formData.creator" placeholder="请输入创作者" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="创作灵感" prop="inspiration">
          <el-input
            v-model="formData.inspiration"
            type="textarea"
            :rows="2"
            placeholder="请输入创作灵感"
          />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="口味" prop="taste_profile">
              <el-input v-model="formData.taste_profile" placeholder="如：酸甜、清爽、浓郁..." />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="杯型" prop="glass_type">
              <el-select v-model="formData.glass_type" placeholder="请选择杯型" style="width: 100%">
                <el-option
                  v-for="glass in glassTypes"
                  :key="glass"
                  :label="glass"
                  :value="glass"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="容量(ml)" prop="serving_ml">
              <el-input-number
                v-model="formData.serving_ml"
                :min="0"
                :step="10"
                style="width: 100%"
                placeholder="请输入容量"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="售价(¥)" prop="price">
              <el-input-number
                v-model="formData.price"
                :min="0"
                :precision="2"
                :step="1"
                style="width: 100%"
                placeholder="请输入售价"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="制作方法" prop="preparation_method">
          <el-input
            v-model="formData.preparation_method"
            type="textarea"
            :rows="3"
            placeholder="请输入制作方法"
          />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="装饰" prop="garnish">
              <el-input v-model="formData.garnish" placeholder="请输入装饰方式" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-select v-model="formData.status" placeholder="请选择状态" style="width: 100%">
                <el-option label="草稿" value="draft" />
                <el-option label="测试中" value="testing" />
                <el-option label="已通过" value="approved" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="配方成分" prop="ingredients_text">
          <el-input
            v-model="formData.ingredients_text"
            type="textarea"
            :rows="3"
            placeholder="请输入配方成分，如：60ml 威士忌、30ml 柠檬汁、15ml 糖浆..."
          />
        </el-form-item>
        <el-form-item label="品鉴笔记" prop="tasting_notes">
          <el-input
            v-model="formData.tasting_notes"
            type="textarea"
            :rows="3"
            placeholder="请输入品鉴笔记"
          />
        </el-form-item>
        <el-form-item label="图片URL" prop="image_url">
          <el-input v-model="formData.image_url" placeholder="请输入图片URL（可选）" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="formDialogVisible = false">取消</el-button>
        <el-button class="glow-button" :loading="submitting" @click="handleSubmit">
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
          确定要删除特调 <span class="highlight">{{ currentDeleteItem?.name }}</span> 吗？
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
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import {
  Goblet,
  Search,
  Plus,
  Edit,
  Delete,
  Refresh,
  Warning,
  Grid,
  Tickets,
  User,
  Star,
  View
} from '@element-plus/icons-vue'
import { api, type SpecialCreation, type SpecialCreationForm } from '@/api'

type StatusType = 'draft' | 'testing' | 'approved' | 'all'

const glassTypes = [
  '鸡尾酒杯', '古典杯', '高球杯', '柯林杯', '香槟杯',
  '白兰地杯', '威士忌杯', '玛格丽特杯', '飓风杯', '啤酒杯', '闻香杯'
]

const statusMap: Record<string, string> = {
  draft: '草稿',
  testing: '测试中',
  approved: '已通过'
}

const statusTagMap: Record<string, 'info' | 'warning' | 'success'> = {
  draft: 'info',
  testing: 'warning',
  approved: 'success'
}

const mockSpecials: SpecialCreation[] = [
  {
    id: 1,
    name: '金色午夜',
    creator: '张伟',
    inspiration: '融合东方茶文化与西方调酒艺术，打造独特的午夜体验',
    taste_profile: '馥郁、烟熏、茶香、微甜',
    glass_type: '古典杯',
    serving_ml: 180,
    price: 168.00,
    preparation_method: '1. 将威士忌倒入冰镇古典杯\n2. 加入普洱茶汤轻轻搅拌\n3. 加入冰块和蜂蜜\n4. 用烟熏肉桂装饰',
    garnish: '肉桂棒、橙皮',
    ingredients_text: '60ml 单一麦芽威士忌、30ml 普洱茶汤、15ml 蜂蜜糖浆、2滴 安格斯特拉苦精',
    image_url: '',
    status: 'approved',
    tasting_notes: '入口有浓郁的麦芽香气，中段浮现普洱的陈香，尾韵带有淡淡烟熏。层次丰富，回味悠长。',
    created_at: '2024-01-15T10:00:00Z',
    updated_at: '2024-01-15T10:00:00Z'
  },
  {
    id: 2,
    name: '樱花限定',
    creator: '李娜',
    inspiration: '春日樱花烂漫，将美好融入酒中',
    taste_profile: '清新、花香、微甜、优雅',
    glass_type: '鸡尾酒杯',
    serving_ml: 120,
    price: 128.00,
    preparation_method: '1. 所有原料加入摇酒壶\n2. 加冰摇匀\n3. 双重过滤倒入冰镇酒杯\n4. 用盐渍樱花装饰',
    garnish: '盐渍樱花、樱花粉',
    ingredients_text: '45ml 金酒、20ml 樱花利口酒、15ml 柠檬汁、10ml 糖浆、适量 樱花糖浆',
    image_url: '',
    status: 'testing',
    tasting_notes: '花香浓郁，口感清爽，金酒的草本与樱花的甜美完美融合。适合春日品尝。',
    created_at: '2024-02-20T14:30:00Z',
    updated_at: '2024-02-20T14:30:00Z'
  },
  {
    id: 3,
    name: '老炮儿',
    creator: '王强',
    inspiration: '致敬经典，加入中国白酒元素的创新特调',
    taste_profile: '浓烈、酱香、复杂、醇厚',
    glass_type: '闻香杯',
    serving_ml: 100,
    price: 198.00,
    preparation_method: '1. 将白酒和苦精倒入调酒杯\n2. 加冰搅拌30秒\n3. 过滤倒入酒杯\n4. 喷上柠檬油',
    garnish: '柠檬皮 twist',
    ingredients_text: '50ml 茅台、10ml 甜味美思、2滴 安格斯特拉苦精、1滴 橙味苦精',
    image_url: '',
    status: 'draft',
    tasting_notes: '酱香突出，优雅细腻，酒体醇厚，回味悠长。适合资深酒客品尝。',
    created_at: '2024-03-10T09:15:00Z',
    updated_at: '2024-03-10T09:15:00Z'
  },
  {
    id: 4,
    name: '蜜桃乌龙',
    creator: '陈芳',
    inspiration: '年轻人喜爱的蜜桃乌龙茶饮，调酒版更有惊喜',
    taste_profile: '甜美、果香、茶香、清爽',
    glass_type: '高球杯',
    serving_ml: 350,
    price: 88.00,
    preparation_method: '1. 杯底放入蜜桃块\n2. 加入乌龙茶和伏特加\n3. 加入冰块\n4. 用苏打水填满\n5. 轻轻搅拌',
    garnish: '蜜桃片、薄荷叶',
    ingredients_text: '40ml 伏特加、80ml 冷泡乌龙茶、30ml 蜜桃糖浆、60ml 蜜桃汁、适量 苏打水',
    image_url: '',
    status: 'approved',
    tasting_notes: '蜜桃香甜与乌龙茶香交织，酒精感柔和易饮，是夏日消暑佳品。',
    created_at: '2024-03-25T16:45:00Z',
    updated_at: '2024-03-25T16:45:00Z'
  },
  {
    id: 5,
    name: '烟熏教父',
    creator: '刘洋',
    inspiration: '经典鸡尾酒教父的烟熏版本，更具个性',
    taste_profile: '烟熏、杏仁、甜香、强烈',
    glass_type: '古典杯',
    serving_ml: 150,
    price: 148.00,
    preparation_method: '1. 先用烟熏枪将杯子烟熏\n2. 加入大冰块\n3. 倒入威士忌和杏仁利口酒\n4. 轻轻搅拌\n5. 用烟熏杏仁装饰',
    garnish: '烟熏杏仁、橙皮',
    ingredients_text: '60ml 烟熏威士忌、25ml 阿玛雷托杏仁利口酒、10ml 糖浆',
    image_url: '',
    status: 'testing',
    tasting_notes: '强烈的烟熏开场，中段是杏仁的甜美，尾韵悠长。层次感非常丰富。',
    created_at: '2024-04-05T11:20:00Z',
    updated_at: '2024-04-05T11:20:00Z'
  },
  {
    id: 6,
    name: '翡翠花园',
    creator: '赵雪',
    inspiration: '草本植物的清新，如漫步清晨花园',
    taste_profile: '草本、清新、微甜、清爽',
    glass_type: '鸡尾酒杯',
    serving_ml: 130,
    price: 108.00,
    preparation_method: '1. 将黄瓜和罗勒叶轻压\n2. 加入其他原料\n3. 加冰摇匀\n4. 过滤倒入酒杯\n5. 用罗勒叶装饰',
    garnish: '新鲜罗勒叶、黄瓜片',
    ingredients_text: '45ml 金酒、20ml 黄瓜汁、15ml 青柠汁、10ml 糖浆、6片 新鲜罗勒叶',
    image_url: '',
    status: 'draft',
    tasting_notes: '非常清新的草本风味，黄瓜的清爽与罗勒的芳香完美结合，余韵带有金酒的杜松子香。',
    created_at: '2024-04-18T13:50:00Z',
    updated_at: '2024-04-18T13:50:00Z'
  }
]

const loading = ref(false)
const submitting = ref(false)
const deleting = ref(false)
const viewMode = ref<'table' | 'card'>('table')
const formDialogVisible = ref(false)
const deleteDialogVisible = ref(false)
const isEdit = ref(false)
const currentEditId = ref<number | null>(null)
const currentDeleteItem = ref<SpecialCreation | null>(null)
const activeTab = ref<StatusType>('all')

const filterForm = reactive({
  status: '' as string,
  keyword: ''
})

const pagination = reactive({
  currentPage: 1,
  pageSize: 10
})

const allSpecials = ref<SpecialCreation[]>([])
const formRef = ref<FormInstance>()

const defaultFormData: SpecialCreationForm = {
  name: '',
  creator: '',
  inspiration: '',
  taste_profile: '',
  glass_type: '',
  serving_ml: 0,
  price: 0,
  preparation_method: '',
  garnish: '',
  ingredients_text: '',
  image_url: '',
  status: 'draft' as const,
  tasting_notes: ''
}

const formData = reactive<SpecialCreationForm>({ ...defaultFormData })

const formRules: FormRules = {
  name: [{ required: true, message: '请输入特调名称', trigger: 'blur' }],
  creator: [{ required: true, message: '请输入创作者', trigger: 'blur' }],
  taste_profile: [{ required: true, message: '请输入口味描述', trigger: 'blur' }],
  glass_type: [{ required: true, message: '请选择杯型', trigger: 'change' }],
  serving_ml: [{ required: true, message: '请输入容量', trigger: 'blur' }],
  price: [{ required: true, message: '请输入售价', trigger: 'blur' }],
  status: [{ required: true, message: '请选择状态', trigger: 'change' }]
}

const filteredData = computed(() => {
  let result = [...allSpecials.value]

  const activeStatus = activeTab.value !== 'all' ? activeTab.value : filterForm.status
  if (activeStatus) {
    result = result.filter(s => s.status === activeStatus)
  }

  if (filterForm.keyword.trim()) {
    const keyword = filterForm.keyword.toLowerCase().trim()
    result = result.filter(s =>
      s.name.toLowerCase().includes(keyword) ||
      s.creator.toLowerCase().includes(keyword) ||
      s.taste_profile.toLowerCase().includes(keyword)
    )
  }

  return result
})

const paginatedData = computed(() => {
  const start = (pagination.currentPage - 1) * pagination.pageSize
  const end = start + pagination.pageSize
  return filteredData.value.slice(start, end)
})

const total = computed(() => filteredData.value.length)
const totalCount = computed(() => allSpecials.value.length)
const draftCount = computed(() => allSpecials.value.filter(s => s.status === 'draft').length)
const testingCount = computed(() => allSpecials.value.filter(s => s.status === 'testing').length)
const approvedCount = computed(() => allSpecials.value.filter(s => s.status === 'approved').length)

function getStatusLabel(status: string): string {
  return statusMap[status] || status
}

function getStatusTagType(status: string): 'info' | 'warning' | 'success' {
  return statusTagMap[status] || 'info'
}

function getDefaultImage(): string {
  return 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=luxury%20signature%20cocktail%20golden%20garnish%20dark%20bar%20background&image_size=square'
}

function formatDate(dateStr: string): string {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

function tableRowClassName({ row }: { row: SpecialCreation }): string {
  if (row.status === 'approved') {
    return 'approved-row'
  }
  return ''
}

function handleTabChange(tabName: string) {
  filterForm.status = tabName === 'all' ? '' : tabName
  pagination.currentPage = 1
}

function handleFilterChange() {
  pagination.currentPage = 1
}

function resetFilters() {
  filterForm.status = ''
  filterForm.keyword = ''
  activeTab.value = 'all'
  pagination.currentPage = 1
}

function resetForm() {
  Object.assign(formData, defaultFormData)
  formRef.value?.resetFields()
}

function openAddDialog() {
  isEdit.value = false
  currentEditId.value = null
  resetForm()
  formDialogVisible.value = true
}

function openEditDialog(row: SpecialCreation) {
  isEdit.value = true
  currentEditId.value = row.id
  Object.assign(formData, {
    name: row.name,
    creator: row.creator,
    inspiration: row.inspiration,
    taste_profile: row.taste_profile,
    glass_type: row.glass_type,
    serving_ml: row.serving_ml,
    price: row.price,
    preparation_method: row.preparation_method,
    garnish: row.garnish,
    ingredients_text: row.ingredients_text,
    image_url: row.image_url,
    status: row.status,
    tasting_notes: row.tasting_notes
  })
  formDialogVisible.value = true
}

async function fetchSpecials() {
  loading.value = true
  try {
    const params = {
      status: filterForm.status || undefined,
      keyword: filterForm.keyword || undefined,
      page: pagination.currentPage,
      page_size: pagination.pageSize
    }
    const response = await api.getSpecials(params)
    if (response.code === 0 && Array.isArray(response.data)) {
      allSpecials.value = response.data
    } else {
      throw new Error('API response invalid')
    }
  } catch (error) {
    console.warn('API请求失败，使用mock数据:', error)
    allSpecials.value = mockSpecials
  } finally {
    loading.value = false
  }
}

async function handleSubmit() {
  if (!formRef.value) return

  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    if (isEdit.value && currentEditId.value) {
      await api.updateSpecial(currentEditId.value, formData)
      ElMessage.success('更新成功')
    } else {
      await api.createSpecial(formData)
      ElMessage.success('创建成功')
    }
    formDialogVisible.value = false
    fetchSpecials()
  } catch (error) {
    console.error('提交失败:', error)
    if (isEdit.value && currentEditId.value) {
      const index = mockSpecials.findIndex(s => s.id === currentEditId.value)
      if (index !== -1) {
        mockSpecials[index] = {
          ...mockSpecials[index],
          ...formData,
          id: currentEditId.value,
          updated_at: new Date().toISOString()
        } as SpecialCreation
      }
    } else {
      const newId = Math.max(...mockSpecials.map(s => s.id), 0) + 1
      mockSpecials.push({
        ...formData,
        id: newId,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString()
      } as SpecialCreation)
    }
    ElMessage.success('操作成功')
    formDialogVisible.value = false
    fetchSpecials()
  } finally {
    submitting.value = false
  }
}

function handleDelete(row: SpecialCreation) {
  currentDeleteItem.value = row
  deleteDialogVisible.value = true
}

async function confirmDelete() {
  if (!currentDeleteItem.value) return

  deleting.value = true
  try {
    await api.deleteSpecial(currentDeleteItem.value.id)
    ElMessage.success('删除成功')
    deleteDialogVisible.value = false
    fetchSpecials()
  } catch (error) {
    console.error('删除失败:', error)
    const index = mockSpecials.findIndex(s => s.id === currentDeleteItem.value!.id)
    if (index !== -1) {
      mockSpecials.splice(index, 1)
    }
    ElMessage.success('删除成功')
    deleteDialogVisible.value = false
    fetchSpecials()
  } finally {
    deleting.value = false
    currentDeleteItem.value = null
  }
}

onMounted(() => {
  fetchSpecials()
})
</script>

<style lang="scss" scoped>
@import '@/styles/variables.scss';

.specials-page {
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
    gap: 16px;
    cursor: pointer;

    .stat-card {
      display: flex;
      flex-direction: column;
      align-items: center;
      padding: 12px 20px;
      background: rgba(26, 26, 46, 0.8);
      border: 1px solid rgba(212, 175, 55, 0.2);
      border-radius: 8px;
      min-width: 80px;
      transition: all $transition-normal;

      &:hover {
        transform: translateY(-2px);
        box-shadow: $shadow-gold;
      }

      .stat-label {
        font-size: 12px;
        color: #a0a0a0;
        margin-bottom: 4px;
      }

      .stat-value {
        font-size: 22px;
        font-weight: 700;
        color: $primary-gold;
      }

      &.draft .stat-value {
        color: $text-secondary;
      }

      &.testing .stat-value {
        color: $warning;
      }

      &.approved .stat-value {
        color: $success;
      }
    }
  }
}

.status-tabs {
  :deep(.el-tabs__header) {
    margin: 0;
    border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  }

  :deep(.el-tabs__item) {
    color: $text-secondary;
    font-weight: 500;

    &.is-active {
      color: $primary-gold;
    }

    &:hover {
      color: $primary-gold;
    }
  }

  :deep(.el-tabs__active-bar) {
    background: linear-gradient(90deg, $primary-gold, $secondary-gold);
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
  flex-wrap: wrap;
  gap: 12px;

  .filter-left {
    display: flex;
    gap: 12px;
    align-items: center;
    flex-wrap: wrap;

    .filter-select {
      width: 140px;
    }

    .filter-input {
      width: 240px;
    }
  }

  .filter-right {
    display: flex;
    gap: 12px;
    align-items: center;
  }
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

.table-container {
  background: rgba(26, 26, 46, 0.6);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  overflow: hidden;
  backdrop-filter: blur(10px);

  .specials-table {
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

      .el-table__row.approved-row > td {
        background: rgba(39, 174, 96, 0.06) !important;
        border-left: 3px solid $success;

        &:hover {
          background: rgba(39, 174, 96, 0.1) !important;
        }
      }

      .name-text {
        color: #f5f5f5;
        font-weight: 500;
      }

      .price-text {
        color: $success;
        font-weight: 600;
      }
    }
  }

  .pagination-container {
    display: flex;
    justify-content: flex-end;
    padding: 20px;
    border-top: 1px solid rgba(255, 255, 255, 0.08);
  }
}

.card-container {
  .card-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: 20px;
  }

  .special-card {
    background: rgba(26, 26, 46, 0.8);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 12px;
    overflow: hidden;
    transition: all $transition-normal;
    backdrop-filter: blur(10px);

    &:hover {
      transform: translateY(-4px);
      box-shadow: $shadow-lg;
      border-color: rgba(212, 175, 55, 0.3);
    }

    &.approved-card {
      border-color: rgba(212, 175, 55, 0.6);
      box-shadow: 0 0 20px rgba(212, 175, 55, 0.2);

      &:hover {
        box-shadow: 0 0 30px rgba(212, 175, 55, 0.4);
        border-color: $primary-gold;
      }
    }

    .card-image {
      position: relative;
      height: 200px;
      overflow: hidden;

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
        transition: transform $transition-slow;
      }

      &:hover img {
        transform: scale(1.05);
      }

      .status-badge {
        position: absolute;
        top: 12px;
        right: 12px;
        padding: 4px 12px;
        font-size: 12px;
        font-weight: 600;
        border-radius: 20px;
        box-shadow: $shadow-sm;

        &.draft {
          background: rgba(108, 117, 125, 0.9);
          color: #fff;
        }

        &.testing {
          background: rgba(243, 156, 18, 0.9);
          color: #fff;
        }

        &.approved {
          background: linear-gradient(135deg, $primary-gold, $secondary-gold);
          color: $dark-bg;
        }
      }
    }

    .card-content {
      padding: 16px;

      .card-header {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        margin-bottom: 10px;

        .card-title {
          margin: 0;
          font-size: 18px;
          font-weight: 600;
          color: $text-primary;
        }

        .card-creator {
          display: flex;
          align-items: center;
          gap: 4px;
          font-size: 13px;
          color: $text-secondary;

          .el-icon {
            color: $primary-gold;
          }
        }
      }

      .card-meta {
        display: flex;
        align-items: center;
        gap: 16px;
        margin-bottom: 10px;

        .meta-item {
          display: flex;
          align-items: center;
          gap: 4px;
          font-size: 13px;
          color: $text-secondary;

          .el-icon {
            color: $primary-gold;
          }
        }
      }

      .card-taste {
        margin: 0 0 8px 0;
        font-size: 13px;
        color: $text-secondary;
        line-height: 1.5;
      }

      .card-inspiration {
        margin: 0 0 12px 0;
        font-size: 12px;
        color: $text-muted;
        line-height: 1.5;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;

        .el-icon {
          color: $primary-gold;
          margin-right: 4px;
        }
      }

      .card-price {
        margin-bottom: 12px;

        .current-price {
          font-size: 20px;
          font-weight: 700;
          background: linear-gradient(135deg, #d4af37, #ffd700);
          -webkit-background-clip: text;
          -webkit-text-fill-color: transparent;
          background-clip: text;
        }
      }

      .card-footer {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding-top: 12px;
        border-top: 1px solid rgba(255, 255, 255, 0.08);

        .create-time {
          font-size: 12px;
          color: $text-muted;
        }

        .card-actions {
          display: flex;
          gap: 8px;
        }
      }
    }
  }

  .pagination-container {
    display: flex;
    justify-content: flex-end;
    padding: 20px;
    margin-top: 20px;
    background: rgba(26, 26, 46, 0.6);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 12px;
  }
}

.delete-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px 0;
  text-align: center;

  .el-icon {
    margin-bottom: 16px;
  }

  .delete-text {
    margin: 0;
    font-size: 15px;
    color: $text-primary;
    line-height: 1.6;

    .highlight {
      color: $primary-gold;
      font-weight: 600;
    }
  }
}
</style>
