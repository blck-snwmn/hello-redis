package main

import (
	"log"
	"net/http"
	"sync"
)

func main() {
	var sg sync.WaitGroup
	for i := 0; i < 2; i++ {
		sg.Add(1)
		go func() {
			defer sg.Done()
			resp, err := http.Get("http://localhost:8080")
			if err != nil {
				log.Println(err)
			}
			defer resp.Body.Close()
		}()
	}
	sg.Wait()
}
