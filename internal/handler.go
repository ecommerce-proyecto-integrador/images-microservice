package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	//"github.com/ValeHenriquez/example-rabbit-go/tasks-server/controllers"
	//"github.com/ValeHenriquez/example-rabbit-go/tasks-server/models"

	"github.com/ecommerce-proyecto-integrador/images-microservice/controllers"
	"github.com/ecommerce-proyecto-integrador/images-microservice/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func generateUniqueImageName() string {
	// Genera un UUID como parte del nombre del archivo.
	uniqueID := uuid.New()

	// Obtiene el timestamp actual para asegurar aún más la unicidad.
	//currentTime := time.Now().UnixNano()

	// Extrae la extensión del nombre de archivo original.
	fileExtension := ".jpg"

	// Combina todos los componentes para formar un nombre único.
	uniqueName := fmt.Sprintf(uniqueID.String(), fileExtension)

	return uniqueName
}

func Handler(d amqp.Delivery, ch *amqp.Channel) {
	//imageController := controllers.NewImageController(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var response models.Response

	actionType := d.Type
	log.Println("Está aquí")
	switch actionType {
	case "UPLOAD_IMAGE":
		log.Println(" [.] Uploading image")

		var data struct {
			Image []byte `json:"image"`
		}
		err := json.Unmarshal(d.Body, &data)
		failOnError(err, "Failed to unmarshal image")
		fileName := generateUniqueImageName()
		image := models.FtpImage{
			Name: fileName,
		}
		imageJson, err := json.Marshal(data)
		failOnError(err, "Failed to marshal image")

		_, err = controllers.UploadImage(image, fileName, data.Image)
		if err != nil {
			response = models.Response{
				Success: "error",
				Message: "Error uploading image",
				Data:    []byte(err.Error()),
			}
		} else {
			response = models.Response{
				Success: "succes",
				Message: "Image uploaded",
				Data:    imageJson,
			}
		}

	case "DELETE_IMAGE":
		log.Println(" [.] Deleting image")

		var data struct {
			FileName string `json:"filename"`
		}
		err := json.Unmarshal(d.Body, &data)
		err = controllers.DeleteImage(data.FileName)

		ImageJson, err := json.Marshal(data)
		if err != nil {
			response = models.Response{
				Success: "error",
				Message: "Error deleting image",
				Data:    []byte(err.Error()),
			}
		} else {
			response = models.Response{
				Success: "succes",
				Message: "Image deleted",
				Data:    ImageJson,
			}
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
