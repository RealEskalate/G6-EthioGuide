"use client"
import { useRouter } from "next/navigation"
import { UserSidebar } from "./UserSidebar"

const adminMenuItems = [
  { iconSrc: "/icons/dashboard.svg", iconAlt: "Dashboard", label: "dashboard", active: true },
  { iconSrc: "/icons/official-notices.svg", iconAlt: "Notices", label: "notices", active: false },
  { iconSrc: "/icons/ai-chat.svg", iconAlt: "AI Chat", label: "AI Chat", active: false },
  { iconSrc: "/icons/discussions.svg", iconAlt: "View Feedback", label: "viewFeedback", active: false },
  { iconSrc: "/icons/user-managemnet.svg", iconAlt: "User Management", label: "userManagement", active: false },
  { iconSrc: "/icons/manage-procedure.svg", iconAlt: "Manage Procedure", label: "procedures", active: false },
]

export default function OrganizationSidebar() {
  const router = useRouter();
  const handleSettingsClick = () => {
    // org settings logic
  }

  const handleLogoutClick = () => {
    // org logout logic
  }

  const handleMenuItemClick = (label: string) => {
    router.push(`/organization/${label}`)
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
