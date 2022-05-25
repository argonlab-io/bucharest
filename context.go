package bucharest

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Context interface {
	context.Context
	ENV() ENV
	GORM() *gorm.DB
	Log() *logrus.Logger
	Redis() *redis.Client
	SQL() *sql.DB
	SQLX() *sqlx.DB
	SetValue(key, val interface{})
	Update(option *ContextOptions)
}

var ErrNoENV = errors.New("ENV is not present in this context")
var ErrNoGORM = errors.New("*gorm.DB is not present in this context")
var ErrNoLogrus = errors.New("*logrus.Logger is not present in this context")
var ErrNoRedis = errors.New("*redis.Client is not present in this context")
var ErrNoSQL = errors.New("*sql.DB is not present in this context")
var ErrNoSQLX = errors.New("*sqlx.DB is not present in this context")

func AddValuesToContext(ctx Context, values MapAny) Context {
	for key, value := range values {
		ctx.SetValue(key, value)
	}
	return ctx
}
