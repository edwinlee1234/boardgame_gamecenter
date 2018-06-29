package main

var success = "success"
var wrong = "error"

// Response 回應資訊格式
type Response struct {
	Status string                   `json:"status"`
	Data   map[string][]interface{} `json:"data"`
}
