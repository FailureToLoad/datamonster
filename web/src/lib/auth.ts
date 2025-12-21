export async function checkAuth(): Promise<boolean> {
  const response = await fetch("/api/me", { credentials: "include" });
  return response.ok;
}
