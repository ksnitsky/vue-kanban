import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/AuthStore'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/auth',
      name: 'auth',
      component: () => import('@/views/AuthView.vue'),
      meta: { requiresAuth: false },
    },
    {
      path: '/',
      name: 'projects',
      component: () => import('@/views/ProjectsView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/projects/:projectId',
      name: 'project',
      component: () => import('@/views/ProjectView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/boards/:boardId',
      name: 'board',
      component: () => import('@/views/BoardView.vue'),
      meta: { requiresAuth: true },
    },
  ],
})

router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()

  if (to.name === 'auth') {
    if (!authStore.user && !authStore.loading) {
      try {
        await authStore.checkAuth()
      } catch (e) {
        next()
        return
      }
    }
    
    if (authStore.isAuthenticated) {
      next({ name: 'projects' })
    } else {
      next()
    }
    return
  }

  if (to.meta.requiresAuth) {
    // Если ещё не загружали пользователя - загрузить
    if (!authStore.user && !authStore.loading) {
      try {
        await authStore.checkAuth()
      } catch (e) {
        // Ошибка авторизации - редирект на auth
        next({ name: 'auth' })
        return
      }
    }
    
    if (!authStore.isAuthenticated) {
      next({ name: 'auth' })
    } else {
      next()
    }
  } else {
    next()
  }
})

export default router
