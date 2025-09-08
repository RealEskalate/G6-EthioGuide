// "use client"

// export const dynamic = 'force-dynamic'

// import { Suspense, useEffect } from "react"
// import { useSearchParams, useRouter } from "next/navigation"
// import { useGetProcedureFlexibleQuery } from "@/app/store/slices/proceduresApi"
// import {
//   ArrowLeft,
//   Calendar,
//   CheckCircle,
//   Clock,
//   DollarSign,
//   Download,
//   Eye,
//   FileText,
//   Share2,
//   ThumbsUp,
// } from "lucide-react"
// import Link from "next/link"
// import { Button } from "@/components/ui/button"
// import { Badge } from "@/components/ui/badge"
// import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
// import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
// import { useGetProcedureFeedbackQuery } from "@/app/store/slices/feedbackApi"
// import { useCreateChecklistMutation } from "@/app/store/slices/checklistsApi"
// import { useSession } from "next-auth/react"
// import { useListDiscussionsQuery } from "@/app/store/slices/discussionsApi"

// function ProcedureDetailInner() {
//   const search = useSearchParams()
//   const id = search.get("id") || ""
//   const router = useRouter()
  
//   const { data: procedure, isLoading, isError } = useGetProcedureFlexibleQuery(id, { skip: !id })
//   const { data: session } = useSession()
//   const [createChecklist, { isLoading: savingChecklist }] = useCreateChecklistMutation()
//   const { data: feedbackData, isLoading: loadingFeedback } = useGetProcedureFeedbackQuery({ procedureId: id, page: 1, limit: 10, token: session?.accessToken || null }, { skip: !id })
  
//   const notFound = !isLoading && !isError && (!procedure || !procedure.id)

//   // Discussions for this procedure
//   const { data: discussions, isLoading: loadingDiscussions } = useListDiscussionsQuery(
//     { procedureIds: id ? [id] : [], page: 1, limit: 5 },
//     { skip: !id }
//   )

//   // Debug logging
//   useEffect(() => {
//     console.log('Feedback debug:', {
//       id,
//       loadingFeedback,
//       feedbackData: feedbackData ? {
//         feedbacks: feedbackData.feedbacks?.length || 0,
//         page: feedbackData.page,
//         total: feedbackData.total
//       } : null
//     })
//   }, [id, loadingFeedback, feedbackData])

//   // Feedback submission is handled on /user/feedback; this page only lists feedback and links out.

//   const totalFees = Array.isArray(procedure?.fees) && (procedure!.fees?.length || 0) > 0
//     ? (() => {
//         const feesArr = procedure!.fees!
//         const total = feesArr.reduce((sum, f) => sum + (Number(f.amount) || 0), 0)
//         const currency = feesArr[0]?.currency || "ETB"
//         return `${total} ${currency}`.trim()
//       })()
//     : null

//   return (
//     <div className="min-h-screen w-full bg-gray-50  relative overflow-hidden">
//       <div className="absolute inset-0 overflow-hidden pointer-events-none">
//         <div className="absolute -top-40 -right-40 w-80 h-80 bg-gradient-to-br from-[#3a6a8d]/10 to-[#2e4d57]/10 rounded-full blur-3xl animate-pulse"></div>
//         <div
//           className="absolute -bottom-40 -left-40 w-80 h-80 bg-gradient-to-tr from-[#5e9c8d]/10 to-[#1c3b2e]/10 rounded-full blur-3xl animate-pulse"
//           style={{ animationDelay: "2s" }}
//         ></div>
//       </div>

//       <div className="relative z-10 p-4 sm:p-6 md:p-8 max-w-7xl mx-auto">
//         <div className="mb-6 md:mb-8">
//           <Link href="/user/home">
//             <Button
//               variant="ghost"
//               className="mb-4 hover:bg-[#3a6a8d]/10 hover:text-[#3a6a8d] transition-all duration-300 rounded-xl px-3 sm:px-4 py-2 text-[#2e4d57]"
//             >
//               <ArrowLeft className="w-4 h-4 mr-2" />
//               Back to Home
//             </Button>
//           </Link>

//           <div className="bg-white/90 backdrop-blur-sm rounded-2xl border border-[#a7b3b9]/30 p-4 sm:p-6 md:p-8 shadow-lg">
//             <div className="flex flex-col lg:flex-row lg:items-start justify-between gap-4 md:gap-6">
//               <div className="flex-1">
//                 <div className="flex flex-wrap items-center gap-2 mb-3 md:mb-4">
//                   {Boolean(procedure?.verified) && (
//                     <Badge className="bg-gradient-to-r from-[#5e9c8d]/20 to-[#1c3b2e]/20 text-[#1c3b2e] border-[#5e9c8d]/30 flex items-center gap-1.5 px-2 sm:px-3 py-1 rounded-full text-xs">
//                       <CheckCircle className="w-3 sm:w-3.5 h-3 sm:h-3.5" />
//                       Verified
//                     </Badge>
//                   )}
//                   {(procedure?.tags ?? []).map((tag, index) => (
//                     <Badge key={index} variant="outline" className="text-[#3a6a8d] border-[#3a6a8d]/30 text-xs">
//                       {tag}
//                     </Badge>
//                   ))}
//                 </div>

//                 <h1 className="text-2xl sm:text-3xl md:text-4xl font-bold text-[#2e4d57] mb-3 md:mb-4 leading-tight">
//                   {isLoading
//                     ? "Loading..."
//                     : notFound
//                     ? "Procedure Not Found"
//                     : (procedure?.title || procedure?.name || "Procedure")}
//                 </h1>

