"use client"

import type React from "react"

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
  MapPin,
  Send,
  Bookmark,
  Play,
  Languages,
  CheckCircle,
  Mic,
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

import { useRouter } from "next/navigation"

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

  const router = useRouter();

  // Sidebar menu items with navigation
  // Removed unused menuItems

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col">
      <div className="flex-1 flex flex-col">
        {/* Chat Messages */}
        <div className="flex-1 overflow-y-auto p-6 space-y-6">
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
                            <Card key={step.id} className="border border-gray-200 hover:shadow-md transition-shadow duration-200">
                              <CardContent className="p-4">
                                <div className="flex items-center space-x-3 mb-3">
                                  <div className="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center">
                                    <IconComponent className="w-4 h-4 text-blue-600" />
                                  </div>
                            <h3 className="font-semibold text-gray-900">Step {step.id}: {step.title}</h3>
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
                        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                          <Card className="border border-gray-200 hover:shadow-lg hover:-translate-y-1 transition-all duration-300 ease-in-out animate-fade-in" style={{ animationDelay: "0.1s" }}>
                            <CardContent className="p-4">
                              <div className="flex items-center space-x-2 mb-2">
                                <FileText className="w-4 h-4 text-blue-600 transition-transform duration-200 hover:scale-110" />
                                <h4 className="font-medium text-gray-900">Required Documents</h4>
                              </div>
                        <p className="text-sm text-gray-600">You&apos;ll need documents for proof of identity and business certificate.</p>
                            </CardContent>
                          </Card>
                          <Card className="border border-gray-200 hover:shadow-lg hover:-translate-y-1 transition-all duration-300 ease-in-out animate-fade-in" style={{ animationDelay: "0.2s" }}>
                            <CardContent className="p-4">
                              <div className="flex items-center space-x-2 mb-2">
                                <DollarSign className="w-4 h-4 text-green-600 transition-transform duration-200 hover:scale-110" />
                                <h4 className="font-medium text-gray-900">Processing Fee</h4>
                              </div>
                              <p className="text-sm text-gray-600">The application fee is 350 ETB, payable at the time of submission.</p>
                            </CardContent>
                          </Card>
                          <Card className="border border-gray-200 hover:shadow-lg hover:-translate-y-1 transition-all duration-300 ease-in-out animate-fade-in" style={{ animationDelay: "0.3s" }}>
                            <CardContent className="p-4">
                              <div className="flex items-center space-x-2 mb-2">
                                <MapPin className="w-4 h-4 text-orange-600 transition-transform duration-200 hover:scale-110" />
                                <h4 className="font-medium text-gray-900">Office Location</h4>
                              </div>
                              <p className="text-sm text-gray-600">Visit your local federal administrative office during business hours (8:00 AM - 5:00 PM).</p>
                            </CardContent>
                          </Card>
                        </div>
                        <div className="flex flex-wrap gap-3 pt-4">
                          <Button className="bg-[#3A6A8D] hover:bg-[#2d5470] text-white">
                            <Bookmark className="w-4 h-4 mr-2" />
                            Save Checklist
                          </Button>
                          <Button variant="outline" className="border-gray-300 bg-transparent hover:bg-blue-100 hover:text-blue-900">
                            <Play className="w-4 h-4 mr-2" />
                            Start Procedure
                          </Button>
                          <Button variant="outline" className="border-gray-300 bg-transparent hover:bg-blue-100 hover:text-blue-900">
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
                  <div className="bg-[#3A6A8D] text-white rounded-lg p-4 max-w-md shadow-sm">
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
        <div className="bg-gray-50 p-4">
          <div className="flex items-center space-x-3 max-w-4xl mx-auto">
            <Button variant="ghost" size="sm" className="p-2 h-10 w-10 rounded-full hover:bg-gray-200">
              <Mic className="w-5 h-5 text-gray-500" />
            </Button>
            <div className="flex-1 bg-white rounded-full px-4 py-3 shadow-sm border border-gray-200">
              <input
                type="text"
                value={inputMessage}
                onChange={(e) => setInputMessage(e.target.value)}
                onKeyDown={handleKeyPress}
                placeholder="Type your question here..."
                className="w-full bg-transparent border-none outline-none text-gray-700 placeholder-gray-400"
              />
            </div>
            <Button
              onClick={handleSendMessage}
              className="bg-[#3A6A8D] hover:bg-[#2d5470] text-white rounded-full p-2 h-10 w-10 flex items-center justify-center"
            >
              <Send className="w-4 h-4" />
            </Button>
          </div>
        </div>
      </div>
    </div>
  )
}
