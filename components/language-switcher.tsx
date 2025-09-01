'use client';

import { useTranslation } from 'react-i18next';
import { Button } from '@/components/ui/button';
import { useEffect } from 'react';

export function LanguageSwitcher() {
  const { i18n } = useTranslation();

  const changeLanguage = (lng: string) => {
    i18n.changeLanguage(lng);
    localStorage.setItem('i18nextLng', lng);
  };

  // Log current language for debugging
  useEffect(() => {
    console.log('Current language:', i18n.language);
  }, [i18n.language]);

  return (
    <div className="flex space-x-2">
      <Button
        variant={i18n.language === 'en' ? 'default' : 'outline'}
        className="bg-primary text-primary-foreground"
        onClick={() => changeLanguage('en')}
      >
        English
      </Button>
      <Button
        variant={i18n.language === 'am' ? 'default' : 'outline'}
        className="bg-primary text-primary-foreground"
        onClick={() => changeLanguage('am')}
      >
        አማርኛ
      </Button>
    </div>
  );
}