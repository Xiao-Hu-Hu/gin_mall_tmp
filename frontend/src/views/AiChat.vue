<template>
  <div class="mall-container ai-chat-page">
    <h2 class="page-title">AI 客服</h2>
    <div class="chat-layout">
      <!-- Sessions -->
      <aside class="session-sidebar">
        <el-button type="primary" size="small" @click="newSession" style="width:100%;margin-bottom:12px">新会话</el-button>
        <div v-for="s in sessions" :key="s.session_id" class="session-item" :class="{ active: sessionId === s.session_id }" @click="loadSession(s.session_id)">
          <div class="session-text">{{ s.last_message || s.session_id }}</div>
          <el-button text type="danger" size="small" @click.stop="deleteSession(s.session_id)"><el-icon><Delete /></el-icon></el-button>
        </div>
        <el-empty v-if="!sessions.length" description="暂无会话" :image-size="40" />
      </aside>
      <!-- Chat -->
      <div class="chat-main">
        <div class="messages" ref="messagesRef">
          <div v-for="(msg, idx) in messages" :key="idx" :class="['msg', msg.role]">
            <div class="bubble">
              <div class="content">{{ msg.content }}</div>
              <!-- Recommendations -->
              <div v-if="msg.recommendations?.length" class="recommendations">
                <div v-for="p in msg.recommendations" :key="p.id" class="rec-card">
                  <router-link :to="`/products/${p.id}`">
                    <img :src="p.img_path" />
                    <div class="rec-info">
                      <div class="name">{{ p.name }}</div>
                      <div class="price">&yen;{{ formatPrice(p.discount_price || p.price) }}</div>
                    </div>
                  </router-link>
                </div>
              </div>
            </div>
          </div>
          <div v-if="chatLoading" class="msg ai"><div class="bubble"><span class="typing">AI 正在思考...</span></div></div>
        </div>
        <div class="input-area">
          <el-input v-model="inputMsg" type="textarea" :rows="2" placeholder="输入消息..." @keydown.enter.ctrl="sendMessage" />
          <el-button type="primary" :loading="chatLoading" @click="sendMessage" :disabled="!inputMsg.trim()">发送</el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { aiChat, getAiSessions, getAiSessionHistory, deleteAiSession } from '../api/ai'
import { formatPrice } from '../utils/format'
import { ElMessage } from 'element-plus'

const messages = ref([])
const sessions = ref([])
const sessionId = ref('')
const inputMsg = ref('')
const chatLoading = ref(false)
const messagesRef = ref(null)

function scrollBottom() { nextTick(() => { if (messagesRef.value) messagesRef.value.scrollTop = messagesRef.value.scrollHeight }) }

async function loadSessions() {
  try {
    const res = await getAiSessions({ limit: 20 })
    const list = res.sessions || res.data?.sessions || []
    sessions.value = list.map(sid => typeof sid === 'string' ? { session_id: sid, last_message: sid.slice(0, 12) + '...' } : sid)
  } catch (e) { /* ignore */ }
}

function newSession() { sessionId.value = ''; messages.value = [] }

async function loadSession(sid) {
  sessionId.value = sid
  messages.value = []
  try {
    const res = await getAiSessionHistory(sid)
    const history = res.history || res.data?.history || []
    messages.value = history.map(m => ({
      role: m.role === 'assistant' ? 'ai' : m.role,
      content: m.content,
      recommendations: m.recommendations || []
    }))
    scrollBottom()
  } catch (e) { /* ignore */ }
}

async function deleteSession(sid) {
  try { await deleteAiSession(sid); if (sessionId.value === sid) newSession(); loadSessions() } catch (e) { /* ignore */ }
}

async function sendMessage() {
  const text = inputMsg.value.trim()
  if (!text) return
  inputMsg.value = ''
  messages.value.push({ role: 'user', content: text })
  scrollBottom()
  chatLoading.value = true
  try {
    const payload = { message: text }
    if (sessionId.value) payload.session_id = sessionId.value
    payload.history = messages.value.slice(0, -1).map(m => ({ role: m.role === 'ai' ? 'assistant' : m.role, content: m.content }))
    const res = await aiChat(payload)
    const data = res.data
    if (data?.session_id) sessionId.value = data.session_id
    messages.value.push({ role: 'ai', content: data?.reply || '抱歉，暂时无法回复', recommendations: data?.recommendations || [] })
    scrollBottom()
    loadSessions()
  } catch (e) {
    messages.value.push({ role: 'ai', content: '请求失败，请稍后重试' })
  } finally { chatLoading.value = false }
}

onMounted(loadSessions)
</script>

<style scoped>
.ai-chat-page { padding: 20px 0 40px; }
.page-title { font-size: 22px; margin-bottom: 20px; }
.chat-layout { display: flex; gap: 16px; height: 600px; background: #fff; border-radius: 8px; overflow: hidden; }
.session-sidebar { width: 220px; background: #f9f9f9; padding: 12px; overflow-y: auto; border-right: 1px solid var(--border); }
.session-item { display: flex; align-items: center; gap: 8px; padding: 8px; border-radius: 6px; cursor: pointer; font-size: 13px; margin-bottom: 4px; }
.session-item:hover { background: #eee; }
.session-item.active { background: #fff0f0; color: var(--primary); }
.session-text { flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.chat-main { flex: 1; display: flex; flex-direction: column; }
.messages { flex: 1; overflow-y: auto; padding: 16px; display: flex; flex-direction: column; gap: 12px; }
.msg { display: flex; }
.msg.user { justify-content: flex-end; }
.msg.ai { justify-content: flex-start; }
.bubble { max-width: 70%; padding: 10px 14px; border-radius: 12px; font-size: 14px; line-height: 1.6; }
.msg.user .bubble { background: var(--primary); color: #fff; border-bottom-right-radius: 4px; }
.msg.ai .bubble { background: #f0f0f0; color: var(--text); border-bottom-left-radius: 4px; }
.typing { color: #999; }
.recommendations { display: flex; flex-wrap: wrap; gap: 8px; margin-top: 10px; }
.rec-card { width: 150px; border: 1px solid var(--border); border-radius: 6px; overflow: hidden; }
.rec-card a { color: var(--text); }
.rec-card img { width: 100%; height: 100px; object-fit: cover; }
.rec-info { padding: 6px; }
.rec-info .name { font-size: 12px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.rec-info .price { font-size: 14px; color: var(--price); font-weight: bold; }
.input-area { padding: 12px; border-top: 1px solid var(--border); display: flex; gap: 8px; align-items: flex-end; }
.input-area .el-input { flex: 1; }
</style>
