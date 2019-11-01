import Vue from 'vue'
import App from './App.vue'
import Home from './components/Home.vue'
import Callback from './components/Callback.vue'
import VueRouter from 'vue-router'

Vue.config.productionTip = false

Vue.use(VueRouter)

const router = new VueRouter({
  routes: [
    { path: '/', name: 'home', component: Home },
    { path: '/callback', name: 'callback', component: Callback},
    { path: '/callback/:state', name: 'callback', component: Callback}
  ]
})

new Vue({
  render: h => h(App),
  router
}).$mount('#app')
