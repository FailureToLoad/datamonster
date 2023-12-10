import axios from "axios";
const ax = axios.create({
  headers: {
    "Accept-Language": "en-US,en;q=0.5",
  },
  validateStatus: (status) => {
    return status >= 200 && status < 500;
  },
});
ax.defaults.withCredentials = true;
ax.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    if (error.response.status === 401) {
      localStorage.removeItem("iid");
    }
    return Promise.reject(error);
  },
);

export function setInterceptor(token: string) {
  const id = ax.interceptors.request.use((config) => {
    config.headers.Authorization = `Bearer ${token}`;
    return config;
  });
  return id;
}

export default ax;
