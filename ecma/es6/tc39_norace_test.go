// +build !race

package es6

import "testing"

func (ctx *tc39TestCtx) runTest(name string, f func(t *testing.T)) {
	ctx.t.Run(name, func(t *testing.T) {
		t.Parallel()
		f(t)
	})
}

func (ctx *tc39TestCtx) flush() {
}
