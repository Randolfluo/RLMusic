<template>
  <div 
    class="ws-debug-wrapper" 
    :style="{ left: position.x + 'px', top: position.y + 'px', width: size.w + 'px', height: size.h + 'px' }"
  >
    <n-card 
      title="WebSocket Debugger" 
      closable 
      @close="$emit('close')"
      class="ws-debug-card"
      size="small"
      :content-style="{ padding: '8px', display: 'flex', flexDirection: 'column', height: '100%' }"
    >
      <template #header-extra>
        <div class="drag-handle" @mousedown="startDrag">
          <n-icon size="18"><Move /></n-icon>
        </div>
      </template>

      <n-tabs type="line" size="small" animated style="flex: 1; display: flex; flex-direction: column;">
        <n-tab-pane name="monitor" tab="Monitor" style="flex: 1; display: flex; flex-direction: column; overflow: hidden;">
          <n-space vertical style="height: 100%; display: flex; flex-direction: column;">
            <!-- Actions & Filter -->
            <n-space align="center" justify="space-between">
               <n-space size="small">
                 <n-button size="tiny" @click="sendHello">Hello</n-button>
                 <n-button size="tiny" @click="sendHeartbeat">Ping</n-button>
                 <n-button size="tiny" @click="sendGetRoomList">Rooms</n-button>
                 <n-button size="tiny" @click="sendTimeSync">Time</n-button>
               </n-space>
               <n-space size="small">
                  <n-checkbox v-model:checked="hideHeartbeat" size="small">Hide Ping</n-checkbox>
                  <n-button size="tiny" type="error" @click="logs = []">Clear</n-button>
               </n-space>
            </n-space>
            
            <!-- Log Window -->
            <div class="log-window" ref="logRef">
              <div v-for="(log, i) in filteredLogs" :key="i" class="log-entry">
                <div class="log-header">
                  <span class="time">[{{ log.time }}]</span>
                  <n-tag size="tiny" :type="log.direction === 'TX' ? 'info' : 'success'">
                    {{ log.direction }}
                  </n-tag>
                  <span class="msg-type">{{ log.type }}</span>
                </div>
                <div class="log-payload">
                   <pre v-if="prettyPrint">{{ JSON.stringify(log.payload, null, 2) }}</pre>
                   <span v-else>{{ JSON.stringify(log.payload) }}</span>
                </div>
              </div>
            </div>
             <n-space justify="end">
                <n-checkbox v-model:checked="prettyPrint" size="small">Pretty Print</n-checkbox>
             </n-space>
          </n-space>
        </n-tab-pane>
        
        <n-tab-pane name="sender" tab="Manual Sender">
          <n-space vertical>
            <n-input v-model:value="customType" placeholder="Message Type (e.g. CHAT)" />
            <n-input
              v-model:value="customPayload"
              type="textarea"
              placeholder='Payload JSON (e.g. {"content": "hello"})'
              :rows="8"
            />
            <n-button type="primary" block @click="sendCustomMessage" :disabled="!isValidJson">
              Send Message
            </n-button>
          </n-space>
        </n-tab-pane>
      </n-tabs>
      
      <!-- Resize Handle -->
      <div class="resize-handle" @mousedown="startResize"></div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue';
import { NCard, NSpace, NButton, NTag, NTabs, NTabPane, NInput, NCheckbox, NIcon, useMessage } from 'naive-ui';
import { Move } from '@vicons/carbon';
import { socket, MsgType, type WSMessage } from '@/core/realtime/socket';
import { userStore } from '@/store';

const user = userStore();
const message = useMessage();

// State
const logs = ref<{ time: string; direction: 'RX' | 'TX'; type: string; payload: any }[]>([]);
const logRef = ref<HTMLElement | null>(null);
const hideHeartbeat = ref(true);
const prettyPrint = ref(true);

// Manual Sender State
const customType = ref('');
const customPayload = ref('{\n  \n}');

// Window State
const position = ref({ x: window.innerWidth - 420, y: window.innerHeight - 520 });
const size = ref({ w: 400, h: 500 });

// Filtered Logs
const filteredLogs = computed(() => {
  if (hideHeartbeat.value) {
    return logs.value.filter(l => l.type !== MsgType.HEARTBEAT && l.type !== 'HEARTBEAT');
  }
  return logs.value;
});

const isValidJson = computed(() => {
  try {
    JSON.parse(customPayload.value);
    return true;
  } catch (e) {
    return false;
  }
});

const addLog = (direction: 'RX' | 'TX', type: string, payload: any) => {
  logs.value.push({
    time: new Date().toLocaleTimeString(),
    direction,
    type,
    payload
  });
  
  // Limit logs to prevent memory issues
  if (logs.value.length > 200) {
      logs.value.shift();
  }
  
  nextTick(() => {
    if (logRef.value) {
      logRef.value.scrollTop = logRef.value.scrollHeight;
    }
  });
};

