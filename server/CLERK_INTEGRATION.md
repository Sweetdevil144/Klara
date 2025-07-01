# Clerk Authentication Integration

This document explains how to use the Clerk authentication system integrated into your AI Notes App.

## Overview

The app now uses **Clerk** for authentication instead of traditional username/password authentication. Key features:

- 🔐 **Google OAuth** integration via Clerk
- 🔑 **JWT Session Tokens** for API authentication
- 👤 **User Profile Management** with API key storage
- 📝 **Secure Note Management** tied to authenticated users
- 🔄 **MongoDB Aggregation** for efficient data retrieval
- 🤖 **AI Chat Integration** with OpenAI/Gemini and memory management
- 🧠 **Mem0 Memory System** for conversational context

## Environment Variables

Ensure your `.env` file contains:

```env
CLERK_SECRET_KEY=sk_live_your_secret_key_here
CLERK_PUBLISHABLE_KEY=pk_live_your_publishable_key_here
JWT_SECRET=your_jwt_secret_here  # Optional, for additional JWT operations
MEM0_API_KEY=your_mem0_api_key_here  # For memory management
```

## User Model Structure

The updated user model supports:

```go
type User struct {
    ID           primitive.ObjectID   `json:"id"`
    ClerkID      string               `json:"clerkId"`      // Clerk user identifier
    Email        string               `json:"email"`
    Username     string               `json:"username"`
    FirstName    string               `json:"firstName"`
    LastName     string               `json:"lastName"`
    OpenAIKey    string               `json:"openaiKey"`    // Encrypted API key
    GeminiKey    string               `json:"geminiKey"`    // Encrypted API key
    CreatedAt    time.Time            `json:"createdAt"`
    UpdatedAt    time.Time            `json:"updatedAt"`
    NoteIds      []primitive.ObjectID `json:"noteIds"`      // References to notes
}
```

## API Endpoints

### Authentication

All protected routes require the `Authorization` header:

```bash
Authorization: Bearer <clerk_session_token>
```

### User Management

#### Create/Sync User Profile

```http
POST /api/v1/user/profile
Content-Type: application/json
Authorization: Bearer <token>

{
  "email": "user@example.com",
  "username": "johndoe",
  "firstName": "John",
  "lastName": "Doe"
}
```

#### Get User Profile

```http
GET /api/v1/user/profile
Authorization: Bearer <token>

Response:
{
  "message": "User profile retrieved successfully",
  "user": {
    "id": "...",
    "clerkId": "user_123",
    "email": "user@example.com",
    "username": "johndoe",
    "firstName": "John",
    "lastName": "Doe",
    "hasOpenaiKey": true,
    "hasGeminiKey": false,
    "createdAt": "2023-12-01T10:00:00Z",
    "noteIds": ["note1", "note2"]
  }
}
```

#### Update API Keys

```http
PUT /api/v1/user/api-keys
Content-Type: application/json
Authorization: Bearer <token>

{
  "openaiKey": "sk-...",
  "geminiKey": "AI..."
}
```

#### Delete API Key

```http
DELETE /api/v1/user/api-keys/{keyType}
Authorization: Bearer <token>

# keyType can be: openai, gemini
```

### Notes Management

#### Create Note

```http
POST /api/v1/notes
Content-Type: application/json
Authorization: Bearer <token>

{
  "title": "My AI Note",
  "content": "This is the content of my note..."
}
```

#### Get My Notes

```http
GET /api/v1/notes
Authorization: Bearer <token>

Response:
{
  "message": "Notes retrieved successfully",
  "notes": [...],
  "count": 5
}
```

#### Get Specific Note

```http
GET /api/v1/notes/{noteId}
Authorization: Bearer <token>
```

#### Update Note

```http
PUT /api/v1/notes/{noteId}
Content-Type: application/json
Authorization: Bearer <token>

{
  "title": "Updated Title",
  "content": "Updated content..."
}
```

#### Delete Note

```http
DELETE /api/v1/notes/{noteId}
Authorization: Bearer <token>
```

### 🤖 AI Chat Management

#### Start/Continue Chat Conversation

```http
POST /api/v1/chat
Content-Type: application/json
Authorization: Bearer <token>

{
  "sessionId": "optional-existing-session-id",
  "message": "Hello, I want to discuss machine learning concepts",
  "model": "openai"  // or "gemini"
}

Response:
{
  "message": "Chat response generated successfully",
  "data": {
    "sessionId": "uuid-session-id",
    "message": "Hello! I'd be happy to discuss machine learning...",
    "role": "assistant",
    "model": "openai",
    "memories": [...],  // Relevant memories retrieved for context
    "createdAt": "2023-12-01T10:00:00Z"
  }
}
```

