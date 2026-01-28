<template>
  <n-result
    class="error"
    status="404"
    title="404 资源不存在"
    description="生活总归带点荒谬"
  >
    <template #footer>
      <n-button type="primary" @click="goHome">回到首页</n-button>
    </template>
  </n-result>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router";
const router = useRouter();

const goHome = () => {
  router.push("/");
};
</script>

<style lang="scss" scoped>
.error {
  height: calc(100vh - 208px);
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
}
</style>