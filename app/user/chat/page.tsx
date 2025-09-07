"use client";

import type React from "react";
import { useRouter } from "next/navigation";
import { useState, useEffect, useRef, useCallback } from "react";
import { useSelector, useDispatch } from 'react-redux';
import { useSession } from "next-auth/react";
import { useCreateChecklistMutation } from '@/app/store/slices/checklistsApi'
import { RootState, AppDispatch } from '@/app/store/store';
import { fetchChatHistory, sendMessage, addUserMessage, clearError, fetchChatById, clearMessages } from '@/app/store/slices/aiChatSlice';
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import ReactMarkdown from 'react-markdown';
import {
  Bot,
  User,
  FileText,
  DollarSign,
  Building,
  Send,
  Bookmark,
  Languages,
  CheckCircle,
  Mic,
  History,
  Clock,
  Play,
  Square,
  Plus,
} from "lucide-react";

// Minimal Web Speech API typings in the global namespace
declare global {
  interface SpeechRecognition {
    lang: string;
    interimResults: boolean;
    maxAlternatives: number;
    continuous: boolean;
  onresult: ((event: unknown) => void) | null;
    onend: (() => void) | null;
  onerror: ((event: unknown) => void) | null;
    start: () => void;
    stop: () => void;
  }
  type SpeechRecognitionConstructor = new () => SpeechRecognition;
  // Extend Window typings to avoid any-casts
  interface Window {
    webkitSpeechRecognition?: SpeechRecognitionConstructor;
    SpeechRecognition?: SpeechRecognitionConstructor;
    responsiveVoice?: {
      speak: (text: string, voice?: string, opts?: { rate?: number; pitch?: number; onend?: () => void; onerror?: () => void }) => boolean;
      cancel: () => void;
      voiceSupport: () => boolean;
    };
  }
}

// Minimal types for SpeechRecognition event payloads
type SRAlternative = { transcript?: string };
type SRResult = { 0?: SRAlternative };
type SRResultList = ArrayLike<SRResult>;

interface Message {
  id: string;
  type: "user" | "assistant";
  content: string;
  timestamp: string;
  procedures?: Procedure[];
}

interface Procedure {
  id: number;
  title: string;
  icon: string;
  items: string[];
  completed?: boolean;
}

// interface ChatHistory {
//   id: string;
//   title: string;
//   lastMessage: string;
//   timestamp: string;
//   messageCount: number;
// }

