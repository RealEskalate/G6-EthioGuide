import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { BiSolidEdit } from "react-icons/bi";
import { FaEye } from "react-icons/fa";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from "@/components/ui/select";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {Trash2 } from "lucide-react";

export default function NoticeManagement() {
  const notices = [
    {
      title: "New Business Registration Requirements",
      detail:
        "Updated documentation requirements for new business applications",
      status: "Published",
      publishDate: "Jan 15, 2024",
      expiryDate: "Mar 15, 2024",
      lastUpdatedAt: "Jan 20, 2024",
    },
    {
      title: "Office Closure - Holiday Schedule",
      detail:
        "Business registration office will be closed during national holidays",
      status: "Published",
      publishDate: "Dec 15, 2023",
      expiryDate: "Dec 15, 2023",
      lastUpdatedAt: "Jan 18, 2024",
    },
    {
      title: "Digital Transformation Initiative",
      detail: "Introduction of new online services for business registration",
      status: "Published",
      publishDate: "Dec 01, 2023",
      expiryDate: "Jan 01, 2024",
      lastUpdatedAt: "Dec 15, 2023",
    },
    {
      title: "Fee Structure Updates",
      detail: "Revised fee schedule for various business registration services",
      status: "Published",
      publishDate: "Jan 10, 2024",
      expiryDate: "Jun 10, 2024",
      lastUpdatedAt: "Jan 12, 2024",
    },
  ];

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-semibold text-primary-dark">
          Notice Management
        </h1>
        <p className="text-muted-foreground text-sm text-neutral">
          Create, publish, and manage official notices for public communication
        </p>
      </div>

      {/* Filters */}
      <div className="flex flex-col md:flex-row gap-4 items-end">
        <Input placeholder="Search notices by title..." className="md:w-1/3" />
        {/* <Select>
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="All Status" />
          </SelectTrigger>
          <SelectContent className="bg-white">
            <SelectItem value="all">All Status</SelectItem>
            <SelectItem value="published">Published</SelectItem>
            <SelectItem value="draft">Draft</SelectItem>
            <SelectItem value="archived">Archived</SelectItem>
          </SelectContent>
        </Select> */}
        <label htmlFor="date" className="text-primary-dark">Publish date</label>
        <Input type="date" className="md:w-[200px]" name="date"/>
      </div>

      {/* Table */}
      <Card className="shadow-sm">
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow className="text-neutral">
                <TableHead>Notice Title</TableHead>
                <TableHead>Publish Date</TableHead>
                <TableHead>Expiry Date</TableHead>
                <TableHead>Last Updated</TableHead>
                <TableHead>Actions</TableHead>
              </TableRow>
            </TableHeader>

            <TableBody>
              {notices.map(
                ({
                  title,
                  detail,
                  status,
                  publishDate,
                  expiryDate,
                  lastUpdatedAt,
                }) => (
                  <TableRow key={title} className="hover:bg-accent">
                    <TableCell>
                      <p className="font-medium">{title}</p>
                      <p className="text-sm text-muted-foreground text-neutral">
                        {detail}
                      </p>
                    </TableCell>
                    <TableCell className="text-neutral">
                      {publishDate}
                    </TableCell>
                    <TableCell className="text-neutral">{expiryDate}</TableCell>
                    <TableCell className="text-neutral">
                      {lastUpdatedAt}
                    </TableCell>
                    <TableCell className="flex space-x-2 mt-3">
                      <FaEye className="w-4 h-4 text-primary cursor-pointer" />
                      <BiSolidEdit className="w-4 h-4 text-primary cursor-pointer" />
                      <Trash2 className="w-4 h-4 text-red-600 cursor-pointer" />
                    </TableCell>
                  </TableRow>
                )
              )}
            </TableBody>
          </Table>

          {/* Pagination */}
          <div className="flex items-center justify-between mt-4">
            <p className="text-sm text-muted-foreground text-neutral">
              Showing 1 to 4 of 12 results
            </p>
            <div className="flex space-x-2">
              <Button variant="outline" size="sm">
                Previous
              </Button>
              <Button size="sm" className="bg-primary text-white">
                1
              </Button>
              <Button variant="outline" size="sm">
                Next
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Create New Notice */}
      <div className="flex justify-end">
        <Button className="bg-primary hover:bg-primary-light text-white px-6 py-1 rounded-full flex items-center space-x-2">
          <span>+ Create New Notice</span>
        </Button>
      </div>
    </div>
  );
}
