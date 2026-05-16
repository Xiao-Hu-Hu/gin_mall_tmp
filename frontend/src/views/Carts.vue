<template>
  <div class="mall-container cart-page">
    <h2 class="page-title">我的购物车</h2>
    <div v-if="cartStore.items.length" class="cart-layout">
      <div class="cart-list">
        <div class="cart-header">
          <el-checkbox :model-value="cartStore.isAllSelected" @change="cartStore.toggleSelectAll()">全选</el-checkbox>
          <span>商品信息</span><span>单价</span><span>数量</span><span>小计</span><span>操作</span>
        </div>
        <div v-for="item in cartStore.items" :key="item.id" class="cart-item">
          <el-checkbox :model-value="cartStore.selectedIds.has(item.id)" @change="cartStore.toggleSelect(item.id)" />
          <div class="item-info">
            <router-link :to="`/products/${item.product_id}`"><img :src="item.img_path" /></router-link>
            <div>
              <router-link :to="`/products/${item.product_id}`" class="item-name">{{ item.name }}</router-link>
              <div class="item-boss">{{ item.boss_name }}</div>
            </div>
          </div>
          <div class="item-price">&yen;{{ formatPrice(item.discount_price || item.price) }}</div>
          <div class="item-qty">
            <el-input-number :model-value="item.num" :min="1" :max="item.maxNum || 99" @change="v => cartStore.updateItem(item.id, v)" size="small" />
          </div>
          <div class="item-subtotal">&yen;{{ formatPrice((parseFloat(item.discount_price || item.price) * item.num).toFixed(2)) }}</div>
          <el-button text type="danger" @click="handleDelete(item.id)">删除</el-button>
        </div>
      </div>
      <div class="cart-summary">
        <div class="summary-row"><span>已选商品</span><span class="val">{{ cartStore.selectedCount }} 件</span></div>
        <div class="summary-row total"><span>合计</span><span class="val price-current">&yen;{{ formatPrice(cartStore.totalPrice.toFixed(2)) }}</span></div>
        <el-button type="danger" size="large" :disabled="!cartStore.selectedCount" @click="$router.push('/checkout')" style="width:100%;margin-top:16px">去结算</el-button>
      </div>
    </div>
    <el-empty v-else description="购物车是空的">
      <el-button type="primary" @click="$router.push('/')">去逛逛</el-button>
    </el-empty>
  </div>
</template>

<script setup>
import { useCartStore } from '../stores/cart'
import { formatPrice } from '../utils/format'
import { ElMessageBox } from 'element-plus'

const cartStore = useCartStore()
cartStore.fetchCart()

function handleDelete(id) {
  ElMessageBox.confirm('确定删除此商品？', '提示', { type: 'warning' }).then(() => cartStore.removeItem(id)).catch(() => {})
}
</script>

<style scoped>
.cart-page { padding: 20px 0 40px; }
.page-title { font-size: 22px; margin-bottom: 20px; }
.cart-layout { display: flex; gap: 20px; }
.cart-list { flex: 1; background: #fff; border-radius: 8px; padding: 16px; }
.cart-header { display: flex; align-items: center; gap: 16px; padding-bottom: 12px; border-bottom: 1px solid var(--border); font-size: 13px; color: var(--text-light); }
.cart-header span:nth-child(2) { flex: 2; }
.cart-header span:nth-child(3), .cart-header span:nth-child(4), .cart-header span:nth-child(5), .cart-header span:nth-child(6) { width: 100px; text-align: center; }
.cart-item { display: flex; align-items: center; gap: 16px; padding: 16px 0; border-bottom: 1px solid #f5f5f5; }
.item-info { flex: 2; display: flex; gap: 12px; align-items: center; }
.item-info img { width: 80px; height: 80px; object-fit: cover; border-radius: 4px; }
.item-name { font-size: 14px; color: var(--text); display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; }
.item-boss { font-size: 12px; color: var(--text-light); margin-top: 4px; }
.item-price { width: 100px; text-align: center; font-size: 14px; }
.item-qty { width: 100px; text-align: center; }
.item-subtotal { width: 100px; text-align: center; font-weight: bold; color: var(--price); }
.cart-summary { width: 280px; background: #fff; border-radius: 8px; padding: 20px; height: fit-content; position: sticky; top: 120px; }
.summary-row { display: flex; justify-content: space-between; padding: 8px 0; font-size: 14px; }
.summary-row.total { border-top: 1px solid var(--border); padding-top: 12px; margin-top: 8px; font-size: 16px; }
.summary-row .val { font-weight: bold; }
</style>
