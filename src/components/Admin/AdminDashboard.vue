<template>
  <div class="admin-dashboard-container">
    <!-- Ambient Background Effects -->
    <div class="ambient-glows">
      <div class="glow-orb orb-1"></div>
      <div class="glow-orb orb-2"></div>
      <div class="glow-orb orb-3"></div>
    </div>
    <div class="noise-overlay"></div>

    <div class="dashboard-content">
      <!-- Header Section -->
      <div class="header-section">
        <div class="header-content">
          <div class="title-wrapper">
            <div class="icon-box">
              <n-icon :component="Permissions" />
            </div>
            <div>
              <h1 class="page-title">管理控制台</h1>
              <p class="page-subtitle">系统维护与智能中心</p>
            </div>
          </div>
          <div class="header-actions">
            <div class="status-badge" v-if="stats">
              <span class="dot"></span>
              系统运行中
            </div>
          </div>
        </div>
      </div>

      <!-- System Monitor Section (Gray) -->
      <div class="dashboard-section" v-if="realtimeStats">
        <div class="section-label">实时监控</div>
        <div class="realtime-monitor-grid">
          <!-- CPU Card -->
          <div class="monitor-card glass-panel span-2 gray-theme">
            <div class="monitor-main">
              <div class="monitor-icon gray">
                <n-icon :component="Cpu" />
              </div>
              <div class="monitor-info">
                <div class="monitor-label">CPU 使用率</div>
                <div class="progress-container">
                  <n-progress
                    type="circle"
                    :percentage="realtimeStats?.cpu_usage || 0"
                    :color="'#64748b'"
                    :rail-color="'rgba(100, 116, 139, 0.1)'"
                    :indicator-placement="'inside'"
                    :stroke-width="10"
                    style="width: 50px; height: 50px;"
                  >
                    <span class="progress-text">
                        <n-number-animation :from="prevCpuUsage" :to="realtimeStats?.cpu_usage || 0" :precision="0" />%
                    </span>
                  </n-progress>
                </div>
              </div>
            </div>
            <div class="monitor-chart">
              <v-chart class="mini-chart" :option="cpuChartOption" autoresize />
            </div>
          </div>

          <!-- Memory Card -->
          <div class="monitor-card glass-panel span-2 gray-theme">
            <div class="monitor-main">
              <div class="monitor-icon gray">
                <n-icon :component="HardDisk" />
              </div>
              <div class="monitor-info">
                <div class="monitor-label">内存占用</div>
                <div class="progress-container">
                  <n-progress
                    type="circle"
                    :percentage="realtimeStats?.mem_usage || 0"
                    :color="'#64748b'"
                    :rail-color="'rgba(100, 116, 139, 0.1)'"
                    :indicator-placement="'inside'"
                    :stroke-width="10"
                    style="width: 50px; height: 50px;"
                  >
                    <span class="progress-text">
                        <n-number-animation :from="prevMemUsage" :to="realtimeStats?.mem_usage || 0" :precision="0" />%
                    </span>
                  </n-progress>
                </div>
              </div>
            </div>
            <div class="monitor-chart">
              <v-chart class="mini-chart" :option="memChartOption" autoresize />
            </div>
          </div>

          <!-- API Call Card -->
          <div class="monitor-card glass-panel gray-theme">
            <div class="monitor-icon gray">
              <n-icon :component="Api" />
            </div>
            <div class="monitor-content">
              <div class="monitor-label">API 调用</div>
              <div class="monitor-value">
                <n-number-animation :from="prevApiCallCount" :to="realtimeStats?.api_call_count || 0" />
              </div>
            </div>
          </div>

          <!-- QPS Card -->
          <div class="monitor-card glass-panel gray-theme">
            <div class="monitor-icon gray">
              <n-icon :component="TrendingUp" />
            </div>
            <div class="monitor-content">
              <div class="monitor-label">QPS</div>
              <div class="monitor-value">
                <n-number-animation :from="0" :to="qps" />
              </div>
            </div>
          </div>

          <!-- Go Routines Card -->
          <div class="monitor-card glass-panel gray-theme">
            <div class="monitor-icon gray">
              <n-icon :component="Connection" />
            </div>
            <div class="monitor-content">
              <div class="monitor-label">Go 协程</div>
              <div class="monitor-value">
                <n-number-animation :from="prevGoRoutines" :to="realtimeStats?.go_routines || 0" />
              </div>
            </div>
          </div>

          <!-- Total Volume Card -->
          <div class="monitor-card glass-panel gray-theme">
            <div class="monitor-icon gray">
              <n-icon :component="CloudStorage" />
            </div>
            <div class="monitor-content">
              <div class="monitor-label">歌曲占用空间</div>
              <div class="monitor-value text-sm">
                {{ formatBytes(realtimeStats?.total_volume || 0) }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Resource Overview Section (Purple) -->
      <div class="dashboard-section" v-if="stats">
        <div class="section-label">资源概览</div>
        <div class="stats-grid">
          <!-- Music Count Card -->
          <div class="stat-card primary-stat purple-theme">
            <div class="stat-content">
              <div class="stat-label">歌曲总数</div>
              <div class="stat-value">{{ stats.song_count.toLocaleString() }}</div>
              <div class="stat-trend">
                <n-icon :component="TrendingUp" />
                <span>曲库增长中</span>
              </div>
            </div>
            <div class="stat-chart-decoration">
              <svg viewBox="0 0 100 40" preserveAspectRatio="none">
                <path d="M0,20 Q25,5 50,20 T100,20 V40 H0 Z" fill="rgba(255,255,255,0.2)" />
              </svg>
            </div>
            <div class="stat-icon-bg">
              <n-icon :component="Music" />
            </div>
          </div>

          <!-- User Card -->
          <div class="stat-card secondary-stat purple-theme" @click="handleUserManage" style="cursor: pointer;">
            <div class="stat-header">
              <div class="stat-icon-small purple">
                <n-icon :component="User" />
              </div>
              <span class="stat-label">用户管理</span>
            </div>
            <div class="stat-value-small">{{ stats.user_count.toLocaleString() }}</div>
            <div class="stat-footer">活跃账号</div>
          </div>

          <!-- Playlist Card -->
          <div class="stat-card secondary-stat purple-theme" @click="handlePlaylistManage" style="cursor: pointer;">
            <div class="stat-header">
              <div class="stat-icon-small purple">
                <n-icon :component="CheckOne" />
              </div>
              <span class="stat-label">公共歌单管理</span>
            </div>
            <div class="stat-value-small">{{ stats.playlist_count.toLocaleString() }}</div>
            <div class="stat-footer">公开收藏集</div>
          </div>
          
          <!-- Artist Card -->
          <div class="stat-card secondary-stat purple-theme">
              <div class="stat-header">
                <div class="stat-icon-small purple">
                  <n-icon :component="People" />
                </div>
                <span class="stat-label">艺术家</span>
              </div>
              <div class="stat-value-small">{{ stats.artist_count.toLocaleString() }}</div>
              <div class="stat-footer">入驻艺人</div>
            </div>
        </div>
        
        <!-- Tools Sub-Section integrated into Business Data -->
        <div class="tools-integration">
          <div class="ops-grid-compact">
              <!-- Scan Music -->
              <div class="op-card-compact glass-panel white-theme">
                <div class="compact-icon green">
                  <n-icon :component="Scan" />
                </div>
                <div class="compact-info">
                  <h3>曲库扫描</h3>
                  <p>同步本地文件与数据库</p>
                </div>
                <n-button class="action-circle-btn" circle secondary type="success" :loading="loading.scan" @click="handleScanMusic">
                  <template #icon><n-icon :component="Play" /></template>
                </n-button>
              </div>

              <!-- Export Data -->
              <div class="op-card-compact glass-panel white-theme">
                <div class="compact-icon teal">
                  <n-icon :component="FileExcel" />
                </div>
                <div class="compact-info">
                  <h3>数据导出</h3>
                  <p>导出数据库为 Excel</p>
                </div>
                <n-button class="action-circle-btn" circle secondary type="info" @click="handleExportExcel">
                  <template #icon><n-icon :component="Download" /></template>
                </n-button>
              </div>

              <!-- API Documentation -->
              <div class="op-card-compact glass-panel white-theme">
                <div class="compact-icon blue">
                  <n-icon :component="DocDetail" />
                </div>
                <div class="compact-info">
                  <h3>API 文档</h3>
                  <p>在线接口文档</p>
                </div>
                <n-button class="action-circle-btn" circle secondary type="info" @click="handleApiDoc">
                  <template #icon><n-icon :component="Link" /></template>
                </n-button>
              </div>
            </div>
        </div>
      </div>

      <!-- AI Feature Section (Gradient) -->
      <div class="dashboard-section ai-section">
        <div class="section-label">AI 智能套件</div>
        <div class="ops-grid">
          <!-- Playlist AI -->
          <div class="op-card glass-panel gradient-card" @click="handleGeneratePlaylistDesc">
            <div class="op-card-glow blue"></div>
            <div class="op-icon-large blue">
              <n-icon :component="MusicList" />
            </div>
            <div class="op-info">
              <h3>歌单洞察</h3>
              <p>为所有公共歌单生成 AI 描述</p>
            </div>
            <div class="op-action">
              <n-button class="glass-btn" size="small" :loading="loading.playlist">
                <template #icon><n-icon :component="Lightning" /></template>
                生成
              </n-button>
            </div>
          </div>

          <!-- Artist AI -->
          <div class="op-card glass-panel gradient-card" @click="handleGenerateArtistDesc">
            <div class="op-card-glow purple"></div>
            <div class="op-icon-large purple">
              <n-icon :component="People" />
            </div>
            <div class="op-info">
              <h3>艺术家画像</h3>
              <p>分析并生成所有艺术家的简介</p>
            </div>
            <div class="op-action">
              <n-button class="glass-btn" size="small" :loading="loading.artist">
                <template #icon><n-icon :component="Lightning" /></template>
                生成
              </n-button>
            </div>
          </div>

          <!-- Album AI -->
          <div class="op-card glass-panel gradient-card" @click="handleGenerateAlbumDesc">
            <div class="op-card-glow orange"></div>
            <div class="op-icon-large orange">
              <n-icon :component="RecordDisc" />
            </div>
            <div class="op-info">
              <h3>专辑纪事</h3>
              <p>为专辑系列创建丰富叙述</p>
            </div>
            <div class="op-action">
              <n-button class="glass-btn" size="small" :loading="loading.album">
                <template #icon><n-icon :component="Lightning" /></template>
                生成
              </n-button>
            </div>
          </div>

          <!-- Playlist Intros -->
          <div class="op-card glass-panel gradient-card" @click="handleGenerateAllIntros">
            <div class="op-card-glow teal"></div>
            <div class="op-icon-large teal">
              <n-icon :component="Voice" />
            </div>
            <div class="op-info">
              <h3>全站开场白</h3>
              <p>批量生成所有公共歌单的开场白</p>
            </div>
            <div class="op-action">
              <n-button class="glass-btn" size="small" :loading="loading.intro">
                <template #icon><n-icon :component="Lightning" /></template>
                生成
              </n-button>
            </div>
          </div>
        </div>
      </div>

      <!-- Danger Zone (Red) -->
      <div class="dashboard-section danger-section">
        <div class="section-label danger">危险区域</div>
        <div class="ops-grid-compact">
          <!-- Reset System -->
          <div class="op-card-compact glass-panel danger-zone">
            <div class="compact-icon red">
              <n-icon :component="Delete" />
            </div>
            <div class="compact-info">
              <h3>恢复出厂设置</h3>
              <p>清除所有数据（不可逆）</p>
            </div>
            <n-popconfirm @positive-click="handleResetSystem" negative-text="取消" positive-text="确认重置">
              <template #trigger>
                <n-button class="action-circle-btn" circle secondary type="error" :loading="loading.reset">
                  <template #icon><n-icon :component="Power" /></template>
                </n-button>
              </template>
              确定要重置整个系统吗？
            </n-popconfirm>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// Force update
import { ref, onMounted, onUnmounted, reactive, computed, provide, watch } from "vue";
import { useRouter } from "vue-router";
import { 
  Permissions, MusicList, People, RecordDisc, 
  Scan, FileExcel, Delete, Music, User, Time, CheckOne,
  TrendingUp, Lightning, Play, Download, Power,
  DocDetail, Link, Cpu, Api, HardDisk, CloudStorage, Connection, Voice
} from "@icon-park/vue-next";
import { useMessage, NIcon, NButton, NPopconfirm, NNumberAnimation, NProgress } from "naive-ui";
import { 
  scanMusic 
} from "@/api/song";
import {
  generatePublicPlaylistDescriptions,
  generateArtistDescriptions,
  generateAlbumDescriptions,
  generateAllPublicPlaylistIntros
} from "@/api/ai";
import { resetSystem, getSystemStats, getSystemStatus, type SystemStats, type SystemStatus } from "@/api/system";
import { ResultCode } from "@/utils/request";

// ECharts
import VChart, { THEME_KEY } from "vue-echarts";
import { use } from "echarts/core";
import { CanvasRenderer } from "echarts/renderers";
import { LineChart } from "echarts/charts";
import { GridComponent, TooltipComponent } from "echarts/components";

use([CanvasRenderer, LineChart, GridComponent, TooltipComponent]);

const router = useRouter();
const message = useMessage();

// Provide ECharts theme
provide(THEME_KEY, "light");

// State
const stats = ref<SystemStats | null>(null);
const realtimeStats = ref<SystemStatus | null>(null);
const prevApiCallCount = ref(0);
const prevSystemUptime = ref(0);
const prevCpuUsage = ref(0);
const prevMemUsage = ref(0);
const prevGoRoutines = ref(0);
const prevTotalVolume = ref(0);
const qps = ref(0);
const intervalId = ref<number | null>(null);

// History data for charts
const cpuHistory = ref<number[]>(new Array(300).fill(0));
const memHistory = ref<number[]>(new Array(300).fill(0));

const loading = reactive({
  playlist: false,
  artist: false,
  album: false,
  scan: false,
  intro: false,
  reset: false
});

const uptimeUnit = ref('auto'); // auto, s, m, h

// Chart Options
const getChartOption = (data: number[], color: string) => {
  return {
    animation: false,
    grid: { left: 0, right: 0, top: 5, bottom: 0 },
    xAxis: { type: 'category', show: false, boundaryGap: false },
    yAxis: { type: 'value', show: false, min: 0, max: 100 },
    series: [
      {
        data: data,
        type: 'line',
        smooth: true,
        showSymbol: false,
        lineStyle: { width: 2, color: color },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0, y: 0, x2: 0, y2: 1,
            colorStops: [{ offset: 0, color: color }, { offset: 1, color: 'rgba(255, 255, 255, 0)' }],
            global: false
          },
          opacity: 0.3
        }
      }
    ],
    tooltip: { 
      trigger: 'axis', 
      formatter: '{c}%',
      position: function (point: any) { return [point[0], '10%']; }
    }
  };
};

