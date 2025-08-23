package gameerr

var (
	ErrSessionBindFailed = New("Q-000", "failed to bind session")
	ErrRoomNotFound      = New("Q-001", "room not found")
	ErrFailedToJoinRoom  = New("Q-002", "failed to join room")
)
