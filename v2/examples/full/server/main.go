package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/honestbank/kp/v2/producer"
	"github.com/honestbank/kp_demo/server/handlers"
	"github.com/honestbank/kp_demo/types"
)

func StartServer() {
	p, err := producer.New[types.RegistrationEvent]("user-registration-complete")
	if err != nil {
		panic(err)
	}
	r := http.NewServeMux()
	r.Handle("/", handlers.RenderForm())
	r.Handle("/save", handlers.Save(p))
	handler := handlers.Tracing(r)
	httpServer := &http.Server{
		Addr:         ":8888",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      handler,
	}
	fmt.Println("listening at: http://localhost:8888")
	err = httpServer.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
