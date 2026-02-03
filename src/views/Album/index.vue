<template>
  <div class="album-detail">
    <h2>Album Detail (Placeholder)</h2>
    <p>Album ID: {{ $route.query.id }}</p>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router';
const route = useRoute();
</script>

<style scoped>
.album-detail {
  padding: 20px;
  color: var(--n-text-color);
}
</style>