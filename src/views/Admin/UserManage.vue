<template>
  <div class="user-manage-container">
    <!-- 动态背景层 -->
    <div class="ambient-background">
      <div class="gradient-blob blob-1"></div>
      <div class="gradient-blob blob-2"></div>
      <div class="gradient-blob blob-3"></div>
      <div class="grain-overlay"></div>
    </div>

    <div class="content-wrapper">
      <!-- 杂志风格标题区 -->
      <header class="editorial-header">
        <div class="header-left">
          <button class="back-btn" @click="router.back()">
            <n-icon :component="Left" />
          </button>
          <div class="title-stack">
            <span class="kicker">Administration</span>
            <h1 class="headline">用户管理</h1>
            <p class="deck">管理全站注册用户与权限</p>
          </div>
        </div>
        <div class="header-right">
          <div class="stat-pill">
            <span class="stat-number">{{ pagination.itemCount }}</span>
            <span class="stat-label">位用户</span>
          </div>
        </div>
      </header>

      <!-- 搜索控制面板 -->
      <section class="control-panel">
        <div class="search-wrapper">
          <div class="search-icon">
            <n-icon :component="Search" />
          </div>
          <input
            v-model="searchText"
            type="text"
            placeholder="搜索用户名..."
            class="search-input"
            @keyup.enter="handleSearch"
          />
          <button class="search-action" @click="handleSearch">
            搜索
          </button>
        </div>

        <div class="view-toggles">
          <button
            class="view-btn"
            :class="{ active: viewMode === 'table' }"
            @click="viewMode = 'table'"
          >
            <n-icon :component="List" />
          </button>
          <button
            class="view-btn"
            :class="{ active: viewMode === 'grid' }"
            @click="viewMode = 'grid'"
          >
            <n-icon :component="GridFour" />
          </button>
        </div>
      </section>

      <!-- 数据概览卡片 -->
      <section class="metrics-bar">
        <div class="metric-item">
          <span class="metric-value">{{ pagination.page }}</span>
          <span class="metric-label">当前页</span>
        </div>
        <div class="metric-divider"></div>
        <div class="metric-item">
          <span class="metric-value">{{ pagination.pageSize }}</span>
          <span class="metric-label">每页显示</span>
        </div>
        <div class="metric-divider"></div>
        <div class="metric-item">
          <span class="metric-value">{{ userList.length }}</span>
          <span class="metric-label">本页用户</span>
        </div>
      </section>

      <!-- 桌面端表格视图 -->
      <Transition name="fade-scale">
        <div v-if="viewMode === 'table' && !isMobile" class="table-view">
          <div class="table-container">
            <n-data-table
              :columns="columns"
              :data="userList"
              :loading="loading"
              :row-key="(row: any) => row.id"
              :row-class-name="'user-row'"
              :scroll-x="900"
            />
          </div>
        </div>
      </Transition>

      <!-- 网格/卡片视图 -->
      <Transition name="fade-scale">
        <div v-if="viewMode === 'grid' || isMobile" class="grid-view">
          <div class="user-cards">
            <div
              v-for="user in userList"
              :key="user.id"
              class="user-card"
              :class="{ 'is-admin': user.user_group === 'admin' }"
            >
              <div class="card-accent"></div>
              <div class="card-content">
                <div class="user-header">
                  <div class="avatar-ring">
                    <div class="avatar">{{ user.username.charAt(0).toUpperCase() }}</div>
                  </div>
                  <div class="user-meta">
                    <h3 class="user-name">{{ user.username }}</h3>
                    <span class="user-id">#{{ user.id }}</span>
                  </div>
                  <div class="role-tag" :class="user.user_group">
                    {{ user.user_group === 'admin' ? '管理员' : '用户' }}
                  </div>
                </div>

                <div class="user-details">
                  <div class="detail-item">
                    <n-icon :component="Mail" />
                    <span>{{ user.email || '未绑定邮箱' }}</span>
                  </div>
                  <div class="detail-item">
                    <n-icon :component="Time" />
                    <span>{{ user.last_login ? formatDate(user.last_login) : '从未登录' }}</span>
                  </div>
                </div>

                <div class="card-actions">
                  <button
                    class="action-btn secondary"
                    @click="confirmRoleChange(user)"
                  >
                    <n-icon :component="Permissions" />
                    {{ user.user_group === 'admin' ? '降级' : '提升' }}
                  </button>
                  <button
                    class="action-btn danger"
                    @click="confirmDelete(user)"
                  >
                    <n-icon :component="Delete" />
                    删除
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </Transition>

      <!-- 分页器 -->
      <div class="pagination-wrapper">
        <n-pagination
          v-model:page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :item-count="pagination.itemCount"
          :show-size-picker="true"
          :page-sizes="pagination.pageSizes"
          @update:page="fetchUsers"
          @update:page-size="handlePageSizeChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, h, onMounted, computed, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import { useMessage, NPopconfirm, NDataTable, NPagination, NIcon, type DataTableColumns } from 'naive-ui';
