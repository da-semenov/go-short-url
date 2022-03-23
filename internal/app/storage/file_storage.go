package storage

import (
	"encoding/gob"
	"errors"
	"io"
	"os"
	"path"
	"sync"
)

type FileStorage struct {
	sync.Mutex
	store          map[string]string
	cfgFileStorage string
	f              *os.File
	encoder        *gob.Encoder
}

type StoreRecord struct {
	Key   string
	Value string
}

func NewFileStorage(fileStorage string) (*FileStorage, error) {
	var s FileStorage
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

func (s *FileStorage) init(filePath string) error {
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

func (s *FileStorage) flush() error {
	for k, v := range s.store {
		rec := StoreRecord{Key: k, Value: v}
		err := s.encoder.Encode(&rec)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *FileStorage) copyStoreToTmp() (string, error) {
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

func (s *FileStorage) Find(key string) (string, error) {
	if val, ok := s.store[key]; ok {
		return val, nil
	} else {
		return "", errors.New("value not found")
	}
}

func (s *FileStorage) FindByUser(id string) ([]UserURLs, error) {
	// TODO
	return nil, errors.New("unexpecting using of method")
}

func (s *FileStorage) Save(key string, value string) error {
	var err error
	s.Lock()
	defer s.Unlock()
	if val, ok := s.store[key]; ok {
		if val != value {
			s.store[key] = value
			rec := StoreRecord{Key: key, Value: value}
			err = s.encoder.Encode(&rec)
		} else {
			s.store[key] = value
		}
	} else {
		s.store[key] = value
		rec := StoreRecord{Key: key, Value: value}
		err = s.encoder.Encode(&rec)
	}
	return err
}
