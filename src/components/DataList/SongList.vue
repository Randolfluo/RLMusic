<template>
  <div class="song-list-component">
    <div class="list-control">
      <div class="left">
        <!-- Slot for left side controls (e.g. Play All button if moved here, or just empty) -->
        <slot name="controls"></slot>
      </div>
      <div class="right">
        <n-button-group size="small">
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

    <n-data-table
      :columns="columns"
      :data="songs"
      :bordered="false"
      :row-props="rowProps"
      :row-class-name="() => 'song-row'"
      :loading="loading"
    />
    
    <div v-if="!loading && songs.length === 0" class="empty">
      <n-empty description="暂无歌曲" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h } from "vue";
import { useRouter } from "vue-router";
import { NButton, NButtonGroup, NIcon, NImage, NTooltip, NDataTable, NEmpty } from "naive-ui";
import { HamburgerButton, Pic } from "@icon-park/vue-next";
import { musicStore, settingStore } from "@/store";

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
  }
});

const router = useRouter();
const music = musicStore();
const setting = settingStore();
const viewMode = ref<'thumbnail' | 'concise'>('thumbnail');

const columns = computed(() => {
  const baseColumns: any[] = [
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
  ];

  if (viewMode.value === 'thumbnail') {
    baseColumns.splice(1, 0, {
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

const formatDuration = (seconds: number) => {
  if (!seconds) return "00:00";
  const m = Math.floor(seconds / 60);
  const s = Math.floor(seconds % 60);
  return `${m.toString().padStart(2, "0")}:${s.toString().padStart(2, "0")}`;
};

const rowProps = (_: any, index: number) => {
  return {
    style: "cursor: pointer;",
    onClick: () => {
        const tracks = mapSongsToPlayer(props.songs);
        music.setPlaylists(tracks);
        music.setPlaySongIndex(index);
        music.setPlayState(true);
    }
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
    justify-content: flex-end;
    align-items: center;
    padding: 0;
    pointer-events: none;
    
    .right {
        pointer-events: auto;
    }
    
    /* 如果 slot 有内容，可能需要额外处理，这里默认 float right */
    .left {
       display: none; /* 暂时隐藏左侧空 slot */
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
