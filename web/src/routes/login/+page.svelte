<script lang="ts">
  import { login } from "$lib/auth.svelte";
</script>

<main>
  <section class="login-container">
    <div class="login-content">
      <h1>Datamonster</h1>
      <p class="tagline">Homelab Data Platform</p>
      <button onclick={login} class="login-button">Sign in with Keycloak</button
      >
    </div>
  </section>
</main>

<style>
  main {
    flex: 1;
    display: flex;
    flex-direction: column;
  }

  .login-container {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 2rem;
  }

  .login-content {
    text-align: center;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1.5rem;
  }

  h1 {
    margin: 0;
    font-size: clamp(2.5rem, 8vw, 4.5rem);
    font-weight: 700;
    background: linear-gradient(135deg, var(--color-accent) 0%, #a78bfa 100%);
    background-clip: text;
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }

  .tagline {
    margin: 0;
    font-size: 1.25rem;
    color: var(--color-text-muted);
    letter-spacing: 0.02em;
  }

  .login-button {
    margin-top: 1rem;
    padding: 0.875rem 2rem;
    font-size: 1rem;
    font-weight: 500;
    color: var(--color-bg);
    background: var(--color-accent);
    border: none;
    border-radius: 0.5rem;
    cursor: pointer;
    transition:
      opacity 0.2s ease,
      transform 0.2s ease;
  }

  .login-button:hover {
    opacity: 0.9;
    transform: translateY(-1px);
  }

  .login-button:active {
    transform: translateY(0);
  }
</style>
