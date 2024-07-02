package share

import (
	"container/list"
	"github.com/BAN1ce/skyTree/pkg/packet"
	"sync"
)

type Queue struct {
	mux   sync.RWMutex
	queue list.List // *packet.Message
}

func NewQueue() *Queue {
	return &Queue{
		queue: list.List{},
	}
}

func (t *Queue) PutBackMessage(message []*packet.Message) {
	t.mux.Lock()
	defer t.mux.Unlock()
	for _, v := range message {
		t.queue.PushFront(v)
	}
}

// AppendMessage appends a message to the queue
func (t *Queue) AppendMessage(message []*packet.Message) {
	t.mux.Lock()
	defer t.mux.Unlock()

	for i := 0; i < len(message); i++ {
		t.queue.PushBack(message[i])
	}
}

func (t *Queue) PopQueue(n int) []*packet.Message {
	t.mux.Lock()
	defer t.mux.Unlock()

	var message []*packet.Message

	for i := 0; i < n; i++ {
		if t.queue.Len() == 0 {
			break
		}
		first := t.queue.Front()
		message = append(message, first.Value.(*packet.Message))
		t.queue.Remove(first)
	}
	return message
}

func (t *Queue) Len() int {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.queue.Len()
}
