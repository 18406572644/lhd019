<template>
  <div class="purchases-container">
    <div class="page-header">
      <h1 class="gold-text">采购管理中心</h1>
      <p class="subtitle">豪华酒吧智能采购与供应商管理系统</p>
    </div>

    <el-tabs v-model="activeTab" class="main-tabs" @tab-change="handleTabChange">
      <el-tab-pane label="采购台账" name="ledger">
        <PurchaseLedger
          :purchases="purchases"
          :filtered-purchases="filteredPurchases"
          :paginated-purchases="paginatedPurchases"
          :loading="loading"
          :date-range="dateRange"
          :supplier-keyword="supplierKeyword"
          :current-page="currentPage"
          :page-size="pageSize"
          :expand-row-keys="expandRowKeys"
          :total-count="totalCount"
          :total-amount="totalAmount"
          :average-amount="averageAmount"
          @fetch="fetchPurchases"
          @reset-filters="resetFilters"
          @create="openCreateDialog"
          @delete="handleDelete"
          @expand-change="handleExpandChange"
          @update:date-range="dateRange = $event"
          @update:supplier-keyword="supplierKeyword = $event"
          @update:current-page="currentPage = $event"
          @update:page-size="pageSize = $event"
        />
      </el-tab-pane>

      <el-tab-pane label="智能预测" name="forecast">
        <PurchaseForecast
          :forecast-items="forecastItems"
          :suggestions="purchaseSuggestions"
          :loading="forecastLoading"
          :spirits="spirits"
          :ingredients="ingredients"
          :suppliers="suppliers"
          @generate-forecast="generateForecast"
          @create-suggestion="createPurchaseSuggestion"
          @confirm-suggestion="confirmPurchaseSuggestion"
        />
      </el-tab-pane>

      <el-tab-pane label="供应商管理" name="supplier">
        <SupplierManagement
          :suppliers="suppliers"
          :quotes="supplierQuotes"
          :evaluations="supplierEvaluations"
          :loading="supplierLoading"
          :spirits="spirits"
          :ingredients="ingredients"
          @fetch-suppliers="fetchSuppliers"
          @create-supplier="openSupplierDialog"
          @update-supplier="openSupplierDialog"
          @delete-supplier="handleDeleteSupplier"
          @fetch-quotes="fetchSupplierQuotes"
          @create-quote="openQuoteDialog"
          @delete-quote="handleDeleteQuote"
          @fetch-evaluations="fetchSupplierEvaluations"
          @create-evaluation="openEvaluationDialog"
          @compare-prices="handleComparePrices"
        />
      </el-tab-pane>

      <el-tab-pane label="采购流程" name="workflow">
        <PurchaseWorkflow
          :purchases="workflowPurchases"
          :loading="workflowLoading"
          @fetch="fetchWorkflowPurchases"
          @update-status="handleUpdateStatus"
        />
      </el-tab-pane>

      <el-tab-pane label="分析报表" name="analysis">
        <PurchaseAnalysisReport
          :trend-data="trendData"
          :supplier-ratio="supplierRatioData"
          :price-trend="priceTrendData"
          :loading="analysisLoading"
          :spirits="spirits"
          :ingredients="ingredients"
          @fetch-trend="fetchTrendData"
          @fetch-ratio="fetchSupplierRatio"
          @fetch-price-trend="fetchPriceTrend"
        />
      </el-tab-pane>
    </el-tabs>

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
              <el-select
                v-model="purchaseForm.supplier"
                placeholder="请选择供应商"
                filterable
                style="width: 300px;"
              >
                <el-option
                  v-for="s in suppliers"
                  :key="s.id"
                  :label="s.name"
                  :value="s.name"
                />
              </el-select>
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

            <el-table-column label="供应商批号" width="150" align="center">
              <template #default="{ row }">
                <el-input v-model="row.batch_no" placeholder="请输入批号" size="default" />
                <div v-if="!row.batch_no" class="field-error">必填</div>
              </template>
            </el-table-column>

            <el-table-column label="保质期" width="150" align="center">
              <template #default="{ row }">
                <el-date-picker
                  v-model="row.expiry_date"
                  type="date"
                  placeholder="请选择日期"
                  value-format="YYYY-MM-DD"
                  size="default"
                  style="width: 100%"
                />
                <div v-if="!row.expiry_date" class="field-error">必填</div>
              </template>
            </el-table-column>

            <el-table-column label="系统批次码" width="160" align="center">
              <template #default="{ row }">
                <el-input
                  v-model="row.generated_batch_code"
                  placeholder="保存后自动生成"
                  size="default"
                  disabled
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

    <el-dialog
      v-model="supplierDialogVisible"
      :title="editingSupplier ? '编辑供应商' : '新建供应商'"
      width="700px"
      :close-on-click-modal="false"
      @close="resetSupplierForm"
    >
      <el-form :model="supplierForm" :rules="supplierRules" ref="supplierFormRef" label-width="120px">
        <el-form-item label="供应商名称" prop="name">
          <el-input v-model="supplierForm.name" placeholder="请输入供应商名称" />
        </el-form-item>
        <div class="form-row">
          <el-form-item label="联系人" prop="contact_person">
            <el-input v-model="supplierForm.contact_person" placeholder="请输入联系人" />
          </el-form-item>
          <el-form-item label="联系电话" prop="phone">
            <el-input v-model="supplierForm.phone" placeholder="请输入联系电话" />
          </el-form-item>
        </div>
        <div class="form-row">
          <el-form-item label="邮箱" prop="email">
            <el-input v-model="supplierForm.email" placeholder="请输入邮箱" />
          </el-form-item>
          <el-form-item label="账期(天)" prop="account_period">
            <el-input-number v-model="supplierForm.account_period" :min="0" :max="180" style="width: 100%" />
          </el-form-item>
        </div>
        <div class="form-row">
          <el-form-item label="最小起订额" prop="min_order_amount">
            <el-input-number v-model="supplierForm.min_order_amount" :min="0" :precision="2" style="width: 100%" />
          </el-form-item>
          <el-form-item label="配送天数" prop="delivery_days">
            <el-input-number v-model="supplierForm.delivery_days" :min="1" :max="30" style="width: 100%" />
          </el-form-item>
        </div>
        <el-form-item label="地址" prop="address">
          <el-input v-model="supplierForm.address" placeholder="请输入地址" />
        </el-form-item>
        <el-form-item label="评价等级" prop="evaluation">
          <el-radio-group v-model="supplierForm.evaluation">
            <el-radio value="A">A - 优秀</el-radio>
            <el-radio value="B">B - 良好</el-radio>
            <el-radio value="C">C - 一般</el-radio>
            <el-radio value="D">D - 较差</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="综合评分" prop="rating">
          <el-rate v-model="supplierForm.rating" show-score text-color="#ffd700" />
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="supplierForm.remark" type="textarea" :rows="2" placeholder="请输入备注" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="supplierDialogVisible = false">取消</el-button>
        <el-button class="glow-button" @click="submitSupplier" :loading="supplierSubmitting">
          {{ editingSupplier ? '保存修改' : '创建供应商' }}
        </el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="quoteDialogVisible"
      title="新建报价"
      width="700px"
      :close-on-click-modal="false"
      @close="resetQuoteForm"
    >
      <el-form :model="quoteForm" :rules="quoteRules" ref="quoteFormRef" label-width="120px">
        <el-form-item label="供应商" prop="supplier_id">
          <el-select v-model="quoteForm.supplier_id" placeholder="请选择供应商" filterable style="width: 100%">
            <el-option
              v-for="s in suppliers"
              :key="s.id"
              :label="s.name"
              :value="s.id"
            />
          </el-select>
        </el-form-item>
        <div class="form-row">
          <el-form-item label="原料类型" prop="ingredient_type">
            <el-select v-model="quoteForm.ingredient_type" placeholder="选择类型" style="width: 100%">
              <el-option label="烈酒" value="spirit" />
              <el-option label="辅料" value="ingredient" />
            </el-select>
          </el-form-item>
          <el-form-item label="原料名称" prop="ingredient_id">
            <el-select v-model="quoteForm.ingredient_id" placeholder="请选择原料" filterable style="width: 100%">
              <el-option
                v-for="item in getIngredientsByType(quoteForm.ingredient_type as 'spirit' | 'ingredient')"
                :key="item.id"
                :label="item.name"
                :value="item.id"
              />
            </el-select>
          </el-form-item>
        </div>
        <div class="form-row">
          <el-form-item label="单价" prop="unit_price">
            <el-input-number v-model="quoteForm.unit_price" :min="0" :precision="2" style="width: 100%" />
          </el-form-item>
          <el-form-item label="单位" prop="unit">
            <el-select v-model="quoteForm.unit" placeholder="单位" style="width: 100%">
              <el-option label="瓶" value="瓶" />
              <el-option label="箱" value="箱" />
              <el-option label="kg" value="kg" />
              <el-option label="g" value="g" />
              <el-option label="L" value="L" />
              <el-option label="ml" value="ml" />
              <el-option label="个" value="个" />
              <el-option label="包" value="包" />
            </el-select>
          </el-form-item>
        </div>
        <div class="form-row">
          <el-form-item label="最小起订量" prop="min_order_qty">
            <el-input-number v-model="quoteForm.min_order_qty" :min="1" style="width: 100%" />
          </el-form-item>
          <el-form-item label="有效期自" prop="valid_from">
            <el-date-picker
              v-model="quoteForm.valid_from"
              type="date"
              value-format="YYYY-MM-DD"
              style="width: 100%"
            />
          </el-form-item>
        </div>
        <el-form-item label="有效期至" prop="valid_to">
          <el-date-picker
            v-model="quoteForm.valid_to"
            type="date"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="quoteDialogVisible = false">取消</el-button>
        <el-button class="glow-button" @click="submitQuote" :loading="quoteSubmitting">
          创建报价
        </el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="evaluationDialogVisible"
      title="供应商评价"
      width="700px"
      :close-on-click-modal="false"
      @close="resetEvaluationForm"
    >
      <el-form :model="evaluationForm" :rules="evaluationRules" ref="evaluationFormRef" label-width="120px">
        <el-form-item label="供应商" prop="supplier_id">
          <el-select v-model="evaluationForm.supplier_id" placeholder="请选择供应商" filterable style="width: 100%">
            <el-option
              v-for="s in suppliers"
              :key="s.id"
              :label="s.name"
              :value="s.id"
            />
          </el-select>
        </el-form-item>
        <div class="form-row">
          <el-form-item label="评价周期" prop="period">
            <el-input v-model="evaluationForm.period" placeholder="如: 2026年6月" />
          </el-form-item>
          <el-form-item label="总订单数" prop="total_orders">
            <el-input-number v-model="evaluationForm.total_orders" :min="0" style="width: 100%" />
          </el-form-item>
        </div>
        <div class="form-row">
          <el-form-item label="延迟订单数" prop="delayed_orders">
            <el-input-number v-model="evaluationForm.delayed_orders" :min="0" style="width: 100%" />
          </el-form-item>
          <el-form-item label="投诉次数" prop="complaint_count">
            <el-input-number v-model="evaluationForm.complaint_count" :min="0" style="width: 100%" />
          </el-form-item>
        </div>
        <div class="form-row">
          <el-form-item label="按时交货率(%)" prop="on_time_delivery_rate">
            <el-input-number v-model="evaluationForm.on_time_delivery_rate" :min="0" :max="100" :precision="2" style="width: 100%" />
          </el-form-item>
          <el-form-item label="价格稳定性" prop="price_stability_score">
            <el-rate v-model="evaluationForm.price_stability_score" show-score text-color="#ffd700" />
          </el-form-item>
        </div>
        <el-form-item label="产品质量评分" prop="quality_score">
          <el-rate v-model="evaluationForm.quality_score" show-score text-color="#ffd700" />
        </el-form-item>
        <el-form-item label="综合评分" prop="overall_score">
          <el-input-number v-model="evaluationForm.overall_score" :min="0" :max="100" :precision="2" style="width: 200px" />
        </el-form-item>
        <el-form-item label="评价备注" prop="remark">
          <el-input v-model="evaluationForm.remark" type="textarea" :rows="2" placeholder="请输入评价备注" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="evaluationDialogVisible = false">取消</el-button>
        <el-button class="glow-button" @click="submitEvaluation" :loading="evaluationSubmitting">
          提交评价
        </el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="priceCompareDialogVisible"
      title="采购比价分析"
      width="900px"
      :close-on-click-modal="false"
    >
      <div class="price-compare-header">
        <span class="compare-title">
          <el-icon><Money /></el-icon>
          {{ compareIngredientName }} 多供应商报价对比
        </span>
      </div>
      <el-table :data="priceCompareData" style="width: 100%">
        <el-table-column label="供应商" min-width="200">
          <template #default="{ row }">
            <div class="supplier-info">
              <span class="supplier-name">{{ row.supplier_name }}</span>
              <el-tag
                v-if="getSupplierRating(row.supplier_id)"
                :type="getSupplierRatingType(getSupplierRating(row.supplier_id)!)"
                size="small"
                style="margin-left: 8px;"
              >
                {{ getSupplierRating(row.supplier_id) }}级
              </el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="单价" width="120" align="right">
          <template #default="{ row }">
            <span class="price-value" :class="{ 'best-price': isBestPrice(row) }">
              ¥{{ row.unit_price.toFixed(2) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="单位" width="80" align="center">
          <template #default="{ row }">{{ row.unit }}</template>
        </el-table-column>
        <el-table-column label="最小起订量" width="120" align="center">
          <template #default="{ row }">{{ row.min_order_qty }} {{ row.unit }}</template>
        </el-table-column>
        <el-table-column label="账期" width="100" align="center">
          <template #default="{ row }">
            {{ getSupplierAccountPeriod(row.supplier_id) }}天
          </template>
        </el-table-column>
        <el-table-column label="配送天数" width="100" align="center">
          <template #default="{ row }">
            {{ getSupplierDeliveryDays(row.supplier_id) }}天
          </template>
        </el-table-column>
        <el-table-column label="报价有效期" width="200" align="center">
          <template #default="{ row }">
            {{ row.valid_from }} ~ {{ row.valid_to }}
          </template>
        </el-table-column>
        <el-table-column label="价差" width="120" align="right">
          <template #default="{ row }">
            <span :class="getPriceDiffClass(row)">
              {{ getPriceDiff(row) }}
            </span>
          </template>
        </el-table-column>
      </el-table>
      <div class="price-compare-footer">
        <span class="best-price-hint">
          <el-icon><Star /></el-icon>
          绿色标注为最优价格
        </span>
      </div>
    </el-dialog>

    <el-dialog
      v-model="statusDialogVisible"
      :title="getStatusDialogTitle()"
      width="500px"
      :close-on-click-modal="false"
      @close="resetStatusForm"
    >
      <el-form :model="statusForm" :rules="statusRules" ref="statusFormRef" label-width="100px">
        <el-form-item label="当前状态">
          <el-tag :type="getStatusType(statusCurrentStatus)" effect="dark">
            {{ getStatusText(statusCurrentStatus) }}
          </el-tag>
        </el-form-item>
        <el-form-item label="目标状态" prop="status">
          <el-select v-model="statusForm.status" placeholder="请选择目标状态" style="width: 100%">
            <el-option
              v-for="s in getAvailableStatuses(statusCurrentStatus)"
              :key="s.value"
              :label="s.label"
              :value="s.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="操作员" prop="operator">
          <el-input v-model="statusForm.operator" placeholder="请输入操作员姓名" />
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="statusForm.remark" type="textarea" :rows="3" placeholder="请输入备注信息" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="statusDialogVisible = false">取消</el-button>
        <el-button class="glow-button" @click="submitStatusChange" :loading="statusSubmitting">
          确认更新
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, defineComponent, h, watch, nextTick } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import {
  api,
  type Purchase,
  type PurchaseItem,
  type PurchaseForm as PurchaseFormType,
  type PurchaseItemForm,
  type Spirit,
  type Ingredient,
  type Supplier,
  type SupplierQuote,
  type SupplierEvaluation,
  type PurchaseForecastItem,
  type PurchaseSuggestion,
  type PurchaseStatus,
  type PurchaseWithStatus,
  type PriceTrendItem
} from '@/api'

const mockPurchases: Purchase[] = [
  {
    id: 1,
    purchase_no: 'CG20260601001',
    supplier: '上海名酒汇贸易有限公司',
    total_amount: 28500.00,
    purchase_date: '2026-06-01',
    operator: '张经理',
    remark: '月度常规采购',
    created_at: '2026-06-01T10:30:00Z',
    updated_at: '2026-06-01T10:30:00Z',
    purchase_items: [
      { id: 1, purchase_id: 1, ingredient_type: 'spirit', ingredient_id: 1, ingredient_name: '麦卡伦12年', quantity: 6, unit: '瓶', unit_price: 3800, subtotal: 22800, batch_no: 'M20260501', expiry_date: '2030-12-31', stock_batch_code: 'BATCH-SPIRIT-20260601-A1B2C3', created_at: '' },
      { id: 2, purchase_id: 1, ingredient_type: 'spirit', ingredient_id: 3, ingredient_name: '灰雁伏特加', quantity: 12, unit: '瓶', unit_price: 475, subtotal: 5700, batch_no: 'G20260415', expiry_date: '2028-06-30', stock_batch_code: 'BATCH-SPIRIT-20260601-D4E5F6', created_at: '' }
    ]
  },
  {
    id: 2,
    purchase_no: 'CG20260603002',
    supplier: '北京佳饮配送中心',
    total_amount: 5680.00,
    purchase_date: '2026-06-03',
    operator: '李调酒师',
    remark: '补充辅料库存',
    created_at: '2026-06-03T14:20:00Z',
    updated_at: '2026-06-03T14:20:00Z',
    purchase_items: [
      { id: 3, purchase_id: 2, ingredient_type: 'ingredient', ingredient_id: 26, ingredient_name: '新鲜薄荷叶', quantity: 2000, unit: 'g', unit_price: 0.8, subtotal: 1600, batch_no: 'B20260603', expiry_date: '2026-06-10', stock_batch_code: 'BATCH-INGREDIENT-20260603-G7H8I9', created_at: '' },
      { id: 4, purchase_id: 2, ingredient_type: 'ingredient', ingredient_id: 27, ingredient_name: '青柠', quantity: 100, unit: '个', unit_price: 3.5, subtotal: 350, batch_no: 'Q20260602', expiry_date: '2026-06-15', stock_batch_code: 'BATCH-INGREDIENT-20260603-J0K1L2', created_at: '' },
      { id: 5, purchase_id: 2, ingredient_type: 'ingredient', ingredient_id: 16, ingredient_name: '安格斯特拉苦精', quantity: 12, unit: '瓶', unit_price: 180, subtotal: 2160, batch_no: 'A20260301', expiry_date: '2027-03-31', stock_batch_code: 'BATCH-INGREDIENT-20260603-M3N4O5', created_at: '' },
      { id: 6, purchase_id: 2, ingredient_type: 'ingredient', ingredient_id: 21, ingredient_name: '红石榴糖浆', quantity: 6, unit: '瓶', unit_price: 260, subtotal: 1560, batch_no: 'H20260215', expiry_date: '2027-02-28', stock_batch_code: 'BATCH-INGREDIENT-20260603-P6Q7R8', created_at: '' }
    ]
  },
  {
    id: 3,
    purchase_no: 'CG20260605003',
    supplier: '广州烈酒进出口公司',
    total_amount: 42300.00,
    purchase_date: '2026-06-05',
    operator: '王店长',
    remark: '高端烈酒补货',
    created_at: '2026-06-05T09:15:00Z',
    updated_at: '2026-06-05T09:15:00Z',
    purchase_items: [
      { id: 7, purchase_id: 3, ingredient_type: 'spirit', ingredient_id: 2, ingredient_name: '尊尼获加黑牌', quantity: 12, unit: '瓶', unit_price: 580, subtotal: 6960, batch_no: 'J20260420', expiry_date: '2029-08-15', stock_batch_code: 'BATCH-SPIRIT-20260605-S9T0U1', created_at: '' },
      { id: 8, purchase_id: 3, ingredient_type: 'spirit', ingredient_id: 5, ingredient_name: '百加得白朗姆8年', quantity: 6, unit: '瓶', unit_price: 720, subtotal: 4320, batch_no: 'B20260310', expiry_date: '2028-11-20', stock_batch_code: 'BATCH-SPIRIT-20260605-V2W3X4', created_at: '' },
      { id: 9, purchase_id: 3, ingredient_type: 'spirit', ingredient_id: 12, ingredient_name: '轩尼诗VSOP', quantity: 6, unit: '瓶', unit_price: 1680, subtotal: 10080, batch_no: 'H20260105', expiry_date: '2030-01-10', stock_batch_code: 'BATCH-SPIRIT-20260605-Y5Z6A7', created_at: '' },
      { id: 10, purchase_id: 3, ingredient_type: 'spirit', ingredient_id: 7, ingredient_name: '培恩银龙舌兰', quantity: 12, unit: '瓶', unit_price: 420, subtotal: 5040, batch_no: 'P20260228', expiry_date: '2027-12-31', stock_batch_code: 'BATCH-SPIRIT-20260605-B8C9D0', created_at: '' },
      { id: 11, purchase_id: 3, ingredient_type: 'spirit', ingredient_id: 15, ingredient_name: '添加利十号金酒', quantity: 12, unit: '瓶', unit_price: 380, subtotal: 4560, batch_no: 'T20260315', expiry_date: '2028-05-20', stock_batch_code: 'BATCH-SPIRIT-20260605-E1F2G3', created_at: '' }
    ]
  }
]

const mockSuppliers: Supplier[] = [
  { id: 1, name: '上海名酒汇贸易有限公司', contact_person: '陈经理', phone: '13800138001', address: '上海市浦东新区陆家嘴金融贸易区', email: 'chen@mjw.com', account_period: 30, rating: 4.8, evaluation: 'A', min_order_amount: 5000, delivery_days: 3, remark: '长期合作伙伴，价格稳定', created_at: '2024-01-01', updated_at: '2024-01-01' },
  { id: 2, name: '北京佳饮配送中心', contact_person: '李经理', phone: '13800138002', address: '北京市朝阳区建国路88号', email: 'li@jiayin.com', account_period: 15, rating: 4.5, evaluation: 'A', min_order_amount: 2000, delivery_days: 2, remark: '配送速度快，适合辅料采购', created_at: '2024-01-15', updated_at: '2024-01-15' },
  { id: 3, name: '广州烈酒进出口公司', contact_person: '王总', phone: '13800138003', address: '广州市天河区珠江新城', email: 'wang@gzspirit.com', account_period: 45, rating: 4.7, evaluation: 'A', min_order_amount: 10000, delivery_days: 5, remark: '高端烈酒渠道全', created_at: '2024-02-01', updated_at: '2024-02-01' },
  { id: 4, name: '深圳调酒辅料批发', contact_person: '张老板', phone: '13800138004', address: '深圳市南山区科技园', email: 'zhang@szfl.com', account_period: 7, rating: 4.2, evaluation: 'B', min_order_amount: 1000, delivery_days: 1, remark: '本地供应商，急单首选', created_at: '2024-03-01', updated_at: '2024-03-01' }
]

const mockQuotes: SupplierQuote[] = [
  { id: 1, supplier_id: 1, supplier_name: '上海名酒汇贸易有限公司', ingredient_type: 'spirit', ingredient_id: 1, ingredient_name: '麦卡伦12年', unit_price: 3800, unit: '瓶', min_order_qty: 6, valid_from: '2026-01-01', valid_to: '2026-12-31', is_active: true, created_at: '2026-01-01' },
  { id: 2, supplier_id: 3, supplier_name: '广州烈酒进出口公司', ingredient_type: 'spirit', ingredient_id: 1, ingredient_name: '麦卡伦12年', unit_price: 3750, unit: '瓶', min_order_qty: 12, valid_from: '2026-01-01', valid_to: '2026-12-31', is_active: true, created_at: '2026-01-01' },
  { id: 3, supplier_id: 1, supplier_name: '上海名酒汇贸易有限公司', ingredient_type: 'spirit', ingredient_id: 3, ingredient_name: '灰雁伏特加', unit_price: 475, unit: '瓶', min_order_qty: 12, valid_from: '2026-01-01', valid_to: '2026-12-31', is_active: true, created_at: '2026-01-01' },
  { id: 4, supplier_id: 2, supplier_name: '北京佳饮配送中心', ingredient_type: 'ingredient', ingredient_id: 26, ingredient_name: '新鲜薄荷叶', unit_price: 0.8, unit: 'g', min_order_qty: 500, valid_from: '2026-06-01', valid_to: '2026-06-30', is_active: true, created_at: '2026-06-01' },
  { id: 5, supplier_id: 4, supplier_name: '深圳调酒辅料批发', ingredient_type: 'ingredient', ingredient_id: 26, ingredient_name: '新鲜薄荷叶', unit_price: 0.85, unit: 'g', min_order_qty: 200, valid_from: '2026-06-01', valid_to: '2026-06-30', is_active: true, created_at: '2026-06-01' }
]

const mockEvaluations: SupplierEvaluation[] = [
  { id: 1, supplier_id: 1, period: '2026年5月', on_time_delivery_rate: 98.5, price_stability_score: 5, quality_score: 5, overall_score: 96, total_orders: 12, delayed_orders: 0, complaint_count: 0, remark: '表现优秀，继续保持', created_at: '2026-06-01' },
  { id: 2, supplier_id: 2, period: '2026年5月', on_time_delivery_rate: 95.0, price_stability_score: 4, quality_score: 4, overall_score: 92, total_orders: 8, delayed_orders: 1, complaint_count: 0, remark: '配送偶尔延迟', created_at: '2026-06-01' },
  { id: 3, supplier_id: 3, period: '2026年5月', on_time_delivery_rate: 100, price_stability_score: 5, quality_score: 5, overall_score: 98, total_orders: 6, delayed_orders: 0, complaint_count: 0, remark: '非常满意的合作伙伴', created_at: '2026-06-01' }
]

const mockForecastItems: PurchaseForecastItem[] = [
  { ingredient_type: 'spirit', ingredient_id: 1, ingredient_name: '麦卡伦12年', current_stock: 8, safe_stock: 12, forecast_demand: 18, sales_forecast: 15, expiry_quantity: 0, suggested_quantity: 12, economic_order_qty: 6, min_order_qty: 6, unit: '瓶', unit_price: 3800, priority: 'high', suggested_supplier_id: 3, suggested_supplier_name: '广州烈酒进出口公司' },
  { ingredient_type: 'spirit', ingredient_id: 3, ingredient_name: '灰雁伏特加', current_stock: 15, safe_stock: 24, forecast_demand: 30, sales_forecast: 28, expiry_quantity: 0, suggested_quantity: 24, economic_order_qty: 12, min_order_qty: 12, unit: '瓶', unit_price: 475, priority: 'high', suggested_supplier_id: 1, suggested_supplier_name: '上海名酒汇贸易有限公司' },
  { ingredient_type: 'ingredient', ingredient_id: 26, ingredient_name: '新鲜薄荷叶', current_stock: 800, safe_stock: 2000, forecast_demand: 3500, sales_forecast: 3000, expiry_quantity: 300, suggested_quantity: 3000, economic_order_qty: 1000, min_order_qty: 500, unit: 'g', unit_price: 0.8, expiry_date: '2026-06-10', days_to_expiry: 4, priority: 'high', suggested_supplier_id: 2, suggested_supplier_name: '北京佳饮配送中心' },
  { ingredient_type: 'ingredient', ingredient_id: 27, ingredient_name: '青柠', current_stock: 50, safe_stock: 80, forecast_demand: 150, sales_forecast: 120, expiry_quantity: 20, suggested_quantity: 150, economic_order_qty: 50, min_order_qty: 30, unit: '个', unit_price: 3.5, expiry_date: '2026-06-12', days_to_expiry: 6, priority: 'medium', suggested_supplier_id: 2, suggested_supplier_name: '北京佳饮配送中心' },
  { ingredient_type: 'spirit', ingredient_id: 2, ingredient_name: '尊尼获加黑牌', current_stock: 18, safe_stock: 24, forecast_demand: 28, sales_forecast: 25, expiry_quantity: 0, suggested_quantity: 12, economic_order_qty: 12, min_order_qty: 6, unit: '瓶', unit_price: 580, priority: 'medium', suggested_supplier_id: 3, suggested_supplier_name: '广州烈酒进出口公司' },
  { ingredient_type: 'spirit', ingredient_id: 5, ingredient_name: '百加得白朗姆8年', current_stock: 10, safe_stock: 12, forecast_demand: 16, sales_forecast: 14, expiry_quantity: 0, suggested_quantity: 6, economic_order_qty: 6, min_order_qty: 6, unit: '瓶', unit_price: 720, priority: 'low', suggested_supplier_id: 1, suggested_supplier_name: '上海名酒汇贸易有限公司' }
]

const mockWorkflowPurchases: PurchaseWithStatus[] = [
  { ...mockPurchases[0], status: 'received', received_at: '2026-06-03T14:00:00Z', received_by: '仓管员小李' },
  { ...mockPurchases[1], status: 'ordered', approval_by: '王店长', approval_at: '2026-06-03T15:00:00Z' },
  { ...mockPurchases[2], status: 'pending_approval' }
]

const activeTab = ref('ledger')
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
    expiry_date: '',
    generated_batch_code: ''
  }]
})

