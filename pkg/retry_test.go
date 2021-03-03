package pkg_test

import (
	"testing"
	"time"

	"github.com/kyma-incubator/milv/pkg"
	"github.com/stretchr/testify/assert"
)

func TestLimit(t *testing.T) {
	//GIVEN
	backoff := 1 * time.Second
	waiter := pkg.NewWaiter(backoff)

	before := time.Now()
	expected := before.Add(backoff)
	//WHEN
	waiter.Wait()

	//THEN
	afterExecution := time.Now()
	assert.True(t, afterExecution.After(expected))
}