//                 {!isError && !notFound && procedure?.summary && (
//                   <p className="text-[#1c3b2e] text-sm sm:text-base md:text-lg leading-relaxed mb-4 md:mb-6">
//                     {procedure.summary}
//                   </p>
//                 )}
//                 {isError && (
//                   <p className="text-red-600 text-sm">Failed to load procedure.</p>
//                 )}

//                 <div className="flex flex-wrap items-center gap-3 md:gap-4 text-xs sm:text-sm text-[#a7b3b9]">
//           {typeof procedure?.views === "number" && (
//                     <span className="flex items-center gap-1.5">
//                       <Eye className="w-4 h-4" />
//             {procedure.views} views
//                     </span>
//                   )}
//           {typeof procedure?.likes === "number" && (
//                     <span className="flex items-center gap-1.5">
//                       <ThumbsUp className="w-4 h-4" />
//             {procedure.likes} likes
//                     </span>
//                   )}
//           {procedure?.updatedAt && (
//                     <span className="flex items-center gap-1.5">
//                       <Calendar className="w-4 h-4" />
//             Updated {new Date(procedure.updatedAt).toLocaleDateString()}
//                     </span>
//                   )}
//                 </div>
//               </div>

//               <div className="flex flex-col sm:flex-row lg:flex-col gap-2 lg:w-48">
//                 <Button className="bg-gradient-to-r from-[#3a6a8d] to-[#2e4d57] hover:from-[#2e4d57] hover:to-[#1c3b2e] text-white transition-all duration-300 hover:scale-105 rounded-xl py-2.5 sm:py-3 font-medium text-sm">
//                   <Download className="w-4 h-4 mr-2" />
//                   Download Guide
//                 </Button>
//                 <Button
//                   variant="outline"
//                   className="border-[#3a6a8d]/30 text-[#3a6a8d] hover:bg-[#3a6a8d]/10 transition-all duration-300 rounded-xl py-2.5 sm:py-3 font-medium text-sm bg-transparent"
//                 >
//                   <Share2 className="w-4 h-4 mr-2" />
//                   Share
//                 </Button>
//                 <Button
//                   disabled={!id || savingChecklist}
//                   onClick={async () => {
//                     if (!id) return
//                     try {
//                       const result = await createChecklist({ procedureId: id, token: session?.accessToken || undefined }).unwrap()
//                       console.log('Checklist created successfully:', result)
//                     } catch (error) {
//                       console.error('Failed to create checklist:', error)
//                       // Still navigate to workspace so user can see current state
//                     } finally {
//                       router.push('/user/workspace')
//                     }
//                   }}
//                   className="bg-[#3a6a8d] hover:bg-[#2e4d57] text-white transition-all duration-300 rounded-xl py-2.5 sm:py-3 font-medium text-sm"
//                 >
//                   {savingChecklist ? 'Saving…' : 'Save Checklist'}
//                 </Button>
//               </div>
//             </div>
//           </div>
//         </div>

//         <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 md:gap-8">
//           <div className="lg:col-span-2 space-y-6 md:space-y-8">
//             <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 md:gap-6">
//               <Card className="bg-white/90 backdrop-blur-sm border-[#a7b3b9]/30 shadow-lg hover:shadow-xl transition-all duration-300 rounded-2xl">
//                 <CardHeader className="pb-3">
//                   <CardTitle className="flex items-center gap-2 text-[#2e4d57] text-base sm:text-lg">
//                     <Clock className="w-5 h-5 text-[#3a6a8d]" />
//                     Processing Time
//                   </CardTitle>
//                 </CardHeader>
//                 <CardContent>
//                   <p className="text-xl sm:text-2xl font-bold text-[#3a6a8d]">
//                     {procedure?.processingTime?.minDays ?? "—"}-{procedure?.processingTime?.maxDays ?? "—"} days
//                   </p>
//                   <p className="text-xs sm:text-sm text-[#a7b3b9] mt-1">Average processing time</p>
//                 </CardContent>
//               </Card>

//               <Card className="bg-white/90 backdrop-blur-sm border-[#a7b3b9]/30 shadow-lg hover:shadow-xl transition-all duration-300 rounded-2xl">
//                 <CardHeader className="pb-3">
//                   <CardTitle className="flex items-center gap-2 text-[#2e4d57] text-base sm:text-lg">
//                     <DollarSign className="w-5 h-5 text-[#5e9c8d]" />
//                     Total Fees
//                   </CardTitle>
//                 </CardHeader>
//                 <CardContent>
//                   <p className="text-xl sm:text-2xl font-bold text-[#5e9c8d]">{totalFees ?? "—"}</p>
//                   <p className="text-xs sm:text-sm text-[#a7b3b9] mt-1">All fees included</p>
//                 </CardContent>
//               </Card>
//             </div>

