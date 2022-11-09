package handlers

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/honestbank/kp/v2/producer"
	"github.com/honestbank/kp_demo/types"
)

//go:embed templates
var fs embed.FS

func RenderForm() http.HandlerFunc {
	f, err := fs.Open("templates/index.html")
	if err != nil {
		panic(err)
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		all, err := ioutil.ReadAll(f)
		if err != nil {
			writer.WriteHeader(500)
			writer.Write([]byte(err.Error()))
		}
		writer.WriteHeader(200)
		writer.Header().Set("content-type", "text/html")
		writer.Write(all)
	}
}

func Save(producer producer.Producer[types.RegistrationEvent]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body types.RegistrationEvent
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			w.WriteHeader(400)
			fmt.Println(err.Error())
			return
		}
		body.RegisteredAt = time.Now().Format(time.RFC3339)
		err = producer.Produce(r.Context(), body)
		if err != nil {
			w.WriteHeader(500)
			fmt.Println(err.Error())
			return
		}
		w.WriteHeader(200)
		_, _ = w.Write([]byte("{\"status\": \"ok\"}"))
	}
}
