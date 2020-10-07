package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

var numberOfRequests = flag.Int("requests", 10, "number of requests")
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

type logInfo struct {
	RequestString  string
	ResponseString string
}

func makeRequest(url string, ch chan<- logInfo, wg *sync.WaitGroup) {
	timestart := time.Now()
	resp, _ := http.Get(url)
	timeend := time.Now()
	fmt.Println(timeend.Sub(timestart))
	body, _ := ioutil.ReadAll(resp.Body)

	loginfo := logInfo{
		RequestString:  url,
		ResponseString: string(body),
	}
	ch <- loginfo
	wg.Done()

}

func main() {
	flag.Parse()
	clientWg := &sync.WaitGroup{}
	var (
		temp         []logInfo
		randomString string
	)
	var uniqueRequests int
	uniqueRequests = (*numberOfRequests * 20) / 100
	rand.Seed(time.Now().UnixNano())
	ch := make(chan logInfo)
	clientWg.Add(*numberOfRequests)

	for i := 0; i < *numberOfRequests; i++ {
		if i < uniqueRequests {
			randomString = randSeq(10)
		} else {
			randomString = "RRHBECYGwy"
		}

		go makeRequest("http://127.0.0.1:8080/"+randomString, ch, clientWg)
	}

	go func() {
		clientWg.Wait()
		close(ch)
	}()

	for i := 0; i < *numberOfRequests; i++ {
		temp = append(temp, <-ch)
	}

	file, _ := json.MarshalIndent(temp, "", "")

	f, err := os.OpenFile("client.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write(file); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
