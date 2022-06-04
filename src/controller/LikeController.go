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

type LikeController struct {
	LikeService service.ILikeService
	validate    *validator.Validate
}

func NewLikeController(likeService service.ILikeService) LikeController {
	config := &validator.Config{TagName: "validate"}

	return LikeController{LikeService: likeService, validate: validator.New(config)}
}

func (c LikeController) GetAllByPostId(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, error := strconv.Atoi(params["postId"])

	if error != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	likes := c.LikeService.GetAllByPostId(uint(id))

	payload, _ := json.Marshal(likes)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(payload))
}

func (c LikeController) Create(w http.ResponseWriter, r *http.Request) {
	var likeDto request.LikeDto

	json.NewDecoder(r.Body).Decode(&likeDto)

	error := c.validate.Struct(likeDto)

	if error != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	newLike, error := c.LikeService.Create(likeDto)

	if error != nil {
		handleLikeError(error, w)

		return
	}

	payload, _ := json.Marshal(newLike)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(payload))

}

func (c LikeController) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userId, error := strconv.Atoi(params["userId"])

	if error != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	postId, error := strconv.Atoi(params["postId"])

	if error != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	c.LikeService.Delete(uint(userId), uint(postId))

	w.WriteHeader(http.StatusNoContent)
}

func handleLikeError(error error, w http.ResponseWriter) http.ResponseWriter {
	w.WriteHeader(http.StatusBadRequest)

	return w
}
