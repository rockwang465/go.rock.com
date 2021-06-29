package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.rock.com/rock-platform/rock/server/clients/k8s"
	"go.rock.com/rock-platform/rock/server/database/api"
	"go.rock.com/rock-platform/rock/server/utils"
	"k8s.io/api/core/v1"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type NodeLabel struct {
	Key   string `json:"key" example:"beta.kubernetes.io/os"`
	Value string `json:"value" example:"linux"`
}

type NodeAnnotation struct {
	Key   string `json:"key" example:"flannel.alpha.coreos.com/backend-type"`
	Value string `json:"value" example:"host-gw"`
}

type ClusterNodeResp struct {
	Name                    string            `json:"name" example:"kubernetes-master1"`
	UID                     string            `json:"uid" example:"3550d3f1-51b4-41e7-ba65-83d029f31e2b"`
	Labels                  []*NodeLabel      `json:"labels"`
	Annotations             []*NodeAnnotation `json:"annotations"`
	PodCIDR                 string            `json:"pod_cidr" example:"10.244.0.0/24"`
	Unschedulable           bool              `json:"unschedulable" example:"false"`
	KernelVersion           string            `json:"kernel_version" example:"4.18.0-193.6.3.el8_2.x86_64"`
	OSImage                 string            `json:"os_image" example:"CentOS Linux 8 (Core)"`
	OS                      string            `json:"os" example:"linux"`
	Architecture            string            `json:"architecture" example:"amd64"`
	ContainerRunTimeVersion string            `json:"container_run_time_version" example:"docker://19.3.4"`
	InternalIP              string            `json:"internal_ip" example:"10.10.10.10"`
	Hostname                string            `json:"hostname" example:"kubernetes-master1"`
	CreatedAt               time.Time         `json:"created_at" example:"2021-02-13T18:12:05+08:00"`
}

type NodeReq struct {
	Id   int64  `json:"id" uri:"id" binding:"required,min=1" example:"1"`
	Name string `json:"name" uri:"name" binding:"required" example:"kubernetes-master1"`
}

type GlobalNodeResp struct {
	Name                    string            `json:"name" example:"kubernetes-master1"`
	ClusterName             string            `json:"cluster_name,omitempty" binding:"required" example:"devops"`
	ClusterId               int               `json:"cluster_id,omitempty" binding:"required" example:"1"`
	UID                     string            `json:"uid" example:"3550d3f1-51b4-41e7-ba65-83d029f31e2b"`
	Labels                  []*NodeLabel      `json:"labels"`
	Annotations             []*NodeAnnotation `json:"annotations"`
	PodCIDR                 string            `json:"pod_cidr" example:"10.244.0.0/24"`
	Unschedulable           bool              `json:"unschedulable" example:"false"`
	KernelVersion           string            `json:"kernel_version" example:"4.18.0-193.6.3.el8_2.x86_64"`
	OSImage                 string            `json:"os_image" example:"CentOS Linux 8 (Core)"`
	OS                      string            `json:"os" example:"linux"`
	Architecture            string            `json:"architecture" example:"amd64"`
	ContainerRunTimeVersion string            `json:"container_run_time_version" example:"docker://19.3.4"`
	InternalIP              string            `json:"internal_ip" example:"10.10.10.10"`
	Hostname                string            `json:"hostname" example:"kubernetes-master1"`
	CreatedAt               time.Time         `json:"created_at" example:"2021-02-13T18:12:05+08:00"`
}

// @Summary Get specific cluster's all nodes
// @Description api for get specific cluster's all nodes
// @Tags CLUSTER
// @Accept json
// @Produce json
// @Param id path integer true "Cluster ID"
// @Success 200 {array} v1.ClusterNodeResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/clusters/{id}/nodes [get]
func (c *Controller) GetClusterNodes(ctx *gin.Context) {
	var idReq IdReq // cluster id
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		panic(err)
	}

	cluster, err := api.GetClusterById(idReq.Id)
	if err != nil {
		panic(err)
	}

	nodeList, err := k8s.GetClusterNodes(cluster.Config)
	if err != nil {
		panic(err)
	}
	nodes, err := formatNodesResp(nodeList.Items)
	if err != nil {
		panic(err)
	}

	c.Logger.Infof("Get specific cluster all nodes by cluster id:%v", idReq.Id)
	ctx.JSON(http.StatusOK, nodes)
}

