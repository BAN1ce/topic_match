package sub_topic

import (
	"github.com/BAN1ce/skyTree/config"
	"github.com/BAN1ce/skyTree/inner/facade"
	"github.com/BAN1ce/skyTree/logger"
	"github.com/BAN1ce/skyTree/pkg/broker"
	"github.com/BAN1ce/skyTree/pkg/broker/client"
	"github.com/BAN1ce/skyTree/pkg/errs"
	"github.com/BAN1ce/skyTree/pkg/packet"
	"github.com/BAN1ce/skyTree/pkg/retry"
	"github.com/eclipse/paho.golang/packets"
	"github.com/google/uuid"
	"github.com/zyedidia/generic/list"
	"sync"
	"time"
)

type WithRetryClient struct {
	client *Client
	queue  *list.List[*packet.Message]
	mux    sync.RWMutex
}

func NewQoSWithRetry(client *Client) *WithRetryClient {
	return &WithRetryClient{
		client: client,
		queue:  list.New[*packet.Message](),
	}
}

func (r *WithRetryClient) Publish(publish *packet.Message) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	if publish.PublishPacket.QoS == broker.QoS0 {
		return r.client.Publish(publish)
	}

	// NoLocal means that the server does not send messages published by the client itself.
	if r.client.meta.NoLocal && r.client.writer.GetID() == publish.SendClientID {
		return nil
	}

	if publish.RetryInfo == nil {
		publish.RetryInfo = &packet.RetryInfo{
			Key:          uuid.NewString(),
			FirstPubTime: time.Now(),
			IntervalTime: config.GetConfig().Retry.GetInterval(),
		}
		r.queue.PushBack(publish)
	}

	var (
		retryTask = retry.NewTask(publish.RetryInfo.Key, publish, r.GetPacketWriter().GetID(), publish.RetryInfo.IntervalTime)
	)

	//  create a retry task
	if err := facade.GetPublishRetry().Create(retryTask); err != nil {
		logger.Logger.Error().Err(err).Msg("WithRetryClient: publish retry failed")
	}

	return r.client.Publish(publish)
}

func (r *WithRetryClient) messageToRetryTask(m *packet.Message) *retry.Task {
	return &retry.Task{
		Key:      uuid.NewString(),
		Data:     m,
		ClientID: r.client.GetPacketWriter().GetID(),
	}

}

func (r *WithRetryClient) HandlePublishAck(pubAck *packets.Puback) (err error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	node := r.queue.Front

	for node != nil {
		if node.Value.PublishPacket.PacketID == pubAck.PacketID {
			if node.Value.PublishPacket.QoS == broker.QoS1 {
				facade.GetPublishRetry().Delete(node.Value.RetryInfo.Key)
				r.queue.Remove(node)
				return r.client.HandlePublishAck(pubAck)
			}
			break
		}
		node = node.Next
	}

	err = errs.ErrPacketIDNotFound
	return
}

func (r *WithRetryClient) GetPacketWriter() client.PacketWriter {
	return r.client.GetPacketWriter()
}

func (r *WithRetryClient) PubRel(message *packet.Message) error {
	// TODO: need retry ?
	return r.client.PubRel(message)
}

func (r *WithRetryClient) GetUnFinishedMessage() []*packet.Message {
	r.mux.RLock()
	defer r.mux.RUnlock()

	var (
		unFinish []*packet.Message
	)
	r.queue.Front.Each(func(val *packet.Message) {
		unFinish = append(unFinish, val)
	})
	return unFinish
}

func (r *WithRetryClient) HandlePublishRec(pubRec *packets.Pubrec) (err error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	var (
		node    = r.queue.Front
		message *packet.Message
	)
	for node != nil {
		if node.Value.PublishPacket.PacketID == pubRec.PacketID {
			message = node.Value
			node.Value.PubReceived = true

			if err := r.client.PubRel(message); err != nil {
				logger.Logger.Warn().Err(err).Msg("WithRetryClient: pubRel failed")
			}
			return r.client.HandlePublishRec(pubRec)
		}
		node = node.Next
	}
	err = errs.ErrPacketIDNotFound
	return err
}

func (r *WithRetryClient) HandelPublishComp(pubComp *packets.Pubcomp) (err error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	node := r.queue.Front
	for node != nil {
		if node.Value.PublishPacket.PacketID == pubComp.PacketID {
			node.Value.PubReceived = true
			r.queue.Remove(node)
			facade.GetPublishRetry().Delete(node.Value.RetryInfo.Key)
			return r.client.HandelPublishComp(pubComp)
		}
		node = node.Next
	}
	err = errs.ErrPacketIDNotFound
	return err
}

func (r *WithRetryClient) Close() error {
	return nil
}
