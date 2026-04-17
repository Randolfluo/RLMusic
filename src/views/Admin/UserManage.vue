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
      <header class="editorial-header" data-animate="fade-up">
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

      <!-- 数据概览卡片 -->
      <section class="stats-row" data-animate="fade-up" data-delay="1">
        <div class="stat-card">
          <div class="stat-icon coral"><n-icon :component="People" /></div>
          <div class="stat-info">
            <span class="stat-value">{{ pagination.itemCount }}</span>
            <span class="stat-name">总用户数</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon purple"><n-icon :component="Permissions" /></div>
          <div class="stat-info">
            <span class="stat-value">{{ adminCount }}</span>
            <span class="stat-name">管理员</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon teal"><n-icon :component="User" /></div>
          <div class="stat-info">
            <span class="stat-value">{{ userCount }}</span>
            <span class="stat-name">普通用户</span>
          </div>
        </div>
      </section>

      <!-- 搜索控制面板 -->
      <section class="control-panel" data-animate="fade-up" data-delay="2">
        <div class="search-wrapper">
          <div class="search-icon">
            <n-icon :component="Search" />
          </div>
          <input
            v-model="searchText"
            type="text"
            placeholder="搜索用户名或邮箱..."
            class="search-input"
            @keyup.enter="handleSearch"
          />
          <button class="search-action" @click="handleSearch">
            搜索
          </button>
        </div>

        <div class="bulk-toggle" v-if="!isMobile">
          <label class="toggle-label">
            <input type="checkbox" v-model="batchMode" />
            <span>批量模式</span>
          </label>
        </div>
      </section>

      <!-- 桌面端自定义表格 -->
      <Transition name="fade-scale">
        <div v-if="!isMobile" class="custom-table-view" data-animate="fade-up" data-delay="3">
          <div class="table-header" :style="tableGridStyle">
            <div class="th checkbox-col" v-if="batchMode">
              <input type="checkbox" :checked="isAllSelected" @change="toggleSelectAll" />
            </div>
            <div class="th user-col">用户</div>
            <div class="th role-col">角色</div>
            <div class="th email-col">邮箱</div>
            <div class="th login-col">上次登录</div>
            <div class="th action-col">操作</div>
          </div>
          <div class="table-body">
            <div
              v-for="(user, index) in userList"
              :key="user.id"
              class="table-row"
              :class="{ 'is-admin': user.user_group === 'admin', 'is-selected': selectedIds.has(user.id) }"
              :style="{ ...tableGridStyle, animationDelay: `${index * 40}ms` }"
            >
              <div class="td checkbox-col" v-if="batchMode">
                <input type="checkbox" :checked="selectedIds.has(user.id)" @change="toggleSelect(user.id)" />
              </div>
              <div class="td user-col">
                <div class="table-user-cell">
                  <div class="table-avatar">{{ user.username.charAt(0).toUpperCase() }}</div>
                  <div class="table-user-info">
                    <span class="table-username">{{ user.username }}</span>
                    <span class="table-user-id">#{{ user.id }}</span>
                  </div>
                </div>
              </div>
              <div class="td role-col">
                <span class="table-role" :class="user.user_group">{{ user.user_group === 'admin' ? '管理员' : '普通用户' }}</span>
              </div>
              <div class="td email-col">
                <div class="table-contact">
                  <n-icon :component="Mail" size="14" />
                  <span>{{ user.email || '未绑定邮箱' }}</span>
                </div>
              </div>
              <div class="td login-col">
                <div class="table-contact">
                  <n-icon :component="Time" size="14" />
                  <span>{{ user.last_login ? formatDate(user.last_login) : '从未登录' }}</span>
                </div>
              </div>
              <div class="td action-col">
                <div class="table-actions">
                  <n-popconfirm @positive-click="handleRoleChange(user)" positive-text="确认" negative-text="取消">
                    <template #trigger>
                      <button class="table-btn role" :title="user.user_group === 'admin' ? '降级' : '提升'">
                        <n-icon :component="Permissions" size="16" />
                      </button>
                    </template>
                    确定要{{ user.user_group === 'admin' ? '降级' : '提升' }}该用户？
                  </n-popconfirm>
                  <n-popconfirm @positive-click="handleDelete(user)" positive-text="确认删除" negative-text="取消">
                    <template #trigger>
                      <button class="table-btn delete" title="删除用户">
                        <n-icon :component="Delete" size="16" />
                      </button>
                    </template>
                    确定删除该用户？此操作不可逆！
                  </n-popconfirm>
                </div>
              </div>
            </div>
          </div>
        </div>
      </Transition>

      <!-- 移动端/网格卡片视图 -->
      <Transition name="fade-scale">
        <div v-if="isMobile" class="grid-view">
          <div class="user-cards">
            <div
              v-for="(user, index) in userList"
              :key="user.id"
              class="user-card"
              :class="{ 'is-admin': user.user_group === 'admin', 'is-selected': selectedIds.has(user.id) }"
              :style="{ animationDelay: `${index * 60}ms` }"
              @click="batchMode && toggleSelect(user.id)"
            >
              <div class="card-shimmer"></div>
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

                <div class="card-actions" @click.stop>
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
              <div class="selection-indicator" v-if="batchMode">
                <div class="check-circle" :class="{ checked: selectedIds.has(user.id) }">
                  <n-icon :component="Check" v-if="selectedIds.has(user.id)" />
                </div>
              </div>
            </div>
          </div>
        </div>
      </Transition>

      <!-- 空状态 -->
      <div v-if="userList.length === 0 && !loading" class="empty-state">
        <div class="empty-icon"><n-icon :component="SearchEmpty" size="48" /></div>
        <p>未找到匹配的用户</p>
      </div>

      <!-- 分页器 -->
      <div class="pagination-wrapper" v-if="userList.length > 0">
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

    <!-- 批量操作浮动栏 -->
    <Transition name="slide-up">
      <div v-if="selectedIds.size > 0" class="batch-bar">
        <div class="batch-info">
          <span class="batch-count">{{ selectedIds.size }}</span>
          <span>位用户已选择</span>
        </div>
        <div class="batch-actions">
          <button class="batch-btn role" @click="batchChangeRole">
            <n-icon :component="Permissions" />
            切换权限
          </button>
          <button class="batch-btn delete" @click="batchDelete">
            <n-icon :component="Delete" />
            批量删除
          </button>
          <button class="batch-btn clear" @click="selectedIds.clear()">取消</button>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import { useMessage, NPagination, NPopconfirm, NIcon } from 'naive-ui';
