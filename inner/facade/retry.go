package facade

import (
	"github.com/BAN1ce/skyTree/config"
	"github.com/BAN1ce/skyTree/inner/metric"
	"github.com/BAN1ce/skyTree/logger"
	"github.com/BAN1ce/skyTree/pkg/retry"
	"sync"
)

type RetrySchedule interface {
	Create(task *retry.Task) error
	Delete(key string)
}

var (
	oncePublishRetry sync.Once
	publishRetry     *PublishRetry
	pubRelRetry      RetrySchedule
	oncePubRelRetry  sync.Once
)

/*

 */
// ----------------------------------  PublishRetry ----------------------------------//
/**

 */

type PublishRetry struct {
	schedule *retry.DelayTaskSchedule
}

func InitPublishRetry(call, timeout func(task *retry.Task) error, option ...retry.Option) RetrySchedule {
	oncePublishRetry.Do(func() {
		publishRetry = newPublishRetry(call, timeout, option...)
		publishRetry.schedule.Start()
	})
	return publishRetry
}

func newPublishRetry(call, timeout func(t *retry.Task) error, option ...retry.Option) *PublishRetry {
	p := &PublishRetry{}
	p.schedule = retry.NewSchedule(config.GetRootContext(), call, timeout, option...)
	return p
}

func (p *PublishRetry) Retry(t *retry.Task) error {
	return nil

}

func (p *PublishRetry) Timeout(t *retry.Task) error {
	return nil
}

func GetPublishRetry() RetrySchedule {
	return publishRetry
}

func (p *PublishRetry) Create(task *retry.Task) error {

	metric.PublishRetryTaskCurrent.Inc()
	metric.PublishRetryTaskCount.WithLabelValues(task.Data.GetFullTopic()).Inc()
	return p.schedule.Create(task)
}

func (p *PublishRetry) Delete(key string) {
	metric.PublishRetryTaskCurrent.Add(-1)
	metric.PublishRetryDelete.Inc()
	logger.Logger.Debug().Msg("PublishRetry Delete")
	p.schedule.Delete(key)
}

//func DeletePublishRetryKey(key string) {
//	GetPubRelRetry().Delete(key)
//}
//
//func SinglePubRelRetry(ctx context.Context, option ...retry.Option) RetrySchedule {
//	oncePubRelRetry.Do(func() {
//		var s = retry.NewSchedule(ctx, option...)
//		s.Start()
//		pubRelRetry = s
//	})
//	return pubRelRetry
//}
//
//func GetPubRelRetry() RetrySchedule {
//	return SinglePubRelRetry(config.GetRootContext())
//}
