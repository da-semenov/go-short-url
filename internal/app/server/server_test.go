package server

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestURLService_GetID(t *testing.T) {

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
			url:  "full_URL",
			want: want{
				path:   "/full_URL",
				scheme: "http",
				host:   "localhost:8080",
			},
			wantErr: false,
		},
		{
			name: "GetID Test 2",
			url:  "",
			want: want{
				path:   "/short_URL",
				scheme: "http",
				host:   "localhost:8080",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewUserService(dbRepoMock, fileRepoMock, "http://localhost:8080/")
			s.encode = MockEncode
			res, _, err := s.GetID(tt.url)
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
