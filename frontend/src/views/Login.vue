<template>
  <div class="login-container">
    <!-- 动态背景 -->
    <div class="background-animation">
      <div class="gradient-orb orb-1"></div>
      <div class="gradient-orb orb-2"></div>
      <div class="gradient-orb orb-3"></div>
    </div>
    
    <!-- 登录卡片 -->
    <div class="login-card glass-effect">
      <div class="login-header">
        <div class="logo-wrapper">
          <div class="logo-icon">
            <div class="logo-shine"></div>
          </div>
        </div>
        <h1 class="title">NextKey</h1>
        <p class="subtitle">卡密验证与云控制系统</p>
      </div>
      
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="0"
        class="login-form"
      >
        <el-form-item prop="username">
          <el-input
            v-model="form.username"
            placeholder="用户名"
            size="large"
            prefix-icon="User"
            class="modern-input"
          />
        </el-form-item>
        
        <el-form-item prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="密码"
            size="large"
            prefix-icon="Lock"
            class="modern-input"
            @keyup.enter="handleLogin"
          />
        </el-form-item>
        
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            class="login-button"
            :loading="loading"
            @click="handleLogin"
          >
            <span v-if="!loading">登录</span>
          </el-button>
        </el-form-item>
      </el-form>
      
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { login } from '@/api/auth'
import { ElMessage } from 'element-plus'

const router = useRouter()
const authStore = useAuthStore()

const formRef = ref(null)
const loading = ref(false)

const form = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    loading.value = true
    try {
      const data = await login(form)
      authStore.setTokens(data.access_token, data.refresh_token, data.expires_in)
      ElMessage.success('登录成功')
      router.push('/')
    } catch (error) {
      console.error(error)
    } finally {
      loading.value = false
    }
  })
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #FF8C42 0%, #FFD93D 50%, #FF6B35 100%);
  position: relative;
  overflow: hidden;
  padding: 20px;
}

/* ==================== 动态背景 ==================== */
.background-animation {
  position: absolute;
  inset: 0;
  overflow: hidden;
  z-index: 0;
}

.gradient-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.5;
  animation: float 20s ease-in-out infinite;
}

.orb-1 {
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, rgba(255, 140, 66, 0.8) 0%, transparent 70%);
  top: -10%;
  left: -10%;
  animation-delay: 0s;
}

.orb-2 {
  width: 500px;
  height: 500px;
  background: radial-gradient(circle, rgba(255, 217, 61, 0.6) 0%, transparent 70%);
  bottom: -15%;
  right: -15%;
  animation-delay: -7s;
}

.orb-3 {
  width: 350px;
  height: 350px;
  background: radial-gradient(circle, rgba(255, 107, 53, 0.4) 0%, transparent 70%);
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  animation-delay: -14s;
}

@keyframes float {
  0%, 100% {
    transform: translate(0, 0) scale(1);
  }
  33% {
    transform: translate(50px, -50px) scale(1.1);
  }
  66% {
    transform: translate(-50px, 50px) scale(0.9);
  }
}

/* ==================== 登录卡片 ==================== */
.login-card {
  width: 100%;
  max-width: 420px;
  padding: 48px 40px;
  border-radius: 0;
  background: rgba(255, 251, 240, 0.98);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 4px solid #FF8C42;
  box-shadow: 8px 8px 0 0 rgba(0, 0, 0, 0.2),
              0 0 0 2px #FFD93D inset;
  position: relative;
  z-index: 1;
  animation: scale-bounce 0.6s var(--ease-bounce);
}

/* ==================== Logo和标题 ==================== */
.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.logo-wrapper {
  display: flex;
  justify-content: center;
  margin-bottom: 24px;
}

.logo-icon {
  width: 64px;
  height: 64px;
  border-radius: 0;
  background: linear-gradient(135deg, #FF8C42 0%, #FF6B35 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 4px 4px 0 0 rgba(0, 0, 0, 0.3);
  border: 3px solid #FFD93D;
  position: relative;
  overflow: hidden;
  animation: pulse-glow 3s ease-in-out infinite;
}

.logo-shine {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, transparent 0%, rgba(255, 255, 255, 0.3) 50%, transparent 100%);
  animation: shine 3s ease-in-out infinite;
}

