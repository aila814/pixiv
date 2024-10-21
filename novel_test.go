package pixiv

import (
	"testing"
)

var token = "Lq8TIOzdmmXFaXd4L2fUvIKKdE47GLAJ5Wdbq3VIoSo"

func TestPixiv(t *testing.T) {
	//var a NovelInfo

	app := NewApp()

	app.SetProxy("socks5://10.10.10.2:1080")
	app.SetRefreshToken("JMfu2g_lzwHS7ojLnUA3mifsBRPQLgbfqQLZ9W1M73E")
	app.SetAccessToken(token)
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
