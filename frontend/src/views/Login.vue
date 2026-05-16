<template>
  <div class="login-page">
    <el-card class="login-card" shadow="always">
      <h2 class="title">欢迎登录 Gin Mall</h2>
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top" @submit.prevent="doLogin">
        <el-form-item label="用户名" prop="user_name">
          <el-input v-model="form.user_name" prefix-icon="User" placeholder="请输入用户名" size="large" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="form.password" type="password" prefix-icon="Lock" placeholder="请输入密码" show-password size="large" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" native-type="submit" :loading="loading" size="large" style="width:100%">登录</el-button>
        </el-form-item>
      </el-form>
      <div class="footer-link">没有账号？<router-link to="/register">立即注册</router-link></div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useCartStore } from '../stores/cart'
import { userLogin } from '../api/auth'
import { ElMessage } from 'element-plus'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()
const cartStore = useCartStore()
const formRef = ref()
const loading = ref(false)
const form = reactive({ user_name: '', password: '' })
const rules = {
  user_name: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

async function doLogin() {
  await formRef.value.validate()
  loading.value = true
  try {
    const res = await userLogin(form)
    if (res.data?.token) {
      auth.setAuth(res.data.token, res.data.user)
      cartStore.fetchCart()
      ElMessage.success('登录成功')
      router.push(route.query.redirect || '/')
    }
  } catch (e) { /* interceptor handles error */ } finally { loading.value = false }
}
</script>

<style scoped>
.login-page { display: flex; justify-content: center; align-items: center; min-height: calc(100vh - 200px); background: linear-gradient(135deg, #fff5f5 0%, #fff0e6 100%); }
.login-card { width: 420px; padding: 20px; }
.title { text-align: center; margin-bottom: 24px; color: var(--primary); }
.footer-link { text-align: center; font-size: 14px; color: var(--text-light); }
</style>
