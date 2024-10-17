package pixiv

// 用户信息
type UserInfo struct {
	UserID      int64  //用户id
	UserName    string //用户名
	Caption     string //用户简介
	Account     string //账户名
	TotalNovels int64  //用户小说数量
}

// 单篇小说信息
type NovelInfo struct {
	// 小说id
	ID       string   // 小说id
	Title    string   // 小说标题
	Length   string   // 小说字数
	Caption  string   // 小说简介
	Date     string   // 发布时间
	Tags     []string // 标签
	SeriesID string   //系列id(如果是在系列中)
}
