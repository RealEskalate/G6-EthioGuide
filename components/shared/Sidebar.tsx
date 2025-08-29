"use client"

import { UserSidebar } from "./UserSidebar"

const defaultMenuItems = [
  { iconSrc: "/icons/dashboard.svg", iconAlt: "Dashboard", label: "Dashboard", active: true, href: "/user/home" },
  { iconSrc: "/icons/workspace.svg", iconAlt: "My workspace", label: "My workspace", active: false },
  { iconSrc: "/icons/ai-chat.svg", iconAlt: "AI Chat", label: "AI Chat", active: false },
  { iconSrc: "/icons/discussions.svg", iconAlt: "Discussions", label: "Discussions", active: false },
  { iconSrc: "/icons/official-notices.svg", iconAlt: "Official Notices", label: "Official Notices", active: false },
]

export function Sidebar() {
  const handleSettingsClick = () => {
    console.log(" Settings clicked")
    // Add settings navigation logic here
  }

  const handleLogoutClick = () => {
    console.log(" Logout clicked")
    // Add logout logic here
  }

  const handleMenuItemClick = (label: string) => {
    console.log(` ${label} clicked`)
    // Add navigation logic here
  }

  const menuItemsWithHandlers = defaultMenuItems.map((item) => ({
    ...item,
    onClick: () => handleMenuItemClick(item.label),
  }))

  return (
    <UserSidebar
      menuItems={menuItemsWithHandlers}
      onSettingsClick={handleSettingsClick}
      onLogoutClick={handleLogoutClick}
    />
  )
}
