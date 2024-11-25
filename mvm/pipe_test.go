package mvm

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMiddlewareSequence(t *testing.T) {
	ctx := Context{}

	times := 100
	var b []int

	mws := generate(times, func(index int) Middleware {
		return testMiddleware{index, &b}
	})

	pip := CreatePipe(withMws(mws...))
	pip.Execute(&ctx)

	var expected []int
	var rev []int
	for i := 0; i < times; i++ {
		expected = append(expected, i)
		rev = append(rev, i)
	}

	slices.Reverse(rev)
	expected = append(expected, rev...)
	assert.Equal(t, expected, b)
}

type testMiddleware struct {
	i int
	b *[]int
}

func (mw testMiddleware) Execute(ctx *Context, next MiddlewareFunc) {
	*mw.b = append(*mw.b, mw.i)
	next(ctx)
	*mw.b = append(*mw.b, mw.i)
}
func TestController(t *testing.T) {
	times := 5
	var m int
	abortIndex := times - 1
	controllers := generate(times, func(index int) Controller {
		return pipTestController{
			t:     t,
			i:     index,
			abort: abortIndex == index,
		}
	})

	pipe := CreatePipe(withPipeControllers(controllers...))

	c := Context{
		Model: &m,
	}

	pipe.Execute(&c)

	assert.Equal(t, abortIndex, m+1)
}

type pipTestController struct {
	t     *testing.T
	abort bool
	i     int
}

func (c pipTestController) Filter(*Context) bool {
	return true
}

func (c pipTestController) Handle(ctx *Context) Result {
	if c.abort {
		ctx.Abort()
		return Result{}
	}
	t := c.t
	p, ok := ctx.Model.(*int)
	assert.True(t, ok)
	if c.i > 0 {
		assert.Equal(t, c.i-1, *p)
	}
	*p = c.i
	return Result{
		Model: p,
	}
}

func generate[T any](n int, f func(index int) T) []T {
	s := make([]T, n)
	for i := 0; i < n; i++ {
		s[i] = f(i)
	}
	return s
}
