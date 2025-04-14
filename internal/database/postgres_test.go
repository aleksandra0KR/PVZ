package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDB(t *testing.T) {
	p := &Postgres{db: &sqlx.DB{}}
	assert.NotNil(t, p.GetDB())
}
