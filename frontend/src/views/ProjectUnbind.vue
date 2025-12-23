<template>
  <div class="unbind-page">
    <div class="background-orbs">
      <div class="orb orb-1"></div>
      <div class="orb orb-2"></div>
      <div class="orb orb-3"></div>
    </div>

    <el-card class="unbind-card">
      <div class="card-header">
        <h1>卡密解绑</h1>
        <p class="subtitle">请输入卡密进行解绑</p>
      </div>

      <el-form :model="form" label-position="top" class="unbind-form">
        <el-form-item label="卡密">
          <el-input v-model="form.cardKey" placeholder="请输入卡密" />
        </el-form-item>
      </el-form>

      <div class="action-row">
        <el-button
          type="primary"
          size="large"
          :loading="loading"
          :disabled="!canSubmit"
          @click="handleUnbind"
        >
          立即解绑
        </el-button>
      </div>

      <el-alert
        v-if="statusMessage"
        :type="statusType"
        show-icon
        class="result-alert"
        :title="statusMessage"
      />
    </el-card>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'
import { ElMessage } from 'element-plus'

const route = useRoute()
const unbindSlug = computed(() => String(route.params.unbindSlug || '').trim())

const form = ref({
  cardKey: ''
})

const loading = ref(false)
const statusMessage = ref('')
const statusType = ref('success')

const canSubmit = computed(() =>
  Boolean(unbindSlug.value && form.value.cardKey.trim())
)

const handleUnbind = async () => {
  if (!unbindSlug.value) {
    ElMessage.error('项目参数缺失')
    return
  }
  if (!form.value.cardKey.trim()) {
    ElMessage.warning('请输入卡密')
    return
  }
  loading.value = true
  statusMessage.value = ''
  try {
    const response = await axios.post('/api/card/unbind-public', {
      unbind_slug: unbindSlug.value,
      card_key: form.value.cardKey.trim()
    })
    const result = response?.data
    if (result?.code !== 0) {
      throw new Error(result?.message || '解绑失败')
    }
    statusType.value = 'success'
    statusMessage.value = result?.data?.message || '解绑成功'
    ElMessage.success(statusMessage.value)
  } catch (error) {
    const message = error?.message || '解绑失败'
    statusType.value = 'error'
    statusMessage.value = message
    ElMessage.error(message)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.unbind-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  background: linear-gradient(135deg, #ffedd5 0%, #fff7ed 40%, #fef3c7 100%);
  position: relative;
  overflow: hidden;
}

.background-orbs {
  position: absolute;
  inset: 0;
  z-index: 0;
}

.orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(70px);
  opacity: 0.5;
  animation: float 16s ease-in-out infinite;
}

.orb-1 {
  width: 360px;
  height: 360px;
  background: radial-gradient(circle, rgba(251, 146, 60, 0.7) 0%, transparent 70%);
  top: -15%;
  left: -10%;
}

.orb-2 {
  width: 420px;
  height: 420px;
  background: radial-gradient(circle, rgba(253, 186, 116, 0.5) 0%, transparent 70%);
  bottom: -20%;
  right: -15%;
  animation-delay: -6s;
}

.orb-3 {
  width: 320px;
  height: 320px;
  background: radial-gradient(circle, rgba(251, 191, 36, 0.4) 0%, transparent 70%);
  top: 45%;
  left: 55%;
  animation-delay: -12s;
}

.unbind-card {
  width: 100%;
  max-width: 520px;
  border-radius: 0;
  border: 3px solid #fb923c;
  box-shadow: 8px 8px 0 0 rgba(0, 0, 0, 0.2),
              0 0 0 2px #fcd34d inset;
  background: rgba(255, 251, 235, 0.98);
  position: relative;
  z-index: 1;
  padding: 24px 28px;
}

.card-header {
  text-align: center;
  margin-bottom: 16px;
}

.card-header h1 {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 6px;
  color: #c2410c;
}

.subtitle {
  font-size: 13px;
  color: var(--color-text-secondary);
}

.unbind-form :deep(.el-input__wrapper) {
  border-radius: var(--radius-sm);
  border: 1px solid var(--color-border-light);
}

.action-row {
  display: flex;
  justify-content: center;
  margin-top: 12px;
}

.action-row :deep(.el-button) {
  width: 100%;
  border-radius: 0;
  font-weight: 600;
  background: linear-gradient(135deg, #fb923c 0%, #f97316 100%);
  border: 2px solid #ea580c;
  box-shadow: 4px 4px 0 0 rgba(0, 0, 0, 0.2);
}

.result-alert {
  margin-top: 16px;
}

@keyframes float {
  0%, 100% {
    transform: translate(0, 0) scale(1);
  }
  50% {
    transform: translate(30px, -30px) scale(1.05);
  }
}

@media (max-width: 768px) {
  .unbind-card {
    padding: 20px;
  }

  .card-header h1 {
    font-size: 24px;
  }
}
</style>
