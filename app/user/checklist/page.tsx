"use client"

import { Suspense, useEffect } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useSession } from 'next-auth/react';
import { useCreateChecklistMutation } from '@/app/store/slices/checklistsApi';
import { useTranslation } from 'react-i18next';

function ChecklistRedirectInner() {
  const { t } = useTranslation('user');
  const router = useRouter();
  const search = useSearchParams();
  const { data: session, status } = useSession();
  const [createChecklist, { isLoading, isError }] = useCreateChecklistMutation();

  useEffect(() => {
    const procedureId = search.get('id');
    if (!procedureId) return;
    if (status === 'loading') return;
    if (status === 'unauthenticated') {
      router.replace('/auth/signin');
      return;
    }
    const run = async () => {
      try {
        const result = await createChecklist({ procedureId, token: session?.accessToken || undefined }).unwrap();
        console.log('Checklist created successfully:', result);
      } catch (error) {
        console.error('Failed to create checklist:', error);
      } finally {
        router.replace('/user/workspace');
      }
    };
    run();
  }, [search, status, session, createChecklist, router]);

  return (
    <div className="min-h-[50vh] flex items-center justify-center text-gray-600">
      {isLoading ? t('checklist_redirect.saving') : isError ? t('checklist_redirect.failed') : t('checklist_redirect.redirecting')}
    </div>
  );
}

export default function ChecklistRedirectPage() {
  const { t } = useTranslation('user');
  return (
    <Suspense fallback={<div className="min-h-[50vh] flex items-center justify-center text-gray-600">{t('checklist_redirect.redirecting')}</div>}>
      <ChecklistRedirectInner />
    </Suspense>
  );
}