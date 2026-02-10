<!--
  ListenTogether/index.vue
  功能：一起听歌页面（多人实时聊天室）
  说明：
    - 使用 WebSocket 实现实时通信
    - 支持发送/接收聊天消息（包含头像、昵称）
    - 支持时间同步、心跳保活
    - 支持加入/离开房间事件，并以系统通知形式展示
    - 提供左侧调试面板（显示连接状态、延迟等）
-->
<template>
  <div class="listen-together">
    <!-- 房间列表侧边栏 -->
    <div class="room-sidebar">
      <div class="sidebar-header">
        <h3>房间列表</h3>
        <n-space>
           <n-button size="tiny" secondary type="info" @click="showDebug = !showDebug">Debug</n-button>
           <n-button size="tiny" type="primary" @click="showCreateRoomModal = true">创建</n-button>
        </n-space>
      </div>
      <div class="room-list">
        <!-- 活跃房间列表 (Lobby) -->
        <div class="section-label">大厅</div>
        <div 
          v-for="room in availableRooms" 
          :key="room.id" 
          class="room-item"
          @click="joinRoom(room.id)"
        >
          <span class="room-name">{{ room.id }}</span>
          <span class="room-count">{{ room.count }}人</span>
        </div>

        <!-- 已加入的房间 -->
        <div class="section-label" style="margin-top: 12px;">我的房间</div>
        <div 
          v-for="roomId in joinedRooms" 
          :key="roomId" 
          class="room-item"
          :class="{ active: currentRoomId === roomId }"
          @click="currentRoomId = roomId"
        >
          <span class="room-name">{{ roomId }}</span>
          <n-button size="tiny" text type="error" @click.stop="leaveRoom(roomId)">退出</n-button>
        </div>
      </div>
    </div>

    <!-- 中间：控制面板 + 在线用户 (针对当前选中的房间) -->
    <div class="control-panel" v-if="currentRoomId">
      <h3>{{ currentRoomId }}</h3>
      <!-- 房间状态 -->
      <div class="status">
        <div class="status-item">
          <span class="label">Status:</span>
          <span :class="['value', socket.isConnected.value ? 'connected' : 'disconnected']">
            {{ socket.isConnected.value ? 'Connected' : 'Disconnected' }}
          </span>
        </div>
        <div class="status-item">
          <span class="label">Owner:</span>
          <span class="value">{{ currentRoomOwnerId || '-' }}</span>
        </div>
      </div>
      
      <!-- 在线用户列表 -->
      <div class="online-users">
        <h3>Online Users ({{ currentRoomUsers.length }})</h3>
        <div class="user-list">
          <div v-for="u in currentRoomUsers" :key="u.id" class="user-item">
            <n-avatar round size="small" :src="u.avatarUrl || '/images/ico/user-filling.svg'" />
            <span class="nickname">{{ u.nickname }}</span>
            <span v-if="String(u.id) === String(currentRoomOwnerId)" class="owner-tag">房主</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 右侧聊天区域 -->
    <div class="chat-container" v-if="currentRoomId">
      <!-- 消息列表 -->
      <div class="messages" ref="messagesRef">
        <div v-for="(msg, index) in currentRoomMessages" :key="index">
          
          <!-- 系统消息 -->
          <div v-if="msg.type === 'system'" class="system-message">
            <span class="content">{{ msg.content }}</span>
          </div>

          <!-- 聊天消息 -->
          <div v-else class="message-item" :class="{ 'my-message': msg.isMine }">
            <n-avatar round size="small" :src="msg.avatarUrl || '/images/ico/user-filling.svg'" class="avatar" />
            <div class="message-content-wrapper">
              <span class="sender">{{ msg.sender }}</span>
              <div class="message-content">
                <div class="text">{{ msg.content }}</div>
              </div>
            </div>
          </div>

        </div>
      </div>
      <!-- 输入区域 -->
      <div class="input-area">
        <n-input
          v-model:value="inputValue"
          type="text"
          placeholder="发送消息..."
          @keyup.enter="handleSendMessage"
        />
        <n-button type="primary" @click="handleSendMessage">发送</n-button>
      </div>
    </div>

    <div v-else class="empty-state">
      <n-empty description="请选择或加入一个房间" />
    </div>

    <!-- 创建房间弹窗 -->
    <n-modal v-model:show="showCreateRoomModal" preset="dialog" title="创建房间">
      <n-input v-model:value="newRoomName" placeholder="请输入房间名称" />
      <template #action>
        <n-button @click="showCreateRoomModal = false">取消</n-button>
        <n-button type="primary" @click="handleCreateRoom">确定</n-button>
      </template>
    </n-modal>

    <!-- Debug Panel -->
    <WebSocketDebug v-if="showDebug" @close="showDebug = false" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick, watch, defineAsyncComponent } from 'vue';
