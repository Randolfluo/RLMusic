<!-- 
  SongList.vue
  功能：歌曲列表组件
  说明：
    - 展示歌曲列表，支持多种视图模式（缩略图、简洁）
    - 提供歌曲播放、跳转歌手/专辑详情等功能
    - 支持分页加载
    - 支持右键菜单操作（播放、喜欢、下载等）
-->
<template>
  <div class="song-list-component">
    <!-- 列表控制栏（右上角视图切换） -->
    <div class="list-control">
      <div class="left">
        <!-- 多选模式下的操作按钮 -->
        <n-button-group v-if="isMultiSelectMode && selectedRowKeys.length > 0" size="small">
          <n-button @click="handleBatchPlay">
            <template #icon>
              <n-icon :component="PlayOne" />
            </template>
            播放
          </n-button>
          <n-button @click="handleBatchAddToPlaylist">
            <template #icon>
              <n-icon :component="FolderPlus" />
            </template>
            添加到
          </n-button>
          <n-button @click="handleBatchDownload">
            <template #icon>
              <n-icon :component="Download" />
            </template>
            下载
          </n-button>
          <n-button @click="handleBatchDelete">
            <template #icon>
              <n-icon :component="Delete" />
            </template>
            删除
          </n-button>
        </n-button-group>
        <span v-if="isMultiSelectMode" style="margin-left: 10px; font-size: 12px; opacity: 0.6;">
          已选择 {{ selectedRowKeys.length }} 项
        </span>
        <slot name="controls"></slot>
      </div>
      <div class="right">
        <n-button-group size="small">
          <n-tooltip trigger="hover">
            <template #trigger>
              <n-button :type="isMultiSelectMode ? 'primary' : 'default'" @click="toggleMultiSelect">
                <template #icon>
                  <n-icon :component="CheckOne" />
                </template>
              </n-button>
            </template>
            {{ isMultiSelectMode ? '退出多选' : '批量操作' }}
          </n-tooltip>
          <n-tooltip trigger="hover">
            <template #trigger>
              <n-button :type="viewMode === 'thumbnail' ? 'primary' : 'default'" @click="viewMode = 'thumbnail'">
                <template #icon>
                  <n-icon :component="Pic" />
                </template>
              </n-button>
            </template>
            缩略图模式
          </n-tooltip>
          <n-tooltip trigger="hover">
            <template #trigger>
              <n-button :type="viewMode === 'concise' ? 'primary' : 'default'" @click="viewMode = 'concise'">
                <template #icon>
                  <n-icon :component="HamburgerButton" />
                </template>
              </n-button>
            </template>
            简洁模式
          </n-tooltip>
        </n-button-group>
      </div>
    </div>

    <!-- 歌曲数据表格 -->
    <n-data-table
      :columns="columns"
      :data="songs"
      :bordered="false"
      :row-props="rowProps"
      :row-class-name="() => 'song-row'"
      :loading="loading"
      :row-key="(row) => row.id"
      v-model:checked-row-keys="selectedRowKeys"
    />
    
    <!-- 空状态显示 -->
    <div v-if="!loading && songs.length === 0" class="empty">
      <n-empty description="暂无歌曲" />
    </div>

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

    <!-- 添加到歌单模态框 -->
    <AddToPlaylistModal
      v-model:show="showAddToPlaylistModal"
      :song-ids="songsToAdd"
      @success="handleAddToPlaylistSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h, nextTick } from "vue";
import { useRouter } from "vue-router";
import { NButton, NButtonGroup, NIcon, NImage, NTooltip, NDataTable, NEmpty, NDropdown, useMessage, useDialog } from "naive-ui";
import { HamburgerButton, Pic, Like, PlayOne, PlayTwo, Download, FolderPlus, Copy, CheckOne, More, CloudStorage, Delete } from "@icon-park/vue-next";
import { musicStore, settingStore } from "@/store";
import AddToPlaylistModal from "@/components/DataModel/AddToPlaylistModal.vue";
import { removeSongsFromPlaylist } from "@/api/playlist";

