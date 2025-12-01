<template>
  <div class="profile">
    <div class="container">
      <h1>个人中心</h1>
      <div class="card">
        <h2>个人信息</h2>
        <div v-if="error" class="alert alert-error">{{ error }}</div>
        <div v-if="success" class="alert alert-success">{{ success }}</div>
        
        <div class="form-group">
          <label>用户名</label>
          <input :value="user?.user_name" type="text" disabled />
        </div>
        
        <div class="form-group">
          <label>昵称</label>
          <input v-model="nickName" type="text" />
          <button @click="updateNickName" class="btn btn-primary" style="margin-top: 10px;">更新昵称</button>
        </div>
        
        <div class="form-group">
          <label>邮箱</label>
          <input v-model="email" type="email" placeholder="请输入邮箱" />
          <button @click="sendEmail" class="btn btn-primary" style="margin-top: 10px;">发送验证邮件</button>
          <div v-if="user?.email" style="margin-top: 10px; color: #666;">
            当前邮箱: {{ user.email }}
          </div>
        </div>
        
        <div class="form-group">
          <label>头像</label>
          <input type="file" accept="image/*" @change="handleAvatarChange" />
          <div v-if="user?.avatar" style="margin-top: 10px;">
            <img :src="getAvatarUrl(user.avatar)" style="width: 100px; height: 100px; object-fit: cover; border-radius: 50%;" />
          </div>
        </div>
        
        <div class="form-group">
          <label>账户余额</label>
          <input v-model="key" type="text" placeholder="请输入16位密钥查看余额" maxlength="16" />
          <button @click="showMoney" class="btn btn-primary" style="margin-top: 10px;">查看余额</button>
          <div v-if="money" style="margin-top: 10px; font-size: 20px; font-weight: bold; color: #28a745;">
            余额: ¥{{ money }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '../stores/user'
import api from '../utils/api'

const userStore = useUserStore()
const user = ref(null)
const nickName = ref('')
const email = ref('')
const key = ref('')
const money = ref('')
const error = ref('')
const success = ref('')

const fetchUser = async () => {
  // 这里需要从store获取用户信息，或者调用API获取
  // 由于后端没有获取当前用户信息的接口，我们假设用户信息在登录时已经保存
  user.value = userStore.user
  if (user.value) {
    nickName.value = user.value.nick_name || ''
    email.value = user.value.email || ''
  }
}

const updateNickName = async () => {
  if (!nickName.value.trim()) {
    error.value = '昵称不能为空'
    return
  }
  
  try {
    const res = await api.put('/user', { nick_name: nickName.value })
    if (res.status === 200) {
      success.value = '昵称更新成功'
      user.value = res.data
      userStore.user = res.data
    } else {
      error.value = res.msg || '更新失败'
    }
  } catch (err) {
    error.value = err.msg || err.error || '更新失败'
  }
}

const sendEmail = async () => {
  if (!email.value.trim()) {
    error.value = '请输入邮箱'
    return
  }
  
  const password = prompt('请输入密码（用于绑定邮箱）:')
  if (!password) return
  
  try {
    const res = await api.post('/user/sending-email', {
      email: email.value,
      password: password,
      operation_type: 1 // 1: 绑定邮箱
    })
    if (res.status === 200) {
      success.value = '验证邮件已发送，请查收'
    } else {
      error.value = res.msg || '发送失败'
    }
  } catch (err) {
    error.value = err.msg || err.error || '发送失败'
  }
}

const handleAvatarChange = async (e) => {
  const file = e.target.files[0]
  if (!file) return
  
  const formData = new FormData()
  formData.append('file', file)
  
  try {
    const res = await api.post('/avatar', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    if (res.status === 200) {
      success.value = '头像上传成功'
      user.value = res.data
      userStore.user = res.data
    } else {
      error.value = res.msg || '上传失败'
    }
  } catch (err) {
    error.value = err.msg || err.error || '上传失败'
  }
}

const showMoney = async () => {
  if (!key.value || key.value.length !== 16) {
    error.value = '密钥必须是16位'
    return
  }
  
  try {
    const res = await api.post('/money', { key: key.value })
    if (res.status === 200) {
      money.value = res.data
    } else {
      error.value = res.msg || '查询失败'
    }
  } catch (err) {
    error.value = err.msg || err.error || '查询失败'
  }
}

const getAvatarUrl = (avatar) => {
  if (!avatar) return ''
  // 如果已经是完整URL，直接返回
  if (avatar.startsWith('http')) {
    return avatar
  }
  // 否则拼接完整URL
  return `http://localhost:3000/static/imgs/avatar/${avatar}`
}

onMounted(() => {
  fetchUser()
})
</script>

