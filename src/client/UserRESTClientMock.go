package client

import (
	"posts-ms/src/dto/response"

	"github.com/stretchr/testify/mock"
)

type UserRESTClientMock struct {
	mock.Mock
}

func (m UserRESTClientMock) GetUser(id int) (*response.UserResponseDTO, error) {
	return &response.UserResponseDTO{Auth0ID: "1", Username: "Username"}, nil
}
