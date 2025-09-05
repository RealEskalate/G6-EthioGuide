"use client"

import { useEffect, useState } from "react"
import Image from "next/image"
import { Search, Plus, MessageSquare } from "lucide-react"
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
  const { data, isLoading, isError } = useGetDiscussionsQuery({ page: 0, limit: 10 })

  useEffect(() => {
    if (data) {
      console.log("Discussions list:", data)
    }
  }, [data])

  const [expandedMap, setExpandedMap] = useState<Record<string, boolean>>({})

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
    <div className="min-h-screen bg-gray-50 p-2 sm:p-4">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between mb-8 gap-4">
          <div>
            <div className="flex items-center gap-3 mb-2">
              <MessageSquare className="h-8 w-8 text-[#3A6A8D]" />
              <h1 className="text-2xl sm:text-3xl font-bold text-gray-900">Community Discussions</h1>
            </div>
            <p className="text-gray-600 text-sm sm:text-base">Join the conversation. Share, ask, and collaborate.</p>
          </div>
          <div className="flex w-full sm:w-auto gap-2">
            <Button
              variant="outline"
              className="border-[#3A6A8D] text-[#3A6A8D] w-full sm:w-auto"
              onClick={() => router.push("/user/my-discussions")}
            >
              My Discussions
            </Button>
            <Button
              className="bg-[#3A6A8D] hover:bg-[#2d5470] text-white w-full sm:w-auto"
              onClick={() => router.push("/user/create-post")}
            >
              <Plus className="h-4 w-4 mr-2" />
              Add Discussion
            </Button>
          </div>
        </div>
        {/* Search and Filters */}
        <Card className="p-4 mb-6">
          <div className="flex flex-col gap-4 w-full mb-2 sm:flex-row">
            <div className="relative flex-1 flex">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
              <input
                type="text"
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
            <div className="flex gap-2 flex-1">
              <Select
                value={selectedCategory}
                onValueChange={setSelectedCategory}
              >
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
                  <SelectValue placeholder="Latest" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="latest">Latest</SelectItem>
                  <SelectItem value="popular">Popular</SelectItem>
                  <SelectItem value="trending">Trending</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
          {/* Quick Tags (API-only). Hidden if API returned no tags. */}
          {searchTags.length > 0 && (
            <div className="flex flex-wrap gap-3 mb-3">
              {searchTags.map((tag, i) => (
                <Badge key={tag} variant="outline" className={tagPillClasses(i)}>
                  {tag}
                </Badge>
              ))}
            </div>
          )}
        </Card>
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
          {/* Main Content */}
          <div className="lg:col-span-9 space-y-6">
            <div className="space-y-6 mt-2">
              {filteredDiscussions.map((discussion, index) => {
                const rowKey = `${discussion.title}-${index}`
                const isExpanded = !!expandedMap[rowKey]
                return (
                  <Card
                    key={discussion.title}
                    className="p-4 sm:p-6 bg-white hover:shadow-lg transition-all duration-300 cursor-pointer animate-in fade-in slide-in-from-bottom-4"
                    style={{ animationDelay: `${index * 100}ms` }}
                  >
                    <CardContent className="p-0">
                      <div className="flex gap-4 flex-col sm:flex-row">
                        <Image
                          src={discussion.avatar || "/placeholder.svg"}
                          alt={discussion.title}
                          width={48}
                          height={48}
                          className="w-12 h-12 rounded-full object-cover mx-auto sm:mx-0"
                        />
                        <div className="flex-1">
                          <h3 className="text-base font-semibold text-gray-900 mb-2">
                            {discussion.title}
                          </h3>
                          <p className={`text-gray-700 mb-4 ${isExpanded ? "" : "line-clamp-2"}`}>
                            {discussion.content}
                          </p>
                          <div className="flex flex-wrap gap-2">
                            {discussion.tags.map((tag: string, i: number) => {
                              const clean = tag.replace(/^#/, "")
                              return (
                                <Badge
                                  key={`${discussion.title}-${clean}-${i}`}
                                  variant="outline"
                                  className={`text-xs ${tagPillClasses(i)}`}
                                >
                                  {clean}
                                </Badge>
                              )
                            })}
                          </div>
                          <div className="flex justify-end pt-2">
                            <Button
                              variant="ghost"
                              size="sm"
                              className="text-[#3A6A8D] hover:bg-blue-100 hover:text-blue-700"
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
            </div>
          </div>

          {/* Sidebar */}
          <div className="lg:col-span-3 space-y-6">
            {/* Popular Tags (from API/dummy) */}
            <Card className="p-4 bg-white hover:shadow-xl transform transition duration-300 hover:scale-105">
              <h3 className="font-semibold text-gray-900 mb-4">Popular Tags</h3>
              <div className="space-y-2">
                {popularTags.map((tag, i) => (
                  <div key={tag.name} className="flex items-center justify-between">
                    <Badge variant="outline" className={tagPillClasses(i)}>
                      {tag.name}
                    </Badge>
                    <span className="text-sm text-gray-500">{tag.count}</span>
                  </div>
                ))}
              </div>
            </Card>
          </div>
        </div>
        {isLoading && <div>Loading...</div>}
        {isError && <div>Failed to load discussions.</div>}
        {!isLoading && !isError && data && (
          <div className="text-sm text-gray-700">
            Total: {data.total} • Page: {data.page} • Limit: {data.limit}
          </div>
        )}
      </div>
    </div>
  )
}