# 调试会话: finance-report-bugs

**状态**: [OPEN]
**创建时间**: 2026-06-06
**会话ID**: finance-report-bugs

## 问题描述

1. 营收趋势表数据未正常渲染展示
2. 成本分析表导出的 PDF 文件标题显示为乱码
3. 分类销售页面的分类名称与产品名称字段均显示乱码
4. 利润核算、支付对账页面接口报错提示 SQL 表不存在

## 可验证假设列表

| ID | 假设描述 | 状态 | 验证方式 |
|----|---------|------|---------|
| H1 | 数据库表 `operating_costs`、`reconciliation_logs` 未正确创建或初始化，导致 SQL 表不存在错误 | PENDING | 检查后端错误日志、SQL 初始化脚本执行情况 |
| H2 | PDF 导出时的中文字符编码问题，jspdf 默认不支持中文，导致标题乱码 | PENDING | 检查 pdfExport.ts 中 jsPDF 配置，添加中文字体支持 |
| H3 | 字段映射错误，前端使用的字段名（如 category_name）与后端返回的字段（如 category）不一致，导致显示乱码 | PENDING | 检查 API 响应数据结构与前端模板绑定的字段 |
| H4 | 营收趋势图数据结构不匹配，daily_data 字段缺失或格式错误，导致 ECharts 无法渲染 | PENDING | 检查 RevenueReport 接口返回的 daily_data 字段 |
| H5 | 前端 API 响应处理错误，`res.data` 与 `res.data.data` 层级混淆，导致数据为空 | PENDING | 检查 api/index.ts 中响应包装结构，对比实际 API 返回 |

## 日志收集

### 插桩点

| 位置 | 日志类型 |
|------|---------|
| api/index.ts - getRevenueReport | API 响应数据结构 |
| RevenueReport.vue - fetchData | 数据接收和处理流程 |
| pdfExport.ts - exportToPDF | PDF 导出时的标题编码 |
| CategorySales.vue - fetchData | 分类销售数据字段映射 |
| api/index.ts - getProfitReport | 利润核算 API 响应和错误 |

## 修复记录

| 问题 | 修复方案 | 状态 |
|------|---------|------|
| TBD | TBD | PENDING |
