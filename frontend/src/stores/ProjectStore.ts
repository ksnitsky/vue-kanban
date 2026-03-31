import { defineStore } from 'pinia'
import { ref } from 'vue'
import { projectApi, type Project } from '@/api/project'

export const useProjectStore = defineStore('project', () => {
  const projects = ref<Project[]>([])
  const currentProject = ref<Project | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchProjects() {
    try {
      loading.value = true
      projects.value = await projectApi.list()
    } catch (e) {
      error.value = 'Failed to fetch projects'
    } finally {
      loading.value = false
    }
  }

  async function fetchProject(id: string) {
    try {
      loading.value = true
      currentProject.value = await projectApi.get(id)
    } catch (e) {
      error.value = 'Failed to fetch project'
    } finally {
      loading.value = false
    }
  }

  async function createProject(name: string, description?: string) {
    try {
      loading.value = true
      const project = await projectApi.create({ name, description })
      projects.value.push(project)
      return project
    } catch (e) {
      error.value = 'Failed to create project'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function deleteProject(id: string) {
    try {
      await projectApi.delete(id)
      projects.value = projects.value.filter((p) => p.id !== id)
    } catch (e) {
      error.value = 'Failed to delete project'
    }
  }

  return {
    projects,
    currentProject,
    loading,
    error,
    fetchProjects,
    fetchProject,
    createProject,
    deleteProject,
  }
})
