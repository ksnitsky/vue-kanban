<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { boardApi, type Board, type Column, type Card } from '@/api/board'
import draggable from 'vuedraggable'
import NewBlock from '@/components/_newBlock.vue'

const route = useRoute()
const router = useRouter()

const board = ref<Board | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)
const dragState = ref<{
  fromColumnId: string | null
  cardId: string | null
  oldIndex: number | null
}>({ fromColumnId: null, cardId: null, oldIndex: null })

const boardId = route.params.boardId as string

onMounted(async () => {
  await fetchBoard()
})

async function fetchBoard() {
  try {
    loading.value = true
    board.value = await boardApi.get(boardId)
  } catch (e) {
    error.value = 'Не удалось загрузить доску'
  } finally {
    loading.value = false
  }
}

const columns = computed({
  get: () => board.value?.columns || [],
  set: async (newColumns) => {
    if (!board.value) return
    const columnIds = newColumns.map((c) => c.id)
    await boardApi.reorderColumns(board.value.id, columnIds)
    board.value.columns = newColumns
  },
})

async function createColumn(title: string) {
  if (!board.value) return
  
  try {
    const column = await boardApi.createColumn({
      board_id: board.value.id,
      title,
    })
    board.value.columns.push({ ...column, cards: [] })
  } catch (e) {
    console.error('Failed to create column')
  }
}

async function deleteColumn(columnId: string) {
  if (!confirm('Удалить эту колонку?')) return
  
  try {
    await boardApi.deleteColumn(columnId)
    if (board.value) {
      board.value.columns = board.value.columns.filter((c) => c.id !== columnId)
    }
  } catch (e) {
    console.error('Failed to delete column')
  }
}

async function createCard(columnId: string, content: string) {
  try {
    const card = await boardApi.createCard({
      column_id: columnId,
      content,
    })
    const column = board.value?.columns.find((c) => c.id === columnId)
    if (column) {
      column.cards.push(card)
    }
  } catch (e) {
    console.error('Failed to create card')
  }
}

async function deleteCard(columnId: string, cardId: string) {
  try {
    await boardApi.deleteCard(cardId)
    const column = board.value?.columns.find((c) => c.id === columnId)
    if (column) {
      column.cards = column.cards.filter((c) => c.id !== cardId)
    }
  } catch (e) {
    console.error('Failed to delete card')
  }
}

function onDragStart(evt: any, columnId: string) {
  dragState.value = {
    fromColumnId: columnId,
    cardId: evt.item?.dataset?.cardId || null,
    oldIndex: evt.oldIndex
  }
}

async function onCardChange(evt: any, columnId: string) {
  if (evt.added) {
    if (dragState.value.fromColumnId && dragState.value.fromColumnId !== columnId) {
      try {
        await boardApi.moveCard({
          card_id: dragState.value.cardId!,
          target_column_id: columnId,
          position: evt.added.newIndex
        })
      } catch (e) {
        console.error('Failed to move card')
        await fetchBoard()
      }
    }
  } else if (evt.moved) {
    const column = board.value?.columns.find(c => c.id === columnId)
    if (column) {
      try {
        const cardIds = column.cards.map(c => c.id)
        await boardApi.reorderCards(columnId, cardIds)
      } catch (e) {
        console.error('Failed to reorder cards')
      }
    }
  }
}

async function onColumnChange(evt: any) {
  if (evt.moved && board.value) {
    try {
      const columnIds = board.value.columns.map(c => c.id)
      await boardApi.reorderColumns(board.value.id, columnIds)
    } catch (e) {
      console.error('Failed to reorder columns')
    }
  }
}

function goBack() {
  if (board.value?.project_id) {
    router.push({ name: 'project', params: { projectId: board.value.project_id } })
  } else {
    router.push({ name: 'projects' })
  }
}
</script>

