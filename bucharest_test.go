package bucharest_test

import (
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

func TestNewContextWithOptions(t *testing.T) {
	ctx := NewContextWithOptions(nil)
	if ctx == nil {
		t.Fatalf("NewContext returned nil")
	}
	select {
	case x := <-ctx.Done():
		t.Errorf("<-c.Done() == %v want nothing (it should block)", x)
	default:
	}
	if got, want := fmt.Sprint(ctx), "bucharest.BuchatrestContext"; got != want {
		t.Errorf("NewContextWithOptions().String() = %q want %q", got, want)
	}
}

func TestNewContextWithOptionsFromParentContext(t *testing.T) {
	ctx := NewContextWithOptions(&ContextOptions{Parent: NewContext()})
	if ctx == nil {
		t.Fatalf("NewContext returned nil")
	}
	select {
	case x := <-ctx.Done():
		t.Errorf("<-c.Done() == %v want nothing (it should block)", x)
	default:
	}
	if got, want := fmt.Sprint(ctx), "bucharest.BuchatrestContext"; got != want {
		t.Errorf("NewContextWithOptions().String() = %q want %q", got, want)
	}
}

func TestNewContextWithOptionsWithNoAddtionalOptions(t *testing.T) {
	ctx := NewContextWithOptions(&ContextOptions{Parent: NewContext()})
	if ctx == nil {
		t.Fatalf("NewContext returned nil")
	}
	select {
	case x := <-ctx.Done():
		t.Errorf("<-c.Done() == %v want nothing (it should block)", x)
	default:
	}
	if got, want := fmt.Sprint(ctx), "bucharest.BuchatrestContext"; got != want {
		t.Errorf("NewContextWithOptions().String() = %q want %q", got, want)
	}

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

	// prepare Gorm
	gorm := &gorm.DB{}

	// Prepare Logrus
	logrus := &logrus.Logger{}

	// Prepare Redis
	redis_ := &redis.Client{}

	// Prepare sql
	sql_ := &sql.DB{}

	// Prepare sqlx
	sqlx_ := &sqlx.DB{}

	ctx := NewContextWithOptions(&ContextOptions{
		Parent: NewContext(),
		ENV:    env,
		GORM:   gorm,
		Logrus: logrus,
		Redis:  redis_,
		SQL:    sql_,
		SQLX:   sqlx_,
	})
	if ctx == nil {
		t.Fatalf("NewContext returned nil")
	}
	select {
	case x := <-ctx.Done():
		t.Errorf("<-c.Done() == %v want nothing (it should block)", x)
	default:
	}
	if got, want := fmt.Sprint(ctx), "bucharest.BuchatrestContext"; got != want {
		t.Errorf("NewContextWithOptions().String() = %q want %q", got, want)
	}

	assert.Equal(t, env, ctx.ENV())
	err = os.Remove(envTempPath)
	assert.NoError(t, err)
	assert.Equal(t, gorm, ctx.GORM())
	assert.Equal(t, logrus, ctx.Log())
	assert.Equal(t, redis_, ctx.Redis())
	assert.Equal(t, sql_, ctx.SQL())
	assert.Equal(t, sqlx_, ctx.SQLX())
}
