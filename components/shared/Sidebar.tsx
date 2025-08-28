"use client"

import { UserSidebar } from "./UserSidebar"
import { usePathname, useRouter } from "next/navigation"


export function Sidebar() {
  const router = useRouter();
  const pathname = usePathname();

// Define menu items with hrefs
  const menuItems = [
    { iconSrc: "/icons/dashboard.svg", iconAlt: "Dashboard", label: "Dashboard", href: "/user/home" },
    { iconSrc: "/icons/workspace.svg", iconAlt: "My workspace", label: "My workspace", href: "/user/workspace" },
    { iconSrc: "/icons/ai-chat.svg", iconAlt: "AI Chat", label: "AI Chat", href: "/user/chat" },
    { iconSrc: "/icons/discussions.svg", iconAlt: "Discussions", label: "Discussions", href: "/user/discussions" },
    { iconSrc: "/icons/official-notices.svg", iconAlt: "Official Notices", label: "Official Notices", href: "/user/notices" },
  ];

  const menuItemsWithHandlers = menuItems.map((item) => ({
    ...item,
    active: pathname === item.href,
    onClick: () => router.push(item.href),
  }));

  const handleSettingsClick = () => {
    // Example: router.push('/settings');
  };
  const handleLogoutClick = () => {
    // Example: router.push('/logout');
  };

  return (
    <UserSidebar
      menuItems={menuItemsWithHandlers}
      onSettingsClick={handleSettingsClick}
      onLogoutClick={handleLogoutClick}
    />
  );
}
