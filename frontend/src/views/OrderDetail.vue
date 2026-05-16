<template>
  <div class="mall-container order-detail-page">
    <el-page-header @back="$router.push('/orders')" content="订单详情" style="margin-bottom:20px" />
    <el-card v-if="order" shadow="never">
      <div class="status-banner" :class="`status-${order.type}`">
        <el-tag :type="formatOrderStatusType(order.type)" size="large">{{ formatOrderStatus(order.type) }}</el-tag>
        <span v-if="order.type === 1">请尽快完成支付</span>
        <span v-else-if="order.type === 2">感谢您的购买</span>
      </div>
      <el-descriptions :column="2" border style="margin-top:16px">
        <el-descriptions-item label="订单号">{{ order.order_num }}</el-descriptions-item>
        <el-descriptions-item label="下单时间">{{ formatTime(order.create_at) }}</el-descriptions-item>
        <el-descriptions-item label="商品">{{ order.product_name }}</el-descriptions-item>
        <el-descriptions-item label="数量">{{ order.num }}</el-descriptions-item>
        <el-descriptions-item label="商品图片"><img :src="order.img_path" style="width:80px;height:80px;object-fit:cover;border-radius:4px" /></el-descriptions-item>
        <el-descriptions-item label="订单金额"><span class="price-current" style="font-size:20px">&yen;{{ formatPrice(order.total_money) }}</span></el-descriptions-item>
        <el-descriptions-item label="收货人">{{ order.address_name }} {{ order.address_phone }}</el-descriptions-item>
        <el-descriptions-item label="收货地址">{{ order.address }}</el-descriptions-item>
        <el-descriptions-item label="商家">{{ order.boss_nickname }}</el-descriptions-item>
      </el-descriptions>
      <!-- Payment -->
      <div v-if="order.type === 1" class="pay-section">
        <el-divider>余额支付</el-divider>
        <el-form inline>
          <el-form-item label="加密密钥">
            <el-input v-model="payKey" type="password" placeholder="请输入16位密钥" maxlength="16" show-password style="width:300px" />
          </el-form-item>
          <el-form-item>
            <el-button type="danger" size="large" :loading="paying" @click="doPay">确认支付 &yen;{{ formatPrice(order.total_money) }}</el-button>
          </el-form-item>
        </el-form>
      </div>
    </el-card>
    <el-skeleton v-else :rows="8" animated />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getOrder } from '../api/order'
import { payOrder } from '../api/payment'
import { formatPrice, formatTime, formatOrderStatus, formatOrderStatusType } from '../utils/format'
import { ElMessage } from 'element-plus'

const route = useRoute()
const order = ref(null)
const payKey = ref('')
const paying = ref(false)

async function loadOrder() {
  const params = {}
  if (route.query.product_id) params.product_id = Number(route.query.product_id)
  if (route.query.boss_id) params.boss_id = Number(route.query.boss_id)
  if (route.query.address_id) params.address_id = Number(route.query.address_id)
  try {
    const res = await getOrder(route.params.id, params)
    order.value = res.data
  } catch (e) { /* ignore */ }
}

async function doPay() {
  if (!payKey.value || payKey.value.length !== 16) return ElMessage.warning('请输入16位密钥')
  paying.value = true
  try {
    await payOrder({
      order_id: order.value.id,
      money: order.value.total_money,
      order_no: order.value.order_num,
      product_id: order.value.product_id,
      boss_id: order.value.boss_id,
      boss_name: order.value.boss_nickname,
      num: order.value.num,
      key: payKey.value,
      pay_time: new Date().toISOString(),
      sign: ''
    })
    ElMessage.success('支付请求已提交，正在处理...')
    // Poll for payment status
    let attempts = 0
    const poll = setInterval(async () => {
      attempts++
      await loadOrder()
      if (order.value?.type === 2) {
        clearInterval(poll)
        ElMessage.success('支付成功！')
      } else if (attempts >= 30) {
        clearInterval(poll)
        ElMessage.info('支付处理中，请稍后刷新查看')
      }
    }, 2000)
  } catch (e) { /* interceptor handles */ } finally { paying.value = false }
}

onMounted(loadOrder)
</script>

<style scoped>
.order-detail-page { padding: 20px 0 40px; }
.status-banner { padding: 16px; border-radius: 8px; display: flex; align-items: center; gap: 12px; }
.status-banner.status-1 { background: #fdf6ec; }
.status-banner.status-2 { background: #f0f9eb; }
.pay-section { margin-top: 20px; text-align: center; }
</style>
