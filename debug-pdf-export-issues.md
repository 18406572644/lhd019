# 调试会话: pdf-export-issues

**状态**: [OPEN]
**创建时间**: 2026-06-06
**会话ID**: pdf-export-issues

## 问题描述

1. 导出的 PDF 文件整体样式错乱
2. 部分图表与列表数据完全未显示
3. 部分列表数据显示不完整、内容被强制截断
4. 导出的报表无法正常阅读与使用

## 可验证假设列表

| ID | 假设描述 | 状态 | 验证方式 |
|----|---------|------|---------|
| H1 | html2canvas 截图时未正确处理元素宽度，导致内容被截断或缩放错误 | PENDING | 检查 wrapper 宽度设置、html2canvas 的 windowWidth/windowHeight 参数 |
| H2 | ECharts 图表使用 canvas 渲染，html2canvas 无法正确截取动态渲染的 canvas 内容 | PENDING | 检查图表是否已完全渲染、是否需要等待 ECharts 渲染完成后再截图 |
| H3 | 缺少分页处理逻辑，长列表和多图表内容被强制截断在单页内 | PENDING | 检查 pdf.addImage 是否支持多页、是否需要按内容高度分页 |
| H4 | 导出时克隆的 DOM 元素丢失了 CSS 样式（如深色背景、金色文字、表格边框） | PENDING | 检查 clonedContent 是否保留了 computed style、是否需要内联样式 |
| H5 | 元素 `position: relative` 或 `z-index` 等 CSS 属性在克隆后失效，导致图表和表格重叠或错位 | PENDING | 检查 wrapper 和 clonedContent 的样式设置 |

## 日志收集

### 插桩点

| 位置 | 日志类型 |
|------|---------|
| pdfExport.ts - exportToPDF | wrapper 尺寸、canvas 尺寸、PDF 尺寸 |
| pdfExport.ts - exportToPDF | html2canvas 完成时间、canvas 数据 |
| 各报表页面 - handleExport | 导出前图表是否已初始化、数据是否存在 |

## 修复记录

| 问题 | 修复方案 | 状态 |
|------|---------|------|
| TBD | TBD | PENDING |
