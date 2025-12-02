<script lang="ts">
  import "../app.css";
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import { checkAuth, getAuthState } from "$lib/auth.svelte";

  const auth = getAuthState();

  let { children } = $props();

  onMount(async () => {
    const authenticated = await checkAuth();
    const currentPath = page.url.pathname;

    if (!authenticated && currentPath !== "/login") {
      goto("/login");
    } else if (authenticated && currentPath === "/login") {
      goto("/");
    }
  });
</script>

{#if auth.isLoading}
  <main class="loading">
    <div class="spinner"></div>
  </main>
{:else}
  {@render children()}
{/if}

<style>
  main.loading {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .spinner {
    width: 2rem;
    height: 2rem;
    border: 2px solid var(--color-border);
    border-top-color: var(--color-accent);
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }
</style>
