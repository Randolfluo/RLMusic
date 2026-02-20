<template>
  <div class="system-stats-view">
    <div class="container">
      <h1>系统信息</h1>
      <SystemStats />
    </div>
  </div>
</template>

<script setup lang="ts">
import SystemStats from "@/components/Home/SystemStats.vue";
</script>

<style lang="scss" scoped>
.system-stats-view {
  width: 100%;
  min-height: calc(100vh - 140px); // Adjust based on header/footer
  padding: 40px 20px;
  display: flex;
  justify-content: center;

  @media (max-width: 768px) {
    padding: 20px 16px;
  }

  .container {
    width: 100%;
    max-width: 1200px;
    
    h1 {
      font-size: 32px;
      font-weight: bold;
      margin-bottom: 30px;
      color: var(--n-text-color);

      @media (max-width: 768px) {
        font-size: 24px;
        margin-bottom: 20px;
      }
    }
  }
}
</style>
