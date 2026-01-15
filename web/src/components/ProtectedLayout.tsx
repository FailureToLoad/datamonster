import { Suspense, type ParentProps } from 'solid-js';
import { query, createAsync } from '@solidjs/router';
import { Get } from '~/lib/request';
import { GlossaryProvider } from '~/context/glossary';
import styles from './ProtectedLayout.module.css';

export const checkAuth = query(async () => {
  const response = await Get('/api/me');
  if (!response.ok) {
    window.location.href = '/auth/login';
  }
}, 'auth');

function LoadingScreen() {
  return (
    <div class={styles.loadingLayout}>
      <span class="loading loading-spinner loading-lg" />
    </div>
  );
}

export default function ProtectedLayout(props: ParentProps) {
  createAsync(() => checkAuth());

  return (
    <Suspense fallback={<LoadingScreen />}>
      <GlossaryProvider>{props.children}</GlossaryProvider>
    </Suspense>
  );
}
