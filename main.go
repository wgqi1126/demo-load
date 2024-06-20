package main

import (
	"context"
	"encoding/json"
	"math"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	http.HandleFunc("/", index)
	poe(http.ListenAndServe(":8800", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	tm := time.Now()

	c, _ := strconv.Atoi(r.URL.Query().Get("c"))
	if c == 0 {
		c = runtime.NumCPU()
	}
	n, _ := strconv.ParseInt(r.URL.Query().Get("n"), 10, 64)
	if n == 0 {
		n = 1000 * 10000
	}
	t, _ := strconv.ParseFloat(r.URL.Query().Get("t"), 64)
	ctx := r.Context()
	if t > 0 {
		ctx, _ = context.WithTimeout(ctx, time.Duration(t*float64(time.Second)))
		n = math.MaxInt64
	}

	cur := int64(0)

	wg := sync.WaitGroup{}
	for i := 0; i < c; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for atomic.LoadInt64(&cur) < n {
				select {
				case <-ctx.Done():
					goto end
				default:
					_ = 10000000 / 322
					atomic.AddInt64(&cur, 1)
				}
			}
		end:
			return
		}()
	}
	wg.Wait()

	spent := time.Since(tm)
	rs := map[string]any{
		"c":     c,
		"n":     n,
		"t":     t,
		"cur":   cur,
		"spent": spent.String(),
		"avg":   (spent / time.Duration(n)).String(),
		"avg_n": int64(float64(cur) / float64(spent) * float64(time.Second)),
	}
	rsB, err := json.Marshal(rs)
	poe(err)
	_, err = w.Write(rsB)
	poe(err)
}

func poe(err error) {
	if err != nil {
		panic(err)
	}
}
