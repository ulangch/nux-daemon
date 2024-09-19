// electron/main.js
const { app, BrowserWindow } = require('electron');
const { join } = require('path');
const { exec } = require('child_process');

// 屏蔽安全警告
// ectron Security Warning (Insecure Content-Security-Policy)
process.env['ELECTRON_DISABLE_SECURITY_WARNINGS'] = 'true';

// 创建浏览器窗口时，调用这个函数。
const createWindow = () => {
  const win = new BrowserWindow({
    width: 500,
    height: 800,
    title: '我的私有云',
    icon: join(__dirname, '../dist/hasky.png'),
  });

  // win.loadURL('http://localhost:3000')
  // development模式
  if (process.env.VITE_DEV_SERVER_URL) {
    win.loadURL(process.env.VITE_DEV_SERVER_URL);
    // 开启调试台
    // win.webContents.openDevTools();
  } else {
    win.loadFile(join(__dirname, '../dist/index.html'));
  }
};

// Electron 会在初始化后并准备
app.whenReady().then(() => {
  createWindow();
  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) createWindow();
  });
});

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') app.quit();
});

function runAsAdmin(command, callback) {
  const execCommand = `powershell.exe -Command "Start-Process cmd -ArgumentList '/c ${command}' -Verb RunAs"`;
  exec(execCommand, (error, stdout, stderr) => {
    if (error) {
      console.error(`Error: ${error.message}`);
      return;
    }
    if (stderr) {
      console.error(`Stderr: ${stderr}`);
      return;
    }
    console.log(`Stdout: ${stdout}`);
    if (callback) callback(stdout);
  });
}

// Expose the function to the renderer process
const { ipcMain } = require('electron');
ipcMain.handle('run-as-admin', async (event, command) => {
  return new Promise((resolve, reject) => {
    runAsAdmin(command, result => {
      resolve(result);
    });
  });
});
