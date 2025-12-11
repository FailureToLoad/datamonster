<template>
  <main class="error">
    <div class="error-content">
      <h1>{{ status }}</h1>
      <p class="message">{{ displayMessage }}</p>
      <RouterLink to="/" class="home-link">Return Home</RouterLink>
    </div>
  </main>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, RouterLink } from 'vue-router'

const route = useRoute()

const errorMessages: Record<number, string> = {
  404: 'Page not found',
  500: 'Internal server error',
  503: 'Service unavailable',
}

const status = computed(() => Number(route.params.status) || 404)

const displayMessage = computed(() => {
  const isDev = import.meta.env.DEV
  if (isDev && route.query.message) {
    return route.query.message as string
  }
  return errorMessages[status.value] ?? 'An error occurred'
})
</script>

<style scoped>
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
