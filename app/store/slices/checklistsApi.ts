import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'
import type { UserProcedureChecklist, PatchChecklistPayload } from '@/app/types/checklist'

interface Paginated<T> { data?: T[]; items?: T[]; results?: T[] }

function extractArray<T>(raw: Paginated<T> | T[]): T[] {
  if (Array.isArray(raw)) return raw
  return raw.data || raw.items || raw.results || []
}

export const checklistsApi = createApi({
  reducerPath: 'checklistsApi',
  baseQuery: fetchBaseQuery({ baseUrl: '/api/v1' }),
  tagTypes: ['MyChecklists', 'Checklist'],
  endpoints: (builder) => ({
    getMyChecklists: builder.query<UserProcedureChecklist[], void>({
      query: () => '/myProcedures',
      transformResponse: (raw: any): UserProcedureChecklist[] => {
        const arr = extractArray<UserProcedureChecklist>(raw)
        return arr.map(ch => {
          const id = (ch as any).id || (ch as any)._id || (ch as any).userProcedureId
          let progress = ch.progress
          if ((progress === undefined || progress === null) && ch.items?.length) {
            const done = ch.items.filter(i => i.completed).length
            progress = Math.round((done / ch.items.length) * 100)
          }
            let status = ch.status as UserProcedureChecklist['status'] | undefined
            if (!status) {
              if (!ch.items || ch.items.length === 0) status = 'NOT_STARTED'
              else if (progress === 100) status = 'COMPLETED'
              else if (progress && progress > 0) status = 'IN_PROGRESS'
              else status = 'NOT_STARTED'
            }
            return { ...ch, id, progress, status }
        })
      },
      providesTags: (result) => result ? [
        ...result.map(r => ({ type: 'Checklist' as const, id: r.id })),
        { type: 'MyChecklists', id: 'LIST' }
      ] : [{ type: 'MyChecklists', id: 'LIST' }]
    }),
    getChecklist: builder.query<UserProcedureChecklist, string>({
      query: (id) => `/checklists/${id}`,
      transformResponse: (raw: any): UserProcedureChecklist => {
        const data = raw?.data || raw
        const idVal = data.id || data._id || data.userProcedureId
        let progress = data.progress
        if ((progress === undefined || progress === null) && data.items?.length) {
          const done = data.items.filter((i: any) => i.completed).length
          progress = Math.round((done / data.items.length) * 100)
        }
        let status = data.status as UserProcedureChecklist['status'] | undefined
        if (!status) {
          if (!data.items || data.items.length === 0) status = 'NOT_STARTED'
          else if (progress === 100) status = 'COMPLETED'
          else if (progress && progress > 0) status = 'IN_PROGRESS'
          else status = 'NOT_STARTED'
        }
        return { ...data, id: idVal, progress, status }
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
