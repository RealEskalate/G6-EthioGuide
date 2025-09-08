// "use client"

// export const dynamic = 'force-dynamic'

// import type React from "react"
// import { Suspense, useEffect, useMemo, useState } from "react"
// import { useSearchParams } from "next/navigation"
// import { useSession } from "next-auth/react"
// import { ArrowLeft, Send } from "lucide-react"
// import { Button } from "@/components/ui/button"
// import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
// import { Input } from "@/components/ui/input"
// import { Label } from "@/components/ui/label"
// import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
// import { Textarea } from "@/components/ui/textarea"
// import { Badge } from "@/components/ui/badge"
// import { useGetProcedureFeedbackQuery, useSubmitProcedureFeedbackMutation } from "@/app/store/slices/feedbackApi"
// import { toast } from "sonner"

// function FeedbackPageInner() {
//   const search = useSearchParams()
//   const procedureId = search.get("id") || ""

//   const [feedbackType, setFeedbackType] = useState("inaccuracy")
//   const [subject, setSubject] = useState("")
//   const [detailedFeedback, setDetailedFeedback] = useState("")
//   const [page, setPage] = useState(1)
//   const limit = 5
//   const [token, setToken] = useState<string | null>(null)
//   const { data: session } = useSession()
//   const [tagsInput, setTagsInput] = useState("")

//   useEffect(() => {
//     try {
//       const t = typeof window !== "undefined" ? window.localStorage.getItem("token") : null
//       if (t) setToken(t)
//     } catch {
//       // ignore
//     }
//   }, [])

//   const {
//     data: history,
//     isFetching,
//     isError,
//     refetch,
//   } = useGetProcedureFeedbackQuery(
//     { procedureId, page, limit, token: (session as { accessToken?: string } | null | undefined)?.accessToken || token },
//     { skip: !procedureId },
//   )
//   const [submitFeedback, { isLoading: isSubmitting }] = useSubmitProcedureFeedbackMutation()

//   const totalPages = useMemo(() => {
//     if (!history) return 1
//     return Math.max(1, Math.ceil((history.total || 0) / (history.limit || limit)))
//   }, [history])

//   const handleSubmit = async (e: React.FormEvent) => {
//     e.preventDefault()
//     if (!procedureId) {
//       toast.error("Open feedback from a procedure to submit.")
//       return
//     }
//   const authToken = (session as { accessToken?: string } | null | undefined)?.accessToken || token
//     if (!authToken) {
//       toast.error("You must be logged in to submit feedback.")
//       return
//     }
//     const content = subject ? `${subject}\n\n${detailedFeedback}`.trim() : detailedFeedback.trim()
//     if (!content) {
//       toast.error("Please enter your feedback.")
//       return
//     }
//     try {
//       const tags = tagsInput.split(',').map(t => t.trim()).filter(Boolean)
//       await submitFeedback({ procedureId, content, type: feedbackType, tags, token: authToken }).unwrap()
//       toast.success("Feedback submitted.")
//       setSubject("")
//       setDetailedFeedback("")
//       setTagsInput("")
//       setPage(1)
//       refetch()
//     } catch (err: unknown) {
//       const e = err as { data?: { error?: unknown; message?: unknown }; message?: unknown }
//       const msg = (e?.data?.error ?? e?.data?.message ?? e?.message ?? 'Failed to submit feedback')
//       toast.error(String(msg))
//     }
//   }

//   return (
//     <div className="min-h-screen bg-gray-50 relative overflow-hidden">
//       <div className="absolute inset-0 overflow-hidden pointer-events-none">
//         <div className="absolute -top-40 -right-40 w-80 h-80 bg-gradient-to-br from-[#5e9c8d]/20 to-[#3a6a8d]/20 rounded-full blur-3xl animate-pulse"></div>
//         <div className="absolute -bottom-40 -left-40 w-80 h-80 bg-gradient-to-tr from-[#2e4d57]/20 to-[#1c3b2e]/20 rounded-full blur-3xl animate-pulse delay-1000"></div>
//       </div>

