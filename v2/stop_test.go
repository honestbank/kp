package v2_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/stretchr/testify/assert"

	backoff_policy "github.com/honestbank/backoff-policy"
	v2 "github.com/honestbank/kp/v2"
	"github.com/honestbank/kp/v2/middlewares/backoff"
)

func TestStopEndsTheBackoffWait(t *testing.T) {
	// The first message fails, so the backoff waits an hour before the second message.
	policy := backoff_policy.NewBackoff(func(count int) time.Duration {
		if count == 0 {
			return 0
		}

		return time.Hour
	})
	processor := v2.New[kafka.Message]()
	processor.AddMiddleware(backoff.NewBackoffMiddleware(policy))
	time.AfterFunc(100*time.Millisecond, processor.Stop)

	calls := 0
	start := time.Now()
	err := processor.Run(func(ctx context.Context, item *kafka.Message) error {
		calls++

		return errors.New("some error")
	})

	assert.NoError(t, err)
	assert.Less(t, time.Since(start), 5*time.Second, "Stop ends the wait, so Run does not wait the hour")
	assert.Equal(t, 2, calls, "the message in progress still runs after Stop")
}
