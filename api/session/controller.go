package session

import (
	"github.com/BAN1ce/skyTree/api/base"
	"github.com/BAN1ce/skyTree/pkg/broker/session"
	"github.com/BAN1ce/skyTree/pkg/model"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	sessions session.Manager
}

func NewController(manager session.Manager) *Controller {
	return &Controller{
		sessions: manager,
	}
}

func (o *Controller) Info(g *gin.Context) {
	var (
		req InfoReq
	)
	if err := g.ShouldBindUri(&req); err != nil {
		g.JSON(400, base.WithCode(400))
		return
	}
	clientSession, ok := o.sessions.ReadClientSession(g.Request.Context(), req.ClientID)
	if !ok {
		g.JSON(404, base.Response{
			Msg:     "client session not found",
			Success: false,
		})
		return
	}
	g.JSON(200, base.WithData(model.ToSessionModel(clientSession)))

}