// Props 定义
const props = defineProps({
  songs: {
    type: Array as () => any[],
    default: () => [],
  },
  loading: {
    type: Boolean,
    default: false,
  },
  page: {
    type: Number,
    default: 1,
  },
  pageSize: {
    type: Number,
    default: 30,
  },
  // 歌单ID，用于删除操作
  playlistId: {
    type: Number,
    default: 0
  },
  // 是否是当前用户的私有歌单（有权修改）
  isOwner: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['refresh']);

const router = useRouter();
const music = musicStore();
const setting = settingStore();
const viewMode = ref<'thumbnail' | 'concise'>('thumbnail');
const message = useMessage();
const dialog = useDialog();
const showDropdown = ref(false);
const dropdownX = ref(0);
const dropdownY = ref(0);
const currentSong = ref<any>(null);
const isMultiSelectMode = ref(false);
const selectedRowKeys = ref<Array<string | number>>([]);

// 添加到歌单模态框状态
const showAddToPlaylistModal = ref(false);
const songsToAdd = ref<number[]>([]);

const openAddToPlaylist = (ids: number[]) => {
  songsToAdd.value = ids;
  showAddToPlaylistModal.value = true;
};

const handleAddToPlaylistSuccess = () => {
    // 如果是多选模式，操作成功后退出多选
    if (isMultiSelectMode.value) {
        toggleMultiSelect();
    }
};

const toggleMultiSelect = () => {
  isMultiSelectMode.value = !isMultiSelectMode.value;
  if (!isMultiSelectMode.value) {
    selectedRowKeys.value = [];
  }
};

// 批量操作处理函数
const handleBatchPlay = () => {
    if (selectedRowKeys.value.length === 0) return;
    // 获取选中的歌曲对象
    const selectedSongs = props.songs.filter(s => selectedRowKeys.value.includes(s.id));
    if (selectedSongs.length > 0) {
        const tracks = mapSongsToPlayer(selectedSongs);
        music.setPlaylists(tracks);
        music.setPlaySongIndex(0);
        music.setPlayState(true);
        message.success(`已开始播放选中的 ${selectedSongs.length} 首歌曲`);
        // 退出多选模式
        toggleMultiSelect();
    }
};

const handleBatchAddToPlaylist = () => {
    if (selectedRowKeys.value.length === 0) return;
    openAddToPlaylist(selectedRowKeys.value.map(id => Number(id)));
};

const handleBatchDownload = () => {
    message.info(`批量下载: ${selectedRowKeys.value.length} 首 (功能开发中)`);
};

const handleBatchDelete = () => {
    if (selectedRowKeys.value.length === 0) return;
    
    if (!props.playlistId || !props.isOwner) {
        message.warning("只有歌单所有者可以删除歌曲");
        return;
    }

    dialog.warning({
        title: "批量删除",
        content: `确定要从歌单中删除选中的 ${selectedRowKeys.value.length} 首歌曲吗？`,
        positiveText: "删除",
        negativeText: "取消",
        onPositiveClick: () => {
            const songIds = selectedRowKeys.value.map(id => Number(id));
            removeSongsFromPlaylist({ playlist_id: props.playlistId, song_ids: songIds })
                .then(() => {
                    message.success("删除成功");
                    toggleMultiSelect();
                    emit('refresh');
                })
                .catch((err) => {
                    message.error(err.message || "删除失败");
                });
        }
    });
};

// 渲染图标辅助函数
const renderIcon = (icon: any, color?: string) => {
  return () => h(NIcon, { color }, { default: () => h(icon) });
};

// 渲染菜单头部 (歌曲信息)
const renderMenuHeader = (song: any) => {
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
      src: song.cover_url || (song.album ? song.album.picUrl : null) || song.picUrl || `/api/song/cover/${song.id}`,
      width: 40,
      height: 40,
      previewDisabled: true,
      style: { borderRadius: '4px', marginRight: '10px' }
    }),
    h('div', { style: { display: 'flex', flexDirection: 'column', overflow: 'hidden' } }, [
      h('span', { style: { fontWeight: 'bold', fontSize: '14px', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis', maxWidth: '180px' } }, song.title),
      h('span', { style: { fontSize: '12px', opacity: 0.8, whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis', maxWidth: '180px' } }, song.artist_name || (song.artists ? song.artists.map((a: any) => a.name).join(' / ') : 'Unknown'))
    ])
  ]);
};

// 右键菜单选项配置
const menuOptions = computed(() => {
    if (!currentSong.value) return [];
    const isLiked = music.getSongIsLike(currentSong.value.id);
    const song = currentSong.value;

    const options = [
        {
            key: 'header',
            type: 'render',
            render: () => renderMenuHeader(song),
            disabled: true // 禁止点击
        },
        {
            label: '立即播放',
            key: 'play',
            icon: renderIcon(PlayOne)
        },
        {
            label: '下一首播放',
            key: 'play-next',
            icon: renderIcon(PlayTwo)
        },
        {
            label: '添加到歌单',
            key: 'add-to-playlist',
            icon: renderIcon(FolderPlus)
        },
        {
            type: 'divider',
            key: 'd1'
        },
        {
            label: isLiked ? '取消喜欢' : '喜欢',
            key: 'like',
            icon: renderIcon(Like, isLiked ? '#d03050' : undefined)
        },
        {
            label: '更多操作',
            key: 'more',
            icon: renderIcon(More),
            children: [
                {
                    label: '下载',
                    key: 'download',
                    icon: renderIcon(Download)
                },
                {
                    label: '复制链接',
                    key: 'copy-link',
                    icon: renderIcon(Copy)
                }
            ]
        }
    ];

    if (props.playlistId > 0 && props.isOwner) {
        options.push(
            {
                type: 'divider',
                key: 'd2'
            },
            {
                label: '从歌单中删除',
                key: 'delete-from-playlist',
                icon: renderIcon(Delete)
            }
        );
    }

    return options;
});

// 处理右键点击事件
const handleContextMenu = (e: MouseEvent, row: any) => {
    e.preventDefault();
    showDropdown.value = false;
    nextTick(() => {
        currentSong.value = row;
        showDropdown.value = true;
        dropdownX.value = e.clientX;
        dropdownY.value = e.clientY;
    });
};

// 点击外部关闭右键菜单
const onClickOutside = () => {
    showDropdown.value = false;
};

// 处理右键菜单选择
const handleSelect = (key: string) => {
    showDropdown.value = false;
    const song = currentSong.value;
    if (!song) return;

    switch (key) {
        case 'play':
            // 播放当前歌曲
            const index = props.songs.findIndex(s => s.id === song.id);
            if (index !== -1) {
                 const tracks = mapSongsToPlayer(props.songs);
                 music.setPlaylists(tracks);
                 music.setPlaySongIndex(index);
                 music.setPlayState(true);
            }
            break;
        case 'play-next':
            // 构造播放器需要的歌曲对象结构
            const track = {
                ...song,
                name: song.title,
                artist: [{ name: song.artist_name, id: song.artist_id }],
                album: { name: song.album_title, id: song.album_id, picUrl: song.cover_url }
            };
            music.addSongToNext(track);
            message.success('已添加到下一首播放');
            break;
        case 'like':
            const isLiked = music.getSongIsLike(song.id);
            music.changeLikeList(song.id, !isLiked);
            break;
        case 'add-to-playlist':
            openAddToPlaylist([song.id]);
            break;
        case 'download':
            message.info('开始下载 (功能开发中)');
            break;
        case 'copy-link':
             const link = `${window.location.origin}/song/${song.id}`;
             navigator.clipboard.writeText(link).then(() => {
                 message.success('链接已复制');
             });
            break;
        case 'delete-from-playlist':
            if (!props.playlistId || !props.isOwner) {
                message.warning("只有歌单所有者可以删除歌曲");
                return;
            }
            dialog.warning({
                title: "删除歌曲",
                content: `确定要从歌单中删除歌曲 "${song.title}" 吗？`,
                positiveText: "删除",
                negativeText: "取消",
                onPositiveClick: () => {
                    removeSongsFromPlaylist({ playlist_id: props.playlistId, song_ids: [song.id] })
                        .then(() => {
                            message.success("删除成功");
                            emit('refresh');
                        })
                        .catch((err) => {
                            message.error(err.message || "删除失败");
                        });
                }
            });
            break;
    }
};

// 表格列配置
const columns = computed(() => {
  const baseColumns: any[] = [];
  
  if (isMultiSelectMode.value) {
    baseColumns.push({
      type: 'selection',
      fixed: 'left'
    });
  }

  baseColumns.push(
    {
      title: "#",
      key: "index",
      width: 60,
      align: 'center',
      render: (_: any, index: number) => h('span', { 
        style: { opacity: 0.6, fontFamily: 'Monaco, monospace', fontSize: '13px' } 
      }, `${index + 1 + (props.page - 1) * props.pageSize}`),
    },
    {
      title: "标题",
      key: "title",
      render: (row: any) => {
        return h('span', {
          style: { cursor: 'pointer', transition: 'all 0.3s', fontSize: '15px', fontWeight: '500' },
          onClick: (e: Event) => {
             e.stopPropagation();
             router.push(`/song/${row.id}`);
          },
          onMouseover: (e: Event) => {
            (e.target as HTMLElement).style.color = setting.themeColor;
            // (e.target as HTMLElement).style.paddingLeft = '4px'; // 增加一点位移效果
          },
          onMouseout: (e: Event) => {
            (e.target as HTMLElement).style.color = 'inherit';
            // (e.target as HTMLElement).style.paddingLeft = '0';
          },
          class: 'song-title-link'
        }, row.title)
      }
    },
    {
      title: "歌手",
      key: "artist_name",
      render: (row: any) => {
        // Handle array of artists or single artist name string
        let artistName = row.artist_name;
        let artistId = row.artist_id;

        if (!artistName && row.artists && row.artists.length > 0) {
             artistName = row.artists.map((a: any) => a.name).join(' / ');
             artistId = row.artists[0]?.id;
        }

        return h('span', {
            class: 'artist-link',
            style: { cursor: 'pointer', transition: 'color 0.3s', fontSize: '13px', opacity: 0.8 },
            onClick: (e: Event) => {
                e.stopPropagation();
                if (artistId) router.push({ name: 'artist', query: { id: artistId } });
            },
            onMouseover: (e: Event) => {
                (e.target as HTMLElement).style.color = setting.themeColor;
                (e.target as HTMLElement).style.opacity = '1';
            },
            onMouseout: (e: Event) => {
                (e.target as HTMLElement).style.color = 'inherit';
                (e.target as HTMLElement).style.opacity = '0.8';
            }
        }, artistName || "Unknown Artist")
      },
    },
    {
      title: "专辑",
      key: "album_title",
      render: (row: any) => {
        // 先尝试获取 album_name（API返回的标准字段），然后是 album_title（旧兼顾），再是对象内的属性
        const albumName = row.album_name || row.album_title || row.album?.name || row.album?.title || "Unknown Album";
        const albumId = row.album_id || row.album?.id;
        
        return h('span', {
           class: 'album-link',
           style: { cursor: 'pointer', transition: 'color 0.3s', fontSize: '13px', opacity: 0.6 },
           onClick: (e: Event) => {
               e.stopPropagation();
               if (albumId) router.push({ path: '/album', query: { id: albumId } });
           },
           onMouseover: (e: Event) => {
               (e.target as HTMLElement).style.color = setting.themeColor;
               (e.target as HTMLElement).style.opacity = '1';
           },
           onMouseout: (e: Event) => {
               (e.target as HTMLElement).style.color = 'inherit';
               (e.target as HTMLElement).style.opacity = '0.6';
           }
        }, albumName);
      },
    },
    {
      title: "时长",
      key: "duration",
      width: 100,
      render: (row: any) => h('span', { style: { opacity: 0.5, fontFamily: 'Monaco, monospace', fontSize: '12px' } }, formatDuration(row.duration)),
    },
  );

  if (viewMode.value === 'thumbnail') {
    // 如果是多选模式，插入位置要后移一位
    const insertIndex = isMultiSelectMode.value ? 2 : 1;
    baseColumns.splice(insertIndex, 0, {
      title: "封面",
      key: "cover",
      width: 80,
      align: 'center',
      render: (row: any) => {
        return h(NImage, {
          src: row.cover_url || (row.album ? row.album.picUrl : null) || row.picUrl || `/api/song/cover/${row.id}`,
          fallbackSrc: '/images/logo/favicon.png',
          width: 48,
          height: 48,
          lazy: true, 
          objectFit: 'cover',
          style: { borderRadius: '6px', verticalAlign: 'middle', boxShadow: '0 2px 6px rgba(0,0,0,0.2)' },
          previewDisabled: true
        });
      }
    });
  }

  return baseColumns;
});

// 格式化时长
const formatDuration = (seconds: number) => {
  if (!seconds) return "00:00";
  const m = Math.floor(seconds / 60);
  const s = Math.floor(seconds % 60);
  return `${m.toString().padStart(2, "0")}:${s.toString().padStart(2, "0")}`;
};

// 行属性配置 (点击播放、右键菜单)
const rowProps = (row: any, index: number) => {
  return {
    style: "cursor: pointer;",
    onClick: (e: Event) => {
        // 多选模式下点击行不触发播放，而是切换选中状态（NDataTable 默认行为可能需要自行处理，或者仅禁止播放）
        if (isMultiSelectMode.value) {
           const id = row.id;
           const idx = selectedRowKeys.value.indexOf(id);
           if (idx > -1) {
             selectedRowKeys.value.splice(idx, 1);
           } else {
             selectedRowKeys.value.push(id);
           }
           return;
        }

        const tracks = mapSongsToPlayer(props.songs);
        music.setPlaylists(tracks);
        music.setPlaySongIndex(index);
        music.setPlayState(true);
    },
    onContextmenu: (e: MouseEvent) => handleContextMenu(e, row)
  };
};

// 映射歌曲到播放器格式
const mapSongsToPlayer = (list: any[]) => {
    return list.map(song => ({
        ...song,
        name: song.title,
        // Ensure artist format is consistent for player
        artist: song.artists || [{ name: song.artist_name, id: song.artist_id }],
        album: song.album || { 
            name: song.album_name || song.album_title, // 优先使用 album_name
            id: song.album_id, 
            picUrl: song.cover_url || `/api/song/cover/${song.id}` 
        },
        picUrl: song.cover_url || `/api/song/cover/${song.id}`
    }));
};
</script>

<style scoped lang="scss">
.song-list-component {
    padding-bottom: 20px;
    position: relative;
}
.list-control {
    position: absolute;
    top: -46px; /* 向上移动，与父组件的标题行或控制行对齐 */
    right: 0;
    z-index: 100;
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0;
    pointer-events: none;
    
    .right {
        pointer-events: auto;
    }
    
    /* 如果 slot 有内容，可能需要额外处理，这里默认 float right */
    .left {
       /* display: none; */ /* 暂时隐藏左侧空 slot */
       pointer-events: auto;
       display: flex;
       align-items: center;
    }
}

:deep(.n-data-table) {
  .n-data-table-th {
    background-color: transparent !important;
    border-bottom: 1px solid rgba(128, 128, 128, 0.2);
    font-weight: normal;
    opacity: 0.7;
    padding: 12px 16px;
  }
  
  .n-data-table-td {
    background-color: transparent !important;
    border-bottom: 1px solid rgba(128, 128, 128, 0.1);
    transition: background-color 0.2s ease;
    padding: 10px 16px; 
    vertical-align: middle;
  }

  .song-row {
     transition: all 0.2s ease;
     &:hover {
        td {
           background-color: rgba(128, 128, 128, 0.1) !important;
        }
     }
  }
  
  // 隐藏最后一个 border
  .n-data-table-tr:last-child .n-data-table-td {
      border-bottom: none;
  }
}
</style>
