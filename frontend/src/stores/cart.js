import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getCarts, addCart, updateCart, deleteCart } from '../api/cart'

export const useCartStore = defineStore('cart', () => {
  const items = ref([])
  const selectedIds = ref(new Set())

  const selectedItems = computed(() => items.value.filter(i => selectedIds.value.has(i.id)))
  const totalCount = computed(() => items.value.reduce((s, i) => s + i.num, 0))
  const selectedCount = computed(() => selectedItems.value.reduce((s, i) => s + i.num, 0))
  const totalPrice = computed(() => selectedItems.value.reduce((s, i) => s + parseFloat(i.discount_price || i.price) * i.num, 0))
  const isAllSelected = computed(() => items.value.length > 0 && items.value.every(i => selectedIds.value.has(i.id)))

  async function fetchCart() {
    try {
      const res = await getCarts()
      items.value = res.data?.item || []
      items.value.forEach(i => { if (i.check) selectedIds.value.add(i.id) })
    } catch (e) { items.value = [] }
  }

  async function addItem(productId, bossId, num = 1) {
    await addCart({ product_id: productId, boss_id: bossId, num })
    await fetchCart()
  }

  async function updateItem(id, num) {
    await updateCart(id, { num })
    await fetchCart()
  }

  async function removeItem(id) {
    await deleteCart(id)
    selectedIds.value.delete(id)
    await fetchCart()
  }

  function toggleSelect(id) {
    if (selectedIds.value.has(id)) selectedIds.value.delete(id)
    else selectedIds.value.add(id)
  }

  function toggleSelectAll() {
    if (isAllSelected.value) selectedIds.value.clear()
    else items.value.forEach(i => selectedIds.value.add(i.id))
  }

  function clearSelected() { selectedIds.value.clear() }

  return { items, selectedIds, selectedItems, totalCount, selectedCount, totalPrice, isAllSelected, fetchCart, addItem, updateItem, removeItem, toggleSelect, toggleSelectAll, clearSelected }
})
