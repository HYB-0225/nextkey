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
      
      <!-- 未激活卡：编辑有效时长 -->
      <el-form-item v-if="!form.activated && form.card_type !== 'permanent'" label="有效时长">
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
        <div style="color: #999; font-size: 12px; margin-top: 5px;">
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
  if (form.value.activated || form.value.card_type === 'permanent') {
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
      duration_value: value,
      duration_unit: unit,
      expire_time: data.expire_at ? new Date(data.expire_at) : null
    }
  }
})

watch(dialogVisible, (val) => {
  emit('update:visible', val)
})

// card_type 自动同步
watch([() => form.value.duration_value, () => form.value.duration_unit], () => {
  if (form.value.activated || form.value.card_type === 'permanent') {
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

