// "use client";

// import { Calendar, CheckCircle, Clock, FileText, AlertCircle, Search } from "lucide-react"
// import { Button } from "@/components/ui/button"
// import { CardContent, Card } from "@/components/ui/card"
// import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
// import { Badge } from "@/components/ui/badge"
// import { useRouter } from "next/navigation"
// import { useSession } from "next-auth/react"
// import { useGetMyChecklistsQuery } from "@/app/store/slices/checklistsApi"
// import { useListProceduresQuery } from "@/app/store/slices/proceduresApi"
// import { useMemo, useState, useEffect } from "react"
// import { motion, useReducedMotion } from "framer-motion"
// import Pagination from "@/components/shared/pagination"
// import { useLazyGetOrgQuery } from "@/app/store/slices/orgsApi"

// export default function WorkspacePage() {
//   const router = useRouter();
//   const { data: session } = useSession();
//   const { data: myChecklists, isLoading: loadingChecklists, error: checklistsError } = useGetMyChecklistsQuery(
//     { token: session?.accessToken || undefined },
//     { skip: !session?.accessToken }
//   );
//   const { data: proceduresData } = useListProceduresQuery({ page: 1, limit: 100 });
//   const [triggerGetOrg] = useLazyGetOrgQuery();

//   // Search & org filter state
//   const [searchInput, setSearchInput] = useState("");
//   const [searchQuery, setSearchQuery] = useState("");
//   const [orgFilter, setOrgFilter] = useState<string>("all");
//   const [orgNameMap, setOrgNameMap] = useState<Record<string, string>>({});

//   const prefersReducedMotion = useReducedMotion();
//   const listItemVariants = prefersReducedMotion
//     ? { hidden: { opacity: 0 }, visible: { opacity: 1 } }
//     : {
//         hidden: { opacity: 0, y: 10, scale: 0.98, filter: "blur(0.2px)" },
//         visible: {
//           opacity: 1,
//           y: 0,
//           scale: 1,
//           filter: "none",
//           transition: { type: "spring" as const, stiffness: 220, damping: 18, mass: 0.9 },
//         },
//       }

//   // Calculate stats from real data
//   const stats = useMemo(() => {
//     const total = myChecklists?.length || 0;
//     const inProgress = myChecklists?.filter(c => c.status === 'IN_PROGRESS').length || 0;
//     const completed = myChecklists?.filter(c => c.status === 'COMPLETED').length || 0;
    
//     return [
//       {
//         title: "Total Procedures",
//         value: total.toString(),
//         icon: FileText,
//         color: "text-blue-600",
//         bgColor: "bg-blue-50",
//       },
//       {
//         title: "In Progress",
//         value: inProgress.toString(),
//         icon: Clock,
//         color: "text-orange-600",
//         bgColor: "bg-orange-50",
//       },
//       {
//         title: "Completed",
//         value: completed.toString(),
//         icon: CheckCircle,
//         color: "text-green-600",
//         bgColor: "bg-green-50",
//       },
//       {
//         title: "Available Procedures",
//         value: (proceduresData?.list?.length || 0).toString(),
//         icon: FileText,
//         color: "text-purple-600",
//         bgColor: "bg-purple-50",
//       },
//     ];
//   }, [myChecklists, proceduresData]);

//   // Transform real checklist data into UI format
//   const procedures = useMemo(() => {
//     if (!myChecklists) return [];
    
//     return myChecklists.map((checklist) => {
//       const procedure = proceduresData?.list?.find(p => p.id === checklist.procedureId);
//   const pAny = (procedure || {}) as Record<string, unknown>;
//   const organizationId: string | undefined = (pAny["organizationId"] as string | undefined) || (pAny["OrganizationID"] as string | undefined) || (pAny["organization_id"] as string | undefined);
//       const organizationName = organizationId ? (orgNameMap[organizationId] || "Organization") : "Organization";
//       const progress = checklist.progress || 0;
      
//       let status = "Not Started";
//       let statusColor = "bg-gray-100 text-gray-800";
//       let buttonText = "Start Now";
//       let buttonVariant: "default" | "outline" = "default";
      
//       if (checklist.status === 'COMPLETED') {
//         status = "Completed";
//         statusColor = "bg-green-100 text-green-800";
//         buttonText = "View Details";
//         buttonVariant = "outline";
//       } else if (checklist.status === 'IN_PROGRESS') {
//         status = "In Progress";
//         statusColor = "bg-orange-100 text-orange-800";
//         buttonText = "Continue";
//         buttonVariant = "default";
//       }
      
//       const completedItems = checklist.items?.filter(item => item.is_checked).length || 0;
//       const totalItems = checklist.items?.length || 0;
      
