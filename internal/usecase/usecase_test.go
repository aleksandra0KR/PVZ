package usecase

import (
	"final/internal/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUseCase(t *testing.T) {
	repo := &repository.Repository{}
	usecase := NewUseCase(repo)
	assert.NotNil(t, usecase)
}
