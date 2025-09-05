import { getToken } from 'next-auth/jwt';
import { NextRequest, NextResponse } from 'next/server';

export async function middleware(request: NextRequest) {
  const { pathname, origin } = request.nextUrl;
  

  // Skip middleware for NextAuth routes and /auth/login
  if (pathname.startsWith('/api/auth') || pathname === '/auth/login') {
    console.log('Skipping middleware for:', { pathname });
    return NextResponse.next();
  }

// Define a type for the token to include user and role
type TokenWithRole = {
  user?: {
    role?: string;
    [key: string]: unknown;
  };
  [key: string]: unknown;
};

// Get the session token
const token = await getToken({ req: request, secret: process.env.NEXTAUTH_SECRET }) as TokenWithRole;

// If no token (unauthenticated), redirect to login
if (!token) {
  console.log('No session found, redirecting to /auth/login', { pathname });
  return NextResponse.redirect(new URL('/auth/login', origin));
}

// Extract role, default to 'user' if undefined
const role = token.user?.role;
console.log("Middleware role:", role);
if (!role) {
  console.warn("No role found on token, forcing login");
  return NextResponse.redirect(new URL('/auth/login', origin));
}


  // Define role-based redirect rules
  const redirects: { [key: string]: string } = {
    admin: '/admin/dashboard',
    user: '/user/home',
    org: '/organization/dashboard',
  };

  // Define protected routes by role
  const protectedRoutes: { [key: string]: string[] } = {
    admin: ['/admin', '/admin/*'],
    user: ['/user', '/user/*'],
    org: ['/organization', '/organization/*'],
  };

  // Check if the current path is a protected route
  const isProtectedRoute = Object.entries(protectedRoutes).some(([r, paths]) =>
    paths.some((path) => {
      if (path.endsWith('/*')) {
        const basePath = path.slice(0, -2);
        return pathname.startsWith(basePath);
      }
      return pathname === path;
    })
  );

  // Redirect authenticated users from root to their role-specific dashboard
  if (pathname === '/') {
    const redirectUrl = redirects[role] || redirects.user;
    console.log('Redirecting to role-specific dashboard:', { role, redirectUrl });
    return NextResponse.redirect(new URL(redirectUrl, origin));
  }

  // Restrict access to protected routes
  if (isProtectedRoute) {
    const allowedPaths = protectedRoutes[role] || [];
    const isAllowed = allowedPaths.some((path) => {
      if (path.endsWith('/*')) {
        const basePath = path.slice(0, -2);
        return pathname.startsWith(basePath);
      }
      return pathname === path;
    });

    if (!isAllowed) {
      console.log('Access denied, redirecting to role-specific dashboard:', {
        role,
        pathname,
        redirectUrl: redirects[role],
      });
      return NextResponse.redirect(new URL(redirects[role] || '/auth/login', origin));
    }
  }

  // Allow the request to proceed
  console.log('Allowing request:', { role, pathname });
  return NextResponse.next();
}

export const config = {
  matcher: [
    '/',
    '/admin/:path*',
    '/user/:path*',
    '/organization/:path*',
  ],
};