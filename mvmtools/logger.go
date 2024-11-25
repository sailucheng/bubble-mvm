package main

import (
	"fmt"
	"io"
	"os"
	"sync"
)

func loggerFactory(verbose bool) *logger {
	var w io.Writer = os.Stdout
	var errw io.Writer = os.Stderr
	var mu sync.Mutex
	if !verbose {
		w = io.Discard
	}

	return &logger{
		mu:      &mu,
		verbose: verbose,
		w:       w,
		errw:    errw,
	}
}
func (l *logger) err(format string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintf(l.errw, format+"\n", args...)
}
func (l *logger) log(format string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintf(l.w, format+"\n", args...)
}

type logger struct {
	verbose bool
	mu      *sync.Mutex
	w       io.Writer
	errw    io.Writer
}

type Logger interface {
	log(format string, args ...any)
	err(format string, args ...any)
}
