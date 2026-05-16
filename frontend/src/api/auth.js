import api from './index'

export const userRegister = (data) => api.post('/user/register', data)
export const userLogin = (data) => api.post('/user/login', data)
export const updateUser = (data) => api.put('/user', data)
export const uploadAvatar = (file) => {
  const fd = new FormData()
  fd.append('file', file)
  return api.post('/avatar', fd, { headers: { 'Content-Type': 'multipart/form-data' } })
}
export const sendEmail = (data) => api.post('/user/sending-email', data)
export const validEmail = () => api.post('/user/valid-email')
export const showMoney = (data) => api.post('/money', data)