func formatNodesResp(nodeList []v1.Node) (*[]ClusterNodeResp, error) {
	nodesResp := []ClusterNodeResp{}
	for _, node := range nodeList {
		nodeResp, err := formatNodeResp(&node)
		if err != nil {
			return nil, err
		}

		nodesResp = append(nodesResp, *nodeResp)
	}
	return &nodesResp, nil
}

func formatNodeResp(node *v1.Node) (*ClusterNodeResp, error) {
	Node := &ClusterNodeResp{
		Name: node.Name,
		UID:  string(node.UID),
		//Labels:                  node.ObjectMeta.Labels,
		//Annotations:             node.ObjectMeta.Annotations,
		PodCIDR:                 node.Spec.PodCIDR,
		Unschedulable:           node.Spec.Unschedulable,
		KernelVersion:           node.Status.NodeInfo.KernelVersion,
		OSImage:                 node.Status.NodeInfo.OSImage,
		OS:                      node.Status.NodeInfo.OperatingSystem,
		Architecture:            node.Status.NodeInfo.Architecture,
		ContainerRunTimeVersion: node.Status.NodeInfo.ContainerRuntimeVersion,
		CreatedAt:               node.CreationTimestamp.Time,
	}

	for _, nodeAddress := range node.Status.Addresses {
		if nodeAddress.Type == "InternalIP" {
			Node.InternalIP = nodeAddress.Address
		}
		if nodeAddress.Type == "Hostname" {
			Node.Hostname = nodeAddress.Address
		}
	}

	labels := []*NodeLabel{} // 指针存储,节省内存,但必须先初始化
	for key, value := range node.ObjectMeta.Labels {
		label := NodeLabel{
			Key:   key,
			Value: value,
		}
		labels = append(labels, &label)
	}
	Node.Labels = labels

	annotations := []*NodeAnnotation{}
	for key, value := range node.ObjectMeta.Annotations {
		annotation := NodeAnnotation{
			Key:   key,
			Value: value,
		}
		annotations = append(annotations, &annotation)
	}
	Node.Annotations = annotations

	return Node, nil

	//上面 formatNodesList 函数中range nodeList中的node信息如下:
	//v1.Node{
	//  TypeMeta:v1.TypeMeta{Kind:"", APIVersion:""},
	//  ObjectMeta:v1.ObjectMeta{
	//    Name:"kubernetes-master",
	//    GenerateName:"", Namespace:"",
	//    SelfLink:"/api/v1/nodes/kubernetes-master",
	//    UID:"3f9ea9ac-5940-11eb-a0db-ac1f6b472cc8",
	//    ResourceVersion:"6048413", Generation:0,
	//    CreationTimestamp:v1.Time{
	//      Time:time.Time{wall:0x0, ext:63746538604, loc:(*time.Location)(0x35011e0)}
	//    },
	//    DeletionTimestamp:(*v1.Time)(nil),
	//    DeletionGracePeriodSeconds:(*int64)(nil),
	//    Labels:map[string]string{
	//      "beta.kubernetes.io/arch":"amd64",
	//      "beta.kubernetes.io/os":"linux",
	//      "kubernetes.io/hostname":"kubernetes-master",
	//      "node-role.kubernetes.io/master":""
	//    },
	//    Annotations:map[string]string{
	//      "flannel.alpha.coreos.com/backend-data":"null",
	//      "flannel.alpha.coreos.com/backend-type":"host-gw",
	//      "flannel.alpha.coreos.com/kube-subnet-manager":"true",
	//      "flannel.alpha.coreos.com/public-ip":"10.151.3.96",
	//      "kubeadm.alpha.kubernetes.io/cri-socket":"/var/run/dockershim.sock",
	//      "node.alpha.kubernetes.io/ttl":"0",
	//      "volumes.kubernetes.io/controller-managed-attach-detach":"true"
	//    },
	//    OwnerReferences:[]v1.OwnerReference(nil),
	//    Initializers:(*v1.Initializers)(nil),
	//    Finalizers:[]string(nil), ClusterName:""
	//  },
	//  Spec:v1.NodeSpec{
	//    PodCIDR:"10.244.0.0/24",
	//    ProviderID:"",
	//    Unschedulable:false,
	//    Taints:[]v1.Taint(nil),
	//    ConfigSource:(*v1.NodeConfigSource)(nil),
	//    DoNotUse_ExternalID:"",
	//  },
	//  Status:v1.NodeStatus{
	//    Capacity:v1.ResourceList{
	//      "cpu":resource.Quantity{i:resource.int64Amount{value:64, scale:0}, d:resource.infDecAmount{Dec:(*inf.Dec)(nil)}, s:"64", Format:"DecimalSI"},
	//      "ephemeral-storage":resource.Quantity{i:resource.int64Amount{value:238656749568, scale:0}, d:resource.infDecAmount{Dec:(*inf.Dec)(nil)}, s:"", Format:"BinarySI"},
	//      "hugepages-1Gi":resource.Quantity{i:resource.int64Amount{value:0, scale:0}, d:resource.infDecAmount{Dec:(*inf.Dec)(nil)}, s:"0", Format:"DecimalSI"},
	//      "hugepages-2Mi":resource.Quantity{i:resource.int64Amount{value:0, scale:0}, d:resource.infDecAmount{Dec:(*inf.Dec)(nil)}, s:"0", Format:"DecimalSI"},
	//      "memory":resource.Quantity{i:resource.int64Amount{value:201234952192, scale:0}, d:resource.infDecAmount{Dec:(*inf.Dec)(nil)}, s:"196518508Ki", Format:"BinarySI"},
	//      "nvidia.com/gpu":resource.Quantity{i:resource.int64Amount{value:4, scale:0}, d:resource.infDecAmount{Dec:(*inf.Dec)(nil)}, s:"4", Format:"DecimalSI"},
	//      "pods":resource.Quantity{i:resource.int64Amount{value:110, scale:0}, d:resource.infDecAmount{Dec:(*inf.Dec)(nil)}, s:"110", Format:"DecimalSI"}
	//	},
	//    Allocatable:v1.ResourceList{
	//      "cpu":resource.Quantity{i:resource.int64Amount{value:64, scale:0}, d:resource.infDecAmount{Dec:(*inf.Dec)(nil)}, s:"64", Format:"DecimalSI"},
	//      "ephemeral-storage":resource.Quantity{i:resource.int64Amount{value:238656749568, scale:0}, d:resource.infDecAmount{Dec:(*inf.Dec)(nil)}, s:"", Format:"BinarySI"},
	//      "hugepages-1Gi":resource.Quantity{i:resource.int64Amount{value:0, scale:0}, d:resource.infDecAmount{Dec:(*inf.Dec)(nil)}, s:"0", Format:"DecimalSI"},
	//      "hugepages-2Mi":resource.Quantity{i:resource.int64Amount{value:0, scale:0}, d:resource.infDecAmount{Dec:(*inf.Dec)(nil)}, s:"0", Format:"DecimalSI"},
	//      "memory":resource.Quantity{i:resource.int64Amount{value:201234952192, scale:0}, d:resource.infDecAmount{Dec:(*inf.Dec)(nil)}, s:"196518508Ki", Format:"BinarySI"},
	//      "nvidia.com/gpu":resource.Quantity{i:resource.int64Amount{value:4, scale:0}, d:resource.infDecAmount{Dec:(*inf.Dec)(nil)}, s:"4", Format:"DecimalSI"},
	//      "pods":resource.Quantity{i:resource.int64Amount{value:110, scale:0}, d:resource.infDecAmount{Dec:(*inf.Dec)(nil)}, s:"110", Format:"DecimalSI"}
	//	},
	//    Phase:"",
	//    Conditions:[]v1.NodeCondition{
	//       v1.NodeCondition{Type:"MemoryPressure", Status:"False", LastHeartbeatTime:v1.Time{Time:time.Time{wall:0x0, ext:63749227334, loc:(*time.Location)(0x35011e0)}}, LastTransitionTime:v1.Time{Time:time.Time{wall:0x0, ext:63746538600, loc:(*time.Location)(0x35011e0)}}, Reason:"KubeletHasSufficientMemory", Message:"kubelet has sufficient memory available"},
	//       v1.NodeCondition{Type:"DiskPressure", Status:"False", LastHeartbeatTime:v1.Time{Time:time.Time{wall:0x0, ext:63749227334, loc:(*time.Location)(0x35011e0)}}, LastTransitionTime:v1.Time{Time:time.Time{wall:0x0, ext:63746538600, loc:(*time.Location)(0x35011e0)}}, Reason:"KubeletHasNoDiskPressure", Message:"kubelet has no disk pressure"},
	//       v1.NodeCondition{Type:"PIDPressure", Status:"False", LastHeartbeatTime:v1.Time{Time:time.Time{wall:0x0, ext:63749227334, loc:(*time.Location)(0x35011e0)}}, LastTransitionTime:v1.Time{Time:time.Time{wall:0x0, ext:63746538600, loc:(*time.Location)(0x35011e0)}}, Reason:"KubeletHasSufficientPID", Message:"kubelet has sufficient PID available"},
	//       v1.NodeCondition{Type:"Ready", Status:"True", LastHeartbeatTime:v1.Time{Time:time.Time{wall:0x0, ext:63749227334, loc:(*time.Location)(0x35011e0)}}, LastTransitionTime:v1.Time{Time:time.Time{wall:0x0, ext:63746539626, loc:(*time.Location)(0x35011e0)}}, Reason:"KubeletReady", Message:"kubelet is posting ready status"}
	//    },
	//    Addresses:[]v1.NodeAddress{
	//      v1.NodeAddress{Type:"InternalIP", Address:"10.151.3.96"},
	//      v1.NodeAddress{Type:"Hostname", Address:"kubernetes-master"}
	//    },
	//    DaemonEndpoints:v1.NodeDaemonEndpoints{KubeletEndpoint:v1.DaemonEndpoint{Port:10250}},
	//    NodeInfo:v1.NodeSystemInfo{
	//      MachineID:"4c618dcc6f4f4c5fb56aa1a3130b00b4",
	//      SystemUUID:"00000000-0000-0000-0000-AC1F6B472CC8",
	//      BootID:"fc2a6d96-873b-499c-8636-f69521471198",
	//      KernelVersion:"3.10.0-693.5.2.el7.x86_64",
	//      OSImage:"CentOS Linux 7 (Core)",
	//      ContainerRuntimeVersion:"docker://18.6.1",
	//      KubeletVersion:"v1.13.2",
	//      KubeProxyVersion:"v1.13.2",
	//      OperatingSystem:"linux",
	//      Architecture:"amd64"
	//    },
	//    Images:[]v1.ContainerImage{
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/nebula-test/engine-video-process-service@sha256:1eb6d0fac4791cfe79a0694ec733bb8afe3881b77e1c5dd7d30c0176d1299917", "registry.sensenebula.io:5000/nebula-test/engine-video-process-service:v2.3.0-master-kestrelV15-cuda11-1442037f"}, SizeBytes:6350696801},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/nebula-test/engine-image-process-service@sha256:436c445de5b938bed709e6b1af8fa2fa16e00252e6e54b1346b29ca6dd966f60", "registry.sensenebula.io:5000/nebula-test/engine-image-process-service:v2.3.0-master-1c21a89-cuda11-t4"}, SizeBytes:4584961403},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/nebula-test/gpu-searcher@sha256:c5baad54c434bb83d5e28efe9846162d6552841d651a38ef62d6dda8b9eaa63d", "registry.sensenebula.io:5000/nebula-test/gpu-searcher:v2.3.0-master-1535a13"}, SizeBytes:3823476870},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/nebula-test/engine-static-feature-db@sha256:8011f826987cc1fa5100e0e7d1394559f850b033349d8faa55b7887925f27b55", "registry.sensenebula.io:5000/nebula-test/engine-static-feature-db:v2.3.0-master-ff39b54"}, SizeBytes:3693668462},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/gitlabci/golang@sha256:86a3ba51f60fcb7d4b33184fb5263c7113a6fa2bb6146c15d20a33341c8a1536", "registry.sensenebula.io:5000/gitlabci/golang:1.9-cuda-gcc49-1"}, SizeBytes:3287755137},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/component/cassandra@sha256:a1af2d8f5e5ac81a724e97a03eac764dc2d1d212aaba29bd0dae41e7903e65b7", "registry.sensenebula.io:5000/component/cassandra:3.11.4"}, SizeBytes:992466356},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/component/kafka@sha256:f17fed4bfb5c0bf2fbb2c3763665ec339d45a49d8272ec0b16693b6e69638227", "registry.sensenebula.io:5000/component/kafka:2.11-1.1.1"}, SizeBytes:956535508}, v1.ContainerImage{Names:[]string{"senselink/database_service:v1.8.0-p"}, SizeBytes:948674937}, v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/elasticsearch/elasticsearch-oss@sha256:e40d418547ce10bbf2c0e8524f71a9f465e11967aca28213edec281bc8ea0fd8", "registry.sensenebula.io:5000/elasticsearch/elasticsearch-oss:6.5.4"}, SizeBytes:644764369},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/nebula-test/engine-image-ingress-service@sha256:c576fd2513c5c6eed600c987a3c09e5eb5ad5d3fec829ed4ddd3c2f96b02fa7b", "registry.sensenebula.io:5000/nebula-test/engine-image-ingress-service:v2.3.0-master-335d3e83"}, SizeBytes:583792459},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/mysql/mysql-agent@sha256:29ebfa3c0790823b5132e26ab93e5c8b11ed99f7cc455c8aa299a11666bb90d6", "registry.sensenebula.io:5000/mysql/mysql-agent:356674d-2020122217"}, SizeBytes:448177276},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/solsson/burrow@sha256:a14a6911386f4523a249f2942bce6053f417b63dec6779f65a5a261aecbc5397", "registry.sensenebula.io:5000/solsson/burrow:v1.0.0"}, SizeBytes:440512053},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/mysql/mysql-router@sha256:bfe13ae7258c8d63601437bd46cc3e3bc871642cb50d8f49bcfc76ac48dce3b3", "registry.sensenebula.io:5000/mysql/mysql-router:8.0.22-6c5292c"}, SizeBytes:433571043},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/component/zookeeper@sha256:94d809d938c5ff8cca5fc3e08514ace7130861ec311ab04ce16d912fb787b1b6", "registry.sensenebula.io:5000/component/zookeeper:1.0-3.4.14"}, SizeBytes:407887087},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/mysql/mysql-server@sha256:5b2e6db9829c81653d14799ccf2bdf143954edef374b6a5119f816f2e1fd4bec", "registry.sensenebula.io:5000/mysql/mysql-server:8.0.22"}, SizeBytes:405009901},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-log@sha256:4ae93b77405ddf8bdd0485c6c3ad0651ed6bf576147e04cb3c4e33b749c9e18c", "registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-log:2.3.0-2.3.0-001-dev-eb650a"}, SizeBytes:367618701},
	//      v1.ContainerImage{Names:[]string{"senselink/java_openapi_service:v2.3.0-p"}, SizeBytes:325448734},
	//      v1.ContainerImage{Names:[]string{"senselink/nebula_service:v2.3.0-p"}, SizeBytes:325448734},
	//      v1.ContainerImage{Names:[]string{"senselink/event_service:v2.3.0-p"}, SizeBytes:325448734},
	//      v1.ContainerImage{Names:[]string{"senselink/tk_service:v2.3.0-p"}, SizeBytes:325448734},
	//      v1.ContainerImage{Names:[]string{"senselink/java_websocket_service:v2.3.0-p"}, SizeBytes:325448686},
	//      v1.ContainerImage{Names:[]string{"senselink/java_inner_service:v2.3.0-p"}, SizeBytes:325448686},
	//      v1.ContainerImage{Names:[]string{"senselink/java_service:v2.3.1-p"}, SizeBytes:325447822},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/kubernetes/nginx-ingress-controller@sha256:3872438389bda2d8f878187527d68b86dbdfb3c73ac67651186f5460e01c9073", "registry.sensenebula.io:5000/kubernetes/nginx-ingress-controller:0.30.0"}, SizeBytes:322915865},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-frontend-comparison@sha256:20bf0bdee7d0061762f0135d063790ac0a29ae6a2bf5ac48a9c64767ea27e42f", "registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-frontend-comparison:2.3.0-2.3.0-001-dev-bab4a5"}, SizeBytes:290748261},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-integrated-adapter@sha256:7076b680f1f6cf14fac5bc2530215808d90ce6d71ac5a3a273b5875bda119f5d", "registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-integrated-adapter:2.3.0-2.3.0-001-dev-adc50a"}, SizeBytes:289918108},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-records-export@sha256:ac044c9aa393fe5b7492ebdafc4ce81b4e22199367979be329009c72c8f41563", "registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-records-export:2.3.0-2.3.0-dev-c0ed51"}, SizeBytes:286284177},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-search-face@sha256:1ae32413a1470ef81f2edaa8c95984020e92e7adfb47b5f13ec63754a06b27c5", "registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-search-face:2.3.0-2.3.0-001-dev-42d666"}, SizeBytes:286133352},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-guest-management@sha256:50198c53f48ce21c0d5f110e3dcb162780df7d97424e9c9f877435e30a5c9cad", "registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-guest-management:2.3.0-2.3.0-001-dev-b70625"}, SizeBytes:286106866},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-message-center@sha256:f94f506c26c81ec017eb08cf0214e06f948553272acf45d8b58def7eed49bf3d", "registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-message-center:2.3.0-2.3.0-001-dev-0a48c2"}, SizeBytes:286050391},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-attendance-management@sha256:52698b3d6ef849641235c6032e64c1354149729dcdd83f92f21807b26ee0f569", "registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-attendance-management:2.3.0-2.3.0-001-dev-c1d3b7"}, SizeBytes:285838577},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-lbs@sha256:4ab692dcbb9af2c90ff906a825348544d5e1189e9e558dbb380499e3c2afa7a1", "registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-lbs:2.3.0-2.3.0-dev-efbdbb"}, SizeBytes:285651208},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-records-management@sha256:11dc0e3caca2cef75aaf682ee8661e1635579fb50809d7fd294f54abcc909270", "registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-records-management:2.3.0-2.3.0-001-dev-e85dab"}, SizeBytes:285651199},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-struct-process-service@sha256:0c329bc5afb7b84cf9369fb6215a09401b7542869bb21e134e273419dd0e1593", "registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-struct-process-service:2.3.0-2.3.0-001-dev-2bcd01"}, SizeBytes:284876028},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-td-result-consume@sha256:cdd5b5bbeac5ba2a6cfd1e149e87f8a423fbdf96a336076b6d8cd116a8db1281", "registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-td-result-consume:2.3.0-2.3.0-001-dev-0fda31"}, SizeBytes:283992189},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-ac-result-consume@sha256:475d74813bfff8371436c40e57344dc82c172fbcb56f8dac1fc5500018f4b876", "registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-ac-result-consume:2.3.0-2.3.0-001-dev-9aba45"}, SizeBytes:283838596},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-lib-auth@sha256:7be05965f38827411a913aee4e15b54df89f88060bed6eee6d5c544755629d8d", "registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-lib-auth:2.3.0-2.3.0-001-dev-069361"}, SizeBytes:283327467},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-watchlist-management@sha256:489d83cfeb14a5d7f2083445f52220b36ac82b657cc31c675b5fec32a4076d43", "registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-watchlist-management:2.3.0-2.3.0-001-dev-251d57"}, SizeBytes:283118293},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-bulk-tool@sha256:420832da90a3a912bcdcc785e657374e30c46838a31a1af1b8b67afe98303f2a", "registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-bulk-tool:2.3.0-2.3.0-001-dev-51d6f8"}, SizeBytes:282984735},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-rule-management@sha256:7acd848b83e55e69d492953dc27846d774e7c8a54f9db09de405e3db56b4084e", "registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-rule-management:2.3.0-2.3.0-001-dev-57cb77"}, SizeBytes:275636655},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-target-export@sha256:142b8852a8e4db366545e6c0cf683b7456ae7e34b017b5d0b5d7cfccddd2170e", "registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-target-export:2.3.0-2.3.0-001-dev-9a6ee5"}, SizeBytes:275594613},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-device-management@sha256:734d633cde53da5a35f9399d86999eb59d9d9e90e1c02c992e056dfb0b60daf7", "registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-device-management:2.3.0-2.3.0-001-dev-b029f7"}, SizeBytes:275557587},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-tools@sha256:dd4ab740601cab7849a211cdfdd05b7ce52cb1bae911eadabe439633ae7eac35", "registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-tools:2.3.0-2.3.0-001-dev-de7021"}, SizeBytes:275065126},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-uums@sha256:94ccaa2be557c6b0cbb8b17d48e96d22f2eaf77a9d651575f0f0fb5823ee9ca7", "registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-uums:2.3.0-2.3.0-001-dev-75df81"}, SizeBytes:274623581},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-oauth2@sha256:9800eefd292fbbffa41f5d3f9a5c0db6931570d330a35d63e78986504b533cd7", "registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-oauth2:2.3.0-2.3.0-001-dev-06ecd5"}, SizeBytes:273501072},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-scheduled@sha256:13dd73926988802a59e73fbe7594b5ecf4d52f67fd8b819f75512882e5e50e38", "registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-scheduled:2.3.0-2.3.0-001-dev-ade00b"}, SizeBytes:273352288},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-dashboard@sha256:4f9f46c210497bd5502eb26f2952356cfc8a0e65d3786ac2c3bdddfd9d0d76dc", "registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-dashboard:2.3.0-2.3.0-001-dev-6f7f65"}, SizeBytes:273275775},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-timezone-management@sha256:74e14ca67211b017920505d5338c6a6608cbdfbbf22c2d54fc7ebbe8eb8ad304", "registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-timezone-management:2.3.0-2.3.0-001-dev-29e59e"}, SizeBytes:272915336},
	//      v1.ContainerImage{Names:[]string{"registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-map-management@sha256:e682411503742b7ff9670b9e023ffaeb459783434afddb1df0c216305f0c2274", "registry.sensenebula.io:5000/sensenebula-guard-std/senseguard-map-management:2.3.0-2.3.0-001-dev-beb8ab"}, SizeBytes:272905821},
	//      v1.ContainerImage{Names:[]string{"senselink/opq_service:v1.8.0-p"}, SizeBytes:268557806}
	//    },
	//    VolumesInUse:[]v1.UniqueVolumeName(nil), VolumesAttached:[]v1.AttachedVolume(nil), Config:(*v1.NodeConfigStatus)(nil)}
	//}
}

