<template>
  <div class="login-page">
    <!-- Animated Background -->
    <div class="bg-decoration">
      <div class="blob blob-1"></div>
      <div class="blob blob-2"></div>
      <div class="blob blob-3"></div>
    </div>

    <!-- Noise Overlay -->
    <div class="noise-overlay"></div>

    <div class="login-card glass-card">
      <div class="back-home">
        <n-button quaternary @click="goHome" class="back-home-btn">
          <span class="back-home-icon">←</span>
          <span>返回首页</span>
        </n-button>
      </div>

      <div class="header">
        <div class="logo-wrapper">
          <img :src="logoSrc" alt="logo" class="logo-img" />
        </div>
        <h2 class="app-title">RLMusic</h2>
        <p class="app-subtitle">聆听生活每一刻</p>
      </div>

      <n-tabs
        animated
        class="custom-tabs"
        type="segment"
        justify-content="space-evenly"
        :pane-style="{
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          paddingTop: '24px',
        }"
      >
        <!-- 登录 -->
        <n-tab-pane name="login" tab="登录">
          <n-form
            class="form-container"
            ref="loginFormRef"
            :model="loginForm"
            :rules="loginRules"
            :show-label="false"
            size="large"
          >
            <n-form-item path="username">
              <n-input
                placeholder="请输入用户名"
                v-model:value="loginForm.username"
                @keydown.enter.prevent
                round
              >
                <template #prefix>
                  <n-icon :component="PersonRound" class="input-icon" />
                </template>
              </n-input>
            </n-form-item>
            <n-form-item path="password">
              <n-input
                type="password"
                show-password-on="click"
                placeholder="请输入密码"
                v-model:value="loginForm.password"
                @keydown.enter.prevent
                round
              >
                <template #prefix>
                  <n-icon :component="PasswordRound" class="input-icon" />
                </template>
              </n-input>
            </n-form-item>

            <n-form-item>
              <div class="options-row">
                <n-checkbox v-model:checked="rememberMe">记住密码</n-checkbox>
                <n-checkbox v-model:checked="autoLogin">自动登录</n-checkbox>
              </div>
            </n-form-item>

            <n-form-item>
              <n-button
                class="submit-btn"
                type="primary"
                round
                @click="handleLogin"
                :loading="loading"
                size="large"
              >
                <template #icon>
                  <n-icon :component="LoginIcon" />
                </template>
                登录
              </n-button>
            </n-form-item>
          </n-form>
        </n-tab-pane>

        <!-- 注册 -->
        <n-tab-pane name="register" tab="注册">
          <n-form
            class="form-container"
            ref="registerFormRef"
            :model="registerForm"
            :rules="registerRules"
            :show-label="false"
            size="large"
          >
            <n-form-item path="username">
              <n-input
                placeholder="用户名"
                v-model:value="registerForm.username"
                round
              >
                <template #prefix>
                  <n-icon :component="PersonRound" class="input-icon" />
                </template>
              </n-input>
            </n-form-item>
            <n-form-item path="password">
              <n-input
                type="password"
                show-password-on="click"
                placeholder="密码 (4-20位)"
                v-model:value="registerForm.password"
                round
              >
                <template #prefix>
                  <n-icon :component="PasswordRound" class="input-icon" />
                </template>
              </n-input>
            </n-form-item>
            <n-form-item path="email">
              <n-input
                placeholder="邮箱"
                v-model:value="registerForm.email"
                round
              >
                <template #prefix>
                  <n-icon :component="EmailRound" class="input-icon" />
                </template>
              </n-input>
            </n-form-item>
            <n-form-item>
              <n-button
                class="submit-btn"
                type="primary"
                round
                @click="handleRegister"
                :loading="loading"
                size="large"
              >
                <template #icon>
                  <n-icon :component="UserAdd" />
                </template>
                注册
              </n-button>
            </n-form-item>
          </n-form>
        </n-tab-pane>
      </n-tabs>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from "vue";
