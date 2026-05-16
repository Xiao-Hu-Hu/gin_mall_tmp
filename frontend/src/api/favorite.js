import api from './index'

export const getFavorites = () => api.get('/favorites')
export const addFavorite = (data) => api.post('/favorites', data)
export const deleteFavorite = (id) => api.delete(`/favorites/${id}`)
