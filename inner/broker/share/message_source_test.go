package share

import (
	"context"
	"fmt"
	"github.com/BAN1ce/skyTree/pkg/mocks/broker"
	"github.com/BAN1ce/skyTree/pkg/packet"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

func TestMessageSource_fillShareTopic(t *testing.T) {
	var (
		fullTopic     = "fullTopic"
		messageSource = &MessageSource{
			fullShareTopic: fullTopic,
		}
		message []*packet.Message
	)

	for i := 0; i < 10; i++ {
		message = append(message, &packet.Message{
			ClientID: "test",
		})
	}

	messageSource.fillShareTopic(message)

	for _, v := range message {
		if v.ShareTopic != fullTopic {
			t.Errorf("topic error")
		}
	}

}

func TestMessageSource_nextMessage(t *testing.T) {
	var (
		mockCtrl              = gomock.NewController(t)
		topic                 = "topic"
		fullTopic             = "$share/fullTopic/" + topic
		mockMessageSource     = broker.NewMockMessageSource(mockCtrl)
		messageSource         = newMessageSource(fullTopic, topic, mockMessageSource)
		messageSourceResponse []*packet.Message
	)
	defer mockCtrl.Finish()

	for i := 0; i < 10; i++ {
		messageSourceResponse = append(messageSourceResponse, &packet.Message{
			ClientID: fmt.Sprintf("%d", i),
		})
	}

	mockMessageSource.EXPECT().NextMessages(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(messageSourceResponse, nil).AnyTimes()

	result, err := messageSource.nextMessage(context.TODO(), 10)
	if err != nil {
		t.Errorf("next message error")
	}
	for i := 0; i < 10; i++ {
		if result[i].ClientID != messageSourceResponse[i].ClientID {
			t.Errorf("message error")
		}
	}
}

func TestMessageSource_ListenMessage(t *testing.T) {
	var (
		mockCtrl              = gomock.NewController(t)
		topic                 = "topic"
		fullTopic             = "$share/fullTopic/" + topic
		mockMessageSource     = broker.NewMockMessageSource(mockCtrl)
		messageSource         = newMessageSource(fullTopic, topic, mockMessageSource)
		messageSourceResponse []*packet.Message
		result                []*packet.Message
		once                  bool
		responseChan          = make(chan *packet.Message, 10)
	)
	defer mockCtrl.Finish()

	for i := 0; i < 10; i++ {
		messageSourceResponse = append(messageSourceResponse, &packet.Message{
			ClientID: fmt.Sprintf("%d", i),
		})
	}

	mockMessageSource.EXPECT().NextMessages(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, n int, startMessageID string, include bool) ([]*packet.Message, error) {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		if !once {
			once = true
			return messageSourceResponse, nil
		} else {
			time.Sleep(1 * time.Second)
			return messageSourceResponse, nil
		}

	}).AnyTimes()

	err := messageSource.ListenMessage(context.TODO(), responseChan)
	if err != nil {
		t.Errorf("listen message error")
	}
	go func() {
		time.Sleep(2 * time.Second)
		messageSource.Close()

	}()
	for m := range responseChan {
		result = append(result, m)
	}

	if len(result) < 20 {
		t.Errorf("listen message error, got %d", len(result))
	} else {
		t.Log("listen message success, got message length", len(result))
	}

}
