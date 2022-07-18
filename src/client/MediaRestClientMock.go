package client

import (
	"context"
	"mime/multipart"

	"github.com/stretchr/testify/mock"
)

type MediaRestClientMock struct {
	mock.Mock
}

func (m MediaRestClientMock) Upload(file multipart.File, ctx context.Context) (uint, error) {
	return uint(1), nil
}
