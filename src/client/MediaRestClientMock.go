package client

import (
	"mime/multipart"

	"github.com/stretchr/testify/mock"
)

type MediaRestClientMock struct {
	mock.Mock
}

func (m MediaRestClientMock) Upload(file multipart.File) (uint, error) {
	return uint(1), nil
}
