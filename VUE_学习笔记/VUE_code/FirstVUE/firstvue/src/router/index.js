/**
 * 配置路由router：url路径跳转
 * vue默认的配置文件都是index.js
 */

//导入vue、vue-router
import Vue from 'vue';
import VueRouter from 'vue-router';

//导入自定义的组件
import Content from "../components/Content";
import Main from "../components/Main";
import Test from "../components/Test";

//安装路由
Vue.use(VueRouter)

//配置路由
export default new VueRouter({
  routes: [
    //Content组件
    {
      //路径
      path: '/content',
      //路由名
      name: 'content',
      //需要跳转到的组件
      component: Content
    },
    //Main组件
    {
      path: '/main',
      name: 'main',
      component: Main
    },
    //Test组件
    {
      path: '/test',
      name: 'test',
      component: Test
    }
  ]
});
