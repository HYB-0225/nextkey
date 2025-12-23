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
  return sidebarCollapsed.value ? '80px' : '200px'
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
  background: #fff;
  color: #1f2937;
  transition: all var(--duration-normal) var(--ease-out);
  box-shadow: 0 10px 30px -10px rgba(0,0,0,0.08);
  position: relative;
  z-index: 100;
  display: flex;
  flex-direction: column;
  height: calc(100vh - 32px) !important;
  margin: 16px 0 16px 16px;
  border-radius: 24px;
  border-right: none;
}

.sidebar.collapsed {
  width: 80px !important; /* 收起时稍微宽一点，好看 */
}

/* Logo */
.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 0 16px;
  background: transparent;
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-out);
  border-bottom: 1px solid #e5e7eb;
}

.logo:hover {
  background: #f3f4f6;
}

.logo-icon {
  width: 32px;
  height: 32px;
  border-radius: 12px;
  background: linear-gradient(135deg, #FF8C42 0%, #FF6B35 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: var(--shadow-sm);
  border: none;
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
  color: #111827;
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
  background: rgba(0, 0, 0, 0.12);
  border-radius: 2px;
}

/* 菜单项基础样式 */
:deep(.el-menu-item) {
  margin: 8px 12px;
  border-radius: 12px; /* 菜单项变圆 */
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  height: 44px;
  line-height: 44px;
  border: none !important; /* 移除像素风边框 */
  color: #4b5563;
}

/* 悬停效果 */
:deep(.el-menu-item:hover) {
  color: #111827;
  background: #f3f4f6 !important;
  box-shadow: 0 6px 16px rgba(0,0,0,0.06);
  transform: translateY(-2px);
}

/* 激活状态 */
:deep(.el-menu-item.is-active) {
  color: #fff;
  background: linear-gradient(90deg, #FF8C42 0%, #FFD93D 100%) !important;
  box-shadow: 0 4px 12px rgba(255, 140, 66, 0.3);
}

/* 侧边栏底部 */
.sidebar-footer {
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  cursor: pointer;
  border-top: 1px solid #e5e7eb;
  background: #f9fafb;
  transition: all var(--duration-fast) var(--ease-out);
  border-radius: 0 0 24px 24px;
}

.sidebar-footer:hover {
  background: #eef2f7;
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
  color: #4b5563;
}

/* ==================== 主内容区 ==================== */
.main-container {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.header {
  background: transparent; /* 透明 */
  border-bottom: none;
  box-shadow: none;
  display: flex;
  align-items: center;
  padding: 16px 24px; /* 给一点顶部呼吸空间 */
  margin-bottom: 0;
  z-index: 10;
  min-height: 64px;
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
  background: transparent; /* 背景色由 body 决定 */
  padding: 24px 32px;
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
    margin: 0;
    border-radius: 0;
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
    padding: 16px;
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

