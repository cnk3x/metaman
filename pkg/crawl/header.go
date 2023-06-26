package crawl

import "net/http"

func HeaderSet(name, value string) Option {
	return func(req *http.Request) (err error) {
		req.Header.Set(name, value)
		return
	}
}

func HeaderWindowsEdge(req *http.Request) (err error) {
	req.Header.Set(HeaderUserAgent, `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.58`)
	req.Header.Set(HeaderAccept, `text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7`)
	req.Header.Set(HeaderAcceptEncoding, `gzip, deflate, br`)
	req.Header.Set(HeaderAcceptLanguage, `zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6,zh-TW;q=0.5`)
	req.Header.Set("Dnt", "1")
	req.Header.Set("Connection", "close")
	return
}

func NoCache(req *http.Request) (err error) {
	req.Header.Set(`Pragma`, `no-cache`)
	req.Header.Set(`Cache-Control`, `no-cache`)
	return
}

func Referer(referer string) Option {
	return HeaderSet("Referer", referer)
}
