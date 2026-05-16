<template>
  <div class="mall-container checkout-page">
    <h2 class="page-title">确认订单</h2>
    <!-- Address -->
    <el-card shadow="never" class="section-card">
      <template #header><span class="card-title">收货地址</span></template>
      <div v-if="selectedAddress" class="selected-addr">
        <div><strong>{{ selectedAddress.name }}</strong> {{ selectedAddress.phone }}</div>
        <div>{{ selectedAddress.address }}</div>
        <el-button text type="primary" @click="showAddrDialog = true">更换地址</el-button>
      </div>
      <div v-else>
        <el-empty description="请选择收货地址" :image-size="60">
          <el-button type="primary" @click="showAddrDialog = true">选择地址</el-button>
        </el-empty>
      </div>
    </el-card>
    <!-- Products -->
    <el-card shadow="never" class="section-card">
      <template #header><span class="card-title">商品信息</span></template>
      <div v-for="item in orderItems" :key="item.product_id" class="order-item">
        <img :src="item.img_path" />
        <div class="item-info">
          <div class="name">{{ item.name }}</div>
          <div class="price">&yen;{{ formatPrice(item.price) }} x {{ item.num }}</div>
        </div>
        <div class="subtotal">&yen;{{ formatPrice((parseFloat(item.price) * item.num).toFixed(2)) }}</div>
      </div>
    </el-card>
    <!-- Summary -->
    <el-card shadow="never" class="section-card">
      <div class="summary-row"><span>商品总额</span><span>&yen;{{ totalPrice }}</span></div>
      <div class="summary-row"><span>运费</span><span>&yen;0.00</span></div>
      <div class="summary-row total"><span>实付金额</span><span class="price-current" style="font-size:24px">&yen;{{ totalPrice }}</span></div>
      <el-button type="danger" size="large" :disabled="!selectedAddress || submitting" :loading="submitting" @click="submitOrder" style="width:200px;margin-top:16px">提交订单</el-button>
    </el-card>
    <!-- Address Dialog -->
    <el-dialog v-model="showAddrDialog" title="选择收货地址" width="500px">
      <div v-for="addr in addresses" :key="addr.id" class="addr-option" :class="{ selected: selectedAddress?.id === addr.id }" @click="selectAddr(addr)">
        <div><strong>{{ addr.name }}</strong> {{ addr.phone }}</div>
        <div style="font-size:13px;color:#666">{{ addr.address }}</div>
      </div>
      <el-empty v-if="!addresses.length" description="暂无地址" :image-size="40" />
      <el-button text type="primary" @click="$router.push('/addresses?mode=select')">管理收货地址</el-button>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useCartStore } from '../stores/cart'
import { getProduct } from '../api/product'
import { getAddresses } from '../api/address'
import { createOrder } from '../api/order'
import { formatPrice } from '../utils/format'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const cartStore = useCartStore()
const orderItems = ref([])
const addresses = ref([])
const selectedAddress = ref(null)
const showAddrDialog = ref(false)
const submitting = ref(false)

const totalPrice = computed(() => orderItems.value.reduce((s, i) => s + parseFloat(i.price) * i.num, 0).toFixed(2))

async function loadItems() {
  if (route.query.productId) {
    const res = await getProduct(route.query.productId)
    const p = res.data
    orderItems.value = [{ product_id: p.id, name: p.name, img_path: p.img_path, price: p.discount_price || p.price, num: Number(route.query.num) || 1, boss_id: p.boss_id }]
  } else {
    orderItems.value = cartStore.selectedItems.map(i => ({ product_id: i.product_id, name: i.name, img_path: i.img_path, price: i.discount_price || i.price, num: i.num, boss_id: i.boss_id, cart_id: i.id }))
  }
}

async function loadAddresses() {
  try {
    const res = await getAddresses()
    addresses.value = res.data || []
    // Check sessionStorage for address selected from Addresses page
    const stored = sessionStorage.getItem('selectedAddress')
    if (stored) {
      selectedAddress.value = JSON.parse(stored)
      sessionStorage.removeItem('selectedAddress')
    } else if (addresses.value.length && !selectedAddress.value) {
      selectedAddress.value = addresses.value[0]
    }
  } catch (e) { /* ignore */ }
}

function selectAddr(addr) { selectedAddress.value = addr; showAddrDialog.value = false }

async function submitOrder() {
  if (!selectedAddress.value) return ElMessage.warning('请选择收货地址')
  submitting.value = true
  try {
    const ids = []
    for (const item of orderItems.value) {
      const res = await createOrder({ product_id: item.product_id, num: item.num, address_id: selectedAddress.value.id })
      if (res.data?.id) ids.push(res.data.id)
    }
    // Remove cart items if cart checkout
    if (!route.query.productId) {
      for (const item of orderItems.value) {
        if (item.cart_id) await cartStore.removeItem(item.cart_id)
      }
    }
    ElMessage.success('订单创建成功')
    if (ids.length === 1) router.push(`/orders/${ids[0]}`)
    else router.push('/orders')
  } catch (e) { /* interceptor handles */ } finally { submitting.value = false }
}

onMounted(loadAddresses)
loadItems()
</script>

<style scoped>
.checkout-page { padding: 20px 0 40px; }
.page-title { font-size: 22px; margin-bottom: 20px; }
.section-card { margin-bottom: 16px; }
.card-title { font-weight: bold; }
.selected-addr { line-height: 1.8; }
.selected-addr .el-button { margin-top: 4px; padding-left: 0; }
.order-item { display: flex; align-items: center; gap: 16px; padding: 12px 0; border-bottom: 1px solid #f5f5f5; }
.order-item:last-child { border-bottom: none; }
.order-item img { width: 70px; height: 70px; object-fit: cover; border-radius: 4px; }
.item-info { flex: 1; }
.item-info .name { font-size: 14px; margin-bottom: 4px; }
.item-info .price { font-size: 13px; color: var(--text-light); }
.subtotal { font-weight: bold; color: var(--price); font-size: 16px; }
.summary-row { display: flex; justify-content: flex-end; gap: 40px; padding: 6px 0; font-size: 14px; }
.summary-row.total { border-top: 1px solid var(--border); padding-top: 12px; margin-top: 8px; font-size: 16px; font-weight: bold; }
.addr-option { padding: 12px; border: 1px solid var(--border); border-radius: 8px; margin-bottom: 8px; cursor: pointer; transition: all .2s; }
.addr-option:hover, .addr-option.selected { border-color: var(--primary); background: #fff5f5; }
</style>
