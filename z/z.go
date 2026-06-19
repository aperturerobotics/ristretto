/*
 * SPDX-FileCopyrightText: © 2017-2025 Istari Digital, Inc.
 * SPDX-License-Identifier: Apache-2.0
 */

package z

import (
	"context"
	"sync"
)

// Key constrains the types ristretto can hash into a cache key. KeyToHash, the
// hashing entry point, is build-tagged: native builds use reflect to resolve
// named key types, and the js browser build uses a reflect-free type switch.
type Key interface {
	~uint64 | ~string | ~[]byte | ~byte | ~int | ~uint | ~int32 | ~uint32 | ~int64
}

var (
	dummyCloserChan <-chan struct{}
	tmpDir          string
)

// Closer holds the two things we need to close a goroutine and wait for it to
// finish: a chan to tell the goroutine to shut down, and a WaitGroup with
// which to wait for it to finish shutting down.
type Closer struct {
	waiting sync.WaitGroup

	ctx    context.Context
	cancel context.CancelFunc
}

// SetTmpDir sets the temporary directory for the temporary buffers.
func SetTmpDir(dir string) {
	tmpDir = dir
}

// NewCloser constructs a new Closer, with an initial count on the WaitGroup.
func NewCloser(initial int) *Closer {
	ret := &Closer{}
	ret.ctx, ret.cancel = context.WithCancel(context.Background())
	ret.waiting.Add(initial)
	return ret
}

// AddRunning Add()'s delta to the WaitGroup.
func (lc *Closer) AddRunning(delta int) {
	lc.waiting.Add(delta)
}

// Ctx can be used to get a context, which would automatically get cancelled when Signal is called.
func (lc *Closer) Ctx() context.Context {
	if lc == nil {
		return context.Background()
	}
	return lc.ctx
}

// Signal signals the HasBeenClosed signal.
func (lc *Closer) Signal() {
	// Todo(ibrahim): Change Signal to return error on next badger breaking change.
	lc.cancel()
}

// HasBeenClosed gets signaled when Signal() is called.
func (lc *Closer) HasBeenClosed() <-chan struct{} {
	if lc == nil {
		return dummyCloserChan
	}
	return lc.ctx.Done()
}

// Done calls Done() on the WaitGroup.
func (lc *Closer) Done() {
	if lc == nil {
		return
	}
	lc.waiting.Done()
}

// Wait waits on the WaitGroup. (It waits for NewCloser's initial value, AddRunning, and Done
// calls to balance out.)
func (lc *Closer) Wait() {
	lc.waiting.Wait()
}

// SignalAndWait calls Signal(), then Wait().
func (lc *Closer) SignalAndWait() {
	lc.Signal()
	lc.Wait()
}

// ZeroOut zeroes out all the bytes in the range [start, end).
func ZeroOut(dst []byte, start, end int) {
	if start < 0 || start >= len(dst) {
		return // BAD
	}
	if end >= len(dst) {
		end = len(dst)
	}
	if end-start <= 0 {
		return
	}
	Memclr(dst[start:end])
	// b := dst[start:end]
	// for i := range b {
	// 	b[i] = 0x0
	// }
}
