package common

import (
	"bufio"
	"fmt"
	apps_v1beta1 "k8s.io/api/apps/v1beta1"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"os"
	"strconv"
)

// 灰度发布
func GrayDeploy(clientset kubernetes.Clientset, deployment apps_v1beta1.Deployment, new_deployment apps_v1beta1.Deployment, service core_v1.Service) {

	// 查看旧的deployment状态
	GetDeploymentStatus(clientset, deployment)
	now_replicas := *deployment.Spec.Replicas

	instanceNum := 0
	for {
		fmt.Println("请输入灰度发布的实例数量：")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		instanceNum, _ = strconv.Atoi(input.Text())
		fmt.Println("您输入的实例数量为：", instanceNum)
		if int32(instanceNum) > now_replicas {
			fmt.Println("灰度实例数量应该小于等于已有实例数量！")
		} else {
			break
		}
	}

	// 创建新的deployment
	replicas := int32(instanceNum)
	new_deployment.Spec.Replicas = &replicas
	ApplyDeployment(clientset, new_deployment)

	// 查看新的deployment状态
	PrintDeploymentStatus(clientset, new_deployment)

	ApplyService(clientset, service)

}

// 更新发布
func UpdateDeploy(clientset kubernetes.Clientset, deployment apps_v1beta1.Deployment, new_deployment apps_v1beta1.Deployment, service core_v1.Service, step int32) {

	// 查看旧的deployment状态
	GetDeploymentStatus(clientset, deployment)
	now_replicas := *deployment.Spec.Replicas

	var new_replicas int32
	var old_replicas int32

	if now_replicas < step {
		fmt.Println("发布pod数量过大，请重新设置！")
		return
	}

	old_replicas = now_replicas
	new_replicas = 0
	for {
		if old_replicas <= new_replicas {
			new_replicas = now_replicas
			old_replicas = int32(0)
		} else {
			new_replicas = new_replicas + step
			old_replicas = old_replicas - step
		}

		// 创建新的deployment
		new_deployment.Spec.Replicas = &new_replicas
		ApplyDeployment(clientset, new_deployment)

		// 查看新的deployment状态
		PrintDeploymentStatus(clientset, new_deployment)

		// 调整service selector
		ApplyService(clientset, service)

		// 缩容旧的deployment
		deployment.Spec.Replicas = &old_replicas
		ApplyDeployment(clientset, deployment)

		//查看新的deployment状态
		PrintDeploymentStatus(clientset, deployment)


		if new_replicas == now_replicas {
			fmt.Println("更新发布完成！")
			break
		} else {
			fmt.Printf("%d台发布完成，是否继续发布（Y/N）", new_replicas)
			input := bufio.NewScanner(os.Stdin)
			input.Scan()
			confirm := input.Text()
			if confirm == "Y"{
				continue
			}else{
				break
			}
		}
	}
}

func RollBack(clientset kubernetes.Clientset, deployment apps_v1beta1.Deployment, new_deployment apps_v1beta1.Deployment) {
	ApplyDeployment(clientset, deployment)
	new_replicas := int32(0)
	new_deployment.Spec.Replicas = &new_replicas
	ApplyDeployment(clientset, new_deployment)
	fmt.Println("回滚完成！")
}
