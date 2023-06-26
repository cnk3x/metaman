package crawl

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type ClientOption func(client *http.Client) (err error)

func NewClient() *http.Client {
	return &http.Client{Transport: retry(5, time.Second, time.Minute, 0, defaultTransport())}
}

var defaultJar, _ = cookiejar.New(nil)

// CookieEnabled enabled cookie for client
func CookieEnabled(enabled bool) ClientOption {
	return func(client *http.Client) (err error) {
		if enabled {
			client.Jar = defaultJar
		} else {
			client.Jar = nil
		}
		return
	}
}

func CookieJar(jar http.CookieJar) ClientOption {
	return func(client *http.Client) (err error) {
		client.Jar = jar
		return
	}
}

func Proxy(proxy func() *url.URL) ClientOption {
	return func(client *http.Client) (err error) {
		ht := client.Transport.(*retryRoundtripper).Next.(*http.Transport)
		ht.Proxy = func(r *http.Request) (*url.URL, error) { return proxy(), nil }
		if ht.TLSClientConfig == nil {
			ht.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		} else {
			ht.TLSClientConfig.InsecureSkipVerify = true
		}
		return
	}
}

// Timeout timeout
func Timeout(timeout time.Duration) ClientOption {
	return func(client *http.Client) (err error) {
		client.Timeout = timeout
		return
	}
}

// Insecure tls insecure skip verify
func Insecure(enabled bool) ClientOption {
	return func(client *http.Client) (err error) {
		ht := client.Transport.(*retryRoundtripper).Next.(*http.Transport)
		if ht.TLSClientConfig == nil {
			ht.TLSClientConfig = &tls.Config{InsecureSkipVerify: enabled}
		} else {
			ht.TLSClientConfig.InsecureSkipVerify = enabled
		}
		return
	}
}

func Retry(max int) ClientOption {
	return func(client *http.Client) (err error) {
		client.Transport.(*retryRoundtripper).MaxRetryCount = max
		return
	}
}

func defaultTransport() *http.Transport {
	var dialer = net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second}
	return &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}
