<template>
  <div>
    <h1>个人信息</h1>
    <!--传入参数（方式一，耦合较高）：获取传入的参数id-->
    <!--{{$route.params.id}}-->

    <!--传入参数（方式二）：获取传入的参数id-->
    {{id}}
  </div>
</template>

<script>
export default {
  name: "Profile",
  //传入参数（方式二）:定义参数：id
  props: ['id'],
  //配置钩子（进入前拦截）：类似于java中的过滤器。to：路由将要跳转的路径信息；from：路径跳转前的路径信息；next：路由的控制参数
  beforeRouteEnter: (to, from, next) => {
    console.log("准备进入Profile个人信息界面");
    //跳入下一个页面：进入页面之前，调用getData方法(vm 就是当前组件的实例, 相当于this)
    next(vm => {
      vm.getDate()
    });
  },
  //配置钩子（离开前拦截）
  beforeRouteLeave: (to, from, next) => {
    console.log("准备离开Profile个人信息界面");
    next();
  },
  methods:{
    //获取testjson.json的数据
    getDate: function(){
      this.axios({
        method: "get",
        url: "http://localhost:8080/static/mock/testjson.json"
      }).then(
        function(response){
          console.log(response)
        }
      )
    }
  }
}
</script>

<style scoped>

</style>
