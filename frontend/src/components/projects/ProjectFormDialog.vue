<template>
  <div class="modern-dialog theme-primary">
    <el-dialog 
      v-model="dialogVisible" 
      :title="title" 
      :width="isMobile ? '95%' : '600px'"
      :fullscreen="isMobile"
      :close-on-click-modal="false"
      @close="handleClose"
      @opened="handleOpened"
    >
    <el-form :model="form" :label-width="isMobile ? '0px' : '120px'" :label-position="isMobile ? 'top' : 'right'">
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
      <el-form-item label="更新地址">
        <el-input v-model="form.update_url" placeholder="可选，用于客户端检查更新" />
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
      <el-divider content-position="left">解绑配置</el-divider>
      <el-form-item label="启用解绑">
        <el-switch v-model="form.enable_unbind" />
      </el-form-item>
      <el-form-item label="验证HWID" v-if="form.enable_unbind">
        <el-switch v-model="form.unbind_verify_hwid" />
        <div class="form-tip">关闭后不验证HWID是否已绑定，但仍需传入HWID以从列表中移除</div>
      </el-form-item>
      <el-form-item label="解绑扣时" v-if="form.enable_unbind">
        <el-input v-model.number="form.unbind_deduct_time" type="number" placeholder="0表示不扣时">
          <template #append>秒</template>
        </el-input>
      </el-form-item>
      <el-form-item label="解绑冷却" v-if="form.enable_unbind">
        <el-input v-model.number="form.unbind_cooldown" type="number">
          <template #append>秒</template>
        </el-input>
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
  update_url: '',
  token_expire: 3600,
  enable_hwid: true,
  enable_ip: true,
  enable_unbind: false,
  unbind_verify_hwid: true,
  unbind_deduct_time: 0,
  unbind_cooldown: 86400,
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
    update_url: '',
    token_expire: 3600,
    enable_hwid: true,
    enable_ip: true,
    enable_unbind: false,
    unbind_verify_hwid: true,
    unbind_deduct_time: 0,
    unbind_cooldown: 86400,
    description: ''
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
:deep(.el-radio) {
  margin-right: 16px;
}

:deep(.el-switch) {
  --el-switch-on-color: #667eea;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
  line-height: 1.4;
}

:deep(.el-divider) {
  margin: 24px 0 16px 0;
}
</style>

