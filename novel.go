package pixiv

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"golang.org/x/net/proxy"
)

type Pixiv struct {
	Accesstoken  string
	Refreshtoken string
	HttpClient   *http.Client
	Mux          sync.Mutex
}

const (
	ApiAddress string = "https://app-api.pixiv.net"
)

func NewApp() *Pixiv {
	return &Pixiv{HttpClient: &http.Client{Timeout: 3000 * time.Millisecond}}
}

// 设置超时时间 毫秒
func (p *Pixiv) SetHttpTimeout(t int) {
	p.Mux.Lock()
	defer p.Mux.Unlock()
	p.HttpClient.Timeout = time.Duration(t) * time.Millisecond
}

// 设置代理
func (p *Pixiv) SetProxy(Proxy string) error {
	p.Mux.Lock()
	defer p.Mux.Unlock()

	proxyUrl, err := url.Parse(Proxy)
	if err != nil {
		return err
	}
	if proxyUrl.Scheme == "http" || proxyUrl.Scheme == "https" {
		// 设置代理
		proxyURL, err := url.Parse(Proxy)
		if err != nil {
			return err
		}
		// 创建一个Transport并设置代理
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
		p.HttpClient.Transport = transport
		return nil
	} else if proxyUrl.Scheme == "socks5" {
		// 解析代理URL
		proxyURL, err := url.Parse(Proxy)
		if err != nil {
			return err
		}
		dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
		if err != nil {
			return err
		}
		transport := &http.Transport{
			//Dial: dialer.Dial,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			},
		}

		p.HttpClient.Transport = transport
		return nil
	} else {
		return fmt.Errorf("未知的代理类型 %s", Proxy)
	}
}

// 测试
func (p *Pixiv) Test() {

	req, err := http.NewRequest("GET", "http://ip.sb", nil)
	if err != nil {
		fmt.Println("错误1", err)
		return
	}
	req.Header.Set("User-Agent", "curl/8.0")
	resq, err := p.HttpClient.Do(req)
	if err != nil {
		fmt.Println("错误2", err)
		return
	}
	body := GetBody(resq.Body)
	fmt.Println(body)

}

// 获取body内容
func GetBody(body io.Reader) string {
	b, err := io.ReadAll(body)
	if err != nil {
		return ""
	}
	return string(b)
}

// 获取系列小说
func (p *Pixiv) GetSeriesNovels(SeriesID string) {

}

// 获取用户小说
func (p *Pixiv) GetUserNovels(UserID string) {

}