const cpuChartOption = computed(() => getChartOption(cpuHistory.value, '#64748b'));
const memChartOption = computed(() => getChartOption(memHistory.value, '#64748b'));

// API Calls
const fetchRealtimeStats = async () => {
  try {
    const res = await getSystemStatus();
    if (res.code === ResultCode.SUCCESS) {
      if (realtimeStats.value) {
        prevApiCallCount.value = realtimeStats.value.api_call_count;
        prevSystemUptime.value = realtimeStats.value.system_uptime;
        prevCpuUsage.value = realtimeStats.value.cpu_usage;
        prevMemUsage.value = realtimeStats.value.mem_usage;
        prevGoRoutines.value = realtimeStats.value.go_routines;
        prevTotalVolume.value = realtimeStats.value.total_volume;

        // Calculate QPS (approximate)
        // Interval is 500ms, so multiply diff by 2
        const diff = res.data.api_call_count - prevApiCallCount.value;
        qps.value = diff > 0 ? diff * 2 : 0;
      }
      realtimeStats.value = res.data;
      
      // Update history
      cpuHistory.value.push(res.data.cpu_usage);
      if (cpuHistory.value.length > 300) cpuHistory.value.shift();
      
      memHistory.value.push(res.data.mem_usage);
      if (memHistory.value.length > 300) memHistory.value.shift();
    }
  } catch (error) {
    console.error("Failed to fetch realtime stats", error);
  }
};

