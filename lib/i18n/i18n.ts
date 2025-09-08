// import i18next from 'i18next';
// import { initReactI18next } from 'react-i18next';
// import LanguageDetector from 'i18next-browser-languagedetector';
// import enCommon from '../../public/locales/en/common.json';
// import amCommon from '../../public/locales/am/common.json';
// import enAuth from '../../public/locales/en/auth.json';
// import amAuth from '../../public/locales/am/auth.json';

// i18next
//   .use(initReactI18next)
//   .use(LanguageDetector)
//   .init({
//     resources: {
//       en: {
//         common: enCommon,
//         auth: enAuth,
//       },
//       am: {
//         common: amCommon,
//         auth: amAuth,
//       },
//     },
//     fallbackLng: 'en',
//     supportedLngs: ['en', 'am'],
//     ns: ['common', 'auth'],
//     defaultNS: 'common',
//     detection: {
//       order: ['localStorage', 'navigator'],
//       lookupLocalStorage: 'i18nextLng',
//     },
//     interpolation: {
//       escapeValue: false,
//     },
//   });

// export default i18next;

import i18next from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';
import enCommon from '../../public/locales/en/common.json';
import amCommon from '../../public/locales/am/common.json';
import enAuth from '../../public/locales/en/auth.json';
import amAuth from '../../public/locales/am/auth.json';
import enUser from '../../public/locales/en/user.json';
import amUser from '../../public/locales/am/user.json';

i18next
  .use(initReactI18next)
  .use(LanguageDetector)
  .init({
    resources: {
      en: {
        common: enCommon,
        auth: enAuth,
        user: enUser,
      },
      am: {
        common: amCommon,
        auth: amAuth,
        user: amUser,
      },
    },
    fallbackLng: 'en',
    supportedLngs: ['en', 'am'],
    ns: ['common', 'auth', 'user'],
    defaultNS: 'common',
    detection: {
      order: ['sessionStorage', 'navigator'],
      lookupSessionStorage: 'i18nextLng',
    },
    interpolation: {
      escapeValue: false,
    },
  });

export default i18next;
