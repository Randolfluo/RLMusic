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
          <n-button v-if="!isMobile" @click="handleBatchDownload">
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
    >
      <template #empty>
        <div class="empty">
          <n-empty :description="emptyDescription" />
        </div>
      </template>
    </n-data-table>

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
import { ref, computed, h, nextTick, onMounted, onUnmounted } from "vue";
import { useRouter } from "vue-router";
import { NButton, NButtonGroup, NIcon, NImage, NTooltip, NDataTable, NEmpty, NDropdown, useMessage, useDialog } from "naive-ui";
import { HamburgerButton, Pic, Like, PlayOne, PlayTwo, PauseOne, Download, FolderPlus, Copy, CheckOne, More, Delete, VolumeNotice } from "@icon-park/vue-next";
import { musicStore, settingStore } from "@/store";
import AddToPlaylistModal from "@/components/DataModel/AddToPlaylistModal.vue";
import { removeSongsFromPlaylist } from "@/api/playlist";
import { getSongCover, resolveCoverUrl } from "@/api/song";

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
  },
  // 空状态描述
  emptyDescription: {
    type: String,
    default: "暂无歌曲"
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
const isMobile = ref(window.innerWidth < 640);

const handleResize = () => {
  isMobile.value = window.innerWidth < 640;
};

onMounted(() => {
  window.addEventListener('resize', handleResize);
});

onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
});

// 添加到歌单模态框状态
const showAddToPlaylistModal = ref(false);
const songsToAdd = ref<number[]>([]);
const getCoverSrc = (coverUrl?: string, id?: number | string, album?: any, picUrl?: string) => {
  const candidate = coverUrl || album?.picUrl || picUrl;
  if (candidate) return resolveCoverUrl(candidate);
  if (id) return getSongCover(id);
  return "/images/logo/logo.png";
};

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
    if (selectedRowKeys.value.length === 0) return;
    const selectedSongs = props.songs.filter(s => selectedRowKeys.value.includes(s.id));
    selectedSongs.forEach((song, index) => {
        setTimeout(() => {
            handleDownload(song);
        }, index * 500);
    });
    toggleMultiSelect();
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

// 渲染操作按钮组
const renderActionButtons = (row: any) => {
    const isLiked = music.getSongIsLike(row.id);
    const buttons = [
        h(NTooltip, { trigger: 'hover', placement: 'top', showArrow: false }, {
            trigger: () => h(NButton, {
                quaternary: true,
                circle: true,
                size: 'small',
                type: isLiked ? 'error' : 'default',
                onClick: (e: Event) => {
                    e.stopPropagation();
                    music.changeLikeList(row.id, !isLiked);
                }
            }, { 
                icon: () => h(NIcon, { component: Like, color: isLiked ? '#d03050' : undefined }) 
            }),
            default: () => isLiked ? '取消喜欢' : '喜欢'
        })
    ];

    // 桌面端显示下载按钮
    if (!isMobile.value) {
        buttons.push(
            h(NTooltip, { trigger: 'hover', placement: 'top', showArrow: false }, {
                trigger: () => h(NButton, {
                    quaternary: true,
                    circle: true,
                    size: 'small',
                    onClick: (e: Event) => {
                        e.stopPropagation();
                        handleDownload(row);
                    }
                }, { 
                    icon: () => h(NIcon, { component: Download }) 
                }),
                default: () => '下载'
            })
        );
    }

    return h('div', { 
        class: 'action-buttons',
        onClick: (e: Event) => e.stopPropagation() // 防止触发行动点击
    }, buttons);
};