//       return {
//         id: checklist.id,
//         title: procedure?.title || procedure?.name || `Procedure ${checklist.procedureId}`,
//         department: organizationName,
//         status,
//         progress,
//         startDate: checklist.createdAt ? new Date(checklist.createdAt).toLocaleDateString() : "Recently",
//         estimatedCompletion: checklist.status === 'COMPLETED' ? "Completed" : "Ongoing",
//         documentsUploaded: `${completedItems}/${totalItems} documents uploaded`,
//         statusColor,
//         buttonText,
//         buttonVariant,
//         organizationId,
//         organizationName,
//       };
//     });
//   }, [myChecklists, proceduresData, orgNameMap]);

//   // Fetch organization names lazily for unique IDs from saved procedures
//   useEffect(() => {
//     if (!myChecklists || !proceduresData?.list) return;
//     const ids = new Set<string>();
//     for (const cl of myChecklists) {
//       const p = proceduresData.list.find((pp) => pp.id === cl.procedureId) as Record<string, unknown> | undefined;
//       const oid: string | undefined = (p?.["organizationId"] as string | undefined) || (p?.["OrganizationID"] as string | undefined) || (p?.["organization_id"] as string | undefined);
//       if (oid) ids.add(oid);
//     }
//     const missing = Array.from(ids).filter((id) => !(id in orgNameMap));
//     if (missing.length === 0) return;
//     (async () => {
//       const entries: Array<[string, string]> = [];
//       for (const id of missing) {
//         try {
//           const res = await triggerGetOrg(id).unwrap();
//           entries.push([id, res.name]);
//         } catch {
//           // noop on failure
//         }
//       }
//       if (entries.length) setOrgNameMap((prev) => ({ ...prev, ...Object.fromEntries(entries) }));
//     })();
//   }, [myChecklists, proceduresData, orgNameMap, triggerGetOrg]);

//   // Build filtered list (apply search + organization filter)
//   const filteredProcedures = useMemo(() => {
//     const q = searchQuery.trim().toLowerCase();
//     return procedures.filter((p) => {
//       const matchesSearch = !q ||
//         p.title.toLowerCase().includes(q) ||
//         (p.organizationName || "").toLowerCase().includes(q) ||
//         (p.documentsUploaded || "").toLowerCase().includes(q);
//       const matchesOrg = orgFilter === "all" || (p.organizationName === orgFilter);
//       return matchesSearch && matchesOrg;
//     });
//   }, [procedures, searchQuery, orgFilter]);

//   // Pagination state (10 per page)
//   const PAGE_SIZE = 10;
//   const [page, setPage] = useState(1);
//   const totalPages = Math.max(1, Math.ceil(filteredProcedures.length / PAGE_SIZE));
//   useEffect(() => {
//     // Clamp page if procedures length changes or becomes smaller
//     if (page > totalPages) setPage(totalPages);
//   }, [totalPages, page]);
//   const paginatedProcedures = useMemo(() => {
//     const start = (page - 1) * PAGE_SIZE;
//     return filteredProcedures.slice(start, start + PAGE_SIZE);
//   }, [filteredProcedures, page]);

//   // Add handleRetry to reload the page
//   function handleRetry() {
//     router.refresh?.() // Next.js 13+
//     // or fallback to window.location.reload()
//     if (!router.refresh) window.location.reload()
//   }

//   return (
//     <div className="min-h-screen w-full bg-gray-50 relative overflow-hidden">
//       <div className="absolute inset-0 overflow-hidden pointer-events-none">
//         <div className="absolute -top-24 -right-24 w-56 h-56 rounded-full blur-3xl" style={{ background: 'radial-gradient(closest-side, rgba(167,179,185,0.10), rgba(167,179,185,0))' }}></div>
//         <div className="absolute -bottom-28 -left-28 w-64 h-64 rounded-full blur-3xl" style={{ background: 'radial-gradient(closest-side, rgba(94,156,141,0.10), rgba(94,156,141,0))' }}></div>
//       </div>

//       <main className="relative z-10 p-4 sm:p-6 md:p-8">
//         <div className="max-w-7xl mx-auto">
//           {/* Header Section */}
//           <div className="mb-6 md:mb-8 w-full">
//             <div className="bg-white/90 backdrop-blur-sm rounded-2xl border border-[#a7b3b9]/30 p-4 sm:p-6 md:p-8 shadow-lg">
//               <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
//                 <div>
//                   <h1 className="text-2xl sm:text-3xl font-bold text-[#2e4d57] mb-2">My Workspace</h1>
//                   <p className="text-[#1c3b2e] text-sm sm:text-base">Track and manage your ongoing procedures</p>
//                 </div>
//                 <div className="inline-flex items-center gap-2 bg-[#3a6a8d]/10 backdrop-blur-sm border border-[#3a6a8d]/30 rounded-full px-3 sm:px-4 py-2">
//                   <CheckCircle className="w-4 h-4 text-[#3a6a8d] animate-pulse" />
//                   <span className="text-xs sm:text-sm font-medium text-[#2e4d57]">Workspace</span>
//                 </div>
//               </div>
//             </div>
//           </div>

