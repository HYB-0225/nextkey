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

