"use client";
import { Card, CardContent } from "@/components/ui/card";
import { BiSolidEdit } from "react-icons/bi";
import { FaEye } from "react-icons/fa";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useState, useEffect } from "react";
import Notice from "@/types/notice";
import Link from "next/link";
import { Trash2 } from "lucide-react";
import Pagination from "../shared/pagination";

export default function NoticeManagement() {
  const [notices, setNotices] = useState<Notice[]>([]);
  const [page, setPage] = useState(1);
  // const [totalNotice, setTotalNotice] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  console.log(notices);
  useEffect(() => {
    const fetchsNotices = async () => {
      try {
        // const res = await fetch(`https://ethio-guide-backend.onrender.com/api/v1/notices?page=${page}&limit=${5}`);
        const res = await fetch(
          `https://ethio-guide-backend.onrender.com/api/v1/notices`
        );

        const data = await res.json();

        setNotices(data.data); // adjust to your API response
        // setTotalNotice(data.total); // if returned
        setTotalPages(Math.ceil(data.total / 5));
      } catch (err) {
        console.error(err);
      }
    };

    fetchsNotices();
  }, [page]); // <-- this will re-run whenever 'page' changes
  // const notices = [
  //   {
  //     title: "New Business Registration Requirements",
  //     body:
  //       "Updated documentation requirements for new business applications",
  //     status: "Published",
  //     createdAt: "Jan 15, 2024",
  //     expiryDate: "Mar 15, 2024",
  //     updatedAt: "Jan 20, 2024",
  //   },
  //   {
  //     title: "Office Closure - Holiday Schedule",
  //     body:
  //       "Business registration office will be closed during national holidays",
  //     status: "Published",
  //     createdAt: "Dec 15, 2023",
  //     expiryDate: "Dec 15, 2023",
  //     updatedAt: "Jan 18, 2024",
  //   },
  //   {
  //     title: "Digital Transformation Initiative",
  //     body: "Introduction of new online services for business registration",
  //     status: "Published",
  //     createdAt: "Dec 01, 2023",
  //     expiryDate: "Jan 01, 2024",
  //     updatedAt: "Dec 15, 2023",
  //   },
  //   {
  //     title: "Fee Structure Updates",
  //     body: "Revised fee schedule for various business registration services",
  //     status: "Published",
  //     createdAt: "Jan 10, 2024",
  //     expiryDate: "Jun 10, 2024",
  //     updatedAt: "Jan 12, 2024",
  //   },
  // ];

  return (
    <div className="p-6 space-y-6 w-full">
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
        <label htmlFor="date" className="text-primary-dark">
          Publish date
        </label>
        <Input type="date" className="md:w-[200px]" name="date" />
      </div>

      {/* Table */}
      <Card className="shadow-sm overflow-x-auto">
        <CardContent>
          <div className="min-w-[700px]">
            <Table>
              <TableHeader>
                <TableRow className="text-neutral">
                  <TableHead>Notice Title</TableHead>
                  <TableHead>Publish Date</TableHead>
                  {/* <TableHead>Expiry Date</TableHead> */}
                  <TableHead>Last Updated</TableHead>
                  <TableHead>Actions</TableHead>
                </TableRow>
              </TableHeader>

              <TableBody>
                {notices.map(
                  ({
                    id,
                    // orgId,
                    title,
                    body,
                    // procedures,
                    createdAt,
                    updatedAt,
                  }) => (
                    <TableRow key={id} className="hover:bg-accent">
                      <TableCell>
                        <p className="font-medium">{title}</p>
                        <p className="text-sm text-muted-foreground text-neutral">
                          {body?.length > 100
                            ? body.slice(0, 100) + "..."
                            : body}
                        </p>
                      </TableCell>
                      <TableCell className="text-neutral">
                        {createdAt}
                      </TableCell>
                      {/* <TableCell className="text-neutral">
                        {updatedAt}
                      </TableCell> */}
                      <TableCell className="text-neutral">
                        {updatedAt}
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
          </div>
        </CardContent>
      </Card>
      {/* Pagination */}
      <Pagination
        page={page}
        totalPages={totalPages}
        onPageChange={(pagenum: number) => setPage(pagenum)}
      />

      {/* <div className="flex items-center justify-between mt-4">
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
          </div> */}

      {/* Create New Notice */}
      <div className="flex justify-end">
        <Link href="/admin/notices/create">
          <Button className="bg-primary hover:bg-primary-light text-white px-6 py-1 rounded-full flex items-center space-x-2">
            <span>+ Create New Notice</span>
          </Button>
        </Link>
      </div>
    </div>
  );
}
