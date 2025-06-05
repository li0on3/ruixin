import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'

// 快速优化版配置 - 无需额外安装插件
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
    // 关闭源码映射
    sourcemap: false,
    // 优化的代码分割
    rollupOptions: {
      output: {
        manualChunks(id) {
          // 将 node_modules 的代码分离
          if (id.includes('node_modules')) {
            // Vue 相关
            if (id.includes('vue') || id.includes('pinia') || id.includes('vue-router')) {
              return 'vue-core';
            }
            // Element Plus
            if (id.includes('element-plus')) {
              return 'element-ui';
            }
            // ECharts
            if (id.includes('echarts')) {
              return 'charts';
            }
            // 其他依赖
            if (id.includes('axios') || id.includes('dayjs') || id.includes('nprogress')) {
              return 'utils';
            }
            // xlsx 单独打包
            if (id.includes('xlsx')) {
              return 'xlsx';
            }
            return 'vendor';
          }
        },
        // 更好的chunk命名
        chunkFileNames: 'js/[name]-[hash].js',
        entryFileNames: 'js/[name]-[hash].js',
        assetFileNames: 'assets/[name]-[hash].[ext]'
      }
    },
    // 降低chunk大小警告限制
    chunkSizeWarningLimit: 1500,
    // 提高构建性能
    reportCompressedSize: false
  },
  server: {
    host: '0.0.0.0',
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:8080',
        changeOrigin: true,
        secure: false
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
      'element-plus',
      'axios',
      'dayjs',
      'nprogress',
      'echarts'
    ]
  }
})