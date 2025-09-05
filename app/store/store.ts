// store/store.ts
import { configureStore } from "@reduxjs/toolkit";
import { apiSlice } from "./slices/workspaceSlice";
import { historyApi } from "./slices/historySlice";
import { discussionsListApi } from "./slices/discussionsGetSlice"; // added

export const store = configureStore({
  reducer: {
    [apiSlice.reducerPath]: apiSlice.reducer,
    [historyApi.reducerPath]: historyApi.reducer,
    [discussionsListApi.reducerPath]: discussionsListApi.reducer, // added
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(
      apiSlice.middleware,
      historyApi.middleware,
      discussionsListApi.middleware // added
    ),
  devTools: true,
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
