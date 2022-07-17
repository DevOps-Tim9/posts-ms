package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"posts-ms/src/dto/request"
	"posts-ms/src/service"
	"posts-ms/src/utils"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v8"
)

type CommentController struct {
	CommentService service.ICommentService
	validate       *validator.Validate
	logger         *logrus.Entry
}

func NewCommentController(commentService service.ICommentService) CommentController {
	config := &validator.Config{TagName: "validate"}
	logger := utils.Logger()

	return CommentController{CommentService: commentService, validate: validator.New(config), logger: logger}
}

func (c CommentController) GetAllByPostId(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "Handle /api/comments/posts/{postId}")

	defer span.Finish()

	c.logger.Info("Getting comments for specified post request received")
	params := mux.Vars(r)

	id, error := strconv.Atoi(params["postId"])

	if error != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	comments := c.CommentService.GetAllByPostId(uint(id), ctx)

	payload, _ := json.Marshal(comments)

	c.logger.Info("Returning list of comments")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(payload))
}

func (c CommentController) Create(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "Handle /api/comments")

	defer span.Finish()

	c.logger.Info("Creating comment request received")

	var commentDto request.CommentDto

	json.NewDecoder(r.Body).Decode(&commentDto)

	error := c.validate.Struct(commentDto)

	if error != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	newLike, error := c.CommentService.Create(commentDto, ctx)

	if error != nil {
		c.logger.Error("Error occured in creating comment")

		AddSystemEvent(time.Now().Format("2006-01-02 15:04:05"), "Comment unsuccessfully created")

		handleCommentError(error, w)

		return
	}

	payload, _ := json.Marshal(newLike)

	c.logger.Info("Comment created successfully")

	AddSystemEvent(time.Now().Format("2006-01-02 15:04:05"), "Comment successfully created")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(payload))
}

func (c CommentController) Delete(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "Handle /api/comments/{id}")

	defer span.Finish()

	c.logger.Info("Deleting comment request received")

	params := mux.Vars(r)

	id, error := strconv.Atoi(params["id"])

	if error != nil {
		c.logger.Error("Error occured in deleting comment")

		AddSystemEvent(time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf("Comment with id %d unsuccessfully deleted", id))

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	c.CommentService.Delete(uint(id), ctx)

	c.logger.Info("Deleting comment was successful")

	AddSystemEvent(time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf("Comment with id %d successfully deleted", id))

	w.WriteHeader(http.StatusNoContent)
}

func handleCommentError(error error, w http.ResponseWriter) http.ResponseWriter {
	w.WriteHeader(http.StatusBadRequest)

	return w
}

func AddSystemEvent(time string, message string) error {
	logger := utils.Logger()
	event := request.EventRequestDTO{
		Timestamp: time,
		Message:   message,
	}

	b, _ := json.Marshal(&event)
	endpoint := os.Getenv("EVENTS_MS")
	logger.Info("Sending system event to events-ms")
	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(b))
	req.Header.Set("content-type", "application/json")

	_, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Debug("Error happened during sending system event")
		return err
	}

	return nil
}
