package retry

import (
	"github.com/BAN1ce/skyTree/pkg/packet"
	"sync/atomic"
	"time"
)

type DelayTask struct {
	key    string
	delay  time.Duration
	data   interface{}
	circle atomic.Int64
}

func newDelayTask(key string, delay time.Duration, data interface{}) *DelayTask {
	return &DelayTask{
		data:  data,
		delay: delay,
		key:   key,
	}
}

type Task struct {
	Key       string
	Data      *packet.Message
	ClientID  string
	DelayTime time.Duration
}

func NewTask(Key string, Data *packet.Message, ClientID string, DelayTime time.Duration) *Task {
	return &Task{
		Key:       Key,
		Data:      Data,
		ClientID:  ClientID,
		DelayTime: DelayTime,
	}
}
