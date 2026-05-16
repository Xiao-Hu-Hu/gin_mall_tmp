import api from './index'

export const getCarts = () => api.get('/carts')
export const getCart = (id, params) => api.get(`/carts/${id}`, { params })
export const addCart = (data) => api.post('/carts', data)
export const updateCart = (id, data) => api.put(`/carts/${id}`, data)
export const deleteCart = (id) => api.delete(`/carts/${id}`)
