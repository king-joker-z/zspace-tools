<template>
  <div class="app">
    <!-- 顶部导航栏 -->
    <header class="navbar">
      <div class="navbar-brand">
        <svg class="logo" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/>
        </svg>
        <span class="brand-text">ZSpace Tools</span>
        <span class="version">v0.1.0</span>
      </div>
      <div class="navbar-status">
        <span :class="['status-dot', { active: status.watching }]"></span>
        <span class="status-text">{{ status.watching ? '监听中' : '已停止' }}</span>
      </div>
    </header>

    <!-- 主内容区 -->
    <main class="main-content">
      <!-- Tab 导航 -->
      <nav class="tabs">
        <button
          v-for="tab in tabs"
          :key="tab.key"
          :class="['tab-btn', { active: currentTab === tab.key }]"
          @click="currentTab = tab.key"
        >
          <span class="tab-icon" v-html="tab.icon"></span>
          {{ tab.label }}
        </button>
      </nav>

      <!-- Dashboard -->
      <div v-if="currentTab === 'dashboard'" class="tab-content">
        <Dashboard
          :status="status"
          :config="config"
          :logs="recentLogs"
          @start-watch="startWatch"
          @stop-watch="stopWatch"
          @scan="manualScan"
        />
      </div>

      <!-- 文件夹配置 -->
      <div v-if="currentTab === 'folders'" class="tab-content">
        <FolderConfig
          :config="config"
          @update="updateConfig"
        />
      </div>

      <!-- 规则编辑 -->
      <div v-if="currentTab === 'rules'" class="tab-content">
        <RuleEditor
          :config="config"
          @update="updateConfig"
        />
      </div>

      <!-- 转换日志 -->
      <div v-if="currentTab === 'logs'" class="tab-content">
        <LogViewer :logs="allLogs" @load-more="loadMoreLogs" />
      </div>
    </main>
  </div>
</template>

<script>
import Dashboard from './views/Dashboard.vue'
import FolderConfig from './components/FolderConfig.vue'
import RuleEditor from './components/RuleEditor.vue'
import LogViewer from './components/LogViewer.vue'

export default {
  name: 'App',
  components: { Dashboard, FolderConfig, RuleEditor, LogViewer },
  data() {
    return {
      currentTab: 'dashboard',
      tabs: [
        { key: 'dashboard', label: '仪表盘', icon: '📊' },
        { key: 'folders', label: '监听文件夹', icon: '📁' },
        { key: 'rules', label: '转换规则', icon: '🔄' },
        { key: 'logs', label: '转换日志', icon: '📋' },
      ],
      config: {
        watch_dirs: ['/watch'],
        rules: [],
        auto_watch: false,
      },
      status: {
        watching: false,
        watch_dirs: [],
        rules: 0,
        log_count: 0,
      },
      recentLogs: [],
      allLogs: [],
      ws: null,
      logPage: 1,
    }
  },
  mounted() {
    this.fetchConfig()
    this.fetchStatus()
    this.fetchLogs()
    this.connectWebSocket()
    // 定时刷新状态
    this._statusTimer = setInterval(() => this.fetchStatus(), 5000)
  },
  beforeUnmount() {
    clearInterval(this._statusTimer)
    if (this.ws) this.ws.close()
  },
  methods: {
    async fetchConfig() {
      try {
        const res = await fetch('/api/config')
        const data = await res.json()
        if (data.code === 0) this.config = data.data
      } catch (e) {
        console.error('获取配置失败', e)
      }
    },
    async fetchStatus() {
      try {
        const res = await fetch('/api/status')
        const data = await res.json()
        if (data.code === 0) this.status = data.data
      } catch (e) {
        console.error('获取状态失败', e)
      }
    },
    async fetchLogs() {
      try {
        const res = await fetch(`/api/logs?page=1&size=50`)
        const data = await res.json()
        if (data.code === 0) {
          this.allLogs = data.data.items || []
          this.recentLogs = this.allLogs.slice(0, 10)
        }
      } catch (e) {
        console.error('获取日志失败', e)
      }
    },
    async loadMoreLogs() {
      this.logPage++
      try {
        const res = await fetch(`/api/logs?page=${this.logPage}&size=50`)
        const data = await res.json()
        if (data.code === 0 && data.data.items) {
          this.allLogs = [...this.allLogs, ...data.data.items]
        }
      } catch (e) {
        console.error('加载更多日志失败', e)
      }
    },
    async updateConfig(newConfig) {
      try {
        const res = await fetch('/api/config', {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(newConfig),
        })
        const data = await res.json()
        if (data.code === 0) {
          this.config = newConfig
          this.fetchStatus()
        }
      } catch (e) {
        console.error('更新配置失败', e)
      }
    },
    async startWatch() {
      try {
        await fetch('/api/watch/start', { method: 'POST' })
        this.fetchStatus()
      } catch (e) {
        console.error('启动监听失败', e)
      }
    },
    async stopWatch() {
      try {
        await fetch('/api/watch/stop', { method: 'POST' })
        this.fetchStatus()
      } catch (e) {
        console.error('停止监听失败', e)
      }
    },
    async manualScan() {
      try {
        const res = await fetch('/api/scan', { method: 'POST' })
        const data = await res.json()
        if (data.code === 0) {
          this.fetchLogs()
          this.fetchStatus()
        }
      } catch (e) {
        console.error('手动扫描失败', e)
      }
    },
    connectWebSocket() {
      const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:'
      const wsUrl = `${protocol}//${location.host}/api/ws/logs`
      this.ws = new WebSocket(wsUrl)

      this.ws.onmessage = (event) => {
        try {
          const entry = JSON.parse(event.data)
          this.recentLogs.unshift(entry)
          this.allLogs.unshift(entry)
          // 限制前端缓存
          if (this.recentLogs.length > 20) this.recentLogs.pop()
          if (this.allLogs.length > 500) this.allLogs.pop()
        } catch (e) {
          // 忽略解析错误
        }
      }

      this.ws.onclose = () => {
        // 自动重连
        setTimeout(() => this.connectWebSocket(), 3000)
      }
    },
  },
}
</script>

