import { SignIn } from '@clerk/nextjs'

export default function Page() {
  return (
    <div className="min-h-screen bg-black flex items-center justify-center relative overflow-hidden">
      {/* Background effects */}
      <div className="absolute inset-0 matrix-bg" />
      <div className="absolute inset-0 opacity-10">
        <div className="absolute inset-0" style={{
          backgroundImage: 'linear-gradient(rgba(255,255,255,0.1) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,0.1) 1px, transparent 1px)',
          backgroundSize: '50px 50px'
        }} />
      </div>

      <div className="relative z-10 w-full max-w-md mx-auto p-6">
        <div className="text-center mb-8">
          <h1 className="text-4xl font-bold text-white mb-2">Welcome Back</h1>
          <p className="text-gray-400">Sign in to continue to Klara</p>
        </div>
        
        <div className="glow-border rounded-lg p-1 bg-gradient-to-r from-white/10 to-transparent">
          <SignIn 
            appearance={{
              elements: {
                rootBox: "mx-auto",
                card: "bg-gray-900/90 shadow-2xl border-0 rounded-lg",
                headerTitle: "text-white",
                headerSubtitle: "text-gray-400",
                formButtonPrimary: "bg-white text-black hover:bg-gray-200 font-semibold",
                formFieldInput: "bg-gray-800 border-gray-700 text-white",
                formFieldLabel: "text-gray-300",
                identityPreviewText: "text-gray-300",
                identityPreviewEditButton: "text-white",
                footer: "hidden"
              }
            }}
          />
        </div>
      </div>
    </div>
  )
} 