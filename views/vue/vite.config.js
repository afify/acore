// views/vue/vite.config.js
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

export default defineConfig({
  root: '.',               // this folder
  plugins: [vue()],
  resolve: {
    alias: { '@': path.resolve(__dirname, 'src') }
  },
  build: {
    outDir: '../static/js', // emit into views/static/js
    emptyOutDir: false,
    rollupOptions: {
      input: {
        index: 'src/main.js'
      },
      output: {
        entryFileNames: '[name].js'
      }
    }
  }
})