// Monitor incoming messages
const handleMessage = (msg: WSMessage) => {
  addLog('RX', msg.type, msg.payload);
};

// Actions
const sendHello = () => {
  const payload = { id: user.getUserData.userId, nickname: user.getUserData.nickname };
  socket.sendHello(payload);
  addLog('TX', MsgType.HELLO, payload);
};

const sendHeartbeat = () => {
  socket.sendHeartbeat();
  addLog('TX', MsgType.HEARTBEAT, { client_now: Date.now() });
};

const sendGetRoomList = () => {
  socket.sendGetRoomList();
  addLog('TX', MsgType.GET_ROOM_LIST, {});
};

const sendTimeSync = () => {
  socket.sendTimeSyncReq();
  addLog('TX', MsgType.TIME_SYNC_REQ, { client_now: Date.now() });
};

const sendCustomMessage = () => {
    if (!customType.value) {
        message.error("Type is required");
        return;
    }
    try {
        const payload = JSON.parse(customPayload.value);
        socket.send(customType.value, payload);
        addLog('TX', customType.value, payload);
        message.success("Sent!");
    } catch(e) {
        message.error("Invalid JSON");
    }
}

// Dragging Logic
const isDragging = ref(false);
const dragOffset = ref({ x: 0, y: 0 });

const startDrag = (e: MouseEvent) => {
  isDragging.value = true;
  dragOffset.value = {
    x: e.clientX - position.value.x,
    y: e.clientY - position.value.y
  };
  document.addEventListener('mousemove', onDrag);
  document.addEventListener('mouseup', stopDrag);
};

const onDrag = (e: MouseEvent) => {
  if (!isDragging.value) return;
  position.value = {
    x: e.clientX - dragOffset.value.x,
    y: e.clientY - dragOffset.value.y
  };
};

const stopDrag = () => {
  isDragging.value = false;
  document.removeEventListener('mousemove', onDrag);
  document.removeEventListener('mouseup', stopDrag);
};

// Resizing Logic
const isResizing = ref(false);
const startResize = (e: MouseEvent) => {
  isResizing.value = true;
  document.addEventListener('mousemove', onResize);
  document.addEventListener('mouseup', stopResize);
  e.preventDefault(); // Prevent text selection
};

const onResize = (e: MouseEvent) => {
  if (!isResizing.value) return;
  const newW = e.clientX - position.value.x;
  const newH = e.clientY - position.value.y;
  size.value = {
    w: Math.max(300, newW),
    h: Math.max(300, newH)
  };
};

const stopResize = () => {
  isResizing.value = false;
  document.removeEventListener('mousemove', onResize);
  document.removeEventListener('mouseup', stopResize);
};

onMounted(() => {
  socket.onMessage(handleMessage);
});

onUnmounted(() => {
  socket.offMessage(handleMessage);
});
</script>

<style scoped>
.ws-debug-wrapper {
  position: fixed;
  z-index: 9999;
  box-shadow: 0 4px 20px rgba(0,0,0,0.2);
  display: flex;
  flex-direction: column;
}

.ws-debug-card {
  height: 100%;
  opacity: 0.98;
}

.drag-handle {
  cursor: move;
  padding: 4px;
  display: flex;
  align-items: center;
  color: var(--n-text-color-3);
}

.drag-handle:hover {
    color: var(--n-primary-color);
}

.log-window {
  flex: 1;
  min-height: 0; /* Important for flex container scrolling */
  overflow-y: auto;
  background: rgba(0, 0, 0, 0.03);
  border: 1px solid var(--n-border-color);
  border-radius: 4px;
  padding: 8px;
  font-family: monospace;
  font-size: 11px;
  margin-top: 8px;
}

.log-entry {
  margin-bottom: 8px;
  border-bottom: 1px dashed rgba(0,0,0,0.1);
  padding-bottom: 4px;
}

.log-header {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 2px;
}

.time {
  color: #999;
}

.msg-type {
  font-weight: bold;
  color: var(--n-text-color);
}

.log-payload {
  white-space: pre-wrap;
  word-break: break-all;
  color: var(--n-text-color-3);
  padding-left: 8px;
}

.log-payload pre {
    margin: 0;
    font-family: inherit;
}

.resize-handle {
  position: absolute;
  bottom: 0;
  right: 0;
  width: 15px;
  height: 15px;
  cursor: nwse-resize;
  background: linear-gradient(135deg, transparent 50%, var(--n-border-color) 50%);
  border-bottom-right-radius: 4px;
}
</style>