const fetchStats = async () => {
  try {
    const res = await getSystemStats();
    if (res.code === ResultCode.SUCCESS) {
      stats.value = res.data;
    }
  } catch (error) {
    console.error("Failed to fetch system stats", error);
  }
};

// Handlers
const toggleUptimeUnit = () => {
  const map: Record<string, string> = { 'auto': 's', 's': 'm', 'm': 'h', 'h': 'auto' };
  uptimeUnit.value = map[uptimeUnit.value];
};

const getUptimeValue = (seconds: number) => {
  if (uptimeUnit.value === 's') return seconds;
  if (uptimeUnit.value === 'm') return Math.floor(seconds / 60);
  if (uptimeUnit.value === 'h') return Math.floor(seconds / 3600);
  return 0; 
};

const getUptimeSuffix = () => {
  const map: Record<string, string> = { 's': '秒', 'm': '分钟', 'h': '小时' };
  return map[uptimeUnit.value] || '';
};

const formatDuration = (seconds: number) => {
  if (!seconds) return "0h";
  const hours = Math.floor(seconds / 3600);
  const days = Math.floor(hours / 24);
  if (days > 0) return `${days}d ${hours % 24}h`;
  return `${hours}h`;
};

const formatBytes = (bytes: number) => {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

const handleGeneratePlaylistDesc = async () => {
  loading.playlist = true;
  try {
    const res = await generatePublicPlaylistDescriptions();
    if (res.code === ResultCode.SUCCESS) {
      message.success(res.message || "任务已在后台开始");
    } else {
      message.error(res.message || "请求失败");
    }
  } catch (error) {
    message.error("发生错误");
  } finally {
    loading.playlist = false;
  }
};

const handleGenerateArtistDesc = async () => {
  loading.artist = true;
  try {
    const res = await generateArtistDescriptions();
    if (res.code === ResultCode.SUCCESS) {
      message.success(res.message || "任务已在后台开始");
    } else {
      message.error(res.message || "请求失败");
    }
  } catch (error) {
    message.error("发生错误");
  } finally {
    loading.artist = false;
  }
};

const handleGenerateAlbumDesc = async () => {
  loading.album = true;
  try {
    const res = await generateAlbumDescriptions();
    if (res.code === ResultCode.SUCCESS) {
      message.success(res.message || "任务已在后台开始");
    } else {
      message.error(res.message || "请求失败");
    }
  } catch (error) {
    message.error("发生错误");
  } finally {
    loading.album = false;
  }
};

const handleGenerateAllIntros = async () => {
  loading.intro = true;
  try {
    const res = await generateAllPublicPlaylistIntros();
    if (res.code === ResultCode.SUCCESS) {
      message.success("公共歌单开场白生成任务已启动");
    } else {
      message.error(res.message || "生成失败");
    }
  } catch (error) {
    message.error("发生错误");
  } finally {
    loading.intro = false;
  }
};

const handleScanMusic = async () => {
  loading.scan = true;
  try {
    const res = await scanMusic();
    if (res.code === ResultCode.SUCCESS) {
      message.success(`扫描完成：新增 ${res.data.added}，更新 ${res.data.updated}`);
      fetchStats();
    } else {
      message.error(res.message || "扫描失败");
    }
  } catch (error) {
    message.error("发生错误");
  } finally {
    loading.scan = false;
  }
};

const handleExportExcel = () => {
  window.open("/api/system/export/excel", "_blank");
};

const handleApiDoc = () => {
  window.open("https://ce9bjycbn4.apifox.cn/416729186e0", "_blank");
};

const handleUserManage = () => {
  router.push("/admin/users");
};

const handlePlaylistManage = () => {
  router.push("/admin/playlists");
};

const handleResetSystem = async () => {
  loading.reset = true;
  try {
    const res = await resetSystem();
    if (res.code === ResultCode.SUCCESS) {
      message.success("系统重置成功");
      setTimeout(() => {
        window.location.reload();
      }, 1000);
    } else {
      message.error(res.message || "重置失败");
    }
  } catch (error) {
    message.error("发生错误");
  } finally {
    loading.reset = false;
  }
};

// Lifecycle
onMounted(() => {
  fetchStats();
  fetchRealtimeStats();
  intervalId.value = setInterval(fetchRealtimeStats, 500) as unknown as number;
});

onUnmounted(() => {
  if (intervalId.value) {
    clearInterval(intervalId.value);
  }
});
</script>

<style lang="scss" scoped>
@import url('https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700;800&display=swap');

.admin-dashboard-container {
  position: relative;
  min-height: 100vh;
  width: 100%;
  padding: 48px;
  background: #f8f9fc;
  font-family: 'Plus Jakarta Sans', sans-serif;
  overflow-x: hidden;
  box-sizing: border-box;

  :global(.dark) & {
    background: #0f1115;
    color: #fff;
  }
}

/* Ambient Background */
.ambient-glows {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 0;
  pointer-events: none;
  overflow: hidden;

  .glow-orb {
    position: absolute;
    border-radius: 50%;
    filter: blur(100px);
    opacity: 0.5;
    animation: float 20s ease-in-out infinite;
  }

  .orb-1 {
    top: -10%;
    left: -10%;
    width: 800px;
    height: 800px;
    background: radial-gradient(circle, rgba(99, 102, 241, 0.25) 0%, rgba(99, 102, 241, 0) 70%);
    animation-delay: 0s;
  }

  .orb-2 {
    bottom: -20%;
    right: -10%;
    width: 600px;
    height: 600px;
    background: radial-gradient(circle, rgba(139, 92, 246, 0.25) 0%, rgba(139, 92, 246, 0) 70%);
    animation-delay: -5s;
  }
  
  .orb-3 {
    top: 30%;
    left: 40%;
    width: 500px;
    height: 500px;
    background: radial-gradient(circle, rgba(236, 72, 153, 0.2) 0%, rgba(236, 72, 153, 0) 70%);
    animation-delay: -10s;
  }
}

@keyframes float {
  0%, 100% { transform: translate(0, 0); }
  33% { transform: translate(30px, -50px); }
  66% { transform: translate(-20px, 20px); }
}

.noise-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/200' opacity='0.03'%3E%3Cfilter id='noiseFilter'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.8' numOctaves='3' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noiseFilter)'/%3E%3C/svg%3E");
  z-index: 1;
  pointer-events: none;
}