//             <Card className="bg-white/90 backdrop-blur-sm border-[#a7b3b9]/30 shadow-lg rounded-2xl">
//               <CardHeader>
//                 <CardTitle className="flex items-center gap-2 text-[#2e4d57] text-lg sm:text-xl">
//                   <FileText className="w-5 sm:w-6 h-5 sm:h-6 text-[#3a6a8d]" />
//                   Step-by-Step Instructions
//                 </CardTitle>
//                 <CardDescription className="text-[#1c3b2e] text-sm sm:text-base">
//                   Follow these steps to complete your {(procedure?.title || procedure?.name || "procedure").toString().toLowerCase()}
//                 </CardDescription>
//               </CardHeader>
//               <CardContent className="space-y-4 md:space-y-6">
//                 {Array.isArray(procedure?.steps) && procedure!.steps.length > 0 ? (
//                   [...(procedure!.steps || [])]
//                     .sort((a, b) => ((a?.order ?? 0) - (b?.order ?? 0)))
//                     .map((s, index) => (
//                     <div
//                       key={index}
//                       className="flex gap-3 md:gap-4 p-3 md:p-4 bg-gradient-to-r from-[#a7b3b9]/10 to-[#5e9c8d]/10 rounded-xl border border-[#a7b3b9]/20 hover:shadow-md transition-all duration-300"
//                     >
//                       <div className="flex-shrink-0 w-8 sm:w-10 h-8 sm:h-10 bg-gradient-to-r from-[#3a6a8d] to-[#2e4d57] text-white rounded-full flex items-center justify-center font-bold text-sm sm:text-base">
//                         {s.order ?? index + 1}
//                       </div>
//                       <div className="flex-1 min-w-0">
//                         <h3 className="font-semibold text-[#2e4d57] mb-1 sm:mb-2 text-sm sm:text-base">{s.title || s.text?.slice(0, 80) || "Step"}</h3>
//                         <p className="text-[#1c3b2e] text-xs sm:text-sm leading-relaxed mb-2">{s.text || s.description}</p>
//                         {(s.estimatedTime || s.time) && (
//                           <span className="inline-flex items-center gap-1 text-xs font-medium text-[#3a6a8d] bg-[#3a6a8d]/10 px-2 py-1 rounded-full">
//                             <Clock className="w-3 h-3" />
//                             {s.estimatedTime || s.time}
//                           </span>
//                         )}
//                       </div>
//                     </div>
//                   ))
//                 ) : (
//                   <div className="text-sm text-[#6b7280] italic">No steps defined for this procedure.</div>
//                 )}
//               </CardContent>
//             </Card>

//             <Tabs defaultValue="feedback" className="w-full">
//               <TabsList className="grid w-full grid-cols-3 bg-white/90 backdrop-blur-sm border border-[#a7b3b9]/30 rounded-xl p-1">
//                 <TabsTrigger
//                   value="feedback"
//                   className="data-[state=active]:bg-gradient-to-r data-[state=active]:from-[#3a6a8d] data-[state=active]:to-[#2e4d57] data-[state=active]:text-white rounded-lg text-xs sm:text-sm"
//                 >
//                   Feedback
//                 </TabsTrigger>
//                 <TabsTrigger
//                   value="notices"
//                   className="data-[state=active]:bg-gradient-to-r data-[state=active]:from-[#3a6a8d] data-[state=active]:to-[#2e4d57] data-[state=active]:text-white rounded-lg text-xs sm:text-sm"
//                 >
//                   Notices
//                 </TabsTrigger>
//                 <TabsTrigger
//                   value="discussions"
//                   className="data-[state=active]:bg-gradient-to-r data-[state=active]:from-[#3a6a8d] data-[state=active]:to-[#2e4d57] data-[state=active]:text-white rounded-lg text-xs sm:text-sm"
//                 >
//                   Discussion
//                 </TabsTrigger>
//               </TabsList>

//               <TabsContent value="feedback" className="mt-4 md:mt-6">
//                 <Card className="bg-white/90 backdrop-blur-sm border-[#a7b3b9]/30 shadow-lg rounded-2xl">
//                   <CardHeader>
//                     <CardTitle className="text-[#2e4d57] text-base sm:text-lg">Recent Feedback</CardTitle>
//                     <CardDescription className="text-[#1c3b2e] text-sm">
//                       What others say about this procedure
//                     </CardDescription>
//                   </CardHeader>
//                   <CardContent className="space-y-3">
//                     {loadingFeedback && <div className="text-center text-gray-500 py-4">Loading feedback...</div>}
//                     {!loadingFeedback && (feedbackData?.feedbacks?.length ?? 0) === 0 && (
//                       <div className="text-center text-gray-500 py-4">No feedback yet. Be the first to share!</div>
//                     )}
//                     {(feedbackData?.feedbacks ?? []).map((f) => (
//                       <div
//                         key={f.id}
//                         className="flex gap-3 p-3 rounded-lg bg-gradient-to-r from-[#a7b3b9]/10 to-[#5e9c8d]/10 border border-[#a7b3b9]/20"
//                       >
//                         <div className="w-8 h-8 rounded-full bg-[#ced4da] flex items-center justify-center text-xs font-semibold overflow-hidden">
//                           <span>{f.userID?.charAt(0)?.toUpperCase() || "U"}</span>
//                         </div>
//                         <div className="flex-1">
//                           <div className="flex items-center gap-2 mb-1">
//                             <span className="text-sm font-medium text-[#2e4d57]">User {f.userID?.slice(-4) || "Anonymous"}</span>
//                             <Badge className="text-xs px-1 py-0 bg-[#e6f4ff] text-[#1c3b2e]">
//                               {f.type?.replace("_", " ")}
//                             </Badge>
//                           </div>
//                           <p className="text-xs text-[#1c3b2e] mb-1">{f.content}</p>
//                           {f.tags?.length ? (
//                             <div className="flex gap-1 mb-2 flex-wrap">
//                               {f.tags.map((t, i) => (
//                                 <span key={i} className="text-xs bg-gray-200 text-gray-600 px-2 py-0.5 rounded">#{t}</span>
//                               ))}
//                             </div>
//                           ) : null}
//                           <div className="flex items-center gap-3 text-xs text-[#9ca3af]">
//                             <span>{f.createdAT ? new Date(f.createdAT).toLocaleDateString() : "Recently"}</span>
//                           </div>
//                         </div>
//                       </div>
//                     ))}
//                     <Button asChild className="w-full bg-gradient-to-r from-[#3a6a8d] to-[#2e4d57] hover:from-[#2e4d57] hover:to-[#1c3b2e] text-white rounded-xl text-sm mt-2">
//                       <a href={`/user/feedback?id=${encodeURIComponent(id)}`}>Add Feedback</a>
//                     </Button>
//                   </CardContent>
//                 </Card>
//               </TabsContent>

