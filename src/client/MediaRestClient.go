package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"posts-ms/src/dto/response"
	"time"

	"github.com/opentracing/opentracing-go"
)

type IMediaClient interface {
	Upload(multipart.File, context.Context) (uint, error)
}

type MediaRESTClient struct {
	endpoint string
}

func NewMediaRESTClient() MediaRESTClient {
	return MediaRESTClient{endpoint: "http://medias-server:8082"}
}

func (c MediaRESTClient) Upload(image multipart.File, ctx context.Context) (uint, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Third service - Send request to media-ms for uploading media")

	defer span.Finish()

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)

	fw, err := writer.CreateFormFile("files", "image.png")

	if err != nil {
		return 0, err
	}

	_, err = io.Copy(fw, image)

	if err != nil {
		return 0, err
	}

	writer.Close()

	req, err := http.NewRequest("POST", c.endpoint+"/api/medias", bytes.NewReader(body.Bytes()))

	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := client.Do(req)

	if err != nil {
		return 0, err
	}

	var media response.MediaDto

	json.NewDecoder(res.Body).Decode(&media)

	return media.Id, err
}
