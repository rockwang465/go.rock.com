package k8s

import (
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetClusterNodes(k8sConf string) (*v1.NodeList, error) {
	clientSet, err := GetK8sClient(k8sConf)
	if err != nil {
		panic(err)
	}

	return clientSet.CoreV1().Nodes().List(metaV1.ListOptions{})
}

func GetClusterNode(k8sConf, nodeName string) (*v1.Node, error) {
	clientSet, err := GetK8sClient(k8sConf)
	if err != nil {
		panic(err)
	}

	return clientSet.CoreV1().Nodes().Get(nodeName, metaV1.GetOptions{})
}
