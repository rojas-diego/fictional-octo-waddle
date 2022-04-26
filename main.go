package main

import (
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		value, err := strconv.Atoi(r.URL.Query().Get("value"))
		if err != nil {
			w.WriteHeader(200)
			return
		}
		w.WriteHeader(value)
	})

	// http.HandleFunc("/allocate", func(w http.ResponseWriter, r *http.Request) {
	// 	mb, err := strconv.Atoi(r.URL.Query().Get("mb"))
	// 	if err != nil {
	// 		mb = 32
	// 		return
	// 	}

	// })

	// http.HandleFunc("/compute", func(w http.ResponseWriter, r *http.Request) {
	// 	duration, err := time.ParseDuration(r.URL.Query().Get("for"))
	// 	if err != nil {
	// 		duration = time.Second * 30
	// 	}
	// 	go func(duration time.Duration) {
	// 		ctx, cancel := context.WithTimeout(context.TODO(), duration)
	// 		defer cancel()

	// 		select {
	// 		case <-ctx.Done():
	// 		}
	// 	}(duration)
	// 	w.WriteHeader(200)
	// })

	http.ListenAndServe(":3000", nil)
}
