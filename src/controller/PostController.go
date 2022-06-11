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

type PostController struct {
	PostService service.IPostService
	validate    *validator.Validate
	logger      *logrus.Entry
}

func NewPostController(postService service.IPostService) PostController {
	config := &validator.Config{TagName: "validate"}
	logger := utils.Logger()

	return PostController{PostService: postService, validate: validator.New(config), logger: logger}
}

func (c PostController) GetAllByUserId(w http.ResponseWriter, r *http.Request) {
	c.logger.Info("Getting all posts for specified user request received")
	params := mux.Vars(r)

	id, error := strconv.Atoi(params["userId"])

	if error != nil {
		c.logger.Error("Error occured in getting posts by user")

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	likes := c.PostService.GetAllByUserId(uint(id))

	payload, _ := json.Marshal(likes)

	c.logger.Info("Returning list of posts for specified user")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(payload))
}

func (c PostController) GetAllByUserIds(w http.ResponseWriter, r *http.Request) {
	c.logger.Info("Getting all posts for specified users request received")

	var search request.SearchPostPageableDto

	json.NewDecoder(r.Body).Decode(&search)

	likes := c.PostService.GetAllByUserIds(search.Ids)

	payload, _ := json.Marshal(likes)

	c.logger.Info("Returning list of posts for specified users")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(payload))
}

func (p PostController) Create(w http.ResponseWriter, r *http.Request) {
	p.logger.Info("Creating post request received")

	r.ParseMultipartForm(32 << 20)

	var postDto request.PostDto

	postDtoJSON := r.Form["post"][0]

	error := p.validate.Struct(postDto)

	error = json.Unmarshal([]byte(postDtoJSON), &postDto)

	if error != nil {
		p.logger.Error("Error occured in creating post")

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	files := r.MultipartForm.File["files"]

	post, err := p.PostService.Create(postDto, files)

	if err != nil {
		p.logger.Error("Error occured in creating post")

		handleMunicipalityError(err, w)

		return
	}

	payload, _ := json.Marshal(post)

	p.logger.Info("Post created successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(payload))
}

func (c PostController) Delete(w http.ResponseWriter, r *http.Request) {
	c.logger.Info("Deleting post request received")

	params := mux.Vars(r)

	id, error := strconv.Atoi(params["id"])

	if error != nil {
		c.logger.Error("Error occured in deleting post")

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	c.PostService.Delete(uint(id))

	c.logger.Info("Post deleted successfully")

	w.WriteHeader(http.StatusNoContent)
}

func handleMunicipalityError(error error, w http.ResponseWriter) http.ResponseWriter {
	w.WriteHeader(http.StatusConflict)

	return w
}
