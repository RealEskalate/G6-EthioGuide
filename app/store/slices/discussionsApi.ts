import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'

export interface DiscussionPost {
  id: string
  title: string
  content: string
  tags: string[]
  procedures: string[]
  userID: string
  createdAt: string
  updatedAt: string
}

export interface DiscussionsListResponse {
  posts: DiscussionPost[]
  page: number
  limit: number
  total: number
}

export const discussionsApi = createApi({
  reducerPath: 'discussionsApi',
  baseQuery: fetchBaseQuery({ baseUrl: '/api/v1' }),
  tagTypes: ['DiscussionsList', 'Discussion'],
  endpoints: (builder) => ({
    listDiscussions: builder.query<
      DiscussionsListResponse,
      { title?: string; procedureIds?: string[]; tags?: string[]; page?: number; limit?: number; sort_by?: string; sort_order?: 'asc' | 'desc' }
    >({
      query: ({ title, procedureIds, tags, page = 1, limit = 10, sort_by, sort_order }) => {
        const params = new URLSearchParams()
        params.set('page', String(page))
        params.set('limit', String(limit))
        if (title) params.set('title', title)
        if (Array.isArray(procedureIds)) {
          for (const pid of procedureIds) params.append('procedure_ids', pid)
        }
        if (Array.isArray(tags)) {
          for (const t of tags) params.append('tags', t)
        }
        if (sort_by) params.set('sort_by', sort_by)
        if (sort_order) params.set('sort_order', sort_order)
        return { url: `/discussions?${params.toString()}` }
      },
      transformResponse: (raw: unknown): DiscussionsListResponse => {
        type Raw = { posts?: unknown; page?: number; limit?: number; total?: number }
        type PostRaw = {
          id?: string; _id?: string; title?: string; content?: string;
          tags?: unknown; procedures?: unknown; userID?: string; user_id?: string;
          createdAt?: string; created_at?: string; updatedAt?: string; updated_at?: string;
        }
        const data = (raw as Raw) || {}
        const rawPosts = (() => {
          const p = (data as Raw).posts
          return Array.isArray(p) ? (p as unknown[]) : []
        })()
        const posts: DiscussionPost[] = rawPosts.map((p) => {
          const pr = p as PostRaw
          return {
            id: pr.id || pr._id || '',
            title: pr.title || '',
            content: pr.content || '',
            tags: Array.isArray(pr.tags as string[]) ? (pr.tags as string[]) : [],
            procedures: Array.isArray(pr.procedures as string[]) ? (pr.procedures as string[]) : [],
            userID: pr.userID || pr.user_id || '',
            createdAt: pr.createdAt || pr.created_at || '',
            updatedAt: pr.updatedAt || pr.updated_at || '',
          }
        })
        return {
          posts,
          page: (data.page as number) ?? 1,
          limit: (data.limit as number) ?? (posts.length || 10),
          total: (data.total as number) ?? posts.length,
        }
      },
      providesTags: (res) => [
        { type: 'DiscussionsList' as const, id: 'LIST' },
        ...((res?.posts ?? []).map(p => ({ type: 'Discussion' as const, id: p.id })))
      ]
    }),
    getDiscussion: builder.query<DiscussionPost, string>({
      query: (id) => ({ url: `/discussions/${encodeURIComponent(id)}` }),
      transformResponse: (raw: unknown): DiscussionPost => {
        const p = (raw || {}) as { id?: string; _id?: string; title?: string; content?: string; tags?: unknown; procedures?: unknown; userID?: string; user_id?: string; createdAt?: string; created_at?: string; updatedAt?: string; updated_at?: string }
        return {
          id: p.id || p._id || '',
          title: p.title || '',
          content: p.content || '',
          tags: Array.isArray(p.tags as string[]) ? (p.tags as string[]) : [],
          procedures: Array.isArray(p.procedures as string[]) ? (p.procedures as string[]) : [],
          userID: p.userID || p.user_id || '',
          createdAt: p.createdAt || p.created_at || '',
          updatedAt: p.updatedAt || p.updated_at || '',
        }
      },
      providesTags: (_res, _e, id) => [{ type: 'Discussion', id }]
    })
  })
})

export const { useListDiscussionsQuery, useGetDiscussionQuery } = discussionsApi
