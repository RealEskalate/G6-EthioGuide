"use client";

import { SessionProvider } from 'next-auth/react';
import { I18nextProvider } from 'react-i18next';
import i18next from '@/lib/i18n/i18n';
import { store } from '@/app/store/store';
import { Provider } from "react-redux";
import { ReactNode } from 'react';
import { Toaster } from "@/components/ui/sonner";

interface ProvidersProps {
  children: ReactNode;
}

export function Providers({ children }: ProvidersProps) {
  return (

    <SessionProvider>
      <I18nextProvider i18n={i18next}>
        <Provider store={store}>
          {children}
          <Toaster />
        </Provider>
      </I18nextProvider>
    </SessionProvider>

//     <Provider store={store}>
//       <SessionProvider>
//         <I18nextProvider i18n={i18next}>
//           {children}
//         </I18nextProvider>
//       </SessionProvider>
//     </Provider>

  );
}