import { Left, Search, Mail, Time, Permissions, Delete, List, GridFour } from '@icon-park/vue-next';
import { adminGetAllUsers, adminUpdateUserRole, adminDeleteUser, type UserInfo } from '@/api/user';
import { ResultCode } from "@/utils/request";

const router = useRouter();
const message = useMessage();

// 响应式窗口大小检测
const windowWidth = ref(window.innerWidth);
const isMobile = computed(() => windowWidth.value < 768);

const viewMode = ref<'table' | 'grid'>('table');

const handleResize = () => {
  windowWidth.value = window.innerWidth;
};

const loading = ref(false);
const searchText = ref('');
const userList = ref<UserInfo[]>([]);

const pagination = reactive({
  page: 1,
  pageSize: 12,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [12, 24, 48],
});

const formatDate = (date: string) => {
  const d = new Date(date);
  return d.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' });
};

const handlePageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize;
  pagination.page = 1;
  fetchUsers();
};

const fetchUsers = async () => {
  loading.value = true;
  try {
    const res = await adminGetAllUsers({
      page: pagination.page,
      limit: pagination.pageSize,
      query: searchText.value
    });
    if (res.code === ResultCode.SUCCESS) {
      userList.value = res.data.list;
      pagination.itemCount = res.data.total;
    } else {
      message.error(res.message || '获取用户列表失败');
    }
  } catch (error) {
    message.error('获取用户列表失败');
  } finally {
    loading.value = false;
  }
};

const columns: DataTableColumns<UserInfo> = [
  {
    title: '用户',
    key: 'username',
    width: 200,
    render(row: UserInfo) {
      return h('div', { class: 'table-user-cell' }, [
        h('div', { class: 'table-avatar' }, row.username.charAt(0).toUpperCase()),
        h('div', { class: 'table-user-info' }, [
          h('span', { class: 'table-username' }, row.username),
          h('span', { class: 'table-user-id' }, `#${row.id}`)
        ])
      ]);
    }
  },
  {
    title: '角色',
    key: 'user_group',
    width: 120,
    align: 'center',
    render(row: UserInfo) {
      return h('span', {
        class: `table-role ${row.user_group}`
      }, row.user_group === 'admin' ? '管理员' : '普通用户');
    }
  },
  {
    title: '邮箱',
    key: 'email',
    width: 240,
    render(row: UserInfo) {
      return h('div', { class: 'table-email' }, [
        h(NIcon, { component: Mail, size: 14 }),
        h('span', row.email || '未绑定邮箱')
      ]);
    }
  },
  {
    title: '上次登录',
    key: 'last_login',
    width: 180,
    render(row: UserInfo) {
      return h('div', { class: 'table-time' }, [
        h(NIcon, { component: Time, size: 14 }),
        h('span', row.last_login ? formatDate(row.last_login) : '从未登录')
      ]);
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 160,
    fixed: 'right',
    align: 'center',
    render(row: UserInfo) {
      return h('div', { class: 'table-actions' }, [
        h(
          NPopconfirm,
          {
            onPositiveClick: () => handleRoleChange(row),
            positiveText: '确认',
            negativeText: '取消'
          },
          {
            trigger: () => h('button', {
              class: 'table-btn role',
              title: row.user_group === 'admin' ? '降级' : '提升'
            }, h(NIcon, { component: Permissions, size: 16 })),
            default: () => `确定要${row.user_group === 'admin' ? '降级' : '提升'}该用户？`
          }
        ),
        h(
          NPopconfirm,
          {
            onPositiveClick: () => handleDelete(row),
            positiveText: '确认删除',
            negativeText: '取消'
          },
          {
            trigger: () => h('button', {
              class: 'table-btn delete',
              title: '删除用户'
            }, h(NIcon, { component: Delete, size: 16 })),
            default: () => '确定删除该用户？此操作不可逆！'
          }
        )
      ]);
    }
  }
];

onMounted(() => {
  fetchUsers();
  window.addEventListener('resize', handleResize);
});

onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
});

