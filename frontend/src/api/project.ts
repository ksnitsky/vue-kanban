import api from './client'

export interface Project {
  id: string
  user_id: string
  name: string
  description: string
  created_at: string
  updated_at: string
}

export interface CreateProjectRequest {
  name: string
  description?: string
}

export const projectApi = {
  list: async (): Promise<Project[]> => {
    const response = await api.get<Project[]>('/projects')
    return response.data
  },

  get: async (id: string): Promise<Project> => {
    const response = await api.get<Project>(`/projects/${id}`)
    return response.data
  },

  create: async (data: CreateProjectRequest): Promise<Project> => {
    const response = await api.post<Project>('/projects', data)
    return response.data
  },

  update: async (id: string, data: CreateProjectRequest): Promise<Project> => {
    const response = await api.put<Project>(`/projects/${id}`, data)
    return response.data
  },

  delete: async (id: string): Promise<void> => {
    await api.delete(`/projects/${id}`)
  },
}
