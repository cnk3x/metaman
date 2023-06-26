package crawl

import "net/http"

type StatusError struct {
	StatusCode int
	Data       []byte
}

func (e *StatusError) Error() string {
	return http.StatusText(e.StatusCode)
}
