export async function checkAuth(): Promise<boolean> {
  const response = await fetch("/api/me",);
  return response.ok;
}
