import { cache, redirect } from '@solidjs/router';
import { Get } from '~/lib/request';

export const getUser = cache(async () => {
  const response = await Get('/api/me');
  if (!response.ok) {
    throw redirect('/');
  }
  return response.json();
}, 'user');

export const checkAuth = cache(async () => {
  const response = await Get('/api/me');
  return response.ok;
}, 'auth-check');
