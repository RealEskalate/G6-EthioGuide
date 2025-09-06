"use client";

import {  Calendar, CheckCircle, Clock, FileText } from "lucide-react"
import { Button } from "@/components/ui/button"
import { CardContent, Card } from "@/components/ui/card"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Badge } from "@/components/ui/badge"
import { useRouter } from "next/navigation"
import { useEffect, useMemo } from "react";
import { useGetMyProceduresQuery } from "@/app/store/slices/workspaceSlice";

export default function WorkspacePage() {
  const router = useRouter();
  const { data, error, isLoading } = useGetMyProceduresQuery({ page: 1, limit: 20 });

  useEffect(() => {
    console.log("API data:", data);
    console.log("API error:", error);
    console.log("API loading:", isLoading);
  }, [data, error, isLoading]);

  // Derive stats from API response
  const { totalCount, inProgressCount, completedCount, documentsCount } = useMemo(() => {
    const items = data?.data ?? [];
    const total = typeof data?.total === "number" ? data.total : items.length;
    const inProg = items.filter((p) => p.status === "IN_PROGRESS").length;
    const completed = items.filter((p) => p.status === "COMPLETED").length;
    // Backend doesn't provide documents metric yet; keep 0 (adjust when available)
    const documents = 0;
    return { totalCount: total, inProgressCount: inProg, completedCount: completed, documentsCount: documents };
  }, [data]);

  const stats = [
    {
      title: "Total Procedures",
      value: String(totalCount),
      icon: FileText,
      color: "text-[#3A6A8D]",
      bgColor: "bg-blue-50",
    },
    {
      title: "In Progress",
      value: String(inProgressCount),
      icon: Clock,
      color: "text-secondary",
      bgColor: "bg-orange-50",
    },
    {
      title: "Completed",
      value: String(completedCount),
      icon: CheckCircle,
      color: "text-[#1C3B2E]",
      bgColor: "bg-green-50",
    },
    {
      title: "Documents",
      value: String(documentsCount), // placeholder until backend provides docs metric
      icon: FileText,
      color: "text-[#1C3B2E]",
      bgColor: "bg-purple-50",
    },
  ]

  // Define gradient backgrounds for the 4 count cards (order matters)
  const statCardBackgrounds = [
    "bg-gradient-to-br from-[#e6f0f5] to-[#d1e7f0]",
    "bg-gradient-to-br from-[#e8f4f2] to-[#d1ede7]",
    "bg-gradient-to-br from-[#e3e8ea] to-[#d6dde0]",
    "bg-gradient-to-br from-[#f0f2f3] to-[#e6eaeb]",
  ]

  // Map backend status to UI labels/colors/button
  const mapStatusUI = (status: "NOT_STARTED" | "IN_PROGRESS" | "COMPLETED") => {
    switch (status) {
      case "COMPLETED":
        return { label: "Completed", color: "bg-green-100 text-green-800", buttonText: "View Details", buttonVariant: "outline" as const };
      case "IN_PROGRESS":
        return { label: "In Progress", color: "bg-orange-100 text-orange-800", buttonText: "Continue", buttonVariant: "default" as const };
      default:
        return { label: "Not Started", color: "bg-gray-100 text-gray-800", buttonText: "Start Now", buttonVariant: "default" as const };
    }
  };

  // Build cards from API response (only API data, no mocks)
  const apiProcedures = useMemo(() => {
    const items = data?.data ?? [];
    return items.map((p, idx) => {
      const ui = mapStatusUI(p.status);
      return {
        id: p.userProcedureId || `${p.procedureId}-${idx}`,
        title: p.procedureTitle,
        department: "—", // backend doesn't provide; keep placeholder
        status: ui.label,
        progress: p.percent ?? 0,
        startDate: undefined,
        estimatedCompletion: undefined,
        documentsUploaded: undefined,
        completedDate: undefined,
        requirements: undefined,
        readyToStart: undefined,
        documentsRequired: undefined,
        addedDate: undefined,
        statusColor: ui.color,
        buttonText: ui.buttonText,
        buttonVariant: ui.buttonVariant,
      };
    });
  }, [data]);

  return (
    <main className="min-h-screen bg-gray-50 p-4 sm:p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header Section */}
        <div className="flex items-start sm:items-center justify-between mb-6 sm:mb-8 gap-3 sm:gap-0 flex-col sm:flex-row">
          <div>
            <h1 className="text-2xl sm:text-3xl font-bold text-gray-900 mb-1 sm:mb-2">My Workspace</h1>
            <p className="text-gray-700">Track and manage your ongoing procedures</p>
          </div>
          
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 sm:gap-6 mb-6 sm:mb-8">
          {stats.map((stat, index) => (
            <Card
              key={index}
              className={`border-0 shadow-sm hover:shadow-md transition-all duration-300 hover:-translate-y-1 cursor-pointer ${statCardBackgrounds[index % statCardBackgrounds.length]}`}
            >
              <CardContent className="p-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-gray-700 mb-1">{stat.title}</p>
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
        <div className="flex flex-col sm:flex-row gap-3 sm:gap-4 mb-4 sm:mb-6">
          <div className="flex items-center gap-2">
            <span className="text-sm font-medium text-gray-700">Status:</span>
            <Select defaultValue="all">
              <SelectTrigger className="w-full sm:w-32">
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
          <div className="flex items-center gap-2">
            <span className="text-sm font-medium text-gray-700">Organization:</span>
            <Select defaultValue="all">
              <SelectTrigger className="w-full sm:w-48">
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

        {/* Procedure Cards */}
        <div className="space-y-4">
          {isLoading && <div className="text-sm text-gray-500">Loading procedures...</div>}
          {error && !isLoading && <div className="text-sm text-red-600">Failed to load procedures.</div>}
          {!isLoading && !error && apiProcedures.length === 0 && (
            <div className="text-sm text-gray-500">No procedures found.</div>
          )}

          {!isLoading && !error && apiProcedures.map((procedure, index) => (
            <Card
              key={procedure.id}
              className="border-0 shadow-sm hover:shadow-lg transition-all duration-300 hover:-translate-y-1 animate-in fade-in slide-in-from-bottom-4 bg-white"
              style={{ animationDelay: `${index * 100}ms` }}
            >
              <CardContent className="p-4 sm:p-6">
                <div className="flex items-start justify-between mb-4 gap-3 sm:gap-4 flex-col sm:flex-row">
                  <div className="flex items-start gap-4 w-full">
                    <div className="w-10 h-10 sm:w-12 sm:h-12 bg-blue-100 rounded-lg flex items-center justify-center transition-transform duration-200 hover:scale-110">
                      <FileText className="w-5 h-5 sm:w-6 sm:h-6 text-blue-600" />
                    </div>
                    <div className="min-w-0">
                      <h3 className="font-semibold text-base sm:text-lg text-gray-900 mb-1 hover:text-[#3A6A8D] transition-colors duration-200 truncate">
                        {procedure.title}
                      </h3>
                      <p className="text-gray-600 text-xs sm:text-sm">{procedure.department}</p>
                    </div>
                  </div>
                  <div className="flex items-center sm:items-center gap-3 w-full sm:w-auto justify-between sm:justify-end">
                    <Badge className={`${procedure.statusColor} transition-all duration-200`}>
                      {procedure.status}
                    </Badge>
                    <Button
                      variant={procedure.buttonVariant}
                      size="sm"
                      className={
                        procedure.buttonVariant === "default"
                          ? "bg-[#3A6A8D] hover:bg-[#2d5470] text-white transition-all duration-200 hover:scale-105 w-full sm:w-auto"
                          : "transition-all duration-200 w-full sm:w-auto"
                      }
                      onClick={() => router.push(`/user/checklist?userProcedureId=${encodeURIComponent(procedure.id)}`)}
                    >
                      {procedure.buttonText}
                    </Button>
                  </div>
                </div>

                <div className="mb-4">
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-sm font-medium text-gray-700">Progress</span>
                    <span className="text-xs sm:text-sm text-gray-600">{procedure.progress}% Complete</span>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2 overflow-hidden">
                    <div
                      className={`h-2 rounded-full transition-all duration-1000 ease-out ${
                        procedure.status === "Completed"
                          ? "bg-[#5E9C8D]"
                          : procedure.status === "In Progress"
                            ? "bg-[#FEF9C3]"
                            : "bg-gray-300"
                      }`}
                      style={{ width: `${procedure.progress}%`, animationDelay: `${index * 200 + 500}ms` }}
                    />
                  </div>
                </div>

                {/* Procedure Details */}
                <div className="flex flex-wrap items-center gap-3 sm:gap-6 text-xs sm:text-sm text-gray-600">
                  {procedure.startDate && (
                    <div className="flex items-center gap-1">
                      <Calendar className="w-4 h-4" />
                      <span>Started: {procedure.startDate}</span>
                    </div>
                  )}
                  {procedure.completedDate && (
                    <div className="flex items-center gap-1">
                      <CheckCircle className="w-4 h-4" />
                      <span>Completed: {procedure.completedDate}</span>
                    </div>
                  )}
                  {procedure.addedDate && (
                    <div className="flex items-center gap-1">
                      <Calendar className="w-4 h-4" />
                      <span>Added: {procedure.addedDate}</span>
                    </div>
                  )}
                  {procedure.estimatedCompletion && (
                    <div className="flex items-center gap-1">
                      <Clock className="w-4 h-4" />
                      <span>Est. completion: {procedure.estimatedCompletion}</span>
                    </div>
                  )}
                  {procedure.documentsUploaded && (
                    <div className="flex items-center gap-1">
                      <FileText className="w-4 h-4" />
                      <span>{procedure.documentsUploaded}</span>
                    </div>
                  )}
                  {procedure.requirements && (
                    <div className="flex items-center gap-1">
                      <CheckCircle className="w-4 h-4" />
                      <span>{procedure.requirements}</span>
                    </div>
                  )}
                  {procedure.readyToStart && (
                    <div className="flex items-center gap-1">
                      <Clock className="w-4 h-4" />
                      <span>{procedure.readyToStart}</span>
                    </div>
                  )}
                  {procedure.documentsRequired && (
                    <div className="flex items-center gap-1">
                      <FileText className="w-4 h-4" />
                      <span>{procedure.documentsRequired}</span>
                    </div>
                  )}
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    </main>
  )
}
  
