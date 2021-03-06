// Copyright 2015 The LUCI Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package clock

import (
	"time"

	"golang.org/x/net/context"
)

type systemTimer struct {
	// ctx is the underlying timer. It starts as nil, and is initialized on Reset.
	ctx context.Context

	// timerC is the timer channel.
	timerC chan TimerResult

	// timerStoppedC is a signal channel used by Stop to alert our monitor that
	// this timer has been manually canceled.
	timerStoppedC chan struct{}
	// timerMonitorResultC returns true if the timer monitor was prematurely
	// terminated, false if not. It is used by our timer monitor to indicate its
	// status when stopped.
	timerMonitorResultC chan bool
}

var _ Timer = (*systemTimer)(nil)

func newSystemTimer(ctx context.Context) Timer {
	return &systemTimer{
		ctx:    ctx,
		timerC: make(chan TimerResult, 1),
	}
}

func (t *systemTimer) GetC() <-chan TimerResult {
	return t.timerC
}

func (t *systemTimer) Reset(d time.Duration) (running bool) {
	running = t.Stop()

	// If our Context is already done, finish immediately.
	if err := t.ctx.Err(); err != nil {
		t.timerC <- TimerResult{Time: time.Now(), Err: t.ctx.Err()}
		return
	}

	// Start a monitor goroutine and our actual timer. Copy our channels, since
	// future stop/reset will change the systemTimer's values and our goroutine
	// should only operate on this round's values.
	t.timerStoppedC = make(chan struct{})
	t.timerMonitorResultC = make(chan bool, 1)

	timerStoppedC, timerMonitorResultC := t.timerStoppedC, t.timerMonitorResultC

	realTimer := time.NewTimer(d)
	go func() {
		defer realTimer.Stop()

		interrupted := false
		defer func() {
			timerMonitorResultC <- interrupted
			close(timerMonitorResultC)
		}()

		select {
		case <-timerStoppedC:
			interrupted = true

		case <-t.ctx.Done():
			t.timerC <- TimerResult{Time: time.Now(), Err: t.ctx.Err()}
		case now := <-realTimer.C:
			t.timerC <- TimerResult{Time: now, Err: t.ctx.Err()}
		}
	}()

	return
}

func (t *systemTimer) Stop() bool {
	if t.timerStoppedC == nil {
		return false
	}
	close(t.timerStoppedC)
	t.timerStoppedC = nil
	return <-t.timerMonitorResultC
}
