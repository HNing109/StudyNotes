<template>
  <div id="app">
<!--    <img src="./assets/logo.png">-->
<!--    <Login></Login>-->
    <!-- 显示视图 -->
    <router-view></router-view>
  </div>
</template>

<script>

import Login from "./view/Login";
export default {
  name: 'App',
  components: {
    // Login
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
