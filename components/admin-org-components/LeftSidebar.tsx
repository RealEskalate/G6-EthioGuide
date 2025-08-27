"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"
import { IoHomeOutline, IoMegaphoneOutline } from "react-icons/io5"
import { MdOutlineFeedback } from "react-icons/md"
import { PiUsersThree } from "react-icons/pi"
import { BiChat } from "react-icons/bi"
import { RiFolderSettingsLine } from "react-icons/ri"

export default function LeftSidebar() {
  const pathname = usePathname()

  const links = [
    { href: "/dashboard", label: "Dashboard", icon: <IoHomeOutline size={18} /> },
    { href: "/admin/notices", label: "Notices", icon: <IoMegaphoneOutline size={18} /> },
    { href: "/feedback", label: "View Feedback", icon: <MdOutlineFeedback size={18} /> },
    { href: "/users", label: "User Management", icon: <PiUsersThree size={18} /> },
    { href: "/procedures", label: "Manage Procedures", icon: <RiFolderSettingsLine size={18} /> },
    { href: "/chats", label: "Chats", icon: <BiChat size={18} /> },
  ]

  return (
    <aside className="text-[#A7B3B9] fixed left-0 top-16 h-[calc(100vh)] w-64 border-l bg-background shadow-lg flex flex-col">
      <nav className="flex-1 p-4 space-y-2">
        {links.map(({ href, label, icon }) => {
          const isActive = pathname === href
          return (
            <Link
              key={href}
              href={href}
              className={`flex items-center gap-3 rounded-md px-3 py-2 transition 
                ${isActive ? "bg- text-accent-foreground" : "hover:bg-[rgba(58,106,141,0.1)] hover:text-accent-foreground"}`}
            >
              {icon} {label}
            </Link>
          )
        })}
      </nav>
    </aside>
  )
}
