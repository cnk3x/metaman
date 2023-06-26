package douban

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/cnk3x/metaman/pkg/crawl"
	"github.com/cnk3x/metaman/pkg/strs"
	"github.com/cnk3x/metaman/scrape/providers/base"

	"github.com/PuerkitoBio/goquery"
)

type douban struct {
}

// Chart implements base.MovieSearcher.
func (*douban) Chart() (out []base.MovieSummary, err error) {
	url := crawl.Url("https://movie.douban.com/chart")
	err = crawl.WindowsEdge().With(url).HTML(func(dom *crawl.DOM) (err error) {
		root := dom.Find(`#content div.article table`)
		root.Each(func(i int, s *crawl.DOM) {
			var it base.MovieSummary
			itemTr := s.Find("tr.item")
			nbg := itemTr.Find("a.nbg")
			it.ID = strings.TrimSuffix(strings.TrimPrefix(nbg.AttrOr("href", ""), `https://movie.douban.com/subject/`), "/")
			it.Poster = nbg.Find("img").AttrOr("src", "")
			it.Poster = strings.ReplaceAll(it.Poster, `s_ratio_poster`, `m_ratio_poster`)
			it.Title = nbg.AttrOr("title", "")
			it.Year, _ = strconv.Atoi(strs.Find(strs.Text(itemTr.Find("td:nth-of-type(2)")), `(\d{4})-\d{2}-\d{2}`, "$1"))
			out = append(out, it)
		})
		return
	})
	return
}

// Search implements base.MovieSearcher.
func (*douban) Search(q string, year int) (out []base.MovieSummary, err error) {
	reqUri := crawl.Url(fmt.Sprintf(`https://www.douban.com/search?cat=1002&q=%s`, url.QueryEscape(q)))
	// reqUri := crawl.Url(fmt.Sprintf(`https://search.douban.com/movie/subject_search?search_text=%s&cat=1002`, url.QueryEscape(q)))
	ref := crawl.Referer("https://movie.douban.com/")
	err = crawl.WindowsEdge().With(reqUri, ref).HTML(func(dom *crawl.DOM) (err error) {
		resultEl := dom.Find(`div.result-list .result`)
		resultEl.Each(func(i int, s *crawl.DOM) {
			var it base.MovieSummary
			it.ID = strs.Sub(s.Find("a.nbg").AttrOr("onclick", ""), "sid: ", ",")
			it.Poster = s.Find("a.nbg>img").AttrOr("src", "")
			it.Poster = strings.ReplaceAll(it.Poster, `s_ratio_poster`, `m_ratio_poster`)
			it.Title = strs.Text(s.Find("h3>a"))
			it.Year, _ = strconv.Atoi(strs.Find(strs.Text(s.Find(`span.subject-cast`)), `\d{4}`, ""))
			out = append(out, it)
		})
		return
	})
	return
}

// Subject implements base.MovieSearcher.
func (*douban) Subject(id string) (out base.MovieDetail, err error) {
	url := crawl.Url(fmt.Sprintf(`https://movie.douban.com/subject/%s/`, id))
	err = crawl.WindowsEdge().With(url).HTML(func(dom *goquery.Selection) error {
		info := dom.Find("#info")
		out.WebSite = strs.Text(info.Find(`span:contains(官方网站)~a`))
		out.Date = strs.Text(info.Find(`[property="v:initialReleaseDate"]`))
		out.Duration, _ = strconv.Atoi(strings.TrimSuffix(strs.Text(info.Find(`[property="v:runtime"]`)), "分钟"))

		cText := strs.Text(info)
		out.IMDB = strs.Find(cText, `IMDb: ([^\s]+)`, "$1")
		out.Country = strs.Find(cText, `国家[^\s]*: ([^\s]+)`, "$1")
		out.Language = strs.Find(cText, `语言: ([^\s]+)`, "$1")
		out.Alias = strings.Split(strs.Find(cText, `又名: ([^\s]+)`, `$1`), "/")

		h1 := dom.Find(`#content>h1`)
		out.Title = strs.Text(h1.Find(`[property="v:itemreviewed"]`))
		out.Year, _ = strconv.Atoi(strs.Unwrap(strs.Text(h1.Find(`span.year`)), "(", ")"))

		out.Poster = dom.Find("#mainpic img").AttrOr("src", "")
		out.Poster = strings.ReplaceAll(out.Poster, `s_ratio_poster`, `m_ratio_poster`)

		intra := dom.Find(`#link-report-intra`)
		out.Summary = strs.Text(intra.Find(`span.all.hidden`))
		if out.Summary == "" {
			out.Summary = strs.Text(intra.Find(`[property="v:summary"]`))
		}

		dom.Find(`#info [property="v:genre"]`).Each(func(i int, s *goquery.Selection) {
			out.Genre = append(out.Genre, strs.Text(s))
		})

		dom.Find(`#info [rel="v:directedBy"]`).Each(func(i int, s *goquery.Selection) {
			out.Celebrities = append(out.Celebrities, base.MovieCelebrity{
				ID:   strs.Unwrap(s.AttrOr("href", ""), "/celebrity/", "/"),
				Name: strs.Text(s),
				Role: "director",
			})
		})

		dom.Find(`#info span:contains(编剧)~span.attrs > a`).Each(func(i int, s *goquery.Selection) {
			out.Celebrities = append(out.Celebrities, base.MovieCelebrity{
				ID:   strs.Unwrap(s.AttrOr("href", ""), "/celebrity/", "/"),
				Name: strs.Text(s),
				Role: "author",
			})
		})

		dom.Find(`#info [rel="v:starring"]`).Each(func(i int, s *goquery.Selection) {
			out.Celebrities = append(out.Celebrities, base.MovieCelebrity{
				ID:   strs.Unwrap(s.AttrOr("href", ""), "/celebrity/", "/"),
				Name: strs.Text(s),
				Role: "starring",
			})
		})

		out.Rating, _ = strconv.ParseFloat(strs.Text(dom.Find(`#interest_sectl [property="v:average"]`)), 32)
		return nil
	})
	return
}

