package main

import (
	"flag"
	"fmt"
	"k8s-client-go/common"
	"k8s-client-go/k8scrdctrl/controller"
	"k8s-client-go/k8scrdctrl/pkg/client/clientset/versioned"
	"k8s-client-go/k8scrdctrl/pkg/client/informers/externalversions"
	"k8s.io/client-go/informers"
	"k8s.io/klog"
	"time"
)

func main() {

	// 日志参数
	klog.InitFlags(nil)
	flag.Set("logtostderr", "1") // 输出日志到stderr
	flag.Parse()

	// 读取admin.conf, 生成客户端基本配置
	restConf, err := common.GetRestConf("./conf/admin.conf")
	if  err != nil {
		fmt.Println(err)
	}

	// 创建CRD的client
	crdClientset, err := versioned.NewForConfig(&restConf)
	if err != nil {
		fmt.Println(err)
	}

	// 创建K8S内置的client
	clientset, err := common.InitClient("./conf/admin.conf")
	if err != nil {
		fmt.Println(err)
	}

	// 内建informer工厂
	informerFactory := informers.NewSharedInformerFactory(&clientset, time.Second * 120)
	// crd Informer工厂
	crdInformerFactory := externalversions.NewSharedInformerFactory(crdClientset, time.Second * 120)

	// POD informer
	podInformer := informerFactory.Core().V1().Pods()
	// nginx informer
	nginxInformer := crdInformerFactory.Mycompany().V1().Nginxes()

	// 创建调度controller
	nginxController := &controller.NginxController{Clientset: &clientset, CrdClientset: crdClientset, PodInformer:podInformer, NginxInformer: nginxInformer}
	nginxController.Start()

	// 等待
	for {
		time.Sleep(1 * time.Second)
	}

	return
}