//               <TabsContent value="notices" className="mt-4 md:mt-6">
//                 <Card className="bg-white/90 backdrop-blur-sm border-[#a7b3b9]/30 shadow-lg rounded-2xl">
//                   <CardHeader>
//                     <CardTitle className="text-[#2e4d57] text-base sm:text-lg">Important Notices</CardTitle>
//                   </CardHeader>
//                   <CardContent className="space-y-3 md:space-y-4">
//                     <div className="flex gap-3 p-3 md:p-4 bg-gradient-to-r from-[#5e9c8d]/10 to-[#1c3b2e]/10 border border-[#5e9c8d]/30 rounded-xl">
//                       <FileText className="w-5 h-5 text-[#5e9c8d] flex-shrink-0 mt-0.5" />
//                       <div>
//                         <h4 className="font-semibold text-[#2e4d57] text-sm sm:text-base">No notices yet</h4>
//                         <p className="text-[#1c3b2e] text-xs sm:text-sm">We will display official updates here.</p>
//                       </div>
//                     </div>
//                   </CardContent>
//                 </Card>
//               </TabsContent>

//               <TabsContent value="discussions" className="mt-4 md:mt-6">
//                 <Card className="bg-white/90 backdrop-blur-sm border-[#a7b3b9]/30 shadow-lg rounded-2xl">
//                   <CardHeader>
//                     <CardTitle className="text-[#2e4d57] text-base sm:text-lg">Community Discussion</CardTitle>
//                     <CardDescription className="text-[#1c3b2e] text-sm">Ask questions and get help from the community</CardDescription>
//                   </CardHeader>
//                   <CardContent className="space-y-4">
//                     {loadingDiscussions && <div className="text-center text-gray-500 py-6">Loading discussions...</div>}
//                     {!loadingDiscussions && (discussions?.posts?.length ?? 0) === 0 && (
//                       <div className="text-center text-gray-500 py-6">No discussions yet. Start the first one!</div>
//                     )}
//                     {(discussions?.posts ?? []).map((post) => (
//                       <div
//                         key={post.id}
//                         className="bg-white border border-[#a7b3b9]/30 rounded-lg p-4 hover:bg-[#f8fafc] transition-all duration-300"
//                       >
//                         <div className="flex items-start gap-3">
//                           <div className="w-10 h-10 rounded-full bg-[#3a6a8d] flex items-center justify-center text-white text-sm font-semibold overflow-hidden">
//                             <span>{post.userID?.slice(-2).toUpperCase() || "U"}</span>
//                           </div>
//                           <div className="flex-1">
//                             <h4 className="font-medium text-[#2e4d57] mb-1 text-sm line-clamp-1">{post.title || "Untitled"}</h4>
//                             <p className="text-sm text-[#1c3b2e] mb-2 line-clamp-2">{post.content}</p>
//                             {post.tags?.length ? (
//                               <div className="flex items-center gap-2 mb-2 flex-wrap">
//                                 {post.tags.slice(0, 4).map((t, i) => (
//                                   <Badge key={i} className="bg-[#dbeafe] text-[#1e40af] text-[10px] px-2 py-0.5">{t}</Badge>
//                                 ))}
//                               </div>
//                             ) : null}
//                             <div className="flex items-center gap-4 text-xs text-[#9ca3af]">
//                               <span>{post.createdAt ? new Date(post.createdAt).toLocaleDateString() : ""}</span>
//                             </div>
//                           </div>
//                         </div>
//                       </div>
//                     ))}
//                     <Button asChild className="w-full bg-gradient-to-r from-[#3a6a8d] to-[#2e4d57] hover:from-[#2e4d57] hover:to-[#1c3b2e] text-white rounded-xl text-sm">
//                       <Link href="/user/discussions">View Discussions</Link>
//                     </Button>
//                   </CardContent>
//                 </Card>
//               </TabsContent>
//             </Tabs>
//           </div>

//           {/* Aside */}
//           <div className="lg:col-span-1 space-y-6">
//             <Card className="bg-white/90 backdrop-blur-sm border-[#a7b3b9]/30 shadow-lg rounded-2xl">
//               <CardHeader>
//                 <CardTitle className="flex items-center gap-2 text-[#2e4d57] text-base sm:text-lg">
//                   <FileText className="w-5 h-5 text-[#5e9c8d]" />
//                   Required Documents
//                 </CardTitle>
//               </CardHeader>
//               <CardContent>
//                 <ul className="space-y-2">
//                   {Array.isArray(procedure?.documentsRequired) && procedure!.documentsRequired.length > 0 ? (
//                     procedure!.documentsRequired.map((doc, index) => (
//                       <li key={index} className="flex items-center gap-2 text-[#1c3b2e] text-sm">
//                         <CheckCircle className="w-4 h-4 text-[#5e9c8d] flex-shrink-0" />
//                         {typeof doc === "string" ? doc : (doc?.name || "Document")}
//                       </li>
//                     ))
//                   ) : (
//                     <li className="text-sm text-[#6b7280] italic">No documents listed.</li>
//                   )}
//                 </ul>
//               </CardContent>
//             </Card>
//           </div>
//         </div>
//       </div>
//     </div>
//   )
// }

// export default function ProcedureDetailPage() {
//   return (
//     <Suspense fallback={<div className="p-6">Loading...</div>}>
//       <ProcedureDetailInner />
//     </Suspense>
//   )
// }


"use client"

export const dynamic = 'force-dynamic'

