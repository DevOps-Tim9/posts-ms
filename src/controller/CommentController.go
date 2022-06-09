package controller

import (
	"encoding/json"
	"net/http"
	"posts-ms/src/dto/request"
	"posts-ms/src/service"
	"strconv"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v8"
)

type CommentController struct {
	CommentService service.ICommentService
	validate       *validator.Validate
}

func NewCommentController(commentService service.ICommentService) CommentController {
	config := &validator.Config{TagName: "validate"}

	return CommentController{CommentService: commentService, validate: validator.New(config)}
}

func (c CommentController) GetAllByPostId(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, error := strconv.Atoi(params["postId"])

	if error != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	comments := c.CommentService.GetAllByPostId(uint(id))

	payload, _ := json.Marshal(comments)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(payload))
}

func (c CommentController) Create(w http.ResponseWriter, r *http.Request) {
	var commentDto request.CommentDto

	json.NewDecoder(r.Body).Decode(&commentDto)

	error := c.validate.Struct(commentDto)

	if error != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	newLike, error := c.CommentService.Create(commentDto)

	if error != nil {
		handleCommentError(error, w)

		return
	}

	payload, _ := json.Marshal(newLike)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(payload))
}

func (c CommentController) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, error := strconv.Atoi(params["id"])

	if error != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	c.CommentService.Delete(uint(id))

	w.WriteHeader(http.StatusNoContent)
}

func handleCommentError(error error, w http.ResponseWriter) http.ResponseWriter {
	w.WriteHeader(http.StatusBadRequest)

	return w
}