//       <div className="relative bg-white/80 backdrop-blur-md border-b border-[#a7b3b9]/30 px-4 py-4 shadow-sm">
//         <div className="flex items-center justify-between max-w-7xl mx-auto">
//           <div className="flex items-center gap-3">
//             <Button
//               variant="ghost"
//               size="sm"
//               className="p-2 hover:bg-[#5e9c8d]/20 transition-all duration-300 rounded-xl"
//               onClick={() => {
//                 window.location.href = procedureId
//                   ? `/user/procedures-detail?id=${encodeURIComponent(procedureId)}`
//                   : "/user/procedures-detail"
//               }}
//             >
//               <ArrowLeft className="h-5 w-5 text-[#2e4d57]" />
//             </Button>
//             <h1 className="text-xl font-semibold text-[#1c3b2e] tracking-tight">Feedback</h1>
//           </div>
//         </div>
//       </div>

//       <div className="relative max-w-7xl mx-auto p-4 sm:p-6 lg:p-8">
//         <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 lg:gap-8">
//           <Card className="bg-white/80 backdrop-blur-md border-[#a7b3b9]/30 shadow-xl hover:shadow-2xl transition-all duration-500 rounded-2xl">
//             <CardHeader className="pb-4">
//               <CardTitle className="text-lg font-semibold text-[#1c3b2e] tracking-tight">Submit Feedback</CardTitle>
//             </CardHeader>
//             <CardContent className="space-y-6">
//               <form onSubmit={handleSubmit} className="space-y-6">
//                 <div className="space-y-2">
//                   <Label htmlFor="feedback-type" className="text-sm font-medium text-[#2e4d57]">
//                     Type of Feedback
//                   </Label>
//                   <Select value={feedbackType} onValueChange={setFeedbackType}>
//                     <SelectTrigger className="bg-white/70 border-[#a7b3b9]/50 hover:border-[#3a6a8d] transition-all duration-300 rounded-xl h-12">
//                       <SelectValue placeholder="Report Issue" />
//                     </SelectTrigger>
//                     <SelectContent className="bg-white/95 backdrop-blur-md border-[#a7b3b9]/30 rounded-xl">
//                       <SelectItem
//                         value="inaccuracy"
//                         className="hover:bg-[#5e9c8d]/20 rounded-lg transition-colors duration-200"
//                       >
//                         Report Inaccuracy
//                       </SelectItem>
//                       <SelectItem
//                         value="improvement"
//                         className="hover:bg-[#5e9c8d]/20 rounded-lg transition-colors duration-200"
//                       >
//                         Improvement
//                       </SelectItem>
//                       <SelectItem
//                         value="other"
//                         className="hover:bg-[#5e9c8d]/20 rounded-lg transition-colors duration-200"
//                       >
//                         Other
//                       </SelectItem>
//                     </SelectContent>
//                   </Select>
//                 </div>

//                 <div className="space-y-2">
//                   <Label htmlFor="subject" className="text-sm font-medium text-[#2e4d57]">
//                     Subject/Procedure/Document
//                   </Label>
//                   <Input
//                     id="subject"
//                     placeholder="Describe your report procedure..."
//                     value={subject}
//                     onChange={(e) => setSubject(e.target.value)}
//                     className="bg-white/70 border-[#a7b3b9]/50 hover:border-[#3a6a8d] focus:border-[#3a6a8d] transition-all duration-300 rounded-xl h-12"
//                     disabled={!procedureId || isSubmitting}
//                   />
//                 </div>

//                 <div className="space-y-2">
//                   <Label htmlFor="detailed-feedback" className="text-sm font-medium text-[#2e4d57]">
//                     Detailed Feedback
//                   </Label>
//                   <Textarea
//                     id="detailed-feedback"
//                     placeholder="Please provide detailed feedback..."
//                     value={detailedFeedback}
//                     onChange={(e) => setDetailedFeedback(e.target.value)}
//                     className="min-h-[120px] bg-white/70 border-[#a7b3b9]/50 hover:border-[#3a6a8d] focus:border-[#3a6a8d] transition-all duration-300 rounded-xl resize-none"
//                     disabled={!procedureId || isSubmitting}
//                   />
//                 </div>

//                 <div className="space-y-2">
//                   <Label htmlFor="tags" className="text-sm font-medium text-[#2e4d57]">
//                     Tags (comma separated)
//                   </Label>
//                   <Input
//                     id="tags"
//                     placeholder="e.g. document, payment, delay"
//                     value={tagsInput}
//                     onChange={(e) => setTagsInput(e.target.value)}
//                     className="bg-white/70 border-[#a7b3b9]/50 hover:border-[#3a6a8d] focus:border-[#3a6a8d] transition-all duration-300 rounded-xl h-12"
//                     disabled={!procedureId || isSubmitting}
//                   />
//                 </div>

