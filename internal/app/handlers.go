package app

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Service interface {
	GetID(url string) (string, error)
	GetURL(id string) (string, error)
	GetShorten(url string) (*ShortenResponse, error)
}

type URLHandler struct {
	service Service
}

func EncodeURLHandler(service Service) *URLHandler {
	var h URLHandler
	h.service = service
	return &h
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}

func (u *URLHandler) DefaultHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Unsupported request type",
		http.StatusMethodNotAllowed)
}

func (u *URLHandler) PostMethodHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(b) == 0 {
		http.Error(w, "Body can't be empty", http.StatusBadRequest)
		return
	}

	res, _ := u.service.GetID(string(b))
	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(res))
}

func (u *URLHandler) PostShortenHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(b) == 0 {
		http.Error(w, "Body can't be empty", http.StatusBadRequest)
		return
	} else {
		var req ShortenRequest
		if err := json.Unmarshal(b, &req); err != nil {
			http.Error(w, "json error", http.StatusBadRequest)
			return
		}
		result, err := u.service.GetShorten(req.URL)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		responseBody, err := json.Marshal(result)
		if err != nil {
			http.Error(w, "Can't serialize response", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(responseBody)
		if err != nil {
			http.Error(w, "Can't write response", http.StatusBadRequest)
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
		http.Error(w, "URL was not found", http.StatusBadRequest)
		return
	}
	w.Header().Add("Location", res)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
