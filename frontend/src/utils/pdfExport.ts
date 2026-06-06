import html2canvas from 'html2canvas'
import jsPDF from 'jspdf'

export interface PDFExportConfig {
  title: string
  subtitle?: string
  filters?: Record<string, string>
  companyName?: string
}

export const exportToPDF = async (elementId: string, config: PDFExportConfig) => {
  const element = document.getElementById(elementId)
  if (!element) {
    throw new Error('Element not found')
  }

  const canvas = await html2canvas(element, {
    scale: 2,
    useCORS: true,
    logging: false,
    backgroundColor: '#1a1a2e'
  })

  const imgData = canvas.toDataURL('image/png')
  const pdf = new jsPDF('l', 'mm', 'a4')
  const pdfWidth = pdf.internal.pageSize.getWidth()
  const pdfHeight = pdf.internal.pageSize.getHeight()
  const imgWidth = canvas.width
  const imgHeight = canvas.height
  const ratio = Math.min(pdfWidth / imgWidth, pdfHeight / imgHeight)
  const imgX = (pdfWidth - imgWidth * ratio) / 2
  const imgY = 30

  pdf.setFillColor(26, 26, 46)
  pdf.rect(0, 0, pdfWidth, pdfHeight, 'F')

  pdf.setFontSize(18)
  pdf.setTextColor(212, 175, 55)
  pdf.text(config.companyName || '调酒台账', pdfWidth / 2, 15, { align: 'center' })

  pdf.setFontSize(14)
  pdf.setTextColor(245, 245, 245)
  pdf.text(config.title, pdfWidth / 2, 22, { align: 'center' })

  if (config.subtitle) {
    pdf.setFontSize(10)
    pdf.setTextColor(160, 160, 160)
    pdf.text(config.subtitle, pdfWidth / 2, 27, { align: 'center' })
  }

  if (config.filters) {
    const filterText = Object.entries(config.filters)
      .map(([key, value]) => `${key}: ${value}`)
      .join(' | ')
    if (filterText) {
      pdf.setFontSize(9)
      pdf.setTextColor(160, 160, 160)
      pdf.text(filterText, pdfWidth / 2, 30, { align: 'center' })
    }
  }

  pdf.addImage(imgData, 'PNG', imgX, imgY, imgWidth * ratio, imgHeight * ratio)

  const pageCount = pdf.internal.getNumberOfPages()
  for (let i = 1; i <= pageCount; i++) {
    pdf.setPage(i)
    pdf.setFontSize(8)
    pdf.setTextColor(100, 100, 100)
    pdf.text(
      `第 ${i} 页 / 共 ${pageCount} 页 | 生成时间: ${new Date().toLocaleString('zh-CN')}`,
      pdfWidth / 2,
      pdfHeight - 10,
      { align: 'center' }
    )
  }

  const fileName = `${config.title}_${new Date().toISOString().split('T')[0]}.pdf`
  pdf.save(fileName)
}

export const formatCurrency = (num: number): string => {
  return '¥' + num.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

export const formatNumber = (num: number): string => {
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + '万'
  }
  return num.toLocaleString('zh-CN', { minimumFractionDigits: 0, maximumFractionDigits: 2 })
}
