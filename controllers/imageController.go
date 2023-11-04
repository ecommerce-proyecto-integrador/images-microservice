package controllers

/*
import (
	"fmt"
	"os"
)

func SaveImage(image []byte, fileName string) error {
	err := os.WriteFile(fmt.Sprintf("./images/%s", fileName), image, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Elimina una imagen del disco local.
func DeleteImage(fileName string) error {
	err := os.Remove(fmt.Sprintf("./images/%s", fileName))
	if err != nil {
		return err
	}
	return nil
}
*/

import (
	"fmt"
	"github.com/gorilla/mux"

	/*amqp "github.com/rabbitmq/amqp091-go"*/
	"net/http"
	"os"
)

/*type ImageController struct {
	conn *amqp.Connection
}

func NewImageController(conn *amqp.Connection) *ImageController {
	return &ImageController{conn: conn}
}*/

func /*(c *ImageController) */ UploadImage(fileName string, imageData []byte) error {
	// Procesa la solicitud para subir la imagen utilizando el nombre del archivo y los datos de la imagen.
	// Guarda el archivo en el sistema de archivos (puedes cambiar la ubicación según tus necesidades).
	imagePath := "../images/" + fileName
	dest, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = dest.Write(imageData)
	if err != nil {
		return err
	}

	// Puedes enviar una respuesta de éxito o notificar a otros servicios si es necesario.
	// Ejemplo de publicación en RabbitMQ:
	// err = c.PublishImageUploadedMessage(fileName)
	// if err != nil {
	//     return err
	// }

	return nil
}

func /*(c *ImageController)*/ GetImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageID := vars["id"]

	imagePath := "../images/" + imageID
	fmt.Println(imagePath)
	// Convierte el ID de la imagen a un tipo adecuado (puede variar según tus necesidades).
	/*id, err := strconv.Atoi(imageID)
	if err != nil {
		http.Error(w, "Invalid image ID", http.StatusBadRequest)
		return
	}*/

	// Supongamos que tus imágenes están almacenadas en la carpeta "images".
	http.StripPrefix("/images/", http.FileServer(http.Dir("images"))).ServeHTTP(w, r)

}

func /*(c *ImageController)*/ GetProductWithImage(w http.ResponseWriter, r *http.Request) {
	// Implementa la lógica para obtener un producto con su imagen correspondiente
}
