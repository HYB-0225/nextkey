import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue')
  },
  {
    path: '/unbind/:unbindSlug',
    name: 'ProjectUnbind',
    component: () => import('@/views/ProjectUnbind.vue')
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
  const loggedIn = !!localStorage.getItem('admin_token')
  if (to.meta.requiresAuth && !loggedIn) {
    return next('/login')
  }
  if (to.path === '/login' && loggedIn) {
    return next('/')
  }
  next()
})

export default router

