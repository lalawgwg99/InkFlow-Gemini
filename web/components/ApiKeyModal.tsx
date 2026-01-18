"use client";

import { useState, useEffect } from "react";
import { Key, Lock } from "lucide-react";

export default function ApiKeyModal({
    onSave,
    forceShow = false,
    onClose
}: {
    onSave: (key: string) => void;
    forceShow?: boolean;
    onClose?: () => void;
}) {
    const [apiKey, setApiKey] = useState("");
    const [isOpen, setIsOpen] = useState(false);

    useEffect(() => {
        if (forceShow) {
            setIsOpen(true);
            return;
        }
        const key = localStorage.getItem("GEMINI_API_KEY");
        if (!key) {
            const timer = setTimeout(() => setIsOpen(true), 0);
            return () => clearTimeout(timer);
        } else {
            // Defers the callback to avoid update during render
            const timer = setTimeout(() => onSave(key), 0);
            return () => clearTimeout(timer);
        }
    }, [onSave, forceShow]);

    const handleSave = () => {
        if (apiKey.trim().length > 10) {
            localStorage.setItem("GEMINI_API_KEY", apiKey);
            onSave(apiKey);
            setIsOpen(false);
            onClose?.();
        } else {
            alert("請輸入有效的 Gemini API Key");
        }
    };

    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-sm p-4">
            <div className="bg-gray-900 border border-gray-700 rounded-2xl p-6 max-w-md w-full shadow-2xl">
                <div className="text-center mb-6">
                    <div className="mx-auto w-12 h-12 bg-purple-600/20 rounded-full flex items-center justify-center mb-4">
                        <Key className="w-6 h-6 text-purple-400" />
                    </div>
                    <h2 className="text-2xl font-bold bg-gradient-to-r from-purple-400 to-pink-400 bg-clip-text text-transparent">
                        輸入您的金鑰
                    </h2>
                    <p className="text-gray-400 mt-2 text-sm">
                        您的 API Key 僅會儲存在瀏覽器中，直接與 Google 通訊，保障隱私安全。
                    </p>
                </div>

                <div className="space-y-4">
                    <div>
                        <label className="block text-xs font-medium text-gray-400 mb-1 ml-1">
                            GEMINI_API_KEY
                        </label>
                        <div className="relative">
                            <input
                                type="password"
                                value={apiKey}
                                onChange={(e) => setApiKey(e.target.value)}
                                placeholder="AIzaSy..."
                                className="w-full bg-gray-800 border border-gray-700 rounded-xl px-4 py-3 pl-10 text-white focus:outline-none focus:ring-2 focus:ring-purple-500 transition"
                            />
                            <Lock className="w-4 h-4 text-gray-500 absolute left-3 top-3.5" />
                        </div>
                    </div>

                    <button
                        onClick={handleSave}
                        className="w-full bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-500 hover:to-pink-500 text-white font-bold py-3 rounded-xl transition transform active:scale-95 shadow-lg"
                    >
                        開始使用
                    </button>

                    <div className="text-center">
                        <a
                            href="https://aistudio.google.com/app/apikey"
                            target="_blank"
                            rel="noopener noreferrer"
                            className="text-xs text-gray-500 hover:text-purple-400 transition underline"
                        >
                            沒有 Key？點此免費獲取
                        </a>
                    </div>
                </div>
            </div>
        </div>
    );
}
