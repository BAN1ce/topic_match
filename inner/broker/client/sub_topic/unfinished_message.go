package sub_topic

import (
	"context"
	"github.com/BAN1ce/skyTree/logger"
	"github.com/BAN1ce/skyTree/pkg/broker"
	"github.com/BAN1ce/skyTree/pkg/packet"
)

func FillUnfinishedMessage(ctx context.Context, message []*packet.Message, source broker.ReadMessageSource) []*packet.Message {

	for i := 0; i < len(message); i++ {
		msg := message[i]
		if m, err := source.NextMessages(ctx, 1, msg.MessageID, true, nil); err != nil || len(m) == 0 {
			logger.Logger.Error().Err(err).Msg("fill unfinished message error")
			continue
		} else {
			msg.PublishPacket = m[0].PublishPacket
		}
	}
	logger.Logger.Debug().Any("message", message).Msg("fill unfinished message")
	return message
}
