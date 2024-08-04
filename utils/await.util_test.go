package utils_test

import (
	"testing"
	"time"

	. "github.com/argonlab-io/bucharest/utils"
	"github.com/stretchr/testify/assert"
)

func TestAwaitErrorTimeout(t *testing.T) {
	err := Await(&AwaitOptions{
		Condition: func() bool {
			return false
		},
		Timeout: 1 * time.Second,
	})
	if assert.Error(t, err) {
		assert.ErrorContains(t, err, "timeout")
	}
}

func TestAwaitCompleted(t *testing.T) {
	err := Await(&AwaitOptions{
		Condition: func() bool {
			return true
		},
	})
	assert.NoError(t, err)
}
