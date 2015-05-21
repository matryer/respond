package main

import (
	"log"
	"net/http"

	"gopkg.in/matryer/respond.v1"
)

const version = "1.0"

func main() {

	fn := http.HandlerFunc(handleSomething)

	opts := &respond.Options{
		Before: func(w http.ResponseWriter, r *http.Request, status int, data interface{}) (int, interface{}) {
			w.Header().Set("X-API-Version", version)
			return status, map[string]interface{}{
				"status": status,
				"data":   data,
			}
		},
		After: func(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
			log.Println("<-", status, data)
		},
	}

	log.Println("check out localhost:8080")
	http.ListenAndServe("localhost:8080", opts.Handler(fn))

}

func handleSomething(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{"message": "Hello"}
	respond.With(w, r, http.StatusOK, data)
}
