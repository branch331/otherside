import Vue from 'vue'
import App from './App.vue'
import VueRouter from 'vue-router'

Vue.config.productionTip = false

Vue.use(VueRouter)

const Home = { template: '<div>Home</div>' }
const Callback = { template: '<div>Callback</div>' }

export const router = new VueRouter({
  mode: 'history',
  base: __dirname,
  routes: [
    { path: '/', name: 'home', component: Home },
    { path: '/callback', name: 'callback', component: Callback}
  ]
})

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
