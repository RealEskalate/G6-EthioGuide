// "use client"

// import {
//   Search,
//   ChevronLeft,
//   ChevronRight,
//   Clock,
//   DollarSign,
//   FileText,
//   ShieldCheck,
//   Tag,
//   Sparkles,
//   TrendingUp,
//   Award,
// } from "lucide-react"
// import { Button } from "@/components/ui/button"
// import { Input } from "@/components/ui/input"
// import Image from "next/image"
// import type { Procedure, ProcedureFee } from "@/app/types/procedure"
// import Link from "next/link"
// import { useState, useMemo, useEffect } from "react"
// import { useListProceduresQuery } from "@/app/store/slices/proceduresApi"
// import { Badge } from "@/components/ui/badge"

// export default function UserHomePage() {
//   const [searchQuery, setSearchQuery] = useState("")
//   // dynamic quick access: show only a window (pageLocal) over fetched list subset
//   const { data: procData, isLoading } = useListProceduresQuery({ page: 1, limit: 30 })
//   const [offset, setOffset] = useState(0) // index offset into list
//   const windowSize = 4
//   const memoList = useMemo<Procedure[]>(() => (Array.isArray(procData?.list) ? (procData!.list as Procedure[]) : []), [procData])
//   const filtered = useMemo<Procedure[]>(() => {
//     const q = searchQuery.trim().toLowerCase()
//     if (!q) return memoList
//     return memoList.filter((p: Procedure) => {
//       const hay = ((p.title || p.name || "") + " " + (p.summary || "")).toLowerCase()
//       return hay.includes(q)
//     })
//   }, [memoList, searchQuery])
//   useEffect(() => {
//     setOffset(0)
//   }, [searchQuery])
//   const sliced = useMemo(() => filtered.slice(offset, offset + windowSize), [filtered, offset])
//   const canPrev = offset > 0
//   const canNext = offset + windowSize < filtered.length

//   return (
//     <div className="min-h-screen w-full bg-gray-50 relative overflow-hidden">
//       <div className="absolute inset-0 overflow-hidden pointer-events-none">
//         <div className="absolute -top-24 -right-24 w-56 h-56 rounded-full blur-3xl" style={{ background: 'radial-gradient(closest-side, rgba(167,179,185,0.10), rgba(167,179,185,0))' }}></div>
//       </div>

//       <div className="relative z-10 p-4 sm:p-6 md:p-8 max-w-7xl mx-auto">
//         <div className="mb-6 md:mb-8 w-full">
//           <div className="bg-white/90 backdrop-blur-sm rounded-2xl border border-[#a7b3b9]/30 p-4 sm:p-6 md:p-8 shadow-lg">
//             <div className="flex flex-col sm:flex-row sm:items-center justify-between mb-4 md:mb-6 gap-4">
//               <div>
//                 <h2 className="text-xl sm:text-2xl font-bold text-[#2e4d57] mb-2">Welcome back!</h2>
//                 <p className="text-[#1c3b2e] text-sm sm:text-base">
//                   Navigate Ethiopian government procedures with ease
//                 </p>
//               </div>
//               <div className="inline-flex items-center gap-2 bg-[#3a6a8d]/10 backdrop-blur-sm border border-[#3a6a8d]/30 rounded-full px-3 sm:px-4 py-2">
//                 <Sparkles className="w-4 h-4 text-[#3a6a8d] animate-pulse" />
//                 <span className="text-xs sm:text-sm font-medium text-[#2e4d57]">Government Services</span>
//               </div>
//             </div>

//             <form onSubmit={(e) => e.preventDefault()} className="relative w-full group">
//               <div className="relative">
//                 <Search className="absolute left-3 sm:left-4 top-1/2 transform -translate-y-1/2 w-4 sm:w-5 h-4 sm:h-5 text-[#a7b3b9] group-focus-within:text-[#3a6a8d] transition-colors duration-300" />
//                 <Input
//                   placeholder="Search government procedures..."
//                   value={searchQuery}
//                   onChange={(e) => setSearchQuery(e.target.value)}
//                   className="pl-10 sm:pl-12 pr-20 sm:pr-24 py-3 sm:py-4 bg-white/80 backdrop-blur-sm border-[#a7b3b9]/30 text-[#2e4d57] focus:ring-2 focus:ring-[#3a6a8d]/50 focus:border-[#3a6a8d] transition-all duration-300 rounded-2xl shadow-lg hover:shadow-xl text-base sm:text-lg"
//                 />
//                 <Button
//                   type="button"
//                   className="absolute right-2 top-1/2 transform -translate-y-1/2 bg-[#3a6a8d] hover:bg-[#2e4d57] text-white px-4 sm:px-6 py-2 text-xs sm:text-sm rounded-xl shadow-lg hover:shadow-xl transition-all duration-300 hover:scale-105 border-0"
//                 >
//                   Search
//                 </Button>
//               </div>
//             </form>
//           </div>
//         </div>

