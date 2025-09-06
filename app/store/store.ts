import { configureStore } from '@reduxjs/toolkit'
import { proceduresApi } from '@/app/store/slices/proceduresApi'
import { checklistsApi } from '@/app/store/slices/checklistsApi'
import { categoriesApi } from '@/app/store/slices/categoriesApi'
import { feedbackApi } from '@/app/store/slices/feedbackApi'
import { discussionsApi } from '@/app/store/slices/discussionsApi'

export const store = configureStore({
	reducer: {
		[proceduresApi.reducerPath]: proceduresApi.reducer,
		[checklistsApi.reducerPath]: checklistsApi.reducer,
		[categoriesApi.reducerPath]: categoriesApi.reducer,
		[feedbackApi.reducerPath]: feedbackApi.reducer,
		[discussionsApi.reducerPath]: discussionsApi.reducer,
	},
	middleware: (getDefault) => getDefault().concat(
		proceduresApi.middleware,
		checklistsApi.middleware,
		categoriesApi.middleware,
		feedbackApi.middleware,
		discussionsApi.middleware
	)
})

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch
