import { i18n } from "i18next";

// List of supported languages
export const supportedLanguages = [
  { code: "en", name: "English", shortCode: "EN" },
  { code: "am", name: "አማርኛ (Amharic)", shortCode: "አማ" },
];

// Initialize language from localStorage or browser settings
export const initializeLanguage = (i18nInstance: i18n): string => {
  const savedLanguage = localStorage.getItem("i18nextLng");
  const browserLanguage = navigator.language.split("-")[0];
  const defaultLanguage = supportedLanguages.find(
    (lang) => lang.code === savedLanguage
  )?.code || supportedLanguages.find(
    (lang) => lang.code === browserLanguage
  )?.code || "en";

  i18nInstance.changeLanguage(defaultLanguage);
  return defaultLanguage;
};

// Save language to localStorage and update i18next
export const changeLanguage = (i18nInstance: i18n, languageCode: string): void => {
  if (supportedLanguages.some((lang) => lang.code === languageCode)) {
    i18nInstance.changeLanguage(languageCode);
    localStorage.setItem("i18nextLng", languageCode);
  }
};

// Get display name for language code
export const getLanguageDisplayName = (languageCode: string): string => {
  const lang = supportedLanguages.find((lang) => lang.code === languageCode);
  return lang ? lang.shortCode : languageCode;
};