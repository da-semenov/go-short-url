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
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Unsupported request type"))
		if err != nil {
			panic("Can't write response")
		}
	}
}

func (z *URLHandler) postMethodHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if string(b) == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("Bad request"))
		if err != nil {
			panic("Can't write response")
		}
		return
	} else {
		res, _ := z.service.GetID(string(b))
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(res))
		return
	}
}

func (z *URLHandler) getMethodHandler(w http.ResponseWriter, r *http.Request) {
	id := r.RequestURI[1:]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		res, err := z.service.GetURL(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				panic("Can't write response")
			}
			return
		}
		w.Header().Set("Location", res)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
}
