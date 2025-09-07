"use client";

import type React from "react";
import { useRouter } from "next/navigation";
import { useState, useEffect } from "react";
import { useSelector, useDispatch } from 'react-redux';
import { useSession } from "next-auth/react";
import { useCreateChecklistMutation } from '@/app/store/slices/checklistsApi'
import { RootState, AppDispatch } from '@/app/store/store';
import { fetchChatHistory, sendMessage, addUserMessage, clearError, fetchChatById } from '@/app/store/slices/aiChatSlice';
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
} from "lucide-react";

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

  useEffect(() => {
    const token = session?.accessToken;
    if (status === 'authenticated' && token) {
      dispatch(fetchChatHistory(token));
    }
  }, [dispatch, session, status]);

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
          if (result.meta.requestStatus === "fulfilled") {
            setSuccessMessage("Message sent successfully!");
            setTimeout(() => setSuccessMessage(""), 3000);
          }
        });
      }
      setInputMessage("");
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
              <h1 className="text-2xl font-bold text-gray-900 mb-1">Chat with Your AI Guide</h1>
              <p className="text-gray-600">Your Guide, Your Chat</p>
            </div>
            <Button
              variant="outline"
              onClick={() => setShowHistory(!showHistory)}
              className="border-gray-300 text-white hover:text-white bg-[#3A6A8D] hover:bg-[#2d5470]"
            >
              <History className="w-4 h-4 mr-2" />
              {showHistory ? "Hide History" : "Show History"}
            </Button>
          </div>
        </div>
        {/* Chat Messages */}
        <div className="flex-1 overflow-y-auto p-6 space-y-6">
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
                      <div className="text-gray-800 prose">
                        <ReactMarkdown>{message.content}</ReactMarkdown>
                      </div>
                      <div className="flex items-center justify-between mt-3">
                        <span className="text-xs text-gray-500">{message.timestamp}</span>
                        <Badge variant="secondary" className="bg-green-100 text-green-700">
                          <CheckCircle className="w-3 h-3 mr-1" />
                          Verified
                        </Badge>
                      </div>
                    </div>
                    <Button
                      variant="outline"
                      className="border-gray-300 bg-transparent hover:bg-blue-100 hover:text-blue-700 text-xs py-1 px-2"
                    >
                      <Languages className="w-3 h-3 mr-1" />
                      Translate
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
            <Button variant="ghost" size="sm" className="p-2 h-10 w-10 rounded-full hover:bg-gray-200">
              <Mic className="w-5 h-5 text-gray-500" />
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
          </div>
        </div>
      </div>
    </div>
  );
}