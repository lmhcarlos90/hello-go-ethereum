// Copyright (c) 2022 Matheus Degiovani
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package testctx

import (
	"context"
	"sync"
	"time"
)

var (
	mtx            sync.Mutex
	defaultTimeout = 3 * time.Minute
	pkgCtx         = context.Background()
)

// Cleaner is an interface that is satisfied by both testing.T and testing.B
// with the needed methods for this package to work.
type Cleaner interface {
	Cleanup(func())
}

func defaults() (parent context.Context, timeout time.Duration) {
	mtx.Lock()
	timeout = defaultTimeout
	parent = pkgCtx
	mtx.Unlock()
	return
}

// SetDefaultTimeout sets the default timeout duration for new contexts at
// the package level. While this is safe for concurrent access, it is usually
// an error to use different timeout values for different tests of the
// same package as it can lead to data races. If a particular test or operation
// requires a different timeout value, it is better to call [WithTimeout] and
// use the returned context explicitly.
func SetDefaultTimeout(t time.Duration) {
	if t < 0 {
		panic("negative duration")
	}
	mtx.Lock()
	defaultTimeout = t
	mtx.Unlock()
}

// New returns a new context that times out at the default timeout and is
// canceled once the test is done.
func New(t Cleaner) context.Context {
	parent, timeout := defaults()
	ctx, cancel := context.WithTimeout(parent, timeout)
	t.Cleanup(cancel)
	return ctx
}

// WithTimeout returns a new context that times out at the specified timeout.
func WithTimeout(t Cleaner, timeout time.Duration) context.Context {
	parent, _ := defaults()
	ctx, cancel := context.WithTimeout(parent, timeout)
	t.Cleanup(cancel)
	return ctx
}

// WithParent returns a new context with the specified parent that times out
// at the default timeout value.
func WithParent(t Cleaner, parent context.Context) context.Context {
	_, timeout := defaults()
	ctx, cancel := context.WithTimeout(parent, timeout)
	t.Cleanup(cancel)
	return ctx
}

// WithCancel returns a new context with the default timeout value and a cancel
// function that can be used for early termination.
func WithCancel(t Cleaner) (context.Context, func()) {
	parent, timeout := defaults()
	ctx, cancel := context.WithTimeout(parent, timeout)
	t.Cleanup(cancel)
	return ctx, cancel
}
