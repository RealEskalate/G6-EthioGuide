import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";
import type { ProceduresResponse } from "@/app/types/myprocedures";

// Read token from env (name: ACCESS_TOKEN or NEXT_PUBLIC_ACCESS_TOKEN)
function readEnvToken(): string | null {
  return (
    process.env.NEXT_PUBLIC_ACCESS_TOKEN || process.env.ACCESS_TOKEN || null
  );
}

// Read token from localStorage or cookie (adjust key names if different)
function readAuthToken(): string | null {
  if (typeof window === "undefined") return null;
  const fromLocalStorage =
    localStorage.getItem("accessToken") ||
    localStorage.getItem("token") ||
    localStorage.getItem("authToken");
  if (fromLocalStorage) return fromLocalStorage;
  const m = document.cookie.match(/(?:^|; )accessToken=([^;]+)/);
  return m ? decodeURIComponent(m[1]) : null;
}

const RAW_BACKEND = (process.env.NEXT_PUBLIC_API_URL || 'https://ethio-guide-backend.onrender.com').replace(/\/$/, '')
const API_BASE = /\/api\/v1$/.test(RAW_BACKEND) ? RAW_BACKEND : `${RAW_BACKEND}/api/v1/`

export const apiSlice = createApi({
  reducerPath: "api",
  baseQuery: fetchBaseQuery({
  baseUrl: API_BASE,
    prepareHeaders: (headers) => {
      // changed: prefer browser (session-persisted) token first
      const token = readAuthToken() || readEnvToken();
      if (token) headers.set("Authorization", `Bearer ${token}`);
      return headers;
    },
  }),
  tagTypes: ["Procedure"],
  endpoints: (builder) => ({
    getMyProcedures: builder.query<
      ProceduresResponse,
      { page?: number; limit?: number }
    >({
      // Direct backend endpoint
      query: ({ page = 1, limit = 20 } = {}) =>
        `myProcedures?page=${page}&limit=${limit}`,
      providesTags: ["Procedure"],
    }),
  }),
});

export const { useGetMyProceduresQuery } = apiSlice;
