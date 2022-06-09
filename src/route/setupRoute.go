package route

import (
	"posts-ms/src/config"

	"github.com/gorilla/mux"
)

func SetupRoutes(container config.ControllerContainer) *mux.Router {
	route := mux.NewRouter()

	routerWithApiAsPrefix := route.PathPrefix("/api").Subrouter()

	routerWithApiAsPrefix.HandleFunc("/posts", container.PostController.Create).Methods("POST")
	routerWithApiAsPrefix.HandleFunc("/posts/{id}", container.PostController.Delete).Methods("DELETE")
	routerWithApiAsPrefix.HandleFunc("/posts/users/{userId}", container.PostController.GetAllByUserId).Methods("GET")
	routerWithApiAsPrefix.HandleFunc("/posts/users", container.PostController.GetAllByUserIds).Methods("POST")

	routerWithApiAsPrefix.HandleFunc("/likes", container.LikeController.Create).Methods("POST")
	routerWithApiAsPrefix.HandleFunc("/likes/users/{userId}/posts/{postId}", container.LikeController.Delete).Methods("DELETE")
	routerWithApiAsPrefix.HandleFunc("/likes/posts/{postId}", container.LikeController.GetAllByPostId).Methods("GET")

	routerWithApiAsPrefix.HandleFunc("/comments", container.CommentController.Create).Methods("POST")
	routerWithApiAsPrefix.HandleFunc("/comments/{id}", container.CommentController.Delete).Methods("DELETE")
	routerWithApiAsPrefix.HandleFunc("/comments/posts/{postId}", container.CommentController.GetAllByPostId).Methods("GET")

	return routerWithApiAsPrefix
}
