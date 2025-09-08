"use client";
import { usePathname, useRouter } from "next/navigation";
import { UserSidebar } from "./UserSidebar";
import { signOut } from "next-auth/react";

const organizationMenuItems = [
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

];

export default function OrganizationSidebar() {
  const router = useRouter();
  const pathname = usePathname(); // 👈 gets current URL path

  const handleSettingsClick = () => {
    // Organization settings logic
  };

  const handleLogoutClick = async () => {
    await signOut({ callbackUrl: "/" }); // Call signOut and redirect to login page
  };

  const menuItemsWithHandlers = organizationMenuItems.map((item) => {
    const isActive = pathname.startsWith(`/organization/${item.label}`);
    return {
      ...item,
      active: isActive, // 👈 mark active based on URL
      onClick: () => router.push(`/organization/${item.label}`),
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