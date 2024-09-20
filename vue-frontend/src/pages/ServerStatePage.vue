<template>
  <div>
    <div style="display: flex; flex-direction: row; margin-bottom: 16px; align-items: center">
      <span class="bar-title">服务状态信息</span>
      <button class="bar-refresh">刷新</button>
    </div>
    <div>
      <span v-if="notice.length > 0">{{ notice }}</span>
    </div>
    <div class="sec" style="align-items: center">
      <div>
        <div style="display: flex; flex-direction: row">
          <span class="sec-title">设备名称：&nbsp;</span>
          <span class="sec-content">{{ deviceName }}</span>
        </div>
        <span class="sec-hint" style="margin-top: 6px">设备名称将会展现在Android/iOS APP上。</span>
      </div>
      <img src="../assets/ic_edit_66.svg" style="width: 22px; height: 22px; margin-left: auto; padding-left: 12px; padding-right: 12px" />
    </div>
    <div class="sec" style="flex-direction: column">
      <div style="display: flex; flex-direction: row">
        <span class="sec-title">设备标识：&nbsp;</span>
        <span class="sec-content">{{ deviceId }}</span>
      </div>
      <span class="sec-hint" style="margin-top: 6px">设备的唯一可识别ID，自动生成不可更改。</span>
    </div>
    <div class="sec" style="align-items: center">
      <div class="sec-sub" style="flex-direction: column">
        <span class="sec-title">本次运行</span>
        <span class="sec-content" style="margin-top: 16px">{{ aliveDuration }}</span>
      </div>
      <div class="sec-sub" style="flex-direction: row">
        <CircularProgress :percentage="cpu.usage" :label="cpu.label" style="margin-left: auto" />
        <CircularProgress :percentage="memory.usage" :label="memory.label" style="margin-left: auto" />
      </div>
    </div>
    <div class="sec" style="align-items: center">
      <div>
        <span class="sec-title">局域网地址</span>
        <span class="href-link" style="margin-top: 6px">{{ deviceUrl }}</span>
        <span class="sec-hint" style="margin-top: 30px">使用APP扫描右侧二维码可主动绑定设备。</span>
      </div>
      <div class="sec-right">
        <div style="width: 100px; height: 100px; margin-left: auto; background-color: #d7d7d7">
          <img style="width: 100%; height: 100%" v-if="qrcode" :src="qrcode" alt="Dynamic Image" />
        </div>
      </div>
    </div>
    <div class="sec" style="align-items: center">
      <div>
        <div style="display: flex; flex-direction: row">
          <span class="sec-title">存储空间：&nbsp;</span>
          <span class="sec-content">{{ disk }}</span>
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
import { onMounted, reactive, ref } from 'vue';
import CircularProgress from '@/components/CircularProgress.vue';
import axios from 'axios';

import log from 'electron-log'

const notice = ref('')
const deviceName = ref('');
const deviceId = ref('');
const aliveDuration = ref('0 天 0 时 0 分 1 秒');
const cpu = reactive({
  label: 'CPU',
  usage: 0,
});
const memory = reactive({
  label: '内存',
  usage: 0,
});
const deviceUrl = ref('');
const qrcode = ref('');
const disk = ref('');

const fetchServiceInfo = async () => {
  try {
    const resp = await axios.get('http://localhost:8080/service/info');
    const data = resp.data;
    log.info('fetch service info, data=${data}')
    // notice.value = data
    if (data.status_code != 0) {
      if (data.status_code == 4) {
        notice.value = '请添加磁盘后刷新页面。'
      } else {
        notice.value = '当前服务不可用，请使用管理员权限打开APP。'
      }
    } else {
      deviceName.value = data.name;
      deviceId.value = data.id;
      cpu.usage = data.cpu;
      memory.usage = data.memory;
      deviceUrl.value = data.url;
      qrcode.value = data.qrcode;
      disk.value = data.disk;
    }
  } catch (error) {
    log.error('fetch service info failed: ', error)
    // notice.value = '当前服务不可用。'
    notice.value = String(error)
  }
};

onMounted(() => {
  fetchServiceInfo()
})

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
