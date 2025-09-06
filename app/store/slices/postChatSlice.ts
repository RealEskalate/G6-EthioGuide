import { apiSlice } from "./workspaceSlice";
import type { ChatGuideResponse, PostChatRequest } from "@/app/types/chat";

export const postChatApi = apiSlice.injectEndpoints({
  endpoints: (builder) => ({
    postChat: builder.mutation<ChatGuideResponse, PostChatRequest>({
      query: (body) => ({
        url: "ai/guide",
        method: "POST",
        body,
      }),
    }),
  }),
  overrideExisting: false,
});

export const { usePostChatMutation } = postChatApi;