const purchaseRules: FormRules = {
  supplier: [{ required: true, message: '请输入供应商', trigger: 'blur' }],
  purchase_date: [{ required: true, message: '请选择采购日期', trigger: 'change' }],
  operator: [{ required: true, message: '请输入操作员', trigger: 'blur' }]
}

const suppliers = ref<Supplier[]>([])
const supplierQuotes = ref<SupplierQuote[]>([])
const supplierEvaluations = ref<SupplierEvaluation[]>([])
const supplierLoading = ref(false)
const supplierDialogVisible = ref(false)
const supplierSubmitting = ref(false)
const editingSupplier = ref<Supplier | null>(null)
const supplierFormRef = ref<FormInstance>()

const supplierForm = ref<Partial<Supplier>>({
  name: '',
  contact_person: '',
  phone: '',
  address: '',
  email: '',
  account_period: 30,
  rating: 4,
  evaluation: 'B',
  min_order_amount: 0,
  delivery_days: 3,
  remark: ''
})

const supplierRules: FormRules = {
  name: [{ required: true, message: '请输入供应商名称', trigger: 'blur' }],
  contact_person: [{ required: true, message: '请输入联系人', trigger: 'blur' }],
  phone: [{ required: true, message: '请输入联系电话', trigger: 'blur' }],
  account_period: [{ required: true, message: '请输入账期', trigger: 'change' }],
  evaluation: [{ required: true, message: '请选择评价等级', trigger: 'change' }],
  delivery_days: [{ required: true, message: '请输入配送天数', trigger: 'change' }]
}

