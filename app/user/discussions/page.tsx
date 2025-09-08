// "use client"

// import { useEffect, useState } from "react"
// import { Search, Plus, MessageSquare, Filter, X, ChevronRight, ChevronLeft } from "lucide-react"
// import { Button } from "@/components/ui/button"
// import { Card, CardContent } from "@/components/ui/card"
// import { Badge } from "@/components/ui/badge"
// import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
// import { useRouter } from "next/navigation"
// import { useGetDiscussionsQuery } from "@/app/store/slices/discussionsSlice"
// import { motion, useReducedMotion } from "framer-motion"

// export default function CommunityPage() {
//   const [searchInput, setSearchInput] = useState("")
//   const [searchQuery, setSearchQuery] = useState("")
//   const [selectedCategory, setSelectedCategory] = useState("all")
//   const router = useRouter()

//   // added: pagination state
//   const [page, setPage] = useState(0)
//   const limit = 10
//   const { data, isLoading, isError } = useGetDiscussionsQuery({ page, limit })

//   // added: mobile sidebar toggle
//   const [mobilePanelOpen, setMobilePanelOpen] = useState(false)

//   // added: compute totalPages once for reuse (desktop + sticky bar)
//   const totalPages =
//     typeof data?.total === "number" && typeof data?.limit === "number" && data.limit > 0
//       ? Math.max(1, Math.ceil(data.total / data.limit))
//       : 1

//   useEffect(() => {
//     if (data) {
//       console.log("Discussions list:", data)
//     }
//   }, [data])

//   // no saved cards count here; show discussions API total instead

//   const [expandedMap, setExpandedMap] = useState<Record<string, boolean>>({})

//   // reset expanded per page
//   useEffect(() => {
//     setExpandedMap({})
//   }, [page])

//   const tagPillClasses = (i: number) => {
//     const styles = [
//       "bg-green-50 text-green-700 border-green-200 hover:bg-green-100 hover:text-green-800",
//       // replaced blue with amber
//       "bg-amber-50 text-amber-700 border-amber-200 hover:bg-amber-100 hover:text-amber-800",
//       "bg-teal-50 text-teal-700 border-teal-200 hover:bg-teal-100 hover:text-teal-800",
//       "bg-indigo-50 text-indigo-700 border-indigo-200 hover:bg-indigo-100 hover:text-indigo-800",
//       "bg-emerald-50 text-emerald-700 border-emerald-200 hover:bg-emerald-100 hover:text-emerald-800",
//       "bg-cyan-50 text-cyan-700 border-cyan-200 hover:bg-cyan-100 hover:text-cyan-800",
//     ]
//     return `cursor-pointer rounded-full ${styles[i % styles.length]}`
//   }

//   const apiDiscussions =
//     Array.isArray(data?.posts)
//       ? data!.posts.map((p) => ({
//           id: p.ID,
//           author: p.UserID || "User",
//           avatar: "/images/profile-photo.jpg",
//           timestamp: new Date(p.CreatedAt || p.UpdatedAt || Date.now()).toLocaleString(),
//           title: p.Title ?? "Untitled",
//           content: p.Content ?? "",
//           tags: Array.isArray(p.Tags) ? p.Tags.map((t) => String(t)) : [],
//         }))
//       : []

//   // only use backend data
//   const discussionsData = apiDiscussions

//   // build tags from backend posts for filters and fallback usage
//   const searchTags = (() => {
//     if (!Array.isArray(data?.posts)) return []
//     const set = new Set<string>()
//     data!.posts.forEach((p) => {
//       const tags = Array.isArray(p.Tags) ? p.Tags : []
//       tags.forEach((t) => {
//         const clean = String(t).replace(/^#/, "").trim()
//         if (clean) set.add(clean)
//       })
//     })
//     return Array.from(set)
//   })()

//   // helper to normalize tags for comparison
//   const normalizeTag = (s: string) => String(s || "").replace(/^#/, "").trim().toLowerCase()

//   // Map backend tag values to display labels to match create-post options
//   const displayTag = (raw: string) => {
//     const clean = String(raw || "").replace(/^#/, "").trim().toLowerCase()
//     if (clean === "business") return "National Id"
//     if (clean === "passport") return "passport"
//     if (clean === "tax") return "tax"
//     return String(raw || "").replace(/^#/, "").trim()
//   }

//   const popularTags = (() => {
//     const counts = new Map<string, number>()
//     discussionsData.forEach((d) =>
//       (d.tags || []).forEach((t: string) => {
//         const clean = String(t).replace(/^#/, "").trim()
//         if (clean) counts.set(clean, (counts.get(clean) || 0) + 1)
//       })
//     )
//     let list = Array.from(counts.entries()).map(([name, count]) => ({ name, count }))
//     // fallback to backend-provided tags only (no dummy tags)
//     if (!list.length) {
//       list = searchTags.map((name) => ({ name, count: 1 }))
//     }
//     return list.sort((a, b) => b.count - a.count).slice(0, 20)
//   })()

//   const filteredDiscussions = discussionsData.filter((discussion) => {
//     const matchesSearch =
//       discussion.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
//       discussion.content.toLowerCase().includes(searchQuery.toLowerCase())
//     // updated: exact tag match using backend tags (no hardcoded categories)
//     const matchesCategory =
//       selectedCategory === "all" ||
//       (Array.isArray(discussion.tags) &&
//         discussion.tags.some((tag: string) => normalizeTag(tag) === normalizeTag(selectedCategory)))
//     return matchesSearch && matchesCategory
//   })

