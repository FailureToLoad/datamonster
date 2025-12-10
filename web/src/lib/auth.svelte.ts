const AUTH_CHECK_URL = '/auth/check';
export const AUTH_LOGIN_URL = '/auth/login';

let isAuthenticated = $state(false);
let isLoading = $state(true);

export function getAuthState() {
    return {
        get isAuthenticated() {
            return isAuthenticated;
        },
        get isLoading() {
            return isLoading;
        }
    };
}

export async function checkAuthentication(): Promise<boolean> {
    isLoading = true;
    try {
        const response = await fetch(AUTH_CHECK_URL, {
            method: 'GET',
            credentials: 'include',
        });
        isAuthenticated = response.ok;
        return response.ok;
    } catch {
        isAuthenticated = false;
        return false;
    } finally {
        isLoading = false;
    }
}
