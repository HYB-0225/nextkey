<template>
  <el-dialog v-model="dialogVisible" title="生成卡密" width="600px" @close="handleClose">
    <el-form :model="form" label-width="120px">
      <el-form-item label="生成方式">
        <el-radio-group v-model="form.mode">
          <el-radio label="batch">批量生成</el-radio>
          <el-radio label="custom">自定义卡密</el-radio>
        </el-radio-group>
      </el-form-item>
      
      <el-form-item label="自定义卡密" v-if="form.mode === 'custom'">
        <el-input v-model="form.card_key" placeholder="输入完整卡密" />
      </el-form-item>
      
      <template v-if="form.mode === 'batch'">
        <el-form-item label="前缀">
          <el-input v-model="form.prefix" />
        </el-form-item>
        <el-form-item label="后缀">
          <el-input v-model="form.suffix" />
        </el-form-item>
        <el-form-item label="生成数量">
          <el-input-number v-model="form.count" :min="1" :max="1000" />
        </el-form-item>
      </template>
      
      <el-form-item label="有效时长">
        <div style="display: flex; gap: 10px; align-items: center;">
          <el-input-number v-model="form.duration_value" :min="0" style="width: 150px;" />
          <el-select v-model="form.duration_unit" style="width: 100px;">
            <el-option label="秒" value="second" />
            <el-option label="天" value="day" />
            <el-option label="周" value="week" />
            <el-option label="月" value="month" />
            <el-option label="季" value="quarter" />
            <el-option label="年" value="year" />
          </el-select>
        </div>
        <div style="color: #999; font-size: 12px; margin-top: 5px;">0表示永久</div>
      </el-form-item>
      
      <el-form-item label="卡密类型">
        <el-input v-model="form.card_type" placeholder="normal" />
      </el-form-item>
      
      <el-form-item label="设备码上限">
        <el-input-number v-model="form.max_hwid" :min="-1" />
        <div style="color: #999; font-size: 12px; margin-top: 5px;">-1表示无限制</div>
      </el-form-item>
      
      <el-form-item label="IP上限">
        <el-input-number v-model="form.max_ip" :min="-1" />
        <div style="color: #999; font-size: 12px; margin-top: 5px;">-1表示无限制</div>
      </el-form-item>
      
      <el-form-item label="备注">
        <el-input v-model="form.note" type="textarea" />
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
  }
})

const emit = defineEmits(['update:visible', 'save'])

const dialogVisible = ref(false)

const form = ref({
  mode: 'batch',
  card_key: '',
  prefix: '',
  suffix: '',
  count: 10,
  duration_value: 30,
  duration_unit: 'day',
  card_type: 'normal',
  max_hwid: -1,
  max_ip: -1,
  note: ''
})

watch(() => props.visible, (val) => {
  dialogVisible.value = val
  if (val) {
    resetForm()
  }
})

watch(dialogVisible, (val) => {
  emit('update:visible', val)
})

const resetForm = () => {
  form.value = {
    mode: 'batch',
    card_key: '',
    prefix: '',
    suffix: '',
    count: 10,
    duration_value: 30,
    duration_unit: 'day',
    card_type: 'normal',
    max_hwid: -1,
    max_ip: -1,
    note: ''
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

:deep(.el-input__wrapper:hover),
:deep(.el-textarea__inner:hover) {
  border-color: var(--color-primary-light);
}

:deep(.el-input.is-focus .el-input__wrapper),
:deep(.el-textarea.is-focus .el-textarea__inner) {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(24, 144, 255, 0.1);
}

:deep(.el-button) {
  border-radius: var(--radius-md);
  font-weight: 500;
  transition: all var(--duration-fast) var(--ease-out);
}

:deep(.el-button:hover) {
  transform: translateY(-1px);
}

:deep(.el-button:active) {
  transform: translateY(0);
}

@media (max-width: 768px) {
  :deep(.el-dialog) {
    width: 90% !important;
    margin: 20px auto;
  }
}
</style>

