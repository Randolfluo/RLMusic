<template>
  <div class="login">
    <div class="title">
      <img src="/images/logo/favicon.png" alt="logo" />
      <span>登录云音乐</span>
    </div>
    <n-tabs
      animated
      class="content"
      type="segment"
      justify-content="space-evenly"
      :pane-style="{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        paddingTop: '30px',
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
        >
          <n-form-item path="username">
            <n-input
              placeholder="请输入用户名"
              v-model:value="loginForm.username"
              @keydown.enter.prevent
            >
              <template #prefix>
                <n-icon :component="PersonRound" />
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
            >
              <template #prefix>
                <n-icon :component="PasswordRound" />
              </template>
            </n-input>
          </n-form-item>
          <n-form-item>
            <n-button style="width: 100%" type="primary" @click="handleLogin" :loading="loading">
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
        >
          <n-form-item path="username">
            <n-input
              placeholder="用户名"
              v-model:value="registerForm.username"
            >
              <template #prefix>
                <n-icon :component="PersonRound" />
              </template>
            </n-input>
          </n-form-item>
          <n-form-item path="password">
            <n-input
              type="password"
              show-password-on="click"
              placeholder="密码 (4-20位)"
              v-model:value="registerForm.password"
            >
              <template #prefix>
                <n-icon :component="PasswordRound" />
              </template>
            </n-input>
          </n-form-item>
          <n-form-item path="email">
            <n-input
              placeholder="邮箱"
              v-model:value="registerForm.email"
            >
              <template #prefix>
                <n-icon :component="EmailRound" />
              </template>
            </n-input>
          </n-form-item>
          <n-form-item>
            <n-button style="width: 100%" type="primary" @click="handleRegister" :loading="loading">
              注册
            </n-button>
          </n-form-item>
        </n-form>
      </n-tab-pane>
    </n-tabs>
  </div>
</template>

<script setup>
import { ref, reactive } from "vue";
import { useRouter } from "vue-router";
import { userStore } from "@/store";
import { login, register } from "@/api/login";
import { PersonRound, PasswordRound, EmailRound } from "@vicons/material";
import { useMessage } from "naive-ui";
import { ResultCode } from "@/utils/request";

const router = useRouter();
const user = userStore();
const message = useMessage();
const loading = ref(false);

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
          localStorage.setItem("token", res.data.token); 
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
.login {
  margin-top: 40px;
  display: flex;
  flex-direction: column;
  align-items: center;
  .title {
    display: flex;
    flex-direction: column;
    align-items: center;
    img {
      width: 80px;
      height: 80px;
      margin-bottom: 20px;
    }
    span {
      font-size: 26px;
      font-weight: bold;
    }
  }
  .content {
    width: 300px;
    margin-top: 30px;
    .form-container {
      width: 100%;
      padding: 0 4px;
      box-sizing: border-box;
    }
    :deep(.n-input) {
      .n-input__prefix {
        margin-right: 8px;
      }
    }
  }
}
</style>
