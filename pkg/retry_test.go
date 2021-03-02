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
	retry := pkg.NewLimiter(backoff)

	before := time.Now()
	expected := before.Add(backoff)
	//WHEN
	retry.Limit()

	//THEN
	afterExecution := time.Now()
	assert.True(t, afterExecution.After(expected))
	assert.WithinDuration(t, expected, afterExecution, 250*time.Millisecond)
}
