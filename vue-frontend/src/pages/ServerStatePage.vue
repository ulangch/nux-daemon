<template>
  <div>
    <div style="display: flex; flex-direction: row; margin-bottom: 16px; align-items: center">
      <span class="bar-title">服务状态信息</span>
      <button class="bar-refresh">刷新</button>
    </div>
    <div class="sec" style="align-items: center">
      <div>
        <div style="display: flex; flex-direction: row">
          <span class="sec-title">设备名称：&nbsp;</span>
          <span class="sec-content">我的私有云</span>
        </div>
        <span class="sec-hint" style="margin-top: 6px">设备名称将会展现在Android/iOS APP上。</span>
      </div>
      <img src="../assets/ic_edit_66.svg" style="width: 22px; height: 22px; margin-left: auto; padding-left: 12px; padding-right: 12px" />
    </div>
    <div class="sec" style="flex-direction: column">
      <div style="display: flex; flex-direction: row">
        <span class="sec-title">设备标识：&nbsp;</span>
        <span class="sec-content">20632f6e-30ae-4770-ab3f-634a4a0d4f6e</span>
      </div>
      <span class="sec-hint" style="margin-top: 6px">设备的唯一可识别ID，自动生成不可更改。</span>
    </div>
    <div class="sec" style="align-items: center">
      <div class="sec-sub" style="flex-direction: column">
        <span class="sec-title">本次运行</span>
        <span class="sec-content" style="margin-top: 16px">0 天 0 时 51 分 47 秒</span>
      </div>
      <div class="sec-sub" style="flex-direction: row">
        <CircularProgress :percentage="cpuUsage.usage" :label="cpuUsage.label" style="margin-left: auto" />
        <CircularProgress :percentage="memoryUsage.usage" :label="memoryUsage.label" style="margin-left: auto" />
      </div>
    </div>
    <div class="sec" style="align-items: center">
      <div>
        <span class="sec-title">局域网地址</span>
        <span class="href-link" style="margin-top: 6px">http://192.168.31.20:8080</span>
        <span class="sec-hint" style="margin-top: 30px">使用APP扫描右侧二维码可主动绑定设备。</span>
      </div>
      <div class="sec-right">
        <div style="width: 100px; height: 100px; margin-left: auto; background-color: #d7d7d7">
          <img style="width: 100%; height: 100%" v-if="false" />
        </div>
      </div>
    </div>
    <div class="sec" style="align-items: center">
      <div>
        <div style="display: flex; flex-direction: row">
          <span class="sec-title">存储空间：&nbsp;</span>
          <span class="sec-content">F:\我的资源库</span>
        </div>
        <span class="sec-hint" style="margin-top: 6px">更改存储空间后，原文件将无法在Android/iOS APP上访问。</span>
      </div>
      <img src="../assets/ic_edit_66.svg" style="width: 22px; height: 22px; margin-left: auto; padding-left: 12px; padding-right: 12px" />
    </div>

    <div style="display: flex; flex-direction: column; align-items: center; width: 100%; margin-top: 40px">
      <button class="start-button">重新启动服务</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue';
import CircularProgress from '@/components/CircularProgress.vue';

const cpuUsage = reactive({
  label: 'CPU',
  usage: 20,
});

const memoryUsage = reactive({
  label: '内存',
  usage: 50,
});
</script>

<style scoped>
.bar-title {
  display: flex;
  margin-left: 4px;
  color: #303030;
  font-size: 15px;
  font-weight: bold;
}
.bar-refresh {
  display: flex;
  margin-left: auto;
  margin-right: 4px;
  color: #007bff;
  font-size: 15px;
  font-weight: bold;
  background-color: #e7e7e7;
  border-radius: 30px;
  padding-top: 4px;
  padding-bottom: 4px;
  padding-left: 18px;
  padding-right: 18px;
  border-width: 0px;
  stroke-width: 0px;
}
.sec {
  display: flex;
  padding: 14px;
  background-color: white;
  margin-top: 12px;
  border-radius: 16px;
}
.sec-sub {
  flex: 1;
  display: flex;
  flex-direction: row;
}
.sec-right {
  display: flex;
  flex-grow: 1;
  flex-direction: row;
  justify-content: right;
}
.sec-title {
  display: flex;
  color: #303030;
  font-weight: 500;
  font-size: 15px;
}
.sec-content {
  display: flex;
  color: #676767;
  font-weight: 400;
  font-size: 15px;
}
.sec-hint {
  display: flex;
  color: #676767;
  font-weight: 320;
  font-size: 11px;
}
.href-link {
  display: flex;
  color: #007bff;
  font-weight: 320;
  font-size: 14px;
}
.start-button {
  width: 100%;
  color: white;
  font-size: 15px;
  font-weight: 500;
  border-radius: 80px;
  border-width: 0px;
  background-color: #007bff;
  height: 50px;
}
</style>
