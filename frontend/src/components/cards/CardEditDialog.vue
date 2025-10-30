<template>
  <div class="modern-dialog theme-info">
    <el-dialog 
      v-model="dialogVisible" 
      title="编辑卡密" 
      :width="isMobile ? '95%' : '600px'"
      :fullscreen="isMobile"
      :close-on-click-modal="false"
      @close="handleClose"
      @opened="handleOpened"
    >
    <el-form :model="form" :label-width="isMobile ? '0px' : '120px'" :label-position="isMobile ? 'top' : 'right'">
      <!-- 未激活卡：编辑有效时长 -->
      <el-form-item v-if="!form.activated" label="有效时长">
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
        
        <div v-if="!form.is_permanent" style="color: #999; font-size: 12px; margin-top: 5px;">
          修改的是卡密的基础时长
        </div>
      </el-form-item>
      
      <!-- 已激活卡：编辑到期时间 -->
      <el-form-item v-if="form.activated && form.card_type !== 'permanent'" label="到期时间">
        <el-date-picker
          v-model="form.expire_time"
          type="datetime"
          placeholder="选择到期时间"
          format="YYYY-MM-DD HH:mm:ss"
          value-format="YYYY-MM-DD HH:mm:ss"
          style="width: 100%;"
        />
        <div style="color: #999; font-size: 12px; margin-top: 5px;">
          直接修改到期时间，系统将自动反算有效时长
        </div>
      </el-form-item>
      
      <el-form-item v-if="form.activated && form.card_type === 'permanent'">
        <el-alert type="info" :closable="false" show-icon>
          永久卡无需设置时长
        </el-alert>
      </el-form-item>
      
      <el-form-item label="卡密类型">
        <el-select v-model="form.card_type" placeholder="根据时长自动推断" style="width: 100%;" :disabled="!form.activated">
          <el-option
            v-for="type in CARD_TYPES"
            :key="type.value"
            :label="type.label"
            :value="type.value"
          />
        </el-select>
        <div style="color: #999; font-size: 12px; margin-top: 5px;">
          <span v-if="!form.activated">系统根据有效时长自动设置</span>
          <span v-else>已激活卡密可手动修改类型</span>
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
      
      <el-form-item label="专属信息">
        <el-input v-model="form.custom_data" type="textarea" :rows="3" />
      </el-form-item>
      
      <el-form-item label="设备码列表">
        <div style="width: 100%;">
          <div style="display: flex; gap: 10px; margin-bottom: 10px;">
            <el-input 
              v-model="newHWID" 
              placeholder="输入设备码" 
              @keyup.enter="handleAddHWID"
              style="flex: 1;"
            />
            <el-button type="primary" @click="handleAddHWID">添加</el-button>
          </div>
          <div v-if="form.hwid_list && form.hwid_list.length > 0" style="display: flex; flex-wrap: wrap; gap: 8px;">
            <el-tag 
              v-for="(hwid, index) in form.hwid_list" 
              :key="index"
              closable
              @close="handleRemoveHWID(index)"
            >
              {{ hwid }}
            </el-tag>
          </div>
          <div v-else style="color: #999; font-size: 12px;">暂无设备码</div>
        </div>
      </el-form-item>
      
      <el-form-item label="IP列表">
        <div style="width: 100%;">
          <div style="display: flex; gap: 10px; margin-bottom: 10px;">
            <el-input 
              v-model="newIP" 
              placeholder="输入IP地址" 
              @keyup.enter="handleAddIP"
              style="flex: 1;"
            />
            <el-button type="primary" @click="handleAddIP">添加</el-button>
          </div>
          <div v-if="form.ip_list && form.ip_list.length > 0" style="display: flex; flex-wrap: wrap; gap: 8px;">
            <el-tag 
              v-for="(ip, index) in form.ip_list" 
              :key="index"
              closable
              @close="handleRemoveIP(index)"
              type="success"
            >
              {{ ip }}
            </el-tag>
          </div>
          <div v-else style="color: #999; font-size: 12px;">暂无IP地址</div>
        </div>
      </el-form-item>
    </el-form>
    
      <template #footer>
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, watch, nextTick, computed } from 'vue'
import { useResponsive } from '@/composables/useResponsive'
import { staggerFormItems } from '@/utils/animations'
import { CARD_TYPES } from '@/constants/cardTypes'
import { unitValueToSeconds, secondsToUnitValue } from '@/composables/useDuration'

const { isMobile } = useResponsive()

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  cardData: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['update:visible', 'save'])

const dialogVisible = ref(false)

const form = ref({
  id: null,
  activated: false,
  is_permanent: false,
  duration_value: 30,
  duration_unit: 'day',
  expire_time: null,
  card_type: 'month',
  max_hwid: -1,
  max_ip: -1,
  note: '',
  custom_data: '',
  hwid_list: [],
  ip_list: []
})

const newHWID = ref('')
const newIP = ref('')

// 实时预览
const durationPreview = computed(() => {
  if (form.value.activated || form.value.is_permanent) {
    return null
  }
  const seconds = unitValueToSeconds(form.value.duration_value, form.value.duration_unit)
  return `${seconds.toLocaleString()}秒`
})

watch(() => props.visible, (val) => {
  dialogVisible.value = val
})

watch(() => props.cardData, (data) => {
  if (data) {
    const { value, unit } = secondsToUnitValue(data.duration || 0)
    form.value = {
      ...data,
      is_permanent: data.duration === 0 || data.card_type === 'permanent',
      duration_value: value || 30,
      duration_unit: unit,
      expire_time: data.expire_at ? new Date(data.expire_at) : null
    }
  }
})

watch(dialogVisible, (val) => {
  emit('update:visible', val)
})

// card_type 自动同步 - 根据时长范围推断 (仅未激活卡)
watch([() => form.value.is_permanent, () => form.value.duration_value, () => form.value.duration_unit], () => {
  if (form.value.activated) {
    return
  }
  
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

const handleClose = () => {
  dialogVisible.value = false
}

const handleSave = () => {
  emit('save', form.value)
}

const handleAddHWID = () => {
  if (newHWID.value.trim()) {
    if (!form.value.hwid_list) {
      form.value.hwid_list = []
    }
    if (!form.value.hwid_list.includes(newHWID.value.trim())) {
      form.value.hwid_list.push(newHWID.value.trim())
      newHWID.value = ''
    }
  }
}

const handleRemoveHWID = (index) => {
  form.value.hwid_list.splice(index, 1)
}

const handleAddIP = () => {
  if (newIP.value.trim()) {
    if (!form.value.ip_list) {
      form.value.ip_list = []
    }
    if (!form.value.ip_list.includes(newIP.value.trim())) {
      form.value.ip_list.push(newIP.value.trim())
      newIP.value = ''
    }
  }
}

const handleRemoveIP = (index) => {
  form.value.ip_list.splice(index, 1)
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

