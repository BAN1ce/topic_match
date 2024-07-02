package share

import (
	"container/list"
	"fmt"
	"github.com/BAN1ce/skyTree/pkg/packet"
	"reflect"
	"testing"
)

func TestNewQueue(t *testing.T) {
	tests := []struct {
		name string
		want *Queue
	}{
		{
			name: "new",
			want: &Queue{
				queue: list.List{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewQueue(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQueue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueue(t *testing.T) {
	var (
		queue   = NewQueue()
		message []*packet.Message
	)

	for i := 0; i < 10; i++ {
		message = append(message, &packet.Message{
			ClientID: fmt.Sprintf("%d", i),
		})
	}

	queue.AppendMessage(message)

	if queue.Len() != 10 {
		t.Errorf("queue len error")
	}

	if len(queue.PopQueue(5)) != 5 {
		t.Errorf("queue pop error")
	}

	if queue.Len() != 5 {
		t.Errorf("queue len error")
	}

	gotMessage := queue.PopQueue(5)

	for i := 0; i < 5; i++ {
		if gotMessage[i].ClientID != fmt.Sprintf("%d", i+5) {
			t.Errorf("queue pop error")
		}
	}
}
