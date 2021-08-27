package repository

import (
	"testing"
)

func TestSetAndRetrieveTargets(t *testing.T) {
	er := &targetsRepository{}
	er.StoreTargets(serviceName, []Target{
		{Endpoint: "localhost:3000/-/exceptions"},
		{Endpoint: "localhost:3001/-/exceptions"},
	})

	retrievedTargets := er.GetTargets()
	if retrievedTargets[serviceName][0].Endpoint != "localhost:3000/-/exceptions" ||
		retrievedTargets[serviceName][1].Endpoint != "localhost:3001/-/exceptions" {
		t.Errorf("Inconsistent target fetch and retrieval")
	}
}
