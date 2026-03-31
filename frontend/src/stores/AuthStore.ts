import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi, type User } from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const ws = ref<WebSocket | null>(null)

  const isAuthenticated = computed(() => user.value !== null)

  async function checkAuth() {
    try {
      loading.value = true
      user.value = await authApi.getMe()
    } catch {
      user.value = null
    } finally {
      loading.value = false
    }
  }

  async function devLogin() {
    try {
      loading.value = true
      user.value = await authApi.devLogin()
    } catch {
      user.value = null
    } finally {
      loading.value = false
    }
  }

  async function logout() {
    await authApi.logout()
    user.value = null
  }

  function connectWebSocket(token: string): Promise<void> {
    return new Promise((resolve, reject) => {
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
      const wsUrl = `${protocol}//${window.location.host}/api/auth/ws?token=${token}`
      
      ws.value = new WebSocket(wsUrl)

      ws.value.onmessage = (event) => {
        const data = JSON.parse(event.data)
        if (data.type === 'auth_success') {
          user.value = data.user
          ws.value?.close()
          resolve()
        }
      }

      ws.value.onerror = (error) => {
        reject(error)
      }

      ws.value.onclose = () => {
        ws.value = null
      }
    })
  }

  return {
    user,
    loading,
    error,
    isAuthenticated,
    checkAuth,
    devLogin,
    logout,
    connectWebSocket,
  }
})
