"use client";

import { SignInButton, SignUpButton, useAuth } from "@clerk/nextjs";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import Image from "next/image";
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

  if (!isLoaded || isSignedIn) {
    return (
      <div className="min-h-screen bg-black flex items-center justify-center">
        <button className="px-8 py-4 bg-gradient-to-r from-black to-white rounded-full text-white text-xl font-medium shadow-lg hover:from-black hover:to-white transition-all duration-300 animate-pulse-glow">
          Loading...
        </button>
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
            <div className="flex items-center justify-center gap-4 m-4">
              <Image
                src="/klara.png"
                alt="Klara Logo"
                width={80}
                height={80}
                className="rounded-lg"
              />
            </div>
            <h1 className="text-6xl md:text-8xl font-bold tracking-tight">
              <span className="bg-gradient-to-r from-white to-white bg-clip-text text-transparent">
                Kl
              </span>
              <span className="text-white">ara</span>
            </h1>
            <p className="text-xl md:text-2xl text-white max-w-2xl mx-auto">
              The future of note-taking. Intelligent, contextual, and infinitely
              powerful.
            </p>
          </div>

          {/* Feature highlights */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mt-16 mb-16">
            <Card className="bg-black border-white p-6 hover:bg-white hover:text-black transition-all duration-300 group">
              <div className="text-center space-y-3">
                <div className="w-12 h-12 mx-auto bg-white rounded-lg flex items-center justify-center group-hover:bg-black">
                  <span className="text-2xl">🧠</span>
                </div>
                <h3 className="text-lg font-semibold text-white group-hover:text-black">AI-Powered</h3>
                <p className="text-white text-sm group-hover:text-black">
                  Intelligent insights and contextual understanding
                </p>
              </div>
            </Card>

            <Card className="bg-black border-white p-6 hover:bg-white hover:text-black transition-all duration-300 group">
              <div className="text-center space-y-3">
                <div className="w-12 h-12 mx-auto bg-white rounded-lg flex items-center justify-center group-hover:bg-black">
                  <span className="text-2xl">⚡</span>
                </div>
                <h3 className="text-lg font-semibold text-white group-hover:text-black">
                  Lightning Fast
                </h3>
                <p className="text-white text-sm group-hover:text-black">
                  Instant search and seamless performance
                </p>
              </div>
            </Card>

            <Card className="bg-black border-white p-6 hover:bg-white hover:text-black transition-all duration-300 group">
              <div className="text-center space-y-3">
                <div className="w-12 h-12 mx-auto bg-white rounded-lg flex items-center justify-center group-hover:bg-black">
                  <span className="text-2xl">🔒</span>
                </div>
                <h3 className="text-lg font-semibold text-white group-hover:text-black">Secure</h3>
                <p className="text-white text-sm group-hover:text-black">
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
                className="bg-white text-black hover:bg-black hover:text-white font-semibold px-8 py-3 text-lg animate-pulse-glow"
              >
                Get Started
              </Button>
            </SignUpButton>
            
            <SignInButton>
              <Button 
                size="lg" 
                className="border border-white bg-black text-white hover:bg-white hover:text-black hover:border-black font-semibold px-8 py-3 text-lg shadow-[0_0_10px_rgba(255,255,255,0.3)] hover:shadow-none transition-all duration-300"
              >
                Sign In
              </Button>
            </SignInButton>
          </div>

          {/* Subtle call to action text */}
          <p className="text-white text-sm mt-8">
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
