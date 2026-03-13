<template>
  <div class="user-manage-container">
    <!-- Ambient Background Effects -->
    <div class="ambient-glows">
      <div class="glow-orb orb-1"></div>
      <div class="glow-orb orb-2"></div>
    </div>
    <div class="noise-overlay"></div>

    <div class="content-wrapper">
      <div class="header-section">
        <div class="title-group">
          <button class="nav-btn" @click="router.back()">
            <n-icon :component="Left" />
          </button>
          <div class="text-content">
            <h1 class="page-title">用户管理</h1>
            <p class="page-subtitle">管理全站注册用户与权限</p>
          </div>
        </div>
        <div class="header-actions">
          <div class="search-bar">
            <n-input v-model:value="searchText" placeholder="搜索用户..." @keyup.enter="handleSearch" class="search-input">
              <template #prefix>
                <n-icon :component="Search" />
              </template>
            </n-input>
          </div>
          <div class="stat-badge">
            <span class="label">Total</span>
            <span class="value">{{ pagination.itemCount }}</span>
          </div>
        </div>
      </div>

      <div class="table-container glass-panel">
        <n-data-table
          :columns="columns"
          :data="userList"
          :loading="loading"
          :pagination="pagination"
          :row-key="(row: any) => row.id"
          :row-class-name="'user-row'"
          :scroll-x="800"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, h, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useMessage, NPopconfirm, NInput, NIcon, NDataTable, type DataTableColumns } from 'naive-ui';
import { Left, Search, User, Mail, Time, Permissions, Delete } from '@icon-park/vue-next';
import { adminGetAllUsers, adminUpdateUserRole, adminDeleteUser, type UserInfo } from '@/api/user'; 
import { ResultCode } from "@/utils/request";

const router = useRouter();
const message = useMessage();

const loading = ref(false);
const searchText = ref('');
const userList = ref<UserInfo[]>([]);

const pagination = reactive({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  prefix: ({ itemCount }: { itemCount?: number }) => `共 ${itemCount || 0} 位用户`,
  onChange: (page: number) => {
    pagination.page = page;
    fetchUsers();
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.pageSize = pageSize;
    pagination.page = 1;
    fetchUsers();
  }
});

