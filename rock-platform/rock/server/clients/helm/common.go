package helm

// helm client

import (
	"go.rock.com/rock-platform/rock/server/database/api"
	"go.rock.com/rock-platform/rock/server/utils"
	"k8s.io/helm/pkg/helm"
)

// get helm client by cluster id
func getHelmClient(clusterId int64) (*helm.Client, error) {
	cluster, err := api.GetClusterById(clusterId) // get cluster by id
	if err != nil {
		return nil, err
	}

	host, err := utils.GenTillerAddr(cluster.Config) // get tiller addr by k8s admin.conf
	if err != nil {
		return nil, err
	}

	// https://juejin.im/post/6844903527882424327 ，此链接有介绍helm client固定写法
	// 使用helm模块，传入k8s集群的IP，生成options。然后将options传入到helm.NewClient中，生成helm客户端(固定写法)
	options := []helm.Option{helm.Host(host), helm.ConnectTimeout(5)}
	return helm.NewClient(options...), nil
}
