import axios from 'axios'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import router from '@/router'

const request = axios.create({
  baseURL: '/',
  timeout: 10000
})

let isRefreshing = false
let failedQueue = []

const processQueue = (error, token = null) => {
  failedQueue.forEach(prom => {
    if (error) {
      prom.reject(error)
    } else {
      prom.resolve(token)
    }
  })
  
  failedQueue = []
}

request.interceptors.request.use(
  config => {
    const authStore = useAuthStore()
    if (authStore.token) {
      config.headers.Authorization = `Bearer ${authStore.token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

request.interceptors.response.use(
  response => {
    const res = response.data
    
    if (res.code !== 0) {
      ElMessage.error(res.message || '请求失败')
      
      if (res.code === 401) {
        const authStore = useAuthStore()
        authStore.logout()
        router.push('/login')
      }
      
      return Promise.reject(new Error(res.message || '请求失败'))
    }
    
    return res.data
  },
  async error => {
    const originalRequest = error.config
    
    // 处理401错误，尝试刷新令牌
    if (error.response?.status === 401 && !originalRequest._retry) {
      if (originalRequest.url?.includes('/admin/refresh')) {
        // 刷新令牌本身失败，直接登出
        const authStore = useAuthStore()
        authStore.logout()
        router.push('/login')
        return Promise.reject(error)
      }
      
      if (isRefreshing) {
        // 正在刷新，将请求加入队列
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject })
        }).then(token => {
          originalRequest.headers.Authorization = `Bearer ${token}`
          return request(originalRequest)
        }).catch(err => {
          return Promise.reject(err)
        })
      }
      
      originalRequest._retry = true
      isRefreshing = true
      
      const authStore = useAuthStore()
      
      try {
        // 尝试刷新令牌
        const response = await axios.post('/admin/refresh', {
          refresh_token: authStore.refreshToken
        })
        
        const { access_token, refresh_token } = response.data.data
        
        authStore.setTokens(access_token, refresh_token)
        
        // 更新失败队列中的请求
        processQueue(null, access_token)
        
        // 重试原始请求
        originalRequest.headers.Authorization = `Bearer ${access_token}`
        return request(originalRequest)
      } catch (refreshError) {
        processQueue(refreshError, null)
        authStore.logout()
        router.push('/login')
        return Promise.reject(refreshError)
      } finally {
        isRefreshing = false
      }
    }
    
    ElMessage.error(error.message || '网络错误')
    return Promise.reject(error)
  }
)

export default request

