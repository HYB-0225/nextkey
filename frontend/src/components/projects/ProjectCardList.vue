<template>
  <div class="project-card-list" ref="cardListRef">
    <div 
      v-for="project in projects" 
      :key="project.id" 
      class="project-card"
      @click="handleCardClick(project)"
    >
      <div class="card-header">
        <div class="card-title">
          <el-checkbox 
            v-model="selectedIds" 
            :label="project.id" 
            @click.stop
            class="card-checkbox"
          />
          <span class="project-name">{{ project.name }}</span>
        </div>
        <el-tag :type="project.mode === 'paid' ? 'success' : 'info'" size="small">
          {{ project.mode === 'paid' ? '付费' : '免费' }}
        </el-tag>
      </div>
      
      <div class="card-body">
        <div class="info-row">
          <span class="label">UUID</span>
          <span class="value text-ellipsis">
            <CopyableText :text="project.uuid" success-message="UUID已复制">
              {{ project.uuid }}
            </CopyableText>
          </span>
        </div>
        <div class="info-row">
          <span class="label">版本</span>
          <span class="value">{{ project.version }}</span>
        </div>
        <div class="info-row">
          <span class="label">配置</span>
          <div class="tags">
            <el-tag v-if="project.enable_hwid" size="small" type="info">机器码</el-tag>
            <el-tag v-if="project.enable_ip" size="small" type="info">IP验证</el-tag>
          </div>
        </div>
      </div>
      
      <div class="card-footer">
        <el-button size="small" @click.stop="$emit('edit', project)">
          <el-icon><Edit /></el-icon>
          编辑
        </el-button>
        <el-button size="small" type="primary" @click.stop="$emit('view-cards', project)">
          <el-icon><Ticket /></el-icon>
          卡密
        </el-button>
        <el-button size="small" type="success" @click.stop="$emit('view-vars', project)">
          <el-icon><Cloudy /></el-icon>
          变量
        </el-button>
        <el-button size="small" type="danger" @click.stop="$emit('delete', project)">
          <el-icon><Delete /></el-icon>
        </el-button>
      </div>
    </div>
    
    <div v-if="loading" class="loading-state">
      <el-icon class="is-loading"><Loading /></el-icon>
      <span>加载中...</span>
    </div>
    
    <div v-if="!loading && projects.length === 0" class="empty-state">
      <el-icon><Box /></el-icon>
      <span>暂无项目</span>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { Edit, Ticket, Cloudy, Delete, Loading, Box } from '@element-plus/icons-vue'
import CopyableText from '@/components/common/CopyableText.vue'
import { staggerScaleIn } from '@/utils/animations'

const props = defineProps({
  projects: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  },
  selectedProjects: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['selection-change', 'edit', 'view-cards', 'view-vars', 'delete'])

const selectedIds = computed({
  get: () => props.selectedProjects.map(p => p.id),
  set: (ids) => {
    const selected = props.projects.filter(p => ids.includes(p.id))
    emit('selection-change', selected)
  }
})

const handleCardClick = (project) => {
  // 点击卡片主体区域时不做任何操作,只有点击按钮才触发
}

// 动画相关
const cardListRef = ref(null)

onMounted(() => {
  animateCards()
})

watch(() => props.projects, () => {
  if (props.projects.length > 0) {
    // 数据更新时重新动画
    setTimeout(() => {
      animateCards()
    }, 50)
  }
})

const animateCards = () => {
  if (cardListRef.value) {
    const cards = cardListRef.value.querySelectorAll('.project-card')
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
.project-card-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.project-card {
  background: #fff;
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border-light);
  overflow: hidden;
  transition: all var(--duration-fast) var(--ease-out);
  box-shadow: var(--shadow-sm);
  cursor: pointer;
}

.project-card:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
}

.project-card:active {
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

.project-name {
  font-weight: 600;
  font-size: 15px;
  color: var(--color-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
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
  align-items: center;
  font-size: 13px;
}

.label {
  color: var(--color-text-secondary);
  font-weight: 500;
  flex-shrink: 0;
  margin-right: 12px;
}

.value {
  color: var(--color-text-primary);
  text-align: right;
  flex: 1;
}

.text-ellipsis {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tags {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
  justify-content: flex-end;
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

.card-footer :deep(.el-button:last-child) {
  flex: 0;
  min-width: 40px;
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

