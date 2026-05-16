import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getCategories, getHotProducts, getCarousels } from '../api/product'

export const useProductStore = defineStore('product', () => {
  const categories = ref([])
  const hotProducts = ref([])
  const carousels = ref([])

  async function fetchCategories() {
    if (categories.value.length) return
    try { const res = await getCategories(); categories.value = res.data?.item || [] } catch (e) { /* ignore */ }
  }

  async function fetchHotProducts(limit = 8) {
    try {
      const res = await getHotProducts({ limit })
      const data = res.data
      hotProducts.value = Array.isArray(data) ? data : (data?.item || [])
    } catch (e) { /* ignore */ }
  }

  async function fetchCarousels() {
    if (carousels.value.length) return
    try { const res = await getCarousels(); carousels.value = res.data?.item || [] } catch (e) { /* ignore */ }
  }

  return { categories, hotProducts, carousels, fetchCategories, fetchHotProducts, fetchCarousels }
})
