import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'


// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  // 前端解决跨域
  server:{
    proxy:{
      '/user':{
        target:"http://127.0.0.1:8080/"
      }
    }
  },
  base: './',// 打包相对路径  yarn build
})


