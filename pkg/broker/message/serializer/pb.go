package serializer

import (
	"bytes"
	"fmt"
	"github.com/BAN1ce/skyTree/pkg/packet"
	"github.com/BAN1ce/skyTree/pkg/pool"
	"github.com/eclipse/paho.golang/packets"
	"google.golang.org/protobuf/proto"
	"time"
)

const ProtoBufVersion = SerialVersion(iota + 1)

type ProtoBufSerializer struct {
}

func (p *ProtoBufSerializer) Encode(pub *packet.Message, buf *bytes.Buffer) error {
	var (
		pubBuf       = pool.ByteBufferPool.Get()
		protoMessage = packet.StorePublishMessage{
			MessageID:   pub.MessageID,
			CreatedTime: time.Now().Unix(),
			ExpiredTime: pub.ExpiredTime,
			SenderID:    pub.SendClientID,
		}
	)
	defer func() {
		pool.ByteBufferPool.Put(pubBuf)
	}()

	if pub.PublishPacket != nil {
		if _, err := pub.PublishPacket.WriteTo(pubBuf); err != nil {
			return err
		}
	}

	protoMessage.PublishPacket = pubBuf.Bytes()
	pbBody, err := proto.Marshal(&protoMessage)
	if err != nil {
		return err
	}
	if n, err := buf.Write(pbBody); err != nil {
		return err
	} else if n != len(pbBody) {
		return fmt.Errorf("write to buffer error, expect %d, got %d", len(pbBody), n)
	}
	return nil

}

func (p *ProtoBufSerializer) Decode(rawData []byte) (*packet.Message, error) {
	var (
		protoMessage = packet.StorePublishMessage{}
		bf           = pool.ByteBufferPool.Get()
	)
	defer pool.ByteBufferPool.Put(bf)
	if err := proto.Unmarshal(rawData, &protoMessage); err != nil {
		return nil, err
	}
	bf.Write(protoMessage.PublishPacket)

	ctl, err := packets.ReadPacket(bf)
	if err != nil {
		return nil, err
	}

	if ctl.Content == nil {
		return nil, fmt.Errorf("empty content")
	}

	if pub, ok := ctl.Content.(*packets.Publish); !ok {
		return nil, fmt.Errorf("invalid content")
	} else {

		return &packet.Message{
			SendClientID:  protoMessage.SenderID,
			MessageID:     protoMessage.MessageID,
			PublishPacket: pub,
			PubReceived:   protoMessage.PubReceived,
			CreatedTime:   protoMessage.CreatedTime,
			ExpiredTime:   protoMessage.ExpiredTime,
		}, nil
	}

}
