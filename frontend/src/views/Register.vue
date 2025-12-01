<template>
  <div class="register">
    <div class="container">
      <div class="card" style="max-width: 400px; margin: 50px auto;">
        <h2 style="text-align: center; margin-bottom: 20px;">用户注册</h2>
        <div v-if="error" class="alert alert-error">{{ error }}</div>
        <div v-if="success" class="alert alert-success">{{ success }}</div>
        <form @submit.prevent="handleRegister">
          <div class="form-group">
            <label>用户名</label>
            <input v-model="form.user_name" type="text" required />
          </div>
          <div class="form-group">
            <label>昵称</label>
            <input v-model="form.nick_name" type="text" required />
          </div>
          <div class="form-group">
            <label>邮箱</label>
            <input v-model="form.email" type="email" required placeholder="请输入常用邮箱" />
          </div>
          <div class="form-group">
            <label>密码</label>
            <input v-model="form.password" type="password" required />
          </div>
          <div class="form-group">
            <label>密钥（16位，用于加密金额）</label>
            <input v-model="form.key" type="text" maxlength="16" required />
            <small style="color: #666; font-size: 12px;">请输入16位密钥，用于加密您的账户余额</small>
          </div>
          <button type="submit" class="btn btn-primary" style="width: 100%;" :disabled="loading">
            {{ loading ? '注册中...' : '注册' }}
          </button>
        </form>
        <p style="text-align: center; margin-top: 15px;">
          已有账号？<router-link to="/login">立即登录</router-link>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '../utils/api'

const router = useRouter()

const form = ref({
  user_name: '',
  nick_name: '',
  email: '',
  password: '',
  key: ''
})

const error = ref('')
const success = ref('')
const loading = ref(false)

const handleRegister = async () => {
  error.value = ''
  success.value = ''
  
  if (form.value.key.length !== 16) {
    error.value = '密钥必须是16位'
    return
  }
  
  loading.value = true
  try {
    const res = await api.post('/user/register', form.value)
    if (res.status === 200) {
      success.value = '注册成功，请登录'
      setTimeout(() => {
        router.push('/login')
      }, 1500)
    } else {
      error.value = res.msg || '注册失败'
    }
  } catch (err) {
    error.value = err.msg || err.error || '注册失败'
  } finally {
    loading.value = false
  }
}
</script>

