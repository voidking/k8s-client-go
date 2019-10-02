package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
)

// 读取service yaml文件
func ReadServiceYaml(filename string) (core_v1.Service) {
	// 读取YAML
	serviceYaml, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return core_v1.Service{}
	}

	// YAML转JSON
	deployJson, err := yaml.ToJSON(serviceYaml)
	if err != nil {
		fmt.Println(err)
		return core_v1.Service{}
	}

	// JSON转struct
	var service core_v1.Service
	err = json.Unmarshal(deployJson, &service)
	if err != nil {
		fmt.Println(err)
		return core_v1.Service{}
	}
	return service
}

func ApplyService(clientset kubernetes.Clientset, new_service core_v1.Service) {

	// 查询k8s是否有该service
	service, err := clientset.CoreV1().Services("default").Get(new_service.Name, meta_v1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			fmt.Println(err)
			return
		}
		// 不存在则创建
		_, err = clientset.CoreV1().Services("default").Create(&new_service)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else { // 已存在则更新
		//service.Spec.Selector = new_service.Spec.Selector
		new_service.ResourceVersion = service.ResourceVersion
		new_service.Spec.ClusterIP = service.Spec.ClusterIP
		_, err = clientset.CoreV1().Services("default").Update(&new_service)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Printf("apply service %s success!", service.Name)
}

