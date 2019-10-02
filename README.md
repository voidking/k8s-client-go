# k8s-client-go
学习client-go用法的demo，改自[k8s-client-go](https://github.com/owenliang/k8s-client-go)。

# 清单

* common: 初始化连接，deployment相关，pod相关，service相关，deploy相关
* login-pre: xterm.js的基本用法, 为后续web ssh访问k8s container做铺垫
* login: xterm.js+client-go remotecommand实现完美web ssh登录container
* klog: client-go的sdk日志配置
* k8scrd: 自定义CRD，利用code generation生成controller骨架代码
* k8scrdctrl: 实现一个类似于replicas的controller，动态管理POD数量

# 参考

[client-go doc](https://godoc.org/k8s.io/client-go/kubernetes)
