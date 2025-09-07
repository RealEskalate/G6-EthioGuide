"use client"

import { useParams, useRouter } from "next/navigation"
import { useEffect, useMemo, useState } from "react"
import { useGetChecklistQuery, usePatchChecklistMutation } from "@/app/store/slices/checklistsApi"
import { useSession } from "next-auth/react"
import { Card, CardContent } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import { Badge } from "@/components/ui/badge"
import { ArrowLeft, CheckCircle, Clock, FileText } from "lucide-react"
import { motion } from "framer-motion"

export default function ChecklistDetailPage() {
  const params = useParams<{ id: string }>()
  const id = decodeURIComponent(params.id)
  const router = useRouter()
  const { data: session } = useSession()
  const { data: checklist, isLoading, refetch, error: checklistError } = useGetChecklistQuery({ id, token: session?.accessToken || undefined }, { skip: !id })
  const [patchChecklist, { isLoading: patching }] = usePatchChecklistMutation()
  const [localItems, setLocalItems] = useState(checklist?.items ?? [])

  useEffect(() => {
    setLocalItems(checklist?.items ?? [])
  }, [checklist?.items])

  const progress = useMemo(() => checklist?.progress ?? 0, [checklist?.progress])
  const status = useMemo(() => checklist?.status ?? "NOT_STARTED", [checklist?.status])

  const toggleItem = async (itemId: string, completed: boolean) => {
    // Optimistic update
    setLocalItems((prev) => prev.map((it) => (it.id === itemId ? { ...it, is_checked: completed } : it)))
    try {
      await patchChecklist({ id: itemId, isChecked: completed, token: session?.accessToken || undefined }).unwrap()
      // Refresh to recalc server-side percent/status if provided
      refetch()
    } catch (error) {
      // revert on error
      setLocalItems((prev) => prev.map((it) => (it.id === itemId ? { ...it, is_checked: !completed } : it)))
      console.error('Failed to update checklist item:', error)
    }
  }

  return (
    <main className="min-h-screen w-full bg-gray-50 relative overflow-hidden p-4 sm:p-6 md:p-8">
      {/* subtle brand orbs like workspace */}
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        <div className="absolute -top-24 -right-24 w-56 h-56 rounded-full blur-3xl" style={{ background: 'radial-gradient(closest-side, rgba(167,179,185,0.10), rgba(167,179,185,0))' }} />
        <div className="absolute -bottom-28 -left-28 w-64 h-64 rounded-full blur-3xl" style={{ background: 'radial-gradient(closest-side, rgba(94,156,141,0.10), rgba(94,156,141,0))' }} />
      </div>

      <div className="relative z-10 max-w-3xl mx-auto space-y-6">
        <div className="flex items-center justify-between">
          <Button variant="outline" onClick={() => router.back()} className="border-[#3a6a8d] text-[#3a6a8d] hover:bg-[#3a6a8d]/10">
            <ArrowLeft className="w-4 h-4 mr-2" /> Back
          </Button>
        </div>

        <Card className="bg-white/80 backdrop-blur-md rounded-2xl border border-[#e5e7eb] shadow-xl relative overflow-hidden">
          <div className="absolute inset-0 bg-gradient-to-r from-[#3a6a8d]/10 via-transparent to-[#5e9c8d]/10" />
          <CardContent className="p-6 relative z-10">
            {isLoading ? (
              <div className="text-gray-600">Loading checklist…</div>
            ) : checklistError ? (
              <div className="text-red-600">
                <p>Failed to load checklist.</p>
                <Button 
                  variant="outline" 
                  onClick={() => refetch()} 
                  className="mt-2"
                >
                  Try Again
                </Button>
              </div>
            ) : checklist && (checklist.items?.length ?? 0) >= 0 ? (
              <div className="space-y-4">
                <div className="flex items-start justify-between">
                  <div className="flex items-start gap-3">
                    <div className="w-12 h-12 rounded-xl flex items-center justify-center" style={{ backgroundColor: '#e6f0f5' }}>
                      <FileText className="w-6 h-6" style={{ color: '#3a6a8d' }} />
                    </div>
                    <div>
                      <h1 className="text-2xl font-semibold text-[#111827]">
                        {checklist.procedureTitle || `Procedure ${checklist.procedureId}`}
                      </h1>
                    </div>
                  </div>
                  <div className="flex items-center gap-3">
                    <Badge className={`${
                      status === 'COMPLETED' ? 'bg-green-100 text-green-800' :
                      status === 'IN_PROGRESS' ? 'bg-orange-100 text-orange-800' : 'bg-gray-100 text-gray-800'
                    } border border-[#e5e7eb]`}>{status.replace('_', ' ')}</Badge>
                    <div className="flex items-center gap-2 text-sm text-[#4b5563]">
                      <Clock className="w-4 h-4" />
                      <span>{checklist.updatedAt ? new Date(checklist.updatedAt).toLocaleString() : '—'}</span>
                    </div>
                  </div>
                </div>

                <div>
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-sm font-medium text-[#111827]">Progress</span>
                    <span className="text-sm text-[#4b5563]">{progress}% Complete</span>
                  </div>
                  <div className="w-full bg-[#e5e7eb] rounded-full h-2 overflow-hidden">
                    {(() => {
                      const barColor =
                        status === 'COMPLETED'
                          ? 'bg-[#5e9c8d]'
                          : status === 'IN_PROGRESS'
                            ? 'bg-gradient-to-r from-[#3a6a8d] to-[#2e4d57]'
                            : 'bg-[#a7b3b9]/50';
                      return (
                        <motion.div
                          className={`h-2 rounded-full ${barColor}`}
                          initial={{ width: 0 }}
                          animate={{ width: `${progress}%` }}
                          transition={{ type: 'spring', stiffness: 180, damping: 24, mass: 0.9 }}
                        />
                      );
                    })()}
                  </div>
                </div>

                <div className="space-y-3">
                  <h2 className="text-lg font-medium text-[#111827]">Checklist Items</h2>
                  {localItems && localItems.length > 0 ? (
                    <ul className="space-y-2">
                      {localItems.map((item) => (
                        <li key={item.id} className="flex items-start gap-3 p-3 rounded-xl border border-[#e5e7eb] bg-white/90 hover:bg-white transition-colors">
                          <Checkbox
                            checked={item.is_checked}
                            onCheckedChange={(v) => toggleItem(item.id, Boolean(v))}
                            className="mt-0.5"
                            disabled={patching}
                          />
                          <div className="flex-1">
                            <div className="flex items-center gap-2">
                              <span className="font-medium text-[#111827]">{item.content || 'Item'}</span>
                              {item.is_checked && (
                                <CheckCircle className="w-4 h-4 text-green-600" />
                              )}
                            </div>
                            {item.type && (
                              <p className="text-sm text-[#4b5563]">Type: {item.type}</p>
                            )}
                          </div>
                        </li>
                      ))}
                    </ul>
                  ) : (
                    <div className="text-[#4b5563] text-sm">No items yet for this checklist.</div>
                  )}
                </div>

                {/* Actions removed per request: no refresh or back-to-workspace */}
              </div>
            ) : (
              <div className="text-red-600">Failed to load checklist.</div>
            )}
          </CardContent>
        </Card>
      </div>
    </main>
  )
}
