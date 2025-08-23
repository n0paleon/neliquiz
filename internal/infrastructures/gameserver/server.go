package gameserver

import (
	"NeliQuiz/internal/infrastructures/workerpool"
	"NeliQuiz/internal/shared/config"
	"context"
	"github.com/sirupsen/logrus"
	pitaya "github.com/topfreegames/pitaya/v3/pkg"
	"github.com/topfreegames/pitaya/v3/pkg/acceptor"
	"github.com/topfreegames/pitaya/v3/pkg/component"
	pitayaConfig "github.com/topfreegames/pitaya/v3/pkg/config"
	"github.com/topfreegames/pitaya/v3/pkg/groups"
	"github.com/topfreegames/pitaya/v3/pkg/serialize/json"
	"net"
	"time"
)

type Server struct {
	cfg        *config.Config
	App        pitaya.Pitaya
	ServerType string
}

var app pitaya.Pitaya

func NewGameServer(cfg *config.Config) *Server {
	conf := configApp()
	serverType := "main"
	builder := pitaya.NewDefaultBuilder(true, serverType, pitaya.Standalone, map[string]string{}, *conf)

	addr := net.JoinHostPort(cfg.GameHost, cfg.GamePort)
	wsAcceptor := acceptor.NewWSAcceptor(addr)
	builder.AddAcceptor(wsAcceptor)
	builder.Groups = groups.NewMemoryGroupService(builder.Config.Groups.Memory)
	builder.Serializer = json.NewSerializer()
	app = builder.Build()

	app.SetDebug(cfg.GameDebug)
	if !cfg.GameDebug {
		pitaya.SetLogger(newCustomLogger())
	}

	return &Server{
		cfg:        cfg,
		App:        app,
		ServerType: serverType,
	}
}

func (s *Server) RegisterComponent(comp component.Component, opts ...component.Option) {
	if s.App == nil {
		logrus.Fatal("App not initialized. Call Initialize() first")
	}
	s.App.Register(comp, opts...)
}

func (s *Server) RegisterComponents(components map[string]component.Component) {
	for name, comp := range components {
		s.RegisterComponent(comp, component.WithName(name))
	}
}

func (s *Server) Start(ctx context.Context) {
	pool := workerpool.GetPool()

	// start game server global ticker
	_ = pool.Submit(func() {
		ticker := time.NewTicker(1 * time.Minute) // 1 minute interval
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				logrus.Infof("context is done, global ticker will be stopped")
				return
			case <-ticker.C:
				total := s.App.GetNumberOfConnectedClients()
				logrus.Infof("Connected clients: %d", total)
			}
		}
	})

	s.App.Start()
}

func (s *Server) Shutdown() {
	s.App.Shutdown()
}

func configApp() *pitayaConfig.PitayaConfig {
	conf := pitayaConfig.NewDefaultPitayaConfig()
	conf.Buffer.Handler.LocalProcess = 15
	conf.Heartbeat.Interval = time.Duration(15 * time.Second)
	conf.Buffer.Agent.Messages = 32
	conf.Handler.Messages.Compression = false
	conf.Cluster.SD.Etcd.Endpoints = []string{}
	return conf
}
