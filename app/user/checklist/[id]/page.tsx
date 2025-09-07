"use client"

import { useParams, useRouter } from "next/navigation"
import { useEffect, useMemo, useState } from "react"
import { useGetChecklistQuery, usePatchChecklistMutation } from "@/app/store/slices/checklistsApi"
import { useSession } from "next-auth/react"
import { Card, CardContent } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import { Badge } from "@/components/ui/badge"
import { ArrowLeft, CheckCircle, Clock } from "lucide-react"

export default function ChecklistDetailPage() {
  const params = useParams<{ id: string }>()
  const id = decodeURIComponent(params.id)
  const router = useRouter()
  const { data: session } = useSession()
  const { data: checklist, isLoading, isFetching, refetch, error: checklistError } = useGetChecklistQuery({ id, token: session?.accessToken || undefined }, { skip: !id })
  const [patchChecklist, { isLoading: patching, error: patchError }] = usePatchChecklistMutation()
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
    <main className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-3xl mx-auto space-y-6">
        <div className="flex items-center justify-between">
          <Button variant="ghost" onClick={() => router.back()} className="hover:bg-gray-100">
            <ArrowLeft className="w-4 h-4 mr-2" /> Back
          </Button>
        </div>

        <Card className="border-0 shadow-sm bg-white">
          <CardContent className="p-6">
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
                  <div>
                    <h1 className="text-2xl font-semibold text-gray-900">
                      {checklist.procedureTitle || `Procedure ${checklist.procedureId}`}
                    </h1>
                    <p className="text-gray-600 text-sm">ID: {checklist.id}</p>
                  </div>
                  <div className="flex items-center gap-3">
                    <Badge className={`${
                      status === 'COMPLETED' ? 'bg-green-100 text-green-800' :
                      status === 'IN_PROGRESS' ? 'bg-orange-100 text-orange-800' : 'bg-gray-100 text-gray-800'
                    }`}>{status.replace('_', ' ')}</Badge>
                    <div className="flex items-center gap-2 text-sm text-gray-600">
                      <Clock className="w-4 h-4" />
                      <span>{checklist.updatedAt ? new Date(checklist.updatedAt).toLocaleString() : '—'}</span>
                    </div>
                  </div>
                </div>

                <div>
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-sm font-medium text-gray-700">Progress</span>
                    <span className="text-sm text-gray-600">{progress}% Complete</span>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2 overflow-hidden">
                    <div
                      className={`h-2 rounded-full ${
                        status === 'COMPLETED' ? 'bg-[#5E9C8D]' : status === 'IN_PROGRESS' ? 'bg-[#FEF9C3]' : 'bg-gray-300'
                      }`}
                      style={{ width: `${progress}%` }}
                    />
                  </div>
                </div>

                <div className="space-y-3">
                  <h2 className="text-lg font-medium text-gray-900">Checklist Items</h2>
                  {localItems && localItems.length > 0 ? (
                    <ul className="space-y-2">
                      {localItems.map((item) => (
                        <li key={item.id} className="flex items-start gap-3 p-3 rounded-md border border-gray-200 bg-white">
                          <Checkbox
                            checked={item.is_checked}
                            onCheckedChange={(v) => toggleItem(item.id, Boolean(v))}
                            className="mt-0.5"
                            disabled={patching}
                          />
                          <div className="flex-1">
                            <div className="flex items-center gap-2">
                              <span className="font-medium text-gray-900">{item.content || 'Item'}</span>
                              {item.is_checked && (
                                <CheckCircle className="w-4 h-4 text-green-600" />
                              )}
                            </div>
                            {item.type && (
                              <p className="text-sm text-gray-600">Type: {item.type}</p>
                            )}
                          </div>
                        </li>
                      ))}
                    </ul>
                  ) : (
                    <div className="text-gray-600 text-sm">No items yet for this checklist.</div>
                  )}
                </div>

                <div className="flex items-center justify-end gap-3 pt-2">
                  <Button variant="outline" onClick={() => refetch()} disabled={isFetching || patching}>
                    Refresh
                  </Button>
                  <Button onClick={() => router.push('/user/workspace')} className="bg-[#3A6A8D] hover:bg-[#2d5470] text-white">
                    Back to Workspace
                  </Button>
                </div>
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
