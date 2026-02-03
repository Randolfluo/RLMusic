<template>
  <Transition name="playlist">
    <div class="play-list" v-show="music.showPlayList" @click.stop>
      <div class="header">
        <div class="left">
          <span class="title">当前播放</span>
          <span class="num">({{ music.getPlaylists.length }})</span>
        </div>
        <div class="right">
          <div class="clear" @click="clearPlayList">
            <n-icon :component="Delete" />
            <span>清空列表</span>
          </div>
          <n-icon
            class="close"
            size="20"
            :component="CloseSmall"
            @click="music.showPlayList = false"
          />
        </div>
      </div>
      <n-scrollbar class="list-scroll">
        <div class="list" v-if="music.getPlaylists[0]">
          <div
            v-for="(item, index) in music.getPlaylists"
            :key="item.id"
            :class="
              index == music.getPlaySongIndex ? 'item active' : 'item'
            "
            @dblclick="music.setPlaySongIndex(index)"
          >
            <div class="name">
              <span class="song-name text-hidden">{{ item.name }}</span>
              <span class="line" v-if="item.artist">-</span>
              <div class="artist text-hidden" v-if="item.artist">
                <span
                  class="ar"
                  v-for="(a, i) in item.artist"
                  :key="a.id"
                >
                  {{ a.name }}
                  <span v-if="i !== item.artist.length - 1">/</span>
                </span>
              </div>
            </div>
            <n-icon
              class="del"
              size="18"
              :component="CloseSmall"
              @click.stop="delSong(index)"
            />
          </div>
        </div>
        <div class="empty" v-else>
          <span>暂无歌曲</span>
        </div>
      </n-scrollbar>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { musicStore } from "@/store";
import { Delete, CloseSmall } from "@icon-park/vue-next";

const music = musicStore();

// 清空列表
const clearPlayList = () => {
  music.setPlaylists([]);
  music.setPlaySongIndex(0);
};

// 删除单曲
const delSong = (index: number) => {
  const list = music.getPlaylists;
  list.splice(index, 1);
  music.setPlaylists(list);
  if (index < music.getPlaySongIndex) {
    music.setPlaySongIndex(music.getPlaySongIndex - 1);
  } else if (index == music.getPlaySongIndex) {
    music.setPlaySongIndex(index >= list.length ? 0 : index);
  }
};
</script>

<style lang="scss" scoped>
.play-list {
  position: fixed;
  bottom: 80px;
  right: 20px;
  width: 350px;
  height: 500px;
  background-color: var(--n-color-modal);
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.1);
  border-radius: 12px;
  z-index: 999;
  display: flex;
  flex-direction: column;
  backdrop-filter: blur(20px);
  overflow: hidden;
  .header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px 20px;
    border-bottom: 1px solid var(--n-border-color);
    .left {
      display: flex;
      align-items: center;
      .title {
        font-size: 16px;
        font-weight: bold;
      }
      .num {
        font-size: 12px;
        margin-left: 4px;
        opacity: 0.6;
        margin-top: 2px;
      }
    }
    .right {
      display: flex;
      align-items: center;
      .clear {
        display: flex;
        align-items: center;
        cursor: pointer;
        opacity: 0.6;
        transition: all 0.3s;
        margin-right: 12px;
        font-size: 13px;
        &:hover {
          opacity: 1;
        }
        .n-icon {
          margin-right: 4px;
        }
      }
      .close {
        cursor: pointer;
        opacity: 0.6;
        transition: all 0.3s;
        &:hover {
          opacity: 1;
        }
      }
    }
  }
  .list-scroll {
    flex: 1;
    .list {
      padding: 8px 0;
      .item {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 10px 20px;
        cursor: pointer;
        transition: all 0.3s;
        border-radius: 8px;
        margin: 0 8px;
        &:hover {
          background-color: var(--n-close-color-hover);
        }
        &.active {
          color: var(--n-primary-color);
          background-color: var(--n-close-color-pressed);
          .name {
            .line {
              color: var(--n-primary-color);
            }
            .artist {
              color: var(--n-primary-color);
            }
          }
        }
        .name {
          flex: 1;
          display: flex;
          align-items: center;
          overflow: hidden;
          font-size: 14px;
          .song-name {
            font-weight: 500;
          }
          .line {
            margin: 0 4px;
            opacity: 0.6;
          }
          .artist {
            font-size: 12px;
            opacity: 0.6;
            display: flex;
            align-items: center;
          }
        }
        .del {
          margin-left: 10px;
          opacity: 0;
          cursor: pointer;
          transition: all 0.3s;
          &:hover {
            color: var(--n-error-color);
          }
        }
        &:hover {
          .del {
            opacity: 1;
          }
        }
      }
    }
    .empty {
      height: 100%;
      display: flex;
      align-items: center;
      justify-content: center;
      opacity: 0.5;
    }
  }
}

.playlist-enter-active,
.playlist-leave-active {
  transition: all 0.3s ease;
}

.playlist-enter-from,
.playlist-leave-to {
  opacity: 0;
  transform: translateY(20px);
}
</style>
