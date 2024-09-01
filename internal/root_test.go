package internal

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// WbScannerTestSuite defines a test suite for WbScanner.
type WbScannerTestSuite struct {
	suite.Suite
	scanner    *WbScanner
	mockServer *httptest.Server
}

// SetupSuite is called before the suite starts running.
func (suite *WbScannerTestSuite) SetupSuite() {
	// Create a mock HTTP server
	suite.mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/test-get":
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"message": "GET request success"}`)
		case r.Method == http.MethodPost && r.URL.Path == "/test-post":
			if r.Header.Get("Content-Type") != "application/json" {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, `{"error": "Invalid Content-Type"}`)
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{"error": "Failed to read request body"}`)
				return
			}
			defer r.Body.Close()

			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"message": "POST request success", "received": %s}`, body)
		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, `{"error": "Not found"}`)
		}
	}))

	// Initialize WbScanner with the mock server's URL
	suite.scanner = NewWbScanner(suite.mockServer.URL, nil)
}

// TearDownSuite is called after the suite has finished running.
func (suite *WbScannerTestSuite) TearDownSuite() {
	if suite.mockServer != nil {
		suite.mockServer.Close()
	}
}

// TestGetRequest tests GET requests.
func (suite *WbScannerTestSuite) TestGetRequest() {
	body, err := suite.scanner.Request("/test-get", http.MethodGet, bytes.Buffer{}, nil, nil)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), `{"message": "GET request success"}`, string(body))
}

// TestPostRequest tests POST requests.
func (suite *WbScannerTestSuite) TestPostRequest() {
	payload := bytes.NewBufferString(`{"title": "test", "body": "content"}`)
	headers := map[string]string{"Authorization": "Bearer testtoken"}

	body, err := suite.scanner.Request("/test-post", http.MethodPost, *payload, nil, headers)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), `{"message": "POST request success", "received": {"title": "test", "body": "content"}}`, string(body))
}

// TestInvalidRoute tests handling of invalid routes.
func (suite *WbScannerTestSuite) TestInvalidRoute() {
	body, err := suite.scanner.Request("/invalid-route", http.MethodGet, bytes.Buffer{}, nil, nil)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), `{"error": "Not found"}`, string(body))
}

// Run the test suite
func TestWbScannerTestSuite(t *testing.T) {
	suite.Run(t, new(WbScannerTestSuite))
}
