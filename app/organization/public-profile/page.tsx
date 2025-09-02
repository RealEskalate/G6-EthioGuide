"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { LuGlobe, LuClock2 } from "react-icons/lu";
import { FaShieldAlt } from "react-icons/fa";
import { MdGroups } from "react-icons/md";

const orgDetail = {
  name: "National Immigration Office",
  type: "Government Agency",
  est: "1995",
  description:
    "Managing immigration processes, visa applications, and citizenship services for all residents and visitors. We ensure secure, efficient, and transparent immigration procedures.",
  contact: {
    phone: "+1 (555) 123-4567",
    email: "info@immigration.gov",
  },
  address: "1234 Government Plaza, Capital City, CC 12345",
};

const quickStats = {
  activeProcedures: 24,
  publishedNotices: 156,
};

const orgMission =
  "The National Immigration Office serves as the primary government agency responsible for managing all aspects of immigration, citizenship, and border control within our nation. Established in 1995, we have been committed to maintaining the highest standards of service while ensuring national security and facilitating legitimate travel and immigration.";

const orgValues = [
  {
    title: "Security First",
    description:
      "Advanced security measures to protect our borders and citizens",
      icon:<FaShieldAlt size={25}/>,
  },
  {
    title: "Efficient Processing",
    description: "Streamlined procedures for faster application processing",
      icon:<LuClock2 size={25}/>,

  },
  {
    title: "Public Service",
    description: "Dedicated to serving immigrants and citizens with excellence",
      icon:<MdGroups size={25}/>,

  },
];

const notices = [
  {
    id: 1,
    type: "URGENT",
    date: "January 15, 2024",
    title: "New Visa Application Requirements Effective February 1st",
    description:
      "Important updates to documentation requirements for all visa categories. Please review the new guidelines before submitting applications...",
  },
  {
    id: 2,
    type: "INFO",
    date: "January 12, 2024",
    title: "Extended Office Hours During Peak Season",
    description:
      "To better serve our applicants during the busy travel season, we are extending our office hours and adding weekend appointments...",
  },
  {
    id: 3,
    type: "UPDATE",
    date: "January 10, 2024",
    title: "Online Portal Maintenance Scheduled",
    description:
      "Our online application portal will undergo scheduled maintenance on January 20th from 2:00 AM to 6:00 AM. Services will be temporarily unavailable...",
  },
];

const procedures = [
  {
    id: 1,
    title: "Tourist Visa Application",
    category: "Tourism",
    time: "5-10 days",
    description: "Standard tourist visa for visitors",
  },
  {
    id: 2,
    title: "Work Permit Application",
    category: "Employment",
    time: "15-30 days",
    description: "Employment authorization for foreign workers",
  },
  {
    id: 3,
    title: "Student Visa Application",
    category: "Education",
    time: "10-20 days",
    description: "Education visa for international students",
  },
];

export default function ImmigrationOfficePage() {
  return (
    <div className="p-6 space-y-8">
      {/* Header Section */}
      <Card className="p-6 bg-gradient-to-r from-primary-dark to-primary-light text-white">
        <div className="flex">
          <div className="mt-2 flex">
            <div className="bg-primary w-20 h-20 flex items-center justify-center rounded">
              <div className="bg-white p-1 rounded w-8 h-8 flex items-center justify-center">
                <LuGlobe className="text-primary-dark size-5" />
              </div>
            </div>
          </div>
          <div>
            <CardHeader>
              <CardTitle className="text-2xl">{orgDetail.name}</CardTitle>
              <p className="text-sm">
                <span className="p-1 bg-primary-light rounded-2xl px-2">
                  {orgDetail.type}
                </span>{" "}
                Est. {orgDetail.est}
              </p>
            </CardHeader>
            <CardContent>
              <p className="mb-4">{orgDetail.description}</p>
              <div className="grid grid-cols-2 gap-4 text-sm">
                <div>
                  <p className="font-semibold">Contact</p>
                  <p>{orgDetail.contact.phone}</p>
                  <p>{orgDetail.contact.email}</p>
                </div>
                <div>
                  <p className="font-semibold">Address</p>
                  <p>{orgDetail.address}</p>
                </div>
              </div>
            </CardContent>
          </div>
        </div>
      </Card>

      <div className="text-primary-dark">
        {/* About Section */}
        <section>
          <h2 className="text-xl font-bold mb-2">About Our Organization</h2>
          <p className="text-sm text-muted-foreground mb-6">{orgMission}</p>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {orgValues.map((value) => (
              <Card key={value.title} className=" bg-accent">
                <div className="flex justify-center text-white">
                  <span className="bg-primary p-5 rounded-full">
                    {value.icon}
                  </span>
                </div>
                <CardContent className="p-6">
                  <h3 className="font-semibold text-center">{value.title}</h3>
                  <p className="text-sm text-muted-foreground text-center text-neutral">
                    {value.description}
                  </p>
                </CardContent>
              </Card>
            ))}
          </div>
        </section>

        <Separator />

        {/* Recent Feedbacks */}
        {/* <section>
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-bold">Recent Feedbacks</h2>
            <Button variant="link">View All</Button>
          </div>
          <div className="space-y-4">
            {notices.map((notice) => (
              <Card key={notice.id}>
                <CardContent className="p-6">
                  <div className="flex items-center gap-2 mb-2">
                    <Badge className="bg-accent">{notice.type}</Badge>
                    <span className="text-xs text-muted-foreground">
                      {notice.date}
                    </span>
                  </div>
                  <h3 className="font-semibold">{notice.title}</h3>
                  <p className="text-sm text-muted-foreground mb-2">
                    {notice.description}
                  </p>
                  <Button variant="link" size="sm">
                    Read More
                  </Button>
                </CardContent>
              </Card>
            ))}
          </div>
        </section> */}

        <Separator />

        {/* Popular Procedures */}
        {/* <section>
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-bold">Popular Procedures</h2>
            <Button variant="link">View All</Button>
          </div>
          <div className="grid gap-4">
            {procedures.map((proc) => (
              <Card key={proc.id}>
                <CardContent className="p-6 flex items-center justify-between">
                  <div>
                    <h3 className="font-semibold">{proc.title}</h3>
                    <p className="text-sm text-muted-foreground">
                      {proc.description}
                    </p>
                    <Badge className="mt-1 bg-accent">{proc.category}</Badge>
                    <p className="text-xs text-muted-foreground">
                      Processing Time: {proc.time}
                    </p>
                  </div>
                  <Button className="bg-primary-light text-white">Start Application</Button>
                </CardContent>
              </Card>
            ))}
          </div>
        </section> */}
      </div>
    </div>
  );
}
