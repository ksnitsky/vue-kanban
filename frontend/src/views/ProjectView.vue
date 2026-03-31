<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useProjectStore } from '@/stores/ProjectStore'
import { boardApi, type Board } from '@/api/board'

const route = useRoute()
const router = useRouter()
const projectStore = useProjectStore()

const boards = ref<Board[]>([])
const showCreateModal = ref(false)
const newBoardName = ref('')

const projectId = route.params.projectId as string

onMounted(async () => {
  await projectStore.fetchProject(projectId)
  await fetchBoards()
})

async function fetchBoards() {
  try {
    boards.value = await boardApi.list(projectId)
  } catch (e) {
    console.error('Failed to fetch boards')
  }
}

async function createBoard() {
  if (!newBoardName.value.trim()) return

  try {
    const board = await boardApi.create({
      project_id: projectId,
      name: newBoardName.value,
    })
    showCreateModal.value = false
    newBoardName.value = ''
    boards.value.push(board)
  } catch (e) {
    console.error('Failed to create board')
  }
}

async function deleteBoard(id: string) {
  if (confirm('Are you sure you want to delete this board?')) {
    await boardApi.delete(id)
    boards.value = boards.value.filter((b) => b.id !== id)
  }
}
</script>

<template>
  <div class="project-view">
    <header>
      <button @click="router.push({ name: 'projects' })" class="btn-back">← Back</button>
      <h1>{{ projectStore.currentProject?.name || 'Project' }}</h1>
      <button @click="showCreateModal = true" class="btn-primary">New Board</button>
    </header>

    <p v-if="projectStore.currentProject?.description" class="description">
      {{ projectStore.currentProject.description }}
    </p>

    <div v-if="boards.length === 0" class="empty">
      <p>No boards yet. Create your first board!</p>
    </div>

    <div v-else class="boards-grid">
      <router-link
        v-for="board in boards"
        :key="board.id"
        :to="{ name: 'board', params: { boardId: board.id } }"
        class="board-card"
      >
        <h2>{{ board.name }}</h2>
        <p>{{ board.columns?.length || 0 }} columns</p>
        <button @click.prevent="deleteBoard(board.id)" class="btn-delete">Delete</button>
      </router-link>
    </div>

    <div v-if="showCreateModal" class="modal-overlay" @click="showCreateModal = false">
      <div class="modal" @click.stop>
        <h2>Create New Board</h2>
        <form @submit.prevent="createBoard">
          <div class="form-group">
            <label for="name">Name</label>
            <input
              id="name"
              v-model="newBoardName"
              type="text"
              placeholder="Board name"
              required
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
.project-view {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

header {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1rem;
}

h1 {
  flex: 1;
  color: #333;
}

.btn-back {
  padding: 0.5rem 1rem;
  background: #f0f0f0;
  border: none;
  border-radius: 0.25rem;
  cursor: pointer;
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

.btn-secondary {
  padding: 0.75rem 1.5rem;
  background: #f0f0f0;
  color: #333;
  border: none;
  border-radius: 0.5rem;
  cursor: pointer;
}

.description {
  color: #666;
  margin-bottom: 2rem;
}

.empty {
  text-align: center;
  padding: 3rem;
  color: #666;
}

.boards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 1.5rem;
}

.board-card {
  background: white;
  padding: 1.5rem;
  border-radius: 0.5rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  text-decoration: none;
  color: inherit;
  transition: transform 0.2s, box-shadow 0.2s;
  position: relative;
}

.board-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

.board-card h2 {
  margin: 0 0 0.5rem;
  color: #333;
}

.board-card p {
  margin: 0;
  color: #666;
  font-size: 0.9rem;
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

.form-group input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 0.25rem;
  font-size: 1rem;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 1.5rem;
}
</style>
