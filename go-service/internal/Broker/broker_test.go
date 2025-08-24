package broker_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/SLANGERES/go-service/internal/Broker"
	"github.com/SLANGERES/go-service/internal/Model"
)

// Fake log entry
var logEntry = model.Logging{
	Id:       "test_id",
	Level:    "info",
	Message:  "test message",
	Service:  "test-service",
	MetaData: Metadata,
}

var Metadata=model.MetaData{
	SourceIp: "0.0.0.0",
	Region: "test-region-01",
}
// Mock connection function
func TestSendDataToQueue_JSONMarshal(t *testing.T) {
	// Just make sure marshaling doesn’t fail
	_, err := json.Marshal(logEntry)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

//Sending data to queue 
func TestSendDataToQueue_Integration(t *testing.T) {

	err := broker.SendDataToQueue(logEntry)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	time.Sleep(1 * time.Second)
}
