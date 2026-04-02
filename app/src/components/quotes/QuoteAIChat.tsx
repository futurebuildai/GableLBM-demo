import { useState, useRef, useEffect } from 'react';
import { Send, Mic, MicOff, Sparkles, ChevronDown, ChevronUp } from 'lucide-react';
import { QuoteService } from '../../services/QuoteService';

interface Message {
    role: 'user' | 'assistant';
    text: string;
}

interface QuoteAIChatProps {
    quoteId: string;
    onQuoteUpdated: () => void;
}

export default function QuoteAIChat({ quoteId, onQuoteUpdated }: QuoteAIChatProps) {
    const [messages, setMessages] = useState<Message[]>([]);
    const [input, setInput] = useState('');
    const [loading, setLoading] = useState(false);
    const [collapsed, setCollapsed] = useState(false);
    const [isRecording, setIsRecording] = useState(false);
    const messagesEndRef = useRef<HTMLDivElement>(null);
    const recognitionRef = useRef<any>(null);

    useEffect(() => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    }, [messages]);

    async function handleSend() {
        const text = input.trim();
        if (!text || loading) return;

        setInput('');
        setMessages(prev => [...prev, { role: 'user', text }]);
        setLoading(true);

        try {
            const result = await QuoteService.aiEditQuote(quoteId, text);
            setMessages(prev => [...prev, { role: 'assistant', text: result.explanation }]);
            onQuoteUpdated();
        } catch (err) {
            setMessages(prev => [...prev, {
                role: 'assistant',
                text: `Error: ${err instanceof Error ? err.message : 'Failed to process command'}`
            }]);
        } finally {
            setLoading(false);
        }
    }

    function toggleVoice() {
        if (isRecording) {
            recognitionRef.current?.stop();
            setIsRecording(false);
            return;
        }

        const SpeechRecognition = (window as any).SpeechRecognition || (window as any).webkitSpeechRecognition;
        if (!SpeechRecognition) return;

        const recognition = new SpeechRecognition();
        recognition.continuous = false;
        recognition.interimResults = true;
        recognition.lang = 'en-US';

        recognition.onresult = (event: any) => {
            let transcript = '';
            for (let i = 0; i < event.results.length; i++) {
                transcript += event.results[i][0].transcript;
            }
            setInput(transcript);
        };

        recognition.onend = () => setIsRecording(false);
        recognition.onerror = () => setIsRecording(false);

        recognitionRef.current = recognition;
        recognition.start();
        setIsRecording(true);
    }

    function handleKeyDown(e: React.KeyboardEvent) {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
            handleSend();
        }
    }

    return (
        <div className="rounded-xl border border-white/10 bg-surface-2 overflow-hidden">
            {/* Header */}
            <button
                onClick={() => setCollapsed(!collapsed)}
                className="w-full flex items-center justify-between px-4 py-3 hover:bg-surface-3 transition-colors"
            >
                <div className="flex items-center gap-2">
                    <Sparkles className="w-4 h-4 text-violet-400" />
                    <span className="text-sm font-semibold text-white">AI Quote Assistant</span>
                </div>
                {collapsed ? <ChevronUp className="w-4 h-4 text-zinc-500" /> : <ChevronDown className="w-4 h-4 text-zinc-500" />}
            </button>

            {!collapsed && (
                <>
                    {/* Messages */}
                    <div className="max-h-64 overflow-y-auto px-4 py-3 space-y-3 border-t border-white/5">
                        {messages.length === 0 && (
                            <p className="text-xs text-zinc-500 italic">
                                Try: "Change qty of 2x4s to 50" or "Add 20 sheets of OSB"
                            </p>
                        )}
                        {messages.map((msg, i) => (
                            <div key={i} className={`flex ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}>
                                <div className={`max-w-[85%] rounded-lg px-3 py-2 text-sm ${
                                    msg.role === 'user'
                                        ? 'bg-blue-600/20 text-blue-300 border border-blue-500/20'
                                        : 'bg-surface-3 text-zinc-300 border border-white/5'
                                }`}>
                                    {msg.text}
                                </div>
                            </div>
                        ))}
                        {loading && (
                            <div className="flex justify-start">
                                <div className="bg-surface-3 border border-white/5 rounded-lg px-3 py-2 text-sm text-zinc-400">
                                    <span className="animate-pulse">Thinking...</span>
                                </div>
                            </div>
                        )}
                        <div ref={messagesEndRef} />
                    </div>

                    {/* Input */}
                    <div className="flex items-center gap-2 px-3 py-2 border-t border-white/5">
                        <button
                            onClick={toggleVoice}
                            className={`p-2 rounded-lg transition-colors ${
                                isRecording
                                    ? 'bg-red-500/20 text-red-400 animate-pulse'
                                    : 'hover:bg-surface-3 text-zinc-500 hover:text-zinc-300'
                            }`}
                            title={isRecording ? 'Stop recording' : 'Voice input'}
                        >
                            {isRecording ? <MicOff className="w-4 h-4" /> : <Mic className="w-4 h-4" />}
                        </button>
                        <input
                            type="text"
                            value={input}
                            onChange={e => setInput(e.target.value)}
                            onKeyDown={handleKeyDown}
                            placeholder="Type a command..."
                            disabled={loading}
                            className="flex-1 bg-transparent text-sm text-white placeholder-zinc-600 outline-none"
                        />
                        <button
                            onClick={handleSend}
                            disabled={!input.trim() || loading}
                            className="p-2 rounded-lg bg-gable-green/20 text-gable-green hover:bg-gable-green/30 transition-colors disabled:opacity-30 disabled:cursor-not-allowed"
                        >
                            <Send className="w-4 h-4" />
                        </button>
                    </div>
                </>
            )}
        </div>
    );
}