const columns: DataTableColumns<UserInfo> = [
  {
    title: '用户',
    key: 'username',
    render(row: UserInfo) {
      return h(
        'div',
        { class: 'user-info-cell' },
        [
          h('div', { class: 'avatar-placeholder' }, row.username.charAt(0).toUpperCase()),
          h('div', { class: 'text-info' }, [
            h('span', { class: 'username' }, row.username),
            h('span', { class: 'user-id' }, `ID: ${row.id}`)
          ])
        ]
      );
    }
  },
  {
    title: '角色',
    key: 'user_group',
    width: 100,
    render(row: UserInfo) {
      const isAdmin = row.user_group === 'admin';
      return h(
        'div',
        { class: `role-badge ${isAdmin ? 'admin' : 'user'}` },
        [
          h(NIcon, { component: isAdmin ? Permissions : User, size: 14 }),
          h('span', isAdmin ? '管理员' : '用户')
        ]
      );
    }
  },
  {
    title: '邮箱',
    key: 'email',
    width: 200,
    render(row: UserInfo) {
      return h('div', { class: 'email-cell' }, [
        h(NIcon, { component: Mail, size: 14 }),
        h('span', row.email || '未绑定')
      ]);
    }
  },
  {
    title: '上次登录',
    key: 'last_login',
    width: 180,
    render(row: UserInfo) {
      return h('div', { class: 'time-cell' }, [
        h(NIcon, { component: Time, size: 14 }),
        h('span', row.last_login ? new Date(row.last_login).toLocaleString() : '从未登录')
      ]);
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 140,
    fixed: 'right',
    render(row: UserInfo) {
      return h('div', { class: 'actions-cell' }, [
        h(
          NPopconfirm,
          {
            onPositiveClick: () => handleRoleChange(row),
            positiveText: '确认',
            negativeText: '取消'
          },
          {
            trigger: () => h(
              'button',
              { 
                class: 'action-btn role', 
                title: row.user_group === 'admin' ? '降级为普通用户' : '提升为管理员' 
              },
              h(NIcon, { component: Permissions, size: 18 })
            ),
            default: () => `确定要将用户 ${row.username} 设为${row.user_group === 'admin' ? '普通用户' : '管理员'}吗？`
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
            trigger: () => h(
              'button',
              { class: 'action-btn delete', title: '删除用户' },
              h(NIcon, { component: Delete, size: 18 })
            ),
            default: () => '确定要删除该用户吗？此操作不可逆！'
          }
        )
      ]);
    }
  }
];

onMounted(() => {
  fetchUsers();
});

const fetchUsers = async () => {
  loading.value = true;
  try {
    const res = await adminGetAllUsers({ page: pagination.page, limit: pagination.pageSize, query: searchText.value });
    if (res.code === ResultCode.SUCCESS) {
      userList.value = res.data.list;
      pagination.itemCount = res.data.total;
    } else {
        message.error(res.message || '获取用户列表失败');
    }
  } catch (error) {
    message.error('获取用户列表失败');
    console.error(error);
  } finally {
    loading.value = false;
  }
};

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
      // 如果当前页只有一条数据且不是第一页，删除后跳转到上一页
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
</script>

<style lang="scss" scoped>
@import url('https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700&display=swap');

.user-manage-container {
  position: relative;
  min-height: 100vh;
  width: 100%;
  padding: 40px;
  background: #f8f9fc;
  font-family: 'Plus Jakarta Sans', sans-serif;
  overflow: hidden;

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

  .glow-orb {
    position: absolute;
    border-radius: 50%;
    filter: blur(100px);
    opacity: 0.4;
  }

  .orb-1 {
    top: -10%;
    right: -5%;
    width: 600px;
    height: 600px;
    background: radial-gradient(circle, rgba(16, 185, 129, 0.3) 0%, rgba(16, 185, 129, 0) 70%);
  }

  .orb-2 {
    bottom: -10%;
    left: -10%;
    width: 500px;
    height: 500px;
    background: radial-gradient(circle, rgba(59, 130, 246, 0.25) 0%, rgba(59, 130, 246, 0) 70%);
  }
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

.content-wrapper {
  position: relative;
  z-index: 2;
  max-width: 1200px;
  margin: 0 auto;
}

/* Header */
.header-section {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 32px;
  flex-wrap: wrap;
  gap: 20px;

  .title-group {
    display: flex;
    align-items: center;
    gap: 20px;

    .nav-btn {
      width: 48px;
      height: 48px;
      border-radius: 14px;
      border: 1px solid rgba(0, 0, 0, 0.05);
      background: rgba(255, 255, 255, 0.6);
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 24px;
      color: #64748b;
      cursor: pointer;
      transition: all 0.2s ease;
      backdrop-filter: blur(8px);

      &:hover {
        background: #fff;
        transform: translateX(-2px);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
        color: #1e293b;
      }
    }

    .text-content {
      .page-title {
        font-size: 32px;
        font-weight: 800;
        color: #1e293b;
        margin: 0;
        line-height: 1.1;
        letter-spacing: -0.02em;
      }

      .page-subtitle {
        font-size: 14px;
        color: #64748b;
        margin: 6px 0 0;
        font-weight: 500;
      }
    }
  }

  .header-actions {
    display: flex;
    align-items: center;
    gap: 16px;
    flex-wrap: wrap;

    .search-bar {
      width: 240px;
      
      :deep(.n-input) {
        border-radius: 12px;
        background-color: rgba(255, 255, 255, 0.6);
        border: 1px solid rgba(0, 0, 0, 0.05);
        
        &:hover, &.n-input--focus {
          background-color: #fff;
          box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
        }
      }
    }

    .stat-badge {
      padding: 8px 16px;
      background: rgba(255, 255, 255, 0.5);
      border: 1px solid rgba(255, 255, 255, 0.6);
      border-radius: 100px;
      display: flex;
      align-items: center;
      gap: 8px;
      backdrop-filter: blur(4px);

      .label {
        font-size: 12px;
        text-transform: uppercase;
        letter-spacing: 0.05em;
        color: #94a3b8;
        font-weight: 700;
      }

      .value {
        font-size: 16px;
        font-weight: 700;
        color: #10b981;
      }
    }
  }
}

/* Glass Panel Table */
.glass-panel {
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.6);
  border-radius: 24px;
  padding: 8px;
  box-shadow: 0 20px 40px -10px rgba(0, 0, 0, 0.05);
  overflow: hidden;
}

/* Custom Table Styles */
:deep(.n-data-table) {
  --n-th-font-weight: 700 !important;
  --n-th-text-color: #64748b !important;
  
  .n-data-table-th {
    background: transparent !important;
    border-bottom: 1px solid rgba(0, 0, 0, 0.05) !important;
    padding: 16px 24px !important;
    font-size: 13px;
    text-transform: uppercase;
    letter-spacing: 0.03em;
  }

  .n-data-table-td {
    background: transparent !important;
    border-bottom: 1px solid rgba(0, 0, 0, 0.03) !important;
    padding: 16px 24px !important;
    transition: background 0.2s;
  }

  .n-data-table-tr:last-child .n-data-table-td {
    border-bottom: none !important;
  }

  .n-data-table-tr:hover .n-data-table-td {
    background: rgba(255, 255, 255, 0.5) !important;
  }
}

/* Cell Renderers */
.user-info-cell {
  display: flex;
  align-items: center;
  gap: 12px;

  .avatar-placeholder {
    width: 40px;
    height: 40px;
    border-radius: 12px;
    background: linear-gradient(135deg, #3b82f6, #8b5cf6);
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 700;
    font-size: 18px;
    box-shadow: 0 4px 10px rgba(59, 130, 246, 0.2);
  }

  .text-info {
    display: flex;
    flex-direction: column;

    .username {
      font-size: 15px;
      font-weight: 700;
      color: #1e293b;
    }

    .user-id {
      font-size: 12px;
      color: #94a3b8;
      font-family: monospace;
    }
  }
}

.role-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 600;

  &.admin {
    background: rgba(245, 158, 11, 0.1);
    color: #d97706;
    border: 1px solid rgba(245, 158, 11, 0.2);
  }

  &.user {
    background: rgba(148, 163, 184, 0.1);
    color: #64748b;
    border: 1px solid rgba(148, 163, 184, 0.2);
  }
}

.email-cell, .time-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #64748b;
  font-size: 13px;
  
  .n-icon {
    color: #94a3b8;
  }
}

.actions-cell {
  display: flex;
  gap: 8px;
  justify-content: center;
}

.action-btn {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  border: none;
  background: transparent;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
  color: #94a3b8;

  &.role:hover {
    background: #fef3c7;
    color: #d97706;
    transform: scale(1.1);
  }

  &.delete:hover {
    background: #fee2e2;
    color: #ef4444;
    transform: scale(1.1);
  }
}

@media (max-width: 768px) {
  .user-manage-container {
    padding: 20px 16px;
  }

  .header-section {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
    
    .header-actions {
      width: 100%;
      justify-content: space-between;
      
      .search-bar {
        width: 100%;
        flex: 1;
      }
    }
  }
}
</style>