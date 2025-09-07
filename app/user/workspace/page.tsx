// "use client";

// import { Calendar, CheckCircle, Clock, FileText, PauseCircle } from "lucide-react"
// import { Button } from "@/components/ui/button"
// import { CardContent, Card } from "@/components/ui/card"
// import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
// import { Badge } from "@/components/ui/badge"
// import { useRouter } from "next/navigation"
// import { useMemo } from "react";

// export default function WorkspacePage() {
//   const router = useRouter();

//   // Your mock procedures data
//   const procedures = [
//     {
//       id: 1,
//       title: "New Passport Application",
//       department: "Immigration Department",
//       status: "In Progress",
//       progress: 60,
//       startDate: "Dec 15, 2024",
//       estimatedCompletion: "Jan 30, 2025",
//       documentsUploaded: "4/6 documents uploaded",
//       statusColor: "bg-orange-100 text-orange-800",
//       buttonText: "Continue",
//       buttonVariant: "default" as const,
//     },
//     {
//       id: 2,
//       title: "Driver's License Renewal",
//       department: "Road Authority",
//       status: "Completed",
//       progress: 100,
//       completedDate: "Dec 10, 2024",
//       requirements: "All requirements met",
//       statusColor: "bg-green-100 text-green-800",
//       buttonText: "View Details",
//       buttonVariant: "outline" as const,
//     },
//     {
//       id: 3,
//       title: "Business License Application",
//       department: "National Bank",
//       status: "Not Started",
//       progress: 0,
//       addedDate: "Dec 20, 2024",
//       readyToStart: "Ready to start",
//       documentsRequired: "0/5 documents uploaded",
//       statusColor: "bg-gray-100 text-gray-800",
//       buttonText: "Start Now",
//       buttonVariant: "default" as const,
//     },
//   ]

//   // Calculate stats from mock procedures
//   const stats = useMemo(() => {
//     const total = procedures.length;
//     const inProgress = procedures.filter(p => p.status === "In Progress").length;
//     const completed = procedures.filter(p => p.status === "Completed").length;
//     const notStarted = procedures.filter(p => p.status === "Not Started").length;

//     return [
//       { title: "Total Procedures", value: String(total), icon: FileText, color: "text-blue-600", bgColor: "bg-blue-50" },
//       { title: "In Progress", value: String(inProgress), icon: Clock, color: "text-orange-600", bgColor: "bg-orange-50" },
//       { title: "Completed", value: String(completed), icon: CheckCircle, color: "text-green-600", bgColor: "bg-green-50" },
//       { title: "Not Started", value: String(notStarted), icon: PauseCircle, color: "text-gray-600", bgColor: "bg-gray-50" },
//     ];
//   }, []);

//   return (
//     <main className="min-h-screen bg-gray-50 p-6">
//       <div className="max-w-7xl mx-auto">
//         {/* Header Section */}
//         <div className="flex items-center justify-between mb-8">
//           <div>
//             <h1 className="text-3xl font-bold text-gray-900 mb-2">My Workspace</h1>
//             <p className="text-neutral">Track and manage your ongoing procedures</p>
//           </div>
//         </div>

//         {/* Stats Cards */}
//         <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
//           {stats.map((stat, index) => (
//             <Card
//               key={index}
//               className="border-0 shadow-sm hover:shadow-md transition-all duration-300 hover:-translate-y-1 cursor-pointer bg-white"
//             >
//               <CardContent className="p-6">
//                 <div className="flex items-center justify-between">
//                   <div>
//                     <p className="text-sm font-medium text-neutral mb-1">{stat.title}</p>
//                     <p className="text-3xl font-bold text-gray-900">
//                       {stat.value}
//                     </p>
//                   </div>
//                   <div className={`p-3 rounded-lg ${stat.bgColor} transition-transform duration-200 hover:scale-110`}>
//                     <stat.icon className={`w-6 h-6 ${stat.color}`} />
//                   </div>
//                 </div>
//               </CardContent>
//             </Card>
//           ))}
//         </div>

