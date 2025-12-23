<template>
  <router-view v-slot="{ Component }">
    <transition name="page" mode="out-in">
      <component :is="Component" />
    </transition>
  </router-view>
</template>

<script setup>
import { onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

onMounted(() => {
  // Initialize token refresh timer if user is logged in
  if (authStore.token && authStore.refreshToken) {
    authStore.scheduleTokenRefresh()
  }
})
</script>

<style>
@import './styles/animations.css';

:root {
  /* --- 核心变量定义 --- */
  --radius-sm: 8px;
  --radius-md: 12px;
  --radius-lg: 16px;
  --radius-xl: 24px;
  
  /* 现代阴影系统 */
  --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.05);
  --shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  --shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
  --shadow-2xl: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  
  /* 颜色映射 */
  --color-bg-primary: #ffffff;
  --color-bg-secondary: #f3f4f6; /* 浅灰背景 */
  --color-bg-tertiary: #f9fafb;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  /* 使用现代无衬线字体 */
  font-family: 'Inter', system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  background-color: var(--color-bg-secondary);
  color: #374151;
}

#app {
  min-height: 100vh;
}

html {
  scroll-behavior: smooth;
}

::selection {
  background-color: var(--color-primary);
  color: white;
}
</style>

