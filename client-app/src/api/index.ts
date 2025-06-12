import axios from 'axios';

// Создаем экземпляр axios
const api = axios.create({
  baseURL: (import.meta.env.K8S_RUN) ? "http://localhost:8080/api/v1" : "http://api.dengovie.ingress/api/v1",
  withCredentials: true, // Включаем передачу cookies (включая httpOnly)
});

// Интерцептор запросов
api.interceptors.request.use((config) => {
  config.withCredentials = true;
  return config;
}, (error) => {
  return Promise.reject(error);
});

export default api;