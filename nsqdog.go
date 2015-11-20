package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/zorkian/go-datadog-api"
)

var hostname = flag.String("N", "", "Hostname to report to Datadog (defaults to hostname)")

func main() {
	flag.Parse()
	if *hostname == "" {
		h, _ := os.Hostname()
		hostname = &h
	}
	args := flag.Args()
	if len(args) < 3 {
		fmt.Printf("usage: %s [ -N hostname (optional, defaults to hostname) ] <api-key> <app-key> [optional, tags for:datadog]\n", os.Args[0])
		fmt.Printf("example:\n%s my-api-key my-app-key footag env:production\n", os.Args[0])
		os.Exit(2)
	}
	client := datadog.NewClient(args[1], args[2])
	extraTags := args[3:]

	resp, err := http.Get("http://localhost:4151/stats?format=json")
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatalf("response was %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	data := new(Response)
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}
	metrics := make([]datadog.Metric, 0)
	m := func(met datadog.Metric) {
		metrics = append(metrics, met)
	}
	for _, topic := range data.Data.Topics {
		tags := []string{"topic:" + topic.TopicName}
		tags = append(tags, extraTags...)
		m(metricInt64("topic.depth", topic.Depth, tags))
		m(metricInt64("topic.backend_depth", topic.BackendDepth, tags))
		m(metricUint64("topic.messages", topic.MessageCount, tags))
		for _, p := range topic.E2eProcessingLatency.Percentiles {
			m(metric(fmt.Sprintf("topic.e2e_latency_%f", p["quantile"]), p["value"], tags))
		}
		for _, channel := range topic.Channels {
			ctags := append(tags, "channel:"+channel.ChannelName)
			m(metricInt64("chan.depth", channel.Depth, ctags))
			m(metricInt64("chan.backend_depth", channel.BackendDepth, ctags))
			m(metricInt("chan.in_flight", channel.InFlightCount, ctags))
			m(metricInt("chan.deferred", channel.DeferredCount, ctags))
			m(metricUint64("chan.messages", channel.MessageCount, ctags))
			m(metricUint64("chan.requeued", channel.RequeueCount, ctags))
			m(metricUint64("chan.timedout", channel.TimeoutCount, ctags))
			m(metricInt("chan.clients", len(channel.Clients), ctags))
			for _, p := range channel.E2eProcessingLatency.Percentiles {
				m(metric(fmt.Sprintf("chan.e2e_latency_%f", p["quantile"]), p["value"], ctags))
			}
			for _, client := range channel.Clients {
				/*
					there's a bunch of useful values here. may want to revise this later.
					Name: (string) (len=14) "mariadb-1-prod",
					ClientID: (string) (len=14) "mariadb-1-prod",
					Hostname: (string) (len=14) "mariadb-1-prod",
					Version: (string) (len=2) "V2",
					RemoteAddress: (string) (len=20) "10.240.156.158:51549",
					UserAgent: (string) (len=28) "nsq_metrics_to_elasticsearch",
				*/
				t := append(ctags, "client:"+client.ClientID+"-"+client.UserAgent, "addr:"+client.RemoteAddress)
				m(metricInt32("client.state", client.State, t))
				m(metricInt64("client.ready", client.ReadyCount, t))
				m(metricInt64("client.in_flight", client.InFlightCount, t))
				m(metricUint64("client.messages", client.MessageCount, t))
				m(metricUint64("client.finished", client.FinishCount, t))
				m(metricUint64("client.requeued", client.RequeueCount, t))
			}
		}
	}
	err = client.PostMetrics(metrics)
	if err != nil {
		log.Fatalf("fatal: %s\n", err)
	}
}
