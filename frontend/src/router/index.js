import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  { path: '/', name: 'Home', component: () => import('../views/Home.vue') },
  { path: '/login', name: 'Login', component: () => import('../views/Login.vue') },
  { path: '/register', name: 'Register', component: () => import('../views/Register.vue') },
  { path: '/products', name: 'Products', component: () => import('../views/Products.vue') },
  { path: '/products/:id', name: 'ProductDetail', component: () => import('../views/ProductDetail.vue') },
  { path: '/cart', name: 'Cart', component: () => import('../views/Carts.vue'), meta: { auth: true } },
  { path: '/checkout', name: 'Checkout', component: () => import('../views/Checkout.vue'), meta: { auth: true } },
  { path: '/orders', name: 'Orders', component: () => import('../views/Orders.vue'), meta: { auth: true } },
  { path: '/orders/:id', name: 'OrderDetail', component: () => import('../views/OrderDetail.vue'), meta: { auth: true } },
  { path: '/addresses', name: 'Addresses', component: () => import('../views/Addresses.vue'), meta: { auth: true } },
  { path: '/favorites', name: 'Favorites', component: () => import('../views/Favorites.vue'), meta: { auth: true } },
  { path: '/profile', name: 'Profile', component: () => import('../views/User.vue'), meta: { auth: true } },
  { path: '/ai-chat', name: 'AiChat', component: () => import('../views/AiChat.vue'), meta: { auth: true } },
  { path: '/seckill', name: 'Seckill', component: () => import('../views/Seckill.vue') },
  { path: '/create-product', name: 'CreateProduct', component: () => import('../views/CreateProduct.vue'), meta: { auth: true } },
]

const router = createRouter({ history: createWebHistory(), routes })

router.beforeEach((to, from, next) => {
  const auth = useAuthStore()
  if (to.meta.auth && !auth.token) {
    next({ path: '/login', query: { redirect: to.fullPath } })
  } else if (['/login', '/register'].includes(to.path) && auth.token) {
    next('/')
  } else {
    next()
  }
})

export default router