//         {/* Filters */}
//         <div className="flex gap-4 mb-6">
//           <div className="flex items-center gap-2">
//             <span className="text-sm font-medium text-gray-700 ">Status:</span>
//             <Select defaultValue="all">
//               <SelectTrigger className="w-32">
//                 <SelectValue />
//               </SelectTrigger>
//               <SelectContent>
//                 <SelectItem value="all">All</SelectItem>
//                 <SelectItem value="in-progress">In Progress</SelectItem>
//                 <SelectItem value="completed">Completed</SelectItem>
//                 <SelectItem value="not-started">Not Started</SelectItem>
//               </SelectContent>
//             </Select>
//           </div>
//           <div className="flex items-center gap-2 ">
//             <span className="text-sm  font-medium text-gray-700">Organization:</span>
//             <Select defaultValue="all">
//               <SelectTrigger className="w-48">
//                 <SelectValue />
//               </SelectTrigger>
//               <SelectContent>
//                 <SelectItem value="all">All Organizations</SelectItem>
//                 <SelectItem value="immigration">Immigration Department</SelectItem>
//                 <SelectItem value="road">Road Authority</SelectItem>
//                 <SelectItem value="bank">National Bank</SelectItem>
//               </SelectContent>
//             </Select>
//           </div>
//         </div>

//         {/* Procedure Cards */}
//         <div className="space-y-4">
//           {procedures.map((procedure, index) => (
//             <Card
//               key={procedure.id}
//               className="border-0 shadow-sm hover:shadow-lg transition-all duration-300 hover:-translate-y-1 animate-in fade-in slide-in-from-bottom-4 bg-white"
//               style={{ animationDelay: `${index * 100}ms` }}
//             >
//               <CardContent className="p-6">
//                 <div className="flex items-start justify-between mb-4">
//                   <div className="flex items-start gap-4">
//                     <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center transition-transform duration-200 hover:scale-110">
//                       <FileText className="w-6 h-6 text-blue-600" />
//                     </div>
//                     <div>
//                       <h3 className="font-semibold text-lg text-gray-900 mb-1 hover:text-[#3A6A8D] transition-colors duration-200">
//                         {procedure.title}
//                       </h3>
//                       <p className="text-gray-600 text-sm">{procedure.department}</p>
//                     </div>
//                   </div>
//                   <div className="flex items-center gap-3">
//                     <Badge className={`${procedure.statusColor} transition-all duration-200`}>
//                       {procedure.status}
//                     </Badge>
//                     <Button
//                       variant={procedure.buttonVariant}
//                       size="sm"
//                       className={
//                         procedure.buttonVariant === "default"
//                           ? "bg-[#3A6A8D] hover:bg-[#2d5470] text-white transition-all duration-200 hover:scale-105"
//                           : "transition-all duration-200"
//                       }
//                       onClick={() => router.push("/user/checklist")}
//                     >
//                       {procedure.buttonText}
//                     </Button>
//                   </div>
//                 </div>

//                 <div className="mb-4">
//                   <div className="flex items-center justify-between mb-2">
//                     <span className="text-sm font-medium text-gray-700">Progress</span>
//                     <span className="text-sm text-gray-600">{procedure.progress}% Complete</span>
//                   </div>
//                   <div className="w-full bg-gray-200 rounded-full h-2 overflow-hidden">
//                     <div
//                       className={`h-2 rounded-full transition-all duration-1000 ease-out ${
//                         procedure.status === "Completed"
//                           ? "bg-[#5E9C8D]"
//                           : procedure.status === "In Progress"
//                             ? "bg-[#FEF9C3]"
//                             : "bg-gray-300"
//                       }`}
//                       style={{
//                         width: `${procedure.progress}%`,
//                         animationDelay: `${index * 200 + 500}ms`,
//                       }}
//                     />
//                   </div>
//                 </div>

