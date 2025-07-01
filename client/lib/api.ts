const API_BASE_URL = 'http://localhost:8080/api/v1'

export interface Note {
  id: string
  title: string
  content: string
  createdAt: string
  updatedAt: string
}

export interface CreateNoteRequest {
  title: string
  content: string
}

export interface UpdateNoteRequest {
  title?: string
  content?: string
}

export interface NoteChatRequest {
  message: string
  model: string
  provider: string
}

export interface NoteChatResponse {
  message: string
  model: string
  noteContext: string
  suggestion?: string
}

export interface SuggestionRequest {
  newTitle?: string
  newContent?: string
}

export interface AIModel {
  id: string
  name: string
  provider: 'openai' | 'gemini'
  free: boolean
}

interface NotesResponse {
  message: string
  notes: Note[]
  count: number
}

interface NoteResponse {
  message: string
  note: Note
}

interface CreateNoteResponse {
  message: string
  noteId: string
  note: Note
}

export const AI_MODELS: AIModel[] = [
  // OpenAI Models
  { id: 'gpt-3.5-turbo', name: 'GPT-3.5 Turbo', provider: 'openai', free: true },
  { id: 'gpt-4', name: 'GPT-4', provider: 'openai', free: false },
  { id: 'gpt-4-turbo', name: 'GPT-4 Turbo', provider: 'openai', free: false },
  { id: 'gpt-4o', name: 'GPT-4o', provider: 'openai', free: false },
  { id: 'gpt-4o-mini', name: 'GPT-4o Mini', provider: 'openai', free: true },
  
  // Gemini Models - Updated to correct model names
  { id: 'gemini-1.5-flash', name: 'Gemini 1.5 Flash', provider: 'gemini', free: true },
  { id: 'gemini-1.5-pro', name: 'Gemini 1.5 Pro', provider: 'gemini', free: false },
  { id: 'gemini-2.0-flash', name: 'Gemini 2.0 Flash', provider: 'gemini', free: true },
]

class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message)
    this.name = 'ApiError'
  }
}

async function apiRequest<T>(
  endpoint: string,
  options: RequestInit = {},
  token?: string
): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...(token && { 'Authorization': `Bearer ${token}` }),
      ...options.headers,
    },
  })

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}))
    throw new ApiError(
      response.status,
      errorData.message || `HTTP ${response.status}: ${response.statusText}`
    )
  }

  return response.json()
}

export const notesApi = {
  // Get all notes for the current user
  async getNotes(token: string): Promise<Note[]> {
    const response = await apiRequest<NotesResponse>('/notes', {}, token)
    return response.notes || []
  },

  // Get a specific note by ID
  async getNote(id: string, token: string): Promise<Note> {
    const response = await apiRequest<NoteResponse>(`/notes/${id}`, {}, token)
    return response.note
  },

  // Create a new note
  async createNote(data: CreateNoteRequest, token: string): Promise<Note> {
    const response = await apiRequest<CreateNoteResponse>('/notes', {
      method: 'POST',
      body: JSON.stringify(data),
    }, token)
    return response.note
  },

  // Update a note
  async updateNote(id: string, data: UpdateNoteRequest, token: string): Promise<Note> {
    const response = await apiRequest<NoteResponse>(`/notes/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }, token)
    return response.note
  },

  // Delete a note
  async deleteNote(id: string, token: string): Promise<void> {
    return apiRequest<void>(`/notes/${id}`, {
      method: 'DELETE',
    }, token)
  },

  // Chat with AI about a specific note
  async chatWithNote(id: string, data: NoteChatRequest, token: string): Promise<NoteChatResponse> {
    return apiRequest<NoteChatResponse>(`/notes/${id}/chat`, {
      method: 'POST',
      body: JSON.stringify(data),
    }, token)
  },

  // Apply AI suggestion to a note
  async applySuggestion(id: string, data: SuggestionRequest, token: string): Promise<Note> {
    return apiRequest<Note>(`/notes/${id}/apply-suggestion`, {
      method: 'POST',
      body: JSON.stringify(data),
    }, token)
  },
}

export const userApi = {
  // Create or sync user profile
  async createProfile(token: string): Promise<any> {
    return apiRequest<any>('/user/profile', {
      method: 'POST',
    }, token)
  },

  // Get user profile
  async getProfile(token: string): Promise<any> {
    return apiRequest<any>('/user/profile', {}, token)
  },

  // Update API keys
  async updateAPIKeys(data: { openaiKey?: string; geminiKey?: string }, token: string): Promise<any> {
    return apiRequest<any>('/user/api-keys', {
      method: 'PUT',
      body: JSON.stringify(data),
    }, token)
  },

  // Delete API key
  async deleteAPIKey(keyType: string, token: string): Promise<void> {
    return apiRequest<void>(`/user/api-keys/${keyType}`, {
      method: 'DELETE',
    }, token)
  },
} 