//   const prefersReducedMotion = useReducedMotion()
//   const itemVariants = prefersReducedMotion
//     ? { hidden: { opacity: 0 }, visible: { opacity: 1, transition: { duration: 0.15 } } }
//     : {
//         hidden: { opacity: 0, y: 12, scale: 0.985, filter: "blur(0.2px)" },
//         visible: (i: number) => ({
//           opacity: 1,
//           y: 0,
//           scale: 1,
//           filter: "none",
//           transition: {
//             type: "spring" as const,
//             stiffness: 220,
//             damping: 18,
//             mass: 0.9,
//             delay: i * 0.05 + 0.06,
//           },
//         }),
//       }

//   return (
//     <div className="min-h-screen bg-gray-50 flex flex-col page-fade">
//       {/* animations (scoped) */}
//       <style jsx>{`
//         .page-fade { animation: fadeIn .45s ease-out both; }
//         .fade-in-up { animation: fadeInUp .5s ease-out both; }
//         .slide-up { animation: slideUp .4s ease-out both; }
//         .card-tilt { transition: transform .25s ease, box-shadow .25s ease; }
//         .card-tilt:hover { transform: translateY(-3px); box-shadow: 0 12px 28px rgba(0,0,0,.08); }
//         @keyframes fadeIn { from { opacity: 0 } to { opacity: 1 } }
//         @keyframes fadeInUp { from { opacity: 0; transform: translateY(10px) } to { opacity: 1; transform: translateY(0) } }
//         @keyframes slideUp { from { transform: translateY(16px); opacity: .6 } to { transform: translateY(0); opacity: 1 } }
//         @media (prefers-reduced-motion: reduce) {
//           .page-fade, .fade-in-up, .slide-up, .card-enter, .card-enter-left, .card-enter-right { animation: none !important; opacity: 1 !important; transform: none !important; filter: none !important; }
//           .card-tilt, .card-tilt:hover { transform: none !important; box-shadow: none !important; }
//         }
//       `}</style>

//       <div className="max-w-7xl mx-auto w-full flex-1">
//         {/* Header */}
//         <div className="bg-white border border-gray-100 rounded-xl p-4 sm:p-5 mb-4 sm:mb-6 shadow-sm fade-in-up">
//           <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
//             <div className="w-full">
//               <div className="flex items-center gap-3 mb-1">
//                 <MessageSquare className="h-7 w-7" style={{ color: '#3a6a8d' }} />
//                 <h1 className="text-xl leading-snug sm:text-3xl font-bold text-[#111827]">
//                   Community Discussions
//                 </h1>
//               </div>
//               <p className="text-[#4b5563] text-xs sm:text-sm md:text-base">
//                 Join the conversation. Share, ask, and collaborate.
//               </p>
//             </div>
//             <div className="hidden sm:flex gap-2">
//               <Button
//                 variant="outline"
//                 className="border-[#3A6A8D] text-[#3A6A8D] hover:bg-[#3A6A8D]/10"
//                 onClick={() => router.push("/user/my-discussions")}
//               >
//                 My Discussions
//               </Button>
//               <Button
//                 className="bg-gradient-to-r from-[#3a6a8d] to-[#2e4d57] hover:from-[#2e4d57] hover:to-[#1c3b2e] text-white"
//                 onClick={() => router.push("/user/create-post")}
//               >
//                 <Plus className="h-4 w-4 mr-2" />
//                 Add Discussion
//               </Button>
//             </div>
//             {/* Mobile quick button to open filters */}
//             <div className="w-full sm:hidden flex gap-2">
//               <Button
//                 variant="outline"
//                 className="flex-1 border-[#3A6A8D] text-[#3A6A8D] hover:bg-[#3A6A8D]/10"
//                 onClick={() => setMobilePanelOpen(true)}
//               >
//                 <Filter className="h-4 w-4 mr-2" />
//                 Filters & Tags
//               </Button>
//             </div>
//           </div>
//         </div>

//         {/* Search & Filters (desktop / tablet) */}
//         <Card className="bg-white border border-gray-100 rounded-xl p-4 mb-4 sm:mb-6 hidden sm:block shadow-sm fade-in-up" style={{ animationDelay: "80ms" }}>
//           <div className="flex flex-col gap-4 w-full mb-2 sm:flex-row">
//             <div className="relative flex-1 flex">
//               {/* ...existing code... */}
//               <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
//               <input
//                 type="text"
//                 // ...existing props...
//                 placeholder="Search discussions..."
//                 value={searchInput}
//                 onChange={(e) => setSearchInput(e.target.value)}
//                 className="w-full pl-10 pr-4 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#3A6A8D] focus:border-transparent"
//               />
//               <Button
//                 type="button"
//                 className="ml-2 px-4 py-2 bg-gradient-to-r from-[#3a6a8d] to-[#2e4d57] hover:from-[#2e4d57] hover:to-[#1c3b2e] text-white"
//                 onClick={() => setSearchQuery(searchInput)}
//               >
//                 Search
//               </Button>
//             </div>
//             {/* Categories only (backend tags); fixed width so search expands */}
//             <div className="flex gap-2 flex-none w-56">
//               <Select value={selectedCategory} onValueChange={setSelectedCategory}>
//                 <SelectTrigger className="w-full border-[#3A6A8D] text-[#3A6A8D] focus:ring-2 focus:ring-[#3A6A8D] focus:border-transparent">
//                   <SelectValue placeholder="All Categories" />
//                 </SelectTrigger>
//                 <SelectContent>
//                   <SelectItem value="all">All Categories</SelectItem>
//                   {searchTags.map((tag) => (
//                     <SelectItem key={tag} value={tag}>
//                       {displayTag(tag)}
//                     </SelectItem>
//                   ))}
//                 </SelectContent>
//               </Select>
//               {/* removed Sort/Latest select */}
//             </div>
//           </div>
//           {searchTags.length > 0 && (
//             <div className="relative z-10 flex gap-2 overflow-x-auto pb-1 scrollbar-thin scrollbar-track-transparent scrollbar-thumb-gray-300">
//               {searchTags.map((tag, i) => (
//                 <Badge key={tag} variant="outline" className={`${tagPillClasses(i)} flex-shrink-0`}>
//                   {displayTag(tag)}
//                 </Badge>
//               ))}
//             </div>
//           )}
//         </Card>

