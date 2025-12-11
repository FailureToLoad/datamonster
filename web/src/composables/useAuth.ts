import { ref } from 'vue'

const AUTH_CHECK_URL = '/auth/check'
export const AUTH_LOGIN_URL = '/auth/login'

const isAuthenticated = ref(false)
const isLoading = ref(true)

export function useAuth() {
  async function checkAuthentication(): Promise<boolean> {
    isLoading.value = true
    try {
      const response = await fetch(AUTH_CHECK_URL, {
        method: 'GET',
        credentials: 'include',
      })
      isAuthenticated.value = response.ok
      return response.ok
    } catch {
      isAuthenticated.value = false
      return false
    } finally {
      isLoading.value = false
    }
  }

  return {
    isAuthenticated,
    isLoading,
    checkAuthentication,
  }
}
