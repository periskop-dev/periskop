package scraper

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type Request struct {
	Target        string
	ResultChannel chan<- []errorAggregate
	WaitGroup     *sync.WaitGroup
}

type Processor struct {
	numWorkers      int
	requestsChannel chan Request
	fetcher         ErrorsFetcher
}

func (p Processor) Run() {
	for i := 1; i <= p.numWorkers; i++ {
		go worker(p)
	}
}

func worker(p Processor) {
	for {
		select {
		case r := <-p.requestsChannel:
			if errorAggregates, err := p.fetcher(r.Target); err == nil {
				r.ResultChannel <- errorAggregates
			} else {
				r.ResultChannel <- make([]errorAggregate, 0)
			}
			r.WaitGroup.Done()
		}
	}
}

func (p Processor) Enqueue(r Request) {
	p.requestsChannel <- r
}

func NewProcessor(numWorkers int) Processor {
	return Processor{
		numWorkers:      numWorkers,
		requestsChannel: make(chan Request),
		fetcher:         defaultErrorsFetcher(),
	}
}

type ErrorsFetcher func(string) ([]errorAggregate, error)

func defaultErrorsFetcher() ErrorsFetcher {
	return func(target string) ([]errorAggregate, error) {

		body, err := fetch(target)
		if err != nil {
			return nil, err
		}

		var responsePayload responsePayload
		if err := json.Unmarshal(body, &responsePayload); err != nil {
			return nil, err
		}

		return responsePayload.ErrorAggregate, nil
	}
}

func fetch(target string) ([]byte, error) {
	var netClient = &http.Client{
		Timeout: time.Second * 30,
	}
	resp, err := netClient.Get(target)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
