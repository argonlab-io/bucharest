package bucharest

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

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
}

var ErrNoENV = errors.New("ENV is not present in this context")
var ErrNoGORM = errors.New("*gorm.DB is not present in this context")
var ErrNoLogrus = errors.New("*logrus.Logger is not present in this context")
var ErrNoRedis = errors.New("*redis.Client is not present in this context")
var ErrNoSQL = errors.New("*sql.DB is not present in this context")
var ErrNoSQLX = errors.New("*sqlx.DB is not present in this context")

// DefaultContext is basically just a mimic of context.Background
// It is a basic object that implement bucharest.Context and context.Context
// For original implementation please see context.Background()
type DefaultContext struct{}

func (d *DefaultContext) String() string {
	return "bucharest.DefaultContext"
}

func (d *DefaultContext) Deadline() (deadline time.Time, ok bool) {
	return
}

func (ctx *DefaultContext) Done() <-chan struct{} {
	return nil
}

func (ctx *DefaultContext) Err() error {
	return nil
}

func (ctx *DefaultContext) Value(key interface{}) interface{} {
	return nil
}

func (ctx *DefaultContext) ENV() ENV {
	panic(ErrNoENV)
}

func (ctx *DefaultContext) GORM() *gorm.DB {
	panic(ErrNoGORM)
}

func (ctx *DefaultContext) Log() *logrus.Logger {
	panic(ErrNoLogrus)
}

func (ctx *DefaultContext) Redis() *redis.Client {
	panic(ErrNoRedis)
}

func (ctx *DefaultContext) SQL() *sql.DB {
	panic(ErrNoSQL)
}

func (ctx *DefaultContext) SQLX() *sqlx.DB {
	panic(ErrNoSQLX)
}

func NewContext() Context {
	return &DefaultContext{}
}

func NewContextWithCancel(parent Context) (ctx Context, cancel context.CancelFunc) {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	c := newCancelCtx(parent)
	propagateCancel(parent, &c)
	return &c, func() { c.cancel(true, context.Canceled) }
}

func newCancelCtx(parent Context) cancelCtx {
	return cancelCtx{Context: parent}
}

var goroutines int32

func propagateCancel(parent Context, child canceler) {
	done := parent.Done()
	if done == nil {
		return // parent is never canceled
	}

	select {
	case <-done:
		// parent is already canceled
		child.cancel(false, parent.Err())
		return
	default:
	}

	p, ok := parentCancelCtx(parent)
	if ok {
		p.mu.Lock()
		if p.err != nil {
			// parent has already been canceled
			child.cancel(false, p.err)
		} else {
			if p.children == nil {
				p.children = make(map[canceler]struct{})
			}
			p.children[child] = struct{}{}
		}
		p.mu.Unlock()
	} else {
		atomic.AddInt32(&goroutines, +1)
		go func() {
			select {
			case <-parent.Done():
				child.cancel(false, parent.Err())
			case <-child.Done():
			}
		}()
	}
}

var cancelCtxKey int

func parentCancelCtx(parent Context) (*cancelCtx, bool) {
	done := parent.Done()
	if done == closedchan || done == nil {
		return nil, false
	}
	p, ok := parent.Value(&cancelCtxKey).(*cancelCtx)
	if !ok {
		return nil, false
	}
	pdone, _ := p.done.Load().(chan struct{})
	if pdone != done {
		return nil, false
	}
	return p, true
}

func removeChild(parent Context, child canceler) {
	p, ok := parentCancelCtx(parent)
	if !ok {
		return
	}
	p.mu.Lock()
	if p.children != nil {
		delete(p.children, child)
	}
	p.mu.Unlock()
}

type canceler interface {
	cancel(removeFromParent bool, err error)
	Done() <-chan struct{}
}

var closedchan = make(chan struct{})

func init() {
	close(closedchan)
}

type cancelCtx struct {
	Context
	mu       sync.Mutex
	done     atomic.Value
	children map[canceler]struct{}
	err      error
}

func (c *cancelCtx) Value(key interface{}) interface{} {
	if key == &cancelCtxKey {
		return c
	}
	return c.Context.Value(key)
}

