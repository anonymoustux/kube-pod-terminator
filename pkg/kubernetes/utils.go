package kubernetes

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"time"
)

func GetConfig(masterUrl, kubeConfigPath string) (*rest.Config, error) {
	config, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeConfigPath)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func GetClientSet(config *rest.Config) (*kubernetes.Clientset, error) {
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientSet, nil
}

func getTerminatingPods(clientSet *kubernetes.Clientset, namespace string) ([]v1.Pod, error) {
	var resultSlice []v1.Pod
	pods, err := clientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, pod := range pods.Items {
		deletionTimestamp := pod.ObjectMeta.DeletionTimestamp
		if deletionTimestamp != nil && deletionTimestamp.Add(2 * time.Millisecond).Before(time.Now()) {
			resultSlice = append(resultSlice, pod)
		}
	}
	return resultSlice, nil
}