import { useRouter } from "vue-router";
import { userStore } from "@/store";
import { login, register } from "@/api/login";
import { PersonRound, PasswordRound, EmailRound } from "@vicons/material";
import { Login as LoginIcon, UserAdd } from "@icon-park/vue-next";
import { useMessage, NIcon } from "naive-ui";
import { ResultCode } from "@/utils/request";
import { aesEncrypt, aesDecrypt } from "@/utils/encrypt";

const router = useRouter();
const user = userStore();
const message = useMessage();
const loading = ref(false);
const rememberMe = ref(false);
const autoLogin = ref(false);

const logoSrc = computed(() => `${import.meta.env.BASE_URL}images/logo/favicon.png`);

const loginFormRef = ref(null);
const registerFormRef = ref(null);

const loginRules = {
  username: { required: true, message: "请输入用户名", trigger: "blur" },
  password: { required: true, message: "请输入密码", trigger: "blur" }
};

const registerRules = {
  username: { required: true, message: "请输入用户名", trigger: "blur" },
  password: [
    { required: true, message: "请输入密码", trigger: "blur" },
    { min: 4, max: 20, message: "长度在 4 到 20 个字符", trigger: "blur" }
  ],
  email: { required: true, message: "请输入邮箱", trigger: "blur" }
};

const loginForm = reactive({
  username: "",
  password: ""
});

const registerForm = reactive({
  username: "",
  password: "",
  email: ""
});

onMounted(() => {
  const savedUser = localStorage.getItem("remember_user");
  const savedPass = localStorage.getItem("remember_pass");
  const savedAutoLogin = localStorage.getItem("auto_login");

  if (savedAutoLogin === 'true') {
    autoLogin.value = true;
  }

  if (savedUser && savedPass) {
    loginForm.username = savedUser;
    try {
      loginForm.password = aesDecrypt(savedPass);
      rememberMe.value = true;
    } catch (e) {
      console.error("解密失败", e);
    }
  }
});

const handleLogin = (e) => {
  e.preventDefault();
  loginFormRef.value?.validate((errors) => {
    if (!errors) {
      loading.value = true;
      const data = { ...loginForm };
      login(data).then(res => {
        loading.value = false;
        if (res.code === ResultCode.SUCCESS) {
          message.success("登录成功");
          const userData = {
            userId: res.data.id,
            nickname: res.data.username,
            email: res.data.email,
            userGroup: res.data.user_group,
            avatarUrl: res.data.avatar || ""
          };
          user.setUserData(userData);
          sessionStorage.setItem("token", res.data.token);

          if (rememberMe.value || autoLogin.value) {
            localStorage.setItem("remember_user", loginForm.username);
            localStorage.setItem("remember_pass", aesEncrypt(loginForm.password));
          } else {
            localStorage.removeItem("remember_user");
            localStorage.removeItem("remember_pass");
          }

          localStorage.setItem("auto_login", String(autoLogin.value));
          router.push("/");
        } else {
          message.error(res.msg || "登录失败");
        }
      }).catch(err => {
        loading.value = false;
        message.error("登录出错");
        console.error(err);
      });
    }
  });
};

const handleRegister = (e) => {
  e.preventDefault();
  registerFormRef.value?.validate((errors) => {
    if (!errors) {
      loading.value = true;
      const data = { ...registerForm };
      register(data).then(res => {
        loading.value = false;
        if (res.code === ResultCode.SUCCESS) {
          message.success("注册成功，请登录");
        } else {
          message.error(res.msg || "注册失败");
        }
      }).catch(err => {
        loading.value = false;
        message.error("注册出错");
        console.error(err);
      });
    }
  });
};

const goHome = () => {
  router.push("/");
};
</script>

