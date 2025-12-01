<template>
  <div class="cart">
    <div class="container">
      <h1>购物车</h1>
      <div v-if="loading" class="loading">加载中...</div>
      <div v-else-if="carts.length === 0" class="empty">购物车为空</div>
      <div v-else>
        <div v-for="cart in carts" :key="cart.id" class="card">
          <div style="display: flex; gap: 20px; align-items: center;">
            <img :src="cart.product.img_path" style="width: 100px; height: 100px; object-fit: cover; border-radius: 4px;" />
            <div style="flex: 1;">
              <h3>{{ cart.product.name }}</h3>
              <p>价格: ¥{{ cart.product.discount_price || cart.product.price }}</p>
              <p>数量: {{ cart.num }}</p>
              <p>小计: ¥{{ (parseFloat(cart.product.discount_price || cart.product.price) * cart.num).toFixed(2) }}</p>
            </div>
            <div>
              <input v-model.number="cart.num" type="number" min="1" style="width: 60px; margin-right: 10px;" />
              <button @click="updateCart(cart)" class="btn btn-primary">更新</button>
              <button @click="deleteCart(cart.id)" class="btn btn-danger">删除</button>
            </div>
          </div>
        </div>
        <div class="card" style="text-align: right;">
          <h2>总计: ¥{{ totalPrice.toFixed(2) }}</h2>
          <button @click="checkout" class="btn btn-success" style="margin-top: 10px;">结算</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '../utils/api'

const router = useRouter()
const carts = ref([])
const loading = ref(true)

const totalPrice = computed(() => {
  return carts.value.reduce((sum, cart) => {
    const price = parseFloat(cart.product.discount_price || cart.product.price)
    return sum + price * cart.num
  }, 0)
})

const fetchCarts = async () => {
  loading.value = true
  try {
    const res = await api.post('/carts')
    if (res.status === 200) {
      carts.value = res.data?.item || []
    }
  } catch (error) {
    console.error('获取购物车失败:', error)
  } finally {
    loading.value = false
  }
}

const updateCart = async (cart) => {
  try {
    const res = await api.put(`/carts/${cart.id}`, { num: cart.num })
    if (res.status === 200) {
      alert('更新成功')
      fetchCarts()
    }
  } catch (error) {
    alert(error.msg || error.error || '更新失败')
  }
}

const deleteCart = async (id) => {
  if (!confirm('确定要删除吗？')) return
  
  try {
    const res = await api.delete(`/carts/${id}`)
    if (res.status === 200) {
      alert('删除成功')
      fetchCarts()
    }
  } catch (error) {
    alert(error.msg || error.error || '删除失败')
  }
}

const checkout = () => {
  if (carts.value.length === 0) {
    alert('购物车为空')
    return
  }
  router.push('/orders')
}

onMounted(() => {
  fetchCarts()
})
</script>

