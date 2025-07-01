"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useAuth } from "@clerk/nextjs";
import {
  X,
  Send,
  Bot,
  RefreshCw,
  Check,
  XIcon,
  ChevronDown,
  Crown,
  Zap,
} from "lucide-react";
import {
  notesApi,
  AI_MODELS,
  AIModel,
  NoteChatResponse,
  Note,
} from "@/lib/api";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";

interface AIChatSidebarProps {
  isOpen: boolean;
  onClose: () => void;
  note: Note | null;
  onNoteUpdate: (updatedNote: Note) => void;
  userToken: string;
}

interface ChatMessage {
  id: string;
  role: "user" | "assistant";
  content: string;
  model?: string;
  suggestion?: string;
  timestamp: Date;
}

export default function AIChatSidebar({
  isOpen,
  onClose,
  note,
  onNoteUpdate,
  userToken,
}: AIChatSidebarProps) {
  const { getToken } = useAuth();
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [input, setInput] = useState("");
  const [selectedModel, setSelectedModel] = useState<AIModel>(AI_MODELS[0]);
  const [loading, setLoading] = useState(false);
  const [showModelSelector, setShowModelSelector] = useState(false);
  const [pendingSuggestion, setPendingSuggestion] =
    useState<NoteChatResponse | null>(null);

  const sendMessage = async () => {
    if (!input.trim() || !note || loading) return;

    const userMessage: ChatMessage = {
      id: Date.now().toString(),
      role: "user",
      content: input,
      timestamp: new Date(),
    };

    setMessages((prev) => [...prev, userMessage]);
    setInput("");
    setLoading(true);

    try {
      const response = await notesApi.chatWithNote(
        note.id,
        {
          message: input,
          model: selectedModel.id,
          provider: selectedModel.provider,
        },
        userToken
      );

      const aiMessage: ChatMessage = {
        id: (Date.now() + 1).toString(),
        role: "assistant",
        content: response.message,
        model: selectedModel.name,
        suggestion: response.suggestion,
        timestamp: new Date(),
      };

      setMessages((prev) => [...prev, aiMessage]);
      setPendingSuggestion(response);
    } catch (error) {
      console.error("Failed to send message:", error);
      const errorMessage: ChatMessage = {
        id: (Date.now() + 1).toString(),
        role: "assistant",
        content: "Sorry, I encountered an error. Please try again.",
        timestamp: new Date(),
      };
      setMessages((prev) => [...prev, errorMessage]);
    } finally {
      setLoading(false);
    }
  };

  const applySuggestion = async () => {
    if (!pendingSuggestion || !note) return;

    try {
      setLoading(true);
      
      // Get a fresh token for applying suggestions
      const freshToken = await getToken();
      
      console.log(
        "Applying suggestion with fresh token:",
        freshToken ? "Token obtained" : "No token"
      );
      console.log("Suggestion content:", pendingSuggestion.suggestion);

      if (!freshToken) {
        throw new Error('Unable to obtain authentication token. Please refresh the page and try again.');
      }

      const updatedNote = await notesApi.applySuggestion(
        note.id,
        {
          newContent: pendingSuggestion.suggestion,
        },
        freshToken // Use fresh token instead of cached userToken
      );

      onNoteUpdate(updatedNote);
      setPendingSuggestion(null);

      console.log("Suggestion applied successfully");
    } catch (error) {
      console.error("Failed to apply suggestion:", error);

      // More specific error logging
      if (error instanceof Error) {
        console.error("Error message:", error.message);
        if (
          error.message.includes("Unauthorized") ||
          error.message.includes("401")
        ) {
          console.error(
            "Authentication error - token may be expired or invalid"
          );
        }
      }
    } finally {
      setLoading(false);
    }
  };

  const rejectSuggestion = () => {
    setPendingSuggestion(null);
  };

  const retryLastMessage = async () => {
    if (messages.length === 0) return;

    const lastUserMessage = messages.filter((m) => m.role === "user").pop();
    if (!lastUserMessage) return;

    setLoading(true);

    // Clear pending suggestion when retrying
    setPendingSuggestion(null);

    try {
      // Get a fresh token for retry
      const freshToken = await getToken();
      
      console.log('Retrying with fresh token:', freshToken ? 'Token obtained' : 'No token');
      
      if (!freshToken) {
        throw new Error('Unable to obtain authentication token. Please refresh the page and try again.');
      }

      const response = await notesApi.chatWithNote(
        note!.id,
        {
          message: lastUserMessage.content,
          model: selectedModel.id,
          provider: selectedModel.provider,
        },
        freshToken // Use fresh token instead of cached userToken
      );

      const aiMessage: ChatMessage = {
        id: Date.now().toString(),
        role: "assistant",
        content: response.message,
        model: selectedModel.name,
        suggestion: response.suggestion,
        timestamp: new Date(),
      };

      // Replace the last AI message instead of appending
      setMessages((prev) => {
        const newMessages = [...prev];
        const lastAiIndex = newMessages
          .map((m) => m.role)
          .lastIndexOf("assistant");
        if (lastAiIndex !== -1) {
          newMessages[lastAiIndex] = aiMessage;
        } else {
          newMessages.push(aiMessage);
        }
        return newMessages;
      });

      setPendingSuggestion(response);
    } catch (error) {
      console.error("Failed to retry message:", error);
      
      // Add specific error message for token issues
      if (error instanceof Error && (
        error.message.includes('Unauthorized') || 
        error.message.includes('401') ||
        error.message.includes('token')
      )) {
        const tokenErrorMessage: ChatMessage = {
          id: Date.now().toString(),
          role: "assistant",
          content: "⚠️ Authentication error. Please refresh the page and try again.",
          timestamp: new Date(),
        };
        setMessages((prev) => [...prev, tokenErrorMessage]);
      }
    } finally {
      setLoading(false);
    }
  };

  if (!isOpen) return null;

  return (
    <>
      {/* Overlay */}
      <div
        className="fixed inset-0 bg-black/30 z-40 lg:hidden"
        onClick={onClose}
      />

      {/* Sidebar */}
      <div className="fixed right-0 top-0 h-full w-96 bg-black border-l border-white z-50 flex flex-col">
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b border-white">
          <div className="flex items-center space-x-2">
            <Bot className="w-5 h-5 text-white" />
            <h2 className="text-lg font-semibold text-white">AI Assistant</h2>
          </div>
          <button
            onClick={onClose}
            className="p-2 text-white hover:bg-white hover:text-black rounded-lg border border-white transition-colors"
          >
            <X size={20} />
          </button>
        </div>

        {/* Model Selector */}
        <div className="p-4 border-b border-white">
          <div className="relative">
            <button
              onClick={() => setShowModelSelector(!showModelSelector)}
              className="w-full flex items-center justify-between p-3 bg-white text-black rounded-lg hover:bg-gray-100 transition-colors border border-black"
            >
              <div className="flex items-center space-x-2">
                {selectedModel.free ? (
                  <Zap className="w-4 h-4 text-black" />
                ) : (
                  <Crown className="w-4 h-4 text-black" />
                )}
                <span className="text-sm font-medium">
                  {selectedModel.name}
                </span>
                <span className="text-xs text-gray-600">
                  ({selectedModel.provider})
                </span>
              </div>
              <ChevronDown size={16} />
            </button>

            {showModelSelector && (
              <div className="absolute top-full left-0 right-0 mt-1 bg-white rounded-lg shadow-lg z-10 max-h-64 overflow-y-auto border border-black">
                {AI_MODELS.map((model) => (
                  <button
                    key={model.id}
                    onClick={() => {
                      setSelectedModel(model);
                      setShowModelSelector(false);
                    }}
                    className={`w-full flex items-center space-x-2 p-3 text-left hover:bg-gray-100 transition-colors ${
                      selectedModel.id === model.id ? "bg-gray-100" : ""
                    }`}
                  >
                    {model.free ? (
                      <Zap className="w-4 h-4 text-black" />
                    ) : (
                      <Crown className="w-4 h-4 text-black" />
                    )}
                    <div className="flex-1">
                      <div className="text-sm font-medium text-black">
                        {model.name}
                      </div>
                      <div className="text-xs text-gray-600">
                        {model.provider} • {model.free ? "Free" : "Paid"}
                      </div>
                    </div>
                  </button>
                ))}
              </div>
            )}
          </div>
        </div>

        {/* Note Context */}
        {note && (
          <div className="p-4 border-b border-white">
            <h3 className="text-sm font-medium text-white mb-2">
              Current Note
            </h3>
            <div className="text-xs text-gray-300">
              <div className="font-medium truncate">
                {note.title || "Untitled"}
              </div>
              <div className="mt-1 line-clamp-2">
                {note.content || "No content"}
              </div>
            </div>
          </div>
        )}

        {/* Suggestion Panel */}
        {pendingSuggestion && pendingSuggestion.suggestion && (
          <div className="p-4 border-b border-white bg-white text-black">
            <div className="flex items-center justify-between mb-3">
              <h3 className="text-sm font-medium text-black">
                ✨ AI Suggestion
              </h3>
              <div className="flex space-x-2">
                <button
                  onClick={applySuggestion}
                  disabled={loading}
                  className="flex items-center px-3 py-1 text-xs font-medium text-white bg-green-600 hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed rounded transition-colors border"
                  title="Apply suggestion"
                >
                  {loading ? (
                    <>
                      <div className="w-3 h-3 border border-white border-t-transparent rounded-full animate-spin mr-1" />
                      Applying...
                    </>
                  ) : (
                    <>
                      <Check size={14} className="mr-1" />
                      Accept
                    </>
                  )}
                </button>
                <button
                  onClick={rejectSuggestion}
                  className="flex items-center px-3 py-1 text-xs font-medium text-white bg-red-600 hover:bg-red-700 rounded transition-colors border"
                  title="Reject suggestion"
                >
                  <XIcon size={14} className="mr-1" />
                  Reject
                </button>
              </div>
            </div>
            <div className="text-xs text-black bg-gray-50 p-3 rounded border border-gray-200 max-h-40 overflow-y-auto">
              <div className="font-medium text-gray-700 mb-2">
                Suggested content:
              </div>
              <div className="prose prose-xs max-w-none text-black">
                <ReactMarkdown remarkPlugins={[remarkGfm]}>
                  {pendingSuggestion.suggestion}
                </ReactMarkdown>
              </div>
            </div>
          </div>
        )}

        {/* Messages */}
        <div className="flex-1 overflow-y-auto p-4 space-y-4">
          {messages.length === 0 && (
            <div className="text-center text-gray-400 mt-8">
              <Bot className="w-12 h-12 mx-auto mb-4 opacity-50" />
              <p className="text-sm">Start a conversation about your note!</p>
            </div>
          )}

          {messages.map((message) => (
            <div
              key={message.id}
              className={`flex ${
                message.role === "user" ? "justify-end" : "justify-start"
              }`}
            >
              <div
                className={`max-w-[80%] rounded-lg p-3 border ${
                  message.role === "user"
                    ? "bg-white text-black border-white"
                    : "bg-gray-900 text-white border-gray-700"
                }`}
              >
                <div className="text-sm">
                  <ReactMarkdown remarkPlugins={[remarkGfm]}>
                    {message.content}
                  </ReactMarkdown>
                </div>
                {message.model && (
                  <div className="text-xs text-gray-400 mt-1">
                    {message.model}
                  </div>
                )}
              </div>
            </div>
          ))}

          {loading && (
            <div className="flex justify-start">
              <div className="bg-gray-900 text-white rounded-lg p-3 border border-gray-700">
                <div className="flex items-center space-x-2">
                  <div className="w-2 h-2 bg-white rounded-full animate-pulse"></div>
                  <div
                    className="w-2 h-2 bg-white rounded-full animate-pulse"
                    style={{ animationDelay: "0.2s" }}
                  ></div>
                  <div
                    className="w-2 h-2 bg-white rounded-full animate-pulse"
                    style={{ animationDelay: "0.4s" }}
                  ></div>
                </div>
              </div>
            </div>
          )}
        </div>

        {/* Input */}
        <div className="p-4 border-t border-white">
          <div className="flex space-x-2">
            <Input
              value={input}
              onChange={(e) => setInput(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && sendMessage()}
              placeholder="Ask AI about your note..."
              className="flex-1 border-black text-white placeholder-gray-500"
              disabled={loading || !note}
            />
            <Button
              onClick={sendMessage}
              disabled={loading || !input.trim() || !note}
              className="bg-white text-black hover:bg-gray-100 border border-black"
            >
              <Send size={16} />
            </Button>
          </div>

          {messages.length > 0 && (
            <div className="flex justify-center mt-2">
              <Button
                onClick={retryLastMessage}
                disabled={loading}
                variant="ghost"
                size="sm"
                className="text-white border border-white hover:border-white hover:shadow-[0_0_10px_rgba(255,255,255,0.5)]"
              >
                <RefreshCw size={14} className="mr-1" />
                Retry
              </Button>
            </div>
          )}
        </div>
      </div>
    </>
  );
}
