"use client";

import { useState } from "react";
import { StyleConfig, STYLES } from "@/lib/styles";
import { generateContent } from "@/lib/gemini";
import { Sparkles, ArrowRight } from "lucide-react";
import { motion } from "framer-motion";

export default function IdeaGenerator({ onSelectIdea }: { onSelectIdea: (idea: string, style: StyleConfig) => void }) {
    const [selectedStyle, setSelectedStyle] = useState<StyleConfig>(STYLES[0]);
    const [ideas, setIdeas] = useState<string[]>([]);
    const [isLoading, setIsLoading] = useState(false);

    const generateIdeas = async () => {
        const apiKey = localStorage.getItem("GEMINI_API_KEY");
        if (!apiKey) {
            alert("請先設定 API Key");
            return;
        }

        setIsLoading(true);
        try {
            const prompt = `
${selectedStyle.prompt}

---
任務：請直接提供 5 個符合此風格的「高病毒傳播潛力」寫作主題。
格式：純文字列表，每行一個主題，不要編號，不要其他廢話。
`;
            const result = await generateContent(apiKey, prompt);
            const ideaList = result.split("\n").filter(line => line.trim().length > 0).slice(0, 5);
            setIdeas(ideaList);
        } catch (error) {
            alert("生成失敗");
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className="space-y-8">
            <div className="text-center space-y-4">
                <h2 className="text-3xl font-bold text-white">
                    選擇您的<span className="text-transparent bg-clip-text bg-gradient-to-r from-yellow-400 to-orange-500">寫作 Muse</span>
                </h2>
                <p className="text-gray-400">不知道寫什麼？讓 AI 幫您尋找百萬流量的靈感。</p>
            </div>

            {/* Style Selector */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                {STYLES.map((style) => (
                    <button
                        key={style.id}
                        onClick={() => setSelectedStyle(style)}
                        className={`p-4 rounded-xl border text-left transition relative overflow-hidden group ${selectedStyle.id === style.id
                                ? "bg-purple-900/20 border-purple-500 ring-2 ring-purple-500/20"
                                : "bg-gray-900 border-gray-800 hover:border-gray-600"
                            }`}
                    >
                        <div className="font-bold text-lg mb-2 text-white">{style.name}</div>
                        <div className="text-xs text-gray-400 line-clamp-2">{style.description}</div>
                        {selectedStyle.id === style.id && (
                            <div className="absolute top-2 right-2 w-2 h-2 bg-purple-500 rounded-full animate-pulse" />
                        )}
                    </button>
                ))}
            </div>

            {/* Generate Button */}
            <div className="text-center">
                <button
                    onClick={generateIdeas}
                    disabled={isLoading}
                    className="group relative inline-flex items-center justify-center px-8 py-3 font-bold text-white transition-all duration-200 bg-gradient-to-r from-indigo-500 via-purple-500 to-pink-500 rounded-full hover:scale-105 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-600 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                    {isLoading ? (
                        "召喚靈感中..."
                    ) : (
                        <>
                            <Sparkles className="w-5 h-5 mr-2 animate-pulse" />
                            生成 5 個靈感
                        </>
                    )}
                    <div className="absolute inset-0 rounded-full ring-2 ring-white/20 group-hover:ring-white/40 transition-all" />
                </button>
            </div>

            {/* Idea Cards */}
            {ideas.length > 0 && (
                <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    className="grid gap-4"
                >
                    {ideas.map((idea, index) => (
                        <div
                            key={index}
                            className="bg-gray-800/50 hover:bg-gray-800 border border-gray-700 hover:border-purple-500/50 p-4 rounded-xl flex items-center justify-between cursor-pointer transition group"
                            onClick={() => onSelectIdea(idea, selectedStyle)}
                        >
                            <span className="text-lg text-gray-200">{idea}</span>
                            <ArrowRight className="w-5 h-5 text-gray-500 group-hover:text-purple-400 transition" />
                        </div>
                    ))}
                </motion.div>
            )}
        </div>
    );
}
