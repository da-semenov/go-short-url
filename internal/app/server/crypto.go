package server

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"time"
)

type CryptoService struct {
	key    []byte
	aesgcm cipher.AEAD
	nonce  []byte
}

func NewCryptoService() (*CryptoService, error) {
	var cs CryptoService
	cs.key, _ = hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
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
	uid := strconv.FormatInt(time.Now().Unix(), 10)
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
	return user, string(token), nil
}

func (s *CryptoService) Validate(token string) (bool, string) {
	res, err := s.decrypt([]byte(token))
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
