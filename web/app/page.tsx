"use client";

import { useState } from "react";
import ApiKeyModal from "@/components/ApiKeyModal";
import IdeaGenerator from "@/components/IdeaGenerator";
import Writer from "@/components/Writer";
import { StyleConfig, STYLES } from "@/lib/styles";
import { BookOpen, Sparkles } from "lucide-react";

export default function Home() {
  const [apiKey, setApiKey] = useState("");
  const [activeTab, setActiveTab] = useState<"idea" | "write">("idea");

  // State for Writer
  const [writerInput, setWriterInput] = useState("");
  const [writerStyle, setWriterStyle] = useState<StyleConfig>(STYLES[0]);

  const handleIdeaSelected = (idea: string, style: StyleConfig) => {
    setWriterInput(idea);
    setWriterStyle(style);
    setActiveTab("write");
  };

  return (
    <main className="min-h-screen bg-black text-gray-100 flex flex-col items-center p-4 md:p-8 relative overflow-hidden">
      {/* Background Decoration */}
      <div className="absolute top-0 left-0 w-[500px] h-[500px] bg-purple-600 rounded-full mix-blend-multiply filter blur-[128px] opacity-20 -translate-x-1/2 -translate-y-1/2 animate-blob pointer-events-none"></div>
      <div className="absolute top-0 right-0 w-[500px] h-[500px] bg-blue-600 rounded-full mix-blend-multiply filter blur-[128px] opacity-20 translate-x-1/2 -translate-y-1/2 animate-blob pointer-events-none animation-delay-2000"></div>

      <ApiKeyModal onSave={setApiKey} />

      {/* Header */}
      <header className="mb-12 text-center z-10">
        <h1 className="text-4xl md:text-6xl font-black mb-4 tracking-tight">
          InkFlow<span className="text-transparent bg-clip-text bg-gradient-to-r from-purple-400 to-pink-400">.AI</span>
        </h1>
        <p className="text-gray-400 max-w-lg mx-auto">
          專為創作者打造的 Gemini 3 Pro 寫作助手。
          <br />靈感生成、風格模仿、多平台輸出。
        </p>
      </header>

      {/* Navigation Tabs */}
      <div className="flex bg-gray-900/80 backdrop-blur-md p-1 rounded-2xl border border-gray-800 mb-8 z-10 shadow-xl">
        <button
          onClick={() => setActiveTab("idea")}
          className={`px-6 py-3 rounded-xl font-bold flex items-center gap-2 transition-all ${activeTab === "idea"
              ? "bg-gray-800 text-white shadow-lg border border-gray-700"
              : "text-gray-400 hover:text-gray-200"
            }`}
        >
          <Sparkles className={`w-4 h-4 ${activeTab === "idea" ? "text-yellow-400" : ""}`} />
          靈感生成
        </button>
        <button
          onClick={() => setActiveTab("write")}
          className={`px-6 py-3 rounded-xl font-bold flex items-center gap-2 transition-all ${activeTab === "write"
              ? "bg-gray-800 text-white shadow-lg border border-gray-700"
              : "text-gray-400 hover:text-gray-200"
            }`}
        >
          <BookOpen className={`w-4 h-4 ${activeTab === "write" ? "text-purple-400" : ""}`} />
          自由寫作
        </button>
      </div>

      {/* Main Content Area */}
      <div className="w-full max-w-5xl z-10">
        {activeTab === "idea" ? (
          <div className="animate-in fade-in slide-in-from-bottom-4 duration-500">
            <IdeaGenerator onSelectIdea={handleIdeaSelected} />
          </div>
        ) : (
          <div className="animate-in fade-in slide-in-from-bottom-4 duration-500">
            <Writer
              styleConfig={writerStyle}
              initialInput={writerInput}
              onBack={() => setActiveTab("idea")}
            />
          </div>
        )}
      </div>

      {/* Footer */}
      <footer className="mt-24 text-gray-600 text-sm">
        Powered by Google Gemini 3 Pro • {new Date().getFullYear()} InkFlow AI
      </footer>
    </main>
  );
}
