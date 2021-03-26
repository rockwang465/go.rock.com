package license

// cactl.go file from license-ca team

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const defaultMasterCAUrl = "https://private.ca.sensetime.com:8443"
const defaultSlaveCAUrl = "https://slave.private.ca.sensetime.com:8443"

const defaultMasterCACert = `
-----BEGIN CERTIFICATE-----
MIIDozCCAougAwIBAgIJAJgOlpYuWlSUMA0GCSqGSIb3DQEBCwUAMGgxCzAJBgNV
BAYTAkNOMRAwDgYDVQQIDAdCRUlKSU5HMRAwDgYDVQQHDAdCRUlKSU5HMRIwEAYD
VQQKDAlTRU5TRVRJTUUxITAfBgNVBAMMGHByaXZhdGUuY2Euc2Vuc2V0aW1lLmNv
bTAeFw0xNzEyMDYwOTM0MzdaFw0yNzEyMDQwOTM0MzdaMGgxCzAJBgNVBAYTAkNO
MRAwDgYDVQQIDAdCRUlKSU5HMRAwDgYDVQQHDAdCRUlKSU5HMRIwEAYDVQQKDAlT
RU5TRVRJTUUxITAfBgNVBAMMGHByaXZhdGUuY2Euc2Vuc2V0aW1lLmNvbTCCASIw
DQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALPbG4PxtqX9TEk720hkxqlY07WB
KWg3MD51jzZzVEDe0LnsD0kmdSt0lA+WvIGwXNXh0TNX9B7zcNwJ+dhj6oEujA+Z
zmd3FpulpJElU0nE/R68LzTa/4bXCIwMmpkKvMbuLdwSNimbSKiO9IGrloCNFTfP
Fskmmp3NbcXkNFQCRseGFUGGJDfsNdSp5qGsTIolpqoBRlHyxsHxqzk3PVkvRZ0u
7ytQKQENbb4w60ukqh45hLX6J0irQfqSY8Bw51gos3OfQ3ur8z3HdFMp+/PxMh4n
rAMvqBLe4d6fBj+oj2Ej27gQZ8aDvV1jWh92rN5A9RKTM3XV90PRGHzMvn0CAwEA
AaNQME4wHQYDVR0OBBYEFBc2fH74sxyPX/N+TbATRDVmcM1+MB8GA1UdIwQYMBaA
FBc2fH74sxyPX/N+TbATRDVmcM1+MAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEL
BQADggEBADScq9hKnAFlGw5gWJoNuTx6FPD2MJ6Zm0/VoD7xNS32nIaVVI0Tt6VH
eZe0JD7Cer4LIPUb5oJTmcR2mUYgBhVLtZKoLRwgH7daRqaI/LOdV8XQR+qRqyj6
iBtYOZmumXqvsW2NsrxV/fAWbXeVZl3bE7YVfbvktBhdFNT05DVEJDu+0QmoClHN
e39TYZbLuUgfBIVZUVItKJfp1NVVX6M5U+/KEzwxShAVOez/S3Jsn+dROKBf6WQn
mLmCh5WMppaIbSjWatz2hBcqarh12gGQgNwyd+zyWbqtCddEdaxNW8WLj1Y8JLxH
rO2hAGzKct7qiBd6mDCBJfSWIVxKU0Q=
-----END CERTIFICATE-----
`

