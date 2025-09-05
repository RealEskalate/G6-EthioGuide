"use client"

import React, { useState } from "react"
import { ArrowLeft, Send } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Textarea } from "@/components/ui/textarea"
import { Badge } from "@/components/ui/badge"

export default function FeedbackPage() {
  const [feedbackType, setFeedbackType] = useState("report-issue")
  const [subject, setSubject] = useState("")
  const [detailedFeedback, setDetailedFeedback] = useState("")

  const feedbackHistory = [
    {
      id: 1,
      type: "Report Issue",
      subject: "Login page not working properly",
      date: "2024-01-15",
      status: "Pending",
      statusColor: "bg-yellow-100 text-yellow-800 border-yellow-200",
    },
    {
      id: 2,
      type: "Suggest Improvement",
      subject: "Add dark mode to the interface",
      date: "2024-01-10",
      status: "Resolved",
      statusColor: "bg-green-100 text-green-800 border-green-200",
    },
  ]

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    console.log({ feedbackType, subject, detailedFeedback })
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="bg-white border-b border-gray-200 px-4 py-3">
        <div className="flex items-center justify-between max-w-7xl mx-auto">
          <div className="flex items-center gap-3">
            <Button
              variant="ghost"
              size="sm"
              className="p-2 hover:bg-gray-100"
              onClick={() => {
                window.location.href = "/user/procedures-detail";
              }}
            >
              <ArrowLeft className="h-5 w-5 text-gray-600" />
            </Button>
            <h1 className="text-xl font-semibold text-gray-900">Feedback</h1>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="max-w-7xl mx-auto p-6">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Submit Feedback Section */}
          <Card className="shadow-sm">
            <CardHeader className="pb-4">
              <CardTitle className="text-lg font-semibold text-gray-900">Submit Feedback</CardTitle>
            </CardHeader>
            <CardContent className="space-y-6">
              <form onSubmit={handleSubmit} className="space-y-6">
                {/* Type of Feedback */}
                <div className="space-y-2">
                  <Label htmlFor="feedback-type" className="text-sm font-medium text-gray-700">
                    Type of Feedback
                  </Label>
                  <Select value={feedbackType} onValueChange={setFeedbackType}>
                    <SelectTrigger className="bg-white border-gray-300 animate-in fade-in duration-500">
                      <SelectValue placeholder="Report Issue" />
                    </SelectTrigger>
                    <SelectContent className="animate-in fade-in duration-500">
                      <SelectItem value="report-issue" className="!bg-white hover:!bg-gray-200 animate-in fade-in duration-700">Report Issue</SelectItem>
                      <SelectItem value="suggest-improvement" className="!bg-white hover:!bg-gray-200 animate-in fade-in duration-800">Suggest Improvement</SelectItem>
                      <SelectItem value="general-feedback" className="!bg-white hover:!bg-gray-200 animate-in fade-in duration-900">General Feedback</SelectItem>
                      <SelectItem value="feature-request" className="!bg-white hover:!bg-gray-200 animate-in fade-in duration-1000">Feature Request</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                {/* Subject/Procedure/Document */}
                <div className="space-y-2">
                  <Label htmlFor="subject" className="text-sm font-medium text-gray-700">
                    Subject/Procedure/Document
                  </Label>
                  <Input
                    id="subject"
                    placeholder="Describe your report procedure..."
                    value={subject}
                    onChange={(e) => setSubject(e.target.value)}
                    className="bg-gray-50 border-gray-300"
                  />
                </div>

                {/* Detailed Feedback */}
                <div className="space-y-2">
                  <Label htmlFor="detailed-feedback" className="text-sm font-medium text-gray-700">
                    Detailed Feedback
                  </Label>
                  <Textarea
                    id="detailed-feedback"
                    placeholder="Please provide detailed feedback..."
                    value={detailedFeedback}
                    onChange={(e) => setDetailedFeedback(e.target.value)}
                    className="min-h-[120px] bg-gray-50 border-gray-300 resize-none"
                  />
                </div>

                <Button type="submit" className="w-full bg-slate-600 hover:bg-slate-700 text-white py-3 rounded-md">
                  <Send className="h-4 w-4 mr-2" />
                  Send Feedback
                </Button>
              </form>
            </CardContent>
          </Card>

          {/* Your Feedback History Section */}
          <Card className="shadow-sm">
            <CardHeader className="pb-4">
              <CardTitle className="text-lg font-semibold text-gray-900">Your Feedback History</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              {feedbackHistory.map((feedback) => (
                <div key={feedback.id} className="border border-gray-200 rounded-lg p-4 space-y-3 bg-white">
                  <div className="flex items-center justify-between">
                    <span className="text-sm font-medium text-gray-700">{feedback.type}</span>
                    <Badge variant="outline" className={`${feedback.statusColor} text-xs px-2 py-1 rounded-full`}>
                      {feedback.status}
                    </Badge>
                  </div>
                  <p className="text-sm text-gray-600">{feedback.subject}</p>
                  <p className="text-xs text-gray-400">{feedback.date}</p>
                </div>
              ))}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}