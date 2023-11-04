package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ecommerce-proyecto-integrador/images-microservice/config"
	"github.com/ecommerce-proyecto-integrador/images-microservice/controllers"
	"github.com/ecommerce-proyecto-integrador/images-microservice/internal"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

const imageStoragePath = "images"

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func getChannel() *amqp.Channel {
	ch := config.GetChannel()
	if ch == nil {
		log.Panic("Failed to get channel")
	}
	return ch
}

func declareQueue(ch *amqp.Channel) amqp.Queue {
	q, err := ch.QueueDeclare(
		"images_queue", // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Failed to declare a queue")
	return q
}

// Establece la calidad de servicio (QoS) para el canal de RabbitMQ.
func setQoS(ch *amqp.Channel) {
	err := ch.Qos(
		1,     // prefetch count: Especifica cuántos mensajes puede recibir un consumidor antes de que se detenga la entrega. En este caso, se establece en 1.
		0,     // prefetch size: No se usa en este caso, se establece como 0.
		false, // global: Indica si estas configuraciones de QoS se aplican a nivel de canal o a nivel de conexión. En este caso, es a nivel de canal (false).
	)

	failOnError(err, "Failed to set QoS")
}

// Registra un consumidor para la cola dada y devuelve un canal de entrega de mensajes.
func registerConsumer(ch *amqp.Channel, q amqp.Queue) <-chan amqp.Delivery {
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")
	return msgs
}

func main() {
	fmt.Println("Images MS starting...")

	// Configura la conexión con RabbitMQ
	godotenv.Load()
	fmt.Println("Loaded env variables...")
	config.SetupRabbitMQ()
	fmt.Println("RabbitMQ Connection configured...")

	ch := getChannel()              // Obtiene un canal de RabbitMQ
	q := declareQueue(ch)           // Declara una cola y obtiene su estructura
	setQoS(ch)                      // Establece la calidad de servicio en el canal
	msgs := registerConsumer(ch, q) // Registra un consumidor para la cola y obtiene un canal de entrega de mensajes

	// Configura el enrutador HTTP
	r := mux.NewRouter()

	// Define tus rutas
	//r.HandleFunc("/images", controllers.UploadImage()).Methods("POST")
	r.HandleFunc("/images/{id}", controllers.GetImage).Methods("GET")
	//r.HandleFunc("/products/{id}", imageController.GetProductWithImage).Methods("GET")

	go func() {
		for d := range msgs {
			internal.Handler(d, ch) // Llama al manejador de mensajes internos con el mensaje y el canal de RabbitMQ
		}
	}()

	// Inicia el servidor HTTP
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8181", nil))
}
