package rabbitmq

import (
	"context"
	"encoding/json"
	"posts-ms/src/dto/request"
	"posts-ms/src/dto/response"
	"time"

	"github.com/gofrs/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/streadway/amqp"
)

type RMQProducer struct {
	ConnectionString string
}

func (r RMQProducer) StartRabbitMQ() (*amqp.Channel, error) {
	connectRabbitMQ, err := amqp.Dial(r.ConnectionString)

	if err != nil {
		return nil, err
	}

	channelRabbitMQ, err := connectRabbitMQ.Channel()

	if err != nil {
		return nil, err
	}

	return channelRabbitMQ, err
}

func DeleteImage(id uint, channel *amqp.Channel, ctx context.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Third service (rabbitmq) - Send request to media-ms for deleting media")

	defer span.Finish()

	uuid, _ := uuid.NewV4()

	media := response.MediaDto{
		Id:  id,
		Url: "",
	}

	payload, _ := json.Marshal(media)

	channel.Publish(
		"DeleteImageOnMedias-MS-exchange",    // exchange
		"DeleteImageOnMedias-MS-routing-key", // routing key
		false,                                // mandatory
		false,                                // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			MessageId:    uuid.String(),
			Timestamp:    time.Now(),
			Body:         payload,
		})
}

func AddNotification(notification *request.NotificationDTO, channel *amqp.Channel, ctx context.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Third service (rabbitmq) - Send request to user-ms for notifying user")

	defer span.Finish()

	uuid, _ := uuid.NewV4()

	payload, _ := json.Marshal(notification)

	channel.Publish(
		"AddNotification-MS-exchange",    // exchange
		"AddNotification-MS-routing-key", // routing key
		false,                            // mandatory
		false,                            // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			MessageId:    uuid.String(),
			Timestamp:    time.Now(),
			Body:         payload,
		})
}
