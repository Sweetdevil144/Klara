# Klara Backend Documentation

## System Overview

Klara is an AI-powered note-taking application backend built with Go, using the Fiber web framework. The system provides intelligent note management with contextual AI assistance, user authentication via Clerk, and memory-enhanced conversations.

## System Architecture

### Core Components

#### **Web Framework & Server**

- **Framework**: GoFiber v2 - High-performance HTTP framework
- **Port**: 8080 (configurable via PORT environment variable)
- **CORS**: Enabled for cross-origin requests
- **Logging**: Request/response logging with latency tracking
- **Error Handling**: Centralized error handling with structured JSON responses

#### **Database Layer**

- **Database**: MongoDB with official Go driver
- **Connection**: Singleton pattern with connection pooling
- **Collections**:
  - `users` - User profiles and API keys
  - `notes` - User notes and content
  - `chat_sessions` - AI conversation sessions
  - `chat_messages` - Individual chat messages

#### **Authentication & Security**

- **Provider**: Clerk Authentication
- **Method**: JWT token verification
- **Middleware**: Bearer token validation with 30-second clock skew tolerance
- **Authorization**: User-scoped data access control
- **Security Features**:
  - Authorized party validation for CSRF protection
  - Automatic token refresh handling
  - Session extension for active users

#### **AI Integration**

- **Providers**: OpenAI and Google Gemini
- **Memory System**: Mem0 AI for contextual conversation memory
- **Features**:
  - Multi-model AI support (GPT-3.5, GPT-4, Gemini 1.5/2.0)
  - Context-aware responses using conversation history
  - Direct, no-fluff AI responses (eliminates conversational padding)
  - Content-preserving note updates (adds to existing content rather than replacing)

#### **Services Architecture**

##### **AIService**

- Handles multi-provider AI communication (OpenAI/Gemini)
- Manages conversation context and memory integration
- Implements content-preserving note enhancement
- Provides direct, actionable AI responses without conversational fluff

##### **MemoryService**

- Integrates with Mem0 AI for persistent conversation memory
- Enables contextual AI responses based on user history
- Supports memory search and retrieval for relevant context

### System Instructions

#### **AI Response Guidelines**

The system implements strict AI response formatting:

