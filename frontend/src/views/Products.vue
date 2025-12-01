<template>
  <div class="products">
    <div class="container">
      <h1>商品列表</h1>
      <div class="search-bar" style="margin-bottom: 20px;">
        <input v-model="searchKeyword" type="text" placeholder="搜索商品..." style="width: 300px; padding: 10px;" />
        <button @click="handleSearch" class="btn btn-primary">搜索</button>
      </div>
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
            <div style="font-size: 12px; color: #666;">库存: {{ product.num }}</div>
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
const loading = ref(true)
const searchKeyword = ref('')

const normalizeImageUrl = (url) => {
  if (!url) return ''
  if (url.startsWith('http')) return url
  // 后端静态资源由 Go 服务在 3000 端口提供
  return `http://8.137.53.3:3000${url.startsWith('/static') ? url : '/static/imgs/product/' + url}`
}

const getProductImage = (url) => normalizeImageUrl(url)

const fetchProducts = async () => {
  loading.value = true
  try {
    const res = await api.get('/products', { params: { page: 1, page_size: 20 } })
    if (res.status === 200) {
      products.value = res.data?.item || []
    }
  } catch (error) {
    console.error('获取商品列表失败:', error)
  } finally {
    loading.value = false
  }
}

const handleSearch = async () => {
  if (!searchKeyword.value.trim()) {
    fetchProducts()
    return
  }
  
  loading.value = true
  try {
    const res = await api.post('/products', { info: searchKeyword.value })
    if (res.status === 200) {
      products.value = res.data?.item || []
    }
  } catch (error) {
    console.error('搜索失败:', error)
  } finally {
    loading.value = false
  }
}

const goToProduct = (id) => {
  router.push(`/product/${id}`)
}

onMounted(() => {
  fetchProducts()
})
</script>

