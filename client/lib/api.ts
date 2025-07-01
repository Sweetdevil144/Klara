/* eslint-disable @typescript-eslint/no-explicit-any */
const API_BASE_URL = `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/v1`;

export interface Note {
  id: string;
  title: string;
  content: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateNoteRequest {
  title: string;
  content: string;
}

export interface UpdateNoteRequest {
  title?: string;
  content?: string;
}

export interface NoteChatRequest {
  message: string;
  model: string;
  provider: string;
}

export interface NoteChatResponse {
  message: string;
  model: string;
  noteContext: string;
  suggestion?: string;
}

export interface SuggestionRequest {
  newTitle?: string;
  newContent?: string;
}

export interface AIModel {
  id: string;
  name: string;
  provider: "openai" | "gemini";
  free: boolean;
}

interface NotesResponse {
  message: string;
  notes: Note[];
  count: number;
}

interface NoteResponse {
  message: string;
  note: Note;
}

interface CreateNoteResponse {
  message: string;
  noteId: string;
  note: Note;
}

export const AI_MODELS: AIModel[] = [
  // OpenAI Models
  {
    id: "gpt-3.5-turbo",
    name: "GPT-3.5 Turbo",
    provider: "openai",
    free: true,
  },
  { id: "gpt-4", name: "GPT-4", provider: "openai", free: false },
  { id: "gpt-4-turbo", name: "GPT-4 Turbo", provider: "openai", free: false },
  { id: "gpt-4o", name: "GPT-4o", provider: "openai", free: false },
  { id: "gpt-4o-mini", name: "GPT-4o Mini", provider: "openai", free: true },

  // Gemini Models - Updated to correct model names
  {
    id: "gemini-1.5-flash",
    name: "Gemini 1.5 Flash",
    provider: "gemini",
    free: true,
  },
  {
    id: "gemini-1.5-pro",
    name: "Gemini 1.5 Pro",
    provider: "gemini",
    free: false,
  },
  {
    id: "gemini-2.0-flash",
    name: "Gemini 2.0 Flash",
    provider: "gemini",
    free: true,
  },
];

class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message);
    this.name = "ApiError";
  }
}

// Session extension utilities
let lastActivityTime = Date.now();
let sessionExtensionInterval: NodeJS.Timeout | null = null;

// Track user activity to extend sessions automatically
function trackUserActivity() {
  lastActivityTime = Date.now();
}

// Extend session automatically for active users
function startSessionExtension(getTokenFn: () => Promise<string | null>) {
  // Clear any existing interval
  if (sessionExtensionInterval) {
    clearInterval(sessionExtensionInterval);
  }

  // Extend session every 5 minutes if user has been active in the last 10 minutes
  sessionExtensionInterval = setInterval(async () => {
    const timeSinceLastActivity = Date.now() - lastActivityTime;
    const tenMinutes = 10 * 60 * 1000;

    if (timeSinceLastActivity < tenMinutes) {
      try {
        // Get a fresh token to extend the session
        await getTokenFn();
        console.log("Session extended automatically due to user activity");
      } catch (error) {
        console.warn("Failed to extend session:", error);
      }
    }
  }, 5 * 60 * 1000); // Every 5 minutes
}

// Stop session extension when user logs out
function stopSessionExtension() {
  if (sessionExtensionInterval) {
    clearInterval(sessionExtensionInterval);
    sessionExtensionInterval = null;
  }
}

// Token management utilities
let tokenRefreshPromise: Promise<string> | null = null;

async function getValidToken(
  getTokenFn: () => Promise<string | null>
): Promise<string> {
  try {
    // Track activity
    trackUserActivity();

    // If we're already refreshing a token, wait for that promise
    if (tokenRefreshPromise) {
      return await tokenRefreshPromise;
    }

    // Get a fresh token
    const token = await getTokenFn();
    if (!token) {
      throw new Error("Unable to obtain authentication token");
    }
    return token;
  } catch (error) {
    console.error("Token refresh failed:", error);
    throw new Error(
      "Authentication failed. Please refresh the page and try again."
    );
  }
}

