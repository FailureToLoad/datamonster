<script lang="ts">
  import { page } from '$app/state';
</script>

<main class="error">
  <div class="error-content">
    <h1>{page.status}</h1>
    <p class="message">{page.error?.message ?? 'An unexpected error occurred'}</p>
    <a href="/" class="home-link">Return Home</a>
  </div>
</main>

<style>
  .error {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 2rem;
  }

  .error-content {
    text-align: center;
  }

  h1 {
    margin: 0;
    font-size: clamp(4rem, 15vw, 8rem);
    font-weight: 700;
    color: var(--color-accent);
  }

  .message {
    margin: 1rem 0 2rem;
    font-size: 1.125rem;
    color: var(--color-text-muted);
  }

  .home-link {
    display: inline-block;
    padding: 0.75rem 1.5rem;
    background: var(--color-bg-elevated);
    border: 1px solid var(--color-border);
    border-radius: 0.5rem;
    color: var(--color-text);
    text-decoration: none;
    transition: border-color 0.2s ease;
  }

  .home-link:hover {
    border-color: var(--color-accent);
  }
</style>
