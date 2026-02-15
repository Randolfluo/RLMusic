<template>
  <div class="stats-container">
    <div class="header">
      <h2 class="section-title">
        <n-icon :component="ChartGraph" class="title-icon" />
        系统概览
      </h2>
      <n-button quaternary circle type="primary" @click="getStats" :loading="loading" class="refresh-btn">
        <template #icon>
          <n-icon :component="Refresh" />
        </template>
      </n-button>
    </div>
    
    <div class="stats-grid">
      <!-- 资源统计 -->
      <div class="stats-card glass-card blue-theme">
        <div class="card-icon">
          <n-icon :component="Music" />
        </div>
        <div class="card-content">
          <n-statistic label="歌曲总数">
            <n-number-animation :from="0" :to="stats.song_count" />
            <template #suffix>首</template>
          </n-statistic>
        </div>
      </div>

      <div class="stats-card glass-card purple-theme">
        <div class="card-icon">
          <n-icon :component="RecordDisc" />
        </div>
        <div class="card-content">
          <n-statistic label="专辑总数">
            <n-number-animation :from="0" :to="stats.album_count" />
            <template #suffix>张</template>
          </n-statistic>
        </div>
      </div>

      <div class="stats-card glass-card pink-theme">
        <div class="card-icon">
          <n-icon :component="People" />
        </div>
        <div class="card-content">
          <n-statistic label="艺术家">
            <n-number-animation :from="0" :to="stats.artist_count" />
            <template #suffix>位</template>
          </n-statistic>
        </div>
      </div>

      <div class="stats-card glass-card orange-theme">
        <div class="card-icon">
          <n-icon :component="MusicList" />
        </div>
        <div class="card-content">
          <n-statistic label="歌单总数">
            <n-number-animation :from="0" :to="stats.playlist_count" />
            <template #suffix>个</template>
          </n-statistic>
        </div>
      </div>

      <!-- 时间与用户 -->
      <div class="stats-card glass-card cyan-theme interactive" @click="toggleUnit('uptime')">
        <div class="card-icon">
          <n-icon :component="Time" />
        </div>
        <div class="card-content">
          <n-statistic label="运行时长">
            <n-number-animation :from="0" :to="calcTime(stats.system_uptime, units.uptime)" />
            <template #suffix>{{ getUnitText(units.uptime) }}</template>
          </n-statistic>
        </div>
      </div>

      <div class="stats-card glass-card teal-theme interactive" @click="toggleUnit('music')">
        <div class="card-icon">
          <n-icon :component="Customer" />
        </div>
        <div class="card-content">
          <n-statistic label="歌曲时长">
            <n-number-animation :from="0" :to="calcTime(stats.music_duration, units.music)" />
            <template #suffix>{{ getUnitText(units.music) }}</template>
          </n-statistic>
        </div>
      </div>

      <div class="stats-card glass-card green-theme interactive" @click="toggleUnit('listen')">
        <div class="card-icon">
          <n-icon :component="Headset" />
        </div>
        <div class="card-content">
          <n-statistic label="听歌时长">
            <n-number-animation :from="0" :to="stats.user_listening_duration ? calcTime(stats.user_listening_duration, units.listen) : 0" />
            <template #suffix>{{ getUnitText(units.listen) }}</template>
          </n-statistic>
        </div>
      </div>

      <div class="stats-card glass-card indigo-theme">
        <div class="card-icon">
          <n-icon :component="User" />
        </div>
        <div class="card-content">
          <n-statistic label="用户总数">
            <n-number-animation :from="0" :to="stats.user_count" />
            <template #suffix>人</template>
          </n-statistic>
        </div>
      </div>

      <!-- 系统性能 -->
      <div class="stats-card glass-card red-theme">
        <div class="card-icon">
          <n-icon :component="Cpu" />
        </div>
        <div class="card-content">
          <n-statistic label="CPU">
            <n-number-animation :from="0" :to="stats.cpu_usage || 0" :precision="1" />
            <template #suffix>%</template>
          </n-statistic>
        </div>
      </div>

      <div class="stats-card glass-card yellow-theme">
        <div class="card-icon">
          <n-icon :component="MicroSd" />
        </div>
        <div class="card-content">
          <n-statistic label="内存">
            <n-number-animation :from="0" :to="stats.mem_usage || 0" :precision="1" />
            <template #suffix>%</template>
          </n-statistic>
        </div>
      </div>

      <div class="stats-card glass-card gray-theme span-2-desktop">
        <div class="card-icon">
          <n-icon :component="Api" />
        </div>
        <div class="card-content">
          <n-statistic label="API 调用">
            <n-number-animation :from="0" :to="stats.api_call_count || 0" />
            <template #suffix>次</template>
          </n-statistic>
        </div>
      </div>
    </div>
    <n-divider class="divider" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from "vue";
