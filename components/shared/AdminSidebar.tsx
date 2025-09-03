"use client";
import { usePathname, useRouter } from "next/navigation";
import { UserSidebar } from "./UserSidebar";

const adminMenuItems = [
  {
    iconSrc: "/icons/dashboard.svg",
    iconAlt: "Dashboard",
    label: "dashboard",
  },
  {
    iconSrc: "/icons/official-notices.svg",
    iconAlt: "Notices",
    label: "notices",
  },
  {
    iconSrc: "/icons/discussions.svg",
    iconAlt: "View Feedbacks",
    label: "feedback",
  },
  {
    iconSrc: "/icons/manage-procedure.svg",
    iconAlt: "Manage Procedures",
    label: "procedures",
  },
  {
    iconSrc: "/icons/building.svg",
    iconAlt: "Create Organization",
    label: "create",
  },
];

export function AdminSidebar() {
  const router = useRouter();
  const pathname = usePathname(); // 👈 gets current URL path

  const handleSettingsClick = () => {
    // Admin settings logic
  };

  const handleLogoutClick = () => {
    // Admin logout logic
  };

  const menuItemsWithHandlers = adminMenuItems.map((item) => {
    const isActive = pathname.startsWith(`/admin/${item.label}`);
    return {
      ...item,
      active: isActive, // 👈 mark active based on URL
      onClick: () => router.push(`/admin/${item.label}`),
    };
  });

  return (
    <UserSidebar
      menuItems={menuItemsWithHandlers}
      onSettingsClick={handleSettingsClick}
      onLogoutClick={handleLogoutClick}
      settingsLabel="Settings"
      logoutLabel="Sign Out"
    />
  );
}