- No conversational phrases ("Here's", "Okay", "I understand")
- No meta-commentary about actions being performed
- Direct, actionable content only
- Content preservation in note updates (append, don't replace)
- Factual, professional tone without unnecessary preamble

#### **Note Update Behavior**

- **Content Preservation**: Always maintain existing note content
- **Additive Updates**: Append new information to existing content
- **Format Consistency**: Maintain original formatting (bullets, numbers, headers)
- **List Extension**: When adding items, continue existing numbering/formatting
- **Contradiction Handling**: Only remove content when directly contradicted

#### **Authentication Flow**

1. Client sends Bearer token in Authorization header
2. Middleware extracts and validates JWT with Clerk
3. User ID extracted from verified claims
4. User context attached to request for downstream handlers
5. All data operations scoped to authenticated user

## API Endpoints

### Public Endpoints

#### Health Check

```bash
GET /api/v1/public/health
Response: {"status": "ok", "message": "Klara API is running"}
```

#### Authentication Debug

```bash
GET /api/v1/public/debug/auth
Headers: Authorization: Bearer <token> (optional)
Response: Authentication status and user information
```

### Protected Endpoints (Require Authentication)

#### User Management

##### **Create/Sync User Profile**

```bash
POST /api/v1/user/profile
Headers: Authorization: Bearer <token>
Body: {"email": "user@example.com", "firstName": "John", "lastName": "Doe"}
Response: Created/updated user profile
```

##### **Get User Profile**

```bash
GET /api/v1/user/profile
Headers: Authorization: Bearer <token>
Response: Complete user profile with API key status
```

##### **Update API Keys**

```bash
PUT /api/v1/user/api-keys
Headers: Authorization: Bearer <token>
Body: {"openaiKey": "sk-...", "geminiKey": "AI..."}
Response: API key update confirmation
```

##### **Delete API Key**

```bash
DELETE /api/v1/user/api-keys/{keyType}
Headers: Authorization: Bearer <token>
Parameters: keyType (openai|gemini)
Response: Deletion confirmation
```

##### **Get User with Notes**

```bash
GET /api/v1/user/with-notes
Headers: Authorization: Bearer <token>
Response: User profile with associated notes
```

#### Notes Management

##### **Create Note**

```bash
POST /api/v1/notes
Headers: Authorization: Bearer <token>
Body: {"title": "Note Title", "content": "Note content"}
Response: Created note object
```

##### **Get All User Notes**

```bash
GET /api/v1/notes
Headers: Authorization: Bearer <token>
Response: {"notes": [...], "count": number}
```

##### **Get Specific Note**

```bash
GET /api/v1/notes/{id}
Headers: Authorization: Bearer <token>
Response: Single note object
```

##### **Update Note**

```bash
PUT /api/v1/notes/{id}
Headers: Authorization: Bearer <token>
Body: {"title": "Updated Title", "content": "Updated content"}
Response: Updated note object
```

##### **Delete Note**

```bash
DELETE /api/v1/notes/{id}
Headers: Authorization: Bearer <token>
Response: Deletion confirmation
```

#### AI Chat Integration

##### **Chat with Note Context**

```bash
POST /api/v1/notes/{id}/chat
Headers: Authorization: Bearer <token>
Body: {
  "message": "User question about the note",
  "model": "gpt-4o-mini",
  "provider": "openai"
}
Response: {
  "message": "AI response",
  "model": "openai",
  "noteContext": "Note title and content",
  "suggestion": "Actionable suggestion"
}
```

##### **Apply AI Suggestion to Note**

```bash
POST /api/v1/notes/{id}/apply-suggestion
Headers: Authorization: Bearer <token>
Body: {"newTitle": "Updated title", "newContent": "Enhanced content"}
Response: Updated note with AI enhancements
```

#### General AI Chat

##### **Start/Continue Chat Session**

```bash
POST /api/v1/chat
Headers: Authorization: Bearer <token>
Body: {
  "sessionId": "optional-uuid",
  "message": "User message",
  "model": "openai"
}
Response: {
  "sessionId": "uuid",
  "message": "AI response",
  "role": "assistant",
  "model": "openai",
  "memories": [...],
  "createdAt": "timestamp"
}
```

##### **Get Chat Sessions**

```bash
GET /api/v1/chat/sessions
Headers: Authorization: Bearer <token>
Response: List of user's chat sessions with metadata
```

##### **Get Chat History**

```bash
GET /api/v1/chat/sessions/{sessionId}
Headers: Authorization: Bearer <token>
Response: Complete message history for session
```

##### **Delete Chat Session**

```bash
DELETE /api/v1/chat/sessions/{sessionId}
Headers: Authorization: Bearer <token>
Response: Deletion confirmation
```

##### **Update Note with Chat Context**

```bash
POST /api/v1/chat/update-note
Headers: Authorization: Bearer <token>
Body: {
  "noteId": "note-uuid",
  "sessionId": "chat-session-uuid",
  "model": "openai",
  "prompt": "Custom enhancement instruction"
}
Response: Note updated with AI enhancements based on chat history
```

## Data Models

### User Model

```json
{
  "id": "ObjectID",
  "clerkId": "string (unique)",
  "email": "string",
  "username": "string",
  "firstName": "string",
  "lastName": "string",
  "openaiKey": "string (encrypted)",
  "geminiKey": "string (encrypted)",
  "createdAt": "timestamp",
  "updatedAt": "timestamp",
  "noteIds": ["ObjectID array"]
}
```

### Note Model

```json
{
  "id": "ObjectID",
  "userId": "ObjectID (reference to User)",
  "title": "string",
  "content": "string",
  "createdAt": "timestamp",
  "updatedAt": "timestamp"
}
```

### Chat Session Model

```json
{
  "id": "ObjectID",
  "sessionId": "string (UUID)",
  "userId": "ObjectID",
  "title": "string (first message preview)",
  "model": "string (openai|gemini)",
  "messageCount": "number",
  "lastActivity": "timestamp",
  "createdAt": "timestamp"
}
```

### Chat Message Model

```json
{
  "id": "ObjectID",
  "sessionId": "string",
  "userId": "ObjectID",
  "role": "string (user|assistant)",
  "content": "string",
  "model": "string",
  "memories": ["Memory objects"],
  "createdAt": "timestamp"
}
```

## Error Handling

### Standard Error Response Format

```json
{
  "error": true,
  "message": "Error description"
}
```

### Common HTTP Status Codes

- **200**: Success
- **400**: Bad Request (invalid input)
- **401**: Unauthorized (invalid/missing token)
- **404**: Not Found (resource doesn't exist)
- **500**: Internal Server Error

### Authentication Errors

```json
{
  "error": "Invalid or expired token",
  "debug": "Token verification failed - please refresh your session"
}
```

## Environment Configuration

### Required Environment Variables

- `CLERK_SECRET_KEY`: Clerk authentication secret
- `CLERK_PUBLISHABLE_KEY`: Clerk public key
- `MONGO_URI`: MongoDB connection string
- `MONGO_DB_NAME`: Database name
- `MEM0_API_KEY`: Mem0 AI service API key
- `PORT`: Server port (optional, defaults to 8080)

### Optional Configuration

- User API keys stored per-user for OpenAI and Gemini
- Custom prompts supported for AI interactions
- Configurable AI model selection per request

## Performance Considerations

### Connection Management

- MongoDB connection pooling
- HTTP client reuse for AI API calls
- 60-second timeout for AI requests

### Memory Management

- Contextual memory search limited to 5 most relevant memories
- Chat history pagination for large sessions
- Efficient note filtering and search

### Security Features

- JWT token validation with clock skew tolerance
- User-scoped data access (no cross-user data leakage)
- Encrypted API key storage
- CORS protection with authorized party validation
