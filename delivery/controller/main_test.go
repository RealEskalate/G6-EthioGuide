package controller_test

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

// toJSON is a common helper to convert any struct to a JSON bytes buffer for requests.
func toJSON(v interface{}) *bytes.Buffer {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(v)
	if err != nil {
		panic(err)
	}
	return &buf
}

// TestMain sets the Gin mode to test for all tests in this package.
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
