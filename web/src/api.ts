import axios from "axios";

function makeAuthHeader(token: string) {
  return {
    headers: { Authorization: `Bearer ${token}` },
  };
}

export async function Get<T>(path: string, token: string) {
  const config = makeAuthHeader(token);
  return axios.get<T>(path, config);
}

export async function Post<T>(path: string, body: any, token: string) {
  const config = makeAuthHeader(token);
  return axios.post<T>(path, body, config);
}
