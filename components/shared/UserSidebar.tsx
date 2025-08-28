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
  settingsLabel?: string
  logoutLabel?: string
  onSettingsClick?: () => void
  onLogoutClick?: () => void
  className?: string
}

export function UserSidebar({
  menuItems,
  settingsLabel = "Settings",
  logoutLabel = "Logout",
  onSettingsClick,
  onLogoutClick,
  className,
}: UserSidebarProps) {
  const [collapsed, setCollapsed] = useState(false)

  return (
    <aside
      className={cn(
        "bg-white border-r border-gray-200 transition-all duration-300 ease-in-out relative",
        collapsed ? "w-20" : "w-64",
        className,
      )}
    >
      <Button
        variant="ghost"
        size="sm"
        onClick={() => setCollapsed(!collapsed)}
        className="absolute -right-3 top-6 z-10 w-6 h-6 rounded-full border border-gray-200 bg-white shadow-sm hover:shadow-md transition-shadow"
        aria-label={collapsed ? "Expand sidebar" : "Collapse sidebar"}
      >
        {collapsed 
          ? <ChevronRight className="w-3 h-3 text-gray-300 hover:text-gray-100" /> 
          : <ChevronLeft className="w-3 h-3 text-gray-300 hover:text-gray-100" />
        }
      </Button>

      <div className="flex flex-col h-full">
        <nav className="p-4 space-y-2 flex-1">
          {menuItems.map((item, index) => (
            <div
              key={index}
              className={cn(
                "flex items-center gap-3 px-3 py-2 rounded-lg cursor-pointer transition-all duration-200",
                item.active
                  ? "text-gray-800 bg-gray-100 shadow-sm"
                  : "text-gray-600 hover:bg-gray-200 hover:text-gray-900 font-medium",
                collapsed ? "justify-center" : "",
              )}
              onClick={item.onClick}
              role="button"
              tabIndex={0}
              onKeyDown={(e) => {
                if (e.key === "Enter" || e.key === " ") {
                  item.onClick?.()
                }
              }}
            >
              <CustomIcon
                src={item.iconSrc}
                alt={item.iconAlt}
                className={cn("flex-shrink-0", collapsed ? "w-8 h-8" : "w-5 h-5")}
              />
              {!collapsed && (
                <span className={cn("transition-opacity duration-200", item.active ? "font-medium" : "")}>
                  {item.iconAlt}
                </span>
              )}
            </div>
          ))}
        </nav>

        <div className="p-4 space-y-2 border-t border-gray-100">
          <div
            className={cn(
              "flex items-center gap-3 px-3 py-2 text-gray-600 hover:bg-gray-200 hover:text-gray-900 font-medium rounded-lg cursor-pointer transition-all duration-200",
              collapsed ? "justify-center" : "",
            )}
            onClick={onSettingsClick}
            role="button"
            tabIndex={0}
            onKeyDown={(e) => {
              if (e.key === "Enter" || e.key === " ") {
                onSettingsClick?.()
              }
            }}
          >
            <CustomIcon
              src="/icons/settings.svg"
              alt="Settings"
              className={cn("flex-shrink-0", collapsed ? "w-8 h-8" : "w-5 h-5")}
            />
            {!collapsed && <span>{settingsLabel}</span>}
          </div>

          <div
            className={cn(
              "flex items-center gap-3 px-3 py-2 text-gray-600 hover:bg-red-100 hover:text-red-700 font-medium rounded-lg cursor-pointer transition-all duration-200",
              collapsed ? "justify-center" : "",
            )}
            onClick={onLogoutClick}
            role="button"
            tabIndex={0}
            onKeyDown={(e) => {
              if (e.key === "Enter" || e.key === " ") {
                onLogoutClick?.()
              }
            }}
          >
            <CustomIcon
              src="/icons/logout.svg"
              alt="Logout"
              className={cn("flex-shrink-0", collapsed ? "w-8 h-8" : "w-5 h-5")}
            />
            {!collapsed && <span>{logoutLabel}</span>}
          </div>
        </div>
      </div>
    </aside>
  )
}
