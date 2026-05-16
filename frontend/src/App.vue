<template>
  <div class="mall-app">
    <MallHeader />
    <main class="mall-main">
      <router-view />
    </main>
    <MallFooter />
    <div class="ai-float-btn" @click="$router.push('/ai-chat')" title="AI 客服">
      <el-icon :size="28"><ChatDotRound /></el-icon>
    </div>
  </div>
</template>

<script setup>
import MallHeader from './components/MallHeader.vue'
import MallFooter from './components/MallFooter.vue'
import { ChatDotRound } from '@element-plus/icons-vue'
</script>

<style>
.mall-main { min-height: calc(100vh - 64px - 76px); }
.ai-float-btn { position: fixed; right: 24px; bottom: 80px; width: 52px; height: 52px; border-radius: 50%; background: linear-gradient(135deg, #ff4d4f, #ff7a45); color: #fff; display: flex; align-items: center; justify-content: center; cursor: pointer; box-shadow: 0 4px 16px rgba(255,77,79,.4); transition: transform .2s; z-index: 999; }
.ai-float-btn:hover { transform: scale(1.1); }
</style>
