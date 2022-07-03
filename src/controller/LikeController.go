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

type LikeController struct {
	LikeService service.ILikeService
	validate    *validator.Validate
	logger      *logrus.Entry
}

func NewLikeController(likeService service.ILikeService) LikeController {
	config := &validator.Config{TagName: "validate"}
	logger := utils.Logger()

	return LikeController{LikeService: likeService, validate: validator.New(config), logger: logger}
}

func (c LikeController) GetAllByPostId(w http.ResponseWriter, r *http.Request) {
	c.logger.Info("Getting all likes for specified post request received")

	params := mux.Vars(r)

	id, error := strconv.Atoi(params["postId"])

	if error != nil {
		w.WriteHeader(http.StatusBadRequest)

		c.logger.Error("Error occured in getting likes for specified post")

		return
	}

	c.logger.Info("Returning list of likes for specified post")

	likes := c.LikeService.GetAllByPostId(uint(id))

	payload, _ := json.Marshal(likes)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(payload))
}

func (c LikeController) Create(w http.ResponseWriter, r *http.Request) {
	c.logger.Info("Creating like request received")
	var likeDto request.LikeDto

	json.NewDecoder(r.Body).Decode(&likeDto)

	error := c.validate.Struct(likeDto)

	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newLike, error := c.LikeService.Create(likeDto)

	if error != nil {
		c.logger.Error("Error occured in creating like")

		handleLikeError(error, w)
		return
	}

	payload, _ := json.Marshal(newLike)

	c.logger.Info("Like created successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(payload))

}

func (c LikeController) Delete(w http.ResponseWriter, r *http.Request) {
	c.logger.Info("Deleting like request received")

	params := mux.Vars(r)

	userId, error := strconv.Atoi(params["userId"])

	if error != nil {
		c.logger.Error("Error occured in deleting like")

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	postId, error := strconv.Atoi(params["postId"])

	if error != nil {
		c.logger.Error("Error occured in deleting like")

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	c.logger.Info("Like deleted successfully")

	c.LikeService.Delete(uint(userId), uint(postId))

	w.WriteHeader(http.StatusNoContent)
}

func handleLikeError(error error, w http.ResponseWriter) http.ResponseWriter {
	w.WriteHeader(http.StatusBadRequest)

	return w
}
