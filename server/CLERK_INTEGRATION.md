# Clerk Authentication Integration

This document explains how to use the Clerk authentication system integrated into your AI Notes App.

## Overview

The app now uses **Clerk** for authentication instead of traditional username/password authentication. Key features:

- 🔐 **Google OAuth** integration via Clerk
- 🔑 **JWT Session Tokens** for API authentication
- 👤 **User Profile Management** with API key storage
- 📝 **Secure Note Management** tied to authenticated users
- 🔄 **MongoDB Aggregation** for efficient data retrieval

## Environment Variables

Ensure your `.env` file contains:

```env
CLERK_SECRET_KEY=sk_live_your_secret_key_here
CLERK_PUBLISHABLE_KEY=pk_live_your_publishable_key_here
JWT_SECRET=your_jwt_secret_here  # Optional, for additional JWT operations
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
    ClaudeKey    string               `json:"claudeKey"`    // Encrypted API key
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
    "hasClaudeKey": true,
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
  "geminiKey": "AI...",
  "claudeKey": "claude-..."
}
```

#### Delete API Key

```http
DELETE /api/v1/user/api-keys/{keyType}
Authorization: Bearer <token>

# keyType can be: openai, gemini, claude
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

// Usage
function UserProfile() {
  const apiCall = useApiCall();

  const createProfile = async (userData: UserData) => {
    const response = await apiCall('/user/profile', {
      method: 'POST',
      body: JSON.stringify(userData),
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
- Easy integration with OpenAI, Gemini, Claude

### 4. **Scalable Design**

- Stateless authentication
- MongoDB aggregation pipelines
- Clerk handles user management complexity

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
3. **User Isolation**: Notes are strictly tied to authenticated users
4. **CORS**: Configure appropriately for your frontend domain

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
- `400`: Bad Request
- `401`: Unauthorized
- `404`: Not Found
- `500`: Internal Server Error

## Next Steps

1. **Frontend Integration**: Connect your React/Next.js app with Clerk
2. **API Key Encryption**: Implement encryption for stored API keys
3. **Rate Limiting**: Add rate limiting middleware
4. **Logging**: Implement comprehensive API logging
5. **Testing**: Add unit and integration tests

## Troubleshooting

### Common Issues

1. **"User not authenticated"**

   - Check if Authorization header is present
   - Verify Clerk session token is valid

2. **"User profile not found"**

   - User needs to create profile first via `/user/profile`
   - Check if ClerkID matches

3. **Database connection errors**
   - Verify MongoDB connection string
   - Check database permissions

For more details, see the [Clerk Go SDK documentation](https://clerk.com/docs/references/go/overview).
