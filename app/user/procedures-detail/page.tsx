"use client"

import { useState } from "react"
import {
  Clock,
  DollarSign,
  Download,
  ThumbsUp,
  FileText,
  Globe,
  MessageCircle,
  Share,
} from "lucide-react"
import Image from "next/image"
import Link from "next/link"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Badge } from "@/components/ui/badge"

export default function ProcedureDetail() {
  const [activeTab, setActiveTab] = useState<"feedback" | "notices" | "discussions">("feedback")

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="flex">
        <main className="flex-1 p-6 pl-20">
          <div className="max-w-4xl space-y-6">
            <div className="bg-white rounded-lg p-6 shadow-sm hover:shadow-md transition-all duration-300 transform hover:-translate-y-1">
              {/* Header Section */}
              <div className="mb-6">
                <h1 className="text-2xl font-bold text-[#111827] mb-1 animate-in fade-in duration-500">
                  Driver&#39;s License Renewal
                </h1>
                <p className="text-[#6b7280] text-sm animate-in fade-in duration-700">
                  Complete guide to renew your driver&#39;s license in Ethiopia
                </p>


              <div className="grid grid-cols-2 gap-6 mb-6">
                <div className="bg-white rounded-lg p-6 text-center border border-[#e5e7eb] hover:border-[#4A90E2] transition-all duration-300 hover:shadow-lg transform hover:-translate-y-1 group">
                  <div className="w-10 h-10 bg-[#dbeafe] rounded-full flex items-center justify-center mx-auto mb-3 group-hover:scale-110 transition-transform duration-300">
                    <Clock className="w-5 h-5 text-[#4A90E2] group-hover:animate-pulse" />
                  </div>
                  <h3 className="font-medium text-[#111827] text-sm mb-1 group-hover:text-[#4A90E2] transition-colors duration-300">
                    Processing Time
                  </h3>
                  <p className="text-[#6b7280] text-sm">2-3 Business</p>
                  <p className="text-[#6b7280] text-sm">Days</p>
                </div>
                <div className="bg-white rounded-lg p-6 text-center border border-[#e5e7eb] hover:border-[#16a34a] transition-all duration-300 hover:shadow-lg transform hover:-translate-y-1 group">
                  <div className="w-10 h-10 bg-[#dcfce7] rounded-full flex items-center justify-center mx-auto mb-3 group-hover:scale-110 transition-transform duration-300">
                    <DollarSign className="w-5 h-5 text-[#16a34a] group-hover:animate-pulse" />
                  </div>
                  <h3 className="font-medium text-[#111827] text-sm mb-1 group-hover:text-[#16a34a] transition-colors duration-300">
                    Total Fees
                  </h3>
                  <p className="text-[#6b7280] text-sm">350 ETB</p>
                </div>
              </div>

              <div>
                <h2 className="text-lg font-medium text-[#111827] mb-4">Required Documents</h2>
                <div className="grid grid-cols-4 gap-4">
                  <div
                    className="bg-[#f3f4f6] rounded-lg p-3 text-center border border-[#e5e7eb] hover:bg-[#e5e7eb] hover:border-[#4A90E2] transition-all duration-300 transform hover:-translate-y-1 hover:shadow-md animate-in fade-in duration-500"
                    style={{ animationDelay: "100ms" }}
                  >
                    <span className="text-sm text-[#6b7280]">2 Passport Photos</span>
                  </div>
                  <div
                    className="bg-[#f3f4f6] rounded-lg p-3 text-center border border-[#e5e7eb] hover:bg-[#e5e7eb] hover:border-[#4A90E2] transition-all duration-300 transform hover:-translate-y-1 hover:shadow-md animate-in fade-in"
                    style={{ animationDelay: "200ms" }}
                  >
                    <span className="text-sm text-[#6b7280]">Original ID Card</span>
                  </div>
                  <div
                    className="bg-[#f3f4f6] rounded-lg p-3 text-center border border-[#e5e7eb] hover:bg-[#e5e7eb] hover:border-[#4A90E2] transition-all duration-300 transform hover:-translate-y-1 hover:shadow-md animate-in fade-in"
                    style={{ animationDelay: "300ms" }}
                  >
                    <span className="text-sm text-[#6b7280]">Current License</span>
                  </div>
                  <div
                    className="bg-[#f3f4f6] rounded-lg p-3 text-center border border-[#e5e7eb] hover:bg-[#e5e7eb] hover:border-[#4A90E2] transition-all duration-300 transform hover:-translate-y-1 hover:shadow-md animate-in fade-in"
                    style={{ animationDelay: "400ms" }}
                  >
                    <span className="text-sm text-[#6b7280]">Application Form</span>
                  </div>
                </div>
              </div>
            </div>

            <div className="bg-white rounded-lg p-6 shadow-sm hover:shadow-md transition-all duration-300 transform hover:-translate-y-1">
              <h2 className="text-lg font-medium text-[#111827] mb-4">Step-by-Step Instructions</h2>
              <div className="space-y-4">
                <div className="flex gap-4 group animate-in fade-in duration-500" style={{ animationDelay: "100ms" }}>
                  <div className="w-8 h-8 bg-[#3A6A8D] text-white rounded-full flex items-center justify-center font-medium text-sm flex-shrink-0 group-hover:scale-110 group-hover:shadow-lg transition-all duration-300">
                    1
                  </div>
                  <div className="flex-1">
                    <h3 className="font-medium text-[#111827] mb-1 group-hover:text-[#3A6A8D] transition-colors duration-300">
                      Fill Application Form
                    </h3>
                    <p className="text-[#6b7280] text-sm mb-2">
                      Download and complete the official renewal form with accurate information.
                    </p>
                    <Button
                      variant="outline"
                      size="sm"
                      className="text-[#3A6A8D] border-[#3A6A8D] bg-white hover:bg-[#3A6A8D] hover:text-white text-xs transition-all duration-300 transform hover:scale-105 hover:shadow-lg"
                    >
                      <Download className="w-3 h-3 mr-1 transition-transform duration-300 group-hover:rotate-12" />
                      Download PDF Form
                    </Button>
                    <Button
                      asChild
                      variant="outline"
                      size="sm"
                      className="text-[#3A6A8D] border-[#3A6A8D] text-xs bg-transparent transition-all duration-300 transform hover:scale-105 hover:shadow-lg"
                    >
                      <a href="/user/notices/1">View Full Notice</a>
                    </Button>
                  </div>
                </div>
                <div className="flex gap-4 group animate-in fade-in duration-500" style={{ animationDelay: "200ms" }}>
                  <div className="w-8 h-8 bg-[#3A6A8D] text-white rounded-full flex items-center justify-center font-medium text-sm flex-shrink-0 group-hover:scale-110 group-hover:shadow-lg transition-all duration-300">
                    2
                  </div>
                  <div className="flex-1">
                    <h3 className="font-medium text-[#111827] mb-1 group-hover:text-[#3A6A8D] transition-colors duration-300">
                      Prepare Required Documents
                    </h3>
                    <p className="text-[#6b7280] text-sm mb-2">
                      Gather all necessary documents and make copies as needed.
                    </p>
                    <div className="bg-[#f0f9ff] p-2 rounded text-xs text-[#0369a1]">
                      💡 AI Tip: Bring originals and copies to avoid delays
                    </div>
                  </div>
                </div>

                <div className="flex gap-4 group animate-in fade-in duration-500" style={{ animationDelay: "300ms" }}>
                  <div className="w-8 h-8 bg-[#3A6A8D] text-white rounded-full flex items-center justify-center font-medium text-sm flex-shrink-0 group-hover:scale-110 group-hover:shadow-lg transition-all duration-300">
                    3
                  </div>
                  <div className="flex-1">
                    <h3 className="font-medium text-[#111827] mb-1 group-hover:text-[#3A6A8D] transition-colors duration-300">
                      Visit License Office
                    </h3>
                    <p className="text-[#6b7280] text-sm">
                      Submit your application and documents at the designated office.
                    </p>
                  </div>
                </div>

                <div className="flex gap-4 group animate-in fade-in duration-500" style={{ animationDelay: "400ms" }}>
                  <div className="w-8 h-8 bg-[#3A6A8D] text-white rounded-full flex items-center justify-center font-medium text-sm flex-shrink-0 group-hover:scale-110 group-hover:shadow-lg transition-all duration-300">
                    4
                  </div>
                  <div className="flex-1">
                    <h3 className="font-medium text-[#111827] mb-1 group-hover:text-[#3A6A8D] transition-colors duration-300">
                      Pay Renewal Fees
                    </h3>
                    <p className="text-[#6b7280] text-sm">
                      Complete payment at the cashier or through online payment if available.
                    </p>
                  </div>
                </div>

                <div className="flex gap-4 group animate-in fade-in duration-500" style={{ animationDelay: "500ms" }}>
                  <div className="w-8 h-8 bg-[#3A6A8D] text-white rounded-full flex items-center justify-center font-medium text-sm flex-shrink-0 group-hover:scale-110 group-hover:shadow-lg transition-all duration-300">
                    5
                  </div>
                  <div className="flex-1">
                    <h3 className="font-medium text-[#111827] mb-1 group-hover:text-[#3A6A8D] transition-colors duration-300">
                      Collect New License
                    </h3>
                    <p className="text-[#6b7280] text-sm">
                      Return after processing period to collect your renewed license.
                    </p>
                  </div>
                </div>
              </div>

              <div className="mt-6">
                <Link href="/user/workspace">
                  <Button className="bg-[#3A6A8D] hover:bg-[#2e4d57] text-white px-6 transition-all duration-300 transform hover:scale-105 hover:shadow-lg active:scale-95">
                    <span className="mr-2 transition-transform duration-300 hover:rotate-12">📋</span>
                    Save Checklist
                  </Button>
                </Link>
              </div>
            </div>

            {/* Feedback Section */}
            <div className="bg-white rounded-lg p-6 shadow-sm hover:shadow-md transition-all duration-300 transform hover:-translate-y-1">
              <div className="flex gap-4 mb-4 border-b border-[#e5e7eb]">
                <button
                  onClick={() => setActiveTab("feedback")}
                  className={`pb-2 px-3 text-xs font-medium transition-all duration-300 ${
                    activeTab === "feedback"
                      ? "text-[#4A90E2] border-b-2 border-[#4A90E2] transform scale-105"
                      : "text-[#6b7280] hover:text-[#111827] hover:scale-105"
                  }`}
                >
                  Feedback
                </button>
                <button
                  onClick={() => setActiveTab("notices")}
                  className={`pb-2 px-3 text-xs font-medium transition-all duration-300 ${
                    activeTab === "notices"
                      ? "text-[#4A90E2] border-b-2 border-[#4A90E2] transform scale-105"
                      : "text-[#6b7280] hover:text-[#111827] hover:scale-105"
                  }`}
                >
                  Notices
                </button>
                <button
                  onClick={() => setActiveTab("discussions")}
                  className={`pb-2 px-3 text-xs font-medium transition-all duration-300 ${
                    activeTab === "discussions"
                      ? "text-[#4A90E2] border-b-2 border-[#4A90E2] transform scale-105"
                      : "text-[#6b7280] hover:text-[#111827] hover:scale-105"
                  }`}
                >
                  Discussions
                </button>
              </div>

              <div className="min-h-[200px]">
                {activeTab === "feedback" && (
                  <div className="space-y-3 animate-in fade-in duration-500">
                    <div className="flex gap-3 p-3 rounded-lg hover:bg-[#f8fafc] transition-all duration-300 transform hover:-translate-y-1">
                      <div className="w-8 h-8 rounded-full bg-[#ced4da] flex items-center justify-center text-xs font-semibold transition-transform duration-300 hover:scale-110 overflow-hidden">
                        <Image src="/sarah-profile.png" alt="Sarah" width={100} height={100} className="w-full h-full object-cover" onError={(e) => { e.currentTarget.style.display = 'none'; }} />
                        <span>SM</span>
                      </div>
                      <div className="flex-1">
                        <div className="flex items-center gap-2 mb-1">
                          <span className="text-sm font-medium text-[#111827]">Sarah M.</span>
                          <Badge className="bg-[#16a34a] text-white text-xs px-1 py-0 animate-pulse">Verified</Badge>
                        </div>
                        <p className="text-xs text-[#6b7280] mb-1">
                          The payment option saved me a lot of time. Highly recommend using it!
                        </p>
                        <div className="flex items-center gap-3 text-xs text-[#9ca3af]">
                          <div className="flex items-center gap-1 hover:text-[#4A90E2] transition-colors duration-300 cursor-pointer">
                            <ThumbsUp className="w-3 h-3 transition-transform duration-300 hover:scale-110" />
                            <span>12</span>
                          </div>
                          <span>3 days ago</span>
                        </div>
                      </div>
                    </div>

                    <div className="flex gap-3 p-3 rounded-lg hover:bg-[#f8fafc] transition-all duration-300 transform hover:-translate-y-1">
                      <div className="w-8 h-8 rounded-full bg-[#ced4da] flex items-center justify-center text-xs font-semibold transition-transform duration-300 hover:scale-110 overflow-hidden">
                        <Image src="/michael-profile.png" alt="Michael" width={100} height={100} className="w-full h-full object-cover" onError={(e) => { e.currentTarget.style.display = 'none'; }} />
                        <span>MT</span>
                      </div>
                      <div className="flex-1">
                        <div className="flex items-center gap-2 mb-1">
                          <span className="text-sm font-medium text-[#111827]">Michael T.</span>
                        </div>
                        <p className="text-xs text-[#6b7280] mb-1">
                          Make sure to bring exact change for fees. They don&#39;t always have change available.
                        </p>
                        <div className="flex items-center gap-3 text-xs text-[#9ca3af]">
                          <div className="flex items-center gap-1 hover:text-[#4A90E2] transition-colors duration-300 cursor-pointer">
                            <ThumbsUp className="w-3 h-3 transition-transform duration-300 hover:scale-110" />
                            <span>8</span>
                          </div>
                          <span>5 days ago</span>
                        </div>
                      </div>
                    </div>
                  </div>
                )}

                {activeTab === "notices" && (
                  <div className="space-y-4 animate-in fade-in duration-500">
                    <div className="bg-white border border-[#e5e7eb] rounded-lg p-4 hover:bg-[#f8fafc] transition-all duration-300 transform hover:-translate-y-1 hover:shadow-md">
                      <div className="flex items-start justify-between mb-3">
                        <div className="flex items-center gap-2">
                          <h4 className="font-medium text-[#111827] text-sm">
                            New Employee Onboarding Process Updates
                          </h4>
                          <Badge className="bg-[#16a34a] text-white text-xs px-2 py-1 animate-pulse">Active</Badge>
                        </div>
                        <div className="flex items-center gap-1 text-[#6b7280]">
                          <ThumbsUp className="w-4 h-4" />
                          <span className="text-sm">12</span>
                        </div>
                      </div>
                      <p className="text-sm text-[#6b7280] mb-3">
                        Updated guidelines for the employee onboarding process, including new documentation requirements
                        and digital workflow procedures.
                      </p>
                      <div className="flex items-center gap-4 text-xs text-[#6b7280] mb-3">
                        <div className="flex items-center gap-1">
                          <FileText className="w-3 h-3" />
                          <span>Published Dec 15, 2024</span>
                        </div>
                        <span>HR Department</span>
                        <div className="flex items-center gap-1">
                          <Globe className="w-3 h-3" />
                          <span>Organization: TechCorp Solutions</span>
                        </div>
                      </div>
                      <div className="flex items-center gap-3">
                        <Button
                          asChild
                          variant="outline"
                          size="sm"
                          className="text-[#3A6A8D] border-[#3A6A8D] text-xs bg-transparent transition-all duration-300 transform hover:scale-105 hover:shadow-lg"
                        >
                          <a href="/user/notices/1">View Full Notice</a>
                        </Button>
                        <Button
                          variant="outline"
                          size="sm"
                          className="text-[#6b7280] border-[#e5e7eb] text-xs bg-transparent transition-all duration-300 transform hover:scale-105 hover:shadow-lg"
                        >
                          <Share className="w-3 h-3 mr-1 transition-transform duration-300 group-hover:rotate-12" />
                          Share
                        </Button>
                        <Button
                          variant="outline"
                          size="sm"
                          className="text-[#6b7280] border-[#e5e7eb] text-xs bg-transparent transition-all duration-300 transform hover:scale-105 hover:shadow-lg"
                        >
                          <FileText className="w-3 h-3 mr-1 transition-transform duration-300 group-hover:rotate-12" />
                          PDF
                        </Button>
                      </div>
                      <Button
                        asChild
                        className="w-full bg-[#3A6A8D] hover:bg-[#2e4d57] text-white text-sm py-2 transition-all duration-300 transform hover:scale-105 hover:shadow-lg active:scale-95"
                      >
                        <a href="/user/notices">View Notices</a>
                      </Button>
                    </div>
                  </div>
                )}

                {activeTab === "discussions" && (
                  <div className="space-y-4 animate-in fade-in duration-500">
                    <div className="bg-white border border-[#e5e7eb] rounded-lg p-4 hover:bg-[#f8fafc] transition-all duration-300 transform hover:-translate-y-1 hover:shadow-md">
                      <div className="flex items-start gap-3">
                        <div className="w-10 h-10 rounded-full bg-[#4A90E2] flex items-center justify-center text-white text-sm font-semibold transition-transform duration-300 hover:scale-110 overflow-hidden">
                          <Image src="/user-profile-illustration.png" alt="User" width={100} height={100} className="w-full h-full object-cover" onError={(e) => { e.currentTarget.style.display = 'none'; }} />
                          <span>AC</span>
                        </div>
                        <div className="flex-1">
                          <div className="flex items-center gap-2 mb-2">
                            <span className="text-sm font-medium text-[#111827]">Alex Chen</span>
                          </div>
                          <h4 className="font-medium text-[#111827] mb-2 text-sm">
                            How to integrate AI tools into daily study routine?
                          </h4>
                          <p className="text-sm text-[#6b7280] mb-3">
                            I&#39;ve been experimenting with various AI tools for studying and note-taking. Would love to
                            hear your experiences and recommendations for the best workflow...
                          </p>
                          <div className="flex items-center gap-2 mb-3">
                            <Badge className="bg-[#dbeafe] text-[#1e40af] text-xs px-2 py-1">AI</Badge>
                            <Badge className="bg-[#f0fdf4] text-[#16a34a] text-xs px-2 py-1">STUDYING</Badge>
                          </div>
                          <div className="flex items-center gap-4 text-xs text-[#9ca3af]">
                            <div className="flex items-center gap-1 hover:text-[#4A90E2] transition-colors duration-300 cursor-pointer">
                              <span>👁</span>
                              <span>24</span>
                            </div>
                            <div className="flex items-center gap-1 hover:text-[#4A90E2] transition-colors duration-300 cursor-pointer">
                              <MessageCircle className="w-3 h-3" />
                              <span>7</span>
                            </div>
                            <div className="flex items-center gap-1 hover:text-[#4A90E2] transition-colors duration-300 cursor-pointer">
                              <ThumbsUp className="w-3 h-3" />
                              <span>12</span>
                            </div>
                            <div className="flex items-center gap-1 hover:text-[#4A90E2] transition-colors duration-300 cursor-pointer">
                              <Share className="w-3 h-3" />
                              <span>3</span>
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>
                    <Button asChild className="w-full bg-[#3A6A8D] hover:bg-[#2e4d57] text-white text-sm py-2 transition-all duration-300 transform hover:scale-105 hover:shadow-lg active:scale-95">
                      <Link href="/user/discussions">View Discussions</Link>
                    </Button>
                  </div>
                )}
              </div>

              {activeTab === "feedback" && (
                <Button asChild className="w-full bg-[#3A6A8D] hover:bg-[#2e4d57] text-white mt-6 text-sm py-2 transition-all duration-300 transform hover:scale-105 hover:shadow-lg active:scale-95">
                  <a href="/user/feedback">Add Feedback</a>
                </Button>
              )}
            </div>
          </div>
        </div>
        </main>
        <aside className="w-96 p-6 space-y-6  mr-48">
          <div className="bg-white rounded-lg p-6 shadow-sm hover:shadow-md transition-all duration-300 transform hover:-translate-y-1">
            <h3 className="font-medium text-[#111827] mb-3">AI Assistant</h3>
            <p className="text-sm text-[#6b7280] mb-4">Need help with any step? Ask me anything!</p>
            <Input
              placeholder="Ask a question..."
              className="mb-4 border-[#e5e7eb] text-sm focus:border-[#4A90E2] focus:ring-2 focus:ring-[#4A90E2]/20 transition-all duration-300"
            />
            <div className="space-y-1 text-sm">
              <button className="text-[#6b7280] hover:text-[#4A90E2] block w-full text-left p-2 rounded hover:bg-[#f8fafc] transition-all duration-300 transform hover:translate-x-1">
                What if I lost my current license?
              </button>
              <button className="text-[#6b7280] hover:text-[#4A90E2] block w-full text-left p-2 rounded hover:bg-[#f8fafc] transition-all duration-300 transform hover:translate-x-1">
                Can I renew online?
              </button>
              <button className="text-[#6b7280] hover:text-[#4A90E2] block w-full text-left p-2 rounded hover:bg-[#f8fafc] transition-all duration-300 transform hover:translate-x-1">
                Office working hours?
              </button>
            </div>
          </div>

        </aside>
      </div>
    </div>
  )
}
