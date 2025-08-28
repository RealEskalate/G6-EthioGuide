"use client";

import { useState } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Calendar, Heart, Share2, FileText } from "lucide-react";

interface Notice {
  id: number;
  title: string;
  status: "Active" | "Upcoming" | "Expired";
  statusColor: string;
  description: string;
  published: string;
  department: string;
  organization: string;
  likes: number;
}

const notices: Notice[] = [
  {
    id: 1,
    title: "New Employee Onboarding Process Updates",
    status: "Active",
    statusColor: "bg-green-100 text-green-800",
    description:
      "Updated guidelines for the employee onboarding process, including new documentation requirements and digital workflow procedures.",
    published: "Dec 15, 2024",
    department: "HR Department",
    organization: "TechCorp Solutions",
    likes: 12,
  },
  {
    id: 2,
    title: "System Maintenance Schedule",
    status: "Upcoming",
    statusColor: "bg-yellow-100 text-yellow-800",
    description:
      "Scheduled system maintenance on December 25th from 2:00 AM to 6:00 AM. All services will be temporarily unavailable during this period.",
    published: "Dec 14, 2024",
    department: "IT Operations",
    organization: "TechCorp Solutions",
    likes: 8,
  },
  {
    id: 3,
    title: "Holiday Policy Updates",
    status: "Expired",
    statusColor: "bg-gray-100 text-gray-500",
    description:
      "Updates to the company holiday policy for 2024, including new national holidays and revised vacation request procedures.",
    published: "Nov 30, 2024",
    department: "HR Department",
    organization: "TechCorp Solutions",
    likes: 25,
  },
];

export default function NoticesPage() {
  const [status, setStatus] = useState("all");
  const [department, setDepartment] = useState("all");

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-5xl mx-auto">
        {/* Header */}
        <div className="flex items-center justify-between mb-8">
          <div>
            <div className="flex items-center gap-3 mb-2">
              <FileText className="h-8 w-8 text-[#3A6A8D]" />
              <h1 className="text-3xl font-bold text-gray-900">Official Notices</h1>
            </div>
            <p className="text-gray-600">Get notices of different organizations</p>
          </div>
        </div>

        {/* Search and Filters */}
        <Card className="p-4 mb-6">
          <div className="flex flex-col sm:flex-row gap-4 w-full">
            <div className="relative flex-1">
              <input
                type="text"
                placeholder="Search notices..."
                className="w-full pl-10 pr-4 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#3A6A8D] focus:border-transparent"
              />
            </div>
            <div className="flex gap-2 flex-1">
              <Select value={status} onValueChange={setStatus}>
                <SelectTrigger className="w-full bg-white rounded-lg border border-gray-200 shadow-sm px-4 py-2 text-gray-900 focus:ring-2 focus:ring-[#3A6A8D] focus:border-transparent transition-all">
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
                <SelectTrigger className="w-full bg-white rounded-lg border border-gray-200 shadow-sm px-4 py-2 text-gray-900 focus:ring-2 focus:ring-[#3A6A8D] focus:border-transparent transition-all">
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
        <div className="space-y-6">
          {notices.map((notice, index) => (
            <Card key={notice.id} className={`bg-white p-6 animate-in fade-in slide-in-from-bottom-4`} style={{ animationDelay: `${index * 100}ms` }}>
              <CardContent className="p-0">
                <div className="flex items-center gap-3 mb-2">
                  <h2 className="text-lg font-semibold text-gray-900 mr-2">
                    {notice.title}
                  </h2>
                  <Badge className={notice.statusColor}>{notice.status}</Badge>
                </div>
                <p className="text-gray-600 mb-2">{notice.description}</p>
                <div className="flex flex-wrap items-center gap-4 text-sm text-gray-500 mb-2">
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
                <div className="flex items-center gap-6 pt-2 text-sm text-gray-500">
                  <a href="#" className="hover:underline flex items-center gap-1">
                    View Full Notice
                  </a>
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