//         {/* Mobile slide-over panel */}
//         {mobilePanelOpen && (
//           <div className="fixed inset-0 z-50 sm:hidden">
//             <div
//               className="absolute inset-0 bg-black/40 backdrop-blur-sm"
//               onClick={() => setMobilePanelOpen(false)}
//             />
//             <div className="absolute bottom-0 left-0 right-0 bg-white rounded-t-2xl shadow-xl max-h-[85vh] flex flex-col slide-up">
//               <div className="flex items-center justify-between px-4 pt-4 pb-2 border-b">
//                 <h2 className="text-base font-semibold text-gray-800">Filters & Tags</h2>
//                 <Button variant="ghost" size="sm" onClick={() => setMobilePanelOpen(false)}>
//                   <X className="h-5 w-5" />
//                 </Button>
//               </div>
//               <div className="p-4 space-y-5 overflow-y-auto">
//                 {/* Mobile search */}
//                 <div>
//                   <div className="relative">
//                     <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
//                     <input
//                       type="text"
//                       placeholder="Search discussions..."
//                       value={searchInput}
//                       onChange={(e) => setSearchInput(e.target.value)}
//                       className="w-full pl-10 pr-3 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#3A6A8D]"
//                     />
//                   </div>
//                   <Button
//                     className="mt-2 w-full bg-[#3A6A8D] hover:bg-[#2d5470] text-white"
//                     onClick={() => {
//                       setSearchQuery(searchInput)
//                       setMobilePanelOpen(false)
//                     }}
//                   >
//                     Apply
//                   </Button>
//                 </div>
//                 {/* Category select (mobile) - backend tags + style like My Discussions */}
//                 <div className="space-y-2">
//                   <label className="text-xs uppercase tracking-wide text-gray-500 font-medium">
//                     Category
//                   </label>
//                   <Select value={selectedCategory} onValueChange={setSelectedCategory}>
//                     <SelectTrigger className="w-full border-[#3A6A8D] text-[#3A6A8D] hover:bg-[#3A6A8D]/5 focus:ring-2 focus:ring-[#3A6A8D] focus:border-transparent">
//                       <SelectValue placeholder="All Categories" />
//                     </SelectTrigger>
//                     <SelectContent>
//                       <SelectItem value="all">All Categories</SelectItem>
//                       {searchTags.map((tag) => (
//                         <SelectItem key={tag} value={tag}>
//                           {displayTag(tag)}
//                         </SelectItem>
//                       ))}
//                     </SelectContent>
//                   </Select>
//                 </div>
//                 {/* Tags */}
//                 {searchTags.length > 0 && (
//                   <div className="space-y-2">
//                     <label className="text-xs uppercase tracking-wide text-gray-500 font-medium">
//                       Popular Tags
//                     </label>
//                     <div className="flex flex-wrap gap-2">
//             {searchTags.map((tag, i) => (
//                         <Badge
//                           key={tag}
//                           variant="outline"
//                           className={`${tagPillClasses(i)} px-3 py-1`}
//                         >
//               {displayTag(tag)}
//                         </Badge>
//                       ))}
//                     </div>
//                   </div>
//                 )}
//                 <div className="pt-2">
//                   <Button
//                     variant="outline"
//                     className="w-full"
//                     onClick={() => {
//                       setSelectedCategory("all")
//                       setSearchInput("")
//                       setSearchQuery("")
//                     }}
//                   >
//                     Reset
//                   </Button>
//                 </div>
//               </div>
//             </div>
//           </div>
//         )}

//         <div className="grid grid-cols-1 lg:grid-cols-12 gap-4 sm:gap-6">
//           {/* Main Content */}
//           <div className="lg:col-span-9 space-y-4 sm:space-y-6">
//             {isLoading && (
//               <div className="space-y-4">
//                 {[0,1,2].map(i=>(
//                   <div key={i} className="border border-gray-100 rounded-xl bg-white p-4 sm:p-6 shadow-sm animate-pulse">
//                     <div className="flex gap-4">
//                       <div className="h-10 w-10 sm:h-12 sm:w-12 rounded-full bg-gray-200" />
//                       <div className="flex-1 space-y-3">
//                         <div className="h-4 w-1/2 bg-gray-200 rounded" />
//                         <div className="h-3 w-full bg-gray-100 rounded" />
//                         <div className="h-3 w-5/6 bg-gray-100 rounded" />
//                       </div>
//                     </div>
//                   </div>
//                 ))}
//               </div>
//             )}

