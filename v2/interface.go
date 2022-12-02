package v2

import (
	"github.com/honestbank/kp/v2/internal/middleware"
)

type KafkaProcessor[MessageType any] interface {
	AddMiddleware(middleware middleware.Middleware[*MessageType, error]) KafkaProcessor[MessageType]
	Stop()
	Run(processor Processor[MessageType]) error
}