const handleDownload = (row: any) => {
    const url = `/api/song/stream/${row.id}`;
    const link = document.createElement('a');
    link.href = url;
    link.download = `${row.title || 'audio'}.${row.format || 'mp3'}`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
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
      src: getCoverSrc(song.cover_url, song.id, song.album, song.picUrl),
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
        ...(isMobile.value ? [] : [{
            label: '下载',
            key: 'download',
            icon: renderIcon(Download)
        }]),
        {
            label: '更多操作',
            key: 'more',
            icon: renderIcon(More),
            children: [
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
        case 'download':
            handleDownload(song);
            break;
        case 'add-to-playlist':
            openAddToPlaylist([song.id]);
            break;
        case 'copy-link':
            const routeData = router.resolve({ name: 'song', params: { id: song.id } });
            const link = `${window.location.origin}${routeData.href}`;
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
      title: "",
      key: "index",
      width: 50,
      align: 'center',
      render: (row: any, index: number) => {
        const isCurrent = currentPlayingSong.value?.id === row.id;
        const playing = isPlaying.value;
        
        return h('div', { class: 'index-cell' }, [
            // 播放状态图标 (当歌曲是当前播放歌曲时显示)
            isCurrent ? h(NIcon, { 
                size: 18, 
                color: setting.themeColor,
                component: playing ? VolumeNotice : PauseOne 
            }) : h('span', { class: 'index-num' }, `${index + 1 + (props.page - 1) * props.pageSize}`),
            
            // 悬浮播放图标 (非当前播放歌曲时，悬浮显示播放)
            !isCurrent ? h(NIcon, { 
                class: 'hover-play-icon',
                size: 18,
                component: PlayOne
            }) : null
        ]);
      },
    },
    {
      title: "标题",
      key: "title",
      render: (row: any) => {
        const isCurrent = currentPlayingSong.value?.id === row.id;
        return h('div', { class: 'title-cell' }, [
            h('span', {
              style: { 
                  cursor: 'pointer', 
                  transition: 'all 0.3s', 
                  fontSize: '15px', 
                  fontWeight: isCurrent ? 'bold' : '500',
                  color: isCurrent ? setting.themeColor : 'inherit'
              },
              onClick: (e: Event) => {
                 e.stopPropagation();
                 router.push({ name: 'song', params: { id: row.id } });
              },
              class: 'song-title-link'
            }, row.title),
            // 操作按钮 (悬浮显示)
            renderActionButtons(row)
        ]);
      }
    },
    {
      title: "歌手",
      key: "artist_name",
      render: (row: any) => {
        // Handle array of artists or single artist name string
        let artistName = row.artist_name;
        let artistId = row.artist_id;

        if (row.artists && row.artists.length > 0) {
            return h('div', {
                class: 'artist-links-container',
                style: { 
                    whiteSpace: 'nowrap', 
                    overflow: 'hidden', 
                    textOverflow: 'ellipsis',
                    maxWidth: '100%'
                }
            }, row.artists.map((artist: any, index: number) => {
                const nodes = [];
                nodes.push(h('span', {
                    class: 'artist-link',
                    style: { cursor: 'pointer', transition: 'color 0.3s', fontSize: '13px', opacity: 0.75 },
                    onClick: (e: Event) => {
                        e.stopPropagation();
                        router.push({ name: 'artist', query: { id: artist.id } });
                    },
                    onMouseover: (e: Event) => {
                        (e.target as HTMLElement).style.color = setting.themeColor;
                        (e.target as HTMLElement).style.opacity = '1';
                    },
                    onMouseout: (e: Event) => {
                        (e.target as HTMLElement).style.color = 'inherit';
                        (e.target as HTMLElement).style.opacity = '0.75';
                    }
                }, artist.name));
                
                if (index < row.artists.length - 1) {
                    nodes.push(h('span', { style: { margin: '0 4px', opacity: 0.5, fontSize: '12px' } }, '/'));
                }
                return nodes;
            }).flat());
        }

        return h('span', {
            class: 'artist-link',
            style: { cursor: 'pointer', transition: 'color 0.3s', fontSize: '13px', opacity: 0.75 },
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
                (e.target as HTMLElement).style.opacity = '0.75';
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
      width: 80,
      align: 'right',
      render: (row: any) => h('span', { style: { opacity: 0.5, fontFamily: 'DM Mono, Monaco, monospace', fontSize: '13px', fontVariantNumeric: 'tabular-nums', fontWeight: 'bold' } }, formatDuration(row.duration)),
    },
  );

  if (viewMode.value === 'thumbnail') {
    // 缩略图模式下：
    // 1. 隐藏原来的标题列和歌手列
    // 2. 插入一个新的合并列，显示封面+标题+歌手

    // 移除原有标题和歌手列
    // 原顺序: [0:index, 1:title, 2:artist, 3:album, 4:duration]
    // 倒序删除以避免索引错乱
    baseColumns.splice(2, 1); // 移除歌手列
    baseColumns.splice(1, 1); // 移除标题列

    // 插入新列
    const insertIndex = isMultiSelectMode.value ? 2 : 1;
    baseColumns.splice(insertIndex, 0, {
      title: "歌曲",
      key: "song_info",
      // width: 'auto', // 自适应宽度
      render: (row: any) => {
        const isCurrent = currentPlayingSong.value?.id === row.id;
        
        // 封面容器 (含播放遮罩)
        const coverNode = h('div', {
            class: 'cover-container',
            style: { 
                position: 'relative', 
                width: '56px', 
                height: '56px',
                borderRadius: '8px',
                overflow: 'hidden',
                boxShadow: '0 4px 10px rgba(0,0,0,0.1)',
                cursor: 'pointer',
                flexShrink: 0
            },
            onClick: (e: Event) => {
                e.stopPropagation();
                // 播放逻辑
                const index = props.songs.findIndex(s => s.id === row.id);
                if (index !== -1) {
                    const tracks = mapSongsToPlayer(props.songs);
                    music.setPlaylists(tracks);
                    music.setPlaySongIndex(index);
                    music.setPlayState(true);
                }
            }
        }, [
            h(NImage, {
                src: getCoverSrc(row.cover_url, row.id, row.album, row.picUrl),
                fallbackSrc: '/images/logo/favicon.png',
                width: 56, 
                height: 56,
                lazy: true, 
                objectFit: 'cover',
                previewDisabled: true,
                style: { width: '100%', height: '100%', display: 'block' }
            }),
            // 播放遮罩
            h('div', { class: 'cover-overlay' }, [
                h(NIcon, { component: PlayOne, size: 24, color: '#fff' })
            ])
        ]);

        // 标题
        const titleNode = h('span', {
          style: { 
              cursor: 'pointer', 
              transition: 'all 0.3s', 
              fontSize: '16px', 
              fontWeight: isCurrent ? 'bold' : '600',
              color: isCurrent ? setting.themeColor : 'var(--n-text-color)',
              marginRight: '8px',
              lineHeight: '1.2',
              whiteSpace: 'nowrap',
              overflow: 'hidden',
              textOverflow: 'ellipsis',
              maxWidth: '100%',
              display: 'inline-block'
          },
          onClick: (e: Event) => {
             e.stopPropagation();
             // 使用 name 跳转更稳健
             router.push({ name: 'song', params: { id: row.id } });
          },
          class: 'song-title-link'
        }, row.title);

        // 歌手
        let artistName = row.artist_name;
        let artistId = row.artist_id;
        
        // Render artist node(s)
        let artistNode;

        if (row.artists && row.artists.length > 0) {
             // Multiple artists
             artistNode = h('div', {
                 style: {
                    whiteSpace: 'nowrap',
                    overflow: 'hidden',
                    textOverflow: 'ellipsis',
                    maxWidth: '100%',
                    display: 'inline-block',
                    fontSize: '14px',
                    fontWeight: '400',
                    color: 'var(--n-text-color)'
                 }
             }, row.artists.map((artist: any, index: number) => {
                const nodes = [];
                nodes.push(h('span', {
                    class: 'artist-link',
                    style: { cursor: 'pointer', transition: 'color 0.3s', opacity: 0.6 },
                    onClick: (e: Event) => {
                        e.stopPropagation();
                        router.push({ name: 'artist', query: { id: artist.id } });
                    },
                    onMouseover: (e: Event) => {
                        (e.target as HTMLElement).style.color = setting.themeColor;
                        (e.target as HTMLElement).style.opacity = '1';
                    },
                    onMouseout: (e: Event) => {
                        (e.target as HTMLElement).style.color = 'inherit';
                        (e.target as HTMLElement).style.opacity = '0.6';
                    }
                }, artist.name));
                
                if (index < row.artists.length - 1) {
                    nodes.push(h('span', { style: { margin: '0 4px', opacity: 0.4 } }, '/'));
                }
                return nodes;
            }).flat());
        } else {
             // Single/Unknown artist
             if (!artistName && row.artists && row.artists.length > 0) {
                 // Fallback if logic above failed (should not happen)
                 artistName = row.artists.map((a: any) => a.name).join(' / ');
                 artistId = row.artists[0]?.id;
             }
             
             artistNode = h('span', {
                class: 'artist-link',
                style: { 
                    cursor: 'pointer', 
                    transition: 'color 0.3s', 
                    fontSize: '14px', 
                    opacity: 0.6,
                    fontWeight: '400',
                    whiteSpace: 'nowrap',
                    overflow: 'hidden',
                    textOverflow: 'ellipsis',
                    maxWidth: '100%',
                    display: 'inline-block'
                },
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
                    (e.target as HTMLElement).style.opacity = '0.6';
                }
            }, artistName || "Unknown");
        }

        // 分隔符
        const dividerNode = h('span', {
            style: { margin: '0 4px', opacity: 0.4, fontSize: '12px' }
        }, '-');

        // 专辑
        const albumName = row.album_name || row.album_title || row.album?.name || row.album?.title || "Unknown Album";
        const albumId = row.album_id || row.album?.id;
        
        const albumNode = h('span', {
           class: 'album-link',
           style: { 
               fontSize: '12px', 
               opacity: 0.5,
               cursor: 'pointer',
               transition: 'color 0.3s',
               whiteSpace: 'nowrap',
               overflow: 'hidden',
               textOverflow: 'ellipsis',
               maxWidth: '100%',
               display: 'inline-block'
           },
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
               (e.target as HTMLElement).style.opacity = '0.5';
           }
        }, albumName);

        let firstLineChildren = [];
        let secondLineChildren = [];

        if (isMobile.value) {
            // 移动端：
            // 第一行：标题
            // 第二行：歌手 (不显示专辑)
            firstLineChildren = [titleNode];
            secondLineChildren = [artistNode];
        } else {
            // 桌面端：
            // 第一行：标题 - 歌手
            // 第二行：专辑 (移除 SQ 标签)
            firstLineChildren = [titleNode, dividerNode, artistNode];
            secondLineChildren = [albumNode];
        }
        
        // 组合：第一行
        const firstLine = h('div', { 
            style: { display: 'flex', alignItems: 'center', marginBottom: '6px', overflow: 'hidden' } 
        }, firstLineChildren);

        // 组合：第二行 
        const secondLine = h('div', {
            style: { display: 'flex', alignItems: 'center', overflow: 'hidden' }
        }, secondLineChildren);

        // 文本容器
        const textContainer = h('div', {
            style: { display: 'flex', flexDirection: 'column', justifyContent: 'center', flex: 1, minWidth: 0 }
        }, [firstLine, secondLine]);


        return h('div', {
            style: { display: 'flex', alignItems: 'center', gap: '16px', width: '100%' }
        }, [coverNode, textContainer, renderActionButtons(row)]);
      }
    });

    // 既然我们在主列里显示了专辑，是否还要保留单独的专辑列？
    // 缩略图模式通常信息密度较高，可以隐藏单独的专辑列，或者保留。
    // 为了美观，我们隐藏单独的专辑列
    // 原专辑列现在是索引 2 (因为前面移除了2个，插入了1个)
    // [index, song_info, album, duration]
    baseColumns.splice(insertIndex + 1, 1); // 移除专辑列
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
    onClick: () => {
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
        if (!tracks || tracks.length === 0) return;
        
        // Calculate correct play index (handling potential filtering in mapSongsToPlayer)
        let playIndex = index;
        if (tracks.length !== props.songs.length) {
            const clickedSong = props.songs[index];
            if (clickedSong && clickedSong.id) {
                // Find which occurrence of this song ID was clicked
                let occurrence = 0;
                for (let i = 0; i < index; i++) {
                    if (props.songs[i] && props.songs[i].id === clickedSong.id) occurrence++;
                }
                
                // Find corresponding index in filtered tracks
                let matchCount = 0;
                for (let i = 0; i < tracks.length; i++) {
                    if (tracks[i].id === clickedSong.id) {
                        if (matchCount === occurrence) {
                            playIndex = i;
                            break;
                        }
                        matchCount++;
                    }
                }
                if (playIndex === -1) playIndex = tracks.findIndex(t => t.id === clickedSong.id);
            }
        }
        
        if (playIndex === -1) playIndex = 0;

        music.setPlaylists(tracks);
        music.setPlaySongIndex(playIndex);
        music.setPlayState(true);
    },
    onContextmenu: (e: MouseEvent) => handleContextMenu(e, row)
  };
};

// 映射歌曲到播放器格式
const mapSongsToPlayer = (list: any[]) => {
    if (!list || !Array.isArray(list)) return [];
    return list.map(song => {
        if (!song) return null;
        return {
            ...song,
            name: song.title || 'Unknown Song',
            // Ensure artist format is consistent for player
            artist: song.artists || (song.artist_name ? [{ name: song.artist_name, id: song.artist_id || 0 }] : []),
            album: song.album || { 
                name: song.album_name || song.album_title || 'Unknown Album', // 优先使用 album_name
                id: song.album_id || 0, 
                picUrl: song.cover_url ? resolveCoverUrl(song.cover_url) : song.id ? getSongCover(song.id) : ''
            },
            picUrl: song.cover_url ? resolveCoverUrl(song.cover_url) : song.id ? getSongCover(song.id) : ''
        };
    }).filter(item => item !== null && item.id);
};

// 当前播放的歌曲
const currentPlayingSong = computed(() => music.getPlaySongData);
const isPlaying = computed(() => music.getPlayState);
</script>

<style scoped lang="scss">
.song-list-component {
    padding-bottom: 20px;
    position: relative;
    background: transparent;
}
.list-control {
    position: absolute;
    top: -56px; /* 调整位置 */
    right: 0;
    z-index: 100;
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0;
    pointer-events: none;
    
    .right {
        pointer-events: auto;
        background: rgba(255, 255, 255, 0.5);
        backdrop-filter: blur(8px);
        padding: 4px;
        border-radius: 8px;
        border: 1px solid rgba(255, 255, 255, 0.2);
    }
    
    .left {
       pointer-events: auto;
       display: flex;
       align-items: center;
    }
}

:deep(.n-data-table) {
  background: transparent !important;
  
  .n-data-table-th {
    background-color: transparent !important;
    border-bottom: 1px solid rgba(0, 0, 0, 0.04);
    font-weight: 600;
    font-size: 13px;
    opacity: 0.5;
    padding: 12px 16px;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }
  
  .n-data-table-td {
    background-color: transparent !important;
    border-bottom: 1px solid rgba(0, 0, 0, 0.02) !important;
    padding: 12px 16px; 
    vertical-align: middle;
    transition: all 0.2s ease;
  }

  .n-data-table-tr {
      border-radius: 12px;
      transition: all 0.2s cubic-bezier(0.25, 0.8, 0.25, 1);
      
      &:hover {
         transform: scale(1.002);
         box-shadow: 0 4px 12px rgba(0, 0, 0, 0.03);
         z-index: 1;
         position: relative;
         
         .n-data-table-td {
             background-color: rgba(255, 255, 255, 0.6) !important;
             backdrop-filter: blur(4px);
         }
      }
  }

  .song-row {
     /* 移除旧的 hover 背景，使用 tr 的新效果 */
     &:hover {
        /* 隐藏序号，显示播放按钮 */
        .index-num {
            display: none;
        }
        .hover-play-icon {
            display: inline-flex !important;
        }
     }
  }

  /* Index Cell Styles */
  .index-cell {
      display: flex;
      justify-content: center;
      align-items: center;
      width: 100%;
      height: 100%;
  }

  .index-num {
      opacity: 0.4;
      font-size: 13px;
      font-family: 'DM Mono', monospace;
      font-weight: 500;
  }

  .hover-play-icon {
      display: none !important;
      opacity: 0.8;
      color: var(--n-color-primary);
      filter: drop-shadow(0 2px 4px rgba(var(--n-color-primary-rgb), 0.3));
  }
}

.action-buttons {
    display: flex;
    gap: 8px;
    margin-left: auto;
    opacity: 0;
    transform: translateX(10px);
    transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
    align-items: center;
    background: rgba(255, 255, 255, 0.6);
    padding: 4px 8px;
    border-radius: 20px;
    backdrop-filter: blur(4px);
}

:deep(.song-row:hover) .action-buttons {
    opacity: 1;
    transform: translateX(0);
}

:deep(.cover-overlay) {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.3);
    display: flex;
    justify-content: center;
    align-items: center;
    opacity: 0;
    transform: scale(0.9);
    transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
    backdrop-filter: blur(2px);
}

:deep(.cover-container:hover) .cover-overlay {
    opacity: 1;
    transform: scale(1);
}

/* 链接样式优化 */
:deep(.song-title-link) {
    position: relative;
    display: inline-block;
    
    &::after {
        content: '';
        position: absolute;
        bottom: -2px;
        left: 0;
        width: 0;
        height: 2px;
        background: var(--n-color-primary);
        transition: width 0.3s ease;
        opacity: 0.5;
    }
    
    &:hover::after {
        width: 100%;
    }
}
</style>
