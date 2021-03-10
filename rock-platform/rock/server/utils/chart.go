package utils

import (
	"bytes"
	"fmt"
	"go.rock.com/rock-platform/rock/server/conf"
	"go.rock.com/rock-platform/rock/server/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type K8sConfig struct {
	Clusters []Clusters `yaml:"clusters"`
}
type Clusters struct {
	Cluster Cluster `yaml:"cluster"`
	Name    string  `yaml:"name"`
}
type Cluster struct {
	CertificateAuthorityData string `yaml:"certificate-authority-data"`
	Server                   string `yaml:"server"`
}
type Context struct {
	Cluster string `yaml:"cluster"`
	User    string `yaml:"user"`
}

// get tiller address
func GenTillerAddr(config string) (string, error) {
	Conf := conf.GetConfig()
	tillerPort := Conf.Viper.GetString("tiller.port")
	clusterIp, err := GetClusterHostFromConfig(config) // get ip address in k8s admin.conf
	if err != nil {
		return "", err
	}
	tillerAddr := fmt.Sprintf("%s:%s", clusterIp, tillerPort) // example: 10.151.3.99:31134
	return tillerAddr, nil
}

// get ip address in k8s admin.conf
func GetClusterHostFromConfig(config string) (string, error) {
	var k8sConfig K8sConfig
	if err := yaml.Unmarshal([]byte(config), &k8sConfig); err != nil { // from string to struct
		e := NewRockError(400, 40000025, fmt.Sprintf("Cluster's config is malformed, message is: %s", err.Error()))
		return "", e
	}

	serverAddr := k8sConfig.Clusters[0].Cluster.Server // get admin.conf server domain (example: http:10.151.3.99:16443)

	Url, err := url.Parse(serverAddr) // use url module to parse domain
	if err != nil {
		e := NewRockError(400, 40000025, fmt.Sprintf("Error occured when parse k8s cluster config server address, message is: %s", err.Error()))
		return "", e
	}
	return Url.Hostname(), nil // Url.Hostname example: 10.151.3.99
}

// Install or upgrade the helm chart to the specified environment
func InstallOrUpgradeChart(repoUrl, chartTgzName, k8sConfig, ns, releaseName, appConfig string) error {
	logger := log.GetLogger()
	tillerAddr, err := GenTillerAddr(k8sConfig)
	if err != nil {
		panic(err)
	}

	binLocation := getBinFileLocation()

	tmpFile, err := ioutil.TempFile("", releaseName)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte(appConfig))
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	chartUrl := fmt.Sprintf("%s/charts/%s", repoUrl, chartTgzName)
	fileParams := fmt.Sprintf("-f %s", tmpFile.Name())
	// helm upgrade --install {releaseName} {chartUrl} {fileParams} --host {tillerAddr} --namespace {ns} --reset-values --force
	cmdStr := fmt.Sprintf("upgrade --install %s %s %s --host %s --namespace %s --reset-values --force", releaseName, chartUrl, fileParams, tillerAddr, ns)
	logger.Logger.Infof("%s %s", binLocation, cmdStr)

	cmdParams := strings.Split(cmdStr, " ")
	cmd := exec.Command(binLocation, cmdParams...)
	var stdout, stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return NewRockError(500, 50000006, fmt.Sprintf("Error occurred when helm install chart: %v", stderr.String()))
	}
	return nil
}

// generate chart name
func GenerateChartName(chartName, namespace string) string {
	return fmt.Sprintf("%s-%s", chartName, namespace) // example: kafka-component
}

// get the chartmuseum repo address by config.yaml
func GetRepoUrl() string {
	config := conf.GetConfig()
	return config.Viper.GetString("chartsRepo.addr") // example: http://10.151.3.75:8080
}

// get the chart tgz package name
func GenChartTgzName(chartName, chartVersion string) string {
	return fmt.Sprintf("%s-%s.tgz", chartName, chartVersion)
}

// get the helm binary location
func getBinFileLocation() string {
	BinName := ""
	switch runtime.GOOS {
	case "windows":
		BinName = "helm.exe"
		return filepath.Join("console", "tools", BinName)
	default:
		BinName = "helm"
	}
	return BinName
}
