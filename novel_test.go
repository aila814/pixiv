package pixiv

import (
	"fmt"
	"github.com/aila814/pixiv/filea"
	"testing"
)

var token = "084SzMXHagZ8BwfKsvmeUflJNg1nDe6DADv1o-t5_vE"

func TestPixiv(t *testing.T) {

	app := NewApp()

	//app.SetProxy("socks5://10.10.10.2:1080")
	app.SetRefreshToken("JMfu2g_lzwHS7ojLnUA3mifsBRPQLgbfqQLZ9W1M73E")
	//app.GetAccessToken()
	//fmt.Println(app.AccessToken)
	//return
	app.SetAccessToken(token)
	//bin, err := app.GetIllustByte("https://i.pximg.net/img-zip-ugoira/img/2024/05/12/23/54/48/118679032_ugoira600x600.zip")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//filea.WriteFileByte("118679032.zip", bin)

	//return

	//_, _, err := loadImage("118679032/000000.jpg")
	//if err != nil {
	//	fmt.Println(fmt.Sprintf("加载图片错误 %v", err))
	//	os.Exit(1)
	//}
	//return
	//怪盗シェイプシフター_112650024_p2.png
	//ピンクラバースーツ_106643910_p0.jpg
	//幽々子の触手もの 03_117906850.gif
	//無題_118679032.gif
	//八雲紫の触手もの_98638017_p0.gif
	//赤城さんの初めて_95047084_p0.gif
	//機ﾇﾁｮ_83991907.gif
	//大妖精の触手もの_93983329_p0.gif
	i, Err := app.GetIllust("93983329")
	if Err.Err != nil {
		fmt.Println("错误：", Err.Err, Err.Code)
		return
	}

	fmt.Println(i.Title)
	//fmt.Println(Err.Body)
	fmt.Println("========================================")

	if i.Type == "ugoira" {
		g, Err := app.GetIllustGif(i.ID)
		if Err.Err != nil {
			fmt.Println("错误：", Err.Err, Err.Code)
			return
		}
		fmt.Println(g.ZipUrl)
		fmt.Println("===========================================================")
		filename := fmt.Sprintf("%s.zip", i.ID)
		var bin []byte
		var err error
		if !filea.FileExists(filename) {
			bin, err = app.GetIllustByte(g.ZipUrl)
			if err != nil {
				fmt.Println("下载gif错误:", err)
				return
			}
			err = filea.WriteFileByte(filename, bin)
			if err != nil {
				fmt.Println("保存zip错误:", err)
				return
			}
		}

		err = filea.Unzip(filename, fmt.Sprintf("./%s", i.ID))
		if err != nil {
			fmt.Println("解压zip错误:", err)
			return
		}
		err = app.ToGif(i.ID, fmt.Sprintf("%s.gif", i.ID), g.Pages)
		if err != nil {
			fmt.Println("转换到gif错误:", err)
			return
		}
	}

	//==================================================
	//12552116   在异世界当萝莉真祖.txt
	//12278185   不洁的茉莉.txt
	//series, Err := app.GetSeriesNovels("12278185", true)
	//if Err.Err != nil {
	//	fmt.Println("错误：", Err.Err, Err.Code)
	//	return
	//}
	//fmt.Println(series)
	//====================================
	//白夜幼女 68413373
	//百！病！不！侵！  17582659
	//userNovel, Err := app.GetUserNovels("684133738945")
	//if Err.Err != nil {
	//	fmt.Println(Err.Err, Err.Code)
	//	return
	//}
	//fmt.Println(userNovel.Novels)
	//====================================
	// app.SetHttpTimeout(5000)
	/*	err := app.SetProxy("socks5://10.10.10.2:1080")
		if err != nil {
			fmt.Println("设置代理错误2", err)
			return
		}*/
	//=====================================
	//获取token
	//err := app.GetAccessToken()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(app.AccessToken)
	//===================================
	//	在邪神祭祀中被变成美少女容器的话，人生大概就完蛋了吧！？  23230041
	//地痞流氓的改邪归正 | E_Vao #pixiv 21075769
	//content, Err := app.GetNovelContent("21075769")
	//if Err.Err != nil {
	//	fmt.Println("错误：", Err.Err, Err.Code)
	//	return
	//}
	//fmt.Println(content.Images)
	//==============================
	//userInfo, Err := app.GetUserDetail("17582659")
	//if Err.Err != nil {
	//	fmt.Println("错误：", Err.Err, Err.Code)
	//	return
	//}
	//fmt.Println(userInfo)
	//novelInfo, Err := app.GetNovelDetail("23230041")
	//if Err.Err != nil {
	//	fmt.Println("错误：", Err.Err, Err.Code)
	//	return
	//}
	//fmt.Println(novelInfo)

	//data, err := app.GetIllustByte("https://i.pximg.net/img-original/img/2024/07/25/03/09/47/120858495_p0.png")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//os.WriteFile("./1.png", data, 0644)

}
