"use client"

import { useState } from "react"
import {
  ArrowLeft,
  Bold,
  Italic,
  List,
  Upload,
  Calendar,
  Clock,
  Eye,
  Copy,
  Bell,
  Facebook,
  Instagram,
  Twitter,
  Youtube
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Textarea } from "@/components/ui/textarea"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import Link from "next/link"

export default function CreateOfficialNotice() {
  const [noticeTitle, setNoticeTitle] = useState("")
  const [noticeDescription, setNoticeDescription] = useState("")

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50">
  {/* Header removed as requested */}

      <div className="max-w-7xl mx-auto px-6 py-8">
        <div className="flex items-center gap-4 mb-8">
          <Link href="/admin/notices">
            <Button variant="ghost" size="sm" className="p-2 hover:bg-slate-100 rounded-lg transition-colors">
              <ArrowLeft className="h-5 w-5 text-slate-600" />
            </Button>
          </Link>
          <h1 className="text-2xl font-semibold text-slate-900 tracking-tight">Create Official Notice</h1>
        </div>
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Main Form */}
          <div className="lg:col-span-2 space-y-8">
            <Card className="bg-white/70 backdrop-blur-sm border-slate-200/60 shadow-lg shadow-slate-200/50">
              <CardContent className="p-6 space-y-6">
                {/* Notice Title */}
                <div className="space-y-3">
                  <label className="text-sm font-semibold text-slate-700 tracking-wide">Notice Title</label>
                  <Input
                    placeholder="Enter notice title..."
                    value={noticeTitle}
                    onChange={(e) => setNoticeTitle(e.target.value)}
                    className="bg-white/80 border-slate-200 focus:border-blue-400 focus:ring-blue-400/20 h-12 text-base placeholder:text-slate-400"
                  />
                </div>

                {/* Notice Description */}
                <div className="space-y-3">
                  <label className="text-sm font-semibold text-slate-700 tracking-wide">Notice Description</label>
                  <div className="bg-white/80 border border-slate-200 rounded-xl overflow-hidden shadow-sm">
                    {/* Rich Text Toolbar */}
                    <div className="flex items-center gap-1 p-3 border-b border-slate-200 bg-slate-50/50">
                      <Button variant="ghost" size="sm" className="p-2 hover:bg-white rounded-lg transition-colors">
                        <Bold className="h-4 w-4 text-slate-600" />
                      </Button>
                      <Button variant="ghost" size="sm" className="p-2 hover:bg-white rounded-lg transition-colors">
                        <Italic className="h-4 w-4 text-slate-600" />
                      </Button>
                      <Button variant="ghost" size="sm" className="p-2 hover:bg-white rounded-lg transition-colors">
                        <List className="h-4 w-4 text-slate-600" />
                      </Button>
                      <Button variant="ghost" size="sm" className="p-2 hover:bg-white rounded-lg transition-colors">
                        <List className="h-4 w-4 text-slate-600" />
                      </Button>
                    </div>
                    <Textarea
                      placeholder="Enter notice description..."
                      value={noticeDescription}
                      onChange={(e) => setNoticeDescription(e.target.value)}
                      className="border-none resize-none min-h-32 focus-visible:ring-0 text-base placeholder:text-slate-400 bg-transparent"
                    />
                  </div>
                </div>

                {/* Attachments */}
                <div className="space-y-3">
                  <label className="text-sm font-semibold text-slate-700 tracking-wide">Attachments</label>
                  <div className="bg-gradient-to-br from-slate-50 to-blue-50/30 border-2 border-dashed border-slate-300 hover:border-blue-400 rounded-xl p-10 transition-colors group cursor-pointer">
                    <div className="text-center">
                      <div className="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-4 group-hover:bg-blue-200 transition-colors">
                        <Upload className="h-8 w-8 text-blue-600" />
                      </div>
                      <p className="text-slate-700 font-medium">
                        Drop files here or{" "}
                        <button className="text-blue-600 hover:text-blue-700 font-semibold underline underline-offset-2">
                          browse
                        </button>
                      </p>
                      <p className="text-sm text-slate-500 mt-2">PDF, DOC, PNG, JPG up to 10MB</p>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card className="bg-white/70 backdrop-blur-sm border-slate-200/60 shadow-lg shadow-slate-200/50">
              <CardContent className="p-6 space-y-6">
                {/* Related Department */}
                <div className="space-y-3">
                  <label className="text-sm font-semibold text-slate-700 tracking-wide">Related Department</label>
                  <Select>
                    <SelectTrigger className="bg-white/80 border-slate-200 focus:border-blue-400 h-12 text-base">
                      <SelectValue placeholder="Select Department" className="text-slate-400" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="hr">HR Department</SelectItem>
                      <SelectItem value="it">IT Department</SelectItem>
                      <SelectItem value="finance">Finance Department</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                {/* Publication Status and Priority */}
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div className="space-y-3">
                    <label className="text-sm font-semibold text-slate-700 tracking-wide">Publication Status</label>
                    <Select defaultValue="active">
                      <SelectTrigger className="bg-white/80 border-slate-200 focus:border-blue-400 h-12 text-base">
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="active">Active</SelectItem>
                        <SelectItem value="draft">Draft</SelectItem>
                        <SelectItem value="archived">Archived</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="space-y-3">
                    <label className="text-sm font-semibold text-slate-700 tracking-wide">Priority Level</label>
                    <Select defaultValue="medium">
                      <SelectTrigger className="bg-white/80 border-slate-200 focus:border-blue-400 h-12 text-base">
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="low">Low</SelectItem>
                        <SelectItem value="medium">Medium</SelectItem>
                        <SelectItem value="high">High</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>

                {/* Publication Date and Time */}
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div className="space-y-3">
                    <label className="text-sm font-semibold text-slate-700 tracking-wide">Publication Date</label>
                    <div className="relative">
                      <Input
                        type="text"
                        placeholder="mm/dd/yyyy"
                        className="bg-white/80 border-slate-200 focus:border-blue-400 h-12 text-base pr-12 placeholder:text-slate-400"
                      />
                      <Calendar className="absolute right-4 top-1/2 transform -translate-y-1/2 h-5 w-5 text-slate-400" />
                    </div>
                  </div>
                  <div className="space-y-3">
                    <label className="text-sm font-semibold text-slate-700 tracking-wide">Publication Time</label>
                    <div className="relative">
                      <Input
                        type="text"
                        placeholder="--:-- --"
                        className="bg-white/80 border-slate-200 focus:border-blue-400 h-12 text-base pr-12 placeholder:text-slate-400"
                      />
                      <Clock className="absolute right-4 top-1/2 transform -translate-y-1/2 h-5 w-5 text-slate-400" />
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            <div className="flex flex-wrap gap-4 pt-2">
              <Button className="bg-[#3A6A8D] hover:bg-[#2d5470] text-white shadow-lg shadow-blue-200 h-12 px-6 font-semibold tracking-wide">
                <Upload className="h-4 w-4 mr-2" />
                Publish Notice
              </Button>
              <Button className="bg-secondary hover:bg-[#5E9C8F] text-white shadow-lg  h-12 px-6 font-semibold tracking-wide">
                <Copy className="h-4 w-4 mr-2" />
                Save Draft
              </Button>
              <Button
                variant="outline"
                className="text-slate-600 border-slate-300 hover:bg-slate-50 h-12 px-6 font-medium bg-transparent"
              >
                Cancel
              </Button>
            </div>
          </div>

          {/* Preview Panel */}
          <div className="space-y-6">
            <Card className="bg-white/80 backdrop-blur-sm border-slate-200/60 shadow-xl shadow-slate-200/50 sticky top-24">
              <CardHeader className="pb-4 border-b border-slate-100">
                <CardTitle className="flex items-center gap-3 text-lg text-slate-800">
                  <div className="w-8 h-8 bg-blue-100 rounded-lg flex items-center justify-center">
                    <Eye className="h-4 w-4 text-blue-600" />
                  </div>
                  Preview
                </CardTitle>
              </CardHeader>
              <CardContent className="p-6 space-y-6">
                <div className="space-y-4">
                  <h3 className="font-bold text-slate-900 text-lg leading-tight">
                    {noticeTitle || "Sample Notice Title"}
                  </h3>
                  <p className="text-slate-600 leading-relaxed">
                    {noticeDescription ||
                      "This is how your notice description will appear to users. The formatting and content will be preserved as entered."}
                  </p>
                  <div className="flex items-center gap-3 pt-2">
                    <span className="text-sm font-medium text-slate-600">HR Department</span>
                    <Badge className="bg-emerald-100 text-emerald-700 hover:bg-emerald-100 font-medium px-3 py-1">
                      Active
                    </Badge>
                  </div>
                </div>

                <div className="border-t border-slate-100 pt-6">
                  <h4 className="font-semibold text-slate-800 mb-4 flex items-center gap-2">
                    <Clock className="h-4 w-4 text-slate-500" />
                    Version History
                  </h4>
                  <div className="space-y-3">
                    <div className="flex justify-between items-center p-3 bg-slate-50 rounded-lg">
                      <span className="text-sm font-medium text-slate-700">Current Draft</span>
                      <span className="text-xs text-slate-500 bg-white px-2 py-1 rounded-full">Just now</span>
                    </div>
                    <div className="flex justify-between items-center p-3 bg-slate-50 rounded-lg">
                      <span className="text-sm font-medium text-slate-700">Version 1.0</span>
                      <span className="text-xs text-slate-500 bg-white px-2 py-1 rounded-full">2 days ago</span>
                    </div>
                  </div>
                </div>

                <div className="space-y-3 pt-4 border-t border-slate-100">
                  <Button
                    variant="outline"
                    className="w-full justify-center bg-white/80 border-slate-200 hover:bg-slate-50 h-11 font-medium"
                  >
                    <Copy className="h-4 w-4 mr-2 text-slate-600" />
                    Duplicate Notice
                  </Button>
                  <Button
                    variant="outline"
                    className="w-full justify-center bg-white/80 border-slate-200 hover:bg-slate-50 h-11 font-medium"
                  >
                    <Bell className="h-4 w-4 mr-2 text-slate-600" />
                    Send Notifications
                  </Button>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>

      
    </div>
  )
}
