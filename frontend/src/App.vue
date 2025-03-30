<template>
  <div class="app-container">
    <header class="header">
      <h1>Go Playground</h1>
    </header>
    <main>
      <router-view />
    </main>
    <footer class="footer">
      <p>Powered by Go {{ goVersion }}</p>
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
}

.app-container {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.header {
  background-color: #00ADD8;
  color: white;
  padding: 1rem;
  text-align: center;
}

main {
  flex: 1;
  padding: 1rem;
}

.footer {
  background-color: #f5f5f5;
  text-align: center;
  padding: 1rem;
  font-size: 0.8rem;
}
</style> 