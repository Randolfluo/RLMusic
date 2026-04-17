<template>
  <div class="stats-container">
    <div class="header">
      <h2 class="section-title">
        <n-icon :component="ChartGraph" class="title-icon" />
        系统概览
      </h2>
      <n-button secondary circle type="primary" @click="getStats" :loading="loading" class="refresh-btn">
        <template #icon>
          <n-icon :component="Refresh" />
        </template>
      </n-button>
    </div>

    <div class="stats-grid">
      <!-- 资源统计 -->
      <div class="stats-card glass-card coral-theme">
        <div class="card-icon">
          <n-icon :component="Music" />
        </div>
        <div class="card-content">
          <n-statistic label="歌曲总数">
            <n-number-animation :from="prevStats.song_count" :to="stats.song_count" />
            <template #suffix>首</template>
          </n-statistic>
        </div>
      </div>

      <div class="stats-card glass-card teal-theme">
        <div class="card-icon">
          <n-icon :component="RecordDisc" />
        </div>
        <div class="card-content">
          <n-statistic label="专辑总数">
            <n-number-animation :from="prevStats.album_count" :to="stats.album_count" />
            <template #suffix>张</template>
          </n-statistic>
        </div>
      </div>

      <div class="stats-card glass-card gold-theme">
        <div class="card-icon">
          <n-icon :component="People" />
        </div>
        <div class="card-content">
          <n-statistic label="艺术家">
            <n-number-animation :from="prevStats.artist_count" :to="stats.artist_count" />
            <template #suffix>位</template>
          </n-statistic>
        </div>
      </div>

      <div class="stats-card glass-card purple-theme">
        <div class="card-icon">
          <n-icon :component="MusicList" />
        </div>
        <div class="card-content">
          <n-statistic label="歌单总数">
            <n-number-animation :from="prevStats.playlist_count" :to="stats.playlist_count" />
            <template #suffix>个</template>
          </n-statistic>
        </div>
      </div>

      <!-- 时间与用户 -->
      <div class="stats-card glass-card blue-theme interactive" @click="toggleUnit('uptime')">
        <div class="card-icon">
          <n-icon :component="Time" />
        </div>
        <div class="card-content">
          <n-statistic label="运行时长">
            <n-number-animation :from="calcTime(prevStats.system_uptime, units.uptime)" :to="calcTime(stats.system_uptime, units.uptime)" />
            <template #suffix>{{ getUnitText(units.uptime) }}</template>
          </n-statistic>
        </div>
      </div>

      <div class="stats-card glass-card green-theme interactive" @click="toggleUnit('music')">
        <div class="card-icon">
          <n-icon :component="Customer" />
        </div>
        <div class="card-content">
          <n-statistic label="歌曲时长">
            <n-number-animation :from="calcTime(prevStats.music_duration, units.music)" :to="calcTime(stats.music_duration, units.music)" />
            <template #suffix>{{ getUnitText(units.music) }}</template>
          </n-statistic>
        </div>
      </div>

      <div class="stats-card glass-card ink-theme interactive" @click="toggleUnit('listen')">
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

      <div class="stats-card glass-card user-theme">
        <div class="card-icon">
          <n-icon :component="User" />
        </div>
        <div class="card-content">
          <n-statistic label="用户总数">
            <n-number-animation :from="prevStats.user_count" :to="stats.user_count" />
            <template #suffix>人</template>
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
  Time, Customer, Headset, User, ChartGraph
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
  user_listening_duration: 0
});

