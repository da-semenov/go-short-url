package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/da-semenov/go-short-url/internal/app/urls"
	"net/http"
)

type CryptoService interface {
	Validate(token string) (bool, string)
	GetNewUserToken() (string, string, error)
}

type UserService interface {
	GetURLsByUser(ctx context.Context, userID string) ([]urls.UserURLs, error)
	SaveUserURL(ctx context.Context, userID string, originalURL string, shortURL string) error
	SaveBatch(ctx context.Context, userID string, src []urls.UserBatch) ([]urls.UserBatchResult, error)
	GetURLByShort(ctx context.Context, userID string, shortURL string) (string, error)
	GetID(url string) (string, string, error)
	Ping(ctx context.Context) bool
}

type UserHandler struct {
	userService   UserService
	cryptoService CryptoService
}

func NewUserHandler(userService UserService, cs CryptoService) *UserHandler {
	var h UserHandler
	h.userService = userService
	h.cryptoService = cs
	return &h
}

func (z *UserHandler) makeCookie() (*http.Cookie, string, error) {
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
	var userID string
	var ok bool
	token, err := r.Cookie("token")
	if err == nil {
		ok, userID = z.cryptoService.Validate(token.Value)
		if !ok {
			fmt.Println("invalid cookie")
		}
	}
	if errors.Is(err, http.ErrNoCookie) || !ok {
		fmt.Println("new cookie")
		var newToken *http.Cookie
		newToken, userID, err = z.makeCookie()
		fmt.Println(userID, " ", newToken)
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
	res, err := z.userService.GetURLsByUser(r.Context(), userID)
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
			http.Error(w, "can't serialize response", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(responseBody)
		if err != nil {
			http.Error(w, "can't write response", http.StatusBadRequest)
			return
		}
	}
}

func (z *UserHandler) PingHandler(w http.ResponseWriter, r *http.Request) {
	if !z.userService.Ping(r.Context()) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (z *UserHandler) DefaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	_, err := w.Write([]byte("unsupported request type"))
	if err != nil {
		http.Error(w, "can't write response", http.StatusBadRequest)
		return
	}
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
	if string(b) == "" {
		http.Error(w, "body can't be empty", http.StatusBadRequest)
		return
	} else {
		resURL, key, _ := z.userService.GetID(string(b))

		err = z.userService.SaveUserURL(r.Context(), userID, string(b), key)
		if errors.As(err, &urls.ErrDuplicateKey) {
			fmt.Println(key)
			w.WriteHeader(http.StatusConflict)
			_, err = w.Write([]byte(resURL))
			if err != nil {
				http.Error(w, "can't write response", http.StatusBadRequest)
				return
			}
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(resURL))
		if err != nil {
			http.Error(w, "can't write response", http.StatusBadRequest)
			return
		}
		return
	}
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
	if string(b) == "" {
		http.Error(w, "body can't be empty", http.StatusBadRequest)
		return
	} else {
		var req urls.ShortenRequest
		if err := json.Unmarshal(b, &req); err != nil {
			http.Error(w, "json error", http.StatusBadRequest)
			return
		}
		resURL, key, err := z.userService.GetID(req.URL)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		result := urls.ShortenResponse{Result: resURL}
		responseBody, err := json.Marshal(result)
		if err != nil {
			http.Error(w, "can't serialize response", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		err = z.userService.SaveUserURL(r.Context(), userID, req.URL, key)
		if errors.Is(err, urls.ErrDuplicateKey) {
			w.WriteHeader(http.StatusConflict)
			_, err = w.Write(responseBody)
			if err != nil {
				http.Error(w, "can't write response", http.StatusBadRequest)
				return
			}
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(responseBody)
		if err != nil {
			http.Error(w, "can't write response", http.StatusBadRequest)
			return
		}
		return
	}
}

func (z *UserHandler) PostShortenBatchHandler(w http.ResponseWriter, r *http.Request) {
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

	var req []urls.UserBatch
	if err := json.Unmarshal(b, &req); err != nil {
		http.Error(w, "json error", http.StatusBadRequest)
		return
	}
	result, err := z.userService.SaveBatch(r.Context(), userID, req)
	if errors.As(err, &urls.ErrDuplicateKey) {
		w.WriteHeader(http.StatusConflict)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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
}

func (z *UserHandler) GetMethodHandler(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "" || r.RequestURI[1:] == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	} else {
		_, err := z.getTokenCookie(w, r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		key := r.RequestURI[1:]
		userID := ""
		res, err := z.userService.GetURLByShort(r.Context(), userID, key)
		if errors.Is(err, urls.ErrNotFound) {
			w.WriteHeader(http.StatusGone)
			return
		}
		if err != nil {
			http.Error(w, "url was not found", http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", res)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
}
