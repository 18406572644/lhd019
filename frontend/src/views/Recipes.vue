<template>
  <div class="recipes-page">
    <div class="page-header">
      <div class="header-title">
        <el-icon :size="28" color="#d4af37"><Goblet /></el-icon>
        <h2 class="title">调酒配方管理</h2>
      </div>
      <div class="header-stats">
        <div class="stat-card">
          <span class="stat-label">总配方数</span>
          <span class="stat-value">{{ total }}</span>
        </div>
        <div class="stat-card signature">
          <span class="stat-label">招牌特调</span>
          <span class="stat-value">{{ signatureCount }}</span>
        </div>
      </div>
    </div>

    <div class="filter-bar">
      <div class="filter-left">
        <el-select
          v-model="filterForm.category"
          placeholder="选择分类"
          clearable
          class="filter-select"
          @change="handleFilterChange"
        >
          <el-option
            v-for="cat in categories"
            :key="cat.value"
            :label="cat.label"
            :value="cat.value"
          />
        </el-select>
        <el-select
          v-model="filterForm.is_signature"
          placeholder="招牌筛选"
          clearable
          class="filter-select"
          @change="handleFilterChange"
        >
          <el-option label="仅招牌特调" :value="true" />
          <el-option label="非招牌特调" :value="false" />
        </el-select>
        <el-input
          v-model="filterForm.keyword"
          placeholder="搜索配方名称、口味..."
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
          新增配方
        </el-button>
      </div>
    </div>

    <div v-if="viewMode === 'table'" class="table-container">
      <el-table
        v-loading="loading"
        :data="paginatedData"
        stripe
        class="recipes-table"
        :row-class-name="tableRowClassName"
      >
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="name" label="名称" min-width="160">
          <template #default="{ row }">
            <div class="cell-name">
              <span class="name-text">{{ row.name }}</span>
              <el-tag
                v-if="row.is_signature"
                type="warning"
                effect="dark"
                size="small"
                class="signature-tag"
              >
                <el-icon><Star /></el-icon>
                招牌
              </el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="category" label="分类" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="getCategoryTagType(row.category)" effect="dark" size="small">
              {{ row.category }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="glass_type" label="杯型" width="110" align="center" />
        <el-table-column prop="price" label="售价" width="100" align="right">
          <template #default="{ row }">
            <span class="price-text">¥{{ row.price.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="cost" label="成本" width="100" align="right">
          <template #default="{ row }">
            <span class="cost-text">¥{{ row.cost.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="毛利率" width="100" align="center">
          <template #default="{ row }">
            <span :class="getProfitClass(row)">
              {{ getProfitMargin(row) }}%
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="difficulty" label="难度" width="90" align="center">
          <template #default="{ row }">
            <el-tag :type="getDifficultyTagType(row.difficulty)" effect="dark" size="small">
              {{ row.difficulty }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="taste_profile" label="口味" min-width="120" />
        <el-table-column label="操作" width="200" align="center" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link :icon="View" size="small" @click="openViewDialog(row)">
              详情
            </el-button>
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
          v-for="recipe in paginatedData"
          :key="recipe.id"
          class="recipe-card"
          :class="{ 'signature-card': recipe.is_signature }"
        >
          <div class="card-image">
            <img
              :src="recipe.image_url || getDefaultImage(recipe.category)"
              :alt="recipe.name"
            />
            <div v-if="recipe.is_signature" class="signature-badge">
              <el-icon><Star /></el-icon>
              招牌特调
            </div>
            <div class="card-category">
              <el-tag :type="getCategoryTagType(recipe.category)" effect="dark" size="small">
                {{ recipe.category }}
              </el-tag>
            </div>
          </div>
          <div class="card-content">
            <h3 class="card-title">{{ recipe.name }}</h3>
            <div class="card-meta">
              <span class="meta-item">
                <el-icon><Goblet /></el-icon>
                {{ recipe.glass_type }}
              </span>
              <span class="meta-item">
                <el-tag :type="getDifficultyTagType(recipe.difficulty)" effect="dark" size="small">
                  {{ recipe.difficulty }}
                </el-tag>
              </span>
            </div>
            <p class="card-taste">{{ recipe.taste_profile }}</p>
            <div class="card-price">
              <span class="current-price">¥{{ recipe.price.toFixed(2) }}</span>
              <span class="cost-price">成本: ¥{{ recipe.cost.toFixed(2) }}</span>
            </div>
            <div class="card-actions">
              <el-button type="primary" link :icon="View" size="small" @click="openViewDialog(recipe)">
                详情
              </el-button>
              <el-button type="primary" link :icon="Edit" size="small" @click="openEditDialog(recipe)">
                编辑
              </el-button>
              <el-button type="danger" link :icon="Delete" size="small" @click="handleDelete(recipe)">
                删除
              </el-button>
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
      :title="isEdit ? '编辑配方' : '新增配方'"
      width="800px"
      :close-on-click-modal="false"
      class="recipe-dialog"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
        class="recipe-form"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="配方名称" prop="name">
              <el-input v-model="formData.name" placeholder="请输入配方名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
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
          </el-col>
        </el-row>
        <el-row :gutter="20">
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
          <el-col :span="12">
            <el-form-item label="难度" prop="difficulty">
              <el-select v-model="formData.difficulty" placeholder="请选择难度" style="width: 100%">
                <el-option label="简单" value="简单" />
                <el-option label="中等" value="中等" />
                <el-option label="困难" value="困难" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
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
          <el-col :span="12">
            <el-form-item label="成本(¥)" prop="cost">
              <el-input-number
                v-model="formData.cost"
                :min="0"
                :precision="2"
                :step="0.1"
                style="width: 100%"
                placeholder="自动计算/手动输入"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="口味" prop="taste_profile">
          <el-input v-model="formData.taste_profile" placeholder="请输入口味描述，如：酸甜、清爽、浓郁..." />
        </el-form-item>
        <el-form-item label="招牌特调" prop="is_signature">
          <el-switch
            v-model="formData.is_signature"
            active-text="是"
            inactive-text="否"
            :active-value="true"
            :inactive-value="false"
          />
        </el-form-item>
        <el-form-item label="图片URL" prop="image_url">
          <el-input v-model="formData.image_url" placeholder="请输入图片URL（可选）" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="2"
            placeholder="请输入配方描述"
          />
        </el-form-item>
        <el-form-item label="制作方法" prop="preparation_method">
          <el-input
            v-model="formData.preparation_method"
            type="textarea"
            :rows="3"
            placeholder="请输入制作方法"
          />
        </el-form-item>

        <el-form-item label="配方成分" required>
          <div class="ingredients-list">
            <div
              v-for="(ingredient, index) in formData.ingredients"
              :key="index"
              class="ingredient-row"
            >
              <el-select
                v-model="ingredient.ingredient_type"
                placeholder="类型"
                class="ingredient-type"
                @change="(val: string) => handleIngredientTypeChange(index, val)"
              >
                <el-option label="基酒" value="spirit" />
                <el-option label="食材" value="ingredient" />
              </el-select>
              <el-select
                v-model="ingredient.ingredient_id"
                placeholder="选择成分"
                class="ingredient-select"
                filterable
                @change="(val: number) => handleIngredientChange(index, val)"
              >
                <el-option
                  v-for="item in getIngredientOptions(ingredient.ingredient_type)"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                />
              </el-select>
              <el-input-number
                v-model="ingredient.amount"
                :min="0"
                :precision="2"
                :step="1"
                placeholder="用量"
                class="ingredient-amount"
                @change="calculateTotalCost"
              />
              <el-select
                v-model="ingredient.unit"
                placeholder="单位"
                class="ingredient-unit"
              >
                <el-option label="ml" value="ml" />
                <el-option label="g" value="g" />
                <el-option label="个" value="个" />
                <el-option label="勺" value="勺" />
                <el-option label="滴" value="滴" />
                <el-option label="片" value="片" />
                <el-option label="块" value="块" />
              </el-select>
              <el-button
                type="danger"
                :icon="Delete"
                circle
                size="small"
                @click="removeIngredient(index)"
              />
            </div>
            <el-button type="primary" :icon="Plus" @click="addIngredient">
              添加成分
            </el-button>
          </div>
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
      v-model="viewDialogVisible"
      title="配方详情"
      width="600px"
      class="view-dialog"
    >
      <div v-if="currentViewRecipe" class="recipe-detail">
        <div class="detail-header">
          <img
            :src="currentViewRecipe.image_url || getDefaultImage(currentViewRecipe.category)"
            :alt="currentViewRecipe.name"
            class="detail-image"
          />
          <div class="detail-info">
            <h2 class="detail-title">
              {{ currentViewRecipe.name }}
              <el-tag
                v-if="currentViewRecipe.is_signature"
                type="warning"
                effect="dark"
                size="small"
                class="signature-tag"
              >
                <el-icon><Star /></el-icon>
                招牌
              </el-tag>
            </h2>
            <div class="detail-tags">
              <el-tag :type="getCategoryTagType(currentViewRecipe.category)" effect="dark">
                {{ currentViewRecipe.category }}
              </el-tag>
              <el-tag effect="dark">
                <el-icon><Goblet /></el-icon>
                {{ currentViewRecipe.glass_type }}
              </el-tag>
              <el-tag :type="getDifficultyTagType(currentViewRecipe.difficulty)" effect="dark">
                {{ currentViewRecipe.difficulty }}
              </el-tag>
            </div>
            <p class="detail-taste">口味: {{ currentViewRecipe.taste_profile }}</p>
            <div class="detail-price">
              <span class="price">售价: ¥{{ currentViewRecipe.price.toFixed(2) }}</span>
              <span class="cost">成本: ¥{{ currentViewRecipe.cost.toFixed(2) }}</span>
              <span :class="getProfitClass(currentViewRecipe)">
                毛利: {{ getProfitMargin(currentViewRecipe) }}%
              </span>
            </div>
          </div>
        </div>

        <el-descriptions :column="1" border class="detail-descriptions">
          <el-descriptions-item label="描述">
            {{ currentViewRecipe.description || '暂无描述' }}
          </el-descriptions-item>
          <el-descriptions-item label="制作方法">
            {{ currentViewRecipe.preparation_method || '暂无制作方法' }}
          </el-descriptions-item>
        </el-descriptions>

        <h4 class="ingredients-title">配方成分</h4>
        <div class="ingredients-detail">
          <div
            v-for="(ingredient, index) in currentViewRecipe.ingredients"
            :key="index"
            class="ingredient-detail-item"
          >
            <span class="ingredient-name">
              <el-tag :type="ingredient.ingredient_type === 'spirit' ? 'primary' : 'success'" size="small" effect="dark">
                {{ ingredient.ingredient_type === 'spirit' ? '基酒' : '食材' }}
              </el-tag>
              {{ ingredient.ingredient_name || getIngredientName(ingredient) }}
            </span>
            <span class="ingredient-amount-detail">
              {{ ingredient.amount }} {{ ingredient.unit }}
            </span>
          </div>
        </div>
      </div>
      <template #footer>
        <el-button @click="viewDialogVisible = false">关闭</el-button>
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
          确定要删除配方 <span class="highlight">{{ currentDeleteItem?.name }}</span> 吗？
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
  Warning,
  Star,
  Grid,
  Tickets,
  View
} from '@element-plus/icons-vue'
import { api, type Recipe, type RecipeForm, type RecipeIngredient, type Spirit, type Ingredient } from '@/api'

const categories = [
  { label: '经典鸡尾酒', value: '经典鸡尾酒' },
  { label: '特调鸡尾酒', value: '特调鸡尾酒' },
  { label: '无酒精饮品', value: '无酒精饮品' }
]

const glassTypes = [
  '鸡尾酒杯', '古典杯', '高球杯', '柯林杯', '香槟杯',
  '白兰地杯', '威士忌杯', '玛格丽特杯', '飓风杯', '啤酒杯'
]

const difficultyMap: Record<string, 'primary' | 'success' | 'warning' | 'danger' | 'info'> = {
  '简单': 'success',
  '中等': 'warning',
  '困难': 'danger'
}

const categoryTagMap: Record<string, 'primary' | 'success' | 'warning' | 'danger' | 'info'> = {
  '经典鸡尾酒': 'primary',
  '特调鸡尾酒': 'warning',
  '无酒精饮品': 'success'
}

const mockRecipes: Recipe[] = [
  {
    id: 1,
    name: '莫吉托',
    category: '经典鸡尾酒',
    glass_type: '高球杯',
    price: 68.00,
    cost: 18.50,
    difficulty: '简单',
    is_signature: false,
    taste_profile: '清爽、薄荷、酸甜',
    description: '古巴经典鸡尾酒，以朗姆酒为基底，搭配新鲜薄荷和青柠',
    preparation_method: '1. 将薄荷叶和青柠块放入杯中轻压\n2. 加入糖浆和朗姆酒\n3. 加入碎冰搅拌\n4. 用苏打水填满\n5. 用薄荷枝和青柠片装饰',
    image_url: '',
    ingredients: [
      { ingredient_type: 'spirit', ingredient_id: 4, ingredient_name: '百加得白朗姆', amount: 60, unit: 'ml', cost: 15 },
      { ingredient_type: 'ingredient', ingredient_id: 1, ingredient_name: '青柠', amount: 2, unit: '个', cost: 2 },
      { ingredient_type: 'ingredient', ingredient_id: 2, ingredient_name: '薄荷叶', amount: 10, unit: '片', cost: 1 },
      { ingredient_type: 'ingredient', ingredient_id: 3, ingredient_name: '糖浆', amount: 15, unit: 'ml', cost: 0.5 }
    ],
    created_at: '2024-01-15T10:00:00Z',
    updated_at: '2024-01-15T10:00:00Z'
  },
  {
    id: 2,
    name: '金色年华',
    category: '特调鸡尾酒',
    glass_type: '鸡尾酒杯',
    price: 128.00,
    cost: 45.00,
    difficulty: '中等',
    is_signature: true,
    taste_profile: '馥郁、果香、微甜',
    description: '本店招牌特调，融合多种稀有威士忌与特制糖浆',
    preparation_method: '1. 将所有原料加入摇酒壶\n2. 加入冰块充分摇匀\n3. 双重过滤倒入冰镇酒杯\n4. 用橙皮油喷香装饰',
    image_url: '',
    ingredients: [
      { ingredient_type: 'spirit', ingredient_id: 1, ingredient_name: '尊尼获加黑牌', amount: 45, unit: 'ml', cost: 27 },
      { ingredient_type: 'spirit', ingredient_id: 8, ingredient_name: '麦卡伦12年', amount: 15, unit: 'ml', cost: 13.5 },
      { ingredient_type: 'ingredient', ingredient_id: 3, ingredient_name: '糖浆', amount: 10, unit: 'ml', cost: 0.5 },
      { ingredient_type: 'ingredient', ingredient_id: 4, ingredient_name: '安格斯特拉苦精', amount: 2, unit: '滴', cost: 4 }
    ],
    created_at: '2024-01-16T10:00:00Z',
    updated_at: '2024-01-16T10:00:00Z'
  },
  {
    id: 3,
    name: '长岛冰茶',
    category: '经典鸡尾酒',
    glass_type: '高球杯',
    price: 88.00,
    cost: 32.00,
    difficulty: '中等',
    is_signature: false,
    taste_profile: '强烈、清爽、酸甜',
    description: '经典烈酒混合饮品，看似茶实则酒精含量颇高',
    preparation_method: '1. 将所有烈酒和柠檬汁倒入杯中\n2. 加入冰块搅拌\n3. 用可乐填满\n4. 轻轻搅拌\n5. 用柠檬片装饰',
    image_url: '',
    ingredients: [
      { ingredient_type: 'spirit', ingredient_id: 3, ingredient_name: '绝对伏特加', amount: 15, unit: 'ml', cost: 4.5 },
      { ingredient_type: 'spirit', ingredient_id: 5, ingredient_name: '添加利伦敦干金酒', amount: 15, unit: 'ml', cost: 6.75 },
      { ingredient_type: 'spirit', ingredient_id: 4, ingredient_name: '百加得白朗姆', amount: 15, unit: 'ml', cost: 3.75 },
      { ingredient_type: 'spirit', ingredient_id: 6, ingredient_name: '奥美加银标龙舌兰', amount: 15, unit: 'ml', cost: 6 },
      { ingredient_type: 'ingredient', ingredient_id: 1, ingredient_name: '青柠汁', amount: 25, unit: 'ml', cost: 1 },
      { ingredient_type: 'ingredient', ingredient_id: 3, ingredient_name: '糖浆', amount: 10, unit: 'ml', cost: 0.5 }
    ],
    created_at: '2024-01-17T10:00:00Z',
    updated_at: '2024-01-17T10:00:00Z'
  },
  {
    id: 4,
    name: '蜜桃气泡',
    category: '无酒精饮品',
    glass_type: '高球杯',
    price: 48.00,
    cost: 12.00,
    difficulty: '简单',
    is_signature: false,
    taste_profile: '甜美、果香、气泡',
    description: '清爽的无酒精饮品，适合不饮酒的客人',
    preparation_method: '1. 将蜜桃糖浆和柠檬汁倒入杯中\n2. 加入冰块\n3. 用苏打水和气泡水填满\n4. 轻轻搅拌\n5. 用蜜桃片和薄荷叶装饰',
    image_url: '',
    ingredients: [
      { ingredient_type: 'ingredient', ingredient_id: 5, ingredient_name: '蜜桃糖浆', amount: 30, unit: 'ml', cost: 3 },
      { ingredient_type: 'ingredient', ingredient_id: 1, ingredient_name: '青柠汁', amount: 15, unit: 'ml', cost: 1 },
      { ingredient_type: 'ingredient', ingredient_id: 6, ingredient_name: '苏打水', amount: 100, unit: 'ml', cost: 2 },
      { ingredient_type: 'ingredient', ingredient_id: 7, ingredient_name: '新鲜蜜桃', amount: 2, unit: '片', cost: 6 }
    ],
    created_at: '2024-01-18T10:00:00Z',
    updated_at: '2024-01-18T10:00:00Z'
  },
  {
    id: 5,
    name: '干马天尼',
    category: '经典鸡尾酒',
    glass_type: '鸡尾酒杯',
    price: 78.00,
    cost: 22.00,
    difficulty: '简单',
    is_signature: false,
    taste_profile: '干烈、芳香、清爽',
    description: '经典中的经典，以金酒和干味美思调制',
    preparation_method: '1. 将冰镇后的酒杯备好\n2. 金酒和干味美思倒入调酒杯\n3. 加入冰块搅拌约30秒\n4. 过滤倒入酒杯\n5. 用橄榄或柠檬 twist 装饰',
    image_url: '',
    ingredients: [
      { ingredient_type: 'spirit', ingredient_id: 5, ingredient_name: '添加利伦敦干金酒', amount: 60, unit: 'ml', cost: 27 },
      { ingredient_type: 'spirit', ingredient_id: 7, ingredient_name: '君度橙味力娇酒', amount: 10, unit: 'ml', cost: 5 }
    ],
    created_at: '2024-01-19T10:00:00Z',
    updated_at: '2024-01-19T10:00:00Z'
  },
  {
    id: 6,
    name: '午夜巴黎',
    category: '特调鸡尾酒',
    glass_type: '古典杯',
    price: 158.00,
    cost: 58.00,
    difficulty: '困难',
    is_signature: true,
    taste_profile: '深邃、烟熏、复杂',
    description: '主厨特调，使用稀有苏格兰威士忌与特制香料',
    preparation_method: '1. 首先用烟熏枪将杯子烟熏\n2. 将威士忌、利口酒和苦精倒入调酒杯\n3. 加入大冰块搅拌\n4. 倒入烟熏过的杯子\n5. 用橙皮和烟熏肉桂装饰',
    image_url: '',
    ingredients: [
      { ingredient_type: 'spirit', ingredient_id: 8, ingredient_name: '麦卡伦12年', amount: 50, unit: 'ml', cost: 45 },
      { ingredient_type: 'spirit', ingredient_id: 10, ingredient_name: '百利甜酒', amount: 15, unit: 'ml', cost: 4.5 },
      { ingredient_type: 'ingredient', ingredient_id: 4, ingredient_name: '安格斯特拉苦精', amount: 3, unit: '滴', cost: 3.5 },
      { ingredient_type: 'ingredient', ingredient_id: 8, ingredient_name: '肉桂棒', amount: 1, unit: '块', cost: 5 }
    ],
    created_at: '2024-01-20T10:00:00Z',
    updated_at: '2024-01-20T10:00:00Z'
  },
  {
    id: 7,
    name: '椰林飘香',
    category: '经典鸡尾酒',
    glass_type: '飓风杯',
    price: 72.00,
    cost: 25.00,
    difficulty: '简单',
    is_signature: false,
    taste_profile: '甜美、椰香、果香',
    description: '热带风情经典，朗姆酒搭配椰奶和菠萝汁',
    preparation_method: '1. 将所有原料倒入搅拌机\n2. 加入碎冰搅拌至顺滑\n3. 倒入飓风杯\n4. 用菠萝片和樱桃装饰',
    image_url: '',
    ingredients: [
      { ingredient_type: 'spirit', ingredient_id: 4, ingredient_name: '百加得白朗姆', amount: 50, unit: 'ml', cost: 12.5 },
      { ingredient_type: 'ingredient', ingredient_id: 9, ingredient_name: '椰奶', amount: 50, unit: 'ml', cost: 3 },
      { ingredient_type: 'ingredient', ingredient_id: 10, ingredient_name: '菠萝汁', amount: 100, unit: 'ml', cost: 4.5 }
    ],
    created_at: '2024-01-21T10:00:00Z',
    updated_at: '2024-01-21T10:00:00Z'
  },
  {
    id: 8,
    name: '薄荷柠檬苏打',
    category: '无酒精饮品',
    glass_type: '高球杯',
    price: 38.00,
    cost: 8.00,
    difficulty: '简单',
    is_signature: false,
    taste_profile: '清爽、薄荷、酸甜',
    description: '经典无酒精饮品，清新解渴',
    preparation_method: '1. 薄荷叶和青柠块放入杯底轻压\n2. 加入糖浆和青柠汁\n3. 加入冰块\n4. 用苏打水填满\n5. 用薄荷叶和青柠片装饰',
    image_url: '',
    ingredients: [
      { ingredient_type: 'ingredient', ingredient_id: 2, ingredient_name: '薄荷叶', amount: 12, unit: '片', cost: 1.5 },
      { ingredient_type: 'ingredient', ingredient_id: 1, ingredient_name: '青柠', amount: 1.5, unit: '个', cost: 1.5 },
      { ingredient_type: 'ingredient', ingredient_id: 3, ingredient_name: '糖浆', amount: 20, unit: 'ml', cost: 1 },
      { ingredient_type: 'ingredient', ingredient_id: 6, ingredient_name: '苏打水', amount: 150, unit: 'ml', cost: 3 }
    ],
    created_at: '2024-01-22T10:00:00Z',
    updated_at: '2024-01-22T10:00:00Z'
  }
]

const mockSpirits: Spirit[] = [
  { id: 1, name: '尊尼获加黑牌', category: '威士忌', brand: 'Johnnie Walker', origin: '苏格兰', alcohol_content: 40.0, volume_ml: 700, unit: '瓶', stock_quantity: 12, min_stock: 5, cost_price: 280.00, selling_price_per_ml: 0.6, description: '苏格兰调和威士忌', created_at: '2024-01-15T10:00:00Z', updated_at: '2024-01-15T10:00:00Z' },
  { id: 3, name: '绝对伏特加', category: '伏特加', brand: 'Absolut', origin: '瑞典', alcohol_content: 40.0, volume_ml: 700, unit: '瓶', stock_quantity: 8, min_stock: 5, cost_price: 120.00, selling_price_per_ml: 0.3, description: '瑞典伏特加', created_at: '2024-01-17T10:00:00Z', updated_at: '2024-01-17T10:00:00Z' },
  { id: 4, name: '百加得白朗姆', category: '朗姆', brand: 'Bacardi', origin: '波多黎各', alcohol_content: 37.5, volume_ml: 700, unit: '瓶', stock_quantity: 15, min_stock: 5, cost_price: 95.00, selling_price_per_ml: 0.25, description: '白朗姆酒', created_at: '2024-01-18T10:00:00Z', updated_at: '2024-01-18T10:00:00Z' },
  { id: 5, name: '添加利伦敦干金酒', category: '金酒', brand: 'Tanqueray', origin: '英国', alcohol_content: 43.1, volume_ml: 700, unit: '瓶', stock_quantity: 2, min_stock: 5, cost_price: 180.00, selling_price_per_ml: 0.45, description: '伦敦干金酒', created_at: '2024-01-19T10:00:00Z', updated_at: '2024-01-19T10:00:00Z' },
  { id: 6, name: '奥美加银标龙舌兰', category: '龙舌兰', brand: 'Olmeca', origin: '墨西哥', alcohol_content: 38.0, volume_ml: 700, unit: '瓶', stock_quantity: 6, min_stock: 5, cost_price: 150.00, selling_price_per_ml: 0.4, description: '银标龙舌兰', created_at: '2024-01-20T10:00:00Z', updated_at: '2024-01-20T10:00:00Z' },
  { id: 7, name: '君度橙味力娇酒', category: '利口酒', brand: 'Cointreau', origin: '法国', alcohol_content: 40.0, volume_ml: 700, unit: '瓶', stock_quantity: 4, min_stock: 3, cost_price: 210.00, selling_price_per_ml: 0.5, description: '橙味利口酒', created_at: '2024-01-21T10:00:00Z', updated_at: '2024-01-21T10:00:00Z' },
  { id: 8, name: '麦卡伦12年', category: '威士忌', brand: 'Macallan', origin: '苏格兰', alcohol_content: 40.0, volume_ml: 700, unit: '瓶', stock_quantity: 1, min_stock: 3, cost_price: 880.00, selling_price_per_ml: 1.8, description: '单一麦芽威士忌', created_at: '2024-01-22T10:00:00Z', updated_at: '2024-01-22T10:00:00Z' },
  { id: 10, name: '百利甜酒', category: '利口酒', brand: 'Baileys', origin: '爱尔兰', alcohol_content: 17.0, volume_ml: 700, unit: '瓶', stock_quantity: 7, min_stock: 5, cost_price: 110.00, selling_price_per_ml: 0.3, description: '奶油利口酒', created_at: '2024-01-24T10:00:00Z', updated_at: '2024-01-24T10:00:00Z' }
]

const mockIngredients: Ingredient[] = [
  { id: 1, name: '青柠', category: '水果', stock_quantity: 50, unit: '个', min_stock: 20, cost_price: 2.5, supplier: '鲜果供应商', remark: '' },
  { id: 2, name: '薄荷叶', category: '香草', stock_quantity: 100, unit: '片', min_stock: 30, cost_price: 0.2, supplier: '香草农场', remark: '' },
  { id: 3, name: '糖浆', category: '糖浆', stock_quantity: 20, unit: '瓶', min_stock: 5, cost_price: 45, supplier: '调酒原料商', remark: '' },
  { id: 4, name: '安格斯特拉苦精', category: '调料', stock_quantity: 15, unit: '瓶', min_stock: 5, cost_price: 120, supplier: '进口调料商', remark: '' },
  { id: 5, name: '蜜桃糖浆', category: '糖浆', stock_quantity: 10, unit: '瓶', min_stock: 3, cost_price: 65, supplier: '调酒原料商', remark: '' },
  { id: 6, name: '苏打水', category: '汽水', stock_quantity: 48, unit: '罐', min_stock: 24, cost_price: 3.5, supplier: '饮料供应商', remark: '' },
  { id: 7, name: '新鲜蜜桃', category: '水果', stock_quantity: 30, unit: '个', min_stock: 10, cost_price: 8, supplier: '鲜果供应商', remark: '' },
  { id: 8, name: '肉桂棒', category: '香料', stock_quantity: 100, unit: '根', min_stock: 30, cost_price: 1.5, supplier: '香料供应商', remark: '' },
  { id: 9, name: '椰奶', category: '果汁', stock_quantity: 24, unit: '罐', min_stock: 12, cost_price: 12, supplier: '进口食品商', remark: '' },
  { id: 10, name: '菠萝汁', category: '果汁', stock_quantity: 24, unit: '罐', min_stock: 12, cost_price: 15, supplier: '进口食品商', remark: '' }
]

const loading = ref(false)
const submitting = ref(false)
const deleting = ref(false)
const viewMode = ref<'table' | 'card'>('table')
const formDialogVisible = ref(false)
const viewDialogVisible = ref(false)
const deleteDialogVisible = ref(false)
const isEdit = ref(false)
const currentEditId = ref<number | null>(null)
const currentDeleteItem = ref<Recipe | null>(null)
const currentViewRecipe = ref<Recipe | null>(null)

const allSpirits = ref<Spirit[]>([])
const allIngredients = ref<Ingredient[]>([])

const filterForm = reactive({
  category: '',
  keyword: '',
  is_signature: undefined as boolean | undefined
})

const pagination = reactive({
  currentPage: 1,
  pageSize: 10
})

const allRecipes = ref<Recipe[]>([])
const formRef = ref<FormInstance>()

const defaultFormData: RecipeForm = {
  name: '',
  category: '',
  glass_type: '',
  price: 0,
  cost: 0,
  difficulty: '简单',
  is_signature: false,
  taste_profile: '',
  description: '',
  preparation_method: '',
  image_url: '',
  ingredients: []
}

const formData = reactive<RecipeForm>({ ...defaultFormData, ingredients: [] })

const formRules: FormRules = {
  name: [{ required: true, message: '请输入配方名称', trigger: 'blur' }],
  category: [{ required: true, message: '请选择分类', trigger: 'change' }],
  glass_type: [{ required: true, message: '请选择杯型', trigger: 'change' }],
  price: [{ required: true, message: '请输入售价', trigger: 'blur' }],
  cost: [{ required: true, message: '请输入成本', trigger: 'blur' }],
  difficulty: [{ required: true, message: '请选择难度', trigger: 'change' }],
  taste_profile: [{ required: true, message: '请输入口味描述', trigger: 'blur' }]
}

const filteredData = computed(() => {
  let result = [...allRecipes.value]
  
  if (filterForm.category) {
    result = result.filter(r => r.category === filterForm.category)
  }
  
  if (filterForm.is_signature !== undefined) {
    result = result.filter(r => r.is_signature === filterForm.is_signature)
  }
  
  if (filterForm.keyword.trim()) {
    const keyword = filterForm.keyword.toLowerCase().trim()
    result = result.filter(r =>
      r.name.toLowerCase().includes(keyword) ||
      r.taste_profile.toLowerCase().includes(keyword)
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

const signatureCount = computed(() =>
  allRecipes.value.filter(r => r.is_signature).length
)

function getCategoryTagType(category: string): 'primary' | 'success' | 'warning' | 'danger' | 'info' {
  return categoryTagMap[category] || 'info'
}

function getDifficultyTagType(difficulty: string): 'primary' | 'success' | 'warning' | 'danger' | 'info' {
  return difficultyMap[difficulty] || 'info'
}

function getProfitMargin(recipe: Recipe): string {
  if (recipe.price <= 0) return '0.00'
  return (((recipe.price - recipe.cost) / recipe.price) * 100).toFixed(1)
}

function getProfitClass(recipe: Recipe): string {
  const margin = (recipe.price - recipe.cost) / recipe.price
  if (margin >= 0.7) return 'profit-high'
  if (margin >= 0.5) return 'profit-medium'
  return 'profit-low'
}

function getDefaultImage(category: string): string {
  const imageMap: Record<string, string> = {
    '经典鸡尾酒': 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=classic%20cocktail%20in%20elegant%20glass%20dark%20luxury%20bar%20background&image_size=square',
    '特调鸡尾酒': 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=signature%20cocktail%20golden%20garnish%20luxury%20bar&image_size=square',
    '无酒精饮品': 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=non%20alcoholic%20mocktail%20fresh%20fruit%20elegant%20glass&image_size=square'
  }
  return imageMap[category] || imageMap['经典鸡尾酒']
}

function getIngredientOptions(type: string): { id: number; name: string }[] {
  if (type === 'spirit') {
    return allSpirits.value.map(s => ({ id: s.id, name: s.name }))
  } else {
    return allIngredients.value.map(i => ({ id: i.id, name: i.name }))
  }
}

function getIngredientName(ingredient: RecipeIngredient): string {
  if (ingredient.ingredient_name) return ingredient.ingredient_name
  if (ingredient.ingredient_type === 'spirit') {
    const spirit = allSpirits.value.find(s => s.id === ingredient.ingredient_id)
    return spirit?.name || '未知基酒'
  } else {
    const ing = allIngredients.value.find(i => i.id === ingredient.ingredient_id)
    return ing?.name || '未知食材'
  }
}

function tableRowClassName({ row }: { row: Recipe }): string {
  if (row.is_signature) {
    return 'signature-row'
  }
  return ''
}

function addIngredient() {
  formData.ingredients.push({
    ingredient_type: 'spirit',
    ingredient_id: 0,
    amount: 0,
    unit: 'ml'
  })
}

function removeIngredient(index: number) {
  formData.ingredients.splice(index, 1)
  calculateTotalCost()
}

function handleIngredientTypeChange(index: number, type: string) {
  formData.ingredients[index].ingredient_id = 0
  formData.ingredients[index].ingredient_name = undefined
  formData.ingredients[index].cost = undefined
  calculateTotalCost()
}

function handleIngredientChange(index: number, id: number) {
  const ingredient = formData.ingredients[index]
  if (ingredient.ingredient_type === 'spirit') {
    const spirit = allSpirits.value.find(s => s.id === id)
    if (spirit) {
      ingredient.ingredient_name = spirit.name
    }
  } else {
    const ing = allIngredients.value.find(i => i.id === id)
    if (ing) {
      ingredient.ingredient_name = ing.name
    }
  }
  calculateTotalCost()
}

function calculateTotalCost() {
  let total = 0
  for (const ing of formData.ingredients) {
    if (ing.ingredient_type === 'spirit') {
      const spirit = allSpirits.value.find(s => s.id === ing.ingredient_id)
      if (spirit && ing.unit === 'ml') {
        total += spirit.selling_price_per_ml * ing.amount
      }
    } else {
      const ingredient = allIngredients.value.find(i => i.id === ing.ingredient_id)
      if (ingredient) {
        const unitCost = ingredient.cost_price / ingredient.stock_quantity
        total += unitCost * ing.amount
      }
    }
  }
  if (total > 0) {
    formData.cost = parseFloat(total.toFixed(2))
  }
}

function handleFilterChange() {
  pagination.currentPage = 1
  fetchRecipes()
}

function resetFilters() {
  filterForm.category = ''
  filterForm.keyword = ''
  filterForm.is_signature = undefined
  pagination.currentPage = 1
  fetchRecipes()
}

function resetForm() {
  Object.assign(formData, defaultFormData)
  formData.ingredients = []
  formRef.value?.resetFields()
}

function openAddDialog() {
  isEdit.value = false
  currentEditId.value = null
  resetForm()
  addIngredient()
  formDialogVisible.value = true
}

function openEditDialog(row: Recipe) {
  isEdit.value = true
  currentEditId.value = row.id
  Object.assign(formData, {
    name: row.name,
    category: row.category,
    glass_type: row.glass_type,
    price: row.price,
    cost: row.cost,
    difficulty: row.difficulty,
    is_signature: row.is_signature,
    taste_profile: row.taste_profile,
    description: row.description,
    preparation_method: row.preparation_method,
    image_url: row.image_url,
    ingredients: JSON.parse(JSON.stringify(row.ingredients))
  })
  formDialogVisible.value = true
}

function openViewDialog(row: Recipe) {
  currentViewRecipe.value = row
  viewDialogVisible.value = true
}

async function fetchRecipes() {
  loading.value = true
  try {
    const params = {
      category: filterForm.category || undefined,
      keyword: filterForm.keyword || undefined,
      is_signature: filterForm.is_signature,
      page: pagination.currentPage,
      page_size: pagination.pageSize
    }
    const response = await api.getRecipes(params)
    if (response.code === 0 && Array.isArray(response.data)) {
      allRecipes.value = response.data
    } else {
      throw new Error('API response invalid')
    }
  } catch (error) {
    console.warn('API请求失败，使用mock数据:', error)
    allRecipes.value = mockRecipes
  } finally {
    loading.value = false
  }
}

async function fetchSpiritsAndIngredients() {
  try {
    const [spiritsRes, ingredientsRes] = await Promise.all([
      api.getSpirits(),
      api.getIngredients()
    ])
    if (spiritsRes.code === 0 && Array.isArray(spiritsRes.data)) {
      allSpirits.value = spiritsRes.data
    } else {
      allSpirits.value = mockSpirits
    }
    if (ingredientsRes.code === 0 && Array.isArray(ingredientsRes.data)) {
      allIngredients.value = ingredientsRes.data
    } else {
      allIngredients.value = mockIngredients
    }
  } catch (error) {
    console.warn('获取基酒和食材失败，使用mock数据:', error)
    allSpirits.value = mockSpirits
    allIngredients.value = mockIngredients
  }
}

async function handleSubmit() {
  if (!formRef.value) return
  
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return
  
  if (formData.ingredients.length === 0) {
    ElMessage.warning('请至少添加一种成分')
    return
  }
  
  for (let i = 0; i < formData.ingredients.length; i++) {
    const ing = formData.ingredients[i]
    if (!ing.ingredient_id || ing.amount <= 0 || !ing.unit) {
      ElMessage.warning(`请完善第 ${i + 1} 种成分的信息`)
      return
    }
  }
  
  submitting.value = true
  try {
    if (isEdit.value && currentEditId.value) {
      await api.updateRecipe(currentEditId.value, formData)
      ElMessage.success('更新成功')
    } else {
      await api.createRecipe(formData)
      ElMessage.success('创建成功')
    }
    formDialogVisible.value = false
    fetchRecipes()
  } catch (error) {
    console.error('提交失败:', error)
    if (isEdit.value && currentEditId.value) {
      const index = mockRecipes.findIndex(r => r.id === currentEditId.value)
      if (index !== -1) {
        mockRecipes[index] = { ...mockRecipes[index], ...formData, id: currentEditId.value, updated_at: new Date().toISOString() } as Recipe
      }
    } else {
      const newId = Math.max(...mockRecipes.map(r => r.id), 0) + 1
      mockRecipes.push({
        ...formData,
        id: newId,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString()
      } as Recipe)
    }
    ElMessage.success('操作成功')
    formDialogVisible.value = false
    fetchRecipes()
  } finally {
    submitting.value = false
  }
}

function handleDelete(row: Recipe) {
  currentDeleteItem.value = row
  deleteDialogVisible.value = true
}

async function confirmDelete() {
  if (!currentDeleteItem.value) return
  
  deleting.value = true
  try {
    await api.deleteRecipe(currentDeleteItem.value.id)
    ElMessage.success('删除成功')
    deleteDialogVisible.value = false
    fetchRecipes()
  } catch (error) {
    console.error('删除失败:', error)
    const index = mockRecipes.findIndex(r => r.id === currentDeleteItem.value!.id)
    if (index !== -1) {
      mockRecipes.splice(index, 1)
    }
    ElMessage.success('删除成功')
    deleteDialogVisible.value = false
    fetchRecipes()
  } finally {
    deleting.value = false
    currentDeleteItem.value = null
  }
}

onMounted(() => {
  fetchRecipes()
  fetchSpiritsAndIngredients()
})
</script>

<style lang="scss" scoped>
@import '@/styles/variables.scss';

.recipes-page {
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

      &.signature .stat-value {
        background: linear-gradient(135deg, #d4af37, #ffd700);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        background-clip: text;
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

  .recipes-table {
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

      .el-table__row.signature-row > td {
        background: rgba(212, 175, 55, 0.08) !important;

        &:hover {
          background: rgba(212, 175, 55, 0.12) !important;
        }
      }

      .cell-name {
        display: flex;
        align-items: center;
        gap: 8px;

        .name-text {
          color: #f5f5f5;
          font-weight: 500;
        }
      }

      .price-text {
        color: $success;
        font-weight: 600;
      }

      .cost-text {
        color: $text-secondary;
      }

      .profit-high {
        color: $success;
        font-weight: 600;
      }

      .profit-medium {
        color: $warning;
        font-weight: 600;
      }

      .profit-low {
        color: $danger;
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
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 20px;
  }

  .recipe-card {
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

    &.signature-card {
      border-color: rgba(212, 175, 55, 0.5);
      box-shadow: 0 0 20px rgba(212, 175, 55, 0.15);

      &:hover {
        box-shadow: 0 0 30px rgba(212, 175, 55, 0.3);
      }
    }

    .card-image {
      position: relative;
      height: 180px;
      overflow: hidden;

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }

      .signature-badge {
        position: absolute;
        top: 12px;
        right: 12px;
        display: flex;
        align-items: center;
        gap: 4px;
        padding: 4px 10px;
        background: linear-gradient(135deg, #d4af37, #ffd700);
        color: $dark-bg;
        font-size: 12px;
        font-weight: 600;
        border-radius: 20px;
        box-shadow: $shadow-md;
      }

      .card-category {
        position: absolute;
        bottom: 12px;
        left: 12px;
      }
    }

    .card-content {
      padding: 16px;

      .card-title {
        margin: 0 0 10px 0;
        font-size: 18px;
        font-weight: 600;
        color: $text-primary;
      }

      .card-meta {
        display: flex;
        align-items: center;
        gap: 12px;
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
        margin: 0 0 12px 0;
        font-size: 13px;
        color: $text-secondary;
        line-height: 1.5;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
      }

      .card-price {
        display: flex;
        align-items: baseline;
        gap: 10px;
        margin-bottom: 12px;

        .current-price {
          font-size: 20px;
          font-weight: 700;
          color: $success;
        }

        .cost-price {
          font-size: 13px;
          color: $text-secondary;
        }
      }

      .card-actions {
        display: flex;
        gap: 8px;
        padding-top: 12px;
        border-top: 1px solid rgba(255, 255, 255, 0.08);
      }
    }
  }

  .pagination-container {
    display: flex;
    justify-content: flex-end;
    padding: 20px 0 0 0;
  }
}

.ingredients-list {
  display: flex;
  flex-direction: column;
  gap: 12px;

  .ingredient-row {
    display: flex;
    gap: 8px;
    align-items: center;

    .ingredient-type {
      width: 100px;
    }

    .ingredient-select {
      flex: 1;
      min-width: 150px;
    }

    .ingredient-amount {
      width: 120px;
    }

    .ingredient-unit {
      width: 90px;
    }
  }
}

.recipe-detail {
  .detail-header {
    display: flex;
    gap: 20px;
    margin-bottom: 20px;

    .detail-image {
      width: 180px;
      height: 180px;
      object-fit: cover;
      border-radius: 12px;
      border: 2px solid rgba(212, 175, 55, 0.3);
    }

    .detail-info {
      flex: 1;

      .detail-title {
        margin: 0 0 12px 0;
        font-size: 24px;
        font-weight: 700;
        color: $text-primary;
        display: flex;
        align-items: center;
        gap: 10px;
      }

      .detail-tags {
        display: flex;
        gap: 8px;
        margin-bottom: 12px;
        flex-wrap: wrap;

        .el-tag {
          display: flex;
          align-items: center;
          gap: 4px;
        }
      }

      .detail-taste {
        margin: 0 0 12px 0;
        color: $text-secondary;
        font-size: 14px;
      }

      .detail-price {
        display: flex;
        gap: 20px;
        align-items: center;

        .price {
          font-size: 18px;
          font-weight: 700;
          color: $success;
        }

        .cost {
          font-size: 14px;
          color: $text-secondary;
        }
      }
    }
  }

  .detail-descriptions {
    margin-bottom: 20px;
  }

  .ingredients-title {
    margin: 0 0 12px 0;
    font-size: 16px;
    font-weight: 600;
    color: $primary-gold;
  }

  .ingredients-detail {
    display: flex;
    flex-direction: column;
    gap: 8px;

    .ingredient-detail-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 10px 14px;
      background: rgba(255, 255, 255, 0.03);
      border-radius: 8px;
      border: 1px solid rgba(255, 255, 255, 0.05);

      .ingredient-name {
        display: flex;
        align-items: center;
        gap: 8px;
        color: $text-primary;
      }

      .ingredient-amount-detail {
        color: $primary-gold;
        font-weight: 600;
      }
    }
  }
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

:deep(.el-radio-group) {
  .el-radio-button {
    .el-radio-button__inner {
      background: rgba(255, 255, 255, 0.05);
      border-color: rgba(255, 255, 255, 0.1);
      color: $text-secondary;

      &:hover {
        color: $primary-gold;
      }
    }

    &.is-active .el-radio-button__inner {
      background: $primary-gold;
      border-color: $primary-gold;
      color: $dark-bg;
    }
  }
}

:deep(.el-dialog) {
  background: $dark-bg-secondary !important;
  border: 1px solid rgba(212, 175, 55, 0.2) !important;
  border-radius: $border-radius-lg !important;

  .el-dialog__header {
    border-bottom: 1px solid rgba(255, 255, 255, 0.08) !important;
    padding: 20px 24px !important;

    .el-dialog__title {
      color: $primary-gold !important;
      font-weight: 600;
    }
  }

  .el-dialog__body {
    padding: 24px !important;
  }

  .el-dialog__footer {
    border-top: 1px solid rgba(255, 255, 255, 0.08) !important;
    padding: 16px 24px !important;
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
    color: $text-secondary;
    line-height: 1.8;
    font-size: 15px;

    .highlight {
      color: $primary-gold;
      font-weight: 600;
    }
  }
}

:deep(.el-button) {
  border-radius: 6px;
  font-weight: 500;
  transition: all 0.2s ease;

  &.el-button--primary {
    background: linear-gradient(135deg, $primary-gold, $secondary-gold);
    border: none;
    color: $dark-bg;

    &:hover {
      background: linear-gradient(135deg, #e0bc44, $primary-gold);
      box-shadow: 0 4px 12px rgba(212, 175, 55, 0.3);
    }
  }

  &.el-button--danger {
    background: $accent-red;
    border: none;

    &:hover {
      background: #ff5e78;
      box-shadow: 0 4px 12px rgba(233, 69, 96, 0.3);
    }
  }

  &.el-button--default {
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.1);
    color: $text-secondary;

    &:hover {
      border-color: $primary-gold;
      color: $primary-gold;
    }
  }
}

:deep(.el-tag) {
  border-radius: 4px;
  font-weight: 500;
}

:deep(.el-pagination) {
  .el-pager li {
    background: rgba(255, 255, 255, 0.05);
    color: $text-secondary;

    &.is-active {
      background: $primary-gold;
      color: $dark-bg;
    }

    &:hover:not(.is-active) {
      color: $primary-gold;
    }
  }

  .btn-prev,
  .btn-next {
    background: rgba(255, 255, 255, 0.05);
    color: $text-secondary;

    &:hover {
      color: $primary-gold;
    }
  }

  .el-pagination__total,
  .el-pagination__jump {
    color: $text-secondary;
  }
}

:deep(.el-descriptions) {
  --el-descriptions-text-color: $text-primary;
  --el-descriptions-label-text-color: $text-secondary;
  --el-descriptions-border-color: rgba(255, 255, 255, 0.1);
}
</style>