func (c *cancelCtx) Done() <-chan struct{} {
	d := c.done.Load()
	if d != nil {
		return d.(chan struct{})
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	d = c.done.Load()
	if d == nil {
		d = make(chan struct{})
		c.done.Store(d)
	}
	return d.(chan struct{})
}

func (c *cancelCtx) Err() error {
	c.mu.Lock()
	err := c.err
	c.mu.Unlock()
	return err
}

type stringer interface {
	String() string
}

func contextName(c Context) string {
	if s, ok := c.(stringer); ok {
		return s.String()
	}
	return reflect.TypeOf(c).String()
}

func (c *cancelCtx) String() string {
	return contextName(c.Context) + ".WithCancel"
}

func (c *cancelCtx) cancel(removeFromParent bool, err error) {
	if err == nil {
		panic("context: internal error: missing cancel error")
	}
	c.mu.Lock()
	if c.err != nil {
		c.mu.Unlock()
		return
	}
	c.err = err
	d, _ := c.done.Load().(chan struct{})
	if d == nil {
		c.done.Store(closedchan)
	} else {
		close(d)
	}

	for child := range c.children {
		child.cancel(false, err)
	}
	c.children = nil
	c.mu.Unlock()

	if removeFromParent {
		removeChild(c.Context, c)
	}
}

func NewContextWithDeadline(parent Context, d time.Time) (Context, context.CancelFunc) {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	if cur, ok := parent.Deadline(); ok && cur.Before(d) {
		return NewContextWithCancel(parent)
	}
	c := &timerCtx{
		cancelCtx: newCancelCtx(parent),
		deadline:  d,
	}
	propagateCancel(parent, c)
	dur := time.Until(d)
	if dur <= 0 {
		c.cancel(true, context.DeadlineExceeded)
		return c, func() { c.cancel(false, context.Canceled) }
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err == nil {
		c.timer = time.AfterFunc(dur, func() {
			c.cancel(true, context.DeadlineExceeded)
		})
	}
	return c, func() { c.cancel(true, context.Canceled) }
}

type timerCtx struct {
	cancelCtx
	timer    *time.Timer
	deadline time.Time
}

func (c *timerCtx) Deadline() (deadline time.Time, ok bool) {
	return c.deadline, true
}

func (c *timerCtx) String() string {
	return contextName(c.cancelCtx.Context) + ".WithDeadline(" +
		c.deadline.String() + " [" +
		time.Until(c.deadline).String() + "])"
}

func (c *timerCtx) cancel(removeFromParent bool, err error) {
	c.cancelCtx.cancel(false, err)
	if removeFromParent {
		// Remove this timerCtx from its parent cancelCtx's children.
		removeChild(c.cancelCtx.Context, c)
	}
	c.mu.Lock()
	if c.timer != nil {
		c.timer.Stop()
		c.timer = nil
	}
	c.mu.Unlock()
}

func NewContextWithTimeout(parent Context, timeout time.Duration) (Context, context.CancelFunc) {
	return NewContextWithDeadline(parent, time.Now().Add(timeout))
}

func NewContextWithValue(parent Context, key, val interface{}) Context {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	if key == nil {
		panic("nil key")
	}
	if !reflect.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}
	return &valueCtx{parent, key, val}
}

type valueCtx struct {
	Context
	key, val interface{}
}

func (c *valueCtx) String() string {
	return contextName(c.Context) + ".WithValue(type " +
		reflect.TypeOf(c.key).String() +
		", val " + stringify(c.val) + ")"
}

func (c *valueCtx) Value(key interface{}) interface{} {
	if c.key == key {
		return c.val
	}
	return c.Context.Value(key)
}

func stringify(v interface{}) string {
	switch s := v.(type) {
	case stringer:
		return s.String()
	case string:
		return s
	}
	return "<not Stringer>"
}

func AddValuesToContext(ctx Context, values MapAny) Context {
	for key, value := range values {
		ctx = NewContextWithValue(ctx, key, value)
	}
	return ctx
}
