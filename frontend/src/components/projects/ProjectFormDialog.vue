<template>
  <el-dialog v-model="dialogVisible" :title="title" width="600px" @close="handleClose">
    <el-form :model="form" label-width="120px">
      <el-form-item label="项目名称">
        <el-input v-model="form.name" />
      </el-form-item>
      <el-form-item label="模式">
        <el-radio-group v-model="form.mode">
          <el-radio label="free">免费</el-radio>
          <el-radio label="paid">付费</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="版本号">
        <el-input v-model="form.version" />
      </el-form-item>
      <el-form-item label="Token有效期">
        <el-input v-model.number="form.token_expire" type="number">
          <template #append>秒</template>
        </el-input>
      </el-form-item>
      <el-form-item label="启用机器码">
        <el-switch v-model="form.enable_hwid" />
      </el-form-item>
      <el-form-item label="启用IP验证">
        <el-switch v-model="form.enable_ip" />
      </el-form-item>
      <el-form-item label="描述">
        <el-input v-model="form.description" type="textarea" />
      </el-form-item>
    </el-form>
    
    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" @click="handleSave">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: '创建项目'
  },
  projectData: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['update:visible', 'save'])

const dialogVisible = ref(false)

const form = ref({
  name: '',
  mode: 'free',
  version: '1.0.0',
  token_expire: 3600,
  enable_hwid: true,
  enable_ip: true,
  description: ''
})

watch(() => props.visible, (val) => {
  dialogVisible.value = val
  if (val) {
    if (props.projectData) {
      form.value = { ...props.projectData }
    } else {
      resetForm()
    }
  }
})

watch(dialogVisible, (val) => {
  emit('update:visible', val)
})

const resetForm = () => {
  form.value = {
    name: '',
    mode: 'free',
    version: '1.0.0',
    token_expire: 3600,
    enable_hwid: true,
    enable_ip: true,
    description: ''
  }
}

const handleClose = () => {
  dialogVisible.value = false
}

const handleSave = () => {
  emit('save', form.value)
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

:deep(.el-form-item__label) {
  font-weight: 500;
  color: var(--color-text-primary);
}

:deep(.el-input__wrapper),
:deep(.el-textarea__inner) {
  border-radius: var(--radius-md);
  transition: all var(--duration-fast) var(--ease-out);
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

