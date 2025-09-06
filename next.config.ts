import type { NextConfig } from "next";

const RAW = process.env.NEXT_PUBLIC_API_URL || 'https://ethio-guide-backend.onrender.com'
const BASE = RAW.replace(/\/$/, '')

const nextConfig: NextConfig = {
  async rewrites() {
    // Proxy /api/v1/* to backend to avoid CORS and 404 when frontend dev server lacks those routes
    const hasApiSuffix = /\/api\/v1$/.test(BASE)
    return [
      {
        source: '/api/v1/:path*',
        destination: hasApiSuffix ? `${BASE}/:path*` : `${BASE}/api/v1/:path*`
      }
    ]
  }
};

export default nextConfig;