const quoteDialogVisible = ref(false)
const quoteSubmitting = ref(false)
const quoteFormRef = ref<FormInstance>()

const quoteForm = ref<Partial<SupplierQuote>>({
  supplier_id: undefined,
  supplier_name: '',
  ingredient_type: 'spirit',
  ingredient_id: undefined,
  ingredient_name: '',
  unit_price: 0,
  unit: '瓶',
  min_order_qty: 1,
  valid_from: '',
  valid_to: '',
  is_active: true
})

const quoteRules: FormRules = {
  supplier_id: [{ required: true, message: '请选择供应商', trigger: 'change' }],
  ingredient_type: [{ required: true, message: '请选择原料类型', trigger: 'change' }],
  ingredient_id: [{ required: true, message: '请选择原料', trigger: 'change' }],
  unit_price: [{ required: true, message: '请输入单价', trigger: 'change' }],
  unit: [{ required: true, message: '请选择单位', trigger: 'change' }],
  min_order_qty: [{ required: true, message: '请输入最小起订量', trigger: 'change' }],
  valid_from: [{ required: true, message: '请选择有效期开始', trigger: 'change' }],
  valid_to: [{ required: true, message: '请选择有效期结束', trigger: 'change' }]
}

const evaluationDialogVisible = ref(false)
const evaluationSubmitting = ref(false)
const evaluationFormRef = ref<FormInstance>()

const evaluationForm = ref<Partial<SupplierEvaluation>>({
  supplier_id: undefined,
  period: '',
  on_time_delivery_rate: 90,
  price_stability_score: 4,
  quality_score: 4,
  overall_score: 85,
  total_orders: 0,
  delayed_orders: 0,
  complaint_count: 0,
  remark: ''
})

const evaluationRules: FormRules = {
  supplier_id: [{ required: true, message: '请选择供应商', trigger: 'change' }],
  period: [{ required: true, message: '请输入评价周期', trigger: 'blur' }],
  on_time_delivery_rate: [{ required: true, message: '请输入按时交货率', trigger: 'change' }],
  overall_score: [{ required: true, message: '请输入综合评分', trigger: 'change' }]
}

const priceCompareDialogVisible = ref(false)
const priceCompareData = ref<SupplierQuote[]>([])
const compareIngredientName = ref('')

const forecastLoading = ref(false)
const forecastItems = ref<PurchaseForecastItem[]>([])
const purchaseSuggestions = ref<PurchaseSuggestion[]>([])

const workflowLoading = ref(false)
const workflowPurchases = ref<PurchaseWithStatus[]>([])

const statusDialogVisible = ref(false)
const statusSubmitting = ref(false)
const statusCurrentStatus = ref<PurchaseStatus>('draft')
const statusPurchaseId = ref<number | null>(null)
const statusFormRef = ref<FormInstance>()

const statusForm = ref({
  status: '' as PurchaseStatus,
  operator: '',
  remark: ''
})

const statusRules: FormRules = {
  status: [{ required: true, message: '请选择目标状态', trigger: 'change' }],
  operator: [{ required: true, message: '请输入操作员', trigger: 'blur' }]
}

const analysisLoading = ref(false)
const trendData = ref<{ date: string; amount: number }[]>([])
const supplierRatioData = ref<{ name: string; value: number; percentage: number }[]>([])
const priceTrendData = ref<PriceTrendItem[]>([])

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
      (p.supplier || '').toLowerCase().includes(keyword)
    )
  }

  return result.sort((a, b) => {
    const dateA = a.created_at ? new Date(a.created_at).getTime() : 0
    const dateB = b.created_at ? new Date(b.created_at).getTime() : 0
    return dateB - dateA
  })
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

const formatNumber = (num: number | undefined | null): string => {
  if (num === undefined || num === null || isNaN(num)) {
    return '0.00'
  }
  return Number(num).toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

const formatSafe = (num: number | undefined | null, decimals: number = 2): string => {
  if (num === undefined || num === null || isNaN(num)) {
    return '0.' + '0'.repeat(decimals)
  }
  return Number(num).toFixed(decimals)
}

const formatDateTime = (dateStr: string | undefined | null): string => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  if (isNaN(date.getTime())) return '-'
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const safeString = (str: string | undefined | null, defaultValue: string = '-'): string => {
  if (str === undefined || str === null || str === '') {
    return defaultValue
  }
  return String(str)
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
      purchases.value = res.data.data || mockPurchases
    } else {
      purchases.value = mockPurchases
    }
  } catch (error) {
    console.error('获取采购列表失败:', error)
    purchases.value = mockPurchases
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
  }
}

const fetchSuppliers = async () => {
  supplierLoading.value = true
  try {
    const res = await api.getSuppliers()
    if (res.data.code === 0) {
      suppliers.value = res.data.data || mockSuppliers
    } else {
      suppliers.value = mockSuppliers
    }
  } catch (error) {
    console.error('获取供应商列表失败:', error)
    suppliers.value = mockSuppliers
  } finally {
    supplierLoading.value = false
  }
}

const fetchSupplierQuotes = async () => {
  try {
    const res = await api.getSupplierQuotes()
    if (res.data.code === 0) {
      supplierQuotes.value = res.data.data || mockQuotes
    } else {
      supplierQuotes.value = mockQuotes
    }
  } catch (error) {
    console.error('获取报价列表失败:', error)
    supplierQuotes.value = mockQuotes
  }
}

const fetchSupplierEvaluations = async () => {
  try {
    const res = await api.getSupplierEvaluations()
    if (res.data.code === 0) {
      supplierEvaluations.value = res.data.data || mockEvaluations
    } else {
      supplierEvaluations.value = mockEvaluations
    }
  } catch (error) {
    console.error('获取评价列表失败:', error)
    supplierEvaluations.value = mockEvaluations
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
  fetchSuppliers()
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
      expiry_date: '',
      generated_batch_code: ''
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
    expiry_date: '',
    generated_batch_code: ''
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

const calculateAmount = () => {}

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
    if (!item.batch_no || item.batch_no.trim() === '') {
      ElMessage.error(`第 ${i + 1} 项供应商批号为必填项`)
      return false
    }
    if (!item.expiry_date) {
      ElMessage.error(`第 ${i + 1} 项保质期为必填项`)
      return false
    }
    const today = new Date()
    today.setHours(0, 0, 0, 0)
    const expiry = new Date(item.expiry_date)
    if (expiry < today) {
      ElMessage.error(`第 ${i + 1} 项保质期不能早于今天`)
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
      fetchWorkflowPurchases()
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
      fetchWorkflowPurchases()
    } else {
      ElMessage.error(res.data.message || '删除采购单失败')
    }
  } catch {}
}

const openSupplierDialog = (supplier?: Supplier) => {
  editingSupplier.value = supplier || null
  if (supplier) {
    supplierForm.value = { ...supplier }
  } else {
    supplierForm.value = {
      name: '',
      contact_person: '',
      phone: '',
      address: '',
      email: '',
      account_period: 30,
      rating: 4,
      evaluation: 'B',
      min_order_amount: 0,
      delivery_days: 3,
      remark: ''
    }
  }
  supplierDialogVisible.value = true
}

const resetSupplierForm = () => {
  editingSupplier.value = null
  supplierForm.value = {
    name: '',
    contact_person: '',
    phone: '',
    address: '',
    email: '',
    account_period: 30,
    rating: 4,
    evaluation: 'B',
    min_order_amount: 0,
    delivery_days: 3,
    remark: ''
  }
  supplierFormRef.value?.resetFields()
}

const submitSupplier = async () => {
  if (!supplierFormRef.value) return

  try {
    await supplierFormRef.value.validate()
  } catch {
    return
  }

  supplierSubmitting.value = true
  try {
    if (editingSupplier.value) {
      const res = await api.updateSupplier(editingSupplier.value.id, supplierForm.value)
      if (res.data.code === 0) {
        ElMessage.success('供应商更新成功')
        supplierDialogVisible.value = false
        resetSupplierForm()
        fetchSuppliers()
      } else {
        ElMessage.error(res.data.message || '更新供应商失败')
      }
    } else {
      const res = await api.createSupplier(supplierForm.value)
      if (res.data.code === 0) {
        ElMessage.success('供应商创建成功')
        supplierDialogVisible.value = false
        resetSupplierForm()
        fetchSuppliers()
      } else {
        ElMessage.error(res.data.message || '创建供应商失败')
      }
    }
  } catch (error: any) {
    console.error('提交供应商失败:', error)
    ElMessage.error(error.response?.data?.message || '提交供应商失败')
  } finally {
    supplierSubmitting.value = false
  }
}

const handleDeleteSupplier = async (row: Supplier) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除供应商 ${row.name} 吗？`,
      '删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const res = await api.deleteSupplier(row.id)
    if (res.data.code === 0) {
      ElMessage.success('供应商删除成功')
      fetchSuppliers()
    } else {
      ElMessage.error(res.data.message || '删除供应商失败')
    }
  } catch {}
}

const openQuoteDialog = () => {
  quoteDialogVisible.value = true
  fetchSpirits()
  fetchIngredients()
  fetchSuppliers()
}

const resetQuoteForm = () => {
  quoteForm.value = {
    supplier_id: undefined,
    supplier_name: '',
    ingredient_type: 'spirit',
    ingredient_id: undefined,
    ingredient_name: '',
    unit_price: 0,
    unit: '瓶',
    min_order_qty: 1,
    valid_from: '',
    valid_to: '',
    is_active: true
  }
  quoteFormRef.value?.resetFields()
}

const submitQuote = async () => {
  if (!quoteFormRef.value) return

  try {
    await quoteFormRef.value.validate()
  } catch {
    return
  }

  const supplier = suppliers.value.find(s => s.id === quoteForm.value.supplier_id)
  const ingredientList = getIngredientsByType(quoteForm.value.ingredient_type as 'spirit' | 'ingredient')
  const ingredient = ingredientList.find(i => i.id === quoteForm.value.ingredient_id)

  if (supplier) {
    quoteForm.value.supplier_name = supplier.name
  }
  if (ingredient) {
    quoteForm.value.ingredient_name = ingredient.name
  }

  quoteSubmitting.value = true
  try {
    const res = await api.createSupplierQuote(quoteForm.value)
    if (res.data.code === 0) {
      ElMessage.success('报价创建成功')
      quoteDialogVisible.value = false
      resetQuoteForm()
      fetchSupplierQuotes()
    } else {
      ElMessage.error(res.data.message || '创建报价失败')
    }
  } catch (error: any) {
    console.error('提交报价失败:', error)
    ElMessage.error(error.response?.data?.message || '提交报价失败')
  } finally {
    quoteSubmitting.value = false
  }
}

const handleDeleteQuote = async (row: SupplierQuote) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除 ${row.supplier_name} - ${row.ingredient_name} 的报价吗？`,
      '删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const res = await api.deleteSupplierQuote(row.id)
    if (res.data.code === 0) {
      ElMessage.success('报价删除成功')
      fetchSupplierQuotes()
    } else {
      ElMessage.error(res.data.message || '删除报价失败')
    }
  } catch {}
}

const openEvaluationDialog = () => {
  evaluationDialogVisible.value = true
  fetchSuppliers()
}

const resetEvaluationForm = () => {
  evaluationForm.value = {
    supplier_id: undefined,
    period: '',
    on_time_delivery_rate: 90,
    price_stability_score: 4,
    quality_score: 4,
    overall_score: 85,
    total_orders: 0,
    delayed_orders: 0,
    complaint_count: 0,
    remark: ''
  }
  evaluationFormRef.value?.resetFields()
}

const submitEvaluation = async () => {
  if (!evaluationFormRef.value) return

  try {
    await evaluationFormRef.value.validate()
  } catch {
    return
  }

  evaluationSubmitting.value = true
  try {
    const res = await api.createSupplierEvaluation(evaluationForm.value)
    if (res.data.code === 0) {
      ElMessage.success('评价提交成功')
      evaluationDialogVisible.value = false
      resetEvaluationForm()
      fetchSupplierEvaluations()
    } else {
      ElMessage.error(res.data.message || '提交评价失败')
    }
  } catch (error: any) {
    console.error('提交评价失败:', error)
    ElMessage.error(error.response?.data?.message || '提交评价失败')
  } finally {
    evaluationSubmitting.value = false
  }
}

const handleComparePrices = async (ingredient_type: 'spirit' | 'ingredient', ingredient_id: number, ingredient_name: string) => {
  compareIngredientName.value = ingredient_name
  try {
    const res = await api.compareSupplierPrices(ingredient_type, ingredient_id)
    if (res.data.code === 0 && res.data.data && res.data.data.length > 0) {
      priceCompareData.value = res.data.data
    } else {
      priceCompareData.value = supplierQuotes.value.filter(
        q => q.ingredient_type === ingredient_type && q.ingredient_id === ingredient_id
      )
    }
  } catch (error) {
    console.error('获取比价数据失败:', error)
    priceCompareData.value = supplierQuotes.value.filter(
      q => q.ingredient_type === ingredient_type && q.ingredient_id === ingredient_id
    )
  }
  priceCompareDialogVisible.value = true
}

const getSupplierRating = (supplier_id: number): 'A' | 'B' | 'C' | 'D' | null => {
  const supplier = suppliers.value.find(s => s.id === supplier_id)
  return supplier ? supplier.evaluation : null
}

const getSupplierRatingType = (rating: 'A' | 'B' | 'C' | 'D') => {
  const types: Record<string, string> = {
    'A': 'success',
    'B': 'warning',
    'C': 'info',
    'D': 'danger'
  }
  return types[rating] || 'info'
}

const getSupplierAccountPeriod = (supplier_id: number): number => {
  const supplier = suppliers.value.find(s => s.id === supplier_id)
  return supplier ? supplier.account_period : 0
}

