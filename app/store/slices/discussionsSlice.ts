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

// helper: try to read current user id from localStorage (adjust keys if needed)
function readCurrentUserId(): string | null {
  if (typeof window === "undefined") return null;
  // direct id keys
  const directKeys = ["userId", "userID", "uid", "id", "currentUserId", "_id"];
  for (const k of directKeys) {
    const val = localStorage.getItem(k);
    if (val && String(val).trim()) return String(val).trim();
  }
  // try common JSON containers
  const jsonKeys = ["user", "profile", "account", "auth", "currentUser"];
  for (const k of jsonKeys) {
    const raw = localStorage.getItem(k);
    if (!raw) continue;
    try {
      const obj = JSON.parse(raw);
      const cand =
        obj?.id ??
        obj?._id ??
        obj?.userId ??
        obj?.userID ??
        obj?.uid ??
        obj?.user?.id ??
        obj?.user?._id ??
        obj?.user?.userId ??
        obj?.user?.userID;
      if (cand) return String(cand);
    } catch {
      /* ignore */
    }
  }
  return null;
}

// helper: read access token from storage/cookie (same keys used elsewhere)
function readAuthToken(): string | null {
  if (typeof window === "undefined") return null;
  const ls =
    localStorage.getItem("accessToken") ||
    localStorage.getItem("token") ||
    localStorage.getItem("authToken");
  if (ls) return ls;
  const m = document.cookie.match(/(?:^|; )accessToken=([^;]+)/);
  return m ? decodeURIComponent(m[1]) : null;
}

// helper: decode JWT and try common user id fields
function readUserIdFromToken(): string | null {
  try {
    const token = readAuthToken();
    if (!token) return null;
    const [, payload] = token.split(".");
    if (!payload) return null;
    const json = JSON.parse(atob(payload));
    const candidates = [
      json.userId,
      json.userID,
      json.user_id, // added
      json.UserId, // added
      json.UserID, // added
      json.id,
      json._id,
      json.sub,
      json.uid,
      json.user?.id,
      json.user?._id,
      json.user?.userId,
      json.user?.userID,
      json.user?.user_id, // added
    ].filter(Boolean);
    return candidates.length ? String(candidates[0]) : null;
  } catch {
    return null;
  }
}

// helper: normalize user id field from a post
function getPostUserId(p: unknown): string | undefined {
  const any = p as Record<string, unknown>;
  const cand =
    any?.UserID ?? // API uses this shape
    any?.userID ?? // tolerance
    any?.userId ?? // tolerance
    any?.user_id ?? // tolerance
    any?.UserId; // tolerance
  return cand ? String(cand) : undefined;
}