export default function ChatPage() {
  const dispatch: AppDispatch = useDispatch();
  const { data: session, status } = useSession();
  const { messages, chatHistory, status: chatStatus, error } = useSelector((state: RootState) => state.aiChat);
  const [inputMessage, setInputMessage] = useState("");
  const [showHistory, setShowHistory] = useState(false);
  const [successMessage, setSuccessMessage] = useState("");
  const router = useRouter();
  const [createChecklist] = useCreateChecklistMutation();
  const [translating, setTranslating] = useState<Record<string, boolean>>({});
  const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'https://ethio-guide-backend.onrender.com/api/v1';

  const isAmharic = (text: string): boolean => /[\u1200-\u137F]/.test(text);

  // Speech recognition (STT) state
  const [isListening, setIsListening] = useState(false);
  const [supportsSTT, setSupportsSTT] = useState(false);
  const recognitionRef = useRef<SpeechRecognition | null>(null);

  // Speech synthesis (TTS) state
  const [speakingMessageId, setSpeakingMessageId] = useState<string | null>(null);
  const utteranceRef = useRef<SpeechSynthesisUtterance | null>(null);
  const [supportsTTS, setSupportsTTS] = useState(false);
  const audioRef = useRef<HTMLAudioElement | null>(null);
  const audioUrlRef = useRef<string | null>(null);
  const [rvReady, setRvReady] = useState(false);
  const rvCheckRef = useRef<number | null>(null);
  const speakingIdRef = useRef<string | null>(null);
  useEffect(() => { speakingIdRef.current = speakingMessageId; }, [speakingMessageId]);
  const lastSpokenIdRef = useRef<string | null>(null);

  // Voice mode and language
  const [voiceMode, setVoiceMode] = useState(false);
  const [voiceLang, setVoiceLang] = useState(() => {
    const nav = (typeof navigator !== 'undefined' && navigator.language) || 'en-US';
    return nav.toLowerCase().startsWith('am') ? 'am-ET' : 'en-US';
  });
  const [hasAmharicVoice, setHasAmharicVoice] = useState(false);

  // Helper to send text (used by voice mode)
  const sendText = useCallback((text: string) => {
    const trimmed = text.trim();
    if (!trimmed) return;
    const token = session?.accessToken;
    const newMessage: Message = {
      id: Date.now().toString(),
      type: "user",
      content: trimmed,
      timestamp: new Date().toLocaleString(),
    };
    dispatch(addUserMessage(newMessage));
    if (token) {
      dispatch(sendMessage({ query: trimmed, token })).then((result) => {
        if (result.meta.requestStatus === "fulfilled") {
          setSuccessMessage("Message sent successfully!");
          setTimeout(() => setSuccessMessage(""), 3000);
        }
      });
    }
    setInputMessage("");
  }, [dispatch, session]);

  useEffect(() => {
    const token = session?.accessToken;
    if (status === 'authenticated' && token) {
      dispatch(fetchChatHistory(token));
    }
  }, [dispatch, session, status]);

  // Initialize Web Speech APIs
  useEffect(() => {
    // Capture audio element for cleanup to avoid ref-churn warning
    const audioElForCleanup = audioRef.current;
    // STT
    const SpeechRecognitionImpl: SpeechRecognitionConstructor | undefined =
      (typeof window !== 'undefined' ? (window.webkitSpeechRecognition || window.SpeechRecognition) : undefined);

    if (SpeechRecognitionImpl) {
      setSupportsSTT(true);
      const recognition = new SpeechRecognitionImpl();
      recognition.lang = voiceLang || 'en-US';
      recognition.interimResults = false;
      recognition.maxAlternatives = 1;
      recognition.continuous = false;

      recognition.onresult = (e: unknown) => {
        const results = (e as { results?: SRResultList }).results;
        const transcript = results
          ? Array.from(results)
              .map((r: SRResult) => r?.[0]?.transcript || '')
              .join(' ')
              .trim()
          : '';
        // Ignore empty or punctuation-only transcripts (e.g., '*')
        if (!transcript || /^[^\p{L}\p{N}]+$/u.test(transcript)) return;
        setInputMessage(transcript);
        if (voiceMode) {
          sendText(transcript);
        }
      };

      recognition.onend = () => {
        setIsListening(false);
      };

  recognition.onerror = () => {
        setIsListening(false);
      };

      recognitionRef.current = recognition;
    }

    // TTS
    if (typeof window !== 'undefined' && 'speechSynthesis' in window) {
      setSupportsTTS(true);
      const computeHasVoice = () => {
        try {
          const voices = window.speechSynthesis.getVoices();
          const exists = voices.some((v) => {
            const lang = (v.lang || '').toLowerCase();
            const name = (v.name || '').toLowerCase();
            return lang === 'am-et' || lang.startsWith('am') || name.includes('amharic') || name.includes('amh');
          });
          setHasAmharicVoice(exists);
        } catch {
          setHasAmharicVoice(false);
        }
      };
      computeHasVoice();
      window.speechSynthesis.onvoiceschanged = computeHasVoice;
    }

    // Load ResponsiveVoice for Amharic cloud fallback
    if (typeof window !== 'undefined') {
      if (!window.responsiveVoice) {
        const script = document.createElement('script');
        // Using public URL; if you have a key, append ?key=YOUR_KEY
        script.src = 'https://code.responsivevoice.org/responsivevoice.js';
        script.async = true;
        script.onload = () => {
          // Poll until voices are ready
          let attempts = 0;
          rvCheckRef.current = window.setInterval(() => {
            attempts += 1;
            try {
              if (window.responsiveVoice && window.responsiveVoice.voiceSupport()) {
                setRvReady(true);
                if (rvCheckRef.current) {
                  clearInterval(rvCheckRef.current);
                  rvCheckRef.current = null;
                }
              }
            } catch {}
            if (attempts > 50 && rvCheckRef.current) { // ~5s
              clearInterval(rvCheckRef.current);
              rvCheckRef.current = null;
            }
          }, 100);
        };
        document.head.appendChild(script);
      } else {
        try {
          if (window.responsiveVoice.voiceSupport()) setRvReady(true);
        } catch {}
      }
    }

    return () => {
      // Cleanup
      if (recognitionRef.current && isListening) {
        try {
          recognitionRef.current.stop();
        } catch {}
      }
      if (utteranceRef.current) {
        window.speechSynthesis.cancel();
      }
      if (rvCheckRef.current) {
        clearInterval(rvCheckRef.current);
        rvCheckRef.current = null;
      }
      // Use captured ref value; don't access audioRef.current directly here
      if (audioElForCleanup && !audioElForCleanup.paused) {
        try { audioElForCleanup.pause(); } catch {}
      }
    };
  }, [isListening, voiceLang, voiceMode, sendText]);

  const toggleListening = () => {
    if (!supportsSTT || !recognitionRef.current) return;
    // Barge-in: cancel TTS if speaking
    if (typeof window !== 'undefined' && window.speechSynthesis?.speaking) {
      window.speechSynthesis.cancel();
      setSpeakingMessageId(null);
      utteranceRef.current = null;
    }
    if (isListening) {
      try { recognitionRef.current.stop(); } catch {}
      setIsListening(false);
    } else {
      try {
        recognitionRef.current.start();
        setIsListening(true);
      } catch {
        setIsListening(false);
      }
    }
  };

  // Helpers for TTS input normalization and cleaning
  const toPlainText = (raw: unknown): string => {
    if (typeof raw === 'string') return raw;
    if (raw && typeof raw === 'object') {
      const obj = raw as { response?: unknown; content?: unknown };
      if (typeof obj.response === 'string') return obj.response;
      if (typeof obj.content === 'string') return obj.content;
    }
    try { return JSON.stringify(raw ?? ''); } catch { return String(raw ?? ''); }
  };
  const cleanForTTS = (s: string): string => s
    .replace(/[`*_#>\-\[\]()`]/g, '')
    .replace(/https?:\/\/\S+/g, '')
    .replace(/\s{2,}/g, ' ')
    .trim();

  const speakMessage = useCallback(async (rawText: unknown, messageId: string) => {
    if (!supportsTTS) return;
    // If already speaking this message, stop
    if (speakingIdRef.current === messageId) {
      try { window.speechSynthesis.cancel(); } catch {}
      if (audioRef.current && !audioRef.current.paused) {
        try { audioRef.current.pause(); } catch {}
      }
      try { window.responsiveVoice?.cancel?.(); } catch {}
      setSpeakingMessageId(null);
      utteranceRef.current = null;
      if (audioUrlRef.current) {
        URL.revokeObjectURL(audioUrlRef.current);
        audioUrlRef.current = null;
      }
      return;
    }
    // Stop current speech
    if (window.speechSynthesis.speaking) {
      try { window.speechSynthesis.cancel(); } catch {}
    }
    if (audioRef.current && !audioRef.current.paused) {
      try { audioRef.current.pause(); } catch {}
      if (audioUrlRef.current) {
        URL.revokeObjectURL(audioUrlRef.current);
        audioUrlRef.current = null;
      }
    }
    try { window.responsiveVoice?.cancel?.(); } catch {}

    // Normalize and clean input text
    const text = toPlainText(rawText);
    const cleaned = cleanForTTS(text);

    // If requesting Amharic, try ResponsiveVoice first, then cloud, then local
    const wantsAmharic = voiceLang.startsWith('am');
    // 1) Use free Google Translate proxy first (no billing)
    if (wantsAmharic && !hasAmharicVoice && audioRef.current) {
      setSpeakingMessageId(messageId);
      fetch('/api/tts', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ text: cleaned || text, lang: 'am' }),
      })
        .then(async (res) => {
          if (!res.ok) throw new Error('TTS request failed');
          const blob = await res.blob();
          if (audioUrlRef.current) {
            URL.revokeObjectURL(audioUrlRef.current);
            audioUrlRef.current = null;
          }
          const url = URL.createObjectURL(blob);
          audioUrlRef.current = url;
          audioRef.current!.src = url;
          try { await audioRef.current!.play(); } catch { setSpeakingMessageId(null); }
          audioRef.current!.onended = () => {
            setSpeakingMessageId(null);
            if (audioUrlRef.current) { URL.revokeObjectURL(audioUrlRef.current); audioUrlRef.current = null; }
          };
        })
        .catch(() => { setSpeakingMessageId(null); });
      return;
    }

    // 2) Try ResponsiveVoice
    if (wantsAmharic && rvReady && window.responsiveVoice) {
      setSpeakingMessageId(messageId);
      try {
        const ok = window.responsiveVoice.speak(cleaned || text, 'Amharic Female', {
          rate: 1.0, pitch: 1.0, onend: () => setSpeakingMessageId(null), onerror: () => setSpeakingMessageId(null)
        });
        if (ok) return;
      } catch { setSpeakingMessageId(null); }
    }

    // 3) Default: use browser speech synthesis
    const utterance = new SpeechSynthesisUtterance(cleaned || text);
    utterance.lang = voiceLang || 'en-US';
    utterance.rate = 1.0;
    utterance.pitch = 1.0;
    const pickVoice = () => {
      const voices = window.speechSynthesis.getVoices();
      const langLower = (utterance.lang || 'en-US').toLowerCase();
      const byExact = voices.find((v) => (v.lang || '').toLowerCase() === 'am-et');
      const byLang = voices.find((v) => (v.lang || '').toLowerCase().startsWith('am'));
      const byName = voices.find((v) => /amharic|amh/gi.test(v.name || ''));
      const byRegion = voices.find((v) => (v.lang || '').toLowerCase().endsWith('-et'));
      const byPrefix = voices.find((v) => (v.lang || '').toLowerCase().startsWith(langLower.slice(0, 2)));
      const preferred = byExact || byLang || byName || byRegion || byPrefix || null;
      if (preferred) utterance.voice = preferred;
    };
    if (window.speechSynthesis.getVoices().length === 0) {
      window.speechSynthesis.onvoiceschanged = pickVoice;
    } else {
      pickVoice();
    }
    utterance.onend = () => {
      setSpeakingMessageId(null);
      utteranceRef.current = null;
    };
    utterance.onerror = () => {
      setSpeakingMessageId(null);
      utteranceRef.current = null;
    };
    utteranceRef.current = utterance;
    setSpeakingMessageId(messageId);
    window.speechSynthesis.speak(utterance);
  }, [supportsTTS, voiceLang, hasAmharicVoice, rvReady]);

  // Auto-speak assistant replies when in voice mode
  useEffect(() => {
    if (!voiceMode || !supportsTTS || messages.length === 0) return;
    const last = messages[messages.length - 1];
    if (last.type === 'assistant' && last.content) {
      if (lastSpokenIdRef.current === last.id) return; // avoid re-speaking same message
      // If listening, stop to avoid overlap
      if (recognitionRef.current && isListening) {
        try { recognitionRef.current.stop(); } catch {}
        setIsListening(false);
      }
      speakMessage(last.content, last.id);
      lastSpokenIdRef.current = last.id;
    }
  }, [messages, voiceMode, supportsTTS, isListening]);

  const handleSendMessage = () => {
    if (inputMessage.trim()) {
      const token = session?.accessToken;
      const newMessage: Message = {
        id: Date.now().toString(),
        type: "user",
        content: inputMessage,
        timestamp: new Date().toLocaleString(),
      };
      dispatch(addUserMessage(newMessage));
      if (token) {
        dispatch(sendMessage({ query: inputMessage, token })).then((result) => {
          // log the raw thunk result and payload
          console.log("Chat sendMessage result:", result);
          if (result?.meta?.requestStatus === "fulfilled") {
            console.log("Chat API payload:", result.payload);
            setSuccessMessage("Message sent successfully!");
            setTimeout(() => setSuccessMessage(""), 3000);
          } else {
            console.error("Chat API error result:", result);
          }
        });
      }
      setInputMessage("");
    }
  };

  const translateMessage = async (message: Message) => {
    const token = session?.accessToken;
    const rawContent: unknown = message.content as unknown;
    const sourceText = typeof rawContent === 'string'
      ? rawContent
      : (typeof rawContent === 'object' && rawContent && (rawContent as { response?: string }).response
        ? (rawContent as { response?: string }).response!
        : JSON.stringify(rawContent));
    const target = isAmharic(sourceText) ? 'en' : 'am';
    if (!token) return;
    setTranslating((prev) => ({ ...prev, [message.id]: true }));
    try {
      // Sanitize and limit payload to avoid backend rejection
      const cleaned = (sourceText || '')
        .replace(/[`*_#>\-\[\]()]/g, '')
        .replace(/https?:\/\/\S+/g, '')
        .replace(/\s{2,}/g, ' ')
        .trim()
        .slice(0, 4000);

      const payload = {
        content: {
          response: cleaned,
          procedures: (message.procedures || []).map((p) => ({ id: String(p.id), name: p.title })),
        },
      };

      const res = await fetch(`${API_BASE_URL}/ai/translate`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json',
          Authorization: `Bearer ${token}`,
          lang: target,
        },
        body: JSON.stringify(payload),
      });

      if (!res.ok) {
        const errText = await res.text().catch(() => '');
        throw new Error(errText || 'Translate failed');
      }

      const translated = await res.json();
      const translatedText = typeof translated === 'string'
        ? translated
        : (translated?.response ?? translated?.content ?? (typeof translated?.content === 'object' ? translated?.content?.response : ''));

      if (translatedText) {
        const newMsg: Message = {
          id: `${message.id}-translated-${Date.now()}`,
          type: 'assistant',
          content: translatedText,
          timestamp: new Date().toLocaleString(),
        } as Message;
        dispatch(addUserMessage(newMsg));
      }
    } catch (e) {
      console.error(e);
      setSuccessMessage('Translation failed');
      setTimeout(() => setSuccessMessage(''), 3000);
    } finally {
      setTranslating((prev) => ({ ...prev, [message.id]: false }));
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSendMessage();
    }
  };

  const handleChatSelect = (chatId: string) => {
    const token = session?.accessToken;
    if (token) {
      dispatch(fetchChatById({ id: chatId, token })).then((result) => {
        if (result.meta.requestStatus === "fulfilled") {
          setSuccessMessage("Chat history loaded successfully!");
          setTimeout(() => setSuccessMessage(""), 3000);
        }
      });
    }
  };

  // show intro until user interacts (no messages yet)
  const isEmpty = messages.length === 0;

  // quick suggestions for first-time users
  const suggestions = [
    "What documents do I need for a business license?",
    "Help me start the tourist visa application.",
    "Show me the steps to renew a residence permit.",
    "Translate the requirements into Amharic.",
  ];
  const handleUseSuggestion = (text: string) => setInputMessage(text);

  // parser: extract Procedure, Required Documents, Steps from assistant text
  const parseGuide = (text: string) => {
    const lines = (text || "").split(/\r?\n/).map(l => l.trim());
    let procedure = "";
    const documents: string[] = [];
    const steps: string[] = [];
    let inDocs = false;
    let inSteps = false;

    for (const raw of lines) {
      const line = raw.replace(/\s+$/g, "");
      if (!line) continue;

      if (/^procedure\s*:/i.test(line)) {
        procedure = line.split(/:/, 2)[1]?.trim() || "";
        inDocs = false; inSteps = false;
        continue;
      }
      if (/^required documents\s*:?/i.test(line)) {
        inDocs = true; inSteps = false;
        continue;
      }
      if (/^steps\s*:?/i.test(line)) {
        inSteps = true; inDocs = false;
        continue;
      }

      // bullets and numbered lines
      const isBullet = /^[-•]\s+/.test(line);
      const isNum = /^\d+[\.\)]\s+/.test(line);

      if (inDocs && (isBullet || isNum)) {
        documents.push(line.replace(/^[-•]\s+/, "").replace(/^\d+[\.\)]\s+/, "").trim());
        continue;
      }
      if (inSteps && (isBullet || isNum)) {
        steps.push(line.replace(/^[-•]\s+/, "").replace(/^\d+[\.\)]\s+/, "").trim());
        continue;
      }
    }

    return {
      hasStructured: Boolean(procedure || documents.length || steps.length),
      procedure,
      documents,
      steps,
    };
  };

  if (status === "loading") {
    return <div>Loading...</div>;
  }

  if (status === "unauthenticated") {
    return <div>Please sign in to access the chat.</div>;
  }

  return (
    <div className="min-h-screen bg-gray-50 flex">
      <div
        className={`bg-white border-r border-gray-200 transition-all duration-300 ${showHistory ? "w-80" : "w-0"} overflow-hidden`}
      >
        <div className="p-4 border-b border-gray-200">
          <div className="flex items-center space-x-2">
            <History className="w-5 h-5 text-[#3A6A8D]" />
            <h2 className="font-semibold text-gray-900">Chat History</h2>
          </div>
        </div>
        <div className="p-4 space-y-3 max-h-[calc(100vh-80px)] overflow-y-auto">
          {chatHistory.map((chat) => (
            <Card
              key={chat.id}
              className="cursor-pointer hover:shadow-md transition-shadow duration-200 border border-gray-200"
              onClick={() => handleChatSelect(chat.id)}
            >
              <CardContent className="p-3">
                <h3 className="font-medium text-gray-900 text-sm mb-1 line-clamp-1">{chat.title}</h3>
                <p className="text-xs text-gray-600 mb-2 line-clamp-2">{chat.lastMessage}</p>
                <div className="flex items-center justify-between text-xs text-gray-500">
                  <div className="flex items-center space-x-1">
                    <Clock className="w-3 h-3" />
                    <span>{chat.timestamp}</span>
                  </div>
                  <span>{chat.messageCount} messages</span>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
      {/* Main Chat Area */}
      <div className="flex-1 flex flex-col">
        {/* Header */}
        <div className="bg-gray-50 border-b border-gray-50 p-6">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-extrabold bg-clip-text text-transparent bg-gradient-to-r from-[#3A6A8D] to-[#6aa5d8] tracking-tight">Your AI Guide</h1>
              <p className="text-sm text-gray-600 mt-1">Ask anything about procedures, fees, and steps — in English or Amharic.</p>
            </div>
            <div className="flex items-center gap-2">
            <Button
              variant="outline"
              onClick={() => setShowHistory(!showHistory)}
                className="border-gray-200 text-[#3A6A8D] hover:text-white bg-white hover:bg-[#3A6A8D] transition-colors rounded-full px-4"
                title={showHistory ? 'Hide previous conversations' : 'Show previous conversations'}
            >
              <History className="w-4 h-4 mr-2" />
                {showHistory ? "Hide History" : "History"}
              </Button>
              <Button
                variant="outline"
                onClick={() => {
                  dispatch(clearMessages());
                  setInputMessage("");
                  setSpeakingMessageId(null);
                  if (typeof window !== 'undefined' && window.speechSynthesis?.speaking) {
                    try { window.speechSynthesis.cancel(); } catch {}
                  }
                  if (audioRef.current && !audioRef.current.paused) {
                    try { audioRef.current.pause(); } catch {}
                  }
                  setSuccessMessage('Started a new chat');
                  setTimeout(() => setSuccessMessage(''), 2000);
                }}
                className="border-[#3A6A8D] text-white bg-[#3A6A8D] hover:bg-[#2d5470] rounded-full px-4"
                title="Start a new chat"
              >
                <Plus className="w-4 h-4 mr-2" /> New Chat
            </Button>
            </div>
          </div>
        </div>
        {/* Chat Messages */}
        <div className="flex-1 overflow-y-auto p-6 space-y-6">
          {/* Subtle intro, always visible */}
          <div className="bg-white/70 backdrop-blur rounded-xl border border-gray-200 p-4 md:p-5 shadow-sm">
            <div className="flex items-start gap-3">
              <div className="w-10 h-10 rounded-full bg-[#3A6A8D] flex items-center justify-center shrink-0">
                <Bot className="w-5 h-5 text-white" />
              </div>
              <div className="flex-1">
                <h2 className="text-sm md:text-base font-semibold text-gray-900">Welcome</h2>
                <p className="text-xs md:text-sm text-gray-600 mt-1">Type or speak your question. Use Auto for hands-free voice. Toggle AM/EN for your preferred language. Click Translate to switch the answer language.</p>
              </div>
            </div>
          </div>
          {/* {chatStatus === 'loading' && <p className="text-gray-500">Loading messages...</p>}
          {error && (
            <p className="text-red-500">
              {error}
              <Button variant="ghost" onClick={() => dispatch(clearError())} className="ml-2 text-sm">Clear</Button>
            </p>
          )} */}
          {successMessage && <p className="text-green-500">{successMessage}</p>}
          {messages.map((message) => (
            <div key={message.id} className="animate-fade-in">
              {message.type === "assistant" ? (
                <div className="flex items-start space-x-3">
                  <div className="flex-shrink-0">
                    <div className="w-8 h-8 bg-[#3A6A8D] rounded-full flex items-center justify-center">
                      <Bot className="w-4 h-4 text-white" />
                    </div>
                  </div>
                  <div className="flex-1 space-y-4">
                    <div className="bg-white rounded-lg p-4 shadow-sm border border-gray-200">
                      <div className="text-gray-900 prose prose-sm md:prose-base lg:prose-lg max-w-none leading-relaxed">
                        {(() => {
                          const raw: unknown = message.content as unknown;
                          const text = typeof raw === 'string' ? raw : (typeof raw === 'object' && raw && (raw as { response?: string }).response ? (raw as { response?: string }).response! : JSON.stringify(raw));
                          return <ReactMarkdown>{String(text)}</ReactMarkdown>;
                        })()}
                      </div>
                      <div className="flex items-center justify-between mt-4">
                        <span className="text-xs text-gray-500">{message.timestamp}</span>
                        <Badge variant="secondary" className="bg-green-100 text-green-700">
                          <CheckCircle className="w-3 h-3 mr-1" />
                          Verified
                        </Badge>
                      </div>
                      {supportsTTS && (
                        <div className="mt-3 flex justify-end">
                          <Button
                            variant="outline"
                            size="sm"
                            className="border-gray-300 bg-transparent hover:bg-blue-100 hover:text-blue-700 text-xs py-1 px-2"
                            onClick={() => speakMessage(message.content, message.id)}
                          >
                            {speakingMessageId === message.id ? (
                              <>
                                <Square className="w-3 h-3 mr-1" /> Stop
                              </>
                            ) : (
                              <>
                                <Play className="w-3 h-3 mr-1" /> Play
                              </>
                            )}
                          </Button>
                        </div>
                      )}
                    </div>
                    <Button
                      variant="outline"
                      className="border-gray-300 bg-transparent hover:bg-blue-100 hover:text-blue-700 text-xs py-1 px-2"
                      onClick={() => translateMessage(message)}
                      disabled={!!translating[message.id]}
                    >
                      <Languages className="w-3 h-3 mr-1" />
                      {(() => {
                        if (translating[message.id]) return 'Translating…';
                        const raw: unknown = message.content as unknown;
                        const text = typeof raw === 'string' ? raw : (typeof raw === 'object' && raw && (raw as { response?: string }).response ? (raw as { response?: string }).response! : JSON.stringify(raw));
                        return isAmharic(text) ? 'Translate to EN' : 'Translate to AM';
                      })()}
                    </Button>
                    {/* Procedures */}
                    {message.procedures && (
                      <div className="space-y-2">
                        {message.procedures.map((procedure) => {
                          const IconComponent = { FileText, DollarSign, Building }[procedure.icon] || FileText;
                          return (
                            <Card
                              key={procedure.id}
                              className="bg-white border-2 border-transparent bg-gradient-to-r from-blue-50 to-indigo-50 rounded-md shadow-xs hover:shadow-sm hover:scale-102 transition-all duration-200 animate-in fade-in"
                            >
                              <CardContent className="p-2">
                                <div className="flex items-center space-x-2 mb-1.5">
                                  <div className="w-5 h-5 bg-indigo-100 rounded-full flex items-center justify-center transform hover:scale-110 transition-transform duration-150">
                                    <IconComponent className="w-2.5 h-2.5 text-indigo-600" />
                                  </div>
                                  <h3 className="font-medium text-gray-900 text-xs font-sans">
                                    Procedure {procedure.id}: {procedure.title}
                                  </h3>
                                </div>
                                <ul className="space-y-0.5 ml-7">
                                  {procedure.items.length > 0 ? (
                                    procedure.items.map((item, index) => (
                                      <li key={index} className="text-gray-700 text-[0.65rem] font-sans flex items-start">
                                        <span className="w-0.75 h-0.75 bg-indigo-400 rounded-full mt-1 mr-1.5 flex-shrink-0"></span>
                                        {item}
                                      </li>
                                    ))
                                  ) : (
                                    <li className="text-gray-500 text-[0.65rem] font-sans italic">No details available</li>
                                  )}
                                </ul>
                                <div className="flex flex-wrap gap-1.5 pt-2">
                                  <Button
                                    className="bg-[#3A6A8D] hover:bg-[#2d5470] text-white text-[0.65rem] font-sans py-0.5 px-1.5 rounded-md transform hover:scale-105 transition-transform duration-150"
                                    onClick={async () => {
                                      try {
                                        await createChecklist({ procedureId: String(procedure.id), token: session?.accessToken || undefined }).unwrap()
                                      } catch {
                                        // ignore error; still navigate to workspace to show current state
                                      } finally {
                                        router.push('/user/workspace')
                                      }
                                    }}
                                  >
                                    <Bookmark className="w-2.5 h-2.5 mr-1" />
                                    Save Checklist
                                  </Button>
                                  <Button
                                    variant="outline"
                                    className="border-indigo-300 bg-transparent hover:bg-indigo-100 hover:text-indigo-700 text-[0.65rem] font-sans py-0.5 px-1.5 rounded-md transform hover:scale-105 transition-transform duration-150"
                                  >
                                    <Play className="w-2.5 h-2.5 mr-1" />
                                    Procedure
                                  </Button>
                                </div>
                              </CardContent>
                            </Card>
                          );
                        })}
                      </div>
                    )}
                  </div>
                </div>
              ) : (
                <div className="flex items-start space-x-3 justify-end">
                  <div className="bg-[#3A6A8D] text-white rounded-lg p-4 max-w-md shadow-sm">
                    <p className="text-sm font-sans">{message.content}</p>
                    <span className="text-xs text-gray-200 mt-2 block">{message.timestamp}</span>
                  </div>
                  <div className="flex-shrink-0">
                    <div className="w-8 h-8 bg-gray-300 rounded-full flex items-center justify-center">
                      <User className="w-4 h-4 text-gray-600" />
                    </div>
                  </div>
                </div>
              )}
            </div>
          ))}
          {chatStatus === 'loading' && <p className="text-gray-500">Loading messages...</p>}
          {error && (
            <p className="text-red-500">
              {error}
              <Button variant="ghost" onClick={() => dispatch(clearError())} className="ml-2 text-sm">Clear</Button>
            </p>
          )}
        </div>
        {/* Input Area */}
        <div className="bg-gray-50 p-4">
          <div className="flex items-center space-x-3 max-w-4xl mx-auto">
            <Button
              variant="ghost"
              size="sm"
              className={`p-2 h-10 w-10 rounded-full ${isListening ? 'bg-red-100 hover:bg-red-200' : 'hover:bg-gray-200'}`}
              onClick={toggleListening}
              disabled={!supportsSTT}
              title={supportsSTT ? (isListening ? 'Stop listening' : 'Start voice input') : 'Voice input not supported'}
            >
              <Mic className={`w-5 h-5 ${isListening ? 'text-red-600' : 'text-gray-500'}`} />
            </Button>
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setVoiceMode((v) => !v)}
              className={`px-3 h-10 rounded-full ${voiceMode ? 'bg-green-100 hover:bg-green-200 text-green-700' : 'hover:bg-gray-200 text-gray-600'}`}
              title="Toggle voice mode (auto send and auto speak)"
            >
              {voiceMode ? 'Auto' : 'Manual'}
            </Button>
            <Button
              variant="ghost"
              size="sm"
              onClick={() => {
                const next = voiceLang.startsWith('am') ? 'en-US' : 'am-ET';
                setVoiceLang(next);
                if (recognitionRef.current) {
                  try { recognitionRef.current.lang = next; } catch {}
                }
              }}
              className="px-3 h-10 rounded-full hover:bg-gray-200 text-gray-700"
              title="Toggle voice language between Amharic and English"
            >
              {voiceLang.startsWith('am') ? 'AM' : 'EN'}
            </Button>
            <div className="flex-1 bg-white rounded-full px-4 py-3 shadow-sm border border-[#3A6A8D] focus-within:ring-2 focus-within:ring-[#3A6A8D] focus-within:border-transparent">
              <input
                type="text"
                value={inputMessage}
                onChange={(e) => setInputMessage(e.target.value)}
                onKeyPress={handleKeyPress}
                placeholder="Type your question here..."
                className="w-full bg-transparent border-none outline-none text-gray-700 placeholder-gray-400 text-sm font-sans"
              />
            </div>
            <Button
              onClick={handleSendMessage}
              className="bg-[#3A6A8D] hover:bg-[#2d5470] text-white rounded-full p-2 h-10 w-10 flex items-center justify-center"
              disabled={chatStatus === 'loading'}
            >
              <Send className="w-4 h-4" />
            </Button>
            {/* Hidden audio element for cloud TTS playback */}
            <audio ref={audioRef} hidden />
          </div>
        </div>
      </div>
    </div>
  );
}