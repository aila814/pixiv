package pixiv

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"golang.org/x/net/proxy"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"sync"
	"time"
)

type Pixiv struct {
	AccessToken  string
	RefreshToken string
	HttpClient   *http.Client
	Mux          sync.Mutex
}

const (
	Version           = "v0.2.0"
	ApiAddress string = "https://app-api.pixiv.net"
)

func NewApp() *Pixiv {
	return &Pixiv{HttpClient: &http.Client{}}
}

// 设置超时时间 毫秒
func (p *Pixiv) SetHttpTimeout(t int) {
	p.Mux.Lock()
	defer p.Mux.Unlock()
	p.HttpClient.Timeout = time.Duration(t) * time.Millisecond
}

/*
设置代理 支持http或socks5
http://127.0.0.1:8080
http://123:123@127.0.0.1:8080
socks5://127.0.0.1:8080
socks5://123:123@127.0.0.1:8080
*/
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

// 设置Refreshtoken
func (p *Pixiv) SetRefreshToken(RefreshToken string) {
	p.Mux.Lock()
	defer p.Mux.Unlock()
	p.RefreshToken = RefreshToken
}

// 设置Refreshtoken
func (p *Pixiv) SetAccessToken(AccessToken string) {
	p.Mux.Lock()
	defer p.Mux.Unlock()
	p.AccessToken = AccessToken
}

// 获取AccessToken
func (p *Pixiv) GetAccessToken() error {
	from := url.Values{}
	from.Add("grant_type", "refresh_token")
	from.Add("client_id", "MOBrBDS8blbauoSck0ZfDbtuzpyT")
	from.Add("client_secret", "lsACyCD94FhDUtGTXi3QzcFE2uU1hqtDaKeqrdwj")
	from.Add("include_policy", "true")
	from.Add("refresh_token", p.RefreshToken)
	req, err := http.NewRequest("POST", "https://oauth.secure.pixiv.net/auth/token", bytes.NewBufferString(from.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("accept-language", "zh_CN")
	req.Header.Set("User-Agent", "PixivAndroidApp/5.0.234")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body := GetBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return errors.New(body)
	}
	p.Mux.Lock()
	p.AccessToken = gjson.Get(body, "access_token").String()
	p.Mux.Unlock()
	return nil
}

// 检测AccessToken是否有效
func (p *Pixiv) TestAccessToken() (bool, error) {
	Url := fmt.Sprintf("%s/v1/illust/recommended?include_privacy_policy=true&filter=for_android&include_ranking_illusts=true", ApiAddress)
	req, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("accept-language", "zh_CN")
	req.Header.Set("User-Agent", "PixivAndroidApp/5.0.234 ")
	req.Header.Set("Authorization", "Bearer "+p.AccessToken)

	resp, err := p.HttpClient.Do(req)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {

		return false, fmt.Errorf("%s", GetBody(resp.Body))
	}
	return true, nil
}

// 转换时间到北京时间格式
func convertTimeToBeijing(timeStr string) string {
	// 解析时间字符串
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return timeStr
	}
	// 将时区调整为东八区
	beijingTime := t.In(time.FixedZone("CST", 8*60*60))

	// 格式化为北京时间格式
	beijingTimeStr := beijingTime.Format("2006-01-02 15:04:05")

	return beijingTimeStr
}

// 获取系列小说
func (p *Pixiv) GetSeriesNovels(SeriesID string, OnlyDetail bool) (SeriesNovel, Error) {
	var (
		Title         string // 系列标题
		Count         int64  // 系列小说数量
		Caption       string //系列简介
		userId        string
		last_order    int
		Err           Error
		NovelInfoList []NovelInfo
		SeriesInfo    SeriesNovel
	)

	for {

		Url := fmt.Sprintf("%s/v2/novel/series?series_id=%s&last_order=%d", ApiAddress, SeriesID, last_order)
		req, err := http.NewRequest("GET", Url, nil)
		if err != nil {
			Err.Err = err
			Err.Code = 0
			return SeriesInfo, Err
		}
		req.Header.Set("accept-language", "zh_CN")
		req.Header.Set("User-Agent", "PixivAndroidApp/5.0.234 ")
		req.Header.Set("Authorization", "Bearer "+p.AccessToken)

		resp, err := p.HttpClient.Do(req)
		if err != nil {
			Err.Err = err
			return SeriesInfo, Err
		}
		body := GetBody(resp.Body)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			Err.Code = resp.StatusCode
			Err.Body = body
			if gjson.Get(body, "error.message").String() != "" {
				Err.Err = errors.New(gjson.Get(body, "error.message").String())
			} else {
				Err.Err = errors.New(gjson.Get(body, "error.user_message").String())
			}

			return SeriesInfo, Err
		}
		if len(Title) == 0 {
			Title = gjson.Get(body, "novel_series_detail.title").String()
			Count = gjson.Get(body, "novel_series_detail.content_count").Int()
			Caption = gjson.Get(body, "novel_series_detail.caption").String()
			userId = gjson.Get(body, "novel_series_detail.user.id").String()
		}
		if OnlyDetail {
			break
		}
		novels := gjson.Get(body, "novels")
		novels.ForEach(func(key, value gjson.Result) bool {
			var tagArry []string
			tags := value.Get("tags")
			tags.ForEach(func(key2, value2 gjson.Result) bool {
				tag := value2.Get("name").String()
				tagArry = append(tagArry, tag)
				return true
			})
			NovelInfoList = append(NovelInfoList, NovelInfo{
				UserID:      value.Get("user.id").String(),
				ID:          value.Get("id").String(),
				Title:       value.Get("title").String(),
				Length:      value.Get("text_length").Int(),
				Caption:     value.Get("caption").String(),
				Date:        convertTimeToBeijing(value.Get("create_date").String()),
				Tags:        tagArry,
				SeriesID:    value.Get("series.id").String(),
				SeriesTitle: value.Get("series.title").String(),
			})
			return true
		})

		if len(novels.Array()) != 30 {
			break
		}
		last_order = last_order + 30
	}

	SeriesInfo.UserID = userId
	SeriesInfo.Title = Title
	SeriesInfo.Caption = Caption
	SeriesInfo.Count = Count
	SeriesInfo.Novels = NovelInfoList
	Err.Code = 200
	Err.Err = nil
	return SeriesInfo, Err

}

