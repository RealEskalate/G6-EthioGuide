import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'
import type { UserProcedureChecklist, PatchChecklistPayload, ChecklistItem } from '@/app/types/checklist'

type PaginatedRaw = { data?: unknown; items?: unknown; results?: unknown }

function extractArray(raw: unknown): unknown[] {
  if (Array.isArray(raw)) return raw
  const obj = (raw || {}) as PaginatedRaw
  if (Array.isArray(obj.data)) return obj.data
  if (Array.isArray(obj.items)) return obj.items
  if (Array.isArray(obj.results)) return obj.results
  return []
}

type ChecklistItemRaw = {
  id?: string; _id?: string
  title?: string; name?: string
  description?: string
  completed?: boolean; done?: boolean
  order?: number; position?: number
  updatedAt?: string; updated_at?: string
}

type ChecklistRaw = {
  id?: string; _id?: string; userProcedureId?: string
  procedureId?: string; procedure_id?: string
  procedureTitle?: string; title?: string; name?: string
  organizationName?: string; organization_name?: string; orgName?: string
  status?: UserProcedureChecklist['status']
  progress?: number
  startedAt?: string; completedAt?: string; createdAt?: string; updatedAt?: string
  started_at?: string; completed_at?: string; created_at?: string; updated_at?: string
  items?: ChecklistItemRaw[]
}

function normalizeItem(raw: unknown): ChecklistItem {
  const r = (raw || {}) as ChecklistItemRaw
  const id = r.id || r._id || crypto.randomUUID()
  return {
    id,
    title: r.title || r.name || 'Item',
    description: r.description,
    completed: Boolean(r.completed ?? r.done ?? false),
    order: r.order ?? r.position,
    updatedAt: r.updatedAt || r.updated_at,
  }
}

function normalizeChecklist(raw: unknown): UserProcedureChecklist {
  const d = (raw || {}) as ChecklistRaw
  const id = d.id || d._id || d.userProcedureId || crypto.randomUUID()
  const items: ChecklistItem[] = Array.isArray(d.items) ? d.items.map(normalizeItem) : []
  // derive progress if missing
  let progress = typeof d.progress === 'number' ? d.progress : undefined
  if ((progress === undefined || progress === null) && items.length) {
    const done = items.filter(i => i.completed).length
    progress = Math.round((done / items.length) * 100)
  }
  // derive status if missing
  let status = d.status
  if (!status) {
    if (items.length === 0) status = 'NOT_STARTED'
    else if (progress === 100) status = 'COMPLETED'
    else if (progress && progress > 0) status = 'IN_PROGRESS'
    else status = 'NOT_STARTED'
  }
  return {
    id,
    procedureId: d.procedureId || d.procedure_id || '',
    procedureTitle: d.procedureTitle || d.title || d.name,
    organizationName: d.organizationName || d.organization_name || d.orgName,
    status,
    progress,
    startedAt: d.startedAt || d.started_at,
    completedAt: d.completedAt || d.completed_at,
    createdAt: d.createdAt || d.created_at,
    updatedAt: d.updatedAt || d.updated_at,
    items,
  }
}

export const checklistsApi = createApi({
  reducerPath: 'checklistsApi',
  baseQuery: fetchBaseQuery({ baseUrl: '/api/v1' }),
  tagTypes: ['MyChecklists', 'Checklist'],
  endpoints: (builder) => ({
    getMyChecklists: builder.query<UserProcedureChecklist[], void>({
      query: () => '/myProcedures',
      transformResponse: (raw: unknown): UserProcedureChecklist[] => {
        const arr = extractArray(raw)
        return arr.map(normalizeChecklist)
      },
      providesTags: (result) => result ? [
        ...result.map(r => ({ type: 'Checklist' as const, id: r.id })),
        { type: 'MyChecklists', id: 'LIST' }
      ] : [{ type: 'MyChecklists', id: 'LIST' }]
    }),
    getChecklist: builder.query<UserProcedureChecklist, string>({
      query: (id) => `/checklists/${id}`,
      transformResponse: (raw: unknown): UserProcedureChecklist => {
        const data = (raw as { data?: unknown })?.data ?? raw
        return normalizeChecklist(data)
      },
      providesTags: (result, error, id) => [{ type: 'Checklist', id }]
    }),
    patchChecklist: builder.mutation<UserProcedureChecklist, { id: string; body: PatchChecklistPayload }>({
      query: ({ id, body }) => ({
        url: `/checklists/${id}`,
        method: 'PATCH',
        body
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: 'Checklist', id },
        { type: 'MyChecklists', id: 'LIST' }
      ]
    })
  })
})

export const { useGetMyChecklistsQuery, useGetChecklistQuery, usePatchChecklistMutation } = checklistsApi
