package mvm

import "log"

var DefaultPipe *Pipe

func init() {
	DefaultPipe = CreatePipe()
}

func RegisterControllers(controllers ...Controller) {
	DefaultPipe.controllers = append(DefaultPipe.controllers, controllers...)
}

func RegisterMiddlewares(mws ...Middleware) {
	DefaultPipe.mws = append(DefaultPipe.mws, mws...)
}

type Middleware interface {
	Execute(ctx *Context, next MiddlewareFunc)
}

type MiddlewareFunc func(ctx *Context)

type Pipe struct {
	controllers []Controller
	mws         []Middleware
}
type opt func(*Pipe)

func withPipeControllers(controllers ...Controller) opt {
	return func(p *Pipe) {
		p.controllers = append(p.controllers, controllers...)
	}
}

func withMws(mws ...Middleware) opt {
	return func(p *Pipe) {
		p.mws = append(p.mws, mws...)
	}
}

func CreatePipe(opts ...opt) *Pipe {
	p := Pipe{
		mws:         make([]Middleware, 0),
		controllers: make([]Controller, 0),
	}
	for _, o := range opts {
		o(&p)
	}
	p.mws = append(p.mws, viewerUpdaterMiddleware{})
	p.mws = append(p.mws, controllerMiddleware(func() []Controller {
		return p.controllers
	}))
	return &p
}

func (pip Pipe) Execute(ctx *Context) {
	index := 0
	var next MiddlewareFunc

	next = func(ctx *Context) {
		if index < len(pip.mws) {
			currentIndex := index
			index++
			pip.mws[currentIndex].Execute(ctx, next)

			if ctx.IsAbort() {
				return
			}
		}
	}

	next(ctx)
}

type viewerUpdaterMiddleware struct{}

func (v viewerUpdaterMiddleware) Execute(ctx *Context, next MiddlewareFunc) {
	current := ctx.Viewer
	if current == nil {
		next(ctx)
		return
	}

	result := current.Update(ctx)
	ctx.Result = &result
	ctx.fillTeaModel()

	next(ctx)
}

// end of pipe,drop next
type controllerMiddleware func() []Controller

func (c controllerMiddleware) Execute(ctx *Context, _ MiddlewareFunc) {
	for _, controller := range c() {
		if controller.Filter(ctx) {
			result := callControllerMethod(ctx, controller)
			populateContext(ctx, result)

			if ctx.IsAbort() {
				return
			}

			result = controller.Handle(ctx)
			populateContext(ctx, result)

			if ctx.IsAbort() {
				return
			}
		}
	}
}
func populateContext(ctx *Context, result Result) {
	if result.Viewer == nil && result.Model == nil {
		//all result has this field
		return
	}
	ctx.Result = &result
	ctx.fillTeaModel()
}
func callControllerMethod(ctx *Context, controller Controller) Result {
	var ret Result
	if len(ctx.events) > 0 {
		for _, e := range ctx.events {
			result, err := triggerControllerMethod(controller, e.name, e.args...)
			if err != nil {
				log.Printf("trigger controller method %s failed: %v", e.name, err)
				result.Err = err
				continue
			}
			ret.Composite(*result)
		}
		return ret
	}
	return nopeResult
}
