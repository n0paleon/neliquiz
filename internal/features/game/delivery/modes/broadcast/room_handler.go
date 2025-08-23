package broadcast

import (
	"NeliQuiz/internal/infrastructures/gameserver"
	"context"
	pitaya "github.com/topfreegames/pitaya/v3/pkg"
	"github.com/topfreegames/pitaya/v3/pkg/component"
	"github.com/topfreegames/pitaya/v3/pkg/timer"
	"time"
)

type Room struct {
	component.Base
	timer *timer.Timer
	app   pitaya.Pitaya
}

func (r *Room) Ping(_ context.Context, _ []byte) (*map[string]int64, error) {
	return &map[string]int64{
		"ts": time.Now().UnixMilli(),
	}, nil
}

func NewBroadcastRoom(server *gameserver.Server) *Room {
	room := &Room{
		app: server.App,
	}

	return room
}
