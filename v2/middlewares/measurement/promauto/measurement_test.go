package promauto_test

import (
	"context"
	"errors"
	"io"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/honestbank/kp/v2/middlewares/measurement/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
)

func TestMeasure(t *testing.T) {
	t.Run("calls next", func(t *testing.T) {
		called := false
		promauto.NewMeasurementMiddleware("integration_test").Process(context.Background(), nil, func(ctx context.Context, msg *kafka.Message) error {
			called = true
			time.Sleep(time.Millisecond * 550)
			return nil
		})
		assert.True(t, called)
	})
	t.Run("works if message errors", func(t *testing.T) {
		promauto.NewMeasurementMiddleware("integration_test").Process(context.Background(), nil, func(ctx context.Context, msg *kafka.Message) error {
			return errors.New("some error")
		})
	})
	t.Run("works even when push gateway url is not set", func(t *testing.T) {
		called := false
		promauto.NewMeasurementMiddleware("integration_test").Process(context.Background(), nil, func(ctx context.Context, msg *kafka.Message) error {
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
		recorder := httptest.NewRecorder()
		handler := promhttp.Handler()
		promauto.NewMeasurementMiddleware("integration_test").Process(context.Background(), nil, func(ctx context.Context, msg *kafka.Message) error {
			time.Sleep(time.Millisecond * 550)

			return nil
		})
		handler.ServeHTTP(recorder, httptest.NewRequest("GET", "/metrics", nil))
		bytes, err := io.ReadAll(recorder.Result().Body)
		assert.NoError(t, err)
		scrapedValues := string(bytes)
		successMatches200 := success200.FindStringSubmatch(scrapedValues)
		successMatches500 := success500.FindStringSubmatch(scrapedValues)
		assert.Len(t, successMatches200, 1)
		assert.True(t, strings.HasSuffix(successMatches200[0], "0"))
		assert.Len(t, successMatches500, 1)
		assert.False(t, strings.HasSuffix(successMatches500[0], "0"))
	})
}
