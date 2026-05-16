import api from './index'

export const getAddresses = () => api.get('/addresses')
export const getAddress = (id) => api.get(`/addresses/${id}`)
export const createAddress = (data) => api.post('/addresses', data)
export const updateAddress = (id, data) => api.put(`/addresses/${id}`, data)
export const deleteAddress = (id) => api.delete(`/addresses/${id}`)
