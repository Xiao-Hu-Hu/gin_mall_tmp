<template>
  <div class="mall-container addresses-page">
    <div class="page-header">
      <h2>收货地址</h2>
      <el-button type="primary" @click="openDialog(null)">新增地址</el-button>
    </div>
    <div v-if="addresses.length" class="addr-list">
      <el-card v-for="addr in addresses" :key="addr.id" shadow="hover" class="addr-card" :class="{ selectable: isSelectMode }" @click="isSelectMode && selectAddr(addr)">
        <div class="addr-info">
          <div><strong>{{ addr.name }}</strong> <span style="color:#999">{{ addr.phone }}</span></div>
          <div style="font-size:13px;color:#666;margin-top:6px">{{ addr.address }}</div>
        </div>
        <div class="addr-actions" v-if="!isSelectMode">
          <el-button text type="primary" @click.stop="openDialog(addr)">编辑</el-button>
          <el-button text type="danger" @click.stop="handleDelete(addr.id)">删除</el-button>
        </div>
      </el-card>
    </div>
    <el-empty v-else description="暂无收货地址">
      <el-button type="primary" @click="openDialog(null)">添加新地址</el-button>
    </el-empty>
    <!-- Dialog -->
    <el-dialog v-model="dialogVisible" :title="editId ? '编辑地址' : '新增地址'" width="480px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="收件人" prop="name"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="手机号" prop="phone"><el-input v-model="form.phone" maxlength="11" /></el-form-item>
        <el-form-item label="详细地址" prop="address"><el-input v-model="form.address" type="textarea" :rows="3" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getAddresses, createAddress, updateAddress, deleteAddress } from '../api/address'
import { ElMessage, ElMessageBox } from 'element-plus'

const route = useRoute()
const router = useRouter()
const addresses = ref([])
const isSelectMode = ref(route.query.mode === 'select')
const dialogVisible = ref(false)
const editId = ref(null)
const saving = ref(false)
const formRef = ref()
const form = reactive({ name: '', phone: '', address: '' })
const rules = { name: [{ required: true, message: '请输入收件人', trigger: 'blur' }], phone: [{ required: true, message: '请输入手机号', trigger: 'blur' }], address: [{ required: true, message: '请输入地址', trigger: 'blur' }] }

async function loadAddresses() { try { const res = await getAddresses(); addresses.value = res.data || [] } catch (e) { /* ignore */ } }

function openDialog(addr) {
  editId.value = addr?.id || null
  form.name = addr?.name || ''
  form.phone = addr?.phone || ''
  form.address = addr?.address || ''
  dialogVisible.value = true
}

async function handleSave() {
  await formRef.value.validate()
  saving.value = true
  try {
    if (editId.value) await updateAddress(editId.value, { ...form })
    else await createAddress({ ...form })
    ElMessage.success('保存成功')
    dialogVisible.value = false
    loadAddresses()
  } catch (e) { /* interceptor handles */ } finally { saving.value = false }
}

function handleDelete(id) {
  ElMessageBox.confirm('确定删除此地址？', '提示', { type: 'warning' }).then(async () => {
    await deleteAddress(id)
    loadAddresses()
  }).catch(() => {})
}

function selectAddr(addr) {
  // Store selected address in sessionStorage for checkout to pick up
  sessionStorage.setItem('selectedAddress', JSON.stringify(addr))
  router.back()
}

onMounted(loadAddresses)
</script>

<style scoped>
.addresses-page { padding: 20px 0 40px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.addr-list { display: flex; flex-direction: column; gap: 12px; }
.addr-card { display: flex; justify-content: space-between; align-items: center; }
.addr-card.selectable { cursor: pointer; }
.addr-card.selectable:hover { border-color: var(--primary); }
.addr-info { flex: 1; }
.addr-actions { display: flex; gap: 4px; }
</style>
