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
          <el-input v-model="form.prefix" placeholder="可选" />
        </el-form-item>
        <el-form-item label="后缀">
          <el-input v-model="form.suffix" placeholder="可选" />
        </el-form-item>
        <el-form-item label="生成数量">
          <el-input-number v-model="form.count" :min="1" :max="1000" />
        </el-form-item>
        
        <el-form-item label="高级选项">
          <el-collapse v-model="advancedOptions" style="width: 100%;">
            <el-collapse-item name="1">
              <template #title>
                <span style="font-size: 14px; color: #606266;">卡密生成配置</span>
              </template>
              
              <el-form :label-width="isMobile ? '0px' : '100px'" :label-position="isMobile ? 'top' : 'right'">
                <el-form-item label="随机长度">
                  <el-input-number v-model="form.length" :min="6" :max="32" />
                  <div style="color: #999; font-size: 12px; margin-top: 5px;">
                    随机部分的字符数量，范围6-32
                  </div>
                </el-form-item>
                
                <el-form-item label="字符类型">
                  <el-radio-group v-model="form.charset_type">
                    <el-radio label="alphanumeric">英文+数字</el-radio>
                    <el-radio label="letters">仅英文</el-radio>
                  </el-radio-group>
                  <div style="color: #999; font-size: 12px; margin-top: 5px;">
                    随机部分使用的字符集
                  </div>
                </el-form-item>
              </el-form>
            </el-collapse-item>
          </el-collapse>
        </el-form-item>
      </template>
      
      <el-form-item label="有效时长">
        <el-radio-group v-model="form.is_permanent" style="margin-bottom: 10px;">
          <el-radio :label="false">限时卡</el-radio>
          <el-radio :label="true">永久卡</el-radio>
        </el-radio-group>
        
        <div v-if="!form.is_permanent" style="display: flex; gap: 10px; align-items: center; flex-wrap: wrap;">
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
        
        <div v-else style="color: #999; font-size: 12px;">
          永久有效，无时间限制
        </div>
      </el-form-item>
      
      <el-form-item label="卡密类型">
        <el-select v-model="form.card_type" placeholder="根据时长自动推断" style="width: 100%;" disabled>
          <el-option
            v-for="type in CARD_TYPES"
            :key="type.value"
            :label="type.label"
            :value="type.value"
          />
        </el-select>
        <div style="color: #999; font-size: 12px; margin-top: 5px;">
          系统根据有效时长自动设置
        </div>
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
  if (form.value.is_permanent) {
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
const advancedOptions = ref([])

const form = ref({
  mode: 'batch',
  card_key: '',
  prefix: '',
  suffix: '',
  count: 10,
  length: 16,
  charset_type: 'alphanumeric',
  is_permanent: false,
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

// card_type 自动同步 - 根据时长范围推断
watch([() => form.value.is_permanent, () => form.value.duration_value, () => form.value.duration_unit], () => {
  if (form.value.is_permanent) {
    form.value.card_type = 'permanent'
    return
  }
  
  const seconds = unitValueToSeconds(form.value.duration_value, form.value.duration_unit)
  
  // 范围推断逻辑
  if (seconds < 604800) {                  // <7天
    form.value.card_type = 'trial'
  } else if (seconds < 5184000) {          // 7天-60天
    form.value.card_type = 'month'
  } else if (seconds < 15552000) {         // 60天-180天
    form.value.card_type = 'quarter'
  } else {                                 // ≥180天
    form.value.card_type = 'year'
  }
})

const resetForm = () => {
  form.value = {
    mode: 'batch',
    card_key: '',
    prefix: '',
    suffix: '',
    count: 10,
    length: 16,
    charset_type: 'alphanumeric',
    is_permanent: false,
    duration_value: 30,
    duration_unit: 'day',
    card_type: 'month',
    max_hwid: -1,
    max_ip: -1,
    note: ''
  }
  advancedOptions.value = []
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

