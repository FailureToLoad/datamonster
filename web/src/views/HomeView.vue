<template>
  <main>
    <section class="hero">
      <div class="hero-content">
        <h1>Datamonster</h1>
        <p class="tagline">Homelab Data Platform</p>
      </div>

      <div class="actions">
        <button v-if="isAuthenticated" class="btn btn-primary" disabled>Add Settlement</button>
        <a v-else :href="AUTH_LOGIN_URL" class="btn btn-primary">Login</a>
      </div>

      <div class="status">
        <span class="status-dot"></span>
        <span class="status-text">System Online</span>
      </div>
    </section>

    <footer>
      <p>&copy; {{ currentYear }} Datamonster</p>
    </footer>
  </main>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAuth, AUTH_LOGIN_URL } from '@/composables/useAuth'

const { isAuthenticated } = useAuth()
const currentYear = computed(() => new Date().getFullYear())
</script>

<style scoped>
main {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.hero {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  gap: 3rem;
}

.hero-content {
  text-align: center;
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
  margin: 0.75rem 0 0;
  font-size: 1.25rem;
  color: var(--color-text-muted);
  letter-spacing: 0.02em;
}

.actions {
  display: flex;
  gap: 1rem;
}

.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.75rem 1.5rem;
  font-size: 1rem;
  font-weight: 500;
  text-decoration: none;
  border-radius: 0.5rem;
  transition: all 0.2s ease;
}

.btn-primary {
  background: var(--color-accent);
  color: var(--color-bg);
}

.btn-primary:hover {
  box-shadow: 0 0 20px var(--color-accent-glow);
}

.status {
  display: flex;
  align-items: center;
  gap: 0.625rem;
  padding: 0.75rem 1.5rem;
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-border);
  border-radius: 9999px;
}

.status-dot {
  width: 0.5rem;
  height: 0.5rem;
  background: var(--color-success);
  border-radius: 50%;
  box-shadow: 0 0 8px var(--color-success);
  animation: pulse 2s ease-in-out infinite;
}

.status-text {
  font-size: 0.875rem;
  color: var(--color-text-muted);
}

@keyframes pulse {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

footer {
  padding: 1.5rem;
  text-align: center;
  border-top: 1px solid var(--color-border);
}

footer p {
  margin: 0;
  font-size: 0.875rem;
  color: var(--color-text-muted);
}
</style>
