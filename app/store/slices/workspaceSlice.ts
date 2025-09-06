import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";
import type { ProceduresResponse } from "@/app/types/myprocedures";

// Read token from env (name: ACCESS_TOKEN or NEXT_PUBLIC_ACCESS_TOKEN)
function readEnvToken(): string | null {
  return (
    process.env.NEXT_PUBLIC_ACCESS_TOKEN || process.env.ACCESS_TOKEN || null
  );
}

// Read token from localStorage/sessionStorage or cookie (more keys covered)
function readAuthToken(): string | null {
  if (typeof window === "undefined") return null;
  const ls =
    localStorage.getItem("accessToken") ||
    localStorage.getItem("token") ||
    localStorage.getItem("authToken") ||
    localStorage.getItem("access_token");
  const ss =
    sessionStorage.getItem("accessToken") ||
    sessionStorage.getItem("token") ||
    sessionStorage.getItem("authToken") ||
    sessionStorage.getItem("access_token");
  if (ls) return ls;
  if (ss) return ss;
  const m = document.cookie.match(/(?:^|; )accessToken=([^;]+)/);
  return m ? decodeURIComponent(m[1]) : null;
}

export const apiSlice = createApi({
  reducerPath: "api",
  baseQuery: fetchBaseQuery({
    baseUrl: "https://ethio-guide-backend.onrender.com/api/v1/",
    credentials: "include", // added
    prepareHeaders: (headers) => {
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
      query: ({ page = 1, limit = 20 } = {}) => {
        // explicit header like discussions slice
        const lsToken =
          typeof window !== "undefined"
            ? localStorage.getItem("accessToken") ||
              localStorage.getItem("token") ||
              localStorage.getItem("authToken")
            : null;
        const envToken =
          process.env.NEXT_PUBLIC_ACCESS_TOKEN ||
          process.env.ACCESS_TOKEN ||
          null;

        return {
          url: `checklists/myprocedures?page=${page}&limit=${limit}`,
          method: "GET",
          headers:
            lsToken || envToken
              ? { Authorization: `Bearer ${lsToken ?? envToken}` }
              : undefined,
        };
      },
      providesTags: ["Procedure"],
    }),
  }),
});

export const { useGetMyProceduresQuery } = apiSlice;
