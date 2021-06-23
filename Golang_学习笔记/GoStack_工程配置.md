# 1、go env配置

- go编译器版本：

  15.12

  

- 配置命令： 

  go env -w GO111MODULE=on
  go env -w GOPROXY=https://goproxy.io,direct

  go env -w GOINSECURE=git.ctyun.cn
  go env -w GOPRIVATE=git.ctyun.cn

   

  ```shell
  查看是否生效：
  windows：C:\Users\Lenovo\AppData\Roaming\go
  
  GO111MODULE=on
  GOINSECURE=git.ctyun.cn
  GOPRIVATE=git.ctyun.cn
  GOPROXY=https://goproxy.io,direct
  ```

  

  使用go env命令（查看配置）：

  ```shell
  C:\Code\gostack>go env
  set GO111MODULE=on
  set GOARCH=amd64
  set GOBIN=
  set GOCACHE=C:\Users\Lenovo\AppData\Local\go-build
  set GOENV=C:\Users\Lenovo\AppData\Roaming\go\env
  set GOEXE=.exe
  set GOFLAGS= -mod=
  set GOHOSTARCH=amd64
  set GOHOSTOS=windows
  set GOINSECURE=git.ctyun.cn
  set GOMODCACHE=C:\Users\Lenovo\go\pkg\mod
  set GONOPROXY=git.ctyun.cn
  set GONOSUMDB=git.ctyun.cn
  set GOOS=windows
  set GOPATH=C:\Users\Lenovo\go
  set GOPRIVATE=git.ctyun.cn
  set GOPROXY=direct
  set GOROOT=C:\Program Files\Go
  set GOSUMDB=sum.golang.org
  set GOTMPDIR=
  set GOTOOLDIR=C:\Program Files\Go\pkg\tool\windows_amd64
  set GCCGO=gccgo
  set AR=ar
  set CC=gcc
  set CXX=g++
  set CGO_ENABLED=1
  set GOMOD=C:\Code\gostack\go.mod
  set CGO_CFLAGS=-g -O2
  set CGO_CPPFLAGS=
  set CGO_CXXFLAGS=-g -O2
  set CGO_FFLAGS=-g -O2
  set CGO_LDFLAGS=-g -O2
  set PKG_CONFIG=pkg-config
  set GOGCCFLAGS=-m64 -mthreads -fno-caret-diagnostics -Qunused-arguments -fmessage-length=0 -fdebug-prefix-map=C:\Users\Lenovo\AppData\Local\Temp\go-build594281098=/tmp/go-build -gno-re
  cord-gcc-switches
  ```

  

- 获取go.mod中配置的包：

  命令：go mod download

  

  Q：若上述命令执行后未出现错误，但Gostack工程依然不能import外部包

  A：方式一：重启Goland

  ​      方式二：删除Gostack工程，重新git clone工程，打开Goland编译器，即可自动导入外部包。
  
  

- **<font color='red'>在获取GitLab中的包时，需要填写自己的Gitlab账号、密码</font>** 

  

# 2、Goland配置

- 无需配置Goland，否则会出现包无法使用的情况

  ![image-20210622101124589](GoStack工程配置.assets/image-20210622101124589.png)
