<template>
  <div class="log-viewer">
    <!-- Toolbar -->
    <div class="toolbar">
      <div class="toolbar-left">
        <n-select
          v-model:value="filterLevel"
          :options="levelOptions"
          size="small"
          style="width: 120px"
        />
        <n-input
          v-model:value="searchText"
          placeholder="搜索日志..."
          size="small"
          clearable
          style="width: 200px"
        >
          <template #prefix>
            <n-icon :component="Search" />
          </template>
        </n-input>
      </div>
      <div class="toolbar-right">
        <n-tag :type="statusType" size="small">{{ statusLabel }}</n-tag>
        <n-button size="small" @click="togglePause">
          <template #icon><n-icon :component="paused ? Play : Pause" /></template>
          {{ paused ? '恢复' : '暂停' }}
        </n-button>
        <n-switch v-model:value="autoScroll" size="small">
          <template #checked>自动滚动</template>
          <template #unchecked>手动滚动</template>
        </n-switch>
        <n-button size="small" @click="clearLogs">
          <template #icon><n-icon :component="Delete" /></template>
          清空
        </n-button>
      </div>
    </div>

    <!-- Log list -->
    <div ref="logContainer" class="log-container" @scroll="onScroll">
      <div v-if="filteredLogs.length === 0" class="log-empty">
        暂无日志
      </div>
      <div
        v-for="(entry, idx) in filteredLogs"
        :key="idx"
        class="log-line"
        :class="'level-' + entry.level"
      >
        <span class="log-time">{{ entry.time }}</span>
        <span class="log-level">{{ entry.level.toUpperCase() }}</span>
        <span class="log-msg">{{ entry.message }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { Search, Play, Pause, Delete } from '@icon-park/vue-next'
import { NSelect, NInput, NButton, NSwitch, NTag, NIcon } from 'naive-ui'
import { apiBaseURL } from '@/utils/request'

interface LogEntry {
  time: string
  level: string
  message: string
}

const logs = ref<LogEntry[]>([])
const filterLevel = ref<string>('')
const searchText = ref('')
const autoScroll = ref(true)
const paused = ref(false)
const statusLabel = ref('连接中...')
const statusType = ref<'warning' | 'success' | 'error'>('warning')
const logContainer = ref<HTMLElement | null>(null)

let abortController: AbortController | null = null

const levelOptions = [
  { label: '全部级别', value: '' },
  { label: 'DEBUG', value: 'DEBUG' },
  { label: 'INFO', value: 'INFO' },
  { label: 'WARN', value: 'WARN' },
  { label: 'ERROR', value: 'ERROR' },
]

const filteredLogs = computed(() => {
  let result = logs.value
  if (filterLevel.value) {
    result = result.filter((e) => e.level.toUpperCase() === filterLevel.value.toUpperCase())
  }
  if (searchText.value) {
    const kw = searchText.value.toLowerCase()
    result = result.filter((e) => e.message.toLowerCase().includes(kw))
  }
  return result
})

function scrollToBottom() {
  if (!autoScroll.value) return
  nextTick(() => {
    const el = logContainer.value
    if (el) el.scrollTop = el.scrollHeight
  })
}

function onScroll() {
  const el = logContainer.value
  if (!el) return
  const atBottom = el.scrollHeight - el.scrollTop - el.clientHeight < 30
  if (atBottom && !autoScroll.value) {
    autoScroll.value = true
  }
}

function togglePause() {
  paused.value = !paused.value
  if (paused.value) {
    disconnectSSE()
    statusLabel.value = '已暂停'
    statusType.value = 'warning'
  } else {
    connectSSE()
  }
}

function clearLogs() {
  logs.value = []
}

async function fetchHistory() {
  try {
    const token = sessionStorage.getItem('token')
    const url = `${apiBaseURL}/system/logs?limit=500`
    const res = await fetch(url, {
      headers: { Authorization: `Bearer ${token}` },
    })
    const data = await res.json()
    if (data.code === 1000) {
      logs.value = data.data.logs || []
      scrollToBottom()
    }
  } catch {
    // ignore
  }
}

async function connectSSE() {
  disconnectSSE()
  abortController = new AbortController()
  const token = sessionStorage.getItem('token')
  const url = `${apiBaseURL}/system/logs/stream`

  statusLabel.value = '连接中...'
  statusType.value = 'warning'

  try {
    const response = await fetch(url, {
      headers: {
        Authorization: `Bearer ${token}`,
        Accept: 'text/event-stream',
      },
      signal: abortController.signal,
    })

    if (!response.ok || !response.body) {
      statusLabel.value = '连接失败'
      statusType.value = 'error'
      return
    }

    statusLabel.value = '实时'
    statusType.value = 'success'

    const reader = response.body.getReader()
    const decoder = new TextDecoder()
    let buffer = ''

    while (true) {
      const { done, value } = await reader.read()
      if (done) break

      buffer += decoder.decode(value, { stream: true })
      const lines = buffer.split('\n')
      buffer = lines.pop() || ''

      for (const line of lines) {
        if (line.startsWith('data: ')) {
          try {
            const entry: LogEntry = JSON.parse(line.slice(6))
            logs.value.push(entry)
            if (logs.value.length > 5000) {
              logs.value.splice(0, logs.value.length - 5000)
            }
          } catch {
            // ignore parse errors
          }
        }
      }
      scrollToBottom()
    }
  } catch (err: any) {
    if (err?.name !== 'AbortError') {
      statusLabel.value = '连接断开'
      statusType.value = 'error'
    }
  }
}

function disconnectSSE() {
  if (abortController) {
    abortController.abort()
    abortController = null
  }
}

watch(filteredLogs, () => {
  if (autoScroll.value) scrollToBottom()
})

onMounted(async () => {
  await fetchHistory()
  if (!paused.value) {
    connectSSE()
  }
})

onUnmounted(() => {
  disconnectSSE()
})
</script>

<style lang="scss" scoped>
.log-viewer {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 64px);
  background: #1a1a2e;
  font-family: 'Consolas', 'Menlo', 'Courier New', monospace;
  font-size: 13px;
  color: #e0e0e0;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #16213e;
  border-bottom: 1px solid #0f3460;
  gap: 12px;
  flex-shrink: 0;
  flex-wrap: wrap;

  .toolbar-left,
  .toolbar-right {
    display: flex;
    align-items: center;
    gap: 8px;
  }
}

