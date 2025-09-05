"use client"

import { useState } from "react"
import { ChevronLeft, ChevronRight } from "lucide-react"
import { Button } from "@/components/ui/button"
import { CustomIcon } from "./CustomIcon"
import { cn } from "@/lib/utils"

interface MenuItem {
  iconSrc: string
  iconAlt: string
  label: string
  active?: boolean
  href?: string
  onClick?: () => void
}

interface UserSidebarProps {
  menuItems: MenuItem[]
  logoutLabel?: string
  onLogoutClick?: () => void
  className?: string
}

export function UserSidebar({
  menuItems,
  logoutLabel = "Logout",
  onLogoutClick,
  className,
}: UserSidebarProps) {
  const [collapsed, setCollapsed] = useState(false)

  return (
    <aside
      className={cn(
        "bg-white border-r border-gray-200 transition-all duration-300 ease-in-out relative flex flex-col w-20",
        collapsed ? "md:w-20" : "md:w-64",
        className,
      )}
    >
      {/* Collapse button hidden on mobile */}
      <Button
        variant="ghost"
        size="sm"
        onClick={() => setCollapsed(!collapsed)}
        className="hidden md:flex absolute -right-3 top-6 z-10 w-6 h-6 rounded-full border border-gray-200 bg-white shadow-sm hover:shadow-md transition-shadow"
        aria-label={collapsed ? "Expand sidebar" : "Collapse sidebar"}
      >
        {collapsed
          ? <ChevronRight className="w-3 h-3 text-gray-400" />
          : <ChevronLeft className="w-3 h-3 text-gray-400" />
        }
      </Button>

      <div className="flex flex-col h-full">
        <nav className="p-3 md:p-4 space-y-2 flex-1">
          {menuItems.map((item, index) => (
            <div
              key={index}
              className={cn(
                "group relative flex items-center gap-3 px-3 py-2 rounded-lg cursor-pointer transition-colors",
                item.active
                  ? "bg-gray-100 text-black shadow-sm"
                  : "text-black hover:bg-gray-100",
                collapsed && "justify-center md:justify-start",
                "md:gap-3",
              )}
              onClick={item.onClick}
              role="button"
              tabIndex={0}
              onKeyDown={(e) => {
                if (e.key === "Enter" || e.key === " ") item.onClick?.()
              }}
            >
              <CustomIcon
                src={item.iconSrc}
                alt={item.iconAlt}
                className={cn(
                  "flex-shrink-0 w-8 h-8 md:w-5 md:h-5 transition-transform",
                  item.active && "drop-shadow-sm",
                  "group-hover:scale-110"
                )}
              />
              {/* Inline label only when expanded (desktop) */}
              {!collapsed && (
                <span className="hidden md:inline text-sm font-medium truncate text-black">
                  {item.label}
                </span>
              )}
              {/* Tooltip (shown when icon-only: always on mobile, or collapsed desktop) */}
              <span
                className={cn(
                  "pointer-events-none absolute left-full top-1/2 -translate-y-1/2 ml-2",
                  "px-2 py-1 rounded-md bg-white text-black border border-gray-200 text-xs font-medium shadow-lg",
                  "opacity-0 translate-x-1 group-hover:opacity-100 group-hover:translate-x-0",
                  "transition-all duration-200 whitespace-nowrap z-50",
                  // hide tooltip only if expanded (desktop)
                  !collapsed ? "md:opacity-0 md:group-hover:opacity-0 md:hidden" : ""
                )}
              >
                {item.label}
              </span>
            </div>
          ))}
        </nav>

        <div className="p-3 md:p-4 border-t border-gray-100">
          <div
            className={cn(
              "group relative flex items-center gap-3 px-3 py-2 rounded-lg cursor-pointer transition-colors",
              "text-black hover:bg-red-50 hover:text-red-600",
              collapsed && "justify-center md:justify-start",
            )}
            onClick={onLogoutClick}
            role="button"
            tabIndex={0}
            onKeyDown={(e) => {
              if (e.key === "Enter" || e.key === " ") onLogoutClick?.()
            }}
          >
            <CustomIcon
              src="/icons/logout.svg"
              alt="Logout"
              className="w-8 h-8 md:w-5 md:h-5 group-hover:scale-110 transition-transform"
            />
            {!collapsed && (
              <span className="hidden md:inline text-sm font-medium text-black">
                {logoutLabel}
              </span>
            )}
            <span
              className={cn(
                "pointer-events-none absolute left-full top-1/2 -translate-y-1/2 ml-2",
                "px-2 py-1 rounded-md bg-white text-black border border-gray-200 text-xs font-medium shadow-lg",
                "opacity-0 translate-x-1 group-hover:opacity-100 group-hover:translate-x-0",
                "transition-all duration-200 whitespace-nowrap z-50",
                !collapsed ? "md:opacity-0 md:group-hover:opacity-0 md:hidden" : ""
              )}
            >
              {logoutLabel}
            </span>
          </div>
        </div>
      </div>
    </aside>
  )
}