// 获取用户小说
func (p *Pixiv) GetUserNovels(UserID string) (UserNovel, Error) {
	var (
		Err           Error
		last_order    int
		NovelInfoList []NovelInfo
		userNovel     UserNovel
		userName      string
		userAccount   string
		Count         int64
	)

	for {
		Url := fmt.Sprintf("%s/v1/user/novels?user_id=%s&offset=%d", ApiAddress, UserID, last_order)
		req, err := http.NewRequest("GET", Url, nil)
		if err != nil {
			Err.Err = err
			return userNovel, Err
		}
		req.Header.Set("accept-language", "zh_CN")
		req.Header.Set("User-Agent", "PixivAndroidApp/5.0.234 ")
		req.Header.Set("Authorization", "Bearer "+p.AccessToken)
		resp, err := p.HttpClient.Do(req)
		if err != nil {
			Err.Err = err
			return userNovel, Err
		}
		body := GetBody(resp.Body)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			Err.Code = resp.StatusCode
			Err.Body = body
			if gjson.Get(body, "error.message").String() != "" {
				Err.Err = errors.New(gjson.Get(body, "error.message").String())
			} else {
				Err.Err = errors.New(gjson.Get(body, "error.user_message").String())
			}
			return userNovel, Err
		}

		if len(userName) == 0 {
			userName = gjson.Get(body, "user.name").String()
			userAccount = gjson.Get(body, "user.account").String()
		}
		novels := gjson.Get(body, "novels")
		if len(novels.Array()) != 0 {
			Count = Count + int64(len(novels.Array()))
		}
		novels.ForEach(func(key, value gjson.Result) bool {
			var tagArry []string
			tags := value.Get("tags")
			tags.ForEach(func(key2, value2 gjson.Result) bool {
				tag := value2.Get("name").String()
				tagArry = append(tagArry, tag)
				return true
			})
			NovelInfoList = append(NovelInfoList, NovelInfo{
				UserID:      value.Get("user.id").String(),
				ID:          value.Get("id").String(),
				Title:       value.Get("title").String(),
				Length:      value.Get("text_length").Int(),
				Caption:     value.Get("caption").String(),
				Date:        convertTimeToBeijing(value.Get("create_date").String()),
				Tags:        tagArry,
				SeriesID:    value.Get("series.id").String(),
				SeriesTitle: value.Get("series.title").String(),
			})
			return true
		})
		if len(novels.Array()) != 30 {
			break
		}
		last_order = last_order + 30
	}
	userNovel.UserName = userName
	userNovel.AccountID = userAccount
	userNovel.Count = Count
	userNovel.Novels = NovelInfoList
	Err.Code = 200
	Err.Err = nil
	return userNovel, Err
}