//         <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 md:gap-8">
//           <div className="lg:col-span-2">
//             <section className="animate-fade-in-up" style={{ animationDelay: "0.2s" }}>
//               <div className="flex items-center gap-3 mb-4 md:mb-6">
//                 <div className="w-7 sm:w-8 h-7 sm:h-8 bg-[#3a6a8d] rounded-lg flex items-center justify-center">
//                   <TrendingUp className="w-3.5 sm:w-4 h-3.5 sm:h-4 text-white" />
//                 </div>
//                 <h2 className="text-lg sm:text-xl font-bold text-[#2e4d57]">Quick Access Procedures</h2>
//               </div>

//               <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 md:gap-6">
//                 {isLoading &&
//                   Array.from({ length: 4 }).map((_, i) => (
//                     <div
//                       key={i}
//                       className="bg-white/80 backdrop-blur-sm p-4 sm:p-6 rounded-2xl border border-[#a7b3b9]/30 animate-pulse shadow-lg"
//                     >
//                       <div className="flex items-start justify-between mb-4">
//                         <div className="w-12 sm:w-14 h-12 sm:h-14 bg-gradient-to-br from-[#a7b3b9] to-[#5e9c8d] rounded-xl" />
//                         <div className="w-16 sm:w-20 h-5 sm:h-6 bg-[#a7b3b9] rounded-full" />
//                       </div>
//                       <div className="h-4 sm:h-5 bg-[#a7b3b9] rounded-lg w-5/6 mb-3" />
//                       <div className="h-3 sm:h-4 bg-[#a7b3b9] rounded-lg w-4/5 mb-4 sm:mb-5" />
//                       <div className="grid grid-cols-2 gap-2 mb-4 sm:mb-5">
//                         <div className="h-6 sm:h-7 bg-[#a7b3b9] rounded-lg" />
//                         <div className="h-6 sm:h-7 bg-[#a7b3b9] rounded-lg" />
//                       </div>
//                       <div className="h-9 sm:h-11 bg-[#a7b3b9] rounded-xl" />
//                     </div>
//                   ))}

//                 {!isLoading &&
//                     sliced.map((p: Procedure, idx) => {
//                     const title = p.title || p.name || "Procedure"
//                     const summary = p.summary || ""
//                       const pt = p.processingTime
//                     const ptText =
//                       pt && (pt.minDays || pt.maxDays) ? `${pt.minDays ?? "—"}–${pt.maxDays ?? "—"} days` : "—"
//                       const feesArr: ProcedureFee[] = Array.isArray(p.fees) ? p.fees : []
//                       const totalFee = feesArr.reduce((sum: number, f: ProcedureFee) => sum + (Number(f.amount) || 0), 0)
//                     const currency = feesArr[0]?.currency || ""
//                     const feeText = feesArr.length === 0 || totalFee === 0 ? "Free" : `${totalFee} ${currency}`.trim()
//                     const tags: string[] = Array.isArray(p.tags) ? p.tags : []
//                     // const updated = p.updatedAt ? new Date(p.updatedAt).toLocaleDateString() : null
//                     const docsCount = Array.isArray(p.documentsRequired) ? p.documentsRequired.length : null
//                     const verified = !!p.verified

//                     return (
//                       <div
//                         key={p.id || idx}
//                         className="bg-white/90 backdrop-blur-sm p-4 sm:p-6 rounded-2xl border border-[#a7b3b9]/30 hover:shadow-xl transition-all duration-300 hover:-translate-y-1 group cursor-pointer relative overflow-hidden"
//                         style={{ animationDelay: `${idx * 0.1}s` }}
//                       >
//                         <div className="absolute inset-0 bg-gradient-to-br from-[#a7b3b9]/5 to-[#5e9c8d]/5 opacity-0 group-hover:opacity-100 transition-opacity duration-300 rounded-2xl"></div>

