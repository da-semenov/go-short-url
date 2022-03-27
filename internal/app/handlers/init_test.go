package handlers

import (
	"errors"
	"github.com/da-semenov/go-short-url/internal/app/urls"

	"os"
	"testing"
)

var userService *UserServiceMock
var cryptoService *CryptoServiceMock
var userHandler *UserHandler

func TestMain(m *testing.M) {
	userService = new(UserServiceMock)
	userService.On("GetID", "full_URL").Return("short_URL", "short_URL", nil)
	userService.On("GetID", "original_URL").Return("short_URL", "short_URL", nil)
	userService.On("GetID", "bad_URL").Return("short_URL", "short_URL", nil)
	userService.On("GetID", "").Return("", "", errors.New("url is empty"))

	userService.On("GetURLsByUser", "user_id").Return("url-for-user-1", nil)

	var d []urls.UserBatch
	d = append(d, urls.UserBatch{CorrelationID: "correlation1", OriginalURL: "original_URL_1"})
	userService.On("SaveBatch", "user_id", d).Return("correlation1", "short_URL_1", nil)

	userService.On("SaveUserURL", "user_id", "original_URL", "short_URL").Return(nil)
	userService.On("SaveUserURL", "user_id", "bad_URL", "short_URL").Return(urls.ErrDuplicateKey)
	userService.On("GetURLByShort", "user_id", "short_URL").Return("original_URL", nil)
	userService.On("GetURLByShort", "", "short_URL").Return("original_URL", nil)
	userService.On("GetURLByShort", "user_id", "badURL").Return("", urls.ErrNotFound)
	userService.On("GetURLByShort", "", "badURL").Return("", urls.ErrNotFound)

	cryptoService := new(CryptoServiceMock)
	cryptoService.On("Validate", "user_id").Return(true, "user_id")

	cryptoService.On("GetNewUserToken").Return("user_id", "valid_user_Token", nil)

	userHandler = NewUserHandler(userService, cryptoService)
	os.Exit(m.Run())
}
