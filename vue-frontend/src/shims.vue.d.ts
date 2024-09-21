declare module '*.vue' {
  import { DefineComponent } from 'vue';
  const component: DefineComponent<{}, {}, any>;
  export default component;
}

interface ElectronLogAPI {
  info: (message: string) => void;
  error: (message: string) => void;
  warn: (message: string) => void;
}

interface ElectronAPI {
  log: ElectronLogAPI;
  selectFolder: () => Promise<string>;
}

interface Window {
  electron: ElectronAPI;
}