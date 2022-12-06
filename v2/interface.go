package v2

import (
	"github.com/honestbank/kp/v2/internal/middleware"
)

type MessageProcessor[MessageType any] interface {
	AddMiddleware(middleware middleware.Middleware[*MessageType, error]) MessageProcessor[MessageType]
	Stop()
	Run(processor Processor[MessageType]) error
}
