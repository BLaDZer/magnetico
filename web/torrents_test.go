package web

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"tgragnato.it/magnetico/persistence"
)

func TestTorrents(t *testing.T) {
	t.Parallel()

	var buffer bytes.Buffer
	if err := torrents().Render(&buffer); err != nil {
		t.Errorf("torrents render: %v", err)
	}
}

func TestTorrentsHandler(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequest("GET", "/torrents", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(torrentsHandler)
	handler.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", res.Status)
	}

	contentType := res.Header.Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("expected Content-Type text/html; got %v", contentType)
	}
}

func TestApiTorrents(t *testing.T) {
	t.Parallel()

	initDb()

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "valid request without optional params",
			queryParams:    "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid request with only lastOrderedValue",
			queryParams:    "lastOrderedValue=123.45",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "`lastOrderedValue`, `lastID` must be supplied altogether, if supplied.",
		},
		{
			name:           "invalid request with only lastID",
			queryParams:    "lastID=123",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "`lastOrderedValue`, `lastID` must be supplied altogether, if supplied.",
		},
		{
			name:           "valid request with lastOrderedValue and lastID",
			queryParams:    "lastOrderedValue=123.45&lastID=123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid request with non-numeric epoch",
			queryParams:    "epoch=abc",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "error while parsing the URL: strconv.ParseInt: parsing \"abc\": invalid syntax",
		},
		{
			name:           "invalid request with non-boolean ascending",
			queryParams:    "ascending=notabool",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "error while parsing the URL: strconv.ParseBool: parsing \"notabool\": invalid syntax",
		},
		{
			name:           "invalid request with non-numeric limit",
			queryParams:    "limit=notanumber",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "error while parsing the URL: strconv.ParseUint: parsing \"notanumber\": invalid syntax",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/api/torrents?"+tt.queryParams, nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			rec := httptest.NewRecorder()
			handler := http.HandlerFunc(apiTorrents)
			handler.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %v; got %v", tt.expectedStatus, res.StatusCode)
			}

			if tt.expectedError != "" {
				if !strings.Contains(rec.Body.String(), tt.expectedError) {
					t.Errorf("expected error %q; got %q", tt.expectedError, rec.Body.String())
				}
			}
		})
	}
}

func TestParseOrderBy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input          string
		expectedOutput persistence.OrderingCriteria
		expectedError  bool
	}{
		{"RELEVANCE", persistence.ByRelevance, false},
		{"TOTAL_SIZE", persistence.ByTotalSize, false},
		{"DISCOVERED_ON", persistence.ByDiscoveredOn, false},
		{"N_FILES", persistence.ByNFiles, false},
		{"UPDATED_ON", persistence.ByUpdatedOn, false},
		{"N_SEEDERS", persistence.ByNSeeders, false},
		{"N_LEECHERS", persistence.ByNLeechers, false},
		{"UNKNOWN", persistence.ByDiscoveredOn, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			output, err := parseOrderBy(tt.input)
			if (err != nil) != tt.expectedError {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
			}
			if output != tt.expectedOutput {
				t.Errorf("expected output: %v, got: %v", tt.expectedOutput, output)
			}
		})
	}
}
