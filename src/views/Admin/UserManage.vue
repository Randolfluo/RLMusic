<template>
  <div class="user-manage-container">
    <div class="header">
      <div class="title-section">
        <n-button circle secondary @click="router.back()" class="back-btn">
          <template #icon><n-icon :component="Left" /></template>
        </n-button>
        <h1>用户管理</h1>
      </div>
      <div class="actions">
        <n-input-group>
          <n-input v-model:value="searchText" placeholder="搜索用户名/邮箱..." @keyup.enter="handleSearch">
            <template #prefix>
              <n-icon :component="Search" />
            </template>
          </n-input>
          <n-button type="primary" @click="handleSearch">搜索</n-button>
        </n-input-group>
      </div>
    </div>

    <div class="content glass-panel">
      <n-data-table
        :columns="columns"
        :data="userList"
        :loading="loading"
        :pagination="pagination"
        :row-key="row => row.id"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, h, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useMessage, NButton, NTag, NSpace, NPopconfirm } from 'naive-ui';
import { Left, Search } from '@icon-park/vue-next';
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

const columns = [
  {
    title: 'ID',
    key: 'id',
    width: 80
  },
  {
    title: '用户名',
    key: 'username',
    render(row: UserInfo) {
      return h(
        'div',
        { style: 'display: flex; align-items: center; gap: 8px;' },
        [
          h(NTag, { type: row.user_group === 'admin' ? 'warning' : 'default', size: 'small' }, { default: () => row.user_group === 'admin' ? '管理员' : '用户' }),
          row.username
        ]
      );
    }
  },
  {
    title: '邮箱',
    key: 'email'
  },
  {
    title: '上次登录',
    key: 'last_login',
    render(row: UserInfo) {
      return row.last_login ? new Date(row.last_login).toLocaleString() : '从未登录';
    }
  },
  {
    title: '操作',
    key: 'actions',
    render(row: UserInfo) {
      return h(NSpace, null, {
        default: () => [
          h(
            NPopconfirm,
            {
              onPositiveClick: () => handleRoleChange(row)
            },
            {
              trigger: () => h(
                NButton,
                {
                  size: 'small',
                  type: row.user_group === 'admin' ? 'default' : 'warning',
                  secondary: true
                },
                { default: () => row.user_group === 'admin' ? '设为用户' : '设为管理员' }
              ),
              default: () => `确定要将用户 ${row.username} 设为${row.user_group === 'admin' ? '普通用户' : '管理员'}吗？`
            }
          ),
          h(
            NPopconfirm,
            {
              onPositiveClick: () => handleDelete(row)
            },
            {
              trigger: () => h(
                NButton,
                { size: 'small', type: 'error', secondary: true },
                { default: () => '删除' }
              ),
              default: () => '确定要删除该用户吗？此操作不可逆！'
            }
          )
        ]
      });
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
.user-manage-container {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
  min-height: 100vh;
  
  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
    
    .title-section {
      display: flex;
      align-items: center;
      gap: 16px;
      
      h1 {
        margin: 0;
        font-size: 24px;
        color: var(--n-text-color);
      }
    }
  }
  
  .glass-panel {
    background: rgba(255, 255, 255, 0.6);
    backdrop-filter: blur(16px);
    border: 1px solid rgba(255, 255, 255, 0.6);
    border-radius: 16px;
    padding: 24px;
    
    :global(.dark) & {
      background: rgba(30, 41, 59, 0.4);
      border-color: rgba(255, 255, 255, 0.05);
    }
  }
}
</style>