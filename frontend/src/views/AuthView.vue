<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/AuthStore'
import { authApi } from '@/api/auth'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const error = ref<string | null>(null)
const authUrl = ref<string | null>(null)
const isDev = import.meta.env.DEV

async function handleLogin() {
  if (isDev) {
    await devLogin()
  } else {
    await startAuthFlow()
  }
}

async function devLogin() {
  try {
    loading.value = true
    error.value = null
    await authStore.devLogin()
    if (authStore.isAuthenticated) {
      router.push({ name: 'projects' })
    }
  } catch (e) {
    error.value = 'Failed to login'
  } finally {
    loading.value = false
  }
}

async function startAuthFlow() {
  try {
    loading.value = true
    error.value = null
    const { token } = await authApi.getToken()
    
    const botUsername = import.meta.env.VITE_TELEGRAM_BOT_USERNAME || 'your_bot'
    authUrl.value = `https://t.me/${botUsername}?start=${token}`

    await authStore.connectWebSocket(token)
    router.push({ name: 'projects' })
  } catch (e) {
    error.value = 'Failed to start authentication'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="auth-view">
    <div class="auth-container">
      <h1>Kanban</h1>
      
      <div v-if="loading" class="loading">
        <p>Connecting...</p>
      </div>

      <div v-else-if="error" class="error">
        <p>{{ error }}</p>
        <button @click="handleLogin" class="btn-primary">Try again</button>
      </div>

      <div v-else-if="authUrl" class="auth-link">
        <p>Click the button below to authenticate via Telegram</p>
        <a :href="authUrl" target="_blank" class="telegram-btn">
          <svg viewBox="0 0 24 24" width="24" height="24">
            <path fill="currentColor" d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm4.64 6.8c-.15 1.58-.8 5.42-1.13 7.19-.14.75-.42 1-.68 1.03-.58.05-1.02-.38-1.58-.75-.88-.58-1.38-.94-2.23-1.5-.99-.65-.35-1.01.22-1.59.15-.15 2.71-2.48 2.76-2.69a.2.2 0 00-.05-.18c-.06-.05-.14-.03-.21-.02-.09.02-1.49.95-4.22 2.79-.4.27-.76.41-1.08.4-.36-.01-1.04-.2-1.55-.37-.63-.2-1.12-.31-1.08-.66.02-.18.27-.36.74-.55 2.92-1.27 4.86-2.11 5.83-2.51 2.78-1.16 3.35-1.36 3.73-1.36.08 0 .27.02.39.12.1.08.13.19.14.27-.01.06.01.24 0 .38z"/>
          </svg>
          Open Telegram
        </a>
        <p class="hint">Waiting for authentication...</p>
      </div>

      <div v-else class="auth-initial">
        <p class="auth-description">Sign in to manage your projects</p>
        <button @click="handleLogin" class="telegram-btn">
          <svg viewBox="0 0 24 24" width="24" height="24">
            <path fill="currentColor" d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm4.64 6.8c-.15 1.58-.8 5.42-1.13 7.19-.14.75-.42 1-.68 1.03-.58.05-1.02-.38-1.58-.75-.88-.58-1.38-.94-2.23-1.5-.99-.65-.35-1.01.22-1.59.15-.15 2.71-2.48 2.76-2.69a.2.2 0 00-.05-.18c-.06-.05-.14-.03-.21-.02-.09.02-1.49.95-4.22 2.79-.4.27-.76.41-1.08.4-.36-.01-1.04-.2-1.55-.37-.63-.2-1.12-.31-1.08-.66.02-.18.27-.36.74-.55 2.92-1.27 4.86-2.11 5.83-2.51 2.78-1.16 3.35-1.36 3.73-1.36.08 0 .27.02.39.12.1.08.13.19.14.27-.01.06.01.24 0 .38z"/>
          </svg>
          Login with Telegram
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.auth-view {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.auth-container {
  background: white;
  padding: 3rem;
  border-radius: 1rem;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
  text-align: center;
  max-width: 400px;
  width: 90%;
}

h1 {
  margin-bottom: 2rem;
  color: #333;
}

.loading {
  color: #666;
}

.error {
  color: #e74c3c;
}

.error button {
  margin-top: 1rem;
  padding: 0.75rem 1.5rem;
  background: #667eea;
  color: white;
  border: none;
  border-radius: 0.5rem;
  cursor: pointer;
  font-weight: 600;
}

.auth-description {
  color: #666;
  margin-bottom: 1.5rem;
}

.telegram-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 1rem 2rem;
  background: #0088cc;
  color: white;
  text-decoration: none;
  border-radius: 0.5rem;
  font-weight: 600;
  border: none;
  cursor: pointer;
  font-size: 1rem;
}

.telegram-btn:hover {
  background: #0077b5;
}

.auth-link .telegram-btn {
  margin: 1rem 0;
}

.hint {
  color: #999;
  font-size: 0.9rem;
  margin-top: 1rem;
}
</style>