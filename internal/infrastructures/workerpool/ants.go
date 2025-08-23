package workerpool

import (
	"log"
	"sync"

	"github.com/panjf2000/ants/v2"
)

var (
	once sync.Once
	pool *ants.Pool
)

// InitPool create ants.Pool instance
func InitPool(size int) *ants.Pool {
	once.Do(func() {
		var err error
		pool, err = ants.NewPool(size,
			ants.WithPreAlloc(true),
		)
		if err != nil {
			log.Fatalf("failed to create worker pool: %v", err)
		}
	})
	return pool
}

// GetPool get existing ants.Pool
func GetPool() *ants.Pool {
	if pool == nil {
		log.Fatal("workerpool is not initialized")
	}
	return pool
}
