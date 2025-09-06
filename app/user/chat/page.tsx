"use client"

import type React from "react"
import { useRouter } from "next/navigation"
import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
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
} from "lucide-react"

interface Message {
  id: string
  type: "user" | "assistant"
  content: string
  timestamp: string
  steps?: Step[]
}

interface Step {
  id: number
  title: string
  icon: React.ElementType
  items: string[]
  completed?: boolean
}

interface ChatHistory {
  id: string
  title: string
  lastMessage: string
  timestamp: string
  messageCount: number
}

export default function ChatPage() {
  const [messages, setMessages] = useState<Message[]>([
    {
      id: "1",
      type: "assistant",
      content:
        "Hello! I'm your AI legal assistant. I can help you navigate Ethiopian legal procedures, business registration, and more. What would you like to know?",
      timestamp: "Dec 04 at 2:24 PM",
    },
    {
      id: "2",
      type: "user",
      content: "How do I register my business in Ethiopia?",
      timestamp: "Dec 04 at 2:25 PM",
    },
    {
      id: "3",
      type: "assistant",
      content: "Here's a comprehensive guide to register your business in Ethiopia:",
      timestamp: "Dec 04 at 2:25 PM",
      steps: [
        {
          id: 1,
          title: "Prepare Required Documents",
          icon: FileText,
          items: [
            "Valid ID or passport copy",
            "Business name reservation certificate",
            "Memorandum of association",
            "Articles of association",
          ],
        },
        {
          id: 2,
          title: "Pay Registration Fees",
          icon: DollarSign,
          items: ["Registration fee: 200 ETB", "Stamp duty: 50 ETB", "Certificate fee: 100 ETB"],
        },
        {
          id: 3,
          title: "Visit Registration Office",
          icon: Building,
          items: [
            "Submit documents at Regional Trade Office",
            "Wait for verification (5-7 business days)",
            "Collect business license certificate",
          ],
        },
      ],
    },
  ])

  const [inputMessage, setInputMessage] = useState("")

  const [chatHistory] = useState<ChatHistory[]>([
    {
      id: "1",
      title: "Business Registration Process",
      lastMessage: "Here's a comprehensive guide to register your business...",
      timestamp: "Dec 04 at 2:25 PM",
      messageCount: 8,
    },
    {
      id: "2",
      title: "Passport Application Requirements",
      lastMessage: "You'll need to prepare the following documents...",
      timestamp: "Dec 03 at 4:15 PM",
      messageCount: 12,
    },
    {
      id: "3",
      title: "Tax Registration Procedures",
      lastMessage: "The tax registration process involves several steps...",
      timestamp: "Dec 02 at 10:30 AM",
      messageCount: 6,
    },
    {
      id: "4",
      title: "Work Permit Application",
      lastMessage: "For foreign nationals seeking employment...",
      timestamp: "Dec 01 at 3:45 PM",
      messageCount: 15,
    },
  ])

  const [showHistory, setShowHistory] = useState(false)

  const handleSendMessage = () => {
    if (inputMessage.trim()) {
      const newMessage: Message = {
        id: Date.now().toString(),
        type: "user",
        content: inputMessage,
        timestamp: new Date().toLocaleString(),
      }
      setMessages([...messages, newMessage])
      setInputMessage("")
    }
  }

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault()
      handleSendMessage()
    }
  }
  const router = useRouter()

  return (
    <div className="min-h-screen bg-gray-50 flex relative">
      {/* animations (scoped) */}
      <style jsx>{`
        .fade-in { animation: fadeIn 360ms ease-out both; }
        .fade-in-up { animation: fadeInUp 520ms ease-out both; }
        .slide-left { animation: slideLeft 360ms ease-out both; }
        .pop-in { animation: popIn 340ms ease-out both; }
        .btn-pop { transition: transform .2s ease; }
        .btn-pop:hover { transform: translateY(-1px) scale(1.01); }
        .btn-pop:active { transform: scale(0.98); }
        .ring-hover { transition: box-shadow .25s ease; }
        .ring-hover:hover { box-shadow: 0 0 0 6px rgba(58,106,141,.15); }
        @keyframes fadeIn { from { opacity: 0 } to { opacity: 1 } }
        @keyframes fadeInUp { from { opacity: 0; transform: translateY(10px) } to { opacity: 1; transform: translateY(0) } }
        @keyframes slideLeft { from { opacity: 0; transform: translateX(-16px) } to { opacity: 1; transform: translateX(0) } }
        @keyframes popIn { 0% { opacity: 0; transform: scale(.98) } 100% { opacity: 1; transform: scale(1) } }
      `}</style>

      {/* Mobile overlay (only when open on small screens) */}
      {showHistory && (
        <div
          className="sm:hidden fixed inset-0 z-40 bg-black/40 backdrop-blur-sm"
          onClick={() => setShowHistory(false)}
        />
      )}

      {/* Sidebar / History Panel */}
      <div
        className={`bg-white border-r border-gray-200 transition-all duration-300 ease-in-out
        ${showHistory ? "w-72 sm:w-80 slide-left" : "w-0"} overflow-hidden
        fixed sm:relative top-0 bottom-0 left-0 z-50 flex flex-col`}
        aria-hidden={!showHistory}
      >
        {/* hide tab order when collapsed */}
        {showHistory && (
          <>
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
          </>
        )}
      </div>

      {/* Main Chat Area */}
      <div className={`flex-1 flex flex-col transition-all duration-300 ${showHistory ? "sm:ml-0" : ""}`}>
        {/* Header */}
        <div className="bg-white border-b border-gray-200 p-4 sm:p-6 sticky top-0 z-30 flex items-center justify-between fade-in-up" style={{ animationDelay: "40ms" }}>
          <div>
            <h1 className="text-lg sm:text-2xl font-bold text-gray-900 mb-1">Chat with Your AI Guide</h1>
            <p className="text-gray-600 text-xs sm:text-sm">Your Guide, Your Chat</p>
          </div>
          <Button
            variant="outline"
            aria-expanded={showHistory}
            onClick={() => setShowHistory((p) => !p)}
            className="border-gray-300 text-white bg-[#3A6A8D] hover:bg-[#2d5470] px-3 py-2 h-auto text-xs sm:text-sm btn-pop"
          >
            <History className="w-4 h-4 mr-1" />
            {showHistory ? "Hide" : "History"}
          </Button>
        </div>

        {/* Messages */}
        <div className="flex-1 overflow-y-auto p-4 sm:p-6 space-y-6">
          {messages.map((message, idx) => (
            <div key={message.id} className="fade-in-up" style={{ animationDelay: `${idx * 70 + 80}ms` }}>
              {message.type === "assistant" ? (
                <div className="flex items-start space-x-3">
                  <div className="flex-shrink-0">
                    <div className="w-8 h-8 bg-[#3A6A8D] rounded-full flex items-center justify-center">
                      <Bot className="w-4 h-4 text-white" />
                    </div>
                  </div>
                  <div className="flex-1 space-y-4">
                    <div className="bg-white rounded-lg p-4 shadow-sm border border-gray-200 pop-in">
                      <p className="text-gray-800">{message.content}</p>
                      <div className="flex items-center justify-between mt-3">
                        <span className="text-xs text-gray-500">{message.timestamp}</span>
                        <Badge variant="secondary" className="bg-green-100 text-green-700">
                          <CheckCircle className="w-3 h-3 mr-1" />
                          Verified
                        </Badge>
                      </div>
                    </div>
                    {message.steps && (
                      <div className="space-y-4">
                        {message.steps.map((step) => {
                          const IconComponent = step.icon
                          return (
                            <Card
                              key={step.id}
                              className="border border-white bg-white hover:shadow-md transition-shadow duration-200"
                            >
                              <CardContent className="p-4">
                                <div className="flex items-center space-x-3 mb-3">
                                  <div className="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center">
                                    <IconComponent className="w-4 h-4 text-blue-600" />
                                  </div>
                                  <h3 className="font-semibold text-gray-900">
                                    Step {step.id}: {step.title}
                                  </h3>
                                </div>
                                <ul className="space-y-2 ml-11">
                                  {step.items.map((item, index) => (
                                    <li key={index} className="text-gray-700 text-sm flex items-start">
                                      <span className="w-1.5 h-1.5 bg-gray-400 rounded-full mt-2 mr-3 flex-shrink-0"></span>
                                      {item}
                                    </li>
                                  ))}
                                </ul>
                              </CardContent>
                            </Card>
                          )
                        })}
                        <div className="flex flex-wrap gap-3 pt-4">
                          <Button className="bg-[#3A6A8D] hover:bg-[#2d5470] text-white" onClick={() => router.push("./workspace")}>
                            <Bookmark className="w-4 h-4 mr-2" />
                            Save Checklist
                          </Button>
                          <Button variant="outline" className="border-gray-300 bg-transparent hover:bg-blue-100 hover:text-blue-700">
                            <Languages className="w-4 h-4 mr-2" />
                            Translate
                          </Button>
                        </div>
                      </div>
                    )}
                  </div>
                </div>
              ) : (
                <div className="flex items-start space-x-3 justify-end">
                  <div className="bg-[#3A6A8D] text-white rounded-lg p-4 max-w-md shadow-sm pop-in">
                    <p>{message.content}</p>
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
        </div>

        {/* Input Area */}
        <div className="bg-gray-50 p-3 sm:p-4">
          <div className="flex items-center gap-2 sm:gap-3 max-w-4xl mx-auto">
            <Button variant="ghost" size="sm" className="p-2 h-10 w-10 rounded-full hover:bg-gray-200 ring-hover">
              <Mic className="w-5 h-5 text-gray-500" />
            </Button>

            <div className="flex-1 bg-white rounded-full px-4 py-3 shadow-sm border border-[#3A6A8D] focus-within:ring-2 focus-within:ring-[#3A6A8D] focus-within:border-transparent">
              <input
                type="text"
                value={inputMessage}
                onChange={(e) => setInputMessage(e.target.value)}
                onKeyPress={handleKeyPress}
                placeholder="Type your question here..."
                className="w-full bg-transparent border-none outline-none text-gray-700 placeholder-gray-400"
              />
            </div>

            <Button
              onClick={handleSendMessage}
              className="bg-[#3A6A8D] hover:bg-[#2d5470] text-white rounded-full p-2 h-10 w-10 flex items-center justify-center btn-pop"
            >
              <Send className="w-4 h-4" />
            </Button>
          </div>
        </div>
      </div>
    </div>
  )
}