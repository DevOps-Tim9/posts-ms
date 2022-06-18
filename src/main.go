package main

import (
	"fmt"
	"net/http"
	"os"
	"posts-ms/src/client"
	"posts-ms/src/config"
	"posts-ms/src/controller"
	"posts-ms/src/rabbitmq"
	"posts-ms/src/repository"
	"posts-ms/src/route"
	"posts-ms/src/service"

	"github.com/rs/cors"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

func main() {
	dataBase, _ := config.SetupDB()

	amqpServerURL := os.Getenv("AMQP_SERVER_URL")

	rabbit := rabbitmq.RMQProducer{
		ConnectionString: amqpServerURL,
	}

	channel, _ := rabbit.StartRabbitMQ()

	defer channel.Close()

	repositoryContainer := initializeRepositories(dataBase)
	serviceContainer := initializeServices(repositoryContainer, channel)
	controllerContainer := initializeControllers(serviceContainer)

	router := route.SetupRoutes(controllerContainer)

	port := os.Getenv("SERVER_PORT")

	http.ListenAndServe(fmt.Sprintf(":%s", port), cors.AllowAll().Handler(router))
}

func initializeControllers(serviceContainer config.ServiceContainer) config.ControllerContainer {
	postController := controller.NewPostController(serviceContainer.PostService)
	likeController := controller.NewLikeController(serviceContainer.LikeService)
	commentController := controller.NewCommentController(serviceContainer.CommentService)

	container := config.NewControllerContainer(
		postController,
		likeController,
		commentController,
	)

	return container
}

func initializeServices(repositoryContainer config.RepositoryContainer, channel *amqp.Channel) config.ServiceContainer {
	mediaClient := client.NewMediaRESTClient()
	userClient := client.NewUserRESTClient()
	postService := service.PostService{
		PostRepository:    repositoryContainer.PostRepository,
		LikeRepository:    repositoryContainer.LikeRepository,
		CommentRepository: repositoryContainer.CommentRepository,
		MediaClient:       mediaClient,
		RabbitMQChannel:   channel,
	}
	likeService := service.LikeService{LikeRepository: repositoryContainer.LikeRepository, PostService: postService, UserRESTClient: userClient, RabbitMQChannel: channel}
	commentService := service.CommentService{CommentRepository: repositoryContainer.CommentRepository, PostService: postService, UserRESTClient: userClient, RabbitMQChannel: channel}

	container := config.NewServiceContainer(
		postService,
		likeService,
		commentService,
	)

	return container
}

func initializeRepositories(dataBase *gorm.DB) config.RepositoryContainer {
	postRepository := repository.PostRepository{Database: dataBase}
	likeRepository := repository.LikeRepository{Database: dataBase}
	commentRepository := repository.CommentRepository{Database: dataBase}

	container := config.NewRepositoryContainer(
		postRepository,
		likeRepository,
		commentRepository,
	)

	return container
}
