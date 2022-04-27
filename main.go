package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"
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

	http.HandleFunc("/allocate", func(w http.ResponseWriter, r *http.Request) {
		mb, err := strconv.Atoi(r.URL.Query().Get("mb"))
		if err != nil {
			mb = 32
		}
		duration, err := time.ParseDuration(r.URL.Query().Get("for"))
		if err != nil {
			duration = time.Second * 30
		}
		go func(mb int, duration time.Duration) {
			mem := make([][]byte, 0)
			for i := 0; i < mb; i++ {
				mem = append(mem, make([]byte, 1024*1024))
				for j := 0; j < 1024*1024; j++ {
					mem[i][j] = byte(1)
				}
			}
			time.Sleep(duration)
			for i := 0; i < mb; i++ {
				for j := 0; j < 1024*1024; j++ {
					mem[i][j] = byte(0)
				}
			}
		}(mb, duration)
	})

	http.HandleFunc("/delay", func(w http.ResponseWriter, r *http.Request) {
		wait, err := time.ParseDuration(r.URL.Query().Get("value"))
		if err != nil {
			wait = time.Millisecond * time.Duration(rand.Intn(500))
		}
		time.Sleep(wait)
	})

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

	// go func() {
	// 	for {
	// 		PrintMemUsage()
	// 		time.Sleep(time.Second)
	// 	}
	// }()
	http.ListenAndServe(":3000", nil)
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