const defaultSlaveCACert = `
-----BEGIN CERTIFICATE-----
MIIDrzCCApegAwIBAgIJAOI2xfBCEdAmMA0GCSqGSIb3DQEBCwUAMG4xCzAJBgNV
BAYTAkNOMRAwDgYDVQQIDAdCRUlKSU5HMRAwDgYDVQQHDAdCRUlKSU5HMRIwEAYD
VQQKDAlTRU5TRVRJTUUxJzAlBgNVBAMMHnNsYXZlLnByaXZhdGUuY2Euc2Vuc2V0
aW1lLmNvbTAeFw0xNzEyMTIxMTA3MjNaFw0yNzEyMTAxMTA3MjNaMG4xCzAJBgNV
BAYTAkNOMRAwDgYDVQQIDAdCRUlKSU5HMRAwDgYDVQQHDAdCRUlKSU5HMRIwEAYD
VQQKDAlTRU5TRVRJTUUxJzAlBgNVBAMMHnNsYXZlLnByaXZhdGUuY2Euc2Vuc2V0
aW1lLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALPbG4PxtqX9
TEk720hkxqlY07WBKWg3MD51jzZzVEDe0LnsD0kmdSt0lA+WvIGwXNXh0TNX9B7z
cNwJ+dhj6oEujA+Zzmd3FpulpJElU0nE/R68LzTa/4bXCIwMmpkKvMbuLdwSNimb
SKiO9IGrloCNFTfPFskmmp3NbcXkNFQCRseGFUGGJDfsNdSp5qGsTIolpqoBRlHy
xsHxqzk3PVkvRZ0u7ytQKQENbb4w60ukqh45hLX6J0irQfqSY8Bw51gos3OfQ3ur
8z3HdFMp+/PxMh4nrAMvqBLe4d6fBj+oj2Ej27gQZ8aDvV1jWh92rN5A9RKTM3XV
90PRGHzMvn0CAwEAAaNQME4wHQYDVR0OBBYEFBc2fH74sxyPX/N+TbATRDVmcM1+
MB8GA1UdIwQYMBaAFBc2fH74sxyPX/N+TbATRDVmcM1+MAwGA1UdEwQFMAMBAf8w
DQYJKoZIhvcNAQELBQADggEBAG8vG7uYYFpgwU6ZG1tVxjhMhMFnI7iIasX6kFrd
7yi8N5T3PnYQfHY2ryCkZK6lkdqOhYjX7QuIptRhKeZtIKzkJZIzC2ImnQImf+ah
WIkhN5pmuaA9rb43NRxnfCwLKbMxnheZnBUnFg/Ty83yYTcDEs2zAjNmiGJKLERn
xIUnoWEiXb/tGTatTPNwmNtWbrfy3AeFP39iRD82FPXtsMve45+EnGpt2WAXjx/q
LSbFMBojo7wGfFUu8rw7RDt9b8XgOgjQNYLUlct4MtIsCFMZJU17gCBJ5DFRTnHC
MFD+L3DkGdtm5sbsgdsVB9F3vhsnFWO8y9E2uusM4G8rnT8=
-----END CERTIFICATE-----
`

const defaultConsoleKey = `
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAxV8OzmENKTgpshVjko98tT8aLeD61g+ujdVj9aviOKKm5Edh
04jPwzr3n41ZMDf+B/xPqe7HWTxpv5Lu4Kqa1JBi0qqZZFShBqxcLQuAjVzNIGZe
sDhIrgm93ubLP8ZwvK+vdZmPbOnx1VKfjeOglZKWS+VrwXG61IM+0iYPTXaOJX/J
QDucSvcXeyfUxnybC4lgQdCQeTfF+nTWunO0a7A9vrbx79uzN5yUy+c5RBNuJhXd
MqXDfWvI46hLqPJS2zTDnG5CVKBWm1N0G3HGtXBa4SiwgAn2g3My8TluIQU85ThQ
Sy0umI//yMy8kXY5lLqxA2n72zfezS7PrHZ8zwIDAQABAoIBAAE7U6NUFbnxIMl8
uq9ad+PFrgslQUt+s48tCr+ov/OsiDAahfDFBM7qGkuDnU/guZQhLfoYhGP5LYvF
hfoe9nJnKEa6S9TFdm/NOZIKZVX8g0c1fFfLMiDr7KRsek4+lcuHqSepuqxqVVkI
d/hxuDnWvVth5idB53GWFBlJpYTNOshNfllkN0+Gwyo5QZRt2aWTnyp8+g0wPq8t
cVI1U5YB6tKABpl7qA25OhZhdHZF0tmwN51rsJ1YoML/PRpu1p3DH5FhJ81IDe0A
U5lpBGkMuSORNftgba/8LfRwyqTn/YYn579aZl6579C4K1ANxWJJnO64zZahDCNV
YrFz4WkCgYEA/drKBwlx+DPCSnX8f0thXNkrq9yk1l1m5+EIcd1oiHPl3iFHRJoo
iHXh5P8RPvVxkzB9LPnnyEx9aMJ6tEONDYkf13K88YcR4JwMcVDdKLSGrpHK0mDc
5xlkd/ueYwfiEKqnCwGOXw4BkPKopL5fyYaHLEn8GwyI1Mcgjw41KTMCgYEAxwoR
XELzw4w0/nnzwBVwUfV76OdFrXylwktq1H7AMMQmfWXTIrBMLkY/f9mkALydZZeq
pm8BjuKKKUicUN98Zh5LK7EQV2ogK3ps0OH6wQNVNnSz5oRRUsEkrM0Uaw15R613
qyZ6Dg3OkDmr5ZyUmVROS45oDPoZlz4PfM63tfUCgYBk8DlCwQuzQIlx6CZFS2jk
bWoDBVH59tuzOfSMqhglocf2Ik9fRNj3IcB3uMBXw2qstywe1SPHrjpzjFkUEoQk
rLCfj3z3oNiH8iS0bg3yYI3pHgmCy4cq0Rr05nUdNYY7UE/pfW3p9/zBcOuDzjry
O+7FuolnC/3gdWlJ2MFkpwKBgDALCxO1CXfbAPOn5iEoS5tM4OLf6B6vJqeWYqv2
CFf9ELlV+be2zDyjMjKfCwoufOOHz2YrBzpBDk5Wu3x95V4U09ow/BvNfwRfoaJt
2YP7VPc3BjGPIL4T5tFbEyGf9/VINsl2GSIJTSHc+dQLjobQJbHxJsZzG/g4v65F
i2x9AoGBANYmI4o2ZGoa1ywxwVoOvQtl3JmVEwy8Wwfp9qv/9bmijjlkKBQFQNWC
udksam8veH8inhBWok1G84T7e46IKcxSVz+26KaG3ViBIyVgrlJkSZ5JVBkZFRaG
JyYxeOfzwryGz1Z6plgVStXq6O9RPuGUeVkT1K5CmeMM0PQlmaG0
-----END RSA PRIVATE KEY-----
`

