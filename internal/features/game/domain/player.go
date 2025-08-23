package domain

type Player struct {
	UserID    string
	Nickname  string
	Points    int
	Answered  bool
	Connected bool
}
