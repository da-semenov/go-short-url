package app

import (
	"encoding/gob"
	"errors"
	"io"
	"os"
	"path"
)

type Storage struct {
	store          map[string]string
	cfgFileStorage string
	f              *os.File
	encoder        *gob.Encoder
}

type StoreRecord struct {
	Key   string
	Value string
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}

func NewStorage(fileStorage string) (*Storage, error) {
	var s Storage
	var tmpPath string
	s.cfgFileStorage = fileStorage
	s.store = make(map[string]string)

	err := os.MkdirAll(path.Dir(s.cfgFileStorage), 0755)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(s.cfgFileStorage); !os.IsNotExist(err) {
		tmpPath, err = s.copyStoreToTmp()
		if err != nil {
			return nil, err
		}
		err = s.init(tmpPath)
		if err != nil {
			return nil, err
		}
	}
	f, err := os.OpenFile(s.cfgFileStorage, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	s.encoder = gob.NewEncoder(f)
	err = s.flush()
	if err != nil {
		return nil, err
	}
	s.f = f
	os.Remove(tmpPath)
	return &s, nil
}

func (s *Storage) init(filePath string) error {
	f, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	gobDecoder := gob.NewDecoder(f)

	tmp := new(StoreRecord)
	for {
		err := gobDecoder.Decode(tmp)
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return err
		}
		s.store[tmp.Key] = tmp.Value
	}
	return nil
}

func (s *Storage) flush() error {
	for k, v := range s.store {
		rec := StoreRecord{Key: k, Value: v}
		err := s.encoder.Encode(&rec)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) copyStoreToTmp() (string, error) {
	in, err := os.Open(s.cfgFileStorage)
	if err != nil {
		return "", err
	}
	defer in.Close()

	out, err := os.CreateTemp(path.Dir(s.cfgFileStorage), "*.tmp")
	dstPath := out.Name()
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return "", err
	}
	return dstPath, out.Close()
}

func (s *Storage) Find(id string) (string, error) {
	if val, ok := s.store[id]; ok {
		return val, nil
	}
	return "", errors.New("id not found")
}

func (s *Storage) Save(id string, value string) error {
	var err error
	if val, ok := s.store[id]; ok {
		if val != value {
			s.store[id] = value
			rec := StoreRecord{Key: id, Value: value}
			err = s.encoder.Encode(&rec)
		} else {
			s.store[id] = value
		}
	} else {
		s.store[id] = value
		rec := StoreRecord{Key: id, Value: value}
		err = s.encoder.Encode(&rec)
	}
	return err
}
