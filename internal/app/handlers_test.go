package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type ServiceMock struct {
	mock.Mock
}

func (s *ServiceMock) GetID(url string) (string, error) {
	args := s.Called(url)
	return args.String(0), args.Error(1)
}

func (s *ServiceMock) GetURL(id string) (string, error) {
	args := s.Called(id)
	return args.String(0), args.Error(1)
}

func (s *ServiceMock) GetShorten(url string) (*ShortenResponse, error) {
	args := s.Called(url)
	return &ShortenResponse{Result: args.String(0)}, args.Error(1)
}

var service *ServiceMock
var handler *URLHandler

func TestMain(m *testing.M) {
	service = new(ServiceMock)

	service.On("GetID", "URL").Return("encode_URL", nil)
	service.On("GetID", "").Return("", errors.New("URL is empty"))

	service.On("GetURL", "encode_URL").Return("URL", nil)
	service.On("GetURL", "xxx").Return("", errors.New("id not found"))

	service.On("GetShorten", "URL").Return("encode_URL", nil)
	service.On("GetShorten", "").Return("encode_URL", errors.New("URL is empty"))

	handler = EncodeURLHandler(service)
	os.Exit(m.Run())
}

func TestURLHandler_postMethodHandler(t *testing.T) {
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
		{name: "POST test 1.Empty body",
			args:  args{requestBody: ""},
			wants: wants{responseCode: http.StatusBadRequest, resultResponse: ""},
		},
		{name: "POST test 2.Encoded URL",
			args:  args{requestBody: "URL"},
			wants: wants{responseCode: http.StatusCreated, resultResponse: "encode_URL"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("POST", "/", strings.NewReader(tt.args.requestBody))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler.PostMethodHandler)
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

func TestURLHandler_getMethodHandler(t *testing.T) {
	type args struct {
		encodeURLid string
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
		{name: "GET test 1.Get URL.",
			args:  args{encodeURLid: "encode_URL"},
			wants: wants{responseCode: http.StatusTemporaryRedirect, resultResponse: "URL"},
		},
		{name: "GET test 2.Empty id.",
			args:  args{encodeURLid: ""},
			wants: wants{responseCode: http.StatusBadRequest, resultResponse: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", fmt.Sprintf("/%s", tt.args.encodeURLid), nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler.GetMethodHandler)
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

func TestURLHandler_defaultHandler(t *testing.T) {
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
		{name: "test 1.PUT method.",
			args:  args{method: "PUT"},
			wants: wants{responseCode: http.StatusMethodNotAllowed, resultResponse: "Unsupported request type"},
		},
		{name: "test 2.PATCH method.",
			args:  args{method: "PATCH"},
			wants: wants{responseCode: http.StatusMethodNotAllowed, resultResponse: "Unsupported request type"},
		},
		{name: "test 3.DELETE method.",
			args:  args{method: "DELETE"},
			wants: wants{responseCode: http.StatusMethodNotAllowed, resultResponse: "Unsupported request type"},
		},
		{name: "test 4.HEAD method.",
			args:  args{method: "HEAD"},
			wants: wants{responseCode: http.StatusMethodNotAllowed, resultResponse: "Unsupported request type"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.args.method, "/", nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler.DefaultHandler)
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.wants.responseCode, res.StatusCode, "Expected status %d, got %d", tt.wants.responseCode, res.StatusCode)

			responseBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("Can't read response body, %e", err)
			}
			assert.Equal(t, "Unsupported request type\n", string(responseBody), "Expected body is %s, got %s", tt.wants.resultResponse, string(responseBody))
		})
	}
}

func TestURLHandler_postMethodShortenHandler(t *testing.T) {
	type args struct {
		request *ShortenRequest
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
		{name: "test 1.Positive.",
			args: args{request: &ShortenRequest{"URL"}},
			wants: wants{responseCode: http.StatusCreated,
				response: "encode_URL"},
		},
		{name: "test 2.Empty body.",
			args: args{request: nil},
			wants: wants{responseCode: http.StatusBadRequest,
				response: ""},
		},
		{name: "test 3.Empty URL.",
			args: args{request: &ShortenRequest{""}},
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
			h := http.HandlerFunc(handler.PostShortenHandler)
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
				var result ShortenResponse
				if err := json.Unmarshal(responseBody, &result); err != nil {
					t.Error("Can't unmarshal", err)
				}
				assert.Equal(t, tt.wants.response, result.Result, "Expected body is %s, got %s", tt.wants.response, result.Result)
			}
		})
	}
}
