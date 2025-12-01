<template>
  <div class="home">
    <div class="container">
      <h1>欢迎来到胡队的袜子铺</h1>
      <div v-if="carousels.length > 0" class="carousel">
        <img v-for="(item, index) in carousels" :key="index" :src="item.img_path" alt="轮播图" />
      </div>
      <h2>热门商品</h2>
      <div v-if="loading" class="loading">加载中...</div>
      <div v-else-if="products.length === 0" class="empty">暂无商品</div>
      <div v-else class="grid">
        <div v-for="product in products" :key="product.id" class="product-card" @click="goToProduct(product.id)">
          <img :src="getProductImage(product.img_path)" :alt="product.name" />
          <div class="product-card-content">
            <div class="product-card-title">{{ product.name }}</div>
            <div class="product-card-price">
              ¥{{ product.discount_price || product.price }}
              <span v-if="product.discount_price" class="original">¥{{ product.price }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '../utils/api'

const router = useRouter()
const products = ref([])
const carousels = ref([])
const loading = ref(true)

const normalizeImageUrl = (url, isProduct = true) => {
  if (!url) return ''
  if (url.startsWith('http')) return url
  if (!isProduct) {
    return `http://8.137.53.3:3000${url.startsWith('/static') ? url : '/static/imgs/product/' + url}`
  }
  return `http://8.137.53.3:3000${url.startsWith('/static') ? url : '/static/imgs/product/' + url}`
}

const getProductImage = (url) => normalizeImageUrl(url, true)

const fetchProducts = async () => {
  try {
    const res = await api.get('/products', { params: { page: 1, page_size: 8 } })
    if (res.status === 200) {
      products.value = res.data?.item || []
    }
  } catch (error) {
    console.error('获取商品列表失败:', error)
  } finally {
    loading.value = false
  }
}

const fetchCarousels = async () => {
  try {
    const res = await api.get('/carousels')
    if (res.status === 200) {
      carousels.value = res.data || []
    }
  } catch (error) {
    console.error('获取轮播图失败:', error)
  }
}

const goToProduct = (id) => {
  router.push(`/product/${id}`)
}

onMounted(() => {
  fetchCarousels()
  fetchProducts()
})
</script>

<style scoped>
.carousel {
  margin: 20px 0;
  display: flex;
  gap: 10px;
  overflow-x: auto;
}

.carousel img {
  min-width: 300px;
  height: 200px;
  object-fit: cover;
  border-radius: 8px;
}
</style>

