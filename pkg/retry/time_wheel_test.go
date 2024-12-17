package retry

import (
	"context"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

func BenchmarkCreateTask(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tm := NewTimeWheel(ctx, time.Second, 1000, func(key string, data interface{}) {
		time.Sleep(10 * time.Millisecond)
	})
	tm.Run()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := tm.CreateTask(strconv.Itoa(i), time.Duration(i)*time.Second, nil); err != nil {
			b.Errorf("create task error: %v", err)
		}
	}
}

func TestDelayTask(t *testing.T) {
	var taskCount = 100
	var result atomic.Int32
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tm := NewTimeWheel(ctx, time.Millisecond*1000, 1000, func(key string, data interface{}) {
		result.Add(1)
	})

	tm.Run()

	for i := 0; i < taskCount; i++ {
		if err := tm.CreateTask(strconv.Itoa(i), time.Duration(1)*time.Second, nil); err != nil {
			t.Errorf("create task error: %v", err)
		}
	}

	time.Sleep(3 * time.Second)

	if result.Load() != 100 {
		t.Errorf("result is not 100")
	}
}
