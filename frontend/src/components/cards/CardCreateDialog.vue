<template>
  <div class="modern-dialog theme-success">
    <el-dialog 
      v-model="dialogVisible" 
      title="生成卡密" 
      :width="isMobile ? '95%' : '600px'"
      :fullscreen="isMobile"
      :close-on-click-modal="false"
      @close="handleClose"
      @opened="handleOpened"
    >
    <el-form :model="form" :label-width="isMobile ? '0px' : '120px'" :label-position="isMobile ? 'top' : 'right'">
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
      
      <el-form-item label="卡密类型">
        <el-select v-model="form.card_type" placeholder="请选择卡密类型" style="width: 100%;">
          <el-option
            v-for="type in CARD_TYPES"
            :key="type.value"
            :label="type.label"
            :value="type.value"
          />
        </el-select>
      </el-form-item>
      
      <el-form-item label="有效时长" v-if="form.card_type !== 'permanent'">
        <div style="display: flex; gap: 10px; align-items: center; flex-wrap: wrap;">
          <div style="display: flex; gap: 10px; align-items: center;">
            <el-input-number v-model="form.duration_value" :min="1" style="width: 150px;" />
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
      </el-form-item>
      
      <el-form-item v-if="form.card_type === 'permanent'">
        <el-alert type="info" :closable="false" show-icon>
          永久卡无需设置时长
        </el-alert>
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
  </div>
</template>

<script setup>
import { ref, watch, nextTick, computed } from 'vue'
import { useResponsive } from '@/composables/useResponsive'
import { staggerFormItems } from '@/utils/animations'
import { CARD_TYPES } from '@/constants/cardTypes'
import { unitValueToSeconds } from '@/composables/useDuration'

const { isMobile } = useResponsive()

// 实时预览
const durationPreview = computed(() => {
  if (form.value.card_type === 'permanent') {
    return null
  }
  const seconds = unitValueToSeconds(form.value.duration_value, form.value.duration_unit)
  return `${seconds.toLocaleString()}秒`
})

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
  card_type: 'month',
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

// card_type 自动同步
watch([() => form.value.duration_value, () => form.value.duration_unit], () => {
  if (form.value.card_type === 'permanent') {
    return
  }
  
  const seconds = unitValueToSeconds(form.value.duration_value, form.value.duration_unit)
  
  const rules = [
    { seconds: 0, type: 'permanent' },
    { seconds: 604800, type: 'trial' },      // 7天
    { seconds: 2592000, type: 'month' },     // 30天
    { seconds: 7776000, type: 'quarter' },   // 90天
    { seconds: 31536000, type: 'year' }      // 365天
  ]
  
  for (const rule of rules) {
    if (seconds === rule.seconds) {
      form.value.card_type = rule.type
      break
    }
  }
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
    card_type: 'month',
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

const handleOpened = () => {
  nextTick(() => {
    const formItems = document.querySelectorAll('.el-form-item')
    if (formItems.length > 0) {
      staggerFormItems(formItems)
    }
  })
}
</script>

<style scoped>
/* 组件特有样式 */
</style>

