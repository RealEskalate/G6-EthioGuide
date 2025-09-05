import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";
import type { ProceduresResponse } from "@/app/types/myprocedures";

export const apiSlice = createApi({
  reducerPath: "api",
  baseQuery: fetchBaseQuery({
    // Back to direct backend URL now that CORS is fixed
    baseUrl: "https://ethio-guide-backend.onrender.com/api/v1/",
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
