import axios from 'axios'
import { useAuthStore } from '../stores/auth'
import { ElMessage } from 'element-plus'

const api = axios.create({
  baseURL: '/api/v1',
  timeout: 30000
})

api.interceptors.request.use(config => {
  const auth = useAuthStore()
  if (auth.token) {
    config.headers['token'] = auth.token
  }
  return config
})

api.interceptors.response.use(
  res => {
    const data = res.data
    if (data.status && data.status !== 200) {
      ElMessage.error(data.error || data.msg || '请求失败')
      return Promise.reject(data)
    }
    return data
  },
  err => {
    if (!err.config?.suppressErrorToast) {
      ElMessage.error(err.message || '网络错误')
    }
    return Promise.reject(err)
  }
)

export default api
