<template>
  <header class="mall-header">
    <div class="mall-container header-inner">
      <router-link to="/" class="logo">Gin Mall</router-link>
      <div class="search-box">
        <el-input v-model="keyword" placeholder="搜索商品" clearable @keyup.enter="onSearch" size="large">
          <template #append>
            <el-button @click="onSearch"><el-icon><Search /></el-icon></el-button>
          </template>
        </el-input>
      </div>
      <div class="header-actions">
        <template v-if="auth.token">
          <el-button text @click="$router.push('/create-product')" style="color:#fff;font-size:13px">发布商品</el-button>
          <el-dropdown trigger="click">
            <span class="user-btn">
              <el-avatar :size="32" :src="auth.user?.avatar" />
              <span class="username">{{ auth.user?.nickname || auth.user?.username }}</span>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="$router.push('/profile')">个人中心</el-dropdown-item>
                <el-dropdown-item @click="$router.push('/orders')">我的订单</el-dropdown-item>
                <el-dropdown-item @click="$router.push('/favorites')">我的收藏</el-dropdown-item>
                <el-dropdown-item @click="$router.push('/addresses')">收货地址</el-dropdown-item>
                <el-dropdown-item @click="$router.push('/ai-chat')">AI 客服</el-dropdown-item>
                <el-dropdown-item divided @click="handleLogout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
          <el-badge :value="cartStore.totalCount" :hidden="cartStore.totalCount === 0" class="cart-badge">
            <el-button text @click="$router.push('/cart')" style="color:#fff;font-size:20px">
              <el-icon><ShoppingCart /></el-icon>
            </el-button>
          </el-badge>
        </template>
        <template v-else>
          <el-button text @click="$router.push('/login')" style="color:#fff">登录</el-button>
          <el-button text @click="$router.push('/register')" style="color:#fff">注册</el-button>
        </template>
      </div>
    </div>
    <div class="category-nav mall-container" v-if="categories.length">
      <router-link to="/products" class="cat-link" :class="{ active: !$route.query.category_id }">全部商品</router-link>
      <router-link v-for="cat in categories" :key="cat.id" :to="`/products?category_id=${cat.id}`" class="cat-link" :class="{ active: $route.query.category_id == cat.id }">
        {{ cat.category_name }}
      </router-link>
    </div>
  </header>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useCartStore } from '../stores/cart'
import { getCategories } from '../api/product'

const router = useRouter()
const auth = useAuthStore()
const cartStore = useCartStore()
const keyword = ref('')
const categories = ref([])

function onSearch() {
  if (keyword.value.trim()) {
    router.push({ path: '/products', query: { info: keyword.value.trim() } })
  }
}

function handleLogout() {
  auth.logout()
  cartStore.items = []
  router.push('/')
}

onMounted(async () => {
  try {
    const res = await getCategories()
    categories.value = res.data?.item || []
  } catch (e) { /* ignore */ }
  if (auth.token) cartStore.fetchCart()
})
</script>

<style scoped>
.mall-header { background: linear-gradient(135deg, #ff4d4f 0%, #ff7a45 100%); position: sticky; top: 0; z-index: 100; }
.header-inner { display: flex; align-items: center; height: 64px; gap: 24px; }
.logo { font-size: 24px; font-weight: bold; color: #fff; white-space: nowrap; }
.search-box { flex: 1; max-width: 500px; }
.search-box :deep(.el-input__wrapper) { border-radius: 20px; }
.header-actions { display: flex; align-items: center; gap: 12px; }
.user-btn { display: flex; align-items: center; gap: 8px; color: #fff; cursor: pointer; }
.username { font-size: 14px; max-width: 80px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.cart-badge :deep(.el-badge__content) { background-color: #fff; color: var(--primary); }
.category-nav { display: flex; gap: 4px; padding: 8px 0; overflow-x: auto; }
.cat-link { padding: 4px 14px; border-radius: 16px; font-size: 13px; color: rgba(255,255,255,.85); white-space: nowrap; transition: all .2s; }
.cat-link:hover, .cat-link.active { background: rgba(255,255,255,.25); color: #fff; }
</style>
