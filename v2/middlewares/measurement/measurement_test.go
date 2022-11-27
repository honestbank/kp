//go:build integration_test

package measurement_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/honestbank/kp/v2/middlewares/measurement"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/assert"
)

func TestMeasure(t *testing.T) {
	t.Run("calls next", func(t *testing.T) {
		called := false
		measurement.NewMeasurementMiddleware("localhost:9091", "integration_test").Process(context.Background(), nil, func(ctx context.Context, msg *kafka.Message) error {
			called = true
			time.Sleep(time.Millisecond * 550)
			return nil
		})
		assert.True(t, called)
	})
	t.Run("works if message errors", func(t *testing.T) {
		measurement.NewMeasurementMiddleware("localhost:9091", "integration_test").Process(context.Background(), nil, func(ctx context.Context, msg *kafka.Message) error {
			return errors.New("some error")
		})
	})
	t.Run("works even when push gateway url is not set", func(t *testing.T) {
		called := false
		measurement.NewMeasurementMiddleware("", "integration_test").Process(context.Background(), nil, func(ctx context.Context, msg *kafka.Message) error {
			called = true
			time.Sleep(time.Millisecond * 550)
			return errors.New("some random error")
		})
		assert.True(t, called)
	})
	time.Sleep(time.Second * 6)
	t.Run("pushes to prometheus", func(t *testing.T) {
		success200, _ := regexp.Compile(`kp_operation_duration_milliseconds_bucket\{.+success.+le="200"\}.+`)
		success500, _ := regexp.Compile(`kp_operation_duration_milliseconds_bucket\{.+success.+le="1000"\}.+`)
		res, err := http.Get("http://localhost:9091/metrics")
		assert.NoError(t, err)
		defer res.Body.Close()
		bytes, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		scrapedValues := string(bytes)
		successMatches200 := success200.FindStringSubmatch(scrapedValues)
		successMatches500 := success500.FindStringSubmatch(scrapedValues)
		assert.Len(t, successMatches200, 1)
		assert.True(t, strings.HasSuffix(successMatches200[0], "0"))
		assert.Len(t, successMatches500, 1)
		assert.True(t, strings.HasSuffix(successMatches500[0], "1"))
	})
}
