package main

import (
	"fmt"
	"net/http"
	"net/http/httptrace"
	"sync"
	"time"

	"github.com/influxdata/tdigest"
)

// number of distinct requests
var n = 1000

var startTime = time.Now()

func now() time.Duration {
	return time.Since(startTime)
}

var results = make(chan time.Duration, n)

func main() {
	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			var dnsStart, dnsDuration time.Duration
			trace := &httptrace.ClientTrace{
				DNSStart: func(_ httptrace.DNSStartInfo) {
					dnsStart = now()
				},
				DNSDone: func(_ httptrace.DNSDoneInfo) {
					dnsDuration = now() - dnsStart
				},
			}
			req, _ := http.NewRequest("GET", "https://example.com", nil)
			req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
			_, err := http.DefaultTransport.RoundTrip(req)
			results <- dnsDuration
			wg.Done()
			if err != nil {
				fmt.Println(err)
				return
			}
		}()
	}
	wg.Wait()
	td := tdigest.NewWithCompression(float64(n))
	var avgDNS float64
	var maxDNS float64
	for i := 0; i < n; i++ {
		dns := <-results
		avgDNS += dns.Seconds()
		td.Add(dns.Seconds(), 1)
		if maxDNS < dns.Seconds() {
			maxDNS = dns.Seconds()
		}
	}
	fmt.Println("avg: ", avgDNS/float64(n))
	fmt.Println("max: ", maxDNS)
	fmt.Println("50th percentile: ", td.Quantile(0.5))
	fmt.Println("75th percentile: ", td.Quantile(0.75))
	fmt.Println("90th percentile: ", td.Quantile(0.9))
	fmt.Println("99th percentile: ", td.Quantile(0.99))
}
