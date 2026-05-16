import api from './index'

export const payOrder = (data) => api.post('/paydown', data)
