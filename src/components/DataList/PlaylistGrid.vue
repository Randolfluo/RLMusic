<template>
  <div class="playlist-grid">
    <n-spin :show="loading">
      <div v-if="!loading && (!playlists || playlists.length === 0)" class="empty">
        <n-empty :description="emptyText" />
      </div>
      <n-grid
        v-else
        x-gap="20"
        y-gap="20"
        :cols="cols"
        responsive="screen"
        :collapsed="collapsed"
        :collapsed-rows="collapsedRows"
      >
        <n-grid-item v-for="item in (playlists || [])" :key="item.id">
          <n-card
            hoverable
            class="playlist-card"
            content-style="padding: 0;"
            @click="onPlaylistClick(item.id)"
            @contextmenu.prevent="handleContextMenu($event, item)"
          >
            <div class="cover-container">
              <n-image
                preview-disabled
                class="cover-img"
                object-fit="cover"
                :src="item.cover_url || '/images/logo/favicon.png'"
                fallback-src="/images/logo/favicon.png"
              />
              <div class="play-count">
                <n-icon :component="Play" size="12" />
                <span>{{ formatCount(item.play_count) }}</span>
              </div>
            </div>
            <div class="info">
              <div class="title">{{ item.title }}</div>
            </div>
          </n-card>
        </n-grid-item>
      </n-grid>
    </n-spin>

    <!-- 右键菜单 -->
    <n-dropdown
      placement="bottom-start"
      trigger="manual"
      :x="dropdownX"
      :y="dropdownY"
      :options="menuOptions"
      :show="showDropdown"
      :on-clickoutside="onClickOutside"
      @select="handleSelect"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, nextTick, h } from "vue";
import { Play, Like, Delete, FolderPlus, More, PlayOne } from "@icon-park/vue-next";
import { useRouter } from "vue-router";
import { NDropdown, NIcon, NImage, useMessage, useDialog } from "naive-ui";
import { useUserDataStore } from "@/store/userData";
import { deletePrivatePlaylist, subscribePlaylist, unsubscribePlaylist, checkIsSubscribed } from "@/api/playlist";
import { ResultCode } from "@/utils/request";

const router = useRouter();
const userStore = useUserDataStore();
const message = useMessage();
const dialog = useDialog();

const props = defineProps({
  playlists: {
    type: Array as () => any[],
    default: () => [],
  },
  loading: {
    type: Boolean,
    default: false,
  },
  emptyText: {
    type: String,
    default: "暂无数据",
  },
  cols: {
    type: String,
    default: "2 s:3 m:4 l:5 xl:6",
  },
  collapsed: {
    type: Boolean,
    default: false,
  },
  collapsedRows: {
    type: Number,
    default: 1,
  },
});

const emit = defineEmits(['refresh']);

const showDropdown = ref(false);
const dropdownX = ref(0);
const dropdownY = ref(0);
const currentPlaylist = ref<any>(null);
const isSubscribed = ref(false);

// 渲染图标辅助函数
const renderIcon = (icon: any, color?: string) => {
  return () => h(NIcon, { color }, { default: () => h(icon) });
};

