package mvm

import (
	"runtime"
	"sync"
)

type ReentrantMutex struct {
	owner     int64
	recursion int64
	mu        sync.Mutex
}

func (m *ReentrantMutex) Lock() {
	gid := getGID()
	if m.owner == gid {
		m.recursion++
		return
	}

	m.mu.Lock()
	m.recursion = 1
	m.owner = gid
}

func (m *ReentrantMutex) Unlock() {
	g := getGID()
	if m.owner != g {
		panic("wrong goroutine unlocking the mutex")
	}

	m.recursion--
	if m.recursion == 0 {
		m.owner = -1
		m.mu.Unlock()
	}
}

func getGID() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	var gid int64
	for _, b := range buf[:n] {
		if '0' <= b && b <= '9' {
			gid = gid*10 + int64(b-'0')
		} else if gid > 0 {
			break
		}
	}
	return gid
}