#### Get All Chat Sessions

```http
GET /api/v1/chat/sessions
Authorization: Bearer <token>

Response:
{
  "message": "Chat sessions retrieved successfully",
  "sessions": [
    {
      "id": "...",
      "sessionId": "uuid-session-id",
      "title": "Hello, I want to discuss machine...",
      "model": "openai",
      "messageCount": 10,
      "lastActivity": "2023-12-01T10:00:00Z",
      "createdAt": "2023-12-01T09:00:00Z"
    }
  ],
  "count": 1
}
```

#### Get Chat History for Session

```http
GET /api/v1/chat/sessions/{sessionId}
Authorization: Bearer <token>

Response:
{
  "message": "Chat history retrieved successfully",
  "sessionId": "uuid-session-id",
  "messages": [
    {
      "id": "...",
      "role": "user",
      "content": "Hello, I want to discuss machine learning",
      "model": "openai",
      "createdAt": "2023-12-01T09:00:00Z"
    },
    {
      "id": "...",
      "role": "assistant",
      "content": "Hello! I'd be happy to discuss machine learning...",
      "model": "openai",
      "createdAt": "2023-12-01T09:00:01Z"
    }
  ],
  "count": 2
}
```

#### Delete Chat Session

```http
DELETE /api/v1/chat/sessions/{sessionId}
Authorization: Bearer <token>

Response:
{
  "message": "Chat session deleted successfully"
}
```

#### 📝 Update Note with AI Based on Chat

```http
POST /api/v1/chat/update-note
Content-Type: application/json
Authorization: Bearer <token>

{
  "noteId": "note-object-id",
  "sessionId": "chat-session-id",
  "model": "openai",  // or "gemini"
  "prompt": "Please focus on the key insights about neural networks"  // Optional custom prompt
}

Response:
{
  "message": "Note updated successfully with AI assistance",
  "note": {
    "id": "note-object-id",
    "title": "Machine Learning Notes",
    "content": "Updated content with insights from the conversation...",
    "updatedAt": "2023-12-01T10:00:00Z"
  }
}
```

## 🧠 Memory Management with Mem0

The system automatically manages conversational memory using Mem0:

### How Memory Works

1. **Automatic Storage**: Every chat message is stored as memory
2. **Context Retrieval**: Relevant memories are retrieved for each new message
3. **Smart Filtering**: Memories are filtered and ranked by relevance
4. **Session Tracking**: Memories are linked to chat sessions for better organization

### Memory Features

- **User-specific**: Each user's memories are isolated
- **Session-aware**: Memories can be grouped by chat sessions
- **Searchable**: Smart semantic search for relevant context
- **Persistent**: Memories persist across sessions for continuity

## Frontend Integration

### React/Next.js Example

```typescript
import { useAuth } from '@clerk/nextjs';

function useApiCall() {
  const { getToken } = useAuth();

  const apiCall = async (endpoint: string, options: RequestInit = {}) => {
    const token = await getToken();

    return fetch(`http://localhost:8080/api/v1${endpoint}`, {
      ...options,
      headers: {
        ...options.headers,
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    });
  };

  return apiCall;
}

// Chat functionality example
function ChatInterface() {
  const apiCall = useApiCall();
  const [sessionId, setSessionId] = useState(null);
  const [messages, setMessages] = useState([]);

  const sendMessage = async (message: string, model: string) => {
    const response = await apiCall('/chat', {
      method: 'POST',
      body: JSON.stringify({
        sessionId,
        message,
        model
      }),
    });

    const data = await response.json();
    if (response.ok) {
      setSessionId(data.data.sessionId);
      setMessages(prev => [...prev,
        { role: 'user', content: message },
        { role: 'assistant', content: data.data.message }
      ]);
    }
    return data;
  };

  const updateNoteWithChat = async (noteId: string) => {
    const response = await apiCall('/chat/update-note', {
      method: 'POST',
      body: JSON.stringify({
        noteId,
        sessionId,
        model: 'openai'
      }),
    });
    return response.json();
  };

  const getChatSessions = async () => {
    const response = await apiCall('/chat/sessions');
    return response.json();
  };

  return (
    // Your React chat component
  );
}

