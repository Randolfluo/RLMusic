<template>
  <div class="private-playlists">
    <div class="header">
      <div class="cover">
        <n-image
          class="cover-img"
          :src="resolveAvatarUrl(userStore.userData.avatarUrl) || '/images/logo/favicon.png'"
          fallback-src="/images/logo/favicon.png"
          object-fit="cover"
          preview-disabled
        />
      </div>
      <div class="info">
        <div class="tag">私有歌单</div>
        <div class="title">{{ userStore.userData.nickname }} 的歌单</div>
        <div class="creator">
          <n-avatar 
            round 
            size="small" 
            :src="resolveAvatarUrl(userStore.userData.avatarUrl) || '/images/logo/favicon.png'" 
            style="margin-right: 8px; vertical-align: middle;"
          />
          <span style="vertical-align: middle;">{{ userStore.userData.nickname }}</span>
        </div>
        <div class="desc">
            共 {{ total }} 个歌单
        </div>
        <n-button type="primary" @click="showCreateModal = true">
          <template #icon>
            <n-icon :component="Plus" />
          </template>
          创建歌单
        </n-button>
      </div>
    </div>

    <n-divider />

    <PlaylistGrid
      :loading="loading"
      :playlists="playlists"
      empty-text="暂无私有歌单"
      :enable-ai-intro="true"
      @refresh="getPlaylists"
      @generate-intro="handleGenerateIntro"
    />

    <div class="pagination-container" style="display: flex; justify-content: center; margin-top: 20px;">
      <Pagination
        v-if="playlists.length > 0"
        :totalCount="total"
        :pageNumber="page"
        :showSizePicker="true"
        @pageNumberChange="onPageChange"
        @pageSizeChange="onPageSizeChange"
      />
    </div>

    <!-- 创建私有歌单弹窗 -->
    <n-modal v-model:show="showCreateModal" title="创建私有歌单" preset="card" style="width: 420px; max-width: 90vw;" :mask-closable="false">
      <n-form :model="createForm" label-placement="left" label-width="80">
        <n-form-item label="歌单名称" required>
          <n-input v-model:value="createForm.title" placeholder="输入歌单名称" maxlength="50" show-count clearable />
        </n-form-item>
        <n-form-item label="歌单描述">
          <n-input v-model:value="createForm.description" type="textarea" :rows="3" placeholder="输入歌单描述（可选）" maxlength="200" show-count clearable />
        </n-form-item>
      </n-form>
      <template #footer>
        <div style="display: flex; justify-content: flex-end; gap: 12px;">
          <n-button @click="showCreateModal = false">取消</n-button>
          <n-button type="primary" :loading="creating" :disabled="!createForm.title.trim()" @click="handleCreatePlaylist">创建</n-button>
        </div>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { getUserPrivatePlaylists, createPrivatePlaylist } from "@/api/playlist";
import { generatePlaylistIntros } from "@/api/ai";
import { ResultCode } from "@/utils/request";
import { useMessage, NImage, NDivider, NAvatar, useDialog, NButton, NIcon, NModal, NForm, NFormItem, NInput } from "naive-ui";
import PlaylistGrid from "@/components/DataList/PlaylistGrid.vue";
import Pagination from "@/components/Pagination/index.vue";
import { useUserDataStore } from "@/store/userData";
import { resolveAvatarUrl } from "@/api/user";
import { Plus } from "@icon-park/vue-next";

const message = useMessage();
const dialog = useDialog();
const userStore = useUserDataStore();
const loading = ref(false);
const playlists = ref<any[]>([]);
const page = ref(1);
const limit = ref(20);
const total = ref(0);
const showCreateModal = ref(false);
const createForm = ref({ title: "", description: "" });
const creating = ref(false);

onMounted(() => {
  getPlaylists();
});

const getPlaylists = async () => {
  loading.value = true;
  try {
    const res = await getUserPrivatePlaylists(page.value, limit.value);
    if (res.code === ResultCode.SUCCESS) {
      if (Array.isArray(res.data)) {
          playlists.value = res.data;
          total.value = res.data.length;
      } else {
          playlists.value = res.data.list;
          total.value = res.data.total;
      }
    }
  } catch (error) {
    message.error("获取歌单失败");
  } finally {
    loading.value = false;
  }
};

const onPageChange = (val: number) => {
  page.value = val;
  getPlaylists();
};

const onPageSizeChange = (val: number) => {
  limit.value = val;
  page.value = 1;
  getPlaylists();
};

const handleGenerateIntro = (playlist: any) => {
  dialog.info({
    title: "生成开场白",
    content: `确定要为歌单 "${playlist.title}" 生成开场白吗？`,
    positiveText: "生成",
    negativeText: "取消",
    onPositiveClick: () => {
      generatePlaylistIntros(playlist.id)
        .then((res) => {
          if (res.code === ResultCode.SUCCESS) {
            message.success("已开始生成开场白，请稍后查看");
          } else {
            message.error(res.message || "生成失败");
          }
        })
        .catch(() => {
          message.error("请求失败");
        });
    },
  });
};

const handleCreatePlaylist = async () => {
  const title = createForm.value.title.trim();
  if (!title) {
    message.warning("歌单名称不能为空");
    return;
  }
  creating.value = true;
  try {
    const res = await createPrivatePlaylist({
      title,
      description: createForm.value.description.trim(),
    });
    if (res.code === ResultCode.SUCCESS) {
      message.success("创建成功");
      createForm.value = { title: "", description: "" };
      showCreateModal.value = false;
      getPlaylists();
    } else {
      message.error(res.message || "创建失败");
    }
  } catch (error) {
    message.error("创建失败");
  } finally {
    creating.value = false;
  }
};
</script>

<style scoped lang="scss">
.private-playlists {
  padding: 24px;
  
  .header {
    display: flex;
    margin-bottom: 24px;
    
    .cover {
      width: 200px;
      height: 200px;
      border-radius: 8px;
      overflow: hidden;
      margin-right: 24px;
      flex-shrink: 0;
      box-shadow: 0 4px 12px rgba(0,0,0,0.1);
      
      .cover-img {
        width: 100%;
        height: 100%;
      }
    }
    
    .info {
      flex: 1;
      display: flex;
      flex-direction: column;
      justify-content: flex-start;
      
      .tag {
        display: inline-block;
        width: fit-content;
        padding: 4px 12px;
        border: 1px solid var(--n-color-primary);
        color: var(--n-color-primary);
        border-radius: 4px;
        font-size: 14px;
        margin-bottom: 12px;
      }
      
      .title {
        font-size: 32px;
        font-weight: bold;
        margin-bottom: 12px;
        line-height: 1.2;
      }
      
      .creator {
        display: flex;
        align-items: center;
        margin-bottom: 16px;
        font-size: 14px;
        opacity: 0.8;
      }
      
      .desc {
        font-size: 14px;
        opacity: 0.6;
        line-height: 1.6;
        margin-bottom: 16px;
        display: -webkit-box;
        -webkit-line-clamp: 3;
        -webkit-box-orient: vertical;
        overflow: hidden;
      }
    }
  }
}
</style>
