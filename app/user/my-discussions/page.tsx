"use client"

import { useEffect, useMemo, useState } from "react"
import Image from "next/image"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Input } from "@/components/ui/input"
import { Textarea } from "@/components/ui/textarea"
import { useRouter } from "next/navigation"
import {
  useGetDiscussionsQuery,
  useUpdateDiscussionMutation,
  useDeleteDiscussionMutation,
  type DiscussionPost,
} from "@/app/store/slices/discussionsSlice"
import { Toaster, toast } from "react-hot-toast"

export default function MyDiscussionsPage() {
  const router = useRouter()
  const { data, isLoading, isError, refetch } = useGetDiscussionsQuery({ page: 0, limit: 20 })
  const [updateDiscussion, { isLoading: isUpdating }] = useUpdateDiscussionMutation()
  const [deleteDiscussion, { isLoading: isDeleting }] = useDeleteDiscussionMutation()

  // Inline edit state
  const [editingId, setEditingId] = useState<string | null>(null)
  const [editTitle, setEditTitle] = useState("")
  const [editContent, setEditContent] = useState("")

  // Tag pill styles
  const tagPillClasses = (i: number) => {
    const styles = [
      "bg-green-50 text-green-700 border-green-200 hover:bg-green-100 hover:text-green-800",
      "bg-blue-50 text-blue-700 border-blue-200 hover:bg-blue-100 hover:text-blue-800",
      "bg-teal-50 text-teal-700 border-teal-200 hover:bg-teal-100 hover:text-teal-800",
      "bg-indigo-50 text-indigo-700 border-indigo-200 hover:bg-indigo-100 hover:text-indigo-800",
      "bg-emerald-50 text-emerald-700 border-emerald-200 hover:bg-emerald-100 hover:text-emerald-800",
      "bg-cyan-50 text-cyan-700 border-cyan-200 hover:bg-cyan-100 hover:text-cyan-800",
    ]
    return `text-xs cursor-pointer rounded-full ${styles[i % styles.length]}`
  }

  // Posts from API
  const posts: DiscussionPost[] = useMemo(() => data?.posts ?? [], [data])

  // Edit helpers
  const onEdit = (p: DiscussionPost) => {
    setEditingId(p.ID)
    setEditTitle(p.Title ?? "")
    setEditContent(p.Content ?? "")
  }
  const onCancel = () => {
    setEditingId(null)
    setEditTitle("")
    setEditContent("")
  }
  const onSave = async (id: string) => {
    if (!editTitle.trim() || !editContent.trim()) return
    try {
      await updateDiscussion({ id, data: { title: editTitle.trim(), content: editContent.trim() } }).unwrap()
      onCancel()
      await refetch()
    } catch (e) {
      console.error("Failed to update discussion:", e)
      alert("Failed to update. Please try again.")
    }
  }

  // Delete with toast confirmation
  const onDelete = async (id: string) => {
    toast.custom(
      (t) => (
        <div className="bg-white border border-gray-200 shadow-lg rounded-md p-3 w-80">
          <div className="text-sm font-semibold text-gray-900 mb-1">Delete discussion?</div>
          <div className="text-xs text-gray-600 mb-3">This action cannot be undone.</div>
          <div className="flex justify-end gap-2">
            <button
              className="px-3 py-1 text-sm border border-gray-200 rounded-md hover:bg-gray-50"
              onClick={() => toast.dismiss(t.id)}
            >
              Cancel
            </button>
            <button
              className="px-3 py-1 text-sm rounded-md bg-red-600 text-white hover:bg-red-700"
              onClick={async () => {
                try {
                  toast.dismiss(t.id)
                  await deleteDiscussion(id).unwrap()
                  toast.success("Discussion deleted")
                  await refetch()
                } catch (e) {
                  console.error("Failed to delete discussion:", e)
                  toast.error("Failed to delete. Please try again.")
                }
              }}
              disabled={isDeleting}
            >
              Delete
            </button>
          </div>
        </div>
      ),
      { duration: Infinity }
    )
  }

  return (
    <div className="min-h-screen bg-gray-50 p-4 sm:p-6">
      <Toaster position="top-right" toastOptions={{ duration: 3500 }} />
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between mb-6 gap-3">
          <h1 className="text-2xl sm:text-3xl font-bold text-gray-900">My Discussions</h1>
          <div className="flex gap-2 w-full sm:w-auto">
            <Button
              variant="outline"
              className="border-[#3A6A8D] text-[#3A6A8D] w-full sm:w-auto"
              onClick={() => router.push("/user/discussions")}
            >
              Back to Discussions
            </Button>
          </div>
        </div>

        {isLoading && <div>Loading...</div>}
        {isError && <div>Failed to load discussions.</div>}

        {/* Cards */}
        {!isLoading && !isError && (
          <div className="space-y-4">
            {posts.map((p, index) => {
              const isEditing = editingId === p.ID
              return (
                <Card
                  key={p.ID}
                  className="p-4 sm:p-6 bg-white hover:shadow-lg transition-all duration-300 animate-in fade-in slide-in-from-bottom-4"
                  style={{ animationDelay: `${index * 80}ms` }}
                >
                  <CardContent className="p-0">
                    <div className="flex gap-4 flex-col sm:flex-row">
                      <Image
                        src={"/images/profile-photo.jpg"}
                        alt={p.Title}
                        width={48}
                        height={48}
                        className="w-12 h-12 rounded-full object-cover mx-auto sm:mx-0"
                      />
                      <div className="flex-1">
                        {!isEditing ? (
                          <>
                            <h3 className="text-base font-semibold text-gray-900 mb-2">{p.Title}</h3>
                            <p className="text-gray-700 mb-4 line-clamp-2">{p.Content}</p>
                          </>
                        ) : (
                          <div className="space-y-2 mb-3">
                            <Input
                              value={editTitle}
                              onChange={(e) => setEditTitle(e.target.value)}
                              placeholder="Edit title"
                            />
                            <Textarea
                              value={editContent}
                              onChange={(e) => setEditContent(e.target.value)}
                              rows={4}
                              placeholder="Edit content"
                            />
                          </div>
                        )}

                        {/* Tags */}
                        <div className="flex flex-wrap gap-2">
                          {(p.Tags ?? []).map((t, i) => {
                            const clean = String(t).replace(/^#/, "")
                            return (
                              <Badge key={`${p.ID}-${clean}-${i}`} variant="outline" className={tagPillClasses(i)}>
                                {clean}
                              </Badge>
                            )
                          })}
                        </div>

                        {/* Actions */}
                        <div className="flex gap-2 justify-end mt-3">
                          {!isEditing ? (
                            <>
                              <Button
                                variant="outline"
                                size="sm"
                                className="border-gray-300"
                                onClick={() => onEdit(p)}
                              >
                                Edit
                              </Button>
                              <Button
                                variant="outline"
                                size="sm"
                                className="border-red-300 text-red-600 hover:bg-red-50"
                                disabled={isDeleting}
                                onClick={() => onDelete(p.ID)}
                              >
                                Delete
                              </Button>
                            </>
                          ) : (
                            <>
                              <Button
                                size="sm"
                                className="bg-[#3A6A8D] hover:bg-[#2d5470] text-white"
                                disabled={isUpdating || !editTitle.trim() || !editContent.trim()}
                                onClick={() => onSave(p.ID)}
                              >
                                Save
                              </Button>
                              <Button variant="outline" size="sm" onClick={onCancel}>
                                Cancel
                              </Button>
                            </>
                          )}
                        </div>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              )
            })}
            {posts.length === 0 && <div className="text-gray-600">No discussions found.</div>}
          </div>
        )}
      </div>
    </div>
  )
}
