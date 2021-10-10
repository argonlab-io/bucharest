package bucharest

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/argonlab-io/bucharest/utils"
	"github.com/stretchr/testify/assert"
)

type otherContext struct {
	Context
}

func contains(m map[canceler]struct{}, key canceler) bool {
	_, ret := m[key]
	return ret
}

const (
	shortDuration    = 1 * time.Millisecond // a reasonable duration to block in a test
	veryLongDuration = 1000 * time.Hour     // an arbitrary upper bound on the test's running time
)

// quiescent returns an arbitrary duration by which the program should have
// completed any remaining work and reached a steady (idle) state.
func quiescent(t *testing.T) time.Duration {
	deadline, ok := t.Deadline()
	if !ok {
		return 5 * time.Second
	}

	const arbitraryCleanupMargin = 1 * time.Second
	return time.Until(deadline) - arbitraryCleanupMargin
}

func XTestNewContext(t *testing.T) {
	c := NewContext()
	if c == nil {
		t.Fatalf("NewContext returned nil")
	}
	select {
	case x := <-c.Done():
		t.Errorf("<-c.Done() == %v want nothing (it should block)", x)
	default:
	}
	if got, want := fmt.Sprint(c), "bucharest.DefaultContext"; got != want {
		t.Errorf("NewContext().String() = %q want %q", got, want)
	}
}

func XTestWithCancel(t *testing.T) {
	c1, cancel := NewContextWithCancel(NewContext())

	if got, want := fmt.Sprint(c1), "bucharest.DefaultContext.WithCancel"; got != want {
		t.Errorf("c1.String() = %q want %q", got, want)
	}

	o := otherContext{c1}
	c2, lintIgnore := NewContextWithCancel(o) // cancel() is to propagate synchronously.
	_ = lintIgnore
	contexts := []Context{c1, o, c2}

	for i, c := range contexts {
		if d := c.Done(); d == nil {
			t.Errorf("c[%d].Done() == %v want non-nil", i, d)
		}
		if e := c.Err(); e != nil {
			t.Errorf("c[%d].Err() == %v want nil", i, e)
		}

		select {
		case x := <-c.Done():
			t.Errorf("<-c.Done() == %v want nothing (it should block)", x)
		default:
		}
	}

	cancel() // Should propagate synchronously.
	for i, c := range contexts {
		select {
		case <-c.Done():
		default:
			t.Errorf("<-c[%d].Done() blocked, but shouldn't have", i)
		}
		if e := c.Err(); e != context.Canceled {
			t.Errorf("c[%d].Err() == %v want %v", i, e, context.Canceled)
		}
	}
}

func XTestParentFinishesChild(t *testing.T) {
	parent, cancel := NewContextWithCancel(NewContext())
	cancelChild, stop := NewContextWithCancel(parent)
	defer stop()
	valueChild := NewContextWithValue(parent, "key", "value")
	timerChild, stop := NewContextWithTimeout(valueChild, veryLongDuration)
	defer stop()

	select {
	case x := <-parent.Done():
		t.Errorf("<-parent.Done() == %v want nothing (it should block)", x)
	case x := <-cancelChild.Done():
		t.Errorf("<-cancelChild.Done() == %v want nothing (it should block)", x)
	case x := <-timerChild.Done():
		t.Errorf("<-timerChild.Done() == %v want nothing (it should block)", x)
	case x := <-valueChild.Done():
		t.Errorf("<-valueChild.Done() == %v want nothing (it should block)", x)
	default:
	}

	// The parent's children should contain the two cancelable children.
	pc := parent.(*cancelCtx)
	cc := cancelChild.(*cancelCtx)
	tc := timerChild.(*timerCtx)
	parent.(*cancelCtx).mu.Lock()
	if len(pc.children) != 2 || !contains(pc.children, cc) || !contains(pc.children, tc) {
		t.Errorf("bad linkage: pc.children = %v, want %v and %v",
			pc.children, cc, tc)
	}
	pc.mu.Unlock()

	if p, ok := parentCancelCtx(cc.Context); !ok || p != pc {
		t.Errorf("bad linkage: parentCancelCtx(cancelChild.Context) = %v, %v want %v, true", p, ok, pc)
	}
	if p, ok := parentCancelCtx(tc.Context); !ok || p != pc {
		t.Errorf("bad linkage: parentCancelCtx(timerChild.Context) = %v, %v want %v, true", p, ok, pc)
	}

	cancel()

	pc.mu.Lock()
	if len(pc.children) != 0 {
		t.Errorf("pc.cancel didn't clear pc.children = %v", pc.children)
	}
	pc.mu.Unlock()

	// parent and children should all be finished.
	check := func(ctx Context, name string) {
		select {
		case <-ctx.Done():
		default:
			t.Errorf("<-%s.Done() blocked, but shouldn't have", name)
		}
		if e := ctx.Err(); e != context.Canceled {
			t.Errorf("%s.Err() == %v want %v", name, e, context.Canceled)
		}
	}
	check(parent, "parent")
	check(cancelChild, "cancelChild")
	check(valueChild, "valueChild")
	check(timerChild, "timerChild")

	// WithCancel should return a canceled context on a canceled parent.
	precanceledChild := NewContextWithValue(parent, "key", "value")
	select {
	case <-precanceledChild.Done():
	default:
		t.Errorf("<-precanceledChild.Done() blocked, but shouldn't have")
	}
	if e := precanceledChild.Err(); e != context.Canceled {
		t.Errorf("precanceledChild.Err() == %v want %v", e, context.Canceled)
	}
}