//                         <div className="relative z-10">
//                           <div className="flex items-start justify-between mb-4">
//                             <div className="w-12 sm:w-14 h-12 sm:h-14 bg-[#3a6a8d]/10 rounded-xl flex items-center justify-center group-hover:scale-110 transition-transform duration-300">
//                               <Image
//                                 src="/icons/business.svg"
//                                 alt={title}
//                                 width={28}
//                                 height={28}
//                                 className="w-6 sm:w-7 h-6 sm:h-7"
//                               />
//                             </div>
//                             {verified && (
//                               <Badge className="bg-[#5e9c8d]/15 text-[#1c3b2e] border border-[#5e9c8d]/30 flex items-center gap-1.5 px-2 sm:px-3 py-1 rounded-full text-xs">
//                                 <ShieldCheck className="w-3 sm:w-3.5 h-3 sm:h-3.5" />
//                                 Verified
//                               </Badge>
//                             )}
//                           </div>

//                           <h3 className="font-semibold text-[#2e4d57] mb-3 line-clamp-2 group-hover:text-[#3a6a8d] transition-colors duration-300 text-base sm:text-lg leading-snug">
//                             {title}
//                           </h3>

//                           {summary && (
//                             <p className="text-xs sm:text-sm text-[#1c3b2e] mb-4 sm:mb-5 line-clamp-2 leading-relaxed">
//                               {summary}
//                             </p>
//                           )}

//                           <div className="flex flex-wrap items-center gap-2 mb-4 sm:mb-5">
//                             <span className="inline-flex items-center gap-1.5 text-xs font-medium text-[#2e4d57] bg-[#3a6a8d]/10 border border-[#3a6a8d]/30 rounded-lg px-2 sm:px-3 py-1.5 sm:py-2">
//                               <Clock className="w-3 sm:w-3.5 h-3 sm:h-3.5 text-[#3a6a8d]" /> {ptText}
//                             </span>
//                             <span className="inline-flex items-center gap-1.5 text-xs font-medium text-[#2e4d57] bg-[#5e9c8d]/10 border border-[#5e9c8d]/30 rounded-lg px-2 sm:px-3 py-1.5 sm:py-2">
//                               <DollarSign className="w-3 sm:w-3.5 h-3 sm:h-3.5 text-[#5e9c8d]" /> {feeText}
//                             </span>
//                             {typeof docsCount === "number" && (
//                               <span className="inline-flex items-center gap-1.5 text-xs font-medium text-[#2e4d57] bg-[#a7b3b9]/15 border border-[#a7b3b9]/40 rounded-lg px-2 sm:px-3 py-1.5 sm:py-2">
//                                 <FileText className="w-3 sm:w-3.5 h-3 sm:h-3.5 text-[#1c3b2e]" /> {docsCount} docs
//                               </span>
//                             )}
//                           </div>

//                           {tags.length > 0 && (
//             <div className="flex items-center gap-1.5 mb-4 sm:mb-5 flex-wrap">
//                               {tags.slice(0, 2).map((t, i) => (
//                                 <span
//                                   key={i}
//               className="inline-flex items-center gap-1 text-xs font-medium text-[#3a6a8d] bg-[#3a6a8d]/10 px-2 sm:px-2.5 py-1 rounded-full"
//                                 >
//                                   <Tag className="w-2.5 sm:w-3 h-2.5 sm:h-3" /> {t}
//                                 </span>
//                               ))}
//                               {tags.length > 2 && (
//                                 <span className="text-xs text-[#a7b3b9] font-medium">+{tags.length - 2} more</span>
//                               )}
//                             </div>
//                           )}

//                           <Link href={`/user/procedures-detail?id=${p.id}`}>
//                             <Button className="w-full bg-[#3a6a8d] hover:bg-[#2e4d57] text-white transition-all duration-300 hover:scale-105 hover:shadow-lg rounded-xl py-2.5 sm:py-3 font-medium text-sm sm:text-base">
//                               Open Detail
//                             </Button>
//                           </Link>
//                         </div>
//                       </div>
//                     )
//                   })}
//               </div>

