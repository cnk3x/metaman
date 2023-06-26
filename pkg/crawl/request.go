package crawl

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func NewRequest() *http.Request {
	req := new(http.Request)
	req.Method = http.MethodGet
	req.Proto = "HTTP/1.1"
	req.ProtoMajor = 1
	req.ProtoMinor = 1
	req.Header = make(http.Header)
	return req
}

func Context(ctx context.Context) Option {
	return func(req *http.Request) (err error) {
		if ctx != nil {
			*req = *req.WithContext(ctx)
		}
		return nil
	}
}

func Url(u string) Option {
	return func(req *http.Request) (err error) {
		if req.URL, err = url.Parse(u); err != nil {
			return
		}
		req.Host = req.URL.Host
		return
	}
}

func Body(body io.Reader) func(req *http.Request) (err error) {
	return func(req *http.Request) (err error) {
		rc, ok := body.(io.ReadCloser)
		if !ok && body != nil {
			rc = io.NopCloser(body)
		}
		req.Body = rc
		if body != nil {
			switch v := body.(type) {
			case *bytes.Buffer:
				req.ContentLength = int64(v.Len())
				buf := v.Bytes()
				req.GetBody = func() (io.ReadCloser, error) {
					r := bytes.NewReader(buf)
					return io.NopCloser(r), nil
				}
			case *bytes.Reader:
				req.ContentLength = int64(v.Len())
				snapshot := *v
				req.GetBody = func() (io.ReadCloser, error) {
					r := snapshot
					return io.NopCloser(&r), nil
				}
			case *strings.Reader:
				req.ContentLength = int64(v.Len())
				snapshot := *v
				req.GetBody = func() (io.ReadCloser, error) {
					r := snapshot
					return io.NopCloser(&r), nil
				}
			default:
				// This is where we'd set it to -1 (at least
				// if body != NoBody) to mean unknown, but
				// that broke people during the Go 1.8 testing
				// period. People depend on it being 0 I
				// guess. Maybe retry later. See Issue 18117.
			}
			// For client requests, Request.ContentLength of 0
			// means either actually 0, or unknown. The only way
			// to explicitly say that the ContentLength is zero is
			// to set the Body to nil. But turns out too much code
			// depends on NewRequest returning a non-nil Body,
			// so we use a well-known ReadCloser variable instead
			// and have the http package also treat that sentinel
			// variable to mean explicitly zero.
			if req.GetBody != nil && req.ContentLength == 0 {
				req.Body = http.NoBody
				req.GetBody = func() (io.ReadCloser, error) { return http.NoBody, nil }
			}
		}
		return nil
	}
}

func FormBody(form url.Values) func(req *http.Request) (err error) {
	return Options(Body(strings.NewReader(form.Encode())), HeaderSet(HeaderContentType, ""))
}

func Options(options ...Option) Option {
	return func(req *http.Request) (err error) {
		return applyOptions(req, options...)
	}
}

func applyOptions(r *http.Request, options ...Option) (err error) {
	for _, apply := range options {
		if err = apply(r); err != nil {
			break
		}
	}
	return
}
