package web

import (
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestInfohashMiddleware(t *testing.T) {
	t.Parallel()

	inputV1 := "1234567890123456789012345678901234567890"
	v1, err := hex.DecodeString(inputV1)
	if err != nil {
		t.Fatalf("error decoding infohash: %v", err)
	}

	inputV2 := "123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0"
	v2, err := hex.DecodeString(inputV2)
	if err != nil {
		t.Fatalf("error decoding infohash v2: %v", err)
	}

	tests := []struct {
		name             string
		urlInfohash      string
		expectedStatus   int
		expectedInfohash []byte
	}{
		{
			name:             "Valid Infohash v1",
			urlInfohash:      inputV1,
			expectedStatus:   http.StatusOK,
			expectedInfohash: v1,
		},
		{
			name:             "Valid Infohash v2",
			urlInfohash:      inputV2,
			expectedStatus:   http.StatusOK,
			expectedInfohash: v2,
		},
		{
			name:             "Invalid Infohash",
			urlInfohash:      "invalidinfohash",
			expectedStatus:   http.StatusBadRequest,
			expectedInfohash: []byte(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{}
			req.SetPathValue("infohash", tt.urlInfohash)

			rr := httptest.NewRecorder()
			handler := infohashMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				infohash := r.Context().Value(InfohashKey).([]byte)
				if !reflect.DeepEqual(infohash, tt.expectedInfohash) {
					t.Errorf("expected infohash %v, got %v", tt.expectedInfohash, infohash)
				}
				w.WriteHeader(http.StatusOK)
			}))

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}