func formatGlobalNodesResp(nodes []v1.Node) ([]*GlobalNodeResp, error) {
	clusterNodes := make([]*GlobalNodeResp, 0)
	for _, node := range nodes {
		clusterNode, err := formatGlobalNodeResp(&node)
		if err != nil {
			return nil, err
		}
		clusterNodes = append(clusterNodes, clusterNode)
	}
	return clusterNodes, nil
}

func formatGlobalNodeResp(node *v1.Node) (*GlobalNodeResp, error) {
	Node := &GlobalNodeResp{
		Name:                    node.Name,
		UID:                     string(node.UID),
		PodCIDR:                 node.Spec.PodCIDR,
		Unschedulable:           node.Spec.Unschedulable,
		KernelVersion:           node.Status.NodeInfo.KernelVersion,
		OSImage:                 node.Status.NodeInfo.OSImage,
		OS:                      node.Status.NodeInfo.OperatingSystem,
		Architecture:            node.Status.NodeInfo.Architecture,
		ContainerRunTimeVersion: node.Status.NodeInfo.ContainerRuntimeVersion,
		CreatedAt:               node.CreationTimestamp.Time,
	}

	var clusterId int
	var err error
	clusterId, err = strconv.Atoi(node.Annotations["console.cluster.id"])
	if err != nil {
		return nil, utils.NewRockError(400, 40000023, fmt.Sprintf("cluster id %s can't be converted int", node.Annotations["console.cluster.id"]))
	}

	Node.ClusterName = node.Annotations["console.cluster.name"] // 将之前保存的cluster.Name和cluster.Id值赋值到这个结构体中
	Node.ClusterId = clusterId

	for _, nodeAddress := range node.Status.Addresses {
		if nodeAddress.Type == "InternalIP" {
			Node.InternalIP = nodeAddress.Address
		}
		if nodeAddress.Type == "Hostname" {
			Node.Hostname = nodeAddress.Address
		}
	}

	labels := []*NodeLabel{} // 指针存储,节省内存,但必须先初始化
	for key, value := range node.Labels {
		label := NodeLabel{
			Key:   key,
			Value: value,
		}
		labels = append(labels, &label)
	}
	Node.Labels = labels

	annotations := []*NodeAnnotation{}
	for key, value := range node.ObjectMeta.Annotations { // 如果是console.cluster开头,为之前记录cluster.Name和cluster.Id的,所以不能保存到数据中,则需要continue忽略掉
		if strings.HasPrefix(key, "console.cluster") {
			continue
		}
		annotation := NodeAnnotation{
			Key:   key,
			Value: value,
		}
		annotations = append(annotations, &annotation)
	}
	Node.Annotations = annotations

	return Node, nil
}

