package scraper

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/soundcloud/periskop-go"
	"github.com/soundcloud/periskop/metrics"
)

const httpClientTimeoutSeconds = 30

type Request struct {
	Target        string
	ResultChannel chan<- responsePayload
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
				r.ResultChannel <- responsePayload{}
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

type ErrorsFetcher func(string) (responsePayload, error)

func defaultErrorsFetcher() ErrorsFetcher {
	return func(target string) (responsePayload, error) {
		body, err := fetch(target)
		if err != nil {
			metrics.ErrorCollector.ReportWithHTTPContext(err, &periskop.HTTPContext{
				RequestMethod: "GET",
				RequestURL:    target,
			}, "scrapped-url-error")
			return responsePayload{}, err
		}

		var rp responsePayload
		if err := json.Unmarshal(body, &rp); err != nil {
			metrics.ErrorCollector.ReportWithHTTPContext(err, &periskop.HTTPContext{
				RequestMethod: "GET",
				RequestURL:    target,
			})
			return responsePayload{}, err
		}
		rp.Instance = target
		return rp, nil
	}
}

func fetch(target string) ([]byte, error) {
	var netClient = &http.Client{
		Timeout: time.Second * httpClientTimeoutSeconds,
	}
	resp, err := netClient.Get(target)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
