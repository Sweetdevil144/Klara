# Klara Frontend Documentation

## System Overview

Klara's frontend is a modern, responsive web application built with Next.js 14, TypeScript, and Tailwind CSS. It provides an intuitive interface for AI-powered note-taking with real-time chat capabilities, seamless authentication, and a sophisticated design system.

## Frontend Architecture

### Core Technologies

#### **Framework & Runtime**

- **Framework**: Next.js 14 with App Router
- **Language**: TypeScript for type safety
- **Styling**: Tailwind CSS with custom design system
- **UI Components**: Shadcn/ui component library
- **Icons**: Lucide React icons

#### **Authentication & State Management**

- **Authentication**: Clerk for user management
- **Session Management**: Automatic token refresh and session extension
- **State**: React hooks for local state management
- **Activity Tracking**: User interaction monitoring for session extension

#### **API Integration**

- **HTTP Client**: Native Fetch API with custom wrapper
- **Error Handling**: Automatic retry logic for token expiration
- **Token Management**: Automatic refresh and validation
- **Request Optimization**: Parallel API calls and caching

### Application Structure

#### **App Router Architecture**

##### **Root Layout** (`app/layout.tsx`)

- Global Clerk provider configuration
- Dark theme enforcement
- Global CSS and font loading
- Metadata configuration

##### **Page Structure**

```bash
app/
├── layout.tsx          # Root layout with Clerk provider
├── page.tsx            # Landing page with authentication
├── notes/
│   └── page.tsx        # Main notes interface
├── profile/
│   └── page.tsx        # User profile and settings
├── sign-in/
│   └── [[...sign-in]]/
│       └── page.tsx    # Clerk sign-in page
└── sign-up/
    └── [[...sign-up]]/
        └── page.tsx    # Clerk sign-up page
```

#### **Component Architecture**

##### **Core Components**

- `AIChatSidebar` - AI chat interface with model selection
- `UI Components` - Reusable Shadcn/ui components
- `Layout Components` - Navigation and structural elements

##### **Component Hierarchy**

```bash
App
├── RootLayout (Clerk Provider)
├── LandingPage (Unauthenticated)
├── NotesPage (Main Application)
│   ├── Sidebar (Notes List)
│   ├── NoteEditor (Content Editor)
│   └── AIChatSidebar (AI Assistant)
└── ProfilePage (User Settings)
```

### Design System

#### **Color Scheme**

