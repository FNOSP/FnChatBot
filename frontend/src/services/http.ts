import axios from 'axios'
import { useAuthStore } from '../store/auth'

// Shared Axios instance with auth header injection.
export const http = axios.create({
  baseURL: 'http://localhost:8080/api',
  timeout: 10000, // 10s - avoid hanging when backend is down
})

http.interceptors.request.use((config) => {
  const auth = useAuthStore()
  if (auth.token) {
    // Attach bearer token for authenticated requests.
    const headers: any = config.headers ?? {}
    headers.Authorization = `Bearer ${auth.token}`
    config.headers = headers
  }
  return config
})

