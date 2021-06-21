package limiter

import "time"

var li = time.NewTicker(2 * time.Second).C

func GetLimiter() <-chan time.Time {
	return li
}