const handleSearch = () => {
  pagination.page = 1;
  fetchUsers();
};

const handleRoleChange = async (row: UserInfo) => {
  const newRole = row.user_group === 'admin' ? 'user' : 'admin';
  try {
    const res = await adminUpdateUserRole(row.id, newRole);
    if (res.code === ResultCode.SUCCESS) {
      message.success('修改权限成功');
      fetchUsers();
    } else {
      message.error(res.message || '修改权限失败');
    }
  } catch (error) {
    message.error('修改权限失败');
  }
};

const handleDelete = async (row: UserInfo) => {
  try {
    const res = await adminDeleteUser(row.id);
    if (res.code === ResultCode.SUCCESS) {
      message.success('删除用户成功');
      if (userList.value.length === 1 && pagination.page > 1) {
        pagination.page--;
      }
      fetchUsers();
    } else {
      message.error(res.message || '删除用户失败');
    }
  } catch (error) {
    message.error('删除用户失败');
  }
};

const confirmRoleChange = (row: UserInfo) => {
  const target = row.user_group === 'admin' ? '普通用户' : '管理员';
  if (window.confirm(`确定要将用户 ${row.username} 设为${target}吗？`)) {
    handleRoleChange(row);
  }
};

const confirmDelete = (row: UserInfo) => {
  if (window.confirm('确定要删除该用户吗？此操作不可逆！')) {
    handleDelete(row);
  }
};
</script>

<style lang="scss" scoped>
// 导入特色字体
@import url('https://fonts.googleapis.com/css2?family=Playfair+Display:wght@600;700&family=Plus+Jakarta+Sans:wght@400;500;600;700&display=swap');

:root {
  // 编辑风格配色 - 温暖米色调
  --bg-primary: #faf8f5;
  --bg-secondary: #f5f2ed;
  --bg-tertiary: #ebe7e0;
  --text-primary: #1a1a1a;
  --text-secondary: #666666;
  --text-muted: #999999;
  --accent-coral: #e07a5f;
  --accent-teal: #3d8b8b;
  --accent-gold: #d4a574;
  --accent-ink: #2c3e50;
  --admin-purple: #7c6fae;
  --user-blue: #5b8db8;
  --danger-red: #c75b5b;
  --shadow-soft: 0 4px 20px rgba(0, 0, 0, 0.06);
  --shadow-medium: 0 8px 30px rgba(0, 0, 0, 0.1);
  --shadow-deep: 0 12px 40px rgba(0, 0, 0, 0.14);
}

.user-manage-container {
  min-height: 100vh;
  background: var(--bg-primary);
  font-family: 'Plus Jakarta Sans', sans-serif;
  position: relative;
  overflow-x: hidden;
}

// 动态背景
.ambient-background {
  position: fixed;
  inset: 0;
  z-index: 0;
  overflow: hidden;
  pointer-events: none;
}

.gradient-blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.5;
  animation: blob-float 20s ease-in-out infinite;

  &.blob-1 {
    width: 600px;
    height: 600px;
    background: linear-gradient(135deg, rgba(224, 122, 95, 0.25), rgba(212, 165, 116, 0.2));
    top: -200px;
    right: -100px;
    animation-delay: 0s;
  }

  &.blob-2 {
    width: 500px;
    height: 500px;
    background: linear-gradient(135deg, rgba(61, 139, 139, 0.2), rgba(91, 141, 184, 0.15));
    bottom: -150px;
    left: -150px;
    animation-delay: -7s;
  }

  &.blob-3 {
    width: 400px;
    height: 400px;
    background: linear-gradient(135deg, rgba(124, 111, 174, 0.15), rgba(224, 122, 95, 0.1));
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    animation-delay: -14s;
  }
}

@keyframes blob-float {
  0%, 100% { transform: translate(0, 0) scale(1); }
  33% { transform: translate(30px, -30px) scale(1.05); }
  66% { transform: translate(-20px, 20px) scale(0.95); }
}

.grain-overlay {
  position: absolute;
  inset: 0;
  opacity: 0.03;
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noise'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noise)'/%3E%3C/svg%3E");
}

