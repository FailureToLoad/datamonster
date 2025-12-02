import { goto } from '$app/navigation';

let isAuthenticated = $state(false);
let isLoading = $state(true);

export async function checkAuth(): Promise<boolean> {
	try {
		const response = await fetch('/auth/check', { credentials: 'include' });
		isAuthenticated = response.ok;
	} catch {
		isAuthenticated = false;
	}
	isLoading = false;
	return isAuthenticated;
}

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

export function login() {
	window.location.href = '/auth/login';
}

export function logout() {
	window.location.href = '/auth/logout';
}

export async function protectedFetch(url: string, options: RequestInit = {}): Promise<Response> {
	const response = await fetch(url, {
		...options,
		credentials: 'include'
	});

	if (response.status === 401) {
		isAuthenticated = false;
		goto('/login');
	}

	return response;
}
