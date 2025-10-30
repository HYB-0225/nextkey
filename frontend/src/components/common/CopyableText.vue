<template>
  <span 
    class="copyable-text"
    :class="{ 'is-copying': isCopying }"
    :title="isMobile ? '长按复制' : '点击复制'"
    @click="handleClick"
    @touchstart="handleTouchStart"
    @touchend="handleTouchEnd"
    @touchmove="handleTouchMove"
  >
    <slot>{{ text }}</slot>
  </span>
</template>

<script setup>
import { ref } from 'vue'
import { copyToClipboard } from '@/utils/copy'
import { useResponsive } from '@/composables/useResponsive'

const { isMobile } = useResponsive()

const props = defineProps({
  text: {
    type: String,
    required: true
  },
  successMessage: {
    type: String,
    default: '复制成功'
  }
})

const isCopying = ref(false)
let longPressTimer = null
let touchMoved = false

const handleCopy = async () => {
  isCopying.value = true
  await copyToClipboard(props.text, props.successMessage)
  
  setTimeout(() => {
    isCopying.value = false
  }, 300)
}

const handleClick = (e) => {
  if (!isMobile.value) {
    e.stopPropagation()
    handleCopy()
  }
}

const handleTouchStart = (e) => {
  if (isMobile.value) {
    touchMoved = false
    longPressTimer = setTimeout(() => {
      if (!touchMoved) {
        e.preventDefault()
        e.stopPropagation()
        handleCopy()
      }
    }, 600)
  }
}

const handleTouchEnd = () => {
  if (longPressTimer) {
    clearTimeout(longPressTimer)
    longPressTimer = null
  }
}

const handleTouchMove = () => {
  touchMoved = true
  if (longPressTimer) {
    clearTimeout(longPressTimer)
    longPressTimer = null
  }
}
</script>

<style scoped>
.copyable-text {
  cursor: pointer;
  transition: all 0.2s ease;
  user-select: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  position: relative;
}

.copyable-text:hover {
  color: var(--el-color-primary);
  text-decoration: underline;
  text-decoration-style: dotted;
}

.copyable-text.is-copying {
  animation: copy-flash 0.3s ease;
}

@keyframes copy-flash {
  0% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
    transform: scale(0.98);
  }
  100% {
    opacity: 1;
    transform: scale(1);
  }
}

/* 移动端样式 */
@media (max-width: 768px) {
  .copyable-text:hover {
    text-decoration: none;
  }
  
  .copyable-text:active {
    color: var(--el-color-primary);
  }
}
</style>

