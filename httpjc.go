package main

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// siteStats describe the statistical info of this site.
type siteStats struct {
	Req int `json:"total"`
	Avg int `json:"average"`
}

// hashInfo stores the result of each request
type hashInfo struct {
	tdur time.Duration
	hstr string
}

// hashWord takes an input string, using SHA512 encoding to generate a base64 encoded string
func hashWord(pws string) string {
	bts := sha512.Sum512([]byte(pws))
	return base64.StdEncoding.EncodeToString([]byte(bts[:]))
}

func main() {
	var mu sync.Mutex
	var total int

	// A map to store hashed information
	hashmap := make(map[int]hashInfo, 10)

	m := http.NewServeMux()
	s := http.Server{Addr: ":8080", Handler: m}

	// Handle / page
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s?\n", r.URL.Path)
	})

	// Handle /hash page
	m.HandleFunc("/hash", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Fprintf(w, "%s\n", r.URL.Path)

		pws := r.PostFormValue("password")

		// Didn't find the "password" parameter.
		if pws == "" {
			fmt.Fprintf(w, "Only password parameter is supported:\n")

			// Show the user what were requested.
			if err := r.ParseForm(); err != nil {
				log.Print(err)
			}
			for k, v := range r.Form {
				fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
			}
		} else {
			// fmt.Fprintf(w, "password = %q\n", pws)

			// Can have concurrent requests
			mu.Lock()
			idx := total
			total++
			mu.Unlock()

			go func(k int, v string) {
				var req hashInfo

				// Artifical delay 5 seconds
				time.Sleep(5 * time.Second)

				// Time the hash function
				t1 := time.Now()
				req.hstr = hashWord(v)
				t2 := time.Now()
				req.tdur = t2.Sub(t1)

				// Create the entry
				hashmap[k] = req
			}(idx, pws)

			// The identifier is a 1-based number
			fmt.Fprintf(w, "%d\n", idx+1)
			// fmt.Fprintf(w, "%v\n", req.tdur)
			// fmt.Fprintf(w, "%q\n", req.hstr)
		}
	})

	// Handle /hash/# pages
	m.HandleFunc("/hash/", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Fprintf(w, "%s\n", r.URL.Path)

		// Get the request count substring
		req := r.URL.Path[len("/hash/"):]

		// Get the identifier
		if idx, err := strconv.Atoi(req); err == nil {
			// fmt.Fprintf(w, "Request number: %d\n", idx)

			// Identifer number is 1-base
			idx--
			// Ready or not
			if hinf, ok := hashmap[idx]; ok {
				fmt.Fprintf(w, "%q\n", hinf.hstr)
			} else {
				fmt.Fprintf(w, "Identifier number %d is not ready.\n", idx+1)
			}
		} else {
			fmt.Fprintf(w, "%v\n", err)
		}
	})

	// Handle /stats page
	m.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Fprintf(w, "/stats\n")

		var stat siteStats
		var tavg time.Duration

		// Wait for all the pending requests done
		for total != len(hashmap) {
		}
		// Total durations of all requests
		for _, req := range hashmap {
			tavg += req.tdur
		}
		// Total number of requests
		stat.Req = total
		// Average time in microseconds
		stat.Avg = int(tavg.Nanoseconds()/1000) / total

		// Return in JSON
		if bytes, err := json.Marshal(stat); err == nil {
			fmt.Fprintf(w, "%v\n", string(bytes))
		} else {
			fmt.Fprintf(w, "%v\n", err)
		}
	})

	// Handle /shutdown page
	m.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Shutdown\n")
		// Use Goroutine to shutdown this server, so the previous string can be returned
		go func() {
			if err := s.Shutdown(context.Background()); err != nil {
				log.Fatal(err)
			}
		}()
	})
	// Launch the server and check for the shutdown
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
	// Wait for all the pending requests
	for total != len(hashmap) {
	}
	log.Printf("Finished")
}
