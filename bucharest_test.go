package bucharest_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	. "github.com/argonlab-io/bucharest"
	"github.com/argonlab-io/bucharest/utils"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const FOUND_NULL_CONTEXT_ERR_MSG = "Found null context."
const BUCHAREST_CONTEXT_TYPE_NAME = "bucharest.BuchatrestContext"
const TEMP_ENV_PATH = "/tmp/.env"
const CONTEXT_DONE_ERR_MSG = "Context is done before expected. Got %v."

func verifyContextTypeName(t *testing.T, ctx Context) {
	if contextTypeName := fmt.Sprint(ctx); contextTypeName != BUCHAREST_CONTEXT_TYPE_NAME {
		t.Errorf("Unexptected context type = %q, want %q.", contextTypeName, BUCHAREST_CONTEXT_TYPE_NAME)
	}
}

func TestNewContextWithOptions(t *testing.T) {
	ctx := NewContextWithOptions(nil)
	if ctx == nil {
		t.Fatal(FOUND_NULL_CONTEXT_ERR_MSG)
	}
	select {
	case x := <-ctx.Done():
		t.Errorf(CONTEXT_DONE_ERR_MSG, x)
	default:
	}
	verifyContextTypeName(t, ctx)
}

func TestNewContextWithOptionsFromParentContext(t *testing.T) {
	ctx := NewContextWithOptions(&ContextOptions{Parent: context.Background()})
	if ctx == nil {
		t.Fatal(FOUND_NULL_CONTEXT_ERR_MSG)
	}
	select {
	case x := <-ctx.Done():
		t.Errorf(CONTEXT_DONE_ERR_MSG, x)
	default:
	}
	verifyContextTypeName(t, ctx)
}

func TestNewContextWithOptionsWithNoAddtionalOptions(t *testing.T) {
	ctx := NewContextWithOptions(&ContextOptions{Parent: context.Background()})
	if ctx == nil {
		t.Fatal(FOUND_NULL_CONTEXT_ERR_MSG)
	}
	select {
	case x := <-ctx.Done():
		t.Errorf(CONTEXT_DONE_ERR_MSG, x)
	default:
	}
	verifyContextTypeName(t, ctx)

	utils.AssertPanic(t, func() { ctx.ENV() }, ErrNoENV)
	utils.AssertPanic(t, func() { ctx.GORM() }, ErrNoGORM)
	utils.AssertPanic(t, func() { ctx.Log() }, ErrNoLogrus)
	utils.AssertPanic(t, func() { ctx.Redis() }, ErrNoRedis)
	utils.AssertPanic(t, func() { ctx.SQL() }, ErrNoSQL)
	utils.AssertPanic(t, func() { ctx.SQLX() }, ErrNoSQLX)
}

func TestNewContextWithOptionsWithAllAddtionalOptions(t *testing.T) {
	// prepare env
	envTempPath := "/tmp/.env"
	tempEnvFile, err := os.Create(envTempPath)
	assert.NoError(t, err)
	assert.NotEmpty(t, tempEnvFile)

	_, err = tempEnvFile.Write(envFile)
	assert.NoError(t, err)

	err = tempEnvFile.Close()
	assert.NoError(t, err)

	env, err := NewENV(envTempPath)
	assert.NoError(t, err)
	assert.NotNil(t, env)

	gorm := &gorm.DB{}
	logrus := &logrus.Logger{}
	redis := &redis.Client{}
	sql := &sql.DB{}
	sqlx := &sqlx.DB{}

	ctx := NewContextWithOptions(&ContextOptions{
		Parent: context.Background(),
		ENV:    env,
		GORM:   gorm,
		Logrus: logrus,
		Redis:  redis,
		SQL:    sql,
		SQLX:   sqlx,
	})
	if ctx == nil {
		t.Fatal(FOUND_NULL_CONTEXT_ERR_MSG)
	}
	select {
	case x := <-ctx.Done():
		t.Errorf("<-c.Done() == %v want nothing (it should block)", x)
	default:
	}
	if contextTypeName := fmt.Sprint(ctx); contextTypeName != BUCHAREST_CONTEXT_TYPE_NAME {
		t.Errorf("NewContextWithOptions().String() = %q want %q", contextTypeName, BUCHAREST_CONTEXT_TYPE_NAME)
	}

	assert.Same(t, env, ctx.ENV())
	err = os.Remove(envTempPath)
	assert.NoError(t, err)
	assert.Same(t, gorm, ctx.GORM())
	assert.Same(t, logrus, ctx.Log())
	assert.Same(t, redis, ctx.Redis())
	assert.Same(t, sql, ctx.SQL())
	assert.Same(t, sqlx, ctx.SQLX())
}

func TestUpdateContext(t *testing.T) {
	// prepare env
	tempEnvFile, err := os.Create(TEMP_ENV_PATH)
	assert.NoError(t, err)
	assert.NotEmpty(t, tempEnvFile)

	_, err = tempEnvFile.Write(envFile)
	assert.NoError(t, err)

	err = tempEnvFile.Close()
	assert.NoError(t, err)

	env, err := NewENV(TEMP_ENV_PATH)
	assert.NoError(t, err)
	assert.NotNil(t, env)

	ctx := NewContextWithOptions(&ContextOptions{
		Parent: context.Background(),
		ENV:    env,
		GORM:   &gorm.DB{},
		Logrus: &logrus.Logger{},
		Redis:  &redis.Client{},
		SQL:    &sql.DB{},
		SQLX:   &sqlx.DB{},
	})
	err = os.Remove(TEMP_ENV_PATH)
	assert.NoError(t, err)

	if ctx == nil {
		t.Fatal(FOUND_NULL_CONTEXT_ERR_MSG)
	}
	select {
	case x := <-ctx.Done():
		t.Errorf(CONTEXT_DONE_ERR_MSG, x)
	default:
	}
	verifyContextTypeName(t, ctx)

	// prepare env
	tempEnvFile, err = os.Create(TEMP_ENV_PATH)
	assert.NoError(t, err)
	assert.NotEmpty(t, tempEnvFile)

	_, err = tempEnvFile.Write([]byte(``))
	assert.NoError(t, err)

	err = tempEnvFile.Close()
	assert.NoError(t, err)

	newENV, err := NewENV(TEMP_ENV_PATH)
	assert.NoError(t, err)
	assert.NotNil(t, newENV)

	gorm := &gorm.DB{}
	logrus := &logrus.Logger{}
	redis := &redis.Client{}
	sql := &sql.DB{}
	sqlx := &sqlx.DB{}

	assert.NotSame(t, newENV, ctx.ENV())
	err = os.Remove(TEMP_ENV_PATH)
	assert.NoError(t, err)
	assert.NotSame(t, gorm, ctx.GORM())
	assert.NotSame(t, logrus, ctx.Log())
	assert.NotSame(t, redis, ctx.Redis())
	assert.NotSame(t, sql, ctx.SQL())
	assert.NotSame(t, sqlx, ctx.SQLX())

	ctx.Update(&ContextOptions{
		ENV:    newENV,
		GORM:   gorm,
		Logrus: logrus,
		Redis:  redis,
		SQL:    sql,
		SQLX:   sqlx,
	})

	assert.Same(t, newENV, ctx.ENV())
	assert.NoError(t, err)
	assert.Same(t, gorm, ctx.GORM())
	assert.Same(t, logrus, ctx.Log())
	assert.Same(t, redis, ctx.Redis())
	assert.Same(t, sql, ctx.SQL())
	assert.Same(t, sqlx, ctx.SQLX())

}
