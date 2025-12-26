export function Get(url: string): Promise<Response> {
  return fetch(url, { credentials: 'include'})
}

export function PostJSON(url: string, data:object): Promise<Response> {
  return fetch(url, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
        credentials: "include",
      })
}

export function PatchJSON(url: string, data: Record<string, number>): Promise<Response> {
  return fetch(url, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
        credentials: "include",
      })
}