//         {/* Error Display */}
//         {checklistsError && (
//           <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
//             <div className="flex items-center gap-2 text-red-800 mb-2">
//               <AlertCircle className="w-5 h-5" />
//               <h3 className="font-semibold">Error loading procedures</h3>
//             </div>
//             <p className="text-red-700 mb-4">
//               {('status' in checklistsError && checklistsError.status === 401) 
//                 ? "Authentication failed. Please check your login credentials."
//                 : "There was a problem loading your procedures. Please try again."
//               }
//             </p>
//             <div className="flex gap-3">
//               <Button 
//                 onClick={handleRetry}
//                 variant="outline" 
//                 className="border-red-300 text-red-700 hover:bg-red-50"
//               >
//                 Try Again
//               </Button>
//               <Button 
//                 onClick={() => router.push("/login")}
//                 className="bg-red-600 hover:bg-red-700"
//               >
//                 Go to Login
//               </Button>
//             </div>
//           </div>
//         )}

//         {/* Stats Cards */}
//         <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
//           {stats.map((stat, index) => (
//             <motion.div key={index} variants={listItemVariants} initial="hidden" animate="visible" transition={{ delay: index * 0.05 + 0.05 }}>
//               <Card className="bg-white/80 backdrop-blur-md rounded-2xl border border-[#e5e7eb] shadow-xl hover:shadow-2xl transition-all duration-700 hover:scale-105 hover:-translate-y-2 cursor-pointer relative overflow-hidden group">
//                 <div className="absolute inset-0 bg-gradient-to-r from-[#3a6a8d]/10 via-transparent to-[#5e9c8d]/10 opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
//                 <CardContent className="p-6 relative z-10">
//                   <div className="flex items-center justify-between">
//                     <div>
//                       <p className="text-sm font-medium text-[#4b5563] mb-1">{stat.title}</p>
//                       <p className="text-3xl font-bold text-[#111827]">{stat.value}</p>
//                     </div>
//                     <div className="w-10 h-10 rounded-lg flex items-center justify-center transition-all duration-500 group-hover:scale-110 group-hover:rotate-6" style={{ backgroundColor: '#e6f0f5' }}>
//                       <stat.icon className="w-6 h-6" style={{ color: '#3a6a8d' }} />
//                     </div>
//                   </div>
//                 </CardContent>
//               </Card>
//             </motion.div>
//           ))}
//         </div>

//             {/* Search & Organization Filter */}
//             <div className="bg-white/90 backdrop-blur-sm rounded-2xl border border-[#e5e7eb] p-4 sm:p-5 mb-6 shadow-md flex flex-col gap-4 sm:flex-row sm:items-center">
//               <div className="relative flex-1 flex min-w-[220px]">
//                 <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
//                 <input
//                   type="text"
//                   placeholder="Search your procedures..."
//                   value={searchInput}
//                   onChange={(e) => setSearchInput(e.target.value)}
//                   className="w-full pl-10 pr-4 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#3A6A8D] focus:border-transparent"
//                 />
//                 <Button
//                   type="button"
//                   className="ml-2 px-4 py-2 bg-[#3A6A8D] hover:bg-[#2d5470] text-white"
//                   onClick={() => { setSearchQuery(searchInput); setPage(1); }}
//                 >
//                   Search
//                 </Button>
//               </div>
//               <div className="flex flex-none w-full sm:w-64">
//                 <Select value={orgFilter} onValueChange={(v) => { setOrgFilter(v); setPage(1); }}>
//                   <SelectTrigger className="w-full border-[#3A6A8D] text-[#3A6A8D] hover:bg-[#3A6A8D]/5 focus:ring-2 focus:ring-[#3A6A8D] focus:border-transparent">
//                     <SelectValue placeholder="All Organizations" />
//                   </SelectTrigger>
//                   <SelectContent>
//                     <SelectItem value="all">All Organizations</SelectItem>
//                     {Array.from(new Set(procedures.map((p) => p.organizationName).filter(Boolean))).map((name) => (
//                       <SelectItem key={name as string} value={name as string}>{name as string}</SelectItem>
//                     ))}
//                   </SelectContent>
//                 </Select>
//               </div>
//             </div>

