package sub_topic

import (
	"github.com/BAN1ce/skyTree/pkg/broker/topic"
	"github.com/BAN1ce/skyTree/pkg/mocks/broker"
	"github.com/BAN1ce/skyTree/pkg/mocks/client"
	"github.com/eclipse/paho.golang/packets"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestQoS2_Publish(t *testing.T) {
	var (
		mockCtrl          = gomock.NewController(t)
		topicMeta         = topic.Meta{}
		mockMessageSource = broker.NewMockReadMessageSource(mockCtrl)
		mockWriter        = client.NewMockPacketWriter(mockCtrl)
		qos2Client        = NewQoS2(&topicMeta, mockWriter, mockMessageSource, nil)
		mockClient        = client.NewMockClient(mockCtrl)
	)

	defer mockCtrl.Finish()

	mockClient.EXPECT().HandelPublishComp(gomock.Any()).DoAndReturn(func(p *packets.Pubcomp) (bool, error) {
		if p.PacketID == 0 {
			return true, nil
		} else {
			return false, nil
		}
	}).AnyTimes()

	qos2Client.client = mockClient

	if ok, err := qos2Client.HandlePublishComp(&packets.Pubcomp{
		Properties: nil,
		PacketID:   0,
		ReasonCode: 0,
	}); err != nil || ok == false {
		t.Errorf("HandlePublishComp error: %v", err)
	}

	if ok, err := qos2Client.HandlePublishComp(&packets.Pubcomp{
		Properties: nil,
		PacketID:   1,
		ReasonCode: 0,
	}); ok {
		t.Errorf("HandlePublishComp error: %v", err)
	}

}

func TestHandlePublishRec(t *testing.T) {
	var (
		mockCtrl          = gomock.NewController(t)
		topicMeta         = topic.Meta{}
		mockMessageSource = broker.NewMockReadMessageSource(mockCtrl)
		mockWriter        = client.NewMockPacketWriter(mockCtrl)
		qos2Client        = NewQoS2(&topicMeta, mockWriter, mockMessageSource, nil)
		mockClient        = client.NewMockClient(mockCtrl)
	)

	defer mockCtrl.Finish()

	mockClient.EXPECT().HandlePublishRec(gomock.Any()).DoAndReturn(func(p *packets.Pubrec) (bool, error) {
		if p.PacketID == 0 {
			return true, nil
		} else {
			return false, nil
		}
	}).AnyTimes()

	qos2Client.client = mockClient

	if ok, err := qos2Client.HandlePublishRec(&packets.Pubrec{
		Properties: nil,
		PacketID:   0,
		ReasonCode: 0,
	}); err != nil || ok == false {
		t.Errorf("HandlePublishComp error: %v", err)
	}

	if ok, err := qos2Client.HandlePublishRec(&packets.Pubrec{
		Properties: nil,
		PacketID:   1,
		ReasonCode: 0,
	}); ok {
		t.Errorf("HandlePublishComp error: %v", err)
	}
}
