package k8s

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// get configMap by k8sConf and namespace
func GetInstanceConfig(k8sConf, ns, instanceName string) (*corev1.ConfigMapList, error) {
	client, err := GetK8sClient(k8sConf)
	if err != nil {
		return nil, err
	}

	labelSelector := generateInstanceLabel(instanceName)
	configMapList, err := client.CoreV1().ConfigMaps(ns).List(metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil, err
	}
	return configMapList, nil
}

// generate label by instance_name
func generateInstanceLabel(instanceName string) string {
	return fmt.Sprintf("app.kubernetes.io/instance=%s", instanceName)
}

// get pods name by namespace, k8s config, instance name
func GetInstancePods(k8sConf, namespace, instanceName string) (*corev1.PodList, error) {
	client, err := GetK8sClient(k8sConf)
	if err != nil {
		return nil, err
	}

	labelSelector := generateInstanceLabel(instanceName)
	podList, err := client.CoreV1().Pods(namespace).List(metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil, err
	}
	return podList, nil
}

// get the instance pod log by pod name and container name
func GetInstanceLog(k8sConf, namespace, podsName, containerName string, entireLog bool) (string, error) {
	client, err := GetK8sClient(k8sConf)
	if err != nil {
		return "", err
	}

	podLogOpts := corev1.PodLogOptions{
		Container: containerName,
	}

	// if not get entire log, then define log size
	if !entireLog {
		var logMaxSize int64 = 1048576 // 1mb
		var tailLines int64 = 5000
		podLogOpts.LimitBytes = &logMaxSize
		podLogOpts.TailLines = &tailLines
	}

	req := client.CoreV1().Pods(namespace).GetLogs(podsName, &podLogOpts)
	//podLogs, err := req.Stream()  // usage from Blog
	//if err != nil {
	//	return "", err
	//}
	////defer podLogs.Close()
	//buf := new(bytes.Buffer)
	//_, err = io.Copy(buf, podLogs)
	//if err != nil {
	//	return "", err
	//}
	//log := buf.String()
	//return log, nil

	resp := req.Do()
	output, err := resp.Raw()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
