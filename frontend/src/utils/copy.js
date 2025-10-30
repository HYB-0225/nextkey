import { ElMessage } from 'element-plus'

/**
 * 复制文本到剪贴板
 * @param {string} text - 要复制的文本
 * @param {string} successMessage - 成功提示消息
 * @returns {Promise<boolean>} 是否复制成功
 */
export async function copyToClipboard(text, successMessage = '复制成功') {
  try {
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(text)
      ElMessage.success(successMessage)
      return true
    } else {
      // 降级方案：使用 document.execCommand (兼容不支持 clipboard API 的环境)
      const textArea = document.createElement('textarea')
      textArea.value = text
      textArea.style.position = 'fixed'
      textArea.style.left = '-999999px'
      textArea.style.top = '-999999px'
      document.body.appendChild(textArea)
      textArea.focus()
      textArea.select()
      
      const successful = document.execCommand('copy')
      textArea.remove()
      
      if (successful) {
        ElMessage.success(successMessage)
        return true
      } else {
        throw new Error('Copy command failed')
      }
    }
  } catch (err) {
    console.error('复制失败:', err)
    ElMessage.error('复制失败')
    return false
  }
}

