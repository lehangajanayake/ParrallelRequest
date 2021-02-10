package main

import (
	"net/http"
	"time"
)

type flags struct{
	FilePath string
	URL string
	Method string
	Requests int
	Limit int
	Log bool
}

type jsonBody map[string]interface{}

//Response : a custom response with response time and error
type Response struct {
	Response *http.Response
	Time time.Duration
	Error error
}