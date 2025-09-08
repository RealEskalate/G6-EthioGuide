import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'
import type { FetchBaseQueryError, FetchArgs, QueryReturnValue } from '@reduxjs/toolkit/query'
import type { Procedure } from '@/app/types/procedure'

interface BackendPagination { page?: number; limit?: number; total?: number }
interface Paginated<T> { data?: T[]; items?: T[]; results?: T[]; hasNext?: boolean; pagination?: BackendPagination }

type AnyRecord = Record<string, unknown>
const isObject = (v: unknown): v is AnyRecord => typeof v === 'object' && v !== null
const get = <T = unknown>(o: unknown, key: string): T | undefined => (isObject(o) ? (o[key] as T) : undefined)
const asString = (v: unknown, fallback = ''): string => (typeof v === 'string' ? v : v == null ? fallback : String(v))
const toNum = (v: unknown, fallback = 0): number => (typeof v === 'number' ? v : typeof v === 'string' && v !== '' ? Number(v) : fallback)
// const asArray = (v: unknown) => Array.isArray(v) ? v : []

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
  baseQuery: fetchBaseQuery({ baseUrl: '/api/v1',
    prepareHeaders: (headers) => {
      if (typeof window !== 'undefined') {
        headers.set('lang', localStorage.getItem('i18nextLng') || 'en')
      }
      return headers;
    }

   }),
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
        const rawObj = raw as unknown
        const arr = (get<Procedure[]>(rawObj, 'data') || get<Procedure[]>(rawObj, 'items') || get<Procedure[]>(rawObj, 'results') || (Array.isArray(raw) ? (raw as Procedure[]) : []))
        const pagination = get<BackendPagination>(rawObj, 'pagination') || {}
        const page = pagination.page ?? 1
  const limit = (pagination.limit ?? arr.length) || 10
        const total = pagination.total ?? (get<boolean>(rawObj, 'hasNext') ? page * limit + 1 : arr.length)
        const totalPages = total && limit ? Math.max(1, Math.ceil(total / limit)) : 1
        // Determine hasNext: prefer explicit flag, else compute
        const hasNextExplicit = get<boolean>(rawObj, 'hasNext')
        const hasNext = typeof hasNextExplicit === 'boolean' ? hasNextExplicit : (page < totalPages)
        return {
          list: arr.map(p => {
            const obj: AnyRecord = (p as unknown as AnyRecord) || {}
            const id = (obj.id || obj._id || obj.procedureId || obj.procedure_id || obj.slug || obj.code || null) as string | null
            const contentBlock: AnyRecord = (obj.content as AnyRecord) || (obj.Content as AnyRecord) || {}

            // Documents normalization
            let documentsRequired: unknown =
              obj['documentsRequired'] ||
              obj['requiredDocuments'] ||
              obj['RequiredDocuments'] ||
              obj['documentRequirements'] ||
              obj['DocumentRequirements'] ||
              obj['documents'] ||
              obj['docs'] ||
              get(obj['content'], 'documentsRequired') ||
              get(obj['content'], 'requiredDocuments') ||
              get(obj['content'], 'documents') ||
              get(obj['content'], 'documents_required') ||
              contentBlock['Documents'] ||
              contentBlock['RequiredDocuments'] ||
              contentBlock['DocumentsRequired'] ||
              contentBlock['documents'] ||
              contentBlock['requiredDocuments'] ||
              contentBlock['documents_required']
            if (!documentsRequired && (contentBlock['Prerequisites'] || contentBlock['prerequisites'])) {
              const preArr = (contentBlock['Prerequisites'] || contentBlock['prerequisites']) as unknown
              if (Array.isArray(preArr)) {
                documentsRequired = (preArr as unknown[]).map((d) => {
                  const dobj = isObject(d) ? (d as AnyRecord) : undefined
                  const name = dobj ? (dobj['name'] || dobj['title'] || 'Document') : asString(d)
                  return { name: asString(name), templateUrl: null }
                })
              }
            }
            if (Array.isArray(documentsRequired)) {
              documentsRequired = (documentsRequired as unknown[]).map((d) => {
                const dobj = isObject(d) ? (d as AnyRecord) : {}
                const name = dobj['name'] || dobj['title'] || dobj['label'] || dobj['filename'] || 'Document'
                const templateUrl = dobj['templateUrl'] || dobj['url'] || null
                return { name: asString(name), templateUrl: (templateUrl as string | null) }
              })
            }

            // Fees normalization (array or object or alternate keys)
            let fees: unknown =
              obj['fees'] ||
              obj['Fees'] ||
              (obj as AnyRecord)['costs'] ||
              (obj as AnyRecord)['Costs'] ||
              get(obj['content'], 'fees') ||
              get(obj['content'], 'costs') ||
              contentBlock['Fees'] ||
              contentBlock['fees'] ||
              contentBlock['Costs'] ||
              contentBlock['costs']
            if (fees && !Array.isArray(fees) && typeof fees === 'object') {
              const feesObj = fees as AnyRecord
              const hasAmount = 'Amount' in feesObj || 'amount' in feesObj
              if (hasAmount) {
                fees = [{
                  amount: toNum(feesObj['Amount'] ?? feesObj['amount'] ?? 0, 0),
                  currency: asString(feesObj['Currency'] ?? feesObj['currency'] ?? ''),
                  label: asString(feesObj['Label'] ?? feesObj['label'] ?? 'Fee')
                }]
              } else {
                fees = Object.values(feesObj)
              }
            }

            // Processing time normalization (supports snake/case variants)
            let processingTime: unknown =
              obj['processingTime'] ||
              (obj as AnyRecord)['processing_time'] ||
              obj['ProcessingTime'] ||
              get(obj['content'], 'processingTime') ||
              get(obj['content'], 'processing_time') ||
              contentBlock['ProcessingTime'] ||
              contentBlock['processingTime'] ||
              contentBlock['processing_time']
            if (processingTime && typeof processingTime === 'number') {
              processingTime = { minDays: processingTime, maxDays: processingTime }
            }
            if (processingTime && typeof processingTime === 'object') {
              const pt = processingTime as AnyRecord
              if ('MinDays' in pt || 'MaxDays' in pt) {
                processingTime = {
                  minDays: toNum(pt['MinDays'] ?? pt['minDays'] ?? null as unknown as number, undefined as unknown as number),
                  maxDays: toNum(pt['MaxDays'] ?? pt['maxDays'] ?? null as unknown as number, undefined as unknown as number)
                }
              } else if ('min_days' in pt || 'max_days' in pt) {
                processingTime = {
                  minDays: toNum(pt['min_days'] ?? pt['minDays'] ?? null as unknown as number, undefined as unknown as number),
                  maxDays: toNum(pt['max_days'] ?? pt['maxDays'] ?? null as unknown as number, undefined as unknown as number)
                }
              }
            }

            const title = obj['title'] || obj['Title'] || obj['name'] || obj['Name'] || get(obj['content'], 'title') || contentBlock['Title'] || contentBlock['Result']
            const summary = obj['summary'] || obj['Summary'] || obj['description'] || obj['Description'] || get(obj['content'], 'summary') || get(obj['content'], 'result') || contentBlock['Summary'] || contentBlock['Result']

            const base: AnyRecord = isObject(p) ? (p as AnyRecord) : {}
            return { ...(base as object), id: id as string, documentsRequired: documentsRequired as unknown[], fees: fees as unknown[], processingTime: processingTime as unknown, title: asString(title), name: asString(base['name'] ?? base['Name'] ?? ''), summary: asString(summary) } as unknown as Procedure
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
      transformResponse: (raw: unknown): Procedure => {
        // Accept various envelope shapes
        let source: unknown = raw
        const rawObj = raw as AnyRecord
        if (isObject(raw) && rawObj?.data && typeof rawObj.data === 'object' && !Array.isArray(rawObj.data)) source = rawObj.data
        if (isObject(raw) && 'procedure' in rawObj) source = rawObj.procedure
        if (isObject(get(raw, 'data')) && isObject((get(raw, 'data') as AnyRecord)['procedure'])) source = (get(raw, 'data') as AnyRecord)['procedure']
        if (Array.isArray(source)) source = (source as unknown[])[0] || {}

        // Normalize capitalized backend keys to expected camelCase
        // (non-mutating: just read from both variants)
        const sObj = (isObject(source) ? (source as AnyRecord) : {})
        const id = sObj['id'] || sObj['ID'] || sObj['_id'] || sObj['procedureId'] || sObj['procedure_id'] || sObj['uuid'] || sObj['slug'] || sObj['code'] || undefined

        // Extract capitalized Content structure if present
        const contentBlock: AnyRecord = (sObj['content'] as AnyRecord) || (sObj['Content'] as AnyRecord) || {}

        // Build steps
        let steps: unknown = sObj['steps']
        const contentSteps = get<unknown>(sObj['content'], 'steps') || contentBlock['Steps']
        if (!steps && contentSteps) {
          const rawSteps = contentSteps
          if (Array.isArray(rawSteps)) {
            steps = rawSteps.map((s, i) => {
              const sobj = isObject(s) ? (s as AnyRecord) : undefined
              return typeof s === 'string' ? ({ order: i + 1, text: s }) : ({ order: toNum(sobj?.['order'] ?? i + 1, i + 1), text: asString(sobj?.['text'] ?? JSON.stringify(s)) })
            })
          } else if (typeof rawSteps === 'object') {
            const entries = Object.entries(rawSteps as Record<string, unknown>)
            steps = entries.map(([k, v], idx) => ({ order: Number(k) || idx + 1, text: typeof v === 'string' ? v : JSON.stringify(v) }))
          }
        }

        // Documents (support many backend variants)
        let documentsRequired: unknown =
          sObj['documentsRequired'] ||
          sObj['requiredDocuments'] ||
          sObj['RequiredDocuments'] ||
          sObj['documentRequirements'] ||
          sObj['DocumentRequirements'] ||
          sObj['documents'] ||
          sObj['docs'] ||
          get(sObj['content'], 'documentsRequired') ||
          get(sObj['content'], 'requiredDocuments') ||
          get(sObj['content'], 'documents') ||
          get(sObj['content'], 'documents_required') ||
          contentBlock['Documents'] ||
          contentBlock['RequiredDocuments'] ||
          contentBlock['DocumentsRequired'] ||
          contentBlock['documents'] ||
          contentBlock['requiredDocuments'] ||
          contentBlock['documents_required']
        // Map prerequisites (capitalized) as documents if no docs provided
        if (!documentsRequired && (contentBlock['Prerequisites'] || contentBlock['prerequisites'])) {
          const preArr = (contentBlock['Prerequisites'] || contentBlock['prerequisites']) as unknown
          if (Array.isArray(preArr)) {
            documentsRequired = (preArr as unknown[]).map((d) => {
              const dobj = isObject(d) ? (d as AnyRecord) : undefined
              const name = dobj ? (dobj['name'] || dobj['title'] || 'Document') : asString(d)
              return { name: asString(name), templateUrl: null }
            })
          }
        }
        if (Array.isArray(documentsRequired)) {
          documentsRequired = (documentsRequired as unknown[]).map((d) => {
            const dobj = isObject(d) ? (d as AnyRecord) : {}
            const name = dobj['name'] || dobj['title'] || dobj['label'] || dobj['filename'] || 'Document'
            const templateUrl = dobj['templateUrl'] || dobj['url'] || null
            return { name: asString(name), templateUrl: (templateUrl as string | null) }
          })
        }

        // Fees (array or object, support alternate keys)
        let fees: unknown =
          sObj['fees'] ||
          sObj['Fees'] ||
          sObj['costs'] ||
          sObj['Costs'] ||
          get(sObj['content'], 'fees') ||
          get(sObj['content'], 'costs') ||
          contentBlock['Fees'] ||
          contentBlock['fees'] ||
          contentBlock['Costs'] ||
          contentBlock['costs']
        if (fees && !Array.isArray(fees) && typeof fees === 'object') {
          // Single object with Amount / Currency / Label (capitalized) => wrap into array
          const fobj = fees as AnyRecord
          const hasAmount = 'Amount' in fobj || 'amount' in fobj
          if (hasAmount && ('Currency' in fobj || 'currency' in fobj)) {
            const single = {
              amount: toNum(fobj['Amount'] ?? fobj['amount'] ?? 0, 0),
              currency: asString(fobj['Currency'] ?? fobj['currency'] ?? ''),
              label: asString(fobj['Label'] ?? fobj['label'] ?? 'Fee')
            }
            fees = [single]
          } else {
            // Convert keyed object {registration: {amount, currency}, ...}
            fees = Object.values(fobj)
          }
        }

        // Processing time: may be at root or inside content (handle snake_case too)
        let processingTime: unknown =
          sObj['processingTime'] ||
          sObj['processing_time'] ||
          sObj['ProcessingTime'] ||
          get(sObj['content'], 'processingTime') ||
          get(sObj['content'], 'processing_time') ||
          contentBlock['ProcessingTime'] ||
          contentBlock['processingTime'] ||
          contentBlock['processing_time']
        if (processingTime && typeof processingTime === 'number') {
          processingTime = { minDays: processingTime, maxDays: processingTime }
        }
        if (processingTime && typeof processingTime === 'object') {
          const pt = processingTime as AnyRecord
          if ('MinDays' in pt || 'MaxDays' in pt) {
            processingTime = {
              minDays: toNum(pt['MinDays'] ?? pt['minDays'] ?? null as unknown as number, undefined as unknown as number),
              maxDays: toNum(pt['MaxDays'] ?? pt['maxDays'] ?? null as unknown as number, undefined as unknown as number)
            }
          } else if ('min_days' in pt || 'max_days' in pt) {
            processingTime = {
              minDays: toNum(pt['min_days'] ?? pt['minDays'] ?? null as unknown as number, undefined as unknown as number),
              maxDays: toNum(pt['max_days'] ?? pt['maxDays'] ?? null as unknown as number, undefined as unknown as number)
            }
          }
        }

        const title = sObj['title'] || sObj['Title'] || sObj['name'] || sObj['Name'] || get(sObj['content'], 'title') || contentBlock['Title'] || contentBlock['Result']
        const summary = sObj['summary'] || sObj['Summary'] || sObj['description'] || sObj['Description'] || get(sObj['content'], 'summary') || get(sObj['content'], 'result') || contentBlock['Summary'] || contentBlock['Result']

        const base: AnyRecord = isObject(source) ? (source as AnyRecord) : {}
  return { ...(base as object), id: id as string, steps: steps as unknown as Procedure['steps'], documentsRequired: documentsRequired as unknown as Procedure['documentsRequired'], fees: fees as unknown as Procedure['fees'], processingTime: processingTime as unknown as Procedure['processingTime'], title: asString(title), name: asString(base['name'] ?? base['Name'] ?? ''), summary: asString(summary) } as unknown as Procedure
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
        let lastError: FetchBaseQueryError | undefined
        for (const p of paths) {
          const res = (await baseFetch(p as unknown as FetchArgs)) as QueryReturnValue<unknown, FetchBaseQueryError>
          if (res.error) { lastError = res.error; continue }
          const raw = res.data
          // Heuristic: if raw is empty object or array length 0, continue
          if (!raw || (Array.isArray(raw) && raw.length === 0) || (typeof raw === 'object' && !Array.isArray(raw) && Object.keys(raw as AnyRecord).length === 0)) {
            continue
          }
          // Reuse existing transform logic by mimicking previous function
          let source: unknown = raw
          const rawInner = raw as AnyRecord
          if (isObject(raw) && rawInner?.data && typeof rawInner.data === 'object' && !Array.isArray(rawInner.data)) source = rawInner.data
          if (isObject(raw) && 'procedure' in rawInner) source = rawInner.procedure
          if (isObject(get(raw, 'data')) && isObject((get(raw, 'data') as AnyRecord)['procedure'])) source = (get(raw, 'data') as AnyRecord)['procedure']
          if (Array.isArray(source)) source = (source as unknown[])[0] || {}
          const sObj = isObject(source) ? (source as AnyRecord) : {}
          const contentBlock: AnyRecord = (sObj['content'] as AnyRecord) || (sObj['Content'] as AnyRecord) || {}
          const resolvedId = sObj['id'] || sObj['ID'] || sObj['_id'] || sObj['procedureId'] || sObj['procedure_id'] || sObj['uuid'] || sObj['slug'] || sObj['code'] || id
          // Steps normalization (copy from main transform)
          let steps: unknown = sObj['steps']
          const cSteps = get<unknown>(sObj['content'], 'steps') || contentBlock['Steps']
          if (!steps && cSteps) {
            const rawSteps = cSteps
            if (Array.isArray(rawSteps)) {
              steps = rawSteps.map((s, i) => {
                const sobj = isObject(s) ? (s as AnyRecord) : undefined
                return typeof s === 'string' ? ({ order: i + 1, text: s }) : ({ order: toNum(sobj?.['order'] ?? i + 1, i + 1), text: asString(sobj?.['text'] ?? JSON.stringify(s)) })
              })
            } else if (typeof rawSteps === 'object') {
              const entries = Object.entries(rawSteps as Record<string, unknown>)
              steps = entries.map(([k, v], idx) => ({ order: Number(k) || idx + 1, text: typeof v === 'string' ? v : JSON.stringify(v) }))
            }
          }
          let documentsRequired: unknown =
            sObj['documentsRequired'] ||
            sObj['requiredDocuments'] ||
            sObj['RequiredDocuments'] ||
            sObj['documentRequirements'] ||
            sObj['DocumentRequirements'] ||
            sObj['documents'] ||
            sObj['docs'] ||
            get(sObj['content'], 'documentsRequired') ||
            get(sObj['content'], 'requiredDocuments') ||
            get(sObj['content'], 'documents') ||
            get(sObj['content'], 'documents_required') ||
            contentBlock['Documents'] ||
            contentBlock['RequiredDocuments'] ||
            contentBlock['DocumentsRequired'] ||
            contentBlock['documents'] ||
            contentBlock['requiredDocuments'] ||
            contentBlock['documents_required']
          if (!documentsRequired && (contentBlock['Prerequisites'] || contentBlock['prerequisites'])) {
            const preArr = (contentBlock['Prerequisites'] || contentBlock['prerequisites']) as unknown
            if (Array.isArray(preArr)) {
              documentsRequired = (preArr as unknown[]).map((d) => {
                const dobj = isObject(d) ? (d as AnyRecord) : undefined
                const name = dobj ? (dobj['name'] || dobj['title'] || 'Document') : asString(d)
                return { name: asString(name), templateUrl: null }
              })
            }
          }
          if (Array.isArray(documentsRequired)) {
            documentsRequired = (documentsRequired as unknown[]).map((d) => {
              const dobj = isObject(d) ? (d as AnyRecord) : {}
              const name = dobj['name'] || dobj['title'] || dobj['label'] || dobj['filename'] || 'Document'
              const templateUrl = dobj['templateUrl'] || dobj['url'] || null
              return { name: asString(name), templateUrl: (templateUrl as string | null) }
            })
          }
          let fees: unknown =
            sObj['fees'] ||
            sObj['Fees'] ||
            sObj['costs'] ||
            sObj['Costs'] ||
            get(sObj['content'], 'fees') ||
            get(sObj['content'], 'costs') ||
            contentBlock['Fees'] ||
            contentBlock['fees'] ||
            contentBlock['Costs'] ||
            contentBlock['costs']
          if (fees && !Array.isArray(fees) && typeof fees === 'object') {
            const fobj = fees as AnyRecord
            const hasAmount = 'Amount' in fobj || 'amount' in fobj
            if (hasAmount) {
              fees = [{
                amount: toNum(fobj['Amount'] ?? fobj['amount'] ?? 0, 0),
                currency: asString(fobj['Currency'] ?? fobj['currency'] ?? ''),
                label: asString(fobj['Label'] ?? fobj['label'] ?? 'Fee')
              }]
            } else {
              fees = Object.values(fobj)
            }
          }
          let processingTime: unknown =
            sObj['processingTime'] ||
            sObj['processing_time'] ||
            sObj['ProcessingTime'] ||
            get(sObj['content'], 'processingTime') ||
            get(sObj['content'], 'processing_time') ||
            contentBlock['ProcessingTime'] ||
            contentBlock['processingTime'] ||
            contentBlock['processing_time']
          if (processingTime && typeof processingTime === 'number') {
            processingTime = { minDays: processingTime, maxDays: processingTime }
          }
          if (processingTime && typeof processingTime === 'object') {
            const pt = processingTime as AnyRecord
            if ('MinDays' in pt || 'MaxDays' in pt) {
              processingTime = { minDays: toNum(pt['MinDays'] ?? pt['minDays'] ?? null as unknown as number, undefined as unknown as number), maxDays: toNum(pt['MaxDays'] ?? pt['maxDays'] ?? null as unknown as number, undefined as unknown as number) }
            } else if ('min_days' in pt || 'max_days' in pt) {
              processingTime = { minDays: toNum(pt['min_days'] ?? pt['minDays'] ?? null as unknown as number, undefined as unknown as number), maxDays: toNum(pt['max_days'] ?? pt['maxDays'] ?? null as unknown as number, undefined as unknown as number) }
            }
          }
          const title = sObj['title'] || sObj['Title'] || sObj['name'] || sObj['Name'] || get(sObj['content'], 'title') || contentBlock['Title'] || contentBlock['Result']
          const summary = sObj['summary'] || sObj['Summary'] || sObj['description'] || sObj['Description'] || get(sObj['content'], 'summary') || get(sObj['content'], 'result') || contentBlock['Summary'] || contentBlock['Result']
          const base: AnyRecord = isObject(source) ? (source as AnyRecord) : {}
          const procedure: Procedure = { ...(base as object) as unknown as Procedure, id: asString(resolvedId), steps: steps as unknown as Procedure['steps'], documentsRequired: documentsRequired as unknown as Procedure['documentsRequired'], fees: fees as unknown as Procedure['fees'], processingTime: processingTime as unknown as Procedure['processingTime'], title: asString(title), name: asString(base['name'] ?? base['Name'] ?? ''), summary: asString(summary) }
          return { data: procedure }
        }
        return { error: lastError ?? ({ status: 404, data: { message: 'Procedure not found via any fallback path' } } as FetchBaseQueryError) }
      },
      providesTags: (result, error, id) => [{ type: 'Procedure', id }]
    })
  })
})

export const { useListProceduresQuery, useGetProcedureQuery, useGetProcedureFlexibleQuery } = proceduresApi
