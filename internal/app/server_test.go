package app

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/url"
	"testing"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) Find(id string) (string, error) {
	args := r.Called(id)
	return args.String(0), nil
}

func (r *RepositoryMock) Save(id string, value string) error {
	return nil
}

func mockEncode(str string) string {
	return str
}

func TestURLService_GetID(t *testing.T) {
	repo := new(RepositoryMock)
	repo.On("Find", "encode_URL").Return("URL")
	repo.On("Find", "").Return("URL")
	repo.On("Save", "URL").Return("encode_URL")

	type want struct {
		path   string
		scheme string
		host   string
	}
	tests := []struct {
		name    string
		url     string
		want    want
		wantErr bool
	}{
		{
			name: "GetID Test 1",
			url:  "URL",
			want: want{
				path:   "/URL",
				scheme: "http",
				host:   "localhost:8080",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(repo)
			s.encode = mockEncode
			res, err := s.GetID(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			u, err := url.Parse(res)
			if err != nil {
				t.Errorf("GetID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				assert.Equal(t, tt.want.scheme, u.Scheme)
				assert.Equal(t, tt.want.host, u.Host)
				assert.Equal(t, tt.want.path, u.Path)
			}
		})
	}
}

func TestURLService_GetURL(t *testing.T) {
	repo := new(RepositoryMock)
	repo.On("Find", "encode_URL").Return("URL", nil)
	repo.On("Save", "URL").Return("encode_URL", nil)

	tests := []struct {
		name    string
		key     string
		want    string
		wantErr bool
	}{
		{
			name:    "GetURL Test 1",
			key:     "encode_URL",
			want:    "URL",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(repo)

			got, err := s.GetURL(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("URL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("URL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
