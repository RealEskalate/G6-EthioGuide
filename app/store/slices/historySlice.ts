import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";
import type { ChatHistoryItem } from "@/app/types/chat";

function readEnvToken(): string | null {
  return (
    process.env.NEXT_PUBLIC_ACCESS_TOKEN || process.env.ACCESS_TOKEN || null
  );
}
function readAuthToken(): string | null {
  if (typeof window === "undefined") return null;
  return localStorage.getItem("accessToken") || null;
}

export const historyApi = createApi({
  reducerPath: "historyApi",
  baseQuery: fetchBaseQuery({
    baseUrl: "https://ethio-guide-backend.onrender.com/api/v1",
    prepareHeaders: (headers) => {
      const token = readAuthToken() || readEnvToken();
      if (token) headers.set("Authorization", `Bearer ${token}`);
      return headers;
    },
  }),
  endpoints: (builder) => ({
    getChatHistory: builder.query<ChatHistoryItem[], void>({
      query: () => "/ai/history",
      transformResponse: (res: unknown): ChatHistoryItem[] => {
        console.log("Raw history response:", res);
        type RawHistoryItem = {
          id?: string;
          _id?: string;
          uuid?: string;
          request?: string;
          procedures?: { name?: string }[];
          response?: string;
          updatedAt?: string;
          createdAt?: string;
        };
        const r = res as { history?: RawHistoryItem[] } | RawHistoryItem[];
        const items =
          Array.isArray(r) || Array.isArray(r?.history)
            ? Array.isArray(r)
              ? r
              : r.history
            : [];
        return (items ?? []).map(
          (it: RawHistoryItem): ChatHistoryItem => ({
            id: String(
              it.id ?? it._id ?? it.uuid ?? Math.random().toString(36).slice(2)
            ),
            title: String(it.request ?? it.procedures?.[0]?.name ?? "Untitled"),
            lastMessage: String(it.response ?? ""),
            timestamp: String(it.updatedAt ?? it.createdAt ?? ""),
            messageCount: Number(it.procedures?.length ?? 0),
          })
        );
      },
    }),
    postTranslate: builder.mutation<
      { translated?: string; lang?: string } | unknown,
      { content: string; lang: string }
    >({
      query: (body) => ({
        url: "/translate",
        method: "POST",
        body,
      }),
    }),
  }),
});

export const {
  useGetChatHistoryQuery,
  useLazyGetChatHistoryQuery,
  usePostTranslateMutation,
} = historyApi;
