<template>
  <div class="sell">
    <div class="container">
      <h1>上架商品</h1>
      <div class="card">
        <div v-if="error" class="alert alert-error">{{ error }}</div>
        <div v-if="success" class="alert alert-success">{{ success }}</div>
        <form @submit.prevent="handleSubmit">
          <div class="form-group">
            <label>商品名称</label>
            <input v-model="form.name" type="text" required />
          </div>
          <div class="form-group">
            <label>商品标题</label>
            <input v-model="form.title" type="text" required />
          </div>
          <div class="form-group">
            <label>商品分类</label>
            <select v-model.number="form.category_id" required>
              <option value="">请选择分类</option>
              <option v-for="cat in categories" :key="cat.id" :value="cat.id">
                {{ cat.name }}
              </option>
            </select>
          </div>
          <div class="form-group">
            <label>商品描述</label>
            <textarea v-model="form.info" required></textarea>
          </div>
          <div class="form-group">
            <label>价格</label>
            <input v-model="form.price" type="number" step="0.01" required />
          </div>
          <div class="form-group">
            <label>折扣价（可选）</label>
            <input v-model="form.discount_price" type="number" step="0.01" />
          </div>
          <div class="form-group">
            <label>库存数量</label>
            <input v-model.number="form.num" type="number" min="1" required />
          </div>
          <div class="form-group">
            <label>商品图片（可多选）</label>
            <input type="file" multiple accept="image/*" @change="handleFileChange" />
            <div v-if="files.length > 0" style="margin-top: 10px;">
              <p>已选择 {{ files.length }} 张图片</p>
            </div>
          </div>
          <button type="submit" class="btn btn-primary" :disabled="loading">
            {{ loading ? '上传中...' : '上架商品' }}
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../utils/api'

const form = ref({
  name: '',
  title: '',
  category_id: '',
  info: '',
  price: '',
  discount_price: '',
  num: 1
})

const categories = ref([])
const files = ref([])
const error = ref('')
const success = ref('')
const loading = ref(false)

const fetchCategories = async () => {
  try {
    const res = await api.get('/categories')
    if (res.status === 200) {
      categories.value = res.data?.item || []
    }
  } catch (error) {
    console.error('获取分类失败:', error)
  }
}

const handleFileChange = (e) => {
  files.value = Array.from(e.target.files)
}

const handleSubmit = async () => {
  if (files.value.length === 0) {
    error.value = '请至少上传一张商品图片'
    return
  }
  
  error.value = ''
  success.value = ''
  loading.value = true
  
  const formData = new FormData()
  formData.append('name', form.value.name)
  formData.append('title', form.value.title)
  formData.append('category_id', form.value.category_id)
  formData.append('info', form.value.info)
  formData.append('price', form.value.price)
  if (form.value.discount_price) {
    formData.append('discount_price', form.value.discount_price)
  }
  formData.append('num', form.value.num)
  
  files.value.forEach(file => {
    formData.append('file', file)
  })
  
  try {
    const res = await api.post('/product', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    if (res.status === 200) {
      success.value = '商品上架成功'
      form.value = {
        name: '',
        title: '',
        category_id: '',
        info: '',
        price: '',
        discount_price: '',
        num: 1
      }
      files.value = []
    } else {
      error.value = res.msg || '上架失败'
    }
  } catch (err) {
    error.value = err.msg || err.error || '上架失败'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchCategories()
})
</script>