.dashboard-content {
  position: relative;
  z-index: 2;
  max-width: 1400px;
  margin: 0 auto;
  padding-bottom: 64px;
  animation: fadeInUp 0.8s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

.header-section {
  margin-bottom: 48px;

  .header-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .title-wrapper {
    display: flex;
    align-items: center;
    gap: 24px;

    .icon-box {
      width: 64px;
      height: 64px;
      background: rgba(255, 255, 255, 0.8);
      border-radius: 20px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 28px;
      color: #4f46e5;
      box-shadow: 0 20px 40px -10px rgba(79, 70, 229, 0.15), 0 0 0 1px rgba(255, 255, 255, 0.5);
      backdrop-filter: blur(10px);
      transition: transform 0.3s ease;

      :global(.dark) & {
        background: rgba(30, 41, 59, 0.8);
        border-color: rgba(255, 255, 255, 0.1);
      }
      
      &:hover {
        transform: scale(1.05) rotate(5deg);
      }
    }

    .page-title {
      font-size: 36px;
      font-weight: 800;
      background: linear-gradient(135deg, #111827 0%, #4b5563 100%);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      margin: 0;
      letter-spacing: -0.03em;
      line-height: 1.1;

      :global(.dark) & { 
        background: linear-gradient(135deg, #f9fafb 0%, #9ca3af 100%);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
      }
    }

    .page-subtitle {
        font-size: 15px;
        font-weight: 500;
        color: #64748b;
        margin: 6px 0 0;
        letter-spacing: 0.01em;

        :global(.dark) & { color: #94a3b8; }
      }
  }

  .status-badge {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 16px;
    background: rgba(16, 185, 129, 0.1);
    border: 1px solid rgba(16, 185, 129, 0.2);
    border-radius: 100px;
    color: #059669;
    font-size: 13px;
    font-weight: 600;
    backdrop-filter: blur(4px);

    .dot {
      width: 8px;
      height: 8px;
      background: #10b981;
      border-radius: 50%;
      box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.2);
      animation: pulse 2s infinite;
    }
  }
}

@keyframes pulse {
  0% { box-shadow: 0 0 0 0 rgba(16, 185, 129, 0.4); }
  70% { box-shadow: 0 0 0 6px rgba(16, 185, 129, 0); }
  100% { box-shadow: 0 0 0 0 rgba(16, 185, 129, 0); }
}

/* Dashboard Section Shared Styles */
.dashboard-section {
  margin-bottom: 64px;
  position: relative;
  z-index: 1;

  .section-label {
    font-size: 12px;
    text-transform: uppercase;
    letter-spacing: 0.1em;
    color: #64748b;
    font-weight: 700;
    margin-bottom: 24px;
    padding-left: 12px;
    border-left: 3px solid;
    line-height: 1;
    border-color: #cbd5e1;
    opacity: 0.8;

    :global(.dark) & { color: #9ca3af; border-color: #475569; }
    
    &.danger { color: #ef4444; border-color: #ef4444; }
  }
  
  &:nth-child(2) .section-label { border-color: #64748b; }
  &:nth-child(3) .section-label { border-color: #8b5cf6; }
  &:nth-child(4) .section-label { border-color: #ec4899; }
}

/* Glass Panel Base */
.glass-panel {
  background: rgba(255, 255, 255, 0.65);
  backdrop-filter: blur(20px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.6);
  box-shadow: 
    0 4px 6px -1px rgba(0, 0, 0, 0.02),
    0 2px 4px -1px rgba(0, 0, 0, 0.02),
    inset 0 0 0 1px rgba(255, 255, 255, 0.5);
  
  :global(.dark) & {
    background: rgba(30, 41, 59, 0.6);
    border-color: rgba(255, 255, 255, 0.05);
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.2);
  }
}

/* Realtime Monitor Grid */
.realtime-monitor-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24px;
  
  .span-2 { grid-column: span 2; }

  @media (max-width: 1200px) {
    grid-template-columns: repeat(2, 1fr);
    .span-2 { grid-column: span 1; }
  }

  @media (max-width: 640px) {
    grid-template-columns: 1fr;
    .span-2 { grid-column: span 1; }
  }
}

.monitor-card {
  display: flex;
  align-items: center;
  gap: 24px;
  padding: 24px;
  border-radius: 24px;
  transition: all 0.4s cubic-bezier(0.16, 1, 0.3, 1);
  position: relative;
  overflow: hidden;

  &.gray-theme {
    .monitor-icon.gray { background: rgba(241, 245, 249, 0.8); color: #64748b; }
  }

  &:hover {
    transform: translateY(-4px);
    background: rgba(255, 255, 255, 0.85);
    box-shadow: 
      0 20px 25px -5px rgba(0, 0, 0, 0.05),
      0 8px 10px -6px rgba(0, 0, 0, 0.01);
    
    .monitor-icon.gray {
      background: rgba(226, 232, 240, 0.8);
      transform: scale(1.05);
      color: #334155;
    }
  }
  
  &.span-2 { justify-content: space-between; padding-right: 0; }

  .monitor-main {
    display: flex;
    align-items: center;
    gap: 20px;
    position: relative;
    z-index: 2;
  }
  
  .monitor-chart {
    flex: 1;
    height: 100%;
    min-width: 100px;
    width: 50%;
    position: absolute;
    right: 0;
    bottom: 0;
    opacity: 0.6;
    pointer-events: none;
    mask-image: linear-gradient(to right, transparent, black 40%);
    -webkit-mask-image: linear-gradient(to right, transparent, black 40%);
    transition: opacity 0.4s ease;
  }
  
  &:hover .monitor-chart { opacity: 1; }
  
  .mini-chart { width: 100%; height: 100%; min-height: 80px; }

  .monitor-icon {
    width: 52px;
    height: 52px;
    border-radius: 16px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 26px;
    flex-shrink: 0;
    transition: all 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
  }

  .monitor-content, .monitor-info {
    display: flex;
    flex-direction: column;
    z-index: 2;
    
    .monitor-label {
      font-size: 13px;
      font-weight: 600;
      text-transform: uppercase;
      letter-spacing: 0.02em;
      color: #94a3b8;
      margin-bottom: 6px;
    }

    .monitor-value {
      font-size: 28px;
      font-weight: 700;
      color: #1e293b;
      line-height: 1.1;
      font-feature-settings: "tnum";
      font-variant-numeric: tabular-nums;
      letter-spacing: -0.02em;

      &.text-sm { font-size: 20px; }
    }
  }
  
  .progress-container { margin-top: 4px; }
  .progress-text { font-size: 13px; font-weight: 700; color: #334155; }
}

/* Stats Grid */
.stats-grid {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 1fr;
  gap: 24px;
  margin-bottom: 32px;

  @media (max-width: 1200px) { grid-template-columns: 1fr 1fr; }
  @media (max-width: 768px) { grid-template-columns: 1fr; }
}

.stat-card {
  @extend .glass-panel;
  padding: 32px;
  position: relative;
  overflow: hidden;
  transition: all 0.4s cubic-bezier(0.16, 1, 0.3, 1);
  border-radius: 28px;

  &:hover {
    transform: translateY(-6px);
    box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.1);
  }

  &.primary-stat {
    background: linear-gradient(120deg, #6366f1 0%, #8b5cf6 50%, #d946ef 100%);
    background-size: 200% 200%;
    color: white;
    border: none;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    animation: gradientFlow 8s ease infinite;
    box-shadow: 0 20px 40px -10px rgba(99, 102, 241, 0.4);

    .stat-label { 
      font-size: 14px;
      font-weight: 600;
      color: rgba(255, 255, 255, 0.9); 
      text-transform: uppercase;
      letter-spacing: 0.05em;
    }
    
    .stat-value { 
      color: white; 
      font-size: 56px; 
      font-weight: 800;
      letter-spacing: -0.03em;
      text-shadow: 0 2px 10px rgba(0,0,0,0.1);
      margin: 16px 0;
    }
    
    .stat-trend { 
      display: flex;
      align-items: center;
      gap: 6px;
      background: rgba(255, 255, 255, 0.2);
      width: fit-content;
      padding: 6px 12px;
      border-radius: 100px;
      font-size: 13px;
      font-weight: 600;
      backdrop-filter: blur(4px);
    }
    
    .stat-chart-decoration {
      position: absolute; bottom: 0; left: 0; width: 100%; height: 80px; opacity: 0.25;
      svg { width: 100%; height: 100%; }
      transition: transform 0.5s;
    }

    .stat-icon-bg {
      position: absolute; right: -30px; top: -30px; font-size: 220px; opacity: 0.1;
      transform: rotate(-15deg);
      transition: all 0.6s cubic-bezier(0.34, 1.56, 0.64, 1);
    }

    &:hover {
      .stat-icon-bg { transform: rotate(0deg) scale(1.1) translate(10px, 10px); opacity: 0.15; }
      .stat-chart-decoration { transform: scaleY(1.1); }
    }
  }

  &.secondary-stat {
    display: flex;
    flex-direction: column;
    justify-content: center;
    cursor: pointer;
    
    .stat-header { display: flex; align-items: center; gap: 16px; margin-bottom: 20px; }

    .stat-icon-small {
      width: 44px; height: 44px; border-radius: 14px; display: flex; align-items: center; justify-content: center; font-size: 22px;
      transition: all 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
      &.purple { background: rgba(139, 92, 246, 0.1); color: #8b5cf6; }
    }

    &:hover .stat-icon-small {
      transform: scale(1.1) rotate(5deg);
      background: rgba(139, 92, 246, 0.2);
      box-shadow: 0 8px 16px -4px rgba(139, 92, 246, 0.3);
    }

    .stat-value-small { 
      font-size: 36px; 
      font-weight: 800; 
      color: #1e293b; 
      line-height: 1; 
      margin-bottom: 8px; 
      letter-spacing: -0.02em;
    }
    
    .stat-footer { font-size: 14px; color: #94a3b8; font-weight: 500; }
  }
}

@keyframes gradientFlow {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

/* Tools Integration */
.tools-integration { margin-top: 32px; }

.ops-grid-compact {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 24px;
}

.op-card-compact {
  @extend .glass-panel;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24px 28px;
  border-radius: 20px;
  transition: all 0.4s cubic-bezier(0.16, 1, 0.3, 1);
  
  &.white-theme { 
    background: #ffffff;
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.02); 
  }
  
  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 20px 30px -10px rgba(0, 0, 0, 0.06);
    border-color: rgba(59, 130, 246, 0.3);
    
    .compact-icon { transform: scale(1.1) rotate(10deg); }
  }

  .compact-icon {
    width: 52px; height: 52px; border-radius: 16px; display: flex; align-items: center; justify-content: center; font-size: 24px; flex-shrink: 0;
    transition: all 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
    
    &.green { background: rgba(16, 185, 129, 0.1); color: #10b981; }
    &.teal { background: rgba(20, 184, 166, 0.1); color: #14b8a6; }
    &.blue { background: rgba(59, 130, 246, 0.1); color: #3b82f6; }
    &.red { background: rgba(239, 68, 68, 0.1); color: #ef4444; }
  }
  
  &:hover .compact-icon.green { background: rgba(16, 185, 129, 0.2); box-shadow: 0 8px 16px -4px rgba(16, 185, 129, 0.2); }
  &:hover .compact-icon.teal { background: rgba(20, 184, 166, 0.2); box-shadow: 0 8px 16px -4px rgba(20, 184, 166, 0.2); }
  &:hover .compact-icon.blue { background: rgba(59, 130, 246, 0.2); box-shadow: 0 8px 16px -4px rgba(59, 130, 246, 0.2); }
  &:hover .compact-icon.red { background: rgba(239, 68, 68, 0.2); box-shadow: 0 8px 16px -4px rgba(239, 68, 68, 0.2); }

  .compact-info {
    flex-grow: 1; margin: 0 20px;
    h3 { font-size: 16px; font-weight: 700; color: #1f2937; margin: 0 0 4px 0; letter-spacing: -0.01em; }
    p { font-size: 13px; color: #6b7280; margin: 0; }
  }
  
  &:hover .compact-info h3 { color: #111827; }

  &.danger-zone {
    border: 1px solid rgba(239, 68, 68, 0.15);
    background: linear-gradient(145deg, #fffafa, #fef2f2);
    &:hover { 
      background: #fef2f2;
      border-color: rgba(239, 68, 68, 0.4); 
      transform: translateY(-4px);
      box-shadow: 0 15px 30px -5px rgba(239, 68, 68, 0.1);
    }
  }
}

.action-circle-btn { 
  width: 44px; height: 44px; flex-shrink: 0;
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
  font-size: 18px;
  
  &:hover {
    transform: scale(1.1);
    box-shadow: 0 8px 20px rgba(0,0,0,0.1);
  }
}

/* AI Section & Ops Grid */
.ops-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 48px 24px;
  position: relative;
  z-index: 3;
}

@media (min-width: 1400px) {
  .ops-grid {
    grid-template-columns: repeat(4, 1fr);
  }
}

@media (max-width: 768px) {
  .ops-grid {
    grid-template-columns: 1fr;
    gap: 24px;
  }
}

.op-card {
  @extend .glass-panel;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  min-height: 320px;
  padding: 40px 32px;
  border-radius: 32px;
  position: relative;
  z-index: 1;
  overflow: hidden;
  transition: all 0.5s cubic-bezier(0.16, 1, 0.3, 1);
  cursor: pointer;
  background: linear-gradient(160deg, rgba(255,255,255,0.7) 0%, rgba(255,255,255,0.4) 100%);

  &:hover {
    transform: translateY(-10px) scale(1.01);
    background: linear-gradient(160deg, rgba(255,255,255,0.9) 0%, rgba(255,255,255,0.7) 100%);
    box-shadow: 
      0 30px 60px -15px rgba(0, 0, 0, 0.1), 
      inset 0 0 0 1px rgba(255, 255, 255, 0.6);
    
    .op-icon-large { transform: scale(1.15) rotate(8deg); }
    .op-info h3 { color: #000; }
  }
  
  &:active { transform: translateY(-4px) scale(0.99); }

  .op-card-glow {
    position: absolute; top: -60px; right: -60px; width: 240px; height: 240px;
    filter: blur(70px); opacity: 0.12; border-radius: 50%;
    transition: all 0.6s ease-out;
    &.blue { background: #3b82f6; }
    &.purple { background: #a855f7; }
    &.orange { background: #f97316; }
    &.teal { background: #14b8a6; }
  }
  
  &:hover .op-card-glow { 
    opacity: 0.25; 
    transform: scale(1.3) translate(-20px, 20px);
  }

  .op-icon-large {
    width: 72px; height: 72px; border-radius: 24px;
    display: flex; align-items: center; justify-content: center;
    font-size: 36px; margin-bottom: 32px;
    transition: all 0.5s cubic-bezier(0.34, 1.56, 0.64, 1);
    box-shadow: 0 10px 20px -5px rgba(0,0,0,0.05);
    
    &.blue { background: rgba(59, 130, 246, 0.1); color: #3b82f6; }
    &.purple { background: rgba(168, 85, 247, 0.1); color: #a855f7; }
    &.orange { background: rgba(249, 115, 22, 0.1); color: #f97316; }
    &.teal { background: rgba(20, 184, 166, 0.1); color: #14b8a6; }
  }
  
  &:hover .op-icon-large.blue { background: rgba(59, 130, 246, 0.2); box-shadow: 0 15px 30px -5px rgba(59, 130, 246, 0.25); }
  &:hover .op-icon-large.purple { background: rgba(168, 85, 247, 0.2); box-shadow: 0 15px 30px -5px rgba(168, 85, 247, 0.25); }
  &:hover .op-icon-large.orange { background: rgba(249, 115, 22, 0.2); box-shadow: 0 15px 30px -5px rgba(249, 115, 22, 0.25); }
  &:hover .op-icon-large.teal { background: rgba(20, 184, 166, 0.2); box-shadow: 0 15px 30px -5px rgba(20, 184, 166, 0.25); }

  .op-info {
    margin-bottom: 32px; flex-grow: 1;
    h3 { font-size: 22px; font-weight: 800; color: #111827; margin: 0 0 12px 0; transition: color 0.3s; letter-spacing: -0.02em; }
    p { font-size: 15px; color: #64748b; line-height: 1.6; margin: 0; font-weight: 500; }
  }
  
  .op-action {
    width: 100%;
  }
}

.glass-btn {
  background: rgba(255, 255, 255, 0.5);
  border: 1px solid rgba(0, 0, 0, 0.05);
  backdrop-filter: blur(4px);
  border-radius: 14px;
  font-weight: 700;
  padding: 8px 16px;
  height: auto;
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
  color: #4b5563;
  
  &:hover { 
    background: white; 
    transform: translateY(-2px); 
    box-shadow: 0 4px 12px rgba(0,0,0,0.08);
    border-color: rgba(0,0,0,0.08);
    color: #111827;
  }
  
  &:active { transform: translateY(0); }
}
</style>