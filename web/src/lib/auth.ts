const AUTH_CHECK_URL = '/auth/check';
export const AUTH_LOGIN_URL = '/auth/login';

export async function checkAuthentication(): Promise<boolean> {
    try {
        const response = await fetch(AUTH_CHECK_URL, {
            method: 'GET',
            credentials: 'include',
        });
        return response.ok;
    } catch {
        return false;
    }
}
