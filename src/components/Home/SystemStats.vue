<template>
  <div class="stats-container">
    <n-grid x-gap="12" y-gap="12" :cols="4">
      <n-gi>
        <n-card size="small" :bordered="false" class="stats-card">
          <n-statistic label="歌曲总数">
            <n-number-animation :from="0" :to="stats.song_count" />
            <template #suffix>首</template>
          </n-statistic>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card size="small" :bordered="false" class="stats-card">
          <n-statistic label="专辑总数">
            <n-number-animation :from="0" :to="stats.album_count" />
             <template #suffix>张</template>
          </n-statistic>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card size="small" :bordered="false" class="stats-card">
          <n-statistic label="艺术家总数">
            <n-number-animation :from="0" :to="stats.artist_count" />
             <template #suffix>位</template>
          </n-statistic>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card size="small" :bordered="false" class="stats-card">
          <n-statistic label="歌单总数">
            <n-number-animation :from="0" :to="stats.playlist_count" />
             <template #suffix>个</template>
          </n-statistic>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card size="small" :bordered="false" class="stats-card">
          <n-statistic label="用户总数">
            <n-number-animation :from="0" :to="stats.user_count" />
             <template #suffix>人</template>
          </n-statistic>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card size="small" :bordered="false" class="stats-card interactive" @click="toggleUnit('uptime')">
          <n-statistic label="系统运行时长">
            <n-number-animation :from="0" :to="calcTime(stats.system_uptime, units.uptime)" />
             <template #suffix>{{ getUnitText(units.uptime) }}</template>
          </n-statistic>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card size="small" :bordered="false" class="stats-card interactive" @click="toggleUnit('music')">
          <n-statistic label="歌曲总时长">
             <n-number-animation :from="0" :to="calcTime(stats.music_duration, units.music)" />
             <template #suffix>{{ getUnitText(units.music) }}</template>
          </n-statistic>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card size="small" :bordered="false" class="stats-card interactive" @click="toggleUnit('listen')">
          <n-statistic label="我的听歌时长">
             <n-number-animation :from="0" :to="calcTime(stats.user_listening_duration, units.listen)" />
             <template #suffix>{{ getUnitText(units.listen) }}</template>
          </n-statistic>
        </n-card>
      </n-gi>
    </n-grid>
    <n-divider />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from "vue";
import { getSystemStats, type SystemStats } from "@/api/system";
import { ResultCode } from "@/utils/request";

const stats = ref<SystemStats>({
  song_count: 0,
  album_count: 0,
  artist_count: 0,
  music_duration: 0,
  playlist_count: 0,
  user_count: 0,
  system_uptime: 0,
  user_listening_duration: 0,
  user_scanned_duration: 0
});

const units = reactive({
  uptime: 'h',
  music: 'm',
  listen: 'm'
});

const toggleUnit = (key: keyof typeof units) => {
  const next: Record<string, string> = { 's': 'm', 'm': 'h', 'h': 's' };
  const val = next[units[key]];
  if (val) units[key] = val;
};

const calcTime = (seconds: number, unit: string) => {
  if (unit === 's') return seconds;
  if (unit === 'm') return Math.floor(seconds / 60);
  if (unit === 'h') return Math.floor(seconds / 3600);
  return 0;
};

const getUnitText = (unit: string) => {
    const map: Record<string, string> = { 's': '秒', 'm': '分钟', 'h': '小时' };
    return map[unit] || unit;
}

onMounted(() => {
  getStats();
});

const getStats = async () => {
    try {
        const res = await getSystemStats();
        if (res.code === ResultCode.SUCCESS) {
            stats.value = res.data;
        }
    } catch (e) {
        console.error(e);
    }
}
</script>

<style scoped lang="scss">
.stats-container {
    margin-bottom: 10px;
    .stats-card {
        background-color: rgba(255, 255, 255, 0.5); // subtle background for light mode
        transition: transform 0.2s, box-shadow 0.2s;
        
        &.interactive {
            cursor: pointer;
            &:active {
                transform: scale(0.98);
            }
            &:hover {
                 background-color: rgba(255, 255, 255, 0.8);
            }
        }

        :deep(.n-statistic .n-statistic-value__content) {
            font-weight: bold;
            color: var(--n-color-primary);
        }
    }
}
</style>
