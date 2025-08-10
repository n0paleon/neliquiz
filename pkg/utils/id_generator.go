package utils

import (
	cryptoRand "crypto/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

// GenerateULID create unique ID based on ULID algorithm with 26 chars length
func GenerateULID() string {
	t := time.Now().UTC()
	entropy := ulid.Monotonic(cryptoRand.Reader, 0)
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}
