import { apiSlice } from "./workspaceSlice";
import type { ChatGuideResponse, PostChatRequest } from "@/app/types/chat";

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

export const postChatApi = apiSlice.injectEndpoints({
  endpoints: (builder) => ({
    postChat: builder.mutation<ChatGuideResponse, PostChatRequest>({
      query: (body) => {
        const token = readToken();
        return {
          url: "ai/guide",
          method: "POST",
          body,
          headers: token ? { Authorization: `Bearer ${token}` } : undefined,
        };
      },
      async onQueryStarted(_arg, { queryFulfilled }) {
        try {
          const res = await queryFulfilled;
          console.log("postChat response:", res.data);
        } catch (err) {
          console.error("postChat error:", err);
        }
      },
    }),
  }),
  overrideExisting: false,
});

export const { usePostChatMutation } = postChatApi;
