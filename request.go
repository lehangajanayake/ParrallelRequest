package main

import(
	"net/http"
	"sync"
	"time"
	"bytes"
	"log"
	"encoding/json"
)

func client()*http.Client {
	return &http.Client{}
}

func newRequest(method string, url string, data jsonBody)(*http.Request, error){
	if data != nil {
		body, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		request, err := http.NewRequest(method, url, bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}
		return request, nil
	}
	request, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(``)))
	if err != nil {
		return nil, err
	}

	return request, nil
}

func request(client *http.Client, request *http.Request, responsechan chan Response, wg *sync.WaitGroup){
	defer wg.Done()
	
	var response Response

	start := time.Now()
	response.Response, response.Error = client.Do(request)
	response.Time = time.Since(start)
	responsechan <- response
}

func parallelRequest(method, URL string, client *http.Client, num int, data []jsonBody)[]Response{
	var wg sync.WaitGroup
	var responses []Response
	responsechan := make(chan Response)
	done := make(chan bool)
	for i := 0; i < num; i++{
		wg.Add(1)
		if len(data) == 0 {
			req, err := newRequest(method, URL, nil)
			if err != nil {
				log.Println("Error creating the request, ", err.Error())
				return nil
			}
			go request(client, req, responsechan, &wg)
			continue
		}else if len(data) == 1 {
			req, err := newRequest(method, URL, data[0])
			if err != nil {
				log.Println("Error creating the request, ", err.Error())
				return nil
			}
			go request(client, req, responsechan, &wg)
			continue
		}
		req, err := newRequest(method, URL, data[i])
		if err != nil {
			log.Println("Error creating the request, ", err.Error())
			return nil
		}
		go request(client, req, responsechan, &wg)
	}
	go func(){
		for v := range responsechan{
			responses = append(responses, v)
		}
		done <- true
	}()
	wg.Wait()
	close(responsechan)
	<- done
	return responses
}