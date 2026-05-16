<template>
  <div class="mall-container detail-page">
    <el-skeleton :loading="loading" animated>
      <template #template>
        <div style="display:flex;gap:40px;padding:20px 0"><el-skeleton-item variant="image" style="width:400px;height:400px" /><div style="flex:1"><el-skeleton-item variant="h1" style="width:60%;height:32px;margin-bottom:16px" /><el-skeleton-item variant="text" style="width:80%;height:20px;margin-bottom:8px" /><el-skeleton-item variant="text" style="width:40%;height:40px" /></div></div>
      </template>
      <template #default>
        <div v-if="product" class="detail-layout">
          <!-- Images -->
          <div class="image-section">
            <img :src="activeImg" class="main-img" />
            <div class="thumb-list" v-if="productImgs.length > 1">
              <img v-for="(img, idx) in productImgs" :key="idx" :src="img.img_path" :class="{ active: activeImgIdx === idx }" @click="activeImgIdx = idx" class="thumb" />
            </div>
          </div>
          <!-- Info -->
          <div class="info-section">
            <h1 class="title">{{ product.title || product.name }}</h1>
            <div class="price-box">
              <span class="price-current">&yen;{{ formatPrice(product.discount_price) }}</span>
              <span class="price-original" v-if="product.discount_price !== product.price">&yen;{{ formatPrice(product.price) }}</span>
            </div>
            <div class="meta"><span>浏览 {{ product.view || 0 }} 次</span><span>库存 {{ product.num }} 件</span></div>
            <div class="desc" v-if="product.info">{{ product.info }}</div>
            <div class="seller" v-if="product.boss_name">
              <img :src="product.boss_avatar" class="seller-avatar" v-if="product.boss_avatar" />
              <span>{{ product.boss_name }}</span>
            </div>
            <div class="quantity-row">
              <span>数量</span>
              <el-input-number v-model="quantity" :min="1" :max="product.num || 99" />
            </div>
            <div class="actions">
              <el-button type="danger" size="large" @click="buyNow">立即购买</el-button>
              <el-button type="primary" size="large" @click="addToCart">加入购物车</el-button>
              <el-button size="large" @click="toggleFav" :icon="isFav ? 'StarFilled' : 'Star'" :type="isFav ? 'warning' : 'default'">
                {{ isFav ? '已收藏' : '收藏' }}
              </el-button>
            </div>
          </div>
        </div>
        <!-- Detail images -->
        <div v-if="productImgs.length" class="detail-images">
          <el-divider content-position="left"><span style="font-size:18px;font-weight:bold">商品详情</span></el-divider>
          <img v-for="(img, idx) in productImgs" :key="idx" :src="img.img_path" style="width:100%;display:block" />
        </div>
      </template>
    </el-skeleton>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useCartStore } from '../stores/cart'
import { getProduct, getProductImgs } from '../api/product'
import { addFavorite, deleteFavorite, getFavorites } from '../api/favorite'
import { formatPrice } from '../utils/format'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const cartStore = useCartStore()
const product = ref(null)
const productImgs = ref([])
const loading = ref(true)
const quantity = ref(1)
const activeImgIdx = ref(0)
const isFav = ref(false)
const favId = ref(null)

const activeImg = computed(() => productImgs.value.length ? productImgs.value[activeImgIdx.value]?.img_path : product.value?.img_path)

async function loadData() {
  loading.value = true
  const id = route.params.id
  try {
    const [pRes, imgRes] = await Promise.all([getProduct(id), getProductImgs(id)])
    product.value = pRes.data
    productImgs.value = imgRes.data?.item || []
    if (product.value?.img_path && !productImgs.value.length) productImgs.value = [{ img_path: product.value.img_path }]
  } catch (e) { /* ignore */ } finally { loading.value = false }
}

async function checkFav() {
  if (!auth.token) return
  try {
    const res = await getFavorites()
    const list = res.data?.item || []
    const found = list.find(f => f.product_id === Number(route.params.id))
    if (found) { isFav.value = true; favId.value = found.id }
  } catch (e) { /* ignore */ }
}

async function addToCart() {
  if (!auth.token) return router.push({ path: '/login', query: { redirect: route.fullPath } })
  await cartStore.addItem(product.value.id, product.value.boss_id, quantity.value)
  ElMessage.success('已加入购物车')
}

function buyNow() {
  if (!auth.token) return router.push({ path: '/login', query: { redirect: route.fullPath } })
  router.push(`/checkout?productId=${product.value.id}&num=${quantity.value}&bossId=${product.value.boss_id}`)
}

async function toggleFav() {
  if (!auth.token) return router.push({ path: '/login', query: { redirect: route.fullPath } })
  if (isFav.value) {
    await deleteFavorite(favId.value)
    isFav.value = false; favId.value = null
    ElMessage.success('已取消收藏')
  } else {
    const res = await addFavorite({ product_id: product.value.id, boss_id: product.value.boss_id })
    isFav.value = true; favId.value = res.data?.id
    ElMessage.success('已收藏')
  }
}

onMounted(() => { loadData(); checkFav() })
</script>

<style scoped>
.detail-page { padding: 20px 0 40px; }
.detail-layout { display: flex; gap: 40px; }
.image-section { width: 400px; flex-shrink: 0; }
.main-img { width: 400px; height: 400px; object-fit: cover; border-radius: 8px; background: #f5f5f5; }
.thumb-list { display: flex; gap: 8px; margin-top: 10px; overflow-x: auto; }
.thumb { width: 70px; height: 70px; object-fit: cover; border-radius: 4px; cursor: pointer; border: 2px solid transparent; }
.thumb.active { border-color: var(--primary); }
.info-section { flex: 1; }
.title { font-size: 22px; font-weight: bold; margin-bottom: 16px; line-height: 1.4; }
.price-box { background: #fef0f0; padding: 16px; border-radius: 8px; margin-bottom: 16px; }
.price-current { font-size: 28px; font-weight: bold; color: var(--price); }
.price-original { font-size: 16px; color: var(--text-light); text-decoration: line-through; margin-left: 12px; }
.meta { display: flex; gap: 20px; font-size: 13px; color: var(--text-light); margin-bottom: 16px; }
.desc { font-size: 14px; color: #666; line-height: 1.8; margin-bottom: 16px; }
.seller { display: flex; align-items: center; gap: 8px; font-size: 14px; color: var(--text-light); margin-bottom: 16px; }
.seller-avatar { width: 28px; height: 28px; border-radius: 50%; }
.quantity-row { display: flex; align-items: center; gap: 12px; margin-bottom: 24px; font-size: 14px; }
.actions { display: flex; gap: 12px; }
.detail-images { margin-top: 40px; }
@media (max-width: 800px) { .detail-layout { flex-direction: column; } .image-section { width: 100%; } .main-img { width: 100%; height: auto; } }
</style>
