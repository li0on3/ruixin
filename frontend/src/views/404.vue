<template>
  <div class="not-found">
    <div class="not-found-content">
      <h1>404</h1>
      <p>抱歉，您访问的页面不存在</p>
      <el-button type="primary" @click="$router.push('/')">
        返回首页
      </el-button>
    </div>
  </div>
</template>

<script setup>
</script>

<style lang="scss" scoped>
.not-found {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f7fa;
  
  .not-found-content {
    text-align: center;
    
    h1 {
      font-size: 120px;
      color: #909399;
      margin: 0;
      font-weight: bold;
    }
    
    p {
      font-size: 20px;
      color: #606266;
      margin: 20px 0 30px;
    }
  }
}
</style>