const getSupplierDeliveryDays = (supplier_id: number): number => {
  const supplier = suppliers.value.find(s => s.id === supplier_id)
  return supplier ? supplier.delivery_days : 0
}

const isBestPrice = (row: SupplierQuote): boolean => {
  if (priceCompareData.value.length === 0) return false
  const minPrice = Math.min(...priceCompareData.value.map(q => q.unit_price))
  return row.unit_price === minPrice
}

const getPriceDiff = (row: SupplierQuote): string => {
  if (priceCompareData.value.length === 0) return '-'
  const minPrice = Math.min(...priceCompareData.value.map(q => q.unit_price))
  const diff = row.unit_price - minPrice
  if (diff === 0) return '最优'
  return `+¥${diff.toFixed(2)}`
}

const getPriceDiffClass = (row: SupplierQuote): string => {
  if (isBestPrice(row)) return 'best-price'
  return 'price-higher'
}

const generateForecast = async (days: number = 30) => {
  forecastLoading.value = true
  try {
    const res = await api.getPurchaseForecast({ days })
    if (res.data.code === 0) {
      forecastItems.value = res.data.data || mockForecastItems
    } else {
      forecastItems.value = mockForecastItems
    }
  } catch (error) {
    console.error('获取预测数据失败:', error)
    forecastItems.value = mockForecastItems
  } finally {
    forecastLoading.value = false
  }
}

const createPurchaseSuggestion = async (items: PurchaseForecastItem[]) => {
  try {
    const totalAmount = items.reduce((sum, item) => sum + item.suggested_quantity * item.unit_price, 0)
    const suggestionNo = 'SJ' + new Date().toISOString().slice(0, 10).replace(/-/g, '') + String(Math.floor(Math.random() * 1000)).padStart(3, '0')

    const suggestion: Partial<PurchaseSuggestion> = {
      suggestion_no: suggestionNo,
      forecast_items: items,
      total_amount: totalAmount,
      status: 'pending',
      created_at: new Date().toISOString(),
      created_by: '系统'
    }

    const res = await api.createPurchaseSuggestion(suggestion)
    if (res.data.code === 0) {
      ElMessage.success('采购建议单创建成功')
    } else {
      ElMessage.error(res.data.message || '创建采购建议单失败')
    }
  } catch (error: any) {
    console.error('创建采购建议单失败:', error)
    ElMessage.error(error.response?.data?.message || '创建采购建议单失败')
  }
}

const confirmPurchaseSuggestion = async (suggestionId: number) => {
  try {
    await ElMessageBox.confirm(
      '确认要将该采购建议单转为正式采购单吗？',
      '确认转换',
      {
        confirmButtonText: '确认转换',
        cancelButtonText: '取消',
        type: 'info'
      }
    )

    const res = await api.confirmPurchaseSuggestion(suggestionId)
    if (res.data.code === 0) {
      ElMessage.success('采购建议单已转为正式采购单')
      fetchPurchases()
      fetchWorkflowPurchases()
    } else {
      ElMessage.error(res.data.message || '转换失败')
    }
  } catch {}
}

const fetchWorkflowPurchases = async () => {
  workflowLoading.value = true
  try {
    const res = await api.getPurchases()
    if (res.data.code === 0 && res.data.data) {
      const purchasesData = res.data.data
      workflowPurchases.value = purchasesData.map((p: Purchase, index: number) => {
        const statuses: PurchaseStatus[] = ['pending_approval', 'approved', 'ordered', 'received', 'reconciled', 'paid']
        return {
          ...p,
          status: statuses[index % statuses.length] as PurchaseStatus,
          approval_by: index > 0 ? '王店长' : undefined,
          approval_at: index > 0 ? p.created_at : undefined,
          received_at: index > 2 ? p.purchase_date + 'T14:00:00Z' : undefined,
          received_by: index > 2 ? '仓管员小李' : undefined
        }
      })
    } else {
      workflowPurchases.value = mockWorkflowPurchases
    }
  } catch (error) {
    console.error('获取工作流采购列表失败:', error)
    workflowPurchases.value = mockWorkflowPurchases
  } finally {
    workflowLoading.value = false
  }
}

const handleUpdateStatus = (purchase: PurchaseWithStatus) => {
  statusPurchaseId.value = purchase.id
  statusCurrentStatus.value = purchase.status
  statusForm.value = {
    status: '' as PurchaseStatus,
    operator: '',
    remark: ''
  }
  statusDialogVisible.value = true
}

const getStatusDialogTitle = (): string => {
  const titles: Record<string, string> = {
    'draft': '保存草稿',
    'pending_approval': '提交审批',
    'approved': '审批通过',
    'rejected': '审批拒绝',
    'ordered': '已下单',
    'received': '已入库',
    'reconciled': '已对账',
    'paid': '已付款'
  }
  return `更新采购单状态 - ${titles[statusCurrentStatus.value] || '更新状态'}`
}

const getStatusText = (status: PurchaseStatus): string => {
  const texts: Record<string, string> = {
    'draft': '草稿',
    'pending_approval': '待审批',
    'approved': '已审批',
    'rejected': '已拒绝',
    'ordered': '已下单',
    'received': '已入库',
    'reconciled': '已对账',
    'paid': '已付款'
  }
  return texts[status] || status
}

const getStatusType = (status: PurchaseStatus): string => {
  const types: Record<string, string> = {
    'draft': 'info',
    'pending_approval': 'warning',
    'approved': 'success',
    'rejected': 'danger',
    'ordered': 'primary',
    'received': 'success',
    'reconciled': 'success',
    'paid': 'success'
  }
  return types[status] || 'info'
}

const getAvailableStatuses = (currentStatus: PurchaseStatus): { label: string; value: PurchaseStatus }[] => {
  const transitions: Record<string, { label: string; value: PurchaseStatus }[]> = {
    'draft': [{ label: '提交审批', value: 'pending_approval' }],
    'pending_approval': [
      { label: '审批通过', value: 'approved' },
      { label: '审批拒绝', value: 'rejected' }
    ],
    'approved': [{ label: '下单采购', value: 'ordered' }],
    'ordered': [{ label: '确认入库', value: 'received' }],
    'received': [{ label: '对账确认', value: 'reconciled' }],
    'reconciled': [{ label: '确认付款', value: 'paid' }],
    'rejected': [{ label: '重新提交', value: 'pending_approval' }]
  }
  return transitions[currentStatus] || []
}

const resetStatusForm = () => {
  statusPurchaseId.value = null
  statusCurrentStatus.value = 'draft'
  statusForm.value = {
    status: '' as PurchaseStatus,
    operator: '',
    remark: ''
  }
  statusFormRef.value?.resetFields()
}

const submitStatusChange = async () => {
  if (!statusFormRef.value || !statusPurchaseId.value) return

  try {
    await statusFormRef.value.validate()
  } catch {
    return
  }

  statusSubmitting.value = true
  try {
    const res = await api.updatePurchaseStatus(statusPurchaseId.value, {
      status: statusForm.value.status,
      remark: statusForm.value.remark,
      operator: statusForm.value.operator
    })
    if (res.data.code === 0) {
      ElMessage.success('状态更新成功')
      statusDialogVisible.value = false
      resetStatusForm()
      fetchWorkflowPurchases()
    } else {
      ElMessage.error(res.data.message || '状态更新失败')
    }
  } catch (error: any) {
    console.error('状态更新失败:', error)
    ElMessage.error(error.response?.data?.message || '状态更新失败')
  } finally {
    statusSubmitting.value = false
  }
}

const fetchTrendData = async (startDate?: string, endDate?: string) => {
  analysisLoading.value = true
  try {
    const res = await api.getPurchaseTrend({ start_date: startDate, end_date: endDate })
    if (res.data.code === 0) {
      trendData.value = res.data.data || []
    } else {
      trendData.value = generateMockTrendData()
    }
  } catch (error) {
    console.error('获取趋势数据失败:', error)
    trendData.value = generateMockTrendData()
  } finally {
    analysisLoading.value = false
  }
}

const fetchSupplierRatio = async (startDate?: string, endDate?: string) => {
  analysisLoading.value = true
  try {
    const res = await api.getSupplierRatio({ start_date: startDate, end_date: endDate })
    if (res.data.code === 0) {
      supplierRatioData.value = res.data.data || []
    } else {
      supplierRatioData.value = generateMockSupplierRatio()
    }
  } catch (error) {
    console.error('获取供应商占比失败:', error)
    supplierRatioData.value = generateMockSupplierRatio()
  } finally {
    analysisLoading.value = false
  }
}

const fetchPriceTrend = async (ingredientId?: number) => {
  analysisLoading.value = true
  try {
    const res = await api.getPriceTrend({ ingredient_id: ingredientId })
    if (res.data.code === 0) {
      priceTrendData.value = res.data.data || []
    } else {
      priceTrendData.value = generateMockPriceTrend()
    }
  } catch (error) {
    console.error('获取价格走势失败:', error)
    priceTrendData.value = generateMockPriceTrend()
  } finally {
    analysisLoading.value = false
  }
}

const generateMockTrendData = () => {
  const data = []
  for (let i = 5; i >= 0; i--) {
    const date = new Date()
    date.setMonth(date.getMonth() - i)
    data.push({
      date: `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`,
      amount: 30000 + Math.random() * 40000
    })
  }
  return data
}

const generateMockSupplierRatio = () => {
  return [
    { name: '上海名酒汇贸易有限公司', value: 85000, percentage: 45.2 },
    { name: '广州烈酒进出口公司', value: 65000, percentage: 34.6 },
    { name: '北京佳饮配送中心', value: 28000, percentage: 14.9 },
    { name: '深圳调酒辅料批发', value: 10000, percentage: 5.3 }
  ]
}

const generateMockPriceTrend = () => {
  const data: PriceTrendItem[] = []
  for (let i = 11; i >= 0; i--) {
    const date = new Date()
    date.setMonth(date.getMonth() - i)
    data.push({
      date: `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`,
      ingredient_id: 1,
      ingredient_name: '麦卡伦12年',
      avg_price: 3750 + Math.random() * 100,
      supplier_id: 1,
      supplier_name: '上海名酒汇贸易有限公司'
    })
  }
  return data
}

const handleTabChange = (tabName: string) => {
  switch (tabName) {
    case 'forecast':
      generateForecast()
      fetchSpirits()
      fetchIngredients()
      fetchSuppliers()
      break
    case 'supplier':
      fetchSuppliers()
      fetchSupplierQuotes()
      fetchSupplierEvaluations()
      fetchSpirits()
      fetchIngredients()
      break
    case 'workflow':
      fetchWorkflowPurchases()
      break
    case 'analysis':
      fetchTrendData()
      fetchSupplierRatio()
      fetchPriceTrend()
      fetchSpirits()
      fetchIngredients()
      break
  }
}

onMounted(() => {
  fetchPurchases()
  fetchSpirits()
  fetchIngredients()
  fetchSuppliers()
})

const { ShoppingCart, Money, TrendCharts, Search, Refresh, Plus, List, Delete, Star, User, Phone, Calendar, Clock, Check, Close, Edit, View, ArrowUp, ArrowDown } = ElementPlusIconsVue