//             {!isLoading && (
//               <div className="space-y-4 sm:space-y-6">
//                 {filteredDiscussions.map((discussion, index) => {
//                   const rowKey = `${discussion.title}-${index}`
//                   const isExpanded = !!expandedMap[rowKey]
//                   return (
//                     <motion.div
//                       key={rowKey}
//                       variants={itemVariants}
//                       initial="hidden"
//                       whileInView="visible"
//                       viewport={{ once: true, amount: 0.2 }}
//                       custom={index}
//                     >
//                       <Card className="group bg-white rounded-2xl border border-[#e5e7eb] shadow-xl relative overflow-hidden ring-1 ring-transparent hover:ring-[#3a6a8d]/20 transition-all duration-300 transform-gpu hover:-translate-y-0.5 hover:shadow-2xl p-3 sm:p-6">
//                         <div className="absolute inset-0 bg-gradient-to-r from-[#3a6a8d]/10 via-transparent to-[#5e9c8d]/10 opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
//                         <CardContent className="relative z-10 p-0">
//                           <div className="flex gap-3 sm:gap-4 flex-col sm:flex-row">
//                             <div className="w-11 h-11 sm:w-12 sm:h-12 rounded-xl bg-[#f3f4f6] flex items-center justify-center mx-auto sm:mx-0 ring-2 ring-white shadow">
//                               <MessageSquare className="w-5 h-5" style={{ color: '#3a6a8d' }} />
//                             </div>
//                             <div className="flex-1">
//                               <div className="flex items-center gap-2 mb-2">
//                                 <h3 className="text-sm sm:text-base font-semibold text-[#111827]">
//                                   {discussion.title}
//                                 </h3>
//                                 <Badge variant="secondary" className="text-[10px] sm:text-xs bg-gray-100 text-gray-700 border border-gray-200">
//                                   Anonymous
//                                 </Badge>
//                               </div>
//                               <p className={`text-[#374151] text-sm sm:text-[15px] mb-3 ${isExpanded ? "" : "line-clamp-2"}`}>
//                                 {discussion.content}
//                               </p>
//                               <div className="flex flex-wrap gap-2">
//                                 {discussion.tags.map((tag: string, i: number) => {
//                                   const clean = tag.replace(/^#/, "")
//                                   return (
//                                     <Badge
//                                       key={`${discussion.title}-${clean}-${i}`}
//                                       variant="outline"
//                                       className={`text-[10px] sm:text-xs ${tagPillClasses(i)}`}
//                                     >
//                                       {displayTag(tag)}
//                                     </Badge>
//                                   )
//                                 })}
//                               </div>
//                               <div className="flex justify-end pt-1 sm:pt-3">
//                                 <Button
//                                   variant="ghost"
//                                   size="sm"
//                                   className="text-[#3A6A8D] hover:bg-[#3A6A8D]/10 hover:text-[#2d5470] text-xs sm:text-sm"
//                                   onClick={(e) => {
//                                     e.stopPropagation()
//                                     setExpandedMap((prev) => ({ ...prev, [rowKey]: !prev[rowKey] }))
//                                   }}
//                                 >
//                                   {isExpanded ? "View Less" : "View More"}
//                                 </Button>
//                               </div>
//                             </div>
//                           </div>
//                         </CardContent>
//                       </Card>
//                     </motion.div>
//                   )
//                 })}
//                 {filteredDiscussions.length === 0 && !isLoading && (
//                   <div className="text-center text-sm text-gray-500 py-10 border border-dashed rounded-lg fade-in-up" style={{ animationDelay: "120ms" }}>
//                     No discussions match your filters.
//                   </div>
//                 )}
//               </div>
//             )}

//             {/* Desktop pagination footer */}
//             {totalPages > 1 && (
//               <div className="hidden sm:flex mt-2 bg-white/90 border border-gray-100 rounded-xl px-3 py-3 items-center justify-between shadow-sm fade-in-up" style={{ animationDelay: "140ms" }}>
//                 <div className="text-sm text-gray-600">
//                   Page {page + 1} of {totalPages}
//                 </div>
//                 <div className="flex gap-2">
//                   <Button
//                     variant="outline"
//                     className="border-gray-300"
//                     disabled={page <= 0 || isLoading}
//                     onClick={() => setPage((p) => Math.max(0, p - 1))}
//                   >
//                     Previous
//                   </Button>
//                   <Button
//                     className="bg-[#3A6A8D] hover:bg-[#2d5470] text-white"
//                     disabled={page >= totalPages - 1 || isLoading}
//                     onClick={() => setPage((p) => p + 1)}
//                   >
//                     Next
//                   </Button>
//                 </div>
//               </div>
//             )}
//           </div>

//           {/* Sidebar (desktop only) */}
//           <div className="lg:col-span-3 space-y-6 hidden lg:block fade-in-up" style={{ animationDelay: "100ms" }}>
//             <Card className="group relative overflow-hidden p-4 bg-white border border-[#e5e7eb] rounded-2xl shadow-xl hover:shadow-2xl transform transition duration-300 hover:-translate-y-0.5">
//               <div className="absolute inset-0 bg-gradient-to-r from-[#3a6a8d]/10 via-transparent to-[#5e9c8d]/10 opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
//               <div className="relative z-10">
//               <h3 className="font-semibold text-gray-900 mb-4 text-sm sm:text-base">Popular Tags</h3>
//               <div className="space-y-2">
//                 {popularTags.map((tag, i) => (
//                   <div key={tag.name} className="flex items-center justify-between text-sm">
//                     <Badge variant="outline" className={tagPillClasses(i)}>
//                       {displayTag(tag.name)}
//                     </Badge>
//                     <span className="text-xs text-gray-500">{tag.count}</span>
//                   </div>
//                 ))}
//               </div>
//               </div>
//             </Card>
//           </div>
//         </div>

//         {/* Status / meta */}
//         <div className="mt-6 text-center sm:text-left fade-in-up" style={{ animationDelay: "160ms" }}>
//           {isLoading && <div className="text-sm text-gray-500">Loading...</div>}
//           {isError && <div className="text-sm text-red-600">Failed to load discussions.</div>}
//           {!isLoading && !isError && data && (
//             <div className="text-xs sm:text-sm text-gray-600">Total: {data.total}</div>
//           )}
//         </div>
//       </div>

