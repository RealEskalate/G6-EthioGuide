// store/store.ts
import { configureStore } from "@reduxjs/toolkit";
import { apiSlice } from "./slices/workspaceSlice";
import { historyApi } from "./slices/historySlice";
// import authReducer from './slices/authSlice';
// import userReducer from './slices/userSlice';

export const store = configureStore({
  reducer: {
    [apiSlice.reducerPath]: apiSlice.reducer,
    [historyApi.reducerPath]: historyApi.reducer,
    // auth: authReducer,
    // user: userReducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(apiSlice.middleware, historyApi.middleware),
  devTools: true,
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
// No changes needed