<style lang="scss" scoped>
.login-page {
  height: 100vh;
  width: 100vw;
  display: flex;
  justify-content: center;
  align-items: center;
  position: relative;
  overflow: hidden;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e8ec 100%);

  .bg-decoration {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    z-index: 0;
    pointer-events: none;
    overflow: hidden;

    .blob {
      position: absolute;
      border-radius: 50%;
      filter: blur(80px);
      opacity: 0.5;
      animation: blob-float 20s infinite ease-in-out;
    }

    .blob-1 {
      width: 600px;
      height: 600px;
      background: linear-gradient(135deg, #8b5cf6 0%, #a78bfa 100%);
      top: -200px;
      right: -200px;
      animation-delay: 0s;
    }

    .blob-2 {
      width: 500px;
      height: 500px;
      background: linear-gradient(135deg, #14b8a6 0%, #5eead4 100%);
      bottom: -150px;
      left: -150px;
      animation-delay: -7s;
    }

    .blob-3 {
      width: 400px;
      height: 400px;
      background: linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%);
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      animation-delay: -14s;
      animation-duration: 25s;
    }
  }

  @keyframes blob-float {
    0%, 100% { transform: translate(0, 0) scale(1); }
    25% { transform: translate(30px, -40px) scale(1.05); }
    50% { transform: translate(-20px, 20px) scale(0.95); }
    75% { transform: translate(20px, 15px) scale(1.02); }
  }

  .noise-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noiseFilter'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noiseFilter)'/%3E%3C/svg%3E");
    opacity: 0.03;
    z-index: 1;
    pointer-events: none;
  }

  .login-card {
    position: relative;
    z-index: 2;
    width: 420px;
    padding: 48px 40px;
    border-radius: 32px;
    background: rgba(255, 255, 255, 0.75);
    backdrop-filter: blur(24px) saturate(180%);
    border: 1px solid rgba(255, 255, 255, 0.5);
    box-shadow:
      0 20px 60px rgba(0, 0, 0, 0.1),
      0 8px 24px rgba(0, 0, 0, 0.05),
      inset 0 0 0 1px rgba(255, 255, 255, 0.6);
    display: flex;
    flex-direction: column;
    align-items: center;
    animation: card-enter 0.8s cubic-bezier(0.16, 1, 0.3, 1);

    @keyframes card-enter {
      0% { opacity: 0; transform: translateY(40px) scale(0.96); }
      100% { opacity: 1; transform: translateY(0) scale(1); }
    }

    .back-home {
      position: absolute;
      top: 20px;
      left: 20px;

      .back-home-btn {
        height: 36px;
        padding: 0 14px;
        border-radius: 12px;
        background: rgba(255, 255, 255, 0.5);
        border: 1px solid rgba(0, 0, 0, 0.08);
        backdrop-filter: blur(8px);
        font-weight: 500;
        color: var(--n-text-color-2);
        transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
        display: inline-flex;
        align-items: center;
        gap: 6px;

        .back-home-icon {
          font-size: 14px;
          transition: transform 0.2s ease;
        }

        &:hover {
          transform: translateY(-2px);
          background: rgba(255, 255, 255, 0.9);
          border-color: rgba(0, 0, 0, 0.12);
          box-shadow: 0 8px 20px rgba(0, 0, 0, 0.08);

          .back-home-icon {
            transform: translateX(-3px);
          }
        }
      }
    }

    .header {
      text-align: center;
      margin-bottom: 32px;

      .logo-wrapper {
        width: 88px;
        height: 88px;
        margin: 0 auto 20px;
        background: linear-gradient(135deg, #fff 0%, #f8f9fa 100%);
        border-radius: 24px;
        padding: 14px;
        box-shadow:
          0 12px 32px rgba(139, 92, 246, 0.15),
          0 4px 12px rgba(0, 0, 0, 0.05);
        display: flex;
        justify-content: center;
        align-items: center;
        transition: all 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);

        &:hover {
          transform: scale(1.08) rotate(5deg);
          box-shadow:
            0 16px 40px rgba(139, 92, 246, 0.2),
            0 6px 16px rgba(0, 0, 0, 0.08);
        }

        .logo-img {
          width: 100%;
          height: 100%;
          object-fit: contain;
        }
      }

      .app-title {
        font-family: 'Plus Jakarta Sans', sans-serif;
        font-size: 32px;
        font-weight: 800;
        margin: 0 0 8px 0;
        background: linear-gradient(135deg, #8b5cf6 0%, #14b8a6 100%);
        -webkit-background-clip: text;
        background-clip: text;
        -webkit-text-fill-color: transparent;
        letter-spacing: -0.02em;
      }

      .app-subtitle {
        font-size: 14px;
        color: var(--n-text-color-3);
        margin: 0;
        letter-spacing: 3px;
        text-transform: uppercase;
        font-weight: 500;
      }
    }

    .custom-tabs {
      width: 100%;

      :deep(.n-tabs-rail) {
        border-radius: 100px;
        background-color: rgba(0, 0, 0, 0.04);
        padding: 4px;
      }

      :deep(.n-tabs-tab) {
        border-radius: 100px;
        font-weight: 600;
        transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
        padding: 10px 24px;

        &.n-tabs-tab--active {
          background: linear-gradient(135deg, #8b5cf6 0%, #7c3aed 100%);
          color: white;
          box-shadow: 0 4px 12px rgba(139, 92, 246, 0.3);
          transform: scale(1.02);
        }

        &:not(.n-tabs-tab--active):hover {
          color: var(--n-color-primary);
        }
      }
    }

    .form-container {
      width: 100%;
      padding: 0 4px;
      box-sizing: border-box;

      :deep(.n-input) {
        background-color: rgba(255, 255, 255, 0.6);
        border: 1px solid rgba(0, 0, 0, 0.06);
        transition: all 0.3s ease;
        border-radius: 100px;
        height: 48px;

        &:hover, &:focus-within {
          background-color: rgba(255, 255, 255, 0.9);
          border-color: rgba(139, 92, 246, 0.3);
          box-shadow: 0 0 0 4px rgba(139, 92, 246, 0.08);
        }

        .n-input__prefix {
          margin-right: 12px;
          color: var(--n-text-color-3);
        }

        .n-input__input-el {
          font-size: 15px;
        }
      }

      .input-icon {
        font-size: 20px;
      }

      .options-row {
        display: flex;
        justify-content: space-between;
        width: 100%;
        padding: 0 8px;
        font-size: 13px;
      }

      .submit-btn {
        width: 100%;
        height: 50px;
        font-size: 16px;
        font-weight: 600;
        letter-spacing: 1px;
        background: linear-gradient(135deg, #8b5cf6 0%, #7c3aed 100%);
        border: none;
        box-shadow: 0 8px 24px rgba(139, 92, 246, 0.35);
        transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);

        &:hover {
          transform: translateY(-2px);
          box-shadow: 0 12px 32px rgba(139, 92, 246, 0.45);
        }

        &:active {
          transform: translateY(1px);
        }
      }
    }
  }
}

// 移动端适配
@media (max-width: 480px) {
  .login-page {
    .login-card {
      width: 90%;
      max-width: 360px;
      padding: 36px 24px;
      border-radius: 24px;

      .header {
        .logo-wrapper {
          width: 72px;
          height: 72px;
          border-radius: 20px;
        }

        .app-title {
          font-size: 26px;
        }

        .app-subtitle {
          font-size: 12px;
          letter-spacing: 2px;
        }
      }
    }
  }
}

/* Dark Mode Support */
:global(.dark) {
  .login-page {
    background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);

    .login-card {
      background: rgba(30, 30, 30, 0.75);
      border-color: rgba(255, 255, 255, 0.08);
      box-shadow:
        0 20px 60px rgba(0, 0, 0, 0.3),
        0 8px 24px rgba(0, 0, 0, 0.2),
        inset 0 0 0 1px rgba(255, 255, 255, 0.05);

      .back-home-btn {
        background: rgba(255, 255, 255, 0.08);
        border-color: rgba(255, 255, 255, 0.1);
        color: rgba(255, 255, 255, 0.7);

        &:hover {
          background: rgba(255, 255, 255, 0.15);
        }
      }

      .header {
        .logo-wrapper {
          background: linear-gradient(135deg, #2a2a3e 0%, #1e1e32 100%);
        }
      }

      .form-container {
        :deep(.n-input) {
          background-color: rgba(255, 255, 255, 0.08);
          border-color: rgba(255, 255, 255, 0.08);

          &:hover, &:focus-within {
            background-color: rgba(255, 255, 255, 0.12);
            border-color: rgba(139, 92, 246, 0.4);
          }
        }
      }
    }
  }
}
</style>