async function apiRequest<T>(
  endpoint: string,
  options: RequestInit = {},
  getTokenFn?: () => Promise<string | null>
): Promise<T> {
  let token: string | undefined;

  if (getTokenFn) {
    token = await getValidToken(getTokenFn);

    if (!sessionExtensionInterval) {
      startSessionExtension(getTokenFn);
    }
  }

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...(token && { Authorization: `Bearer ${token}` }),
      ...options.headers,
    },
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));

    // Handle token expiration specifically
    if (
      response.status === 401 &&
      errorData.debug?.includes("Token verification failed")
    ) {
      console.warn("Token expired, attempting refresh...");

      // Clear any cached token refresh promise
      tokenRefreshPromise = null;

      // If we have a token function, try to refresh
      if (getTokenFn) {
        try {
          const newToken = await getValidToken(getTokenFn);
          // Retry the request with the new token
          const retryResponse = await fetch(`${API_BASE_URL}${endpoint}`, {
            ...options,
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${newToken}`,
              ...options.headers,
            },
          });

          if (retryResponse.ok) {
            return retryResponse.json();
          }
        } catch (refreshError) {
          console.error("Token refresh failed on retry:", refreshError);
        }
      }
    }

    throw new ApiError(
      response.status,
      errorData.message || `HTTP ${response.status}: ${response.statusText}`
    );
  }

  return response.json();
}

// Export session management utilities
export const sessionManager = {
  startSessionExtension,
  stopSessionExtension,
  trackUserActivity,
};

export const notesApi = {
  // Get all notes for the current user
  async getNotes(getTokenFn: () => Promise<string | null>): Promise<Note[]> {
    const response = await apiRequest<NotesResponse>("/notes", {}, getTokenFn);
    return response.notes || [];
  },

  // Get a specific note by ID
  async getNote(
    id: string,
    getTokenFn: () => Promise<string | null>
  ): Promise<Note> {
    const response = await apiRequest<NoteResponse>(
      `/notes/${id}`,
      {},
      getTokenFn
    );
    return response.note;
  },

  // Create a new note
  async createNote(
    data: CreateNoteRequest,
    getTokenFn: () => Promise<string | null>
  ): Promise<Note> {
    const response = await apiRequest<CreateNoteResponse>(
      "/notes",
      {
        method: "POST",
        body: JSON.stringify(data),
      },
      getTokenFn
    );
    return response.note;
  },

  // Update a note
  async updateNote(
    id: string,
    data: UpdateNoteRequest,
    getTokenFn: () => Promise<string | null>
  ): Promise<Note> {
    const response = await apiRequest<NoteResponse>(
      `/notes/${id}`,
      {
        method: "PUT",
        body: JSON.stringify(data),
      },
      getTokenFn
    );
    return response.note;
  },

  // Delete a note
  async deleteNote(
    id: string,
    getTokenFn: () => Promise<string | null>
  ): Promise<void> {
    return apiRequest<void>(
      `/notes/${id}`,
      {
        method: "DELETE",
      },
      getTokenFn
    );
  },

  // Chat with AI about a specific note
  async chatWithNote(
    id: string,
    data: NoteChatRequest,
    getTokenFn: () => Promise<string | null>
  ): Promise<NoteChatResponse> {
    return apiRequest<NoteChatResponse>(
      `/notes/${id}/chat`,
      {
        method: "POST",
        body: JSON.stringify(data),
      },
      getTokenFn
    );
  },

  // Apply AI suggestion to a note
  async applySuggestion(
    id: string,
    data: SuggestionRequest,
    getTokenFn: () => Promise<string | null>
  ): Promise<Note> {
    return apiRequest<Note>(
      `/notes/${id}/apply-suggestion`,
      {
        method: "POST",
        body: JSON.stringify(data),
      },
      getTokenFn
    );
  },
};

export const userApi = {
  // Create or sync user profile
  async createProfile(getTokenFn: () => Promise<string | null>): Promise<any> {
    return apiRequest<any>(
      "/user/profile",
      {
        method: "POST",
      },
      getTokenFn
    );
  },

  // Get user profile
  async getProfile(getTokenFn: () => Promise<string | null>): Promise<any> {
    return apiRequest<any>("/user/profile", {}, getTokenFn);
  },

  // Update API keys
  async updateAPIKeys(
    data: { openaiKey?: string; geminiKey?: string },
    getTokenFn: () => Promise<string | null>
  ): Promise<any> {
    return apiRequest<any>(
      "/user/api-keys",
      {
        method: "PUT",
        body: JSON.stringify(data),
      },
      getTokenFn
    );
  },

  // Delete API key
  async deleteAPIKey(
    keyType: string,
    getTokenFn: () => Promise<string | null>
  ): Promise<void> {
    return apiRequest<void>(
      `/user/api-keys/${keyType}`,
      {
        method: "DELETE",
      },
      getTokenFn
    );
  },
};
