package broadcast

import (
	"NeliQuiz/internal/infrastructures/gameserver"
	"context"
	"github.com/sirupsen/logrus"
	pitaya "github.com/topfreegames/pitaya/v3/pkg"
	"github.com/topfreegames/pitaya/v3/pkg/component"
	"github.com/topfreegames/pitaya/v3/pkg/timer"
	"strconv"
	"time"
)

type Room struct {
	component.Base
	timer        *timer.Timer
	app          pitaya.Pitaya
	groupName    string
	frontendType string
}

func (r *Room) AfterInit() {
	r.timer = pitaya.NewTimer(10*time.Second, func() {
		r.updateOnlineUsers()
	})
}

func (r *Room) Ping(ctx context.Context, _ []byte) (*map[string]int64, error) {
	s := r.app.GetSessionFromCtx(ctx)
	uid := strconv.Itoa(int(s.ID()))
	_ = s.Bind(ctx, uid)
	err := r.app.GroupAddMember(ctx, r.groupName, uid)
	if err == nil {
		r.updateOnlineUsers()
	}

	_ = s.OnClose(func() {
		_ = r.app.GroupRemoveMember(ctx, r.groupName, uid)
		r.updateOnlineUsers()
	})

	return &map[string]int64{
		"ts": time.Now().UnixMilli(),
	}, nil
}

func (r *Room) updateOnlineUsers() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	total := r.app.GetNumberOfConnectedClients()

	msg := &map[string]int64{
		"total": total,
	}
	if err := r.app.GroupBroadcast(ctx, r.frontendType, r.groupName, "onOnlineUsersUpdate", msg); err != nil {
		logrus.Warnf("failed to broadcast onOnlineUsersUpdate: %v", err)
	}
}

func NewBroadcastRoom(server *gameserver.Server) *Room {
	room := &Room{
		app:          server.App,
		groupName:    "broadcast_room",
		frontendType: server.ServerType,
	}

	if err := room.app.GroupCreate(context.Background(), room.groupName); err != nil {
		logrus.Fatalf("failed to create broadcast room: %v", err)
	}

	return room
}
