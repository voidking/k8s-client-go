package main

import (
	"fmt"
	"k8s-client-go/common"
	"k8s-client-go/k8scrd/pkg/client/clientset/versioned"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {

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

	// 获取CRD的nginx对象
	nginx, err := crdClientset.MycompanyV1().Nginxes("default").Get("nginx", v1.GetOptions{})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(nginx)

	return
}
