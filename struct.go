package pixiv

type Error struct {
	// 403 速率过快
	// 404不存在
	// 400 认证失效或无效的令牌
	Code int    // http状态码
	Body string // http内容
	Err  error  //错误信息
}

// 用户信息
type UserInfo struct {
	UserID      int64  //用户id
	UserName    string //用户名
	Caption     string //用户简介
	Account     string //账户名
	TotalNovels int64  //用户小说数量
}

// 用户小说
type UserNovel struct {
	UserName string // 用户名
	Account  string //用户账号
	Count    int64  //用户小说数量
	Novels   []NovelInfo
}

// 系列小说
type SeriesNovel struct {
	UserID  string
	Title   string // 系列标题
	Caption string //系列简介
	Count   int64  //系列小说数量
	Novels  []NovelInfo
}

// 单篇小说信息
type NovelInfo struct {
	UserID   string   //用户id
	ID       string   // 小说id
	Title    string   // 小说标题
	Length   string   // 小说字数
	Caption  string   // 小说简介
	Date     string   // 发布时间
	Tags     []string // 标签
	SeriesID string   //系列id(如果是在系列中)
}

// 小说正文
type NovelContent struct {
	RawnContent string            // 小说json正文
	Content     string            // 小说正文
	Images      map[string]string // 小说插图
}
