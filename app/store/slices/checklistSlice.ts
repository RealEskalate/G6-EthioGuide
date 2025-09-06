import { apiSlice } from "./workspaceSlice";
import type { ChecklistResponse } from "@/app/types/checklist";

export const checklistApi = apiSlice.injectEndpoints({
  endpoints: (builder) => ({
    getChecklist: builder.query<ChecklistResponse, string>({
      // Calls backend directly via apiSlice baseUrl
      query: (userProcedureId) => {
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
          url: `checklists/${encodeURIComponent(userProcedureId)}`,
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
  overrideExisting: false,
});

export const { useGetChecklistQuery } = checklistApi;
