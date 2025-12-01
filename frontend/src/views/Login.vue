<template>
  <div class="login">
    <div class="container">
      <div class="card" style="max-width: 400px; margin: 50px auto;">
        <h2 style="text-align: center; margin-bottom: 20px;">用户登录</h2>
        <div v-if="error" class="alert alert-error">{{ error }}</div>
        <form @submit.prevent="handleLogin">
          <div class="form-group">
            <label>用户名</label>
            <input v-model="form.user_name" type="text" required />
          </div>
          <div class="form-group">
            <label>密码</label>
            <input v-model="form.password" type="password" required />
          </div>
          <button type="submit" class="btn btn-primary" style="width: 100%;" :disabled="loading">
            {{ loading ? '登录中...' : '登录' }}
          </button>
        </form>
        <p style="text-align: center; margin-top: 15px;">
          还没有账号？<router-link to="/register">立即注册</router-link>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import api from '../utils/api'

const router = useRouter()
const userStore = useUserStore()

const form = ref({
  user_name: '',
  password: ''
})

const error = ref('')
const loading = ref(false)

const handleLogin = async () => {
  error.value = ''
  loading.value = true
  try {
    const res = await api.post('/user/login', form.value)
    if (res.status === 200 && res.data) {
      const { user, token } = res.data
      userStore.login(user, token)
      router.push('/')
    } else {
      error.value = res.msg || '登录失败'
    }
  } catch (err) {
    error.value = err.msg || err.error || '登录失败，请检查用户名和密码'
  } finally {
    loading.value = false
  }
}
</script>

