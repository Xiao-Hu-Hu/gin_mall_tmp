<template>
  <div>
    <!-- Carousel -->
    <div class="mall-container" style="margin-top:16px">
      <el-carousel height="400px" v-if="productStore.carousels.length">
        <el-carousel-item v-for="c in productStore.carousels" :key="c.id">
          <router-link :to="`/products/${c.product_id}`">
            <img :src="c.img_path" style="width:100%;height:400px;object-fit:cover;border-radius:8px" />
          </router-link>
        </el-carousel-item>
      </el-carousel>
    </div>

    <!-- Categories -->
    <div class="mall-container">
      <div class="section-title"><span>商品分类</span></div>
      <div class="cat-grid">
        <router-link v-for="cat in productStore.categories" :key="cat.id" :to="`/products?category_id=${cat.id}`" class="cat-item">
          <div class="cat-icon">{{ cat.category_name.charAt(0) }}</div>
          <span>{{ cat.category_name }}</span>
        </router-link>
      </div>
    </div>

    <!-- Hot Products -->
    <div class="mall-container">
      <div class="section-title">
        <span>热门推荐</span>
        <router-link to="/products" class="more">查看更多 &gt;</router-link>
      </div>
      <div class="product-grid">
        <router-link v-for="p in productStore.hotProducts" :key="p.id" :to="`/products/${p.id}`" class="product-card">
          <img :src="p.img_path" :alt="p.name" />
          <div class="info">
            <div class="name">{{ p.name }}</div>
            <div class="price-row">
              <span class="price-current">&yen;{{ formatPrice(p.discount_price) }}</span>
              <span class="price-original" v-if="p.discount_price !== p.price">&yen;{{ formatPrice(p.price) }}</span>
            </div>
          </div>
        </router-link>
      </div>
    </div>

    <!-- Rankings -->
    <div class="mall-container">
      <div class="section-title"><span>排行榜</span></div>
      <el-tabs v-model="activeRankTab" @tab-change="onTabChange">
        <el-tab-pane label="浏览排行" name="view" />
        <el-tab-pane label="销量排行" name="purchase" />
        <el-tab-pane label="收藏排行" name="favorite" />
      </el-tabs>
      <div class="rank-list" v-if="rankData.length">
        <router-link v-for="(p, idx) in rankData" :key="p.id" :to="`/products/${p.id}`" class="rank-item">
          <span class="rank-num" :class="{ top3: idx < 3 }">{{ idx + 1 }}</span>
          <img :src="p.img_path" class="rank-img" />
          <div class="rank-info">
            <div class="name">{{ p.name }}</div>
            <span class="price-current">&yen;{{ formatPrice(p.discount_price) }}</span>
          </div>
        </router-link>
      </div>
      <el-empty v-else description="暂无排行数据" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useProductStore } from '../stores/product'
import { getViewRank, getPurchaseRank, getFavoriteRank } from '../api/product'
import { formatPrice } from '../utils/format'

const productStore = useProductStore()
const activeRankTab = ref('view')
const rankData = ref([])

const rankFns = { view: getViewRank, purchase: getPurchaseRank, favorite: getFavoriteRank }

async function onTabChange(tab) {
  try {
    const res = await rankFns[tab]({ limit: 10 })
    const data = res.data
    rankData.value = Array.isArray(data) ? data : (data?.item || [])
  } catch (e) { rankData.value = [] }
}

onMounted(async () => {
  productStore.fetchCarousels()
  productStore.fetchCategories()
  productStore.fetchHotProducts()
  await onTabChange('view')
})
</script>

<style scoped>
.cat-grid { display: flex; flex-wrap: wrap; gap: 16px; margin-bottom: 30px; }
.cat-item { display: flex; flex-direction: column; align-items: center; gap: 8px; padding: 16px 20px; background: #fff; border-radius: 8px; cursor: pointer; transition: all .2s; min-width: 90px; color: var(--text); }
.cat-item:hover { box-shadow: 0 4px 12px rgba(0,0,0,.1); transform: translateY(-2px); }
.cat-icon { width: 48px; height: 48px; border-radius: 50%; background: linear-gradient(135deg, #ff7875, #ff4d4f); color: #fff; display: flex; align-items: center; justify-content: center; font-size: 20px; font-weight: bold; }
.rank-list { display: flex; flex-direction: column; gap: 8px; margin-bottom: 30px; }
.rank-item { display: flex; align-items: center; gap: 12px; padding: 10px 16px; background: #fff; border-radius: 8px; color: var(--text); transition: box-shadow .2s; }
.rank-item:hover { box-shadow: 0 2px 8px rgba(0,0,0,.08); }
.rank-num { width: 24px; height: 24px; border-radius: 4px; background: #e8e8e8; display: flex; align-items: center; justify-content: center; font-size: 13px; font-weight: bold; color: #999; }
.rank-num.top3 { background: var(--primary); color: #fff; }
.rank-img { width: 60px; height: 60px; object-fit: cover; border-radius: 4px; }
.rank-info { flex: 1; }
.rank-info .name { font-size: 14px; margin-bottom: 4px; }
</style>