//               <div className="flex items-center justify-end gap-3 mt-4 md:mt-6">
//                 <Button
//                   variant="ghost"
//                   size="sm"
//                   disabled={!canPrev}
//                   onClick={() => setOffset((o) => Math.max(0, o - windowSize))}
//                   className="hover:bg-[#3a6a8d]/10 hover:text-[#3a6a8d] transition-all duration-300 rounded-xl px-3 sm:px-4 py-2 text-[#2e4d57]"
//                 >
//                   <ChevronLeft className="w-4 h-4" />
//                 </Button>
//                 <Button
//                   variant="ghost"
//                   size="sm"
//                   disabled={!canNext}
//                   onClick={() => setOffset((o) => Math.min(Math.max(0, filtered.length - windowSize), o + windowSize))}
//                   className="hover:bg-[#3a6a8d]/10 hover:text-[#3a6a8d] transition-all duration-300 rounded-xl px-3 sm:px-4 py-2 text-[#2e4d57]"
//                 >
//                   <ChevronRight className="w-4 h-4" />
//                 </Button>
//               </div>
//             </section>
//           </div>

//           <div className="lg:col-span-1">
//             <section className="w-full animate-fade-in-up" style={{ animationDelay: "0.3s" }}>
//               <div className="flex items-center gap-3 mb-4 md:mb-6">
//                 <div className="w-7 sm:w-8 h-7 sm:h-8 bg-[#5e9c8d] rounded-lg flex items-center justify-center">
//                   <Award className="w-3.5 sm:w-4 h-3.5 sm:h-4 text-white" />
//                 </div>
//                 <h2 className="text-lg sm:text-xl font-bold text-[#2e4d57]">Recent Activity</h2>
//               </div>

//               <div className="bg-white/90 backdrop-blur-sm rounded-2xl border border-[#a7b3b9]/30 p-4 sm:p-6 shadow-lg">
//                 <div className="space-y-3 sm:space-y-4">
//                   {[
//                     {
//                       status: "completed",
//                       title: "Passport Renewal",
//                       statusText: "Completed",
//                       time: "2 hours ago",
//                       color: "green",
//                     },
//                     {
//                       status: "progress",
//                       title: "Business License Application",
//                       statusText: "In Progress",
//                       time: "1 day ago",
//                       color: "blue",
//                     },
//                     {
//                       status: "completed",
//                       title: "Vehicle Registration",
//                       statusText: "Completed",
//                       time: "1 week ago",
//                       color: "green",
//                     },
//                   ].map((item, index) => (
//                     <div
//                       key={index}
//                       className="flex items-center gap-3 py-2 sm:py-3 px-2 sm:px-3 border border-[#a7b3b9]/20 last:border-b-0 hover:bg-[#a7b3b9]/10 rounded-xl transition-all duration-300 group cursor-pointer"
//                     >
//                       <div
//                         className={`w-8 sm:w-10 h-8 sm:h-10 rounded-lg flex items-center justify-center transition-all duration-300 group-hover:scale-110 ${
//                           item.color === "green" ? "bg-[#5e9c8d]/20" : "bg-[#3a6a8d]/20"
//                         }`}
//                       >
//                         <div
//                           className={`w-4 sm:w-5 h-4 sm:h-5 flex items-center justify-center text-xs font-bold rounded-md ${
//                             item.color === "green" ? "text-[#1c3b2e] bg-[#5e9c8d]/20" : "text-[#3a6a8d] bg-[#3a6a8d]/20"
//                           }`}
//                         >
//                           {item.status === "completed" ? "✓" : "⏳"}
//                         </div>
//                       </div>
//                       <div className="flex-1 min-w-0">
//                         <p className="font-medium text-[#2e4d57] text-xs sm:text-sm group-hover:text-[#3a6a8d] transition-colors duration-300 truncate">
//                           {item.title}
//                         </p>
//                         <p
//                           className={`text-xs font-medium ${
//                             item.color === "green" ? "text-[#1c3b2e]" : "text-[#3a6a8d]"
//                           }`}
//                         >
//                           {item.statusText}
//                         </p>
//                         <p className="text-xs text-[#a7b3b9] font-medium">{item.time}</p>
//                       </div>
//                     </div>
//                   ))}
//                 </div>

