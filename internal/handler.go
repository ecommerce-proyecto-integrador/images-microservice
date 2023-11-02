package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	//"github.com/ValeHenriquez/example-rabbit-go/tasks-server/controllers"
	//"github.com/ValeHenriquez/example-rabbit-go/tasks-server/models"

	"github.com/ecommerce-proyecto-integrador/images-microservice/mod/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func Handler(d amqp.Delivery, ch *amqp.Channel) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var response models.Response

	actionType := d.Type
	switch actionType {
	case "UPLOAD_IMAGE":
		log.Println(" [.] Uploading image")

		var data struct {
			Image []byte `json:"image"`
		}
		err := json.Unmarshal(d.Body, &data)
		failOnError(err, "Failed to unmarshal image")

		fileName := fmt.Sprintf("%d.jpg", time.Now().Unix())
		err = controllers.saveImage(data.Image, fileName)
		failOnError(err, "Failed to save image")

		response = models.Response{
			Success: "success",
			Message: "Image uploaded",
			Data:    []byte(fileName),
		}

	case "DELETE_IMAGE":
		log.Println(" [.] Deleting image")

		var data struct {
			FileName string `json:"filename"`
		}
		err := json.Unmarshal(d.Body, &data)
		failOnError(err, "Failed to unmarshal filename")

		err = controllers.deleteImage(data.FileName)
		failOnError(err, "Failed to delete image")

		response = models.Response{
			Success: "success",
			Message: "Image deleted",
		}

	default:
		response = models.Response{
			Success: "error",
			Message: "Unknown action",
		}
	}

	responseJSON, err := json.Marshal(response)
	failOnError(err, "Failed to marshal response")

	err = ch.PublishWithContext(ctx,
		"",        // exchange
		d.ReplyTo, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: d.CorrelationId,
			Body:          responseJSON,
		})
	failOnError(err, "Failed to publish a message")

	d.Ack(false)
}
