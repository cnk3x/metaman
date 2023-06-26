package crawl

import (
	"bufio"
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/andybalholm/brotli"
	"golang.org/x/text/encoding/htmlindex"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type ProcessOption func(resp *http.Response) (err error)
type ProcessMw func(next ProcessOption) ProcessOption
type DOM = goquery.Selection

func AutoDecode(next ProcessOption) ProcessOption {
	return func(resp *http.Response) (err error) {
		if !resp.Uncompressed {
			switch cEncoding := resp.Header.Get(HeaderContentEncoding); cEncoding {
			case "gzip":
				resp.Body, err = gzip.NewReader(resp.Body)
			case "deflate":
				resp.Body = flate.NewReader(resp.Body)
			case "br":
				resp.Body = io.NopCloser(brotli.NewReader(resp.Body))
			default:
				err = fmt.Errorf("unsupport encoding %q", cEncoding)
			}
			if err != nil {
				return
			}
			body := resp.Body
			defer body.Close()
			resp.Uncompressed = true
			resp.Header.Del(HeaderContentEncoding)
		}
		return next(resp)
	}
}

func ToUTF8(next ProcessOption) ProcessOption {
	return func(resp *http.Response) (err error) {
		if _, params, _ := mime.ParseMediaType(resp.Header.Get(HeaderContentType)); len(params) > 0 {
			if charset := params[ParamCharset]; charset != "" && charset != "UTF-8" {
				if codec, err := htmlindex.Get(charset); err == nil && codec != unicode.UTF8 {
					resp.Body = io.NopCloser(transform.NewReader(resp.Body, codec.NewDecoder()))
				}
			}
		}
		return next(resp)
	}
}

func StatusOK(next ProcessOption) ProcessOption {
	return func(resp *http.Response) (err error) {
		if resp.StatusCode != 200 {
			data, _ := io.ReadAll(resp.Body)
			err = &StatusError{StatusCode: resp.StatusCode, Data: data}
			return
		}
		return next(resp)
	}
}

func HTML(process func(dom *DOM) error) ProcessOption {
	return func(resp *http.Response) (err error) {
		var doc *goquery.Document
		if doc, err = goquery.NewDocumentFromReader(resp.Body); err != nil {
			return
		}
		err = process(doc.Selection)
		return
	}
}

func Download(saveTo string) ProcessOption {
	return func(resp *http.Response) (err error) {
		r := bufio.NewReaderSize(resp.Body, 32*1024)
		var fo *os.File
		fo, err = os.Create(saveTo)
		if err != nil {
			return
		}
		defer fo.Close()

		w := bufio.NewWriter(fo)
		defer w.Flush()
		_, err = io.Copy(w, r)
		return
	}
}
