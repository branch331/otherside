<template>
  <div id="app">
    <img alt="Vue logo" src="./assets/logo.png">

      <p>Current route name: {{ $route.name }}</p>
      <ul>
        <li><router-link :to="{ name: 'home' }">home</router-link></li>
        <li><router-link :to="{ name: 'callback' }">callback</router-link></li>
      </ul>

    <LocalEvents/>
  </div>
</template>

<script>
import LocalEvents from './components/LocalEvents.vue'

export default {
  name: 'app',
  components: {
    LocalEvents
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
