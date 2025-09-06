import { apiSlice } from "./workspaceSlice";
import type { ChecklistResponse } from "@/app/types/checklist";

function readToken(): string | null {
  if (typeof window !== "undefined") {
    const ls =
      localStorage.getItem("accessToken") ||
      localStorage.getItem("access_token") ||
      localStorage.getItem("token") ||
      localStorage.getItem("authToken");
    const ss =
      sessionStorage.getItem("accessToken") ||
      sessionStorage.getItem("access_token") ||
      sessionStorage.getItem("token") ||
      sessionStorage.getItem("authToken");
    const cookieMatch =
      typeof document !== "undefined"
        ? document.cookie.match(/(?:^|; )accessToken=([^;]+)/)
        : null;
    const cookieToken = cookieMatch ? decodeURIComponent(cookieMatch[1]) : null;
    return ls || ss || cookieToken || null;
  }
  return (
    process.env.NEXT_PUBLIC_ACCESS_TOKEN || process.env.ACCESS_TOKEN || null
  );
}

export const checklistApi = apiSlice.injectEndpoints({
  endpoints: (builder) => ({
    getChecklist: builder.query<ChecklistResponse, string>({
      // Calls backend directly via apiSlice baseUrl
      query: (userProcedureId) => {
        const token = readToken();
        return {
          url: `checklists/${encodeURIComponent(userProcedureId)}`,
          method: "GET",
          headers: token ? { Authorization: `Bearer ${token}` } : undefined,
        };
      },
      providesTags: ["Procedure"],
    }),
  }),
  overrideExisting: false,
});

export const { useGetChecklistQuery } = checklistApi;