const defaultConsoleCert = `
-----BEGIN CERTIFICATE-----
MIIDdDCCAlwCCQD6c+kzCWZJDTANBgkqhkiG9w0BAQsFADB8MQswCQYDVQQGEwJD
TjEQMA4GA1UECAwHQmVpamluZzEQMA4GA1UEBwwHQmVpamluZzESMBAGA1UECgwJ
U2Vuc2V0aW1lMRIwEAYDVQQLDAlTZW5zZXRpbWUxITAfBgNVBAMMGHByaXZhdGUu
Y2Euc2Vuc2V0aW1lLmNvbTAeFw0xODA3MTcwODU3MTdaFw0yODA3MTQwODU3MTda
MHwxCzAJBgNVBAYTAkNOMRAwDgYDVQQIDAdCZWlqaW5nMRAwDgYDVQQHDAdCZWlq
aW5nMRIwEAYDVQQKDAlTZW5zZXRpbWUxEjAQBgNVBAsMCVNlbnNldGltZTEhMB8G
A1UEAwwYcHJpdmF0ZS5jYS5zZW5zZXRpbWUuY29tMIIBIjANBgkqhkiG9w0BAQEF
AAOCAQ8AMIIBCgKCAQEAxV8OzmENKTgpshVjko98tT8aLeD61g+ujdVj9aviOKKm
5Edh04jPwzr3n41ZMDf+B/xPqe7HWTxpv5Lu4Kqa1JBi0qqZZFShBqxcLQuAjVzN
IGZesDhIrgm93ubLP8ZwvK+vdZmPbOnx1VKfjeOglZKWS+VrwXG61IM+0iYPTXaO
JX/JQDucSvcXeyfUxnybC4lgQdCQeTfF+nTWunO0a7A9vrbx79uzN5yUy+c5RBNu
JhXdMqXDfWvI46hLqPJS2zTDnG5CVKBWm1N0G3HGtXBa4SiwgAn2g3My8TluIQU8
5ThQSy0umI//yMy8kXY5lLqxA2n72zfezS7PrHZ8zwIDAQABMA0GCSqGSIb3DQEB
CwUAA4IBAQDFPS+zoXkDOAy6Y7dI0kwjpQlqZhKjPni1LAuCbNbebpGve9RlZQPr
p0fu1zD8vwgnAwqOCJSPtdw2SphRQCmwOjkEazfHZezve3eJ3hIaMXO89Kwn14ye
SGg9l/1/cwCun61kiQzvW2tMK8KxUtvO62TLmKZMmu6iMj1Koi98TsBnDHhMfpv1
c1UTgXKLq+W7vJrMAwlQqbGI6xjxGG4AHpEkrK23qDlGqkJ1uNXtNf5+dQe14+j4
MopZ2DjS2c2+Z0GfWW6D9IxoZVqq0eBdfLskc6THEh/JwZWbpIYL7k/zQELiB6HO
UUuPtAJTIpTsOPk8nMmG93jLTQYERIa6
-----END CERTIFICATE-----
`

