package license

//
//import (
//	"bytes"
//	"crypto/tls"
//	"crypto/x509"
//	"encoding/json"
//	"fmt"
//	"go.rock.com/rock-platform/rock/server/utils"
//	"io"
//	"io/ioutil"
//	"net/http"
//	"net/url"
//	"strconv"
//	"strings"
//	"time"
//)
//
//type client struct {
//	client *http.Client
//	addr   string
//}
//
//// New returns a client at the specified url.
//func New(uri string) Client {
//	return &client{http.DefaultClient, strings.TrimSuffix(uri, "/")}
//}
//
//// NewClient returns a client at the specified url.
//func NewClient(uri string, cli *http.Client) Client {
//	fmt.Println("uri: --->", uri)
//	fmt.Println("cli:--->")
//	fmt.Printf("%#v\n", cli)
//	return &client{cli, strings.TrimSuffix(uri, "/")}
//}
//
//// helper function for making an http GET request.
//func (c *client) get(rawurl string, out interface{}) error {
//	return c.do(rawurl, "GET", nil, out)
//}
//
//// helper function for making an http POST request.
//func (c *client) post(rawurl string, in, out interface{}) error {
//	return c.do(rawurl, "POST", in, out)
//}
//
//// helper function for making an http PUT request.
//func (c *client) put(rawurl string, in, out interface{}) error {
//	return c.do(rawurl, "PUT", in, out)
//}
//
//// helper function for making an http PATCH request.
//func (c *client) patch(rawurl string, in, out interface{}) error {
//	return c.do(rawurl, "PATCH", in, out)
//}
//
//// helper function for making an http DELETE request.
//func (c *client) delete(rawurl string) error {
//	return c.do(rawurl, "DELETE", nil, nil)
//}
//
//// helper function to make an http request
//func (c *client) do(rawurl, method string, in, out interface{}) error {
//	body, err := c.open(rawurl, method, in, out)
//	if err != nil {
//		return err
//	}
//	defer body.Close()
//	if out != nil {
//		return json.NewDecoder(body).Decode(out)
//	}
//	return nil
//}
//
//// helper function to open an http request
//func (c *client) open(rawurl, method string, in, out interface{}) (io.ReadCloser, error) {
//	uri, err := url.Parse(rawurl)
//	if err != nil {
//		return nil, err
//	}
//	fmt.Println("open uri.String()----------->")
//	fmt.Println(uri.String()) // https://10.151.3.87:8443/status
//	req, err := http.NewRequest(method, uri.String(), nil)
//	if err != nil {
//		return nil, err
//	}
//	fmt.Println("open http.NewRequest uri:")
//	fmt.Println(uri.String())
//	fmt.Println("in")
//	fmt.Println(in)
//	//if in != nil {  // in = nil
//	//	decoded, derr := json.Marshal(in)
//	//	if derr != nil {
//	//		return nil, derr
//	//	}
//	//	buf := bytes.NewBuffer(decoded)
//	//	req.Body = ioutil.NopCloser(buf)
//	//	req.ContentLength = int64(len(decoded))
//	//	req.Header.Set("Content-Length", strconv.Itoa(len(decoded)))
//	//	req.Header.Set("Content-Type", "application/json")
//	//	fmt.Println("in != nil !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! , req:")
//	//	fmt.Printf("%#v\n", req)
//	//}
//	resp, err := c.client.Do(req)
//	if err != nil {
//		return nil, err
//	}
//	//if resp.StatusCode > http.StatusPermanentRedirect {
//    // fmt.Println("resp.StatusCode > 308 !!!!!!!!!!!!!!!")
//	//	defer resp.Body.Close()
//	//	out, err := ioutil.ReadAll(resp.Body)
//	//	if err != nil {
//	//		e := &utils.RockError{
//	//			HttpCode: resp.StatusCode,
//	//			ErrCode:  resp.StatusCode*100000 + 1,
//	//			Message:  err.Error(),
//	//		}
//	//		return nil, e
//	//	}
//	//	err = &utils.RockError{
//	//		HttpCode: resp.StatusCode,
//	//		ErrCode:  resp.StatusCode*100000 + 1,
//	//		Message:  string(out),
//	//	}
//	//	return nil, err
//	//}
//	return resp.Body, nil
//}
//
//func GetLicenseClient(config string) (Client, error) {
//	caUrl, err := utils.GetLicenseCaUrl(config) // https://10.151.3.94:8443
//	if err != nil {
//		return nil, err
//	}
//
//	client, err := createHTTPSClient(defaultMasterCACert, defaultConsoleKey, defaultConsoleCert)
//	if err != nil {
//		return nil, err
//	}
//
//	return NewClient(caUrl, client), nil
//}
//
//// create a https client with server and client certs
//func createHTTPSClient(serverCert, clientKey, clientCert string) (*http.Client, error) {
//	// add client key and cert
//	cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
//	if err != nil {
//		return nil, err
//	}
//
//	// add server cert
//	serverCertPool := x509.NewCertPool()
//	if ok := serverCertPool.AppendCertsFromPEM([]byte(serverCert)); !ok {
//		return nil, utils.NewRockError(400, 40000004, "load CA cert error")
//	}
//
//	tlsConfig := &tls.Config{
//		Certificates:       []tls.Certificate{cert},
//		RootCAs:            serverCertPool,
//		InsecureSkipVerify: true,
//	}
//
//	tlsConfig.BuildNameToCertificate()
//	transport := &http.Transport{TLSClientConfig: tlsConfig}
//	httpClient := &http.Client{
//		Transport: transport,
//		Timeout:   time.Second * 10,
//	} //http request timeout, 3 seconds
//
//	return httpClient, nil
//}
//
//type Client interface {
//	SetClient(*http.Client)
//	SetAddress(string)
//	Status() (*CAStatus, error)
//	C2v() (*C2vResp, error)
//	Fingerprint() (*FingerprintResp, error)
//	Online() (*ActiveResp, error)
//	Offline(string) (*ActiveResp, error)
//	ClientLicenses() (*ClientLicResp, error)
//}