<template>
  <div class="board-page">
    <header class="board-header">
      <button @click="goBack" class="btn-back">← Назад</button>
      <h1>{{ board?.name || 'Доска' }}</h1>
    </header>

    <div v-if="loading" class="loading">Загрузка...</div>
    <div v-else-if="error" class="error">{{ error }}</div>

    <div v-else class="board-wrapper">
      <draggable
        class="columns"
        handle=".handle"
        :list="columns"
        group="columns"
        item-key="id"
        animation="250"
        ghost-class="ghost"
        @change="onColumnChange"
      >
        <template #item="{ element: column }">
          <div class="column-item">
            <div class="column-header">
              <h2 class="handle">{{ column.title }}</h2>
              <button @click="deleteColumn(column.id)" class="btn-delete">×</button>
            </div>
            
            <draggable
              tag="ul"
              :list="column.cards"
              group="cards"
              item-key="id"
              class="cards-list"
              animation="250"
              ghost-class="ghost"
              @start="(e: any) => onDragStart(e, column.id)"
              @change="(e: any) => onCardChange(e, column.id)"
            >
              <template #item="{ element: card }">
                <li class="card-item" :data-card-id="card.id">
                  {{ card.content }}
                  <button @click="deleteCard(column.id, card.id)" class="btn-card-delete">×</button>
                </li>
              </template>
            </draggable>

            <NewBlock
              :is-col="false"
              placeholder="Добавить карточку"
              @create-card="(content: string) => createCard(column.id, content)"
            />
          </div>
        </template>
      </draggable>

      <div class="column-item new-column">
        <NewBlock
          :is-col="true"
          placeholder="Добавить колонку"
          @create-col="createColumn"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.board-page {
  width: 100%;
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: url('@/assets/background.webp') no-repeat center;
  background-size: cover;
}

.board-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem 1.25rem;
  background: rgba(255, 255, 255, 0.95);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.board-header h1 {
  margin: 0;
  font-size: 1.25rem;
  color: #333;
}

.btn-back {
  padding: 0.5rem 1rem;
  background: #f0f0f0;
  border: none;
  border-radius: 0.25rem;
  cursor: pointer;
  font-size: 0.875rem;
}

.btn-back:hover {
  background: #e0e0e0;
}

.loading,
.error {
  text-align: center;
  padding: 3rem;
  color: #666;
}

.error {
  color: #e74c3c;
}

.board-wrapper {
  flex: 1;
  padding: 1.25rem;
  overflow-x: auto;
  display: flex;
  gap: 0.75rem;
}

.columns {
  display: flex;
  gap: 0.75rem;
}

.column-item {
  min-width: 18rem;
  max-width: 18rem;
  height: fit-content;
  max-height: calc(100vh - 5rem);
  background-color: #DFE3E6;
  border-radius: 0.1875rem;
  display: flex;
  flex-direction: column;
}

.column-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem;
}

.column-header h2 {
  margin: 0;
  font-size: 1rem;
  font-weight: bold;
  cursor: move;
  flex: 1;
}

.btn-delete {
  background: none;
  border: none;
  font-size: 1.25rem;
  color: #6B808C;
  cursor: pointer;
  padding: 0;
  line-height: 1;
}

.btn-delete:hover {
  color: #e74c3c;
}

.cards-list {
  list-style: none;
  margin: 0;
  padding: 0 0.75rem;
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  scrollbar-width: none;
}

.cards-list::-webkit-scrollbar {
  display: none;
}

.card-item {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 0.75rem;
  background-color: #FFFFFF;
  border-radius: 0.1875rem;
  box-shadow: 0 1px 4px rgba(9, 45, 66, 0.25);
  cursor: move;
  gap: 0.5rem;
}

.card-item:hover .btn-card-delete {
  opacity: 1;
}

.btn-card-delete {
  background: none;
  border: none;
  font-size: 1rem;
  color: #999;
  cursor: pointer;
  padding: 0;
  opacity: 0;
  transition: opacity 0.2s;
}

.btn-card-delete:hover {
  color: #e74c3c;
}

.new-column {
  background: transparent;
}

.ghost {
  opacity: 0.5;
  background-color: #C9D2D9 !important;
}
</style>