const (
	// Master presents Master CA
	Master ServerType = iota
	// Slave presents Slave CA
	Slave
)

// ServerType defines master or slave server type
type ServerType uint

// CACtl is license-ca control interface
type CACtl struct {
	masterCACert string                      // master ca cert, signed from same key as slave ca cert
	slaveCACert  string                      // slave ca cert, signed from same key as master ca cert
	clientCert   string                      // client cert, signed by client key
	clientKey    string                      // client key
	caURL        map[ServerType]string       // ca urls
	httpClient   map[ServerType]*http.Client // https client
}

// CAStatus carries CA status from CA admin api
type CAStatus struct {
	Mode        string                 // default is voucher mode, if not, you should to change code(rock add)
	Server      uint                   // 0 is master, 1 is slave
	Disable     bool                   // is ca disabled, by soft start/stop ca, if disabled, ca can't supply nomal service
	IsActive    bool                   // master or standby ca
	ActiveLimit int32                  // cluster total active limit
	AloneTime   int32                  // ca alone time, uint seconds, 0 means forever
	DongleTime  int64                  // dongle timestamp
	Status      string                 // ca status, "alive" or "dead", means whether ca is in alive
	AuthID      string                 // cluster license sn
	Product     string                 // product name
	DongleID    string                 // dongle id
	ExpiredAt   string                 // expire time
	Company     string                 // company name
	FeatureIds  []uint64               // feature ids
	Quotas      map[string]quotaLimit  // cluster quotas, used and total
	Consts      map[string]interface{} // cluster consts, value type will be int32 or string
	Devices     []caDeviceInfo         // the quotas that devices have taken
}

// CAStatusRequest is ca request message to license-ca
type CAStatusRequest struct {
}

// ServiceControlRequest is request message to license-ca
type ServiceControlRequest struct {
	Disable bool `json:"disable,omitempty"`
}

// ServiceControlResponse is response from license-ca
type ServiceControlResponse struct {
	Success bool `json:"success,omitempty"`
}

// ActiveResponse is status from ca
type ActiveResponse struct {
	StatusCode    string `json:"status_code,omitempty"`
	StatusMessage string `json:"status_message,omitempty"`
}

// HardwareInfoResponse contains finger_print, dongle c2v
type HardwareInfoResponse struct {
	FingerPrint string `json:"finger_print,omitempty"`
	C2V         string `json:"c2v,omitempty"`
}

// ClientLicResponse ...
type ClientLicResponse struct {
	Licenses []string `json:"licenses,omitempty"`
}

type onlineActiveRequest struct {
	Action string `json:"action,omitempty"`
}

type hardwareInfoRequest struct {
	Type int32 `json:"type,omitempty"`
}

type offlineActiveRequest struct {
	V2C string `json:"v2c,omitempty"`
}

type clientLicRequest struct{}

// caStatusResponse is ca response message from license-ca, UpdateAt is different from rpc
type caStatusResponse struct {
	Disable    bool     `json:"disable,omitempty"`
	IsActive   bool     `json:"is_active,omitempty"`
	AloneTime  int32    `json:"alone_time,omitempty"`
	Status     string   `json:"status,omitempty"`
	Cluster    string   `json:"cluster,omitempty"`
	UpdatedAt  string   `json:"updated_at,omitempty"`
	CaJson     string   `json:"ca_json,omitempty"`
	Product    string   `json:"product,omitempty"`
	DongleId   string   `json:"dongle_id,omitempty"`
	ExpiredAt  string   `json:"expired_at,omitempty"`
	DongleTime string   `json:"dongle_time,omitempty"`
	Company    string   `json:"company,omitempty"`
	Slaves     []string `json:"slaves,omitempty"`
	FeatureIds []string `json:"feature_ids,omitempty"`
}

type caValues interface {
	getConsts() map[string]interface{}
	getQuotas() map[string][2]int32
	getDevices() []caDeviceInfo
	getLimit() int32
	getActive() int32
}

type caDeviceInfo struct {
	UdID       string           `json:"udid,omitempty"`
	QuotaUsage map[string]int32 `json:"quota_usage,omitempty"`
}

type dongleCAValues struct {
	Limit   int32                  `json:"limit,omitempty"`
	Active  int32                  `json:"active,omitempty"`
	Quotas  map[string][2]int32    `json:"quotas,omitempty"`
	Consts  map[string]interface{} `json:"consts,omitempty"`
	Devices []caDeviceInfo         `json:"devices",omitempty`
}

func (d *dongleCAValues) getConsts() map[string]interface{} {
	return d.Consts
}

