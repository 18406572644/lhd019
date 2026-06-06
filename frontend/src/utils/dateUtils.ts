export const getDateRangeByPeriod = (period?: string): { startDate: string; endDate: string } => {
  const now = new Date()
  let startDate = new Date()
  let endDate = new Date()

  switch (period) {
    case 'today':
      startDate = new Date(now.getFullYear(), now.getMonth(), now.getDate())
      endDate = new Date(now.getFullYear(), now.getMonth(), now.getDate())
      break
    case 'week':
      const dayOfWeek = now.getDay()
      const diffToMonday = dayOfWeek === 0 ? -6 : 1 - dayOfWeek
      startDate = new Date(now.getFullYear(), now.getMonth(), now.getDate() + diffToMonday)
      endDate = new Date(startDate.getTime() + 6 * 24 * 60 * 60 * 1000)
      break
    case 'month':
      startDate = new Date(now.getFullYear(), now.getMonth(), 1)
      endDate = new Date(now.getFullYear(), now.getMonth() + 1, 0)
      break
    case 'quarter':
      const quarter = Math.floor(now.getMonth() / 3)
      startDate = new Date(now.getFullYear(), quarter * 3, 1)
      endDate = new Date(now.getFullYear(), quarter * 3 + 3, 0)
      break
    case 'year':
      startDate = new Date(now.getFullYear(), 0, 1)
      endDate = new Date(now.getFullYear(), 11, 31)
      break
    default:
      startDate = new Date(now.getFullYear(), now.getMonth(), 1)
      endDate = new Date(now.getFullYear(), now.getMonth() + 1, 0)
  }

  return {
    startDate: formatDate(startDate),
    endDate: formatDate(endDate)
  }
}

export const formatDate = (date: Date): string => {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

export const getPeriodLabel = (period?: string): string => {
  const labels: Record<string, string> = {
    today: '今日',
    week: '本周',
    month: '本月',
    quarter: '本季度',
    year: '本年',
    custom: '自定义'
  }
  return labels[period || ''] || '本月'
}
