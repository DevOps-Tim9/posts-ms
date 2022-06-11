package controller

import (
	"encoding/json"
	"net/http"
	"posts-ms/src/dto/request"
	"posts-ms/src/service"
	"posts-ms/src/utils"
	"strconv"

	"github.com/gorilla/mux"
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
	c.logger.Info("Getting comments for specified post request received")
	params := mux.Vars(r)

	id, error := strconv.Atoi(params["postId"])

	if error != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	comments := c.CommentService.GetAllByPostId(uint(id))

	payload, _ := json.Marshal(comments)

	c.logger.Info("Returning list of comments")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(payload))
}

func (c CommentController) Create(w http.ResponseWriter, r *http.Request) {
	c.logger.Info("Creating comment request received")

	var commentDto request.CommentDto

	json.NewDecoder(r.Body).Decode(&commentDto)

	error := c.validate.Struct(commentDto)

	if error != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	newLike, error := c.CommentService.Create(commentDto)

	if error != nil {
		c.logger.Error("Error occured in creating comment")
		handleCommentError(error, w)

		return
	}

	payload, _ := json.Marshal(newLike)

	c.logger.Info("Comment created successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(payload))
}

func (c CommentController) Delete(w http.ResponseWriter, r *http.Request) {
	c.logger.Info("Deleting comment request received")

	params := mux.Vars(r)

	id, error := strconv.Atoi(params["id"])

	if error != nil {
		c.logger.Error("Error occured in deleting comment")

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	c.CommentService.Delete(uint(id))

	c.logger.Info("Deleting comment was successful")

	w.WriteHeader(http.StatusNoContent)
}

func handleCommentError(error error, w http.ResponseWriter) http.ResponseWriter {
	w.WriteHeader(http.StatusBadRequest)

	return w
}
