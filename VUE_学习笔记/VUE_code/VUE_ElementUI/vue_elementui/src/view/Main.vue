<template>
  <el-container>
    <el-aside width="200px">
      <!-- 默认展开第一个el-submenu-->
      <el-menu :default-openeds="['1']">
        <el-menu-item index="0">
          <router-link to="/goHome" tag="li">首页</router-link>
        </el-menu-item>
        <el-submenu index="1">
          <template slot="title"><i class="el-icon-caret-right"></i>用户管理</template>
          <el-menu-item-group>
            <el-menu-item index="/">
              <!--插入跳转链接-->
              <!-- <router-link to="/user/profile">个人信息</router-link>-->
              <!--传入参数（方式一，耦合较高）：name为组件名字（调用index.js中name为Profile的组件）， params中的id为参数-->
              <router-link v-bind:to="{name: 'Profile', params: {id:1}}" tag="li">个人信息</router-link>
            </el-menu-item>
            <el-menu-item index="1-2">
              <router-link to="/user/list" tag="li">用户列表</router-link>
            </el-menu-item>
          </el-menu-item-group>
        </el-submenu>
        <el-submenu index="2">
          <template slot="title"><i class="el-icon-menu"></i>内容管理</template>
          <el-menu-item-group>
            <el-menu-item index="2-1">分类列表</el-menu-item>
            <el-menu-item index="2-2">内容列表</el-menu-item>
          </el-menu-item-group>
        </el-submenu>

      </el-menu>
    </el-aside>

    <el-container>
      <el-header style="text-align: right; font-size: 12px">
        <el-dropdown>
          <i class="el-icon-setting" style="margin-right: 15px"></i>
          <el-dropdown-menu slot="dropdown">
            <el-dropdown-item>个人信息</el-dropdown-item>
            <el-dropdown-item>退出登录</el-dropdown-item>
          </el-dropdown-menu>
        </el-dropdown>
      </el-header>

      <el-main>
        <!--在这里展示视图-->
        <router-view />
      </el-main>

    </el-container>

  </el-container>
</template>

<script>
export default {
  name: "Main",
}
</script>

<style scoped>
.el-header {
  background-color: #B3C0D1;
  color: #333;
  line-height: 60px;
}

.el-aside {
  color: #333;
}
</style>

