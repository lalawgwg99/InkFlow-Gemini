"use client";

import { useState } from "react";
import { generateContent } from "@/lib/gemini";
import { StyleConfig } from "@/lib/styles";
import { Copy, RefreshCw, PenTool } from "lucide-react";
import { motion } from "framer-motion";

interface WriterProps {
    styleConfig: StyleConfig;
    initialInput?: string;
    onBack: () => void;
}

export default function Writer({ styleConfig, initialInput = "", onBack }: WriterProps) {
    const [input, setInput] = useState(initialInput);
    const [output, setOutput] = useState("");
    const [isLoading, setIsLoading] = useState(false);
    const [mode, setMode] = useState<"normal" | "thread">("normal");

    const handleGenerate = async () => {
        const apiKey = localStorage.getItem("GEMINI_API_KEY");
        if (!apiKey) {
            alert("請先設定 API Key");
            return;
        }

        setIsLoading(true);
        try {
            let fullPrompt = `
${styleConfig.prompt}

---
使用者的主題/內容:
${input}
---
`;
            if (mode === "thread") {
                fullPrompt += `
請將上述內容改寫為 Twitter/X 貼文串 (Thread)。
格式要求：
- 每條推文不超過 280 字
- 加上編號 (1/5, 2/5...)
- 第一條要有強烈鉤子
- 最後一條要有 CTA
`;
            }

            const result = await generateContent(apiKey, fullPrompt);
            setOutput(result);
        } catch (error) {
            alert("生成失敗，請檢查 API Key 或稍後再試");
        } finally {
            setIsLoading(false);
        }
    };

    const copyToClipboard = () => {
        navigator.clipboard.writeText(output);
        alert("已複製到剪貼簿！");
    };

    return (
        <div className="space-y-6">
            <div className="flex items-center justify-between">
                <button onClick={onBack} className="text-gray-400 hover:text-white transition">
                    ← 返回選擇
                </button>
                <span className="text-sm px-3 py-1 bg-purple-900/50 text-purple-300 rounded-full border border-purple-700">
                    當前風格: {styleConfig.name}
                </span>
            </div>

            <div className="grid md:grid-cols-2 gap-6">
                {/* Input Area */}
                <div className="bg-gray-900 rounded-2xl p-6 border border-gray-800">
                    <label className="block text-sm font-medium text-gray-400 mb-2">
                        輸入主題或草稿
                    </label>
                    <textarea
                        value={input}
                        onChange={(e) => setInput(e.target.value)}
                        className="w-full h-64 bg-gray-800 border border-gray-700 rounded-xl p-4 text-white focus:ring-2 focus:ring-purple-500 transition resize-none"
                        placeholder="例如：為什麼大部分人無法自律..."
                    />

                    <div className="flex justify-between items-center mt-4">
                        <div className="flex gap-2">
                            <button
                                onClick={() => setMode("normal")}
                                className={`px-3 py-1 rounded-lg text-sm transition ${mode === "normal" ? "bg-gray-700 text-white" : "text-gray-500 hover:text-gray-300"}`}
                            >
                                一般文章
                            </button>
                            <button
                                onClick={() => setMode("thread")}
                                className={`px-3 py-1 rounded-lg text-sm transition ${mode === "thread" ? "bg-gray-700 text-white" : "text-gray-500 hover:text-gray-300"}`}
                            >
                                X 貼文串
                            </button>
                        </div>

                        <button
                            onClick={handleGenerate}
                            disabled={isLoading || !input}
                            className="flex items-center gap-2 bg-gradient-to-r from-purple-600 to-pink-600 text-white px-6 py-2 rounded-xl font-bold hover:shadow-lg hover:shadow-purple-500/30 disabled:opacity-50 transition"
                        >
                            {isLoading ? <RefreshCw className="animate-spin w-4 h-4" /> : <PenTool className="w-4 h-4" />}
                            開始撰寫
                        </button>
                    </div>
                </div>

                {/* Output Area */}
                <motion.div
                    initial={{ opacity: 0, x: 20 }}
                    animate={{ opacity: 1, x: 0 }}
                    className="bg-gray-900 rounded-2xl p-6 border border-gray-800 relative"
                >
                    <label className="block text-sm font-medium text-gray-400 mb-2">
                        AI 生成結果
                    </label>
                    <div className="w-full h-64 bg-black/50 border border-gray-800 rounded-xl p-4 text-gray-300 overflow-y-auto whitespace-pre-wrap">
                        {output || <span className="text-gray-600 italic">等待生成...</span>}
                    </div>

                    {output && (
                        <button
                            onClick={copyToClipboard}
                            className="absolute top-6 right-6 p-2 bg-gray-800 hover:bg-gray-700 rounded-lg text-gray-300 transition"
                            title="複製內容"
                        >
                            <Copy className="w-4 h-4" />
                        </button>
                    )}
                </motion.div>
            </div>
        </div>
    );
}
