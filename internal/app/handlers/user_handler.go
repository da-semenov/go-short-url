package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

type CryptoService interface {
	Validate(token string) (bool, string)
	GetNewUserToken() (string, string, error)
}

type UserService interface {
	GetURLsByUser(userID string) ([]string, error)
	Ping() bool
}

type UserHandler struct {
	service       UserService
	cryptoService CryptoService
}

func NewUserHandler(service UserService, cs CryptoService) *UserHandler {
	var h UserHandler
	h.service = service
	h.cryptoService = cs
	return &h
}

func (z *UserHandler) bakeCookie() (*http.Cookie, error) {
	var c http.Cookie
	_, token, err := z.cryptoService.GetNewUserToken()
	if err != nil {
		return nil, err
	}
	c.Name = "token"
	c.Value = token
	return &c, nil
}

func (z *UserHandler) getTokenCookie(w http.ResponseWriter, r *http.Request) (*http.Cookie, error) {
	token, err := r.Cookie("token")
	ok, _ := z.cryptoService.Validate(token.Value)
	if errors.Is(err, http.ErrNoCookie) || !ok {
		newToken, err := z.bakeCookie()
		if err != nil {
			return nil, err
		}
		http.SetCookie(w, newToken)
		token = newToken
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil, err
	}
	return token, nil
}

func (z *UserHandler) GetUserURLsHandler(w http.ResponseWriter, r *http.Request) {
	token, err := z.getTokenCookie(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if ok, userID := z.cryptoService.Validate(token.Value); ok {
		res, err := z.service.GetURLsByUser(userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		} else {
			result := make([]string, len(res))
			for i := range res {
				// TODO:
				result[i] = res[i]
			}
			responseBody, err := json.Marshal(result)
			if err != nil {
				panic("Can't serialize response")
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, err = w.Write(responseBody)
			if err != nil {
				panic("Can't write response")
			}
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}

func (z *UserHandler) PingHandler(w http.ResponseWriter, r *http.Request) {
	if !z.service.Ping() {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}