- **Primary**: Black and White contrast theme
- **Background**: Pure black (#000000)
- **Foreground**: Pure white (#FFFFFF)
- **Accents**: Gray scale variations
- **Interactive**: Hover state inversions (black ↔ white)

#### **Typography**

- **Font**: System fonts with fallbacks
- **Hierarchy**: Clear heading and body text distinction
- **Readability**: High contrast for accessibility

#### **Layout Patterns**

- **Responsive Design**: Mobile-first approach
- **Grid System**: Tailwind CSS grid and flexbox
- **Spacing**: Consistent spacing scale
- **Breakpoints**: Standard responsive breakpoints

### Key Features & Components

#### **Landing Page** (`app/page.tsx`)

##### **Features**

- Hero section with gradient animations
- Feature showcase cards with hover effects
- Dual authentication options (Sign In/Sign Up)
- Responsive design with mobile optimization
- Loading states and smooth transitions

##### **Visual Elements**

- Matrix-style background effects
- Gradient text animations
- Interactive hover states
- Floating UI elements for visual appeal

#### **Notes Interface** (`app/notes/page.tsx`)

##### **Layout Structure**

- **Left Sidebar**: Notes list with search and filtering
- **Main Content**: Note editor with live preview
- **Right Sidebar**: AI chat interface (toggleable)

##### Features

- Real-time note editing with auto-save
- Search functionality across all notes
- Note creation, update, and deletion
- Responsive sidebar with mobile hamburger menu
- Activity tracking for session management

##### **State Management**

```typescript
// Core state structure
const [notes, setNotes] = useState<Note[]>([]);
const [selectedNote, setSelectedNote] = useState<Note | null>(null);
const [searchQuery, setSearchQuery] = useState("");
const [loading, setLoading] = useState(true);
const [error, setError] = useState<string | null>(null);
```

#### **AI Chat Sidebar** (`components/AIChatSidebar.tsx`)

##### Features : AI Sidebar

- Multi-model AI support (OpenAI, Gemini)
- Model selection with pricing indicators
- Real-time chat interface
- Markdown rendering for AI responses
- Note suggestion application
- Conversation context awareness

##### **AI Models Supported**

```typescript
const AI_MODELS: AIModel[] = [
  // OpenAI Models
  {
    id: "gpt-3.5-turbo",
    name: "GPT-3.5 Turbo",
    provider: "openai",
    free: true,
  },
  { id: "gpt-4", name: "GPT-4", provider: "openai", free: false },
  { id: "gpt-4o", name: "GPT-4o", provider: "openai", free: false },

  // Gemini Models
  {
    id: "gemini-1.5-flash",
    name: "Gemini 1.5 Flash",
    provider: "gemini",
    free: true,
  },
  {
    id: "gemini-2.0-flash",
    name: "Gemini 2.0 Flash",
    provider: "gemini",
    free: true,
  },
];
```

#### **Profile Management** (`app/profile/page.tsx`)

##### Features : Profile Management

- User profile information display
- API key management (OpenAI, Gemini)
- Secure key storage with visibility toggles
- Profile synchronization with Clerk
- Success/error feedback systems

##### **Security Features**

- API key masking in UI
- Secure key transmission
- Individual key deletion options
- Profile data validation

### API Integration Layer (`lib/api.ts`)

#### **Core Features**

##### **Automatic Token Management**

- Token refresh on expiration
- Retry logic for failed requests
- Session extension for active users
- Activity tracking integration

##### **Session Management**

```typescript
// Session extension utilities
const sessionManager = {
  startSessionExtension, // Begin automatic session renewal
  stopSessionExtension, // Stop session renewal
  trackUserActivity, // Record user interactions
};
```

##### **API Request Wrapper**

```typescript
async function apiRequest<T>(
  endpoint: string,
  options: RequestInit = {},
  getTokenFn?: () => Promise<string | null>
): Promise<T>;
```

#### **API Modules**

##### **Notes API** (`notesApi`)

- `getNotes()` - Retrieve all user notes
- `getNote(id)` - Get specific note
- `createNote(data)` - Create new note
- `updateNote(id, data)` - Update existing note
- `deleteNote(id)` - Delete note
- `chatWithNote(id, data)` - AI chat about note
- `applySuggestion(id, data)` - Apply AI suggestions

##### **User API** (`userApi`)

- `createProfile()` - Create/sync user profile
- `getProfile()` - Retrieve user profile
- `updateAPIKeys(data)` - Update AI API keys
- `deleteAPIKey(keyType)` - Remove specific API key

### Authentication & Security

#### **Clerk Integration**

##### **Authentication Flow**

1. User accesses protected route
2. Middleware checks authentication status
3. Redirects to sign-in if unauthenticated
4. Provides user context and token access

##### **Protected Routes**

- `/notes` - Main application interface
- `/profile` - User settings and API keys

##### **Token Management**

- Automatic token retrieval via `getToken()`
- Token refresh on expiration
- Session extension for active users
- Secure API key storage

#### **Session Extension**

##### **Activity Tracking**

- Mouse movement detection
- Keyboard input monitoring
- Click event tracking
- Scroll event monitoring

##### **Extension Logic**

- Refresh token every 5 minutes for active users
- Activity window: 10 minutes
- Automatic cleanup on user logout
- Background session maintenance

### Responsive Design

#### **Breakpoint Strategy**

- **Mobile**: < 768px (hamburger menu, stacked layout)
- **Tablet**: 768px - 1024px (collapsible sidebars)
- **Desktop**: > 1024px (full three-panel layout)

#### **Mobile Optimizations**

- Touch-friendly interface elements
- Swipe gestures for navigation
- Optimized keyboard interactions
- Responsive typography scaling

#### **Layout Adaptations**

- Collapsible sidebars on smaller screens
- Overlay modals for mobile interactions
- Adaptive spacing and sizing
- Touch-optimized button sizes

### Performance Optimizations

#### **Code Splitting**

- Page-level code splitting with Next.js
- Component lazy loading where appropriate
- Dynamic imports for heavy components

#### **API Optimizations**

- Request deduplication
- Automatic retry with exponential backoff
- Parallel API calls for data fetching
- Efficient state updates

#### **User Experience**

- Loading states for all async operations
- Optimistic updates for immediate feedback
- Error boundaries for graceful error handling
- Smooth transitions and animations

### Error Handling

#### **Error Boundaries**

- Component-level error catching
- Graceful degradation
- User-friendly error messages
- Automatic error reporting

#### **API Error Handling**

- Automatic token refresh on 401 errors
- Retry logic for network failures
- User feedback for API errors
- Fallback states for failed requests

#### **User Feedback**

- Toast notifications for actions
- Loading indicators for async operations
- Error messages with actionable guidance
- Success confirmations for completed actions

### Accessibility

#### **WCAG Compliance**

- High contrast color scheme
- Keyboard navigation support
- Screen reader compatibility
- Focus management

#### **Interactive Elements**

- Clear focus indicators
- Accessible button labels
- Semantic HTML structure
- ARIA attributes where needed

### Development Features

#### **Type Safety**

- TypeScript throughout the application
- Strict type checking enabled
- Interface definitions for all API responses
- Generic type utilities for reusability

#### **Code Quality**

- ESLint configuration
- Prettier code formatting
- Component composition patterns
- Custom hooks for reusable logic

This architecture provides a scalable, maintainable, and user-friendly frontend that seamlessly integrates with the Klara backend to deliver an exceptional AI-powered note-taking experience.
