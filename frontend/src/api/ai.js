import api from './index'

export const aiChat = (data) => api.post('/ai/customer-service/chat', data)
export const getAiSessions = (params) => api.get('/ai/sessions', { params })
export const getAiSessionHistory = (sessionId) => api.get(`/ai/sessions/${sessionId}`, { suppressErrorToast: true })
export const deleteAiSession = (sessionId) => api.delete(`/ai/sessions/${sessionId}`)
