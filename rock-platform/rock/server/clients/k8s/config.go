package k8s

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	engineConfigmap = "license-config"
)

var engineNamespace = "engine"

// update license-config configmap
func UpdateConfigmapWithLicense(k8sConf, mode, data string) error {
	client, err := GetK8sClient(k8sConf)
	if err != nil {
		return err
	}

	if mode == "dongle" {
		engineNamespace = "nebula"
	}

	configmapCli := client.CoreV1().ConfigMaps(engineNamespace)
	licenseCM, err := configmapCli.Get(engineConfigmap, metav1.GetOptions{})

	licenseData := map[string]string{
		"client_license.lic": data,
	}

	if err != nil {
		if errors.IsNotFound(err) { // if not license-config configmap, create it
			fmt.Println(err.Error())
			cm := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name: engineConfigmap,
				},
				Data: licenseData,
			}
			_, err = configmapCli.Create(cm)
		}
		return err
	}

	//updateCM := &corev1.ConfigMap{
	//	Data:licenseData,
	//}

	licenseCM.Data = licenseData
	_, err = configmapCli.Update(licenseCM)
	return err
}
