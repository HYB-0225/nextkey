import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue')
  },
  {
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'),
    redirect: '/projects',
    meta: { requiresAuth: true },
    children: [
      {
        path: '/projects',
        name: 'Projects',
        component: () => import('@/views/Projects.vue')
      },
      {
        path: '/cards/:projectId?',
        name: 'Cards',
        component: () => import('@/views/Cards.vue')
      },
      {
        path: '/cloudvars/:projectId?',
        name: 'CloudVars',
        component: () => import('@/views/CloudVars.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    next('/login')
  } else if (to.path === '/login' && authStore.isLoggedIn) {
    next('/')
  } else {
    next()
  }
})

export default router