//                 {/* Procedure Details */}
//                 <div className="flex items-center gap-6 text-sm text-gray-600">
//                   {procedure.startDate && (
//                     <div className="flex items-center gap-1">
//                       <Calendar className="w-4 h-4" />
//                       <span>Started: {procedure.startDate}</span>
//                     </div>
//                   )}
//                   {procedure.completedDate && (
//                     <div className="flex items-center gap-1">
//                       <CheckCircle className="w-4 h-4" />
//                       <span>Completed: {procedure.completedDate}</span>
//                     </div>
//                   )}
//                   {procedure.addedDate && (
//                     <div className="flex items-center gap-1">
//                       <Calendar className="w-4 h-4" />
//                       <span>Added: {procedure.addedDate}</span>
//                     </div>
//                   )}
//                   {procedure.estimatedCompletion && (
//                     <div className="flex items-center gap-1">
//                       <Clock className="w-4 h-4" />
//                       <span>Est. completion: {procedure.estimatedCompletion}</span>
//                     </div>
//                   )}
//                   {procedure.documentsUploaded && (
//                     <div className="flex items-center gap-1">
//                       <FileText className="w-4 h-4" />
//                       <span>{procedure.documentsUploaded}</span>
//                     </div>
//                   )}
//                   {procedure.requirements && (
//                     <div className="flex items-center gap-1">
//                       <CheckCircle className="w-4 h-4" />
//                       <span>{procedure.requirements}</span>
//                     </div>
//                   )}
//                   {procedure.readyToStart && (
//                     <div className="flex items-center gap-1">
//                       <Clock className="w-4 h-4" />
//                       <span>{procedure.readyToStart}</span>
//                     </div>
//                   )}
//                   {procedure.documentsRequired && (
//                     <div className="flex items-center gap-1">
//                       <FileText className="w-4 h-4" />
//                       <span>{procedure.documentsRequired}</span>
//                     </div>
//                   )}
//                 </div>
//               </CardContent>
//             </Card>
//           ))}
//         </div>
//       </div>
//     </main>
//   )
// }
"use client";

import { Calendar, CheckCircle, Clock, FileText, PauseCircle, AlertCircle } from "lucide-react"
import { Button } from "@/components/ui/button"
import { CardContent, Card } from "@/components/ui/card"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Badge } from "@/components/ui/badge"
import { useRouter } from "next/navigation"
import { useMemo, useEffect } from "react";
import { useGetMyProceduresQuery } from "@/app/store/slices/workspaceSlice";

// Define the expected response type
interface ProcedureItem {
  id: string;
  percent: number;
  procedure_id: string;
  status: string;
  updated_at: string;
  user_id: string;
}

