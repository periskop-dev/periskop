package scraper

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/soundcloud/periskop/repository"
)

func TestCombineLastErrorsSortsByTimestamp(t *testing.T) {
	firstContent, _ := ioutil.ReadFile("sample-response1.json")
	secondContent, _ := ioutil.ReadFile("sample-response2.json")

	var firstResponsePayload responsePayload
	var secondResponsePayload responsePayload

	json.Unmarshal(firstContent, &firstResponsePayload)   // nolint[errcheck]
	json.Unmarshal(secondContent, &secondResponsePayload) // nolint[errcheck]

	firstOccurrences := firstResponsePayload.ErrorAggregate[0].LatestErrors
	secondOccurrences := secondResponsePayload.ErrorAggregate[0].LatestErrors

	result := combineLastErrors(firstOccurrences, secondOccurrences)

	expectedUUIDs := []string{"uuid4", "uuid2", "uuid3", "uuid1"}

	for i, element := range result {
		if element.UUID != expectedUUIDs[i] {
			t.Errorf("Expected %s, Found %s", expectedUUIDs[i], element.UUID)
		}
	}
}

func TestScrapeCombine(t *testing.T) {
	var targetErrorsCount = make(targetErrorsCountMap)
	var errorAggregates = make(errorAggregateMap)
	errorInstancesAccumulator := make(errorInstancesAccumulatorMap)
	repo := repository.NewMemoryRepository()

	firstContent, _ := ioutil.ReadFile("sample-response1.json")
	var rp responsePayload
	json.Unmarshal(firstContent, &rp) // nolint[errcheck]
	rp.Target = "test"

	errorAggregates.combine("test", &repo, rp, targetErrorsCount, errorInstancesAccumulator)

	count := targetErrorsCount["test"]["com.soundcloud.Foon@e28e036e"]
	if count != 2 {
		t.Errorf("Expected 2 element, Found %d", count)
	}

	countErrorInstances := len(errorInstancesAccumulator["com.soundcloud.Foon@e28e036e"])
	if countErrorInstances != 2 {
		t.Errorf("Expected 2 element, Found %d", countErrorInstances)
	}

	rp.ErrorAggregate[0].TotalCount = 4
	errorAggregates.combine("test", &repo, rp, targetErrorsCount, errorInstancesAccumulator)

	count = targetErrorsCount["test"]["com.soundcloud.Foon@e28e036e"]
	if count != 4 {
		t.Errorf("Expected 4 element, Found %d", count)
	}

	countErrorInstances = len(errorInstancesAccumulator["com.soundcloud.Foon@e28e036e"])
	if countErrorInstances != 4 {
		t.Errorf("Expected 2 element, Found %d", countErrorInstances)
	}
}

func TestScapeCombineNotUpdate(t *testing.T) {
	var targetErrorsCount = make(targetErrorsCountMap)
	var errorAggregates = make(errorAggregateMap)
	errorInstancesAccumulator := make(errorInstancesAccumulatorMap)
	repo := repository.NewMemoryRepository()

	firstContent, _ := ioutil.ReadFile("sample-response1.json")
	var rp responsePayload
	json.Unmarshal(firstContent, &rp) // nolint[errcheck]
	rp.Target = "test"

	errorAggregates.combine("test", &repo, rp, targetErrorsCount, errorInstancesAccumulator)

	rp.ErrorAggregate[0].TotalCount = 1
	errorAggregates.combine("test", &repo, rp, targetErrorsCount, errorInstancesAccumulator)

	count := targetErrorsCount["test"]["com.soundcloud.Foon@e28e036e"]
	if count != 2 {
		t.Errorf("Expected 2 element, Found %d", count)
	}

	countErrorInstances := len(errorInstancesAccumulator["com.soundcloud.Foon@e28e036e"])
	if countErrorInstances != 2 {
		t.Errorf("Expected 2 element, Found %d", countErrorInstances)
	}
}

func TestScapeCombineCreatedAt(t *testing.T) {
	var targetErrorsCount = make(targetErrorsCountMap)
	var errorAggregates = make(errorAggregateMap)
	errorInstancesAccumulator := make(errorInstancesAccumulatorMap)
	repo := repository.NewMemoryRepository()

	firstContent, _ := ioutil.ReadFile("sample-response1.json")
	var rp responsePayload
	json.Unmarshal(firstContent, &rp) // nolint[errcheck]
	rp.Target = "test1"
	errorAggregates.combine("test1", &repo, rp, targetErrorsCount, errorInstancesAccumulator)

	secondContent, _ := ioutil.ReadFile("sample-response2.json")
	json.Unmarshal(secondContent, &rp) // nolint[errcheck]
	rp.Target = "test2"
	errorAggregates.combine("test2", &repo, rp, targetErrorsCount, errorInstancesAccumulator)

	createdAtHour := errorAggregates["com.soundcloud.Foon@e28e036e"].CreatedAt.Hour()

	if createdAtHour != 15 {
		t.Errorf("Expected 15h, Found %d", createdAtHour)
	}
}