@keyframes pulse-glow {
  0%, 100% {
    box-shadow: 4px 4px 0 0 rgba(0, 0, 0, 0.3);
  }
  50% {
    box-shadow: 6px 6px 0 0 rgba(0, 0, 0, 0.4);
  }
}

@keyframes shine {
  0% {
    transform: translateX(-100%) rotate(45deg);
  }
  100% {
    transform: translateX(200%) rotate(45deg);
  }
}

.title {
  font-size: 32px;
  font-weight: 700;
  margin-bottom: 8px;
  background: linear-gradient(135deg, #FF8C42 0%, #FF6B35 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  text-shadow: 2px 2px 0 rgba(255, 217, 61, 0.3);
}

.subtitle {
  color: var(--color-text-secondary);
  font-size: 14px;
  font-weight: 500;
}

/* ==================== 表单 ==================== */
.login-form {
  margin-bottom: 24px;
}

.login-form :deep(.el-form-item__content) {
  margin-left: 0 !important;
}

.login-form :deep(.modern-input) {
  width: 100%;
  transition: all var(--duration-normal) var(--ease-out);
}

.login-form :deep(.modern-input .el-input__wrapper) {
  background: rgba(255, 255, 255, 0.9);
  border: 2px solid var(--color-border);
  border-radius: 0;
  box-shadow: 2px 2px 0 0 rgba(0, 0, 0, 0.1),
              inset 1px 1px 0 0 rgba(255, 217, 61, 0.3);
  transition: all var(--duration-fast) steps(2);
  padding: 12px 16px;
  width: 100%;
}

.login-form :deep(.modern-input .el-input__wrapper:hover) {
  background: rgba(255, 255, 255, 1);
  box-shadow: 3px 3px 0 0 rgba(0, 0, 0, 0.15),
              inset 1px 1px 0 0 rgba(255, 217, 61, 0.5);
}

.login-form :deep(.modern-input.is-focus .el-input__wrapper) {
  background: #fff;
  border-color: #FF8C42;
  box-shadow: 0 0 0 4px rgba(255, 140, 66, 0.2),
              2px 2px 0 0 rgba(0, 0, 0, 0.1);
  transform: translate(-1px, -1px);
}

.login-form :deep(.modern-input .el-input__inner) {
  font-size: 15px;
  color: var(--color-text-primary);
}

.login-form :deep(.modern-input .el-input__inner::placeholder) {
  color: var(--color-text-tertiary);
}

/* ==================== 登录按钮 ==================== */
.login-button {
  width: 100%;
  height: 48px;
  font-size: 16px;
  font-weight: 600;
  background: linear-gradient(135deg, #FF8C42 0%, #FF6B35 100%);
  border: 3px solid #FF6B35;
  border-radius: 0;
  box-shadow: 4px 4px 0 0 rgba(0, 0, 0, 0.3);
  transition: all var(--duration-fast) steps(2);
  position: relative;
  overflow: hidden;
}

.login-button::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.2) 0%, transparent 100%);
  opacity: 0;
  transition: opacity var(--duration-normal) var(--ease-out);
}

.login-button:hover:not(:disabled) {
  transform: translate(2px, 2px);
  box-shadow: 2px 2px 0 0 rgba(0, 0, 0, 0.3);
}

.login-button:hover:not(:disabled)::before {
  opacity: 1;
}

.login-button:active:not(:disabled) {
  transform: translate(4px, 4px);
  box-shadow: 0 0 0 0 rgba(0, 0, 0, 0.3);
}

/* ==================== 响应式 ==================== */
@media (max-width: 768px) {
  .login-card {
    padding: 36px 24px;
    max-width: 100%;
  }
  
  .title {
    font-size: 28px;
  }
  
  .subtitle {
    font-size: 13px;
  }
  
  .logo-icon {
    width: 56px;
    height: 56px;
  }
}

@media (max-width: 480px) {
  .login-container {
    padding: 12px;
  }
  
  .login-card {
    padding: 28px 20px;
  }
  
  .gradient-orb {
    filter: blur(60px);
  }
}
</style>

