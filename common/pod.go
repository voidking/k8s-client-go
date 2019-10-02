package common

import (
	"fmt"
	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetPodsByNamespace(clientset kubernetes.Clientset, namespace string) (core_v1.PodList, error){
	// 获取命名空间下的所有POD
	podsList, err := clientset.CoreV1().Pods(namespace).List(meta_v1.ListOptions{});
	if err != nil {
		return core_v1.PodList{}, err
	}
	return *podsList, nil
}

//
func GetContainerLog(clientset kubernetes.Clientset, namespace string, podName string, containerName string){
	// 获取日志请求
	req := clientset.CoreV1().Pods(namespace).GetLogs(podName, &core_v1.PodLogOptions{Container: containerName})
	// req.Stream()也可以实现Do的效果

	// 发送请求
	res := req.Do()
	if res.Error() != nil {
		fmt.Println(res.Error())
		return
	}

	// 获取结果
	logs, err := res.Raw()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("容器输出:", string(logs))
}
