<template>
  <div class="modern-dialog theme-info">
    <el-dialog 
      v-model="dialogVisible" 
      :title="title" 
      :width="isMobile ? '95%' : '600px'"
      :fullscreen="isMobile"
      :close-on-click-modal="false"
      @close="handleClose"
      @opened="handleOpened"
    >
      <el-form :model="form" :label-width="isMobile ? '0px' : '100px'" :label-position="isMobile ? 'top' : 'right'">
        <el-form-item label="变量名">
          <el-input v-model="form.key" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="值">
          <el-input v-model="form.value" type="textarea" :rows="5" />
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
import { ref, watch, nextTick } from 'vue'
import { useResponsive } from '@/composables/useResponsive'
import { staggerFormItems } from '@/utils/animations'

const { isMobile } = useResponsive()

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: '添加变量'
  },
  isEdit: {
    type: Boolean,
    default: false
  },
  varData: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['update:visible', 'save'])

const dialogVisible = ref(false)

const form = ref({
  key: '',
  value: ''
})

watch(() => props.visible, (val) => {
  dialogVisible.value = val
  if (val) {
    if (props.varData) {
      form.value = { ...props.varData }
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
    key: '',
    value: ''
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
