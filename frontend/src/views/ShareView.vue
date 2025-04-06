<template>
  <div class="share-view">
    <div v-if="loading" class="loading">
      加载中...
    </div>
    <div v-else-if="error" class="error">
      {{ error }}
    </div>
    <div v-else class="share-content">
      <div class="share-header">
        <h1>{{ share.title || '未命名代码片段' }}</h1>
        <div class="share-meta">
          <div v-if="share.description" class="description">
            {{ share.description }}
          </div>
          <div class="info">
            <span>Go {{ share.version }}</span>
            <span>创建于 {{ formatDate(share.created_at) }}</span>
            <span v-if="share.author">作者: {{ share.author }}</span>
            <span>浏览次数: {{ share.views }}</span>
          </div>
          <div v-if="share.expires_at" class="expires">
            将在 {{ formatExpires(share.expires_at) }} 后过期
          </div>
        </div>
      </div>
      <div class="code-container">
        <pre class="code-block"><code>{{ share.code }}</code></pre>
      </div>
      <div class="actions">
        <button @click="runCode" class="btn btn-primary" :disabled="isRunning || isCopying">
          <i class="fas fa-play"></i> {{ isRunning ? '运行中...' : '运行' }}
        </button>
        <button @click="copyCode" class="btn btn-secondary" :disabled="isRunning || isCopying">
          <i class="fas fa-copy"></i> {{ copied ? '已复制' : '复制代码' }}
        </button>
      </div>
      <div v-if="output" class="output">
        <h3>运行结果</h3>
        <pre>{{ output }}</pre>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ShareView',
  data() {
    return {
      loading: true,
      error: null,
      share: null,
      output: '',
      isRunning: false,
      copied: false,
      isCopying: false
    }
  },
  async created() {
    await this.loadShare()
  },
  methods: {
    async loadShare() {
      const shareId = this.$route.params.id
      try {
        const response = await fetch(`/api/share/${shareId}`)
        if (!response.ok) {
          if (response.status === 404) {
            throw new Error('分享不存在或已过期')
          }
          throw new Error('加载失败')
        }
        this.share = await response.json()
        // 增加访问次数
        fetch(`/api/share/${shareId}/view`, { method: 'POST' }).catch(console.error)
      } catch (error) {
        this.error = error.message
      } finally {
        this.loading = false
      }
    },
    async runCode() {
      if (this.isRunning || this.isCopying) return;
      this.isRunning = true;
      try {
        const response = await fetch('/api/run', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            code: this.share.code,
            version: this.share.version
          })
        });
        const data = await response.json();
        this.output = data.error || data.output;
      } catch (error) {
        this.output = '运行失败: ' + error.message;
      } finally {
        this.isRunning = false;
      }
    },
    async copyCode() {
      if (this.isRunning || this.isCopying) return;
      this.isCopying = true;
      try {
        await navigator.clipboard.writeText(this.share.code);
        this.copied = true;
        setTimeout(() => {
          this.copied = false;
        }, 2000);
      } catch (error) {
        console.error('复制失败:', error);
      } finally {
        this.isCopying = false;
      }
    },
    formatDate(date) {
      return new Date(date).toLocaleString()
    },
    formatExpires(date) {
      const expires = new Date(date)
      const now = new Date()
      const diff = expires - now
      const days = Math.floor(diff / (1000 * 60 * 60 * 24))
      const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
      if (days > 0) {
        return `${days}天${hours}小时`
      }
      return `${hours}小时`
    }
  }
}
</script>

<style scoped>
.share-view {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  box-sizing: border-box;
  width: 100%;
}

.loading, .error {
  text-align: center;
  padding: 40px;
  font-size: 18px;
}

.error {
  color: #dc3545;
}

.share-header {
  margin-bottom: 20px;
}

.share-header h1 {
  margin: 0 0 10px 0;
  font-size: 24px;
  word-break: break-word;
}

.share-meta {
  color: #666;
  font-size: 14px;
}

.description {
  margin: 10px 0;
  white-space: pre-line;
  word-break: break-word;
}

.info {
  display: flex;
  flex-wrap: wrap;
  gap: 10px 20px;
  margin: 10px 0;
}

.expires {
  color: #dc3545;
  margin-top: 5px;
}

.code-container {
  background: #f8f9fa;
  border-radius: 4px;
  padding: 15px;
  margin: 20px 0;
  overflow-x: auto;
  position: relative;
  max-width: 100%;
  border: 1px solid #e9ecef;
}

.code-block {
  margin: 0;
  font-family: 'Fira Code', 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.5;
  white-space: pre-wrap;
  word-break: break-all;
  word-wrap: break-word;
  tab-size: 4;
  -moz-tab-size: 4;
  -o-tab-size: 4;
  color: #333;
}

.code-block code {
  display: block;
  width: 100%;
}

/* 基础语法高亮 */
.code-block .keyword {
  color: #0000ff;
}

.code-block .string {
  color: #a31515;
}

.code-block .comment {
  color: #008000;
}

.code-block .function {
  color: #795e26;
}

.actions {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin: 20px 0;
}

.btn {
  padding: 10px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  white-space: nowrap;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 5px;
}

.btn i {
  font-size: 14px;
}

.btn-primary {
  background-color: #007bff;
  color: white;
}

.btn-secondary {
  background-color: #6c757d;
  color: white;
}

.output {
  background: #f8f9fa;
  border-radius: 4px;
  padding: 15px;
  margin-top: 20px;
  width: 100%;
  overflow-x: auto;
}

.output h3 {
  margin: 0 0 10px 0;
  font-size: 18px;
}

.output pre {
  margin: 0;
  font-family: monospace;
  white-space: pre-wrap;
  word-break: break-word;
  font-size: 14px;
  line-height: 1.4;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .share-view {
    padding: 15px 10px;
  }

  .share-header h1 {
    font-size: 20px;
  }

  .info {
    flex-direction: column;
    gap: 5px;
  }

  .code-container {
    padding: 10px;
    margin: 15px 0;
    font-size: 0;  /* 修复移动端可能存在的间隙问题 */
  }

  .code-block {
    font-size: 16px;
    line-height: 1.4;
    overflow-wrap: anywhere; /* 改进移动端文本换行 */
  }

  .actions {
    justify-content: space-between;
  }

  .btn {
    flex: 1;
    text-align: center;
    padding: 12px 10px;
    font-size: 14px;
  }

  .btn i {
    margin-right: 3px;
  }

  .output pre {
    font-size: 16px;
  }
}

@media (max-width: 480px) {
  .share-view {
    padding: 10px 5px;
  }

  .share-header h1 {
    font-size: 18px;
  }

  .share-meta {
    font-size: 13px;
  }

  .code-container {
    padding: 8px;
    border-radius: 3px;
  }

  .actions {
    flex-direction: column;
    width: 100%;
  }

  .btn {
    width: 100%;
    margin-bottom: 5px;
    padding: 12px 0;
    font-size: 14px;
  }

  .output {
    padding: 10px;
  }

  .output h3 {
    font-size: 16px;
  }
}
</style> 