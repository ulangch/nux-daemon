<template>
  <div id="app">
    <div class="container">
      <!-- <nav>
        <router-link to="/" class="nav-item">
          <img class="nav-icon" src="./assets/vue.svg" />
          <span class="nav-text">服务状态信息</span>
        </router-link>
        <router-link to="/setting" class="nav-item" style="margin-top: 10px">
          <img class="nav-icon" src="./assets/vue.svg" />
          <span class="nav-text">管理存储空间</span>
        </router-link>
      </nav> -->
      <main>
        <!-- <router-view /> -->
        <ServerStatePage/>
      </main>
      
    </div>
  </div>
</template>

<script setup lang="ts">
import ServerStatePage from './pages/ServerStatePage.vue';
</script>

<style>
.container {
  display: flex;
  height: 100vh;
  background: linear-gradient(to bottom right, #d0ddff, #ffdfc7, #d9f3ff, #ffdabd);
}

nav {
  display: flex;
  flex-direction: column;
  width: 200px;
  /* background-color: #efefef; */
  /* background-image: none; */
  background-color: rgba(255, 255, 255, 0.4);
  border-radius: 16px;
  margin: 16px;
  padding-top: 20px;
  padding-left: 12px;
  padding-right: 12px;
}

.nav-item {
  display: flex;
  align-items: center;
  padding: 14px;
  text-decoration: none;
  color: #303030;
  border-radius: 16px;
}

.nav-icon {
  margin-left: 10px;
  width: 20px;
  height: 20px;
}

.nav-text {
  /* margin-block-start: 20px; */
  margin-left: 12px;
  font-size: 17px;
}

.nav-item.router-link-exact-active {
  color: #ffffff;
  background-color: #3a99ff; /* 蓝色 */
  border-radius: 16px; /* 圆角 */
}

main {
  flex-grow: 1;
  /* background-color: #e6e6e6; */
  background-color: rgba(255, 255, 255, 0.517);
  border-radius: 16px;
  margin-top: 16px;
  margin-left: 16px;
  margin-right: 16px;
  margin-bottom: 16px;
  padding-left: 12px;
  padding-top: 12px;
  padding-left: 12px;
  padding-right: 12px;
}
</style>
