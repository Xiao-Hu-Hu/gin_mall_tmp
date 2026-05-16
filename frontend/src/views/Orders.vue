<template>
  <div class="mall-container orders-page">
    <h2 class="page-title">我的订单</h2>
    <el-tabs v-model="statusFilter" @tab-change="loadOrders">
      <el-tab-pane label="全部" name="all" />
      <el-tab-pane label="待支付" name="1" />
      <el-tab-pane label="已支付" name="2" />
    </el-tabs>
    <div v-if="orders.length" class="order-list">
      <el-card v-for="order in orders" :key="order.id" shadow="never" class="order-card">
        <div class="order-header">
          <span>订单号: {{ order.order_num }}</span>
          <span>{{ formatTime(order.create_at) }}</span>
          <el-tag :type="formatOrderStatusType(order.type)" size="small">{{ formatOrderStatus(order.type) }}</el-tag>
        </div>
        <div class="order-body" @click="$router.push(`/orders/${order.id}?product_id=${order.product_id}&boss_id=${order.boss_id}&address_id=${order.address_id}`)">
          <img :src="order.img_path" />
          <div class="order-info">
            <div class="name">{{ order.product_name }}</div>
            <div class="meta">x{{ order.num }} &yen;{{ formatPrice(order.discount_price || order.money) }}</div>
          </div>
          <div class="order-total">&yen;{{ formatPrice(order.total_money) }}</div>
        </div>
        <div class="order-footer" v-if="order.type === 1">
          <el-button type="danger" size="small" @click="$router.push(`/orders/${order.id}?product_id=${order.product_id}&boss_id=${order.boss_id}&address_id=${order.address_id}`)">去支付</el-button>
          <el-button text type="info" size="small" @click="handleDelete(order.id)">取消订单</el-button>
        </div>
      </el-card>
    </div>
    <el-empty v-else description="暂无订单" />
    <div class="pagination" v-if="total > pageSize">
      <el-pagination background layout="prev, pager, next" :total="total" :page-size="pageSize" :current-page="pageNum" @current-change="p => { pageNum = p; loadOrders() }" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getOrders, deleteOrder } from '../api/order'
import { formatPrice, formatTime, formatOrderStatus, formatOrderStatusType } from '../utils/format'
import { ElMessageBox } from 'element-plus'

const orders = ref([])
const statusFilter = ref('all')
const pageNum = ref(1)
const pageSize = 10
const total = ref(0)

async function loadOrders() {
  const params = { page_num: pageNum.value, page_size: pageSize }
  if (statusFilter.value !== 'all') params.type = Number(statusFilter.value)
  try { const res = await getOrders(params); orders.value = res.data?.item || []; total.value = res.data?.total || 0 } catch (e) { orders.value = [] }
}

function handleDelete(id) {
  ElMessageBox.confirm('确定取消此订单？', '提示', { type: 'warning' }).then(async () => {
    await deleteOrder(id)
    loadOrders()
  }).catch(() => {})
}

onMounted(loadOrders)
</script>

<style scoped>
.orders-page { padding: 20px 0 40px; }
.page-title { font-size: 22px; margin-bottom: 20px; }
.order-card { margin-bottom: 12px; }
.order-header { display: flex; align-items: center; gap: 16px; font-size: 13px; color: var(--text-light); margin-bottom: 12px; }
.order-body { display: flex; align-items: center; gap: 16px; cursor: pointer; }
.order-body img { width: 80px; height: 80px; object-fit: cover; border-radius: 4px; }
.order-info { flex: 1; }
.order-info .name { font-size: 14px; margin-bottom: 4px; }
.order-info .meta { font-size: 13px; color: var(--text-light); }
.order-total { font-size: 18px; font-weight: bold; color: var(--price); }
.order-footer { display: flex; justify-content: flex-end; gap: 8px; margin-top: 12px; padding-top: 12px; border-top: 1px solid #f5f5f5; }
.pagination { display: flex; justify-content: center; margin-top: 20px; }
</style>