const formatNumberLocal = (num: number | undefined | null): string => {
  if (num === undefined || num === null || isNaN(num)) return '0.00'
  return Number(num).toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

const formatSafeLocal = (num: number | undefined | null, decimals: number = 2): string => {
  if (num === undefined || num === null || isNaN(num)) return '0.' + '0'.repeat(decimals)
  return Number(num).toFixed(decimals)
}

const safeStringLocal = (str: string | undefined | null, defaultValue: string = '-'): string => {
  if (str === undefined || str === null || str === '') return defaultValue
  return String(str)
}

const PurchaseLedger = defineComponent({
  name: 'PurchaseLedger',
  props: {
    purchases: { type: Array as () => Purchase[], required: true },
    filteredPurchases: { type: Array as () => Purchase[], required: true },
    paginatedPurchases: { type: Array as () => Purchase[], required: true },
    loading: { type: Boolean, required: true },
    dateRange: { type: Array as () => string[], required: true },
    supplierKeyword: { type: String, required: true },
    currentPage: { type: Number, required: true },
    pageSize: { type: Number, required: true },
    expandRowKeys: { type: Array as () => number[], required: true },
    totalCount: { type: Number, required: true },
    totalAmount: { type: Number, required: true },
    averageAmount: { type: Number, required: true }
  },
  emits: ['fetch', 'reset-filters', 'create', 'delete', 'expand-change', 'update:dateRange', 'update:supplierKeyword', 'update:currentPage', 'update:pageSize'],
  setup(props, { emit }) {
    return () => h('div', { class: 'ledger-container' }, [
      h('div', { class: 'summary-cards' }, [
        h('el-card', { class: 'glass-card summary-card' }, {
          default: () => h('div', { class: 'card-content' }, [
            h('div', { class: 'card-icon', style: 'background: linear-gradient(135deg, #d4af37, #c9a227);' }, [
              h('el-icon', { size: 28 }, () => h(ShoppingCart))
            ]),
            h('div', { class: 'card-info' }, [
              h('p', { class: 'card-label' }, '采购单数'),
              h('p', { class: 'card-value gold-text' }, props.totalCount),
              h('p', { class: 'card-trend trend-up' }, [
                h('el-icon', () => h(TrendCharts)),
                `共 ${props.filteredPurchases.length} 条记录`
              ])
            ])
          ])
        }),
        h('el-card', { class: 'glass-card summary-card' }, {
          default: () => h('div', { class: 'card-content' }, [
            h('div', { class: 'card-icon', style: 'background: linear-gradient(135deg, #27ae60, #219a52);' }, [
              h('el-icon', { size: 28 }, () => h(Money))
            ]),
            h('div', { class: 'card-info' }, [
              h('p', { class: 'card-label' }, '采购总金额'),
              h('p', { class: 'card-value', style: 'color: #27ae60;' }, '¥' + formatNumberLocal(props.totalAmount)),
              h('p', { class: 'card-trend trend-up' }, [
                h('el-icon', () => h(ShoppingCart)),
                `平均 ¥${formatNumberLocal(props.averageAmount)}/单`
              ])
            ])
          ])
        })
      ]),
      h('el-card', { class: 'glass-card' }, {
        default: () => [
          h('div', { class: 'filter-section' }, [
            h('div', { class: 'filter-left' }, [
              h('el-date-picker', {
                modelValue: props.dateRange,
                type: 'daterange',
                rangeSeparator: '至',
                startPlaceholder: '开始日期',
                endPlaceholder: '结束日期',
                valueFormat: 'YYYY-MM-DD',
                class: 'filter-item',
                'onUpdate:modelValue': (val: string[]) => emit('update:dateRange', val)
              }),
              h('el-input', {
                modelValue: props.supplierKeyword,
                placeholder: '供应商关键词',
                clearable: true,
                class: 'filter-item',
                'onUpdate:modelValue': (val: string) => emit('update:supplierKeyword', val),
                onKeyup: { enter: () => emit('fetch') }
              }, {
                prefix: () => h('el-icon', () => h(Search))
              }),
              h('el-button', { class: 'glow-button', onClick: () => emit('fetch') }, {
                default: () => [h('el-icon', () => h(Search)), ' 查询']
              }),
              h('el-button', { onClick: () => emit('reset-filters') }, {
                default: () => [h('el-icon', () => h(Refresh)), ' 重置']
              })
            ]),
            h('div', { class: 'filter-right' }, [
              h('el-button', { class: 'glow-button', onClick: () => emit('create') }, {
                default: () => [h('el-icon', () => h(Plus)), ' 新建采购']
              })
            ])
          ]),
          h('el-table', {
            data: props.paginatedPurchases,
            style: 'width: 100%',
            rowKey: 'id',
            vLoading: props.loading,
            expandRowKeys: props.expandRowKeys,
            onExpandChange: (row: Purchase, expandedRows: Purchase[]) => emit('expand-change', row, expandedRows)
          }, {
            default: () => [
              h('el-table-column', { type: 'expand', width: 50 }, {
                default: ({ row }: { row: Purchase }) => h('div', { class: 'purchase-items-container' }, [
                  h('div', { class: 'purchase-items-header' }, [
                    h('span', { class: 'items-title' }, [h('el-icon', () => h(List)), ' 采购明细']),
                    h('span', { class: 'items-count' }, `共 ${row.purchase_items?.length || 0} 项`)
                  ]),
                  h('el-table', { data: row.purchase_items || [], size: 'small', style: 'width: 100%' }, {
                    default: () => [
                      h('el-table-column', { prop: 'ingredient_type', label: '类型', width: 100, align: 'center' }, {
                        default: ({ row }: { row: PurchaseItem }) => h('el-tag', {
                          type: row.ingredient_type === 'spirit' ? 'warning' : 'success',
                          effect: 'dark',
                          size: 'small'
                        }, { default: () => row.ingredient_type === 'spirit' ? '烈酒' : '辅料' })
                      }),
                      h('el-table-column', { prop: 'ingredient_name', label: '原料名称', minWidth: 150 }),
                      h('el-table-column', { prop: 'quantity', label: '数量', width: 100, align: 'center' }),
                      h('el-table-column', { prop: 'unit', label: '单位', width: 80, align: 'center' }),
                      h('el-table-column', { prop: 'unit_price', label: '单价', width: 120, align: 'right' }, {
                        default: ({ row }: { row: PurchaseItem }) => '¥' + formatSafeLocal(row.unit_price)
                      }),
                      h('el-table-column', { prop: 'subtotal', label: '小计', width: 120, align: 'right' }, {
                        default: ({ row }: { row: PurchaseItem }) => h('span', { class: 'subtotal' }, '¥' + formatSafeLocal(row.subtotal))
                      }),
                      h('el-table-column', { prop: 'batch_no', label: '供应商批号', width: 130, align: 'center' }, {
                        default: ({ row }: { row: PurchaseItem }) => row.batch_no ? h('span', { class: 'batch-no' }, row.batch_no) : h('span', { class: 'text-muted' }, '-')
                      }),
                      h('el-table-column', { prop: 'expiry_date', label: '有效期', width: 110, align: 'center' }, {
                        default: ({ row }: { row: PurchaseItem }) => row.expiry_date ? h('span', { class: 'expiry-date' }, row.expiry_date) : h('span', { class: 'text-muted' }, '-')
                      }),
                      h('el-table-column', { prop: 'stock_batch_code', label: '系统批次码', width: 160, align: 'center' }, {
                        default: ({ row }: { row: PurchaseItem }) => row.stock_batch_code ? h('span', { class: 'batch-code' }, row.stock_batch_code) : h('span', { class: 'text-muted' }, '-')
                      })
                    ]
                  })
                ])
              }),
              h('el-table-column', { label: '采购单号', minWidth: 180 }, {
                default: ({ row }: { row: Purchase }) => h('span', { class: 'purchase-no' }, safeStringLocal(row.purchase_no))
              }),
              h('el-table-column', { label: '供应商', minWidth: 150 }, {
                default: ({ row }: { row: Purchase }) => h('span', { class: 'supplier-text' }, safeStringLocal(row.supplier))
              }),
              h('el-table-column', { label: '采购日期', width: 120, align: 'center' }, {
                default: ({ row }: { row: Purchase }) => safeStringLocal(row.purchase_date)
              }),
              h('el-table-column', { label: '总金额', width: 130, align: 'right' }, {
                default: ({ row }: { row: Purchase }) => h('span', { class: 'total-amount' }, '¥' + formatSafeLocal(row.total_amount))
              }),
              h('el-table-column', { label: '操作员', width: 100, align: 'center' }, {
                default: ({ row }: { row: Purchase }) => h('span', { class: 'operator-text' }, safeStringLocal(row.operator))
              }),
              h('el-table-column', { prop: 'remark', label: '备注', minWidth: 150 }, {
                default: ({ row }: { row: Purchase }) => row.remark ? h('span', { class: 'remark' }, row.remark) : h('span', { class: 'text-muted' }, '-')
              }),
              h('el-table-column', { label: '操作', width: 100, align: 'center', fixed: 'right' }, {
                default: ({ row }: { row: Purchase }) => h('el-button', {
                  type: 'danger',
                  link: true,
                  size: 'small',
                  onClick: () => emit('delete', row)
                }, { default: () => [h('el-icon', () => h(Delete)), ' 删除'] })
              })
            ]
          }),
          h('div', { class: 'pagination-section' }, [
            h('el-pagination', {
              modelValue: props.currentPage,
              pageSize: props.pageSize,
              pageSizes: [10, 20, 50, 100],
              total: props.filteredPurchases.length,
              layout: 'total, sizes, prev, pager, next, jumper',
              background: true,
              'onUpdate:current-page': (val: number) => emit('update:currentPage', val),
              'onUpdate:page-size': (val: number) => emit('update:pageSize', val)
            })
          ])
        ]
      })
    ])
  }
})

const PurchaseForecast = defineComponent({
  name: 'PurchaseForecast',
  props: {
    forecastItems: { type: Array as () => PurchaseForecastItem[], required: true },
    suggestions: { type: Array as () => PurchaseSuggestion[], required: true },
    loading: { type: Boolean, required: true },
    spirits: { type: Array as () => Spirit[], required: true },
    ingredients: { type: Array as () => Ingredient[], required: true },
    suppliers: { type: Array as () => Supplier[], required: true }
  },
  emits: ['generate-forecast', 'create-suggestion', 'confirm-suggestion'],
  setup(props, { emit }) {
    const selectedItems = ref<number[]>([])
    const forecastDays = ref(30)

    const highPriorityItems = computed(() => props.forecastItems.filter(i => i.priority === 'high'))
    const mediumPriorityItems = computed(() => props.forecastItems.filter(i => i.priority === 'medium'))
    const lowPriorityItems = computed(() => props.forecastItems.filter(i => i.priority === 'low'))

    const totalSuggestedAmount = computed(() => {
      return props.forecastItems.reduce((sum, item) => sum + item.suggested_quantity * item.unit_price, 0)
    })

    const getPriorityType = (priority: string) => {
      const types: Record<string, string> = { high: 'danger', medium: 'warning', low: 'success' }
      return types[priority] || 'info'
    }

    const getPriorityText = (priority: string) => {
      const texts: Record<string, string> = { high: '高', medium: '中', low: '低' }
      return texts[priority] || priority
    }

    const toggleSelectItem = (index: number) => {
      const idx = selectedItems.value.indexOf(index)
      if (idx > -1) {
        selectedItems.value.splice(idx, 1)
      } else {
        selectedItems.value.push(index)
      }
    }

    const selectAll = () => {
      if (selectedItems.value.length === props.forecastItems.length) {
        selectedItems.value = []
      } else {
        selectedItems.value = props.forecastItems.map((_, i) => i)
      }
    }

    const createSuggestionFromSelected = () => {
      if (selectedItems.value.length === 0) {
        ElMessage.warning('请先选择需要生成建议的原料')
        return
      }
      const items = selectedItems.value.map(i => props.forecastItems[i])
      emit('create-suggestion', items)
      selectedItems.value = []
    }

    return () => h('div', { class: 'forecast-container' }, [
      h('div', { class: 'forecast-header' }, [
        h('div', { class: 'header-left' }, [
          h('el-input-number', {
            modelValue: forecastDays.value,
            min: 7,
            max: 90,
            step: 7,
            label: '预测天数',
            'onUpdate:modelValue': (val: number) => forecastDays.value = val
          }),
          h('el-button', { class: 'glow-button', onClick: () => emit('generate-forecast', forecastDays.value), loading: props.loading }, {
            default: () => [h('el-icon', () => h(Refresh)), ' 生成预测']
          })
        ]),
        h('div', { class: 'header-right' }, [
          h('el-button', { onClick: selectAll }, { default: () => selectedItems.value.length === props.forecastItems.length ? '取消全选' : '全选' }),
          h('el-button', { class: 'glow-button', onClick: createSuggestionFromSelected, disabled: selectedItems.value.length === 0 }, {
            default: () => [h('el-icon', () => h(Plus)), ' 生成采购建议']
          })
        ])
      ]),

      h('div', { class: 'forecast-summary' }, [
        h('el-card', { class: 'glass-card summary-card' }, {
          default: () => h('div', { class: 'card-content' }, [
            h('div', { class: 'card-icon', style: 'background: linear-gradient(135deg, #e74c3c, #c0392b);' }, [
              h('el-icon', { size: 28 }, () => h(ArrowUp))
            ]),
            h('div', { class: 'card-info' }, [
              h('p', { class: 'card-label' }, '高优先级'),
              h('p', { class: 'card-value', style: 'color: #e74c3c;' }, highPriorityItems.value.length),
              h('p', { class: 'card-trend' }, '需要立即采购')
            ])
          ])
        }),
        h('el-card', { class: 'glass-card summary-card' }, {
          default: () => h('div', { class: 'card-content' }, [
            h('div', { class: 'card-icon', style: 'background: linear-gradient(135deg, #f39c12, #e67e22);' }, [
              h('el-icon', { size: 28 }, () => h(Clock))
            ]),
            h('div', { class: 'card-info' }, [
              h('p', { class: 'card-label' }, '中优先级'),
              h('p', { class: 'card-value', style: 'color: #f39c12;' }, mediumPriorityItems.value.length),
              h('p', { class: 'card-trend' }, '建议近期采购')
            ])
          ])
        }),
        h('el-card', { class: 'glass-card summary-card' }, {
          default: () => h('div', { class: 'card-content' }, [
            h('div', { class: 'card-icon', style: 'background: linear-gradient(135deg, #27ae60, #219a52);' }, [
              h('el-icon', { size: 28 }, () => h(ArrowDown))
            ]),
            h('div', { class: 'card-info' }, [
              h('p', { class: 'card-label' }, '低优先级'),
              h('p', { class: 'card-value', style: 'color: #27ae60;' }, lowPriorityItems.value.length),
              h('p', { class: 'card-trend' }, '库存充足')
            ])
          ])
        }),
        h('el-card', { class: 'glass-card summary-card' }, {
          default: () => h('div', { class: 'card-content' }, [
            h('div', { class: 'card-icon', style: 'background: linear-gradient(135deg, #d4af37, #c9a227);' }, [
              h('el-icon', { size: 28 }, () => h(Money))
            ]),
            h('div', { class: 'card-info' }, [
              h('p', { class: 'card-label' }, '预计总金额'),
              h('p', { class: 'card-value gold-text' }, '¥' + formatNumberLocal(totalSuggestedAmount.value)),
              h('p', { class: 'card-trend' }, `共 ${props.forecastItems.length} 项原料`)
            ])
          ])
        })
      ]),

      h('el-card', { class: 'glass-card' }, {
        default: () => [
          h('div', { class: 'section-header' }, [
            h('span', { class: 'section-title' }, [h('el-icon', () => h(List)), ' 采购需求预测']),
            h('span', { class: 'section-subtitle' }, `已选择 ${selectedItems.value.length} 项`)
          ]),
          h('el-table', {
            data: props.forecastItems,
            style: 'width: 100%',
            vLoading: props.loading,
            rowClassName: ({ row }: { row: PurchaseForecastItem }) => `priority-row-${row.priority}`
          }, {
            default: () => [
              h('el-table-column', { type: 'selection', width: 50, align: 'center' }, {
                default: ({ $index }: { $index: number }) => h('el-checkbox', {
                  modelValue: selectedItems.value.includes($index),
                  onChange: () => toggleSelectItem($index)
                })
              }),
              h('el-table-column', { label: '优先级', width: 100, align: 'center' }, {
                default: ({ row }: { row: PurchaseForecastItem }) => h('el-tag', {
                  type: getPriorityType(row.priority),
                  effect: 'dark',
                  size: 'small'
                }, { default: () => getPriorityText(row.priority) })
              }),
              h('el-table-column', { label: '原料类型', width: 100, align: 'center' }, {
                default: ({ row }: { row: PurchaseForecastItem }) => h('el-tag', {
                  type: row.ingredient_type === 'spirit' ? 'warning' : 'success',
                  effect: 'dark',
                  size: 'small'
                }, { default: () => row.ingredient_type === 'spirit' ? '烈酒' : '辅料' })
              }),
              h('el-table-column', { prop: 'ingredient_name', label: '原料名称', minWidth: 150 }),
              h('el-table-column', { label: '当前库存', width: 100, align: 'center' }, {
                default: ({ row }: { row: PurchaseForecastItem }) => `${row.current_stock} ${row.unit}`
              }),
              h('el-table-column', { label: '安全库存', width: 100, align: 'center' }, {
                default: ({ row }: { row: PurchaseForecastItem }) => `${row.safe_stock} ${row.unit}`
              }),
              h('el-table-column', { label: '销售预测', width: 100, align: 'center' }, {
                default: ({ row }: { row: PurchaseForecastItem }) => `${row.sales_forecast} ${row.unit}`
              }),
              h('el-table-column', { label: '即将过期', width: 100, align: 'center' }, {
                default: ({ row }: { row: PurchaseForecastItem }) => row.expiry_quantity > 0
                  ? h('span', { style: 'color: #e74c3c;' }, `${row.expiry_quantity} ${row.unit}`)
                  : h('span', { class: 'text-muted' }, '-')
              }),
              h('el-table-column', { label: '建议采购量', width: 120, align: 'center' }, {
                default: ({ row }: { row: PurchaseForecastItem }) => h('span', { class: 'suggested-qty' }, `${row.suggested_quantity} ${row.unit}`)
              }),
              h('el-table-column', { label: '经济批量', width: 100, align: 'center' }, {
                default: ({ row }: { row: PurchaseForecastItem }) => `${row.economic_order_qty} ${row.unit}`
              }),
              h('el-table-column', { label: '最小起订', width: 100, align: 'center' }, {
                default: ({ row }: { row: PurchaseForecastItem }) => `${row.min_order_qty} ${row.unit}`
              }),
              h('el-table-column', { label: '预计金额', width: 120, align: 'right' }, {
                default: ({ row }: { row: PurchaseForecastItem }) => h('span', { class: 'total-amount' }, '¥' + formatSafeLocal(row.suggested_quantity * row.unit_price))
              }),
              h('el-table-column', { label: '建议供应商', minWidth: 180 }, {
                default: ({ row }: { row: PurchaseForecastItem }) => h('div', { class: 'supplier-suggest' }, [
                  h('el-icon', () => h(User)),
                  h('span', row.suggested_supplier_name || '暂无')
                ])
              }),
              h('el-table-column', { label: '保质期提醒', width: 120, align: 'center' }, {
                default: ({ row }: { row: PurchaseForecastItem }) => {
                  if (row.days_to_expiry !== undefined && row.days_to_expiry <= 7) {
                    return h('el-tag', { type: 'danger', size: 'small' }, { default: () => `${row.days_to_expiry}天后过期` })
                  }
                  if (row.days_to_expiry !== undefined && row.days_to_expiry <= 15) {
                    return h('el-tag', { type: 'warning', size: 'small' }, { default: () => `${row.days_to_expiry}天后过期` })
                  }
                  return h('span', { class: 'text-muted' }, '-')
                }
              })
            ]
          })
        ]
      })
    ])
  }
})

const SupplierManagement = defineComponent({
  name: 'SupplierManagement',
  props: {
    suppliers: { type: Array as () => Supplier[], required: true },
    quotes: { type: Array as () => SupplierQuote[], required: true },
    evaluations: { type: Array as () => SupplierEvaluation[], required: true },
    loading: { type: Boolean, required: true },
    spirits: { type: Array as () => Spirit[], required: true },
    ingredients: { type: Array as () => Ingredient[], required: true }
  },
  emits: ['fetch-suppliers', 'create-supplier', 'update-supplier', 'delete-supplier', 'fetch-quotes', 'create-quote', 'delete-quote', 'fetch-evaluations', 'create-evaluation', 'compare-prices'],
  setup(props, { emit }) {
    const activeSubTab = ref('list')
    const supplierKeyword = ref('')
    const evaluationFilter = ref('')

    const filteredSuppliers = computed(() => {
      let result = [...props.suppliers]
      if (supplierKeyword.value) {
        const keyword = supplierKeyword.value.toLowerCase()
        result = result.filter(s => s.name.toLowerCase().includes(keyword)
          || s.contact_person.toLowerCase().includes(keyword)
          || s.phone.includes(keyword))
      }
      if (evaluationFilter.value) {
        result = result.filter(s => s.evaluation === evaluationFilter.value)
      }
      return result
    })

    const getEvaluationType = (evaluation: string) => {
      const types: Record<string, string> = { A: 'success', B: 'warning', C: 'info', D: 'danger' }
      return types[evaluation] || 'info'
    }

    const getEvaluationText = (evaluation: string) => {
      const texts: Record<string, string> = { A: '优秀', B: '良好', C: '一般', D: '较差' }
      return texts[evaluation] || evaluation
    }

    return () => h('div', { class: 'supplier-container' }, [
      h('div', { class: 'supplier-tabs' }, [
        h('el-tabs', { modelValue: activeSubTab.value, onTabChange: (val: string) => activeSubTab.value = val }, {
          default: () => [
            h('el-tab-pane', { label: '供应商档案', name: 'list' }),
            h('el-tab-pane', { label: '报价管理', name: 'quotes' }),
            h('el-tab-pane', { label: '供应商评价', name: 'evaluations' })
          ]
        })
      ]),

      activeSubTab.value === 'list' && h('div', [
        h('div', { class: 'filter-section' }, [
          h('div', { class: 'filter-left' }, [
            h('el-input', {
              modelValue: supplierKeyword.value,
              placeholder: '搜索供应商名称/联系人/电话',
              clearable: true,
              class: 'filter-item',
              style: 'width: 300px;',
              'onUpdate:modelValue': (val: string) => supplierKeyword.value = val
            }, { prefix: () => h('el-icon', () => h(Search)) }),
            h('el-select', {
              modelValue: evaluationFilter.value,
              placeholder: '评价等级',
              clearable: true,
              class: 'filter-item',
              style: 'width: 150px;',
              'onUpdate:modelValue': (val: string) => evaluationFilter.value = val
            }, {
              default: () => [
                h('el-option', { label: 'A级', value: 'A' }),
                h('el-option', { label: 'B级', value: 'B' }),
                h('el-option', { label: 'C级', value: 'C' }),
                h('el-option', { label: 'D级', value: 'D' })
              ]
            }),
            h('el-button', { class: 'glow-button', onClick: () => emit('fetch-suppliers') }, { default: () => [h('el-icon', () => h(Refresh)), ' 刷新'] })
          ]),
          h('div', { class: 'filter-right' }, [
            h('el-button', { class: 'glow-button', onClick: () => emit('create-supplier') }, { default: () => [h('el-icon', () => h(Plus)), ' 新建供应商'] })
          ])
        ]),

        h('el-card', { class: 'glass-card' }, {
          default: () => [
            h('el-table', { data: filteredSuppliers.value, style: 'width: 100%', vLoading: props.loading }, {
              default: () => [
                h('el-table-column', { label: '供应商名称', minWidth: 220 }, {
                  default: ({ row }: { row: Supplier }) => h('div', { class: 'supplier-name-cell' }, [
                    h('span', { class: 'supplier-name' }, row.name),
                    h('el-tag', { type: getEvaluationType(row.evaluation), size: 'small', style: 'margin-left: 8px;' }, { default: () => row.evaluation + '级' })
                  ])
                }),
                h('el-table-column', { label: '联系人', width: 120 }, {
                  default: ({ row }: { row: Supplier }) => h('div', { class: 'contact-info' }, [
                    h('el-icon', () => h(User)), h('span', row.contact_person)
                  ])
                }),
                h('el-table-column', { label: '联系电话', width: 140 }, {
                  default: ({ row }: { row: Supplier }) => h('div', { class: 'contact-info' }, [
                    h('el-icon', () => h(Phone)), h('span', row.phone)
                  ])
                }),
                h('el-table-column', { label: '账期', width: 100, align: 'center' }, {
                  default: ({ row }: { row: Supplier }) => `${row.account_period}天`
                }),
                h('el-table-column', { label: '最小起订额', width: 120, align: 'right' }, {
                  default: ({ row }: { row: Supplier }) => '¥' + formatNumberLocal(row.min_order_amount)
                }),
                h('el-table-column', { label: '配送天数', width: 100, align: 'center' }, {
                  default: ({ row }: { row: Supplier }) => `${row.delivery_days}天`
                }),
                h('el-table-column', { label: '综合评分', width: 120, align: 'center' }, {
                  default: ({ row }: { row: Supplier }) => h('el-rate', { modelValue: row.rating, disabled: true, showText: true, textColor: '#ffd700' })
                }),
                h('el-table-column', { label: '地址', minWidth: 200 }, { default: ({ row }: { row: Supplier }) => safeStringLocal(row.address) }),
                h('el-table-column', { label: '备注', minWidth: 150 }, { default: ({ row }: { row: Supplier }) => safeStringLocal(row.remark) }),
                h('el-table-column', { label: '操作', width: 200, align: 'center', fixed: 'right' }, {
                  default: ({ row }: { row: Supplier }) => h('div', { class: 'action-buttons' }, [
                    h('el-button', { type: 'primary', link: true, size: 'small', onClick: () => emit('update-supplier', row) }, {
                      default: () => [h('el-icon', () => h(Edit)), ' 编辑']
                    }),
                    h('el-button', { type: 'success', link: true, size: 'small', onClick: () => {
                      const ingredientList = [...props.spirits, ...props.ingredients.map(i => ({ ...i, ingredient_type: 'ingredient' }))]
                      const firstIngredient = ingredientList[0]
                      if (firstIngredient) {
                        emit('compare-prices', firstIngredient.ingredient_type || 'spirit', firstIngredient.id, firstIngredient.name)
                      }
                    } }, {
                      default: () => [h('el-icon', () => h(Money)), ' 比价']
                    }),
                    h('el-button', { type: 'danger', link: true, size: 'small', onClick: () => emit('delete-supplier', row) }, {
                      default: () => [h('el-icon', () => h(Delete)), ' 删除']
                    })
                  ])
                })
              ]
            })
          ]
        })
      ]),

      activeSubTab.value === 'quotes' && h('div', [
        h('div', { class: 'filter-section' }, [
          h('div', { class: 'filter-left' }, [
            h('el-button', { class: 'glow-button', onClick: () => emit('fetch-quotes') }, { default: () => [h('el-icon', () => h(Refresh)), ' 刷新'] })
          ]),
          h('div', { class: 'filter-right' }, [
            h('el-button', { class: 'glow-button', onClick: () => emit('create-quote') }, { default: () => [h('el-icon', () => h(Plus)), ' 新建报价'] })
          ])
        ]),
        h('el-card', { class: 'glass-card' }, {
          default: () => h('el-table', { data: props.quotes, style: 'width: 100%' }, {
            default: () => [
              h('el-table-column', { label: '供应商', minWidth: 200 }, { default: ({ row }: { row: SupplierQuote }) => row.supplier_name }),
              h('el-table-column', { label: '原料类型', width: 100, align: 'center' }, {
                default: ({ row }: { row: SupplierQuote }) => h('el-tag', {
                  type: row.ingredient_type === 'spirit' ? 'warning' : 'success',
                  effect: 'dark',
                  size: 'small'
                }, { default: () => row.ingredient_type === 'spirit' ? '烈酒' : '辅料' })
              }),
              h('el-table-column', { label: '原料名称', minWidth: 150 }, { default: ({ row }: { row: SupplierQuote }) => row.ingredient_name }),
              h('el-table-column', { label: '单价', width: 120, align: 'right' }, {
                default: ({ row }: { row: SupplierQuote }) => h('span', { class: 'price-value' }, '¥' + formatSafeLocal(row.unit_price))
              }),
              h('el-table-column', { label: '单位', width: 80, align: 'center' }, { default: ({ row }: { row: SupplierQuote }) => row.unit }),
              h('el-table-column', { label: '最小起订', width: 120, align: 'center' }, {
                default: ({ row }: { row: SupplierQuote }) => `${row.min_order_qty} ${row.unit}`
              }),
              h('el-table-column', { label: '报价有效期', width: 220, align: 'center' }, {
                default: ({ row }: { row: SupplierQuote }) => `${row.valid_from} ~ ${row.valid_to}`
              }),
              h('el-table-column', { label: '操作', width: 100, align: 'center' }, {
                default: ({ row }: { row: SupplierQuote }) => h('el-button', {
                  type: 'danger',
                  link: true,
                  size: 'small',
                  onClick: () => emit('delete-quote', row)
                }, { default: () => [h('el-icon', () => h(Delete)), ' 删除'] })
              })
            ]
          })
        })
      ]),

      activeSubTab.value === 'evaluations' && h('div', [
        h('div', { class: 'filter-section' }, [
          h('div', { class: 'filter-left' }, [
            h('el-button', { class: 'glow-button', onClick: () => emit('fetch-evaluations') }, { default: () => [h('el-icon', () => h(Refresh)), ' 刷新'] })
          ]),
          h('div', { class: 'filter-right' }, [
            h('el-button', { class: 'glow-button', onClick: () => emit('create-evaluation') }, { default: () => [h('el-icon', () => h(Plus)), ' 新建评价'] })
          ])
        ]),
        h('el-card', { class: 'glass-card' }, {
          default: () => h('el-table', { data: props.evaluations, style: 'width: 100%' }, {
            default: () => [
              h('el-table-column', { label: '供应商', minWidth: 200 }, {
              default: ({ row }: { row: SupplierEvaluation }) => {
                const supplier = props.suppliers.find(s => s.id === row.supplier_id)
                return supplier ? supplier.name : '未知供应商'
              }
            }),
              h('el-table-column', { label: '评价周期', width: 120, align: 'center' }, { default: ({ row }: { row: SupplierEvaluation }) => row.period }),
              h('el-table-column', { label: '按时交货率', width: 120, align: 'center' }, {
                default: ({ row }: { row: SupplierEvaluation }) => h('div', { class: 'rate-cell' }, [
                  h('div', { class: 'rate-value' }, `${row.on_time_delivery_rate.toFixed(1)}%`),
                  h('el-progress', {
                    percentage: row.on_time_delivery_rate,
                    showText: false,
                    strokeWidth: 6,
                    color: row.on_time_delivery_rate >= 95 ? '#27ae60' : row.on_time_delivery_rate >= 85 ? '#f39c12' : '#e74c3c'
                  })
                ])
              }),
              h('el-table-column', { label: '价格稳定性', width: 140, align: 'center' }, {
                default: ({ row }: { row: SupplierEvaluation }) => h('el-rate', { modelValue: row.price_stability_score, disabled: true, showText: false })
              }),
              h('el-table-column', { label: '产品质量', width: 140, align: 'center' }, {
                default: ({ row }: { row: SupplierEvaluation }) => h('el-rate', { modelValue: row.quality_score, disabled: true, showText: false })
              }),
              h('el-table-column', { label: '综合评分', width: 120, align: 'center' }, {
                default: ({ row }: { row: SupplierEvaluation }) => h('span', { class: 'overall-score' }, row.overall_score.toFixed(1))
              }),
              h('el-table-column', { label: '总订单数', width: 100, align: 'center' }, { default: ({ row }: { row: SupplierEvaluation }) => row.total_orders }),
              h('el-table-column', { label: '延迟订单', width: 100, align: 'center' }, {
                default: ({ row }: { row: SupplierEvaluation }) => row.delayed_orders > 0
                  ? h('span', { style: 'color: #e74c3c;' }, row.delayed_orders)
                  : h('span', row.delayed_orders)
              }),
              h('el-table-column', { label: '投诉次数', width: 100, align: 'center' }, {
                default: ({ row }: { row: SupplierEvaluation }) => row.complaint_count > 0
                  ? h('span', { style: 'color: #e74c3c;' }, row.complaint_count)
                  : h('span', row.complaint_count)
              }),
              h('el-table-column', { label: '评价备注', minWidth: 200 }, { default: ({ row }: { row: SupplierEvaluation }) => safeStringLocal(row.remark) })
            ]
          })
        })
      ])
    ])
  }
})

const PurchaseWorkflow = defineComponent({
  name: 'PurchaseWorkflow',
  props: {
    purchases: { type: Array as () => PurchaseWithStatus[], required: true },
    loading: { type: Boolean, required: true }
  },
  emits: ['fetch', 'update-status'],
  setup(props, { emit }) {
    const statusFilter = ref('')

    const statusOptions = [
      { value: '', label: '全部状态' },
      { value: 'draft', label: '草稿' },
      { value: 'pending_approval', label: '待审批' },
      { value: 'approved', label: '已审批' },
      { value: 'rejected', label: '已拒绝' },
      { value: 'ordered', label: '已下单' },
      { value: 'received', label: '已入库' },
      { value: 'reconciled', label: '已对账' },
      { value: 'paid', label: '已付款' }
    ]

    const getStatusText = (status: PurchaseStatus): string => {
      const texts: Record<string, string> = {
        'draft': '草稿',
        'pending_approval': '待审批',
        'approved': '已审批',
        'rejected': '已拒绝',
        'ordered': '已下单',
        'received': '已入库',
        'reconciled': '已对账',
        'paid': '已付款'
      }
      return texts[status] || status
    }

    const getStatusType = (status: PurchaseStatus): string => {
      const types: Record<string, string> = {
        'draft': 'info',
        'pending_approval': 'warning',
        'approved': 'success',
        'rejected': 'danger',
        'ordered': 'primary',
        'received': 'success',
        'reconciled': 'success',
        'paid': 'success'
      }
      return types[status] || 'info'
    }

    const filteredPurchases = computed(() => {
      if (!statusFilter.value) return props.purchases
      return props.purchases.filter(p => p.status === statusFilter.value)
    })

    const workflowSteps = [
      { key: 'draft', label: '起草', icon: Edit },
      { key: 'pending_approval', label: '待审批', icon: Clock },
      { key: 'approved', label: '已审批', icon: Check },
      { key: 'ordered', label: '已下单', icon: ShoppingCart },
      { key: 'received', label: '已入库', icon: List },
      { key: 'reconciled', label: '已对账', icon: Money },
      { key: 'paid', label: '已付款', icon: Star }
    ]

    const getStepStatus = (purchaseStatus: string, stepKey: string) => {
      const statusOrder = ['draft', 'pending_approval', 'approved', 'ordered', 'received', 'reconciled', 'paid']
      const purchaseIndex = statusOrder.indexOf(purchaseStatus)
      const stepIndex = statusOrder.indexOf(stepKey)

      if (stepIndex < purchaseIndex || (stepIndex === purchaseIndex && purchaseStatus !== 'rejected')) {
        return 'finish'
      }
      if (stepIndex === purchaseIndex) {
        return 'process'
      }
      return 'wait'
    }

    return () => h('div', { class: 'workflow-container' }, [
      h('div', { class: 'workflow-header' }, [
        h('div', { class: 'header-left' }, [
          h('el-select', {
            modelValue: statusFilter.value,
            placeholder: '选择状态筛选',
            style: 'width: 180px;',
            'onUpdate:modelValue': (val: string) => statusFilter.value = val
          }, {
            default: () => statusOptions.map(opt => h('el-option', { key: opt.value, label: opt.label, value: opt.value }))
          }),
          h('el-button', { class: 'glow-button', onClick: () => emit('fetch'), loading: props.loading }, {
            default: () => [h('el-icon', () => h(Refresh)), ' 刷新']
          })
        ])
      ]),

      h('div', { class: 'workflow-steps-overview' }, [
        h('div', { class: 'steps-container' }, [
          workflowSteps.map((step, index) => h('div', { key: step.key, class: 'step-overview-item' }, [
            h('div', { class: 'step-overview-icon' }, [h('el-icon', { size: 20 }, () => h(step.icon))]),
            h('span', { class: 'step-overview-label' }, step.label),
            index < workflowSteps.length - 1 && h('div', { class: 'step-overview-line' })
          ]))
        ])
      ]),

      h('el-card', { class: 'glass-card' }, {
        default: () => h('el-table', { data: filteredPurchases.value, style: 'width: 100%', vLoading: props.loading }, {
          default: () => [
            h('el-table-column', { label: '采购单号', minWidth: 180 }, {
              default: ({ row }: { row: PurchaseWithStatus }) => h('span', { class: 'purchase-no' }, safeStringLocal(row.purchase_no))
            }),
            h('el-table-column', { label: '供应商', minWidth: 200 }, {
              default: ({ row }: { row: PurchaseWithStatus }) => safeStringLocal(row.supplier)
            }),
            h('el-table-column', { label: '采购日期', width: 120, align: 'center' }, {
              default: ({ row }: { row: PurchaseWithStatus }) => safeStringLocal(row.purchase_date)
            }),
            h('el-table-column', { label: '总金额', width: 130, align: 'right' }, {
              default: ({ row }: { row: PurchaseWithStatus }) => h('span', { class: 'total-amount' }, '¥' + formatSafeLocal(row.total_amount))
            }),
            h('el-table-column', { label: '当前状态', width: 120, align: 'center' }, {
              default: ({ row }: { row: PurchaseWithStatus }) => h('el-tag', {
                type: getStatusType(row.status),
                effect: 'dark',
                size: 'small'
              }, { default: () => getStatusText(row.status) })
            }),
            h('el-table-column', { label: '审批人', width: 120, align: 'center' }, {
              default: ({ row }: { row: PurchaseWithStatus }) => safeStringLocal(row.approval_by)
            }),
            h('el-table-column', { label: '入库人', width: 120, align: 'center' }, {
              default: ({ row }: { row: PurchaseWithStatus }) => safeStringLocal(row.received_by)
            }),
            h('el-table-column', { label: '流程进度', minWidth: 400 }, {
              default: ({ row }: { row: PurchaseWithStatus }) => h('div', { class: 'workflow-timeline' },
                workflowSteps.map((step, index) => h('div', {
                  key: step.key,
                  class: `timeline-item timeline-${getStepStatus(row.status, step.key)}`
                }, [
                  h('div', { class: 'timeline-dot' }, [
                    h('el-icon', { size: 14 }, () => h(step.icon))
                  ]),
                  index < workflowSteps.length - 1 && h('div', {
                    class: `timeline-line timeline-line-${getStepStatus(row.status, workflowSteps[index + 1].key)}`
                  })
                ]))
              )
            }),
            h('el-table-column', { label: '操作', width: 150, align: 'center', fixed: 'right' }, {
              default: ({ row }: { row: PurchaseWithStatus }) => {
                if (row.status === 'paid') {
                  return h('el-tag', { type: 'success', size: 'small' }, { default: () => '已完成' })
                }
                return h('el-button', {
                  type: 'primary',
                  link: true,
                  size: 'small',
                  onClick: () => emit('update-status', row)
                }, { default: () => [h('el-icon', () => h(Edit)), ' 更新状态'] })
              }
            })
          ]
        })
      })
    ])
  }
})

const PurchaseAnalysisReport = defineComponent({
  name: 'PurchaseAnalysisReport',
  props: {
    trendData: { type: Array as () => { date: string; amount: number }[], required: true },
    supplierRatio: { type: Array as () => { name: string; value: number; percentage: number }[], required: true },
    priceTrend: { type: Array as () => PriceTrendItem[], required: true },
    loading: { type: Boolean, required: true },
    spirits: { type: Array as () => Spirit[], required: true },
    ingredients: { type: Array as () => Ingredient[], required: true }
  },
  emits: ['fetch-trend', 'fetch-ratio', 'fetch-price-trend'],
  setup(props, { emit }) {
    const trendChartRef = ref<HTMLElement>()
    const ratioChartRef = ref<HTMLElement>()
    const priceChartRef = ref<HTMLElement>()
    const selectedIngredientId = ref<number>()
    const dateRange = ref<string[]>([])

    let trendChart: echarts.ECharts | null = null
    let ratioChart: echarts.ECharts | null = null
    let priceChart: echarts.ECharts | null = null

    const allIngredients = computed(() => {
      return [
        ...props.spirits.map(s => ({ id: s.id, name: s.name, type: 'spirit' })),
        ...props.ingredients.map(i => ({ id: i.id, name: i.name, type: 'ingredient' }))
      ]
    })

    const initTrendChart = () => {
      if (!trendChartRef.value) return
      trendChart = echarts.init(trendChartRef.value)
      updateTrendChart()
    }

    const updateTrendChart = () => {
      if (!trendChart) return
      const option = {
        tooltip: {
          trigger: 'axis',
          formatter: (params: any) => {
            const data = params[0]
            return `${data.name}<br/>采购金额: ¥${formatNumberLocal(data.value)}`
          }
        },
        grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
        xAxis: {
          type: 'category',
          boundaryGap: false,
          data: props.trendData.map(d => d.date),
          axisLine: { lineStyle: { color: 'rgba(255,255,255,0.2)' } },
          axisLabel: { color: 'rgba(255,255,255,0.7)' }
        },
        yAxis: {
          type: 'value',
          axisLine: { lineStyle: { color: 'rgba(255,255,255,0.2)' } },
          axisLabel: {
            color: 'rgba(255,255,255,0.7)',
            formatter: (value: number) => '¥' + (value / 1000) + 'k'
          },
          splitLine: { lineStyle: { color: 'rgba(255,255,255,0.05)' } }
        },
        series: [{
          name: '采购金额',
          type: 'line',
          smooth: true,
          data: props.trendData.map(d => d.amount),
          lineStyle: { color: '#d4af37', width: 3 },
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(212,175,55,0.4)' },
              { offset: 1, color: 'rgba(212,175,55,0.05)' }
            ])
          },
          itemStyle: { color: '#d4af37' }
        }]
      }
      trendChart.setOption(option)
    }

    const initRatioChart = () => {
      if (!ratioChartRef.value) return
      ratioChart = echarts.init(ratioChartRef.value)
      updateRatioChart()
    }

    const updateRatioChart = () => {
      if (!ratioChart) return
      const colors = ['#d4af37', '#3498db', '#27ae60', '#e74c3c', '#f39c12', '#9b59b6']
      const option = {
        tooltip: {
          trigger: 'item',
          formatter: '{b}: ¥{c} ({d}%)'
        },
        legend: {
          orient: 'vertical',
          right: '5%',
          top: 'center',
          textStyle: { color: 'rgba(255,255,255,0.7)' }
        },
        series: [{
          name: '供应商占比',
          type: 'pie',
          radius: ['40%', '70%'],
          center: ['35%', '50%'],
          avoidLabelOverlap: false,
          itemStyle: { borderRadius: 8, borderColor: 'rgba(0,0,0,0.3)', borderWidth: 2 },
          label: { show: false, position: 'center' },
          emphasis: {
            label: { show: true, fontSize: 16, fontWeight: 'bold', color: '#d4af37' }
          },
          labelLine: { show: false },
          data: props.supplierRatio.map((d, i) => ({
            value: d.value,
            name: d.name,
            itemStyle: { color: colors[i % colors.length] }
          }))
        }]
      }
      ratioChart.setOption(option)
    }

    const initPriceChart = () => {
      if (!priceChartRef.value) return
      priceChart = echarts.init(priceChartRef.value)
      updatePriceChart()
    }

    const updatePriceChart = () => {
      if (!priceChart || props.priceTrend.length === 0) return
      const groupedData: Record<string, { supplier: string; prices: number[] }[]> = {}
      props.priceTrend.forEach(item => {
        if (!groupedData[item.date]) groupedData[item.date] = []
        groupedData[item.date].push({ supplier: item.supplier_name, prices: [item.avg_price] })
      })

      const dates = [...new Set(props.priceTrend.map(d => d.date))].sort()
      const suppliers = [...new Set(props.priceTrend.map(d => d.supplier_name))]
      const colors = ['#d4af37', '#3498db', '#27ae60', '#e74c3c']

      const series = suppliers.map((supplier, idx) => ({
        name: supplier,
        type: 'line',
        smooth: true,
        data: dates.map(date => {
          const item = props.priceTrend.find(d => d.date === date && d.supplier_name === supplier)
          return item ? item.avg_price : null
        }),
        lineStyle: { color: colors[idx % colors.length], width: 2 },
        itemStyle: { color: colors[idx % colors.length] }
      }))

      const option = {
        tooltip: {
          trigger: 'axis',
          formatter: (params: any) => {
            let result = params[0].name + '<br/>'
            params.forEach((p: any) => {
              if (p.value) result += `${p.marker} ${p.seriesName}: ¥${p.value.toFixed(2)}<br/>`
            })
            return result
          }
        },
        legend: {
          data: suppliers,
          textStyle: { color: 'rgba(255,255,255,0.7)' },
          bottom: 0
        },
        grid: { left: '3%', right: '4%', bottom: '15%', top: '10%', containLabel: true },
        xAxis: {
          type: 'category',
          data: dates,
          axisLine: { lineStyle: { color: 'rgba(255,255,255,0.2)' } },
          axisLabel: { color: 'rgba(255,255,255,0.7)' }
        },
        yAxis: {
          type: 'value',
          axisLine: { lineStyle: { color: 'rgba(255,255,255,0.2)' } },
          axisLabel: { color: 'rgba(255,255,255,0.7)', formatter: '¥{value}' },
          splitLine: { lineStyle: { color: 'rgba(255,255,255,0.05)' } }
        },
        series
      }
      priceChart.setOption(option)
    }

    watch(() => props.trendData, () => updateTrendChart(), { deep: true })
    watch(() => props.supplierRatio, () => updateRatioChart(), { deep: true })
    watch(() => props.priceTrend, () => updatePriceChart(), { deep: true })

    const handleResize = () => {
      trendChart?.resize()
      ratioChart?.resize()
      priceChart?.resize()
    }

    onMounted(() => {
      nextTick(() => {
        initTrendChart()
        initRatioChart()
        initPriceChart()
        window.addEventListener('resize', handleResize)
      })
    })

    const handleIngredientChange = () => {
      emit('fetch-price-trend', selectedIngredientId.value)
    }

    const handleDateRangeChange = () => {
      if (dateRange.value && dateRange.value.length === 2) {
        emit('fetch-trend', dateRange.value[0], dateRange.value[1])
        emit('fetch-ratio', dateRange.value[0], dateRange.value[1])
      }
    }

    return () => h('div', { class: 'analysis-container' }, [
      h('div', { class: 'analysis-header' }, [
        h('div', { class: 'header-left' }, [
          h('el-date-picker', {
            modelValue: dateRange.value,
            type: 'daterange',
            rangeSeparator: '至',
            startPlaceholder: '开始日期',
            endPlaceholder: '结束日期',
            valueFormat: 'YYYY-MM-DD',
            'onUpdate:modelValue': (val: string[]) => { dateRange.value = val; handleDateRangeChange() }
          }),
          h('el-select', {
            modelValue: selectedIngredientId.value,
            placeholder: '选择原料查看价格走势',
            filterable: true,
            style: 'width: 250px;',
            'onUpdate:modelValue': (val: number) => { selectedIngredientId.value = val; handleIngredientChange() }
          }, {
            default: () => allIngredients.value.map(item => h('el-option', {
              key: item.id,
              label: item.name,
              value: item.id
            }))
          })
        ])
      ]),

      h('div', { class: 'analysis-grid' }, [
        h('el-card', { class: 'glass-card chart-card' }, {
          default: () => [
            h('div', { class: 'chart-header' }, [
              h('h3', { class: 'chart-title' }, [h('el-icon', () => h(TrendCharts)), ' 采购趋势分析']),
              h('el-button', { size: 'small', onClick: () => emit('fetch-trend') }, {
                default: () => [h('el-icon', () => h(Refresh)), ' 刷新']
              })
            ]),
            h('div', { ref: trendChartRef, class: 'chart-container trend-chart' })
          ]
        }),

        h('el-card', { class: 'glass-card chart-card' }, {
          default: () => [
            h('div', { class: 'chart-header' }, [
              h('h3', { class: 'chart-title' }, [h('el-icon', () => h(User)), ' 供应商占比']),
              h('el-button', { size: 'small', onClick: () => emit('fetch-ratio') }, {
                default: () => [h('el-icon', () => h(Refresh)), ' 刷新']
              })
            ]),
            h('div', { ref: ratioChartRef, class: 'chart-container ratio-chart' })
          ]
        }),

        h('el-card', { class: 'glass-card chart-card full-width' }, {
          default: () => [
            h('div', { class: 'chart-header' }, [
              h('h3', { class: 'chart-title' }, [h('el-icon', () => h(Money)), ' 原料价格走势']),
              h('el-button', { size: 'small', onClick: () => emit('fetch-price-trend', selectedIngredientId.value) }, {
                default: () => [h('el-icon', () => h(Refresh)), ' 刷新']
              })
            ]),
            h('div', { ref: priceChartRef, class: 'chart-container price-chart' })
          ]
        })
      ])
    ])
  }
})
</script>