import { Left, Search, Mail, Time, Permissions, Delete, People, User, Check } from '@icon-park/vue-next';
import { adminGetAllUsers, adminUpdateUserRole, adminDeleteUser, type UserInfo } from '@/api/user';
import { ResultCode } from "@/utils/request";

const SearchEmpty = Search;

const router = useRouter();
const message = useMessage();

const windowWidth = ref(window.innerWidth);
const isMobile = computed(() => windowWidth.value < 768);

const handleResize = () => {
  windowWidth.value = window.innerWidth;
};

const loading = ref(false);
const searchText = ref('');
const userList = ref<UserInfo[]>([]);
const batchMode = ref(false);
const selectedIds = ref<Set<number>>(new Set());

const pagination = reactive({
  page: 1,
  pageSize: 12,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [12, 24, 48],
});

const adminCount = computed(() => userList.value.filter(u => u.user_group === 'admin').length);
const userCount = computed(() => userList.value.filter(u => u.user_group !== 'admin').length);
const isAllSelected = computed(() => userList.value.length > 0 && userList.value.every(u => selectedIds.value.has(u.id)));

const tableGridStyle = computed(() => {
  const cols = batchMode.value
    ? '48px minmax(180px, 2fr) 100px minmax(140px, 1.5fr) 140px 110px'
    : 'minmax(180px, 2fr) 100px minmax(140px, 1.5fr) 140px 110px';
  return { gridTemplateColumns: cols };
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
  selectedIds.value.clear();
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

const handleSearch = () => {
  pagination.page = 1;
  fetchUsers();
};

const toggleSelect = (id: number) => {
  const s = new Set(selectedIds.value);
  if (s.has(id)) s.delete(id);
  else s.add(id);
  selectedIds.value = s;
};

const toggleSelectAll = () => {
  if (isAllSelected.value) {
    selectedIds.value = new Set();
  } else {
    selectedIds.value = new Set(userList.value.map(u => u.id));
  }
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

const batchChangeRole = async () => {
  const ids = Array.from(selectedIds.value);
  if (!ids.length) return;
  const first = userList.value.find(u => u.id === ids[0]);
  const newRole = first?.user_group === 'admin' ? 'user' : 'admin';
  if (!window.confirm(`确定要批量切换 ${ids.length} 位用户的权限为${newRole === 'admin' ? '管理员' : '普通用户'}吗？`)) return;
  let success = 0;
  for (const id of ids) {
    try {
      const res = await adminUpdateUserRole(id, newRole);
      if (res.code === ResultCode.SUCCESS) success++;
    } catch (e) { /* ignore */ }
  }
  message.success(`成功修改 ${success} 位用户的权限`);
  selectedIds.value.clear();
  fetchUsers();
};

const batchDelete = async () => {
  const ids = Array.from(selectedIds.value);
  if (!ids.length) return;
  if (!window.confirm(`确定要批量删除 ${ids.length} 位用户吗？此操作不可逆！`)) return;
  let success = 0;
  for (const id of ids) {
    try {
      const res = await adminDeleteUser(id);
      if (res.code === ResultCode.SUCCESS) success++;
    } catch (e) { /* ignore */ }
  }
  message.success(`成功删除 ${success} 位用户`);
  selectedIds.value.clear();
  if (userList.value.length === ids.length && pagination.page > 1) pagination.page--;
  fetchUsers();
};

onMounted(() => {
  fetchUsers();
  window.addEventListener('resize', handleResize);
});

onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
});
</script>

<style lang="scss" scoped>
@import url('https://fonts.googleapis.com/css2?family=DM+Serif+Display:ital@0;1&family=Outfit:wght@300;400;500;600;700&display=swap');

:root {
  --bg-cream: #fdfbf7;
  --bg-warm: #f7f3ec;
  --bg-card: #ffffff;
  --text-ink: #1f140e;
  --text-body: #5c4d41;
  --text-muted: #998e84;
  --accent-coral: #c75b39;
  --accent-teal: #2d6a4f;
  --accent-purple: #7c6fae;
  --accent-gold: #c9a227;
  --border-light: rgba(31, 20, 14, 0.08);
  --shadow-sm: 0 2px 8px rgba(31, 20, 14, 0.04);
  --shadow-md: 0 8px 24px rgba(31, 20, 14, 0.08);
  --shadow-lg: 0 16px 40px rgba(31, 20, 14, 0.12);
}

.user-manage-container {
  min-height: 100vh;
  background: var(--bg-cream);
  font-family: 'Outfit', sans-serif;
  position: relative;
  overflow-x: hidden;
  padding-bottom: 120px;
}

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
  filter: blur(90px);
  opacity: 0.45;
  animation: blob-float 22s ease-in-out infinite;

  &.blob-1 {
    width: 650px;
    height: 650px;
    background: linear-gradient(135deg, rgba(199, 91, 57, 0.22), rgba(201, 162, 39, 0.14));
    top: -220px;
    right: -120px;
    animation-delay: 0s;
  }

  &.blob-2 {
    width: 550px;
    height: 550px;
    background: linear-gradient(135deg, rgba(45, 106, 79, 0.18), rgba(124, 111, 174, 0.12));
    bottom: -180px;
    left: -180px;
    animation-delay: -8s;
  }

  &.blob-3 {
    width: 450px;
    height: 450px;
    background: linear-gradient(135deg, rgba(124, 111, 174, 0.14), rgba(199, 91, 57, 0.1));
    top: 45%;
    left: 55%;
    transform: translate(-50%, -50%);
    animation-delay: -16s;
  }
}

