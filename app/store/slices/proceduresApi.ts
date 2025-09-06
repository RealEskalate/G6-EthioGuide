import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'
import type { Procedure } from '@/app/types/procedure'

interface BackendPagination { page?: number; limit?: number; total?: number }
interface Paginated<T> { data?: T[]; items?: T[]; results?: T[]; hasNext?: boolean; pagination?: BackendPagination }

export interface ProcedureListArgs {
  page?: number;
  limit?: number;
  name?: string; // search by name (backend expects 'name')
  organizationID?: string;
  groupID?: string;
  minProcessingDays?: number;
  maxProcessingDays?: number;
  sortBy?: string; // createdAt, fee, processingTime
  sortOrder?: 'ASC' | 'DESC';
}

export interface ProcedureListResponse {
  list: Procedure[];
  page: number;
  limit: number;
  total: number;
  totalPages: number;
  hasNext: boolean;
}

export const proceduresApi = createApi({
  reducerPath: 'proceduresApi',
  // Use relative path so Next.js rewrite (in next.config.ts) proxies to external backend server-side (avoids CORS)
  baseQuery: fetchBaseQuery({ baseUrl: '/api/v1' }),
  tagTypes: ['Procedure', 'Procedures'],
  endpoints: (builder) => ({
    listProcedures: builder.query<ProcedureListResponse, ProcedureListArgs | void>({
      query: (args) => {
        const params = new URLSearchParams()
        if (args?.page) params.set('page', String(args.page))
        if (args?.limit) params.set('limit', String(args.limit))
        if (args?.name) params.set('name', args.name)
        if (args?.organizationID) params.set('organizationID', args.organizationID)
        if (args?.groupID) params.set('groupID', args.groupID)
        if (typeof args?.minProcessingDays === 'number') params.set('minProcessingDays', String(args.minProcessingDays))
        if (typeof args?.maxProcessingDays === 'number') params.set('maxProcessingDays', String(args.maxProcessingDays))
        if (args?.sortBy) params.set('sortBy', args.sortBy)
        if (args?.sortOrder) params.set('sortOrder', args.sortOrder)
        return `/procedures${params.size ? `?${params.toString()}` : ''}`
      },
      transformResponse: (raw: Paginated<Procedure> | Procedure[]): ProcedureListResponse => {
        const arr = (raw as Paginated<Procedure>).data || (raw as Paginated<Procedure>).items || (raw as Paginated<Procedure>).results || (Array.isArray(raw) ? raw as Procedure[] : [])
        const pagination = (raw as Paginated<Procedure>).pagination || {}
        const page = pagination.page ?? 1
  const limit = (pagination.limit ?? arr.length) || 10
        const total = pagination.total ?? ((raw as Paginated<Procedure>).hasNext ? page * limit + 1 : arr.length)
        const totalPages = total && limit ? Math.max(1, Math.ceil(total / limit)) : 1
        // Determine hasNext: prefer explicit flag, else compute
        const hasNextExplicit = (raw as Paginated<Procedure>).hasNext
        const hasNext = typeof hasNextExplicit === 'boolean' ? hasNextExplicit : (page < totalPages)
        return {
          list: arr.map(p => {
            const anyP: any = p
            const id = anyP.id || anyP._id || anyP.procedureId || anyP.procedure_id || anyP.slug || anyP.code || null
            const contentBlock = anyP.content || anyP.Content || {}

            // Documents normalization
            let documentsRequired: any =
              anyP.documentsRequired ||
              anyP.requiredDocuments ||
              anyP.RequiredDocuments ||
              anyP.documentRequirements ||
              anyP.DocumentRequirements ||
              anyP.documents ||
              anyP.docs ||
              anyP?.content?.documentsRequired ||
              (anyP?.content as any)?.requiredDocuments ||
              anyP?.content?.documents ||
              (anyP?.content as any)?.documents_required ||
              contentBlock?.Documents ||
              (contentBlock as any)?.RequiredDocuments ||
              (contentBlock as any)?.DocumentsRequired ||
              (contentBlock as any)?.documents ||
              (contentBlock as any)?.requiredDocuments ||
              (contentBlock as any)?.documents_required
            if (!documentsRequired && (contentBlock?.Prerequisites || (contentBlock as any)?.prerequisites)) {
              const preArr = (contentBlock as any).Prerequisites || (contentBlock as any).prerequisites
              if (Array.isArray(preArr)) {
                documentsRequired = preArr.map((d: any) => ({ name: typeof d === 'string' ? d : (d.name || d.title || 'Document'), templateUrl: null }))
              }
            }
            if (Array.isArray(documentsRequired)) {
              documentsRequired = documentsRequired.map((d: any) => ({ name: d.name || d.title || d.label || d.filename || 'Document', templateUrl: d.templateUrl || d.url || null }))
            }

            // Fees normalization (array or object or alternate keys)
            let fees: any =
              anyP.fees ||
              anyP.Fees ||
              (anyP as any).costs ||
              (anyP as any).Costs ||
              anyP?.content?.fees ||
              (anyP?.content as any)?.costs ||
              contentBlock?.Fees ||
              (contentBlock as any)?.fees ||
              (contentBlock as any)?.Costs ||
              (contentBlock as any)?.costs
            if (fees && !Array.isArray(fees) && typeof fees === 'object') {
              const hasAmount = 'Amount' in fees || 'amount' in fees
              if (hasAmount) {
                fees = [{
                  amount: (fees as any).Amount ?? (fees as any).amount ?? 0,
                  currency: (fees as any).Currency ?? (fees as any).currency ?? '',
                  label: (fees as any).Label ?? (fees as any).label ?? 'Fee'
                }]
              } else {
                fees = Object.values(fees)
              }
            }

            // Processing time normalization (supports snake/case variants)
            let processingTime: any =
              anyP.processingTime ||
              (anyP as any).processing_time ||
              anyP.ProcessingTime ||
              anyP?.content?.processingTime ||
              (anyP?.content as any)?.processing_time ||
              contentBlock?.ProcessingTime ||
              (contentBlock as any)?.processingTime ||
              (contentBlock as any)?.processing_time
            if (processingTime && typeof processingTime === 'number') {
              processingTime = { minDays: processingTime, maxDays: processingTime }
            }
            if (processingTime && typeof processingTime === 'object') {
              if ('MinDays' in processingTime || 'MaxDays' in processingTime) {
                processingTime = {
                  minDays: (processingTime as any).MinDays ?? (processingTime as any).minDays ?? null,
                  maxDays: (processingTime as any).MaxDays ?? (processingTime as any).maxDays ?? null
                }
              } else if ('min_days' in (processingTime as any) || 'max_days' in (processingTime as any)) {
                processingTime = {
                  minDays: (processingTime as any).min_days ?? (processingTime as any).minDays ?? null,
                  maxDays: (processingTime as any).max_days ?? (processingTime as any).maxDays ?? null
                }
              }
            }

            const title = anyP.title || anyP.Title || anyP.name || anyP.Name || anyP?.content?.title || contentBlock?.Title || (contentBlock as any)?.Result
            const summary = anyP.summary || anyP.Summary || anyP.description || anyP.Description || anyP?.content?.summary || anyP?.content?.result || contentBlock?.Summary || (contentBlock as any)?.Result

            return { ...anyP, id, documentsRequired, fees, processingTime, title, name: (anyP.name || anyP.Name), summary }
          }).filter(p => p.id),
          page,
          limit,
            total,
          totalPages,
          hasNext
        }
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
        if (raw?.data?.procedure) source = raw.data.procedure
        if (Array.isArray(source)) source = source[0] || {}

        // Normalize capitalized backend keys to expected camelCase
        // (non-mutating: just read from both variants)
        const id = source?.id || source?.ID || source?._id || source?.procedureId || source?.procedure_id || source?.uuid || source?.slug || source?.code || undefined

        // Extract capitalized Content structure if present
        const contentBlock = source.content || source.Content || {}

        // Build steps
        let steps = source.steps
        if (!steps && (source?.content?.steps || contentBlock?.Steps)) {
          const rawSteps = source.content?.steps || contentBlock.Steps
          if (Array.isArray(rawSteps)) {
            steps = rawSteps.map((s: any, i: number) => typeof s === 'string' ? ({ order: i + 1, text: s }) : ({ order: s.order || i + 1, text: s.text || JSON.stringify(s) }))
          } else if (typeof rawSteps === 'object') {
            steps = Object.entries(rawSteps).map(([k, v]: any, idx) => ({ order: Number(k) || idx + 1, text: typeof v === 'string' ? v : JSON.stringify(v) }))
          }
        }

        // Documents (support many backend variants)
        let documentsRequired =
          source.documentsRequired ||
          source.requiredDocuments ||
          source.RequiredDocuments ||
          source.documentRequirements ||
          source.DocumentRequirements ||
          source.documents ||
          source.docs ||
          source?.content?.documentsRequired ||
          source?.content?.requiredDocuments ||
          source?.content?.documents ||
          (source?.content as any)?.documents_required ||
          contentBlock?.Documents ||
          contentBlock?.RequiredDocuments ||
          contentBlock?.DocumentsRequired ||
          (contentBlock as any)?.documents ||
          (contentBlock as any)?.requiredDocuments ||
          (contentBlock as any)?.documents_required
        // Map prerequisites (capitalized) as documents if no docs provided
        if (!documentsRequired && (contentBlock?.Prerequisites || contentBlock?.prerequisites)) {
          const preArr = contentBlock.Prerequisites || contentBlock.prerequisites
          if (Array.isArray(preArr)) {
            documentsRequired = preArr.map((d: any) => ({ name: typeof d === 'string' ? d : (d.name || d.title || 'Document'), templateUrl: null }))
          }
        }
        if (Array.isArray(documentsRequired)) {
          documentsRequired = documentsRequired.map((d: any) => ({ name: d.name || d.title || d.label || d.filename || 'Document', templateUrl: d.templateUrl || d.url || null }))
        }

        // Fees (array or object, support alternate keys)
        let fees: any =
          source.fees ||
          source.Fees ||
          (source as any).costs ||
          (source as any).Costs ||
          source?.content?.fees ||
          (source?.content as any)?.costs ||
          contentBlock?.Fees ||
          (contentBlock as any)?.fees ||
          (contentBlock as any)?.Costs ||
          (contentBlock as any)?.costs
        if (fees && !Array.isArray(fees) && typeof fees === 'object') {
          // Single object with Amount / Currency / Label (capitalized) => wrap into array
          const keys = Object.keys(fees)
          const hasAmount = 'Amount' in fees || 'amount' in fees
          if (hasAmount && ('Currency' in fees || 'currency' in fees)) {
            const single = {
              amount: (fees as any).Amount ?? (fees as any).amount ?? 0,
              currency: (fees as any).Currency ?? (fees as any).currency ?? '',
              label: (fees as any).Label ?? (fees as any).label ?? 'Fee'
            }
            fees = [single]
          } else {
            // Convert keyed object {registration: {amount, currency}, ...}
            fees = Object.values(fees)
          }
        }

        // Processing time: may be at root or inside content (handle snake_case too)
        let processingTime: any =
          source.processingTime ||
          (source as any).processing_time ||
          source.ProcessingTime ||
          source?.content?.processingTime ||
          (source?.content as any)?.processing_time ||
          contentBlock?.ProcessingTime ||
          (contentBlock as any)?.processingTime ||
          (contentBlock as any)?.processing_time
        if (processingTime && typeof processingTime === 'number') {
          processingTime = { minDays: processingTime, maxDays: processingTime }
        }
        if (processingTime && typeof processingTime === 'object') {
          if ('MinDays' in processingTime || 'MaxDays' in processingTime) {
            processingTime = {
              minDays: (processingTime as any).MinDays ?? (processingTime as any).minDays ?? null,
              maxDays: (processingTime as any).MaxDays ?? (processingTime as any).maxDays ?? null
            }
          } else if ('min_days' in (processingTime as any) || 'max_days' in (processingTime as any)) {
            processingTime = {
              minDays: (processingTime as any).min_days ?? (processingTime as any).minDays ?? null,
              maxDays: (processingTime as any).max_days ?? (processingTime as any).maxDays ?? null
            }
          }
        }

        const title = source.title || source.Title || source.name || source.Name || source?.content?.title || contentBlock?.Title || contentBlock?.Result
        const summary = source.summary || source.Summary || source.description || source.Description || source?.content?.summary || source?.content?.result || contentBlock?.Summary || contentBlock?.Result

        return { ...source, id, steps, documentsRequired, fees, processingTime, title, name: (source.name || source.Name), summary, __raw: raw }
      },
      providesTags: (result, error, id) => [{ type: 'Procedure', id }]
    })
    ,
    // Fallback query that attempts multiple backend path variants if the standard one returns nothing
    getProcedureFlexible: builder.query<Procedure, string>({
      // We'll implement custom logic via queryFn to try several endpoints
      // Order: /procedures/{id} -> /procedure/{id} -> /procedures?id={id}
      async queryFn(id, _api, _extra, baseFetch) {
        const paths = [`/procedures/${id}`, `/procedure/${id}`, `/procedures?id=${encodeURIComponent(id)}`]
        let lastError: any = null
        for (const p of paths) {
          const res: any = await baseFetch(p)
          if (res.error) { lastError = res.error; continue }
          const raw = res.data
          // Heuristic: if raw is empty object or array length 0, continue
          if (!raw || (Array.isArray(raw) && raw.length === 0) || (typeof raw === 'object' && !Array.isArray(raw) && Object.keys(raw).length === 0)) {
            continue
          }
          // Reuse existing transform logic by mimicking previous function
          let source: any = raw
          if (raw?.data && typeof raw.data === 'object' && !Array.isArray(raw.data)) source = raw.data
          if (raw?.procedure) source = raw.procedure
          if (raw?.data?.procedure) source = raw.data.procedure
          if (Array.isArray(source)) source = source[0] || {}
          const contentBlock = source.content || source.Content || {}
          const resolvedId = source?.id || source?.ID || source?._id || source?.procedureId || source?.procedure_id || source?.uuid || source?.slug || source?.code || id
          // Steps normalization (copy from main transform)
          let steps = source.steps
          if (!steps && (source?.content?.steps || contentBlock?.Steps)) {
            const rawSteps = source.content?.steps || contentBlock?.Steps
            if (Array.isArray(rawSteps)) {
              steps = rawSteps.map((s: any, i: number) => typeof s === 'string' ? ({ order: i + 1, text: s }) : ({ order: s.order || i + 1, text: s.text || JSON.stringify(s) }))
            } else if (typeof rawSteps === 'object') {
              steps = Object.entries(rawSteps).map(([k, v]: any, idx) => ({ order: Number(k) || idx + 1, text: typeof v === 'string' ? v : JSON.stringify(v) }))
            }
          }
          let documentsRequired: any =
            source.documentsRequired ||
            (source as any).requiredDocuments ||
            (source as any).RequiredDocuments ||
            (source as any).documentRequirements ||
            (source as any).DocumentRequirements ||
            source.documents ||
            source.docs ||
            source?.content?.documentsRequired ||
            (source?.content as any)?.requiredDocuments ||
            source?.content?.documents ||
            (source?.content as any)?.documents_required ||
            contentBlock?.Documents ||
            (contentBlock as any)?.RequiredDocuments ||
            (contentBlock as any)?.DocumentsRequired ||
            (contentBlock as any)?.documents ||
            (contentBlock as any)?.requiredDocuments ||
            (contentBlock as any)?.documents_required
          if (!documentsRequired && (contentBlock?.Prerequisites || contentBlock?.prerequisites)) {
            const preArr = contentBlock.Prerequisites || contentBlock.prerequisites
            if (Array.isArray(preArr)) {
              documentsRequired = preArr.map((d: any) => ({ name: typeof d === 'string' ? d : (d.name || d.title || 'Document'), templateUrl: null }))
            }
          }
          if (Array.isArray(documentsRequired)) {
            documentsRequired = documentsRequired.map((d: any) => ({ name: d.name || d.title || d.label || d.filename || 'Document', templateUrl: d.templateUrl || d.url || null }))
          }
          let fees: any =
            source.fees ||
            source.Fees ||
            (source as any).costs ||
            (source as any).Costs ||
            source?.content?.fees ||
            (source?.content as any)?.costs ||
            contentBlock?.Fees ||
            (contentBlock as any)?.fees ||
            (contentBlock as any)?.Costs ||
            (contentBlock as any)?.costs
          if (fees && !Array.isArray(fees) && typeof fees === 'object') {
            const hasAmount = 'Amount' in fees || 'amount' in fees
            if (hasAmount) {
              fees = [{
                amount: (fees as any).Amount ?? (fees as any).amount ?? 0,
                currency: (fees as any).Currency ?? (fees as any).currency ?? '',
                label: (fees as any).Label ?? (fees as any).label ?? 'Fee'
              }]
            } else {
              fees = Object.values(fees)
            }
          }
          let processingTime: any =
            source.processingTime ||
            (source as any).processing_time ||
            source.ProcessingTime ||
            source?.content?.processingTime ||
            (source?.content as any)?.processing_time ||
            contentBlock?.ProcessingTime ||
            (contentBlock as any)?.processingTime ||
            (contentBlock as any)?.processing_time
          if (processingTime && typeof processingTime === 'number') {
            processingTime = { minDays: processingTime, maxDays: processingTime }
          }
          if (processingTime && typeof processingTime === 'object' && ('MinDays' in processingTime || 'MaxDays' in processingTime)) {
            processingTime = { minDays: (processingTime as any).MinDays ?? (processingTime as any).minDays ?? null, maxDays: (processingTime as any).MaxDays ?? (processingTime as any).maxDays ?? null }
          } else if (processingTime && typeof processingTime === 'object' && ('min_days' in (processingTime as any) || 'max_days' in (processingTime as any))) {
            processingTime = { minDays: (processingTime as any).min_days ?? (processingTime as any).minDays ?? null, maxDays: (processingTime as any).max_days ?? (processingTime as any).maxDays ?? null }
          }
          const title = source.title || source.Title || source.name || source.Name || source?.content?.title || contentBlock?.Title || contentBlock?.Result
          const summary = source.summary || source.Summary || source.description || source.Description || source?.content?.summary || source?.content?.result || contentBlock?.Summary || contentBlock?.Result
          const procedure: Procedure = { ...source, id: resolvedId, steps, documentsRequired, fees, processingTime, title, name: (source.name || source.Name), summary, __raw: raw }
          return { data: procedure }
        }
        return { error: lastError || { status: 404, data: { message: 'Procedure not found via any fallback path' } } }
      },
      providesTags: (result, error, id) => [{ type: 'Procedure', id }]
    })
  })
})

export const { useListProceduresQuery, useGetProcedureQuery, useGetProcedureFlexibleQuery } = proceduresApi