// 内容包装器
.content-wrapper {
  position: relative;
  z-index: 1;
  max-width: 1400px;
  margin: 0 auto;
  padding: 48px 32px;

  @media (max-width: 768px) {
    padding: 24px 16px;
  }
}

// 编辑风格标题区
.editorial-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 40px;
  padding-bottom: 24px;
  border-bottom: 1px solid var(--bg-tertiary);

  @media (max-width: 768px) {
    flex-direction: column;
    align-items: flex-start;
    gap: 20px;
    margin-bottom: 28px;
  }
}

.header-left {
  display: flex;
  align-items: flex-start;
  gap: 20px;
}

.back-btn {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  border: 1px solid var(--bg-tertiary);
  background: var(--bg-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
  color: var(--text-secondary);

  &:hover {
    background: var(--bg-tertiary);
    transform: translateX(-2px);
    color: var(--text-primary);
  }
}

.title-stack {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.kicker {
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.15em;
  color: var(--accent-coral);
}

.headline {
  font-family: 'Playfair Display', serif;
  font-size: 42px;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1.1;
  margin: 0;

  @media (max-width: 768px) {
    font-size: 32px;
  }
}

.deck {
  font-size: 14px;
  color: var(--text-secondary);
  margin: 8px 0 0;
  max-width: 300px;
  line-height: 1.5;
}

.stat-pill {
  display: flex;
  align-items: baseline;
  gap: 6px;
  padding: 12px 20px;
  background: linear-gradient(135deg, var(--accent-ink), #3d4f5f);
  border-radius: 100px;
  color: white;

  .stat-number {
    font-size: 24px;
    font-weight: 700;
  }

  .stat-label {
    font-size: 13px;
    opacity: 0.8;
  }
}

// 控制面板
.control-panel {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  gap: 16px;

  @media (max-width: 768px) {
    flex-direction: column;
    align-items: stretch;
  }
}

.search-wrapper {
  flex: 1;
  max-width: 400px;
  display: flex;
  align-items: center;
  background: var(--bg-secondary);
  border: 1px solid var(--bg-tertiary);
  border-radius: 14px;
  padding: 4px;
  transition: all 0.3s ease;

  &:focus-within {
    border-color: var(--accent-teal);
    box-shadow: 0 0 0 3px rgba(61, 139, 139, 0.1);
  }

  @media (max-width: 768px) {
    max-width: none;
  }
}

.search-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-muted);
}

.search-input {
  flex: 1;
  border: none;
  background: transparent;
  padding: 10px 8px;
  font-size: 14px;
  color: var(--text-primary);
  outline: none;

  &::placeholder {
    color: var(--text-muted);
  }
}

.search-action {
  padding: 10px 18px;
  background: var(--accent-ink);
  color: white;
  border: none;
  border-radius: 10px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    background: #3d4f5f;
    transform: translateY(-1px);
  }
}

.view-toggles {
  display: flex;
  gap: 8px;
}

