package handlers

import (
	"encoding/json"
	"errors"
	"github.com/da-semenov/go-short-url/internal/app/urls"
	"net/http"
)

type CryptoService interface {
	Validate(token string) (bool, string)
	GetNewUserToken() (string, string, error)
}

type UserService interface {
	GetURLsByUser(userID string) ([]urls.UserURLs, error)
	Save(userID string, originalURL string, shortURL string) error
	Ping() bool
}

type UserHandler struct {
	userService   UserService
	service       Service
	cryptoService CryptoService
}

func NewUserHandler(userService UserService, service Service, cs CryptoService) *UserHandler {
	var h UserHandler
	h.userService = userService
	h.service = service
	h.cryptoService = cs
	return &h
}

func (z *UserHandler) bakeCookie() (*http.Cookie, string, error) {
	var c http.Cookie
	userID, token, err := z.cryptoService.GetNewUserToken()
	if err != nil {
		return nil, "", err
	}
	c.Name = "token"
	c.Value = token
	return &c, userID, nil
}

func (z *UserHandler) getTokenCookie(w http.ResponseWriter, r *http.Request) (string, error) {
	token, err := r.Cookie("token")
	ok, userID := z.cryptoService.Validate(token.Value)
	if errors.Is(err, http.ErrNoCookie) || !ok {
		var newToken *http.Cookie
		newToken, userID, err = z.bakeCookie()
		if err != nil {
			return "", err
		}
		http.SetCookie(w, newToken)
		token = newToken
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return "", err
	}
	return userID, nil
}

func (z *UserHandler) GetUserURLsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := z.getTokenCookie(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res, err := z.userService.GetURLsByUser(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(res) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	} else {
		responseBody, err := json.Marshal(res)
		if err != nil {
			panic("can't serialize response")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(responseBody)
		if err != nil {
			panic("can't write response")
		}
	}
}

func (z *UserHandler) PingHandler(w http.ResponseWriter, r *http.Request) {
	if !z.userService.Ping() {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (z *UserHandler) DefaultHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Unsupported request type",
		http.StatusMethodNotAllowed)
}

func (z *UserHandler) PostMethodHandler(w http.ResponseWriter, r *http.Request) {
	b, err := getRequestBody(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userID, err := z.getTokenCookie(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(b) == 0 {
		http.Error(w, "body can't be empty", http.StatusBadRequest)
		return
	}
	res, _ := z.service.GetID(string(b))
	err = z.userService.Save(userID, string(b), res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(res))
	if err != nil {
		panic("Can't write response")
	}
	return
}

func (z *UserHandler) PostShortenHandler(w http.ResponseWriter, r *http.Request) {
	b, err := getRequestBody(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userID, err := z.getTokenCookie(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(b) == 0 {
		http.Error(w, "Body can't be empty", http.StatusBadRequest)
		return
	} else {
		var req urls.ShortenRequest
		if err := json.Unmarshal(b, &req); err != nil {
			http.Error(w, "json error", http.StatusBadRequest)
			return
		}
		res, err := z.service.GetID(req.URL)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		err = z.userService.Save(userID, req.URL, res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		result := urls.ShortenResponse{Result: res}
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
