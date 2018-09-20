package main

import (
	"io/ioutil"
	"net/http"
	"time"
)

const httpGetTimeout = 15000

func getRequest(url string, timeout int) (content []byte, err error) {
	client := &http.Client{Timeout: time.Duration(timeout) * time.Millisecond}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	content, err = ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return nil, err
	}
	return content, nil
}
