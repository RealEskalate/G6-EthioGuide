import { apiSlice } from "./workspaceSlice";
import type { ChecklistResponse } from "@/app/types/checklist";

export const checklistApi = apiSlice.injectEndpoints({
  endpoints: (builder) => ({
    getChecklist: builder.query<ChecklistResponse, string>({
      //   Calls backend directly via apiSlice baseUrl (https://ethio-guide-backend.onrender.com/api/v1/)
      query: (userProcedureId) =>
        `checklists/${encodeURIComponent(userProcedureId)}`,
      providesTags: ["Procedure"],
    }),
  }),
  overrideExisting: false,
});

export const { useGetChecklistQuery } = checklistApi;
