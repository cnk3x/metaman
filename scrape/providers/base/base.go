package base

type MovieSearcher interface {
	//排行榜
	Chart() (out []MovieSummary, err error)
	//搜索
	Search(keyword string, year int) (out []MovieSummary, err error)
	//详情
	Subject(id string) (out MovieDetail, err error)
	//影片图集
	Photos(id string) (out []Photo, err error)
	//下载图片
	DownloadPhoto(id string, uri string, saveTo string) (err error)
}
