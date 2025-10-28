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
import { ref, watch, nextTick } from 'vue'
import { useResponsive } from '@/composables/useResponsive'
import { staggerFormItems } from '@/utils/animations'

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
  duration_value: 30,
  duration_unit: 'day',
  card_type: 'normal',
  max_hwid: -1,
  max_ip: -1,
  note: '',
  custom_data: '',
  hwid_list: [],
  ip_list: []
})

const newHWID = ref('')
const newIP = ref('')

watch(() => props.visible, (val) => {
  dialogVisible.value = val
})

watch(() => props.cardData, (data) => {
  if (data) {
    form.value = { ...data }
  }
})

watch(dialogVisible, (val) => {
  emit('update:visible', val)
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

