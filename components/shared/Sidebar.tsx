"use client";

import { useEffect } from "react";
import { UserSidebar } from "./UserSidebar";
import { usePathname, useRouter } from "next/navigation";
import { signOut, getSession } from "next-auth/react";

export function Sidebar() {
  const router = useRouter();
  const pathname = usePathname();

  // persist access token so RTK slices can read it
  useEffect(() => {
    let alive = true;
    (async () => {
      try {
        const session = await getSession();
        if (!alive || !session) return;
        const token =
          (session as unknown as { accessToken?: string })?.accessToken ||
          (session as unknown as { user?: { accessToken?: string } })?.user?.accessToken;
        if (token) {
          localStorage.setItem("accessToken", token);
          document.cookie = `accessToken=${encodeURIComponent(token)}; path=/; max-age=${60 * 60 * 12}`;
        }
      } catch {
        // ignore
      }
    })();
    return () => {
      alive = false;
    };
  }, []);

  const menuItems = [
    {
      iconSrc: "/icons/dashboard.svg",
      iconAlt: "Dashboard",
      label: "Dashboard",
      href: "/user/home",
    },
    {
      iconSrc: "/icons/workspace.svg",
      iconAlt: "Workspace",
      label: "Workspace",
      href: "/user/workspace",
    },
    {
      iconSrc: "/icons/ai-chat.svg",
      iconAlt: "AI Chat",
      label: "AI Chat",
      href: "/user/chat",
    },
    {
      iconSrc: "/icons/discussions.svg",
      iconAlt: "Discussions",
      label: "Discussions",
      href: "/user/discussions",
    },
    {
      iconSrc: "/icons/official-notices.svg",
      iconAlt: "Official Notices",
      label: "Official Notices",
      href: "/user/notices",
    },
  ];

  const menuItemsWithHandlers = menuItems.map((item) => ({
    ...item,
    active: pathname === item.href,
    onClick: () => router.push(item.href),
  }));

  const handleSettingsClick = () => {
    // Example: router.push('/settings');
  };
  const handleLogoutClick = async () => {
    await signOut({ callbackUrl: "/" }); // Call signOut and redirect to login page
  };

  return (
    <UserSidebar
      menuItems={menuItemsWithHandlers}
      onSettingsClick={handleSettingsClick}
      onLogoutClick={handleLogoutClick}
    />
  );
}