func XTestChildFinishesFirst(t *testing.T) {
	cancelable, stop := NewContextWithCancel(NewContext())
	defer stop()
	for _, parent := range []Context{NewContext(), cancelable} {
		child, cancel := NewContextWithCancel(parent)

		select {
		case x := <-parent.Done():
			t.Errorf("<-parent.Done() == %v want nothing (it should block)", x)
		case x := <-child.Done():
			t.Errorf("<-child.Done() == %v want nothing (it should block)", x)
		default:
		}

		cc := child.(*cancelCtx)
		pc, pcok := parent.(*cancelCtx) // pcok == false when parent == NewContext()
		if p, ok := parentCancelCtx(cc.Context); ok != pcok || (ok && pc != p) {
			t.Errorf("bad linkage: parentCancelCtx(cc.Context) = %v, %v want %v, %v", p, ok, pc, pcok)
		}

		if pcok {
			pc.mu.Lock()
			if len(pc.children) != 1 || !contains(pc.children, cc) {
				t.Errorf("bad linkage: pc.children = %v, cc = %v", pc.children, cc)
			}
			pc.mu.Unlock()
		}

		cancel()

		if pcok {
			pc.mu.Lock()
			if len(pc.children) != 0 {
				t.Errorf("child's cancel didn't remove self from pc.children = %v", pc.children)
			}
			pc.mu.Unlock()
		}

		// child should be finished.
		select {
		case <-child.Done():
		default:
			t.Errorf("<-child.Done() blocked, but shouldn't have")
		}
		if e := child.Err(); e != context.Canceled {
			t.Errorf("child.Err() == %v want %v", e, context.Canceled)
		}

		// parent should not be finished.
		select {
		case x := <-parent.Done():
			t.Errorf("<-parent.Done() == %v want nothing (it should block)", x)
		default:
		}
		if e := parent.Err(); e != nil {
			t.Errorf("parent.Err() == %v want nil", e)
		}
	}
}

func testDeadline(c Context, name string, t *testing.T) {
	t.Helper()
	d := quiescent(t)
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-timer.C:
		t.Fatalf("%s: context not timed out after %v", name, d)
	case <-c.Done():
	}
	if e := c.Err(); e != context.DeadlineExceeded {
		t.Errorf("%s: c.Err() == %v; want %v", name, e, context.DeadlineExceeded)
	}
}

func XTestDeadline(t *testing.T) {
	t.Parallel()

	c, lintIgnore := NewContextWithDeadline(NewContext(), time.Now().Add(shortDuration))
	_ = lintIgnore
	if got, prefix := fmt.Sprint(c), "bucharest.DefaultContext.WithDeadline("; !strings.HasPrefix(got, prefix) {
		t.Errorf("c.String() = %q want prefix %q", got, prefix)
	}
	testDeadline(c, "WithDeadline", t)

	c, lintIgnore = NewContextWithDeadline(NewContext(), time.Now().Add(shortDuration))
	_ = lintIgnore
	o := otherContext{c}
	testDeadline(o, "WithDeadline+otherContext", t)

	c, lintIgnore = NewContextWithDeadline(NewContext(), time.Now().Add(shortDuration))
	_ = lintIgnore
	o = otherContext{c}
	c, lintIgnore = NewContextWithDeadline(o, time.Now().Add(veryLongDuration))
	_ = lintIgnore
	testDeadline(c, "WithDeadline+otherContext+WithDeadline", t)

	c, lintIgnore = NewContextWithDeadline(NewContext(), time.Now().Add(-shortDuration))
	_ = lintIgnore
	testDeadline(c, "WithDeadline+inthepast", t)

	c, lintIgnore = NewContextWithDeadline(NewContext(), time.Now())
	_ = lintIgnore
	testDeadline(c, "WithDeadline+now", t)
}

