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
        <pre><code>{{ share.code }}</code></pre>
      </div>
      <div class="actions">
        <button @click="runCode" class="btn btn-primary" :disabled="isRunning">
          {{ isRunning ? '运行中...' : '运行' }}
        </button>
        <button @click="copyCode" class="btn btn-secondary">
          {{ copied ? '已复制' : '复制代码' }}
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
      copied: false
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
      this.isRunning = true
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
        })
        const data = await response.json()
        this.output = data.error || data.output
      } catch (error) {
        this.output = '运行失败: ' + error.message
      } finally {
        this.isRunning = false
      }
    },
    async copyCode() {
      try {
        await navigator.clipboard.writeText(this.share.code)
        this.copied = true
        setTimeout(() => {
          this.copied = false
        }, 2000)
      } catch (error) {
        console.error('复制失败:', error)
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
}

.share-meta {
  color: #666;
  font-size: 14px;
}

.description {
  margin: 10px 0;
  white-space: pre-line;
}

.info {
  display: flex;
  gap: 20px;
  margin: 10px 0;
}

.expires {
  color: #dc3545;
}

.code-container {
  background: #f8f9fa;
  border-radius: 4px;
  padding: 20px;
  margin: 20px 0;
  overflow-x: auto;
}

.code-container pre {
  margin: 0;
  font-family: 'Fira Code', monospace;
  font-size: 14px;
}

.actions {
  display: flex;
  gap: 10px;
  margin: 20px 0;
}

.btn {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.btn-primary {
  background-color: #007bff;
  color: white;
}

.btn-primary:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

.btn-secondary {
  background-color: #6c757d;
  color: white;
}

.output {
  background: #f8f9fa;
  border-radius: 4px;
  padding: 20px;
  margin-top: 20px;
}

.output h3 {
  margin: 0 0 10px 0;
  font-size: 18px;
}

.output pre {
  margin: 0;
  font-family: monospace;
  white-space: pre-wrap;
}
</style> 