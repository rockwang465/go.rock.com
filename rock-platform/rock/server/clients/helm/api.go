package helm

import (
	"fmt"
	"go.rock.com/rock-platform/rock/server/log"
	"go.rock.com/rock-platform/rock/server/utils"
	"k8s.io/helm/pkg/helm"
	rls "k8s.io/helm/pkg/proto/hapi/services"
	"strings"
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
	logger := log.GetLogger()
	logger.Debugf("helm delete release %v in cluster id %v", releaseName, clusterId)
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
	logger := log.GetLogger()
	logger.Debugf("helm delete release %v in cluster id %v", releaseName, clusterId)
	return nil
}

// remove helm service in cluster by cluster id and releaseName
func DeleteRelease(clusterId int64, releaseName string) (*rls.UninstallReleaseResponse, error) {
	client, err := getHelmClient(clusterId)
	if err != nil {
		return nil, err
	}

	resp, err := client.DeleteRelease(releaseName, helm.DeletePurge(true))
	if err != nil {
		str := err.Error()
		if strings.Contains(str, "not found") { // 当helm服务不存在了，但仍然调当前接口删除，则报错: rpc error: code = Unknown desc = release: "aurora-system-service-default" not found
			return resp, nil
		}
		return resp, err
	}
	logger := log.GetLogger()
	logger.Debugf("helm delete release %v in cluster id %v", releaseName, clusterId)
	return resp, nil
}
