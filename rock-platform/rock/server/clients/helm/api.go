package helm

import (
	"fmt"
	"go.rock.com/rock-platform/rock/server/utils"
	"k8s.io/helm/pkg/helm"
)

// check cluster tiller healthy by cluster id
func PingTillerServer(clusterId int64) error {
	client, err := getHelmClient(clusterId)
	if err != nil {
		return err
	}

	if err := client.PingTiller(); err != nil {
		e := utils.NewRockError(500, 50000005, fmt.Sprintf("Tiller server of cluster %v isn't healthy, please check it first", clusterId))
		return e
	}
	return nil
}

// remove helm service in cluster when FAILED or DELETED status
func DeleteReleaseIfFailedOrDeleted(clusterId int64, releaseName string) error {
	client, err := getHelmClient(clusterId)
	if err != nil {
		return err
	}
	status, _ := client.ReleaseStatus(releaseName)
	if status == nil {
		return nil
	}

	statusStr := status.GetInfo().GetStatus().GetCode().String()
	if statusStr == "FAILED" || statusStr == "DELETED" {
		_, err := client.DeleteRelease(releaseName)
		if err != nil {
			return err
		}
	}
	return nil
}

// remove helm service in cluster when DEPLOYED status
func DeleteManualInstallReleaseIfExist(clusterId int64, releaseName string) error {
	client, err := getHelmClient(clusterId)
	if err != nil {
		return err
	}
	status, _ := client.ReleaseStatus(releaseName) // status equal nil, when release in not found
	if status == nil {
		return nil
	}

	statusStr := status.GetInfo().GetStatus().GetCode().String()
	if statusStr == "DEPLOYED" { // status equal DEPLOYED
		_, err := client.DeleteRelease(releaseName, helm.DeletePurge(true))
		if err != nil {
			return err
		}
	}
	return nil
}
