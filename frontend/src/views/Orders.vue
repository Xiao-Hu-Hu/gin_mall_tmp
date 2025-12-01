<template>
  <div class="orders">
    <div class="container">
      <h1>订单管理</h1>
      
      <!-- 创建订单表单 -->
      <div v-if="showCreateForm" class="card">
        <h2>创建订单</h2>
        <div v-if="createError" class="alert alert-error">{{ createError }}</div>
        <form @submit.prevent="createOrder">
          <div class="form-group">
            <label>选择地址</label>
            <select v-model="orderForm.address_id" required>
              <option value="">请选择地址</option>
              <option v-for="addr in addresses" :key="addr.id" :value="addr.id">
                {{ addr.name }} - {{ addr.phone }} - {{ addr.address }}
              </option>
            </select>
            <button type="button" @click="showAddressForm = true" class="btn btn-secondary" style="margin-top: 10px;">
              添加新地址
            </button>
          </div>
          <div v-if="orderForm.product_id" class="form-group">
            <label>商品ID</label>
            <input v-model="orderForm.product_id" type="number" readonly />
          </div>
          <div class="form-group">
            <label>数量</label>
            <input v-model.number="orderForm.num" type="number" min="1" required />
          </div>
          <button type="submit" class="btn btn-primary" :disabled="creating">
            {{ creating ? '创建中...' : '创建订单' }}
          </button>
          <button type="button" @click="showCreateForm = false" class="btn btn-secondary">取消</button>
        </form>
      </div>

      <!-- 添加地址表单 -->
      <div v-if="showAddressForm" class="card">
        <h2>添加地址</h2>
        <form @submit.prevent="addAddress">
          <div class="form-group">
            <label>收货人姓名</label>
            <input v-model="addressForm.name" type="text" required />
          </div>
          <div class="form-group">
            <label>电话</label>
            <input v-model="addressForm.phone" type="text" required />
          </div>
          <div class="form-group">
            <label>地址</label>
            <textarea v-model="addressForm.address" required></textarea>
          </div>
          <button type="submit" class="btn btn-primary">添加</button>
          <button type="button" @click="showAddressForm = false" class="btn btn-secondary">取消</button>
        </form>
      </div>

      <!-- 订单列表 -->
      <div v-if="loading" class="loading">加载中...</div>
      <div v-else-if="orders.length === 0" class="empty">暂无订单</div>
      <div v-else>
        <div v-for="order in orders" :key="order.id" class="card">
          <div style="display: flex; justify-content: space-between; align-items: start;">
            <div>
              <h3>订单号: {{ order.order_num }}</h3>
              <p>商品: {{ order.product_name }}</p>
              <p>数量: {{ order.num }}</p>
              <p>金额: ¥{{ order.total_money }}</p>
              <p>状态: {{ getOrderStatus(order.type) }}</p>
              <p>创建时间: {{ order.create_at }}</p>
            </div>
            <div>
              <button v-if="order.type === 1" @click="payOrder(order)" class="btn btn-success">支付</button>
              <button @click="deleteOrder(order.id)" class="btn btn-danger">删除</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import api from '../utils/api'

const route = useRoute()

const orders = ref([])
const addresses = ref([])
const loading = ref(true)
const showCreateForm = ref(false)
const showAddressForm = ref(false)
const createError = ref('')
const creating = ref(false)

const orderForm = ref({
  product_id: null,
  boss_id: null,
  num: 1,
  address_id: ''
})

const addressForm = ref({
  name: '',
  phone: '',
  address: ''
})

const getOrderStatus = (type) => {
  const statusMap = {
    1: '待支付',
    2: '已支付',
    3: '已发货',
    4: '已完成'
  }
  return statusMap[type] || '未知'
}

const fetchOrders = async () => {
  loading.value = true
  try {
    const res = await api.post('/orders')
    if (res.status === 200) {
      orders.value = res.data?.item || []
    }
  } catch (error) {
    console.error('获取订单列表失败:', error)
  } finally {
    loading.value = false
  }
}

const fetchAddresses = async () => {
  try {
    const res = await api.get('/addresses')
    if (res.status === 200) {
      // 后端直接返回地址数组
      addresses.value = res.data?.item || res.data || []
    }
  } catch (error) {
    console.error('获取地址列表失败:', error)
  }
}

const addAddress = async () => {
  try {
    const res = await api.post('/addresses', addressForm.value)
    if (res.status === 200) {
      alert('地址添加成功')
      showAddressForm.value = false
      addressForm.value = { name: '', phone: '', address: '' }
      fetchAddresses()
    }
  } catch (error) {
    alert(error.msg || error.error || '添加地址失败')
  }
}

const createOrder = async () => {
  createError.value = ''
  creating.value = true
  try {
    const res = await api.post('/orders', orderForm.value)
    if (res.status === 200) {
      alert('订单创建成功')
      showCreateForm.value = false
      orderForm.value = { product_id: null, boss_id: null, num: 1, address_id: '' }
      fetchOrders()
    } else {
      createError.value = res.msg || '创建订单失败'
    }
  } catch (error) {
    createError.value = error.msg || error.error || '创建订单失败'
  } finally {
    creating.value = false
  }
}

const payOrder = async (order) => {
  const key = prompt('请输入您的16位密钥:')
  if (!key || key.length !== 16) {
    alert('密钥必须是16位')
    return
  }
  
  try {
    const res = await api.post('/paydown', {
      order_id: order.id,
      money: parseFloat(order.money),
      order_no: order.order_num,
      product_id: order.product_id,
      boss_id: order.boss_id,
      num: order.num,
      key: key
    })
    if (res.status === 200) {
      alert('支付成功')
      fetchOrders()
    } else {
      alert(res.msg || res.error || '支付失败')
    }
  } catch (error) {
    alert(error.msg || error.error || '支付失败')
  }
}

const deleteOrder = async (id) => {
  if (!confirm('确定要删除这个订单吗？')) return
  
  try {
    const res = await api.delete(`/orders/${id}`)
    if (res.status === 200) {
      alert('删除成功')
      fetchOrders()
    }
  } catch (error) {
    alert(error.msg || error.error || '删除失败')
  }
}

onMounted(() => {
  fetchOrders()
  fetchAddresses()
  
  // 如果从商品详情页跳转过来，显示创建订单表单
  if (route.query.product_id) {
    orderForm.value.product_id = parseInt(route.query.product_id)
    orderForm.value.boss_id = parseInt(route.query.boss_id)
    orderForm.value.num = parseInt(route.query.num) || 1
    showCreateForm.value = true
  }
})
</script>

