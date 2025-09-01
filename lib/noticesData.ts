export interface Notice {
  id: number;
  title: string;
  status: "Active" | "Upcoming" | "Expired";
  statusColor: string;
  description: string;
  published: string;
  department: string;
  organization: string;
  likes: number;
  author?: string;
  keyChanges?: string[];
  securityUpdates?: string[];
  attachments?: { name: string; size: string }[];
}

export const notices: Notice[] = [
  {
    id: 1,
    title: "New Employee Onboarding Process Updates",
    status: "Active",
    statusColor: "bg-green-100 text-green-800",
    description:
      "Updated guidelines for the employee onboarding process, including new documentation requirements and digital workflow procedures.",
    published: "Dec 15, 2024",
    department: "HR Department",
    organization: "TechCorp Solutions",
    likes: 12,
    author: "HR Department",
    keyChanges: [
      "All remote work arrangements must be pre-approved through the new digital workflow system",
      "Mandatory use of company VPN for all work-related activities",
      "Weekly check-ins with direct supervisors are now required for remote workers",
      "New cybersecurity training completion required within 30 days",
    ],
    securityUpdates: [
      "All employees must enable two-factor authentication on company accounts",
      "Complete the updated security awareness training by January 31st, 2024",
    ],
    attachments: [
      { name: "Remote Work Policy 2024.pdf", size: "2.4 MB" },
      { name: "Security Guidelines Checklist.docx", size: "1.8 MB" },
    ],
  },
  {
    id: 2,
    title: "System Maintenance Schedule",
    status: "Upcoming",
    statusColor: "bg-yellow-100 text-yellow-800",
    description:
      "Scheduled system maintenance on December 25th from 2:00 AM to 6:00 AM. All services will be temporarily unavailable during this period.",
    published: "Dec 14, 2024",
    department: "IT Operations",
    organization: "TechCorp Solutions",
    likes: 8,
    author: "IT Operations",
  },
  {
    id: 3,
    title: "Holiday Policy Updates",
    status: "Expired",
    statusColor: "bg-gray-100 text-gray-500",
    description:
      "Updates to the company holiday policy for 2024, including new national holidays and revised vacation request procedures.",
    published: "Nov 30, 2024",
    department: "HR Department",
    organization: "TechCorp Solutions",
    likes: 25,
    author: "HR Department",
  },
];
