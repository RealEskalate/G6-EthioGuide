import { configureStore } from '@reduxjs/toolkit'
import { proceduresApi } from '@/app/store/slices/proceduresApi'
import { checklistsApi } from '@/app/store/slices/checklistsApi'
import { categoriesApi } from '@/app/store/slices/categoriesApi'

export const store = configureStore({
	reducer: {
		[proceduresApi.reducerPath]: proceduresApi.reducer,
		[checklistsApi.reducerPath]: checklistsApi.reducer,
		[categoriesApi.reducerPath]: categoriesApi.reducer,
	},
	middleware: (getDefault) => getDefault().concat(
		proceduresApi.middleware,
		checklistsApi.middleware,
		categoriesApi.middleware
	)
})

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch
