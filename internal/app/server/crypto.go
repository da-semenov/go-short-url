package server

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"strconv"
	"strings"
	"time"
)

type CryptoService struct {
	key    []byte
	aesgcm cipher.AEAD
	nonce  []byte
}

func NewCryptoService() (*CryptoService, error) {
	var cs CryptoService
	cs.key = []byte(strings.Repeat("a", aes.BlockSize))

	aesblock, err := aes.NewCipher(cs.key)
	if err != nil {
		return nil, err
	}

	cs.aesgcm, err = cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}

	cs.nonce = make([]byte, cs.aesgcm.NonceSize())
	_, err = rand.Read(cs.nonce)
	if err != nil {
		return nil, err
	}

	return &cs, nil
}

func (s *CryptoService) generateUserID() (string, error) {
	uid := strconv.FormatInt(time.Now().UnixNano(), 10)
	return uid, nil
}

func (s *CryptoService) GetNewUserToken() (string, string, error) {
	user, err := s.generateUserID()
	if err != nil {
		return "", "", nil
	}
	token, err := s.encrypt([]byte(user))
	if err != nil {
		return "", "", nil
	}
	stringToken := base64.StdEncoding.EncodeToString(token)
	return user, stringToken, nil
}

func (s *CryptoService) Validate(token string) (bool, string) {
	t, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return false, ""
	}
	res, err := s.decrypt(t)
	if err != nil {
		return false, ""
	}
	return true, res
}

func (s *CryptoService) decrypt(src []byte) (string, error) {
	res, err := s.aesgcm.Open(nil, s.nonce, src, nil)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func (s *CryptoService) encrypt(userID []byte) ([]byte, error) {
	dst := s.aesgcm.Seal(nil, s.nonce, userID, nil)
	return dst, nil
}
