/* eslint-disable @typescript-eslint/no-explicit-any */
'use client'

import { useAuth, useUser, UserButton } from '@clerk/nextjs'
import { useRouter } from 'next/navigation'
import { useEffect, useState } from 'react'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Separator } from '@/components/ui/separator'
import { ArrowLeft, Key, Eye, EyeOff, Save, User, Mail } from 'lucide-react'
import { userApi } from '@/lib/api'

export default function ProfilePage() {
  const { isLoaded, isSignedIn, getToken } = useAuth()
  const { user } = useUser()
  const router = useRouter()
  const [profile, setProfile] = useState<any>(null)
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState<string | null>(null)
  
  // API Keys state
  const [openaiKey, setOpenaiKey] = useState("")
  const [geminiKey, setGeminiKey] = useState("")
  const [showOpenaiKey, setShowOpenaiKey] = useState(false)
  const [showGeminiKey, setShowGeminiKey] = useState(false)

  useEffect(() => {
    if (isLoaded && !isSignedIn) {
      router.push('/')
    }
  }, [isLoaded, isSignedIn, router])

  useEffect(() => {
    if (isLoaded && isSignedIn) {
      loadProfile()
    }
  }, [isLoaded, isSignedIn])

  const loadProfile = async () => {
    try {
      setLoading(true)
      setError(null)
      const profileData = await userApi.getProfile(getToken)
      setProfile(profileData)
      setOpenaiKey(profileData.openaiKey || "")
      setGeminiKey(profileData.geminiKey || "")
    } catch (error) {
      console.error("Failed to load profile:", error)
      setError(`Failed to load profile: ${error instanceof Error ? error.message : 'Unknown error'}`)
    } finally {
      setLoading(false)
    }
  }

  const handleUpdateAPIKeys = async () => {
    try {
      setSaving(true)
      setError(null)
      setSuccess(null)
      
      const updatedProfile = await userApi.updateAPIKeys(
        {
          openaiKey: openaiKey || undefined,
          geminiKey: geminiKey || undefined,
        },
        getToken
      )
      setProfile(updatedProfile)
      setSuccess("API keys updated successfully")
    } catch (error) {
      console.error("Failed to update API keys:", error)
      setError(`Failed to update API keys: ${error instanceof Error ? error.message : 'Unknown error'}`)
    } finally {
      setSaving(false)
    }
  }

  const handleDeleteAPIKey = async (keyType: string) => {
    try {
      setSaving(true)
      setError(null)
      setSuccess(null)
      
      await userApi.deleteAPIKey(keyType, getToken)
      
      // Update local state
      if (keyType === "openai") {
        setOpenaiKey("")
        setProfile((prev: any) => ({ ...prev, openaiKey: null }))
      } else if (keyType === "gemini") {
        setGeminiKey("")
        setProfile((prev: any) => ({ ...prev, geminiKey: null }))
      }
      
      setSuccess(`${keyType.toUpperCase()} API key deleted successfully`)
    } catch (error) {
      console.error(`Failed to delete ${keyType} API key:`, error)
      setError(`Failed to delete ${keyType} API key: ${error instanceof Error ? error.message : 'Unknown error'}`)
    } finally {
      setSaving(false)
    }
  }

  if (!isLoaded || loading) {
    return (
      <div className="min-h-screen bg-white flex items-center justify-center">
        <div className="animate-pulse text-black text-xl">Loading...</div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-white text-black">
      {/* Header */}
      <div className="border-b border-black">
        <div className="max-w-6xl mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-4">
              <Button
                onClick={() => router.push('/notes')}
                variant="ghost"
                size="sm"
                className="text-black hover:bg-black hover:text-white border border-black"
              >
                <ArrowLeft size={16} className="mr-2" />
                Back to Notes
              </Button>
              <h1 className="text-2xl font-bold text-black">Profile Settings</h1>
            </div>
            <UserButton 
              appearance={{
                elements: {
                  avatarBox: "w-8 h-8 border border-black"
                }
              }}
            />
          </div>
        </div>
      </div>

      <div className="max-w-4xl mx-auto px-4 py-8">
        {/* Success/Error Messages */}
        {error && (
          <div className="mb-6 p-4 bg-black text-black border border-black rounded-lg">
            {error}
          </div>
        )}
        
        {success && (
          <div className="mb-6 p-4 bg-white text-black border border-black rounded-lg">
            {success}
          </div>
        )}

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Profile Information */}
          <Card className="p-6 border border-black bg-white">
            <div className="flex items-center space-x-3 mb-6">
              <User className="w-6 h-6 text-black" />
              <h2 className="text-xl font-semibold text-black">Profile Information</h2>
            </div>
            
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-2 text-black">Full Name</label>
                <div className="text-lg text-black">
                  {user?.firstName} {user?.lastName || 'N/A'}
                </div>
              </div>
              
              <div>
                <label className="block text-sm font-medium mb-2 text-black">Username</label>
                <div className="text-lg text-black">
                  {user?.username || profile?.username || 'N/A'}
                </div>
              </div>
              
              <div>
                <label className="block text-sm font-medium mb-2 text-black">Email</label>
                <div className="text-lg flex items-center space-x-2 text-black">
                  <Mail className="w-4 h-4 text-black" />
                  <span>{user?.primaryEmailAddress?.emailAddress || profile?.email || 'N/A'}</span>
                </div>
              </div>
              
              <Separator className="bg-black" />
              
              <div>
                <label className="block text-sm font-medium mb-2 text-black">Member Since</label>
                <div className="text-black">
                  {profile?.createdAt ? 
                    new Date(profile.createdAt).toLocaleDateString('en-US', {
                      year: 'numeric',
                      month: 'long',
                      day: 'numeric'
                    }) : 'N/A'
                  }
                </div>
              </div>
            </div>
          </Card>

          {/* API Keys Management */}
          <Card className="p-6 border border-black bg-white">
            <div className="flex items-center space-x-3 mb-6">
              <Key className="w-6 h-6 text-black" />
              <h2 className="text-xl font-semibold text-black">AI API Keys</h2>
            </div>
            
            <div className="space-y-6">
              {/* Current Status */}
              <div className="space-y-3">
                <h3 className="font-medium text-black">Current Status</h3>
                <div className="flex flex-wrap gap-4">
                  <div className="flex-1 min-w-[280px] flex flex-col sm:flex-row sm:items-center sm:justify-between p-4 border border-black rounded space-y-3 sm:space-y-0">
                    <span className="text-sm text-black font-medium">OpenAI</span>
                    <div className="flex items-center space-x-3">
                      <span className={`text-xs px-3 py-1 rounded flex-shrink-0 ${
                        profile?.hasOpenaiKey 
                          ? 'bg-black text-white' 
                          : 'bg-white text-black border border-black'
                      }`}>
                        {profile?.hasOpenaiKey ? 'Connected' : 'Not Set'}
                      </span>
                      {profile?.hasOpenaiKey && (
                        <Button
                          onClick={() => handleDeleteAPIKey('openai')}
                          disabled={saving}
                          variant="outline"
                          size="sm"
                          className="text-black hover:bg-black hover:text-white px-3 py-1 text-xs border border-black flex-shrink-0"
                        >
                          Remove
                        </Button>
                      )}
                    </div>
                  </div>
                  
                  <div className="flex-1 min-w-[280px] flex flex-col sm:flex-row sm:items-center sm:justify-between p-4 border border-black rounded space-y-3 sm:space-y-0">
                    <span className="text-sm text-black font-medium">Gemini</span>
                    <div className="flex items-center space-x-3">
                      <span className={`text-xs px-3 py-1 rounded flex-shrink-0 ${
                        profile?.hasGeminiKey 
                          ? 'bg-black text-white' 
                          : 'bg-white text-black border border-black'
                      }`}>
                        {profile?.hasGeminiKey ? 'Connected' : 'Not Set'}
                      </span>
                      {profile?.hasGeminiKey && (
                        <Button
                          onClick={() => handleDeleteAPIKey('gemini')}
                          disabled={saving}
                          variant="outline"
                          size="sm"
                          className="text-black hover:bg-black hover:text-white px-3 py-1 text-xs border border-black flex-shrink-0"
                        >
                          Remove
                        </Button>
                      )}
                    </div>
                  </div>
                </div>
              </div>

              <Separator className="bg-black" />

              {/* Add/Update API Keys */}
              <div className="space-y-4">
                <h3 className="font-medium text-black">Add or Update API Keys</h3>
                
                {/* OpenAI Key */}
                <div>
                  <label className="block text-sm font-medium mb-2 text-black">OpenAI API Key</label>
                  <div className="relative">
                    <Input
                      type={showOpenaiKey ? 'text' : 'password'}
                      value={openaiKey}
                      onChange={(e) => setOpenaiKey(e.target.value)}
                      placeholder="sk-..."
                      className="pr-10 border-black bg-white text-black"
                    />
                    <button
                      type="button"
                      onClick={() => setShowOpenaiKey(!showOpenaiKey)}
                      className="absolute right-2 top-1/2 transform -translate-y-1/2 text-black hover:text-gray-600"
                    >
                      {showOpenaiKey ? <EyeOff size={16} /> : <Eye size={16} />}
                    </button>
                  </div>
                </div>

                {/* Gemini Key */}
                <div>
                  <label className="block text-sm font-medium mb-2 text-black">Gemini API Key</label>
                  <div className="relative">
                    <Input
                      type={showGeminiKey ? 'text' : 'password'}
                      value={geminiKey}
                      onChange={(e) => setGeminiKey(e.target.value)}
                      placeholder="AI..."
                      className="pr-10 border-black bg-white text-black"
                    />
                    <button
                      type="button"
                      onClick={() => setShowGeminiKey(!showGeminiKey)}
                      className="absolute right-2 top-1/2 transform -translate-y-1/2 text-black hover:text-gray-600"
                    >
                      {showGeminiKey ? <EyeOff size={16} /> : <Eye size={16} />}
                    </button>
                  </div>
                </div>

                <Button
                  onClick={handleUpdateAPIKeys}
                  disabled={saving || (!openaiKey.trim() && !geminiKey.trim())}
                  className="w-full bg-black text-white hover:bg-gray-800"
                >
                  {saving ? (
                    <div className="flex items-center space-x-2">
                      <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                      <span>Saving...</span>
                    </div>
                  ) : (
                    <div className="flex items-center space-x-2">
                      <Save size={16} />
                      <span>Save API Keys</span>
                    </div>
                  )}
                </Button>
              </div>

              {/* Info */}
              <div className="text-xs text-gray-600 space-y-1">
                <p>• Your API keys are encrypted and stored securely</p>
                <p>• These keys are used for AI chat features in your notes</p>
                <p>• You can remove keys at any time</p>
              </div>
            </div>
          </Card>
        </div>
      </div>
    </div>
  )
} 