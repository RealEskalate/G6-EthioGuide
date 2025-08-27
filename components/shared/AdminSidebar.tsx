"use client"

import { UserSidebar } from "./UserSidebar"

const adminMenuItems = [
  { iconSrc: "/icons/dashboard.svg", iconAlt: "Dashboard", label: "Dashboard", active: true },
  { iconSrc: "/icons/official-notices.svg", iconAlt: "Notices", label: "Notices", active: false },
  { iconSrc: "/icons/ai-chat.svg", iconAlt: "AI Chat", label: "AI Chat", active: false },
  { iconSrc: "/icons/discussions.svg", iconAlt: "View Feedback", label: "View Feedback", active: false },
  { iconSrc: "/icons/user-managemnet.svg", iconAlt: "User Management", label: "User Management", active: false },
  { iconSrc: "/icons/manage-procedure.svg", iconAlt: "Manage Procedure", label: "Manage Procedure", active: false },
]

export function AdminSidebar() {
  const handleSettingsClick = () => {
    // Admin settings logic
  }

  const handleLogoutClick = () => {
    // Admin logout logic
  }

  const handleMenuItemClick = (label: string) => {
    // Admin navigation logic
  }

  const menuItemsWithHandlers = adminMenuItems.map((item) => ({
    ...item,
    onClick: () => handleMenuItemClick(item.label),
  }))

  return (
    <UserSidebar
      menuItems={menuItemsWithHandlers}
      onSettingsClick={handleSettingsClick}
      onLogoutClick={handleLogoutClick}
      settingsLabel="Settings"
      logoutLabel="Sign Out"
    />
  )
}
