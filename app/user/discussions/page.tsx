"use client"

import { useState } from "react"
import { Search, Plus, MessageSquare, Heart, Eye, Share2, Pin } from "lucide-react"
import { Flag } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"

export default function CommunityPage() {
  const [searchQuery, setSearchQuery] = useState("")

  const discussions = [
    {
      id: 1,
      author: "Alex Chen",
      avatar: "/user-avatar.png",
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
      avatar: "/user-avatar.png",
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
      avatar: "/user-avatar.png",
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

  const popularTags = [
    { name: "#AI", count: 45 },
    { name: "#StudyTips", count: 32 },
    { name: "#Events", count: 28 },
    { name: "#Notes", count: 24 },
    { name: "#Research", count: 19 },
  ]

  const topContributors = [
    { name: "Emma Wilson", posts: "42 contributions", avatar: "/user-avatar.png" },
    { name: "David Kim", posts: "38 contributions", avatar: "/user-avatar.png" },
    { name: "Lisa Chang", posts: "35 contributions", avatar: "/user-avatar.png" },
  ]

  const pinnedDiscussions = [{ title: "Community Guidelines", author: "Sarah Johnson" }]

  const quickTags = ["#AI", "#passport", "#tax", "#business"]

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="flex items-center justify-between mb-8">
          <div>
            <div className="flex items-center gap-3 mb-2">
              <MessageSquare className="h-8 w-8 text-[#3A6A8D]" />
              <h1 className="text-3xl font-bold text-gray-900">Community Discussions</h1>
            </div>
            <p className="text-gray-600">Join the conversation. Share, ask, and collaborate.</p>
          </div>
          <Button className="bg-[#3A6A8D] hover:bg-[#2d5470] text-white">
            <Plus className="h-4 w-4 mr-2" />
            Add Discussion
          </Button>
        </div>
        {/* Search and Filters */}
        <Card className="p-4 mb-6">
          <div className="flex flex-col sm:flex-row gap-4 w-full">
            <div className="relative flex-1">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
              <input
                type="text"
                placeholder="Search discussions..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full pl-10 pr-4 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#3A6A8D] focus:border-transparent"
              />
            </div>
            <div className="flex gap-2 flex-1">
              <Select defaultValue="all">
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
        </Card>

        {/* Quick Tags */}
        <div className="flex flex-wrap gap-3 mb-6">
          {quickTags.map((tag) => (
            <Badge
              key={tag}
              variant="outline"
              className="cursor-pointer hover:bg-blue-100 hover:text-blue-700 transition-colors"
            >
              {tag}
            </Badge>
          ))}
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
          {/* Main Content */}
          <div className="lg:col-span-9 space-y-6">
            <div className="space-y-6 mt-2">
              {discussions.map((discussion, index) => (
                <Card
                  key={discussion.id}
                  className={`p-6 bg-white hover:shadow-lg transition-all duration-300 cursor-pointer animate-in fade-in slide-in-from-bottom-4`}
                  style={{ animationDelay: `${index * 100}ms` }}
                >
                  <CardContent className="p-0">
                    <div className="flex gap-4">
                      <img
                        src={discussion.avatar || "/placeholder.svg"}
                        alt={discussion.author}
                        className="w-12 h-12 rounded-full object-cover"
                      />
                      <div className="flex-1">
                        <div className="flex items-center gap-2 mb-2">
                          <span className="text-lg font-semibold text-gray-900">{discussion.author}</span>
                          {discussion.isModerator && (
                            <Badge variant="secondary" className="text-xs">
                              Moderator
                            </Badge>
                          )}
                          {discussion.isPinned && <Pin className="h-4 w-4 text-[#3A6A8D]" />}
                          <span className="text-gray-500 text-sm">{discussion.timestamp}</span>
                        </div>

                        <h3 className="text-base font-semibold text-gray-900 mb-2 ">
                          {discussion.title}
                        </h3>

                        <p className="text-gray-600 mb-4 line-clamp-2">{discussion.content}</p>

                        <div className="flex items-center justify-between">
                          <div className="flex flex-wrap gap-2 pb-6">
                            {discussion.tags.map((tag) => (
                              <Badge key={tag} variant="outline" className="text-xs cursor-pointer  hover:bg-blue-100 hover:text-blue-700">
                                {tag}
                              </Badge>
                            ))}
                          </div>

                          <div className="flex items-center gap-4 text-sm text-gray-500">
                            <div className="flex items-center gap-1">
                              <Heart className="h-4 w-4" />
                              <span>{discussion.likes}</span>
                            </div>
                            <div className="flex items-center gap-1">
                              <MessageSquare className="h-4 w-4" />
                              <span>{discussion.replies}</span>
                            </div>
                            <div className="flex items-center gap-1">
                              <Eye className="h-4 w-4" />
                              <span>{discussion.views}</span>
                            </div>
                            <div className="flex items-center gap-1">
                              <Share2 className="h-4 w-4" />
                              <span>{discussion.shares}</span>
                            </div>
                            <button title="Flag this discussion" className="flex items-center gap-1 text-gray-500 hover:text-red-500 transition-colors">
                              <Flag className="h-4 w-4" />
                            </button>
                          </div>
                        </div>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          </div>

          {/* Sidebar */}
          <div className="lg:col-span-3 space-y-6">
            {/* Popular Tags */}
            <Card className="p-4 bg-white">
              <h3 className="font-semibold text-gray-900 mb-4">Popular Tags</h3>
              <div className="space-y-2">
                {popularTags.map((tag) => (
                  <div key={tag.name} className="flex items-center justify-between">
                    <Badge
                      variant="outline"
                      className="cursor-pointer  hover:bg-blue-100 hover:text-blue-700"
                    >
                      {tag.name}
                    </Badge>
                    <span className="text-sm text-gray-500">{tag.count}</span>
                  </div>
                ))}
              </div>
            </Card>

            {/* Top Contributors */}
            <Card className="p-4 bg-white">
              <h3 className="font-semibold text-gray-900 mb-4">Top Contributors</h3>
              <div className="space-y-3">
                {topContributors.map((contributor) => (
                  <div key={contributor.name} className="flex items-center gap-3">
                    <img
                      src={contributor.avatar || "/placeholder.svg"}
                      alt={contributor.name}
                      className="w-8 h-8 rounded-full object-cover"
                    />
                    <div>
                      <p className="font-medium text-gray-900 text-sm">{contributor.name}</p>
                      <p className="text-xs text-gray-500">{contributor.posts}</p>
                    </div>
                  </div>
                ))}
              </div>
            </Card>

            {/* Pinned Discussions */}
            <Card className="p-4 bg-white">
              <h3 className="font-semibold text-gray-900 mb-4">Pinned Discussions</h3>
              <div className="space-y-2">
                {pinnedDiscussions.map((discussion, index) => (
                  <div key={index} className="p-3 bg-blue-50 rounded-lg border-l-4 border-[#3A6A8D]">
                    <p className="font-medium text-gray-900 text-sm">{discussion.title}</p>
                    <p className="text-xs text-gray-500">by {discussion.author}</p>
                  </div>
                ))}
              </div>
            </Card>
          </div>
        </div>
      </div>
    </div>
  )
}
