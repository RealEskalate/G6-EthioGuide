import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'
import type { Category } from '@/app/types/category'

interface Paginated<T> { data?: T[]; page?: number; limit?: number; total?: number }

export const categoriesApi = createApi({
  reducerPath: 'categoriesApi',
  baseQuery: fetchBaseQuery({ baseUrl: '/api/v1' }),
  tagTypes: ['Categories'],
  endpoints: (builder) => ({
    listCategories: builder.query<{ list: Category[]; page: number; limit: number; total: number }, { page?: number; limit?: number; sortOrder?: 'asc' | 'desc'; parentID?: string; organizationID?: string; title?: string } | void>({
      query: (args) => {
        const params = new URLSearchParams()
        if (args?.page) params.set('page', String(args.page))
        if (args?.limit) params.set('limit', String(args.limit))
        if (args?.sortOrder) params.set('sortOrder', args.sortOrder)
        if (args?.parentID) params.set('parentID', args.parentID)
        if (args?.organizationID) params.set('organizationID', args.organizationID)
        if (args?.title) params.set('title', args.title)
        return `/categories${params.size ? `?${params.toString()}` : ''}`
      },
      transformResponse: (raw: Paginated<Category>): { list: Category[]; page: number; limit: number; total: number } => {
        const list = raw.data || []
        return { list: list.map(c => ({ ...c, id: (c as any).id || (c as any)._id })), page: raw.page || 1, limit: raw.limit || list.length || 10, total: raw.total || list.length }
      },
      providesTags: () => [ { type: 'Categories', id: 'LIST' } ]
    })
  })
})

export const { useListCategoriesQuery } = categoriesApi