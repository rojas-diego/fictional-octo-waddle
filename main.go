package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

var commonHttpResponseCode = []int{
	http.StatusOK,
	http.StatusCreated,
	http.StatusBadRequest,
	http.StatusUnauthorized,
	http.StatusAccepted,
	http.StatusNotFound,
	http.StatusInternalServerError,
}

var (
// runtimeInfoLock   sync.Mutex = sync.Mutex{}
// runtimeInfoLastTS time.Time  = time.Time{}
// runtimeInfoCpuStats *cpu.Stats = nil
)

func main() {
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		value, err := strconv.Atoi(r.URL.Query().Get("value"))
		if err != nil {
			w.WriteHeader(commonHttpResponseCode[rand.Intn(len(commonHttpResponseCode))])
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

	http.HandleFunc("/fibonacci", func(w http.ResponseWriter, r *http.Request) {
		value, err := strconv.Atoi(r.URL.Query().Get("value"))
		if err != nil {
			value = 50000000
		}
		for j := 0; j < value; j++ {
			a := 0
			b := 1
			c := b
			for {
				c = b
				b = a + b
				if b >= value {
					break
				}
				a = c
			}
		}
		w.WriteHeader(200)
	})

	http.HandleFunc("/runtime_info", func(w http.ResponseWriter, r *http.Request) {

		infostat, err := cpu.Info()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for i := range infostat {
			println(infostat[i].Mhz)
		}

		// now := time.Now()
		cpuinfo, err := cpu.Percent(time.Second, false)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		percpuinfo, err := cpu.Percent(time.Second, true)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// runtimeInfoLock.Lock()
		// defer runtimeInfoLock.Unlock()
		// if runtimeInfoCpuStats == nil {
		// 	w.WriteHeader(http.StatusOK)
		// 	w.Write([]byte(`"Called for the first time"`))
		// 	runtimeInfoCpuStats = cpuinfo
		// 	runtimeInfoLastTS = now
		// 	return
		// }

		// deltaNs := now.Sub(runtimeInfoLastTS).Nanoseconds()
		// cpuDeltaNs := cpuinfo.Total - runtimeInfoCpuStats.Total
		// totalPercent := (float64(cpuDeltaNs) / float64(deltaNs)) * float64(100)

		// runtimeInfoLastTS = now
		// runtimeInfoCpuStats = cpuinfo

		type cpuInfoResponse struct {
			TotalPercent []float64 `json:"total_percent"`
			PerCpu       []float64 `json:"per_cpu"`
		}

		body, err := json.MarshalIndent(cpuInfoResponse{
			TotalPercent: cpuinfo,
			PerCpu:       percpuinfo,
		}, "", "    ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(body)

		// w.WriteHeader(http.StatusOK)
		// w.Write([]byte(fmt.Sprintf(`{"
		// 	"total_percent": %f,
		// 	"total": %d,
		// 	"delta_ns": %d,
		// 	"cpu_delta_ns": %d,
		// "}`, totalPercent, cpuinfo.Total, deltaNs, cpuDeltaNs)))
	})

	http.ListenAndServe(":3000", nil)
}