func XTestTimeout(t *testing.T) {
	t.Parallel()

	c, lintIgnore := NewContextWithTimeout(NewContext(), shortDuration)
	_ = lintIgnore
	if got, prefix := fmt.Sprint(c), "bucharest.DefaultContext.WithDeadline("; !strings.HasPrefix(got, prefix) {
		t.Errorf("c.String() = %q want prefix %q", got, prefix)
	}
	testDeadline(c, "WithTimeout", t)

	c, lintIgnore = NewContextWithTimeout(NewContext(), shortDuration)
	_ = lintIgnore
	o := otherContext{c}
	testDeadline(o, "WithTimeout+otherContext", t)

	c, lintIgnore = NewContextWithTimeout(NewContext(), shortDuration)
	_ = lintIgnore
	o = otherContext{c}
	c, lintIgnore = NewContextWithTimeout(o, veryLongDuration)
	_ = lintIgnore
	testDeadline(c, "WithTimeout+otherContext+WithTimeout", t)
}

func XTestCanceledTimeout(t *testing.T) {
	c, lintIgnore := NewContextWithTimeout(NewContext(), time.Second)
	_ = lintIgnore
	o := otherContext{c}
	c, cancel := NewContextWithTimeout(o, veryLongDuration)
	cancel() // Should propagate synchronously.
	select {
	case <-c.Done():
	default:
		t.Errorf("<-c.Done() blocked, but shouldn't have")
	}
	if e := c.Err(); e != context.Canceled {
		t.Errorf("c.Err() == %v want %v", e, context.Canceled)
	}
}

type key1 int
type key2 int

var k1 = key1(1)
var k2 = key2(1) // same int as k1, different type
var k3 = key2(3) // same type as k2, different int

func XTestValues(t *testing.T) {
	check := func(c Context, nm, v1, v2, v3 string) {
		if v, ok := c.Value(k1).(string); ok == (len(v1) == 0) || v != v1 {
			t.Errorf(`%s.Value(k1).(string) = %q, %t want %q, %t`, nm, v, ok, v1, len(v1) != 0)
		}
		if v, ok := c.Value(k2).(string); ok == (len(v2) == 0) || v != v2 {
			t.Errorf(`%s.Value(k2).(string) = %q, %t want %q, %t`, nm, v, ok, v2, len(v2) != 0)
		}
		if v, ok := c.Value(k3).(string); ok == (len(v3) == 0) || v != v3 {
			t.Errorf(`%s.Value(k3).(string) = %q, %t want %q, %t`, nm, v, ok, v3, len(v3) != 0)
		}
	}

	c0 := NewContext()
	check(c0, "c0", "", "", "")

	c1 := NewContextWithValue(NewContext(), k1, "c1k1")
	check(c1, "c1", "c1k1", "", "")

	if got, want := fmt.Sprint(c1), `bucharest.DefaultContext.WithValue(type bucharest.key1, val c1k1)`; got != want {
		t.Errorf("c.String() = %q want %q", got, want)
	}

	c2 := NewContextWithValue(c1, k2, "c2k2")
	check(c2, "c2", "c1k1", "c2k2", "")

	c3 := NewContextWithValue(c2, k3, "c3k3")
	check(c3, "c2", "c1k1", "c2k2", "c3k3")

	c4 := NewContextWithValue(c3, k1, nil)
	check(c4, "c4", "", "c2k2", "c3k3")

	o0 := otherContext{NewContext()}
	check(o0, "o0", "", "", "")

	o1 := otherContext{NewContextWithValue(NewContext(), k1, "c1k1")}
	check(o1, "o1", "c1k1", "", "")

	o2 := NewContextWithValue(o1, k2, "o2k2")
	check(o2, "o2", "c1k1", "o2k2", "")

	o3 := otherContext{c4}
	check(o3, "o3", "", "c2k2", "c3k3")

	o4 := NewContextWithValue(o3, k3, nil)
	check(o4, "o4", "", "c2k2", "")
}

