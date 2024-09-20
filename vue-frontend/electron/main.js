// electron/main.js
const { app, BrowserWindow, Menu } = require('electron');
const { join } = require('path');
const { spawn } = require('child_process')
const path = require('path');

// 屏蔽安全警告
process.env['ELECTRON_DISABLE_SECURITY_WARNINGS'] = 'true';

// 创建浏览器窗口时，调用这个函数。
const createWindow = () => {
  const win = new BrowserWindow({
    width: 500,
    height: 800,
    title: '我的私有云',
    icon: join(__dirname, '../dist/hasky.png'),
  });

  Menu.setApplicationMenu(null);

  // win.loadURL('http://localhost:3000')
  if (process.env.VITE_DEV_SERVER_URL) {
    win.loadURL(process.env.VITE_DEV_SERVER_URL);
    // 开启调试台
    // win.webContents.openDevTools();
  } else {
    win.loadFile(join(__dirname, '../dist/index.html'));
  }
};

const log = require('electron-log');
log.transports.file.resolvePathFn = () => path.join(app.getPath('logs'), 'main.log')

let daemon

const startDaemon = () => {
  let daemonPath
  if (process.env.VITE_DEV_SERVER_URL) {
    daemonPath = path.join(app.getAppPath(), 'public', 'nas-daemon.exe');
  } else {
    daemonPath = path.join(process.resourcesPath, 'public', 'nas-daemon.exe');
  }
  daemon = spawn(daemonPath, [], { stdio: ['ignore', 'pipe', 'pipe'] });
  // daemon = spawn('powershell.exe', ['-Command', `Start-Process -FilePath ${daemonPath} -Verb RunAs`], {stdio: 'inherit'})

  daemon.stdout.on('data', (data) => {
    log.info(data.toString('utf8'));
  });
  daemon.stderr.on('data', (data) => {
    log.error(data.toString('utf8'));
  });
  daemon.on('error', (error) => {
    log.error('Start daemon failed: ' + error.message);
  });
  daemon.on('close', (code) => {
    log.info('daemon exist with code: ' + code)
  });
};

// Electron 会在初始化后并准备
app.whenReady().then(() => {
  startDaemon();
  createWindow();
  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) createWindow();
  });
});

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') app.quit();
});

app.on('will-quit', () => {
  if (daemon) {
    daemon.kill()
  }
});