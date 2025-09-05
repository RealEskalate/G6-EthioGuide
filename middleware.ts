// import { getToken } from "next-auth/jwt";
// import { NextRequest, NextResponse } from "next/server";

// export async function middleware(request: NextRequest) {
//   const { pathname, origin } = request.nextUrl;

//   // Skip middleware for NextAuth routes and /auth/login
//   if (pathname.startsWith("/api/auth") || pathname === "/auth/login") {
//     console.log("Skipping middleware for:", { pathname });
//     return NextResponse.next();
//   }

//   // Define a type for the token to include user and role
//   type TokenWithRole = {
//     user?: {
//       role?: string;
//       [key: string]: unknown;
//     };
//     [key: string]: unknown;
//   };

//   // Get the session token
//   const token = (await getToken({
//     req: request,
//     secret: process.env.NEXTAUTH_SECRET,
//   })) as TokenWithRole;

//   // Define role-based redirect rules
//   const redirects: { [key: string]: string } = {
//     admin: "/admin/dashboard",
//     user: "/user/home",
//     org: "/organization/dashboard",
//   };

//   // Define protected routes by role
//   const protectedRoutes: { [key: string]: string[] } = {
//     admin: ["/admin", "/admin/*"],
//     user: ["/user", "/user/*"],
//     org: ["/organization", "/organization/*"],
//   };

//   // Check if the current path is a protected route (excluding '/')
//   // We're making '/' accessible to unauthenticated users here initially
//   const isProtectedRoute = Object.entries(protectedRoutes).some(([r, paths]) =>
//     paths.some((path) => {
//       if (path.endsWith("/*")) {
//         const basePath = path.slice(0, -2);
//         return pathname.startsWith(basePath);
//       }
//       return pathname === path;
//     })
//   );

//   // --- Start Modified Logic ---

//   // If no token (unauthenticated)
//   if (!token) {
//     // If the user is trying to access a protected route (but not '/'), redirect to login
//     if (isProtectedRoute) {
//       // This now only covers /admin, /user, /organization paths
//       console.log(
//         "No session found for protected route, redirecting to /auth/login",
//         { pathname }
//       );
//       return NextResponse.redirect(new URL("/auth/login", origin));
//     }
//     // If it's not a protected route (like '/'), allow access
//     console.log("No session found, but allowing access to public route:", {
//       pathname,
//     });
//     return NextResponse.next();
//   }

//   // If there is a token (authenticated)
//   const role = token.user?.role;
//   console.log("Middleware role:", role);

//   if (!role) {
//     console.warn(
//       "No role found on token, forcing login for authenticated user without role"
//     );
//     return NextResponse.redirect(new URL("/auth/login", origin));
//   }

//   // If authenticated and trying to access the root path '/', allow it.
//   if (pathname === "/") {
//     console.log("Authenticated user accessing root path:", { role, pathname });
//     return NextResponse.next();
//   }

//   // If authenticated and trying to access other protected routes
//   if (isProtectedRoute) {
//     const allowedPaths = protectedRoutes[role] || [];
//     const isAllowed = allowedPaths.some((path) => {
//       if (path.endsWith("/*")) {
//         const basePath = path.slice(0, -2);
//         return pathname.startsWith(basePath);
//       }
//       return pathname === path;
//     });

//     if (!isAllowed) {
//       console.log(
//         "Access denied for authenticated user, redirecting to role-specific dashboard:",
//         {
//           role,
//           pathname,
//           redirectUrl: redirects[role],
//         }
//       );
//       return NextResponse.redirect(
//         new URL(redirects[role] || "/auth/login", origin)
//       );
//     }
//   }

//   // Allow the request to proceed if not caught by any redirect rules
//   console.log("Allowing request:", { role, pathname });
//   return NextResponse.next();
// }

// export const config = {
//   matcher: ["/", "/admin/:path*", "/user/:path*", "/organization/:path*"],
// };

import { getToken } from "next-auth/jwt";
import { NextRequest, NextResponse } from "next/server";

export async function middleware(request: NextRequest) {
  const { pathname, origin } = request.nextUrl;

  // Skip middleware for NextAuth routes and /auth/login
  if (pathname.startsWith("/api/auth") || pathname === "/auth/login") {
    console.log("Skipping middleware for:", { pathname });
    return NextResponse.next();
  }

  type TokenWithRole = {
    user?: {
      role?: string;
      [key: string]: unknown;
    };
    [key: string]: unknown;
  };

  const token = (await getToken({
    req: request,
    secret: process.env.NEXTAUTH_SECRET,
  })) as TokenWithRole;

  const redirects: { [key: string]: string } = {
    admin: "/admin/dashboard",
    user: "/user/home",
    org: "/organization/dashboard",
  };

  const protectedRoutes: { [key: string]: string[] } = {
    admin: ["/admin", "/admin/*"],
    user: ["/user", "/user/*"],
    org: ["/organization", "/organization/*"],
  };

  const isProtectedRoute = Object.entries(protectedRoutes).some(([r, paths]) =>
    paths.some((path) => {
      if (path.endsWith("/*")) {
        const basePath = path.slice(0, -2);
        return pathname.startsWith(basePath);
      }
      return pathname === path;
    })
  );

  // --- Start Modified Logic ---

  if (!token) {
    if (isProtectedRoute) {
      console.log(
        "No session found for protected route, redirecting to /auth/login",
        { pathname }
      );
      return NextResponse.redirect(new URL("/auth/login", origin));
    }
    console.log("No session found, but allowing access to public route:", {
      pathname,
    });
    return NextResponse.next();
  }

  const role = token.user?.role;
  console.log("Middleware role:", role);

  if (!role) {
    console.warn(
      "No role found on token, forcing login for authenticated user without role"
    );
    return NextResponse.redirect(new URL("/auth/login", origin));
  }

  // --- NEW LOGIC FOR REDIRECTING FROM '/' AFTER LOGIN ---
  if (pathname === "/") {
    const desiredRedirectPath = redirects[role];
    if (desiredRedirectPath && desiredRedirectPath !== "/") {
      console.log(
        `Authenticated user (${role}) on root path, redirecting to role-specific dashboard: ${desiredRedirectPath}`
      );
      return NextResponse.redirect(new URL(desiredRedirectPath, origin));
    }
    // If no specific redirect for the role, or if it's already '/', just allow it.
    console.log("Authenticated user accessing root path:", { role, pathname });
    return NextResponse.next();
  }
  // --- END NEW LOGIC ---

  // If authenticated and trying to access other protected routes
  if (isProtectedRoute) {
    const allowedPaths = protectedRoutes[role] || [];
    const isAllowed = allowedPaths.some((path) => {
      if (path.endsWith("/*")) {
        const basePath = path.slice(0, -2);
        return pathname.startsWith(basePath);
      }
      return pathname === path;
    });

    if (!isAllowed) {
      console.log(
        "Access denied for authenticated user, redirecting to role-specific dashboard:",
        {
          role,
          pathname,
          redirectUrl: redirects[role],
        }
      );
      return NextResponse.redirect(
        new URL(redirects[role] || "/auth/login", origin)
      );
    }
  }

  console.log("Allowing request:", { role, pathname });
  return NextResponse.next();
}

export const config = {
  matcher: ["/", "/admin/:path*", "/user/:path*", "/organization/:path*"],
};