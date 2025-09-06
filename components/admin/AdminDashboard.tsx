"use client";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { IoMegaphoneOutline } from "react-icons/io5";
// import { MdOutlineFeedback } from "react-icons/md";
// import { FaUsers } from "react-icons/fa6";
import { Plus, Megaphone, FileText, MessageSquare } from "lucide-react";
import { useRouter } from "next/navigation";

export default function AdminDashboard({
  totalProcedures,
  totalNotices,
}: {
  totalProcedures: number;
  totalNotices: number;
}) {
  const route = useRouter();
  const stats = [
    {
      data: totalProcedures,
      description: "Procedures Managed",
      icon: (
        <div className="bg-gray-100 p-3 rounded-2xl ">
          <FileText className="w-6 h-6  text-[#3A6A8D] mb-2" />
        </div>
      ),
    },
    {
      data: totalNotices,
      description: "Active Notices",
      icon: (
        <div className="bg-gray-100 p-3 rounded-2xl ">
          <IoMegaphoneOutline className="w-6 h-6 text-[#5E9C8D] mb-2" />
        </div>
      ),
    },
    // {
    //   data: 23,
    //   description: "Pending Feedback",
    //   icon: (
    //     <div className="bg-gray-100 p-3 rounded-2xl ">
    //       <MdOutlineFeedback className="w-6 h-6 text-[#1C3B2E] mb-2" />
    //     </div>
    //   ),
    // },
    // {
    //   data: 1284,
    //   description: "User Interactions",
    //   icon: (
    //     <div className="bg-gray-100 p-3 rounded-2xl ">
    //       <FaUsers className="w-6 h-6 text-[#1C3B2E] mb-2" />
    //     </div>
    //   ),
    // },
  ];

  return (
    <div className="p-6 space-y-6 w-full">
      {/* Welcome Section */}
      <div className="bg-gradient-to-r from-primary-dark to-primary-light text-white p-6 rounded-2xl shadow py-10">
        <h1 className="text-2xl font-semibold">Welcome back, Admin</h1>
      </div>

      {/* Actions */}
      <div className="flex space-x-4 text-white">
        <Button
          className="flex items-center space-x-2 bg-[#3A6A8D] hover:bg-[#5C87A3]"
          onClick={() => route.push("/admin/addNewProcedures")}
        >
          <Plus className="w-4 h-4" />
          <span>Add New Procedure</span>
        </Button>
        <Link href="/admin/notices/create">
          <Button className="flex items-center space-x-2 bg-[#5E9C8D] hover:bg-[#7FB4A6]">
            <Megaphone className="w-4 h-4" />
            <span>Create Notice</span>
          </Button>
        </Link>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        {stats.map(({ data, description, icon }) => (
          <Card
            key={data}
            className="shadow-sm border-gray-50 hover:scale-105 transition-transform duration-300 cursor-pointer"
          >
            <CardContent className="flex items-center justify-between p-4 ">
              {icon}

              <p className="text-2xl font-bold">{data}</p>
            </CardContent>
            <CardContent>
              <p className="text-sm text-muted-foreground">{description}</p>
            </CardContent>
          </Card>
        ))}
      </div>
       {/* Recent Activity and Quick Overview */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-10">
        {/* Recent Activity */}
        <Card className="md:col-span-2 shadow-sm border-gray-50">
          <CardHeader>
            <CardTitle>Recent Activity</CardTitle>
          </CardHeader>
          <CardContent>
            <ul className="space-y-4">
              <li className="flex items-center space-x-3">
                <div className="bg-gray-200 p-2 rounded-full">
                  <Plus className="w-5 h-5 text-blue-600" />
                </div>
                <div>
                  <p className="text-sm font-medium">
                    New procedure added: Passport Renewal
                  </p>
                  <p className="text-xs text-muted-foreground">2 hours ago</p>
                </div>
              </li>
              <li className="flex items-center space-x-3">
                <div className="bg-gray-200 p-2 rounded-full">
                  <Megaphone className="w-5 h-5 text-[#5E9C8D] " />
                </div>
                <div>
                  <p className="text-sm font-medium">
                    Notice published: Office closed for holiday
                  </p>
                  <p className="text-xs text-muted-foreground">1 day ago</p>
                </div>
              </li>
              <li className="flex items-center space-x-3">
                <div className="bg-gray-200 p-2 rounded-full">
                  <MessageSquare className="w-5 h-5 text-[#3A6A8D]" />
                </div>
                <div>
                  <p className="text-sm font-medium">
                    User feedback received on Driver&#39;s License procedure
                  </p>
                  <p className="text-xs text-muted-foreground">3 days ago</p>
                </div>
              </li>
              <li className="flex items-center space-x-3">
                <div className="bg-gray-200 p-2 rounded-full">
                  <FileText className="w-5 h-5 text-gray-600" />
                </div>
                <div>
                  <p className="text-sm font-medium">
                    Procedure updated: Visa Application Process
                  </p>
                  <p className="text-xs text-muted-foreground">5 days ago</p>
                </div>
              </li>
            </ul>
            <Button variant="link" className="mt-4">
              View all activities
            </Button>
          </CardContent>
        </Card>

        {/* Quick Overview */}
        {/* <Card className="shadow-sm border-gray-50">
          <CardHeader>
            <CardTitle>Quick Overview</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              <p className="text-sm">
                <span className="text-blue-600">●</span> Active Procedures: 89
              </p>
              <p className="text-sm">
                <span className="text-[#5E9C8D]">●</span> Draft Procedures: 12
              </p>
              <p className="text-sm">
                <span className="text-gray-600">●</span> Archived: 26
              </p>
            </div>
            <div className="mt-4">
              <hr />
              <h4 className="text-sm font-medium mb-2 mt-5">Top Procedures</h4>
              <ul className="space-y-1 text-sm">
                <li>Passport Application - 342 views</li>
                <li>Visa Renewal - 198 views</li>
                <li>Work Permit - 156 views</li>
              </ul>
            </div>
          </CardContent>
        </Card> */}
      </div>
    </div>
  );
}