//                 <div className="mt-4 sm:mt-6 pt-3 sm:pt-4 border-t border-[#a7b3b9]/20 flex justify-center">
//                   <Link href="/user/workspace">
//                     <Button
//                       variant="link"
//                       className="text-[#3a6a8d] hover:text-[#2e4d57] transition-colors duration-300 font-medium text-xs sm:text-sm hover:bg-[#3a6a8d]/10 px-3 sm:px-4 py-2 rounded-xl"
//                     >
//                       View All Activities →
//                     </Button>
//                   </Link>
//                 </div>
//               </div>
//             </section>
//           </div>
//         </div>
//       </div>
//     </div>
//   )
// }


"use client"

import {
  Search,
  ChevronLeft,
  ChevronRight,
  Clock,
  DollarSign,
  FileText,
  ShieldCheck,
  Tag,
  Sparkles,
  TrendingUp,
  Award,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import Image from "next/image"
import type { Procedure, ProcedureFee } from "@/app/types/procedure"
import Link from "next/link"
import { useState, useMemo, useEffect } from "react"
import { useListProceduresQuery } from "@/app/store/slices/proceduresApi"
import { Badge } from "@/components/ui/badge"
import { useTranslation } from "react-i18next"

export default function UserHomePage() {
  const { t } = useTranslation("user")
  const [searchQuery, setSearchQuery] = useState("")
  const { data: procData, isLoading } = useListProceduresQuery({ page: 1, limit: 30 })
  const [offset, setOffset] = useState(0)
  const windowSize = 4
  const memoList = useMemo<Procedure[]>(() => (Array.isArray(procData?.list) ? (procData!.list as Procedure[]) : []), [procData])
  const filtered = useMemo<Procedure[]>(() => {
    const q = searchQuery.trim().toLowerCase()
    if (!q) return memoList
    return memoList.filter((p: Procedure) => {
      const hay = ((p.title || p.name || "") + " " + (p.summary || "")).toLowerCase()
      return hay.includes(q)
    })
  }, [memoList, searchQuery])
  useEffect(() => {
    setOffset(0)
  }, [searchQuery])
  const sliced = useMemo(() => filtered.slice(offset, offset + windowSize), [filtered, offset])
  const canPrev = offset > 0
  const canNext = offset + windowSize < filtered.length

  return (
    <div className="min-h-screen w-full bg-gray-50 relative overflow-hidden">
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        <div className="absolute -top-24 -right-24 w-56 h-56 rounded-full blur-3xl" style={{ background: 'radial-gradient(closest-side, rgba(167,179,185,0.10), rgba(167,179,185,0))' }}></div>
      </div>

      <div className="relative z-10 p-4 sm:p-6 md:p-8 max-w-7xl mx-auto">
        <div className="mb-6 md:mb-8 w-full">
          <div className="bg-white/90 backdrop-blur-sm rounded-2xl border border-[#a7b3b9]/30 p-4 sm:p-6 md:p-8 shadow-lg">
            <div className="flex flex-col sm:flex-row sm:items-center justify-between mb-4 md:mb-6 gap-4">
              <div>
                <h2 className="text-xl sm:text-2xl font-bold text-[#2e4d57] mb-2">{t("welcome.title")}</h2>
                <p className="text-[#1c3b2e] text-sm sm:text-base">{t("welcome.description")}</p>
              </div>
              <div className="inline-flex items-center gap-2 bg-[#3a6a8d]/10 backdrop-blur-sm border border-[#3a6a8d]/30 rounded-full px-3 sm:px-4 py-2">
                <Sparkles className="w-4 h-4 text-[#3a6a8d] animate-pulse" />
                <span className="text-xs sm:text-sm font-medium text-[#2e4d57]">{t("welcome.services")}</span>
              </div>
            </div>

            <form onSubmit={(e) => e.preventDefault()} className="relative w-full group">
              <div className="relative">
                <Search className="absolute left-3 sm:left-4 top-1/2 transform -translate-y-1/2 w-4 sm:w-5 h-4 sm:h-5 text-[#a7b3b9] group-focus-within:text-[#3a6a8d] transition-colors duration-300" />
                <Input
                  placeholder={t("search.placeholder")}
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-10 sm:pl-12 pr-20 sm:pr-24 py-3 sm:py-4 bg-white/80 backdrop-blur-sm border-[#a7b3b9]/30 text-[#2e4d57] focus:ring-2 focus:ring-[#3a6a8d]/50 focus:border-[#3a6a8d] transition-all duration-300 rounded-2xl shadow-lg hover:shadow-xl text-base sm:text-lg"
                />
                <Button
                  type="button"
                  className="absolute right-2 top-1/2 transform -translate-y-1/2 bg-[#3a6a8d] hover:bg-[#2e4d57] text-white px-4 sm:px-6 py-2 text-xs sm:text-sm rounded-xl shadow-lg hover:shadow-xl transition-all duration-300 hover:scale-105 border-0"
                >
                  {t("search.button")}
                </Button>
              </div>
            </form>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 md:gap-8">
          <div className="lg:col-span-2">
            <section className="animate-fade-in-up" style={{ animationDelay: "0.2s" }}>
              <div className="flex items-center gap-3 mb-4 md:mb-6">
                <div className="w-7 sm:w-8 h-7 sm:h-8 bg-[#3a6a8d] rounded-lg flex items-center justify-center">
                  <TrendingUp className="w-3.5 sm:w-4 h-3.5 sm:h-4 text-white" />
                </div>
                <h2 className="text-lg sm:text-xl font-bold text-[#2e4d57]">{t("procedures.title")}</h2>
              </div>

              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 md:gap-6">
                {isLoading &&
                  Array.from({ length: 4 }).map((_, i) => (
                    <div
                      key={i}
                      className="bg-white/80 backdrop-blur-sm p-4 sm:p-6 rounded-2xl border border-[#a7b3b9]/30 animate-pulse shadow-lg"
                    >
                      <div className="flex items-start justify-between mb-4">
                        <div className="w-12 sm:w-14 h-12 sm:h-14 bg-gradient-to-br from-[#a7b3b9] to-[#5e9c8d] rounded-xl" />
                        <div className="w-16 sm:w-20 h-5 sm:h-6 bg-[#a7b3b9] rounded-full" />
                      </div>
                      <div className="h-4 sm:h-5 bg-[#a7b3b9] rounded-lg w-5/6 mb-3" />
                      <div className="h-3 sm:h-4 bg-[#a7b3b9] rounded-lg w-4/5 mb-4 sm:mb-5" />
                      <div className="grid grid-cols-2 gap-2 mb-4 sm:mb-5">
                        <div className="h-6 sm:h-7 bg-[#a7b3b9] rounded-lg" />
                        <div className="h-6 sm:h-7 bg-[#a7b3b9] rounded-lg" />
                      </div>
                      <div className="h-9 sm:h-11 bg-[#a7b3b9] rounded-xl" />
                    </div>
                  ))}

                {!isLoading &&
                  sliced.map((p: Procedure, idx) => {
                    // Debug problematic procedure data
                    if (!p) {
                      console.error("Procedure is undefined at index:", idx);
                      return null;
                    }

                    const title = p.title || p.name || t("procedures.default_title")
                    const summary = p.summary || ""
                    const pt = p.processingTime
                    const ptText =
                      pt && (pt.minDays || pt.maxDays) ? `${pt.minDays ?? "—"}–${pt.maxDays ?? "—"} ${t("procedures.days")}` : "—"
                    const feesArr: ProcedureFee[] = Array.isArray(p.fees) ? p.fees : [];
                    const totalFee = feesArr.length > 0 ? feesArr.reduce((sum: number, f: ProcedureFee) => sum + (Number(f.amount) || 0), 0) : 0;
                    const currency = feesArr[0]?.currency || ""
                    const feeText = feesArr.length === 0 || totalFee === 0 ? t("procedures.free") : `${totalFee} ${currency}`.trim()
                    const tags: string[] = Array.isArray(p.tags) ? p.tags : []
                    const docsCount = Array.isArray(p.documentsRequired) ? p.documentsRequired.length : null
                    const verified = !!p.verified

                    return (
                      <div
                        key={p.id || idx}
                        className="bg-white/90 backdrop-blur-sm p-4 sm:p-6 rounded-2xl border border-[#a7b3b9]/30 hover:shadow-xl transition-all duration-300 hover:-translate-y-1 group cursor-pointer relative overflow-hidden"
                        style={{ animationDelay: `${idx * 0.1}s` }}
                      >
                        <div className="absolute inset-0 bg-gradient-to-br from-[#a7b3b9]/5 to-[#5e9c8d]/5 opacity-0 group-hover:opacity-100 transition-opacity duration-300 rounded-2xl"></div>

                        <div className="relative z-10">
                          <div className="flex items-start justify-between mb-4">
                            <div className="w-12 sm:w-14 h-12 sm:h-14 bg-[#3a6a8d]/10 rounded-xl flex items-center justify-center group-hover:scale-110 transition-transform duration-300">
                              <Image
                                src="/icons/business.svg"
                                alt={title}
                                width={28}
                                height={28}
                                className="w-6 sm:w-7 h-6 sm:h-7"
                              />
                            </div>
                            {verified && (
                              <Badge className="bg-[#5e9c8d]/15 text-[#1c3b2e] border border-[#5e9c8d]/30 flex items-center gap-1.5 px-2 sm:px-3 py-1 rounded-full text-xs">
                                <ShieldCheck className="w-3 sm:w-3.5 h-3 sm:h-3.5" />
                                {t("procedures.verified")}
                              </Badge>
                            )}
                          </div>

                          <h3 className="font-semibold text-[#2e4d57] mb-3 line-clamp-2 group-hover:text-[#3a6a8d] transition-colors duration-300 text-base sm:text-lg leading-snug">
                            {title}
                          </h3>

                          {summary && (
                            <p className="text-xs sm:text-sm text-[#1c3b2e] mb-4 sm:mb-5 line-clamp-2 leading-relaxed">
                              {summary}
                            </p>
                          )}

                          <div className="flex flex-wrap items-center gap-2 mb-4 sm:mb-5">
                            <span className="inline-flex items-center gap-1.5 text-xs font-medium text-[#2e4d57] bg-[#3a6a8d]/10 border border-[#3a6a8d]/30 rounded-lg px-2 sm:px-3 py-1.5 sm:py-2">
                              <Clock className="w-3 sm:w-3.5 h-3 sm:h-3.5 text-[#3a6a8d]" /> {ptText}
                            </span>
                            <span className="inline-flex items-center gap-1.5 text-xs font-medium text-[#2e4d57] bg-[#5e9c8d]/10 border border-[#5e9c8d]/30 rounded-lg px-2 sm:px-3 py-1.5 sm:py-2">
                              <DollarSign className="w-3 sm:w-3.5 h-3 sm:h-3.5 text-[#5e9c8d]" /> {feeText}
                            </span>
                            {typeof docsCount === "number" && (
                              <span className="inline-flex items-center gap-1.5 text-xs font-medium text-[#2e4d57] bg-[#a7b3b9]/15 border border-[#a7b3b9]/40 rounded-lg px-2 sm:px-3 py-1.5 sm:py-2">
                                <FileText className="w-3 sm:w-3.5 h-3 sm:h-3.5 text-[#1c3b2e]" /> {docsCount} {t("procedures.docs")}
                              </span>
                            )}
                          </div>

                          {tags.length > 0 && (
                            <div className="flex items-center gap-1.5 mb-4 sm:mb-5 flex-wrap">
                              {tags.slice(0, 2).map((tag, i) => (
                                <span
                                  key={i}
                                  className="inline-flex items-center gap-1 text-xs font-medium text-[#3a6a8d] bg-[#3a6a8d]/10 px-2 sm:px-2.5 py-1 rounded-full"
                                >
                                  <Tag className="w-2.5 sm:w-3 h-2.5 sm:h-3" /> {tag}
                                </span>
                              ))}
                              {tags.length > 2 && (
                                <span className="text-xs text-[#a7b3b9] font-medium">+{tags.length - 2} {t("procedures.more")}</span>
                              )}
                            </div>
                          )}

                          <Link href={`/user/procedures-detail?id=${p.id}`}>
                            <Button className="w-full bg-[#3a6a8d] hover:bg-[#2e4d57] text-white transition-all duration-300 hover:scale-105 hover:shadow-lg rounded-xl py-2.5 sm:py-3 font-medium text-sm sm:text-base">
                              {t("procedures.button")}
                            </Button>
                          </Link>
                        </div>
                      </div>
                    )
                  })}
              </div>

              <div className="flex items-center justify-end gap-3 mt-4 md:mt-6">
                <Button
                  variant="ghost"
                  size="sm"
                  disabled={!canPrev}
                  onClick={() => setOffset((o) => Math.max(0, o - windowSize))}
                  className="hover:bg-[#3a6a8d]/10 hover:text-[#3a6a8d] transition-all duration-300 rounded-xl px-3 sm:px-4 py-2 text-[#2e4d57]"
                >
                  <ChevronLeft className="w-4 h-4" />
                </Button>
                <Button
                  variant="ghost"
                  size="sm"
                  disabled={!canNext}
                  onClick={() => setOffset((o) => Math.min(Math.max(0, filtered.length - windowSize), o + windowSize))}
                  className="hover:bg-[#3a6a8d]/10 hover:text-[#3a6a8d] transition-all duration-300 rounded-xl px-3 sm:px-4 py-2 text-[#2e4d57]"
                >
                  <ChevronRight className="w-4 h-4" />
                </Button>
              </div>
            </section>
          </div>

          <div className="lg:col-span-1">
            <section className="w-full animate-fade-in-up" style={{ animationDelay: "0.3s" }}>
              <div className="flex items-center gap-3 mb-4 md:mb-6">
                <div className="w-7 sm:w-8 h-7 sm:h-8 bg-[#5e9c8d] rounded-lg flex items-center justify-center">
                  <Award className="w-3.5 sm:w-4 h-3.5 sm:h-4 text-white" />
                </div>
                <h2 className="text-lg sm:text-xl font-bold text-[#2e4d57]">{t("activity.title")}</h2>
              </div>

              <div className="bg-white/90 backdrop-blur-sm rounded-2xl border border-[#a7b3b9]/30 p-4 sm:p-6 shadow-lg">
                <div className="space-y-3 sm:space-y-4">
                  {[
                    {
                      status: "completed",
                      title: t("activity.items.passport_renewal"),
                      statusText: t("activity.status.completed"),
                      time: t("activity.times.two_hours"),
                      color: "green",
                    },
                    {
                      status: "progress",
                      title: t("activity.items.business_license"),
                      statusText: t("activity.status.in_progress"),
                      time: t("activity.times.one_day"),
                      color: "blue",
                    },
                    {
                      status: "completed",
                      title: t("activity.items.vehicle_registration"),
                      statusText: t("activity.status.completed"),
                      time: t("activity.times.one_week"),
                      color: "green",
                    },
                  ].map((item, index) => (
                    <div
                      key={index}
                      className="flex items-center gap-3 py-2 sm:py-3 px-2 sm:px-3 border border-[#a7b3b9]/20 last:border-b-0 hover:bg-[#a7b3b9]/10 rounded-xl transition-all duration-300 group cursor-pointer"
                    >
                      <div
                        className={`w-8 sm:w-10 h-8 sm:h-10 rounded-lg flex items-center justify-center transition-all duration-300 group-hover:scale-110 ${
                          item.color === "green" ? "bg-[#5e9c8d]/20" : "bg-[#3a6a8d]/20"
                        }`}
                      >
                        <div
                          className={`w-4 sm:w-5 h-4 sm:h-5 flex items-center justify-center text-xs font-bold rounded-md ${
                            item.color === "green" ? "text-[#1c3b2e] bg-[#5e9c8d]/20" : "text-[#3a6a8d] bg-[#3a6a8d]/20"
                          }`}
                        >
                          {item.status === "completed" ? "✓" : "⏳"}
                        </div>
                      </div>
                      <div className="flex-1 min-w-0">
                        <p className="font-medium text-[#2e4d57] text-xs sm:text-sm group-hover:text-[#3a6a8d] transition-colors duration-300 truncate">
                          {item.title}
                        </p>
                        <p
                          className={`text-xs font-medium ${
                            item.color === "green" ? "text-[#1c3b2e]" : "text-[#3a6a8d]"
                          }`}
                        >
                          {item.statusText}
                        </p>
                        <p className="text-xs text-[#a7b3b9] font-medium">{item.time}</p>
                      </div>
                    </div>
                  ))}
                </div>

                <div className="mt-4 sm:mt-6 pt-3 sm:pt-4 border-t border-[#a7b3b9]/20 flex justify-center">
                  <Link href="/user/workspace">
                    <Button
                      variant="link"
                      className="text-[#3a6a8d] hover:text-[#2e4d57] transition-colors duration-300 font-medium text-xs sm:text-sm hover:bg-[#3a6a8d]/10 px-3 sm:px-4 py-2 rounded-xl"
                    >
                      {t("activity.view_all")}
                    </Button>
                  </Link>
                </div>
              </div>
            </section>
          </div>
        </div>
      </div>
    </div>
  )
}