// 渲染菜单头部 (歌单信息)
const renderMenuHeader = (playlist: any) => {
  return h('div', {
    style: {
      display: 'flex',
      alignItems: 'center',
      padding: '4px 8px 8px 8px',
      borderBottom: '1px solid var(--n-divider-color)',
      marginBottom: '4px',
      cursor: 'default'
    }
  }, [
    h(NImage, {
      src: playlist.cover_url || '/images/logo/favicon.png',
      fallbackSrc: '/images/logo/favicon.png',
      width: 40,
      height: 40,
      previewDisabled: true,
      style: { borderRadius: '4px', marginRight: '10px' }
    }),
    h('div', { style: { display: 'flex', flexDirection: 'column', overflow: 'hidden' } }, [
      h('span', { style: { fontWeight: 'bold', fontSize: '14px', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis', maxWidth: '180px' } }, playlist.title),
      h('span', { style: { fontSize: '12px', opacity: 0.8, whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis', maxWidth: '180px' } }, `Created by User ${playlist.owner_id}`)
    ])
  ]);
};

const menuOptions = computed(() => {
    if (!currentPlaylist.value) return [];
    const playlist = currentPlaylist.value;
    const isOwner = userStore.userLogin && userStore.userData.userId === playlist.owner_id;
    
    // 如果是私有歌单，且不是拥有者，理论上应该看不到（除非有分享机制，但目前GetUserPrivatePlaylists只返回自己的）
    // 但这里也可能是公共歌单
    
    const options: any[] = [
        {
            key: 'header',
            type: 'render',
            render: () => renderMenuHeader(playlist),
            disabled: true
        },
        {
            label: '查看详情',
            key: 'detail',
            icon: renderIcon(More)
        }
    ];

    // 如果是公共歌单，或者不是自己的私有歌单（未来扩展），可以收藏
    if (userStore.userLogin && !isOwner) {
        options.push({
            label: isSubscribed.value ? '取消收藏' : '收藏',
            key: 'subscribe',
            icon: renderIcon(Like, isSubscribed.value ? '#d03050' : undefined)
        });
    }

    // 如果是拥有者，可以删除
    if (isOwner) {
        options.push({
            type: 'divider',
            key: 'd1'
        });
        options.push({
            label: '删除歌单',
            key: 'delete',
            icon: renderIcon(Delete)
        });
    }

    return options;
});

const handleContextMenu = async (e: MouseEvent, item: any) => {
    e.preventDefault();
    showDropdown.value = false;
    currentPlaylist.value = item;
    
    // 检查收藏状态 (如果是他人歌单)
    if (userStore.userLogin && userStore.userData.userId !== item.owner_id) {
        try {
            const res = await checkIsSubscribed(item.id);
            if (res.code === ResultCode.SUCCESS) {
                isSubscribed.value = res.data.is_subscribed;
            } else {
                isSubscribed.value = false;
            }
        } catch (e) {
            isSubscribed.value = false;
        }
    }

    nextTick(() => {
        showDropdown.value = true;
        dropdownX.value = e.clientX;
        dropdownY.value = e.clientY;
    });
};

const onClickOutside = () => {
    showDropdown.value = false;
};

const handleSelect = (key: string) => {
    showDropdown.value = false;
    const playlist = currentPlaylist.value;
    if (!playlist) return;

    switch (key) {
        case 'detail':
            router.push(`/playlist/${playlist.id}`);
            break;
        case 'subscribe':
            handleSubscribe(playlist);
            break;
        case 'delete':
            handleDelete(playlist);
            break;
    }
};

const handleSubscribe = async (playlist: any) => {
    try {
        let res;
        if (isSubscribed.value) {
            res = await unsubscribePlaylist(playlist.id);
        } else {
            res = await subscribePlaylist(playlist.id);
        }
        
        if (res.code === ResultCode.SUCCESS) {
            message.success(isSubscribed.value ? "已取消收藏" : "收藏成功");
        } else {
            message.error(res.message || "操作失败");
        }
    } catch (e) {
        message.error("操作失败");
    }
};

const handleDelete = (playlist: any) => {
    dialog.warning({
        title: "删除歌单",
        content: `确定要删除歌单 "${playlist.title}" 吗？此操作不可恢复。`,
        positiveText: "删除",
        negativeText: "取消",
        onPositiveClick: () => {
            deletePrivatePlaylist(playlist.id)
                .then(() => {
                    message.success("删除成功");
                    // 触发刷新事件，通知父组件重新获取列表
                    emit('refresh');
                    // 如果是在公共列表页删除了私有歌单（理论上不会，但如果混合显示），或者在私有列表页
                    // 父组件需要监听 refresh 事件
                })
                .catch((err) => {
                    message.error(err.message || "删除失败");
                });
        }
    });
};

const onPlaylistClick = (id: number) => {
  router.push(`/playlist/${id}`);
};

const formatCount = (count: number) => {
  if (!count) return 0;
  if (count > 10000) return (count / 10000).toFixed(1) + "万";
  return count;
};
</script>

<style scoped lang="scss">
.playlist-grid {
  .playlist-card {
    border-radius: 12px;
    overflow: hidden;
    cursor: pointer;
    transition: transform 0.3s;
    
    &:hover {
      transform: translateY(-5px);
    }
    
    .cover-container {
      position: relative;
      width: 100%;
      padding-top: 100%;
      background-color: #f5f5f5;
      
      .cover-img {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        
        :deep(img) {
            width: 100%;
            height: 100%;
            object-fit: cover;
            display: block;
        }
      }
      
      .play-count {
        position: absolute;
        right: 8px;
        top: 8px;
        background: rgba(0, 0, 0, 0.4);
        backdrop-filter: blur(4px);
        color: #fff;
        padding: 2px 8px;
        border-radius: 12px;
        font-size: 12px;
        display: flex;
        align-items: center;
        gap: 4px;
      }
    }
    
    .info {
      padding: 12px 10px;
      .title {
        font-size: 14px;
        line-height: 1.4;
        height: 40px;
        overflow: hidden;
        text-overflow: ellipsis;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        font-weight: 500;
        color: var(--n-text-color);
      }
    }
  }
}
</style>
