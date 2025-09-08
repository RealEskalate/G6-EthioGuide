import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'
// import type { FetchBaseQueryError } from '@reduxjs/toolkit/query'
import type { 
  UserProcedureChecklist, 
  CreateChecklistResponse, 
  // GetMyProceduresResponse, 
  GetChecklistItemsResponse 
} from '@/app/types/checklist'

type UnknownRecord = Record<string, unknown>;

// (removed unused Paginated type)

function extractUnknownArray(raw: unknown): unknown[] {
  if (Array.isArray(raw)) return raw;
  if (raw && typeof raw === "object") {
    const obj = raw as UnknownRecord;
    // Common envelope keys first
    for (const key of [
      "data",
      "items",
      "results",
      "list",
      "myProcedures",
      "checklists",
      "userProcedures",
      "procedures",
    ]) {
      const v = obj[key];
      if (Array.isArray(v)) return v;
      if (
        v &&
        typeof v === "object" &&
        Array.isArray((v as UnknownRecord)["data"])
      )
        return (v as UnknownRecord)["data"] as unknown[];
    }
    // Fallback: find the first array of objects inside any property
    for (const val of Object.values(obj)) {
      if (Array.isArray(val) && val.some((e) => e && typeof e === "object"))
        return val as unknown[];
      if (val && typeof val === "object") {
        const inner = extractUnknownArray(val);
        if (inner.length) return inner;
      }
    }
  }
  return [];
}

function readString(obj: UnknownRecord, key: string): string | undefined {
  const v = obj[key];
  return typeof v === "string" ? v : undefined;
}

function readNumber(obj: UnknownRecord, key: string): number | undefined {
  const v = obj[key];
  return typeof v === "number" ? v : undefined;
}

function deriveChecklistFromMyProcedures(u: unknown): UserProcedureChecklist {
  const obj = u && typeof u === "object" ? (u as UnknownRecord) : {};
  const id = readString(obj, "id") || "";
  const procedureId = readString(obj, "procedure_id") || "";
  const percent = readNumber(obj, "percent") || 0;
  const statusStr = readString(obj, "status") || "NOT_STARTED";
  const updatedAt = readString(obj, "updated_at") || "";
  const userId = readString(obj, "user_id") || "";

  // Convert status string to our enum
  let status: UserProcedureChecklist["status"] = "NOT_STARTED";
  if (statusStr === "COMPLETED" || percent === 100) {
    status = "COMPLETED";
  } else if (statusStr === "IN_PROGRESS" || percent > 0) {
    status = "IN_PROGRESS";
  }

  return {
    id,
    procedureId,
    status,
    progress: percent,
    updatedAt,
    items: [],
  };
}

function deriveChecklistFromItems(
  items: GetChecklistItemsResponse[]
): UserProcedureChecklist {
  if (!items || items.length === 0) {
    return {
      id: "",
      procedureId: "",
      status: "NOT_STARTED",
      progress: 0,
      items: [],
    };
  }

  const userProcedureId = items[0]?.user_procedure_id || "";
  const completedCount = items.filter((item) => item.is_checked).length;
  const progress = Math.round((completedCount / items.length) * 100);

  let status: UserProcedureChecklist["status"] = "NOT_STARTED";
  if (progress === 100) {
    status = "COMPLETED";
  } else if (progress > 0) {
    status = "IN_PROGRESS";
  }

  const checklistItems = items.map((item) => ({
    id: item.id,
    content: item.content,
    is_checked: item.is_checked,
    type: item.type,
    user_procedure_id: item.user_procedure_id,
  }));

  return {
    id: userProcedureId,
    procedureId: "",
    status,
    progress,
    items: checklistItems,
  };
}

// Use the proxy path to avoid CORS issues
const API_BASE = "/api/v1";

const rawBaseQuery = fetchBaseQuery({
  baseUrl: API_BASE,
  prepareHeaders: (headers) => {
    if (typeof window !== "undefined") {
      const token =
        localStorage.getItem("accessToken") ||
        localStorage.getItem("token") ||
        localStorage.getItem("authToken");
      if (token) headers.set("Authorization", `Bearer ${token}`);
    }
    return headers;
  },
});

