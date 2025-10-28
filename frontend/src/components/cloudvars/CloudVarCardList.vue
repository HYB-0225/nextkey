<template>
  <div class="cloudvar-card-list" ref="cardListRef">
    <div 
      v-for="cloudvar in cloudVars" 
      :key="cloudvar.id" 
      class="cloudvar-card"
    >
      <div class="card-header">
        <div class="card-title">
          <el-checkbox 
            v-model="selectedIds" 
            :label="cloudvar.id" 
            @click.stop
            class="card-checkbox"
          />
          <span class="var-key">{{ cloudvar.key }}</span>
        </div>
      </div>
      
      <div class="card-body">
        <div class="var-value">
          <span class="label">值:</span>
          <div class="value-content">{{ cloudvar.value }}</div>
        </div>
      </div>
      
      <div class="card-footer">
        <el-button size="small" type="primary" @click="$emit('edit', cloudvar)">
          <el-icon><Edit /></el-icon>
          编辑
        </el-button>
        <el-button size="small" type="danger" @click="$emit('delete', cloudvar)">
          <el-icon><Delete /></el-icon>
          删除
        </el-button>
      </div>
    </div>
    
    <div v-if="loading" class="loading-state">
      <el-icon class="is-loading"><Loading /></el-icon>
      <span>加载中...</span>
    </div>
    
    <div v-if="!loading && cloudVars.length === 0" class="empty-state">
      <el-icon><Cloudy /></el-icon>
      <span>暂无云变量</span>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { Edit, Delete, Loading, Cloudy } from '@element-plus/icons-vue'
import { staggerScaleIn } from '@/utils/animations'

const props = defineProps({
  cloudVars: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  },
  selectedVars: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['selection-change', 'edit', 'delete'])

const selectedIds = computed({
  get: () => props.selectedVars.map(v => v.id),
  set: (ids) => {
    const selected = props.cloudVars.filter(v => ids.includes(v.id))
    emit('selection-change', selected)
  }
})

// 动画相关
const cardListRef = ref(null)

onMounted(() => {
  animateCards()
})

watch(() => props.cloudVars, () => {
  if (props.cloudVars.length > 0) {
    setTimeout(() => {
      animateCards()
    }, 50)
  }
})

const animateCards = () => {
  if (cardListRef.value) {
    const cards = cardListRef.value.querySelectorAll('.cloudvar-card')
    if (cards.length > 0) {
      staggerScaleIn(cards, {
        stagger: 0.05,
        ease: 'back.out(1.7)',
      })
    }
  }
}
</script>

<style scoped>
.cloudvar-card-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.cloudvar-card {
  background: #fff;
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border-light);
  overflow: hidden;
  transition: all var(--duration-fast) var(--ease-out);
  box-shadow: var(--shadow-sm);
}

.cloudvar-card:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
}

.cloudvar-card:active {
  transform: scale(0.98);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: var(--color-bg-secondary);
  border-bottom: 1px solid var(--color-border-light);
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  min-width: 0;
}

.card-checkbox {
  flex-shrink: 0;
}

.var-key {
  font-family: 'Courier New', monospace;
  font-size: 14px;
  font-weight: 600;
  color: var(--color-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.card-body {
  padding: 12px 16px;
}

.var-value {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.label {
  color: var(--color-text-secondary);
  font-weight: 500;
  font-size: 12px;
}

.value-content {
  color: var(--color-text-primary);
  font-size: 13px;
  background: var(--color-bg-secondary);
  padding: 8px 12px;
  border-radius: var(--radius-md);
  word-break: break-all;
  line-height: 1.5;
}

.card-footer {
  display: flex;
  gap: 8px;
  padding: 12px 16px;
  background: var(--color-bg-secondary);
  border-top: 1px solid var(--color-border-light);
}

.card-footer :deep(.el-button) {
  flex: 1;
  font-size: 12px;
}

.loading-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  color: var(--color-text-secondary);
  gap: 8px;
}

.loading-state :deep(.el-icon),
.empty-state :deep(.el-icon) {
  font-size: 32px;
}

:deep(.el-checkbox__label) {
  display: none;
}
</style>

