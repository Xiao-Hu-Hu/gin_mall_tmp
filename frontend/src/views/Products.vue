<template>
  <div class="mall-container products-page">
    <div class="products-layout">
      <!-- Sidebar -->
      <aside class="sidebar">
        <h3>全部分类</h3>
        <div class="cat-list">
          <div class="cat-item" :class="{ active: !categoryId }" @click="setCategory(null)">全部商品</div>
          <div v-for="cat in categories" :key="cat.id" class="cat-item" :class="{ active: categoryId === cat.id }" @click="setCategory(cat.id)">
            {{ cat.category_name }}
          </div>
        </div>
      </aside>
      <!-- Main -->
      <div class="main-content">
        <div class="top-bar">
          <span v-if="searchInfo">搜索 "{{ searchInfo }}" 的结果 ({{ total }}件)</span>
          <span v-else-if="currentCatName">{{ currentCatName }} ({{ total }}件)</span>
          <span v-else>全部商品 ({{ total }}件)</span>
        </div>
        <el-skeleton :loading="loading" animated>
          <template #template>
            <div class="product-grid"><el-skeleton-item v-for="i in 8" :key="i" variant="rect" style="height:280px;border-radius:8px" /></div>
          </template>
          <template #default>
            <div v-if="products.length" class="product-grid">
              <router-link v-for="p in products" :key="p.id" :to="`/products/${p.id}`" class="product-card">
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
            <el-empty v-else description="暂无商品" />
          </template>
        </el-skeleton>
        <div class="pagination" v-if="total > pageSize">
          <el-pagination background layout="prev, pager, next" :total="total" :page-size="pageSize" :current-page="pageNum" @current-change="onPageChange" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getProducts, searchProducts } from '../api/product'
import { getCategories } from '../api/product'
import { formatPrice } from '../utils/format'

const route = useRoute()
const router = useRouter()
const products = ref([])
const categories = ref([])
const loading = ref(false)
const total = ref(0)
const pageNum = ref(1)
const pageSize = 12
const categoryId = ref(null)
const searchInfo = ref('')

const currentCatName = computed(() => categories.value.find(c => c.id === categoryId.value)?.category_name || '')

async function loadProducts() {
  loading.value = true
  try {
    let res
    if (searchInfo.value) {
      res = await searchProducts({ info: searchInfo.value, page_num: pageNum.value, page_size: pageSize })
    } else {
      const params = { page_num: pageNum.value, page_size: pageSize }
      if (categoryId.value) params.category_id = categoryId.value
      res = await getProducts(params)
    }
    products.value = res.data?.item || []
    total.value = res.data?.total || 0
  } catch (e) { products.value = []; total.value = 0 } finally { loading.value = false }
}

function setCategory(id) {
  categoryId.value = id
  pageNum.value = 1
  searchInfo.value = ''
  router.push({ query: id ? { category_id: id } : {} })
}

function onPageChange(p) { pageNum.value = p; loadProducts() }

onMounted(async () => {
  try { const res = await getCategories(); categories.value = res.data?.item || [] } catch (e) { /* ignore */ }
  if (route.query.category_id) categoryId.value = Number(route.query.category_id)
  if (route.query.info) searchInfo.value = route.query.info
  loadProducts()
})

watch(() => route.query, (q) => {
  if (q.category_id) categoryId.value = Number(q.category_id)
  else if (!q.info) categoryId.value = null
  if (q.info) searchInfo.value = q.info
  else searchInfo.value = ''
  pageNum.value = 1
  loadProducts()
})
</script>

<style scoped>
.products-page { padding-top: 20px; padding-bottom: 40px; }
.products-layout { display: flex; gap: 20px; }
.sidebar { width: 200px; flex-shrink: 0; background: #fff; border-radius: 8px; padding: 16px; height: fit-content; position: sticky; top: 120px; }
.sidebar h3 { font-size: 16px; margin-bottom: 12px; padding-bottom: 8px; border-bottom: 1px solid var(--border); }
.cat-list { display: flex; flex-direction: column; gap: 4px; }
.cat-item { padding: 8px 12px; border-radius: 6px; cursor: pointer; font-size: 14px; color: var(--text); transition: all .2s; }
.cat-item:hover { background: #fff0f0; color: var(--primary); }
.cat-item.active { background: var(--primary); color: #fff; }
.main-content { flex: 1; }
.top-bar { margin-bottom: 16px; font-size: 15px; color: var(--text-light); }
.pagination { display: flex; justify-content: center; margin-top: 24px; }
</style>
