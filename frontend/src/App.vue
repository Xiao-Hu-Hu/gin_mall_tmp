<template>
  <div id="app">
    <header class="header">
      <div class="header-content">
        <router-link to="/" class="logo">个人商城</router-link>
        <nav class="nav">
          <router-link to="/">首页</router-link>
          <router-link to="/products">商品</router-link>
          <template v-if="userStore.isLoggedIn">
            <router-link to="/cart">购物车</router-link>
            <router-link to="/orders">订单</router-link>
            <router-link to="/sell">上架商品</router-link>
            <router-link to="/profile">个人中心</router-link>
            <span>{{ userStore.user?.nick_name || userStore.user?.user_name }}</span>
            <button @click="handleLogout" class="btn btn-secondary">退出</button>
          </template>
          <template v-else>
            <router-link to="/login">登录</router-link>
            <router-link to="/register">注册</router-link>
          </template>
        </nav>
      </div>
    </header>
    <main>
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { useUserStore } from './stores/user'
import { useRouter } from 'vue-router'

const userStore = useUserStore()
const router = useRouter()

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}
</script>

