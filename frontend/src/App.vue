<template>
  <div id="app">
    <AppHeader />
    <Home />
    <AppFooter />
    <!--
    <router-view></router-view>
    -->
  </div>
</template>

<script>
import Home from './components/Home.vue'
import AppHeader from './components/AppHeader.vue'
import AppFooter from './components/AppFooter.vue'

export default {
  name: 'app',
  components: {
    Home,
    AppHeader,
    AppFooter
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
}
</style>
