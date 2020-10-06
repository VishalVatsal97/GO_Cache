package main

import (
	"io/ioutil"
	"net/http"
	"time"
	"math/rand"
	"encoding/json"
	"sync"
	"os"
	"log"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

type LogInfo struct {
	Request_string string
	Response_string string
}

func MakeRequest(url string, ch chan<-LogInfo, wg *sync.WaitGroup) {
	
	resp, _ := http.Get(url)
	body, _ := ioutil.ReadAll(resp.Body)
	
	loginfo := LogInfo {
		Request_string : url,
		Response_string : string(body),
	}
	ch <- loginfo
	wg.Done()
	
}

func main() {
	
	clientWg := &sync.WaitGroup{}
	var (
		temp []LogInfo
		random_string string
	)

	rand.Seed(time.Now().UnixNano())
	ch := make(chan LogInfo)
	clientWg.Add(10)
	
	for i := 0 ; i < 10 ; i++ {
		if i < 8 {
			random_string = "RRHBECYGwy"
		} else {
			random_string = randSeq(10)
		}
		
		go MakeRequest("http://127.0.0.1:8080/" + random_string, ch, clientWg)
	}
	
	go func() {
        clientWg.Wait()
        close(ch)
    }()
	
	for i := 0 ; i < 10 ;i++ {
		temp  = append(temp,<-ch)
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