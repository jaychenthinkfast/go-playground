<template>
  <div class="home">
    <Playground />
  </div>
</template>

<script>
import Playground from '@/components/Playground.vue'

export default {
  name: 'HomeView',
  components: {
    Playground
  }
}
</script>

<style scoped>
.home {
  height: 100vh;
  display: flex;
  flex-direction: column;
}
</style> 