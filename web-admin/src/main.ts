import { createApp } from 'vue'
import { createPinia } from 'pinia'
import Antd from 'ant-design-vue'
import App from './App.vue'
import router from './router'
import { useAuthStore } from './stores/auth'
import 'ant-design-vue/dist/reset.css'
import './assets/styles/main.css'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
app.use(Antd)

// Initialize auth store from localStorage
const authStore = useAuthStore()
authStore.initAuth()

app.mount('#app')
