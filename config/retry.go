package config

import "time"

type Retry struct {
	MaxTimes int
	Timeout  int
	Interval int
}

func (r Retry) GetMaxTimes() int {
	return r.MaxTimes
}

func (r Retry) GetInterval() time.Duration {
	return time.Duration(r.Interval) * time.Second
}

func (r Retry) GetTimeout() time.Duration {
	return time.Duration(r.Timeout) * time.Second
}