import { getSystemStats, type SystemStats } from "@/api/system";
import { ResultCode } from "@/utils/request";
import { 
  Refresh, Music, RecordDisc, People, MusicList, 
  Time, Customer, Headset, User, Cpu, MicroSd, Api, ChartGraph 
} from "@icon-park/vue-next";

const loading = ref(false);
const stats = ref<SystemStats>({
  song_count: 0,
  album_count: 0,
  artist_count: 0,
  music_duration: 0,
  playlist_count: 0,
  user_count: 0,
  system_uptime: 0,
  user_listening_duration: 0,
  user_scanned_duration: 0,
  cpu_usage: 0,
  mem_usage: 0,
  api_call_count: 0
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

const getStats = async () => {
    loading.value = true;
    try {
        const res = await getSystemStats();
        if (res.code === ResultCode.SUCCESS) {
            stats.value = res.data;
        }
    } catch (e) {
        console.error(e);
    } finally {
        loading.value = false;
    }
}

onMounted(() => {
  getStats();
});
</script>

<style scoped lang="scss">
.stats-container {
  margin-bottom: 24px;
  animation: fadeIn 0.5s ease-out;

  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;

    .section-title {
      font-size: 18px;
      font-weight: 700;
      color: var(--n-text-color);
      display: flex;
      align-items: center;
      gap: 8px;
      margin: 0;

      .title-icon {
        color: var(--n-primary-color);
      }
    }

    .refresh-btn {
      transition: transform 0.3s ease;
      &:hover {
        transform: rotate(180deg);
      }
    }
  }

  .stats-grid {
    display: grid;
    grid-template-columns: repeat(1, 1fr);
    gap: 16px;
    width: 100%;
    
    @media (min-width: 640px) {
      grid-template-columns: repeat(2, 1fr);
    }
    
    @media (min-width: 1024px) {
      grid-template-columns: repeat(4, 1fr);
      
      .span-2-desktop {
        grid-column: span 2;
      }
    }
  }

  .stats-card {
    display: flex;
    align-items: center;
    padding: 16px;
    border-radius: 16px;
    background: var(--n-card-color);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
    transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
    min-height: 100px;
    border: 1px solid rgba(255, 255, 255, 0.1);
    box-sizing: border-box;

    &:hover {
      transform: translateY(-4px);
      box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
      z-index: 1;
    }

    &.interactive {
      cursor: pointer;
      &:active {
        transform: scale(0.98);
      }
    }

    .card-icon {
      width: 48px;
      height: 48px;
      border-radius: 12px;
      display: flex;
      justify-content: center;
      align-items: center;
      font-size: 24px;
      margin-right: 16px;
      flex-shrink: 0;
    }

    .card-content {
      flex: 1;
      min-width: 0;

      :deep(.n-statistic) {
        .n-statistic__label {
          font-size: 12px;
          color: var(--n-text-color-3);
          margin-bottom: 4px;
        }
        .n-statistic-value__content {
          font-weight: 800;
          font-size: 20px;
          font-family: 'SF Pro Display', -apple-system, BlinkMacSystemFont, Roboto, sans-serif;
        }
        .n-statistic-value__suffix {
          font-size: 12px;
          margin-left: 4px;
          color: var(--n-text-color-3);
        }
      }
    }

    // Color Themes
    &.blue-theme { .card-icon { background: rgba(33, 150, 243, 0.1); color: #2196f3; } }
    &.purple-theme { .card-icon { background: rgba(156, 39, 176, 0.1); color: #9c27b0; } }
    &.pink-theme { .card-icon { background: rgba(233, 30, 99, 0.1); color: #e91e63; } }
    &.orange-theme { .card-icon { background: rgba(255, 152, 0, 0.1); color: #ff9800; } }
    &.cyan-theme { .card-icon { background: rgba(0, 188, 212, 0.1); color: #00bcd4; } }
    &.teal-theme { .card-icon { background: rgba(0, 150, 136, 0.1); color: #009688; } }
    &.green-theme { .card-icon { background: rgba(76, 175, 80, 0.1); color: #4caf50; } }
    &.indigo-theme { .card-icon { background: rgba(63, 81, 181, 0.1); color: #3f51b5; } }
    &.red-theme { .card-icon { background: rgba(244, 67, 54, 0.1); color: #f44336; } }
    &.yellow-theme { .card-icon { background: rgba(255, 235, 59, 0.1); color: #fbc02d; } }
    &.gray-theme { .card-icon { background: rgba(158, 158, 158, 0.1); color: #9e9e9e; } }
  }

  .divider {
    margin-top: 24px;
    opacity: 0.6;
  }
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
