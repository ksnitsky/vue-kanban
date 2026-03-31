import api from './client'

export interface AuthToken {
  token: string
}

export interface User {
  id: string
  telegram_id: number
  username: string
  first_name: string
  last_name: string
  photo_url: string
}

export const authApi = {
  getToken: async (): Promise<AuthToken> => {
    const response = await api.get<AuthToken>('/auth/token')
    return response.data
  },

  getMe: async (): Promise<User> => {
    const response = await api.get<User>('/auth/me')
    return response.data
  },

  logout: async (): Promise<void> => {
    await api.post('/auth/logout')
  },

  devLogin: async (): Promise<User> => {
    const response = await api.post<User>('/auth/dev-login')
    return response.data
  },
}
