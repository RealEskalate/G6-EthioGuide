"use client";

import { useState } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Calendar, Heart, Share2, FileText } from "lucide-react";

import { notices } from "@/lib/noticesData";

import { useRouter } from "next/navigation";

export default function NoticesPage() {
  const [status, setStatus] = useState("all");
  const [department, setDepartment] = useState("all");
  const router = useRouter();

  return (
    <div className="min-h-screen bg-gray-50 p-4 sm:p-6">
      <div className="max-w-5xl mx-auto">
        {/* Header */}
        <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between mb-6 sm:mb-8 gap-3 sm:gap-0">
          <div>
            <div className="flex items-center gap-3 mb-2">
              <FileText className="h-7 w-7 sm:h-8 sm:w-8 text-[#3A6A8D]" />
              <h1 className="text-2xl sm:text-3xl font-bold text-gray-900">Official Notices</h1>
            </div>
            <p className="text-gray-600 text-sm">Get notices of different organizations</p>
          </div>
          {/* ...optional right-side actions... */}
        </div>

        {/* Search and Filters */}
        <Card className="p-4 mb-6">
          <div className="flex flex-col sm:flex-row gap-3 sm:gap-4 w-full">
            <div className="relative flex-1">
              <input
                type="text"
                placeholder="Search notices..."
                className="w-full pl-10 pr-4 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#3A6A8D] focus:border-transparent"
              />
              {/* optional search icon slot */}
              {/* ...existing icon if any... */}
            </div>
            <div className="flex gap-2 flex-1">
              <Select value={status} onValueChange={setStatus}>
                <SelectTrigger className="w-full">
                  <SelectValue placeholder="All Status" />
                </SelectTrigger>
                <SelectContent className="rounded-lg border border-gray-200 shadow-md bg-white">
                  <SelectItem value="all">All Status</SelectItem>
                  <SelectItem value="active">Active</SelectItem>
                  <SelectItem value="upcoming">Upcoming</SelectItem>
                  <SelectItem value="expired">Expired</SelectItem>
                </SelectContent>
              </Select>
              <Select value={department} onValueChange={setDepartment}>
                <SelectTrigger className="w-full">
                  <SelectValue placeholder="All Departments" />
                </SelectTrigger>
                <SelectContent className="rounded-lg border border-gray-200 shadow-md bg-white">
                  <SelectItem value="all">All Departments</SelectItem>
                  <SelectItem value="hr">HR Department</SelectItem>
                  <SelectItem value="it">IT Operations</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
        </Card>

        {/* Notices List */}
        <div className="space-y-4 sm:space-y-6">
          {notices.map((notice) => (
            <Card
              key={notice.id}
              className="bg-white p-4 sm:p-6 hover:shadow-lg transition-all duration-300"
            >
              <CardContent className="p-0">
                <div className="flex items-center gap-3 mb-2">
                  <h2 className="text-base sm:text-lg font-semibold text-gray-900 mr-2">
                    {notice.title}
                  </h2>
                  <Badge className={notice.statusColor}>{notice.status}</Badge>
                </div>
                <p className="text-gray-600 text-sm sm:text-base mb-2">{notice.description}</p>
                <div className="flex flex-wrap items-center gap-3 sm:gap-4 text-xs sm:text-sm text-gray-500 mb-2">
                  <div className="flex items-center gap-1">
                    <Calendar className="w-4 h-4" />
                    <span>Published: {notice.published}</span>
                  </div>
                  <Badge variant="outline" className="text-xs cursor-pointer hover:bg-blue-100 hover:text-blue-700">
                    {notice.department}
                  </Badge>
                  <span>
                    Organization: <span className="font-semibold text-gray-700">{notice.organization}</span>
                  </span>
                </div>
                <div className="flex flex-wrap items-center gap-3 sm:gap-6 pt-2 text-xs sm:text-sm text-gray-500">
                  <button
                    className="hover:underline flex items-center gap-1 text-blue-700"
                    onClick={() => router.push(`/user/notices/${notice.id}`)}
                  >
                    View Full Notice
                  </button>
                  <button className="flex items-center gap-1 hover:text-blue-700">
                    <Share2 className="w-4 h-4" /> Share
                  </button>
                  <button className="flex items-center gap-1 hover:text-blue-700">
                    <FileText className="w-4 h-4" /> PDF
                  </button>
                  <div className="flex items-center gap-1 ml-auto">
                    <Heart className="w-4 h-4" />
                    <span>{notice.likes}</span>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    </div>
  );
}
