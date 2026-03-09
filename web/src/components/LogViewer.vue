<template>
  <div class="log-viewer">
    <div class="card">
      <div class="card-header">
        <h3>转换日志</h3>
        <p class="card-desc">实时显示文件后缀转换记录</p>
      </div>

      <div v-if="logs.length === 0" class="empty-state">
        <span class="empty-icon">📋</span>
        <p>暂无转换日志</p>
        <p class="empty-hint">配置规则并启动监听或手动扫描后，转换记录将在此显示</p>
      </div>

      <div v-else>
        <div class="log-table">
          <div class="log-header">
            <span class="col-status">状态</span>
            <span class="col-message">详情</span>
            <span class="col-rule">规则</span>
            <span class="col-time">时间</span>
          </div>
          <div v-for="(log, i) in logs" :key="i" class="log-row">
            <span class="col-status">
              <span :class="['status-tag', 'tag-' + log.status]">
                {{ statusText(log.status) }}
              </span>
            </span>
            <span class="col-message" :title="log.message">{{ log.message }}</span>
            <span class="col-rule">
              <code>{{ log.from }} → {{ log.to }}</code>
            </span>
            <span class="col-time">{{ log.time }}</span>
          </div>
        </div>

        <button class="load-more-btn" @click="$emit('load-more')">
          加载更多
        </button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'LogViewer',
  props: ['logs'],
  emits: ['load-more'],
  methods: {
    statusText(status) {
      const map = { success: '成功', error: '失败', skipped: '跳过' }
      return map[status] || status
    }
  }
}
</script>

<style scoped>
.card {
  background: white;
  border-radius: var(--radius);
  border: 1px solid var(--gray-200);
  padding: 20px;
  box-shadow: var(--shadow-sm);
}

.card-header { margin-bottom: 16px; }
.card-header h3 { font-size: 16px; font-weight: 600; color: var(--gray-900); }
.card-desc { font-size: 13px; color: var(--gray-500); margin-top: 4px; }

.empty-state { text-align: center; padding: 40px; color: var(--gray-400); }
.empty-icon { font-size: 40px; display: block; margin-bottom: 12px; }
.empty-hint { font-size: 12px; margin-top: 4px; }

.log-table {
  border: 1px solid var(--gray-200);
  border-radius: var(--radius-sm);
  overflow: hidden;
}

.log-header {
  display: flex;
  padding: 10px 14px;
  background: var(--gray-50);
  border-bottom: 1px solid var(--gray-200);
  font-size: 12px;
  font-weight: 600;
  color: var(--gray-500);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.log-row {
  display: flex;
  padding: 10px 14px;
  border-bottom: 1px solid var(--gray-100);
  font-size: 13px;
  transition: background 0.15s;
}

.log-row:last-child { border-bottom: none; }
.log-row:hover { background: var(--gray-50); }

.col-status { width: 60px; flex-shrink: 0; }
.col-message { flex: 1; min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; color: var(--gray-700); }
.col-rule { width: 140px; flex-shrink: 0; }
.col-time { width: 140px; flex-shrink: 0; color: var(--gray-400); text-align: right; }

.col-rule code {
  font-size: 12px;
  background: var(--gray-100);
  padding: 2px 6px;
  border-radius: 4px;
  color: var(--gray-600);
}

.status-tag {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 10px;
  font-weight: 500;
}

.tag-success { background: #DCFCE7; color: #166534; }
.tag-error { background: #FEE2E2; color: #991B1B; }
.tag-skipped { background: #FEF3C7; color: #92400E; }

.load-more-btn {
  width: 100%;
  padding: 10px;
  border: none;
  background: var(--gray-50);
  border-radius: var(--radius-sm);
  color: var(--primary);
  font-size: 14px;
  cursor: pointer;
  margin-top: 12px;
  transition: all 0.2s;
}

.load-more-btn:hover { background: var(--primary-light); }

@media (max-width: 640px) {
  .col-rule, .col-time { display: none; }
  .log-header .col-rule, .log-header .col-time { display: none; }
}
</style>
