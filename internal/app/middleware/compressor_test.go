package middleware

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func stubHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func TestCompressGET(t *testing.T) {
	type args struct {
		acceptEncoding string
		method         string
		pattern        string
	}
	type wants struct {
		responseCode    int
		contentType     string
		contentEncoding string
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{name: "GET compress test 1 ",
			args:  args{acceptEncoding: "gzip", method: "GET", pattern: "/shortURL"},
			wants: wants{responseCode: http.StatusOK, contentType: "application/json", contentEncoding: "gzip"},
		},
		{name: "GET compress test 2 ",
			args:  args{acceptEncoding: "deflate", method: "GET", pattern: "/shortURL"},
			wants: wants{responseCode: http.StatusOK, contentType: "application/json", contentEncoding: ""},
		},
		{name: "GET compress test 3 ",
			args:  args{acceptEncoding: "", method: "GET", pattern: "/shortURL"},
			wants: wants{responseCode: http.StatusOK, contentType: "application/json", contentEncoding: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.args.method, tt.args.pattern, nil)
			request.Header.Set("Accept-Encoding", tt.args.acceptEncoding)
			h := GzipHandle(http.HandlerFunc(stubHandler))
			w := httptest.NewRecorder()
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.wants.responseCode, res.StatusCode, "Expected status %d, got %d", tt.wants.responseCode, res.StatusCode)
			assert.Equal(t, tt.wants.contentType, res.Header.Get("Content-Type"), "Expected Content-Type %d, got %d", tt.wants.responseCode, res.StatusCode)
			assert.Equal(t, tt.wants.contentEncoding, res.Header.Get("Content-Encoding"), "Expected Content-Encoding %d, got %d", tt.wants.responseCode, res.StatusCode)
		})
	}
}
