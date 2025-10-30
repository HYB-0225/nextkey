import { ref, computed, watch } from 'vue'
import { unitValueToSeconds, secondsToUnitValue, formatDuration } from './useDuration'

export function useDurationEditor(cardData) {
  const isActivated = computed(() => cardData.value?.activated || false)
  
  // 有效时长模式（未激活卡）
  const durationValue = ref(30)
  const durationUnit = ref('day')
  
  // 到期时间模式（已激活卡）
  const expireTime = ref(null)
  
  // 实时预览
  const durationPreview = computed(() => {
    if (isActivated.value) {
      return null
    }
    const seconds = unitValueToSeconds(durationValue.value, durationUnit.value)
    return `${formatDuration(seconds)} = ${seconds.toLocaleString()}秒`
  })
  
  // card_type 自动推荐
  const getSuggestedCardType = (seconds) => {
    if (seconds === 0) return 'permanent'
    
    const rules = [
      { seconds: 604800, type: 'trial' },      // 7天
      { seconds: 2592000, type: 'month' },     // 30天
      { seconds: 7776000, type: 'quarter' },   // 90天
      { seconds: 31536000, type: 'year' }      // 365天
    ]
    
    for (const rule of rules) {
      if (seconds === rule.seconds) {
        return rule.type
      }
    }
    
    return null
  }
  
  const suggestedCardType = computed(() => {
    if (isActivated.value) {
      return null
    }
    const seconds = unitValueToSeconds(durationValue.value, durationUnit.value)
    return getSuggestedCardType(seconds)
  })
  
  // 从卡密数据初始化
  const initFromCard = (card) => {
    if (!card) return
    
    if (card.activated) {
      // 已激活：显示到期时间
      expireTime.value = card.expire_at ? new Date(card.expire_at) : null
    } else {
      // 未激活：显示有效时长
      const { value, unit } = secondsToUnitValue(card.duration || 0)
      durationValue.value = value
      durationUnit.value = unit
    }
  }
  
  // 获取提交数据
  const getSubmitData = () => {
    if (isActivated.value) {
      // 已激活：提交到期时间
      return {
        expire_at: expireTime.value ? expireTime.value.toISOString() : null
      }
    } else {
      // 未激活：提交时长
      return {
        duration: unitValueToSeconds(durationValue.value, durationUnit.value)
      }
    }
  }
  
  return {
    isActivated,
    durationValue,
    durationUnit,
    expireTime,
    durationPreview,
    suggestedCardType,
    initFromCard,
    getSubmitData
  }
}