export default function WorkspacePage() {
  const router = useRouter();
  const { data, error, isLoading, refetch } = useGetMyProceduresQuery({});

  // Debug: Check what's in localStorage
  useEffect(() => {
    if (typeof window !== "undefined") {
      console.log("LocalStorage tokens:", {
        accessToken: localStorage.getItem("accessToken"),
        access_token: localStorage.getItem("access_token"),
        token: localStorage.getItem("token"),
        authToken: localStorage.getItem("authToken")
      });
    }
  }, []);

  // Log detailed error information
  useEffect(() => {
    if (error) {
      console.error("Detailed error information:", {
        error,
        status: 'status' in error ? error.status : 'No status',
        data: 'data' in error ? error.data : 'No data',
        message: 'message' in error ? error.message : 'No message'
      });
    }
  }, [error]);

  // Function to handle retry with potential reauthentication
  const handleRetry = () => {
    // Try to get the latest token
    if (typeof window !== "undefined") {
      const token = localStorage.getItem("accessToken") || 
                   localStorage.getItem("access_token") || 
                   localStorage.getItem("token") || 
                   localStorage.getItem("authToken");
      
      if (!token) {
        // Redirect to login if no token found
        router.push("/login");
        return;
      }
      
      console.log("Retrying with token:", token ? "Found" : "Not found");
    }
    
    refetch();
  };

  // Transform API data to match the expected format
  const procedures = useMemo(() => {
    if (!data || !Array.isArray(data)) return [];
    
    return data.map((procedure: ProcedureItem) => ({
      id: procedure.id,
      title: `Procedure ${procedure.procedure_id || procedure.id.substring(0, 8)}`,
      department: getDepartmentFromId(procedure.procedure_id),
      status: getStatusText(procedure.status),
      progress: procedure.percent,
      startDate: formatDate(procedure.updated_at),
      estimatedCompletion: calculateEstimatedCompletion(procedure.updated_at, procedure.percent),
      documentsUploaded: `${Math.floor(procedure.percent / 20)}/5 documents uploaded`,
      statusColor: getStatusColor(procedure.status),
      buttonText: getButtonText(procedure.status),
      buttonVariant: getButtonVariant(procedure.status),
    }));
  }, [data]);

  // Helper functions (keep the same as before)
  const getStatusText = (status: string) => {
    switch (status) {
      case "completed": return "Completed";
      case "in-progress": return "In Progress";
      case "pending": return "Not Started";
      default: return "Not Started";
    }
  };

  const getDepartmentFromId = (id: string) => {
    if (!id) return "Government Department";
    if (id.includes("passport") || id.includes("immigration")) return "Immigration Department";
    if (id.includes("driver") || id.includes("license")) return "Road Authority";
    if (id.includes("business") || id.includes("bank")) return "National Bank";
    return "Government Department";
  };

  const formatDate = (dateString: string) => {
    if (!dateString) return "Not specified";
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };

  const calculateEstimatedCompletion = (startDate: string, percent: number) => {
    if (!startDate || percent >= 100) return "Completed";
    const start = new Date(startDate);
    const daysPassed = (new Date().getTime() - start.getTime()) / (1000 * 60 * 60 * 24);
    const totalDaysEstimated = daysPassed / (percent / 100);
    const daysRemaining = totalDaysEstimated - daysPassed;
    const completionDate = new Date();
    completionDate.setDate(completionDate.getDate() + daysRemaining);
    return formatDate(completionDate.toISOString());
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case "completed": return "bg-green-100 text-green-800";
      case "in-progress": return "bg-orange-100 text-orange-800";
      default: return "bg-gray-100 text-gray-800";
    }
  };

  const getButtonText = (status: string) => {
    switch (status) {
      case "completed": return "View Details";
      case "in-progress": return "Continue";
      default: return "Start Now";
    }
  };

  const getButtonVariant = (status: string) => {
    switch (status) {
      case "completed": return "outline" as const;
      default: return "default" as const;
    }
  };

  // Calculate stats from API data
  const stats = useMemo(() => {
    if (!data || !Array.isArray(data)) return [
      { title: "Total Procedures", value: "0", icon: FileText, color: "text-blue-600", bgColor: "bg-blue-50" },
      { title: "In Progress", value: "0", icon: Clock, color: "text-orange-600", bgColor: "bg-orange-50" },
      { title: "Completed", value: "0", icon: CheckCircle, color: "text-green-600", bgColor: "bg-green-50" },
      { title: "Not Started", value: "0", icon: PauseCircle, color: "text-gray-600", bgColor: "bg-gray-50" },
    ];
    
    const total = data.length;
    const inProgress = data.filter((p: ProcedureItem) => p.status === "in-progress").length;
    const completed = data.filter((p: ProcedureItem) => p.status === "completed").length;
    const notStarted = data.filter((p: ProcedureItem) => 
      !p.status || p.status === "pending" || p.status === "not-started"
    ).length;

    return [
      { title: "Total Procedures", value: String(total), icon: FileText, color: "text-blue-600", bgColor: "bg-blue-50" },
      { title: "In Progress", value: String(inProgress), icon: Clock, color: "text-orange-600", bgColor: "bg-orange-50" },
      { title: "Completed", value: String(completed), icon: CheckCircle, color: "text-green-600", bgColor: "bg-green-50" },
      { title: "Not Started", value: String(notStarted), icon: PauseCircle, color: "text-gray-600", bgColor: "bg-gray-50" },
    ];
  }, [data]);

  if (isLoading) {
    return (
      <main className="min-h-screen bg-gray-50 p-6">
        <div className="max-w-7xl mx-auto">
          <div className="flex items-center justify-between mb-8">
            <div>
              <h1 className="text-3xl font-bold text-gray-900 mb-2">My Workspace</h1>
              <p className="text-neutral">Loading your procedures...</p>
            </div>
          </div>
        </div>
      </main>
    );
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
        {error && (
          <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
            <div className="flex items-center gap-2 text-red-800 mb-2">
              <AlertCircle className="w-5 h-5" />
              <h3 className="font-semibold">Error loading procedures</h3>
            </div>
            <p className="text-red-700 mb-4">
              {('status' in error && error.status === 401) 
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

        {/* Rest of your component remains the same */}
        {/* ... */}
      </div>
    </main>
  );
}