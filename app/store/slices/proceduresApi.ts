import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'
import type { Procedure } from '@/app/types/procedure'

interface Paginated<T> { data?: T[]; items?: T[]; results?: T[]; hasNext?: boolean }

export const proceduresApi = createApi({
  reducerPath: 'proceduresApi',
  // Use relative path so Next.js rewrite (in next.config.ts) proxies to external backend server-side (avoids CORS)
  baseQuery: fetchBaseQuery({ baseUrl: '/api/v1' }),
  tagTypes: ['Procedure', 'Procedures'],
  endpoints: (builder) => ({
    listProcedures: builder.query<{ list: Procedure[]; hasNext: boolean }, { page?: number; limit?: number; q?: string } | void>({
      query: (args) => {
        const params = new URLSearchParams()
        if (args?.page) params.set('page', String(args.page))
        if (args?.limit) params.set('limit', String(args.limit))
        if (args?.q) params.set('q', args.q)
        return `/procedures${params.size ? `?${params.toString()}` : ''}`
      },
      transformResponse: (raw: Paginated<Procedure> | Procedure[]): { list: Procedure[]; hasNext: boolean } => {
        const arr = (raw as Paginated<Procedure>).data || (raw as Paginated<Procedure>).items || (raw as Paginated<Procedure>).results || (Array.isArray(raw) ? raw as Procedure[] : [])
        const hasNext = Boolean((raw as Paginated<Procedure>).hasNext)
        return { list: arr.map(p => ({ ...p, id: (p as any).id || (p as any)._id })), hasNext }
      },
      providesTags: (result) => result ? [
        ...result.list.map(p => ({ type: 'Procedure' as const, id: p.id })),
        { type: 'Procedures', id: 'LIST' }
      ] : [{ type: 'Procedures', id: 'LIST' }]
    }),
    getProcedure: builder.query<Procedure, string>({
      query: (id) => `/procedures/${id}`,
      transformResponse: (raw: any): Procedure => {
        // Accept various envelope shapes
        let source = raw
        if (raw?.data && typeof raw.data === 'object' && !Array.isArray(raw.data)) source = raw.data
        if (raw?.procedure) source = raw.procedure
        if (Array.isArray(source)) source = source[0] || {}

        const id = source?.id || source?._id || source?.procedureId || source?.slug || undefined

        // Build steps
        let steps = source.steps
        if (!steps && source?.content?.steps) {
          const rawSteps = source.content.steps
          if (Array.isArray(rawSteps)) {
            steps = rawSteps.map((s: any, i: number) => typeof s === 'string' ? ({ order: i + 1, text: s }) : ({ order: s.order || i + 1, text: s.text || JSON.stringify(s) }))
          } else if (typeof rawSteps === 'object') {
            steps = Object.entries(rawSteps).map(([k, v]: any, idx) => ({ order: Number(k) || idx + 1, text: typeof v === 'string' ? v : JSON.stringify(v) }))
          }
        }

        // Documents
        let documentsRequired = source.documentsRequired || source.documents || source.docs || source?.content?.documentsRequired || source?.content?.documents
        if (Array.isArray(documentsRequired)) {
          documentsRequired = documentsRequired.map((d: any) => ({ name: d.name || d.title || d.label || d.filename || 'Document', templateUrl: d.templateUrl || d.url || null }))
        }

        // Fees (array or object)
        let fees = source.fees || source?.content?.fees
        if (fees && !Array.isArray(fees) && typeof fees === 'object') {
          // Convert keyed object {registration: {amount, currency}, ...}
            fees = Object.values(fees)
        }

        // Processing time: may be at root or inside content
        let processingTime = source.processingTime || source?.content?.processingTime
        if (processingTime && typeof processingTime === 'number') {
          processingTime = { minDays: processingTime, maxDays: processingTime }
        }

        const title = source.title || source.name || source?.content?.title
        const summary = source.summary || source.description || source?.content?.summary || source?.content?.result

        return { ...source, id, steps, documentsRequired, fees, processingTime, title, name: source.name, summary }
      },
      providesTags: (result, error, id) => [{ type: 'Procedure', id }]
    })
  })
})

export const { useListProceduresQuery, useGetProcedureQuery } = proceduresApi