@keyframes blob-float {
  0%, 100% { transform: translate(0, 0) scale(1); }
  33% { transform: translate(25px, -35px) scale(1.06); }
  66% { transform: translate(-15px, 25px) scale(0.94); }
}

.grain-overlay {
  position: absolute;
  inset: 0;
  opacity: 0.025;
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noise'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.85' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noise)'/%3E%3C/svg%3E");
}

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

[data-animate="fade-up"] {
  opacity: 0;
  transform: translateY(24px);
  animation: fadeUp 0.7s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}
[data-animate="fade-up"][data-delay="1"] { animation-delay: 80ms; }
[data-animate="fade-up"][data-delay="2"] { animation-delay: 160ms; }
[data-animate="fade-up"][data-delay="3"] { animation-delay: 240ms; }

@keyframes fadeUp {
  to { opacity: 1; transform: translateY(0); }
}

.editorial-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 32px;
  padding-bottom: 24px;
  border-bottom: 1px solid var(--border-light);

  @media (max-width: 768px) {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
    margin-bottom: 24px;
  }
}

.header-left {
  display: flex;
  align-items: flex-start;
  gap: 18px;
}

.back-btn {
  width: 46px;
  height: 46px;
  border-radius: 14px;
  border: 1px solid var(--border-light);
  background: var(--bg-warm);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.25s ease;
  color: var(--text-body);

  &:hover {
    background: #fff;
    transform: translateX(-3px);
    color: var(--text-ink);
    box-shadow: var(--shadow-sm);
  }
}

