import Vue from 'vue'
import App from './App.vue'
import Home from './components/Home.vue'
import VueRouter from 'vue-router'
import BootstrapVue from 'bootstrap-vue'
import PortalVue from 'portal-vue'
import VueMeta from 'vue-meta'

require('./assets/styles.css')
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

Vue.config.productionTip = false

Vue.use(VueRouter)
Vue.use(BootstrapVue)
Vue.use(PortalVue)
Vue.use(VueMeta, {
  // optional pluginOptions
  refreshOnceOnNavigation: true
})

const router = new VueRouter({
  routes: [
    { path: '/', name: 'home', component: Home }
  ]
})

new Vue({
  render: h => h(App),
  router
}).$mount('#app')
