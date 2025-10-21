module.exports = {
  devServer: {
    proxy: {
      '/upload': {
        target: 'http://localhost:8080',
        changeOrigin: true
      },
      '/stream': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  }
}