package mvm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetG(t *testing.T) {
	for i := 0; i < 100; i++ {
		go func() {
			gid := getGID()
			assert.Greater(t, gid, int64(0))
		}()
	}
}