.log-container {
  flex: 1;
  overflow-y: auto;
  padding: 4px 0;

  &::-webkit-scrollbar {
    width: 6px;
  }
  &::-webkit-scrollbar-track {
    background: transparent;
  }
  &::-webkit-scrollbar-thumb {
    background: #0f3460;
    border-radius: 3px;
  }
}

.log-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #666;
}

.log-line {
  display: flex;
  gap: 12px;
  padding: 2px 16px;
  line-height: 1.6;
  white-space: nowrap;

  &:hover {
    background: rgba(255, 255, 255, 0.03);
  }

  .log-time {
    color: #888;
    flex-shrink: 0;
    width: 160px;
  }

  .log-level {
    flex-shrink: 0;
    width: 52px;
    font-weight: bold;
    text-align: center;
  }

  .log-msg {
    overflow: hidden;
    text-overflow: ellipsis;
  }

  &.level-DEBUG .log-level { color: #888; }
  &.level-DEBUG .log-msg   { color: #999; }

  &.level-INFO .log-level { color: #4fc3f7; }
  &.level-INFO .log-msg   { color: #e0e0e0; }

  &.level-WARN {
    background: rgba(255, 152, 0, 0.06);
    .log-level { color: #ffb74d; }
    .log-msg   { color: #ffe0b2; }
  }

  &.level-ERROR {
    background: rgba(244, 67, 54, 0.08);
    .log-level { color: #ef5350; }
    .log-msg   { color: #ffcdd2; }
  }
}
</style>
