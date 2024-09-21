package server

import (
	"crypto/tls"
	"net/http"

	"github.com/KelpGF/Go-Deploy-Cloud-Run/internal/handlers"
)

type ServerHttp struct{}

func (*ServerHttp) Run() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Kelp Weather By ZipeCode!"))
	})

	http.HandleFunc("/zip-code/weather", handlers.WeatherByCepHandler)

	http.ListenAndServe(":8080", nil)
}
