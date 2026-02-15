<template>
  <!-- 关于本站 -->
  <n-modal
    class="s-modal about-modal"
    v-model:show="showAboutModal"
    preset="card"
    title="关于本站"
    :bordered="false"
    transform-origin="center"
  >
    <div class="copyright">
      <div class="header-section">
        <div class="logo-circle">
          <img src="/images/logo/favicon.png" alt="logo" class="modal-logo" />
        </div>
        <div class="desc">
          <h2 class="name">{{ packageJson.name }}</h2>
          <n-tag size="small" type="primary" round class="version-tag">
            v{{ packageJson.version }}
          </n-tag>
        </div>
      </div>
      
      <div class="content-body">
        <p class="slogan">让音乐回归纯粹，享受本地播放的极致体验。</p>
        
        <div class="info-grid">
          <div class="info-item">
            <span class="label">作者</span>
            <n-a
              :href="packageJson.home"
              target="_blank"
              class="value link"
              v-html="packageJson.author"
            />
          </div>
          <div class="info-item">
            <span class="label">版权</span>
            <span class="value">© 2026 - {{ new Date().getFullYear() }}</span>
          </div>
          <div class="info-item" v-if="icp">
            <span class="label">备案</span>
            <n-a
              class="value link"
              href="https://beian.miit.gov.cn/"
              target="_blank"
              v-html="icp"
            />
          </div>
        </div>
      </div>

      <div class="footer-actions">
        <n-button
          class="github-btn"
          type="primary"
          secondary
          round
          size="large"
          @click="jumpUrl(packageJson.github)"
        >
          <template #icon>
            <n-icon :component="GithubOne" />
          </template>
          访问 GitHub 仓库
        </n-button>
      </div>
    </div>
  </n-modal>
</template>

<script setup>
import { GithubOne } from "@icon-park/vue-next";
import packageJson from "@/../package.json";

// 关于本站数据
const showAboutModal = ref(false);
const icp = ref(import.meta.env.VITE_ICP ? import.meta.env.VITE_ICP : null);

// 链接跳转
const jumpUrl = (url) => {
  window.open(url);
};

// 开启本站数据弹窗
const openAboutSite = () => {
  showAboutModal.value = true;
};

// 暴露方法
defineExpose({
  openAboutSite,
});
</script>

<style lang="scss" scoped>
.about-modal {
  width: 460px;
  border-radius: 20px;
  overflow: hidden;
  
  :deep(.n-card-header) {
    padding: 20px 24px;
    border-bottom: 1px solid var(--n-divider-color);
  }
  
  :deep(.n-card__content) {
    padding: 0;
  }
}

.copyright {
  display: flex;
  flex-direction: column;
  background-color: var(--n-card-color); // Keep standard background for readability
  
  .header-section {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 32px 24px 24px;
    background: linear-gradient(to bottom, rgba(var(--n-primary-color-rgb), 0.05), transparent);
    
    .logo-circle {
      width: 80px;
      height: 80px;
      background: white;
      border-radius: 20px;
      display: flex;
      align-items: center;
      justify-content: center;
      box-shadow: 0 8px 24px rgba(0,0,0,0.08);
      margin-bottom: 16px;
      
      .modal-logo {
        width: 50px;
        height: 50px;
      }
    }
    
    .desc {
      text-align: center;
      
      .name {
        font-size: 24px;
        font-weight: 800;
        margin: 0 0 8px 0;
        background: linear-gradient(120deg, var(--n-text-color) 0%, var(--n-text-color-3) 100%);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
      }
      
      .version-tag {
        font-weight: 600;
        padding: 0 8px;
      }
    }
  }
  
  .content-body {
    padding: 0 32px;
    text-align: center;
    
    .slogan {
      font-size: 15px;
      color: var(--n-text-color-2);
      margin-bottom: 24px;
      line-height: 1.6;
    }
    
    .info-grid {
      display: grid;
      grid-template-columns: repeat(3, 1fr);
      gap: 12px;
      margin-bottom: 24px;
      
      .info-item {
        display: flex;
        flex-direction: column;
        align-items: center;
        background: var(--n-action-color);
        padding: 12px 8px;
        border-radius: 12px;
        
        .label {
          font-size: 12px;
          color: var(--n-text-color-3);
          margin-bottom: 4px;
        }
        
        .value {
          font-size: 13px;
          font-weight: 600;
          color: var(--n-text-color);
          
          &.link {
            color: var(--n-primary-color);
            text-decoration: none;
            
            &:hover {
              text-decoration: underline;
            }
          }
        }
      }
    }
  }
  
  .footer-actions {
    padding: 24px;
    display: flex;
    justify-content: center;
    border-top: 1px solid var(--n-divider-color);
    
    .github-btn {
      width: 100%;
      height: 44px;
      font-weight: 600;
    }
  }
}
</style>
