import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('admin_token') || '')
  const refreshToken = ref(localStorage.getItem('admin_refresh_token') || '')
  
  const isLoggedIn = computed(() => !!token.value)
  
  function setToken(newToken) {
    token.value = newToken
    localStorage.setItem('admin_token', newToken)
  }
  
  function setRefreshToken(newRefreshToken) {
    refreshToken.value = newRefreshToken
    localStorage.setItem('admin_refresh_token', newRefreshToken)
  }
  
  function setTokens(accessToken, newRefreshToken) {
    setToken(accessToken)
    setRefreshToken(newRefreshToken)
  }
  
  function logout() {
    token.value = ''
    refreshToken.value = ''
    localStorage.removeItem('admin_token')
    localStorage.removeItem('admin_refresh_token')
  }
  
  return {
    token,
    refreshToken,
    isLoggedIn,
    setToken,
    setRefreshToken,
    setTokens,
    logout
  }
})

