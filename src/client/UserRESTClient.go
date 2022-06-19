package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"posts-ms/src/dto/response"
)

type IUserRESTClient interface {
	GetUser(id int) (*response.UserResponseDTO, error)
}

type UserRESTClient struct{}

func NewUserRESTClient() UserRESTClient {
	return UserRESTClient{}
}

func (c UserRESTClient) GetUser(id int) (*response.UserResponseDTO, error) {
	user := response.UserResponseDTO{}
	endpoint := fmt.Sprintf("http://%s/users/%d", os.Getenv("USER_SERVICE_DOMAIN"), id)

	req, _ := http.NewRequest("GET", endpoint, nil)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Println(res.StatusCode)
		return nil, err
	} else {
		b, _ := io.ReadAll(res.Body)
		errr := json.Unmarshal(b, &user)
		if errr != nil {
			fmt.Println(errr.Error())
		}
	}
	return &user, nil
}
