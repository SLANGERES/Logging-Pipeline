package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	model "github.com/SLANGERES/go-service/internal/Model"
	"github.com/SLANGERES/go-service/internal/Routers"
)
func TestHealthz(t *testing.T){
	router:=router.Router()
	w:=httptest.NewRecorder()

	req,_:=http.NewRequest("GET","/health",nil)
	router.ServeHTTP(w,req)
	if w.Code!=http.StatusOK{
		t.Fatalf("expected %d and got %d",http.StatusOK,w.Code)
	}
}


//! currently failing because  rabbit mq server is not running locally 
func TestPostLogging(t *testing.T) {
    // Initialize router
    r := router.Router()

    // Prepare test pyload
    var postRequest model.Logging
    postRequest.Id = "test_id"
    postRequest.Level = "test"
    postRequest.Message = "testing the handler in go"
    postRequest.Service = "testing"
    postRequest.TimeStamp = time.Now()
    postRequest.MetaData.Region = "testing-region"
    postRequest.MetaData.SourceIp = "0.0.0.0"

    // Marshal payload into JSON
    body, err := json.Marshal(postRequest)
    if err != nil {
        t.Fatalf("failed to marshal request: %v", err)
    }

    // Create request with JSON body
    req, _ := http.NewRequest("POST", "/log", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")

    // Record response
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    // Assert
    if w.Code != http.StatusOK {
        t.Fatalf("expected %d, got %d, response: %s", http.StatusOK, w.Code, w.Body.String())
    }
}
