<template>
  <div class="mall-container favorites-page">
    <h2 class="page-title">我的收藏</h2>
    <div v-if="favorites.length" class="product-grid">
      <div v-for="fav in favorites" :key="fav.product_id" class="product-card">
        <router-link :to="`/products/${fav.product_id}`">
          <img :src="fav.img_path" :alt="fav.name" />
        </router-link>
        <div class="info">
          <router-link :to="`/products/${fav.product_id}`" class="name">{{ fav.name }}</router-link>
          <div class="price-row">
            <span class="price-current">&yen;{{ formatPrice(fav.discount_price) }}</span>
            <span class="price-original" v-if="fav.discount_price !== fav.price">&yen;{{ formatPrice(fav.price) }}</span>
          </div>
          <div class="fav-actions">
            <el-button text type="primary" size="small" @click="addToCart(fav)">加入购物车</el-button>
            <el-button text type="danger" size="small" @click="handleRemove(fav)">取消收藏</el-button>
          </div>
        </div>
      </div>
    </div>
    <el-empty v-else description="暂无收藏商品">
      <el-button type="primary" @click="$router.push('/')">去逛逛</el-button>
    </el-empty>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useCartStore } from '../stores/cart'
import { getFavorites, deleteFavorite } from '../api/favorite'
import { formatPrice } from '../utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'

const cartStore = useCartStore()
const favorites = ref([])

async function loadFavorites() { try { const res = await getFavorites(); favorites.value = res.data?.item || [] } catch (e) { favorites.value = [] } }

async function addToCart(fav) {
  await cartStore.addItem(fav.product_id, fav.boss_id, 1)
  ElMessage.success('已加入购物车')
}

function handleRemove(fav) {
  ElMessageBox.confirm('确定取消收藏？', '提示', { type: 'warning' }).then(async () => {
    await deleteFavorite(fav.id)
    loadFavorites()
  }).catch(() => {})
}

onMounted(loadFavorites)
</script>

<style scoped>
.favorites-page { padding: 20px 0 40px; }
.page-title { font-size: 22px; margin-bottom: 20px; }
.fav-actions { display: flex; gap: 4px; margin-top: 8px; }
.name { font-size: 14px; color: var(--text); display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; line-height: 1.4; height: 2.8em; margin-bottom: 6px; }
</style>
