import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '../utils/api'

export const useUserStore = defineStore('user', () => {
  const user = ref(null)
  const token = ref(localStorage.getItem('token') || '')

  const isLoggedIn = computed(() => !!token.value)

  const login = async (userData, userToken) => {
    user.value = userData
    token.value = userToken
    localStorage.setItem('token', userToken)
    api.defaults.headers.common['token'] = userToken
  }

  const logout = () => {
    user.value = null
    token.value = ''
    localStorage.removeItem('token')
    delete api.defaults.headers.common['token']
  }

  // 初始化时设置token
  if (token.value) {
    api.defaults.headers.common['token'] = token.value
  }

  return {
    user,
    token,
    isLoggedIn,
    login,
    logout
  }
})