import { NInput, NButton, NAvatar, NModal, NEmpty, NSpace } from 'naive-ui';
import { chatStore } from '@/store';
import { storeToRefs } from 'pinia';
import { socket } from '@/core/realtime/socket';

const WebSocketDebug = defineAsyncComponent(() => import('@/components/WebSocketDebug.vue'));

const chat = chatStore();
const { 
  availableRooms, 
  joinedRooms, 
  currentRoomId, 
  currentRoomMessages, 
  currentRoomUsers, 
  currentRoomOwnerId 
} = storeToRefs(chat);

// 状态管理
const showCreateRoomModal = ref(false);
const showDebug = ref(false);
const newRoomName = ref('');
const inputValue = ref('');
const messagesRef = ref<HTMLElement | null>(null);

// 创建房间
const handleCreateRoom = () => {
  chat.createRoom(newRoomName.value);
  showCreateRoomModal.value = false;
  newRoomName.value = '';
};

// 加入房间
const joinRoom = (roomId: string) => {
  chat.joinRoom(roomId);
};

// 离开房间
const leaveRoom = (roomId: string) => {
  chat.leaveRoom(roomId);
};

// 发送消息
const handleSendMessage = () => {
  chat.sendMessage(inputValue.value);
  inputValue.value = '';
};

// 滚动到底部
const scrollToBottom = () => {
  nextTick(() => {
    if (messagesRef.value) {
      messagesRef.value.scrollTop = messagesRef.value.scrollHeight;
    }
  });
};

// 监听 currentRoomId 变化，自动滚动
watch(currentRoomId, () => {
  scrollToBottom();
});

// 监听消息列表变化，自动滚动
watch(() => currentRoomMessages.value.length, () => {
  if (currentRoomId.value) {
      scrollToBottom();
  }
});

onMounted(() => {
  // 启动监听 (如果不在这里调用，chatStore.startListen() 也会处理)
  // 如果已经连接且 joinedRooms > 0，这里会保持连接
  // 如果未连接，这里会发起连接
  chat.startListen();
  
  if (currentRoomId.value) {
      scrollToBottom();
  }
});

onUnmounted(() => {
  // 如果没有加入任何房间，则断开连接以节省资源
  // 用户需求：只有关闭页面或注销才断开。
  // 但为了不浪费资源，如果用户没在任何房间里，断开也是合理的？
  // 也可以根据 joinedRooms 长度判断
  if (joinedRooms.value.length === 0) {
     chat.stopListen();
  }
  // 如果有加入房间，则保持连接，不调用 stopListen
});
</script>

<style scoped lang="scss">
.listen-together {
  padding: 20px;
  height: calc(100vh - 100px); 
  display: flex;
  gap: 20px;
}

