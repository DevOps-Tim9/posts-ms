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

type PostController struct {
	PostService service.IPostService
	validate    *validator.Validate
}

func NewPostController(postService service.IPostService) PostController {
	config := &validator.Config{TagName: "validate"}

	return PostController{PostService: postService, validate: validator.New(config)}
}

func (c PostController) GetAllByUserId(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, error := strconv.Atoi(params["userId"])

	if error != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	likes := c.PostService.GetAllByUserId(uint(id))

	payload, _ := json.Marshal(likes)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(payload))
}

func (c PostController) GetAllByUserIds(w http.ResponseWriter, r *http.Request) {
	var search request.SearchPostPageableDto

	json.NewDecoder(r.Body).Decode(&search)

	likes := c.PostService.GetAllByUserIds(search.Ids)

	payload, _ := json.Marshal(likes)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(payload))
}

func (p PostController) Create(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)

	var postDto request.PostDto

	postDtoJSON := r.Form["post"][0]

	error := p.validate.Struct(postDto)

	error = json.Unmarshal([]byte(postDtoJSON), &postDto)

	if error != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	files := r.MultipartForm.File["files"]

	post, err := p.PostService.Create(postDto, files)

	if err != nil {
		handleMunicipalityError(err, w)

		return
	}

	payload, _ := json.Marshal(post)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(payload))
}

func (c PostController) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, error := strconv.Atoi(params["id"])

	if error != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	c.PostService.Delete(uint(id))

	w.WriteHeader(http.StatusNoContent)
}

func handleMunicipalityError(error error, w http.ResponseWriter) http.ResponseWriter {
	w.WriteHeader(http.StatusConflict)

	return w
}