// 获取小说简介
func (p *Pixiv) GetNovelDetail(NovelID string) (NovelInfo, Error) {
	var (
		info NovelInfo
		Err  Error
	)
	Url := fmt.Sprintf("%s/v2/novel/detail?novel_id=%s", ApiAddress, NovelID)
	req, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		Err.Err = err
		return info, Err
	}
	req.Header.Set("accept-language", "zh_CN")
	req.Header.Set("User-Agent", "PixivAndroidApp/5.0.234 ")
	req.Header.Set("Authorization", "Bearer "+p.AccessToken)

	resp, err := p.HttpClient.Do(req)
	if err != nil {
		Err.Err = err
		return info, Err
	}
	body := GetBody(resp.Body)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		Err.Code = resp.StatusCode
		Err.Body = body
		if gjson.Get(body, "error.message").String() != "" {
			Err.Err = errors.New(gjson.Get(body, "error.message").String())
		} else {
			Err.Err = errors.New(gjson.Get(body, "error.user_message").String())
		}
		return info, Err
	}
	// 标签
	gjson.Get(body, "novel.tags").ForEach(func(key, value gjson.Result) bool {
		tag := value.Get("name").String()
		info.Tags = append(info.Tags, tag)
		return true
	})
	// 简介
	info.Caption = gjson.Get(body, "novel.caption").String()
	info.ID = NovelID
	//字数
	info.Length = gjson.Get(body, "novel.text_length").Int()
	//发布时间
	info.Date = convertTimeToBeijing(gjson.Get(body, "novel.create_date").String())
	//标题
	info.Title = gjson.Get(body, "novel.title").String()
	// 用户id
	info.UserID = gjson.Get(body, "novel.user.id").String()
	info.SeriesID = gjson.Get(body, "series.id").String()
	info.SeriesTitle = gjson.Get(body, "series.title").String()
	return info, Err
}

// 获取用户简介
func (p *Pixiv) GetUserDetail(UserID string) (UserInfo, Error) {
	var (
		info UserInfo
		Err  Error
	)

	Url := fmt.Sprintf("%s/v1/user/detail?filter=for_android&user_id=%s", ApiAddress, UserID)
	req, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		Err.Err = err
		return info, Err
	}
	req.Header.Set("accept-language", "zh_CN")
	req.Header.Set("User-Agent", "PixivAndroidApp/5.0.234 ")
	req.Header.Set("Authorization", "Bearer "+p.AccessToken)

	resp, err := p.HttpClient.Do(req)
	if err != nil {
		Err.Err = err
		return info, Err
	}
	body := GetBody(resp.Body)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		Err.Code = resp.StatusCode
		Err.Body = body
		if gjson.Get(body, "error.message").String() != "" {
			Err.Err = errors.New(gjson.Get(body, "error.message").String())
		} else {
			Err.Err = errors.New(gjson.Get(body, "error.user_message").String())
		}
		return info, Err
	}

	info.UserID = gjson.Get(body, "user.id").Int()
	info.UserName = gjson.Get(body, "user.name").String()
	info.AccountID = gjson.Get(body, "user.account").String()
	info.Caption = gjson.Get(body, "user.comment").String()
	info.TotalNovels = gjson.Get(body, "profile.total_novels").Int()
	Err.Code = resp.StatusCode
	Err.Body = body
	return info, Err
}

// 获取小说内容/正文
func (p *Pixiv) GetNovelContent(NovelID string) (NovelContent, Error) {
	var (
		RawnContent string
		content     NovelContent
		Err         Error
	)

	Url := fmt.Sprintf("%s/webview/v2/novel?id=%s", ApiAddress, NovelID)
	req, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		Err.Err = err
		return content, Err
	}
	req.Header.Set("accept-language", "zh_CN")
	req.Header.Set("User-Agent", "PixivAndroidApp/5.0.234 ")
	req.Header.Set("Authorization", "Bearer "+p.AccessToken)

	resp, err := p.HttpClient.Do(req)
	if err != nil {
		Err.Err = err
		return content, Err
	}
	body := GetBody(resp.Body)
	defer resp.Body.Close()
	if len(body) == 0 {
		Err.Code = resp.StatusCode
		Err.Body = body
		Err.Err = errors.New("小说原始正文内容为空")
		return content, Err
	}
	if resp.StatusCode != http.StatusOK {
		Err.Code = resp.StatusCode
		Err.Body = body
		if gjson.Get(body, "error.message").String() != "" {
			Err.Err = errors.New(gjson.Get(body, "error.message").String())
		} else {
			Err.Err = errors.New(gjson.Get(body, "error.user_message").String())
		}
		return content, Err
	}

	// 匹配正文
	re := regexp.MustCompile("novel:(.*),")
	math := re.FindStringSubmatch(body)
	if len(math) != 2 {
		Err.Err = errors.New("匹配正文失败")
		Err.Body = body
		Err.Code = resp.StatusCode
		return content, Err
	}
	RawnContent = math[1]
	// 小说正文
	TextContent := gjson.Get(RawnContent, "text").String()
	content.Images = make(map[string]string)
	// 小说插图
	gjson.Get(RawnContent, "images").ForEach(func(key, value gjson.Result) bool {
		image := value.Get("urls.original").String()
		content.Images[key.String()] = image
		return true
	})
	content.RawnContent = RawnContent
	content.Content = TextContent
	Err.Code = resp.StatusCode
	Err.Body = body
	Err.Err = nil
	return content, Err
}
