const { contextBridge, ipcRenderer } = require("electron");
const log = require('electron-log');
// log.transports.file.resolvePath = () => path.join(app.getPath('logs'), 'main.log')

contextBridge.exposeInMainWorld("electron", {
  log: {
      info: (message) => log.info(message),
      error: (message) => log.error(message),
      warn: (message) => log.warn(message)
  },
  selectFolder: () => ipcRenderer.invoke('select-folder')
})