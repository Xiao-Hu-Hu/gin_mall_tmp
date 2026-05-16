<template>
  <div class="mall-container profile-page">
    <h2 class="page-title">个人中心</h2>
    <el-tabs tab-position="left" v-model="activeTab">
      <!-- 个人信息 -->
      <el-tab-pane label="个人信息" name="info">
        <el-card shadow="never">
          <el-descriptions :column="1" border>
            <el-descriptions-item label="用户名">{{ auth.user?.username }}</el-descriptions-item>
            <el-descriptions-item label="邮箱">{{ auth.user?.email || '未绑定' }}</el-descriptions-item>
            <el-descriptions-item label="状态">{{ auth.user?.status || '正常' }}</el-descriptions-item>
          </el-descriptions>
          <el-divider />
          <el-form inline>
            <el-form-item label="修改昵称"><el-input v-model="nickName" placeholder="新昵称" /></el-form-item>
            <el-form-item><el-button type="primary" :loading="saving1" @click="updateNickName">保存</el-button></el-form-item>
          </el-form>
        </el-card>
      </el-tab-pane>
      <!-- 头像设置 -->
      <el-tab-pane label="头像设置" name="avatar">
        <el-card shadow="never">
          <div style="display:flex;align-items:center;gap:20px;margin-bottom:20px">
            <el-avatar :size="80" :src="auth.user?.avatar" />
            <div>
              <el-upload :auto-upload="false" :limit="1" :on-change="onAvatarChange" accept="image/*" :show-file-list="false">
                <el-button type="primary">选择新头像</el-button>
              </el-upload>
              <div v-if="avatarFile" style="margin-top:12px">
                <el-button type="success" :loading="saving2" @click="doUploadAvatar">上传头像</el-button>
              </div>
            </div>
          </div>
        </el-card>
      </el-tab-pane>
      <!-- 账户安全 -->
      <el-tab-pane label="账户安全" name="email">
        <el-card shadow="never">
          <h4 style="margin-bottom:16px">邮箱验证</h4>
          <el-form label-width="100px">
            <el-form-item label="邮箱"><el-input v-model="emailForm.email" /></el-form-item>
            <el-form-item label="密码"><el-input v-model="emailForm.password" type="password" show-password /></el-form-item>
            <el-form-item label="操作类型">
              <el-select v-model="emailForm.operation_type">
                <el-option :value="1" label="绑定邮箱" />
                <el-option :value="2" label="解绑邮箱" />
                <el-option :value="3" label="修改密码" />
              </el-select>
            </el-form-item>
            <el-form-item><el-button type="primary" :loading="saving3" @click="doSendEmail">发送验证邮件</el-button></el-form-item>
          </el-form>
        </el-card>
      </el-tab-pane>
      <!-- 余额 -->
      <el-tab-pane label="我的余额" name="money">
        <el-card shadow="never">
          <el-form inline>
            <el-form-item label="密钥"><el-input v-model="moneyKey" type="password" placeholder="16位密钥" maxlength="16" show-password /></el-form-item>
            <el-form-item><el-button type="primary" :loading="saving4" @click="doShowMoney">查看余额</el-button></el-form-item>
          </el-form>
          <div v-if="moneyResult" style="margin-top:16px">
            <el-result icon="success" :title="`余额: ¥${moneyResult.user_money}`" :sub-title="moneyResult.user_name" />
          </div>
        </el-card>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useAuthStore } from '../stores/auth'
import { updateUser, uploadAvatar, sendEmail, showMoney } from '../api/auth'
import { ElMessage } from 'element-plus'

const auth = useAuthStore()
const activeTab = ref('info')
const nickName = ref(auth.user?.nickname || '')
const avatarFile = ref(null)
const saving1 = ref(false)
const saving2 = ref(false)
const saving3 = ref(false)
const saving4 = ref(false)
const emailForm = reactive({ email: '', password: '', operation_type: 1 })
const moneyKey = ref('')
const moneyResult = ref(null)

async function updateNickName() {
  if (!nickName.value) return ElMessage.warning('请输入昵称')
  saving1.value = true
  try {
    const res = await updateUser({ nick_name: nickName.value })
    if (res.data) auth.setAuth(auth.token, { ...auth.user, ...res.data })
    ElMessage.success('修改成功')
  } catch (e) { /* ignore */ } finally { saving1.value = false }
}

function onAvatarChange(file) { avatarFile.value = file.raw }

async function doUploadAvatar() {
  saving2.value = true
  try {
    const res = await uploadAvatar(avatarFile.value)
    if (res.data) auth.setAuth(auth.token, { ...auth.user, ...res.data })
    avatarFile.value = null
    ElMessage.success('头像上传成功')
  } catch (e) { /* ignore */ } finally { saving2.value = false }
}

async function doSendEmail() {
  saving3.value = true
  try { await sendEmail(emailForm); ElMessage.success('验证邮件已发送') } catch (e) { /* ignore */ } finally { saving3.value = false }
}

async function doShowMoney() {
  if (!moneyKey.value || moneyKey.value.length !== 16) return ElMessage.warning('请输入16位密钥')
  saving4.value = true
  try { const res = await showMoney({ key: moneyKey.value }); moneyResult.value = res.data } catch (e) { /* ignore */ } finally { saving4.value = false }
}
</script>

<style scoped>
.profile-page { padding: 20px 0 40px; }
.page-title { font-size: 22px; margin-bottom: 20px; }
</style>