//             {/* Procedure Cards */}
//             <div className="space-y-4">
//               {loadingChecklists ? (
//                 <div className="text-center py-8">
//                   <div className="text-gray-600">Loading your procedures...</div>
//                 </div>
//               ) : checklistsError ? (
//                 <div className="text-center py-8">
//                   <div className="text-red-600">Failed to load procedures. Please try again.</div>
//                 </div>
//               ) : procedures.length === 0 ? (
//                 <div className="text-center py-8">
//                   <div className="text-gray-600">No procedures found. Start by creating a checklist from the procedures list.</div>
//                   <Button 
//                     className="mt-4 bg-[#3A6A8D] hover:bg-[#2d5470] text-white"
//                     onClick={() => router.push('/user/procedures-list')}
//                   >
//                     Browse Procedures
//                   </Button>
//                 </div>
//               ) : (
//                 paginatedProcedures.map((procedure, index) => (
//                 <motion.div key={procedure.id} variants={listItemVariants} initial="hidden" animate="visible" transition={{ delay: index * 0.06 + 0.12 }}>
//                   <Card className="bg-white/80 backdrop-blur-md rounded-2xl border border-[#e5e7eb] shadow-xl hover:shadow-2xl transition-all duration-700 hover:scale-105 hover:-translate-y-2 relative overflow-hidden group">
//                     <div className="absolute inset-0 bg-gradient-to-r from-[#3a6a8d]/10 via-transparent to-[#5e9c8d]/10 opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
//                     <CardContent className="p-6 relative z-10">
//                     <div className="flex items-start justify-between mb-4">
//                       <div className="flex items-start gap-4">
//                         <div className="w-12 h-12 rounded-xl flex items-center justify-center transition-all duration-500 group-hover:scale-110 group-hover:rotate-3" style={{ backgroundColor: '#e6f0f5' }}>
//                           <FileText className="w-6 h-6" style={{ color: '#3a6a8d' }} />
//                         </div>
//                         <div>
//                           <h3 className="font-semibold text-lg text-[#111827] mb-1 group-hover:text-[#3a6a8d] transition-colors duration-300">
//                             {procedure.title}
//                           </h3>
//                           <p className="text-[#4b5563] text-sm">{procedure.department}</p>
//                         </div>
//                       </div>
//                       <div className="flex items-center gap-3">
//                         <Badge className={`${procedure.statusColor} border border-[#e5e7eb] transition-all duration-200`}>
//                           {procedure.status}
//                         </Badge>
//                         <Button
//                           variant={procedure.buttonVariant}
//                           size="sm"
//                           className={
//                             procedure.buttonVariant === "default"
//                               ? "bg-gradient-to-r from-[#3a6a8d] to-[#2e4d57] hover:from-[#2e4d57] hover:to-[#1c3b2e] text-white transition-all duration-500 hover:scale-105 hover:shadow-lg"
//                               : "transition-all duration-300"
//                           }
//                           onClick={() => router.push(`/user/checklist/${encodeURIComponent(procedure.id)}`)}
//                         >
//                           {procedure.buttonText}
//                         </Button>
//                       </div>
//                     </div>

//                     <div className="mb-4">
//                       <div className="flex items-center justify-between mb-2">
//                         <span className="text-sm font-medium text-[#111827]">Progress</span>
//                         <span className="text-sm text-[#4b5563]">{procedure.progress}% Complete</span>
//                       </div>
//                       <div className="w-full bg-[#e5e7eb] rounded-full h-2 overflow-hidden">
//                         {(() => {
//                           const barColor =
//                             procedure.status === "Completed"
//                               ? "bg-[#5e9c8d]"
//                               : procedure.status === "In Progress"
//                                 ? "bg-gradient-to-r from-[#3a6a8d] to-[#2e4d57]"
//                                 : "bg-[#a7b3b9]/50";
//                           return (
//                             <motion.div
//                               className={`h-2 rounded-full ${barColor}`}
//                               initial={{ width: 0 }}
//                               animate={{ width: `${procedure.progress}%` }}
//                               transition={{ type: "spring", stiffness: 180, damping: 24, mass: 0.9 }}
//                             />
//                           );
//                         })()}
//                       </div>
//                     </div>

//                     {/* Procedure Details */}
//                     <div className="flex items-center gap-6 text-sm text-[#4b5563]">
//                       {procedure.startDate && (
//                         <div className="flex items-center gap-1">
//                           <Calendar className="w-4 h-4" />
//                           <span>Started: {procedure.startDate}</span>
//                         </div>
//                       )}
//                       {/* removed unsupported fields: completedDate, addedDate */}
//                       {procedure.estimatedCompletion && (
//                         <div className="flex items-center gap-1">
//                           <Clock className="w-4 h-4" />
//                           <span>Est. completion: {procedure.estimatedCompletion}</span>
//                         </div>
//                       )}
//                       {procedure.documentsUploaded && (
//                         <div className="flex items-center gap-1">
//                           <FileText className="w-4 h-4" />
//                           <span>{procedure.documentsUploaded}</span>
//                         </div>
//                       )}
//                       {/* removed unsupported fields: requirements, readyToStart, documentsRequired */}
//                     </div>
//                   </CardContent>
//                 </Card>
//                 </motion.div>
//                 ))
//               )}
//             </div>

