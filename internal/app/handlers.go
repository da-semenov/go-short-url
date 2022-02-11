package app

import (
	"io"
	"log"
	"net/http"
)

type Service interface {
	GetID(url string) (string, error)
	GetURL(id string) (string, error)
}

type URLHandler struct {
	service Service
}

func EncodeURLHandler(service Service) *URLHandler {
	var h URLHandler
	h.service = service
	return &h
}

func (u *URLHandler) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		u.getMethodHandler(w, r)
	case http.MethodPost:
		u.postMethodHandler(w, r)
	default:
		http.Error(w, "Only GET and POST methods are allowed",
			http.StatusMethodNotAllowed)
		return
	}
}

func (u *URLHandler) postMethodHandler(w http.ResponseWriter, r *http.Request) {
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
	_, err = w.Write([]byte(res))
}

func (u *URLHandler) getMethodHandler(w http.ResponseWriter, r *http.Request) {
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
	return
}
