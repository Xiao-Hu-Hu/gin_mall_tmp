import api from './index'

export const getProducts = (params) => api.get('/products', { params })
export const getProduct = (id) => api.get(`/products/${id}`)
export const searchProducts = (data) => api.post('/products', data)
export const createProduct = (data) => {
  const fd = new FormData()
  for (const [k, v] of Object.entries(data)) {
    if (k === 'file') {
      if (Array.isArray(v)) v.forEach(f => fd.append('file', f))
    } else if (v !== undefined && v !== null) {
      fd.append(k, v)
    }
  }
  return api.post('/product', fd, { headers: { 'Content-Type': 'multipart/form-data' } })
}
export const getProductImgs = (id) => api.get(`/imgs/${id}`)
export const getCategories = () => api.get('/categories')
export const getCarousels = () => api.get('/carousels')
export const getHotProducts = (params) => api.get('/products/hot', { params })
export const getViewRank = (params) => api.get('/products/view-rank', { params })
export const getPurchaseRank = (params) => api.get('/products/purchase-rank', { params })
export const getFavoriteRank = (params) => api.get('/products/favorite-rank', { params })
