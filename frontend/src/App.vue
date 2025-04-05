<template>
  <div class="app-container">
    <header class="header">
      <h1>Go Playground</h1>
    </header>
    <main>
      <router-view />
    </main>
    <footer class="footer">
      <p>Powered by Go <a href="https://chenjie.info" target="_blank">chenjie.info</a></p>
    </footer>
  </div>
</template>

<script>
export default {
  name: 'App',
  data() {
    return {
      goVersion: '1.24'
    }
  }
}
</script>

<style>
body {
  margin: 0;
  padding: 0;
  font-family: Arial, sans-serif;
  overflow-x: hidden;
}

* {
  box-sizing: border-box;
}

html, body {
  height: 100%;
  width: 100%;
}

.app-container {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  width: 100%;
}

.header {
  background-color: #00ADD8;
  color: white;
  padding: 1rem;
  text-align: center;
}

.header h1 {
  margin: 0;
  font-size: 1.8rem;
}

main {
  flex: 1;
  padding: 1rem;
  width: 100%;
  max-width: 100%;
  display: flex;
  flex-direction: column;
}

/* 确保router-view占满可用空间 */
main > * {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.footer {
  background-color: #f5f5f5;
  text-align: center;
  padding: 1rem;
  font-size: 0.8rem;
  margin-top: auto; /* 确保footer始终在底部 */
}

/* 媒体查询 */
@media (max-width: 768px) {
  main {
    padding: 0.5rem;
    padding-bottom: 1rem;
  }
  
  .header h1 {
    font-size: 1.5rem;
  }
  
  .footer {
    padding: 0.5rem;
  }
}

@media (max-width: 480px) {
  main {
    padding: 0.25rem;
    padding-bottom: 1rem;
  }
  
  .header {
    padding: 0.75rem;
  }
  
  .header h1 {
    font-size: 1.2rem;
  }
}
</style> 