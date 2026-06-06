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

  const wrapper = document.createElement('div')
  wrapper.style.position = 'absolute'
  wrapper.style.left = '-9999px'
  wrapper.style.top = '0'
  wrapper.style.width = element.offsetWidth + 'px'
  wrapper.style.backgroundColor = '#1a1a2e'
  wrapper.style.padding = '20px'
  wrapper.style.fontFamily = '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif'

  const header = document.createElement('div')
  header.style.textAlign = 'center'
  header.style.marginBottom = '20px'
  header.style.paddingBottom = '15px'
  header.style.borderBottom = '1px solid rgba(212, 175, 55, 0.3)'

  const companyEl = document.createElement('div')
  companyEl.style.fontSize = '18px'
  companyEl.style.color = '#d4af37'
  companyEl.style.fontWeight = 'bold'
  companyEl.style.marginBottom = '8px'
  companyEl.textContent = config.companyName || '调酒台账'

  const titleEl = document.createElement('div')
  titleEl.style.fontSize = '16px'
  titleEl.style.color = '#f5f5f5'
  titleEl.style.fontWeight = '600'
  titleEl.style.marginBottom = '6px'
  titleEl.textContent = config.title

  const subtitleEl = document.createElement('div')
  subtitleEl.style.fontSize = '12px'
  subtitleEl.style.color = '#a0a0a0'
  subtitleEl.style.marginBottom = '4px'
  subtitleEl.textContent = config.subtitle || ''

  header.appendChild(companyEl)
  header.appendChild(titleEl)
  if (config.subtitle) {
    header.appendChild(subtitleEl)
  }

  if (config.filters) {
    const filterText = Object.entries(config.filters)
      .map(([key, value]) => `${key}: ${value}`)
      .join(' | ')
    if (filterText) {
      const filterEl = document.createElement('div')
      filterEl.style.fontSize = '10px'
      filterEl.style.color = '#a0a0a0'
      filterEl.textContent = filterText
      header.appendChild(filterEl)
    }
  }

  const clonedContent = element.cloneNode(true) as HTMLElement
  clonedContent.style.position = 'relative'
  clonedContent.style.left = 'auto'
  clonedContent.style.top = 'auto'

  wrapper.appendChild(header)
  wrapper.appendChild(clonedContent)
  document.body.appendChild(wrapper)

  try {
    const canvas = await html2canvas(wrapper, {
      scale: 2,
      useCORS: true,
      logging: false,
      backgroundColor: '#1a1a2e',
      windowWidth: wrapper.offsetWidth,
      windowHeight: wrapper.offsetHeight
    })

    const imgData = canvas.toDataURL('image/png')
    const pdf = new jsPDF('l', 'mm', 'a4')
    const pdfWidth = pdf.internal.pageSize.getWidth()
    const pdfHeight = pdf.internal.pageSize.getHeight()
    const imgWidth = canvas.width
    const imgHeight = canvas.height
    const ratio = Math.min(pdfWidth / imgWidth, pdfHeight / imgHeight)
    const imgX = (pdfWidth - imgWidth * ratio) / 2
    const imgY = 10

    pdf.setFillColor(26, 26, 46)
    pdf.rect(0, 0, pdfWidth, pdfHeight, 'F')

    pdf.addImage(imgData, 'PNG', imgX, imgY, imgWidth * ratio, imgHeight * ratio)

    const pageCount = pdf.internal.getNumberOfPages()
    for (let i = 1; i <= pageCount; i++) {
      pdf.setPage(i)
      pdf.setFontSize(8)
      pdf.setTextColor(100, 100, 100)
      const now = new Date()
      const year = now.getFullYear()
      const month = String(now.getMonth() + 1).padStart(2, '0')
      const day = String(now.getDate()).padStart(2, '0')
      const hours = String(now.getHours()).padStart(2, '0')
      const minutes = String(now.getMinutes()).padStart(2, '0')
      const seconds = String(now.getSeconds()).padStart(2, '0')
      const timeStr = `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
      pdf.text(
        `Page ${i} / ${pageCount} | Generated: ${timeStr}`,
        pdfWidth / 2,
        pdfHeight - 10,
        { align: 'center' }
      )
    }

    const safeTitle = config.title.replace(/[<>:"/\\|?*\u0000-\u001F]/g, '_')
    const fileName = `${safeTitle}_${new Date().toISOString().split('T')[0]}.pdf`
    pdf.save(fileName)
  } finally {
    document.body.removeChild(wrapper)
  }
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
