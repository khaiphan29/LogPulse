package log_handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

   "github.com/gin-gonic/gin"

   "github.com/khaiphan29/logpulse/internal/api/router"
	"github.com/khaiphan29/logpulse/internal/api/parsing"
	"github.com/stretchr/testify/assert"
)

var testRouter *gin.Engine

func init() {
   // Set Gin to test mode
   gin.SetMode(gin.TestMode)

   // Initialize the router
   testRouter = gin.New()
   router.SetupRouter(testRouter)
}

func TestGetLogHandler(t *testing.T) {
   req, _ := http.NewRequest("GET", "/logs", nil) // Create a new HTTP GET request
	w := httptest.NewRecorder()    // Where the HTTP response will be stored

	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message": "This is LOG"}`, w.Body.String())
}


func TestPostLogHandler(t *testing.T) {
	// Iterate through test cases
	for _, testCase := range postTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Serialize the payload to JSON (if it's not already a string)
			var requestBody []byte
			var err error

         // Process the payload based on its type
			if str, ok := testCase.payload.(string); ok {
				requestBody = []byte(str)
			} else {
				requestBody, err = json.Marshal(testCase.payload)
				assert.NoError(t, err, "Failed to marshal payload")
			}

			// Create a new HTTP POST request
			req, _ := http.NewRequest("POST", "/logs", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			// Record the response
			w := httptest.NewRecorder()
			testRouter.ServeHTTP(w, req)

			// Assert HTTP status code
			assert.Equal(t, testCase.expectedStatus, w.Code)

			// Assert response body
			assert.JSONEq(t, testCase.expectedBody, w.Body.String())
		})
	}
}

var postTestCases = []struct {
   name           string                 // Test case name
   payload        interface{}            // Input payload
   expectedStatus int                    // Expected HTTP status code
   expectedBody   string                 // Expected response body
}{
   {
      name: "Valid Payload",
      payload: parser.LogPayload{
         LogID:       "12345",
         Timestamp:   time.Now(),
         LogLevel:    "INFO",
         Message:     "This is a test log message.",
         Metadata:    map[string]interface{}{"key1": "value1"},
         Source:      "unit-test",
         Environment: "test",
         Type:        "application",
      },
      expectedStatus: http.StatusOK,
      expectedBody:   `{"message": "Log received"}`,
   },
   {
      name: "Missing Required Field",
      payload: map[string]interface{}{
         "logId": "12345",
      }, // Missing required fields such as "Timestamp", "LogLevel", etc.
      expectedStatus: http.StatusBadRequest,
      expectedBody:   `{"error": "Invalid JSON"}`,
   },
   {
      name:           "Invalid Timestamp",
      payload:        `{"logId":"12345","timestamp":"invalid-timestamp","logLevel":"INFO","message":"Test message","source":"unit-test","type":"application"}`,
      expectedStatus: http.StatusBadRequest,
      expectedBody:   `{"error":"Invalid JSON"}`, // Default Gin error for invalid data binding
   },
   {
      name:           "Invalid LogLevel",
      payload:        `{"logId":"12345","timestamp":"2025-04-27T11:00:00Z","logLevel":"INVALID","message":"Test message","source":"unit-test","type":"application"}`,
      expectedStatus: http.StatusBadRequest,
      expectedBody:   `{"error":"Invalid log level"}`,
   },
   {
      name:           "Malformed JSON",
      payload:        "{malformed-json", // Invalid JSON string
      expectedStatus: http.StatusBadRequest,
      expectedBody:   `{"error": "Invalid JSON"}`,
   },
}