// Photos implements base.MovieSearcher.
func (*douban) Photos(id string) (out []base.Photo, err error) {
	//官方剧照
	//https://movie.douban.com/subject/35964279/photos?type=S&subtype=o
	//海报
	//https://movie.douban.com/subject/35964279/photos?type=R&subtype=a

	url := crawl.Url(fmt.Sprintf(`https://movie.douban.com/subject/%s/photos?type=R&subtype=o`, id))
	err = crawl.WindowsEdge().With(url).HTML(func(dom *crawl.DOM) (err error) {
		dom.Find("ul.poster-col3>li[data-id]").Each(func(i int, s *crawl.DOM) {
			s.Find("div.name").Children().Remove()
			it := base.Photo{
				ID:      s.AttrOr("data-id", ""),
				MovieID: id,
				Type:    "R",
				ImgUrl:  s.Find("div.cover img").AttrOr("src", ""),
				Size:    strs.Text(s.Find("div.prop")),
				Title:   strs.Text(s.Find("div.name")),
			}
			out = append(out, it)
		})
		return
	})
	if err != nil {
		return
	}

	url = crawl.Url(fmt.Sprintf(`https://movie.douban.com/subject/%s/photos?type=S&subtype=o`, id))
	err = crawl.WindowsEdge().With(url).HTML(func(dom *crawl.DOM) (err error) {
		dom.Find("ul.poster-col3>li[data-id]").Each(func(i int, s *crawl.DOM) {
			s.Find("div.name").Children().Remove()
			it := base.Photo{
				ID:      s.AttrOr("data-id", ""),
				MovieID: id,
				Type:    "S",
				ImgUrl:  s.Find("div.cover img").AttrOr("src", ""),
				Size:    strs.Text(s.Find("div.prop")),
				Title:   strs.Text(s.Find("div.name")),
			}
			out = append(out, it)
		})
		return
	})

	return
}

// DownloadPhoto implements base.MovieSearcher.
func (*douban) DownloadPhoto(id string, uri string, saveTo string) (err error) {
	return
}

func (*douban) Suggest(q string) {
	//https://movie.douban.com/j/subject_suggest?q=%E6%83%8A%E5%A4%A9%E8%90%A5%E6%95%91
	// Accept: */*
	// Accept-Encoding: gzip, deflate, br
	// Accept-Language: zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6,zh-TW;q=0.5
	// Connection: keep-alive
	// Dnt: 1
	// Host: movie.douban.com
	// Referer: https://movie.douban.com/
	// User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.58
	// X-Requested-With: XMLHttpRequest

	// 	[
	//     {
	//         "episode": "",
	//         "img": "https://img1.doubanio.com\/view\/photo\/s_ratio_poster\/public\/p2892053237.jpg",
	//         "title": "惊天营救2",
	//         "url": "https:\/\/movie.douban.com\/subject\/35056376\/?suggest=%E6%83%8A%E5%A4%A9%E8%90%A5%E6%95%91",
	//         "type": "movie",
	//         "year": "2023",
	//         "sub_title": "Extraction 2",
	//         "id": "35056376"
	//     },
	//     {
	//         "episode": "",
	//         "img": "https://img9.doubanio.com\/view\/photo\/s_ratio_poster\/public\/p2594557845.jpg",
	//         "title": "惊天营救",
	//         "url": "https:\/\/movie.douban.com\/subject\/30314127\/?suggest=%E6%83%8A%E5%A4%A9%E8%90%A5%E6%95%91",
	//         "type": "movie",
	//         "year": "2020",
	//         "sub_title": "Extraction",
	//         "id": "30314127"
	//     },
	//     {
	//         "episode": "",
	//         "img": "https://img9.doubanio.com\/view\/photo\/s_ratio_poster\/public\/p2890135815.jpg",
	//         "title": "惊天救援",
	//         "url": "https:\/\/movie.douban.com\/subject\/35215919\/?suggest=%E6%83%8A%E5%A4%A9%E8%90%A5%E6%95%91",
	//         "type": "movie",
	//         "year": "2023",
	//         "sub_title": "惊天救援",
	//         "id": "35215919"
	//     }
	// ]
}

var _ base.MovieSearcher = (*douban)(nil)

func MovieSearcher() base.MovieSearcher {
	return &douban{}
}
