import api from './index'

export const createOrder = (data) => api.post('/orders', data)
export const getOrders = (params) => api.get('/orders', { params })
export const getOrder = (id, params) => api.get(`/orders/${id}`, { params })
export const deleteOrder = (id) => api.delete(`/orders/${id}`)
