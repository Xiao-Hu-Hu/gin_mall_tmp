import axios from 'axios'

// 获取环境变量中的API基础URL，如果没有则使用默认值
const baseURL = import.meta.env.VITE_API_BASE_URL ? 
  `${import.meta.env.VITE_API_BASE_URL}/api/v1` : 
  '/api/v1'

const api = axios.create({
  baseURL: baseURL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
api.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.token = token
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  response => {
    return response.data
  },
  error => {
    if (error.response) {
      const { status, data } = error.response
      if (status === 401) {
        localStorage.removeItem('token')
        window.location.href = '/login'
      }
      return Promise.reject(data || error)
    }
    return Promise.reject(error)
  }
)

export default api

