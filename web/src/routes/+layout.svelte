<script lang="ts">
  import type { Snippet } from 'svelte';
  import '../app.css';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { page } from '$app/state';
  import { checkAuthentication, AUTH_LOGIN_URL } from '$lib/auth';

  let { children }: { children: Snippet } = $props();

  let isLoading = $state(true);
  let isAuthenticated = $state(false);

  const publicPaths = ['/login'];

  onMount(async () => {
    const authenticated = await checkAuthentication();
    isAuthenticated = authenticated;
    
    if (!authenticated && !publicPaths.includes(page.url.pathname)) {
      window.location.href = AUTH_LOGIN_URL;
      return;
    }
    
    isLoading = false;
  });
</script>

{#if isLoading}
<div id="app">
  <main class="loading" role="status" aria-busy="true" aria-live="polite">
    <div class="loading-spinner" aria-hidden="true"></div>
    <p>Checking authentication...</p>
  </main>
</div>
{:else}
<div id="app">
  {@render children()}
</div>
{/if}

<style>
  :global(#app) {
    min-height: 100vh;
    display: flex;
    flex-direction: column;
  }

  .loading {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 1rem;
  }

  .loading p {
    color: var(--color-text-muted);
    font-size: 0.875rem;
  }

  .loading-spinner {
    width: 2.5rem;
    height: 2.5rem;
    border: 3px solid var(--color-border);
    border-top-color: var(--color-accent);
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }
</style>
