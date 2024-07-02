package sub_topic

import (
	"context"
	"github.com/BAN1ce/skyTree/logger"
	"github.com/BAN1ce/skyTree/pkg/broker"
	"github.com/BAN1ce/skyTree/pkg/packet"
	"go.uber.org/zap"
)

func FillUnfinishedMessage(ctx context.Context, message []*packet.Message, source broker.ReadMessageSource) []*packet.Message {

	for i := 0; i < len(message); i++ {
		msg := message[i]
		if m, err := source.NextMessages(ctx, 1, msg.MessageID, true); err != nil || len(m) == 0 {
			logger.Logger.Error("fill unfinished message error", zap.Error(err))
			continue
		} else {
			msg.PublishPacket = m[0].PublishPacket
		}
	}
	logger.Logger.Debug("fill unfinished message", zap.Int("message", len(message)))
	return message
}
