package bucharest

import (
	"context"
	"database/sql"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ContextOptions struct {
	ENV    ENV
	GORM   *gorm.DB
	Logrus *logrus.Logger
	Parent context.Context
	SQL    *sql.DB
	SQLX   *sqlx.DB
	Redis  *redis.Client
}

type BuchatrestContext struct {
	context.Context
	env     ENV
	gorm_   *gorm.DB
	logrus_ *logrus.Logger
	redis_  *redis.Client
	sql_    *sql.DB
	sqlx_   *sqlx.DB
}

func NewContextWithOptions(options *ContextOptions) Context {
	if options == nil {
		options = &ContextOptions{}
	}
	if options.Parent == nil {
		options.Parent = context.Background()
	}
	return &BuchatrestContext{
		Context: options.Parent,
		env:     options.ENV,
		gorm_:   options.GORM,
		logrus_: options.Logrus,
		redis_:  options.Redis,
		sql_:    options.SQL,
		sqlx_:   options.SQLX,
	}
}

func (d *BuchatrestContext) String() string {
	return "bucharest.BuchatrestContext"
}
func (ctx *BuchatrestContext) ENV() ENV {
	if ctx.env == nil {
		panic(ErrNoENV)
	}
	return ctx.env
}

func (ctx *BuchatrestContext) GORM() *gorm.DB {
	if ctx.gorm_ == nil {
		panic(ErrNoGORM)
	}
	return ctx.gorm_
}

func (ctx *BuchatrestContext) Log() *logrus.Logger {
	if ctx.logrus_ == nil {
		panic(ErrNoLogrus)
	}
	return ctx.logrus_
}

func (ctx *BuchatrestContext) Redis() *redis.Client {
	if ctx.redis_ == nil {
		panic(ErrNoRedis)
	}
	return ctx.redis_
}

func (ctx *BuchatrestContext) SQL() *sql.DB {
	if ctx.sql_ == nil {
		panic(ErrNoSQL)
	}
	return ctx.sql_
}

func (ctx *BuchatrestContext) SQLX() *sqlx.DB {
	if ctx.sqlx_ == nil {
		panic(ErrNoSQLX)
	}
	return ctx.sqlx_
}

func (ctx *BuchatrestContext) SetValue(key, val interface{}) {
	ctx.Context = context.WithValue(ctx.Context, key, val)
}
