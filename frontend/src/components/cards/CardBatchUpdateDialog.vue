<template>
  <el-dialog v-model="dialogVisible" title="批量修改卡密" width="600px" @close="handleClose">
    <el-alert type="info" :closable="false" style="margin-bottom: 20px;">
      已选择 {{ selectedCount }} 个卡密,只修改填写的字段
    </el-alert>
    <el-form :model="form" :label-width="isMobile ? '0px' : '120px'" :label-position="isMobile ? 'top' : 'right'">
      <el-form-item label="卡密类型">
        <el-checkbox v-model="form.update_card_type" style="margin-bottom: 10px;">
          修改类型
        </el-checkbox>
        <el-select v-if="form.update_card_type" v-model="form.card_type" placeholder="请选择卡密类型" style="width: 100%;">
          <el-option
            v-for="type in CARD_TYPES"
            :key="type.value"
            :label="type.label"
            :value="type.value"
          />
        </el-select>
      </el-form-item>
      
      <el-form-item label="有效时长" v-if="form.card_type !== 'permanent' || !form.update_card_type">
        <el-checkbox v-model="form.update_duration" style="margin-bottom: 10px;">
          修改时长
        </el-checkbox>
        <div v-if="form.update_duration">
          <div style="display: flex; gap: 10px; align-items: center; flex-wrap: wrap; margin-bottom: 5px;">
            <div style="display: flex; gap: 10px; align-items: center;">
              <el-input-number v-model="form.duration_value" :min="form.card_type === 'permanent' ? 0 : 1" style="width: 150px;" />
              <el-select v-model="form.duration_unit" style="width: 100px;">
                <el-option label="秒" value="second" />
                <el-option label="天" value="day" />
                <el-option label="周" value="week" />
                <el-option label="月" value="month" />
                <el-option label="季" value="quarter" />
                <el-option label="年" value="year" />
              </el-select>
            </div>
            <div v-if="durationPreview" style="color: #999; font-size: 12px;">
              {{ durationPreview }}
            </div>
          </div>
          <el-alert type="warning" :closable="false" show-icon style="margin-top: 5px;">
            批量修改将重置所有已激活卡密的到期时间
          </el-alert>
        </div>
      </el-form-item>
      
      <el-form-item v-if="form.update_card_type && form.card_type === 'permanent'">
        <el-alert type="info" :closable="false" show-icon>
          永久卡无需设置时长
        </el-alert>
      </el-form-item>
      
      <el-form-item label="设备码上限">
        <el-checkbox v-model="form.update_max_hwid" style="margin-bottom: 10px;">
          修改设备码上限
        </el-checkbox>
        <el-input-number v-if="form.update_max_hwid" v-model="form.max_hwid" :min="-1" />
      </el-form-item>
      
      <el-form-item label="IP上限">
        <el-checkbox v-model="form.update_max_ip" style="margin-bottom: 10px;">
          修改IP上限
        </el-checkbox>
        <el-input-number v-if="form.update_max_ip" v-model="form.max_ip" :min="-1" />
      </el-form-item>
      
      <el-form-item label="备注">
        <el-checkbox v-model="form.update_note" style="margin-bottom: 10px;">
          修改备注
        </el-checkbox>
        <el-input v-if="form.update_note" v-model="form.note" type="textarea" />
      </el-form-item>
    </el-form>
    
    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" @click="handleSave">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { useResponsive } from '@/composables/useResponsive'
import { CARD_TYPES } from '@/constants/cardTypes'
import { unitValueToSeconds } from '@/composables/useDuration'

const { isMobile } = useResponsive()

// 实时预览
const durationPreview = computed(() => {
  if (!form.value.update_duration || form.value.card_type === 'permanent') {
    return null
  }
  const seconds = unitValueToSeconds(form.value.duration_value, form.value.duration_unit)
  return `${seconds.toLocaleString()}秒`
})

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  selectedCount: {
    type: Number,
    default: 0
  }
})

const emit = defineEmits(['update:visible', 'save'])

const dialogVisible = ref(false)

const form = ref({
  update_duration: false,
  duration_value: 30,
  duration_unit: 'day',
  update_card_type: false,
  card_type: 'month',
  update_max_hwid: false,
  max_hwid: -1,
  update_max_ip: false,
  max_ip: -1,
  update_note: false,
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
    update_duration: false,
    duration_value: 30,
    duration_unit: 'day',
    update_card_type: false,
    card_type: 'month',
    update_max_hwid: false,
    max_hwid: -1,
    update_max_ip: false,
    max_ip: -1,
    update_note: false,
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

:deep(.el-alert) {
  border-radius: var(--radius-md);
  border: none;
}

:deep(.el-alert--info) {
  background: rgba(102, 126, 234, 0.1);
  color: var(--color-primary);
}

@media (max-width: 768px) {
  :deep(.el-dialog) {
    width: 90% !important;
    margin: 20px auto;
  }
  
  :deep(.el-dialog__body) {
    padding: 16px;
  }
  
  :deep(.el-form-item) {
    margin-bottom: 16px;
  }
  
  :deep(.el-form-item__label) {
    margin-bottom: 6px;
  }
}
</style>

