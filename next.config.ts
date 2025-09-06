import type { NextConfig } from "next";

// Ensure we have a backend URL (client code also reads NEXT_PUBLIC_API_URL directly)
const RAW_BACKEND = process.env.NEXT_PUBLIC_API_URL || ''
// Normalize (strip trailing slash)
const NORMALIZED = RAW_BACKEND.replace(/\/$/, '')

const nextConfig: NextConfig = {
  turbopack: {
    // Force correct workspace root to avoid picking parent lockfile
    root: __dirname,
  },
  async rewrites() {
    // If an external backend is configured, proxy /api/v1/* to it so that
    // frontend calls to relative /api/v1 still work without CORS issues in dev.
    if (NORMALIZED && NORMALIZED.startsWith('http')) {
      // If the env already ends with /api/v1 we keep mapping path after that; otherwise append.
      const hasApiSuffix = /\/api\/v1$/.test(NORMALIZED)
      return [
        {
          source: '/api/v1/:path*',
          destination: hasApiSuffix ? `${NORMALIZED}/:path*` : `${NORMALIZED}/api/v1/:path*`
        }
      ]
    }
    return []
  }
};

export default nextConfig;
