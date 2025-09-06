"use client";

import {  Calendar, CheckCircle, Clock, FileText } from "lucide-react"
import { Button } from "@/components/ui/button"
import { CardContent, Card } from "@/components/ui/card"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Badge } from "@/components/ui/badge"
import { useRouter } from "next/navigation"
import { useMemo, useEffect } from "react"; // added useEffect
import { useGetMyProceduresQuery } from "@/app/store/slices/workspaceSlice"; // added

export default function WorkspacePage() {
  const router = useRouter();

  // log token once on mount
  useEffect(() => {
    if (typeof window === "undefined") return;
    const lsToken =
      localStorage.getItem("accessToken") ||
      localStorage.getItem("token") ||
      localStorage.getItem("authToken") ||
      localStorage.getItem("access_token");
    const cookieMatch = document.cookie.match(/(?:^|; )accessToken=([^;]+)/);
    const cookieToken = cookieMatch ? decodeURIComponent(cookieMatch[1]) : null;
    const envToken = process.env.NEXT_PUBLIC_ACCESS_TOKEN || process.env.ACCESS_TOKEN || null;
    const token = lsToken || cookieToken || envToken || null;
    console.log("Workspace token:", token);
  }, []);

  // fetch user procedures
  const { data, isLoading, isError } = useGetMyProceduresQuery({ page: 1, limit: 20 }); // added
  const items = data?.data ?? []; // added

  // derive stats from API
  const stats = useMemo(() => {
    const total = items.length || Number(data?.total ?? 0);
    const inProgress = items.filter(p => p.status === "IN_PROGRESS").length;
    const completed = items.filter(p => p.status === "COMPLETED").length;
    return [
      { title: "Total Procedures", value: String(total), icon: FileText, color: "text-blue-600", bgColor: "bg-blue-50" },
      { title: "In Progress", value: String(inProgress), icon: Clock, color: "text-orange-600", bgColor: "bg-orange-50" },
      { title: "Completed", value: String(completed), icon: CheckCircle, color: "text-green-600", bgColor: "bg-green-50" },
      { title: "Documents", value: "—", icon: FileText, color: "text-purple-600", bgColor: "bg-purple-50" }, // unknown from API
    ];
  }, [items, data]); // added

  // map API items to card-friendly shape
  const cards = useMemo(() => {
    const toBadge = (s: string) =>
      s === "COMPLETED"
        ? "bg-green-100 text-green-800"
        : s === "IN_PROGRESS"
          ? "bg-orange-100 text-orange-800"
          : "bg-gray-100 text-gray-800";
    const toLabel = (s: string) =>
      s === "COMPLETED" ? "Completed" : s === "IN_PROGRESS" ? "In Progress" : "Not Started";
    const toButton = (s: string) =>
      s === "COMPLETED" ? { text: "View Details", variant: "outline" as const } : { text: "Continue", variant: "default" as const };

    return items.map((p) => ({
      id: p.userProcedureId,
      title: p.procedureTitle,
      department: "", // not provided by API
      status: toLabel(p.status),
      statusRaw: p.status,
      progress: Math.max(0, Math.min(100, Number(p.percent ?? 0))),
      statusColor: toBadge(p.status),
      ...toButton(p.status),
    }));
  }, [items]); // added

  return (
    <main className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header Section */}
        <div className="flex items-center justify-between mb-8">
          <div>
            <h1 className="text-3xl font-bold text-gray-900 mb-2">My Workspace</h1>
            <p className="text-neutral">Track and manage your ongoing procedures</p>
          </div>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          {stats.map((stat, index) => (
            <Card
              key={index}
              className="border-0 shadow-sm hover:shadow-md transition-all duration-300 hover:-translate-y-1 cursor-pointer bg-white"
            >
              <CardContent className="p-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-neutral mb-1">{stat.title}</p>
                    <p className="text-3xl font-bold text-gray-900">{stat.value}</p>
                  </div>
                  <div className={`p-3 rounded-lg ${stat.bgColor} transition-transform duration-200 hover:scale-110`}>
                    <stat.icon className={`w-6 h-6 ${stat.color}`} />
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>

        {/* Filters */}
        <div className="flex gap-4 mb-6">
          <div className="flex items-center gap-2">
            <span className="text-sm font-medium text-gray-700 ">Status:</span>
            <Select defaultValue="all">
              <SelectTrigger className="w-32">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All</SelectItem>
                <SelectItem value="in-progress">In Progress</SelectItem>
                <SelectItem value="completed">Completed</SelectItem>
                <SelectItem value="not-started">Not Started</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div className="flex items-center gap-2 ">
            <span className="text-sm  font-medium text-gray-700">Organization:</span>
            <Select defaultValue="all">
              <SelectTrigger className="w-48">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Organizations</SelectItem>
                <SelectItem value="immigration">Immigration Department</SelectItem>
                <SelectItem value="road">Road Authority</SelectItem>
                <SelectItem value="bank">National Bank</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>

        {/* Status messages */}
        {isLoading && <div className="text-sm text-gray-600 mb-4">Loading your procedures...</div>}
        {!isLoading && isError && <div className="text-sm text-red-600 mb-4">Failed to load procedures.</div>}

        {/* Procedure Cards */}
        <div className="space-y-4">
          {!isLoading && !isError && cards.map((procedure, index) => (
            <Card
              key={procedure.id}
              className="border-0 shadow-sm hover:shadow-lg transition-all duration-300 hover:-translate-y-1 animate-in fade-in slide-in-from-bottom-4 bg-white"
              style={{ animationDelay: `${index * 100}ms` }}
            >
              <CardContent className="p-6">
                <div className="flex items-start justify-between mb-4">
                  <div className="flex items-start gap-4">
                    <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center transition-transform duration-200 hover:scale-110">
                      <FileText className="w-6 h-6 text-blue-600" />
                    </div>
                    <div>
                      <h3 className="font-semibold text-lg text-gray-900 mb-1 hover:text-[#3A6A8D] transition-colors duration-200">
                        {procedure.title}
                      </h3>
                      {procedure.department && <p className="text-gray-600 text-sm">{procedure.department}</p>}
                    </div>
                  </div>
                  <div className="flex items-center gap-3">
                    <Badge className={`${procedure.statusColor} transition-all duration-200`}>
                      {procedure.status}
                    </Badge>
                    <Button
                      variant={procedure.variant}
                      size="sm"
                      className={procedure.variant === "default"
                        ? "bg-[#3A6A8D] hover:bg-[#2d5470] text-white transition-all duration-200 hover:scale-105"
                        : "transition-all duration-200"}
                      onClick={() => router.push("/user/checklist")}
                    >
                      {procedure.text}
                    </Button>
                  </div>
                </div>

                <div className="mb-4">
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-sm font-medium text-gray-700">Progress</span>
                    <span className="text-sm text-gray-600">{procedure.progress}% Complete</span>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2 overflow-hidden">
                    <div
                      className={`h-2 rounded-full transition-all duration-1000 ease-out ${
                        procedure.statusRaw === "COMPLETED"
                          ? "bg-[#5E9C8D]"
                          : procedure.statusRaw === "IN_PROGRESS"
                            ? "bg-[#FEF9C3]"
                            : "bg-gray-300"
                      }`}
                      style={{ width: `${procedure.progress}%`, animationDelay: `${index * 200 + 500}ms` }}
                    />
                  </div>
                </div>

                {/* Procedure Details (placeholder fields not in API kept guarded) */}
                <div className="flex items-center gap-6 text-sm text-gray-600">
                  {/* ...existing optional meta fields if available... */}
                </div>
              </CardContent>
            </Card>
          ))}
          {!isLoading && !isError && cards.length === 0 && (
            <div className="text-sm text-gray-600">No procedures yet.</div>
          )}
        </div>
      </div>
    </main>
  )
}