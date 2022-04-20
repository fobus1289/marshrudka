package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	var wg = &sync.WaitGroup{}
	var counter int32 = 0
	wg.Add(100000)
	var old = time.Now()
	log.Println("Start")

	for j := 0; j < 100; j++ {
		for i := 0; i < 1000; i++ {
			go func() {

				response, err := http.Get("http://localhost:8080/")
				wg.Done()
				if err != nil {
					atomic.AddInt32(&counter, 1)
					return
				}
				var in interface{}

				if err := json.NewDecoder(response.Body).Decode(&in); err != nil {
					log.Println(err)
				}
			}()
		}
	}

	wg.Wait()
	log.Println(counter)
	log.Println(time.Now().Sub(old))
}
