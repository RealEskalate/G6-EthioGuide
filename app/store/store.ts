// store/store.ts
import { configureStore } from "@reduxjs/toolkit";
import { apiSlice } from "./slices/workspaceSlice";
import userReducer from "./slices/userSlice";
import { historyApi } from "./slices/historySlice";
import { discussionsListApi } from "./slices/discussionsGetSlice";
// i18n bootstrap (prevents "NO_I18NEXT_INSTANCE" warnings app-wide)
import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import { noticesApi } from "./slices/noticesSlice"; // added

if (!i18n.isInitialized) {
  i18n.use(initReactI18next).init({
    resources: {},
    lng: "en",
    fallbackLng: "en",
    interpolation: { escapeValue: false },
  });
}

export const store = configureStore({
  reducer: {
    user: userReducer,
    [apiSlice.reducerPath]: apiSlice.reducer,
    [historyApi.reducerPath]: historyApi.reducer,
    [discussionsListApi.reducerPath]: discussionsListApi.reducer,
    [noticesApi.reducerPath]: noticesApi.reducer, // added
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(
      apiSlice.middleware,
      historyApi.middleware,
      discussionsListApi.middleware,
      noticesApi.middleware // added
    ),
  devTools: true,
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
