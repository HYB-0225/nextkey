<template>
  <div class="card-list-mobile" ref="cardListRef">
    <div v-if="cards.length > 0" class="select-all-bar">
      <el-checkbox 
        v-model="isAllSelected" 
        :indeterminate="isIndeterminate"
        @change="handleSelectAll"
      >
        全选当前页
      </el-checkbox>
      <span class="select-info">已选 {{ selectedIds.length }} / {{ cards.length }}</span>
    </div>
    
    <div 
      v-for="card in cards" 
      :key="card.id" 
      class="card-item"
    >
      <div class="card-header">
        <div class="card-title">
          <el-checkbox 
            v-model="selectedIds" 
            :label="card.id" 
            @click.stop
            class="card-checkbox"
          />
          <CopyableText :text="card.card_key" success-message="卡密已复制" class="card-key-wrapper">
            <span class="card-key">{{ card.card_key }}</span>
          </CopyableText>
        </div>
        <el-tag :type="card.activated ? 'success' : 'info'" size="small">
          {{ card.activated ? '已激活' : '未激活' }}
        </el-tag>
      </div>
      
      <div class="card-body">
        <div class="info-row">
          <span class="label">时长</span>
          <span class="value">{{ formatDuration(card.duration) }}</span>
        </div>
        <div class="info-row">
          <span class="label">类型</span>
          <span class="value">{{ card.card_type }}</span>
        </div>
        <div class="info-row">
          <span class="label">设备</span>
          <span class="value">{{ card.hwid_list?.length || 0 }} / {{ card.max_hwid === -1 ? '∞' : card.max_hwid }}</span>
        </div>
        <div class="info-row">
          <span class="label">IP</span>
          <span class="value">{{ card.ip_list?.length || 0 }} / {{ card.max_ip === -1 ? '∞' : card.max_ip }}</span>
        </div>
        <div v-if="card.note" class="info-row">
          <span class="label">备注</span>
          <span class="value note-text">{{ card.note }}</span>
        </div>
      </div>
      
      <div class="card-footer">
        <el-button size="small" @click="$emit('edit', card)">
          <el-icon><Edit /></el-icon>
          编辑
        </el-button>
        <el-button size="small" type="primary" @click="$emit('view', card)">
          <el-icon><View /></el-icon>
          详情
        </el-button>
        <el-button size="small" type="danger" @click="$emit('delete', card)">
          <el-icon><Delete /></el-icon>
        </el-button>
      </div>
    </div>
    
    <div v-if="loading" class="loading-state">
      <el-icon class="is-loading"><Loading /></el-icon>
      <span>加载中...</span>
    </div>
    
    <div v-if="!loading && cards.length === 0" class="empty-state">
      <el-icon><Ticket /></el-icon>
      <span>暂无卡密</span>
    </div>
    
    <el-pagination
      v-if="total > 0"
      class="mobile-pagination"
      :current-page="page"
      :page-size="pageSize"
      :total="total"
      :small="true"
      layout="prev, pager, next"
      @current-change="handlePageChange"
    />
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { Edit, View, Delete, Loading, Ticket } from '@element-plus/icons-vue'
import CopyableText from '@/components/common/CopyableText.vue'
import { formatDuration } from '@/composables/useDuration'
import { staggerScaleIn } from '@/utils/animations'

const props = defineProps({
  cards: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  },
  page: {
    type: Number,
    default: 1
  },
  pageSize: {
    type: Number,
    default: 20
  },
  total: {
    type: Number,
    default: 0
  },
  selectedCards: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['selection-change', 'page-change', 'edit', 'view', 'delete'])

const selectedIds = computed({
  get: () => props.selectedCards.map(c => c.id),
  set: (ids) => {
    const selected = props.cards.filter(c => ids.includes(c.id))
    emit('selection-change', selected)
  }
})

const isAllSelected = computed(() => {
  return props.cards.length > 0 && selectedIds.value.length === props.cards.length
})

const isIndeterminate = computed(() => {
  const selectedCount = selectedIds.value.length
  return selectedCount > 0 && selectedCount < props.cards.length
})

const handleSelectAll = (checked) => {
  if (checked) {
    emit('selection-change', [...props.cards])
  } else {
    emit('selection-change', [])
  }
}

const handlePageChange = (newPage) => {
  emit('page-change', newPage)
}

// 动画相关
const cardListRef = ref(null)

onMounted(() => {
  animateCards()
})

watch(() => props.cards, () => {
  if (props.cards.length > 0) {
    setTimeout(() => {
      animateCards()
    }, 50)
  }
})

const animateCards = () => {
  if (cardListRef.value) {
    const cards = cardListRef.value.querySelectorAll('.card-item')
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
.card-list-mobile {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.select-all-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: var(--color-bg-secondary);
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border-light);
  margin-bottom: 4px;
}

.select-info {
  font-size: 13px;
  color: var(--color-text-secondary);
  font-weight: 500;
}

.card-item {
  background: #fff;
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border-light);
  overflow: hidden;
  transition: all var(--duration-fast) var(--ease-out);
  box-shadow: var(--shadow-sm);
}

.card-item:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
}

.card-item:active {
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

.card-key-wrapper {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.card-key {
  font-family: 'Courier New', monospace;
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.card-body {
  padding: 12px 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  font-size: 13px;
  gap: 8px;
}

.label {
  color: var(--color-text-secondary);
  font-weight: 500;
  flex-shrink: 0;
}

.value {
  color: var(--color-text-primary);
  text-align: right;
}

.note-text {
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
}

.card-footer {
  display: flex;
  gap: 6px;
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

.mobile-pagination {
  margin-top: 16px;
  display: flex;
  justify-content: center;
}

:deep(.el-checkbox__label) {
  display: none;
}

:deep(.el-pagination.is-small .btn-prev),
:deep(.el-pagination.is-small .btn-next),
:deep(.el-pagination.is-small .el-pager li) {
  min-width: 28px;
  height: 28px;
  line-height: 28px;
}
</style>