<style lang="scss" scoped>
.purchases-container {
  padding: 20px;
  min-height: 100vh;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);

  .page-header {
    margin-bottom: 20px;

    .page-title {
      font-size: 28px;
      font-weight: bold;
      color: #d4af37;
      margin: 0;
      text-shadow: 0 2px 10px rgba(212, 175, 55, 0.3);
    }

    .page-subtitle {
      color: rgba(255, 255, 255, 0.6);
      margin-top: 5px 0 0 0;
      font-size: 14px;
    }
  }

  .main-tabs {
    :deep(.el-tabs__nav-wrap::after) {
      background: rgba(255, 255, 255, 0.1);
    }

    :deep(.el-tabs__item) {
      color: rgba(255, 255, 255, 0.7);
      font-size: 15px;
      padding: 0 24px;
      height: 48px;
      line-height: 48px;
      transition: all 0.3s ease;

      &.is-active {
        color: #d4af37;
        font-weight: 600;
      }

      &:hover:not(.is-active) {
        color: #f0c865;
      }
    }

    :deep(.el-tabs__active-bar) {
      background: linear-gradient(90deg, #d4af37, #f0c865);
      height: 3px;
      border-radius: 2px;
    }

    :deep(.el-tabs__content) {
      padding-top: 20px;
    }
  }

  .glass-card {
    background: rgba(255, 255, 255, 0.05);
    backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 12px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
    transition: all 0.3s ease;

    &:hover {
      box-shadow: 0 12px 40px rgba(0, 0, 0, 0.3);
    }

    :deep(.el-card__header) {
      border-bottom: 1px solid rgba(255, 255, 255, 0.08);
      background: rgba(255, 255, 255, 0.02);
    }
  }

  .glow-button {
    background: linear-gradient(135deg, #d4af37, #b8941f);
    border: none;
    color: #1a1a2e;
    font-weight: 600;
    transition: all 0.3s ease;

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 5px 20px rgba(212, 175, 55, 0.4);
    }

    &:active {
      transform: translateY(0);
    }
  }

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;

    .header-left {
      display: flex;
      gap: 12px;
      align-items: center;
    }
  }

  .action-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    flex-wrap: wrap;
    gap: 12px;

    .filter-row {
      display: flex;
      gap: 12px;
      flex-wrap: wrap;
      align-items: center;
    }

    .action-buttons {
      display: flex;
      gap: 12px;
    }
  }

  .gold-text {
    color: #d4af37;
    font-weight: 600;
  }

  .total-amount {
    color: #d4af37;
    font-weight: 600;
    font-size: 14px;
  }

  .purchase-no {
    font-family: 'Courier New', monospace;
    color: rgba(255, 255, 255, 0.9);
    font-weight: 500;
  }

  .contact-info {
    display: flex;
    align-items: center;
    gap: 6px;

    .el-icon {
      color: #d4af37;
    }
  }

  .priority-high {
    background: linear-gradient(90deg, rgba(231, 76, 60, 0.15), transparent);
    border-left: 4px solid #e74c3c;
  }

  .priority-medium {
    background: linear-gradient(90deg, rgba(243, 156, 18, 0.15), transparent);
    border-left: 4px solid #f39c12;
  }

  .priority-low {
    background: linear-gradient(90deg, rgba(39, 174, 96, 0.15), transparent);
    border-left: 4px solid #27ae60;
  }

  :deep(.el-table) {
    background: transparent;
    color: rgba(255, 255, 255, 0.85);

    th.el-table__cell {
      background: rgba(255, 255, 255, 0.05) !important;
      color: #d4af37 !important;
      font-weight: 600;
    }

    td.el-table__cell {
      background: transparent;
      border-bottom: 1px solid rgba(255, 255, 255, 0.05);
    }

    .el-table__row:hover > td {
      background: rgba(255, 255, 255, 0.03);
    }

    &::before {
      background-color: rgba(255, 255, 255, 0.1);
    }
  }

  :deep(.el-input__wrapper) {
    background: rgba(255, 255, 255, 0.05);
    box-shadow: 0 0 0 1px rgba(255, 255, 255, 0.1);

    &.is-focus {
      box-shadow: 0 0 0 2px rgba(212, 175, 55, 0.3);
    }

    .el-input__inner {
      color: rgba(255, 255, 255, 0.9);

      &::placeholder {
        color: rgba(255, 255, 255, 0.4);
      }
    }
  }

  :deep(.el-select) {
    .el-input__wrapper {
      background: rgba(255, 255, 255, 0.05);
    }

    .el-select__placeholder {
      color: rgba(255, 255, 255, 0.4);
    }
  }

  :deep(.el-date-editor) {
    .el-input__wrapper {
      background: rgba(255, 255, 255, 0.05);
    }
  }

  :deep(.el-pagination) {
    margin-top: 16px;

    .el-pager li {
      background: rgba(255, 255, 255, 0.05);
      color: rgba(255, 255, 255, 0.8);

      &.is-active {
        background: #d4af37;
        color: #1a1a2e;
      }

      &:hover:not(.is-active) {
        color: #d4af37;
      }
    }

    .btn-prev, .btn-next {
      background: rgba(255, 255, 255, 0.05);
      color: rgba(255, 255, 255, 0.8);

      &:hover {
        color: #d4af37;
      }
    }

    .el-pagination__total,
    .el-pagination__jump {
      color: rgba(255, 255, 255, 0.7);
    }
  }

  :deep(.el-dialog) {
    background: linear-gradient(135deg, #1a1a2e, #16213e);
    border: 1px solid rgba(212, 175, 55, 0.3);
    border-radius: 12px;

    .el-dialog__header {
      border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    }

    .el-dialog__title {
      color: #d4af37;
    }

    .el-dialog__body {
      color: rgba(255, 255, 255, 0.9);
    }

    .el-form-item__label {
      color: rgba(255, 255, 255, 0.9);
    }

    .el-input__inner {
      color: rgba(255, 255, 255, 0.9);
    }
  }

  :deep(.el-drawer) {
    background: linear-gradient(180deg, #1a1a2e, #16213e);

    .el-drawer__header {
      border-bottom: 1px solid rgba(255, 255, 255, 0.1);

      .el-drawer__title {
        color: #d4af37;
      }
    }
  }

  .forecast-stats {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 16px;
    margin-bottom: 20px;

    .stat-card {
      padding: 20px;
      border-radius: 12px;
      display: flex;
      align-items: center;
      gap: 16px;
      transition: transform 0.3s ease;

      &:hover {
        transform: translateY(-3px);
      }

      &.high {
        background: linear-gradient(135deg, rgba(231, 76, 60, 0.2), rgba(231, 76, 60, 0.05));
        border: 1px solid rgba(231, 76, 60, 0.3);
      }

      &.medium {
        background: linear-gradient(135deg, rgba(243, 156, 18, 0.2), rgba(243, 156, 18, 0.05));
        border: 1px solid rgba(243, 156, 18, 0.3);
      }

      &.low {
        background: linear-gradient(135deg, rgba(39, 174, 96, 0.2), rgba(39, 174, 96, 0.05));
        border: 1px solid rgba(39, 174, 96, 0.3);
      }

      .stat-icon {
        width: 56px;
        height: 56px;
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 28px;

        .high & {
          background: rgba(231, 76, 60);
        }

        .medium & {
          background: rgba(243, 156, 18);
        }

        .low & {
          background: rgba(39, 174, 96);
        }
      }

      .stat-content {
        flex: 1;

        .stat-value {
          font-size: 32px;
          font-weight: bold;
          color: white;
        }

        .stat-label {
          font-size: 14px;
          color: rgba(255, 255, 255, 0.7);
          margin-top: 4px;
        }
      }
    }
  }

  .supplier-tabs {
    :deep(.el-tabs__nav-wrap::after) {
      background: rgba(255, 255, 255, 0.1);
    }

    :deep(.el-tabs__item) {
      color: rgba(255, 255, 255, 0.7);

      &.is-active {
        color: #d4af37;
      }
    }

    :deep(.el-tabs__active-bar) {
      background: #d4af37;
    }
  }

  .supplier-rating {
    display: flex;
    align-items: center;
    gap: 8px;

    .el-rate {
      .el-rate__item {
        font-size: 16px;
      }
    }
  }

  .comparison-table {
    .best-price {
      color: #27ae60;
      font-weight: bold;
    }

    .price-diff {
      color: #e74c3c;
      font-size: 12px;
    }
  }

  .workflow-container {
    .workflow-header {
      margin-bottom: 20px;
    }

    .workflow-steps-overview {
      background: rgba(255, 255, 255, 0.03);
      border-radius: 12px;
      padding: 20px;
      margin-bottom: 20px;
      border: 1px solid rgba(255, 255, 255, 0.08);

      .steps-container {
        display: flex;
        justify-content: space-between;
        align-items: center;
        position: relative;
      }

      .step-overview-item {
        display: flex;
        flex-direction: column;
        align-items: center;
        position: relative;
        flex: 1;

        &:last-child {
          flex: 0 0 auto;
        }

        .step-overview-icon {
          width: 48px;
          height: 48px;
          border-radius: 50%;
          background: linear-gradient(135deg, #d4af37, #b8941f);
          display: flex;
          align-items: center;
          justify-content: center;
          color: #1a1a2e;
          margin-bottom: 8px;
          box-shadow: 0 4px 15px rgba(212, 175, 55, 0.4);
        }

        .step-overview-label {
          color: rgba(255, 255, 255, 0.8);
          font-size: 13px;
          font-weight: 500;
        }

        .step-overview-line {
          position: absolute;
          top: 24px;
          left: 50%;
          width: 100%;
          height: 2px;
          background: linear-gradient(90deg, #d4af37, rgba(212, 175, 55, 0.3));
          z-index: 0;
        }
      }
    }

    .workflow-timeline {
      display: flex;
      align-items: center;
      gap: 4px;
      padding: 8px 0;

      .timeline-item {
        display: flex;
        align-items: center;
        position: relative;

        .timeline-dot {
          width: 24px;
          height: 24px;
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          font-size: 12px;
          transition: all 0.3s ease;
        }

        &.timeline-finish .timeline-dot {
          background: #27ae60;
          color: white;
        }

        &.timeline-process .timeline-dot {
          background: #d4af37;
          color: #1a1a2e;
          animation: pulse 2s infinite;
        }

        &.timeline-wait .timeline-dot {
          background: rgba(255, 255, 255, 0.15);
          color: rgba(255, 255, 255, 0.5);
        }

        .timeline-line {
          width: 24px;
          height: 2px;
          margin-left: 4px;

          &.timeline-line-finish {
            background: #27ae60;
          }

          &.timeline-line-process {
            background: linear-gradient(90deg, #d4af37, rgba(255, 255, 255, 0.15));
          }

          &.timeline-line-wait {
            background: rgba(255, 255, 255, 0.15);
          }
        }
      }
    }
  }

  @keyframes pulse {
    0%, 100% {
      box-shadow: 0 0 0 0 rgba(212, 175, 55, 0.4);
    }
    50% {
      box-shadow: 0 0 0 8px rgba(212, 175, 55, 0);
    }
  }

  .analysis-container {
    .analysis-header {
      margin-bottom: 20px;

      .header-left {
        display: flex;
        gap: 12px;
        align-items: center;
      }
    }

    .analysis-grid {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 20px;

      .chart-card {
        .chart-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: 16px;

          .chart-title {
            font-size: 16px;
            font-weight: 600;
            color: #d4af37;
            margin: 0;
            display: flex;
            align-items: center;
            gap: 8px;
          }
        }

        .chart-container {
          width: 100%;
        }

        .trend-chart {
          height: 320px;
        }

        .ratio-chart {
          height: 320px;
        }

        .price-chart {
          height: 360px;
        }

        &.full-width {
          grid-column: 1 / -1;
        }
      }
    }
  }

  .ledger-container {
    .summary-section {
      margin-bottom: 20px;

      .ledger-header {
        h3 {
          font-size: 18px;
          font-weight: 600;
          color: #d4af37;
          margin: 0 0 16px 0;
          display: flex;
          align-items: center;
          gap: 8px;
        }

        .summary-grid {
          display: grid;
          grid-template-columns: repeat(4, 1fr);
          gap: 16px;

          .summary-item {
            text-align: center;
            padding: 16px;
            background: rgba(255, 255, 255, 0.03);
            border-radius: 8px;
            border: 1px solid rgba(255, 255, 255, 0.08);

            .summary-value {
              font-size: 24px;
              font-weight: bold;
              color: #d4af37;
            }

            .summary-label {
              font-size: 13px;
              color: rgba(255, 255, 255, 0.6);
              margin-top: 4px;
            }
          }
        }
      }
    }

    .detail-section {
      padding: 16px 0;

      .detail-row {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 16px;
        margin-bottom: 12px;

        .detail-item {
          display: flex;

          .detail-label {
            width: 120px;
            color: rgba(255, 255, 255, 0.6);
          }

          .detail-value {
            flex: 1;
            color: rgba(255, 255, 255, 0.9);
          }
        }
      }
    }
  }

  .empty-state {
    text-align: center;
    padding: 60px 20px;
    color: rgba(255, 255, 255, 0.5);

    .empty-icon {
      font-size: 48px;
      margin-bottom: 16px;
      opacity: 0.5;
    }

    .empty-text {
      font-size: 16px;
    }
  }

  @media (max-width: 768px) {
    .forecast-stats {
      grid-template-columns: 1fr;
    }

    .analysis-grid {
      grid-template-columns: 1fr;
    }

    .ledger-container .summary-section .summary-grid {
      grid-template-columns: 1fr 1fr;
    }

    .detail-row {
      grid-template-columns: 1fr !important;
    }
  }
}
</style>

