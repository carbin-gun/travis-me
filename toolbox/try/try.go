package try

import "errors"

var MaxTries = 10
var ErrMaxTriesReached = errors.New("exceeded retry limit")

type TryFunc func(attempt int) (needRetry bool, err error)

// Do keeps trying the function until the needRetry returns false, or some error is returned.
//
func Do(fn TryFunc) error {
	var err error
	var needRetry bool
	var attempt = 1
	for {
		needRetry, err = fn(attempt)
		if !needRetry || err == nil {
			break
		}
		attempt++
		if attempt > MaxTries {
			return ErrMaxTriesReached
		}
	}
	return err
}

// IsMaxRetries checks whether the error is due to hitting the maximum number of retries or not.
func IsMaxTry(err error) bool {
	return err == ErrMaxTriesReached
}
