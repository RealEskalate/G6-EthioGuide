"use client"
import { Search, ChevronLeft, ChevronRight } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent } from "@/components/ui/card"
import Image from "next/image"
import Link from "next/link"
import { useState, useMemo, useEffect } from "react"
import { useSearchParams } from "next/navigation"
import { useListCategoriesQuery } from '@/app/store/slices/categoriesApi'
import { useListProceduresQuery } from '@/app/store/slices/proceduresApi'
import type { Procedure } from '@/app/types/procedure'

function badgeColorFromIndex(i: number) {
  const palette = [
    'bg-blue-100 text-blue-700',
    'bg-green-100 text-green-700',
    'bg-purple-100 text-purple-700',
    'bg-orange-100 text-orange-700',
    'bg-red-100 text-red-700',
    'bg-teal-100 text-teal-700'
  ];
  return palette[i % palette.length];
}

export default function ProceduresPage() {
  const searchParams = useSearchParams()
  
  // Pagination state (simple)
  const [page, setPage] = useState(1)
  const limit = 10
  const [search, setSearch] = useState("");
  const [category, setCategory] = useState("all");
  
  // Initialize from URL params
  useEffect(() => {
    const q = searchParams.get('q')
    const cat = searchParams.get('category')
    if (q) setSearch(q)
    if (cat) setCategory(cat)
  }, [searchParams])
  
  // API calls
  const { data: procData, isLoading: loadingProcedures, isError: procsError } = useListProceduresQuery({ 
    page, 
    limit, 
    name: search || undefined,
    sortBy: 'createdAt',
    sortOrder: 'DESC'
  })
  const { data: catData, isLoading: loadingCats, isError: catsError } = useListCategoriesQuery();

  const categories = useMemo(() => {
    const list = catData?.list || []
    return list.map((c, idx) => ({ id: c.id, title: c.title, badge: badgeColorFromIndex(idx) }))
  }, [catData])

  const backendProcedures = useMemo(() => (procData?.list || []).map((p: Procedure, idx) => {
    // Derive fields to fit existing UI expectations
    const title = p.title || p.name || 'Untitled Procedure'
    const description = p.summary || ''
    const updated = p.updatedAt ? new Date(p.updatedAt).toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' }) : '—'
    // Use first tag as category if available
    const tag = (p.tags && p.tags.length > 0) ? p.tags[0] : 'General'
    return {
      id: p.id,
      name: title,
      description,
      category: tag,
      categoryColor: badgeColorFromIndex(idx),
      processingTime: p.processingTime ? `${p.processingTime.minDays ?? '?'}-${p.processingTime.maxDays ?? '?' } days` : '—',
      lastUpdated: updated,
  icon: '/icons/manage-procedure.svg',
      iconBg: 'bg-blue-50',
      bookmarked: false
    }
  }), [procData])

  // Since search is now handled by the API, we only need to filter by category
  const filteredProcedures = backendProcedures.filter((procedure) => {
    if (category === 'all') return true
    return procedure.category === category
  })

  const updates = [
    {
      type: "warning",
      title: "New Document Requirements",
      description: "Updated documentation requirements for all visa applications effective January 2025.",
      date: "Dec 20, 2024",
      icon: "/icons/document-warning.svg",
    },
    {
      type: "success",
      title: "Holiday Office Hours",
      description: "Immigration Office will be closed during the holiday season. Plan your applications accordingly.",
      date: "Dec 18, 2024",
      icon: "/icons/office-hours-info.svg",
    },
    {
      type: "info",
      title: "Online Portal Upgrade",
      description: "New features added to our online application system for faster processing.",
      date: "Dec 15, 2024",
      icon: "/icons/portal-upgrade.svg",
    },
  ]

  return (
    <div className="relative w-full h-full min-h-screen">
      <div className="absolute inset-0 w-full h-full bg-gray-50 -z-10"></div>
      <div className="flex">
        {/* Main Content */}
        <main className="flex-1 p-8">
          <div className="max-w-7xl mx-auto">
            {/* Page Header */}
            <div className="mb-8">
              <h1 className="text-2xl font-semibold text-[#111827] mb-2">Procedures Offered by Immigration Office</h1>
              <p className="text-[#4b5563]">
                Browse through all available immigration procedures and services. Find detailed information about visa
                applications, work permits, citizenship processes, and more.
              </p>
            </div>

            {/* added: quick redirect to Workspace */}
            <div className="mb-4 flex justify-end">
              <Link href="/user/workspace">
                <Button className="bg-[#3A6A8D] hover:bg-[#2e4d57] text-white">
                  Go to Workspace
                </Button>
              </Link>
            </div>

            {/* Search and Filters */}
            <div className="flex items-center justify-between gap-4 mb-8">
              <div className="flex items-center gap-4">
                <div className="relative flex-1 max-w-md">
                  <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-[#4b5563]" />
                  <Input
                    placeholder="Search procedures by keyword..."
                    className="pl-10 bg-white border-[#d1d5db]"
                    value={search}
                    onChange={e => setSearch(e.target.value)}
                  />
                </div>
                <Select value={category} onValueChange={setCategory} disabled={loadingCats || catsError}>
                  <SelectTrigger className="w-60 bg-white border-[#d1d5db] hover:!bg-gray-100 font-medium">
                    <SelectValue placeholder={loadingCats ? 'Loading...' : 'All Categories'} />
                  </SelectTrigger>
                  <SelectContent className="bg-white max-h-72 overflow-y-auto">
                    <SelectItem value="all" className="hover:!bg-gray-100 font-medium border-0">All Categories</SelectItem>
                    {categories.map(c => (
                      <SelectItem key={c.id} value={c.title} className="hover:!bg-gray-100 font-medium border-0">{c.title}</SelectItem>
                    ))}
                    {(!loadingCats && categories.length === 0) && (
                      <div className="px-3 py-2 text-xs text-gray-500">No categories</div>
                    )}
                  </SelectContent>
                </Select>
              </div>
            </div>

            {/* Procedures Table */}
            <div className="bg-white rounded-lg border border-[#e5e7eb] overflow-hidden mb-8">
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead className="bg-[#f3f4f6] border-b border-[#e5e7eb]">
                    <tr>
                      <th className="text-left py-3 px-6 text-sm font-medium text-[#4b5563]">Procedure Name</th>
                      <th className="text-left py-3 px-6 text-sm font-medium text-[#4b5563]">Category</th>
                      <th className="text-left py-3 px-6 text-sm font-medium text-[#4b5563]">Processing Time</th>
                      <th className="text-left py-3 px-6 text-sm font-medium text-[#4b5563]">Last Updated</th>
                      <th className="text-left py-3 px-6 text-sm font-medium text-[#4b5563]">Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    {loadingProcedures && (
                      <tr><td colSpan={5} className="py-8 text-center text-sm text-gray-500">Loading procedures...</td></tr>
                    )}
                    {(!loadingProcedures && procsError) && (
                      <tr><td colSpan={5} className="py-8 text-center text-sm text-red-600">Failed to load procedures.</td></tr>
                    )}
                    {(!loadingProcedures && !procsError && filteredProcedures.length === 0) && (
                      <tr><td colSpan={5} className="py-8 text-center text-sm text-gray-500">No procedures found.</td></tr>
                    )}
                    {filteredProcedures.map((procedure, index) => (
                      <tr
                        key={procedure.id}
                        className={(index !== filteredProcedures.length - 1 ? "border-b border-[#f3f4f6] " : "") + " hover:bg-gray-100 transition-colors"}
                      >
                        <td className="py-4 px-6">
                          <div className="flex items-center gap-3">
                            <div
                              className={`w-10 h-10 ${procedure.iconBg} rounded-lg flex items-center justify-center`}
                            >
                              <Image
                                src={procedure.icon || "/placeholder.svg"}
                                alt={procedure.name || 'Procedure'}
                                width={20}
                                height={20}
                                className="w-5 h-5"
                              />
                            </div>
                            <div>
                              <div className="font-medium text-sm text-[#111827]">{procedure.name}</div>
                        {/* Footer is now part of the shared layout */}
                              <div className="text-sm text-[#4b5563]">{procedure.description}</div>
                            </div>
                          </div>
                        </td>
                        <td className="py-4 px-6">
                          <Badge className={`${procedure.categoryColor || 'bg-gray-100 text-gray-700'} border-0`}>{procedure.category}</Badge>
                        </td>
                        <td className="py-4 px-6 text-[#4b5563]">{procedure.processingTime}</td>
                        <td className="py-4 px-6 text-[#4b5563]">{procedure.lastUpdated}</td>
                        <td className="py-4 px-6">
                          <Link href={`/user/procedures-detail?id=${procedure.id}`}>
                            <Button className="bg-[#3a6a8d] hover:bg-[#2e4d57] text-white">Start</Button>
                          </Link>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>

            {/* Related Updates & Notices */}
            <div className="mb-8">
              <h2 className="text-xl font-semibold text-[#111827] mb-6">Related Updates & Notices</h2>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                {updates.map((update, index) => (
                  <Card key={index} className="border-[#e5e7eb]">
                    <CardContent className="p-6">
                      <div className="flex items-start gap-3 mb-3">
                        <div
                          className={`w-8 h-8 rounded-full flex items-center justify-center ${
                            update.type === "warning"
                              ? "bg-[#fef2f2]"
                              : update.type === "success"
                                ? "bg-[#dcfce7]"
                                : "bg-[#dbeafe]"
                          }`}
                        >
                          <Image
                            src={update.icon || "/placeholder.svg"}
                            alt={update.title}
                            width={14}
                            height={14}
                            className="w-3.5 h-3.5"
                          />
                        </div>
                        <div className="flex-1">
                          <h3 className="font-medium text-[#111827] mb-2">{update.title}</h3>
                          <p className="text-sm text-[#4b5563] mb-3">{update.description}</p>
                          <div className="flex items-center justify-between">
                            <span className="text-xs text-[#a7b3b9]">{update.date}</span>
                            <Link href={`/user/notices?id=${index}`}>
                              <Button variant="link" className="text-[#3a6a8d] p-0 h-auto text-sm">
                                Read more
                              </Button>
                            </Link>
                          </div>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            </div>

            {/* Pagination */}
            <div className="flex items-center justify-center gap-2">
              <Button variant="ghost" size="sm" disabled={page===1 || loadingProcedures} onClick={()=> setPage(p=> Math.max(1, p-1))}>
                <ChevronLeft className="w-4 h-4 text-[#4b5563]" />
              </Button>
              <Button className="bg-[#3A6A8D] text-white w-8 h-8 p-0" disabled>{page}</Button>
              <Button variant="ghost" size="sm" disabled={!procData?.hasNext || loadingProcedures} onClick={()=> setPage(p=> p+1)}>
                <ChevronRight className="w-4 h-4 text-[#4b5563]" />
              </Button>
            </div>
          </div>
        </main>
      </div>

    </div>
  )
}
