import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'
import type { Procedure, ProcedureDocument, ProcedureFee, ProcedureProcessingTime, ProcedureStep } from '@/app/types/procedure'

// Small helpers for type-safe narrowing without using 'any'
const hasKey = <K extends string>(o: unknown, k: K): o is Record<K, unknown> =>
  typeof o === 'object' && o !== null && k in (o as Record<string, unknown>)
const isErrorResult = (res: unknown): res is { error: unknown } => hasKey(res, 'error')

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
        const arr = (raw as Paginated<Procedure>).data || (raw as Paginated<Procedure>).items || (raw as Paginated<Procedure>).results || (Array.isArray(raw) ? (raw as Procedure[]) : [])
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
            // tolerate alternative id fields from backend without introducing any
            const id = (p as unknown as { id?: string; _id?: string; procedureId?: string; slug?: string; code?: string }).id
              || (p as unknown as { _id?: string })._id
              || (p as unknown as { procedureId?: string }).procedureId
              || (p as unknown as { slug?: string }).slug
              || (p as unknown as { code?: string }).code
              || undefined
            return { ...p, id: id as string }
          }).filter((p): p is Procedure => Boolean(p.id)),
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
        const isRecord = (v: unknown): v is Record<string, unknown> => typeof v === 'object' && v !== null && !Array.isArray(v)
        const getRecord = (o: unknown, key: string): Record<string, unknown> | undefined => {
          if (!isRecord(o)) return undefined
          const v = o[key]
          return isRecord(v) ? v : undefined
        }
        const pickStr = (o: unknown, keys: string[]): string | undefined => {
          if (!isRecord(o)) return undefined
          for (const k of keys) {
            const v = o[k]
            if (typeof v === 'string') return v
          }
          return undefined
        }
        const pickNum = (o: unknown, keys: string[]): number | undefined => {
          if (!isRecord(o)) return undefined
          for (const k of keys) {
            const v = o[k]
            if (typeof v === 'number') return v
          }
          return undefined
        }

        // Accept various envelope shapes
        let source: unknown = raw
        const dataObj = getRecord(source, 'data')
        if (dataObj) source = dataObj
        const procObj = getRecord(source, 'procedure')
        if (procObj) source = procObj
        const nestedData = getRecord(getRecord(source, 'data'), 'procedure')
        if (nestedData) source = nestedData
        if (Array.isArray(source)) source = source[0] ?? {}

        const contentBlock: Record<string, unknown> = getRecord(source, 'content') || getRecord(source, 'Content') || {}

  const idPick = pickStr(source, ['id', 'ID', '_id', 'procedureId', 'procedure_id', 'uuid', 'slug', 'code'])

        // Steps
        let steps: ProcedureStep[] | undefined
        const directSteps = (isRecord(source) ? source['steps'] : undefined) as unknown
        const altSteps = (isRecord(contentBlock) ? (contentBlock['steps'] ?? contentBlock['Steps']) : undefined) as unknown
        const rawSteps = directSteps ?? altSteps
        if (Array.isArray(rawSteps)) {
      steps = rawSteps.map((s, i) => {
            if (typeof s === 'string') return { order: i + 1, text: s }
            if (isRecord(s)) {
              const order = typeof s.order === 'number' ? s.order : i + 1
              const text = typeof s.text === 'string' ? s.text : (typeof s.description === 'string' ? s.description : JSON.stringify(s))
              const title = typeof s.title === 'string' ? s.title : undefined
              const description = typeof s.description === 'string' ? s.description : undefined
        const estimatedTime = typeof s.estimatedTime === 'number' ? String(s.estimatedTime) : (typeof s.estimatedTime === 'string' ? s.estimatedTime : undefined)
        const time = typeof s.time === 'number' ? String(s.time) : (typeof s.time === 'string' ? s.time : undefined)
        return { order, text, title, description, estimatedTime, time }
            }
            return { order: i + 1, text: JSON.stringify(s) }
          })
        } else if (isRecord(rawSteps)) {
          steps = Object.entries(rawSteps).map(([k, v], idx) => ({ order: Number(k) || idx + 1, text: typeof v === 'string' ? v : JSON.stringify(v) }))
        }

        // Documents
        let documentsRequired: ProcedureDocument[] | undefined
        const docsCandidate = (isRecord(source) ? (source['documentsRequired'] ?? source['documents'] ?? source['docs']) : undefined)
          ?? (isRecord(getRecord(source, 'content')) ? (getRecord(source, 'content')!['documentsRequired'] ?? getRecord(source, 'content')!['documents']) : undefined)
          ?? (contentBlock['Documents'] ?? contentBlock['documents'])
        if (Array.isArray(docsCandidate)) {
          documentsRequired = docsCandidate.map((d) => {
            if (typeof d === 'string') return { name: d, templateUrl: null }
            if (isRecord(d)) {
              const name = pickStr(d, ['name', 'title', 'label', 'filename']) || 'Document'
              const templateUrl = pickStr(d, ['templateUrl', 'url']) || null
              return { name, templateUrl }
            }
            return { name: 'Document', templateUrl: null }
          })
        } else {
          const pre = (contentBlock['Prerequisites'] ?? contentBlock['prerequisites']) as unknown
          if (Array.isArray(pre)) {
            documentsRequired = pre.map((d) => typeof d === 'string' ? { name: d, templateUrl: null } : { name: (isRecord(d) ? (pickStr(d, ['name', 'title']) || 'Document') : 'Document'), templateUrl: null })
          }
        }

        // Fees
        let fees: ProcedureFee[] | undefined
        const feesRaw = (isRecord(source) ? source['fees'] : undefined) ?? (getRecord(source, 'content')?.['fees']) ?? (contentBlock['Fees'] as unknown)
        if (Array.isArray(feesRaw)) {
          fees = feesRaw.filter(isRecord).map((f) => ({
            amount: pickNum(f, ['amount', 'Amount']) ?? 0,
            currency: pickStr(f, ['currency', 'Currency']) || '',
            label: pickStr(f, ['label', 'Label']) || 'Fee'
          }))
        } else if (isRecord(feesRaw)) {
          const hasAmount = typeof feesRaw['Amount'] === 'number' || typeof feesRaw['amount'] === 'number'
          if (hasAmount) {
            fees = [{
              amount: pickNum(feesRaw, ['Amount', 'amount']) ?? 0,
              currency: pickStr(feesRaw, ['Currency', 'currency']) || '',
              label: pickStr(feesRaw, ['Label', 'label']) || 'Fee'
            }]
          } else {
            fees = Object.values(feesRaw).filter(isRecord).map((f) => ({
              amount: pickNum(f, ['amount', 'Amount']) ?? 0,
              currency: pickStr(f, ['currency', 'Currency']) || '',
              label: pickStr(f, ['label', 'Label']) || 'Fee'
            }))
          }
        }

        // Processing time
        let processingTime: ProcedureProcessingTime | undefined
        const ptRaw = (isRecord(source) ? source['processingTime'] : undefined) ?? (getRecord(source, 'content')?.['processingTime']) ?? (contentBlock['ProcessingTime'] as unknown)
        if (typeof ptRaw === 'number') {
          processingTime = { minDays: ptRaw, maxDays: ptRaw }
        } else if (isRecord(ptRaw)) {
          processingTime = {
            minDays: pickNum(ptRaw, ['MinDays', 'minDays']),
            maxDays: pickNum(ptRaw, ['MaxDays', 'maxDays'])
          }
        }

  const name = pickStr(source, ['name', 'Name'])
  const titlePick = pickStr(source, ['title', 'Title', 'name', 'Name']) || pickStr(contentBlock, ['title', 'Title', 'Result'])
  const title = titlePick || name || 'Untitled Procedure'
  const summary = pickStr(source, ['summary', 'Summary', 'description', 'Description']) || pickStr(contentBlock, ['summary', 'Summary', 'result', 'Result'])
        const tags = (isRecord(source) && Array.isArray(source['tags'])) ? (source['tags'] as string[]) : undefined
        const verified = (isRecord(source) && typeof source['verified'] === 'boolean') ? (source['verified'] as boolean) : undefined
        const updatedAt = pickStr(source, ['updatedAt', 'updated_at'])
        const views = pickNum(source, ['views'])
        const likes = pickNum(source, ['likes'])

  const id = idPick ?? (typeof crypto !== 'undefined' && 'randomUUID' in crypto ? crypto.randomUUID() : 'unknown')
  return { id, steps, documentsRequired, fees, processingTime, title, name, summary, tags, verified, updatedAt, views, likes }
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
        let lastError: unknown = null
        for (const p of paths) {
          const res = await baseFetch(p)
          if (isErrorResult(res)) { lastError = (res as { error: unknown }).error; continue }
          const raw = (res as { data?: unknown }).data
          // Heuristic: if raw is empty object or array length 0, continue
          if (!raw || (Array.isArray(raw) && raw.length === 0) || (typeof raw === 'object' && raw !== null && !Array.isArray(raw) && Object.keys(raw as Record<string, unknown>).length === 0)) {
            continue
          }
          // reuse helpers
          const isRecord = (v: unknown): v is Record<string, unknown> => typeof v === 'object' && v !== null && !Array.isArray(v)
          const getRecord = (o: unknown, key: string): Record<string, unknown> | undefined => {
            if (!isRecord(o)) return undefined
            const v = o[key]
            return isRecord(v) ? v : undefined
          }
          const pickStr = (o: unknown, keys: string[]): string | undefined => {
            if (!isRecord(o)) return undefined
            for (const k of keys) { const v = o[k]; if (typeof v === 'string') return v }
            return undefined
          }
          const pickNum = (o: unknown, keys: string[]): number | undefined => {
            if (!isRecord(o)) return undefined
            for (const k of keys) { const v = o[k]; if (typeof v === 'number') return v }
            return undefined
          }
          let source: unknown = raw
          const dataObj = getRecord(source, 'data'); if (dataObj) source = dataObj
          const procObj = getRecord(source, 'procedure'); if (procObj) source = procObj
          const nested = getRecord(getRecord(source, 'data'), 'procedure'); if (nested) source = nested
          if (Array.isArray(source)) source = source[0] ?? {}
          const contentBlock: Record<string, unknown> = getRecord(source, 'content') || getRecord(source, 'Content') || {}
          const resolvedId = pickStr(source, ['id', 'ID', '_id', 'procedureId', 'procedure_id', 'uuid', 'slug', 'code']) || id
          // Steps
          let steps: ProcedureStep[] | undefined
          const stepsRaw = (isRecord(source) ? source['steps'] : undefined) ?? (contentBlock['steps'] ?? contentBlock['Steps'])
          if (Array.isArray(stepsRaw)) {
            steps = stepsRaw.map((s, i) => {
              if (typeof s === 'string') return { order: i + 1, text: s }
              if (isRecord(s)) {
                const order = typeof s.order === 'number' ? s.order : i + 1
                const text = typeof s.text === 'string' ? s.text : (typeof s.description === 'string' ? s.description : JSON.stringify(s))
                const estimatedTime = typeof s.estimatedTime === 'number' ? String(s.estimatedTime) : (typeof s.estimatedTime === 'string' ? s.estimatedTime : undefined)
                const time = typeof s.time === 'number' ? String(s.time) : (typeof s.time === 'string' ? s.time : undefined)
                return { order, text, estimatedTime, time }
              }
              return { order: i + 1, text: JSON.stringify(s) }
            })
          } else if (isRecord(stepsRaw)) {
            steps = Object.entries(stepsRaw).map(([k, v], idx) => ({ order: Number(k) || idx + 1, text: typeof v === 'string' ? v : JSON.stringify(v) }))
          }
          // Documents
          let documentsRequired: ProcedureDocument[] | undefined
          const docsRaw = (isRecord(source) ? (source['documentsRequired'] ?? source['documents'] ?? source['docs']) : undefined) ?? (contentBlock['Documents'] ?? contentBlock['documents'])
          if (Array.isArray(docsRaw)) {
            documentsRequired = docsRaw.map((d) => typeof d === 'string' ? ({ name: d, templateUrl: null }) : (isRecord(d) ? ({ name: pickStr(d, ['name', 'title', 'label', 'filename']) || 'Document', templateUrl: pickStr(d, ['templateUrl', 'url']) || null }) : ({ name: 'Document', templateUrl: null })))
          }
          // Fees
          let fees: ProcedureFee[] | undefined
          const feesRaw2 = (isRecord(source) ? source['fees'] : undefined) ?? (contentBlock['Fees'])
          if (Array.isArray(feesRaw2)) {
            fees = feesRaw2.filter(isRecord).map((f) => ({ amount: pickNum(f, ['amount', 'Amount']) ?? 0, currency: pickStr(f, ['currency', 'Currency']) || '', label: pickStr(f, ['label', 'Label']) || 'Fee' }))
          } else if (isRecord(feesRaw2)) {
            const hasAmount = typeof feesRaw2['Amount'] === 'number' || typeof feesRaw2['amount'] === 'number'
            fees = hasAmount ? [{ amount: pickNum(feesRaw2, ['Amount', 'amount']) ?? 0, currency: pickStr(feesRaw2, ['Currency', 'currency']) || '', label: pickStr(feesRaw2, ['Label', 'label']) || 'Fee' }] : Object.values(feesRaw2).filter(isRecord).map((f) => ({ amount: pickNum(f, ['amount', 'Amount']) ?? 0, currency: pickStr(f, ['currency', 'Currency']) || '', label: pickStr(f, ['label', 'Label']) || 'Fee' }))
          }
          // Processing time
          let processingTime: ProcedureProcessingTime | undefined
          const ptRaw = (isRecord(source) ? source['processingTime'] : undefined) ?? (contentBlock['ProcessingTime'])
          if (typeof ptRaw === 'number') processingTime = { minDays: ptRaw, maxDays: ptRaw }
          else if (isRecord(ptRaw)) processingTime = { minDays: pickNum(ptRaw, ['MinDays', 'minDays']), maxDays: pickNum(ptRaw, ['MaxDays', 'maxDays']) }

          const name = pickStr(source, ['name', 'Name'])
          const titlePick = pickStr(source, ['title', 'Title', 'name', 'Name']) || pickStr(contentBlock, ['title', 'Title', 'Result'])
          const title = titlePick || name || 'Untitled Procedure'
          const summary = pickStr(source, ['summary', 'Summary', 'description', 'Description']) || pickStr(contentBlock, ['summary', 'Summary', 'result', 'Result'])
          const tags = (isRecord(source) && Array.isArray(source['tags'])) ? (source['tags'] as string[]) : undefined
          const verified = (isRecord(source) && typeof source['verified'] === 'boolean') ? (source['verified'] as boolean) : undefined
          const updatedAt = pickStr(source, ['updatedAt', 'updated_at'])
          const views = pickNum(source, ['views'])
          const likes = pickNum(source, ['likes'])
          const procedure: Procedure = { id: (resolvedId || id) as string, steps, documentsRequired, fees, processingTime, title, name, summary, tags, verified, updatedAt, views, likes }
          return { data: procedure }
        }
        // Ensure error conforms to FetchBaseQueryError likeness
        const fallbackError = { status: 404, data: { message: 'Procedure not found via any fallback path' } }
        const err = (lastError && typeof lastError === 'object') ? (lastError as never) : (fallbackError as never)
        return { error: err }
      },
      providesTags: (result, error, id) => [{ type: 'Procedure', id }]
    })
  })
})

export const { useListProceduresQuery, useGetProcedureQuery, useGetProcedureFlexibleQuery } = proceduresApi
