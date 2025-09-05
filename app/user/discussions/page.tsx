"use client"

import { useEffect, useState } from "react"
import Image from "next/image"
import { Search, Plus, MessageSquare, Filter, X, ChevronRight, ChevronLeft } from "lucide-react" // added icons
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { useRouter } from "next/navigation"
import { useGetDiscussionsQuery } from "@/app/store/slices/discussionsSlice"

export default function CommunityPage() {
  const [searchInput, setSearchInput] = useState("")
  const [searchQuery, setSearchQuery] = useState("")
  const [selectedCategory, setSelectedCategory] = useState("all")
  const router = useRouter()

  // added: pagination state
  const [page, setPage] = useState(0)
  const limit = 10
  const { data, isLoading, isError } = useGetDiscussionsQuery({ page, limit })

  // added: mobile sidebar toggle
  const [mobilePanelOpen, setMobilePanelOpen] = useState(false)

  // added: compute totalPages once for reuse (desktop + sticky bar)
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

  // reset expanded per page
  useEffect(() => {
    setExpandedMap({})
  }, [page])

  const tagPillClasses = (i: number) => {
    const styles = [
      "bg-green-50 text-green-700 border-green-200 hover:bg-green-100 hover:text-green-800",
      "bg-blue-50 text-blue-700 border-blue-200 hover:bg-blue-100 hover:text-blue-800",
      "bg-teal-50 text-teal-700 border-teal-200 hover:bg-teal-100 hover:text-teal-800",
      "bg-indigo-50 text-indigo-700 border-indigo-200 hover:bg-indigo-100 hover:text-indigo-800",
      "bg-emerald-50 text-emerald-700 border-emerald-200 hover:bg-emerald-100 hover:text-emerald-800",
      "bg-cyan-50 text-cyan-700 border-cyan-200 hover:bg-cyan-100 hover:text-cyan-800",
    ]
    return `cursor-pointer rounded-full ${styles[i % styles.length]}`
  }

  const discussions = [
    {
      id: 1,
      author: "Alex Chen",
      avatar: "/images/profile-photo.jpg",
      timestamp: "2h ago",
      title: "How to integrate AI tools into daily study routine?",
      content:
        "I've been experimenting with various AI tools for studying and note-taking. Would love to hear your experiences and recommendations for the best workflow...",
      tags: ["#AI", "#StudyTips"],
      likes: 24,
      replies: 12,
      views: 156,
      shares: 8,
    },
    {
      id: 2,
      author: "Sarah Johnson",
      avatar: "/images/profile-photo.jpg",
      timestamp: "4h ago",
      title: "📌 Welcome to the Community Guidelines",
      content:
        "Please take a moment to read our community guidelines to ensure a positive and productive environment for everyone. Let's build an amazing learning community together!",
      tags: ["#Guidelines", "#Pinned"],
      likes: 56,
      replies: 8,
      views: 234,
      shares: 15,
      isPinned: true,
      isModerator: true,
    },
    {
      id: 3,
      author: "Mike Rodriguez",
      avatar: "/images/profile-photo.jpg",
      timestamp: "6h ago",
      title: "Best note-taking apps for university students?",
      content:
        "Looking for recommendations on digital note-taking apps that work well for lectures, research, and collaboration. Currently using Notion but wondering if there are better alternatives...",
      tags: ["#Notes", "#Apps"],
      likes: 18,
      replies: 15,
      views: 89,
      shares: 4,
    },
  ]

  const apiDiscussions =
    Array.isArray(data?.posts)
      ? data!.posts.map((p) => ({
          id: p.ID,
          author: p.UserID || "User",
          avatar: "/images/profile-photo.jpg",
          timestamp: new Date(p.CreatedAt || p.UpdatedAt || Date.now()).toLocaleString(),
          title: p.Title ?? "Untitled",
          content: p.Content ?? "",
          tags: Array.isArray(p.Tags) ? p.Tags.map((t) => String(t)) : [],
        }))
      : []

  const discussionsData = apiDiscussions.length ? apiDiscussions : discussions

  const quickTags = [
    "#AI",
    "#StudyTips",
    "#Guidelines",
    "#Pinned",
    "#Notes",
    "#Apps",
    "#Business",
    "#Technology",
    "#Collaboration",
    "#Learning",
    "#University",
    "#Productivity",
    "#Motivation",
    "#Career",
    "#Networking",
    "#Events",
    "#Resources",
    "#Advice",
    "#Support",
    "#Community",
  ]

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
      list = quickTags.map((t) => ({ name: t.replace(/^#/, ""), count: 1 }))
    }
    return list.sort((a, b) => b.count - a.count).slice(0, 20)
  })()

  const searchTags = (() => {
    if (!Array.isArray(data?.posts)) return []
    const set = new Set<string>()
    data!.posts.forEach((p) => {
      const tags = Array.isArray(p.Tags) ? p.Tags : []
      tags.forEach((t) => {
        const clean = String(t).replace(/^#/, "").trim()
        if (clean) set.add(clean)
      })
    })
    return Array.from(set)
  })()

  const filteredDiscussions = discussionsData.filter((discussion) => {
    const matchesSearch =
      discussion.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
      discussion.content.toLowerCase().includes(searchQuery.toLowerCase())
    const matchesCategory =
      selectedCategory === "all" ||
      (selectedCategory === "ai" && discussion.tags.some((tag: string) => tag.toLowerCase().includes("ai"))) ||
      (selectedCategory === "study" && discussion.tags.some((tag: string) => tag.toLowerCase().includes("study"))) ||
      (selectedCategory === "business" && discussion.tags.some((tag: string) => tag.toLowerCase().includes("business")))
    return matchesSearch && matchesCategory
  })

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col">
      <div className="max-w-7xl mx-auto w-full flex-1">
        {/* Header */}
        {/* ...existing header wrapper... */}
        <div className="bg-white/90 border border-gray-100 rounded-xl p-4 sm:p-5 mb-4 sm:mb-6 shadow-sm">
          <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
            <div className="w-full">
              <div className="flex items-center gap-3 mb-1">
                <MessageSquare className="h-7 w-7 text-[#3A6A8D]" />
                <h1 className="text-xl leading-snug sm:text-3xl font-bold text-gray-900">
                  Community Discussions
                </h1>
              </div>
              <p className="text-gray-600 text-xs sm:text-sm md:text-base">
                Join the conversation. Share, ask, and collaborate.
              </p>
            </div>
            <div className="hidden sm:flex gap-2">
              <Button
                variant="outline"
                className="border-[#3A6A8D] text-[#3A6A8D] hover:bg-[#3A6A8D]/5"
                onClick={() => router.push("/user/my-discussions")}
              >
                My Discussions
              </Button>
              <Button
                className="bg-[#3A6A8D] hover:bg-[#2d5470] text-white"
                onClick={() => router.push("/user/create-post")}
              >
                <Plus className="h-4 w-4 mr-2" />
                Add Discussion
              </Button>
            </div>
            {/* Mobile quick button to open filters */}
            <div className="w-full sm:hidden flex gap-2">
              <Button
                variant="outline"
                className="flex-1 border-[#3A6A8D] text-[#3A6A8D] hover:bg-[#3A6A8D]/5"
                onClick={() => setMobilePanelOpen(true)}
              >
                <Filter className="h-4 w-4 mr-2" />
                Filters & Tags
              </Button>
            </div>
          </div>
        </div>

        {/* Search & Filters (desktop / tablet) */}
        <Card className="p-4 mb-4 sm:mb-6 hidden sm:block">
          {/* ...existing search/filter block unchanged... */}
          <div className="flex flex-col gap-4 w-full mb-2 sm:flex-row">
            <div className="relative flex-1 flex">
              {/* ...existing code... */}
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
              <input
                type="text"
                // ...existing props...
                placeholder="Search discussions..."
                value={searchInput}
                onChange={(e) => setSearchInput(e.target.value)}
                className="w-full pl-10 pr-4 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#3A6A8D] focus:border-transparent"
              />
              <Button
                type="button"
                className="ml-2 px-4 py-2 bg-[#3A6A8D] hover:bg-[#2d5470] text-white"
                onClick={() => setSearchQuery(searchInput)}
              >
                Search
              </Button>
            </div>
            {/* ...existing selects... */}
            <div className="flex gap-2 flex-1">
              {/* ...existing category select... */}
              <Select value={selectedCategory} onValueChange={setSelectedCategory}>
                <SelectTrigger className="w-full">
                  <SelectValue placeholder="All Categories" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Categories</SelectItem>
                  <SelectItem value="ai">AI & Technology</SelectItem>
                  <SelectItem value="study">Study Tips</SelectItem>
                  <SelectItem value="business">Business</SelectItem>
                </SelectContent>
              </Select>
              <Select defaultValue="latest">
                <SelectTrigger className="w-full">
                  <SelectValue placeholder="Sort" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="latest">Latest</SelectItem>
                  <SelectItem value="popular">Popular</SelectItem>
                  <SelectItem value="trending">Trending</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
          {searchTags.length > 0 && (
            <div className="flex gap-2 overflow-x-auto pb-1 scrollbar-thin scrollbar-track-transparent scrollbar-thumb-gray-300">
              {searchTags.map((tag, i) => (
                <Badge key={tag} variant="outline" className={`${tagPillClasses(i)} flex-shrink-0`}>
                  {tag}
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
              <div className="absolute bottom-0 left-0 right-0 bg-white rounded-t-2xl shadow-xl max-h-[85vh] flex flex-col">
                <div className="flex items-center justify-between px-4 pt-4 pb-2 border-b">
                  <h2 className="text-base font-semibold text-gray-800">Filters & Tags</h2>
                  <Button variant="ghost" size="sm" onClick={() => setMobilePanelOpen(false)}>
                    <X className="h-5 w-5" />
                  </Button>
                </div>
                <div className="p-4 space-y-5 overflow-y-auto">
                  {/* Mobile search */}
                  <div>
                    <div className="relative">
                      <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
                      <input
                        type="text"
                        placeholder="Search discussions..."
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
                      Apply
                    </Button>
                  </div>
                  {/* Category select */}
                  <div className="space-y-2">
                    <label className="text-xs uppercase tracking-wide text-gray-500 font-medium">
                      Category
                    </label>
                    <Select value={selectedCategory} onValueChange={setSelectedCategory}>
                      <SelectTrigger className="w-full">
                        <SelectValue placeholder="All Categories" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="all">All Categories</SelectItem>
                        <SelectItem value="ai">AI & Technology</SelectItem>
                        <SelectItem value="study">Study Tips</SelectItem>
                        <SelectItem value="business">Business</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  {/* Tags */}
                  {searchTags.length > 0 && (
                    <div className="space-y-2">
                      <label className="text-xs uppercase tracking-wide text-gray-500 font-medium">
                        Popular Tags
                      </label>
                      <div className="flex flex-wrap gap-2">
                        {searchTags.map((tag, i) => (
                          <Badge
                            key={tag}
                            variant="outline"
                            className={`${tagPillClasses(i)} px-3 py-1`}
                          >
                            {tag}
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
                      Reset
                    </Button>
                  </div>
                </div>
              </div>
            </div>
        )}

        <div className="grid grid-cols-1 lg:grid-cols-12 gap-4 sm:gap-6">
          {/* Main Content */}
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
                    <Card
                      key={rowKey}
                      className="p-3 sm:p-6 bg-white border border-gray-100 rounded-xl shadow-sm hover:shadow-md transition-all"
                    >
                      <CardContent className="p-0">
                        <div className="flex gap-3 sm:gap-4 flex-col sm:flex-row">
                          <Image
                            src={discussion.avatar || "/placeholder.svg"}
                            alt={discussion.title}
                            width={44}
                            height={44}
                            className="w-11 h-11 sm:w-12 sm:h-12 rounded-full object-cover mx-auto sm:mx-0 ring-2 ring-white shadow"
                          />
                          <div className="flex-1">
                            <h3 className="text-sm sm:text-base font-semibold text-gray-900 mb-2">
                              {discussion.title}
                            </h3>
                            <p className={`text-gray-700 text-sm sm:text-[15px] mb-3 ${isExpanded ? "" : "line-clamp-2"}`}>
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
                                    {clean}
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
                                {isExpanded ? "View Less" : "View More"}
                              </Button>
                            </div>
                          </div>
                        </div>
                      </CardContent>
                    </Card>
                  )
                })}
                {filteredDiscussions.length === 0 && !isLoading && (
                  <div className="text-center text-sm text-gray-500 py-10 border border-dashed rounded-lg">
                    No discussions match your filters.
                  </div>
                )}
              </div>
            )}

            {/* Desktop pagination footer */}
            {totalPages > 1 && (
              <div className="hidden sm:flex mt-2 bg-white/90 border border-gray-100 rounded-xl px-3 py-3 items-center justify-between shadow-sm">
                <div className="text-sm text-gray-600">
                  Page {page + 1} of {totalPages}
                </div>
                <div className="flex gap-2">
                  <Button
                    variant="outline"
                    className="border-gray-300"
                    disabled={page <= 0 || isLoading}
                    onClick={() => setPage((p) => Math.max(0, p - 1))}
                  >
                    Previous
                  </Button>
                  <Button
                    className="bg-[#3A6A8D] hover:bg-[#2d5470] text-white"
                    disabled={page >= totalPages - 1 || isLoading}
                    onClick={() => setPage((p) => p + 1)}
                  >
                    Next
                  </Button>
                </div>
              </div>
            )}
          </div>

            {/* Sidebar (desktop only) */}
            <div className="lg:col-span-3 space-y-6 hidden lg:block">
              <Card className="p-4 bg-white hover:shadow-xl transform transition duration-300 hover:scale-[1.02]">
                <h3 className="font-semibold text-gray-900 mb-4 text-sm sm:text-base">Popular Tags</h3>
                <div className="space-y-2">
                  {popularTags.map((tag, i) => (
                    <div key={tag.name} className="flex items-center justify-between text-sm">
                      <Badge variant="outline" className={tagPillClasses(i)}>
                        {tag.name}
                      </Badge>
                      <span className="text-xs text-gray-500">{tag.count}</span>
                    </div>
                  ))}
                </div>
              </Card>
            </div>
        </div>

        {/* Status / meta */}
        <div className="mt-6 text-center sm:text-left">
          {isLoading && <div className="text-sm text-gray-500">Loading...</div>}
          {isError && <div className="text-sm text-red-600">Failed to load discussions.</div>}
          {!isLoading && !isError && data && (
            <div className="text-xs sm:text-sm text-gray-600">
              Total: {data.total} • Page: {data.page} • Limit: {data.limit}
            </div>
          )}
        </div>

        {/* Footer */}
        <footer className="mt-10 sm:mt-14 text-center text-[11px] sm:text-xs text-gray-500 border-t pt-6 pb-24 sm:pb-6">
          © {new Date().getFullYear()} EthioGuide. All rights reserved.
        </footer>
      </div>

      {/* Sticky mobile action bar */}
      <div className="sm:hidden fixed bottom-0 left-0 right-0 z-40 bg-white/95 backdrop-blur border-t shadow-lg px-2 py-2 flex items-center justify-between">
        <Button
          variant="outline"
          size="sm"
          className="flex-1 mx-1 text-[#3A6A8D] border-[#3A6A8D]"
          onClick={() => router.push("/user/my-discussions")}
        >
          Me
        </Button>
        <Button
          size="sm"
          className="flex-1 mx-1 bg-[#3A6A8D] hover:bg-[#2d5470] text-white"
          onClick={() => router.push("/user/create-post")}
        >
          <Plus className="h-4 w-4 mr-1" /> New
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