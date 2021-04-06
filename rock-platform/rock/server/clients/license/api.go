package license

//
//import (
//	"fmt"
//	"net/http"
//)
//
//const (
//	status         = "%s/status"
//	c2v            = "%s/hdinfo/2"
//	fingerprint    = "%s/hdinfo/1"
//	online         = "%s/online"
//	offline        = "%s/offline"
//	clientLicenses = "%s/clics"
//)
//
//// SetClient sets the http.Client.
//func (c *client) SetClient(client *http.Client) {
//	c.client = client
//}
//
//// SetAddress sets the server address.
//func (c *client) SetAddress(addr string) {
//	c.addr = addr
//}
//
//// Returns license's current status
//func (c *client) Status() (*CAStatus, error) {
//	var out *CAStatus
//	uri := fmt.Sprintf(status, c.addr)
//	fmt.Println("Status() uri:----->", uri)
//	err := c.get(uri, &out)
//	return out, err
//}
//
//// Return c2v file content
//func (c *client) C2v() (*C2vResp, error) {
//	var out *C2vResp
//	uri := fmt.Sprintf(c2v, c.addr)
//	err := c.get(uri, &out)
//	return out, err
//}
//
//// Return fingerprint file content
//func (c *client) Fingerprint() (*FingerprintResp, error) {
//	var out *FingerprintResp
//	uri := fmt.Sprintf(fingerprint, c.addr)
//	err := c.get(uri, &out)
//	return out, err
//}
//
//// Active license online
//func (c *client) Online() (*ActiveResp, error) {
//	var out *ActiveResp
//	in := &onlineActiveReq{
//		Action: "activate",
//	}
//	uri := fmt.Sprintf(online, c.addr)
//	err := c.post(uri, in, &out)
//	return out, err
//}
//
//// Active license online
//func (c *client) Offline(v2cContent string) (*ActiveResp, error) {
//	var out *ActiveResp
//	in := &offlineActiveReq{
//		V2C: v2cContent,
//	}
//	uri := fmt.Sprintf(offline, c.addr)
//	err := c.post(uri, in, &out)
//	return out, err
//}
//
//// Get client licenses
//func (c *client) ClientLicenses() (*ClientLicResp, error) {
//	var out *ClientLicResp
//	uri := fmt.Sprintf(clientLicenses, c.addr)
//	err := c.get(uri, &out)
//	return out, err
//}
