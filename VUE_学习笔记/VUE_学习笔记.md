# 1、VUE开发环境配置

参考链接（含创建、启动项目）：https://www.cnblogs.com/dybk/p/10643067.html

## 1.1、安装node.js、npm

- node官网：https://nodejs.org/en/download/   （新版的node安装后自带npm）

- 在nodejs的安装目录下，新建node_global和node_cache两个文件夹，然后在node commande prompt命令工具中配置 这两个文件的绝对位置。

  - npm config set prefix "D:\Program Files\nodejs\node_global"

  - npm config set cache "D:\Program Files\nodejs\node_cache"

- 配置环境变量

  - 新增用户变量：在Path中新增  D:\Program Files\nodejs\node_global”
  - 新增系统变量NODE_PATH：创建NODE_PATH，变量值为D:\Program Files\nodejs\node_global\node_modules

- （可不安装）安装cnpm，这个工具node已经自带，但是没有配置环境变量，因此无法在node commande prompt中使用（自己仿照上一步，添加环境变量之后就可以使用）

- 安装VUE、VUE命令行工具、webpack、vue/cli工具：

  - 安装VUE：  -g全局安装

    cnpm install vue -g

  - 安装VUE命令行工具：

    cnpm install vue-cli -g

  - 安装webpack：

    npm install webpack -g      或者     cnpm install webpack -g

  - vue/cli工具：

    npm install --g @vue/cli   或者  cnpm install --g @vue/cli   

    （若安装错误，以管理员身份运行命令，先清除缓存：npm clean cache –force） 

    

- 测试安装否成功：（v注意大小写）

  - node -v
  - npm -V
  - vue -V
  - webpack -v
  
  
  
- 使用命令创建VUE工程：

  - 下载模板：cd至需要存放工程的目录，创建mytest工程，vue init webpack mytest
  - 安装模块：cd至工程目录的根目录，npm install
  - 测试运行：npm run dev



----

# 2、WebStorm开发VUE项目

参考网站：

- 配置、启动vue项目：https://blog.csdn.net/jenybom/article/details/90213374
- 创建项目（该方式不适用于2021的webstorm，无法创建）：https://blog.csdn.net/qq_37350706/article/details/86591102

- 启动github下载的代码：

  - 在项目所在的根目录：下载工程依赖模块， npm install

  - 配置相应的启动参数：

    ![image-20210518171648037](VUE_学习笔记.assets/image-20210518171648037.png)

    ![image-20210518171808478](VUE_学习笔记.assets/image-20210518171808478.png)

    

  - 启动后的样子：

    ![image-20210518171926537](VUE_学习笔记.assets/image-20210518171926537.png)

    



---

# 3、VUE基础

**参考笔记**：[狂神笔记]https://blog.csdn.net/okForrest27/article/details/106849246/

## 3.1、前端三要素

- HTML：超文本标记语言，决定网页的结构、内容

- CSS：层叠样式表，设计网页的表现样式

- JavaScript：无需编译，浏览器可直接解释运行，用于控制网页的行为

  - 标准：按照[ECMAScript] 标准的开发方式，简称是ES,

    ES4 (内部,未征式发布)
    ES5 (全浏览器支持)
    ES6 (常用，当前主流版本: webpack打包成为ES5支持! )

  - 框架：
    - **jQuery**
    - **React**
    - **Vue**
    - **Axios** :前端通信框架;因为Vue 的边界很明确，就是为了处理DOM,所以并不具备通信能力，此时就需要额外使用一个通信框架与服务器交互;当然也可以直接选择使用jQuery提供的AJAX通信功能

- UI框架：

  - Ant-Design:阿里巴巴出品，基于React的UI框架
  - ElementUI、 iview、 ice: 饿了么出品，基于Vue的UI框架
  - Bootstrap: Twitter推出的一个用于前端开发的开源工具包
  - AmazeUI:又叫"妹子UI"，一款HTML5跨屏前端框架.
    JavaScript 构建工具
  - Babel: JS编译工具，主要用于浏览器不支持的ES新特性，比如用于编译TypeScript
  - WebPack: 模块打包器，主要作用是打包、压缩、合并及按序加载。能把各种资源，如 JS、JSX、ES6、SASS、LESS、图片等**都作为模块来处理和使用**。

- 实现三端统一开发的框架（windows、android、ios共用一套代码）
  - 打包方式：
    - 云打包: **HBuild -> HBuildX， DCloud出品; API Cloud**
    - 本地打包: **Cordova** (前身是PhoneGap)
  - 三端统一的框架：**NodeJS **，其框架、工具如下：
    - Express: NodeJS框架
    - Koa: Express简化版
    - NPM: 项目综合管理工具，类似于Maven
    - YARN: NPM的替代方案，类似于Maven和Gradle的关系



----

## 3.2、VUE入门

