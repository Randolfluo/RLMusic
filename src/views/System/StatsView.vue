<template>
  <div class="system-stats-view">
    <!-- Background Decoration - 温暖米色调 -->
    <div class="bg-decoration">
      <div class="blob blob-1"></div>
      <div class="blob blob-2"></div>
      <div class="blob blob-3"></div>
    </div>

    <div class="container">
      <div class="page-header">
        <div class="header-badge">
          <n-icon :component="Server" />
          System
        </div>
        <h1 class="page-title">系统信息</h1>
        <p class="page-subtitle">平台音乐与用户数据概览</p>
      </div>

      <div class="stats-wrapper glass-panel">
        <SystemStats />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import SystemStats from "@/components/Home/SystemStats.vue";
import { Server } from "@icon-park/vue-next";
import { NIcon } from "naive-ui";
</script>

<style lang="scss" scoped>
@import url('https://fonts.googleapis.com/css2?family=Playfair+Display:wght@600;700&family=Plus+Jakarta+Sans:wght@400;500;600;700&display=swap');

.system-stats-view {
  width: 100%;
  min-height: calc(100vh - 140px);
  padding: 40px 20px;
  display: flex;
  justify-content: center;
  position: relative;
  overflow: hidden;
  overflow-x: hidden;
  background: #faf8f5;
  font-family: 'Plus Jakarta Sans', sans-serif;

  @media (max-width: 768px) {
    padding: 20px 16px;
  }

  /* Animated Background - 温暖色调 */
  .bg-decoration {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100vh;
    z-index: -1;
    pointer-events: none;
    overflow: hidden;

    .blob {
      position: absolute;
      border-radius: 50%;
      filter: blur(100px);
      opacity: 0.5;
      animation: blob-float 20s infinite ease-in-out;
    }

    .blob-1 {
      width: 600px;
      height: 600px;
      background: linear-gradient(135deg, rgba(224, 122, 95, 0.25), rgba(212, 165, 116, 0.2));
      top: -200px;
      right: -100px;
      animation-delay: 0s;
    }

    .blob-2 {
      width: 500px;
      height: 500px;
      background: linear-gradient(135deg, rgba(61, 139, 139, 0.2), rgba(91, 141, 184, 0.15));
      bottom: -150px;
      left: -100px;
      animation-delay: -5s;
    }

    .blob-3 {
      width: 400px;
      height: 400px;
      background: linear-gradient(135deg, rgba(124, 111, 174, 0.15), rgba(224, 122, 95, 0.1));
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      animation: blob-float-center 20s infinite ease-in-out;
      animation-delay: -10s;
      opacity: 0.3;
    }
  }

  @keyframes blob-float {
    0%, 100% { transform: translate(0, 0) scale(1); }
    25% { transform: translate(30px, -30px) scale(1.05); }
    50% { transform: translate(-20px, 20px) scale(0.95); }
    75% { transform: translate(15px, 15px) scale(1.02); }
  }

  @keyframes blob-float-center {
    0%, 100% { transform: translate(-50%, -50%) scale(1); }
    25% { transform: translate(calc(-50% + 30px), calc(-50% - 30px)) scale(1.05); }
    50% { transform: translate(calc(-50% - 20px), calc(-50% + 20px)) scale(0.95); }
    75% { transform: translate(calc(-50% + 15px), calc(-50% + 15px)) scale(1.02); }
  }
}

.container {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  position: relative;
  z-index: 1;
  animation: fade-in-up 0.6s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes fade-in-up {
  from { opacity: 0; transform: translateY(30px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Page Header */
.page-header {
  text-align: center;
  margin-bottom: 40px;

  @media (max-width: 768px) {
    margin-bottom: 24px;
  }

  .header-badge {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    font-size: 12px;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 1.5px;
    padding: 8px 16px;
    border-radius: 100px;
    background: linear-gradient(135deg, #3d8b8b 0%, #5b8db8 100%);
    color: white;
    margin-bottom: 20px;
    box-shadow: 0 8px 20px rgba(61, 139, 139, 0.35);
    animation: badge-pulse 3s ease-in-out infinite;

    @keyframes badge-pulse {
      0%, 100% { transform: scale(1); box-shadow: 0 8px 20px rgba(61, 139, 139, 0.35); }
      50% { transform: scale(1.02); box-shadow: 0 12px 30px rgba(61, 139, 139, 0.45); }
    }
  }

  .page-title {
    font-family: 'Playfair Display', serif;
    font-size: 42px;
    font-weight: 700;
    margin: 0 0 12px 0;
    background: linear-gradient(135deg, #e07a5f 0%, #d4a574 50%, #3d8b8b 100%);
    -webkit-background-clip: text;
    background-clip: text;
    -webkit-text-fill-color: transparent;
    line-height: 1.2;
    letter-spacing: -0.02em;

    @media (max-width: 768px) {
      font-size: 32px;
    }
  }

  .page-subtitle {
    font-size: 16px;
    color: #666666;
    margin: 0;
    opacity: 0.8;
  }
}

/* Glass Panel Wrapper */
.stats-wrapper {
  background: #f5f2ed;
  border-radius: 32px;
  border: 1px solid #ebe7e0;
  box-shadow:
    0 4px 20px rgba(0, 0, 0, 0.06),
    inset 0 0 0 1px rgba(255, 255, 255, 0.8);
  padding: 40px;
  transition: all 0.4s cubic-bezier(0.25, 0.8, 0.25, 1);

  &:hover {
    box-shadow:
      0 8px 30px rgba(0, 0, 0, 0.1),
      inset 0 0 0 1px rgba(255, 255, 255, 0.9);
  }

  @media (max-width: 768px) {
    padding: 24px 16px;
    border-radius: 24px;
  }
}

/* Dark Mode Support */
:global(.dark) {
  .system-stats-view {
    background: #1a1a1a;

    .stats-wrapper {
      background: rgba(30, 30, 30, 0.8);
      border-color: rgba(255, 255, 255, 0.1);
      box-shadow:
        0 4px 20px rgba(0, 0, 0, 0.3),
        inset 0 0 0 1px rgba(255, 255, 255, 0.05);
    }
  }
}
</style>