func XTestAllocs(t *testing.T, testingShort func() bool, testingAllocsPerRun func(int, func()) float64) {
	bg := NewContext()
	for _, test := range []struct {
		desc       string
		f          func()
		limit      float64
		gccgoLimit float64
	}{
		{
			desc:       "NewContext()",
			f:          func() { NewContext() },
			limit:      0,
			gccgoLimit: 0,
		},
		{
			desc: fmt.Sprintf("WithValue(bg, %v, nil)", k1),
			f: func() {
				c := NewContextWithValue(bg, k1, nil)
				c.Value(k1)
			},
			limit:      3,
			gccgoLimit: 3,
		},
		{
			desc: "WithTimeout(bg, 1*time.Nanosecond)",
			f: func() {
				c, lintIgnore := NewContextWithTimeout(bg, 1*time.Nanosecond)
				_ = lintIgnore
				<-c.Done()
			},
			limit:      12,
			gccgoLimit: 15,
		},
		{
			desc: "WithCancel(bg)",
			f: func() {
				c, cancel := NewContextWithCancel(bg)
				cancel()
				<-c.Done()
			},
			limit:      5,
			gccgoLimit: 8,
		},
		{
			desc: "WithTimeout(bg, 5*time.Millisecond)",
			f: func() {
				c, cancel := NewContextWithTimeout(bg, 5*time.Millisecond)
				cancel()
				<-c.Done()
			},
			limit:      8,
			gccgoLimit: 25,
		},
	} {
		limit := test.limit
		if runtime.Compiler == "gccgo" {
			// gccgo does not yet do escape analysis.
			// TODO(iant): Remove this when gccgo does do escape analysis.
			limit = test.gccgoLimit
		}
		numRuns := 100
		if testingShort() {
			numRuns = 10
		}
		if n := testingAllocsPerRun(numRuns, test.f); n > limit {
			t.Errorf("%s allocs = %f want %d", test.desc, n, int(limit))
		}
	}
}

func XTestSimultaneousCancels(t *testing.T) {
	root, cancel := NewContextWithCancel(NewContext())
	m := map[Context]context.CancelFunc{root: cancel}
	q := []Context{root}
	// Create a tree of contexts.
	for len(q) != 0 && len(m) < 100 {
		parent := q[0]
		q = q[1:]
		for i := 0; i < 4; i++ {
			ctx, cancel := NewContextWithCancel(parent)
			m[ctx] = cancel
			q = append(q, ctx)
		}
	}
	// Start all the cancels in a random order.
	var wg sync.WaitGroup
	wg.Add(len(m))
	for _, cancel := range m {
		go func(cancel context.CancelFunc) {
			cancel()
			wg.Done()
		}(cancel)
	}

	d := quiescent(t)
	stuck := make(chan struct{})
	timer := time.AfterFunc(d, func() { close(stuck) })
	defer timer.Stop()

	// Wait on all the contexts in a random order.
	for ctx := range m {
		select {
		case <-ctx.Done():
		case <-stuck:
			buf := make([]byte, 10<<10)
			n := runtime.Stack(buf, true)
			t.Fatalf("timed out after %v waiting for <-ctx.Done(); stacks:\n%s", d, buf[:n])
		}
	}
	// Wait for all the cancel functions to return.
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-stuck:
		buf := make([]byte, 10<<10)
		n := runtime.Stack(buf, true)
		t.Fatalf("timed out after %v waiting for cancel functions; stacks:\n%s", d, buf[:n])
	}
}

func XTestInterlockedCancels(t *testing.T) {
	parent, cancelParent := NewContextWithCancel(NewContext())
	child, cancelChild := NewContextWithCancel(parent)
	go func() {
		<-parent.Done()
		cancelChild()
	}()
	cancelParent()
	d := quiescent(t)
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-child.Done():
	case <-timer.C:
		buf := make([]byte, 10<<10)
		n := runtime.Stack(buf, true)
		t.Fatalf("timed out after %v waiting for child.Done(); stacks:\n%s", d, buf[:n])
	}
}

func XTestLayersCancel(t *testing.T) {
	testLayers(t, time.Now().UnixNano(), false)
}

func XTestLayersTimeout(t *testing.T) {
	testLayers(t, time.Now().UnixNano(), true)
}