const prevStats = ref<SystemStats>({
  song_count: 0,
  album_count: 0,
  artist_count: 0,
  music_duration: 0,
  playlist_count: 0,
  user_count: 0,
  system_uptime: 0,
  user_listening_duration: 0
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
            prevStats.value = { ...stats.value };
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
@import url('https://fonts.googleapis.com/css2?family=Playfair+Display:wght@600;700&family=Plus+Jakarta+Sans:wght@400;500;600;700&display=swap');

.stats-container {
  margin-bottom: 24px;
  animation: fadeIn 0.5s ease-out;
  font-family: 'Plus Jakarta Sans', sans-serif;

  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;

    .section-title {
      font-size: 18px;
      font-weight: 700;
      color: #1a1a1a;
      display: flex;
      align-items: center;
      gap: 8px;
      margin: 0;
      font-family: 'Plus Jakarta Sans', sans-serif;

      .title-icon {
        color: #e07a5f;
      }
    }

    .refresh-btn {
      transition: transform 0.3s ease;
      background: #f5f2ed !important;
      border-color: #ebe7e0 !important;
      color: #666666 !important;

      &:hover {
        transform: rotate(180deg);
        background: #ebe7e0 !important;
      }
    }
  }

  .stats-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 10px;
    width: 100%;

    @media (min-width: 640px) {
      gap: 16px;
    }

    @media (max-width: 768px) {
      grid-template-columns: 1fr;
      gap: 12px;
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
    padding: 20px;
    border-radius: 20px;
    background: rgba(255, 255, 255, 0.6);
    backdrop-filter: blur(16px) saturate(180%);
    border: 1px solid rgba(255, 255, 255, 0.5);
    box-shadow:
      0 4px 20px rgba(0, 0, 0, 0.04),
      inset 0 0 0 1px rgba(255, 255, 255, 0.5);
    transition: all 0.4s cubic-bezier(0.25, 0.8, 0.25, 1);
    min-height: 100px;
    box-sizing: border-box;

    &:hover {
      transform: translateY(-6px) scale(1.02);
      box-shadow:
        0 8px 30px rgba(0, 0, 0, 0.08),
        inset 0 0 0 1px rgba(255, 255, 255, 0.6);
      z-index: 1;
    }

    &.interactive {
      cursor: pointer;
      &:active {
        transform: scale(0.98);
      }
    }

    .card-icon {
      width: 52px;
      height: 52px;
      border-radius: 16px;
      display: flex;
      justify-content: center;
      align-items: center;
      font-size: 26px;
      margin-right: 16px;
      flex-shrink: 0;
      transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
      box-shadow: 0 4px 12px rgba(0,0,0,0.1);
    }

    &:hover .card-icon {
      transform: scale(1.1) rotate(5deg);
    }

    .card-content {
      flex: 1;
      min-width: 0;

      :deep(.n-statistic) {
        .n-statistic__label {
          font-size: 12px;
          color: #666666;
          margin-bottom: 4px;
        }
        .n-statistic-value__content {
          font-weight: 700;
          font-size: 20px;
          color: #1a1a1a;
          font-family: 'Playfair Display', serif;
        }
        .n-statistic-value__suffix {
          font-size: 12px;
          margin-left: 4px;
          color: #999999;
        }
      }
    }

    @media (max-width: 640px) {
      padding: 12px;
      border-radius: 14px;
      min-height: 72px;

      .card-icon {
        width: 36px;
        height: 36px;
        border-radius: 10px;
        font-size: 18px;
        margin-right: 10px;
      }

      .card-content {
        :deep(.n-statistic) {
          .n-statistic__label {
            font-size: 11px;
            margin-bottom: 2px;
          }
          .n-statistic-value__content {
            font-size: 16px;
          }
          .n-statistic-value__suffix {
            font-size: 10px;
            margin-left: 2px;
          }
        }

        .statistic-text {
          font-weight: 700;
          font-size: 16px;
          color: #1a1a1a;
          font-family: 'Playfair Display', serif;
        }
      }
    }

    // Color Themes - 温暖米色调配色
    &.coral-theme { .card-icon { background: linear-gradient(135deg, #e07a5f 0%, #d4a574 100%); color: white; } }
    &.teal-theme { .card-icon { background: linear-gradient(135deg, #3d8b8b 0%, #5b8db8 100%); color: white; } }
    &.gold-theme { .card-icon { background: linear-gradient(135deg, #d4a574 0%, #e07a5f 100%); color: white; } }
    &.purple-theme { .card-icon { background: linear-gradient(135deg, #7c6fae 0%, #a89bc9 100%); color: white; } }
    &.blue-theme { .card-icon { background: linear-gradient(135deg, #5b8db8 0%, #3d8b8b 100%); color: white; } }
    &.green-theme { .card-icon { background: linear-gradient(135deg, #11998e 0%, #38d39f 100%); color: white; } }
    &.ink-theme { .card-icon { background: linear-gradient(135deg, #2c3e50 0%, #3d8b8b 100%); color: white; } }
    &.user-theme { .card-icon { background: linear-gradient(135deg, #3d8b8b 0%, #7c6fae 100%); color: white; } }
  }

  .divider {
    margin-top: 24px;
    opacity: 0.6;
    border-color: #ebe7e0;
  }
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Dark Mode Support */
:global(.dark) {
  .stats-container {
    .section-title {
      color: #ffffff;
    }

    .stats-card {
      background: rgba(40, 40, 40, 0.6);
      border-color: rgba(255, 255, 255, 0.1);
      box-shadow:
        0 4px 20px rgba(0, 0, 0, 0.2),
        inset 0 0 0 1px rgba(255, 255, 255, 0.05);

      &:hover {
        background: rgba(50, 50, 50, 0.7);
        box-shadow:
          0 8px 30px rgba(0, 0, 0, 0.3),
          inset 0 0 0 1px rgba(255, 255, 255, 0.08);
      }
    }

    :deep(.n-statistic-value__content) {
      color: #ffffff !important;
    }

    .statistic-text {
      color: #ffffff !important;
    }
  }
}
</style>
