package repository

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRepository(t *testing.T) {
	repo := NewRepository(nil)
	assert.NotNil(t, repo)
}
