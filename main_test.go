// code/middleware/main_test.go
package main // Must match the package of the code being tested

import (
	"encoding/json" // Needed for json validation
	"net/http"      // Needed for HTTP status codes
	"net/http/httptest" // Provides utilities for testing HTTP handlers
	"strings"       // Might be useful for string assertions, though not strictly needed here
	"testing"       // The standard Go testing package
)

// TestHealthHandler tests the healthHandler function.
func TestHealthHandler(t *testing.T) {
	// Create a mock HTTP request. We don't need a specific URL or body for this handler.
	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		// If request creation fails, report a fatal error and stop the test.
		t.Fatalf("Could not create request: %v", err)
	}

	// httptest.NewRecorder acts as a mock http.ResponseWriter.
	// It records the status code, headers, and body written by the handler.
	rr := httptest.NewRecorder()

	// Call the handler function directly, passing the mock recorder and request.
	healthHandler(rr, req)

	// --- Assertions ---

	// Check the status code.
	expectedStatus := http.StatusOK // Expecting a 200 OK
	if status := rr.Code; status != expectedStatus {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, expectedStatus)
	}

	// Check the response body.
	expectedBody := "ok" // Expecting the body to be "ok"
	if rr.Body.String() != expectedBody {
		t.Errorf("Handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}

// TestMessage_JSON tests the JSON method of the Message struct.
func TestMessage_JSON(t *testing.T) {
	// Define a test case for the Message struct.
	msg := Message{
		Text:    "Test message",
		Details: "Some details here",
	}

	// Expected JSON output. Note the backticks for a raw string literal.
	expectedJSON := `{"text":"Test message","details":"Some details here"}`

	// Call the method being tested.
	actualJSON, err := msg.JSON()

	// --- Assertions ---

	// Check if there was an error during JSON marshaling.
	if err != nil {
		t.Fatalf("Message.JSON returned an unexpected error: %v", err)
	}

	// For comparing JSON strings, it's often safer to unmarshal and compare the structures,
	// as whitespace or key order might differ but the JSON is semantically the same.
	// However, for simple cases like this, direct string comparison is usually sufficient
	// if you control the expected format precisely. Let's do a simple string compare first.
	// A more robust way is shown below.

	if actualJSON != expectedJSON {
		t.Errorf("Message.JSON returned incorrect JSON:\nGot: %s\nWant: %s",
			actualJSON, expectedJSON)
	}

	// --- More Robust JSON Comparison (Optional but Recommended) ---
	// This method is better as it handles potential differences in formatting (like spacing).

	var actualMap, expectedMap map[string]interface{}

	// Unmarshal the actual and expected JSON strings into maps.
	err = json.Unmarshal([]byte(actualJSON), &actualMap)
	if err != nil {
		t.Fatalf("Failed to unmarshal actual JSON '%s': %v", actualJSON, err)
	}
	err = json.Unmarshal([]byte(expectedJSON), &expectedMap)
	if err != nil {
		t.Fatalf("Failed to unmarshal expected JSON '%s': %v", expectedJSON, err)
	}

	// Compare the unmarshaled maps. This is a more robust way to check if the JSON content is the same.
	// For simple maps like this, direct comparison works. For complex nested structures,
	// you might need a deeper comparison function or a library like testify/assert's Equal.
	if !mapsEqual(actualMap, expectedMap) {
		t.Errorf("Message.JSON returned semantically incorrect JSON:\nGot: %s\nWant: %s",
			actualJSON, expectedJSON)
	}
}

// Helper function for comparing simple maps (for the robust JSON test)
func mapsEqual(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || bv != v {
			return false
		}
	}
	return true
}
