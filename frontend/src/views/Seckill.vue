<template>
  <div class="mall-container seckill-page">
    <h2 class="page-title">秒杀活动</h2>
    <el-card v-if="activity" shadow="never">
      <div class="activity-header">
        <div class="status-badge" :class="`status-${activity.Status}`">
          {{ activity.Status === 0 ? '未开始' : activity.Status === 1 ? '进行中' : '已结束' }}
        </div>
        <div class="countdown" v-if="countdownText">
          <span>{{ activity.Status === 1 ? '距结束' : '距开始' }}</span>
          <span class="time">{{ countdownText }}</span>
        </div>
      </div>
      <div class="seckill-info">
        <div class="seckill-price-box">
          <div class="seckill-price">&yen;{{ formatPrice(activity.SeckillPrice) }}</div>
          <div class="original-price">原价 &yen;{{ productPrice }}</div>
        </div>
        <div class="stock-bar">
          <el-progress :percentage="stockPercent" :color="stockPercent > 80 ? '#ff4d4f' : '#409eff'" />
          <span>剩余 {{ remainingStock }} / {{ activity.TotalStock }} 件</span>
        </div>
        <el-button type="danger" size="large" :disabled="activity.Status !== 1 || remainingStock <= 0 || seckillLoading" :loading="seckillLoading" @click="doSeckill" class="seckill-btn">
          {{ activity.Status === 0 ? '未开始' : activity.Status === 2 ? '已结束' : remainingStock <= 0 ? '已售罄' : '立即抢购' }}
        </el-button>
        <div v-if="seckillResult" class="seckill-result">
          <el-alert :title="seckillResult" :type="seckillResult.includes('成功') ? 'success' : 'info'" show-icon :closable="false" />
        </div>
      </div>
    </el-card>
    <div v-else-if="loading" style="text-align:center;padding:40px"><el-icon class="is-loading" :size="32"><Loading /></el-icon></div>
    <!-- Activity ID selector for demo -->
    <el-card shadow="never" style="margin-top:16px">
      <el-form inline>
        <el-form-item label="活动ID"><el-input-number v-model="activityId" :min="1" /></el-form-item>
        <el-form-item><el-button type="primary" @click="loadActivity">查看活动</el-button></el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { getSeckillActivity, createSeckillOrder, getSeckillOrderStatus } from '../api/seckill'
import { getProduct } from '../api/product'
import { formatPrice } from '../utils/format'
import { ElMessage } from 'element-plus'

const router = useRouter()
const auth = useAuthStore()
const activity = ref(null)
const activityId = ref(1)
const loading = ref(false)
const remainingStock = ref(0)
const productPrice = ref('0')
const seckillLoading = ref(false)
const seckillResult = ref('')
const countdownText = ref('')
let timer = null

const stockPercent = computed(() => {
  if (!activity.value) return 0
  return Math.round(((activity.value.TotalStock - remainingStock.value) / activity.value.TotalStock) * 100)
})

function startCountdown() {
  if (timer) clearInterval(timer)
  timer = setInterval(() => {
    if (!activity.value) return
    const now = Date.now()
    const end = new Date(activity.value.EndAt).getTime()
    const start = new Date(activity.value.StartAt).getTime()
    let target, label
    if (activity.value.Status === 1) { target = end } else if (activity.value.Status === 0) { target = start } else { countdownText.value = ''; return }
    const diff = Math.max(0, target - now)
    const h = Math.floor(diff / 3600000)
    const m = Math.floor((diff % 3600000) / 60000)
    const s = Math.floor((diff % 60000) / 1000)
    countdownText.value = `${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
  }, 1000)
}

async function loadActivity() {
  loading.value = true
  try {
    const res = await getSeckillActivity(activityId.value)
    activity.value = res.data?.activity || res.data?.Activity
    remainingStock.value = res.data?.remaining_stock || 0
    if (activity.value?.ProductID) {
      try { const pRes = await getProduct(activity.value.ProductID); productPrice.value = pRes.data?.price || '0' } catch (e) { /* ignore */ }
    }
    startCountdown()
  } catch (e) { activity.value = null } finally { loading.value = false }
}

async function doSeckill() {
  if (!auth.token) return router.push('/login')
  seckillLoading.value = true
  seckillResult.value = ''
  try {
    const res = await createSeckillOrder({ activity_id: activity.value.ID })
    seckillResult.value = res.msg || '请求已提交'
    // Poll for order status
    let attempts = 0
    const poll = setInterval(async () => {
      attempts++
      try {
        const statusRes = await getSeckillOrderStatus(activity.value.ID)
        if (statusRes.data) {
          clearInterval(poll)
          seckillResult.value = '抢购成功！订单已创建'
          loadActivity()
        }
      } catch (e) { /* still waiting */ }
      if (attempts >= 15) { clearInterval(poll); seckillResult.value = seckillResult.value || '请稍后查看订单' }
    }, 2000)
  } catch (e) { /* interceptor handles */ } finally { seckillLoading.value = false }
}

onMounted(loadActivity)
onUnmounted(() => { if (timer) clearInterval(timer) })
</script>

<style scoped>
.seckill-page { padding: 20px 0 40px; }
.page-title { font-size: 22px; margin-bottom: 20px; }
.activity-header { display: flex; align-items: center; gap: 16px; margin-bottom: 24px; }
.status-badge { padding: 4px 16px; border-radius: 16px; font-size: 14px; font-weight: bold; }
.status-badge.status-0 { background: #e8e8e8; color: #999; }
.status-badge.status-1 { background: #ff4d4f; color: #fff; animation: pulse 2s infinite; }
.status-badge.status-2 { background: #e8e8e8; color: #999; }
@keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: .7; } }
.countdown { font-size: 14px; }
.countdown .time { font-size: 24px; font-weight: bold; color: var(--primary); margin-left: 8px; font-family: monospace; }
.seckill-info { text-align: center; padding: 20px; }
.seckill-price-box { margin-bottom: 20px; }
.seckill-price { font-size: 40px; font-weight: bold; color: var(--price); }
.original-price { font-size: 14px; color: var(--text-light); text-decoration: line-through; margin-top: 4px; }
.stock-bar { max-width: 400px; margin: 0 auto 20px; font-size: 13px; color: var(--text-light); }
.seckill-btn { width: 200px; height: 48px; font-size: 18px; }
.seckill-result { margin-top: 16px; max-width: 400px; margin-left: auto; margin-right: auto; }
</style>