.title-stack {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.kicker {
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.18em;
  color: var(--accent-coral);
}

.headline {
  font-family: 'DM Serif Display', serif;
  font-size: 44px;
  font-weight: 400;
  color: var(--text-ink);
  line-height: 1.05;
  margin: 0;

  @media (max-width: 768px) {
    font-size: 34px;
  }
}

.deck {
  font-size: 14px;
  color: var(--text-body);
  margin: 6px 0 0;
  max-width: 280px;
  line-height: 1.5;
}

.stat-pill {
  display: flex;
  align-items: baseline;
  gap: 6px;
  padding: 12px 22px;
  background: linear-gradient(135deg, #2c2420, #4a3f38);
  border-radius: 100px;
  color: white;
  box-shadow: var(--shadow-md);

  .stat-number {
    font-size: 26px;
    font-weight: 600;
  }

  .stat-label {
    font-size: 13px;
    opacity: 0.85;
  }
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  margin-bottom: 28px;

  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }
}

.stat-card {
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(12px);
  border-radius: 20px;
  padding: 18px 20px;
  display: flex;
  align-items: center;
  gap: 14px;
  border: 1px solid var(--border-light);
  box-shadow: var(--shadow-sm);
  transition: transform 0.25s ease, box-shadow 0.25s ease;

  &:hover {
    transform: translateY(-3px);
    box-shadow: var(--shadow-md);
  }
}

.stat-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;

  &.coral { background: rgba(199, 91, 57, 0.12); color: var(--accent-coral); }
  &.purple { background: rgba(124, 111, 174, 0.14); color: var(--accent-purple); }
  &.teal { background: rgba(45, 106, 79, 0.12); color: var(--accent-teal); }
}

.stat-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.stat-value {
  font-size: 22px;
  font-weight: 600;
  color: var(--text-ink);
}

.stat-name {
  font-size: 12px;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.control-panel {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 22px;
  gap: 16px;

  @media (max-width: 768px) {
    flex-direction: column;
    align-items: stretch;
  }
}

.search-wrapper {
  flex: 1;
  max-width: 420px;
  display: flex;
  align-items: center;
  background: rgba(255, 255, 255, 0.8);
  border: 1px solid var(--border-light);
  border-radius: 16px;
  padding: 4px;
  transition: all 0.25s ease;

  &:focus-within {
    border-color: rgba(199, 91, 57, 0.35);
    box-shadow: 0 0 0 4px rgba(199, 91, 57, 0.08);
    background: #fff;
  }

  @media (max-width: 768px) {
    max-width: none;
  }
}

.search-icon {
  width: 42px;
  height: 42px;
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
  color: var(--text-ink);
  outline: none;
  font-family: 'Outfit', sans-serif;

  &::placeholder {
    color: var(--text-muted);
  }
}

.search-action {
  padding: 10px 20px;
  background: var(--text-ink);
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    background: #3d2e25;
    transform: translateY(-1px);
  }
}