func testLayers(t *testing.T, seed int64, testTimeout bool) {
	t.Parallel()

	r := rand.New(rand.NewSource(seed))
	errorf := func(format string, a ...interface{}) {
		t.Errorf(fmt.Sprintf("seed=%d: %s", seed, format), a...)
	}
	const (
		minLayers = 30
	)
	type value int
	var (
		vals      []*value
		cancels   []context.CancelFunc
		numTimers int
		ctx       = NewContext()
	)
	for i := 0; i < minLayers || numTimers == 0 || len(cancels) == 0 || len(vals) == 0; i++ {
		switch r.Intn(3) {
		case 0:
			v := new(value)
			ctx = NewContextWithValue(ctx, v, v)
			vals = append(vals, v)
		case 1:
			var cancel context.CancelFunc
			ctx, cancel = NewContextWithCancel(ctx)
			cancels = append(cancels, cancel)
		case 2:
			var cancel context.CancelFunc
			d := veryLongDuration
			if testTimeout {
				d = shortDuration
			}
			ctx, cancel = NewContextWithTimeout(ctx, d)
			cancels = append(cancels, cancel)
			numTimers++
		}
	}
	checkValues := func(when string) {
		for _, key := range vals {
			if val := ctx.Value(key).(*value); key != val {
				errorf("%s: ctx.Value(%p) = %p want %p", when, key, val, key)
			}
		}
	}
	if !testTimeout {
		select {
		case <-ctx.Done():
			errorf("ctx should not be canceled yet")
		default:
		}
	}
	if s, prefix := fmt.Sprint(ctx), "bucharest.DefaultContext."; !strings.HasPrefix(s, prefix) {
		t.Errorf("ctx.String() = %q want prefix %q", s, prefix)
	}
	t.Log(ctx)
	checkValues("before cancel")
	if testTimeout {
		d := quiescent(t)
		timer := time.NewTimer(d)
		defer timer.Stop()
		select {
		case <-ctx.Done():
		case <-timer.C:
			errorf("ctx should have timed out after %v", d)
		}
		checkValues("after timeout")
	} else {
		cancel := cancels[r.Intn(len(cancels))]
		cancel()
		select {
		case <-ctx.Done():
		default:
			errorf("ctx should be canceled")
		}
		checkValues("after cancel")
	}
}

func XTestCancelRemoves(t *testing.T) {
	checkChildren := func(when string, ctx Context, want int) {
		if got := len(ctx.(*cancelCtx).children); got != want {
			t.Errorf("%s: context has %d children, want %d", when, got, want)
		}
	}

	ctx, lintIgnore := NewContextWithCancel(NewContext())
	_ = lintIgnore
	checkChildren("after creation", ctx, 0)
	_, cancel := NewContextWithCancel(ctx)
	checkChildren("with WithCancel child ", ctx, 1)
	cancel()
	checkChildren("after canceling WithCancel child", ctx, 0)

	ctx, lintIgnore = NewContextWithCancel(NewContext())
	_ = lintIgnore
	checkChildren("after creation", ctx, 0)
	_, cancel = NewContextWithTimeout(ctx, 60*time.Minute)
	checkChildren("with WithTimeout child ", ctx, 1)
	cancel()
	checkChildren("after canceling WithTimeout child", ctx, 0)
}

func XTestWithCancelCanceledParent(t *testing.T) {
	parent, pcancel := NewContextWithCancel(NewContext())
	pcancel()

	c, lintIgnore := NewContextWithCancel(parent)
	_ = lintIgnore
	select {
	case <-c.Done():
	default:
		t.Errorf("child not done immediately upon construction")
	}
	if got, want := c.Err(), context.Canceled; got != want {
		t.Errorf("child not canceled; got = %v, want = %v", got, want)
	}
}

func XTestWithValueChecksKey(t *testing.T) {

	panicVal := recoveredValue(func() {
		_ = NewContextWithValue(NewContext(), []byte("foo"), "bar")
	})
	if panicVal == nil {
		t.Error("expected panic")
	}
	panicVal = recoveredValue(func() {
		_ = NewContextWithValue(NewContext(), nil, "bar")
	})
	if got, want := fmt.Sprint(panicVal), "nil key"; got != want {
		t.Errorf("panic = %q; want %q", got, want)
	}
}

func XTestInvalidDerivedFail(t *testing.T) {
	panicVal := recoveredValue(func() {
		_, lintIgnore := NewContextWithCancel(nil)
		_ = lintIgnore
	})
	if panicVal == nil {
		t.Error("expected panic")
	}
	panicVal = recoveredValue(func() {
		_, lintIgnore := NewContextWithDeadline(nil, time.Now().Add(shortDuration))
		_ = lintIgnore
	})
	if panicVal == nil {
		t.Error("expected panic")
	}
	panicVal = recoveredValue(func() { _ = NewContextWithValue(nil, "foo", "bar") })
	if panicVal == nil {
		t.Error("expected panic")
	}
}