// Notes management example
function NotesManager() {
  const apiCall = useApiCall();

  const createNote = async (noteData: NoteData) => {
    const response = await apiCall('/notes', {
      method: 'POST',
      body: JSON.stringify(noteData),
    });
    return response.json();
  };

  const getMyNotes = async () => {
    const response = await apiCall('/notes');
    return response.json();
  };

  return (
    // Your React component
  );
}
```

## Architecture Benefits

### 1. **Secure Authentication**

- No password storage in your database
- Google OAuth handled by Clerk
- JWT tokens for stateless API authentication

### 2. **Efficient Data Structure**

- Users store only note IDs, not full note content
- MongoDB aggregation for efficient note retrieval
- Automatic relationship management

### 3. **API Key Management**

- Secure storage of AI service API keys
- Per-user API key management
- Easy integration with OpenAI, Gemini

### 4. **Scalable Design**

- Stateless authentication
- MongoDB aggregation pipelines
- Clerk handles user management complexity

### 5. **🤖 AI-Powered Features**

- **Memory-Enhanced Conversations**: Context-aware AI responses
- **Smart Note Updates**: AI analyzes chat history and updates notes intelligently
- **Multi-Model Support**: OpenAI and Gemini integration
- **Fine-Tuned Prompts**: Specialized prompts for note updating tasks

## MongoDB Aggregation

The system uses MongoDB aggregation to efficiently join user and note data:

```javascript
// Example aggregation pipeline for getting user with notes
[
  { $match: { clerkId: "user_123" } },
  {
    $lookup: {
      from: "notes",
      localField: "noteIds",
      foreignField: "_id",
      as: "notes",
    },
  },
];
```

## Security Considerations

1. **API Keys**: Stored encrypted in database
2. **Session Tokens**: Validated on each request
3. **User Isolation**: Notes and chats are strictly tied to authenticated users
4. **Memory Privacy**: Mem0 memories are user-specific and isolated
5. **CORS**: Configure appropriately for your frontend domain

## Error Handling

The API returns consistent error responses:

```json
{
  "error": "User not authenticated",
  "message": "Authorization header is required"
}
```

Common HTTP status codes:

- `200`: Success
- `201`: Created
- `400`: Bad Request (missing API keys, invalid model, etc.)
- `401`: Unauthorized
- `404`: Not Found
- `500`: Internal Server Error

## 🚀 Workflow Examples

### Typical AI Note-Taking Workflow

1. **Start Conversation**:

   ```http
   POST /api/v1/chat
   { "message": "Let's discuss React hooks", "model": "openai" }
   ```

2. **Continue Discussion**:

   ```http
   POST /api/v1/chat
   { "sessionId": "existing-id", "message": "What about useEffect?", "model": "openai" }
   ```

3. **Update Notes with Insights**:

   ```http
   POST /api/v1/chat/update-note
   { "noteId": "my-react-notes", "sessionId": "chat-id", "model": "openai" }
   ```

4. **Review Updated Notes**:

   ```http
   GET /api/v1/notes/my-react-notes
   ```

## Next Steps

1. **Frontend Integration**: Connect your React/Next.js app with Clerk
2. **API Key Encryption**: Implement encryption for stored API keys
3. **Rate Limiting**: Add rate limiting middleware for AI calls
4. **Logging**: Implement comprehensive API logging
5. **Testing**: Add unit and integration tests
6. **Memory Management**: Implement memory cleanup and archiving
7. **Advanced Prompts**: Create more specialized prompts for different note types

## Troubleshooting

### Common Issues

1. **"User not authenticated"**

   - Check if Authorization header is present
   - Verify Clerk session token is valid

2. **"User profile not found"**

   - User needs to create profile first via `/user/profile`
   - Check if ClerkID matches

3. **"No [model] API key found"**

   - User needs to add their OpenAI/Gemini API key via `/user/api-keys`
   - Verify API key is valid

4. **"AI API call failed"**

   - Check API key validity
   - Verify sufficient credits/quota in AI service
   - Check network connectivity to AI services

5. **"Memory service errors"**

   - Verify MEM0_API_KEY is set correctly
   - Check Mem0 service status and credits

6. **Database connection errors**
   - Verify MongoDB connection string
   - Check database permissions

For more details, see the [Clerk Go SDK documentation](https://clerk.com/docs/references/go/overview).
