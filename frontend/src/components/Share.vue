<!-- Share.vue -->
<template>
  <div class="share-dialog">
    <div v-if="!shareUrl" class="share-form">
      <h3>分享代码</h3>
      <div class="form-group">
        <label for="title">标题</label>
        <input
          id="title"
          v-model="shareData.title"
          type="text"
          class="form-control"
          placeholder="给你的代码起个标题"
        />
      </div>
      <div class="form-group">
        <label for="description">描述</label>
        <textarea
          id="description"
          v-model="shareData.description"
          class="form-control"
          rows="3"
          placeholder="添加一些描述"
        ></textarea>
      </div>
      <div class="form-group">
        <label for="expires">过期时间</label>
        <select id="expires" v-model="shareData.expires_in" class="form-control">
          <option value="">永不过期</option>
          <option value="1h">1小时</option>
          <option value="24h">1天</option>
          <option value="168h">7天</option>
          <option value="720h">30天</option>
        </select>
      </div>
      <div class="actions">
        <button @click="$emit('close')" class="btn btn-secondary">取消</button>
        <button @click="shareCode" class="btn btn-primary" :disabled="isSharing">
          {{ isSharing ? '分享中...' : '分享' }}
        </button>
      </div>
    </div>
    <div v-else class="share-result">
      <h3>分享成功！</h3>
      <div class="url-display">
        <input
          ref="urlInput"
          type="text"
          readonly
          :value="shareUrl"
          class="form-control"
        />
        <button @click="copyUrl" class="btn btn-primary">
          {{ copied ? '已复制' : '复制链接' }}
        </button>
      </div>
      <div class="expires-info" v-if="expiresAt">
        链接将在 {{ formatExpiresAt }} 后过期
      </div>
      <div class="actions">
        <button @click="$emit('close')" class="btn btn-secondary">关闭</button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Share',
  props: {
    code: {
      type: String,
      required: true
    },
    version: {
      type: String,
      required: true
    }
  },
  data() {
    return {
      shareData: {
        title: '',
        description: '',
        expires_in: '168h' // 默认7天
      },
      isSharing: false,
      shareUrl: '',
      expiresAt: null,
      copied: false
    }
  },
  computed: {
    formatExpiresAt() {
      if (!this.expiresAt) return ''
      const expires = new Date(this.expiresAt)
      const now = new Date()
      const diff = expires - now
      const days = Math.floor(diff / (1000 * 60 * 60 * 24))
      const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
      if (days > 0) {
        return `${days}天${hours}小时`
      }
      return `${hours}小时`
    }
  },
  methods: {
    async shareCode() {
      this.isSharing = true
      try {
        const response = await fetch('/api/share', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            code: this.code,
            version: this.version,
            ...this.shareData
          })
        })

        if (!response.ok) {
          throw new Error('分享失败')
        }

        const data = await response.json()
        this.shareUrl = window.location.origin + data.url
        this.expiresAt = data.expires_at
      } catch (error) {
        console.error('分享失败:', error)
        alert('分享失败，请稍后重试')
      } finally {
        this.isSharing = false
      }
    },
    async copyUrl() {
      try {
        await navigator.clipboard.writeText(this.shareUrl)
        this.copied = true
        setTimeout(() => {
          this.copied = false
        }, 2000)
      } catch (error) {
        console.error('复制失败:', error)
        // 回退方案：选择文本
        this.$refs.urlInput.select()
        document.execCommand('copy')
      }
    }
  }
}
</script>

<style scoped>
.share-dialog {
  padding: 20px;
  max-width: 500px;
  margin: 0 auto;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
}

.form-control {
  width: 100%;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.actions {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
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

.share-result {
  text-align: center;
}

.url-display {
  display: flex;
  gap: 10px;
  margin: 20px 0;
}

.url-display input {
  flex: 1;
}

.expires-info {
  color: #666;
  margin-bottom: 20px;
}
</style> 