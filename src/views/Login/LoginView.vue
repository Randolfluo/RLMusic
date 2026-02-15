<template>
  <div class="login-container">
    <div class="login-background"></div>
    <div class="login-card glass-effect">
      <div class="header">
        <div class="logo-wrapper">
          <img src="/images/logo/favicon.png" alt="logo" class="logo-img" />
        </div>
        <h2 class="app-title">云音乐</h2>
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
          paddingTop: '20px',
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
import { ref, reactive, onMounted } from "vue";
import { useRouter } from "vue-router";
import { userStore } from "@/store";
import { login, register } from "@/api/login";
import { PersonRound, PasswordRound, EmailRound } from "@vicons/material";
import { useMessage } from "naive-ui";
import { ResultCode } from "@/utils/request";
import { aesEncrypt, aesDecrypt } from "@/utils/encrypt"; // 引入加密/解密工具

const router = useRouter();
const user = userStore();
const message = useMessage();
const loading = ref(false);
const rememberMe = ref(false); // 记住密码
const autoLogin = ref(false);  // 自动登录

const loginFormRef = ref(null);
const registerFormRef = ref(null);

// Validate rules
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

// Data
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

// Handlers
const handleLogin = (e) => {
  e.preventDefault();
  loginFormRef.value?.validate((errors) => {
    if (!errors) {
      loading.value = true;
      // 克隆表单数据
      const data = { ...loginForm };
      login(data).then(res => {
        loading.value = false;
        if (res.code === ResultCode.SUCCESS) {
          message.success("登录成功");
          const userData = {
             userId: res.data.id,
             nickname: res.data.username,
             email: res.data.email,
             avatarUrl: res.data.avatar || "" 
          };
          user.setUserData(userData);
          // 使用 sessionStorage 存储 token，这样每次关闭应用后 token 会自动清除，实现“每次打开由于没有token需要重新登录”
          sessionStorage.setItem("token", res.data.token); 
          
          // 如果勾选了自动登录，则强制记住密码
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
      // 克隆表单数据
      const data = { ...registerForm };
      register(data).then(res => {
        loading.value = false;
        if (res.code === ResultCode.SUCCESS) {
          message.success("注册成功，请登录");
          // Clear form not implemented, user can switch tab
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
</script>

<style lang="scss" scoped>
.login-container {
  height: 100vh;
  width: 100vw;
  display: flex;
  justify-content: center;
  align-items: center;
  position: relative;
  overflow: hidden;
  background: #f0f2f5;

  .login-background {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    opacity: 0.8;
    z-index: 0;
    
    &::before {
      content: '';
      position: absolute;
      top: -10%;
      left: -10%;
      width: 50%;
      height: 50%;
      border-radius: 50%;
      background: radial-gradient(circle, rgba(255,255,255,0.2) 0%, rgba(255,255,255,0) 70%);
      animation: float 10s infinite ease-in-out;
    }
    
    &::after {
      content: '';
      position: absolute;
      bottom: -10%;
      right: -10%;
      width: 50%;
      height: 50%;
      border-radius: 50%;
      background: radial-gradient(circle, rgba(255,255,255,0.2) 0%, rgba(255,255,255,0) 70%);
      animation: float 15s infinite ease-in-out reverse;
    }
  }

  .login-card {
    position: relative;
    z-index: 1;
    width: 400px;
    padding: 40px;
    border-radius: 20px;
    background: rgba(255, 255, 255, 0.9);
    backdrop-filter: blur(20px);
    box-shadow: 0 15px 35px rgba(0, 0, 0, 0.2);
    display: flex;
    flex-direction: column;
    align-items: center;
    animation: slideUp 0.8s cubic-bezier(0.2, 0.8, 0.2, 1);

    &.glass-effect {
      background: rgba(255, 255, 255, 0.85);
      border: 1px solid rgba(255, 255, 255, 0.5);
    }

    .header {
      text-align: center;
      margin-bottom: 30px;

      .logo-wrapper {
        width: 80px;
        height: 80px;
        margin: 0 auto 15px;
        background: #fff;
        border-radius: 20px;
        padding: 10px;
        box-shadow: 0 8px 16px rgba(0,0,0,0.1);
        display: flex;
        justify-content: center;
        align-items: center;
        transition: transform 0.3s ease;

        &:hover {
          transform: scale(1.05) rotate(5deg);
        }

        .logo-img {
          width: 100%;
          height: 100%;
          object-fit: contain;
        }
      }

      .app-title {
        font-size: 28px;
        font-weight: 800;
        color: #333;
        margin: 0 0 5px;
        background: linear-gradient(45deg, #333, #666);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
      }

      .app-subtitle {
        font-size: 14px;
        color: #888;
        margin: 0;
        letter-spacing: 2px;
      }
    }

    .custom-tabs {
      width: 100%;
      
      :deep(.n-tabs-rail) {
        border-radius: 30px;
        background-color: rgba(0, 0, 0, 0.05);
      }
      
      :deep(.n-tabs-tab) {
        border-radius: 30px;
        font-weight: bold;
        transition: all 0.3s ease;
        
        &.n-tabs-tab--active {
          color: #333;
          font-size: 16px;
        }
      }
    }

    .form-container {
      width: 100%;
      padding: 0 10px;
      box-sizing: border-box;
      
      :deep(.n-input) {
        background-color: rgba(245, 247, 250, 0.8);
        border: 1px solid transparent;
        transition: all 0.3s ease;
        
        &:hover, &:focus-within {
          background-color: #fff;
          border-color: var(--n-caret-color);
          box-shadow: 0 0 0 3px rgba(100, 100, 255, 0.1);
        }
        
        .n-input__prefix {
          margin-right: 12px;
          color: #999;
        }
      }

      .input-icon {
        font-size: 20px;
      }

      .options-row {
        display: flex;
        justify-content: space-between;
        width: 100%;
        padding: 0 5px;
      }

      .submit-btn {
        width: 100%;
        height: 44px;
        font-size: 16px;
        font-weight: bold;
        letter-spacing: 2px;
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        border: none;
        box-shadow: 0 4px 15px rgba(100, 100, 255, 0.4);
        transition: all 0.3s ease;

        &:hover {
          transform: translateY(-2px);
          box-shadow: 0 8px 20px rgba(100, 100, 255, 0.5);
        }

        &:active {
          transform: translateY(1px);
        }
      }
    }
  }
}

@keyframes float {
  0% { transform: translateY(0px) rotate(0deg); }
  50% { transform: translateY(-20px) rotate(5deg); }
  100% { transform: translateY(0px) rotate(0deg); }
}

@keyframes slideUp {
  0% { opacity: 0; transform: translateY(50px); }
  100% { opacity: 1; transform: translateY(0); }
}

// 移动端适配
@media (max-width: 480px) {
  .login-container {
    .login-card {
      width: 90%;
      padding: 20px;
      margin: 20px;
    }
  }
}
</style>
