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
import { musicStore } from "@/store";

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
const viewMode = ref<'thumbnail' | 'concise'>('thumbnail');

const columns = computed(() => {
  const baseColumns: any[] = [
    {
      title: "#",
      key: "index",
      width: 60,
      render: (_: any, index: number) => index + 1 + (props.page - 1) * props.pageSize,
    },
    {
      title: "标题",
      key: "title",
      render: (row: any) => {
        return h('span', {
          style: { cursor: 'pointer' },
          onClick: (e: Event) => {
             e.stopPropagation();
             router.push(`/song/${row.id}`);
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
            style: { cursor: 'pointer', transition: 'color 0.3s' },
            onClick: (e: Event) => {
                e.stopPropagation();
                if (artistId) router.push({ name: 'artist', query: { id: artistId } });
            },
            onMouseover: (e: Event) => (e.target as HTMLElement).style.color = 'var(--n-color-primary)',
            onMouseout: (e: Event) => (e.target as HTMLElement).style.color = 'inherit'
        }, artistName || "Unknown Artist")
      },
    },
    {
      title: "专辑",
      key: "album_title",
      render: (row: any) => {
        const albumName = row.album_title || row.album?.name || row.album?.title || "Unknown Album";
        const albumId = row.album_id || row.album?.id;
        
        return h('span', {
           class: 'album-link',
           style: { cursor: 'pointer', transition: 'color 0.3s' },
           onClick: (e: Event) => {
               e.stopPropagation();
               if (albumId) router.push({ path: '/album', query: { id: albumId } });
           },
           onMouseover: (e: Event) => (e.target as HTMLElement).style.color = 'var(--n-color-primary)',
           onMouseout: (e: Event) => (e.target as HTMLElement).style.color = 'inherit'
        }, albumName);
      },
    },
    {
      title: "时长",
      key: "duration",
      width: 100,
      render: (row: any) => formatDuration(row.duration),
    },
  ];

  if (viewMode.value === 'thumbnail') {
    baseColumns.splice(1, 0, {
      title: "封面",
      key: "cover",
      width: 80,
      render: (row: any) => {
        return h(NImage, {
          src: row.cover_url || (row.album ? row.album.picUrl : null) || row.picUrl || `/api/song/cover/${row.id}`,
          fallbackSrc: '/images/logo/favicon.png',
          width: 50,
          height: 50,
          objectFit: 'cover',
          style: { borderRadius: '4px', verticalAlign: 'middle' },
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

const rowProps = (row: any, index: number) => {
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
            name: song.album_title, 
            id: song.album_id, 
            picUrl: song.cover_url || `/api/song/cover/${song.id}` 
        },
        picUrl: song.cover_url || `/api/song/cover/${song.id}`
    }));
};
</script>

<style scoped lang="scss">
.list-control {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
    padding-right: 12px;
}

:deep(.n-data-table .n-data-table-td) {
  background-color: transparent !important;
}
:deep(.n-data-table .n-data-table-tr:hover .n-data-table-td) {
  background-color: rgba(255, 255, 255, 0.05) !important;
}
</style>