.view-btn {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  border: 1px solid var(--bg-tertiary);
  background: var(--bg-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
  color: var(--text-muted);

  &:hover, &.active {
    background: var(--accent-ink);
    color: white;
    border-color: var(--accent-ink);
  }

  &.active {
    box-shadow: var(--shadow-soft);
  }
}

// 数据指标栏
.metrics-bar {
  display: flex;
  align-items: center;
  gap: 24px;
  padding: 16px 24px;
  background: var(--bg-secondary);
  border-radius: 16px;
  margin-bottom: 24px;

  @media (max-width: 768px) {
    justify-content: space-around;
    gap: 0;
  }
}

.metric-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.metric-value {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  font-family: 'Playfair Display', serif;
}

.metric-label {
  font-size: 12px;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.metric-divider {
  width: 1px;
  height: 32px;
  background: var(--bg-tertiary);

  @media (max-width: 768px) {
    display: none;
  }
}

// 表格视图样式
.table-view {
  margin-bottom: 24px;
}

.table-container {
  background: var(--bg-secondary);
  border-radius: 20px;
  padding: 8px;
  box-shadow: var(--shadow-soft);
}

:deep(.n-data-table) {
  --n-th-font-weight: 600 !important;
  --n-th-text-color: var(--text-secondary) !important;
  --n-border-color: transparent !important;

  .n-data-table-th {
    background: transparent !important;
    padding: 16px 20px !important;
    font-size: 12px;
    text-transform: uppercase;
    letter-spacing: 0.08em;
  }

  .n-data-table-td {
    background: transparent !important;
    padding: 16px 20px !important;
    border-bottom: 1px solid var(--bg-tertiary) !important;
  }

  .n-data-table-tr:last-child .n-data-table-td {
    border-bottom: none !important;
  }

  .n-data-table-tr:hover .n-data-table-td {
    background: rgba(255, 255, 255, 0.5) !important;
  }
}

.table-user-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.table-avatar {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  background: linear-gradient(135deg, var(--accent-coral), var(--accent-gold));
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 600;
  font-size: 16px;
}

.table-user-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.table-username {
  font-weight: 600;
  color: var(--text-primary);
  font-size: 14px;
}

.table-user-id {
  font-size: 12px;
  color: var(--text-muted);
  font-family: monospace;
}

.table-role {
  display: inline-block;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;

  &.admin {
    background: rgba(124, 111, 174, 0.15);
    color: var(--admin-purple);
  }

  &.user {
    background: rgba(91, 141, 184, 0.15);
    color: var(--user-blue);
  }
}

.table-email, .table-time {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--text-secondary);
  font-size: 13px;

  .n-icon {
    color: var(--text-muted);
  }
}

.table-actions {
  display: flex;
  gap: 8px;
  justify-content: center;
}

.table-btn {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  border: none;
  background: transparent;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
  color: var(--text-muted);

  &.role:hover {
    background: rgba(124, 111, 174, 0.1);
    color: var(--admin-purple);
  }

  &.delete:hover {
    background: rgba(199, 91, 91, 0.1);
    color: var(--danger-red);
  }
}

// 网格/卡片视图
.grid-view {
  margin-bottom: 24px;
}

.user-cards {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
  gap: 20px;

  @media (max-width: 768px) {
    grid-template-columns: 1fr;
    gap: 16px;
  }
}

.user-card {
  position: relative;
  background: var(--bg-secondary);
  border-radius: 20px;
  overflow: hidden;
  transition: all 0.3s ease;
  box-shadow: var(--shadow-soft);

  &:hover {
    transform: translateY(-4px);
    box-shadow: var(--shadow-medium);
  }

  &.is-admin {
    .card-accent {
      background: linear-gradient(180deg, var(--admin-purple), transparent);
    }

    .avatar-ring {
      background: linear-gradient(135deg, var(--admin-purple), #a89bc9);
    }
  }
}

.card-accent {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(180deg, var(--accent-coral), transparent);
}

.card-content {
  padding: 20px;
}

.user-header {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-bottom: 16px;
}

.avatar-ring {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--accent-coral), var(--accent-gold));
  padding: 3px;
  flex-shrink: 0;
}

.avatar {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background: var(--bg-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  font-weight: 700;
  color: var(--text-primary);
}

.user-meta {
  flex: 1;
  min-width: 0;
}

.user-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.user-id {
  font-size: 12px;
  color: var(--text-muted);
  font-family: monospace;
}

.role-tag {
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;

  &.admin {
    background: rgba(124, 111, 174, 0.15);
    color: var(--admin-purple);
  }

  &.user {
    background: rgba(91, 141, 184, 0.15);
    color: var(--user-blue);
  }
}

.user-details {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 16px;
}

.detail-item {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 13px;
  color: var(--text-secondary);

  .n-icon {
    color: var(--text-muted);
    flex-shrink: 0;
  }

  span {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}

.card-actions {
  display: flex;
  gap: 10px;
}

.action-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 10px 16px;
  border-radius: 10px;
  border: none;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;

  &.secondary {
    background: var(--bg-tertiary);
    color: var(--text-secondary);

    &:hover {
      background: var(--accent-ink);
      color: white;
    }
  }

  &.danger {
    background: rgba(199, 91, 91, 0.1);
    color: var(--danger-red);

    &:hover {
      background: var(--danger-red);
      color: white;
    }
  }
}

// 分页器
.pagination-wrapper {
  display: flex;
  justify-content: center;
  padding-top: 8px;
}

// 过渡动画
.fade-scale-enter-active,
.fade-scale-leave-active {
  transition: all 0.3s ease;
}

.fade-scale-enter-from,
.fade-scale-leave-to {
  opacity: 0;
  transform: scale(0.98);
}
</style>
