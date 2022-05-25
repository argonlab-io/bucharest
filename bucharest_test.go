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
	ctx := NewContextWithOptions(&ContextOptions{Parent: context.Background()})
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
	ctx := NewContextWithOptions(&ContextOptions{Parent: context.Background()})
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
		Parent: context.Background(),
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

	assert.Same(t, env, ctx.ENV())
	err = os.Remove(envTempPath)
	assert.NoError(t, err)
	assert.Same(t, gorm, ctx.GORM())
	assert.Same(t, logrus, ctx.Log())
	assert.Same(t, redis_, ctx.Redis())
	assert.Same(t, sql_, ctx.SQL())
	assert.Same(t, sqlx_, ctx.SQLX())
}

func TestUpdateContext(t *testing.T) {
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

	ctx := NewContextWithOptions(&ContextOptions{
		Parent: context.Background(),
		ENV:    env,
		GORM:   &gorm.DB{},
		Logrus: &logrus.Logger{},
		Redis:  &redis.Client{},
		SQL:    &sql.DB{},
		SQLX:   &sqlx.DB{},
	})
	err = os.Remove(envTempPath)
	assert.NoError(t, err)

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

	// prepare env
	envTempPath = "/tmp/.env"
	tempEnvFile, err = os.Create(envTempPath)
	assert.NoError(t, err)
	assert.NotEmpty(t, tempEnvFile)

	_, err = tempEnvFile.Write([]byte(``))
	assert.NoError(t, err)

	err = tempEnvFile.Close()
	assert.NoError(t, err)

	newENV, err := NewENV(envTempPath)
	assert.NoError(t, err)
	assert.NotNil(t, newENV)

	// prepare Gorm
	gorm_ := &gorm.DB{}

	// Prepare Logrus
	logrus_ := &logrus.Logger{}

	// Prepare Redis
	redis_ := &redis.Client{}

	// Prepare sql
	sql_ := &sql.DB{}

	// Prepare sqlx
	sqlx_ := &sqlx.DB{}

	assert.NotSame(t, newENV, ctx.ENV())
	err = os.Remove(envTempPath)
	assert.NoError(t, err)
	assert.NotSame(t, gorm_, ctx.GORM())
	assert.NotSame(t, logrus_, ctx.Log())
	assert.NotSame(t, redis_, ctx.Redis())
	assert.NotSame(t, sql_, ctx.SQL())
	assert.NotSame(t, sqlx_, ctx.SQLX())
}
