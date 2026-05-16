export function formatPrice(priceStr) {
  const num = parseFloat(priceStr)
  if (isNaN(num)) return '0.00'
  return num.toFixed(2)
}

export function formatTime(ts) {
  if (!ts) return ''
  const d = new Date(ts * 1000)
  const pad = n => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

export function formatOrderStatus(type) {
  const map = { 1: '待支付', 2: '已支付', 3: '已取消' }
  return map[type] || '未知'
}

export function formatOrderStatusType(type) {
  const map = { 1: 'warning', 2: 'success', 3: 'info' }
  return map[type] || 'info'
}
