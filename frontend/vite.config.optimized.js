import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
import { visualizer } from 'rollup-plugin-visualizer'

export default defineConfig({
  plugins: [
    vue(),
    AutoImport({
      resolvers: [ElementPlusResolver()],
      imports: ['vue', 'vue-router', 'pinia']
    }),
    Components({
      resolvers: [ElementPlusResolver()]
    }),
    // 打包分析插件（可选）
    visualizer({
      open: true,
      gzipSize: true,
      brotliSize: true
    })
  ],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src')
    }
  },
  build: {
    // 生产环境构建优化
    target: 'es2015',
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: true,
        drop_debugger: true,
        pure_funcs: ['console.log', 'console.info', 'console.debug']
      },
      format: {
        comments: false
      }
    },
    // 关闭源码映射以减小体积
    sourcemap: false,
    // 代码分割优化
    rollupOptions: {
      // 配置外部化依赖（使用CDN）
      external: [
        'vue',
        'element-plus',
        'echarts',
        'dayjs'
      ],
      output: {
        // 配置全局变量名称（对应CDN引入的全局变量）
        globals: {
          vue: 'Vue',
          'element-plus': 'ElementPlus',
          echarts: 'echarts',
          dayjs: 'dayjs'
        },
        // 手动配置代码分割
        manualChunks(id) {
          // 将node_modules中的代码单独打包
          if (id.includes('node_modules')) {
            // 第三方库按需分组
            if (id.includes('@element-plus/icons-vue')) {
              return 'element-icons';
            }
            if (id.includes('axios')) {
              return 'axios';
            }
            if (id.includes('pinia') || id.includes('vue-router')) {
              return 'vue-ecosystem';
            }
            if (id.includes('xlsx')) {
              return 'xlsx';
            }
            // 其他第三方库
            return 'vendor';
          }
          // 将大型组件单独分包
          if (id.includes('src/views/statistics')) {
            return 'statistics';
          }
          if (id.includes('src/views/orders')) {
            return 'orders';
          }
          if (id.includes('src/components/Enhanced')) {
            return 'enhanced-components';
          }
        },
        // 控制chunk命名
        chunkFileNames: (chunkInfo) => {
          const facadeModuleId = chunkInfo.facadeModuleId 
            ? chunkInfo.facadeModuleId.split('/').pop().split('.')[0] 
            : 'chunk';
          return `js/${facadeModuleId}-[hash].js`;
        },
        entryFileNames: 'js/[name]-[hash].js',
        assetFileNames: (assetInfo) => {
          // 分类资源文件
          const info = assetInfo.name.split('.');
          const ext = info[info.length - 1];
          if (/png|jpe?g|svg|gif|tiff|bmp|ico/i.test(ext)) {
            return `images/[name]-[hash][extname]`;
          } else if (/woff|woff2|eot|ttf|otf/i.test(ext)) {
            return `fonts/[name]-[hash][extname]`;
          } else if (ext === 'css') {
            return `css/[name]-[hash][extname]`;
          }
          return `assets/[name]-[hash][extname]`;
        }
      }
    },
    // CSS代码分割
    cssCodeSplit: true,
    // 启用 gzip 压缩提示
    reportCompressedSize: false,
    // 设置 chunk 大小警告限制
    chunkSizeWarningLimit: 500,
    // 设置资源内联阈值（小于4KB的资源会被内联）
    assetsInlineLimit: 4096
  },
  // 优化依赖预构建
  optimizeDeps: {
    include: [
      'vue-router',
      'pinia',
      '@element-plus/icons-vue',
      'axios',
      'nprogress',
      'vue-countup-v3'
    ],
    // 排除已经通过CDN加载的依赖
    exclude: [
      'vue',
      'element-plus',
      'echarts',
      'dayjs'
    ]
  },
  server: {
    host: '0.0.0.0',
    port: 3000,
    // 开启压缩
    compress: true,
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
    },
    // 提取CSS到单独的文件
    extract: true,
    // 开启CSS源映射（开发环境）
    sourceMap: false
  },
  // 定义全局常量
  define: {
    // 在生产环境中启用Vue的生产模式
    __VUE_PROD_DEVTOOLS__: false
  }
})