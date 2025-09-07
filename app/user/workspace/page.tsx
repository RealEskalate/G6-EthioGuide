"use client";

import { Calendar, CheckCircle, Clock, FileText, PauseCircle, AlertCircle } from "lucide-react"
import { Button } from "@/components/ui/button"
import { CardContent, Card } from "@/components/ui/card"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Badge } from "@/components/ui/badge"
import { useRouter } from "next/navigation"
import { useSession } from "next-auth/react"
import { useGetMyChecklistsQuery } from "@/app/store/slices/checklistsApi"
import { useListProceduresQuery } from "@/app/store/slices/proceduresApi"
import { useMemo } from "react"

export default function WorkspacePage() {
  const router = useRouter();
  const { data: session } = useSession();
  const { data: myChecklists, isLoading: loadingChecklists, error: checklistsError } = useGetMyChecklistsQuery(
    { token: session?.accessToken || undefined },
    { skip: !session?.accessToken }
  );
  const { data: proceduresData } = useListProceduresQuery({ page: 1, limit: 100 });

  // Calculate stats from real data
  const stats = useMemo(() => {
    const total = myChecklists?.length || 0;
    const inProgress = myChecklists?.filter(c => c.status === 'IN_PROGRESS').length || 0;
    const completed = myChecklists?.filter(c => c.status === 'COMPLETED').length || 0;
    
    return [
      {
        title: "Total Procedures",
        value: total.toString(),
        icon: FileText,
        color: "text-blue-600",
        bgColor: "bg-blue-50",
      },
      {
        title: "In Progress",
        value: inProgress.toString(),
        icon: Clock,
        color: "text-orange-600",
        bgColor: "bg-orange-50",
      },
      {
        title: "Completed",
        value: completed.toString(),
        icon: CheckCircle,
        color: "text-green-600",
        bgColor: "bg-green-50",
      },
      {
        title: "Available Procedures",
        value: (proceduresData?.list?.length || 0).toString(),
        icon: FileText,
        color: "text-purple-600",
        bgColor: "bg-purple-50",
      },
    ];
  }, [myChecklists, proceduresData]);

  // Transform real checklist data into UI format
  const procedures = useMemo(() => {
    if (!myChecklists) return [];
    
    return myChecklists.map((checklist) => {
      const procedure = proceduresData?.list?.find(p => p.id === checklist.procedureId);
      const progress = checklist.progress || 0;
      
      let status = "Not Started";
      let statusColor = "bg-gray-100 text-gray-800";
      let buttonText = "Start Now";
      let buttonVariant: "default" | "outline" = "default";
      
      if (checklist.status === 'COMPLETED') {
        status = "Completed";
        statusColor = "bg-green-100 text-green-800";
        buttonText = "View Details";
        buttonVariant = "outline";
      } else if (checklist.status === 'IN_PROGRESS') {
        status = "In Progress";
        statusColor = "bg-orange-100 text-orange-800";
        buttonText = "Continue";
        buttonVariant = "default";
      }
      
      const completedItems = checklist.items?.filter(item => item.is_checked).length || 0;
      const totalItems = checklist.items?.length || 0;
      
      return {
        id: checklist.id,
        title: procedure?.title || procedure?.name || `Procedure ${checklist.procedureId}`,
        department: "Immigration Department",
        status,
        progress,
        startDate: checklist.createdAt ? new Date(checklist.createdAt).toLocaleDateString() : "Recently",
        estimatedCompletion: checklist.status === 'COMPLETED' ? "Completed" : "Ongoing",
        documentsUploaded: `${completedItems}/${totalItems} documents uploaded`,
        statusColor,
        buttonText,
        buttonVariant,
      };
    });
  }, [myChecklists, proceduresData]);

  // Add handleRetry to reload the page
  function handleRetry() {
    router.refresh?.() // Next.js 13+
    // or fallback to window.location.reload()
    if (!router.refresh) window.location.reload()
  }

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

        {/* Error Display */}
        {checklistsError && (
          <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
            <div className="flex items-center gap-2 text-red-800 mb-2">
              <AlertCircle className="w-5 h-5" />
              <h3 className="font-semibold">Error loading procedures</h3>
            </div>
            <p className="text-red-700 mb-4">
              {('status' in checklistsError && checklistsError.status === 401) 
                ? "Authentication failed. Please check your login credentials."
                : "There was a problem loading your procedures. Please try again."
              }
            </p>
            <div className="flex gap-3">
              <Button 
                onClick={handleRetry}
                variant="outline" 
                className="border-red-300 text-red-700 hover:bg-red-50"
              >
                Try Again
              </Button>
              <Button 
                onClick={() => router.push("/login")}
                className="bg-red-600 hover:bg-red-700"
              >
                Go to Login
              </Button>
            </div>
          </div>
        )}

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
                    <p className="text-3xl font-bold text-gray-900">
                      {stat.value}
                    </p>
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

            {/* Procedure Cards */}
            <div className="space-y-4">
              {loadingChecklists ? (
                <div className="text-center py-8">
                  <div className="text-gray-600">Loading your procedures...</div>
                </div>
              ) : checklistsError ? (
                <div className="text-center py-8">
                  <div className="text-red-600">Failed to load procedures. Please try again.</div>
                </div>
              ) : procedures.length === 0 ? (
                <div className="text-center py-8">
                  <div className="text-gray-600">No procedures found. Start by creating a checklist from the procedures list.</div>
                  <Button 
                    className="mt-4 bg-[#3A6A8D] hover:bg-[#2d5470] text-white"
                    onClick={() => router.push('/user/procedures-list')}
                  >
                    Browse Procedures
                  </Button>
                </div>
              ) : (
                procedures.map((procedure, index) => (
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
                          <p className="text-gray-600 text-sm">{procedure.department}</p>
                        </div>
                      </div>
                      <div className="flex items-center gap-3">
                        <Badge className={`${procedure.statusColor} transition-all duration-200`}>
                          {procedure.status}
                        </Badge>
                        <Button
                          variant={procedure.buttonVariant}
                          size="sm"
                          className={
                            procedure.buttonVariant === "default"
                              ? "bg-[#3A6A8D] hover:bg-[#2d5470] text-white transition-all duration-200 hover:scale-105"
                              : "transition-all duration-200"
                          }
                          onClick={() => router.push(`/user/checklist/${encodeURIComponent(procedure.id)}`)}
                        >
                          {procedure.buttonText}
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
                            procedure.status === "Completed"
                              ? "bg-[#5E9C8D]"
                              : procedure.status === "In Progress"
                                ? "bg-[#FEF9C3]"
                                : "bg-gray-300"
                          }`}
                          style={{
                            width: `${procedure.progress}%`,
                            animationDelay: `${index * 200 + 500}ms`,
                          }}
                        />
                      </div>
                    </div>

                    {/* Procedure Details */}
                    <div className="flex items-center gap-6 text-sm text-gray-600">
                      {procedure.startDate && (
                        <div className="flex items-center gap-1">
                          <Calendar className="w-4 h-4" />
                          <span>Started: {procedure.startDate}</span>
                        </div>
                      )}
                      {/* removed unsupported fields: completedDate, addedDate */}
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
                      {/* removed unsupported fields: requirements, readyToStart, documentsRequired */}
                    </div>
                  </CardContent>
                </Card>
                ))
              )}
            </div>
          </div>
        </main>
     
  )
}