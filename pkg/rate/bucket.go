package rate

import (
	"context"
	"github.com/BAN1ce/skyTree/logger"
)

type Bucket struct {
	ch  chan struct{}
	num int
}

// NewBucket create a bucket with num tokens
// if num <= 0, bucket is unlimited
func NewBucket(num int) *Bucket {
	var (
		b = &Bucket{
			num: num,
		}
	)
	if num <= 0 {
		b.ch = make(chan struct{})
		close(b.ch)
		return b
	}
	b.ch = make(chan struct{}, num)
	for i := 0; i < num; i++ {
		b.ch <- struct{}{}
	}
	return b
}

func (b *Bucket) GetToken(ctx context.Context) {
	select {
	case <-ctx.Done():
		logger.Logger.Debug().Msg("context done")
	case <-b.ch:

	}

}

// PutToken put a token into bucket, this operation is not concurrent safe
func (b *Bucket) PutToken() {
	if b.num <= 0 {
		return
	}
	select {
	case b.ch <- struct{}{}:
	default:
		logger.Logger.Error().Int("bucket count", len(b.ch)).Msg("put token into bucket failed, bucket is full, token will be dropped. It's abnormal, please check the code.")
	}
}
