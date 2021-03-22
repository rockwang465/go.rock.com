package k8s

import (
	"fmt"
	"go.rock.com/rock-platform/rock/server/utils"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
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

// The scale usage doc:
// https://github.com/5bug/one-cd/blob/master/deployer/scale.go

// get instance replicas number
//func GetInstanceScale(k8sConf, namespace, chartName string) (*autov1.Scale, error) { // only support k8s 1.13
func GetInstanceScale(k8sConf, namespace, chartName string) (*int32, error) {
	client, err := GetK8sClient(k8sConf)
	if err != nil {
		return nil, err
	}

	// get K8S version
	version, err := client.ServerVersion()
	if err != nil {
		return nil, err
	}
	//fmt.Printf("%#v\n", version) // &version.Info{Major:"1", Minor:"18", GitVersion:"v1.18.5", GitCommit:"e6503f8d8f769ace2f338794c914a96fc335df0f", GitTreeState:"clean", BuildDate:"2020-06-26T03:39:24Z", GoVersion:"go1.13.9", Compiler:"gc", Platform:"linux/amd64"}
	majorVersion, err := strconv.Atoi(version.Major)
	if err != nil {
		panic(err)
	}
	minorVersion, err := strconv.Atoi(version.Minor)
	if err != nil {
		panic(err)
	}

	// check the K8S version
	if majorVersion == 1 && minorVersion < 16 { // support 1.13 ...
		scale, err := client.ExtensionsV1beta1().Deployments(namespace).GetScale(chartName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		return &scale.Spec.Replicas, nil
	} else if majorVersion == 1 && minorVersion > 16 { // support 1.16 1.17 1.18 ...
		scale, err := client.AppsV1().Deployments(namespace).GetScale(chartName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		return &scale.Spec.Replicas, nil
	} else {
		err := utils.NewRockError(404, 40400013, fmt.Sprintf("Not found this K8S version %v", version.GitVersion))
		return nil, err
	}

	// scale:
	// {
	//    "metadata": {
	//        "name": "aurora-system-service",
	//        "namespace": "default",
	//        "selfLink": "/apis/apps/v1/namespaces/default/deployments/aurora-system-service/scale",
	//        "uid": "4d4129ca-f510-4146-88d3-5532d5194dc7",
	//        "resourceVersion": "16971",
	//        "creationTimestamp": "2021-03-18T08:11:39Z"
	//    },
	//    "spec": {
	//        "replicas": 1
	//    },
	//    "status": {
	//        "replicas": 1,
	//        "selector": "app.kubernetes.io/instance=aurora-system-service-default,app.kubernetes.io/name=aurora-system-service"
	//    }
	//}
}

// update instance replicas number
//func UpdateInstanceScale(k8sConf, namespace, chartName string, scale *autov1.Scale) (*autov1.Scale, error) {
func UpdateInstanceScale(k8sConf, namespace, chartName string, number int32) (*int32, error) {
	client, err := GetK8sClient(k8sConf)
	if err != nil {
		return nil, err
	}

	// get K8S version
	version, err := client.ServerVersion()
	if err != nil {
		return nil, err
	}
	//fmt.Printf("%#v\n",version)  // &version.Info{Major:"1", Minor:"18", GitVersion:"v1.18.5", GitCommit:"e6503f8d8f769ace2f338794c914a96fc335df0f", GitTreeState:"clean", BuildDate:"2020-06-26T03:39:24Z", GoVersion:"go1.13.9", Compiler:"gc", Platform:"linux/amd64"}

	majorVersion, err := strconv.Atoi(version.Major)
	if err != nil {
		panic(err)
	}
	minorVersion, err := strconv.Atoi(version.Minor)
	if err != nil {
		panic(err)
	}

	// check the K8S version
	if majorVersion == 1 && minorVersion < 16 { // support 1.13 ...
		Scale := &v1beta1.Scale{Spec: v1beta1.ScaleSpec{Replicas: number}} // generate a scale
		Scale.Name = chartName
		Scale.Namespace = namespace

		scale, err := client.ExtensionsV1beta1().Deployments(namespace).UpdateScale(chartName, Scale)
		if err != nil {
			return nil, err
		}
		return &scale.Spec.Replicas, nil
	} else if majorVersion == 1 && minorVersion > 16 { // support 1.16 1.17 1.18 ...
		Scale := &autoscalingv1.Scale{Spec: autoscalingv1.ScaleSpec{Replicas: number}} // generate a scale
		Scale.Name = chartName
		Scale.Namespace = namespace

		scale, err := client.AppsV1().Deployments(namespace).UpdateScale(chartName, Scale)
		if err != nil {
			return nil, err
		}
		return &scale.Spec.Replicas, nil
	} else {
		err := utils.NewRockError(404, 40400013, fmt.Sprintf("Not found this K8S version %v", version.GitVersion))
		return nil, err
	}
}
