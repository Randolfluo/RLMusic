<template>
  <Provider>
    <AppContent />
  </Provider>
</template>

<script setup lang="ts">
import Provider from "@/components/Provider/index.vue";
import AppContent from "./AppContent.vue";
</script>

<style>
/* Global Design System */
@import url('https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700;800&family=Playfair+Display:wght@400;500;600;700&display=swap');

:root {
  /* Primary Colors - Warm Editorial Beige/Cream Theme */
  --rl-primary: #e07a5f;
  --rl-primary-light: #f4a584;
  --rl-primary-dark: #c45c3f;

  /* Secondary Colors - Teal */
  --rl-secondary: #3d8b8b;
  --rl-secondary-light: #5ba8a8;
  --rl-accent: #d4a574;
  --rl-accent-light: #e8c9a0;

  /* Background Colors - Warm Beige */
  --rl-bg-primary: #faf8f5;
  --rl-bg-secondary: #f5f2ed;
  --rl-bg-glass: rgba(245, 242, 237, 0.85);
  --rl-bg-glass-dark: rgba(40, 40, 40, 0.8);

  /* Text Colors */
  --rl-text-primary: #1a1a1a;
  --rl-text-secondary: #666666;
  --rl-text-muted: #999999;

  /* Spacing */
  --rl-space-xs: 4px;
  --rl-space-sm: 8px;
  --rl-space-md: 16px;
  --rl-space-lg: 24px;
  --rl-space-xl: 32px;
  --rl-space-2xl: 48px;

  /* Border Radius */
  --rl-radius-sm: 8px;
  --rl-radius-md: 12px;
  --rl-radius-lg: 20px;
  --rl-radius-xl: 28px;
  --rl-radius-full: 9999px;

  /* Shadows */
  --rl-shadow-sm: 0 2px 8px rgba(0, 0, 0, 0.04);
  --rl-shadow-md: 0 4px 16px rgba(0, 0, 0, 0.06);
  --rl-shadow-lg: 0 8px 32px rgba(0, 0, 0, 0.08);
  --rl-shadow-glass: 0 8px 32px rgba(0, 0, 0, 0.08), inset 0 0 0 1px rgba(255, 255, 255, 0.5);

  /* Transitions */
  --rl-transition-fast: 0.15s cubic-bezier(0.4, 0, 0.2, 1);
  --rl-transition-base: 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  --rl-transition-slow: 0.5s cubic-bezier(0.16, 1, 0.3, 1);
  --rl-transition-bounce: 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
}

/* Global Animations */
@keyframes rl-fade-in {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes rl-fade-in-up {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes rl-fade-in-down {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes rl-float {
  0%, 100% { transform: translateY(0px); }
  50% { transform: translateY(-10px); }
}

@keyframes rl-pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

@keyframes rl-shimmer {
  0% { background-position: -200% 0; }
  100% { background-position: 200% 0; }
}

@keyframes rl-blob {
  0%, 100% { transform: translate(0, 0) scale(1); border-radius: 60% 40% 30% 70% / 60% 30% 70% 40%; }
  25% { transform: translate(20px, -30px) scale(1.1); border-radius: 30% 60% 70% 40% / 50% 60% 30% 60%; }
  50% { transform: translate(-20px, 20px) scale(0.9); border-radius: 50% 60% 30% 60% / 30% 60% 70% 40%; }
  75% { transform: translate(10px, 10px) scale(1.05); border-radius: 60% 40% 60% 30% / 70% 30% 50% 60%; }
}

@keyframes rl-gradient-shift {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

/* Utility Classes */
.rl-page-container {
  padding: 32px;
  max-width: 1400px;
  margin: 0 auto;
  min-height: 100vh;
}

.rl-glass-card {
  background: var(--rl-bg-glass);
  backdrop-filter: blur(20px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.5);
  border-radius: var(--rl-radius-xl);
  box-shadow: var(--rl-shadow-glass);
}

.rl-glass-panel {
  background: rgba(255, 255, 255, 0.5);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: var(--rl-radius-lg);
}

.rl-gradient-text {
  background: linear-gradient(135deg, #e07a5f 0%, #d4a574 50%, #3d8b8b 100%);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.rl-section-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 24px;
}

.rl-section-title {
  font-family: 'Plus Jakarta Sans', sans-serif;
  font-size: 24px;
  font-weight: 700;
  color: var(--rl-text-primary);
  margin: 0;
}

.rl-badge {
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  padding: 4px 10px;
  border-radius: var(--rl-radius-full);
  background: rgba(224, 122, 95, 0.1);
  color: var(--rl-primary);
}

/* Background Decorations */
.rl-bg-decoration {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  z-index: -1;
  pointer-events: none;
  overflow: hidden;
}

.rl-blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.35;
  animation: rl-blob 20s infinite ease-in-out;
}

.rl-blob-purple {
  background: linear-gradient(135deg, #7c6fae 0%, #a89bc9 100%);
}

.rl-blob-teal {
  background: linear-gradient(135deg, #3d8b8b 0%, #5b8db8 100%);
}

.rl-blob-orange {
  background: linear-gradient(135deg, #e07a5f 0%, #d4a574 100%);
}

.rl-noise-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noiseFilter'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noiseFilter)'/%3E%3C/svg%3E");
  opacity: 0.02;
  z-index: 0;
  pointer-events: none;
}

/* Page Transitions */
.rl-page-enter {
  animation: rl-fade-in-up 0.6s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

.rl-page-enter-delay-1 {
  animation: rl-fade-in-up 0.6s cubic-bezier(0.16, 1, 0.3, 1) 0.1s backwards;
}

.rl-page-enter-delay-2 {
  animation: rl-fade-in-up 0.6s cubic-bezier(0.16, 1, 0.3, 1) 0.2s backwards;
}

/* Button Styles */
.rl-btn-primary {
  background: linear-gradient(135deg, var(--rl-primary) 0%, var(--rl-primary-dark) 100%);
  border: none;
  color: white;
  font-weight: 600;
  padding: 12px 24px;
  border-radius: var(--rl-radius-full);
  transition: all var(--rl-transition-bounce);
  box-shadow: 0 4px 16px rgba(224, 122, 95, 0.3);
}

.rl-btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(224, 122, 95, 0.4);
}

/* Card Hover Effects */
.rl-hover-lift {
  transition: all var(--rl-transition-bounce);
}

.rl-hover-lift:hover {
  transform: translateY(-8px);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
}

/* Responsive */
@media (max-width: 768px) {
  :root {
    --rl-space-xl: 24px;
    --rl-space-2xl: 32px;
  }

  .rl-page-container {
    padding: 20px;
  }
}
</style>
