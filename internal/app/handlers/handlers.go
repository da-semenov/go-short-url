package handlers

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/da-semenov/go-short-url/internal/app/urls"
	"io"
	"net/http"
	"strings"
)

type ContextKey string

var contextKeyUID = ContextKey("uid")

type Service interface {
	GetID(url string) (string, error)
	GetURL(id string) (string, error)
}

type URLHandler struct {
	service       Service
	cryptoService CryptoService
}

func EncodeURLHandler(service Service) *URLHandler {
	var h URLHandler
	h.service = service
	return &h
}

func decompress(data []byte) ([]byte, error) {
	var res bytes.Buffer
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	_, err = res.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("failed decompress data: %v", err)
	}
	return res.Bytes(), nil
}

func getRequestBody(r *http.Request) ([]byte, error) {
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return nil, err
	}
	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		unzipBody, err := decompress(b)
		if err != nil {
			panic(err)
		}
		return unzipBody, nil
	}
	return b, nil
}

func (u *URLHandler) DefaultHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "unsupported request type",
		http.StatusMethodNotAllowed)
}

func (u *URLHandler) PostMethodHandler(w http.ResponseWriter, r *http.Request) {
	b, err := getRequestBody(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(b) == 0 {
		http.Error(w, "body can't be empty", http.StatusBadRequest)
		return
	} else {
		res, _ := u.service.GetID(string(b))
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(res))
		if err != nil {
			panic("Can't write response")
		}
		return
	}
}

func (u *URLHandler) PostShortenHandler(w http.ResponseWriter, r *http.Request) {
	b, err := getRequestBody(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(b) == 0 {
		http.Error(w, "body can't be empty", http.StatusBadRequest)
		return
	} else {
		var req urls.ShortenRequest
		if err := json.Unmarshal(b, &req); err != nil {
			http.Error(w, "json error", http.StatusBadRequest)
			return
		}
		res, err := u.service.GetID(req.URL)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		result := urls.ShortenResponse{Result: res}
		responseBody, err := json.Marshal(result)
		if err != nil {
			http.Error(w, "can't serialize response", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(responseBody)
		if err != nil {
			http.Error(w, "can't write response", http.StatusBadRequest)
			return
		}
		return
	}
}

func (u *URLHandler) GetMethodHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Error(w, "id was not provided", http.StatusBadRequest)
		return
	}
	id := r.URL.Path[1:]
	res, err := u.service.GetURL(id)
	if err != nil {
		http.Error(w, "url was not found", http.StatusBadRequest)
		return
	}
	w.Header().Add("Location", res)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
