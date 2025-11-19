//go:build unit

package role_arn

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/stretchr/testify/assert"
)

func approx(d1, d2 time.Duration) bool {
	delta := d1 - d2
	if delta < 0 {
		delta = -delta
	}
	return delta <= 200*time.Millisecond
}

func TestComputeRetryContext_NoDeadline(t *testing.T) {
	ctx := context.Background()
	retryTimeout, _ := computeRetryTimeout(ctx, 2*time.Minute)

	assert.Equal(t, 2*time.Minute, retryTimeout)
}

func TestComputeRetryContext_WithDeadline_Plenty(t *testing.T) {
	// caller provides 5 minutes; retryTimeout should be ~5 minutes
	parent := context.Background()
	deadline := time.Now().Add(5 * time.Minute)
	ctxWithDL, cancelParent := context.WithDeadline(parent, deadline)
	defer cancelParent()

	retryTimeout, _ := computeRetryTimeout(ctxWithDL, 2*time.Minute)

	// expect ~5 minutes
	expected := 5 * time.Minute
	assert.True(t, approx(expected, retryTimeout), "expected approximately %v, got %v", expected, retryTimeout)
}

func TestComputeRetryContext_ExpiredDeadline(t *testing.T) {
	ctxWithExpiredDL, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()

	_, err := computeRetryTimeout(ctxWithExpiredDL, 2*time.Minute)

	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestClassifyRetryError(t *testing.T) {
	tests := []struct {
		err         error
		isRetryable bool
	}{
		{err: rest.Error{Code: 400, Message: "bad"}, isRetryable: true},
		{err: rest.Error{Code: 404, Message: "notfound"}, isRetryable: true},
		{err: rest.Error{Code: 401, Message: "unauth"}, isRetryable: false},
		{err: rest.Error{Code: 500, Message: "server"}, isRetryable: true},
		{err: errors.New("network failure"), isRetryable: false},
	}

	for _, tc := range tests {
		t.Run(tc.err.Error(), func(t *testing.T) {
			retryError := classifyRetryError(tc.err)
			assert.Equal(t, tc.isRetryable, retryError.Retryable)
		})
	}
}
