package crawl

import (
	"bytes"
	"crypto/x509"
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func retry(max int, minWait, maxWait, maxJitter time.Duration, next http.RoundTripper) *retryRoundtripper {
	if next == nil {
		next = http.DefaultTransport
	}
	return &retryRoundtripper{
		Next:          next,
		MaxRetryCount: max,
		minWait:       minWait,
		maxWait:       maxWait,
		maxJitter:     maxJitter,
	}
}

// retryRoundtripper is the roundtripper that will wrap around the actual http.Transport roundtripper
// to enrich the http client with retry functionality.
type retryRoundtripper struct {
	Next          http.RoundTripper
	MaxRetryCount int
	minWait       time.Duration
	maxWait       time.Duration
	maxJitter     time.Duration
}

// RoundTrip implements the actual roundtripper interface (http.RoundTripper).
func (r *retryRoundtripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var (
		resp         *http.Response
		err          error
		dataBuffer   *bytes.Reader
		statusCode   int
		attemptCount = 1
		maxAttempts  = r.MaxRetryCount + 1
	)

	for {
		statusCode = 0

		// if request provides GetBody() we use it as Body,
		// because GetBody can be retrieved arbitrary times for retry
		if req.GetBody != nil {
			bodyReadCloser, _ := req.GetBody()
			req.Body = bodyReadCloser
		} else if req.Body != nil {
			// we need to store the complete body, since we need to reset it if a retry happens
			// but: not very efficient because:
			// a) huge stream data size will all be buffered completely in the memory
			//    imagine: 1GB stream data would work efficiently with io.Copy, but has to be buffered completely in memory
			// b) unnecessary if first attempt succeeds
			// a solution would be to at least support more types for GetBody()

			// store it for the first time
			if dataBuffer == nil {
				data, err := io.ReadAll(req.Body)
				req.Body.Close()
				if err != nil {
					return nil, err
				}
				dataBuffer = bytes.NewReader(data)
				req.ContentLength = int64(dataBuffer.Len())
				req.Body = io.NopCloser(dataBuffer)
			}

			// reset the request body
			dataBuffer.Seek(0, io.SeekStart)
		}

		if resp, err = r.Next.RoundTrip(req); resp != nil {
			statusCode = resp.StatusCode
		}

		if !shouldRetry(statusCode, err) {
			return resp, err
		}

		backoff := r.exponentialBackoff(attemptCount)

		// no need to wait if we do not have retries left
		attemptCount++
		if attemptCount > maxAttempts {
			break
		}

		// we won't need the response anymore, drain (up to a maximum) and close it
		drainAndCloseBody(resp, 16384)

		timer := time.NewTimer(backoff)
		select {
		case <-req.Context().Done():
			// context was canceled, return context error
			return nil, req.Context().Err()
		case <-timer.C:
		}
	}

	// no more attempts, return the last response / error
	return resp, err
}

func shouldRetry(statusCode int, err error) bool {
	// check if error is of type temporary
	t, ok := err.(interface{ Temporary() bool })
	if ok && t.Temporary() {
		return true
	}

	// we cannot know all errors, so we filter errors that should NOT be retried
	switch e := err.(type) {
	case *url.Error:
		switch {
		case
			e.Op == "parse",
			strings.Contains(e.Err.Error(), "stopped after"),
			strings.Contains(e.Error(), "unsupported protocol scheme"),
			strings.Contains(e.Error(), "no Host in request URL"):
			return false
		}
		// check inner error of url.Error
		switch e.Err.(type) {
		case // this errors will not likely change when retrying
			x509.UnknownAuthorityError,
			x509.CertificateInvalidError,
			x509.ConstraintViolationError:
			return false
		}
	case error: // generic error, check for strings if nothing found, retry
		return true
	case nil: // no error, continue
	}

	// most of the codes should not be retried, so we filter status codes that SHOULD be retried
	switch statusCode {
	case // status codes that should be retried
		http.StatusRequestTimeout,
		http.StatusConflict,
		http.StatusLocked,
		http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
		http.StatusInsufficientStorage:
		return true
	case 0: // means we did not get a response. we need to retry
		return true
	default: // on all other status codes we should not retry (e.g. 200, 401 etc.)
		return false
	}
}

// ExponentialBackoff increases the backoff exponentially by multiplying the minWait with 2^attemptCount
//
// minWait: the initial backoff
//
// maxWait: sets an upper bound on the maximum time to wait between two requests. set to 0 for no upper bound
//
// maxJitter: random interval [0, maxJitter) added to the exponential backoff
//
// Example:
//
//	minWait = 1 * time.Seconds
//	maxWait = 60 * time.Seconds
//	maxJitter = 0 * time.Seconds
//
//	Backoff will be: 1, 2, 4, 8, 16, 32, 60, 60, ...
func (r *retryRoundtripper) exponentialBackoff(attemptCount int) time.Duration {
	if r.minWait < 0 {
		r.minWait = 0
	}
	if r.maxJitter < 0 {
		r.maxJitter = 0
	}
	if r.maxWait < r.minWait {
		r.maxWait = 0
	}
	nextWait := time.Duration(math.Pow(2, float64(attemptCount-1)))*r.minWait + randJitter(r.maxJitter)
	if r.maxWait > 0 {
		return minDuration(nextWait, r.maxWait)
	}
	return nextWait
}

func drainAndCloseBody(resp *http.Response, maxBytes int64) {
	if resp != nil {
		io.CopyN(io.Discard, resp.Body, maxBytes)
		resp.Body.Close()
	}
}

// minDuration returns the minimum of two durations
func minDuration(duration1 time.Duration, duration2 time.Duration) time.Duration {
	if duration1 < duration2 {
		return duration1
	}
	return duration2
}

// randJitter returns a random duration in the interval [0, maxJitter)
//
// if maxJitter is <= 0, a duration of 0 is returned
func randJitter(maxJitter time.Duration) time.Duration {
	if maxJitter <= 0 {
		return 0
	}

	return time.Duration(rand.Intn(int(maxJitter)))
}
