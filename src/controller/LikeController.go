package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "Handle /api/likes/posts/{postId}")

	defer span.Finish()

	c.logger.Info("Getting all likes for specified post request received")

	params := mux.Vars(r)

	id, error := strconv.Atoi(params["postId"])

	if error != nil {
		w.WriteHeader(http.StatusBadRequest)

		c.logger.Error("Error occured in getting likes for specified post")

		return
	}

	c.logger.Info("Returning list of likes for specified post")

	likes := c.LikeService.GetAllByPostId(uint(id), ctx)

	payload, _ := json.Marshal(likes)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(payload))
}

func (c LikeController) Create(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "Handle /api/likes")

	defer span.Finish()

	c.logger.Info("Creating like request received")
	var likeDto request.LikeDto

	json.NewDecoder(r.Body).Decode(&likeDto)

	error := c.validate.Struct(likeDto)

	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newLike, error := c.LikeService.Create(likeDto, ctx)

	if error != nil {
		c.logger.Error("Error occured in creating like")

		AddSystemEvent(time.Now().Format("2006-01-02 15:04:05"), "Like unsuccessfully created")

		handleLikeError(error, w)
		return
	}

	payload, _ := json.Marshal(newLike)

	c.logger.Info("Like created successfully")

	AddSystemEvent(time.Now().Format("2006-01-02 15:04:05"), "Like successfully created")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(payload))

}

func (c LikeController) Delete(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "Handle /api/likes/users/{userId}/posts/{postId}")

	defer span.Finish()

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

	AddSystemEvent(time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf("Like for post with id %d of user with id %d successfully deleted", postId, userId))

	c.LikeService.Delete(uint(userId), uint(postId), ctx)

	w.WriteHeader(http.StatusNoContent)
}

func handleLikeError(error error, w http.ResponseWriter) http.ResponseWriter {
	w.WriteHeader(http.StatusBadRequest)

	return w
}
