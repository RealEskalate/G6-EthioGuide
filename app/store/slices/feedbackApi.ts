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
			transformResponse: (raw: any): FeedbackListResponse => {
				// Expected shapes:
				// 1) { feedbacks: [...], page, limit, total }
				// 2) { data: [...], pagination: { page, limit, total } }
				// 3) { feedbacks: { feedbacks: [...], page, limit, total } }
				console.log('Feedback API raw response:', raw)
				const container = (raw && !Array.isArray(raw.feedbacks) && raw.feedbacks && typeof raw.feedbacks === 'object') ? raw.feedbacks : raw
				const arr = Array.isArray(container?.feedbacks)
					? container.feedbacks
					: Array.isArray(container?.data)
						? container.data
						: Array.isArray(raw?.feedbacks)
							? raw.feedbacks
							: []
				const page = container?.page ?? container?.pagination?.page ?? raw?.page ?? raw?.pagination?.page ?? 1
				const limit = container?.limit ?? container?.pagination?.limit ?? raw?.limit ?? raw?.pagination?.limit ?? (arr.length || 5)
				const total = container?.total ?? container?.pagination?.total ?? raw?.total ?? raw?.pagination?.total ?? arr.length
				const feedbacks: FeedbackItem[] = arr.map((f: any) => ({
					id: f.id || f._id || f.feedback_id || crypto.randomUUID(),
					content: f.content || f.body || f.message || '',
					type: f.type || 'general',
					status: f.status || 'new',
					likeCount: f.like_count ?? f.likeCount ?? 0,
					dislikeCount: f.dislike_count ?? f.dislikeCount ?? 0,
					procedureID: f.procedure_id || f.procedureID || f.procedureId || '',
					createdAT: f.created_at || f.createdAt || f.createdAT,
					updatedAT: f.updated_at || f.updatedAt || f.updatedAT,
					userID: f.user_id || f.userID,
					tags: Array.isArray(f.tags) ? f.tags : undefined,
					adminResponse: f.admin_response || f.adminResponse,
					viewCount: f.view_count ?? f.viewCount
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
					if (cleanedTags && cleanedTags.length) (payload as any).Tags = cleanedTags
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
