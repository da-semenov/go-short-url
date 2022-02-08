package app

import (
	"encoding/base64"
	"errors"
)

type storageMap struct {
	store map[string]string
}

func NewStorage() *storageMap {
	var s storageMap
	s.store = make(map[string]string)
	return &s
}

func (s *storageMap) encode(str string) string {
	sha := base64.StdEncoding.EncodeToString([]byte(str))
	return sha
}

func (s *storageMap) GetID(url string) (string, error) {
	const baseURL string = "http://localhost:8080/"
	id := s.encode(url)
	s.store[id] = url
	return baseURL + id, nil
}

func (s *storageMap) GetURL(id string) (string, error) {
	if val, ok := s.store[id]; ok {
		return val, nil
	}
	return "", errors.New("id not found")
}
