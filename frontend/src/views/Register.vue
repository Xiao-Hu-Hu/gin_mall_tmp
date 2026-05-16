<template>
  <div class="register-page">
    <el-card class="register-card" shadow="always">
      <h2 class="title">注册新账号</h2>
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top" @submit.prevent="doRegister">
        <el-form-item label="用户名" prop="user_name">
          <el-input v-model="form.user_name" prefix-icon="User" placeholder="3-20个字符" size="large" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="form.password" type="password" prefix-icon="Lock" placeholder="至少6位" show-password size="large" />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input v-model="form.confirmPassword" type="password" prefix-icon="Lock" placeholder="再次输入密码" show-password size="large" />
        </el-form-item>
        <el-form-item label="昵称" prop="nick_name">
          <el-input v-model="form.nick_name" placeholder="您的昵称" size="large" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email" prefix-icon="Message" placeholder="your@email.com" size="large" />
        </el-form-item>
        <el-form-item label="加密密钥" prop="key">
          <el-input v-model="form.key" placeholder="请设置16位密钥" maxlength="16" show-word-limit size="large" />
          <div class="form-tip">用于余额加密保护，请牢记此密钥</div>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" native-type="submit" :loading="loading" size="large" style="width:100%">注册</el-button>
        </el-form-item>
      </el-form>
      <div class="footer-link">已有账号？<router-link to="/login">立即登录</router-link></div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { userRegister } from '../api/auth'
import { ElMessage } from 'element-plus'

const router = useRouter()
const formRef = ref()
const loading = ref(false)
const form = reactive({ user_name: '', password: '', confirmPassword: '', nick_name: '', email: '', key: '' })

const validateConfirm = (rule, value, callback) => {
  if (value !== form.password) callback(new Error('两次密码不一致'))
  else callback()
}

const rules = {
  user_name: [{ required: true, message: '请输入用户名', trigger: 'blur' }, { min: 3, max: 20, message: '3-20个字符', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }, { min: 6, message: '至少6位', trigger: 'blur' }],
  confirmPassword: [{ required: true, message: '请确认密码', trigger: 'blur' }, { validator: validateConfirm, trigger: 'blur' }],
  nick_name: [{ required: true, message: '请输入昵称', trigger: 'blur' }],
  email: [{ required: true, message: '请输入邮箱', trigger: 'blur' }, { type: 'email', message: '邮箱格式不正确', trigger: 'blur' }],
  key: [{ required: true, message: '请设置密钥', trigger: 'blur' }, { len: 16, message: '密钥必须为16位', trigger: 'blur' }]
}

async function doRegister() {
  await formRef.value.validate()
  loading.value = true
  try {
    await userRegister({ user_name: form.user_name, password: form.password, nick_name: form.nick_name, email: form.email, key: form.key })
    ElMessage.success('注册成功，请登录')
    router.push('/login')
  } catch (e) { /* interceptor handles error */ } finally { loading.value = false }
}
</script>

<style scoped>
.register-page { display: flex; justify-content: center; align-items: center; min-height: calc(100vh - 200px); background: linear-gradient(135deg, #fff5f5 0%, #fff0e6 100%); padding: 20px 0; }
.register-card { width: 420px; padding: 20px; }
.title { text-align: center; margin-bottom: 24px; color: var(--primary); }
.footer-link { text-align: center; font-size: 14px; color: var(--text-light); }
.form-tip { font-size: 12px; color: var(--text-light); margin-top: 4px; }
</style>
