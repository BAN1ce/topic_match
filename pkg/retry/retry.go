package retry

import (
	"context"
	"github.com/BAN1ce/skyTree/logger"
	"time"
)

func NewSchedule(ctx context.Context, call, timeout func(task *Task) error, options ...Option) *DelayTaskSchedule {
	s := NewDelayTaskSchedule(ctx, call, timeout, options...)
	return s
}

type Option func(*DelayTaskSchedule)

func WithInterval(interval time.Duration) Option {
	return func(r *DelayTaskSchedule) {
		r.interval = interval
	}
}
func WithSlotNum(slotNum int) Option {
	return func(r *DelayTaskSchedule) {
		r.slotNum = slotNum
	}
}

func WithName(name string) Option {
	return func(schedule *DelayTaskSchedule) {
		schedule.scheduleName = name
	}
}

type DelayTaskSchedule struct {
	scheduleName string
	interval     time.Duration
	slotNum      int
	tw           *TimeWheel
	timeoutFunc  func(t *Task) error
	callFunc     func(t *Task) error
}

func NewDelayTaskSchedule(ctx context.Context, callFunc, timeoutFunc func(t *Task) error, options ...Option) *DelayTaskSchedule {
	r := &DelayTaskSchedule{
		callFunc:     callFunc,
		timeoutFunc:  timeoutFunc,
		interval:     1 * time.Second,
		slotNum:      3000,
		scheduleName: "publish retry",
	}
	for _, option := range options {
		option(r)
	}
	if r.slotNum == 0 {
		r.slotNum = 1000
	}
	if r.interval == 0 {
		r.interval = time.Second
	}
	r.tw = NewTimeWheel(ctx, r.interval, r.slotNum, func(key string, data interface{}) {
		if t, ok := data.(*Task); ok {

			if t.Data.RetryInfo != nil {
				if t.Data.RetryInfo.IsTimeout() {
					if r.timeoutFunc != nil {
						if err := r.timeoutFunc(t); err != nil {
							logger.Logger.Error().Err(err).Msg("timeoutFunc error")
						}
					}
					return
				}
			}

			if err := r.callFunc(t); err != nil {
				logger.Logger.Error().Err(err).Msg("callFunc error")
			}
			//if err := r.tw.CreateTask(key, t.IntervalTime, data); err != nil {
			//	logger.Logger.Debug().Msg("recreate delay task for retry again")
			//}
			return

		}

	})
	return r
}

func (d *DelayTaskSchedule) Start() {
	go d.tw.Run()
}

func (d *DelayTaskSchedule) Create(task *Task) error {
	logger.Logger.Debug().Str("key", task.Key).Any("data", task.Data).Msg("DelayTaskSchedule Create")
	return d.tw.CreateTask(task.Key, task.DelayTime, task)
}

func (d *DelayTaskSchedule) Delete(key string) {
	logger.Logger.Debug().Str("key", key).Msg("DelayTaskSchedule Delete")
	d.tw.DeleteTask(key)
}
