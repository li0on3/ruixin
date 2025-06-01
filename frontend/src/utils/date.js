import dayjs from 'dayjs'

// 格式化日期
export function formatDate(date, format = 'YYYY-MM-DD HH:mm:ss') {
  if (!date) return ''
  return dayjs(date).format(format)
}

// 格式化日期（只显示日期）
export function formatDateOnly(date) {
  return formatDate(date, 'YYYY-MM-DD')
}

// 格式化时间（只显示时间）
export function formatTimeOnly(date) {
  return formatDate(date, 'HH:mm:ss')
}

// 获取相对时间
export function fromNow(date) {
  if (!date) return ''
  
  const now = dayjs()
  const target = dayjs(date)
  const diff = now.diff(target, 'second')
  
  if (diff < 60) {
    return '刚刚'
  } else if (diff < 3600) {
    return `${Math.floor(diff / 60)}分钟前`
  } else if (diff < 86400) {
    return `${Math.floor(diff / 3600)}小时前`
  } else if (diff < 2592000) {
    return `${Math.floor(diff / 86400)}天前`
  } else {
    return formatDateOnly(date)
  }
}

// 获取日期范围
export function getDateRange(type) {
  let end = dayjs().endOf('day')
  let start
  
  switch (type) {
    case 'today':
      start = dayjs().startOf('day')
      break
    case 'yesterday':
      start = dayjs().subtract(1, 'day').startOf('day')
      end = dayjs().subtract(1, 'day').endOf('day')
      break
    case 'week':
      start = dayjs().subtract(6, 'day').startOf('day')
      break
    case 'month':
      start = dayjs().subtract(29, 'day').startOf('day')
      break
    case 'quarter':
      start = dayjs().subtract(89, 'day').startOf('day')
      break
    default:
      start = dayjs().startOf('day')
  }
  
  return [start.toDate(), end.toDate()]
}