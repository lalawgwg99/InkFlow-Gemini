import { GoogleGenerativeAI } from "@google/generative-ai";

export async function generateContent(apiKey: string, prompt: string, modelName: string = "gemini-3-flash") {
    try {
        const genAI = new GoogleGenerativeAI(apiKey);
        const model = genAI.getGenerativeModel({ model: modelName });

        const result = await model.generateContent(prompt);
        const response = await result.response;
        return response.text();
    } catch (error) {
        console.error("Gemini API Error:", error);
        throw error;
    }
}
