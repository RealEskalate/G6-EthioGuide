import { apiSlice } from "@/app/store/slices/workspaceSlice";

export interface CreateDiscussionPayload {
  title: string;
  content: string;
  tags?: string[];
  attachments?: string[];
  procedureId?: string;
  userProcedureId?: string;
}

export const discussionsApi = apiSlice.injectEndpoints({
  endpoints: (builder) => ({
    createDiscussion: builder.mutation<unknown, CreateDiscussionPayload>({
      query: (body) => ({
        url: "discussions",
        method: "POST",
        body,
      }),
    }),
  }),
  overrideExisting: false,
});

export const { useCreateDiscussionMutation } = discussionsApi;
