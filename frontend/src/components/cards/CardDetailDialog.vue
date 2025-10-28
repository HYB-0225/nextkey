<template>
  <el-dialog v-model="dialogVisible" title="卡密详情" width="700px" @close="handleClose">
    <el-descriptions :column="isMobile ? 1 : 2" border v-if="card">
      <el-descriptions-item label="卡密">{{ card.card_key }}</el-descriptions-item>
      <el-descriptions-item label="状态">
        <el-tag :type="card.activated ? 'success' : 'info'">
          {{ card.activated ? '已激活' : '未激活' }}
        </el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="时长">{{ formatDuration(card.duration) }}</el-descriptions-item>
      <el-descriptions-item label="类型">{{ card.card_type }}</el-descriptions-item>
      <el-descriptions-item label="设备码列表" :span="2">
        <el-tag v-for="hwid in card.hwid_list" :key="hwid" style="margin-right: 5px;">
          {{ hwid }}
        </el-tag>
        <span v-if="!card.hwid_list || card.hwid_list.length === 0">无</span>
      </el-descriptions-item>
      <el-descriptions-item label="IP列表" :span="2">
        <el-tag v-for="ip in card.ip_list" :key="ip" style="margin-right: 5px;">
          {{ ip }}
        </el-tag>
        <span v-if="!card.ip_list || card.ip_list.length === 0">无</span>
      </el-descriptions-item>
      <el-descriptions-item label="备注" :span="2">{{ card.note || '无' }}</el-descriptions-item>
      <el-descriptions-item label="专属信息" :span="2">{{ card.custom_data || '无' }}</el-descriptions-item>
    </el-descriptions>
    
    <template #footer>
      <el-button @click="handleClose">关闭</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'
import { formatDuration } from '@/composables/useDuration'
import { useResponsive } from '@/composables/useResponsive'

const { isMobile } = useResponsive()

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  card: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['update:visible'])

const dialogVisible = ref(false)

watch(() => props.visible, (val) => {
  dialogVisible.value = val
})

watch(dialogVisible, (val) => {
  emit('update:visible', val)
})

const handleClose = () => {
  dialogVisible.value = false
}
</script>

<style scoped>
:deep(.el-dialog) {
  border-radius: var(--radius-lg);
  overflow: hidden;
}

:deep(.el-dialog__header) {
  background: var(--color-bg-tertiary);
  padding: 20px 24px;
  border-bottom: 1px solid var(--color-border-light);
}

:deep(.el-dialog__body) {
  padding: 24px;
}

:deep(.el-tag) {
  border-radius: var(--radius-sm);
  font-weight: 500;
}

:deep(.el-button) {
  border-radius: var(--radius-md);
  font-weight: 500;
  transition: all var(--duration-fast) var(--ease-out);
}

@media (max-width: 768px) {
  :deep(.el-dialog) {
    width: 90% !important;
    margin: 20px auto;
  }
  
  :deep(.el-dialog__body) {
    padding: 16px;
  }
  
  :deep(.el-descriptions__label) {
    width: 100px !important;
  }
}
</style>

