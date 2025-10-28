<template>
  <el-dialog v-model="dialogVisible" title="批量创建项目" width="700px" @close="handleClose">
    <el-alert type="info" :closable="false" style="margin-bottom: 15px;">
      请输入JSON格式的项目数据,格式: [{"name": "项目名", "mode": "free", ...}, ...]
    </el-alert>
    <el-input
      v-model="jsonData"
      type="textarea"
      :rows="12"
      placeholder='[
  {
    "name": "项目1",
    "mode": "free",
    "version": "1.0.0",
    "token_expire": 3600,
    "enable_hwid": true,
    "enable_ip": true,
    "description": "描述"
  }
]'
    />
    
    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" @click="handleSave">确定创建</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:visible', 'save'])

const dialogVisible = ref(false)
const jsonData = ref('')

watch(() => props.visible, (val) => {
  dialogVisible.value = val
  if (val) {
    jsonData.value = ''
  }
})

watch(dialogVisible, (val) => {
  emit('update:visible', val)
})

const handleClose = () => {
  dialogVisible.value = false
}

const handleSave = () => {
  emit('save', jsonData.value)
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

:deep(.el-alert) {
  border-radius: var(--radius-md);
  border: none;
}

:deep(.el-alert--info) {
  background: rgba(102, 126, 234, 0.1);
  color: var(--color-primary);
}

:deep(.el-textarea__inner) {
  border-radius: var(--radius-md);
  font-family: 'Consolas', 'Monaco', monospace;
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
}
</style>