//                 <Button
//                   type="submit"
//                   disabled={!procedureId || isSubmitting}
//                   className="w-full bg-gradient-to-r from-[#3a6a8d] to-[#2e4d57] hover:from-[#2e4d57] hover:to-[#1c3b2e] text-white py-3 rounded-xl font-medium transition-all duration-300 transform hover:scale-[1.02] disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none"
//                 >
//                   <Send className="h-4 w-4 mr-2" />
//                   {isSubmitting ? "Submitting..." : "Send Feedback"}
//                 </Button>
//                 {!procedureId && (
//                   <p className="text-xs text-[#2e4d57]/70 text-center">
//                     Open the feedback page from a procedure to submit feedback.
//                   </p>
//                 )}
//               </form>
//             </CardContent>
//           </Card>

//           <Card className="bg-white/80 backdrop-blur-md border-[#a7b3b9]/30 shadow-xl hover:shadow-2xl transition-all duration-500 rounded-2xl">
//             <CardHeader className="pb-4">
//               <CardTitle className="text-lg font-semibold text-[#1c3b2e] tracking-tight">
//                 Your Feedback History
//               </CardTitle>
//             </CardHeader>
//             <CardContent className="space-y-4">
//               {!procedureId && (
//                 <div className="text-sm text-[#2e4d57]/70 text-center py-8">
//                   No procedure selected. Open feedback from a specific procedure to view its feedback.
//                 </div>
//               )}
//               {procedureId && isFetching && (
//                 <div className="text-sm text-[#2e4d57]/70 text-center py-8">Loading...</div>
//               )}
//               {procedureId && !isFetching && isError && (
//                 <div className="text-sm text-red-600 text-center py-8">Failed to load feedback history.</div>
//               )}
//               {procedureId && !isFetching && history && history.feedbacks.length === 0 && (
//                 <div className="text-sm text-[#2e4d57]/70 text-center py-8">
//                   No feedback yet. Be the first to submit.
//                 </div>
//               )}
//               {procedureId && history && history.feedbacks.length > 0 && (
//                 <>
//                   <div className="space-y-3">
//                     {history.feedbacks.map((f) => (
//                       <div
//                         key={f.id}
//                         className="bg-white/60 backdrop-blur-sm border border-[#a7b3b9]/30 rounded-xl p-4 space-y-3 hover:bg-white/80 transition-all duration-300"
//                       >
//                         <div className="flex items-center justify-between">
//                           <span className="text-sm font-medium text-[#2e4d57] capitalize">
//                             {f.type?.replace(/_/g, " ")}
//                           </span>
//                           <Badge
//                             variant="outline"
//                             className="text-xs px-3 py-1 rounded-full border-[#5e9c8d] text-[#1c3b2e] bg-[#5e9c8d]/10"
//                           >
//                             {f.status || "new"}
//                           </Badge>
//                         </div>
//                         <p className="text-sm text-[#2e4d57]/80 whitespace-pre-line leading-relaxed">
//                           {(f.content || "").slice(0, 220)}
//                           {(f.content || "").length > 220 ? "…" : ""}
//                         </p>
//                         <p className="text-xs text-[#a7b3b9]">
//                           {f.createdAT ? new Date(f.createdAT).toLocaleString() : ""}
//                         </p>
//                       </div>
//                     ))}
//                   </div>
//                   <div className="flex items-center justify-between pt-4 border-t border-[#a7b3b9]/20">
//                     <Button
//                       variant="outline"
//                       size="sm"
//                       disabled={page <= 1}
//                       onClick={() => setPage((p) => Math.max(1, p - 1))}
//                       className="border-[#a7b3b9]/50 text-[#2e4d57] hover:bg-[#5e9c8d]/20 rounded-lg transition-all duration-200"
//                     >
//                       Previous
//                     </Button>
//                     <span className="text-xs text-[#2e4d57]/70 font-medium">
//                       Page {history.page} of {totalPages}
//                     </span>
//                     <Button
//                       variant="outline"
//                       size="sm"
//                       disabled={history.page >= totalPages}
//                       onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
//                       className="border-[#a7b3b9]/50 text-[#2e4d57] hover:bg-[#5e9c8d]/20 rounded-lg transition-all duration-200"
//                     >
//                       Next
//                     </Button>
//                   </div>
//                 </>
//               )}
//             </CardContent>
//           </Card>
//         </div>
//       </div>
//     </div>
//   )
// }

