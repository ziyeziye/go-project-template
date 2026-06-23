package fetcher

import (
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

var client *resty.Client

type Fetcher struct {
	Url     string // 请求地址
	Request *resty.Request
	Method  string
}

func SetProxyIP() string {
	if getIp := viper.GetString("PROXY_API"); getIp != "" {
		r := NewRequest()
		get, err := r.Get(getIp)

		if err != nil {
			log.Println("Get Proxy IP Error: " + err.Error())
			return ""
		}

		ip := gjson.Get(get.String(), "proxy").String()
		if ip != "" {
			client.SetProxy("http://" + ip)
		}
		return ip
	}
	return ""
}

func NewRequest() *resty.Request {
	GetClient()
	return client.R()
}

func Get(url string, req ...*resty.Request) *Fetcher {
	return newFetcher(url, "GET", req...)
}

func Post(url string, req ...*resty.Request) *Fetcher {
	return newFetcher(url, "POST", req...)
}

func (f *Fetcher) Do() (*resty.Response, error) {
	return fetch(f.Url, f.Method, f.Request)
}

func (f *Fetcher) DoRetry(retry uint, stop func(*resty.Response) bool) (*resty.Response, error) {
	if stop == nil {
		stop = func(resp *resty.Response) bool {
			return resp.StatusCode() == 200
		}
	}

	for retries := range retry {
		resp, err := f.Do()

		if err == nil && stop(resp) {
			return resp, nil
		}

		fmt.Printf("获取数据失败 (尝试 %d/%d): %v, %s\n", retries+1, retry, err, f.Url)
	}

	return nil, fmt.Errorf("failed to fetch data")
}

func newFetcher(url string, method string, req ...*resty.Request) *Fetcher {
	f := &Fetcher{
		Url:    url,
		Method: method,
	}

	if len(req) > 0 {
		f.Request = req[0]
	} else {
		f.Request = NewRequest()
	}
	return f
}

// 网页内容抓取函数
func fetch(url string, method string, request *resty.Request) (*resty.Response, error) {
	var resp *resty.Response
	var err error
	switch method {
	case "POST":
		resp, err = request.Post(url)
	case "GET":
		resp, err = request.Get(url)
	default:
		return nil, fmt.Errorf("method %s not support", method)
	}

	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetClient() *resty.Client {
	if client == nil {
		client = resty.New()
	}
	return client
}
