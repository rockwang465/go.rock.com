package k8s

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// restart all pods at engineNamespace namespace
func RestartPodsWithLicense(k8sConf string, timeout int64) error {
	client, err := GetK8sClient(k8sConf)
	if err != nil {
		return err
	}

	podListOption := metav1.ListOptions{ // configuration timeout
		TimeoutSeconds: &timeout,
	}
	return client.CoreV1().Pods(engineNamespace).DeleteCollection(&metav1.DeleteOptions{}, podListOption)
}