//       {/* Sticky mobile action bar */}
//       <div className="sm:hidden fixed bottom-0 left-0 right-0 z-40 bg-white/95 backdrop-blur border-t shadow-lg px-2 py-2 flex items-center justify-between fade-in-up" style={{ animationDelay: "180ms" }}>
//         <Button
//           variant="outline"
//           size="sm"
//           className="flex-1 mx-1 text-[#3A6A8D] border-[#3A6A8D]"
//           onClick={() => router.push("/user/my-discussions")}
//         >
//           Me
//         </Button>
//         <Button
//           size="sm"
//           className="flex-1 mx-1 bg-[#3A6A8D] hover:bg-[#2d5470] text-white"
//           onClick={() => router.push("/user/create-post")}
//         >
//           <Plus className="h-4 w-4 mr-1" /> New
//         </Button>
//         {totalPages > 1 && (
//           <div className="flex flex-1 mx-1 gap-1">
//             <Button
//               variant="outline"
//               size="sm"
//               className="flex-1"
//               disabled={page <= 0}
//               onClick={() => setPage((p) => Math.max(0, p - 1))}
//             >
//               <ChevronLeft className="h-4 w-4" />
//             </Button>
//             <Button
//               variant="outline"
//               size="sm"
//               className="flex-1"
//               disabled={page >= totalPages - 1}
//               onClick={() => setPage((p) => p + 1)}
//             >
//               <ChevronRight className="h-4 w-4" />
//             </Button>
//           </div>
//         )}
//       </div>
//     </div>
//   )
// }


"use client"

import { useEffect, useState } from "react"
import { Search, Plus, MessageSquare, Filter, X, ChevronRight, ChevronLeft } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { useRouter } from "next/navigation"
import { useGetDiscussionsQuery } from "@/app/store/slices/discussionsSlice"
import { motion, useReducedMotion } from "framer-motion"
import { useTranslation } from "react-i18next"

