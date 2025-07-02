"use client";

import { useAuth, UserButton } from "@clerk/nextjs";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import {
  Menu,
  X,
  Plus,
  Search,
  FileText,
  MessageCircle,
  Trash2,
  Edit,
  Bot,
  User,
  Code,
  Eye,
} from "lucide-react";
import { notesApi, Note, sessionManager } from "@/lib/api";
import AIChatSidebar from "@/components/AIChatSidebar";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";

export default function NotesPage() {
  const { isLoaded, isSignedIn, getToken } = useAuth();
  const router = useRouter();
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [aiChatOpen, setAiChatOpen] = useState(false);
  const [notes, setNotes] = useState<Note[]>([]);
  const [selectedNote, setSelectedNote] = useState<Note | null>(null);
  const [searchQuery, setSearchQuery] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [editingTitle, setEditingTitle] = useState(false);
  const [viewMode, setViewMode] = useState<"code" | "preview">("code");
  const [userToken] = useState<string>("");

  useEffect(() => {
    if (isLoaded && !isSignedIn) {
      router.push("/");
    }
  }, [isLoaded, isSignedIn, router]);

  useEffect(() => {
    if (isLoaded && isSignedIn) {
      loadNotes();

      // Start session extension for active users
      sessionManager.startSessionExtension(getToken);

      // Track user activity on page interactions
      const trackActivity = () => sessionManager.trackUserActivity();

      // Add activity listeners
      document.addEventListener("click", trackActivity);
      document.addEventListener("keydown", trackActivity);
      document.addEventListener("scroll", trackActivity);
      document.addEventListener("mousemove", trackActivity);

      // Cleanup on unmount
      return () => {
        sessionManager.stopSessionExtension();
        document.removeEventListener("click", trackActivity);
        document.removeEventListener("keydown", trackActivity);
        document.removeEventListener("scroll", trackActivity);
        document.removeEventListener("mousemove", trackActivity);
      };
    }
  }, [isLoaded, isSignedIn, getToken]);

  const loadNotes = async () => {
    try {
      setLoading(true);
      setError(null);
      const notes = await notesApi.getNotes(getToken);
      setNotes(notes);

      if (notes.length > 0 && !selectedNote) {
        setSelectedNote(notes[0]);
      }
    } catch (error) {
      console.error("Failed to load notes:", error);
      setError(
        `Failed to load notes: ${
          error instanceof Error ? error.message : "Unknown error"
        }`
      );
      setNotes([]);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateNote = async () => {
    try {
      const newNote = await notesApi.createNote(
        {
          title: "Untitled Note",
          content: "",
        },
        getToken
      );
      setNotes((prev) => [newNote, ...prev]);
      setSelectedNote(newNote);
    } catch (error) {
      console.error("Failed to create note:", error);
      setError("Failed to create note");
    }
  };

  const handleUpdateNote = async (updatedNote: Note) => {
    try {
      const note = await notesApi.updateNote(
        updatedNote.id,
        {
          title: updatedNote.title,
          content: updatedNote.content,
        },
        getToken
      );
      setNotes((prev) => prev.map((n) => (n.id === note.id ? note : n)));
      setSelectedNote(note);
    } catch (error) {
      console.error("Failed to update note:", error);
      setError("Failed to update note");
    }
  };

  const handleDeleteNote = async (noteId: string) => {
    try {
      await notesApi.deleteNote(noteId, getToken);
      setNotes((prev) => prev.filter((n) => n.id !== noteId));
      if (selectedNote?.id === noteId) {
        setSelectedNote(null);
      }
    } catch (error) {
      console.error("Failed to delete note:", error);
      setError("Failed to delete note");
    }
  };

  const filteredNotes = (Array.isArray(notes) ? notes : []).filter(
    (note) =>
      note.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
      note.content.toLowerCase().includes(searchQuery.toLowerCase())
  );

  if (!isLoaded || loading) {
    return (
      <div className="min-h-screen bg-white flex items-center justify-center">
        <div className="animate-pulse text-black text-xl">Loading...</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-white text-black flex">
      {/* Left Sidebar */}
      <div
        className={`
        ${sidebarOpen ? "translate-x-0" : "-translate-x-full"}
        fixed inset-y-0 left-0 z-50 w-80 bg-black border-r border-white transform transition-transform duration-300 ease-in-out
        lg:translate-x-0 lg:static lg:inset-0
      `}
      >
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
              <Search
                className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
                size={16}
              />
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
              onClick={handleCreateNote}
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
                        ? "bg-white text-black border-white"
                        : "bg-black text-white border-gray-700 hover:bg-gray-900"
                    }`}
                    onClick={() => setSelectedNote(note)}
                  >
                    <div className="flex items-start justify-between">
                      <div className="flex-1 min-w-0">
                        <h3 className="font-medium truncate text-sm mb-1">
                          {note.title || "Untitled"}
                        </h3>
                        <p className="text-xs text-gray-400 truncate mb-2">
                          {note.content || "No content"}
                        </p>
                        <p className="text-xs text-gray-500">
                          {new Date(note.updatedAt).toLocaleDateString()}
                        </p>
                      </div>
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          handleDeleteNote(note.id);
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
                  onClick={() => router.push("/profile")}
                >
                  <User size={16} />
                </Button>
              </div>
              <UserButton
                appearance={{
                  elements: {
                    avatarBox: "w-8 h-8 border border-white",
                  },
                }}
              />
            </div>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div
        className={`flex-1 flex flex-col ${
          aiChatOpen ? "mr-96" : ""
        } transition-all duration-300`}
      >
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
              {selectedNote
                ? selectedNote.title || "Untitled"
                : "Select a note"}
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
                      onChange={(e) =>
                        setSelectedNote({
                          ...selectedNote,
                          title: e.target.value,
                        })
                      }
                      onBlur={() => {
                        setEditingTitle(false);
                        handleUpdateNote(selectedNote);
                      }}
                      onKeyDown={(e) => {
                        if (e.key === "Enter") {
                          setEditingTitle(false);
                          handleUpdateNote(selectedNote);
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
                        {selectedNote.title || "Untitled Note"}
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

                {/* View Mode Toggle */}
                <div className="flex items-center justify-between mb-4">
                  <div className="flex items-center space-x-2">
                    <Button
                      onClick={() => setViewMode("code")}
                      variant={viewMode === "code" ? "default" : "ghost"}
                      size="sm"
                      className={`${
                        viewMode === "code"
                          ? "bg-black text-white hover:bg-white hover:text-black"
                          : "text-black hover:bg-black hover:text-white"
                      } border border-black`}
                    >
                      <Code size={16} className="mr-1" />
                      Code
                    </Button>
                    <Button
                      onClick={() => setViewMode("preview")}
                      variant={viewMode === "preview" ? "default" : "ghost"}
                      size="sm"
                      className={`${
                        viewMode === "preview"
                          ? "bg-black text-white hover:bg-white hover:text-black"
                          : "text-black hover:bg-black hover:text-white"
                      } border border-black`}
                    >
                      <Eye size={16} className="mr-1" />
                      Preview
                    </Button>
                  </div>
                </div>

                {/* Content */}
                <div className="min-h-[400px]">
                  {viewMode === "code" ? (
                    <textarea
                      value={selectedNote.content}
                      onChange={(e) =>
                        setSelectedNote({
                          ...selectedNote,
                          content: e.target.value,
                        })
                      }
                      onBlur={() => handleUpdateNote(selectedNote)}
                      className="w-full h-full min-h-[400px] bg-transparent border-none resize-none text-black placeholder-gray-400 focus:outline-none text-base leading-relaxed p-0 font-mono"
                      placeholder="Start writing your thoughts..."
                    />
                  ) : (
                    <div className="w-full min-h-[400px] text-black text-base leading-relaxed prose prose-lg max-w-none">
                      <ReactMarkdown
                        remarkPlugins={[remarkGfm]}
                        components={{
                          h1: ({ children }) => (
                            <h1 className="text-3xl font-bold mb-4 text-black">
                              {children}
                            </h1>
                          ),
                          h2: ({ children }) => (
                            <h2 className="text-2xl font-semibold mb-3 text-black">
                              {children}
                            </h2>
                          ),
                          h3: ({ children }) => (
                            <h3 className="text-xl font-medium mb-2 text-black">
                              {children}
                            </h3>
                          ),
                          p: ({ children }) => (
                            <p className="mb-3 text-black">{children}</p>
                          ),
                          ul: ({ children }) => (
                            <ul className="list-disc ml-6 mb-3 text-black">
                              {children}
                            </ul>
                          ),
                          ol: ({ children }) => (
                            <ol className="list-decimal ml-6 mb-3 text-black">
                              {children}
                            </ol>
                          ),
                          li: ({ children }) => (
                            <li className="mb-1 text-black">{children}</li>
                          ),
                          blockquote: ({ children }) => (
                            <blockquote className="border-l-4 border-gray-300 pl-4 italic text-gray-700 mb-3">
                              {children}
                            </blockquote>
                          ),
                          code: ({ children }) => (
                            <code className="bg-gray-100 px-2 py-1 rounded font-mono text-sm text-black">
                              {children}
                            </code>
                          ),
                          pre: ({ children }) => (
                            <pre className="bg-gray-100 p-4 rounded-lg overflow-x-auto mb-3 font-mono text-sm text-black">
                              {children}
                            </pre>
                          ),
                          a: ({ href, children }) => (
                            <a
                              href={href}
                              className="text-blue-600 hover:text-blue-800 underline"
                              target="_blank"
                              rel="noopener noreferrer"
                            >
                              {children}
                            </a>
                          ),
                          strong: ({ children }) => (
                            <strong className="font-bold text-black">
                              {children}
                            </strong>
                          ),
                          em: ({ children }) => (
                            <em className="italic text-black">{children}</em>
                          ),
                          table: ({ children }) => (
                            <table className="border-collapse border border-gray-300 mb-3 w-full">
                              {children}
                            </table>
                          ),
                          th: ({ children }) => (
                            <th className="border border-gray-300 px-4 py-2 bg-gray-100 font-semibold text-black">
                              {children}
                            </th>
                          ),
                          td: ({ children }) => (
                            <td className="border border-gray-300 px-4 py-2 text-black">
                              {children}
                            </td>
                          ),
                        }}
                      >
                        {selectedNote.content || "*No content to preview*"}
                      </ReactMarkdown>
                    </div>
                  )}
                </div>

                {/* Footer */}
                <div className="mt-6 pt-4 border-t border-black flex items-center justify-between text-sm text-gray-600">
                  <span>
                    Last updated:{" "}
                    {new Date(selectedNote.updatedAt).toLocaleString()}
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
                      onClick={() => handleDeleteNote(selectedNote.id)}
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
                <h3 className="text-xl font-medium text-gray-600 mb-2">
                  No note selected
                </h3>
                <p className="text-gray-500 mb-4">
                  Choose a note from the sidebar or create a new one
                </p>
                <Button
                  onClick={handleCreateNote}
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
        onNoteUpdate={handleUpdateNote}
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
  );
}
