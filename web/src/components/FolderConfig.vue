<template>
  <div class="folder-config">
    <div class="card">
      <div class="card-header">
        <h3>监听文件夹</h3>
        <p class="card-desc">配置需要监听的文件夹路径（容器内路径，需通过 volume 映射）</p>
      </div>

      <div class="folder-list">
        <div v-for="(dir, i) in dirs" :key="i" class="folder-item">
          <span class="folder-icon">📁</span>
          <input
            type="text"
            class="folder-input"
            v-model="dirs[i]"
            placeholder="/watch"
            @change="save"
          />
          <button class="icon-btn danger" @click="removeDir(i)" title="删除">
            <span>✕</span>
          </button>
        </div>
      </div>

      <button class="add-btn" @click="addDir">
        <span>＋</span> 添加文件夹
      </button>
    </div>

    <div class="card">
      <div class="card-header">
        <h3>自动监听</h3>
        <p class="card-desc">启用后，服务启动时自动开始文件监听</p>
      </div>
      <label class="toggle-row">
        <span>启动时自动监听</span>
        <label class="toggle">
          <input type="checkbox" v-model="autoWatch" @change="save" />
          <span class="toggle-slider"></span>
        </label>
      </label>
    </div>
  </div>
</template>

<script>
export default {
  name: 'FolderConfig',
  props: ['config'],
  emits: ['update'],
  data() {
    return {
      dirs: [...(this.config?.watch_dirs || ['/watch'])],
      autoWatch: this.config?.auto_watch || false,
    }
  },
  watch: {
    config: {
      handler(val) {
        if (val) {
          this.dirs = [...(val.watch_dirs || ['/watch'])]
          this.autoWatch = val.auto_watch || false
        }
      },
      deep: true,
    }
  },
  methods: {
    addDir() {
      this.dirs.push('/watch')
    },
    removeDir(i) {
      if (this.dirs.length <= 1) return
      this.dirs.splice(i, 1)
      this.save()
    },
    save() {
      this.$emit('update', {
        ...this.config,
        watch_dirs: this.dirs.filter(d => d.trim()),
        auto_watch: this.autoWatch,
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

.folder-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 12px;
}

.folder-item {
  display: flex;
  align-items: center;
  gap: 10px;
}

.folder-icon { font-size: 20px; }

.folder-input {
  flex: 1;
  padding: 10px 14px;
  border: 1px solid var(--gray-300);
  border-radius: var(--radius-sm);
  font-size: 14px;
  font-family: 'SF Mono', SFMono-Regular, Menlo, Consolas, monospace;
  outline: none;
  transition: border-color 0.2s;
}

.folder-input:focus {
  border-color: var(--primary);
  box-shadow: 0 0 0 3px var(--primary-light);
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

.icon-btn.danger {
  background: #FEE2E2;
  color: var(--danger);
}
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

/* Toggle 开关 */
.toggle-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 14px;
  color: var(--gray-700);
}

.toggle {
  position: relative;
  display: inline-block;
  width: 48px;
  height: 28px;
}

.toggle input {
  opacity: 0;
  width: 0;
  height: 0;
}

.toggle-slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--gray-300);
  border-radius: 14px;
  transition: all 0.3s;
}

.toggle-slider::before {
  content: '';
  position: absolute;
  height: 22px;
  width: 22px;
  left: 3px;
  bottom: 3px;
  background: white;
  border-radius: 50%;
  transition: all 0.3s;
  box-shadow: var(--shadow-sm);
}

.toggle input:checked + .toggle-slider {
  background: var(--primary);
}

.toggle input:checked + .toggle-slider::before {
  transform: translateX(20px);
}
</style>
