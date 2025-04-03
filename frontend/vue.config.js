module.exports = {
  devServer: {
    port: 3003,
    proxy: {
      '/api/share': {
        target: 'http://share-service:3002',
        changeOrigin: true,
        ws: true
      },
      '/api/execute': {
        target: 'http://share-service:3002',
        changeOrigin: true,
        ws: true
      },
      '/api/go1.24': {
        target: 'http://backend-go124:3001',
        changeOrigin: true,
        ws: true,
        pathRewrite: {
          '^/api/go1.24': '/api'
        }
      },
      '/api/go1.23': {
        target: 'http://backend-go123:3001',
        changeOrigin: true,
        ws: true,
        pathRewrite: {
          '^/api/go1.23': '/api'
        }
      },
      '/api/go1.22': {
        target: 'http://backend-go122:3001',
        changeOrigin: true,
        ws: true,
        pathRewrite: {
          '^/api/go1.22': '/api'
        }
      },
      '/api': {
        target: 'http://backend-go124:3001',
        changeOrigin: true,
        ws: true
      }
    }
  },
  configureWebpack: {
    resolve: {
      extensions: ['.js', '.vue', '.json']
    }
  }
} 