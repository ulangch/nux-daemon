<template>
  <div>
    <div style="display: flex; flex-direction: row; align-items: center">
      <span class="bar-title">设备状态信息</span>
      <button class="bar-refresh" @click="handleRefresh">刷新</button>
    </div>
    <div style="display: flex; flex-direction: column;">
      <span v-if="notice" class="top-notice">{{ notice }}</span>
    </div>
    <div class="sec" style="align-items: center; margin-top: 16px;">
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
          <img style="width: 100%; height: 100%" v-if="qrcode" :src="qrcode" />
        </div>
      </div>
    </div>
    <div class="sec" style="align-items: center">
      <div>
        <div style="display: flex; flex-direction: row; align-items: center;">
          <span class="sec-title">存储空间：&nbsp;</span>
          <span class="sec-content" v-if="disk">{{ disk }}</span>
          <span class="sec-content-notice" v-if="diskNotice">{{ diskNotice }}</span>
        </div>
        <span class="sec-hint" style="margin-top: 6px">更改存储空间后，原文件将无法在Android/iOS APP上访问。</span>
      </div>
      <img src="../assets/ic_edit_66.svg" style="width: 22px; height: 22px; margin-left: auto; padding-left: 12px; padding-right: 12px" @click="selectFolder" />
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
const qrcode = ref('')
const disk = ref('');
const diskNotice = ref('');

const handleRefresh = () => {
  fetchServiceInfo()
}

const fetchServiceInfo = async () => {
  try {
    const resp = await axios.get('http://localhost:8080/service/info');
    const data = resp.data;
    const info = data.info
    if (!data || data.status_code != 0 || !info) {
       notice.value = '当前服务不可用，请使用管理员权限打开APP。'
    } else {
      deviceName.value = info.name;
      deviceId.value = info.id;
      cpu.usage = info.cpu;
      memory.usage = info.memory;
      deviceUrl.value = info.url;
      qrcode.value = `file://${info.qrcode.replaceAll('\\', '/')}`
      if (info.disk) {
        notice.value = ''
        disk.value = info.disk;
        diskNotice.value = ''
      } else {
        notice.value = '当前还没有为设备添加存储空间，Android/iOS APP将无法正常访问，请点击【存储空间】右侧【编辑】图标添加。'
        disk.value = ''
        diskNotice.value = '点击右侧【编辑】图标添加存储空间。'
      }
    }
  } catch (error) {
    window.electron.log.error('fetch service info failed')
    notice.value = '当前服务不可用，请使用管理员权限打开APP。'
  }
};

const selectFolder = async () => {
  try {
    const path = await window.electron.selectFolder()
    if (path) {
      await axios.post(`http://localhost:8080/service/update_disk?path=${encodeURIComponent(path)}`)
    }
    fetchServiceInfo()
  } catch (error) {
    window.electron.log.error('select folder failed')
  }
}

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
.top-notice {
  color: #ff0000;
  font-size: 13px;
  margin-left: 4px;
  margin-top: 12px;
  font-weight: 400;
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
.sec-content-notice {
  display: flex;
  color: #ff0000;
  font-weight: 400;
  font-size: 13px;
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