export default function CommunityPage() {
  const { t } = useTranslation("user")
  const [searchInput, setSearchInput] = useState("")
  const [searchQuery, setSearchQuery] = useState("")
  const [selectedCategory, setSelectedCategory] = useState("all")
  const router = useRouter()

  const [page, setPage] = useState(0)
  const limit = 10
  const { data, isLoading, isError } = useGetDiscussionsQuery({ page, limit })

  const [mobilePanelOpen, setMobilePanelOpen] = useState(false)

  const totalPages =
    typeof data?.total === "number" && typeof data?.limit === "number" && data.limit > 0
      ? Math.max(1, Math.ceil(data.total / data.limit))
      : 1

  useEffect(() => {
    if (data) {
      console.log("Discussions list:", data)
    }
  }, [data])

  const [expandedMap, setExpandedMap] = useState<Record<string, boolean>>({})

  useEffect(() => {
    setExpandedMap({})
  }, [page])

  const tagPillClasses = (i: number) => {
    const styles = [
      "bg-green-50 text-green-700 border-green-200 hover:bg-green-100 hover:text-green-800",
      "bg-amber-50 text-amber-700 border-amber-200 hover:bg-amber-100 hover:text-amber-800",
      "bg-teal-50 text-teal-700 border-teal-200 hover:bg-teal-100 hover:text-teal-800",
      "bg-indigo-50 text-indigo-700 border-indigo-200 hover:bg-indigo-100 hover:text-indigo-800",
      "bg-emerald-50 text-emerald-700 border-emerald-200 hover:bg-emerald-100 hover:text-emerald-800",
      "bg-cyan-50 text-cyan-700 border-cyan-200 hover:bg-cyan-100 hover:text-cyan-800",
    ]
    return `cursor-pointer rounded-full ${styles[i % styles.length]}`
  }

  const apiDiscussions =
    Array.isArray(data?.posts)
      ? data.posts.map((p, idx) => {
          if (!p) {
            console.error("Post is undefined at index:", idx)
            return null
          }
          return {
            id: p.ID,
            author: p.UserID || t("community.default_author"),
            avatar: "/images/profile-photo.jpg",
            timestamp: new Date(p.CreatedAt || p.UpdatedAt || Date.now()).toLocaleString(),
            title: p.Title ?? t("community.default_title"),
            content: p.Content ?? "",
            tags: Array.isArray(p.Tags) ? p.Tags.map((t) => String(t)) : [],
          }
        }).filter((p): p is NonNullable<typeof p> => p !== null)
      : []

  const discussionsData = apiDiscussions

  const searchTags = (() => {
    if (!Array.isArray(data?.posts)) return []
    const set = new Set<string>()
    data.posts.forEach((p, idx) => {
      if (!p) {
        console.error("Post is undefined at index:", idx)
        return
      }
      const tags = Array.isArray(p.Tags) ? p.Tags : []
      tags.forEach((t) => {
        const clean = String(t).replace(/^#/, "").trim()
        if (clean) set.add(clean)
      })
    })
    return Array.from(set)
  })()

  const normalizeTag = (s: string) => String(s || "").replace(/^#/, "").trim().toLowerCase()

  const displayTag = (raw: string) => {
    const clean = String(raw || "").replace(/^#/, "").trim().toLowerCase()
    if (clean === "business") return t("community.tags.national_id")
    if (clean === "passport") return t("community.tags.passport")
    if (clean === "tax") return t("community.tags.tax")
    return String(raw || "").replace(/^#/, "").trim()
  }

  const popularTags = (() => {
    const counts = new Map<string, number>()
    discussionsData.forEach((d) =>
      (d.tags || []).forEach((t: string) => {
        const clean = String(t).replace(/^#/, "").trim()
        if (clean) counts.set(clean, (counts.get(clean) || 0) + 1)
      })
    )
    let list = Array.from(counts.entries()).map(([name, count]) => ({ name, count }))
    if (!list.length) {
      list = searchTags.map((name) => ({ name, count: 1 }))
    }
    return list.sort((a, b) => b.count - a.count).slice(0, 20)
  })()

  const filteredDiscussions = discussionsData.filter((discussion) => {
    const matchesSearch =
      discussion.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
      discussion.content.toLowerCase().includes(searchQuery.toLowerCase())
    const matchesCategory =
      selectedCategory === "all" ||
      (Array.isArray(discussion.tags) &&
        discussion.tags.some((tag: string) => normalizeTag(tag) === normalizeTag(selectedCategory)))
    return matchesSearch && matchesCategory
  })

  const prefersReducedMotion = useReducedMotion()
  const itemVariants = prefersReducedMotion
    ? { hidden: { opacity: 0 }, visible: { opacity: 1, transition: { duration: 0.15 } } }
    : {
        hidden: { opacity: 0, y: 12, scale: 0.985, filter: "blur(0.2px)" },
        visible: (i: number) => ({
          opacity: 1,
          y: 0,
          scale: 1,
          filter: "none",
          transition: {
            type: "spring" as const,
            stiffness: 220,
            damping: 18,
            mass: 0.9,
            delay: i * 0.05 + 0.06,
          },
        }),
      }

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col page-fade">
      <style jsx>{`
        .page-fade { animation: fadeIn .45s ease-out both; }
        .fade-in-up { animation: fadeInUp .5s ease-out both; }
        .slide-up { animation: slideUp .4s ease-out both; }
        .card-tilt { transition: transform .25s ease, box-shadow .25s ease; }
        .card-tilt:hover { transform: translateY(-3px); box-shadow: 0 12px 28px rgba(0,0,0,.08); }
        @keyframes fadeIn { from { opacity: 0 } to { opacity: 1 } }
        @keyframes fadeInUp { from { opacity: 0; transform: translateY(10px) } to { opacity: 1; transform: translateY(0) } }
        @keyframes slideUp { from { transform: translateY(16px); opacity: .6 } to { transform: translateY(0); opacity: 1 } }
        @media (prefers-reduced-motion: reduce) {
          .page-fade, .fade-in-up, .slide-up, .card-enter, .card-enter-left, .card-enter-right { animation: none !important; opacity: 1 !important; transform: none !important; filter: none !important; }
          .card-tilt, .card-tilt:hover { transform: none !important; box-shadow: none !important; }
        }
      `}</style>

      <div className="max-w-7xl mx-auto w-full flex-1">
        {/* Header */}
        <div className="bg-white border border-gray-100 rounded-xl p-4 sm:p-5 mb-4 sm:mb-6 shadow-sm fade-in-up">
          <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
            <div className="w-full">
              <div className="flex items-center gap-3 mb-1">
                <MessageSquare className="h-7 w-7" style={{ color: '#3a6a8d' }} />
                <h1 className="text-xl leading-snug sm:text-3xl font-bold text-[#111827]">
                  {t("community.title")}
                </h1>
              </div>
              <p className="text-[#4b5563] text-xs sm:text-sm md:text-base">
                {t("community.description")}
              </p>
            </div>
            <div className="hidden sm:flex gap-2">
              <Button
                variant="outline"
                className="border-[#3A6A8D] text-[#3A6A8D] hover:bg-[#3A6A8D]/10"
                onClick={() => router.push("/user/my-discussions")}
              >
                {t("community.actions.my_discussions")}
              </Button>
              <Button
                className="bg-gradient-to-r from-[#3a6a8d] to-[#2e4d57] hover:from-[#2e4d57] hover:to-[#1c3b2e] text-white"
                onClick={() => router.push("/user/create-post")}
              >
                <Plus className="h-4 w-4 mr-2" />
                {t("community.actions.add_discussion")}
              </Button>
            </div>
            <div className="w-full sm:hidden flex gap-2">
              <Button
                variant="outline"
                className="flex-1 border-[#3A6A8D] text-[#3A6A8D] hover:bg-[#3A6A8D]/10"
                onClick={() => setMobilePanelOpen(true)}
              >
                <Filter className="h-4 w-4 mr-2" />
                {t("community.actions.filters")}
              </Button>
            </div>
          </div>
        </div>

        {/* Search & Filters (desktop / tablet) */}
        <Card className="bg-white border border-gray-100 rounded-xl p-4 mb-4 sm:mb-6 hidden sm:block shadow-sm fade-in-up" style={{ animationDelay: "80ms" }}>
          <div className="flex flex-col gap-4 w-full mb-2 sm:flex-row">
            <div className="relative flex-1 flex">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
              <input
                type="text"
                placeholder={t("community.search.placeholder")}
                value={searchInput}
                onChange={(e) => setSearchInput(e.target.value)}
                className="w-full pl-10 pr-4 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#3A6A8D] focus:border-transparent"
              />
              <Button
                type="button"
                className="ml-2 px-4 py-2 bg-gradient-to-r from-[#3a6a8d] to-[#2e4d57] hover:from-[#2e4d57] hover:to-[#1c3b2e] text-white"
                onClick={() => setSearchQuery(searchInput)}
              >
                {t("community.search.button")}
              </Button>
            </div>
            <div className="flex gap-2 flex-none w-56">
              <Select value={selectedCategory} onValueChange={setSelectedCategory}>
                <SelectTrigger className="w-full border-[#3A6A8D] text-[#3A6A8D] focus:ring-2 focus:ring-[#3A6A8D] focus:border-transparent">
                  <SelectValue placeholder={t("community.filter.all")} />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">{t("community.filter.all")}</SelectItem>
                  {searchTags.map((tag) => (
                    <SelectItem key={tag} value={tag}>
                      {displayTag(tag)}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>
          {searchTags.length > 0 && (
            <div className="relative z-10 flex gap-2 overflow-x-auto pb-1 scrollbar-thin scrollbar-track-transparent scrollbar-thumb-gray-300">
              {searchTags.map((tag, i) => (
                <Badge key={tag} variant="outline" className={`${tagPillClasses(i)} flex-shrink-0`}>
                  {displayTag(tag)}
                </Badge>
              ))}
            </div>
          )}
        </Card>

        {/* Mobile slide-over panel */}
        {mobilePanelOpen && (
          <div className="fixed inset-0 z-50 sm:hidden">
            <div
              className="absolute inset-0 bg-black/40 backdrop-blur-sm"
              onClick={() => setMobilePanelOpen(false)}
            />
            <div className="absolute bottom-0 left-0 right-0 bg-white rounded-t-2xl shadow-xl max-h-[85vh] flex flex-col slide-up">
              <div className="flex items-center justify-between px-4 pt-4 pb-2 border-b">
                <h2 className="text-base font-semibold text-gray-800">{t("community.filters.title")}</h2>
                <Button variant="ghost" size="sm" onClick={() => setMobilePanelOpen(false)}>
                  <X className="h-5 w-5" />
                </Button>
              </div>
              <div className="p-4 space-y-5 overflow-y-auto">
                <div>
                  <div className="relative">
                    <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
                    <input
                      type="text"
                      placeholder={t("community.search.placeholder")}
                      value={searchInput}
                      onChange={(e) => setSearchInput(e.target.value)}
                      className="w-full pl-10 pr-3 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#3A6A8D]"
                    />
                  </div>
                  <Button
                    className="mt-2 w-full bg-[#3A6A8D] hover:bg-[#2d5470] text-white"
                    onClick={() => {
                      setSearchQuery(searchInput)
                      setMobilePanelOpen(false)
                    }}
                  >
                    {t("community.filters.apply")}
                  </Button>
                </div>
                <div className="space-y-2">
                  <label className="text-xs uppercase tracking-wide text-gray-500 font-medium">
                    {t("community.filters.category")}
                  </label>
                  <Select value={selectedCategory} onValueChange={setSelectedCategory}>
                    <SelectTrigger className="w-full border-[#3A6A8D] text-[#3A6A8D] hover:bg-[#3A6A8D]/5 focus:ring-2 focus:ring-[#3A6A8D] focus:border-transparent">
                      <SelectValue placeholder={t("community.filter.all")} />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="all">{t("community.filter.all")}</SelectItem>
                      {searchTags.map((tag) => (
                        <SelectItem key={tag} value={tag}>
                          {displayTag(tag)}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
                {searchTags.length > 0 && (
                  <div className="space-y-2">
                    <label className="text-xs uppercase tracking-wide text-gray-500 font-medium">
                      {t("community.filters.popular_tags")}
                    </label>
                    <div className="flex flex-wrap gap-2">
                      {searchTags.map((tag, i) => (
                        <Badge
                          key={tag}
                          variant="outline"
                          className={`${tagPillClasses(i)} px-3 py-1`}
                        >
                          {displayTag(tag)}
                        </Badge>
                      ))}
                    </div>
                  </div>
                )}
                <div className="pt-2">
                  <Button
                    variant="outline"
                    className="w-full"
                    onClick={() => {
                      setSelectedCategory("all")
                      setSearchInput("")
                      setSearchQuery("")
                    }}
                  >
                    {t("community.filters.reset")}
                  </Button>
                </div>
              </div>
            </div>
          </div>
        )}

        <div className="grid grid-cols-1 lg:grid-cols-12 gap-4 sm:gap-6">
          <div className="lg:col-span-9 space-y-4 sm:space-y-6">
            {isLoading && (
              <div className="space-y-4">
                {[0,1,2].map(i=>(
                  <div key={i} className="border border-gray-100 rounded-xl bg-white p-4 sm:p-6 shadow-sm animate-pulse">
                    <div className="flex gap-4">
                      <div className="h-10 w-10 sm:h-12 sm:w-12 rounded-full bg-gray-200" />
                      <div className="flex-1 space-y-3">
                        <div className="h-4 w-1/2 bg-gray-200 rounded" />
                        <div className="h-3 w-full bg-gray-100 rounded" />
                        <div className="h-3 w-5/6 bg-gray-100 rounded" />
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}

            {!isLoading && (
              <div className="space-y-4 sm:space-y-6">
                {filteredDiscussions.map((discussion, index) => {
                  const rowKey = `${discussion.title}-${index}`
                  const isExpanded = !!expandedMap[rowKey]
                  return (
                    <motion.div
                      key={rowKey}
                      variants={itemVariants}
                      initial="hidden"
                      whileInView="visible"
                      viewport={{ once: true, amount: 0.2 }}
                      custom={index}
                    >
                      <Card className="group bg-white rounded-2xl border border-[#e5e7eb] shadow-xl relative overflow-hidden ring-1 ring-transparent hover:ring-[#3a6a8d]/20 transition-all duration-300 transform-gpu hover:-translate-y-0.5 hover:shadow-2xl p-3 sm:p-6">
                        <div className="absolute inset-0 bg-gradient-to-r from-[#3a6a8d]/10 via-transparent to-[#5e9c8d]/10 opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
                        <CardContent className="relative z-10 p-0">
                          <div className="flex gap-3 sm:gap-4 flex-col sm:flex-row">
                            <div className="w-11 h-11 sm:w-12 sm:h-12 rounded-xl bg-[#f3f4f6] flex items-center justify-center mx-auto sm:mx-0 ring-2 ring-white shadow">
                              <MessageSquare className="w-5 h-5" style={{ color: '#3a6a8d' }} />
                            </div>
                            <div className="flex-1">
                              <div className="flex items-center gap-2 mb-2">
                                <h3 className="text-sm sm:text-base font-semibold text-[#111827]">
                                  {discussion.title}
                                </h3>
                                <Badge variant="secondary" className="text-[10px] sm:text-xs bg-gray-100 text-gray-700 border border-gray-200">
                                  {t("community.author_anonymous")}
                                </Badge>
                              </div>
                              <p className={`text-[#374151] text-sm sm:text-[15px] mb-3 ${isExpanded ? "" : "line-clamp-2"}`}>
                                {discussion.content}
                              </p>
                              <div className="flex flex-wrap gap-2">
                                {discussion.tags.map((tag: string, i: number) => {
                                  const clean = tag.replace(/^#/, "")
                                  return (
                                    <Badge
                                      key={`${discussion.title}-${clean}-${i}`}
                                      variant="outline"
                                      className={`text-[10px] sm:text-xs ${tagPillClasses(i)}`}
                                    >
                                      {displayTag(tag)}
                                    </Badge>
                                  )
                                })}
                              </div>
                              <div className="flex justify-end pt-1 sm:pt-3">
                                <Button
                                  variant="ghost"
                                  size="sm"
                                  className="text-[#3A6A8D] hover:bg-[#3A6A8D]/10 hover:text-[#2d5470] text-xs sm:text-sm"
                                  onClick={(e) => {
                                    e.stopPropagation()
                                    setExpandedMap((prev) => ({ ...prev, [rowKey]: !prev[rowKey] }))
                                  }}
                                >
                                  {isExpanded ? t("community.actions.view_less") : t("community.actions.view_more")}
                                </Button>
                              </div>
                            </div>
                          </div>
                        </CardContent>
                      </Card>
                    </motion.div>
                  )
                })}
                {filteredDiscussions.length === 0 && !isLoading && (
                  <div className="text-center text-sm text-gray-500 py-10 border border-dashed rounded-lg fade-in-up" style={{ animationDelay: "120ms" }}>
                    {t("community.empty")}
                  </div>
                )}
              </div>
            )}

            {totalPages > 1 && (
              <div className="hidden sm:flex mt-2 bg-white/90 border border-gray-100 rounded-xl px-3 py-3 items-center justify-between shadow-sm fade-in-up" style={{ animationDelay: "140ms" }}>
                <div className="text-sm text-gray-600">
                  {t("community.pagination", { current: page + 1, total: totalPages })}
                </div>
                <div className="flex gap-2">
                  <Button
                    variant="outline"
                    className="border-gray-300"
                    disabled={page <= 0 || isLoading}
                    onClick={() => setPage((p) => Math.max(0, p - 1))}
                  >
                    {t("community.actions.previous")}
                  </Button>
                  <Button
                    className="bg-[#3A6A8D] hover:bg-[#2d5470] text-white"
                    disabled={page >= totalPages - 1 || isLoading}
                    onClick={() => setPage((p) => p + 1)}
                  >
                    {t("community.actions.next")}
                  </Button>
                </div>
              </div>
            )}
          </div>

          <div className="lg:col-span-3 space-y-6 hidden lg:block fade-in-up" style={{ animationDelay: "100ms" }}>
            <Card className="group relative overflow-hidden p-4 bg-white border border-[#e5e7eb] rounded-2xl shadow-xl hover:shadow-2xl transform transition duration-300 hover:-translate-y-0.5">
              <div className="absolute inset-0 bg-gradient-to-r from-[#3a6a8d]/10 via-transparent to-[#5e9c8d]/10 opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
              <div className="relative z-10">
                <h3 className="font-semibold text-gray-900 mb-4 text-sm sm:text-base">{t("community.popular_tags")}</h3>
                <div className="space-y-2">
                  {popularTags.map((tag, i) => (
                    <div key={tag.name} className="flex items-center justify-between text-sm">
                      <Badge variant="outline" className={tagPillClasses(i)}>
                        {displayTag(tag.name)}
                      </Badge>
                      <span className="text-xs text-gray-500">{tag.count}</span>
                    </div>
                  ))}
                </div>
              </div>
            </Card>
          </div>
        </div>

        <div className="mt-6 text-center sm:text-left fade-in-up" style={{ animationDelay: "160ms" }}>
          {isLoading && <div className="text-sm text-gray-500">{t("community.loading")}</div>}
          {isError && <div className="text-sm text-red-600">{t("community.errors.load_failed")}</div>}
          {!isLoading && !isError && data && (
            <div className="text-xs sm:text-sm text-gray-600">{t("community.total", { count: data.total })}</div>
          )}
        </div>
      </div>

      <div className="sm:hidden fixed bottom-0 left-0 right-0 z-40 bg-white/95 backdrop-blur border-t shadow-lg px-2 py-2 flex items-center justify-between fade-in-up" style={{ animationDelay: "180ms" }}>
        <Button
          variant="outline"
          size="sm"
          className="flex-1 mx-1 text-[#3A6A8D] border-[#3A6A8D]"
          onClick={() => router.push("/user/my-discussions")}
        >
          {t("community.actions.my_discussions_short")}
        </Button>
        <Button
          size="sm"
          className="flex-1 mx-1 bg-[#3A6A8D] hover:bg-[#2d5470] text-white"
          onClick={() => router.push("/user/create-post")}
        >
          <Plus className="h-4 w-4 mr-1" /> {t("community.actions.add_discussion_short")}
        </Button>
        {totalPages > 1 && (
          <div className="flex flex-1 mx-1 gap-1">
            <Button
              variant="outline"
              size="sm"
              className="flex-1"
              disabled={page <= 0}
              onClick={() => setPage((p) => Math.max(0, p - 1))}
            >
              <ChevronLeft className="h-4 w-4" />
            </Button>
            <Button
              variant="outline"
              size="sm"
              className="flex-1"
              disabled={page >= totalPages - 1}
              onClick={() => setPage((p) => p + 1)}
            >
              <ChevronRight className="h-4 w-4" />
            </Button>
          </div>
        )}
      </div>
    </div>
  )
}
