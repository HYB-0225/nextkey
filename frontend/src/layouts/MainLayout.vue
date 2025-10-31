<template>
  <el-container class="layout-container">
    <!-- 侧边栏 -->
    <el-aside 
      :width="sidebarWidth" 
      :class="['sidebar', { 'collapsed': sidebarCollapsed, 'mobile-hidden': isMobile && !mobileMenuOpen }]"
    >
      <div class="logo" @click="toggleSidebar">
        <div class="logo-icon">
          <div class="logo-gradient"></div>
        </div>
        <transition name="fade">
          <span v-if="!sidebarCollapsed" class="logo-text">NextKey</span>
        </transition>
      </div>
      
      <el-menu
        :default-active="activeMenu"
        router
        :collapse="sidebarCollapsed"
        class="sidebar-menu"
      >
        <el-menu-item index="/projects" class="menu-item">
          <el-icon><Box /></el-icon>
          <template #title>
            <span>项目管理</span>
          </template>
        </el-menu-item>
        <el-menu-item index="/cards" class="menu-item">
          <el-icon><Ticket /></el-icon>
          <template #title>
            <span>卡密管理</span>
          </template>
        </el-menu-item>
        <el-menu-item index="/cloudvars" class="menu-item">
          <el-icon><Cloudy /></el-icon>
          <template #title>
            <span>云变量</span>
          </template>
        </el-menu-item>
      </el-menu>
      
      <!-- 侧边栏底部折叠按钮 -->
      <div class="sidebar-footer" @click="toggleSidebar" v-if="!isMobile">
        <el-icon :class="['toggle-icon', { 'rotated': sidebarCollapsed }]">
          <DArrowLeft />
        </el-icon>
        <transition name="fade">
          <span v-if="!sidebarCollapsed" class="toggle-text">收起</span>
        </transition>
      </div>
    </el-aside>
    
    <!-- 移动端遮罩层 -->
    <transition name="fade">
      <div 
        v-if="isMobile && mobileMenuOpen" 
        class="mobile-overlay"
        @click="closeMobileMenu"
      ></div>
    </transition>
    
    <el-container class="main-container">
      <el-header class="header">
        <div class="header-content">
          <div class="header-left">
            <!-- 移动端菜单按钮 -->
            <button 
              v-if="isMobile" 
              class="mobile-menu-btn"
              @click="toggleMobileMenu"
            >
              <el-icon><Menu /></el-icon>
            </button>
            
            <!-- 桌面端折叠按钮 -->
            <button 
              v-else 
              class="collapse-btn"
              @click="toggleSidebar"
            >
              <el-icon><Expand /></el-icon>
            </button>
            
            <h2 class="page-title">{{ pageTitle }}</h2>
          </div>
          
          <div class="header-right">
            <el-button 
              type="danger" 
              plain 
              @click="handleLogout"
              :icon="isMobile ? '' : undefined"
            >
              {{ isMobile ? '' : '退出登录' }}
              <el-icon v-if="isMobile"><SwitchButton /></el-icon>
            </el-button>
          </div>
        </div>
      </el-header>
      
      <el-main class="main-content">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useResponsive } from '@/composables/useResponsive'
import { ElMessageBox } from 'element-plus'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const { isMobile, isTablet } = useResponsive()

const sidebarCollapsed = ref(false)
const mobileMenuOpen = ref(false)

// 监听路由变化,自动关闭移动端侧边栏
watch(() => route.path, () => {
  if (isMobile.value) {
    mobileMenuOpen.value = false
  }
})

const sidebarWidth = computed(() => {
  if (isMobile.value) return '240px'
  return sidebarCollapsed.value ? '64px' : '200px'
})

const activeMenu = computed(() => route.path)

const pageTitle = computed(() => {
  const titles = {
    '/projects': '项目管理',
    '/cards': '卡密管理',
    '/cloudvars': '云变量管理'
  }
  return titles[route.path] || 'NextKey'
})

const toggleSidebar = () => {
  if (!isMobile.value) {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }
}

const toggleMobileMenu = () => {
  mobileMenuOpen.value = !mobileMenuOpen.value
}

const closeMobileMenu = () => {
  mobileMenuOpen.value = false
}

const handleLogout = () => {
  ElMessageBox.confirm('确定要退出登录吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    authStore.logout()
    router.push('/login')
  }).catch(() => {})
}
</script>

<style scoped>
.layout-container {
  height: 100vh;
  overflow: hidden;
}

/* ==================== 侧边栏 ==================== */
.sidebar {
  background: linear-gradient(180deg, #3d2817 0%, #2a1810 100%);
  color: #fff;
  transition: all var(--duration-normal) var(--ease-out);
  box-shadow: 4px 0 0 0 rgba(0, 0, 0, 0.2);
  position: relative;
  z-index: 100;
  display: flex;
  flex-direction: column;
}

.sidebar.collapsed {
  width: 64px !important;
}

/* Logo */
.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 0 16px;
  background: rgba(0, 0, 0, 0.2);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-out);
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.logo:hover {
  background: rgba(0, 0, 0, 0.3);
}

.logo-icon {
  width: 32px;
  height: 32px;
  border-radius: 0;
  background: linear-gradient(135deg, #FF8C42 0%, #FF6B35 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 2px 2px 0 0 rgba(0, 0, 0, 0.4);
  border: 2px solid #FFD93D;
  flex-shrink: 0;
  position: relative;
  overflow: hidden;
}

.logo-gradient {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, transparent 0%, rgba(255, 255, 255, 0.2) 100%);
}

