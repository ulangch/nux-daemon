<template>
    <div v-if="visible" class="modal-overlay">
      <div class="modal">
        <span class="modal-title">{{ title }}</span>
        <input v-model="inputText" placeholder="输入设备名" @keydown.enter="confirm"/>
        <div class="modal-buttons">
          <button class="button-cancel" @click="cancel">取消</button>
          <button class="button-confirm" @click="confirm">确认</button>
        </div>
      </div>
    </div>
  </template>
  
  <script lang="ts">
  import { defineComponent, ref, watch } from 'vue'
  
  export default defineComponent({
    name: 'EditNameDialog',
    props: {
      visible: {
        type: Boolean,
        required: true
      },
      title: {
        type: String,
        default: '弹框'
      }
    },
    emits: ['confirm', 'cancel'],
    setup(props, { emit }) {
      const inputText = ref('')
  
      const confirm = () => {
        emit('confirm', inputText.value)
      }
  
      const cancel = () => {
        emit('cancel')
      }
  
      watch(() => props.visible, (newVal) => {
        if (!newVal) {
          inputText.value = ''
        }
      })
  
      return {
        inputText,
        confirm,
        cancel
      }
    }
  })
  </script>
  
  <style scoped>
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
  }
  
  .modal {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    background: #fff;
    width: 70%;
    padding-top: 20px;
    padding-bottom: 30px;
    padding-left: 30px;
    padding-right: 30px;
    margin-bottom: 20%;
    border-radius: 16px;
    text-align: center;
  }
  .modal-title {
    font-size: 16px;
    color: #303030;
    font-weight: bold;
  }
  
  input {
    margin-top: 25px;
    padding-top: 14px;
    padding-bottom: 14px;
    padding-left: 4%;
    padding-right: 4%;
    width: 92%;
    border-width: 0px;
    border-radius: 8px;
    background-color: #ededed;
    color: #303030;
  }

  input:focus {
    border-color: #007bff;
  }

  .modal-buttons {
    display: flex;
    flex-direction: row;
    margin-top: 30px;
    align-items: center;
    justify-content: center;
    width: 100%;
  }
  
  .button-cancel {
    width: 100%;
    padding: 12px;
    margin-right: 12px;
    border-radius: 8px;
    border-width: 0px;
    color: #303030;
    background-color: #ededed;
    font-size: 14px;
  }
  .button-confirm {
    width: 100%;
    padding: 12px;
    margin-left: 12px;
    border-radius: 8px;
    border-width: 0px;
    background-color: #007bff;
    color: white;
    font-size: 14px;
  }
  </style>