.bulk-toggle {
  .toggle-label {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 13px;
    color: var(--text-body);
    cursor: pointer;
    user-select: none;
    padding: 8px 12px;
    background: rgba(255,255,255,0.6);
    border-radius: 10px;
    border: 1px solid var(--border-light);
    transition: all 0.2s ease;

    &:hover {
      background: #fff;
    }

    input {
      accent-color: var(--accent-coral);
    }
  }
}

.custom-table-view {
  background: rgba(255, 255, 255, 0.65);
  backdrop-filter: blur(16px);
  border-radius: 24px;
  border: 1px solid var(--border-light);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
  margin-bottom: 24px;
}

.table-header {
  display: grid;
  align-items: center;
  padding: 14px 20px;
  background: rgba(247, 243, 236, 0.6);
  border-bottom: 1px solid var(--border-light);
}

.th {
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--text-muted);
}

.checkbox-col {
  display: flex;
  align-items: center;
  justify-content: center;

  input[type="checkbox"] {
    width: 18px;
    height: 18px;
    accent-color: var(--accent-coral);
    cursor: pointer;
  }
}

.table-body {
  .table-row {
    display: grid;
    align-items: center;
    padding: 14px 20px;
    border-bottom: 1px solid var(--border-light);
    transition: background 0.2s ease;
    opacity: 0;
    transform: translateY(12px);
    animation: rowIn 0.5s cubic-bezier(0.16, 1, 0.3, 1) forwards;

    &:last-child {
      border-bottom: none;
    }

    &:hover {
      background: rgba(255, 255, 255, 0.9);
    }

    &.is-admin {
      background: rgba(124, 111, 174, 0.03);
    }

    &.is-selected {
      background: rgba(199, 91, 57, 0.05);
    }
  }
}

@keyframes rowIn {
  to { opacity: 1; transform: translateY(0); }
}

.table-user-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.table-avatar {
  width: 42px;
  height: 42px;
  border-radius: 12px;
  background: linear-gradient(135deg, var(--accent-coral), var(--accent-gold));
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 600;
  font-size: 16px;
  flex-shrink: 0;
}

.table-user-info {
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.table-username {
  font-weight: 500;
  color: var(--text-ink);
  font-size: 14px;
}

.table-user-id {
  font-size: 11px;
  color: var(--text-muted);
  font-family: 'DM Mono', monospace;
}

.table-role {
  display: inline-block;
  padding: 5px 10px;
  border-radius: 20px;
  font-size: 11px;
  font-weight: 600;

  &.admin {
    background: rgba(124, 111, 174, 0.12);
    color: var(--accent-purple);
  }

  &.user {
    background: rgba(45, 106, 79, 0.1);
    color: var(--accent-teal);
  }
}

.table-contact {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--text-body);
  font-size: 13px;

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

.table-actions {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
}

.table-btn {
  width: 34px;
  height: 34px;
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
    color: var(--accent-purple);
  }

  &.delete:hover {
    background: rgba(199, 91, 57, 0.1);
    color: var(--accent-coral);
  }
}

.grid-view {
  margin-bottom: 24px;
}

.user-cards {
  display: grid;
  grid-template-columns: 1fr;
  gap: 14px;
}

.user-card {
  position: relative;
  background: rgba(255, 255, 255, 0.8);
  border-radius: 22px;
  overflow: hidden;
  transition: all 0.35s cubic-bezier(0.34, 1.56, 0.64, 1);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-light);
  opacity: 0;
  transform: translateY(16px);
  animation: cardIn 0.55s cubic-bezier(0.16, 1, 0.3, 1) forwards;

  &:active {
    transform: scale(0.985);
  }

  &.is-selected {
    box-shadow: 0 0 0 2px var(--accent-coral), var(--shadow-md);
  }

  &.is-admin {
    .card-accent {
      background: linear-gradient(90deg, var(--accent-purple), #a89bc9);
    }

    .avatar-ring {
      background: linear-gradient(135deg, var(--accent-purple), #a89bc9);
    }
  }
}

@keyframes cardIn {
  to { opacity: 1; transform: translateY(0); }
}

.card-shimmer {
  position: absolute;
  inset: 0;
  background: linear-gradient(105deg, transparent 40%, rgba(255,255,255,0.6) 50%, transparent 60%);
  transform: translateX(-100%);
  transition: transform 0s;
  pointer-events: none;
}

.user-card:hover .card-shimmer {
  transform: translateX(100%);
  transition: transform 0.8s ease;
}

.card-accent {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, var(--accent-coral), var(--accent-gold));
}

