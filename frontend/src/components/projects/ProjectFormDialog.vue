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
      <el-divider content-position="left">加密配置</el-divider>
      <el-form-item label="加密方案">
        <el-select 
          v-model="form.encryption_scheme" 
          :loading="loadingSchemes"
          :disabled="!projectData && loadingSchemes"
          placeholder="选择加密方案"
        >
          <el-option 
            v-for="scheme in encryptionSchemes" 
            :key="scheme.scheme"
            :label="scheme.name"
            :value="scheme.scheme"
          >
            <div class="scheme-option">
              <span class="scheme-name">{{ scheme.name }}</span>
              <div class="scheme-tags">
                <el-tag 
                  :type="getSecurityTagType(scheme.security_level)"
                  size="small"
                  effect="plain"
                >
                  {{ getSecurityLabel(scheme.security_level) }}
                </el-tag>
                <el-tag 
                  :type="getPerformanceTagType(scheme.performance)"
                  size="small"
                  effect="plain"
                >
                  {{ getPerformanceLabel(scheme.performance) }}
                </el-tag>
              </div>
            </div>
          </el-option>
        </el-select>
        <div class="form-tip" v-if="!projectData">创建后将自动生成加密密钥</div>
        <div class="form-tip" v-else>修改加密方案将生成新密钥，需在客户端同步更新</div>
      </el-form-item>
      <el-form-item v-if="projectData" label="加密密钥">
        <el-input v-model="form.encryption_key" readonly>
          <template #append>
            <el-button @click="copyKey">复制</el-button>
          </template>
        </el-input>
        <div class="form-tip">客户端需要此密钥进行加密通信</div>
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
import { ref, watch, nextTick, onMounted } from 'vue'
import { useResponsive } from '@/composables/useResponsive'
import { staggerFormItems } from '@/utils/animations'
import { copyToClipboard } from '@/utils/copy'
import { getEncryptionSchemes, updateProjectEncryption } from '@/api/crypto'
import { ElMessage, ElMessageBox } from 'element-plus'

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
const encryptionSchemes = ref([])
const loadingSchemes = ref(false)
const originalScheme = ref('')

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
  description: '',
  encryption_scheme: 'aes-256-gcm'
})

watch(() => props.visible, (val) => {
  dialogVisible.value = val
  if (val) {
    if (props.projectData) {
      form.value = { ...props.projectData }
      originalScheme.value = props.projectData.encryption_scheme
    } else {
      resetForm()
      originalScheme.value = ''
    }
  }
})

onMounted(() => {
  loadEncryptionSchemes()
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
    description: '',
    encryption_scheme: 'aes-256-gcm'
  }
}

const loadEncryptionSchemes = async () => {
  try {
    loadingSchemes.value = true
    const res = await getEncryptionSchemes()
    if (res && res.code === 0) {
      encryptionSchemes.value = res.data || []
    }
  } catch (error) {
    console.error('加载加密方案失败:', error)
  } finally {
    loadingSchemes.value = false
  }
}

const handleClose = () => {
  dialogVisible.value = false
}

const handleSave = async () => {
  // 如果是编辑模式且加密方案发生变化，需要确认
  if (props.projectData && form.value.encryption_scheme !== originalScheme.value) {
    try {
      await ElMessageBox.confirm(
        '修改加密方案将生成新的密钥，客户端需要同步更新。确定要修改吗？',
        '确认修改',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
      )
      
      // 单独更新加密方案
      try {
        const res = await updateProjectEncryption(props.projectData.id, {
          encryption_scheme: form.value.encryption_scheme
        })
        if (res.code === 0) {
          form.value.encryption_key = res.data.encryption_key
          ElMessage.success('加密方案已更新')
        }
      } catch (error) {
        ElMessage.error('更新加密方案失败')
        return
      }
    } catch {
      // 用户取消，恢复原值
      form.value.encryption_scheme = originalScheme.value
      return
    }
  }
  
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

const copyKey = () => {
  if (form.value.encryption_key) {
    copyToClipboard(form.value.encryption_key, '加密密钥已复制')
  }
}

const getSecurityTagType = (level) => {
  const types = {
    'secure': 'success',
    'weak': 'warning',
    'insecure': 'danger'
  }
  return types[level] || 'info'
}

const getSecurityLabel = (level) => {
  const labels = {
    'secure': '安全',
    'weak': '弱',
    'insecure': '不安全'
  }
  return labels[level] || level
}

const getPerformanceTagType = (performance) => {
  const types = {
    'fast': 'primary',
    'medium': 'info',
    'slow': 'warning'
  }
  return types[performance] || 'info'
}

const getPerformanceLabel = (performance) => {
  const labels = {
    'fast': '快速',
    'medium': '中等',
    'slow': '慢速'
  }
  return labels[performance] || performance
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

.scheme-option {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  padding: 4px 0;
}

.scheme-name {
  font-weight: 500;
  flex: 1;
}

.scheme-tags {
  display: flex;
  gap: 6px;
  align-items: center;
}

.scheme-tags .el-tag {
  font-size: 11px;
  padding: 0 6px;
  height: 20px;
  line-height: 20px;
}
</style>

