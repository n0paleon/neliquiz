package broadcast

import (
	"NeliQuiz/internal/features/game/delivery/gameerr"
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
	timer      *timer.Timer
	app        pitaya.Pitaya
	groupName  string
	serverType string
}

func (r *Room) AfterInit() {
	r.timer = pitaya.NewTimer(2*time.Second, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := r.app.GroupBroadcast(ctx, r.serverType, r.groupName, "onPing", time.Now().UnixMilli()); err != nil {
			logrus.Errorf("failed to send group broadcast: %v", err)
		}
	})
}

func (r *Room) Join(ctx context.Context, _ []byte) (*map[string]string, error) {
	s := r.app.GetSessionFromCtx(ctx)
	if err := s.Bind(ctx, strconv.Itoa(int(s.ID()))); err != nil {
		logrus.Warnf("failed to bind session: %v", err)
	}

	if err := r.app.GroupAddMember(ctx, r.groupName, s.UID()); err != nil {
		logrus.Errorf("failed to add group member: %v", err)

		err2 := gameerr.ErrFailedToJoinRoom
		return nil, pitaya.Error(err2, err2.Code, map[string]string{
			"message": err2.Message,
		})
	}

	_ = s.OnClose(func() {
		_ = r.app.GroupRemoveMember(ctx, r.groupName, s.UID())
	})

	return &map[string]string{
		"message": "joined",
	}, nil
}

func NewBroadcastRoom(server *gameserver.Server) *Room {
	room := &Room{
		app:        server.App,
		groupName:  "broadcast_room",
		serverType: server.ServerType,
	}

	if err := room.app.GroupCreate(context.Background(), room.groupName); err != nil {
		logrus.Fatalf("failed to create room: %v", err)
	}

	return room
}
