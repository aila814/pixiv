package pixiv

import (
	"fmt"
	"testing"
)

func TestPixiv(t *testing.T) {
	//var a NovelInfo

	app := NewApp()

	// app.SetHttpTimeout(5000)
	err := app.SetProxy("socks5://10.10.10.2:1080")
	//err := app.SetProxy("socks5://123:123@10.10.10.3:9200")
	if err != nil {
		fmt.Println("设置代理错误2", err)
		return
	}
	app.Test()
}
