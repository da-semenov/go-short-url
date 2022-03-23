package server

import (
	"testing"
)

var fileRepoMock *FileRepositoryMock
var dbRepoMock *DBRepositoryMock

func TestMain(m *testing.M) {
	fileRepoMock = new(FileRepositoryMock)
	fileRepoMock.On("Find", "short_URL").Return("full_URL")
	fileRepoMock.On("Find", "").Return("full_URL")
	fileRepoMock.On("SaveUserURL", "full_URL").Return("short_URL")
	dbRepoMock = new(DBRepositoryMock)
}
