package app

import "errors"

type Storage struct {
	store map[string]string
}

func NewStorage() *Storage {
	var s Storage
	s.store = make(map[string]string)
	return &s
}

func (s *Storage) Find(id string) (string, error) {
	if val, ok := s.store[id]; ok {
		return val, nil
	}
	return "", errors.New("id not found")
}

func (s *Storage) Save(id string, value string) error {
	s.store[id] = value
	return nil
}
