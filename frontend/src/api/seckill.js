import api from './index'

export const getSeckillActivity = (id) => api.get(`/seckill/${id}`)
export const createSeckillOrder = (data) => api.post('/seckill/order', data)
export const getSeckillOrderStatus = (activityId) => api.get(`/seckill/order/${activityId}`)
