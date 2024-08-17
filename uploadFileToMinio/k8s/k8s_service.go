package k8s

import (
	"context"
	"fmt"
	"os"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func PodInfo() *v1.Pod {
	var config *rest.Config
	var err error

	// 尝试从环境变量获取配置
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err)
		}
	} else {
		// 如果不在集群内，则尝试从本地 kubeconfig 文件加载配置
		home := homedir.HomeDir()
		kubeconfigPath := home + string(os.PathSeparator) + ".kube" + string(os.PathSeparator) + "config"
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			panic(err)
		}
	}

	// 创建一个新的 Kubernetes 客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// 获取当前 Pod 的名称和命名空间
	podName := os.Getenv("POD_NAME")
	namespace := os.Getenv("POD_NAMESPACE")

	if podName == "" || namespace == "" {
		fmt.Println("环境变量 POD_NAME 或 POD_NAMESPACE 未设置")
		return nil
	}

	// 获取当前 Pod 的详细信息
	pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	// 打印 Pod 的信息
	fmt.Printf("Pod Name: %s\n", pod.Name)
	fmt.Printf("Pod Namespace: %s\n", pod.Namespace)
	fmt.Printf("Pod Status: %s\n", pod.Status.Phase)
	fmt.Printf("Pod Labels: %v\n", pod.Labels)
	return pod
}
