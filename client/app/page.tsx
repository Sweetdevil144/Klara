"use client";

import { SignInButton, SignUpButton, useAuth } from "@clerk/nextjs";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";

export default function Home() {
  const { isLoaded, isSignedIn } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (isLoaded && isSignedIn) {
      router.push("/notes");
    }
  }, [isLoaded, isSignedIn, router]);

  if (!isLoaded) {
    return (
      <div className="min-h-screen bg-black flex items-center justify-center">
        <div className="animate-pulse-glow text-white text-xl">Loading...</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-black relative overflow-hidden">
      {/* Subtle background pattern */}
      <div className="absolute inset-0">
        <div
          className="absolute inset-0 opacity-5"
          style={{
            backgroundImage:
              "linear-gradient(rgba(255,255,255,0.1) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,0.1) 1px, transparent 1px)",
            backgroundSize: "50px 50px",
          }}
        />
      </div>

      <main className="relative z-10 flex flex-col items-center justify-center min-h-screen px-4">
        {/* Hero Section */}
        <div className="text-center space-y-8 max-w-4xl mx-auto animate-fade-in-up">
          {/* Logo/Title */}
          <div className="space-y-4">
            <h1 className="text-6xl md:text-8xl font-bold tracking-tight">
              <span className="bg-gradient-to-r from-white to-gray-400 bg-clip-text text-transparent">
                Kl
              </span>
              <span className="text-white">ara</span>
            </h1>
            <p className="text-xl md:text-2xl text-gray-400 max-w-2xl mx-auto">
              The future of note-taking. Intelligent, contextual, and infinitely
              powerful.
            </p>
          </div>

          {/* Feature highlights */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mt-16 mb-16">
            <Card className="bg-gray-900 border-gray-800 p-6 hover:bg-gray-800 transition-all duration-300">
              <div className="text-center space-y-3">
                <div className="w-12 h-12 mx-auto bg-gray-800 rounded-lg flex items-center justify-center">
                  <span className="text-2xl">🧠</span>
                </div>
                <h3 className="text-lg font-semibold text-white">AI-Powered</h3>
                <p className="text-gray-400 text-sm">
                  Intelligent insights and contextual understanding
                </p>
              </div>
            </Card>

            <Card className="bg-gray-900 border-gray-800 p-6 hover:bg-gray-800 transition-all duration-300">
              <div className="text-center space-y-3">
                <div className="w-12 h-12 mx-auto bg-gray-800 rounded-lg flex items-center justify-center">
                  <span className="text-2xl">⚡</span>
                </div>
                <h3 className="text-lg font-semibold text-white">
                  Lightning Fast
                </h3>
                <p className="text-gray-400 text-sm">
                  Instant search and seamless performance
                </p>
              </div>
            </Card>

            <Card className="bg-gray-900 border-gray-800 p-6 hover:bg-gray-800 transition-all duration-300">
              <div className="text-center space-y-3">
                <div className="w-12 h-12 mx-auto bg-gray-800 rounded-lg flex items-center justify-center">
                  <span className="text-2xl">🔒</span>
                </div>
                <h3 className="text-lg font-semibold text-white">Secure</h3>
                <p className="text-gray-400 text-sm">
                  Your thoughts, protected and private
                </p>
              </div>
            </Card>
          </div>

          {/* CTA Buttons */}
          <div className="flex flex-col sm:flex-row gap-4 justify-center items-center">
            <SignUpButton>
              <Button 
                size="lg" 
                className="bg-white text-black hover:bg-gray-200 font-semibold px-8 py-3 text-lg animate-pulse-glow"
              >
                Get Started
              </Button>
            </SignUpButton>
            
            <SignInButton>
              <Button 
                variant="outline" 
                size="lg" 
                className="border-gray-600 text-white hover:bg-gray-900 font-semibold px-8 py-3 text-lg"
              >
                Sign In
              </Button>
            </SignInButton>
          </div>

          {/* Subtle call to action text */}
          <p className="text-gray-500 text-sm mt-8">
            Join the future of productivity. Start your journey today.
          </p>
        </div>

        {/* Floating elements for visual appeal */}
        <div className="absolute top-20 left-10 w-2 h-2 bg-white/20 rounded-full animate-pulse" />
        <div
          className="absolute top-40 right-20 w-1 h-1 bg-white/30 rounded-full animate-pulse"
          style={{ animationDelay: "1s" }}
        />
        <div
          className="absolute bottom-20 left-20 w-1.5 h-1.5 bg-white/20 rounded-full animate-pulse"
          style={{ animationDelay: "2s" }}
        />
        <div
          className="absolute bottom-40 right-10 w-1 h-1 bg-white/25 rounded-full animate-pulse"
          style={{ animationDelay: "0.5s" }}
        />
      </main>
    </div>
  );
}
