<template>
  <div class="dashboard">
    <!-- 状态卡片行 -->
    <div class="stat-cards">
      <div class="stat-card">
        <div class="stat-icon" :style="{ background: status.watching ? '#DCFCE7' : '#FEE2E2' }">
          <span>{{ status.watching ? '👁️' : '⏸️' }}</span>
        </div>
        <div class="stat-info">
          <div class="stat-label">监听状态</div>
          <div :class="['stat-value', status.watching ? 'text-success' : 'text-danger']">
            {{ status.watching ? '运行中' : '已停止' }}
          </div>
        </div>
        <button
          :class="['control-btn', status.watching ? 'btn-danger' : 'btn-primary']"
          @click="status.watching ? $emit('stop-watch') : $emit('start-watch')"
        >
          {{ status.watching ? '停止' : '启动' }}
        </button>
      </div>

      <div class="stat-card">
        <div class="stat-icon" style="background: #E3F2FF">
          <span>📁</span>
        </div>
        <div class="stat-info">
          <div class="stat-label">监听文件夹</div>
          <div class="stat-value">{{ config.watch_dirs?.length || 0 }} 个</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon" style="background: #FFF3E0">
          <span>🔄</span>
        </div>
        <div class="stat-info">
          <div class="stat-label">转换规则</div>
          <div class="stat-value">{{ config.rules?.length || 0 }} 条</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon" style="background: #F3E5F5">
          <span>📋</span>
        </div>
        <div class="stat-info">
          <div class="stat-label">转换记录</div>
          <div class="stat-value">{{ status.log_count || 0 }} 条</div>
        </div>
      </div>
    </div>

    <!-- 手动扫描 -->
    <div class="card scan-card">
      <div class="card-header">
        <h3>手动扫描</h3>
        <p class="card-desc">一键扫描所有监听文件夹，按规则批量转换文件后缀</p>
      </div>
      <button class="scan-btn" @click="handleScan" :disabled="scanning">
        <span v-if="scanning" class="spinner"></span>
        <span v-else>🔍</span>
        {{ scanning ? '扫描中...' : '立即扫描' }}
      </button>
    </div>

    <!-- 规则预览 -->
    <div class="card" v-if="config.rules && config.rules.length > 0">
      <div class="card-header">
        <h3>当前规则</h3>
      </div>
      <div class="rule-list">
        <div v-for="(rule, i) in config.rules" :key="i" class="rule-item">
          <span :class="['rule-badge', rule.enabled ? 'badge-active' : 'badge-inactive']">
            {{ rule.enabled ? '启用' : '禁用' }}
          </span>
          <span class="rule-from">{{ rule.from }}</span>
          <span class="rule-arrow">→</span>
          <span class="rule-to">{{ rule.to }}</span>
        </div>
      </div>
    </div>

    <!-- 最近日志 -->
    <div class="card">
      <div class="card-header">
        <h3>最近转换</h3>
      </div>
      <div v-if="logs.length === 0" class="empty-state">
        <span class="empty-icon">📭</span>
        <p>暂无转换记录</p>
      </div>
      <div v-else class="log-list">
        <div v-for="(log, i) in logs" :key="i" class="log-item">
          <span :class="['log-status', 'status-' + log.status]">
            {{ log.status === 'success' ? '✅' : log.status === 'error' ? '❌' : '⏭️' }}
          </span>
          <div class="log-info">
            <div class="log-message">{{ log.message }}</div>
            <div class="log-time">{{ log.time }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Dashboard',
  props: ['status', 'config', 'logs'],
  emits: ['start-watch', 'stop-watch', 'scan'],
  data() {
    return { scanning: false }
  },
  methods: {
    async handleScan() {
      this.scanning = true
      this.$emit('scan')
      setTimeout(() => { this.scanning = false }, 2000)
    }
  }
}
</script>

<style scoped>
.stat-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 12px;
  margin-bottom: 20px;
}

.stat-card {
  background: white;
  border-radius: var(--radius);
  padding: 16px;
  display: flex;
  align-items: center;
  gap: 12px;
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--gray-200);
}

.stat-icon {
  width: 44px;
  height: 44px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  flex-shrink: 0;
}

.stat-info {
  flex: 1;
  min-width: 0;
}

.stat-label {
  font-size: 12px;
  color: var(--gray-500);
}

.stat-value {
  font-size: 18px;
  font-weight: 600;
  color: var(--gray-900);
}

.text-success { color: var(--success); }
.text-danger { color: var(--danger); }

.control-btn {
  padding: 6px 14px;
  border: none;
  border-radius: var(--radius-sm);
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  flex-shrink: 0;
}

.btn-primary {
  background: var(--primary);
  color: white;
}
.btn-primary:hover { background: var(--primary-dark); }

.btn-danger {
  background: var(--danger);
  color: white;
}
.btn-danger:hover { background: #E5342A; }

/* 卡片 */
.card {
  background: white;
  border-radius: var(--radius);
  border: 1px solid var(--gray-200);
  padding: 20px;
  margin-bottom: 16px;
  box-shadow: var(--shadow-sm);
}

.card-header {
  margin-bottom: 16px;
}

.card-header h3 {
  font-size: 16px;
  font-weight: 600;
  color: var(--gray-900);
}

.card-desc {
  font-size: 13px;
  color: var(--gray-500);
  margin-top: 4px;
}

/* 扫描按钮 */
.scan-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.scan-btn {
  padding: 10px 24px;
  border: none;
  border-radius: var(--radius-sm);
  background: var(--primary);
  color: white;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.scan-btn:hover { background: var(--primary-dark); }
.scan-btn:disabled { opacity: 0.6; cursor: not-allowed; }

.spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255,255,255,0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* 规则列表 */
.rule-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.rule-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  background: var(--gray-50);
  border-radius: var(--radius-sm);
}

.rule-badge {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 10px;
  font-weight: 500;
}

.badge-active {
  background: #DCFCE7;
  color: #166534;
}

.badge-inactive {
  background: var(--gray-200);
  color: var(--gray-500);
}

.rule-from, .rule-to {
  font-family: 'SF Mono', SFMono-Regular, Menlo, Consolas, monospace;
  font-size: 14px;
  font-weight: 500;
}

.rule-from { color: var(--danger); }
.rule-arrow { color: var(--gray-400); }
.rule-to { color: var(--success); }

/* 日志列表 */
.log-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.log-item {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 10px 12px;
  background: var(--gray-50);
  border-radius: var(--radius-sm);
}

.log-status {
  font-size: 16px;
  flex-shrink: 0;
  margin-top: 2px;
}

.log-info {
  flex: 1;
  min-width: 0;
}

.log-message {
  font-size: 13px;
  color: var(--gray-700);
  word-break: break-all;
}

.log-time {
  font-size: 11px;
  color: var(--gray-400);
  margin-top: 2px;
}

.empty-state {
  text-align: center;
  padding: 32px;
  color: var(--gray-400);
}

.empty-icon {
  font-size: 32px;
  display: block;
  margin-bottom: 8px;
}

@media (max-width: 640px) {
  .stat-cards {
    grid-template-columns: repeat(2, 1fr);
  }
  .scan-card {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }
  .scan-btn {
    justify-content: center;
  }
}
</style>
