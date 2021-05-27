<template>
  <div id="app">
    <h1>Vue-router测试</h1>
    <!--  使用自己的组件：Content、Main、Test  -->
    <router-link to="/main">首页</router-link>
    <router-link to="/test">测试</router-link>
    <router-link to="/content">内容</router-link>
    <!--  用于展示上面组件的视图：否则页面不会显示组件的内容  -->
    <router-view></router-view>
  </div>
</template>

<script>

export default {
  name: 'App',
  components: {
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
