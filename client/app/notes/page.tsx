'use client'

import { useAuth, UserButton } from '@clerk/nextjs'
import { useRouter } from 'next/navigation'
import { useEffect, useState } from 'react'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Separator } from '@/components/ui/separator'
import { Menu, X, Plus, Search, FileText, Settings, MessageCircle, Trash2, Edit, Bot, User } from 'lucide-react'
import { notesApi, Note } from '@/lib/api'
import AIChatSidebar from '@/components/AIChatSidebar'

export default function NotesPage() {
  const { isLoaded, isSignedIn, getToken } = useAuth()
  const router = useRouter()
  const [sidebarOpen, setSidebarOpen] = useState(false)
  const [aiChatOpen, setAiChatOpen] = useState(false)
  const [notes, setNotes] = useState<Note[]>([])
  const [selectedNote, setSelectedNote] = useState<Note | null>(null)
  const [searchQuery, setSearchQuery] = useState('')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [editingTitle, setEditingTitle] = useState(false)
  const [editingContent, setEditingContent] = useState(false)
  const [userToken, setUserToken] = useState<string>('')

  useEffect(() => {
    if (isLoaded && !isSignedIn) {
      router.push('/')
    }
  }, [isLoaded, isSignedIn, router])

  useEffect(() => {
    if (isLoaded && isSignedIn) {
      loadNotes()
    }
  }, [isLoaded, isSignedIn])

  const loadNotes = async () => {
    try {
      setLoading(true)
      setError(null)
      const token = await getToken()
      if (!token) throw new Error('No auth token')
      
      setUserToken(token)
      console.log('Fetching notes with token:', token ? 'Token present' : 'No token')
      const fetchedNotes = await notesApi.getNotes(token)
      console.log('Fetched notes:', fetchedNotes)
      const notesArray = Array.isArray(fetchedNotes) ? fetchedNotes : []
      setNotes(notesArray)
      
      if (notesArray.length > 0 && !selectedNote) {
        setSelectedNote(notesArray[0])
      }
    } catch (err) {
      console.error('Failed to load notes:', err)
      setError(`Failed to load notes: ${err instanceof Error ? err.message : 'Unknown error'}`)
      setNotes([])
    } finally {
      setLoading(false)
    }
  }

  const createNote = async () => {
    try {
      const token = await getToken()
      if (!token) throw new Error('No auth token')
      
      const newNote = await notesApi.createNote({
        title: 'Untitled Note',
        content: ''
      }, token)
      
      setNotes(prev => [newNote, ...prev])
      setSelectedNote(newNote)
      setEditingTitle(true)
    } catch (err) {
      console.error('Failed to create note:', err)
      setError('Failed to create note')
    }
  }

  const updateNote = async (id: string, updates: { title?: string; content?: string }) => {
    try {
      const token = await getToken()
      if (!token) throw new Error('No auth token')
      
      const updatedNote = await notesApi.updateNote(id, updates, token)
      
      setNotes(prev => prev.map(note => 
        note.id === id ? updatedNote : note
      ))
      
      if (selectedNote?.id === id) {
        setSelectedNote(updatedNote)
      }
    } catch (err) {
      console.error('Failed to update note:', err)
      setError('Failed to update note')
    }
  }

  const deleteNote = async (id: string) => {
    try {
      const token = await getToken()
      if (!token) throw new Error('No auth token')
      
      await notesApi.deleteNote(id, token)
      
      setNotes(prev => prev.filter(note => note.id !== id))
      
      if (selectedNote?.id === id) {
        const remainingNotes = notes.filter(note => note.id !== id)
        setSelectedNote(remainingNotes.length > 0 ? remainingNotes[0] : null)
      }
    } catch (err) {
      console.error('Failed to delete note:', err)
      setError('Failed to delete note')
    }
  }

  const handleNoteUpdate = (updatedNote: Note) => {
    setNotes(prev => prev.map(note => 
      note.id === updatedNote.id ? updatedNote : note
    ))
    setSelectedNote(updatedNote)
  }

  const filteredNotes = (Array.isArray(notes) ? notes : []).filter(note =>
    note.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
    note.content.toLowerCase().includes(searchQuery.toLowerCase())
  )

  if (!isLoaded || loading) {
    return (
      <div className="min-h-screen bg-white flex items-center justify-center">
        <div className="animate-pulse text-black text-xl">Loading...</div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-white text-black flex">
      {/* Left Sidebar */}
      <div className={`
        ${sidebarOpen ? 'translate-x-0' : '-translate-x-full'}
        fixed inset-y-0 left-0 z-50 w-80 bg-black border-r border-white transform transition-transform duration-300 ease-in-out
        lg:translate-x-0 lg:static lg:inset-0
      `}>
        <div className="flex flex-col h-full">
          {/* Sidebar Header */}
          <div className="flex items-center justify-between p-4 border-b border-white">
            <h1 className="text-xl font-bold text-white">Klara</h1>
            <button
              onClick={() => setSidebarOpen(false)}
              className="lg:hidden p-2 text-white hover:bg-white hover:text-black rounded-lg border border-white transition-colors"
            >
              <X size={20} />
            </button>
          </div>

          {/* Search */}
          <div className="p-4 border-b border-white">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" size={16} />
              <Input
                placeholder="Search notes..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-10 bg-white border-white text-black placeholder-gray-500 focus:border-gray-300"
              />
            </div>
          </div>

          {/* New Note Button */}
          <div className="p-4 border-b border-white">
            <Button
              onClick={createNote}
              className="w-full bg-white text-black hover:bg-gray-100 font-medium border border-white"
            >
              <Plus size={16} className="mr-2" />
              New Note
            </Button>
          </div>

          {/* Notes List */}
          <div className="flex-1 overflow-y-auto">
            {error && (
              <div className="p-4 text-white text-sm bg-black border-b border-white">
                {error}
              </div>
            )}
            
            <div className="p-2">
              {filteredNotes.length === 0 ? (
                <div className="text-center py-8 text-gray-400">
                  <FileText size={48} className="mx-auto mb-4 opacity-50" />
                  <p>No notes found</p>
                </div>
              ) : (
                filteredNotes.map((note) => (
                  <Card
                    key={note.id}
                    className={`mb-2 p-3 cursor-pointer transition-all duration-200 border ${
                      selectedNote?.id === note.id
                        ? 'bg-white text-black border-white'
                        : 'bg-black text-white border-gray-700 hover:bg-gray-900'
                    }`}
                    onClick={() => setSelectedNote(note)}
                  >
                    <div className="flex items-start justify-between">
                      <div className="flex-1 min-w-0">
                        <h3 className="font-medium truncate text-sm mb-1">
                          {note.title || 'Untitled'}
                        </h3>
                        <p className="text-xs text-gray-400 truncate mb-2">
                          {note.content || 'No content'}
                        </p>
                        <p className="text-xs text-gray-500">
                          {new Date(note.updatedAt).toLocaleDateString()}
                        </p>
                      </div>
                      <button
                        onClick={(e) => {
                          e.stopPropagation()
                          deleteNote(note.id)
                        }}
                        className="p-1 text-gray-500 hover:text-white hover:bg-gray-800 rounded transition-colors"
                      >
                        <Trash2 size={14} />
                      </button>
                    </div>
                  </Card>
                ))
              )}
            </div>
          </div>

          {/* Sidebar Footer */}
          <div className="p-4 border-t border-white">
            <div className="flex items-center justify-between">
              <div className="flex space-x-2">
                <Button 
                  variant="ghost" 
                  size="sm" 
                  className="text-white hover:bg-white hover:text-black border border-white"
                  onClick={() => setAiChatOpen(true)}
                >
                  <MessageCircle size={16} />
                </Button>
                <Button 
                  variant="ghost" 
                  size="sm" 
                  className="text-white hover:bg-white hover:text-black border border-white"
                  onClick={() => router.push('/profile')}
                >
                  <User size={16} />
                </Button>
              </div>
              <UserButton 
                appearance={{
                  elements: {
                    avatarBox: "w-8 h-8 border border-white"
                  }
                }}
              />
            </div>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className={`flex-1 flex flex-col ${aiChatOpen ? 'mr-96' : ''} transition-all duration-300`}>
        {/* Top Bar */}
        <div className="h-14 bg-white border-b border-black flex items-center px-4">
          <button
            onClick={() => setSidebarOpen(true)}
            className="lg:hidden p-2 text-black hover:bg-black hover:text-white rounded-lg border border-black transition-colors mr-4"
          >
            <Menu size={20} />
          </button>
          
          <div className="flex-1 flex items-center justify-between">
            <h2 className="text-lg font-medium text-black">
              {selectedNote ? (selectedNote.title || 'Untitled') : 'Select a note'}
            </h2>
            
            {selectedNote && (
              <Button
                onClick={() => setAiChatOpen(true)}
                variant="ghost"
                size="sm"
                className="text-black hover:bg-black hover:text-white border border-black"
              >
                <Bot size={16} className="mr-2" />
                AI Chat
              </Button>
            )}
          </div>
        </div>

        {/* Note Content */}
        <div className="flex-1 overflow-hidden">
          {selectedNote ? (
            <div className="h-full p-6 overflow-y-auto bg-white">
              <div className="max-w-4xl mx-auto">
                {/* Title */}
                <div className="mb-6">
                  {editingTitle ? (
                    <Input
                      value={selectedNote.title}
                      onChange={(e) => setSelectedNote({...selectedNote, title: e.target.value})}
                      onBlur={() => {
                        setEditingTitle(false)
                        updateNote(selectedNote.id, { title: selectedNote.title })
                      }}
                      onKeyDown={(e) => {
                        if (e.key === 'Enter') {
                          setEditingTitle(false)
                          updateNote(selectedNote.id, { title: selectedNote.title })
                        }
                      }}
                      className="text-3xl font-bold bg-transparent border-none p-0 text-black placeholder-gray-400 focus:ring-0 focus:border-none"
                      placeholder="Note title..."
                      autoFocus
                    />
                  ) : (
                    <div className="flex items-center group">
                      <h1 
                        className="text-3xl font-bold text-black cursor-pointer flex-1"
                        onClick={() => setEditingTitle(true)}
                      >
                        {selectedNote.title || 'Untitled Note'}
                      </h1>
                      <button
                        onClick={() => setEditingTitle(true)}
                        className="opacity-0 group-hover:opacity-100 p-2 text-gray-600 hover:text-black transition-all"
                      >
                        <Edit size={16} />
                      </button>
                    </div>
                  )}
                </div>

                <Separator className="bg-black mb-6" />

                {/* Content */}
                <div className="min-h-[400px]">
                  <textarea
                    value={selectedNote.content}
                    onChange={(e) => setSelectedNote({...selectedNote, content: e.target.value})}
                    onBlur={() => updateNote(selectedNote.id, { content: selectedNote.content })}
                    className="w-full h-full min-h-[400px] bg-transparent border-none resize-none text-black placeholder-gray-400 focus:outline-none text-base leading-relaxed p-0"
                    placeholder="Start writing your thoughts..."
                  />
                </div>

                {/* Footer */}
                <div className="mt-6 pt-4 border-t border-black flex items-center justify-between text-sm text-gray-600">
                  <span>
                    Last updated: {new Date(selectedNote.updatedAt).toLocaleString()}
                  </span>
                  <div className="flex space-x-4">
                    <Button 
                      variant="ghost" 
                      size="sm" 
                      className="text-black hover:bg-black hover:text-white border border-black"
                      onClick={() => setAiChatOpen(true)}
                    >
                      <Bot size={16} className="mr-1" />
                      AI Chat
                    </Button>
                    <Button 
                      variant="ghost" 
                      size="sm" 
                      className="text-black hover:bg-black hover:text-white border border-black"
                      onClick={() => deleteNote(selectedNote.id)}
                    >
                      Delete
                    </Button>
                  </div>
                </div>
              </div>
            </div>
          ) : (
            <div className="flex items-center justify-center h-full bg-white">
              <div className="text-center">
                <FileText size={64} className="mx-auto mb-4 text-gray-400" />
                <h3 className="text-xl font-medium text-gray-600 mb-2">No note selected</h3>
                <p className="text-gray-500 mb-4">Choose a note from the sidebar or create a new one</p>
                <Button
                  onClick={createNote}
                  className="bg-black text-white hover:bg-gray-800 font-medium border border-black"
                >
                  <Plus size={16} className="mr-2" />
                  Create Note
                </Button>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* AI Chat Sidebar */}
      <AIChatSidebar
        isOpen={aiChatOpen}
        onClose={() => setAiChatOpen(false)}
        note={selectedNote}
        onNoteUpdate={handleNoteUpdate}
        userToken={userToken}
      />

      {/* Mobile Overlay for Left Sidebar */}
      {sidebarOpen && (
        <div 
          className="fixed inset-0 bg-black/50 z-40 lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}
    </div>
  )
} 