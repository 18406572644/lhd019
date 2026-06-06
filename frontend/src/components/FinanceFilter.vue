<template>
  <div class="finance-filter glass-card">
    <el-form :inline="true" :model="filterForm" class="filter-form">
      <el-form-item label="周期">
        <el-select v-model="filterForm.period" placeholder="选择周期" style="width: 150px" @change="handlePeriodChange">
          <el-option label="今日" value="today" />
          <el-option label="本周" value="week" />
          <el-option label="本月" value="month" />
          <el-option label="本季度" value="quarter" />
          <el-option label="本年" value="year" />
          <el-option label="自定义" value="custom" />
        </el-select>
      </el-form-item>
      <el-form-item v-if="filterForm.period === 'custom'" label="日期">
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          value-format="YYYY-MM-DD"
          @change="handleDateChange"
        />
      </el-form-item>
      <el-form-item v-if="showCategory" label="分类">
        <el-select v-model="filterForm.category" placeholder="全部分类" style="width: 150px" clearable @change="emitChange">
          <el-option v-for="cat in categories" :key="cat" :label="cat" :value="cat" />
        </el-select>
      </el-form-item>
      <el-form-item v-if="showPaymentMethod" label="支付方式">
        <el-select v-model="filterForm.payment_method" placeholder="全部方式" style="width: 150px" clearable @change="emitChange">
          <el-option v-for="method in paymentMethods" :key="method" :label="method" :value="method" />
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="handleSearch">
          <el-icon><Search /></el-icon>
          查询
        </el-button>
        <el-button @click="handleReset">重置</el-button>
        <el-button type="success" @click="handleExport">
          <el-icon><Download /></el-icon>
          导出PDF
        </el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { Search, Download } from '@element-plus/icons-vue'
import type { FinanceFilterParams } from '@/api'
import { getDateRangeByPeriod } from '@/utils/dateUtils'

interface Props {
  showCategory?: boolean
  showPaymentMethod?: boolean
  categories?: string[]
  paymentMethods?: string[]
}

const props = withDefaults(defineProps<Props>(), {
  showCategory: false,
  showPaymentMethod: false,
  categories: () => [],
  paymentMethods: () => []
})

const emit = defineEmits<{
  (e: 'change', params: FinanceFilterParams): void
  (e: 'export'): void
}>()

const dateRange = ref<string[]>([])
const filterForm = reactive<FinanceFilterParams>({
  period: 'month'
})

const handlePeriodChange = () => {
  if (filterForm.period !== 'custom') {
    const dates = getDateRangeByPeriod(filterForm.period)
    filterForm.start_date = dates.startDate
    filterForm.end_date = dates.endDate
    dateRange.value = [dates.startDate, dates.endDate]
    emitChange()
  }
}

const handleDateChange = (val: string[]) => {
  if (val && val.length === 2) {
    filterForm.start_date = val[0]
    filterForm.end_date = val[1]
    emitChange()
  }
}

const handleSearch = () => {
  emitChange()
}

const handleReset = () => {
  filterForm.period = 'month'
  filterForm.category = ''
  filterForm.payment_method = ''
  const dates = getDateRangeByPeriod('month')
  filterForm.start_date = dates.startDate
  filterForm.end_date = dates.endDate
  dateRange.value = []
  emitChange()
}

const handleExport = () => {
  emit('export')
}

const emitChange = () => {
  emit('change', { ...filterForm })
}

watch(() => props.categories, (newVal) => {
  if (newVal.length > 0 && !filterForm.category) {
    // Optionally set default
  }
}, { immediate: true })
</script>

<style lang="scss" scoped>
.finance-filter {
  margin-bottom: 24px;
  padding: 20px;

  .filter-form {
    display: flex;
    flex-wrap: wrap;
    gap: 16px;
    margin: 0;
  }
}
</style>
