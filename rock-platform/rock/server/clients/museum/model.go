package museum

//{"access-control-process":[
//  {
//    "name":"access-control-process",
//    "version":"v2.3.0-master-wjw-limit-ys-2331fd4",
//    "description":"A Helm chart for Kubernetes",
//    "apiVersion":"v1",
//    "appVersion":"1.0.0",
//    "urls":["charts/access-control-process-v2.3.0-master-wjw-limit-ys-2331fd4.tgz"],
//    "created":"2021-01-14T07:21:06.693437843Z",
//    "digest":"d6bfb660c186c0f3a85aecf5d890e863b61f6b68908a95d5f339b0bf186b54b2"
//  },
// ...
// ],
// {"mysql-operator":[...],
// ...
// }

type ChartVersion struct {
	Name        string       `json:"name"`
	Version     string       `json:"version"`
	Description string       `json:"description"`
	Keywords    []string     `json:"keywords"`
	Maintainers []Maintainer `json:"maintainers"`
	ApiVersion  string       `json:"apiVersion"`
	AppVersion  string       `json:"appVersion"`
	Urls        []string     `json:"urls"`
	Created     string       `json:"created"`
	Digest      string       `json:"digest"`
}

type Maintainer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ChartMapper map[string][]*ChartVersion

type ChartDetail struct {
	Name     string          `json:"name"`
	Versions []*ChartVersion `json:"version"`
}
