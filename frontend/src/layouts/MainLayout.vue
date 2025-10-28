<template>
  <el-container class="layout-container">
    <el-aside width="200px">
      <div class="logo">NextKey</div>
      <el-menu
        :default-active="activeMenu"
        router
        class="el-menu-vertical"
      >
        <el-menu-item index="/projects">
          <el-icon><Box /></el-icon>
          <span>项目管理</span>
        </el-menu-item>
        <el-menu-item index="/cards">
          <el-icon><Ticket /></el-icon>
          <span>卡密管理</span>
        </el-menu-item>
        <el-menu-item index="/cloudvars">
          <el-icon><Cloudy /></el-icon>
          <span>云变量</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    
    <el-container>
      <el-header>
        <div class="header-content">
          <h2>{{ pageTitle }}</h2>
          <el-button type="danger" plain @click="handleLogout">退出登录</el-button>
        </div>
      </el-header>
      
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { ElMessageBox } from 'element-plus'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const activeMenu = computed(() => route.path)

const pageTitle = computed(() => {
  const titles = {
    '/projects': '项目管理',
    '/cards': '卡密管理',
    '/cloudvars': '云变量管理'
  }
  return titles[route.path] || 'NextKey'
})

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
}

.el-aside {
  background-color: #001529;
  color: #fff;
}

.logo {
  height: 60px;
  line-height: 60px;
  text-align: center;
  font-size: 20px;
  font-weight: bold;
  color: #fff;
  background-color: #002140;
}

.el-menu {
  border: none;
  background-color: #001529;
}

:deep(.el-menu-item) {
  color: rgba(255, 255, 255, 0.65);
}

:deep(.el-menu-item:hover),
:deep(.el-menu-item.is-active) {
  color: #fff;
  background-color: #1890ff !important;
}

.el-header {
  background-color: #fff;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  align-items: center;
  padding: 0 20px;
}

.header-content {
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.el-main {
  background-color: #f0f2f5;
  padding: 20px;
}
</style>