func (d *dongleCAValues) getQuotas() map[string][2]int32 {
	return d.Quotas
}

func (d *dongleCAValues) getDevices() []caDeviceInfo {
	return d.Devices
}

func (d *dongleCAValues) getLimit() int32 {
	return d.Limit
}

func (d *dongleCAValues) getActive() int32 {
	return d.Active
}

type voucherCAValues struct {
	Limit   int32                  `json:"limit,omitempty"`
	Active  int32                  `json:"active,omitempty"`
	Quotas  map[string][2]int32    `json:"ext_quotas,omitempty"`
	Consts  map[string]interface{} `json:"ext_consts,omitempty"`
	Devices []caDeviceInfo         `json:"devices,omitempty"`
}

func (d *voucherCAValues) getConsts() map[string]interface{} {
	return d.Consts
}

func (d *voucherCAValues) getQuotas() map[string][2]int32 {
	return d.Quotas
}

func (d *voucherCAValues) getDevices() []caDeviceInfo {
	return d.Devices
}

func (d *voucherCAValues) getLimit() int32 {
	return d.Limit
}

func (d *voucherCAValues) getActive() int32 {
	return d.Active
}

type quotaLimit struct {
	Used  int32 // used quotas
	Total int32 // total quotas
}

// NewServiceCtl creates a CACtl
func NewServiceCtl(masterCAUrl, slaveCAUrl string) (*CACtl, error) {
	ctl := &CACtl{
		masterCACert: defaultMasterCACert,
		slaveCACert:  defaultSlaveCACert,
		clientCert:   defaultConsoleCert,
		clientKey:    defaultConsoleKey,
	}
	ctl.caURL = make(map[ServerType]string)
	//ctl.caURL[Master] = defaultMasterCAUrl
	//ctl.caURL[Master] = defaultSlaveCAUrl
	ctl.caURL[Master] = masterCAUrl // "https://10.151.5.136:8443"
	ctl.caURL[Slave] = slaveCAUrl   // "https://10.151.5.137:8443"
	var err error
	ctl.httpClient = make(map[ServerType]*http.Client, 2)
	ctl.httpClient[Master], err = createHTTPSClient(ctl.masterCACert, ctl.clientKey, ctl.clientCert)
	if err != nil {
		return nil, err
	}
	ctl.httpClient[Slave], err = createHTTPSClient(ctl.slaveCACert, ctl.clientKey, ctl.clientCert)
	if err != nil {
		return nil, err
	}
	return ctl, nil
}

func newCAValues(authAddr string) caValues {
	switch authAddr {
	case "dongle":
		return new(dongleCAValues)
	case "voucher":
		return new(voucherCAValues)
	default:
		return nil
	}
}

// GetCAStatus get CA status from license-ca
// serverType: Master, get status from master license-ca and Slave gets from slave license-ca
func (ctl *CACtl) GetCAStatus(serverType ServerType, authAddr string) (*CAStatus, error) {
	req, err := http.NewRequest("GET", ctl.caURL[serverType]+"/status", nil)
	if err != nil {
		return nil, err
	}
	resp, err := ctl.httpClient[serverType].Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code is %d", resp.StatusCode)
	}
	caRet, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	caResp := new(caStatusResponse)
	err = json.Unmarshal(caRet, caResp)
	if err != nil {
		return nil, err
	}
	caVal := newCAValues(authAddr)
	if caVal == nil {
		return nil, fmt.Errorf("not support auth addr: ", authAddr)
	}
	err = json.Unmarshal([]byte(caResp.CaJson), caVal)
	if err != nil {
		return nil, err
	}
	quotas := make(map[string]quotaLimit, len(caVal.getQuotas()))
	for k, v := range caVal.getQuotas() {
		quotas[k] = quotaLimit{Total: v[1], Used: v[0]}
	}
	dongleTime, _ := strconv.ParseInt(caResp.DongleTime, 10, 64)
	featureIDs := make([]uint64, 0)
	for _, featureID := range caResp.FeatureIds {
		id, _ := strconv.ParseUint(featureID, 10, 64)
		featureIDs = append(featureIDs, id)
	}
	status := &CAStatus{
		Mode:        authAddr,
		Server:      uint(serverType),
		Status:      caResp.Status,
		Disable:     caResp.Disable,
		IsActive:    caResp.IsActive,
		AuthID:      caResp.Cluster,
		Product:     caResp.Product,
		DongleID:    caResp.DongleId,
		ExpiredAt:   caResp.ExpiredAt,
		DongleTime:  dongleTime,
		AloneTime:   caResp.AloneTime,
		Company:     caResp.Company,
		FeatureIds:  featureIDs,
		ActiveLimit: caVal.getLimit(),
		Quotas:      quotas,
		Consts:      caVal.getConsts(),
		Devices:     caVal.getDevices(),
	}
	return status, nil
}

