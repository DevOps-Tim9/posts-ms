package config

import (
	"posts-ms/src/controller"
	"posts-ms/src/repository"
	"posts-ms/src/service"
)

type ControllerContainer struct {
	PostController    controller.PostController
	LikeController    controller.LikeController
	CommentController controller.CommentController
}

type ServiceContainer struct {
	PostService    service.IPostService
	LikeService    service.ILikeService
	CommentService service.CommentService
}

type RepositoryContainer struct {
	PostRepository    repository.IPostRepository
	LikeRepository    repository.ILikeRepository
	CommentRepository repository.ICommentRepository
}

func NewControllerContainer(
	postController controller.PostController,
	likeController controller.LikeController,
	commentController controller.CommentController,
) ControllerContainer {
	return ControllerContainer{
		PostController:    postController,
		LikeController:    likeController,
		CommentController: commentController,
	}
}

func NewServiceContainer(
	postService service.IPostService,
	likeService service.ILikeService,
	commentService service.CommentService,
) ServiceContainer {
	return ServiceContainer{
		PostService:    postService,
		LikeService:    likeService,
		CommentService: commentService,
	}
}

func NewRepositoryContainer(
	postRepository repository.IPostRepository,
	likeRepository repository.ILikeRepository,
	commentRepository repository.ICommentRepository,
) RepositoryContainer {
	return RepositoryContainer{
		PostRepository:    postRepository,
		LikeRepository:    likeRepository,
		CommentRepository: commentRepository,
	}
}
