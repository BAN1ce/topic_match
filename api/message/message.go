package message

import (
	"github.com/BAN1ce/skyTree/api/base"
	"github.com/BAN1ce/skyTree/pkg/broker"
	"github.com/gin-gonic/gin"
	"time"
)

type Controller struct {
	store broker.TopicMessageStore
}

func NewController(store broker.TopicMessageStore) *Controller {
	return &Controller{
		store: store,
	}
}

func (c *Controller) Get(g *gin.Context) {
	var (
		req     GetRequest
		rspData GetResponseData
	)
	if err := g.ShouldBindQuery(&req); err != nil {
		_ = g.AbortWithError(400, err)
		return
	}
	messages, err := c.store.ReadTopicMessage(g.Request.Context(), req.Topic, req.Start, req.Limit)
	if err != nil {
		g.JSON(200, base.WithError(err))
		return
	}
	rspData.Total, _ = c.store.TopicMessageTotal(g.Request.Context(), req.Topic)

	for _, m := range messages {
		rspData.Message = append(rspData.Message, Message{
			ID:          m.MessageID,
			CreatedTime: time.Unix(m.CreatedTime, 0),
			SenderID:    m.SendClientID,
			ExpiredTime: time.Unix(m.ExpiredTime, 0),
		})
	}
	g.JSON(200, base.WithData(rspData))
}
