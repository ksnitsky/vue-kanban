import api from './client'

export interface Card {
  id: string
  column_id: string
  content: string
  position: number
  created_at: string
  updated_at: string
}

export interface Column {
  id: string
  board_id: string
  title: string
  position: number
  cards: Card[]
  created_at: string
  updated_at: string
}

export interface Board {
  id: string
  project_id: string
  name: string
  columns: Column[]
  created_at: string
  updated_at: string
}

export interface CreateBoardRequest {
  project_id: string
  name: string
}

export interface CreateColumnRequest {
  board_id: string
  title: string
}

export interface CreateCardRequest {
  column_id: string
  content: string
}

export interface MoveCardRequest {
  card_id: string
  target_column_id: string
  position: number
}

export const boardApi = {
  list: async (projectId: string): Promise<Board[]> => {
    const response = await api.get<Board[]>(`/projects/${projectId}/boards`)
    return response.data
  },

  get: async (id: string): Promise<Board> => {
    const response = await api.get<Board>(`/boards/${id}`)
    return response.data
  },

  create: async (data: CreateBoardRequest): Promise<Board> => {
    const response = await api.post<Board>('/boards', data)
    return response.data
  },

  update: async (id: string, name: string): Promise<Board> => {
    const response = await api.put<Board>(`/boards/${id}`, { name })
    return response.data
  },

  delete: async (id: string): Promise<void> => {
    await api.delete(`/boards/${id}`)
  },

  createColumn: async (data: CreateColumnRequest): Promise<Column> => {
    const response = await api.post<Column>('/columns', data)
    return response.data
  },

  updateColumn: async (id: string, title: string): Promise<Column> => {
    const response = await api.put<Column>(`/columns/${id}`, { title })
    return response.data
  },

  deleteColumn: async (id: string): Promise<void> => {
    await api.delete(`/columns/${id}`)
  },

  reorderColumns: async (boardId: string, columnIds: string[]): Promise<void> => {
    await api.put('/columns/reorder', { board_id: boardId, column_ids: columnIds })
  },

  createCard: async (data: CreateCardRequest): Promise<Card> => {
    const response = await api.post<Card>('/cards', data)
    return response.data
  },

  updateCard: async (id: string, content: string): Promise<Card> => {
    const response = await api.put<Card>(`/cards/${id}`, { content })
    return response.data
  },

  deleteCard: async (id: string): Promise<void> => {
    await api.delete(`/cards/${id}`)
  },

  moveCard: async (data: MoveCardRequest): Promise<void> => {
    await api.put('/cards/move', data)
  },

  reorderCards: async (columnId: string, cardIds: string[]): Promise<void> => {
    await api.put('/cards/reorder', { column_id: columnId, card_ids: cardIds })
  },
}
