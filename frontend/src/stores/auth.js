import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { refreshToken as refreshTokenAPI } from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('admin_token') || '')
  const refreshToken = ref(localStorage.getItem('admin_refresh_token') || '')
  const tokenExpiresAt = ref(null)
  const refreshTimer = ref(null)
  
  const isLoggedIn = computed(() => !!token.value)
  
  function setToken(newToken) {
    token.value = newToken
    localStorage.setItem('admin_token', newToken)
  }
  
  function setRefreshToken(newRefreshToken) {
    refreshToken.value = newRefreshToken
    localStorage.setItem('admin_refresh_token', newRefreshToken)
  }
  
  function clearRefreshTimer() {
    if (refreshTimer.value) {
      clearTimeout(refreshTimer.value)
      refreshTimer.value = null
    }
  }
  
  async function performTokenRefresh() {
    if (!refreshToken.value) return
    
    try {
      const response = await refreshTokenAPI({
        refresh_token: refreshToken.value
      })
      
      const { access_token, refresh_token } = response
      setTokens(access_token, refresh_token)
    } catch (error) {
      console.error('Auto refresh failed:', error)
      logout()
    }
  }
  
  function scheduleTokenRefresh() {
    clearRefreshTimer()
    
    if (!tokenExpiresAt.value) return
    
    const now = Date.now()
    const expiresAt = tokenExpiresAt.value
    const timeUntilExpiry = expiresAt - now
    
    // Refresh 3 minutes (180 seconds) before expiration
    const refreshTime = timeUntilExpiry - 180000
    
    if (refreshTime > 0) {
      refreshTimer.value = setTimeout(() => {
        performTokenRefresh()
      }, refreshTime)
    } else if (timeUntilExpiry > 0) {
      // If less than 3 minutes remaining, refresh immediately
      performTokenRefresh()
    }
  }
  
  function setTokens(accessToken, newRefreshToken, expiresIn = 900) {
    setToken(accessToken)
    setRefreshToken(newRefreshToken)
    
    // Calculate expiration time (expiresIn is in seconds)
    tokenExpiresAt.value = Date.now() + (expiresIn * 1000)
    
    // Schedule the next refresh
    scheduleTokenRefresh()
  }
  
  function logout() {
    clearRefreshTimer()
    token.value = ''
    refreshToken.value = ''
    tokenExpiresAt.value = null
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
    logout,
    scheduleTokenRefresh
  }
})

