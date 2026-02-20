<template>
  <div class="artists">
    <span 
      v-for="(item, index) in artistsData" 
      :key="item ? item.id : index" 
      class="artist"
    >
      <template v-if="item">
        <span class="name" @click.stop="toArtist(item.id)">
            {{ item.name }}
        </span>
        <span 
          class="split" 
          v-if="index < artistsData.length - 1 && artistsData[index + 1]"
        > / </span>
      </template>
    </span>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router";

const router = useRouter();

const props = defineProps<{
  artistsData: { id: number; name: string }[]
}>();

const toArtist = (id: number) => {
  if (id) {
    router.push({ path: "/artist", query: { id } });
  }
};
</script>

<style lang="scss" scoped>
.artists {
  display: inline-flex;
  align-items: center;
  .artist {
    display: inline-flex;
    align-items: center;
    .name {
      cursor: pointer;
      &:hover {
        color: var(--n-color-primary);
      }
    }
    .split {
      margin: 0 4px;
      opacity: 0.6;
    }
  }
}
</style>