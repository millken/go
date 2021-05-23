package ippool

import "time"

func timeIt(f func() error) (int64, error) {
	startAt := time.Now()
	err := f()
	endAt := time.Now()
	return endAt.UnixNano() - startAt.UnixNano(), err
}
