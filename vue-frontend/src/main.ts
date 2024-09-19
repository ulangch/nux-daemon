import { createApp } from 'vue'
import './style.css'
import App from './App.vue'


import {createRouter, createWebHistory } from 'vue-router';
import DiskSettingPage from './pages/DiskSettingPage.vue';
import ServerStatePage from './pages/ServerStatePage.vue';

const routes = [
  {path: "/", component: ServerStatePage},
  {path: "/state", component: ServerStatePage},
  {path: "/setting", component: DiskSettingPage}
]
const router = createRouter({
  history: createWebHistory(),
  routes
})

const app = createApp(App)
app.use(router)

app.mount('#app')
