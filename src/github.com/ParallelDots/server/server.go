package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/ParallelDots/cache"
)

var enableCaching = flag.Bool("cache", false, "enable/disable caching default disabled")
var c = cache.NewCache()

func handler(w http.ResponseWriter, r *http.Request) {
	response := "Hi there, I love " + r.URL.Path[1:]

	if *enableCaching {
		if val, isPresent := c.FindResponse(r.URL.Path[1:]); isPresent {
			fmt.Println("Found in cache")
			fmt.Fprintf(w, val)
			return
		}
		time.Sleep(30 * time.Second)
		c.AddToCache(r.URL.Path[1:], response)
	} else {
		time.Sleep(30 * time.Second)
	}
	fmt.Fprintf(w, response)
	return
}

func startHTTPServer(wg *sync.WaitGroup) *http.Server {
	srv := &http.Server{Addr: ":8080"}

	http.HandleFunc("/", handler)

	go func() {
		defer wg.Done() // let main know we are done cleaning up

		// always returns error. ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// returning reference so caller can call Shutdown()
	return srv
}

func main() {
	flag.Parse()
	log.Printf("main: starting HTTP server")
	fmt.Println(*enableCaching)
	httpServerExitDone := &sync.WaitGroup{}

	httpServerExitDone.Add(1)
	runtime.GOMAXPROCS(runtime.NumCPU())

	if *enableCaching {
		err := c.LoadFromFile("servercache.gob")
		if err != nil {
			fmt.Println("Did not find any file for the cache")
		}
	}

	srv := startHTTPServer(httpServerExitDone)

	time.Sleep(3600 * time.Second)

	log.Printf("main: stopping HTTP server")

	if err := srv.Shutdown(context.Background()); err != nil {
		panic(err)
	}

	httpServerExitDone.Wait()

	if *enableCaching {
		c.SaveToFile("servercache.gob")
	}

	log.Printf("main: done. exiting")
}
