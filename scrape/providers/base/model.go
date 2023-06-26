package base

// 影片摘要
type MovieSummary struct {
	ID     string `json:"id,omitempty"`     //编号
	Title  string `json:"title,omitempty"`  //标题
	Year   int    `json:"year,omitempty"`   //年份
	Poster string `json:"poster,omitempty"` //封面
}

// 影片详情
type MovieDetail struct {
	MovieSummary
	Summary     string           `json:"summary,omitempty"`     //简介
	Date        string           `json:"date,omitempty"`        //上映日期
	Celebrities []MovieCelebrity `json:"celebrities,omitempty"` //演职员
	Rating      float64          `json:"rating,omitempty"`      //评分
	Genre       []string         `json:"genre,omitempty"`       //类型
	Country     string           `json:"country,omitempty"`     //国家
	Language    string           `json:"language,omitempty"`    //语言
	Alias       []string         `json:"alias,omitempty"`       //别名
	Duration    int              `json:"duration,omitempty"`    //片长
	WebSite     string           `json:"web_site,omitempty"`    //网站
	Douban      string           `json:"douban,omitempty"`      //豆瓣编号
	IMDB        string           `json:"imdb,omitempty"`        //IMDB
	TMDB        string           `json:"tmdb,omitempty"`        //TMDB
}

// 演职员
type MovieCelebrity struct {
	ID     string `json:"id,omitempty"`     //编号
	Name   string `json:"name,omitempty"`   //姓名
	Role   string `json:"role,omitempty"`   //角色
	Poster string `json:"poster,omitempty"` //图片
}

// 名人
type Celebrity struct {
	ID         string `json:"id,omitempty"`         //编号
	Name       string `json:"name,omitempty"`       //姓名
	Poster     string `json:"poster,omitempty"`     //图片
	Sex        string `json:"sex,omitempty"`        //性别
	Birthday   string `json:"birthday,omitempty"`   //出生日
	Birthplace string `json:"birthplace,omitempty"` //出生地
	Job        string `json:"job,omitempty"`        //职业
	Alias      string `json:"alias,omitempty"`      //别名
	Summary    string `json:"summary,omitempty"`    //简介
	Douban     string `json:"douban,omitempty"`     //豆瓣编号
	IMDB       string `json:"imdb,omitempty"`       //IMDB
	TMDB       string `json:"tmdb,omitempty"`       //TMDB
}

// https://movie.douban.com/subject/35964279/photos?type=R

// 照片
type Photo struct {
	ID      string `json:"id,omitempty"`       //编号
	MovieID string `json:"movie_id,omitempty"` //所属影片编号
	Type    string `json:"type,omitempty"`     //C: 剧照, R: 海报
	ImgUrl  string `json:"img_url,omitempty"`  //照片链接
	Size    string `json:"size,omitempty"`     //照片尺寸
	Title   string `json:"title,omitempty"`    //照片标题
	Douban  string `json:"douban,omitempty"`   //豆瓣编号
	IMDB    string `json:"imdb,omitempty"`     //IMDB
	TMDB    string `json:"tmdb,omitempty"`     //TMDB
}