// CAControl sends start or stop cmd to license-ca, change ca service status
// serverType: Master, control master license-ca and Slave control slave license-ca
// true means stop license-ca and false means start server
func (ctl *CACtl) CAControl(serverType ServerType, disable bool) (*ServiceControlResponse, error) {
	scReq := &ServiceControlRequest{Disable: disable}
	reqBody, err := json.Marshal(scReq)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", ctl.caURL[serverType]+"/control", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := ctl.httpClient[serverType].Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code is %d", resp.StatusCode)
	}
	caRet, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	caResp := new(ServiceControlResponse)
	err = json.Unmarshal(caRet, caResp)
	if err != nil {
		return nil, err
	}
	return caResp, nil
}

// OnlineActivate sends a action command to ca activate dongle online
func (ctl *CACtl) OnlineActivate(serverType ServerType, action string) (*ActiveResponse, error) {
	olReq := &onlineActiveRequest{Action: action}
	reqBody, err := json.Marshal(olReq)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", ctl.caURL[serverType]+"/online", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := ctl.httpClient[serverType].Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code is %d", resp.StatusCode)
	}
	caRet, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	caResp := new(ActiveResponse)
	err = json.Unmarshal(caRet, caResp)
	if err != nil {
		return nil, err
	}
	return caResp, nil
}

// HardwareInfo gets fingerprint or c2v from ca
func (ctl *CACtl) HardwareInfo(serverType ServerType, Type int32) (*HardwareInfoResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(ctl.caURL[serverType]+"/hdinfo/%d", Type), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := ctl.httpClient[serverType].Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code is %d, %s", resp.StatusCode, resp.Status)
	}
	caRet, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	caResp := new(HardwareInfoResponse)
	err = json.Unmarshal(caRet, caResp)
	if err != nil {
		return nil, err
	}
	return caResp, nil
}

// OfflineActivate sends v2c file to ca activate dongle offline
func (ctl *CACtl) OfflineActivate(serverType ServerType, v2c string) (*ActiveResponse, error) {
	offReq := &offlineActiveRequest{V2C: v2c}
	reqBody, err := json.Marshal(offReq)
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(reqBody))
	req, err := http.NewRequest("POST", ctl.caURL[serverType]+"/offline", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := ctl.httpClient[serverType].Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code is %d", resp.StatusCode)
	}
	caRet, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	caResp := new(ActiveResponse)
	err = json.Unmarshal(caRet, caResp)
	if err != nil {
		return nil, err
	}
	return caResp, nil
}

// GetClientLics get client licenses from license-ca
// serverType: Master, get status from master license-ca and Slave gets from slave license-ca
func (ctl *CACtl) GetClientLics(serverType ServerType) (*ClientLicResponse, error) {
	req, err := http.NewRequest("GET", ctl.caURL[serverType]+"/clics", nil)
	if err != nil {
		return nil, err
	}
	resp, err := ctl.httpClient[serverType].Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code is %d", resp.StatusCode)
	}
	caRet, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	caResp := new(ClientLicResponse)
	err = json.Unmarshal(caRet, caResp)
	if err != nil {
		return nil, err
	}

	return caResp, nil
}

// create a https client with server and client certs
func createHTTPSClient(serverCert, clientKey, clientCert string) (*http.Client, error) {
	// add client key and cert
	cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
	if err != nil {
		return nil, err
	}

	// add server cert
	serverCertPool := x509.NewCertPool()
	if ok := serverCertPool.AppendCertsFromPEM([]byte(serverCert)); !ok {
		return nil, errors.New("load ca cert error")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      serverCertPool,
		// InsecureSkipVerify控制客户端是否验证服务器的证书链和主机名。
		// 如果InsecureSkipVerify(非安全的跳过验证)为true，那么TLS接受服务器提供的任何证书以及该证书中的任何主机名。
		// 在这种模式下，TLS容易受到中间人攻击。这应该只用于测试。
		InsecureSkipVerify: true,
	}

	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	httpClient := &http.Client{Transport: transport, Timeout: time.Second * 10} //http request timeout, 3 seconds
	return httpClient, nil
}
