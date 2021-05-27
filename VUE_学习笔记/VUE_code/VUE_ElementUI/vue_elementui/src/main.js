// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'

//使用路由
import VueRouter from "vue-router";

//自动扫描路由
import router from './router';

//（完整引入element）导入element ui、element css
import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';

//Axios
import axios from 'axios';
import VueAxios from 'vue-axios';

//使用配置
Vue.use(VueRouter);
Vue.use(ElementUI);
Vue.use(VueAxios, axios);

Vue.config.productionTip = false

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  //element-ui中要求使用的方式
  render: h => h(App)
  // components: { App },
  // template: '<App/>'
})
