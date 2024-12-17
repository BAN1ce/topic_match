package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/BAN1ce/skyTree/inner/broker/server/tcp"
	"github.com/BAN1ce/skyTree/logger"
	"github.com/BAN1ce/skyTree/pkg/errs"
	"net"
	"strings"
	"sync"
	"time"
)

type Listener interface {
	Accept() (net.Conn, error)
	Close() error
	Listen() error
	Name() string
}

type Server struct {
	listener []Listener
	mux      sync.RWMutex
	wg       sync.WaitGroup
	conn     chan net.Conn
	started  bool
	cancel   context.CancelFunc
}

func NewServer(adders []string) *Server {
	var listener []Listener

	for _, address := range adders {
		protocol, addr, err := getProtocolAndAddress(address)
		if err != nil {
			logger.Logger.Fatal().Err(err).Msg("get protocol and address error")
		}
		switch protocol {
		case "tcp":
			tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
			if err != nil {
				logger.Logger.Fatal().Err(err).Msg("resolve tcp addr error")
			}
			listener = append(listener, tcp.NewListener(tcpAddr))

		default:
			logger.Logger.Fatal().Msg("Unsupported protocol:" + protocol)
		}

	}

	return &Server{
		listener: listener,
		conn:     make(chan net.Conn),
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.started {
		return errs.ErrServerStarted
	}
	if len(s.listener) == 0 {
		return errs.ErrListenerIsNil
	}
	ctx, s.cancel = context.WithCancel(ctx)
	s.wg.Add(len(s.listener))

	s.startListener(ctx)
	s.started = true

	return nil
}

func (s *Server) startListener(ctx context.Context) {
	for _, l := range s.listener {
		go func(l Listener) {
			defer s.wg.Done()
			if err := l.Listen(); err != nil {
				logger.Logger.Fatal().Err(err).Msg("server listen failed")
			}
			logger.Logger.Info().Str("Listener", l.Name()).Msg("listen success")
			for {

				select {
				case <-ctx.Done():
					return
				default:
					conn, err := l.Accept()
					if err != nil {
						if errors.Is(err, net.ErrClosed) {
							logger.Logger.Info().Str("Listener", l.Name()).Msg("listener closed")
							return
						}
						// TODO: graceful shutdown should not output error loggers
						logger.Logger.Error().Err(err).Msg("accept error")
						continue
					}
					if conn != nil {
						s.conn <- conn
					}
				}
			}
		}(l)
	}
}

func (s *Server) Close() error {
	s.mux.Lock()
	defer s.mux.Unlock()
	if !s.started {
		return errs.ErrServerNotStarted
	}
	s.cancel()

	for _, l := range s.listener {
		logger.Logger.Info().Str("Listener", l.Name()).Msg("closing listener")
		if err := l.Close(); err != nil {
			return err
		}
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	go func() {
		s.wg.Wait()
		cancel()
	}()

	<-ctx.Done()
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return errs.ErrCloseListenerTimeout
	}
	s.started = false
	return nil
}

func (s *Server) Conn() <-chan net.Conn {
	return s.conn
}

func getProtocolAndAddress(address string) (string, string, error) {
	parts := strings.SplitN(address, "://", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid address format: %s", address)
	}
	return parts[0], parts[1], nil
}