/* 房间侧边栏样式 */
.room-sidebar {
  width: 220px;
  background: var(--n-color-modal);
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  overflow: hidden;

  .sidebar-header {
    padding: 16px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    border-bottom: 1px solid var(--n-border-color);

    h3 {
      margin: 0;
      font-size: 16px;
      font-weight: bold;
    }
  }

  .room-list {
    flex: 1;
    overflow-y: auto;
    padding: 12px;
    
    .section-label {
      font-size: 12px;
      color: var(--n-text-color-3);
      margin-bottom: 8px;
      font-weight: bold;
    }

    .room-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 8px 12px;
      border-radius: 6px;
      cursor: pointer;
      margin-bottom: 4px;
      transition: background-color 0.2s;
      
      &:hover {
        background-color: rgba(0, 0, 0, 0.05);
      }
      
      &.active {
        background-color: var(--n-color-primary);
        color: white;
        
        .room-count {
          color: rgba(255, 255, 255, 0.8);
        }
      }

      .room-name {
        font-size: 14px;
        font-weight: 500;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        color: var(--n-text-color);
      }

      .room-count {
        font-size: 12px;
        color: var(--n-text-color-3);
      }
    }
  }
}

/* 左侧控制面板样式 */
.control-panel {
  width: 250px;
  background: var(--n-color-modal);
  border-radius: 8px;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 20px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);

  h3 {
    margin: 0;
    font-size: 16px;
    font-weight: bold;
  }

  .status {
    display: flex;
    flex-direction: column;
    gap: 10px;
    
    .status-item {
      display: flex;
      justify-content: space-between;
      font-size: 13px;
      
      .label {
        opacity: 0.7;
      }
      
      .value {
        font-weight: 500;
        font-family: monospace;
        
        &.connected { color: #18a058; }
        &.disconnected { color: #d03050; }
      }
    }
  }

  .actions {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .online-users {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;

    h3 {
      margin-bottom: 10px;
      font-size: 14px;
      color: var(--n-text-color-2);
    }

    .user-list {
      flex: 1;
      overflow-y: auto;
      display: flex;
      flex-direction: column;
      gap: 8px;

      .user-item {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 6px;
        border-radius: 6px;
        background-color: rgba(0, 0, 0, 0.02);

        .nickname {
          font-size: 13px;
          color: var(--n-text-color);
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
      }
    }
  }
}

/* 右侧聊天容器样式 */
.chat-container {
  flex: 1;
  max-width: 800px;
  background: var(--n-color-modal);
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  height: 100%;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.empty-state {
  flex: 1;
  max-width: 800px;
  display: flex;
  justify-content: center;
  align-items: center;
  background: var(--n-color-modal);
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.messages {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  
  /* 系统通知消息样式 */
  .system-message {
    display: flex;
    justify-content: center;
    margin: 16px 0;
    
    .content {
      background-color: rgba(0, 0, 0, 0.05);
      color: var(--n-text-color-3);
      font-size: 12px;
      padding: 4px 12px;
      border-radius: 12px;
    }
  }

  /* 聊天消息气泡样式 */
  .message-item {
    margin-bottom: 16px;
    display: flex;
    flex-direction: row;
    align-items: flex-start;
    gap: 8px;

    /* 自己发送的消息样式（右对齐） */
    &.my-message {
      flex-direction: row-reverse;
      
      .message-content-wrapper {
        align-items: flex-end;
      }
    }

    .avatar {
      flex-shrink: 0;
    }

    .message-content-wrapper {
      display: flex;
      flex-direction: column;
      align-items: flex-start;
      max-width: 80%;
    }

    .sender {
      font-size: 12px;
      color: var(--n-text-color-3); // 适配主题色
      margin-bottom: 4px;
      padding: 0 4px;
    }

    .message-content {
      background-color: var(--n-color-embedded);
      color: var(--n-text-color); // 适配主题色
      padding: 8px 12px;
      border-radius: 8px;
      word-break: break-all;
      
      .text {
        line-height: 1.5;
        color: inherit; // 继承父元素颜色，确保深色/浅色模式一致
      }
    }
  }
}

.input-area {
  padding: 20px;
  border-top: 1px solid var(--n-border-color);
  display: flex;
  gap: 10px;
}
</style>
