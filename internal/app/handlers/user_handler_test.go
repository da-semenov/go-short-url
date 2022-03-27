package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/da-semenov/go-short-url/internal/app/urls"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestURLHandler_GetUserURLsHandler(t *testing.T) {
	type args struct {
		shortURLKey string
	}
	type wants struct {
		responseCode   int
		resultResponse string
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{name: "Test 1.Get URL.",
			args:  args{shortURLKey: "short_URL"},
			wants: wants{responseCode: http.StatusOK, resultResponse: "full_URL"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "/user/urls", nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(userHandler.GetUserURLsHandler)
			request.AddCookie(&http.Cookie{Name: "token", Value: "user_id"})

			h.ServeHTTP(w, request)
			res := w.Result()

			defer res.Body.Close()
			assert.Equal(t, tt.wants.responseCode, res.StatusCode, "Expected status %d, got %d", tt.wants.responseCode, res.StatusCode)
		})
	}
}

func TestUserHandler_getTokenCookie(t *testing.T) {
	type args struct {
		w        http.ResponseWriter
		r        *http.Request
		userName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Test 1. Get Token Cookie.",
			args:    args{w: httptest.NewRecorder(), r: httptest.NewRequest("GET", "/user/urls", nil), userName: "user_id"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userHandler.getTokenCookie(tt.args.w, tt.args.r)

			if (err != nil) != tt.wantErr {
				t.Errorf("getTokenCookie() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.args.userName, got, "Expected cookie.name is %s, got %s", got, tt.args.userName)
		})
	}
}

func TestUserHandler_getTokenCookieHeader(t *testing.T) {
	type args struct {
		w          http.ResponseWriter
		r          *http.Request
		cookieName string
		cookieVal  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Test 2. Get Token Cookie. Check Header.",
			args:    args{w: httptest.NewRecorder(), r: httptest.NewRequest("GET", "/user/urls", nil), cookieName: "token", cookieVal: "valid_user_Token"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userHandler.getTokenCookie(tt.args.w, tt.args.r)

			if (err != nil) != tt.wantErr {
				t.Errorf("getTokenCookie() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)
			http.SetCookie(tt.args.w, &http.Cookie{Name: "some Name", Value: "Some value1"})
			http.SetCookie(tt.args.w, &http.Cookie{Name: "some Name", Value: "Some value2"})
			http.SetCookie(tt.args.w, &http.Cookie{Name: "some Name", Value: "Some value3"})

			tmp := tt.args.w.Header().Get("Set-Cookie")
			assert.NotEmpty(t, tmp, "Can't got cookie from response")
			parsedCookie := strings.Split(tmp, "=")
			assert.ElementsMatch(t, parsedCookie, []string{tt.args.cookieName, tt.args.cookieVal}, "ParsedCookie does not match expected")
		})
	}
}

func TestUserHandler_PostShortenBatchHandler(t *testing.T) {
	type args struct {
		requestBody string
	}
	type wants struct {
		responseCode int
		contentType  string
		responseBody string
	}
	tests := []struct {
		name  string
		wants wants
		args  args
	}{
		{name: "Test 1. ShortenBatchHandler.",
			wants: wants{
				responseCode: http.StatusCreated,
				contentType:  "application/json",
				responseBody: "",
			},
			args: args{requestBody: "[{\"correlation_id\": \"correlation1\",\"original_URL\": \"original_URL_1\"}]"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBody := []byte(tt.args.requestBody)

			request := httptest.NewRequest("POST", "/api/shorten/batch", bytes.NewReader(requestBody))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(userHandler.PostShortenBatchHandler)

			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, tt.wants.responseCode, res.StatusCode, "Expected status %d, got %d", tt.wants.responseCode, res.StatusCode)

		})
	}
}

func TestUserHandler_PostMethodHandler(t *testing.T) {
	type args struct {
		requestBody string
	}
	type wants struct {
		responseCode   int
		resultResponse string
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{name: "Test 1. Empty body.",
			args:  args{requestBody: ""},
			wants: wants{responseCode: http.StatusBadRequest, resultResponse: ""},
		},
		{name: "Test 2. Positive.",
			args:  args{requestBody: "original_URL"},
			wants: wants{responseCode: http.StatusCreated, resultResponse: "short_URL"},
		},
		{name: "Test 3. Negative.",
			args:  args{requestBody: "bad_URL"},
			wants: wants{responseCode: http.StatusConflict, resultResponse: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("POST", "/", strings.NewReader(tt.args.requestBody))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(userHandler.PostMethodHandler)
			h.ServeHTTP(w, request)
			res := w.Result()
			fmt.Println(res)

			assert.Equal(t, tt.wants.responseCode, res.StatusCode, "Expected status %d, got %d", tt.wants.responseCode, res.StatusCode)

			if res.StatusCode == http.StatusCreated {
				responseBody, err := io.ReadAll(res.Body)
				defer res.Body.Close()
				if err != nil {
					t.Errorf("Can't read response body, %e", err)
				}
				assert.Equal(t, tt.wants.resultResponse, string(responseBody), "Expected body is %s, got %s", tt.wants.resultResponse, string(responseBody))
			}
		})
	}
}

func TestUserHandler_postShortenHandler(t *testing.T) {
	type args struct {
		request *urls.ShortenRequest
	}
	type wants struct {
		responseCode int
		response     string
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{name: "Test 1. Positive.",
			args: args{request: &urls.ShortenRequest{URL: "original_URL"}},
			wants: wants{responseCode: http.StatusCreated,
				response: "short_URL"},
		},
		{name: "Test 2. Empty body.",
			args: args{request: nil},
			wants: wants{responseCode: http.StatusBadRequest,
				response: ""},
		},
		{name: "Test 3. Object with empty url.",
			args: args{request: &urls.ShortenRequest{URL: ""}},
			wants: wants{responseCode: http.StatusBadRequest,
				response: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requestBody []byte
			if tt.args.request != nil {
				requestBody, _ = json.Marshal(tt.args.request)
			}
			request := httptest.NewRequest("POST", "/api/shorten", bytes.NewReader(requestBody))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(userHandler.PostShortenHandler)
			h.ServeHTTP(w, request)
			res := w.Result()
			fmt.Println(res)

			assert.Equal(t, tt.wants.responseCode, res.StatusCode, "Expected status %d, got %d", tt.wants.responseCode, res.StatusCode)

			if res.StatusCode == http.StatusCreated {
				responseBody, err := io.ReadAll(res.Body)
				defer res.Body.Close()
				if err != nil {
					t.Errorf("Can't read response body, %e", err)
				}
				var result urls.ShortenResponse
				if err := json.Unmarshal(responseBody, &result); err != nil {
					t.Error("Can't unmarshal", err)
				}
				assert.Equal(t, tt.wants.response, result.Result, "Expected body is %s, got %s", tt.wants.response, result.Result)
			}
		})
	}
}

func TestUserHandler_postShortenHandler2(t *testing.T) {
	type args struct {
		requestBody string
	}
	type wants struct {
		responseCode int
		response     string
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{name: "Test 4. Empty object.",
			args: args{requestBody: "{}"},
			wants: wants{responseCode: http.StatusBadRequest,
				response: ""},
		},
		{name: "Test 5. Empty object.",
			args: args{requestBody: ""},
			wants: wants{responseCode: http.StatusBadRequest,
				response: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("POST", "/api/shorten", strings.NewReader(tt.args.requestBody))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(userHandler.PostShortenHandler)
			h.ServeHTTP(w, request)
			res := w.Result()
			fmt.Println(res)

			assert.Equal(t, tt.wants.responseCode, res.StatusCode, "Expected status %d, got %d", tt.wants.responseCode, res.StatusCode)

			if res.StatusCode == http.StatusCreated {
				responseBody, err := io.ReadAll(res.Body)
				defer res.Body.Close()
				if err != nil {
					t.Errorf("Can't read response body, %e", err)
				}
				var result urls.ShortenResponse
				if err := json.Unmarshal(responseBody, &result); err != nil {
					t.Error("Can't unmarshal", err)
				}
				assert.Equal(t, tt.wants.response, result.Result, "Expected body is %s, got %s", tt.wants.response, result.Result)
			}
		})
	}
}

func TestUserHandler_getMethodHandler(t *testing.T) {
	type args struct {
		shortURLKey string
	}
	type wants struct {
		responseCode   int
		resultResponse string
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{name: "Test 1. Get. Positive.",
			args:  args{shortURLKey: "short_URL"},
			wants: wants{responseCode: http.StatusTemporaryRedirect, resultResponse: "original_URL"},
		},
		{name: "Test 2. Get. Negative.",
			args:  args{shortURLKey: ""},
			wants: wants{responseCode: http.StatusBadRequest, resultResponse: ""},
		},
		{name: "Test 3. Get. Negative. Unexists short url.",
			args:  args{shortURLKey: "badURL"},
			wants: wants{responseCode: http.StatusGone, resultResponse: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", fmt.Sprintf("/%s", tt.args.shortURLKey), nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(userHandler.GetMethodHandler)
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.wants.responseCode, res.StatusCode, "Expected status %d, got %d", tt.wants.responseCode, res.StatusCode)
			if res.StatusCode == tt.wants.responseCode {
				assert.Equal(t, tt.wants.resultResponse, res.Header.Get("Location"))
			}
		})
	}
}

