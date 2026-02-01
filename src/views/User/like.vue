<template>
  <div class="user-like">
    <n-card title="我的上传 / 喜欢">
      <n-empty description="暂无数据" />
      <div style="margin-top: 20px; text-align: center;">
         <n-button @click="$router.back()">返回</n-button>
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
</script>

<style scoped>
.user-like {
  padding: 40px 20px;
  max-width: 800px;
  margin: 0 auto;
}
</style>
