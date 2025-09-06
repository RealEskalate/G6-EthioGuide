"use client"

import { useEffect, useState } from "react"
import { AlertTriangle, Info, MessageCircle, Calendar, FileX } from "lucide-react"
import { Button } from "@/components/ui/button"
import { DropdownMenu, DropdownMenuContent, DropdownMenuTrigger } from "@/components/ui/dropdown-menu"
import Image from "next/image"

const notifications = [
  {
    id: 1,
    type: "warning",
    icon: AlertTriangle,
    title: "Document Expiring Soon",
    description: "Your security certificate expires in 3 days",
    time: "2 hours ago",
    unread: true,
    bgColor: "bg-orange-100",
    iconColor: "text-orange-600",
  },
  {
    id: 2,
    type: "info",
    icon: Info,
    title: "New Notice Posted",
    description: "Updated safety guidelines have been published",
    time: "4 hours ago",
    unread: true,
    bgColor: "bg-blue-100",
    iconColor: "text-blue-600",
  },
  {
    id: 3,
    type: "success",
    icon: MessageCircle,
    title: "Response to Your Feedback",
    description: "Admin replied to your suggestion about workflow",
    time: "6 hours ago",
    unread: true,
    bgColor: "bg-green-100",
    iconColor: "text-green-600",
  },
  {
    id: 4,
    type: "reminder",
    icon: Calendar,
    title: "Procedure Reminder",
    description: "Monthly compliance review is due tomorrow",
    time: "8 hours ago",
    unread: true,
    bgColor: "bg-purple-100",
    iconColor: "text-purple-600",
  },
  {
    id: 5,
    type: "error",
    icon: FileX,
    title: "Contract Expiring",
    description: "Vendor agreement expires next week",
    time: "1 day ago",
    unread: true,
    bgColor: "bg-red-100",
    iconColor: "text-red-600",
  },
]

export function NotificationDropdown() {
  const [unreadCount] = useState(12)
  const [mounted, setMounted] = useState(false)
  useEffect(() => setMounted(true), [])

  if (!mounted) {
    return (
      <Button variant="ghost" size="sm" className="relative p-2 rounded-full" aria-hidden>
        <Image src="/icons/notifications.svg" alt="Notifications" width={20} height={20} className="w-5 h-5" />
        {unreadCount > 0 && (
          <div className="absolute -top-1 -right-1 w-3 h-3 bg-red-500 rounded-full"></div>
        )}
      </Button>
    )
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" size="sm" className="relative p-2 rounded-full hover:bg-gray-100">
          <Image src="/icons/notifications.svg" alt="Notifications" width={20} height={20} className="w-5 h-5" />
          {unreadCount > 0 && (
            <div className="absolute -top-1 -right-1 w-3 h-3 bg-red-500 rounded-full animate-pulse"></div>
          )}
        </Button>
      </DropdownMenuTrigger>
  <DropdownMenuContent align="end" className="w-80 p-0 max-h-96 overflow-hidden bg-white bg-opacity-100">
        <div className="p-4 border-b border-gray-100">
          <div className="flex items-center justify-between">
            <h3 className="font-semibold text-gray-900">Notifications</h3>
            <span className="text-sm text-blue-600 font-medium">{unreadCount} new</span>
          </div>
        </div>

        <div className="overflow-y-auto" style={{ maxHeight: "calc(24rem - 8rem)" }}>
          {notifications.map((notification) => {
            const IconComponent = notification.icon
            return (
              <div
                key={notification.id}
                className="p-4 hover:bg-gray-50 border-b border-gray-50 last:border-b-0 transition-colors"
              >
                <div className="flex gap-3">
                  <div
                    className={`w-8 h-8 ${notification.bgColor} rounded-full flex items-center justify-center flex-shrink-0`}
                  >
                    <IconComponent className={`w-4 h-4 ${notification.iconColor}`} />
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-start justify-between gap-2">
                      <p className="font-medium text-gray-900 text-sm">{notification.title}</p>
                      {notification.unread && (
                        <div className="w-2 h-2 bg-blue-500 rounded-full flex-shrink-0 mt-1"></div>
                      )}
                    </div>
                    <p className="text-sm text-gray-600 mt-1">{notification.description}</p>
                    <p className="text-xs text-gray-500 mt-2">{notification.time}</p>
                  </div>
                </div>
              </div>
            )
          })}
        </div>

        <div className="p-4 border-t border-gray-100 flex items-center justify-between">
          <Button variant="ghost" size="sm" className="text-gray-600 hover:text-gray-800 p-0">
            Mark all as read
          </Button>
          <Button variant="ghost" size="sm" className="text-blue-600 hover:text-blue-800 p-0">
            View all notifications
          </Button>
        </div>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
