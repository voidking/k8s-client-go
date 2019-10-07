package main

import (
	"fmt"
	"k8s-client-go/common"
	"k8s.io/client-go/kubernetes"
)

func main() {

	// 初始化k8s客户端
	clientset, err := common.InitClient("./conf/admin.conf")
	if err != nil {
		fmt.Println(err)
		return
	}

	//testDeployment(clientset)
	testDeploy(clientset)
}

func testPod(clientset kubernetes.Clientset){
	podList, err := common.GetPodsByNamespace(clientset, "default")
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(podList)

	common.GetContainerLog(clientset, "default", "nginx-689c895cb5-d7zp9", "nginx")
}

func testDeployment(clientset kubernetes.Clientset){
	deployment := common.ReadDeploymentYaml("./conf/nginx.yaml")
	//common.ApplyDeployment(clientset, deployment)
	success, reasons, err := common.GetDeploymentStatus(clientset, deployment)
	fmt.Println(success)
	fmt.Println(reasons)
	fmt.Println(err)
	//common.PrintDeploymentStatus(clientset, deployment)
}

func testService(clientset kubernetes.Clientset){
	service := common.ReadServiceYaml("./conf/nginx-service.yaml")
	common.ApplyService(clientset,service)

	// 修改service的label，分流给灰度pod
	service.Spec.Selector = map[string]string{
		"app": "nginx",
		//"track": "stable",
	}
}

func testDeploy(clientset kubernetes.Clientset){
	deployment := common.ReadDeploymentYaml("./conf/nginx.yaml")
	//new_deployment := common.ReadDeploymentYaml("./conf/new-nginx.yaml")
	//service := common.ReadServiceYaml("./conf/nginx-service.yaml")

	//common.GrayDeploy(clientset, deployment, new_deployment, service)
	//common.UpdateDeploy(clientset, deployment, new_deployment, service, 1)
	//common.RollBack(clientset, deployment, new_deployment)
	//common.GrayDeploy2(clientset, deployment, "voidking/nginx:v2.0", 1)
	common.UpdateDeploy2(clientset, deployment, "voidking/nginx:v2.0")
}
