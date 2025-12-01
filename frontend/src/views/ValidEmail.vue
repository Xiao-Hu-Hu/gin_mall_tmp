<template>
  <div class="valid-email">
    <div class="container">
      <div class="card" style="max-width: 500px; margin: 50px auto;">
        <h2 style="text-align: center;">邮箱验证</h2>
        <div v-if="loading" class="loading">验证中...</div>
        <div v-else-if="success" class="alert alert-success">
          {{ success }}
        </div>
        <div v-else-if="error" class="alert alert-error">
          {{ error }}
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
const loading = ref(true)
const success = ref('')
const error = ref('')

const validateEmail = async () => {
  const token = route.params.token
  if (!token) {
    error.value = '缺少验证token'
    loading.value = false
    return
  }
  
  try {
    const res = await api.post('/user/valid-email', {}, {
      headers: {
        'Authorization': token
      }
    })
    if (res.status === 200) {
      success.value = '邮箱验证成功！'
    } else {
      error.value = res.msg || '验证失败'
    }
  } catch (err) {
    error.value = err.msg || err.error || '验证失败'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  validateEmail()
})
</script>