.card-content {
  padding: 18px;
}

.user-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 14px;
}

.avatar-ring {
  width: 52px;
  height: 52px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--accent-coral), var(--accent-gold));
  padding: 2.5px;
  flex-shrink: 0;
}

.avatar {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-ink);
}

.user-meta {
  flex: 1;
  min-width: 0;
}

.user-name {
  font-size: 15px;
  font-weight: 500;
  color: var(--text-ink);
  margin: 0 0 2px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.user-id {
  font-size: 11px;
  color: var(--text-muted);
  font-family: 'DM Mono', monospace;
}

.role-tag {
  padding: 5px 10px;
  border-radius: 20px;
  font-size: 10px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.04em;

  &.admin {
    background: rgba(124, 111, 174, 0.12);
    color: var(--accent-purple);
  }

  &.user {
    background: rgba(45, 106, 79, 0.1);
    color: var(--accent-teal);
  }
}

.user-details {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 14px;
}

.detail-item {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 13px;
  color: var(--text-body);

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
  padding: 10px 14px;
  border-radius: 12px;
  border: none;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;

  &.secondary {
    background: var(--bg-warm);
    color: var(--text-body);

    &:hover {
      background: var(--text-ink);
      color: white;
    }
  }

  &.danger {
    background: rgba(199, 91, 57, 0.08);
    color: var(--accent-coral);

    &:hover {
      background: var(--accent-coral);
      color: white;
    }
  }
}

.selection-indicator {
  position: absolute;
  top: 14px;
  right: 14px;

  .check-circle {
    width: 22px;
    height: 22px;
    border-radius: 50%;
    border: 2px solid var(--border-light);
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    font-size: 12px;
    transition: all 0.2s ease;

    &.checked {
      background: var(--accent-coral);
      border-color: var(--accent-coral);
    }
  }
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: var(--text-muted);
  gap: 12px;

  .empty-icon {
    opacity: 0.5;
  }

  p {
    font-size: 14px;
  }
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  padding-top: 8px;
}

// 批量操作栏
.batch-bar {
  position: fixed;
  bottom: 28px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 100;
  background: rgba(31, 20, 14, 0.95);
  backdrop-filter: blur(16px);
  color: white;
  padding: 14px 22px;
  border-radius: 18px;
  display: flex;
  align-items: center;
  gap: 20px;
  box-shadow: var(--shadow-lg);
  border: 1px solid rgba(255,255,255,0.08);

  @media (max-width: 768px) {
    left: 16px;
    right: 16px;
    transform: none;
    flex-direction: column;
    gap: 12px;
    padding: 14px 16px;
    border-radius: 20px;
  }
}

.batch-info {
  display: flex;
  align-items: baseline;
  gap: 6px;
  font-size: 14px;
  opacity: 0.9;

  .batch-count {
    font-size: 20px;
    font-weight: 600;
    color: var(--accent-gold);
  }
}

.batch-actions {
  display: flex;
  gap: 10px;

  @media (max-width: 768px) {
    width: 100%;
  }
}

.batch-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 9px 16px;
  border-radius: 12px;
  border: none;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;

  @media (max-width: 768px) {
    flex: 1;
    justify-content: center;
  }

  &.role {
    background: rgba(124, 111, 174, 0.25);
    color: #e8e4f7;

    &:hover {
      background: rgba(124, 111, 174, 0.4);
    }
  }

  &.delete {
    background: rgba(199, 91, 57, 0.25);
    color: #ffe8e3;

    &:hover {
      background: rgba(199, 91, 57, 0.4);
    }
  }

  &.clear {
    background: rgba(255,255,255,0.1);
    color: rgba(255,255,255,0.8);

    &:hover {
      background: rgba(255,255,255,0.18);
    }
  }
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

.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 0.35s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.slide-up-enter-from,
.slide-up-leave-to {
  opacity: 0;
  transform: translateX(-50%) translateY(20px);

  @media (max-width: 768px) {
    transform: translateY(20px);
  }
}
</style>
