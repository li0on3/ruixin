import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'

export default defineConfig({
  plugins: [
    vue(),
    AutoImport({
      resolvers: [ElementPlusResolver()],
      imports: ['vue', 'vue-router', 'pinia']
    }),
    Components({
      resolvers: [ElementPlusResolver()]
    })
  ],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
      // 配置 Vue 以支持运行时编译
      'vue': 'vue/dist/vue.esm-bundler.js'
    }
  },
  build: {
    // 生产环境构建优化
    target: 'es2015',
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: true,
        drop_debugger: true
      }
    },
    // 代码分割优化
    rollupOptions: {
      output: {
        manualChunks: {
          // 将 Vue 生态系统打包为独立块
          vue: ['vue', 'vue-router', 'pinia'],
          // Element Plus 独立打包
          'element-plus': ['element-plus'],
          // ECharts 独立打包（如果使用）
          echarts: ['echarts'],
          // 工具库独立打包
          utils: ['axios', 'dayjs', 'nprogress']
        }
      }
    },
    // 启用 gzip 压缩提示
    reportCompressedSize: false,
    // 设置 chunk 大小警告限制
    chunkSizeWarningLimit: 1000
  },
  server: {
    host: '0.0.0.0',
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:8080',
        changeOrigin: true,
        secure: false,
        logLevel: 'debug'
      }
    }
  },
  css: {
    preprocessorOptions: {
      scss: {
        api: 'modern-compiler',
        silenceDeprecations: ['import']
      }
    }
  },
  // 优化依赖预构建
  optimizeDeps: {
    include: [
      'vue',
      'vue-router', 
      'pinia',
      'element-plus/es',
      'axios',
      'dayjs',
      'nprogress'
    ]
  }
})