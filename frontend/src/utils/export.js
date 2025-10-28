import { formatDuration } from '@/composables/useDuration'

export function exportToJSON(data, filename = 'cards_export') {
  const jsonStr = JSON.stringify(data, null, 2)
  const blob = new Blob([jsonStr], { type: 'application/json' })
  downloadFile(blob, `${filename}.json`)
}

export function exportToTXT(data, filename = 'cards_export') {
  const lines = data.map(card => card.card_key)
  const text = lines.join('\n')
  const blob = new Blob([text], { type: 'text/plain;charset=utf-8' })
  downloadFile(blob, `${filename}.txt`)
}

export function exportToCSV(data, filename = 'cards_export') {
  const headers = [
    'ID',
    '卡密',
    '激活状态',
    '有效时长',
    '类型',
    '备注',
    '设备码上限',
    '已用设备码',
    'IP上限',
    '已用IP',
    '专属信息',
    '创建时间'
  ]
  
  const rows = data.map(card => [
    card.id,
    card.card_key,
    card.activated ? '已激活' : '未激活',
    formatDuration(card.duration),
    card.card_type,
    card.note || '',
    card.max_hwid === -1 ? '无限制' : card.max_hwid,
    card.hwid_list?.length || 0,
    card.max_ip === -1 ? '无限制' : card.max_ip,
    card.ip_list?.length || 0,
    card.custom_data || '',
    formatDateTime(card.created_at)
  ])
  
  const csvContent = [
    headers.join(','),
    ...rows.map(row => row.map(cell => `"${String(cell).replace(/"/g, '""')}"`).join(','))
  ].join('\n')
  
  const blob = new Blob(['\ufeff' + csvContent], { type: 'text/csv;charset=utf-8' })
  downloadFile(blob, `${filename}.csv`)
}

function downloadFile(blob, filename) {
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}

function formatDateTime(dateStr) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

