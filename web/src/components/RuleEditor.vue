<template>
  <div class="rule-editor">
    <div class="card">
      <div class="card-header">
        <h3>转换规则</h3>
        <p class="card-desc">配置文件后缀转换规则（仅改后缀，不转换文件格式）</p>
      </div>

      <div v-if="rules.length === 0" class="empty-state">
        <span class="empty-icon">🔄</span>
        <p>暂无转换规则</p>
        <p class="empty-hint">点击下方按钮添加第一条规则</p>
      </div>

      <div v-else class="rule-list">
        <div v-for="(rule, i) in rules" :key="i" class="rule-item">
          <label class="toggle-mini">
            <input type="checkbox" v-model="rule.enabled" @change="save" />
            <span class="toggle-slider-mini"></span>
          </label>

          <div class="rule-fields">
            <div class="field-group">
              <label class="field-label">原后缀</label>
              <input
                type="text"
                class="rule-input"
                v-model="rule.from"
                placeholder=".gif"
                @change="save"
              />
            </div>
            <span class="rule-arrow">→</span>
            <div class="field-group">
              <label class="field-label">目标后缀</label>
              <input
                type="text"
                class="rule-input"
                v-model="rule.to"
                placeholder=".jpg"
                @change="save"
              />
            </div>
          </div>

          <button class="icon-btn danger" @click="removeRule(i)" title="删除规则">
            <span>✕</span>
          </button>
        </div>
      </div>

      <button class="add-btn" @click="addRule">
        <span>＋</span> 添加规则
      </button>
    </div>

    <!-- 常用规则快捷添加 -->
    <div class="card">
      <div class="card-header">
        <h3>常用规则</h3>
        <p class="card-desc">一键添加常用的后缀转换规则</p>
      </div>
      <div class="preset-list">
        <button
          v-for="preset in presets"
          :key="preset.from + preset.to"
          class="preset-btn"
          @click="addPreset(preset)"
        >
          <span class="preset-from">{{ preset.from }}</span>
          <span class="preset-arrow">→</span>
          <span class="preset-to">{{ preset.to }}</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'RuleEditor',
  props: ['config'],
  emits: ['update'],
  data() {
    return {
      rules: (this.config?.rules || []).map(r => ({ ...r })),
      presets: [
        { from: '.gif', to: '.jpg' },
        { from: '.jpeg', to: '.jpg' },
        { from: '.png', to: '.jpg' },
        { from: '.tiff', to: '.tif' },
        { from: '.htm', to: '.html' },
        { from: '.mpeg', to: '.mp4' },
      ],
    }
  },
  watch: {
    config: {
      handler(val) {
        if (val) {
          this.rules = (val.rules || []).map(r => ({ ...r }))
        }
      },
      deep: true,
    }
  },
  methods: {
    addRule() {
      this.rules.push({ from: '', to: '', enabled: true })
    },
    removeRule(i) {
      this.rules.splice(i, 1)
      this.save()
    },
    addPreset(preset) {
      // 检查是否已存在
      const exists = this.rules.some(r => r.from === preset.from && r.to === preset.to)
      if (exists) return
      this.rules.push({ ...preset, enabled: true })
      this.save()
    },
    save() {
      this.$emit('update', {
        ...this.config,
        rules: this.rules.filter(r => r.from && r.to),
      })
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
  margin-bottom: 16px;
  box-shadow: var(--shadow-sm);
}

.card-header { margin-bottom: 16px; }
.card-header h3 { font-size: 16px; font-weight: 600; color: var(--gray-900); }
.card-desc { font-size: 13px; color: var(--gray-500); margin-top: 4px; }

.empty-state { text-align: center; padding: 24px; color: var(--gray-400); }
.empty-icon { font-size: 32px; display: block; margin-bottom: 8px; }
.empty-hint { font-size: 12px; margin-top: 4px; }

.rule-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 12px;
}

.rule-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: var(--gray-50);
  border-radius: var(--radius-sm);
  border: 1px solid var(--gray-200);
}

.rule-fields {
  flex: 1;
  display: flex;
  align-items: flex-end;
  gap: 10px;
}

.field-group {
  flex: 1;
}

.field-label {
  display: block;
  font-size: 11px;
  color: var(--gray-500);
  margin-bottom: 4px;
  font-weight: 500;
}

.rule-input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid var(--gray-300);
  border-radius: 6px;
  font-size: 14px;
  font-family: 'SF Mono', SFMono-Regular, Menlo, Consolas, monospace;
  outline: none;
  transition: border-color 0.2s;
}

.rule-input:focus {
  border-color: var(--primary);
  box-shadow: 0 0 0 3px var(--primary-light);
}

.rule-arrow {
  color: var(--gray-400);
  font-size: 18px;
  margin-bottom: 2px;
}

.icon-btn {
  width: 36px;
  height: 36px;
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  transition: all 0.2s;
  flex-shrink: 0;
}

.icon-btn.danger { background: #FEE2E2; color: var(--danger); }
.icon-btn.danger:hover { background: #FECACA; }

.add-btn {
  width: 100%;
  padding: 10px;
  border: 2px dashed var(--gray-300);
  border-radius: var(--radius-sm);
  background: transparent;
  color: var(--gray-500);
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
}

.add-btn:hover {
  border-color: var(--primary);
  color: var(--primary);
  background: var(--primary-light);
}

/* Mini Toggle */
.toggle-mini {
  position: relative;
  display: inline-block;
  width: 40px;
  height: 24px;
  flex-shrink: 0;
}

.toggle-mini input { opacity: 0; width: 0; height: 0; }

.toggle-slider-mini {
  position: absolute;
  cursor: pointer;
  top: 0; left: 0; right: 0; bottom: 0;
  background: var(--gray-300);
  border-radius: 12px;
  transition: all 0.3s;
}

.toggle-slider-mini::before {
  content: '';
  position: absolute;
  height: 18px; width: 18px;
  left: 3px; bottom: 3px;
  background: white;
  border-radius: 50%;
  transition: all 0.3s;
  box-shadow: var(--shadow-sm);
}

.toggle-mini input:checked + .toggle-slider-mini { background: var(--primary); }
.toggle-mini input:checked + .toggle-slider-mini::before { transform: translateX(16px); }

/* 预设规则 */
.preset-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.preset-btn {
  padding: 8px 14px;
  border: 1px solid var(--gray-200);
  border-radius: 20px;
  background: white;
  cursor: pointer;
  font-size: 13px;
  display: flex;
  align-items: center;
  gap: 6px;
  transition: all 0.2s;
}

.preset-btn:hover {
  border-color: var(--primary);
  background: var(--primary-light);
}

.preset-from {
  font-family: 'SF Mono', SFMono-Regular, Menlo, Consolas, monospace;
  color: var(--danger);
  font-weight: 500;
}

.preset-arrow { color: var(--gray-400); }

.preset-to {
  font-family: 'SF Mono', SFMono-Regular, Menlo, Consolas, monospace;
  color: var(--success);
  font-weight: 500;
}

@media (max-width: 640px) {
  .rule-fields { flex-direction: column; gap: 6px; }
  .rule-arrow { display: none; }
}
</style>
