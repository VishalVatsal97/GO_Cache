package main

import (
	"fmt"
	"context"
	"net/http"
	"log"
	"time"
	"github.com/ParallelDots/cache"
	//"runtime"
	"sync"
)

var c = cache.NewCache()

func handler(w http.ResponseWriter,r *http.Request) {
	if val,isPresent := c.FindResponse(r.URL.Path[1:]); isPresent {
		fmt.Println("Found in cache")
		fmt.Fprintf(w,val)
	} else {
		fmt.Println("Going to sleep")
		time.Sleep(10*time.Second)
		fmt.Fprintf(w,"Hi there, I love %s", r.URL.Path[1:])
		c.AddToCache(r.URL.Path[1:],"Hi there, I love " + r.URL.Path[1:])
		fmt.Printf("Hi there, I love %s!\n", r.URL.Path[1:])
	}
}

func startHttpServer(wg *sync.WaitGroup) *http.Server {
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
    log.Printf("main: starting HTTP server")

    httpServerExitDone := &sync.WaitGroup{}

	httpServerExitDone.Add(1)
    
    err := c.LoadFromFile("servercache.gob")
    
    if err!=nil {
        fmt.Println("Did not find any file for the cache")
    }
    
    srv := startHttpServer(httpServerExitDone)

    //log.Printf("main: serving for 10 seconds")

    time.Sleep(300 * time.Second)

    log.Printf("main: stopping HTTP server")

    if err := srv.Shutdown(context.Background()); err != nil {
        panic(err)
    }

	httpServerExitDone.Wait()
	
	c.SaveToFile("servercache.gob")

    log.Printf("main: done. exiting")
}