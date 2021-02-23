package museum

import (
	"bytes"
	"encoding/json"
	"go.rock.com/rock-platform/rock/server/conf"
	"go.rock.com/rock-platform/rock/server/utils"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type client struct {
	client *http.Client
	addr   string
}

// 构造函数-多态实例化
// NewClient returns a client at the specified url.
func NewClient(cli *http.Client, chartRepoAddr string) Client {
	return &client{client: cli, addr: strings.TrimSuffix(chartRepoAddr, "/")}
}

// helper function for making an http GET request.
func (c *client) get(rawUrl string, out interface{}) error {
	return c.do(rawUrl, "GET", nil, out)
}

// helper function for making an http DELETE request.
func (c *client) delete(rawUrl string) error {
	return c.do(rawUrl, "DELETE", nil, nil)
}

// helper function to make an http request
func (c *client) do(rawUrl, method string, in, out interface{}) error {
	body, err := c.open(rawUrl, method, in, out)
	if err != nil {
		return err
	}

	defer body.Close()
	if out != nil {
		return json.NewDecoder(body).Decode(out) // 将body序列化到out中保存(指针类型的out,所以不需要return)
	}
	return nil
}

// helper function to open an http request
func (c *client) open(rawUrl, method string, in, out interface{}) (io.ReadCloser, error) {
	uri, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, uri.String(), nil)
	if err != nil {
		return nil, err
	}
	if in != nil {
		decoded, err := json.Marshal(in)
		if err != nil {
			return nil, err
		}
		buf := bytes.NewBuffer(decoded)
		req.Body = ioutil.NopCloser(buf)
		req.ContentLength = int64(len(decoded))
		req.Header.Set("Content-Length", strconv.Itoa(len(decoded)))
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > http.StatusPermanentRedirect {
		defer resp.Body.Close()
		out, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			e := &utils.RockError{
				HttpCode: resp.StatusCode,
				ErrCode:  resp.StatusCode*100000 + 1,
				Message:  err.Error(),
			}
			return nil, e
		}
		err = &utils.RockError{
			HttpCode: resp.StatusCode,
			ErrCode:  resp.StatusCode*100000 + 1,
			Message:  string(out),
		}
		return nil, err
	}
	return resp.Body, nil
}

// generate chartmuseum client
func GetMuseumClient() Client {
	config := conf.GetConfig()
	chartRepoAddr := config.Viper.GetString("chartsRepo.addr")
	client := &http.Client{}
	return NewClient(client, chartRepoAddr)
}

// 多态 interface
type Client interface {
	Charts() (*ChartMapper, error)
	Versions(name string) ([]*ChartVersion, error)
	Version(name, version string) (*ChartVersion, error)
	DeleteVersion(name, version string) error
}
