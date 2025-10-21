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
	_, cancel, retryTimeout, _ := computeRetryContext(ctx, time.Minute, 2*time.Minute)
	defer cancel()

	assert.Equal(t, 2*time.Minute, retryTimeout)
}

func TestComputeRetryContext_WithDeadline_Plenty(t *testing.T) {
	// caller provides 5 minutes; timeoutDeadlineBuffer is 1 minute -> retryTimeout should be ~4 minutes
	parent := context.Background()
	deadline := time.Now().Add(5 * time.Minute)
	ctxWithDL, cancelParent := context.WithDeadline(parent, deadline)
	defer cancelParent()

	_, cancel, retryTimeout, _ := computeRetryContext(ctxWithDL, time.Minute, 2*time.Minute)
	defer cancel()

	// expect ~4 minutes
	expected := 4 * time.Minute
	assert.True(t, approx(expected, retryTimeout), "expected approximately %v, got %v", expected, retryTimeout)
}

func TestComputeRetryContext_ExpiredDeadline(t *testing.T) {
	ctxWithExpiredDL, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()

	_, _, _, err := computeRetryContext(ctxWithExpiredDL, time.Minute, 2*time.Minute)

	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestClassifyRetryError(t *testing.T) {
	tests := []struct {
		err            error
		expectedSubstr string
	}{
		{err: rest.Error{Code: 400, Message: "bad"}, expectedSubstr: "not yet usable (HTTP 400)"},
		{err: rest.Error{Code: 404, Message: "notfound"}, expectedSubstr: "not yet usable (HTTP 404)"},
		{err: rest.Error{Code: 401, Message: "unauth"}, expectedSubstr: "unusable (HTTP 401)"},
		{err: rest.Error{Code: 500, Message: "server"}, expectedSubstr: "not yet usable (HTTP 500)"},
		{err: errors.New("network failure"), expectedSubstr: "not yet usable"},
	}

	for _, tc := range tests {
		t.Run(tc.expectedSubstr, func(t *testing.T) {
			re := classifyRetryError(tc.err)
			msg := re.Err.Error()

			assert.Contains(t, msg, tc.expectedSubstr)
		})
	}
}
