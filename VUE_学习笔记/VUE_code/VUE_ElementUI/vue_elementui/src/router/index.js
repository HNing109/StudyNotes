//导入vue
import Vue from 'vue';
import VueRouter from 'vue-router'

//导入自定义组件
import Main from '../view/Main';
import Login from "../view/Login";
//导入子组件：路由嵌套
import Profile from "../view/user/Profile";
import UserList from "../view/user/UserList";
import NotFound from "../view/error/NotFound";

//使用路由
Vue.use(VueRouter);

//配置路由
export default new VueRouter({
  // hash（默认值）：路径带 # 符号，如 http://localhost/#/login
  // history：路径不带 # 符号，如 http://localhost/login
  mode: 'history',
  routes:[
    {
      path: '/main',
      component: Main,
      //写入子模块（嵌套路由）
      children: [
        // {
        //   //传入参数（方式一，耦合较高）：  :id 为传递的参数，使用此url需要传入一个参数
        //   path: '/user/profile/:id',
        //   name: 'Profile',
        //   component: Profile
        // },
        {
          //传入参数（方式二）：  :id 为传递的参数，使用此url需要传入一个参数
          path: '/user/profile/:id',
          name: 'Profile',
          component: Profile,
          //打开参数属性
          props: true
        },
        {
          path: '/user/list',
          component: UserList
        }
      ]
    },
    {
      path: '/login',
      component: Login
    },
    //重定向：给url取别名
    {
      path: '/goHome',
      redirect: '/main'
    },
    //配置404路径
    {
      path: '*',
      component: NotFound
    }
  ]
});