// small helper to extract readable error message
function extractErrorMessage(
  e: unknown,
  fallback = "Something went wrong."
): string {
  if (!e) return fallback;
  if (typeof e === "string") return e;
  if (typeof e === "object") {
    const errObj = e as Record<string, unknown>;
    return (
      (errObj?.data as { message?: string })?.message ||
      (errObj?.error as string) ||
      (errObj?.message as string) ||
      (errObj?.status as string) ||
      fallback
    );
  }
  return fallback;
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
      {
        page?: number;
        limit?: number;
        userId?: string;
        selfOnly?: boolean;
      } | void
    >({
      query: (args) => {
        const page = args?.page ?? 0;
        const limit = args?.limit ?? 10;
        const userId = args && "userId" in args ? args.userId : undefined;
        const params = new URLSearchParams({
          page: String(page),
          limit: String(limit),
        });
        if (userId) params.set("userId", userId);
        return `discussions?${params.toString()}`;
      },
      transformResponse: (res: unknown, _meta, arg): DiscussionsList => {
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

        if ("posts" in box && Array.isArray(box.posts)) posts = box.posts;
        if ("total" in box && typeof box.total === "number") total = box.total;
        if ("page" in box && typeof box.page === "number") page = box.page;
        if ("limit" in box && typeof box.limit === "number") limit = box.limit;
        else if (posts.length) limit = posts.length;

        // strict client-side filter for "My Discussions"
        type DiscussionsArg = {
          page?: number;
          limit?: number;
          userId?: string;
          selfOnly?: boolean;
        } | void;
        const wantsSelfOnly =
          Boolean((arg as DiscussionsArg & { selfOnly?: boolean })?.selfOnly) ||
          (typeof window !== "undefined" &&
            window.location?.pathname?.includes("/user/my-discussions"));

        const passedUserId = (arg as DiscussionsArg & { userId?: string })
          ?.userId as string | undefined;

        let currentUserId: string | undefined = passedUserId;
        if (!currentUserId && wantsSelfOnly) {
          currentUserId =
            readCurrentUserId() ?? readUserIdFromToken() ?? undefined;
        }

        if (wantsSelfOnly) {
          if (currentUserId) {
            posts = posts.filter(
              (p) => getPostUserId(p) === String(currentUserId)
            );
            total = posts.length;
          } else {
            // if we can't determine the user id, don't show other users' posts
            posts = [];
            total = 0;
          }
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
      // ensure tags are preserved if not explicitly provided
      async queryFn(
        arg,
        api,
        _extraOptions,
        fetchWithBQ
      ): Promise<
        import("@reduxjs/toolkit/query").QueryReturnValue<
          DiscussionDetail,
          import("@reduxjs/toolkit/query").FetchBaseQueryError,
          import("@reduxjs/toolkit/query").FetchBaseQueryMeta | undefined
        >
      > {
        let preservedTags: string[] | undefined;
        try {
          const state = api.getState() as unknown;
          const root =
            typeof state === "object" && state !== null
              ? (state as Record<string, unknown>)[apiSlice.reducerPath]
              : undefined;
          const entries = Object.values(
            (root as { queries?: unknown })?.queries ?? {}
          );
          for (const entry of entries as unknown[]) {
            if (
              (entry as { endpointName?: string })?.endpointName ===
                "getDiscussions" &&
              Array.isArray(
                (entry as { data?: { posts?: unknown[] } })?.data?.posts
              )
            ) {
              const found = (
                entry as { data?: { posts?: unknown[] } }
              ).data?.posts?.find(
                (p) => String((p as { ID?: string }).ID) === String(arg.id)
              );
              if (
                found &&
                Array.isArray((found as { Tags?: unknown[] }).Tags)
              ) {
                preservedTags = (found as { Tags?: unknown[] }).Tags!.map(
                  (t: unknown) => String(t)
                );
                break;
              }
            }
          }
        } catch {
          // ignore cache scan errors
        }

        const tagsToSend =
          typeof arg.data.tags !== "undefined" ? arg.data.tags : preservedTags;
        const body: Partial<{
          title: string;
          content: string;
          tags: string[];
          procedures: string[];
        }> = { ...arg.data };
        if (typeof tagsToSend !== "undefined") body.tags = tagsToSend;

        const result = await fetchWithBQ({
          url: `discussions/${arg.id}`,
          method: "PATCH",
          body,
        });

        // Ensure the returned data is typed as DiscussionDetail and never undefined
        if (result.error) {
          return {
            error: result.error,
            data: undefined,
            meta: result.meta,
          };
        }
        if (result.data) {
          return {
            error: undefined,
            data: result.data as DiscussionDetail,
            meta: result.meta,
          };
        }
        // fallback: should not happen, but for type safety
        return {
          error: {
            status: "CUSTOM_ERROR",
            data: "No data returned from updateDiscussion",
            error: "No data returned from updateDiscussion",
          },
          data: undefined,
          meta: result.meta,
        };
      },
      async onQueryStarted(_arg, { queryFulfilled }) {
        try {
          await queryFulfilled;
          if (typeof window !== "undefined") {
            const { toast } = await import("react-hot-toast");
            toast.success("Changes saved successfully.", {
              icon: "✅",
              style: {
                background: "#f0fdf4",
                color: "#065f46",
                border: "1px solid #86efac",
              },
            });
          }
        } catch (e) {
          if (typeof window !== "undefined") {
            const { toast } = await import("react-hot-toast");
            toast.error(extractErrorMessage(e, "Failed to save changes."), {
              icon: "⚠️",
              style: {
                background: "#fef2f2",
                color: "#991b1b",
                border: "1px solid #fecaca",
              },
            });
          }
        }
      },
    }),
    // added: delete discussion
    deleteDiscussion: builder.mutation<{ success?: boolean } | unknown, string>(
      {
        query: (id) => ({
          url: `discussions/${id}`,
          method: "DELETE",
        }),
        async onQueryStarted(id, { queryFulfilled }) {
          try {
            await queryFulfilled;
            if (typeof window !== "undefined") {
              const { toast } = await import("react-hot-toast");
              toast.success("Discussion deleted.", {
                icon: "🗑️",
                style: {
                  background: "#f0fdf4",
                  color: "#065f46",
                  border: "1px solid #86efac",
                },
              });
            }
          } catch (e) {
            if (typeof window !== "undefined") {
              const { toast } = await import("react-hot-toast");
              toast.error(
                extractErrorMessage(e, "Failed to delete discussion."),
                {
                  icon: "⚠️",
                  style: {
                    background: "#fef2f2",
                    color: "#991b1b",
                    border: "1px solid #fecaca",
                  },
                }
              );
            }
          }
        },
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
