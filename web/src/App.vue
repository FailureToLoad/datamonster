<template>
  <div v-if="isLoading" id="app">
    <main class="loading" role="status" aria-busy="true" aria-live="polite">
      <div class="loading-spinner" aria-hidden="true"></div>
      <p>Checking authentication...</p>
    </main>
  </div>
  <div v-else id="app">
    <RouterView />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { RouterView } from 'vue-router'
import { useAuth } from './composables/useAuth'

const { isLoading, checkAuthentication } = useAuth()

onMounted(() => {
  checkAuthentication()
})
</script>

<style scoped>
#app {
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
