package api

import (
	"context"
	client2 "github.com/BAN1ce/skyTree/api/client"
	"github.com/BAN1ce/skyTree/api/message"
	"github.com/BAN1ce/skyTree/api/session"
	"github.com/BAN1ce/skyTree/api/topic"
	_ "github.com/BAN1ce/skyTree/docs"
	broker2 "github.com/BAN1ce/skyTree/pkg/broker"
	session2 "github.com/BAN1ce/skyTree/pkg/broker/session"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	_ "net/http/pprof"
	"time"
)

type API struct {
	addr       string
	httpServer *gin.Engine
	server     *http.Server
	apiV1      *gin.RouterGroup
	*Component
}

type Component struct {
	TopicManager   topic.Manager
	ClientManager  client2.Manager
	MessageStore   broker2.MessageStore
	SessionManager session2.Manager
}

func NewAPI(addr string, component *Component) *API {
	api := &API{
		addr:      addr,
		Component: component,
	}
	return api
}

func (a *API) Start(ctx context.Context) error {
	r := gin.New()
	gin.SetMode(gin.ReleaseMode)
	a.httpServer = r
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})
	a.route()

	a.server = &http.Server{
		Addr:    a.addr,
		Handler: a.httpServer,
	}

	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// handle error
		}
	}()

	<-ctx.Done()
	return a.Close()
}

func (a *API) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return a.server.Shutdown(ctx)
}

func (a *API) route() {
	a.httpServer.GET("/debug/pprof/*filepath", gin.WrapH(http.DefaultServeMux))
	a.httpServer.GET("metrics", gin.WrapH(promhttp.Handler()))
	//a.httpServer.GET("/metrics", prometheusHandler())
	a.apiV1 = a.httpServer.Group("/api/v1")
	a.topic()
	a.client()
	a.session()
	a.message()
}

func (a *API) client() {
	ctr := client2.NewController(a.Component.ClientManager)
	a.apiV1.GET("/clients/:client_id", ctr.Info)
}

func (a *API) session() {
	ctr := session.NewController(a.SessionManager)
	a.apiV1.GET("/sessions/:client_id", ctr.Info)
}

func (a *API) topic() {
	ctr := topic.NewController(a.TopicManager)
	a.apiV1.GET("/topics/:topic", ctr.Info)
}

func (a *API) message() {
	ctr := message.NewController(a.MessageStore)
	a.apiV1.GET("/messages", ctr.Get)

}
func (a *API) Name() string {
	return "api"
}