export const checklistsApi = createApi({
  reducerPath: "checklistsApi",
  baseQuery: rawBaseQuery,
  tagTypes: ["MyChecklists", "Checklist"],
  endpoints: (builder) => ({
    getMyChecklists: builder.query<
      UserProcedureChecklist[],
      { token?: string | null } | void
    >({
      query: (arg) => ({
        url: "/checklists/myProcedures",
        headers:
          arg && arg.token
            ? { Authorization: `Bearer ${arg.token}` }
            : undefined,
            'lang': localStorage.getItem("i18nextLng") || "en",
      }),
      transformResponse: (raw: unknown): UserProcedureChecklist[] => {
        const arr = extractUnknownArray(raw);
        return arr.map(deriveChecklistFromMyProcedures);
      },
      providesTags: (result) =>
        result
          ? [
              ...result.map((r) => ({ type: "Checklist" as const, id: r.id })),
              { type: "MyChecklists", id: "LIST" },
            ]
          : [{ type: "MyChecklists", id: "LIST" }],
    }),
    getChecklist: builder.query<
      UserProcedureChecklist,
      { id: string; token?: string | null }
    >({
      query: ({ id, token }) => ({
        url: `/checklists/${id}`,
        headers: token ? { Authorization: `Bearer ${token}` } : undefined,
      }),
      transformResponse: (raw: unknown): UserProcedureChecklist => {
        // Backend returns an array of items: [{ id, content, is_checked, type, user_procedure_id }]
        const arr = extractUnknownArray(raw);
        const items: GetChecklistItemsResponse[] = arr.map((item: unknown) => {
          const obj =
            item && typeof item === "object" ? (item as UnknownRecord) : {};
          return {
            id: readString(obj, "id") || "",
            content: readString(obj, "content") || "",
            is_checked: Boolean(obj["is_checked"]),
            type: readString(obj, "type") || "",
            user_procedure_id: readString(obj, "user_procedure_id") || "",
          };
        });
        return deriveChecklistFromItems(items);
      },
      providesTags: (result) =>
        result ? [{ type: "Checklist", id: result.id }] : [],
    }),
    createChecklist: builder.mutation<
      CreateChecklistResponse,
      { procedureId: string; token?: string | null }
    >({
      query: ({ procedureId, token }) => ({
        url: "/checklists",
        method: "POST",
        body: { procedure_id: procedureId },
        headers: token ? { Authorization: `Bearer ${token}` } : undefined,
      }),
      transformResponse: (raw: unknown): CreateChecklistResponse => {
        const obj =
          raw && typeof raw === "object" ? (raw as UnknownRecord) : {};
        return {
          id: readString(obj, "id") || "",
          percent: readNumber(obj, "percent") || 0,
          procedure_id: readString(obj, "procedure_id") || "",
          status: readString(obj, "status") || "",
          updated_at: readString(obj, "updated_at") || "",
          user_id: readString(obj, "user_id") || "",
        };
      },
      invalidatesTags: [{ type: "MyChecklists", id: "LIST" }],
    }),
    patchChecklist: builder.mutation<
      GetChecklistItemsResponse[],
      { id: string; isChecked?: boolean; token?: string | null }
    >({
      query: ({ id, isChecked, token }) => ({
        url: `/checklists/${id}`,
        method: "PATCH",
        body:
          typeof isChecked === "boolean"
            ? { is_checked: isChecked }
            : undefined,
        headers: token ? { Authorization: `Bearer ${token}` } : undefined,
      }),
      transformResponse: (raw: unknown): GetChecklistItemsResponse[] => {
        const arr = extractUnknownArray(raw);
        return arr.map((item: unknown) => {
          const obj =
            item && typeof item === "object" ? (item as UnknownRecord) : {};
          return {
            id: readString(obj, "id") || "",
            content: readString(obj, "content") || "",
            is_checked: Boolean(obj["is_checked"]),
            type: readString(obj, "type") || "",
            user_procedure_id: readString(obj, "user_procedure_id") || "",
          };
        });
      },
      invalidatesTags: (result, error, { id }) => [
        { type: "Checklist", id },
        { type: "MyChecklists", id: "LIST" },
      ],
    }),
  }),
});

export const {
  useGetMyChecklistsQuery,
  useGetChecklistQuery,
  useCreateChecklistMutation,
  usePatchChecklistMutation,
} = checklistsApi;
