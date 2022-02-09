package app

import (
	"io"
	"net/http"
)

type Service interface {
	GetID(url string) (string, error)
	GetURL(id string) (string, error)
}

type URLHandler struct {
	service Service
}

func ShortURLHandler() *URLHandler {
	var h URLHandler
	h.service = NewStorage()
	return &h
}

func (z *URLHandler) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		z.getMethodHandler(w, r)
	case http.MethodPost:
		z.postMethodHandler(w, r)
	default:
		http.Error(w, "Only GET and POST methods are allowed",
			http.StatusMethodNotAllowed)
		return
	}
}

func (z *URLHandler) postMethodHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(b) == 0 {
		http.Error(w, "Body can't be empty", http.StatusBadRequest)
		return
	} else {
		res, _ := z.service.GetID(string(b))
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(res))
		return
	}
}

func (z *URLHandler) getMethodHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Error(w, "id was not provided", http.StatusBadRequest)
		return
	}
	id := r.URL.Path[1:]
	res, err := z.service.GetURL(id)
	if err != nil {
		http.Error(w, "URL was not found", http.StatusBadRequest)
		return
	}
	w.Header().Add("Location", res)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
