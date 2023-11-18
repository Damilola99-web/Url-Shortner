package routes

import (
	"testing"
	"time"
)

func TestRequest(t *testing.T) {
	// Create an instance of the request struct with sample values
	req := request{
		URL:         "https://example.com",
		CustomShort: "custom",
		Expiry:      time.Hour,
	}

	// Assert that the values of the struct fields are set correctly
	if req.URL != "https://example.com" {
		t.Errorf("Expected URL to be 'https://example.com', got '%s'", req.URL)
	}

	if req.CustomShort != "custom" {
		t.Errorf("Expected CustomShort to be 'custom', got '%s'", req.CustomShort)
	}

	if req.Expiry != time.Hour {
		t.Errorf("Expected Expiry to be '%s', got '%s'", time.Hour, req.Expiry)
	}
}

// Run the test function
func TestMain(m *testing.M) {
	// Run the tests
	m.Run()
}