package crawl

import (
	"encoding/json"
	"net/http"
	"time"
)

const (
	HeaderUserAgent       = "User-Agent"
	HeaderAccept          = "Accept"
	HeaderAcceptLanguage  = "Accept-Language"
	HeaderContentType     = "Content-Type"
	HeaderAcceptEncoding  = "Accept-Encoding"
	HeaderContentEncoding = "Content-Encoding"
	ParamCharset          = "charset"
)

type Option func(req *http.Request) (err error)

type Crawl struct {
	options []Option
	chains  []ProcessMw
	cOpts   []ClientOption
}

func WindowsEdge() *Crawl {
	return (&Crawl{}).With(HeaderWindowsEdge).Chains(AutoDecode, ToUTF8, StatusOK).Client(Timeout(time.Second * 15))
}

func (c *Crawl) Reset() *Crawl {
	c.options = c.options[:0]
	c.chains = c.chains[:0]
	return c.With(HeaderWindowsEdge)
}

func (c *Crawl) With(rOpts ...Option) *Crawl {
	c.options = append(c.options, rOpts...)
	return c
}

func (c *Crawl) Client(options ...ClientOption) *Crawl {
	c.cOpts = append(c.cOpts, options...)
	return c
}

func (c *Crawl) Chains(mws ...ProcessMw) *Crawl {
	c.chains = append(c.chains, mws...)
	return c
}

func (c *Crawl) Process(process ProcessOption) (err error) {
	req := &http.Request{
		Method:     http.MethodGet,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
	}

	for _, apply := range c.options {
		if err = apply(req); err != nil {
			return
		}
	}

	client := new(http.Client)
	client.Transport = retry(0, time.Second, time.Minute, 0, defaultTransport())
	for _, apply := range c.cOpts {
		if err = apply(client); err != nil {
			return
		}
	}

	var resp *http.Response
	if resp, err = client.Do(req); err != nil {
		return
	}

	body := resp.Body
	defer body.Close()

	for _, chain := range c.chains {
		process = chain(process)
	}

	err = process(resp)
	return
}

func (c *Crawl) HTML(process func(dom *DOM) error) (err error) {
	return c.Process(HTML(process))
}

func (c *Crawl) JSON(value any) (err error) {
	return c.Process(func(resp *http.Response) (err error) {
		err = json.NewDecoder(resp.Body).Decode(value)
		return
	})
}

func (c *Crawl) Download(saveTo string) (err error) {
	return c.Process(Download(saveTo))
}
