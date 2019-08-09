package scraper

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestCombineErrorsSortsByTimestamp(t *testing.T) {
	firstContent, _ := ioutil.ReadFile("sample-response1.json")
	secondContent, _ := ioutil.ReadFile("sample-response2.json")

	var firstResponsePayload responsePayload
	var secondResponsePayload responsePayload

	json.Unmarshal(firstContent, &firstResponsePayload)  // nolint[errcheck]
	json.Unmarshal(secondContent, &secondResponsePayload) // nolint[errcheck]

	firstOccurrences := firstResponsePayload.ErrorAggregate[0].LatestErrors
	secondOccurrences := secondResponsePayload.ErrorAggregate[0].LatestErrors

	result := combine(firstOccurrences, secondOccurrences)

	expectedUUIDs := []string{"uuid4", "uuid2", "uuid3", "uuid1"}

	for i, element := range result {
		if element.UUID != expectedUUIDs[i] {
			t.Errorf("Expected %s, Found %s", expectedUUIDs[i], element.UUID)
		}
	}
}
