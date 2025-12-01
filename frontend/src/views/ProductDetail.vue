<template>
  <div class="product-detail">
    <div class="container">
      <div v-if="loading" class="loading">加载中...</div>
      <div v-else-if="product" class="product-detail-content">
        <div class="product-images">
          <img :src="product.img_path" :alt="product.name" />
          <div v-if="productImages.length > 0" class="product-images-list">
            <img v-for="(img, index) in productImages" :key="index" :src="img.img_path" :alt="product.name" />
          </div>
        </div>
        <div class="product-info">
          <h1>{{ product.name }}</h1>
          <div class="product-price">
            <span class="current-price">¥{{ product.discount_price || product.price }}</span>
            <span v-if="product.discount_price" class="original-price">¥{{ product.price }}</span>
          </div>
          <div class="product-meta">
            <p><strong>商品描述:</strong> {{ product.info }}</p>
            <p><strong>库存:</strong> {{ product.num }}</p>
            <p><strong>商家:</strong> {{ product.boss_name }}</p>
          </div>
          <div class="product-actions">
            <div class="form-group">
              <label>数量</label>
              <input v-model.number="quantity" type="number" min="1" :max="product.num" style="width: 100px;" />
            </div>
            <button @click="addToCart" class="btn btn-primary" :disabled="!userStore.isLoggedIn">
              加入购物车
            </button>
            <button @click="buyNow" class="btn btn-success" :disabled="!userStore.isLoggedIn">
              立即购买
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import api from '../utils/api'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const product = ref(null)
const productImages = ref([])
const loading = ref(true)
const quantity = ref(1)

const fetchProduct = async () => {
  try {
    const res = await api.get(`/products/${route.params.id}`)
    if (res.status === 200) {
      product.value = res.data
      fetchProductImages()
    }
  } catch (error) {
    console.error('获取商品详情失败:', error)
  } finally {
    loading.value = false
  }
}

const fetchProductImages = async () => {
  try {
    const res = await api.get(`/imgs/${route.params.id}`)
    if (res.status === 200) {
      productImages.value = res.data || []
    }
  } catch (error) {
    console.error('获取商品图片失败:', error)
  }
}

const addToCart = async () => {
  if (!userStore.isLoggedIn) {
    router.push('/login')
    return
  }
  
  try {
    const res = await api.post('/carts', {
      product_id: product.value.id,
      boss_id: product.value.boss_id,
      num: quantity.value
    })
    if (res.status === 200) {
      alert('已加入购物车')
    }
  } catch (error) {
    alert(error.msg || error.error || '加入购物车失败')
  }
}

const buyNow = async () => {
  if (!userStore.isLoggedIn) {
    router.push('/login')
    return
  }
  
  router.push({
    path: '/orders',
    query: {
      product_id: product.value.id,
      boss_id: product.value.boss_id,
      num: quantity.value
    }
  })
}

onMounted(() => {
  fetchProduct()
})
</script>

<style scoped>
.product-detail-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 40px;
  margin-top: 20px;
}

.product-images img {
  width: 100%;
  max-height: 500px;
  object-fit: contain;
  border-radius: 8px;
}

.product-images-list {
  display: flex;
  gap: 10px;
  margin-top: 10px;
}

.product-images-list img {
  width: 100px;
  height: 100px;
  object-fit: cover;
  border-radius: 4px;
  cursor: pointer;
}

.product-info h1 {
  margin-bottom: 20px;
}

.product-price {
  margin-bottom: 20px;
}

.current-price {
  font-size: 32px;
  font-weight: bold;
  color: #dc3545;
}

.original-price {
  font-size: 18px;
  color: #999;
  text-decoration: line-through;
  margin-left: 10px;
}

.product-meta {
  margin-bottom: 30px;
  line-height: 1.8;
}

.product-actions {
  display: flex;
  gap: 15px;
  align-items: flex-end;
}
</style>

