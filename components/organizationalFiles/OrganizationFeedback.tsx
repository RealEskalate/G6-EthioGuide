'use client'
import { useState } from "react"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Pagination, PaginationContent, PaginationItem, PaginationNext, PaginationPrevious } from "@/components/ui/pagination"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Badge } from "@/components/ui/badge"
import { FaThumbsUp, FaEye, FaReply, FaCheck } from "react-icons/fa"

const feedbackData = [
  {
    id: 1,
    feedback: "The visa application process is too complicated and takes forever...",
    detail: "User feedback regarding processing delays",
    procedure: "Visa Application",
    date: "Jan 15, 2025",
    upvotes: 12,
    status: "New",
  },
  {
    id: 2,
    feedback: "Great improvement in the online portal! Much easier to navigate now.",
    detail: "Positive feedback on system updates",
    procedure: "Work Permit",
    date: "Jan 14, 2025",
    upvotes: 8,
    status: "Reviewed",
  },
  {
    id: 3,
    feedback: "Document requirements are unclear. Need better guidelines.",
    detail: "Suggestion for documentation improvement",
    procedure: "Citizenship",
    date: "Jan 13, 2025",
    upvotes: 5,
    status: "Action Taken",
  },
]

export default function OrganizationFeedback() {
  const [search, setSearch] = useState("")

  const filteredFeedback = feedbackData.filter((item) =>
    item.feedback.toLowerCase().includes(search.toLowerCase())
  )

  const getStatusBadge = (status: string) => {
    switch (status) {
      case "New":
        return <Badge className="bg-blue-100 text-primary">New</Badge>
      case "Reviewed":
        return <Badge className="bg-yellow-100 text-yellow-600">Reviewed</Badge>
      case "Action Taken":
        return <Badge className="bg-green-100 text-secondary">Action Taken</Badge>
      default:
        return <Badge variant="secondary">{status}</Badge>
    }
  }

  return (
    <div className="p-6 text-primary-dark">
      {/* Header */}
      <h1 className="text-2xl font-bold">Organization Feedback</h1>
      <p className="text-sm text-muted-foreground mb-6">
        Manage and respond to user feedback for your procedures
      </p>

      {/* Filters */}
      <div className="flex gap-4 mb-6">
        <Input
          placeholder="Search feedback..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
        <Select>
          <SelectTrigger className="w-[280px]">
            <SelectValue placeholder="All Procedures" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All Procedures</SelectItem>
            <SelectItem value="visa">Visa Application</SelectItem>
            <SelectItem value="work">Work Permit</SelectItem>
            <SelectItem value="citizenship">Citizenship</SelectItem>
          </SelectContent>
        </Select>
        <Select>
          <SelectTrigger className="w-[250px]">
            <SelectValue placeholder="All Status" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All Status</SelectItem>
            <SelectItem value="new">New</SelectItem>
            <SelectItem value="reviewed">Reviewed</SelectItem>
            <SelectItem value="action">Action Taken</SelectItem>
          </SelectContent>
        </Select>
        <Select>
          <SelectTrigger className="w-[250px]">
            <SelectValue placeholder="Sort by Date" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="newest">Newest</SelectItem>
            <SelectItem value="oldest">Oldest</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {/* Table */}
      <div className="border rounded-lg">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Feedback</TableHead>
              <TableHead>Procedure</TableHead>
              <TableHead>Date</TableHead>
              <TableHead>Upvotes</TableHead>
              <TableHead>Status</TableHead>
              <TableHead className="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {filteredFeedback.map((item) => (
              <TableRow key={item.id}>
                <TableCell>
                  <p className="font-medium">{item.feedback}</p>
                  <p className="text-sm text-muted-foreground text-neutral">{item.detail}</p>
                </TableCell>
                <TableCell className="text-secondary-dark">{item.procedure}</TableCell>
                <TableCell className="text-neutral">{item.date}</TableCell>
                <TableCell className="flex items-center gap-2 text-secondary-light">
                  <FaThumbsUp /> {item.upvotes}
                </TableCell>
                <TableCell>{getStatusBadge(item.status)}</TableCell>
                <TableCell className="flex gap-3 justify-end text-gray-600">
                  <FaEye className="cursor-pointer" />
                  <FaReply className="cursor-pointer" />
                  <FaCheck className="cursor-pointer" />
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>

      {/* Pagination */}
      <div className="flex justify-between items-center mt-4">
        <p className="text-sm text-muted-foreground">
          Showing 1 to {filteredFeedback.length} of {feedbackData.length} results
        </p>
        <Pagination>
          <PaginationContent>
            <PaginationItem>
              <PaginationPrevious href="#" />
            </PaginationItem>
            <PaginationItem>
              <Button variant="outline" size="sm">1</Button> {/*need to add a state to manage page number here*/}
            </PaginationItem>
            <PaginationItem>
              <PaginationNext href="#" />
            </PaginationItem>
          </PaginationContent>
        </Pagination>
      </div>
    </div>
  )
}
