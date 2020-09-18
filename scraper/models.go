package scraper

import "time"

type responsePayload struct {
	ErrorAggregate []errorAggregate `json:"aggregated_errors"`
	Instance       string           `json:"instance"`
}

type errorAggregate struct {
	AggregationKey string             `json:"aggregation_key"`
	TotalCount     int                `json:"total_count"`
	Severity       string             `json:"severity"`
	LatestErrors   []errorWithContext `json:"latest_errors"`
}

type errorWithContext struct {
	Error       errorInstance `json:"error"`
	UUID        string        `json:"uuid"`
	Timestamp   time.Time     `json:"timestamp"`
	Severity    string        `json:"severity"`
	HTTPContext *httpContext  `json:"http_context"`
}

type errorInstance struct {
	Class      string         `json:"class"`
	Message    string         `json:"message"`
	Stacktrace []string       `json:"stacktrace"`
	Cause      *errorInstance `json:"cause"`
}

type httpContext struct {
	RequestMethod  string            `json:"request_method"`
	RequestURL     string            `json:"request_url"`
	RequestHeaders map[string]string `json:"request_headers"`
	RequestBody    string            `json:"request_body"`
}

type errorOccurrences []errorWithContext

func (e errorOccurrences) Len() int {
	return len(e)
}
func (e errorOccurrences) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}
func (e errorOccurrences) Less(i, j int) bool {
	return e[i].Timestamp.After(e[j].Timestamp)
}
