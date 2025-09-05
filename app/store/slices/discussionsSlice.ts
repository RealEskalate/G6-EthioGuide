import { apiSlice } from "@/app/store/slices/workspaceSlice";

export interface CreateDiscussionPayload {
  title: string;
  content: string;
  tags?: string[];
  attachments?: string[];
  procedureId?: string;
  userProcedureId?: string;
  procedures?: string[]; // optional in type
}

// add: GET types
export interface DiscussionPost {
  ID: string;
  UserID: string;
  Title: string;
  Content: string;
  Procedures: string[] | null;
  Tags: string[];
  CreatedAt: string; // ISO
  UpdatedAt: string; // ISO
}
export interface DiscussionsList {
  posts: DiscussionPost[];
  total: number;
  page: number;
  limit: number;
}

// add: single discussion detail type (matches backend schema)
export interface DiscussionDetail {
  content: string;
  createdAt: string;
  id: string;
  procedures: string[];
  tags: string[];
  title: string;
  updatedAt: string;
  userID: string;
}

export const discussionsApi = apiSlice.injectEndpoints({
  endpoints: (builder) => ({
    createDiscussion: builder.mutation<unknown, CreateDiscussionPayload>({
      query: (body) => {
        // prefer session-persisted token in localStorage; fallback to env
        const lsToken =
          typeof window !== "undefined"
            ? localStorage.getItem("accessToken") || undefined
            : undefined;
        const envToken =
          process.env.NEXT_PUBLIC_ACCESS_TOKEN ||
          process.env.ACCESS_TOKEN ||
          undefined;

        return {
          url: "discussions",
          method: "POST",
          body: {
            ...body,
            procedures: body.procedures ?? [], // ensure [] is posted
          },
          headers:
            lsToken || envToken
              ? { Authorization: `Bearer ${lsToken ?? envToken}` }
              : undefined,
        };
      },
    }),
    // add: GET /discussions
    getDiscussions: builder.query<
      DiscussionsList,
      { page?: number; limit?: number } | void
    >({
      query: (args) => {
        const page = args?.page ?? 0;
        const limit = args?.limit ?? 10;
        return `discussions?page=${page}&limit=${limit}`;
      },
      transformResponse: (res: unknown): DiscussionsList => {
        // backend might wrap in { Posts: { posts, total, page, limit } }
        const r = res as
          | {
              Posts?: {
                posts?: DiscussionPost[];
                total?: number;
                page?: number;
                limit?: number;
              };
            }
          | {
              posts?: DiscussionPost[];
              total?: number;
              page?: number;
              limit?: number;
            };

        const box =
          (
            r as {
              Posts?: {
                posts?: DiscussionPost[];
                total?: number;
                page?: number;
                limit?: number;
              };
            }
          )?.Posts ?? r;
        let posts: DiscussionPost[] = [];
        let total = 0;
        let page = 0;
        let limit = 10;

        if ("posts" in box && Array.isArray(box.posts)) {
          posts = box.posts;
        }
        if ("total" in box && typeof box.total === "number") {
          total = box.total;
        }
        if ("page" in box && typeof box.page === "number") {
          page = box.page;
        }
        if ("limit" in box && typeof box.limit === "number") {
          limit = box.limit;
        } else if (posts.length) {
          limit = posts.length;
        }

        return { posts, total, page, limit };
      },
    }),
    // add: GET /discussions/{id}
    getDiscussionById: builder.query<DiscussionDetail, string>({
      query: (id) => `discussions/${id}`,
    }),
    // added: update discussion
    updateDiscussion: builder.mutation<
      DiscussionDetail,
      {
        id: string;
        data: Partial<{
          title: string;
          content: string;
          tags: string[];
          procedures: string[];
        }>;
      }
    >({
      query: ({ id, data }) => ({
        url: `discussions/${id}`,
        method: "PATCH",
        body: data,
      }),
    }),
    // added: delete discussion
    deleteDiscussion: builder.mutation<{ success?: boolean } | unknown, string>(
      {
        query: (id) => ({
          url: `discussions/${id}`,
          method: "DELETE",
        }),
      }
    ),
  }),
  overrideExisting: true,
});

export const {
  useCreateDiscussionMutation,
  useGetDiscussionsQuery,
  useLazyGetDiscussionByIdQuery, // added earlier
  useUpdateDiscussionMutation, // added
  useDeleteDiscussionMutation, // added
} = discussionsApi;
