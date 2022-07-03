package main

import (
	"fmt"
	"net/http"
	"os"
	"posts-ms/src/client"
	"posts-ms/src/config"
	setupJaeger "posts-ms/src/config/jaeger"
	"posts-ms/src/controller"
	"posts-ms/src/rabbitmq"
	"posts-ms/src/repository"
	"posts-ms/src/route"
	"posts-ms/src/service"
	"posts-ms/src/utils"

	"github.com/opentracing/opentracing-go"
	"github.com/rs/cors"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

func main() {
	logger := utils.Logger()

	logger.Info("Connecting with DB")

	dataBase, _ := config.SetupDB()

	tracer, trCloser, err := setupJaeger.InitJaeger()

	if err != nil {
		logger.Debug(err.Error())

		fmt.Printf("error init jaeger %v", err)
	} else {
		defer trCloser.Close()

		opentracing.SetGlobalTracer(tracer)
	}

	amqpServerURL := os.Getenv("AMQP_SERVER_URL")

	logger.Info("Connecting on RabbitMq")

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

	logger.Info("Starting server")

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
		Logger:            utils.Logger(),
	}
	likeService := service.LikeService{LikeRepository: repositoryContainer.LikeRepository, PostService: postService, UserRESTClient: userClient, RabbitMQChannel: channel, Logger: utils.Logger()}
	commentService := service.CommentService{CommentRepository: repositoryContainer.CommentRepository, PostService: postService, UserRESTClient: userClient, RabbitMQChannel: channel, Logger: utils.Logger()}

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