//             {/* Pagination */}
//             {totalPages > 1 && (
//               <div className="mt-6">
//                 <Pagination page={page} totalPages={totalPages} onPageChange={setPage} />
//               </div>
//             )}
//           </div>
//         </main>
//     </div>
//   )
// }

"use client";

import { Calendar, CheckCircle, Clock, FileText, AlertCircle, Search } from "lucide-react"
import { Button } from "@/components/ui/button"
import { CardContent, Card } from "@/components/ui/card"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Badge } from "@/components/ui/badge"
import { useRouter } from "next/navigation"
import { useSession } from "next-auth/react"
import { useGetMyChecklistsQuery } from "@/app/store/slices/checklistsApi"
import { useListProceduresQuery } from "@/app/store/slices/proceduresApi"
import { useMemo, useState, useEffect } from "react"
import { motion, useReducedMotion } from "framer-motion"
import Pagination from "@/components/shared/pagination"
import { useLazyGetOrgQuery } from "@/app/store/slices/orgsApi"
import { useTranslation } from "react-i18next"

export default function WorkspacePage() {
  const { t } = useTranslation("user")
  const router = useRouter();
  const { data: session } = useSession();
  const { data: myChecklists, isLoading: loadingChecklists, error: checklistsError } = useGetMyChecklistsQuery(
    { token: session?.accessToken || undefined },
    { skip: !session?.accessToken }
  );
  const { data: proceduresData } = useListProceduresQuery({ page: 1, limit: 100 });
  const [triggerGetOrg] = useLazyGetOrgQuery();

  // Search & org filter state
  const [searchInput, setSearchInput] = useState("");
  const [searchQuery, setSearchQuery] = useState("");
  const [orgFilter, setOrgFilter] = useState<string>("all");
  const [orgNameMap, setOrgNameMap] = useState<Record<string, string>>({});

  const prefersReducedMotion = useReducedMotion();
  const listItemVariants = prefersReducedMotion
    ? { hidden: { opacity: 0 }, visible: { opacity: 1 } }
    : {
        hidden: { opacity: 0, y: 10, scale: 0.98, filter: "blur(0.2px)" },
        visible: {
          opacity: 1,
          y: 0,
          scale: 1,
          filter: "none",
          transition: { type: "spring" as const, stiffness: 220, damping: 18, mass: 0.9 },
        },
      }

  // Calculate stats from real data
  const stats = useMemo(() => {
    const total = Array.isArray(myChecklists) ? myChecklists.length : 0;
    const inProgress = Array.isArray(myChecklists) ? myChecklists.filter(c => c.status === 'IN_PROGRESS').length : 0;
    const completed = Array.isArray(myChecklists) ? myChecklists.filter(c => c.status === 'COMPLETED').length : 0;
    
    return [
      {
        title: t("workspace.stats.total"),
        value: total.toString(),
        icon: FileText,
        color: "text-blue-600",
        bgColor: "bg-blue-50",
      },
      {
        title: t("workspace.stats.in_progress"),
        value: inProgress.toString(),
        icon: Clock,
        color: "text-orange-600",
        bgColor: "bg-orange-50",
      },
      {
        title: t("workspace.stats.completed"),
        value: completed.toString(),
        icon: CheckCircle,
        color: "text-green-600",
        bgColor: "bg-green-50",
      },
      {
        title: t("workspace.stats.available"),
        value: (Array.isArray(proceduresData?.list) ? proceduresData.list.length : 0).toString(),
        icon: FileText,
        color: "text-purple-600",
        bgColor: "bg-purple-50",
      },
    ];
  }, [myChecklists, proceduresData, t]);

  // Transform real checklist data into UI format
  const procedures = useMemo(() => {
    if (!Array.isArray(myChecklists)) return [];
    
    return myChecklists.map((checklist, idx) => {
      if (!checklist) {
        console.error("Checklist is undefined at index:", idx);
        return null;
      }
      const procedure = Array.isArray(proceduresData?.list) ? proceduresData.list.find(p => p.id === checklist.procedureId) : undefined;
      const pAny = (procedure || {}) as Record<string, unknown>;
      const organizationId: string | undefined = (pAny["organizationId"] as string | undefined) || (pAny["OrganizationID"] as string | undefined) || (pAny["organization_id"] as string | undefined);
      const organizationName = organizationId ? (orgNameMap[organizationId] || t("workspace.default_organization")) : t("workspace.default_organization");
      const progress = checklist.progress || 0;
      
      let status = t("workspace.status.not_started");
      let statusColor = "bg-gray-100 text-gray-800";
      let buttonText = t("workspace.actions.start");
      let buttonVariant: "default" | "outline" = "default";
      
      if (checklist.status === 'COMPLETED') {
        status = t("workspace.status.completed");
        statusColor = "bg-green-100 text-green-800";
        buttonText = t("workspace.actions.view_details");
        buttonVariant = "outline";
      } else if (checklist.status === 'IN_PROGRESS') {
        status = t("workspace.status.in_progress");
        statusColor = "bg-orange-100 text-orange-800";
        buttonText = t("workspace.actions.continue");
        buttonVariant = "default";
      }
      
      const completedItems = Array.isArray(checklist.items) ? checklist.items.filter(item => item.is_checked).length : 0;
      const totalItems = Array.isArray(checklist.items) ? checklist.items.length : 0;
      
      return {
        id: checklist.id,
        title: procedure?.title || procedure?.name || t("workspace.default_procedure", { id: checklist.procedureId }),
        department: organizationName,
        status,
        progress,
        startDate: checklist.createdAt ? new Date(checklist.createdAt).toLocaleDateString() : t("workspace.default_date"),
        estimatedCompletion: checklist.status === 'COMPLETED' ? t("workspace.status.completed") : t("workspace.status.ongoing"),
        documentsUploaded: t("workspace.documents_uploaded", { completed: completedItems, total: totalItems }),
        statusColor,
        buttonText,
        buttonVariant,
        organizationId,
        organizationName,
      };
    }).filter((p): p is NonNullable<typeof p> => p !== null);
  }, [myChecklists, proceduresData, orgNameMap, t]);

  // Fetch organization names lazily for unique IDs from saved procedures
  useEffect(() => {
    if (!Array.isArray(myChecklists) || !Array.isArray(proceduresData?.list)) return;
    const ids = new Set<string>();
    for (const cl of myChecklists) {
      const p = proceduresData.list.find((pp) => pp.id === cl.procedureId) as Record<string, unknown> | undefined;
      const oid: string | undefined = (p?.["organizationId"] as string | undefined) || (p?.["OrganizationID"] as string | undefined) || (p?.["organization_id"] as string | undefined);
      if (oid) ids.add(oid);
    }
    const missing = Array.from(ids).filter((id) => !(id in orgNameMap));
    if (missing.length === 0) return;
    (async () => {
      const entries: Array<[string, string]> = [];
      for (const id of missing) {
        try {
          const res = await triggerGetOrg(id).unwrap();
          entries.push([id, res.name]);
        } catch {
          // noop on failure
        }
      }
      if (entries.length) setOrgNameMap((prev) => ({ ...prev, ...Object.fromEntries(entries) }));
    })();
  }, [myChecklists, proceduresData, orgNameMap, triggerGetOrg]);

  // Build filtered list (apply search + organization filter)
  const filteredProcedures = useMemo(() => {
    const q = searchQuery.trim().toLowerCase();
    return procedures.filter((p) => {
      const matchesSearch = !q ||
        p.title.toLowerCase().includes(q) ||
        (p.organizationName || "").toLowerCase().includes(q) ||
        (p.documentsUploaded || "").toLowerCase().includes(q);
      const matchesOrg = orgFilter === "all" || (p.organizationName === orgFilter);
      return matchesSearch && matchesOrg;
    });
  }, [procedures, searchQuery, orgFilter]);

  // Pagination state (10 per page)
  const PAGE_SIZE = 10;
  const [page, setPage] = useState(1);
  const totalPages = Math.max(1, Math.ceil(filteredProcedures.length / PAGE_SIZE));
  useEffect(() => {
    if (page > totalPages) setPage(totalPages);
  }, [totalPages, page]);
  const paginatedProcedures = useMemo(() => {
    const start = (page - 1) * PAGE_SIZE;
    return filteredProcedures.slice(start, start + PAGE_SIZE);
  }, [filteredProcedures, page]);

  // Add handleRetry to reload the page
  function handleRetry() {
    router.refresh?.() // Next.js 13+
    if (!router.refresh) window.location.reload()
  }

  return (
    <div className="min-h-screen w-full bg-gray-50 relative overflow-hidden">
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        <div className="absolute -top-24 -right-24 w-56 h-56 rounded-full blur-3xl" style={{ background: 'radial-gradient(closest-side, rgba(167,179,185,0.10), rgba(167,179,185,0))' }}></div>
        <div className="absolute -bottom-28 -left-28 w-64 h-64 rounded-full blur-3xl" style={{ background: 'radial-gradient(closest-side, rgba(94,156,141,0.10), rgba(94,156,141,0))' }}></div>
      </div>

      <main className="relative z-10 p-4 sm:p-6 md:p-8">
        <div className="max-w-7xl mx-auto">
          {/* Header Section */}
          <div className="mb-6 md:mb-8 w-full">
            <div className="bg-white/90 backdrop-blur-sm rounded-2xl border border-[#a7b3b9]/30 p-4 sm:p-6 md:p-8 shadow-lg">
              <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                <div>
                  <h1 className="text-2xl sm:text-3xl font-bold text-[#2e4d57] mb-2">{t("workspace.title")}</h1>
                  <p className="text-[#1c3b2e] text-sm sm:text-base">{t("workspace.description")}</p>
                </div>
                <div className="inline-flex items-center gap-2 bg-[#3a6a8d]/10 backdrop-blur-sm border border-[#3a6a8d]/30 rounded-full px-3 sm:px-4 py-2">
                  <CheckCircle className="w-4 h-4 text-[#3a6a8d] animate-pulse" />
                  <span className="text-xs sm:text-sm font-medium text-[#2e4d57]">{t("workspace.label")}</span>
                </div>
              </div>
            </div>
          </div>

          {/* Error Display */}
          {checklistsError && (
            <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
              <div className="flex items-center gap-2 text-red-800 mb-2">
                <AlertCircle className="w-5 h-5" />
                <h3 className="font-semibold">{t("workspace.errors.title")}</h3>
              </div>
              <p className="text-red-700 mb-4">
                {('status' in checklistsError && checklistsError.status === 401) 
                  ? t("workspace.errors.auth_failed")
                  : t("workspace.errors.load_failed")
                }
              </p>
              <div className="flex gap-3">
                <Button 
                  onClick={handleRetry}
                  variant="outline" 
                  className="border-red-300 text-red-700 hover:bg-red-50"
                >
                  {t("workspace.actions.retry")}
                </Button>
                <Button 
                  onClick={() => router.push("/login")}
                  className="bg-red-600 hover:bg-red-700"
                >
                  {t("workspace.actions.login")}
                </Button>
              </div>
            </div>
          )}

          {/* Stats Cards */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
            {stats.map((stat, index) => (
              <motion.div key={index} variants={listItemVariants} initial="hidden" animate="visible" transition={{ delay: index * 0.05 + 0.05 }}>
                <Card className="bg-white/80 backdrop-blur-md rounded-2xl border border-[#e5e7eb] shadow-xl hover:shadow-2xl transition-all duration-700 hover:scale-105 hover:-translate-y-2 cursor-pointer relative overflow-hidden group">
                  <div className="absolute inset-0 bg-gradient-to-r from-[#3a6a8d]/10 via-transparent to-[#5e9c8d]/10 opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
                  <CardContent className="p-6 relative z-10">
                    <div className="flex items-center justify-between">
                      <div>
                        <p className="text-sm font-medium text-[#4b5563] mb-1">{stat.title}</p>
                        <p className="text-3xl font-bold text-[#111827]">{stat.value}</p>
                      </div>
                      <div className="w-10 h-10 rounded-lg flex items-center justify-center transition-all duration-500 group-hover:scale-110 group-hover:rotate-6" style={{ backgroundColor: '#e6f0f5' }}>
                        <stat.icon className="w-6 h-6" style={{ color: '#3a6a8d' }} />
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </motion.div>
            ))}
          </div>

          {/* Search & Organization Filter */}
          <div className="bg-white/90 backdrop-blur-sm rounded-2xl border border-[#e5e7eb] p-4 sm:p-5 mb-6 shadow-md flex flex-col gap-4 sm:flex-row sm:items-center">
            <div className="relative flex-1 flex min-w-[220px]">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
              <input
                type="text"
                placeholder={t("workspace.search.placeholder")}
                value={searchInput}
                onChange={(e) => setSearchInput(e.target.value)}
                className="w-full pl-10 pr-4 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#3A6A8D] focus:border-transparent"
              />
              <Button
                type="button"
                className="ml-2 px-4 py-2 bg-[#3A6A8D] hover:bg-[#2d5470] text-white"
                onClick={() => { setSearchQuery(searchInput); setPage(1); }}
              >
                {t("workspace.search.button")}
              </Button>
            </div>
            <div className="flex flex-none w-full sm:w-64">
              <Select value={orgFilter} onValueChange={(v) => { setOrgFilter(v); setPage(1); }}>
                <SelectTrigger className="w-full border-[#3A6A8D] text-[#3A6A8D] hover:bg-[#3A6A8D]/5 focus:ring-2 focus:ring-[#3A6A8D] focus:border-transparent">
                  <SelectValue placeholder={t("workspace.filter.all")} />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">{t("workspace.filter.all")}</SelectItem>
                  {Array.from(new Set(procedures.map((p) => p.organizationName).filter(Boolean))).map((name) => (
                    <SelectItem key={name as string} value={name as string}>{name as string}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>

          {/* Procedure Cards */}
          <div className="space-y-4">
            {loadingChecklists ? (
              <div className="text-center py-8">
                <div className="text-gray-600">{t("workspace.loading")}</div>
              </div>
            ) : checklistsError ? (
              <div className="text-center py-8">
                <div className="text-red-600">{t("workspace.errors.load_failed")}</div>
              </div>
            ) : procedures.length === 0 ? (
              <div className="text-center py-8">
                <div className="text-gray-600">{t("workspace.empty")}</div>
                <Button 
                  className="mt-4 bg-[#3A6A8D] hover:bg-[#2d5470] text-white"
                  onClick={() => router.push('/user/procedures-list')}
                >
                  {t("workspace.actions.browse")}
                </Button>
              </div>
            ) : (
              paginatedProcedures.map((procedure, index) => (
                <motion.div key={procedure.id} variants={listItemVariants} initial="hidden" animate="visible" transition={{ delay: index * 0.06 + 0.12 }}>
                  <Card className="bg-white/80 backdrop-blur-md rounded-2xl border border-[#e5e7eb] shadow-xl hover:shadow-2xl transition-all duration-700 hover:scale-105 hover:-translate-y-2 relative overflow-hidden group">
                    <div className="absolute inset-0 bg-gradient-to-r from-[#3a6a8d]/10 via-transparent to-[#5e9c8d]/10 opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
                    <CardContent className="p-6 relative z-10">
                      <div className="flex items-start justify-between mb-4">
                        <div className="flex items-start gap-4">
                          <div className="w-12 h-12 rounded-xl flex items-center justify-center transition-all duration-500 group-hover:scale-110 group-hover:rotate-3" style={{ backgroundColor: '#e6f0f5' }}>
                            <FileText className="w-6 h-6" style={{ color: '#3a6a8d' }} />
                          </div>
                          <div>
                            <h3 className="font-semibold text-lg text-[#111827] mb-1 group-hover:text-[#3a6a8d] transition-colors duration-300">
                              {procedure.title}
                            </h3>
                            <p className="text-[#4b5563] text-sm">{procedure.department}</p>
                          </div>
                        </div>
                        <div className="flex items-center gap-3">
                          <Badge className={`${procedure.statusColor} border border-[#e5e7eb] transition-all duration-200`}>
                            {procedure.status}
                          </Badge>
                          <Button
                            variant={procedure.buttonVariant}
                            size="sm"
                            className={
                              procedure.buttonVariant === "default"
                                ? "bg-gradient-to-r from-[#3a6a8d] to-[#2e4d57] hover:from-[#2e4d57] hover:to-[#1c3b2e] text-white transition-all duration-500 hover:scale-105 hover:shadow-lg"
                                : "transition-all duration-300"
                            }
                            onClick={() => router.push(`/user/checklist/${encodeURIComponent(procedure.id)}`)}
                          >
                            {procedure.buttonText}
                          </Button>
                        </div>
                      </div>

                      <div className="mb-4">
                        <div className="flex items-center justify-between mb-2">
                          <span className="text-sm font-medium text-[#111827]">{t("workspace.progress")}</span>
                          <span className="text-sm text-[#4b5563]">{t("workspace.progress_value", { percent: procedure.progress })}</span>
                        </div>
                        <div className="w-full bg-[#e5e7eb] rounded-full h-2 overflow-hidden">
                          {(() => {
                            const barColor =
                              procedure.status === t("workspace.status.completed")
                                ? "bg-[#5e9c8d]"
                                : procedure.status === t("workspace.status.in_progress")
                                  ? "bg-gradient-to-r from-[#3a6a8d] to-[#2e4d57]"
                                  : "bg-[#a7b3b9]/50";
                            return (
                              <motion.div
                                className={`h-2 rounded-full ${barColor}`}
                                initial={{ width: 0 }}
                                animate={{ width: `${procedure.progress}%` }}
                                transition={{ type: "spring", stiffness: 180, damping: 24, mass: 0.9 }}
                              />
                            );
                          })()}
                        </div>
                      </div>

                      <div className="flex items-center gap-6 text-sm text-[#4b5563]">
                        {procedure.startDate && (
                          <div className="flex items-center gap-1">
                            <Calendar className="w-4 h-4" />
                            <span>{t("workspace.start_date", { date: procedure.startDate })}</span>
                          </div>
                        )}
                        {procedure.estimatedCompletion && (
                          <div className="flex items-center gap-1">
                            <Clock className="w-4 h-4" />
                            <span>{t("workspace.estimated_completion", { status: procedure.estimatedCompletion })}</span>
                          </div>
                        )}
                        {procedure.documentsUploaded && (
                          <div className="flex items-center gap-1">
                            <FileText className="w-4 h-4" />
                            <span>{procedure.documentsUploaded}</span>
                          </div>
                        )}
                      </div>
                    </CardContent>
                  </Card>
                </motion.div>
              ))
            )}
          </div>

          {totalPages > 1 && (
            <div className="mt-6">
              <Pagination page={page} totalPages={totalPages} onPageChange={setPage} />
            </div>
          )}
        </div>
      </main>
    </div>
  )
}