func recoveredValue(fn func()) (v interface{}) {
	defer func() { v = recover() }()
	fn()
	return
}

func XTestDeadlineExceededSupportsTimeout(t *testing.T) {
	i, ok := context.DeadlineExceeded.(interface {
		Timeout() bool
	})
	if !ok {
		t.Fatal("context.DeadlineExceeded does not support Timeout interface")
	}
	if !i.Timeout() {
		t.Fatal("wrong value for timeout")
	}
}

type myCtx struct {
	Context
}

type myDoneCtx struct {
	Context
}

func (d *myDoneCtx) Done() <-chan struct{} {
	c := make(chan struct{})
	return c
}

func XTestCustomContextGoroutines(t *testing.T) {
	g := atomic.LoadInt32(&goroutines)
	checkNoGoroutine := func() {
		t.Helper()
		now := atomic.LoadInt32(&goroutines)
		if now != g {
			t.Fatalf("%d goroutines created", now-g)
		}
	}
	checkCreatedGoroutine := func() {
		t.Helper()
		now := atomic.LoadInt32(&goroutines)
		if now != g+1 {
			t.Fatalf("%d goroutines created, want 1", now-g)
		}
		g = now
	}

	_, cancel0 := NewContextWithCancel(&myDoneCtx{NewContext()})
	cancel0()
	checkCreatedGoroutine()

	_, cancel0 = NewContextWithTimeout(&myDoneCtx{NewContext()}, veryLongDuration)
	cancel0()
	checkCreatedGoroutine()

	checkNoGoroutine()
	defer checkNoGoroutine()

	ctx1, cancel1 := NewContextWithCancel(NewContext())
	defer cancel1()
	checkNoGoroutine()

	ctx2 := &myCtx{ctx1}
	ctx3, cancel3 := NewContextWithCancel(ctx2)
	defer cancel3()
	checkNoGoroutine()

	_, cancel3b := NewContextWithCancel(&myDoneCtx{ctx2})
	defer cancel3b()
	checkCreatedGoroutine() // ctx1 is not providing Done, must not be used

	ctx4, cancel4 := NewContextWithTimeout(ctx3, veryLongDuration)
	defer cancel4()
	checkNoGoroutine()

	ctx5, cancel5 := NewContextWithCancel(ctx4)
	defer cancel5()
	checkNoGoroutine()

	cancel5()
	checkNoGoroutine()

	_, cancel6 := NewContextWithTimeout(ctx5, veryLongDuration)
	defer cancel6()
	checkNoGoroutine()
}

func XTestCallDefaultContextOptions(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatalf("NewContext returned nil")
	}
	select {
	case x := <-ctx.Done():
		t.Errorf("<-c.Done() == %v want nothing (it should block)", x)
	default:
	}
	if got, want := fmt.Sprint(ctx), "bucharest.DefaultContext"; got != want {
		t.Errorf("NewContext().String() = %q want %q", got, want)
	}

	utils.AssertPanic(t, func() { ctx.ENV() }, ErrNoENV)
	utils.AssertPanic(t, func() { ctx.GORM() }, ErrNoGORM)
	utils.AssertPanic(t, func() { ctx.Log() }, ErrNoLogrus)
	utils.AssertPanic(t, func() { ctx.Redis() }, ErrNoRedis)
	utils.AssertPanic(t, func() { ctx.SQL() }, ErrNoSQL)
	utils.AssertPanic(t, func() { ctx.SQLX() }, ErrNoSQLX)
}

func XTestAddValuesToContext(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatalf("NewContext returned nil")
	}
	select {
	case x := <-ctx.Done():
		t.Errorf("<-c.Done() == %v want nothing (it should block)", x)
	default:
	}
	if got, want := fmt.Sprint(ctx), "bucharest.DefaultContext"; got != want {
		t.Errorf("NewContext().String() = %q want %q", got, want)
	}

	key1 := 1
	value1 := "one"
	key2 := "two"
	value2 := 2

	ctx = AddValuesToContext(ctx, MapAny{
		key1: value1,
		key2: value2,
	})

	assert.Equal(t, value1, ctx.Value(key1))
	assert.Equal(t, value2, ctx.Value(key2))
}
