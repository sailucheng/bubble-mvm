package mvm

import (
	"fmt"
	"log"
	"reflect"
	"sync"
	"sync/atomic"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	ContextPool = sync.Pool{
		New: func() any {
			return &Context{}
		},
	}
	ControllerMethodCache = CreateMethodInvokerCache()
)

type event struct {
	name string
	args []any
}

// Context holds the current state of the application during pipeline processa
// It contains fields for message handling, the business logic model, the rendering logic,
// and the result of the previous operations.
// This struct is used as a container to manage the application's flow and communication between components.
type Context struct {
	// @see tea.Msg for more details about message processing.
	// Msg represents the current message being processed in the application.
	// It is of type tea.Msg and drives the application's logic flow.
	Msg tea.Msg
	// TeaModel is used internally for handling tea bubble processing.
	// It serves the same purpose as the model in tea bubble, but with the additional MVM pipeline logic.
	// In most cases, you don't need to worry about this field.
	TeaModel Model
	// Model represents the business logic model.
	// The view should be rendered based on the state of this field.
	Model any
	// Viewer is responsible for the rendering logic of the interface.
	// It is used to display the view corresponding to the Model.
	Viewer Viewer
	// aborted is used to indicate if the process was aborted.
	aborted atomic.Bool
	// Result stores the result returned by the previous pipeline or controller logic.
	// It holds the result of the application’s current state after processing.
	Result *Result
	events []event
}

// Cmd creates and returns a Result with the provided command `cmd`.
// It also includes the current Viewer and Model from the Context.
func (ctx *Context) Cmd(cmd tea.Cmd) Result {
	return Result{
		Cmd:    cmd,
		Viewer: ctx.Viewer,
		Model:  ctx.Model,
	}
}

// NoAction returns a Result with the current Viewer and Model,
// but without any command or error. This can be used when no action is needed.
func (ctx *Context) NoAction() Result {
	return Result{
		Viewer: ctx.Viewer,
		Model:  ctx.Model,
	}
}

// Propagate returns the current Result from the Context if it exists.
// If the Context's Result is nil, it returns a default Result created
// using the NoAction method.
func (ctx *Context) Propagate() Result {
	if ctx.Result != nil {
		return *ctx.Result
	}
	return ctx.NoAction()
}

// View creates a Result with a new Viewer `v`, leaving the Model and Command unchanged.
// This can be used to render a new view while keeping the same model and command.
func (ctx *Context) View(v Viewer) Result {
	return Result{
		Viewer: v,
		Cmd:    nil,
		Model:  ctx.Model,
	}
}

// WithError creates a Result containing an error `err` along with the current Viewer,
// and Model. This is useful for handling error scenarios and propagating them through the system.
func (ctx *Context) WithError(err error) Result {
	return Result{
		Viewer: ctx.Viewer,
		Err:    err,
		Model:  ctx.Model,
	}
}

// Redirect creates a Result with a new Viewer `v` and a new Model `model`.
// The existing command and error from the previous result are preserved if they exist.
func (ctx *Context) Redirect(model any, v Viewer) Result {
	result := Result{
		Viewer: v,
		Model:  model,
		Cmd:    nil,
	}
	if ctx.Result != nil {
		result.Err = ctx.Result.Err
		result.Cmd = ctx.Result.Cmd
	}
	return result
}

// Quit creates a Result that signals the termination of the program by returning a Quit command.
// It also ensures the current Viewer is preserved for the result.
func (ctx *Context) Quit() Result {
	ctx.Abort()
	var v Viewer = ctx.Viewer

	if ctx.Result != nil {
		v = ctx.Result.Viewer
	}
	return Result{
		Viewer: v,
		Cmd:    tea.Quit,
		Model:  ctx.Model,
	}
}

// None returns a predefined "no-op" Result that indicates no operation is performed.
// This is useful when no action is required or if you need a placeholder result.
func (ctx *Context) None() Result {
	return nopeResult
}

func (ctx *Context) Abort() {
	ctx.aborted.Store(true)
}

func (ctx *Context) IsAbort() bool {
	return ctx.aborted.Load()
}

func (ctx *Context) Trigger(method string, args ...any) {
	ctx.addCallback(method, args...)
}

func (ctx *Context) addCallback(method string, args ...any) {
	ctx.events = append(ctx.events, event{method, args})
}

func triggerControllerMethod(controller Controller, method string, args ...any) (*Result, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Method invocation panic: %v", r)
		}
	}()

	invoker, _, err := ControllerMethodCache.GetOrAdd(controller, method)
	if err != nil {
		return nil, err
	}
	values := MakeReflectValues(args...)

	var rets []reflect.Value
	if invoker.Type().IsVariadic() {
		rets = invoker.CallSlice(values)
	} else {
		rets = invoker.Call(values)
	}

	if len(rets) == 0 {
		return &Result{}, nil
	}

	if result, ok := rets[0].Interface().(Result); ok {
		return &result, nil
	}
	return nil, fmt.Errorf("the first parameter of the 'trigger' must be of type 'mvm.Result")
}

func (ctx *Context) reset() {
	ctx.aborted.Store(false)
	ctx.TeaModel = Model{}
	ctx.Msg = nil
	ctx.Result = nil
	ctx.Model = nil
	ctx.Viewer = nil
	if ctx.events != nil && cap(ctx.events) < 256 {
		ctx.events = ctx.events[:0]
	} else {
		ctx.events = nil
	}
}

func applyStates(ctx *Context) {
	if ctx.Result == nil {
		return
	}
	if ctx.Result.Model != nil {
		ctx.Model = ctx.Result.Model
	}
	if ctx.Result.Viewer != nil {
		ctx.Viewer = ctx.Result.Viewer
	}

	ctx.TeaModel.Model = ctx.Model
	ctx.TeaModel.viewer = ctx.Viewer
}

// fillTeaModel assigns the current Model and Viewer to the TeaModel internal.
// This method ensures that the TeaModel is populated with the current state.
func (ctx *Context) fillTeaModel() {
	// Only fill TeaModel if Result exists
	applyStates(ctx)
}
