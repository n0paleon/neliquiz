package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConfig(t *testing.T) {
	cfg := New(".env.test")
	assert.NotNil(t, cfg)
	assert.NotEmpty(t, cfg.DBUser)
}

func TestNewConfigWithoutInit(t *testing.T) {
	cfg := GetConfig(".env.test")
	assert.NotNil(t, cfg)
	assert.NotEmpty(t, cfg.DBUser)
}