// export default function Page() {
//   return (
//     <Suspense fallback={<div className="p-4 text-center text-gray-500">Loading...</div>}>
//       <FeedbackPageInner />
//     </Suspense>
//   )
// }


"use client"

export const dynamic = 'force-dynamic'

import type React from "react"
import { Suspense, useEffect, useMemo, useState } from "react"
import { useSearchParams } from "next/navigation"
import { useSession } from "next-auth/react"
import { ArrowLeft, Send } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Textarea } from "@/components/ui/textarea"
import { Badge } from "@/components/ui/badge"
import { useGetProcedureFeedbackQuery, useSubmitProcedureFeedbackMutation } from "@/app/store/slices/feedbackApi"
import { toast } from "sonner"
import { useTranslation } from "react-i18next"

// Fallback component with translation support
function FeedbackFallback() {
  const { t } = useTranslation("user")
  return <div className="p-4 text-center text-gray-500">{t("loading")}</div>
}

function FeedbackPageInner() {
  const { t } = useTranslation("user")
  const search = useSearchParams()
  const procedureId = search ? search.get("id") || "" : ""
  const [feedbackType, setFeedbackType] = useState("inaccuracy")
  const [subject, setSubject] = useState("")
  const [detailedFeedback, setDetailedFeedback] = useState("")
  const [page, setPage] = useState(1)
  const limit = 5
  const [token, setToken] = useState<string | null>(null)
  const { data: session } = useSession()
  const [tagsInput, setTagsInput] = useState("")

  useEffect(() => {
    try {
      const t = typeof window !== "undefined" ? window.sessionStorage.getItem("token") : null
      if (t) setToken(t)
    } catch {
      // ignore
    }
  }, [])

  const {
    data: history,
    isFetching,
    isError,
    refetch,
  } = useGetProcedureFeedbackQuery(
    { procedureId, page, limit, token: (session as { accessToken?: string } | null | undefined)?.accessToken || token },
    { skip: !procedureId || !token },
  )
  const [submitFeedback, { isLoading: isSubmitting }] = useSubmitProcedureFeedbackMutation()

  const totalPages = useMemo(() => {
    if (!history) return 1
    return Math.max(1, Math.ceil((history.total || 0) / (history.limit || limit)))
  }, [history])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!procedureId) {
      toast.error(t("feedback.errors.no_procedure"))
      return
    }
    const authToken = (session as { accessToken?: string } | null | undefined)?.accessToken || token
    if (!authToken) {
      toast.error(t("feedback.errors.not_authenticated"))
      return
    }
    const content = subject ? `${subject}\n\n${detailedFeedback}`.trim() : detailedFeedback.trim()
    if (!content) {
      toast.error(t("feedback.errors.no_content"))
      return
    }
    try {
      const tags = tagsInput.split(',').map(t => t.trim()).filter(Boolean)
      await submitFeedback({ procedureId, content, type: feedbackType, tags, token: authToken }).unwrap()
      toast.success(t("feedback.success"))
      setSubject("")
      setDetailedFeedback("")
      setTagsInput("")
      setPage(1)
      refetch()
    } catch (err: unknown) {
      const e = err as { data?: { error?: unknown; message?: unknown }; message?: unknown }
      const msg = (e?.data?.error ?? e?.data?.message ?? e?.message ?? t("feedback.errors.submit_failed"))
      toast.error(String(msg))
    }
  }

  return (
    <div className="min-h-screen bg-gray-50 relative overflow-hidden">
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        <div className="absolute -top-40 -right-40 w-80 h-80 bg-gradient-to-br from-[#5e9c8d]/20 to-[#3a6a8d]/20 rounded-full blur-3xl animate-pulse"></div>
        <div className="absolute -bottom-40 -left-40 w-80 h-80 bg-gradient-to-tr from-[#2e4d57]/20 to-[#1c3b2e]/20 rounded-full blur-3xl animate-pulse delay-1000"></div>
      </div>

      <div className="relative bg-white/80 backdrop-blur-md border-b border-[#a7b3b9]/30 px-4 py-4 shadow-sm">
        <div className="flex items-center justify-between max-w-7xl mx-auto">
          <div className="flex items-center gap-3">
            <Button
              variant="ghost"
              size="sm"
              className="p-2 hover:bg-[#5e9c8d]/20 transition-all duration-300 rounded-xl"
              onClick={() => {
                window.location.href = procedureId
                  ? `/user/procedures-detail?id=${encodeURIComponent(procedureId)}`
                  : "/user/procedures-detail"
              }}
            >
              <ArrowLeft className="h-5 w-5 text-[#2e4d57]" />
            </Button>
            <h1 className="text-xl font-semibold text-[#1c3b2e] tracking-tight">{t("feedback.title")}</h1>
          </div>
        </div>
      </div>

      <div className="relative max-w-7xl mx-auto p-4 sm:p-6 lg:p-8">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 lg:gap-8">
          <Card className="bg-white/80 backdrop-blur-md border-[#a7b3b9]/30 shadow-xl hover:shadow-2xl transition-all duration-500 rounded-2xl">
            <CardHeader className="pb-4">
              <CardTitle className="text-lg font-semibold text-[#1c3b2e] tracking-tight">{t("feedback.submit_title")}</CardTitle>
            </CardHeader>
            <CardContent className="space-y-6">
              <form onSubmit={handleSubmit} className="space-y-6">
                <div className="space-y-2">
                  <Label htmlFor="feedback-type" className="text-sm font-medium text-[#2e4d57]">
                    {t("feedback.form.type_label")}
                  </Label>
                  <Select value={feedbackType} onValueChange={setFeedbackType}>
                    <SelectTrigger className="bg-white/70 border-[#a7b3b9]/50 hover:border-[#3a6a8d] transition-all duration-300 rounded-xl h-12">
                      <SelectValue placeholder={t("feedback.form.type_placeholder")} />
                    </SelectTrigger>
                    <SelectContent className="bg-white/95 backdrop-blur-md border-[#a7b3b9]/30 rounded-xl">
                      <SelectItem
                        value="inaccuracy"
                        className="hover:bg-[#5e9c8d]/20 rounded-lg transition-colors duration-200"
                      >
                        {t("feedback.form.types.inaccuracy")}
                      </SelectItem>
                      <SelectItem
                        value="improvement"
                        className="hover:bg-[#5e9c8d]/20 rounded-lg transition-colors duration-200"
                      >
                        {t("feedback.form.types.improvement")}
                      </SelectItem>
                      <SelectItem
                        value="other"
                        className="hover:bg-[#5e9c8d]/20 rounded-lg transition-colors duration-200"
                      >
                        {t("feedback.form.types.other")}
                      </SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="subject" className="text-sm font-medium text-[#2e4d57]">
                    {t("feedback.form.subject_label")}
                  </Label>
                  <Input
                    id="subject"
                    placeholder={t("feedback.form.subject_placeholder")}
                    value={subject}
                    onChange={(e) => setSubject(e.target.value)}
                    className="bg-white/70 border-[#a7b3b9]/50 hover:border-[#3a6a8d] focus:border-[#3a6a8d] transition-all duration-300 rounded-xl h-12"
                    disabled={!procedureId || isSubmitting}
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="detailed-feedback" className="text-sm font-medium text-[#2e4d57]">
                    {t("feedback.form.detailed_label")}
                  </Label>
                  <Textarea
                    id="detailed-feedback"
                    placeholder={t("feedback.form.detailed_placeholder")}
                    value={detailedFeedback}
                    onChange={(e) => setDetailedFeedback(e.target.value)}
                    className="min-h-[120px] bg-white/70 border-[#a7b3b9]/50 hover:border-[#3a6a8d] focus:border-[#3a6a8d] transition-all duration-300 rounded-xl resize-none"
                    disabled={!procedureId || isSubmitting}
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="tags" className="text-sm font-medium text-[#2e4d57]">
                    {t("feedback.form.tags_label")}
                  </Label>
                  <Input
                    id="tags"
                    placeholder={t("feedback.form.tags_placeholder")}
                    value={tagsInput}
                    onChange={(e) => setTagsInput(e.target.value)}
                    className="bg-white/70 border-[#a7b3b9]/50 hover:border-[#3a6a8d] focus:border-[#3a6a8d] transition-all duration-300 rounded-xl h-12"
                    disabled={!procedureId || isSubmitting}
                  />
                </div>

                <Button
                  type="submit"
                  disabled={!procedureId || isSubmitting}
                  className="w-full bg-gradient-to-r from-[#3a6a8d] to-[#2e4d57] hover:from-[#2e4d57] hover:to-[#1c3b2e] text-white py-3 rounded-xl font-medium transition-all duration-300 transform hover:scale-[1.02] disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none"
                >
                  <Send className="h-4 w-4 mr-2" />
                  {isSubmitting ? t("feedback.form.submitting") : t("feedback.form.submit_button")}
                </Button>
                {!procedureId && (
                  <p className="text-xs text-[#2e4d57]/70 text-center">
                    {t("feedback.form.no_procedure")}
                  </p>
                )}
              </form>
            </CardContent>
          </Card>

          <Card className="bg-white/80 backdrop-blur-md border-[#a7b3b9]/30 shadow-xl hover:shadow-2xl transition-all duration-500 rounded-2xl">
            <CardHeader className="pb-4">
              <CardTitle className="text-lg font-semibold text-[#1c3b2e] tracking-tight">
                {t("feedback.history_title")}
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              {!procedureId && (
                <div className="text-sm text-[#2e4d57]/70 text-center py-8">
                  {t("feedback.history.no_procedure")}
                </div>
              )}
              {procedureId && isFetching && (
                <div className="text-sm text-[#2e4d57]/70 text-center py-8">{t("feedback.history.loading")}</div>
              )}
              {procedureId && !isFetching && isError && (
                <div className="text-sm text-red-600 text-center py-8">{t("feedback.history.error")}</div>
              )}
              {procedureId && !isFetching && history && (!history.feedbacks || history.feedbacks.length === 0) && (
                <div className="text-sm text-[#2e4d57]/70 text-center py-8">
                  {t("feedback.history.empty")}
                </div>
              )}
              {procedureId && history && Array.isArray(history.feedbacks) && history.feedbacks.length > 0 && (
                <>
                  <div className="space-y-3">
                    {history.feedbacks.map((f, idx) => {
                      if (!f) {
                        console.error("Feedback is undefined at index:", idx);
                        return null;
                      }
                      return (
                        <div
                          key={f.id || idx}
                          className="bg-white/60 backdrop-blur-sm border border-[#a7b3b9]/30 rounded-xl p-4 space-y-3 hover:bg-white/80 transition-all duration-300"
                        >
                          <div className="flex items-center justify-between">
                            <span className="text-sm font-medium text-[#2e4d57] capitalize">
                              {t(`feedback.form.types.${f.type}`) || f.type?.replace(/_/g, " ")}
                            </span>
                            <Badge
                              variant="outline"
                              className="text-xs px-3 py-1 rounded-full border-[#5e9c8d] text-[#1c3b2e] bg-[#5e9c8d]/10"
                            >
                              {t(`feedback.history.status.${f.status}`) || f.status || t("feedback.history.status.new")}
                            </Badge>
                          </div>
                          <p className="text-sm text-[#2e4d57]/80 whitespace-pre-line leading-relaxed">
                            {(f.content || "").slice(0, 220)}
                            {(f.content || "").length > 220 ? "…" : ""}
                          </p>
                          <p className="text-xs text-[#a7b3b9]">
                            {f.createdAT ? new Date(f.createdAT).toLocaleString() : ""}
                          </p>
                        </div>
                      )
                    })}
                  </div>
                  <div className="flex items-center justify-between pt-4 border-t border-[#a7b3b9]/20">
                    <Button
                      variant="outline"
                      size="sm"
                      disabled={page <= 1}
                      onClick={() => setPage((p) => Math.max(1, p - 1))}
                      className="border-[#a7b3b9]/50 text-[#2e4d57] hover:bg-[#5e9c8d]/20 rounded-lg transition-all duration-200"
                    >
                      {t("feedback.history.previous")}
                    </Button>
                    <span className="text-xs text-[#2e4d57]/70 font-medium">
                      {t("feedback.history.page", { current: history.page, total: totalPages })}
                    </span>
                    <Button
                      variant="outline"
                      size="sm"
                      disabled={history.page >= totalPages}
                      onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                      className="border-[#a7b3b9]/50 text-[#2e4d57] hover:bg-[#5e9c8d]/20 rounded-lg transition-all duration-200"
                    >
                      {t("feedback.history.next")}
                    </Button>
                  </div>
                </>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}

export default function Page() {
  return (
    <Suspense fallback={<FeedbackFallback />}>
      <FeedbackPageInner />
    </Suspense>
  )
}
