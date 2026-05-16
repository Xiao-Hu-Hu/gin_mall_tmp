<template>
  <div class="mall-container create-product-page">
    <h2 class="page-title">发布商品</h2>
    <el-card shadow="never">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="120px" style="max-width:700px">
        <el-form-item label="商品名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入商品名称" />
        </el-form-item>
        <el-form-item label="商品分类" prop="category_id">
          <el-select v-model="form.category_id" placeholder="请选择分类">
            <el-option v-for="cat in categories" :key="cat.id" :label="cat.category_name" :value="cat.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="商品标题" prop="title">
          <el-input v-model="form.title" placeholder="商品展示标题" />
        </el-form-item>
        <el-form-item label="商品描述">
          <el-input v-model="form.info" type="textarea" :rows="4" placeholder="详细描述商品信息" />
        </el-form-item>
        <el-form-item label="原价" prop="price">
          <el-input v-model="form.price" placeholder="0.00" style="width:200px">
            <template #prepend>¥</template>
          </el-input>
        </el-form-item>
        <el-form-item label="折扣价" prop="discount_price">
          <el-input v-model="form.discount_price" placeholder="0.00" style="width:200px">
            <template #prepend>¥</template>
          </el-input>
        </el-form-item>
        <el-form-item label="库存数量" prop="num">
          <el-input-number v-model="form.num" :min="0" :max="99999" />
        </el-form-item>
        <el-form-item label="是否上架">
          <el-switch v-model="form.on_sale" active-text="上架" inactive-text="下架" />
        </el-form-item>
        <el-form-item label="商品图片" prop="file">
          <el-upload
            :auto-upload="false"
            :on-change="onFileChange"
            :file-list="fileList"
            list-type="picture-card"
            :limit="5"
            accept="image/*"
            :on-exceed="() => ElMessage.warning('最多上传5张图片')"
          >
            <el-icon><Plus /></el-icon>
          </el-upload>
          <div class="form-tip">支持 jpg/png 格式，最多5张，第一张为主图</div>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" size="large" :loading="submitting" @click="handleSubmit">发布商品</el-button>
          <el-button size="large" @click="$router.back()">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getCategories, createProduct } from '../api/product'
import { ElMessage } from 'element-plus'

const router = useRouter()
const formRef = ref()
const submitting = ref(false)
const categories = ref([])
const fileList = ref([])
const files = []

const form = reactive({
  name: '', category_id: undefined, title: '', info: '',
  price: '', discount_price: '', num: 100, on_sale: true
})

const rules = {
  name: [{ required: true, message: '请输入商品名称', trigger: 'blur' }],
  category_id: [{ required: true, message: '请选择分类', trigger: 'change' }],
  title: [{ required: true, message: '请输入商品标题', trigger: 'blur' }],
  price: [{ required: true, message: '请输入价格', trigger: 'blur' }],
  discount_price: [{ required: true, message: '请输入折扣价', trigger: 'blur' }],
  num: [{ required: true, message: '请输入库存', trigger: 'blur' }]
}

function onFileChange(file) {
  files.push(file.raw)
}

async function handleSubmit() {
  await formRef.value.validate()
  if (!files.length) return ElMessage.warning('请至少上传一张商品图片')
  submitting.value = true
  try {
    await createProduct({ ...form, file: files })
    ElMessage.success('商品发布成功')
    router.push('/products')
  } catch (e) { /* interceptor handles */ } finally { submitting.value = false }
}

onMounted(async () => {
  try { const res = await getCategories(); categories.value = res.data?.item || [] } catch (e) { /* ignore */ }
})
</script>

<style scoped>
.create-product-page { padding: 20px 0 40px; }
.page-title { font-size: 22px; margin-bottom: 20px; }
.form-tip { font-size: 12px; color: var(--text-light); margin-top: 8px; }
</style>
