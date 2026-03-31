<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useProjectStore } from '@/stores/ProjectStore'
import { useAuthStore } from '@/stores/AuthStore'

const router = useRouter()
const projectStore = useProjectStore()
const authStore = useAuthStore()

const showCreateModal = ref(false)
const newProjectName = ref('')
const newProjectDescription = ref('')

onMounted(async () => {
  await projectStore.fetchProjects()
})

async function createProject() {
  if (!newProjectName.value.trim()) return
  
  try {
    const project = await projectStore.createProject(
      newProjectName.value,
      newProjectDescription.value
    )
    showCreateModal.value = false
    newProjectName.value = ''
    newProjectDescription.value = ''
    router.push({ name: 'project', params: { projectId: project.id } })
  } catch (e) {
    console.error('Failed to create project')
  }
}

async function deleteProject(id: string) {
  if (confirm('Are you sure you want to delete this project?')) {
    await projectStore.deleteProject(id)
  }
}

async function logout() {
  await authStore.logout()
  router.push({ name: 'auth' })
}
</script>

<template>
  <div class="projects-view">
    <header>
      <h1>My Projects</h1>
      <div class="header-actions">
        <button @click="showCreateModal = true" class="btn-primary">New Project</button>
        <button @click="logout" class="btn-secondary">Logout</button>
      </div>
    </header>

    <div v-if="projectStore.loading" class="loading">Loading...</div>

    <div v-else-if="projectStore.projects.length === 0" class="empty">
      <p>No projects yet. Create your first project!</p>
    </div>

    <div v-else class="projects-grid">
      <router-link
        v-for="project in projectStore.projects"
        :key="project.id"
        :to="{ name: 'project', params: { projectId: project.id } }"
        class="project-card"
      >
        <h2>{{ project.name }}</h2>
        <p v-if="project.description">{{ project.description }}</p>
        <button @click.prevent="deleteProject(project.id)" class="btn-delete">Delete</button>
      </router-link>
    </div>

    <div v-if="showCreateModal" class="modal-overlay" @click="showCreateModal = false">
      <div class="modal" @click.stop>
        <h2>Create New Project</h2>
        <form @submit.prevent="createProject">
          <div class="form-group">
            <label for="name">Name</label>
            <input
              id="name"
              v-model="newProjectName"
              type="text"
              placeholder="Project name"
              required
            />
          </div>
          <div class="form-group">
            <label for="description">Description (optional)</label>
            <textarea
              id="description"
              v-model="newProjectDescription"
              placeholder="Project description"
            />
          </div>
          <div class="modal-actions">
            <button type="button" @click="showCreateModal = false" class="btn-secondary">
              Cancel
            </button>
            <button type="submit" class="btn-primary">Create</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>
.projects-view {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

h1 {
  color: #333;
}

.header-actions {
  display: flex;
  gap: 1rem;
}

.btn-primary {
  padding: 0.75rem 1.5rem;
  background: #667eea;
  color: white;
  border: none;
  border-radius: 0.5rem;
  cursor: pointer;
  font-weight: 600;
}

.btn-primary:hover {
  background: #5a6fd6;
}

.btn-secondary {
  padding: 0.75rem 1.5rem;
  background: #f0f0f0;
  color: #333;
  border: none;
  border-radius: 0.5rem;
  cursor: pointer;
}

.btn-secondary:hover {
  background: #e0e0e0;
}

.loading,
.empty {
  text-align: center;
  padding: 3rem;
  color: #666;
}

.projects-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
}

.project-card {
  background: white;
  padding: 1.5rem;
  border-radius: 0.5rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  text-decoration: none;
  color: inherit;
  transition: transform 0.2s, box-shadow 0.2s;
  position: relative;
}

.project-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

.project-card h2 {
  margin: 0 0 0.5rem;
  color: #333;
}

.project-card p {
  margin: 0;
  color: #666;
}

.btn-delete {
  position: absolute;
  top: 1rem;
  right: 1rem;
  padding: 0.25rem 0.5rem;
  background: #e74c3c;
  color: white;
  border: none;
  border-radius: 0.25rem;
  cursor: pointer;
  font-size: 0.8rem;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal {
  background: white;
  padding: 2rem;
  border-radius: 0.5rem;
  width: 90%;
  max-width: 500px;
}

.modal h2 {
  margin: 0 0 1.5rem;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 600;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 0.25rem;
  font-size: 1rem;
}

.form-group textarea {
  min-height: 100px;
  resize: vertical;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 1.5rem;
}
</style>
