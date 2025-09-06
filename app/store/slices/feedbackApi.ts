import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'

export interface SubmitFeedbackArgs {
	procedureId: string
	content: string
	type: string // inaccuracy | feature_request | general | ...
	tags?: string[]
	token?: string | null
}

export interface UpdateFeedbackArgs {
	feedbackId: string
	adminResponse?: string
	status?: string
	token: string
}

export interface FeedbackItem {
	id: string
	content: string
	type: string
	status: string
	likeCount: number
	dislikeCount: number
	procedureID: string
	createdAT?: string
	updatedAT?: string
	userID?: string
	tags?: string[]
	adminResponse?: string
	viewCount?: number
}

interface FeedbackListResponse {
	feedbacks: FeedbackItem[]
	page: number
	limit: number
	total: number
}

export const feedbackApi = createApi({
	reducerPath: 'feedbackApi',
	baseQuery: fetchBaseQuery({ baseUrl: '/api/v1' }),
	tagTypes: ['FeedbackList', 'FeedbackItem'],
	endpoints: (builder) => ({
		getProcedureFeedback: builder.query<
			FeedbackListResponse,
			{ procedureId: string; page?: number; limit?: number; token?: string | null; status?: string }
		>({
			query: ({ procedureId, page = 1, limit = 5, token, status }) => {
				const params = new URLSearchParams()
				params.set('page', String(page))
				params.set('limit', String(limit))
				if (status) params.set('status', status)
				const headers: Record<string,string> = { 'Accept': 'application/json' }
				if (token) {
					const bearer = token.startsWith('Bearer ') ? token : `Bearer ${token}`
					headers['Authorization'] = bearer
				}
				return { url: `/procedures/${encodeURIComponent(procedureId)}/feedback?${params.toString()}`, headers }
			},
			transformResponse: (raw: unknown): FeedbackListResponse => {
				// Expected shapes:
				// 1) { feedbacks: [...], page, limit, total }
				// 2) { data: [...], pagination: { page, limit, total } }
				// 3) { feedbacks: { feedbacks: [...], page, limit, total } }
				console.log('Feedback API raw response:', raw)
				const isObject = (v: unknown): v is Record<string, unknown> => typeof v === 'object' && v !== null
				const get = <T = unknown>(o: unknown, key: string): T | undefined => (isObject(o) ? (o[key] as T) : undefined)
				const containerCandidate = (isObject(raw) && !Array.isArray(get(raw, 'feedbacks')) && isObject(get(raw, 'feedbacks')))
					? (get<Record<string, unknown>>(raw, 'feedbacks') as unknown)
					: raw
				const containerObj = isObject(containerCandidate) ? containerCandidate : {}
				const fbField = get<unknown[]>(containerObj, 'feedbacks')
				const dataField = get<unknown[]>(containerObj, 'data')
				const rawFbField = get<unknown[]>(raw, 'feedbacks')
				const arr: Array<Record<string, unknown>> = Array.isArray(fbField)
					? (fbField as Array<Record<string, unknown>>)
					: Array.isArray(dataField)
						? (dataField as Array<Record<string, unknown>>)
						: Array.isArray(rawFbField)
							? (rawFbField as Array<Record<string, unknown>>)
							: []
				const containerPage = get<number>(containerObj, 'page') ?? get<number>(get(containerObj, 'pagination'), 'page')
				const rawPage = get<number>(raw, 'page') ?? get<number>(get(raw, 'pagination'), 'page')
				const page = containerPage ?? rawPage ?? 1
				const containerLimit = get<number>(containerObj, 'limit') ?? get<number>(get(containerObj, 'pagination'), 'limit')
				const rawLimit = get<number>(raw, 'limit') ?? get<number>(get(raw, 'pagination'), 'limit')
				const limit = containerLimit ?? rawLimit ?? (arr.length || 5)
				const containerTotal = get<number>(containerObj, 'total') ?? get<number>(get(containerObj, 'pagination'), 'total')
				const rawTotal = get<number>(raw, 'total') ?? get<number>(get(raw, 'pagination'), 'total')
				const total = containerTotal ?? rawTotal ?? arr.length
				const asString = (v: unknown, fallback = ''): string => (typeof v === 'string' ? v : v == null ? fallback : String(v))
				const toNum = (v: unknown, fallback = 0): number => (typeof v === 'number' ? v : typeof v === 'string' && v !== '' ? Number(v) : fallback)
				const feedbacks: FeedbackItem[] = arr.map((f) => ({
					id: asString(f.id ?? f._id ?? (f as Record<string, unknown>)['feedback_id'] ?? crypto.randomUUID()),
					content: asString(f.content ?? (f as Record<string, unknown>)['body'] ?? (f as Record<string, unknown>)['message'] ?? ''),
					type: asString(f.type ?? 'general'),
					status: asString(f.status ?? 'new'),
					likeCount: toNum((f as Record<string, unknown>)['like_count'] ?? f['likeCount'] ?? 0),
					dislikeCount: toNum((f as Record<string, unknown>)['dislike_count'] ?? f['dislikeCount'] ?? 0),
					procedureID: asString((f as Record<string, unknown>)['procedure_id'] ?? f['procedureID'] ?? f['procedureId'] ?? ''),
					createdAT: asString((f as Record<string, unknown>)['created_at'] ?? f['createdAt'] ?? f['createdAT'] ?? ''),
					updatedAT: asString((f as Record<string, unknown>)['updated_at'] ?? f['updatedAt'] ?? f['updatedAT'] ?? ''),
					userID: asString((f as Record<string, unknown>)['user_id'] ?? f['userID'] ?? ''),
					tags: Array.isArray((f as Record<string, unknown>)['tags']) ? ((f as Record<string, unknown>)['tags'] as string[]) : undefined,
					adminResponse: asString((f as Record<string, unknown>)['admin_response'] ?? f['adminResponse'] ?? ''),
					viewCount: toNum((f as Record<string, unknown>)['view_count'] ?? f['viewCount'] ?? undefined, undefined as unknown as number)
				}))
				const result = { feedbacks, page, limit, total }
				console.log('Feedback API transformed result:', result)
				return result
			},
			providesTags: (result, _err, arg) => [
				{ type: 'FeedbackList', id: arg.procedureId },
				...(result?.feedbacks || []).map((f) => ({ type: 'FeedbackItem' as const, id: f.id }))
			]
		}),
		submitProcedureFeedback: builder.mutation<FeedbackItem, SubmitFeedbackArgs>({
				query: ({ procedureId, content, type, tags, token }) => {
					const headers: Record<string,string> = { 'Content-Type': 'application/json', 'Accept': 'application/json' }
					if (token) {
						const bearer = token.startsWith('Bearer ') ? token : `Bearer ${token}`
						headers['Authorization'] = bearer
					}
					const cleanedTags = Array.isArray(tags) ? tags.map(t => t.trim()).filter(Boolean) : undefined
					// Map UI values to backend enums (supporting common misspellings)
					const mapToEnum = (v: string) => {
						const s = String(v || '').toLowerCase().trim().replace(/\s+/g, '_')
						if (['inaccuracy','inacuuracy','inacuracy','incorrect','error','issue'].includes(s)) return 'inaccuracy'
						if (['improvement','inmprovement','improvment','enhancement','suggestion'].includes(s)) return 'improvement'
						return 'other'
					}
					const payload: Record<string, unknown> = { Content: content, Type: mapToEnum(type), ProcedureID: procedureId }
					if (cleanedTags && cleanedTags.length) payload['Tags'] = cleanedTags
					return {
						url: `/procedures/${encodeURIComponent(procedureId)}/feedback`,
						method: 'POST',
						headers,
						body: payload
					}
				},
					onQueryStarted: async (arg, { dispatch, queryFulfilled }) => {
						try {
							const { data: created } = await queryFulfilled
							// Update page 1 cache for the same token (if any)
							dispatch(
								feedbackApi.util.updateQueryData(
									'getProcedureFeedback',
									{ procedureId: arg.procedureId, page: 1, limit: 5, token: arg.token ?? null },
									(draft) => {
										if (!draft) return
										draft.feedbacks = [created as FeedbackItem, ...(draft.feedbacks || [])]
										draft.total = (draft.total || 0) + 1
									}
								)
							)
							// Also update cache with no token key (in case)
							dispatch(
								feedbackApi.util.updateQueryData(
									'getProcedureFeedback',
									{ procedureId: arg.procedureId, page: 1, limit: 5, token: null },
									(draft) => {
										if (!draft) return
										draft.feedbacks = [created as FeedbackItem, ...(draft.feedbacks || [])]
										draft.total = (draft.total || 0) + 1
									}
								)
							)
						} catch {
							// ignore
						}
					},
			invalidatesTags: (_res, _err, arg) => [
				{ type: 'FeedbackList', id: arg.procedureId }
			]
		}),
		updateFeedback: builder.mutation<{ message: string }, UpdateFeedbackArgs>({
			query: ({ feedbackId, adminResponse, status, token }) => {
				const headers: Record<string,string> = { 'Content-Type': 'application/json' }
				const bearer = token.startsWith('Bearer ') ? token : `Bearer ${token}`
				headers['Authorization'] = bearer
				return {
					url: `/feedback/${encodeURIComponent(feedbackId)}`,
					method: 'PATCH',
					headers,
					body: { admin_response: adminResponse, status }
				}
			},
			invalidatesTags: (_res, _err, arg) => [
				{ type: 'FeedbackItem', id: arg.feedbackId }
			]
		})
	})
})

export const { useGetProcedureFeedbackQuery, useSubmitProcedureFeedbackMutation, useUpdateFeedbackMutation } = feedbackApi