.logo-text {
  font-size: 18px;
  font-weight: 700;
  color: #fff;
  letter-spacing: 0.5px;
}

/* 菜单 */
.sidebar-menu {
  border: none;
  background: transparent;
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
}

.sidebar-menu::-webkit-scrollbar {
  width: 4px;
}

.sidebar-menu::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 2px;
}

/* 菜单项基础样式 */
:deep(.el-menu-item) {
  color: rgba(255, 255, 255, 0.65);
  margin: 4px 8px;
  border-radius: 0;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  height: 48px;
  line-height: 48px;
  border: 2px solid transparent;
  position: relative;
  overflow: hidden;
}

/* 悬停效果 */
:deep(.el-menu-item:hover) {
  color: #fff;
  background: rgba(255, 140, 66, 0.15) !important;
  border-color: rgba(255, 140, 66, 0.3);
  transform: translateX(2px);
}

/* 激活状态 - 流淌效果 */
:deep(.el-menu-item.is-active) {
  color: #fff;
  background: linear-gradient(90deg, #FF8C42 0%, #FFD93D 100%) !important;
  box-shadow: 2px 2px 0 0 rgba(0, 0, 0, 0.3);
  transform: translateX(0);
  border: 2px solid #FF6B35;
  animation: menu-activate 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

/* 激活动画 - 从左到右流淌 */
@keyframes menu-activate {
  0% {
    background: linear-gradient(90deg, transparent 0%, transparent 100%);
    box-shadow: none;
  }
  50% {
    background: linear-gradient(90deg, #FF8C42 0%, transparent 50%, transparent 100%);
  }
  100% {
    background: linear-gradient(90deg, #FF8C42 0%, #FFD93D 100%);
    box-shadow: 2px 2px 0 0 rgba(0, 0, 0, 0.3);
  }
}

/* 激活指示器 */
:deep(.el-menu-item.is-active::before) {
  content: '';
  position: absolute;
  left: -2px;
  top: 50%;
  transform: translateY(-50%);
  width: 4px;
  height: 24px;
  background: #FFD93D;
  border-radius: 0;
  box-shadow: 1px 0 0 0 #FF6B35;
  animation: indicator-slide 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

/* 指示器滑动动画 */
@keyframes indicator-slide {
  0% {
    height: 0;
    opacity: 0;
  }
  50% {
    height: 32px;
    opacity: 0.5;
  }
  100% {
    height: 24px;
    opacity: 1;
  }
}

/* 激活后的流光效果 */
:deep(.el-menu-item.is-active::after) {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent 0%, rgba(255, 255, 255, 0.3) 50%, transparent 100%);
  animation: shimmer 0.8s ease-out;
}

@keyframes shimmer {
  0% {
    left: -100%;
  }
  100% {
    left: 100%;
  }
}

/* 侧边栏底部 */
.sidebar-footer {
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  cursor: pointer;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
  background: rgba(0, 0, 0, 0.2);
  transition: all var(--duration-fast) var(--ease-out);
}

.sidebar-footer:hover {
  background: rgba(0, 0, 0, 0.3);
}

.toggle-icon {
  font-size: 18px;
  transition: transform var(--duration-normal) var(--ease-out);
}

.toggle-icon.rotated {
  transform: rotate(180deg);
}

.toggle-text {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.65);
}

/* ==================== 主内容区 ==================== */
.main-container {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.header {
  background: #fff;
  border-bottom: 1px solid var(--color-border-light);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  display: flex;
  align-items: center;
  padding: 0 24px;
  z-index: 10;
  height: 64px;
}

.header-content {
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.mobile-menu-btn,
.collapse-btn {
  width: 40px;
  height: 40px;
  border: none;
  background: transparent;
  border-radius: var(--radius-md);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all var(--duration-fast) var(--ease-out);
  color: var(--color-text-secondary);
}

.mobile-menu-btn:hover,
.collapse-btn:hover {
  background: var(--color-bg-tertiary);
  color: var(--color-text-primary);
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.main-content {
  background: var(--color-bg-secondary);
  padding: 24px;
  overflow-y: auto;
  flex: 1;
  display: flex;
  justify-content: center;
}

.main-content > * {
  width: 100%;
  max-width: 1400px;
}

/* ==================== 平板适配 ==================== */
@media (min-width: 769px) and (max-width: 1023px) {
  .main-content {
    padding: 20px;
  }
}

/* ==================== 移动端适配 ==================== */
@media (max-width: 768px) {
  .sidebar {
    position: fixed;
    left: 0;
    top: 0;
    height: 100vh;
    transform: translateX(0);
    transition: transform var(--duration-normal) cubic-bezier(0.4, 0, 0.2, 1);
    box-shadow: 4px 0 16px rgba(0, 0, 0, 0.2);
  }
  
  .sidebar.mobile-hidden {
    transform: translateX(-100%);
  }
  
  .mobile-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.5);
    z-index: 99;
    animation: fade-in var(--duration-fast) ease-out;
  }
  
  .header {
    padding: 0 16px;
  }
  
  .page-title {
    font-size: 18px;
  }
  
  .main-content {
    padding: 16px;
  }
  
  /* 优化触摸体验 */
  .sidebar-menu {
    -webkit-overflow-scrolling: touch;
  }
  
  /* 移动端菜单项点击反馈 */
  :deep(.el-menu-item:active) {
    opacity: 0.8;
  }
}

/* ==================== 过渡动画 ==================== */
.fade-enter-active,
.fade-leave-active {
  transition: opacity var(--duration-fast) var(--ease-out);
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>