// @Summary Get a specific cluster node
// @Description api for get a specific cluster node
// @Tags CLUSTER
// @Accept json
// @Produce json
// @Param id path integer true "Cluster ID"
// @Param name path string true "Node name"
// @Success 200 {object} v1.ClusterNodeResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/clusters/{id}/nodes/{name} [get]
func (c *Controller) GetClusterNode(ctx *gin.Context) {
	var nodeReq NodeReq // cluster id + kubernetes node name
	if err := ctx.ShouldBindUri(&nodeReq); err != nil {
		panic(err)
	}

	cluster, err := api.GetClusterById(nodeReq.Id)
	if err != nil {
		panic(err)
	}

	node, err := k8s.GetClusterNode(cluster.Config, nodeReq.Name)
	if err != nil {
		panic(err)
	}
	Node, err := formatNodeResp(node)
	if err != nil {
		panic(err)
	}

	c.Logger.Infof("Get specific cluster node by cluster id(%v) and node name(%v)", nodeReq.Id, node.Name)
	ctx.JSON(http.StatusOK, Node)
}

// @Summary Get cluster's all nodes
// @Description api for get cluster's all nodes
// @Tags NODE
// @Accept json
// @Produce json
// @Success 200 {object} v1.GlobalNodeResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/nodes [get]
func (c *Controller) GetGlobalNodes(ctx *gin.Context) {
	clusters, err := api.GetClustersWithoutPagination()
	if err != nil {
		panic(err)
	}

	nodes := []v1.Node{}
	for _, cluster := range clusters {
		nodeList, err := k8s.GetClusterNodes(cluster.Config) // 通过admin.conf获取单个集群的节点信息
		if err != nil {
			c.Logger.Warnf("Get cluster(%v)'s node failed, please check it", cluster.Name)
			continue
		}

		for _, node := range nodeList.Items { // 保存cluster.Name 和cluster.Id,方便后面单个node进行保存数据
			node.Annotations["console.cluster.name"] = cluster.Name                // cluster.Name为集群名称，如 10.151.3.99-devops-env
			node.Annotations["console.cluster.id"] = strconv.Itoa(int(cluster.Id)) // cluster.Id为集群id
		}
		nodes = append(nodes, nodeList.Items...)
	}

	resp, err := formatGlobalNodesResp(nodes)
	if err != nil {
		panic(err)
	}
	c.Logger.Infof("Get all nodes, the nodes length is %v", len(resp))
	ctx.JSON(http.StatusOK, resp)
}