import { Suspense, useEffect } from "react"
import { useSearchParams, useRouter } from "next/navigation"
import { useTranslation } from "react-i18next"
import { useGetProcedureFlexibleQuery } from "@/app/store/slices/proceduresApi"
import {
  ArrowLeft,
  Calendar,
  CheckCircle,
  Clock,
  DollarSign,
  Download,
  Eye,
  FileText,
  Share2,
  ThumbsUp,
} from "lucide-react"
import Link from "next/link"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card"
import { useGetProcedureFeedbackQuery } from "@/app/store/slices/feedbackApi"
import { useCreateChecklistMutation } from "@/app/store/slices/checklistsApi"
import { useSession } from "next-auth/react"
import { useListDiscussionsQuery } from "@/app/store/slices/discussionsApi"

function ProcedureDetailInner() {
  const { t } = useTranslation("user")
  const search = useSearchParams()
  const id = search.get("id") || ""
  const router = useRouter()
  
  const { data: procedure, isLoading, isError } = useGetProcedureFlexibleQuery(id, { skip: !id })
  const { data: session } = useSession()
  const [createChecklist, { isLoading: savingChecklist }] = useCreateChecklistMutation()
  const { data: feedbackData, isLoading: loadingFeedback } = useGetProcedureFeedbackQuery(
    { procedureId: id, page: 1, limit: 10, token: session?.accessToken || null },
    { skip: !id }
  )
  
  const { data: discussions, isLoading: loadingDiscussions } = useListDiscussionsQuery(
    { procedureIds: id ? [id] : [], page: 1, limit: 5 },
    { skip: !id }
  )

  const notFound = !isLoading && !isError && (!procedure || !procedure.id)

  useEffect(() => {
    console.log('Procedure debug:', { id, procedure, isLoading, isError })
    console.log('Feedback debug:', {
      id,
      loadingFeedback,
      feedbackData: feedbackData ? {
        feedbacks: feedbackData.feedbacks?.length || 0,
        page: feedbackData.page,
        total: feedbackData.total
      } : null
    })
    console.log('Discussions debug:', {
      id,
      loadingDiscussions,
      discussions: discussions ? {
        posts: discussions.posts?.length || 0,
        page: discussions.page,
        total: discussions.total
      } : null
    })
  }, [id, procedure, isLoading, isError, loadingFeedback, feedbackData, loadingDiscussions, discussions])

  const totalFees = Array.isArray(procedure?.fees) && (procedure!.fees?.length || 0) > 0
    ? (() => {
        const feesArr = procedure!.fees!
        const total = feesArr.reduce((sum, f) => sum + (Number(f.amount) || 0), 0)
        const currency = feesArr[0]?.currency || "ETB"
        return t("procedure_detail.fees_value", { total, currency })
      })()
    : t("procedure_detail.no_fees")

  return (
    <div className="min-h-screen w-full bg-gray-50 relative overflow-hidden">
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        <div className="absolute -top-40 -right-40 w-80 h-80 bg-gradient-to-br from-[#3a6a8d]/10 to-[#2e4d57]/10 rounded-full blur-3xl animate-pulse"></div>
        <div
          className="absolute -bottom-40 -left-40 w-80 h-80 bg-gradient-to-tr from-[#5e9c8d]/10 to-[#1c3b2e]/10 rounded-full blur-3xl animate-pulse"
          style={{ animationDelay: "2s" }}
        ></div>
      </div>

      <div className="relative z-10 p-4 sm:p-6 md:p-8 max-w-7xl mx-auto">
        <div className="mb-6 md:mb-8">
          <Link href="/user/home">
            <Button
              variant="ghost"
              className="mb-4 hover:bg-[#3a6a8d]/10 hover:text-[#3a6a8d] transition-all duration-300 rounded-xl px-3 sm:px-4 py-2 text-[#2e4d57]"
            >
              <ArrowLeft className="w-4 h-4 mr-2" />
              {t("procedure_detail.back")}
            </Button>
          </Link>

          <div className="bg-white/90 backdrop-blur-sm rounded-2xl border border-[#a7b3b9]/30 p-4 sm:p-6 md:p-8 shadow-lg">
            <div className="flex flex-col lg:flex-row lg:items-start justify-between gap-4 md:gap-6">
              <div className="flex-1">
                <div className="flex flex-wrap items-center gap-2 mb-3 md:mb-4">
                  {Boolean(procedure?.verified) && (
                    <Badge className="bg-gradient-to-r from-[#5e9c8d]/20 to-[#1c3b2e]/20 text-[#1c3b2e] border-[#5e9c8d]/30 flex items-center gap-1.5 px-2 sm:px-3 py-1 rounded-full text-xs">
                      <CheckCircle className="w-3 sm:w-3.5 h-3 sm:h-3.5" />
                      {t("procedure_detail.verified")}
                    </Badge>
                  )}
                  {(procedure?.tags ?? []).map((tag, index) => (
                    <Badge key={index} variant="outline" className="text-[#3a6a8d] border-[#3a6a8d]/30 text-xs">
                      {tag}
                    </Badge>
                  ))}
                </div>

                <h1 className="text-2xl sm:text-3xl md:text-4xl font-bold text-[#2e4d57] mb-3 md:mb-4 leading-tight">
                  {isLoading
                    ? t("procedure_detail.loading")
                    : notFound
                    ? t("procedure_detail.not_found")
                    : (procedure?.title || procedure?.name || t("procedure_detail.default_procedure"))}
                </h1>

                {!isError && !notFound && procedure?.summary && (
                  <p className="text-[#1c3b2e] text-sm sm:text-base md:text-lg leading-relaxed mb-4 md:mb-6">
                    {procedure.summary}
                  </p>
                )}
                {isError && (
                  <p className="text-red-600 text-sm">{t("procedure_detail.errors.load_failed")}</p>
                )}

                <div className="flex flex-wrap items-center gap-3 md:gap-4 text-xs sm:text-sm text-[#a7b3b9]">
                  {typeof procedure?.views === "number" && (
                    <span className="flex items-center gap-1.5">
                      <Eye className="w-4 h-4" />
                      {t("procedure_detail.views", { count: procedure.views })}
                    </span>
                  )}
                  {typeof procedure?.likes === "number" && (
                    <span className="flex items-center gap-1.5">
                      <ThumbsUp className="w-4 h-4" />
                      {t("procedure_detail.likes", { count: procedure.likes })}
                    </span>
                  )}
                  {procedure?.updatedAt && (
                    <span className="flex items-center gap-1.5">
                      <Calendar className="w-4 h-4" />
                      {t("procedure_detail.updated_at", { date: new Date(procedure.updatedAt).toLocaleDateString() })}
                    </span>
                  )}
                </div>
              </div>

              <div className="flex flex-col sm:flex-row lg:flex-col gap-2 lg:w-48">
                <Button className="bg-gradient-to-r from-[#3a6a8d] to-[#2e4d57] hover:from-[#2e4d57] hover:to-[#1c3b2e] text-white transition-all duration-300 hover:scale-105 rounded-xl py-2.5 sm:py-3 font-medium text-sm">
                  <Download className="w-4 h-4 mr-2" />
                  {t("procedure_detail.actions.download_guide")}
                </Button>
                <Button
                  variant="outline"
                  className="border-[#3a6a8d]/30 text-[#3a6a8d] hover:bg-[#3a6a8d]/10 transition-all duration-300 rounded-xl py-2.5 sm:py-3 font-medium text-sm bg-transparent"
                >
                  <Share2 className="w-4 h-4 mr-2" />
                  {t("procedure_detail.actions.share")}
                </Button>
                <Button
                  disabled={!id || savingChecklist}
                  onClick={async () => {
                    if (!id) return
                    try {
                      const result = await createChecklist({ procedureId: id, token: session?.accessToken || undefined }).unwrap()
                      console.log('Checklist created successfully:', result)
                    } catch (error) {
                      console.error('Failed to create checklist:', error)
                    } finally {
                      router.push('/user/workspace')
                    }
                  }}
                  className="bg-[#3a6a8d] hover:bg-[#2e4d57] text-white transition-all duration-300 rounded-xl py-2.5 sm:py-3 font-medium text-sm"
                >
                  {savingChecklist ? t("procedure_detail.actions.saving") : t("procedure_detail.actions.save_checklist")}
                </Button>
              </div>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 md:gap-8">
          <div className="lg:col-span-2 space-y-6 md:space-y-8">
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 md:gap-6">
              <Card className="bg-white/90 backdrop-blur-sm border-[#a7b3b9]/30 shadow-lg hover:shadow-xl transition-all duration-300 rounded-2xl">
                <CardHeader className="pb-3">
                  <CardTitle className="flex items-center gap-2 text-[#2e4d57] text-base sm:text-lg">
                    <Clock className="w-5 h-5 text-[#3a6a8d]" />
                    {t("procedure_detail.processing_time")}
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-xl sm:text-2xl font-bold text-[#3a6a8d]">
                    {procedure?.processingTime?.minDays && procedure?.processingTime?.maxDays
                      ? t("procedure_detail.processing_time_value", {
                          min: procedure.processingTime.minDays,
                          max: procedure.processingTime.maxDays
                        })
                      : t("procedure_detail.no_processing_time")}
                  </p>
                  <p className="text-xs sm:text-sm text-[#a7b3b9] mt-1">{t("procedure_detail.processing_time_description")}</p>
                </CardContent>
              </Card>

              <Card className="bg-white/90 backdrop-blur-sm border-[#a7b3b9]/30 shadow-lg hover:shadow-xl transition-all duration-300 rounded-2xl">
                <CardHeader className="pb-3">
                  <CardTitle className="flex items-center gap-2 text-[#2e4d57] text-base sm:text-lg">
                    <DollarSign className="w-5 h-5 text-[#5e9c8d]" />
                    {t("procedure_detail.total_fees")}
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-xl sm:text-2xl font-bold text-[#5e9c8d]">{totalFees}</p>
                  <p className="text-xs sm:text-sm text-[#a7b3b9] mt-1">{t("procedure_detail.total_fees_description")}</p>
                </CardContent>
              </Card>
            </div>

            <Card className="bg-white/90 backdrop-blur-sm border-[#a7b3b9]/30 shadow-lg rounded-2xl">
              <CardHeader>
                <CardTitle className="flex items-center gap-2 text-[#2e4d57] text-lg sm:text-xl">
                  <FileText className="w-5 sm:w-6 h-5 sm:h-6 text-[#3a6a8d]" />
                  {t("procedure_detail.steps_title")}
                </CardTitle>
                <CardDescription className="text-[#1c3b2e] text-sm sm:text-base">
                  {t("procedure_detail.steps_description", {
                    procedure: (procedure?.title || procedure?.name || t("procedure_detail.default_procedure")).toString().toLowerCase()
                  })}
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4 md:space-y-6">
                {Array.isArray(procedure?.steps) && procedure!.steps.length > 0 ? (
                  [...(procedure!.steps || [])]
                    .sort((a, b) => ((a?.order ?? 0) - (b?.order ?? 0)))
                    .map((s, index) => (
                      <div
                        key={index}
                        className="flex gap-3 md:gap-4 p-3 md:p-4 bg-gradient-to-r from-[#a7b3b9]/10 to-[#5e9c8d]/10 rounded-xl border border-[#a7b3b9]/20 hover:shadow-md transition-all duration-300"
                      >
                        <div className="flex-shrink-0 w-8 sm:w-10 h-8 sm:h-10 bg-gradient-to-r from-[#3a6a8d] to-[#2e4d57] text-white rounded-full flex items-center justify-center font-bold text-sm sm:text-base">
                          {s.order ?? index + 1}
                        </div>
                        <div className="flex-1 min-w-0">
                          <h3 className="font-semibold text-[#2e4d57] mb-1 sm:mb-2 text-sm sm:text-base">
                            {s.title || s.text?.slice(0, 80) || t("procedure_detail.default_step")}
                          </h3>
                          <p className="text-[#1c3b2e] text-xs sm:text-sm leading-relaxed mb-2">{s.text || s.description || t("procedure_detail.no_step_description")}</p>
                          {(s.estimatedTime || s.time) && (
                            <span className="inline-flex items-center gap-1 text-xs font-medium text-[#3a6a8d] bg-[#3a6a8d]/10 px-2 py-1 rounded-full">
                              <Clock className="w-3 h-3" />
                              {s.estimatedTime || s.time}
                            </span>
                          )}
                        </div>
                      </div>
                    ))
                ) : (
                  <div className="text-sm text-[#6b7280] italic">{t("procedure_detail.empty_steps")}</div>
                )}
              </CardContent>
            </Card>

            <Tabs defaultValue="feedback" className="w-full">
              <TabsList className="grid w-full grid-cols-3 bg-white/90 backdrop-blur-sm border border-[#a7b3b9]/30 rounded-xl p-1">
                <TabsTrigger
                  value="feedback"
                  className="data-[state=active]:bg-gradient-to-r data-[state=active]:from-[#3a6a8d] data-[state=active]:to-[#2e4d57] data-[state=active]:text-white rounded-lg text-xs sm:text-sm"
                >
                  {t("procedure_detail.tabs.feedback")}
                </TabsTrigger>
                <TabsTrigger
                  value="notices"
                  className="data-[state=active]:bg-gradient-to-r data-[state=active]:from-[#3a6a8d] data-[state=active]:to-[#2e4d57] data-[state=active]:text-white rounded-lg text-xs sm:text-sm"
                >
                  {t("procedure_detail.tabs.notices")}
                </TabsTrigger>
                <TabsTrigger
                  value="discussions"
                  className="data-[state=active]:bg-gradient-to-r data-[state=active]:from-[#3a6a8d] data-[state=active]:to-[#2e4d57] data-[state=active]:text-white rounded-lg text-xs sm:text-sm"
                >
                  {t("procedure_detail.tabs.discussions")}
                </TabsTrigger>
              </TabsList>

              <TabsContent value="feedback" className="mt-4 md:mt-6">
                <Card className="bg-white/90 backdrop-blur-sm border-[#a7b3b9]/30 shadow-lg rounded-2xl">
                  <CardHeader>
                    <CardTitle className="text-[#2e4d57] text-base sm:text-lg">{t("procedure_detail.feedback_title")}</CardTitle>
                    <CardDescription className="text-[#1c3b2e] text-sm">{t("procedure_detail.feedback_description")}</CardDescription>
                  </CardHeader>
                  <CardContent className="space-y-3">
                    {loadingFeedback && <div className="text-center text-gray-500 py-4">{t("procedure_detail.feedback_loading")}</div>}
                    {!loadingFeedback && (feedbackData?.feedbacks?.length ?? 0) === 0 && (
                      <div className="text-center text-gray-500 py-4">{t("procedure_detail.feedback_empty")}</div>
                    )}
                    {(feedbackData?.feedbacks ?? []).map((f) => (
                      <div
                        key={f.id}
                        className="flex gap-3 p-3 rounded-lg bg-gradient-to-r from-[#a7b3b9]/10 to-[#5e9c8d]/10 border border-[#a7b3b9]/20"
                      >
                        <div className="w-8 h-8 rounded-full bg-[#ced4da] flex items-center justify-center text-xs font-semibold overflow-hidden">
                          <span>{f.userID?.charAt(0)?.toUpperCase() || t("procedure_detail.default_user_initial")}</span>
                        </div>
                        <div className="flex-1">
                          <div className="flex items-center gap-2 mb-1">
                            <span className="text-sm font-medium text-[#2e4d57]">
                              {t("procedure_detail.feedback_user", { id: f.userID?.slice(-4) || t("procedure_detail.anonymous_user") })}
                            </span>
                            <Badge className="text-xs px-1 py-0 bg-[#e6f4ff] text-[#1c3b2e]">
                              {t(`feedback.types.${f.type?.toLowerCase().replace(" ", "_")}`) || f.type?.replace("_", " ")}
                            </Badge>
                          </div>
                          <p className="text-xs text-[#1c3b2e] mb-1">{f.content || t("procedure_detail.no_feedback_content")}</p>
                          {f.tags?.length ? (
                            <div className="flex gap-1 mb-2 flex-wrap">
                              {f.tags.map((t, i) => (
                                <span key={i} className="text-xs bg-gray-200 text-gray-600 px-2 py-0.5 rounded">#{t}</span>
                              ))}
                            </div>
                          ) : null}
                          <div className="flex items-center gap-3 text-xs text-[#9ca3af]">
                            <span>
                              {f.createdAT ? new Date(f.createdAT).toLocaleDateString() : t("procedure_detail.recently")}
                            </span>
                          </div>
                        </div>
                      </div>
                    ))}
                    <Button asChild className="w-full bg-gradient-to-r from-[#3a6a8d] to-[#2e4d57] hover:from-[#2e4d57] hover:to-[#1c3b2e] text-white rounded-xl text-sm mt-2">
                      <a href={`/user/feedback?id=${encodeURIComponent(id)}`}>{t("procedure_detail.actions.add_feedback")}</a>
                    </Button>
                  </CardContent>
                </Card>
              </TabsContent>

              <TabsContent value="notices" className="mt-4 md:mt-6">
                <Card className="bg-white/90 backdrop-blur-sm border-[#a7b3b9]/30 shadow-lg rounded-2xl">
                  <CardHeader>
                    <CardTitle className="text-[#2e4d57] text-base sm:text-lg">{t("procedure_detail.notices_title")}</CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-3 md:space-y-4">
                    <div className="flex gap-3 p-3 md:p-4 bg-gradient-to-r from-[#5e9c8d]/10 to-[#1c3b2e]/10 border border-[#5e9c8d]/30 rounded-xl">
                      <FileText className="w-5 h-5 text-[#5e9c8d] flex-shrink-0 mt-0.5" />
                      <div>
                        <h4 className="font-semibold text-[#2e4d57] text-sm sm:text-base">{t("procedure_detail.no_notices")}</h4>
                        <p className="text-[#1c3b2e] text-xs sm:text-sm">{t("procedure_detail.no_notices_description")}</p>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </TabsContent>

              <TabsContent value="discussions" className="mt-4 md:mt-6">
                <Card className="bg-white/90 backdrop-blur-sm border-[#a7b3b9]/30 shadow-lg rounded-2xl">
                  <CardHeader>
                    <CardTitle className="text-[#2e4d57] text-base sm:text-lg">{t("procedure_detail.discussions_title")}</CardTitle>
                    <CardDescription className="text-[#1c3b2e] text-sm">{t("procedure_detail.discussions_description")}</CardDescription>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    {loadingDiscussions && <div className="text-center text-gray-500 py-6">{t("procedure_detail.discussions_loading")}</div>}
                    {!loadingDiscussions && (discussions?.posts?.length ?? 0) === 0 && (
                      <div className="text-center text-gray-500 py-6">{t("procedure_detail.discussions_empty")}</div>
                    )}
                    {(discussions?.posts ?? []).map((post) => (
                      <div
                        key={post.id}
                        className="bg-white border border-[#a7b3b9]/30 rounded-lg p-4 hover:bg-[#f8fafc] transition-all duration-300"
                      >
                        <div className="flex items-start gap-3">
                          <div className="w-10 h-10 rounded-full bg-[#3a6a8d] flex items-center justify-center text-white text-sm font-semibold overflow-hidden">
                            <span>{post.userID?.slice(-2).toUpperCase() || t("procedure_detail.default_user_initial")}</span>
                          </div>
                          <div className="flex-1">
                            <h4 className="font-medium text-[#2e4d57] mb-1 text-sm line-clamp-1">
                              {post.title || t("procedure_detail.default_discussion_title")}
                            </h4>
                            <p className="text-sm text-[#1c3b2e] mb-2 line-clamp-2">{post.content || t("procedure_detail.no_discussion_content")}</p>
                            {post.tags?.length ? (
                              <div className="flex items-center gap-2 mb-2 flex-wrap">
                                {post.tags.slice(0, 4).map((t, i) => (
                                  <Badge key={i} className="bg-[#dbeafe] text-[#1e40af] text-[10px] px-2 py-0.5">{t}</Badge>
                                ))}
                              </div>
                            ) : null}
                            <div className="flex items-center gap-4 text-xs text-[#9ca3af]">
                              <span>{post.createdAt ? new Date(post.createdAt).toLocaleDateString() : t("procedure_detail.recently")}</span>
                            </div>
                          </div>
                        </div>
                      </div>
                    ))}
                    <Button asChild className="w-full bg-gradient-to-r from-[#3a6a8d] to-[#2e4d57] hover:from-[#2e4d57] hover:to-[#1c3b2e] text-white rounded-xl text-sm">
                      <Link href="/user/discussions">{t("procedure_detail.actions.view_discussions")}</Link>
                    </Button>
                  </CardContent>
                </Card>
              </TabsContent>
            </Tabs>
          </div>

          <div className="lg:col-span-1 space-y-6">
            <Card className="bg-white/90 backdrop-blur-sm border-[#a7b3b9]/30 shadow-lg rounded-2xl">
              <CardHeader>
                <CardTitle className="flex items-center gap-2 text-[#2e4d57] text-base sm:text-lg">
                  <FileText className="w-5 h-5 text-[#5e9c8d]" />
                  {t("procedure_detail.documents_title")}
                </CardTitle>
              </CardHeader>
              <CardContent>
                <ul className="space-y-2">
                  {Array.isArray(procedure?.documentsRequired) && procedure!.documentsRequired.length > 0 ? (
                    procedure!.documentsRequired.map((doc, index) => (
                      <li key={index} className="flex items-center gap-2 text-[#1c3b2e] text-sm">
                        <CheckCircle className="w-4 h-4 text-[#5e9c8d] flex-shrink-0" />
                        {typeof doc === "string" ? doc : (doc?.name || t("procedure_detail.default_document"))}
                      </li>
                    ))
                  ) : (
                    <li className="text-sm text-[#6b7280] italic">{t("procedure_detail.empty_documents")}</li>
                  )}
                </ul>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  )
}

export default function ProcedureDetailPage() {
  const { t } = useTranslation("user")
  return (
    <Suspense fallback={<div className="p-6">{t("procedure_detail.loading")}</div>}>
      <ProcedureDetailInner />
    </Suspense>
  )
}
