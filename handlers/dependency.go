package handlers

import (
	"github.com/allegro/bigcache/v3"
	"github.com/jmoiron/sqlx"
)

type Dependency struct {
	DB        *sqlx.DB
	Memory    *bigcache.BigCache
	JWTSecret []byte
}