- MVVM模式：

  - Mode层：模型层，对应JavaScript对象
  - View层：视图层，也就是DOM（HTML中的标签元素）
  - ViewMode层（**核心部分**）：即，Vue.js。连接视图和数据的中间件，监听、观察数据变化，并对视图进行更新。

  **MVVM的数据传输流程：**

  ![image-20210524093317091](VUE_学习笔记.assets/image-20210524093317091.png)



- VUE的生命周期：

  开始创建**→**初始化数据**→**编译模板**→**挂载DOM→渲染→更新→渲染→卸载

  

-  基本语法：（**<font color='red'>以下示例，在HTML文件中运行，实际开发则在.vue文件，使用vue-cli创建.vue模板文件</font>**）

  - if-else语法： 

    ```HTML
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <title>Title</title>
        <!--1.导入vue.js-->
        <script src="https://cdn.jsdelivr.net/npm/vue@2.5.21/dist/vue.min.js"></script>
    </head>
    <body>
    <div id="app">
        <h1 v-if="ok">Yes</h1>
        <h1 v-else>No</h1>
    </div>
    <script>
        //声明变量：let（局部变量）、var（全局变量）
        let vm = new Vue({
            el: "#app",
            data: {
                ok: true
            }
        });
    </script>
    </body>
    </html>
    ```

  -  for

    ``` HTML
    <div id="app">
        <!--读取vm对象data部分数组中对象的message数据-->
        <li v-for="text in allTexts">
            {{text.message}}
        </li>
    </div>
    <script>
        let vm = new Vue({
            el: "#app",
            data: {
                allTexts: [
                    {message: "1"},
                    {message: "2"},
                    {message: "3"},
                ]
            }
        });
    </script>
    ```

  - 事件绑定（eg：click鼠标点击事件）

    ````HTML
    <div id="app">
        <!-- v-on:click为button按键绑定方法 -->
        <button v-on:click="sayHi()">点我</button>
    </div>
    <script>
        let vm = new Vue({
            el: "#app",
            data: {
                message: "xxxxxx"
            },
            methods: {
                //event为传入事件，调用syHi方法
                sayHi: function (event) {
                    alert(this.message);
                }
            }
        });
    </script>
    ````

    

  - 双向数据绑定：

    数据、View视图层相互绑定：两者的数据是实时同步的。

    ````HTML
    <div id="app">
        <select v-model="selected">
            <option value="" disabled>--请选择--</option>
            <option>A</option>
            <option>B</option>
            <option>C</option>
          </select>
          <p>当前选中的是：{{selected}}</p>
    </div>
    <script>
    	var vm = new Vue({
        el: "#app",
        //1.数据
        data: {
          //双向绑定：页面选择后，直接在页面显示当前选择的数据
          selected: ''
          }
        });
    </script>
    ````

    

  - 组件（VUE的模板功能--**核心部分**）：

    相当于HTML文件中的JSTL标签，自定义的模板可通过引用标签调用，实现模板的重复套用。

    需要使用Vue.component()方法组测组件。

    ```html
    <div id="app">
    	<!-- 遍历items数组,获取message参数，传给模板中的element  -->
        <chris v-for="item in items" v-bind:element="item.message"></chris>
    </div>
    <script>
        //自定义VUE组件  （组件需要放在script的头部，否则无法使用）  chris为自定义组件的名称
        Vue.component("chris", {
            //参数：使用v-bind传入参数
            props: ["element"],
            //模板：获取上面的参数element
            template: '<li> {{element}} </li>'
        });
        
        var vm = new Vue({
            el: "#app",
            //1.数据
            data: {
              //数组
              items: [
                //对象：message
                {message: "chris"},
                {message: "FYJ"},
                {message: "Vickey"}
              ],
            }
        });
    </script>
    ```

    

  - Axios通信组件

     Axios是一个开源的可以用在浏览器端和NodeJS 的异步通信框架，主要作用就是实现AJAX异步通信，其功能特点如下:

    - 从浏览器中创建XMLHttpRequests
    - 从node.js创建http请求
    - 支持Promise API [JS中链式编程]
    - 拦截请求和响应
    - 转换请求数据和响应数据
    - 取消请求
    - 自动转换JSON数据
    - 客户端支持防御XSRF (跨站请求伪造)

    ```html
    <div id = "app1">
      <div>
        {{info.name}}
      </div>
      <div>
        {{info.links}}
      </div>
      <a v-bind:href="info.links.address">跳转{{info.links.name}}</a>
    </div>
    
    <script>
      //Axios异步通信：替代JQuery、Ajax，读取json文件
      var vm1 = new Vue({
        el: "#app1",
    
        //data方法：使用该方法返回Json格式的数据
        data () {
          return {
            info: {
              name: null,
              links: [
                {
                  name: null,
                  address: null
                }
              ]
            }
          }
        },
    
        //钩子函数：读取Json文件的数据，返回给浏览器（ES6特性）
        mounted () {
          axios.get('../../data/testjson.json').then(response=>(this.info=response.data))
        }
      });
    </script>
    ```

    

  - **计算属性：**

     计算属性和调用方法不同：

    - 计算属性：可缓存数据（**仅当计算属性中所调用方法内的数据发生变化时，才会再次调用方法计算数据，否则直接返回上次的计算结果，等同于redis缓存**）
    - 调用方法：每次都会调用方法进行计算。导致处理速度变慢。

    ```html
    <div>
        <span>当前时间（method获取）：{{getTimeByMethods()}}</span> <br/>
        <span>当前时间（computed获取）：{{getTimeByComputed}}</span>
    </div>
    <script>
        var vm=new Vue({
            //方法：需要绑定在Vue对象内  (不具备缓存，需要实时计算)
            methods:{
              getTimeByMethods: function(){
                return Date.now();
              }
            },
    
            //计算属性: (具备缓存，不需要实时计算，有点像redis缓存，仅当该方法中的数据发生变化时，才会再次进行计算)
            computed:{
              getTimeByComputed: function(){
                this.type;
                return Date.now();
              }
            }
        });
    </script>
    ```

    

  - 插槽：slot

    可结合模板使用，将数据循环读出

    - 缩写：

      v:bind: 可以缩写为一个:
      v-on: 可以缩写为一个@ 

     ```html
     <div id="app">
       <mylist>
         <!--使用自定义的slot_title，为slot_title组件中的title属性赋值-->
         <slot_title slot="slot_title" v-bind:title="title"></slot_title>
         <!--使用自定义的slot_itms-->
         <slot_items slot="slot_items" v-for="item in arr" v-bind:item = item></slot_items>
       </mylist>
     </div>
     
     <script>
       <!--自定义组件、模板-->
       Vue.component("mylist",{
         template:
           '<div>' +
             //使用插槽：结合模板使用
             '<slot name="slot_title"></slot>'+
             '<ul>' +
               '<slot name="slot_items"></slot>'+
             '</ul>'+
           '</div>'
       });
     
       Vue.component("slot_title",{
         //属性
         props: ['title'],
         template: '<div>{{title}}</div>'
       });
     
       Vue.component("slot_items",{
         props: ['item'],
         template: '<div>{{item}}</div>'
       });
     
       var v = new Vue({
         el: "#app",
         data:{
           title: "chris-title",
           arr: ['aaaa', 'bbbb', 'cccc']
         }
       });
     </script>
     ```

    

  - 自定义事件：

    - 使用this.$emit('自定义事件名', 参数)

    - 可用于设定数据分发、触发事件的方式

    - 结构：

      ![image-20210524103958500](VUE_学习笔记.assets/image-20210524103958500.png)

    ```html
    <div id="app">
      <mylist>
        <!--使用自定义的slot_title-->
        <slot_title slot="slot_title" v-bind:title="title"></slot_title>
          
        <!--使用自定义的slot_itms ：
        需要为slot_items对应的item、index、remove赋值(这边的v-on:remove为this.$emit(）定义的事件名)
        -->
        <slot_items slot="slot_items"
                    v-for="(item, index) in arr"
                    v-bind:item = "item"
                    v-bind:index = "index"
                    v-on:remove="removeItem(index)"
        ></slot_items>
      </mylist>
    </div>
    
    <script>
      //slot 插槽 这个组件要定义在前面，不然出不来数据
      <!--自定义组件、模板， myList为自定义组件的名称-->
      Vue.component("mylist",{
        template:
          '<div>' +
          //使用插槽：结合模板使用
          '<slot name="slot_title"></slot>'+
          '<ul>' +
          '<slot name="slot_items"></slot>'+
          '</ul>'+
          '</div>'
      });
    
      Vue.component("slot_title",{
        props: ['title'],
        template: '<div>{{title}}</div>'
      });
    
      Vue.component("slot_items",{
        props: ['item', 'index'],
          
        //为每个元素添加删除按钮
        template:'<div>'+
          '{{index}} - {{item}} '+
          ' <button style="margin: 5px" v-on:click="removeFunc">删除</button>'+
          '</div>',
          
        methods: {
          removeFunc: function(index){
            // this.$emit('事件',参数) 自定义事件分发（远程调用方法）
            this.$emit('remove', index)
          }
        }
      });
    
      var v = new Vue({
        el: "#app",
        data:{
          title: "chris-title",
          arr: ['aaaa', 'bbbb', 'cccc']
        },
        methods:{
          removeItem: function(index){
            // 一次删除arr数组的一个元素
            this.arr.splice(index, 1)
            //浏览器控制台打印信息
            console.log("成功删除了: " + this.arr[index])
          }
        }
      });
    
    </script>
    ```

    

  -  

  -  



----

## 3.2、使用VUE开发

-  
-  
-  













