import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import './style.css'

// 导入路由配置
import router from './router/index.js'

// 创建Pinia实例
const pinia = createPinia()

// 创建Vue应用
const app = createApp(App)

// 使用路由和状态管理
app.use(router)
app.use(pinia)

// 挂载应用
app.mount('#app')