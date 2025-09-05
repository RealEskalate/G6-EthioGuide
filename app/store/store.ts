import { configureStore } from '@reduxjs/toolkit'
import { setupListeners } from '@reduxjs/toolkit/query'
import { feedbackApi } from './slices/feedbackApi'

export const store = configureStore({
	reducer: {
		[feedbackApi.reducerPath]: feedbackApi.reducer,
	},
	middleware: (getDefault) => getDefault().concat(feedbackApi.middleware)
})

export type AppStore = typeof store
export type RootState = ReturnType<AppStore['getState']>
export type AppDispatch = AppStore['dispatch']

setupListeners(store.dispatch)

