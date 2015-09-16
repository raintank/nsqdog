package main

import "github.com/nsqio/nsq/nsqd"

type Response struct {
	StatusCode int    `json:"status_code"`
	StatusTxt  string `json:"status_txt"`
	Data       Header `json:"data"`
}

type Header struct {
	Version   string            `json:"version"`
	Health    string            `json:"health"`
	StartTime int64             `json:"start_time"`
	Topics    []nsqd.TopicStats `json:"topics"`
}
