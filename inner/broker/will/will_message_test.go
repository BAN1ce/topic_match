package will

import (
	"github.com/BAN1ce/skyTree/inner/broker/state"
	"github.com/BAN1ce/skyTree/inner/mocks/event"
	"github.com/BAN1ce/skyTree/logger"
	session2 "github.com/BAN1ce/skyTree/pkg/broker/session"
	"github.com/BAN1ce/skyTree/pkg/mocks/broker"
	"github.com/BAN1ce/skyTree/pkg/mocks/session"
	"github.com/BAN1ce/skyTree/pkg/mocks/store"
	"github.com/golang/mock/gomock"
	"testing"
)

func init() {

	logger.Load()

}

func TestMessageMonitor_publishWill(t *testing.T) {
	var (
		mockCtl          = gomock.NewController(t)
		mockStore        = store.NewMockKeyStore(mockCtl)
		mockSubCenter    = broker.NewMockSubCenter(mockCtl)
		mockStoreWrapper = store.NewMockWrapper(mockCtl)
		mockSession      = session.NewMockManager(mockCtl)
		newState         = state.NewState(mockStore)
		eventDriver      = event.NewMockDriver(mockCtl)
		payload          = []byte("hello")
		matchTopic       = map[string]int32{"/will": 1, "/#": 2}
	)
	defer mockCtl.Finish()

	mockSubCenter.EXPECT().GetAllMatchTopics("/will").Return(matchTopic).AnyTimes()

	mockStoreWrapper.EXPECT().StorePublishPacket(matchTopic, gomock.Any()).Return("", nil).AnyTimes()

	eventDriver.EXPECT().Emit(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	monitor := NewWillMessageMonitor(eventDriver, newState, mockSubCenter, mockStoreWrapper, mockSession)

	monitor.publishWill(&session2.WillMessage{
		Topic:       "/will",
		QoS:         0,
		Property:    nil,
		Retain:      false,
		Payload:     payload,
		DelayTaskID: "",
		CreatedTime: "",
		ClientID:    "",
	})

}
