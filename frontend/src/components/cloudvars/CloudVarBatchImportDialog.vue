<template>
  <el-dialog v-model="dialogVisible" title="批量导入变量" width="700px" @close="handleClose">
    <el-alert type="info" :closable="false" style="margin-bottom: 15px;">
      请输入JSON格式的变量数据,格式: [{"key": "变量名", "value": "值"}, ...]
    </el-alert>
    <el-input
      v-model="jsonData"
      type="textarea"
      :rows="12"
      placeholder='[
  {"key": "api_url", "value": "https://api.example.com"},
  {"key": "app_name", "value": "MyApp"},
  {"key": "version", "value": "1.0.0"}
]'
    />
    
    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" @click="handleSave">确定导入</el-button>
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
  animation: scale-bounce 0.4s var(--ease-bounce);
}

:deep(.el-dialog__header) {
  background: linear-gradient(135deg, #FF8C42 0%, #FF6B35 100%);
  padding: 20px 24px;
  border-bottom: none;
}

:deep(.el-dialog__title) {
  color: #fff;
  font-weight: 600;
}

:deep(.el-dialog__headerbtn .el-dialog__close) {
  color: #fff;
}

:deep(.el-dialog__body) {
  padding: 24px;
}

:deep(.el-alert) {
  border-radius: var(--radius-md);
  border: none;
}

:deep(.el-alert--info) {
  background: rgba(255, 140, 66, 0.1);
  color: var(--color-primary);
}

:deep(.el-textarea__inner) {
  border-radius: var(--radius-md);
  font-family: 'Consolas', 'Monaco', monospace;
  transition: all var(--duration-normal) var(--ease-out);
}

:deep(.el-textarea__inner:hover) {
  border-color: var(--color-primary-light);
}

:deep(.el-textarea.is-focus .el-textarea__inner) {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(255, 140, 66, 0.1);
  transform: translateY(-1px);
}

:deep(.el-button) {
  border-radius: var(--radius-md);
  font-weight: 500;
  transition: all var(--duration-fast) var(--ease-out);
}

:deep(.el-button:hover:not(:disabled)) {
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

:deep(.el-button:active:not(:disabled)) {
  transform: translateY(0);
}

:deep(.el-button--primary) {
  background: linear-gradient(135deg, #FF8C42 0%, #FF6B35 100%);
  border: none;
}

@media (max-width: 768px) {
  :deep(.el-dialog) {
    width: 90% !important;
    margin: 20px auto;
  }
}
</style>