func TestUserHandler_DefaultHandler(t *testing.T) {
	type args struct {
		method string
	}
	type wants struct {
		responseCode   int
		resultResponse string
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{name: "Test 1. Other http method.",
			args:  args{method: "PUT"},
			wants: wants{responseCode: http.StatusBadRequest, resultResponse: "Unsupported request type"},
		},
		{name: "Test 2. Other http method.",
			args:  args{method: "PATCH"},
			wants: wants{responseCode: http.StatusBadRequest, resultResponse: "Unsupported request type"},
		},
		{name: "Test 3. Other http method.",
			args:  args{method: "DELETE"},
			wants: wants{responseCode: http.StatusBadRequest, resultResponse: "Unsupported request type"},
		},
		{name: "Test 4. Other http method.",
			args:  args{method: "HEAD"},
			wants: wants{responseCode: http.StatusBadRequest, resultResponse: "Unsupported request type"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.args.method, "/", nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(userHandler.DefaultHandler)
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.wants.responseCode, res.StatusCode, "Expected status %d, got %d", tt.wants.responseCode, res.StatusCode)
			responseBody, err := io.ReadAll(res.Body)

			if err != nil {
				t.Errorf("Can't read response body, %e", err)
			}
			assert.Equal(t, "Unsupported request type", string(responseBody), "Expected body is %s, got %s", tt.wants.resultResponse, string(responseBody))
		})
	}
}
