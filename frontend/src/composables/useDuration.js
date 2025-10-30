export const TIME_UNITS = {
  second: { label: '秒', value: 1 },
  day: { label: '天', value: 86400 },
  week: { label: '周', value: 604800 },
  month: { label: '月', value: 2592000 },
  quarter: { label: '季', value: 7776000 },
  year: { label: '年', value: 31536000 }
}

export function formatDuration(seconds) {
  if (seconds === 0) return '永久'
  
  const units = [
    { name: '年', value: 31536000 },
    { name: '季', value: 7776000 },
    { name: '月', value: 2592000 },
    { name: '周', value: 604800 },
    { name: '天', value: 86400 }
  ]
  
  for (const unit of units) {
    if (seconds % unit.value === 0) {
      return `${seconds / unit.value}${unit.name}`
    }
  }
  
  return `${seconds}秒`
}

export function secondsToUnitValue(seconds) {
  if (seconds === 0) return { value: 0, unit: 'day' }
  
  const units = [
    { key: 'year', value: 31536000 },
    { key: 'quarter', value: 7776000 },
    { key: 'month', value: 2592000 },
    { key: 'week', value: 604800 },
    { key: 'day', value: 86400 }
  ]
  
  for (const unit of units) {
    if (seconds % unit.value === 0) {
      return { value: seconds / unit.value, unit: unit.key }
    }
  }
  
  return { value: seconds, unit: 'second' }
}

export function unitValueToSeconds(value, unit) {
  return value * (TIME_UNITS[unit]?.value || 1)
}

export function formatExpireTime(card) {
  // 未激活
  if (!card.activated) {
    return '未激活'
  }
  
  // 永久卡
  if (card.duration === 0) {
    return '永久'
  }
  
  // 没有过期时间
  if (!card.expire_at) {
    return '未激活'
  }
  
  const expireDate = new Date(card.expire_at)
  const now = new Date()
  const diffMs = expireDate - now
  
  // 格式化完整时间
  const dateStr = expireDate.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false
  }).replace(/\//g, '-')
  
  // 已过期
  if (card.is_expired || diffMs < 0) {
    return `已过期 (${dateStr})`
  }
  
  // 计算相对时间
  const diffSeconds = Math.floor(diffMs / 1000)
  const diffMinutes = Math.floor(diffSeconds / 60)
  const diffHours = Math.floor(diffMinutes / 60)
  const diffDays = Math.floor(diffHours / 24)
  
  let relativeTime = ''
  if (diffDays > 0) {
    relativeTime = `还剩${diffDays}天`
  } else if (diffHours > 0) {
    relativeTime = `还剩${diffHours}小时`
  } else if (diffMinutes > 0) {
    relativeTime = `还剩${diffMinutes}分钟`
  } else {
    relativeTime = `还剩${diffSeconds}秒`
  }
  
  return `${relativeTime} (${dateStr})`
}