<style>
/* ==================== 全局样式 ==================== */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

:root {
  --primary: #007AFF;
  --primary-light: #E3F2FF;
  --primary-dark: #0056CC;
  --success: #34C759;
  --warning: #FF9500;
  --danger: #FF3B30;
  --gray-50: #F9FAFB;
  --gray-100: #F3F4F6;
  --gray-200: #E5E7EB;
  --gray-300: #D1D5DB;
  --gray-400: #9CA3AF;
  --gray-500: #6B7280;
  --gray-600: #4B5563;
  --gray-700: #374151;
  --gray-800: #1F2937;
  --gray-900: #111827;
  --radius: 12px;
  --radius-sm: 8px;
  --shadow-sm: 0 1px 2px rgba(0,0,0,0.05);
  --shadow: 0 1px 3px rgba(0,0,0,0.1), 0 1px 2px rgba(0,0,0,0.06);
  --shadow-md: 0 4px 6px rgba(0,0,0,0.07), 0 2px 4px rgba(0,0,0,0.06);
  --shadow-lg: 0 10px 15px rgba(0,0,0,0.1), 0 4px 6px rgba(0,0,0,0.05);
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'SF Pro Text', 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  background: var(--gray-50);
  color: var(--gray-800);
  line-height: 1.6;
  -webkit-font-smoothing: antialiased;
}

/* ==================== 导航栏 ==================== */
.navbar {
  background: white;
  padding: 0 24px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid var(--gray-200);
  position: sticky;
  top: 0;
  z-index: 100;
  backdrop-filter: blur(20px);
  background: rgba(255, 255, 255, 0.85);
}

.navbar-brand {
  display: flex;
  align-items: center;
  gap: 10px;
}

.logo {
  width: 28px;
  height: 28px;
  color: var(--primary);
}

.brand-text {
  font-size: 18px;
  font-weight: 600;
  color: var(--gray-900);
}

.version {
  font-size: 11px;
  color: var(--gray-400);
  background: var(--gray-100);
  padding: 2px 8px;
  border-radius: 10px;
}

.navbar-status {
  display: flex;
  align-items: center;
  gap: 8px;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--gray-400);
  transition: all 0.3s;
}

.status-dot.active {
  background: var(--success);
  box-shadow: 0 0 8px rgba(52, 199, 89, 0.5);
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.status-text {
  font-size: 14px;
  color: var(--gray-500);
}

/* ==================== Tab 导航 ==================== */
.main-content {
  max-width: 960px;
  margin: 0 auto;
  padding: 20px 16px;
}

.tabs {
  display: flex;
  gap: 4px;
  background: var(--gray-100);
  padding: 4px;
  border-radius: var(--radius);
  margin-bottom: 24px;
  overflow-x: auto;
}

.tab-btn {
  flex: 1;
  min-width: 0;
  padding: 10px 16px;
  border: none;
  background: transparent;
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  color: var(--gray-500);
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  white-space: nowrap;
}

.tab-btn:hover {
  color: var(--gray-700);
}

.tab-btn.active {
  background: white;
  color: var(--primary);
  box-shadow: var(--shadow-sm);
}

.tab-icon {
  font-size: 16px;
}

.tab-content {
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(4px); }
  to { opacity: 1; transform: translateY(0); }
}

/* ==================== 响应式 ==================== */
@media (max-width: 640px) {
  .navbar {
    padding: 0 16px;
  }

  .main-content {
    padding: 16px 12px;
  }

  .tab-btn {
    padding: 8px 12px;
    font-size: 13px;
  }

  .tab-icon {
    display: none;
  }
}
</style>
