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
import { configureStore } from "@reduxjs/toolkit";
import { setupListeners } from "@reduxjs/toolkit/query";
import { apiSlice } from "./slices/workspaceSlice";
import userReducer from "./slices/userSlice";
import { historyApi } from "./slices/historySlice";
import { discussionsListApi } from "./slices/discussionsGetSlice";
import aiChatReducer from "./slices/aiChatSlice";
import { noticesApi } from "./slices/noticesSlice";
import { feedbackApi } from "./slices/feedbackApi";

export const store = configureStore({
  reducer: {
    user: userReducer,
    aiChat: aiChatReducer,
    [apiSlice.reducerPath]: apiSlice.reducer,
    [historyApi.reducerPath]: historyApi.reducer,
    [discussionsListApi.reducerPath]: discussionsListApi.reducer,
    [noticesApi.reducerPath]: noticesApi.reducer,
    [feedbackApi.reducerPath]: feedbackApi.reducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(
      apiSlice.middleware,
      historyApi.middleware,
      discussionsListApi.middleware,
      noticesApi.middleware,
      feedbackApi.middleware
    ),
  devTools: true,
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

setupListeners(store.dispatch);
