"use client"
import { Search, Bookmark, Star, ChevronLeft, ChevronRight } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent } from "@/components/ui/card"
import Image from "next/image"
import Link from "next/link"
import { useState } from "react"

export default function ProceduresPage() {

  const procedures = [
    {
      id: 1,
      name: "Tourist Visa Application",
      description: "Apply for short-term tourist visa for leisure travel",
      category: "Visa Applications",
      categoryColor: "bg-blue-100 text-blue-700",
      processingTime: "5-7 business days",
      lastUpdated: "Dec 15, 2024",
      icon: "/icons/business.svg",
      iconBg: "bg-blue-50",
      bookmarked: false,
    },
    {
      id: 2,
      name: "Work Permit Application",
      description: "Apply for employment authorization and work permit",
      category: "Work Permits",
      categoryColor: "bg-gray-100 text-gray-700",
      processingTime: "10-15 business days",
      lastUpdated: "Dec 10, 2024",
      icon: "/icons/workspace.svg",
      iconBg: "bg-gray-50",
      bookmarked: true,
    },
    {
      id: 3,
      name: "Citizenship Application",
      description: "Apply for naturalization and citizenship status",
      category: "Citizenship",
      categoryColor: "bg-orange-100 text-orange-700",
      processingTime: "6-8 months",
      lastUpdated: "Dec 8, 2024",
      icon: "/icons/citizenship.svg",
      iconBg: "bg-yellow-50",
      bookmarked: false,
    },
    {
      id: 4,
      name: "Family Reunification Visa",
      description: "Bring family members to join you in the country",
      category: "Family Reunification",
      categoryColor: "bg-red-100 text-red-700",
      processingTime: "3-4 months",
      lastUpdated: "Dec 5, 2024",
      icon: "/icons/family-love.svg",
      iconBg: "bg-green-50",
      bookmarked: false,
    },
    {
      id: 5,
      name: "Student Visa Application",
      description: "Apply for student visa for educational purposes",
      category: "Student Permits",
      categoryColor: "bg-blue-100 text-blue-700",
      processingTime: "4-6 weeks",
      lastUpdated: "Dec 1, 2024",
      icon: "/icons/student-graduate.svg",
      iconBg: "bg-blue-50",
      bookmarked: false,
    },
    {
      id: 6,
      name: "Residence Permit Renewal",
      description: "Renew your existing residence permit",
      category: "Renewals",
      categoryColor: "bg-purple-100 text-purple-700",
      processingTime: "2-3 weeks",
      lastUpdated: "Nov 28, 2024",
      icon: "/icons/passport.svg",
      iconBg: "bg-gray-50",
      bookmarked: true,
    },
  ]

  const [search, setSearch] = useState("");
  const [category, setCategory] = useState("all");

  const filteredProcedures = procedures.filter((procedure) => {
    const matchesSearch = procedure.name.toLowerCase().includes(search.toLowerCase()) || procedure.description.toLowerCase().includes(search.toLowerCase());
    const matchesCategory = category === "all" ||
      (category === "visa" && procedure.category === "Visa Applications") ||
      (category === "work" && procedure.category === "Work Permits") ||
      (category === "citizenship" && procedure.category === "Citizenship");
    return matchesSearch && matchesCategory;
  });

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
                <Select value={category} onValueChange={setCategory}>
                  <SelectTrigger className="w-48 bg-white border-[#d1d5db] hover:!bg-gray-100 font-medium">
                    <SelectValue placeholder="All Categories" />
                  </SelectTrigger>
                  <SelectContent className="bg-white">
                    <SelectItem value="all" className="hover:!bg-gray-100 font-medium border-0">All Categories</SelectItem>
                    <SelectItem value="visa" className="hover:!bg-gray-100 font-medium border-0">Visa Applications</SelectItem>
                    <SelectItem value="work" className="hover:!bg-gray-100 font-medium border-0">Work Permits</SelectItem>
                    <SelectItem value="citizenship" className="hover:!bg-gray-100 font-medium border-0">Citizenship</SelectItem>
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
                    {filteredProcedures.map((procedure, index) => (
                      <tr
                        key={procedure.id}
                        className={(index !== procedures.length - 1 ? "border-b border-[#f3f4f6] " : "") + " hover:bg-gray-100 transition-colors"}
                      >
                        <td className="py-4 px-6">
                          <div className="flex items-center gap-3">
                            <div
                              className={`w-10 h-10 ${procedure.iconBg} rounded-lg flex items-center justify-center`}
                            >
                              <Image
                                src={procedure.icon || "/placeholder.svg"}
                                alt={procedure.name}
                                width={20}
                                height={20}
                                className="w-5 h-5"
                              />
                            </div>
                            <div>
                        {/* Footer is now part of the shared layout */}
                              <div className="text-sm text-[#4b5563]">{procedure.description}</div>
                            </div>
                          </div>
                        </td>
                        <td className="py-4 px-6">
                          <Badge className={`${procedure.categoryColor} border-0`}>{procedure.category}</Badge>
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
              <Button variant="ghost" size="sm">
                <ChevronLeft className="w-4 h-4 text-[#4b5563] hover:!bg-gray-200" />
              </Button>
              <Button className="bg-[#3a6a8d] text-white w-8 h-8 p-0">1</Button>
              <Button variant="ghost" className="w-8 h-8 p-0 hover:bg-gray-200">
                2
              </Button>
              <Button variant="ghost" className="w-8 h-8 p-0 hover:bg-gray-200">
                3
              </Button>
              <span className="text-[#4b5563] mx-2">...</span>
              <Button variant="ghost" className="w-8 h-8 p-0 hover:bg-gray-200">
                12
              </Button>
              <Button variant="ghost" size="sm" className="hover:bg-gray-200">
                <ChevronRight className="w-4 h-4 text-[#4b5563]" />
              </Button>
            </div>
          </div>
        </main>
      </div>

    </div>
  )
}
